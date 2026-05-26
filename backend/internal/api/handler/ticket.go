package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/domain/order"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/payment"
)

type TicketHandler struct {
	pool         *pgxpool.Pool
	paymentRouter *payment.Router
}

func NewTicketHandler(pool *pgxpool.Pool, pr *payment.Router) *TicketHandler {
	return &TicketHandler{pool: pool, paymentRouter: pr}
}

func generateCode() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *TicketHandler) Purchase(c *gin.Context) {
	eventID := c.Param("id")
	userID, _ := c.Get("user_id")

	var req order.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify event exists and is published
	var eventStatus string
	err := h.pool.QueryRow(context.Background(),
		`SELECT status FROM events WHERE id = $1`, eventID).Scan(&eventStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	if eventStatus != "published" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event is not available for purchase"})
		return
	}

	tx, err := h.pool.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	defer tx.Rollback(context.Background())

	var totalAmount int32
	for _, item := range req.Items {
		var ttPrice, ttQty, ttMax int32
		var ttStatus string
		err := tx.QueryRow(context.Background(),
			`SELECT price, quantity, max_per_order, status FROM ticket_types WHERE id = $1 AND event_id = $2`,
			item.TicketTypeID, eventID).Scan(&ttPrice, &ttQty, &ttMax, &ttStatus)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("ticket type %s not found", item.TicketTypeID)})
			return
		}
		if ttStatus != "active" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("ticket type %s is not available", item.TicketTypeID)})
			return
		}
		if item.Quantity > ttMax {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("max %d tickets per order for this type", ttMax)})
			return
		}

		// check remaining tickets
		var soldCount int32
		tx.QueryRow(context.Background(),
			`SELECT COALESCE(SUM(oi.quantity), 0) FROM order_items oi
			 JOIN orders o ON oi.order_id = o.id
			 WHERE oi.ticket_type_id = $1 AND o.status IN ('pending', 'paid')`, item.TicketTypeID).Scan(&soldCount)
		if soldCount+item.Quantity > ttQty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not enough tickets available"})
			return
		}

		subtotal := ttPrice * item.Quantity
		totalAmount += subtotal
	}

	// create order
	var orderID string
	err = tx.QueryRow(context.Background(),
		`INSERT INTO orders (id, user_id, event_id, total_amount, status, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, 'pending', NOW(), NOW())
		 RETURNING id`, userID, eventID, totalAmount).Scan(&orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	// create order items
	for _, item := range req.Items {
		var unitPrice int32
		tx.QueryRow(context.Background(),
			`SELECT price FROM ticket_types WHERE id = $1`, item.TicketTypeID).Scan(&unitPrice)
		subtotal := unitPrice * item.Quantity
		_, err = tx.Exec(context.Background(),
			`INSERT INTO order_items (id, order_id, ticket_type_id, quantity, unit_price, subtotal)
			 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`,
			orderID, item.TicketTypeID, item.Quantity, unitPrice, subtotal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order items"})
			return
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finalize order"})
		return
	}

	// initiate payment via mock provider
	paymentReq := &payment.PaymentRequest{
		OrderID:     orderID,
		Amount:      totalAmount,
		Currency:    "TWD",
		Description: fmt.Sprintf("Order %s", orderID),
		ReturnURL:   fmt.Sprintf("/orders/%s", orderID),
	}

	result, err := h.paymentRouter.Pay(context.Background(), payment.ProviderMock, paymentReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment initiation failed"})
		return
	}

	// for mock provider, auto-complete payment
	if result.Status == "completed" {
		h.completeOrder(context.Background(), orderID, "mock", result.ProviderTxID)
	}

	// fetch created order
	var o order.Order
	h.pool.QueryRow(context.Background(),
		`SELECT id, user_id, event_id, total_amount, status, COALESCE(payment_method,''), paid_at, created_at, updated_at
		 FROM orders WHERE id = $1`, orderID).
		Scan(&o.ID, &o.UserID, &o.EventID, &o.TotalAmount, &o.Status, &o.PaymentMethod, &o.PaidAt, &o.CreatedAt, &o.UpdatedAt)

	resp := order.PurchaseResponse{
		Order: o,
		Payment: order.PaymentInfo{
			Provider:    string(result.Provider),
			RedirectURL: result.RedirectURL,
			PaymentURL:  result.PaymentURL,
		},
	}

	// include ticket codes if payment was auto-completed
	if o.Status == "paid" {
		tickets, _ := h.getTicketsForOrder(context.Background(), orderID)
		resp.Tickets = tickets
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *TicketHandler) completeOrder(ctx context.Context, orderID, provider, providerTxID string) {
	// update order status
	h.pool.Exec(ctx,
		`UPDATE orders SET status = 'paid', payment_method = $2, paid_at = NOW(), updated_at = NOW() WHERE id = $1`,
		orderID, provider)

	// create transaction record
	h.pool.Exec(ctx,
		`INSERT INTO transactions (id, order_id, provider, amount, fee, provider_tx_id, status, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, (SELECT total_amount FROM orders WHERE id = $1), 0, $3, 'completed', NOW(), NOW())`,
		orderID, provider, providerTxID)

	// generate tickets
	rows, err := h.pool.Query(ctx,
		`SELECT oi.ticket_type_id, oi.quantity, o.event_id FROM order_items oi
		 JOIN orders o ON oi.order_id = o.id WHERE oi.order_id = $1`, orderID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ticketTypeID, eventID string
		var quantity int32
		rows.Scan(&ticketTypeID, &quantity, &eventID)

		for i := int32(0); i < quantity; i++ {
			code := generateCode()
			secret := generateCode()
			h.pool.Exec(ctx,
				`INSERT INTO tickets (id, order_id, ticket_type_id, event_id, code, qr_secret, status, created_at)
				 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'active', NOW())`,
				orderID, ticketTypeID, eventID, code, secret)
		}
	}
}

func (h *TicketHandler) getTicketsForOrder(ctx context.Context, orderID string) ([]order.TicketInfo, error) {
	rows, err := h.pool.Query(ctx,
		`SELECT id, code, qr_secret FROM tickets WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []order.TicketInfo
	for rows.Next() {
		var t order.TicketInfo
		rows.Scan(&t.ID, &t.Code, &t.QRSegret)
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (h *TicketHandler) ListOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := h.pool.Query(context.Background(),
		`SELECT o.id, o.user_id, o.event_id, o.total_amount, o.status, COALESCE(o.payment_method,''), o.paid_at, o.created_at, o.updated_at,
		        e.title, v.name
		 FROM orders o
		 JOIN events e ON o.event_id = e.id
		 JOIN venues v ON e.venue_id = v.id
		 WHERE o.user_id = $1
		 ORDER BY o.created_at DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list orders"})
		return
	}
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var id, userID2, eventID, status, paymentMethod, eventTitle, venueName string
		var totalAmount int32
		var paidAt, createdAt, updatedAt interface{}
		if err := rows.Scan(&id, &userID2, &eventID, &totalAmount, &status, &paymentMethod, &paidAt, &createdAt, &updatedAt, &eventTitle, &venueName); err != nil {
			continue
		}
		orders = append(orders, map[string]interface{}{
			"id": id, "event_id": eventID, "total_amount": totalAmount, "status": status,
			"payment_method": paymentMethod, "paid_at": paidAt,
			"created_at": createdAt, "updated_at": updatedAt,
			"event_title": eventTitle, "venue_name": venueName,
		})
	}
	if orders == nil {
		orders = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, orders)
}

func (h *TicketHandler) ListTickets(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := h.pool.Query(context.Background(),
		`SELECT t.id, t.order_id, t.ticket_type_id, t.event_id, t.code, t.status, t.used_at, t.created_at,
		        e.title, v.name, tt.name
		 FROM tickets t
		 JOIN orders o ON t.order_id = o.id
		 JOIN events e ON t.event_id = e.id
		 JOIN venues v ON e.venue_id = v.id
		 JOIN ticket_types tt ON t.ticket_type_id = tt.id
		 WHERE o.user_id = $1
		 ORDER BY t.created_at DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tickets"})
		return
	}
	defer rows.Close()

	var tickets []map[string]interface{}
	for rows.Next() {
		var id, orderID, ticketTypeID, eventID, code, status, eventTitle, venueName, ttName string
		var usedAt, createdAt interface{}
		if err := rows.Scan(&id, &orderID, &ticketTypeID, &eventID, &code, &status, &usedAt, &createdAt, &eventTitle, &venueName, &ttName); err != nil {
			continue
		}
		tickets = append(tickets, map[string]interface{}{
			"id": id, "order_id": orderID, "event_id": eventID, "code": code,
			"status": status, "used_at": usedAt, "created_at": createdAt,
			"event_title": eventTitle, "venue_name": venueName, "ticket_type": ttName,
		})
	}
	if tickets == nil {
		tickets = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, tickets)
}
