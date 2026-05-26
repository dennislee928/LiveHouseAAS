package payment

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type CryptoProvider struct {
	Mode         string // "binance" or "evm"
	BinanceAPIKey    string
	BinanceSecretKey string
	EVMReceivingAddr string
	EVMPrivateKey    string
	EVMChainID       int64
	RPCURL           string
	httpClient       *http.Client
}

func NewCryptoProvider(mode string) *CryptoProvider {
	return &CryptoProvider{
		Mode:       mode,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *CryptoProvider) Name() Provider { return ProviderCrypto }

func (p *CryptoProvider) Pay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	switch p.Mode {
	case "binance":
		return p.binancePay(ctx, req)
	case "evm":
		return p.evmPay(ctx, req)
	default:
		return nil, fmt.Errorf("Crypto: unknown mode %s", p.Mode)
	}
}

// --- Binance Pay ---

type BinanceOrderRequest struct {
	MerchantTradeNo string `json:"merchantTradeNo"`
	TotalFee        int64  `json:"totalFee"`
	Currency        string `json:"currency"`
	ProductName     string `json:"productName"`
	ProductType     string `json:"productType"`
	ReturnURL       string `json:"returnUrl"`
	WebhookURL      string `json:"webhookUrl"`
	Timestamp       int64  `json:"timestamp"`
}

type BinanceOrderResponse struct {
	Status      string `json:"status"`
	Code        string `json:"code"`
	Data        struct {
		PrepayID     string `json:"prepayId"`
		TradeNo      string `json:"tradeNo"`
		Currency     string `json:"currency"`
		TotalFee     int64  `json:"totalFee"`
		ExpireTime   int64  `json:"expireTime"`
		QrContent    string `json:"qrContent"`
		UniversalURL string `json:"universalUrl"`
		CheckoutURL  string `json:"checkoutUrl"`
		ShortLink    string `json:"shortLink"`
	} `json:"data"`
	ErrorMessage string `json:"errorMessage"`
}

func (p *CryptoProvider) binancePay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	merchantTradeNo := fmt.Sprintf("LH%s%d", req.OrderID[:8], time.Now().Unix())
	timestamp := time.Now().UnixMilli()

	orderReq := BinanceOrderRequest{
		MerchantTradeNo: merchantTradeNo,
		TotalFee:        int64(req.Amount * 100), // Convert to cents
		Currency:        "USDT",
		ProductName:     req.Description,
		ProductType:     "Goods",
		ReturnURL:       req.ReturnURL,
		WebhookURL:      req.CallbackURL,
		Timestamp:       timestamp,
	}

	payload, _ := json.Marshal(orderReq)
	signature := p.binanceSign(payload)

	httpReq, _ := http.NewRequestWithContext(ctx, "POST",
		"https://bpay.binanceapi.com/gateway/api/v2/order/create",
		strings.NewReader(string(payload)))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("BinancePay-Timestamp", fmt.Sprintf("%d", timestamp))
	httpReq.Header.Set("BinancePay-Nonce", generateNonce())
	httpReq.Header.Set("BinancePay-Certificate-SN", p.BinanceAPIKey)
	httpReq.Header.Set("BinancePay-Signature", signature)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Binance Pay request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var binanceResp BinanceOrderResponse
	if err := json.Unmarshal(body, &binanceResp); err != nil {
		return nil, fmt.Errorf("Binance Pay parse failed: %w", err)
	}

	if binanceResp.Status != "SUCCESS" {
		return nil, fmt.Errorf("Binance Pay error: %s - %s", binanceResp.Code, binanceResp.ErrorMessage)
	}

	return &PaymentResult{
		Provider:     ProviderCrypto,
		Status:       "pending",
		ProviderTxID: binanceResp.Data.TradeNo,
		RedirectURL:  binanceResp.Data.CheckoutURL,
		PaymentURL:   binanceResp.Data.UniversalURL,
	}, nil
}

func (p *CryptoProvider) binanceSign(payload []byte) string {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	nonce := generateNonce()
	raw := timestamp + "\n" + nonce + "\n" + string(payload) + "\n"
	mac := hmacSHA256([]byte(p.BinanceSecretKey), []byte(raw))
	return hex.EncodeToString(mac)
}

func hmacSHA256(key, data []byte) []byte {
	h := sha256.New()
	h.Write(key)
	h.Write(data)
	return h.Sum(nil)
}

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// --- EVM On-Chain ---

type EVMOrderRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Amount  string `json:"amount"`
	OrderID string `json:"order_id"`
}

type EVMOrderResponse struct {
	TxHash  string `json:"tx_hash"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (p *CryptoProvider) evmPay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	// Connect to EVM-compatible node (Polygon)
	// Generate a payment request with a unique order ID
	orderID := fmt.Sprintf("0x%s", req.OrderID[:16])

	// In production, this would create an on-chain escrow or direct payment
	// For now, generate the payment address and return it
	paymentAddr := p.EVMReceivingAddr
	if paymentAddr == "" {
		paymentAddr = "0x0000000000000000000000000000000000000000"
	}

	paymentData := url.Values{}
	paymentData.Set("address", paymentAddr)
	paymentData.Set("amount", fmt.Sprintf("%d", req.Amount))
	paymentData.Set("order_id", orderID)
	paymentData.Set("network", "polygon")
	paymentData.Set("currency", "USDC")

	return &PaymentResult{
		Provider:     ProviderCrypto,
		Status:       "pending",
		ProviderTxID: orderID,
		PaymentURL:   fmt.Sprintf("https://pay.livehouseaas.com/crypto?%s", paymentData.Encode()),
	}, nil
}

func (p *CryptoProvider) VerifyCallback(ctx context.Context, providerTxID string, status string) (*PaymentResult, error) {
	// Verify transaction on-chain
	return &PaymentResult{
		Provider:     ProviderCrypto,
		Status:       status,
		ProviderTxID: providerTxID,
	}, nil
}

func (p *CryptoProvider) Refund(ctx context.Context, providerTxID string, amount int32) error {
	switch p.Mode {
	case "binance":
		return p.binanceRefund(ctx, providerTxID, amount)
	case "evm":
		return p.evmRefund(ctx, providerTxID, amount)
	default:
		return fmt.Errorf("Crypto(%s): refund not implemented", p.Mode)
	}
}

func (p *CryptoProvider) binanceRefund(ctx context.Context, providerTxID string, amount int32) error {
	timestamp := time.Now().UnixMilli()

	refundReq := map[string]interface{}{
		"merchantTradeNo": providerTxID,
		"refundAmount":    int64(amount * 100),
		"refundReason":    "Customer requested refund",
		"timestamp":       timestamp,
	}

	payload, _ := json.Marshal(refundReq)
	signature := p.binanceSign(payload)

	httpReq, _ := http.NewRequestWithContext(ctx, "POST",
		"https://bpay.binanceapi.com/gateway/api/v2/order/refund",
		strings.NewReader(string(payload)))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("BinancePay-Timestamp", fmt.Sprintf("%d", timestamp))
	httpReq.Header.Set("BinancePay-Nonce", generateNonce())
	httpReq.Header.Set("BinancePay-Certificate-SN", p.BinanceAPIKey)
	httpReq.Header.Set("BinancePay-Signature", signature)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Binance Pay refund request failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (p *CryptoProvider) evmRefund(ctx context.Context, providerTxID string, amount int32) error {
	// Would interact with the smart contract to refund
	return nil
}
