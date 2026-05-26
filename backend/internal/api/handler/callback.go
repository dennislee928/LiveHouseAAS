package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/payment"
)

type CallbackHandler struct {
	pool        *pgxpool.Pool
	paymentRouter *payment.Router
}

func NewCallbackHandler(pool *pgxpool.Pool, pr *payment.Router) *CallbackHandler {
	return &CallbackHandler{pool: pool, paymentRouter: pr}
}

type CallbackRequest struct {
	Provider     string `json:"provider" binding:"required"`
	ProviderTxID string `json:"provider_tx_id" binding:"required"`
	Status       string `json:"status" binding:"required,oneof=completed failed"`
}

func (h *CallbackHandler) PaymentCallback(c *gin.Context) {
	var req CallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.paymentRouter.VerifyCallback(context.Background(),
		payment.Provider(req.Provider), req.ProviderTxID, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification failed"})
		return
	}

	if result.Status == "completed" {
		// update order status
		h.pool.Exec(context.Background(),
			`UPDATE orders SET status = 'paid', payment_method = $2, paid_at = NOW(), updated_at = NOW()
			 WHERE id = (SELECT order_id FROM transactions WHERE provider_tx_id = $1 LIMIT 1)
			 OR id = $1`,
			req.ProviderTxID, req.Provider)

		// update transaction
		h.pool.Exec(context.Background(),
			`UPDATE transactions SET status = 'completed', updated_at = NOW()
			 WHERE provider_tx_id = $1`, req.ProviderTxID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "callback processed", "status": result.Status})
}

func (h *CallbackHandler) ECPayNotify(c *gin.Context) {
	c.Request.ParseForm()
	params := make(map[string]string)
	for k, v := range c.Request.Form {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	merchantTradeNo := params["MerchantTradeNo"]
	rtnCode := params["RtnCode"]
	_ = merchantTradeNo

	if rtnCode == "1" {
		// payment successful
		h.pool.Exec(context.Background(),
			`UPDATE orders SET status = 'paid', updated_at = NOW()
			 WHERE id = $1`, merchantTradeNo)
	}

	c.String(http.StatusOK, "1|OK")
}

func (h *CallbackHandler) NewebPayNotify(c *gin.Context) {
	var callback payment.NewebPayCallback
	if err := c.ShouldBindJSON(&callback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid callback"})
		return
	}

	if callback.Status == "SUCCESS" {
		// decrypt TradeInfo to get order details
		// update order status
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *CallbackHandler) Refund(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderID := c.Param("orderId")

	var orderStatus, orderUserID string
	var totalAmount int32
	err := h.pool.QueryRow(context.Background(),
		`SELECT status, user_id, total_amount FROM orders WHERE id = $1`, orderID).
		Scan(&orderStatus, &orderUserID, &totalAmount)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if orderUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your order"})
		return
	}

	if orderStatus != "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order is not paid"})
		return
	}

	var provider, providerTxID string
	h.pool.QueryRow(context.Background(),
		`SELECT provider, COALESCE(provider_tx_id, '') FROM transactions WHERE order_id = $1 AND status = 'completed'`,
		orderID).Scan(&provider, &providerTxID)

	if providerTxID == "" {
		providerTxID = orderID
	}

	err = h.paymentRouter.Refund(context.Background(), payment.Provider(provider), providerTxID, totalAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refund failed: " + err.Error()})
		return
	}

	h.pool.Exec(context.Background(),
		`UPDATE orders SET status = 'refunded', updated_at = NOW() WHERE id = $1`, orderID)
	h.pool.Exec(context.Background(),
		`UPDATE transactions SET status = 'refunded', updated_at = NOW() WHERE order_id = $1`, orderID)
	h.pool.Exec(context.Background(),
		`UPDATE tickets SET status = 'cancelled' WHERE order_id = $1`, orderID)

	c.JSON(http.StatusOK, gin.H{"message": "refund processed"})
}
