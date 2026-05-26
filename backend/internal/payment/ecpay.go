package payment

import (
	"context"
	"crypto/hmac"
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

type ECPayProvider struct {
	MerchantID  string
	HashKey     string
	HashIV      string
	Sandbox     bool
	httpClient  *http.Client
}

func NewECPayProvider(merchantID, hashKey, hashIV string, sandbox bool) *ECPayProvider {
	return &ECPayProvider{
		MerchantID: merchantID,
		HashKey:    hashKey,
		HashIV:     hashIV,
		Sandbox:    sandbox,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *ECPayProvider) Name() Provider { return ProviderECPay }

func (p *ECPayProvider) Pay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	merchantTradeNo := fmt.Sprintf("LH%s%d", req.OrderID[:8], time.Now().Unix())
	totalAmount := req.Amount
	tradeDesc := req.Description
	if tradeDesc == "" {
		tradeDesc = "LiveHouseAAS Ticket"
	}

	form := url.Values{}
	form.Set("MerchantID", p.MerchantID)
	form.Set("MerchantTradeNo", merchantTradeNo)
	form.Set("MerchantTradeDate", time.Now().Format("2006/01/02 15:04:05"))
	form.Set("PaymentType", "aio")
	form.Set("TotalAmount", fmt.Sprintf("%d", totalAmount))
	form.Set("TradeDesc", tradeDesc)
	form.Set("ItemName", "LiveHouse Ticket Purchase")
	form.Set("ReturnURL", req.ReturnURL)
	form.Set("ChoosePayment", "ALL")
	form.Set("EncryptType", "1")
	if req.CallbackURL != "" {
		form.Set("OrderResultURL", req.CallbackURL)
	}
	form.Set("ClientBackURL", req.ReturnURL)

	// Generate CheckMacValue
	checkMac := p.generateCheckMacValue(form)
	form.Set("CheckMacValue", checkMac)

	var baseURL string
	if p.Sandbox {
		baseURL = "https://payment-stage.ecpay.com.tw/Cashier/AioCheckOut/V5"
	} else {
		baseURL = "https://payment.ecpay.com.tw/Cashier/AioCheckOut/V5"
	}

	return &PaymentResult{
		Provider:    ProviderECPay,
		Status:      "pending",
		RedirectURL: baseURL,
		PaymentURL:  baseURL + "?" + form.Encode(),
	}, nil
}

func (p *ECPayProvider) generateCheckMacValue(form url.Values) string {
	// Sort parameters by key
	keys := make([]string, 0, len(form))
	for k := range form {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("HashKey=%s", p.HashKey))
	for _, k := range keys {
		v := form.Get(k)
		if v != "" && k != "CheckMacValue" {
			sb.WriteString(fmt.Sprintf("&%s=%s", k, v))
		}
	}
	sb.WriteString(fmt.Sprintf("&HashIV=%s", p.HashIV))

	raw := sb.String()
	// URL decode once
	raw = url.QueryEscape(raw)

	mac := hmac.New(sha256.New, []byte(p.HashKey))
	mac.Write([]byte(raw))
	checksum := strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))

	return checksum
}

func (p *ECPayProvider) VerifyCheckMacValue(params map[string]string) bool {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	expected := p.generateCheckMacValue(form)
	return expected == params["CheckMacValue"]
}

type ECPayCallback struct {
	MerchantTradeNo string `json:"MerchantTradeNo"`
	TradeNo         string `json:"TradeNo"`
	RtnCode         int    `json:"RtnCode"`
	RtnMsg          string `json:"RtnMsg"`
	TradeAmt        int32  `json:"TradeAmt"`
	PaymentDate     string `json:"PaymentDate"`
	PaymentType     string `json:"PaymentType"`
	CheckMacValue   string `json:"CheckMacValue"`
}

func (p *ECPayProvider) VerifyCallback(ctx context.Context, providerTxID string, status string) (*PaymentResult, error) {
	return &PaymentResult{
		Provider:     ProviderECPay,
		Status:       status,
		ProviderTxID: providerTxID,
	}, nil
}

func (p *ECPayProvider) Refund(ctx context.Context, providerTxID string, amount int32) error {
	merchantTradeNo := fmt.Sprintf("RF%s%d", providerTxID[:8], time.Now().Unix())

	form := url.Values{}
	form.Set("MerchantID", p.MerchantID)
	form.Set("MerchantTradeNo", merchantTradeNo)
	form.Set("TradeNo", providerTxID)
	form.Set("RefundAmount", fmt.Sprintf("%d", amount))
	form.Set("TotalAmount", fmt.Sprintf("%d", amount))
	form.Set("CheckMacValue", p.generateCheckMacValue(form))

	var baseURL string
	if p.Sandbox {
		baseURL = "https://payment-stage.ecpay.com.tw/Cashier/AioCheckOut/V5"
	} else {
		baseURL = "https://payment.ecpay.com.tw/Cashier/AioCheckOut/V5"
	}

	resp, err := p.httpClient.PostForm(baseURL, form)
	if err != nil {
		return fmt.Errorf("ECPay refund request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("ECPay refund parse failed: %w", err)
	}

	if rtnCode, ok := result["RtnCode"].(float64); ok && int(rtnCode) == 1 {
		return nil
	}

	return fmt.Errorf("ECPay refund failed: %v", result["RtnMsg"])
}
