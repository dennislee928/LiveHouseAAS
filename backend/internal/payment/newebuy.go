package payment

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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

type NewebPayProvider struct {
	MerchantID  string
	HashKey     string
	HashIV      string
	Sandbox     bool
	httpClient  *http.Client
}

func NewNewebPayProvider(merchantID, hashKey, hashIV string, sandbox bool) *NewebPayProvider {
	return &NewebPayProvider{
		MerchantID: merchantID,
		HashKey:    hashKey,
		HashIV:     hashIV,
		Sandbox:    sandbox,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *NewebPayProvider) Name() Provider { return ProviderNewebPay }

func (p *NewebPayProvider) Pay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	merchantOrderNo := fmt.Sprintf("LH%s%d", req.OrderID[:8], time.Now().Unix())
	timeStamp := time.Now().Unix()

	tradeInfo := map[string]interface{}{
		"MerchantID":      p.MerchantID,
		"MerchantOrderNo": merchantOrderNo,
		"Amt":             req.Amount,
		"ItemDesc":        "LiveHouse Ticket Purchase",
		"Email":           "",
		"TimeStamp":       timeStamp,
		"Version":         "2.0",
		"ReturnURL":       req.ReturnURL,
		"NotifyURL":       req.CallbackURL,
		"ClientBackURL":   req.ReturnURL,
	}

	tradeInfoJSON, _ := json.Marshal(tradeInfo)
	encrypted := p.encryptAES(string(tradeInfoJSON))

	// Generate SHA256 checksum
	checkValue := fmt.Sprintf("HashKey=%s&%s&HashIV=%s", p.HashKey, encrypted, p.HashIV)
	sha256Hash := p.sha256Hex(checkValue)

	var baseURL string
	if p.Sandbox {
		baseURL = "https://ccore.newebpay.com/MPG/mpg_gateway"
	} else {
		baseURL = "https://core.newebpay.com/MPG/mpg_gateway"
	}

	return &PaymentResult{
		Provider:    ProviderNewebPay,
		Status:      "pending",
		RedirectURL: baseURL,
		PaymentURL:  fmt.Sprintf("%s?MerchantID=%s&TradeInfo=%s&TradeSha=%s&Version=2.0", baseURL, p.MerchantID, encrypted, sha256Hash),
	}, nil
}

func (p *NewebPayProvider) encryptAES(plaintext string) string {
	key := []byte(p.HashKey)
	iv := []byte(p.HashIV)

	block, _ := aes.NewCipher(key)
	padded := pkcs7Pad([]byte(plaintext), aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func (p *NewebPayProvider) decryptAES(ciphertext string) (string, error) {
	key := []byte(p.HashKey)
	iv := []byte(p.HashIV)

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, _ := aes.NewCipher(key)
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	// Remove PKCS7 padding
	padLen := int(decrypted[len(decrypted)-1])
	return string(decrypted[:len(decrypted)-padLen]), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	padding := make([]byte, padLen)
	for i := range padding {
		padding[i] = byte(padLen)
	}
	return append(data, padding...)
}

func (p *NewebPayProvider) sha256Hex(input string) string {
	hash := sha256.Sum256([]byte(input))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

type NewebPayCallback struct {
	Status    string `json:"Status"`
	Message   string `json:"Message"`
	Result    string `json:"Result"`
	MerchantID string `json:"MerchantID"`
	TradeInfo  string `json:"TradeInfo"`
	TradeSha   string `json:"TradeSha"`
	Version    string `json:"Version"`
}

func (p *NewebPayProvider) VerifyCallback(ctx context.Context, providerTxID string, status string) (*PaymentResult, error) {
	return &PaymentResult{
		Provider:     ProviderNewebPay,
		Status:       status,
		ProviderTxID: providerTxID,
	}, nil
}

func (p *NewebPayProvider) Refund(ctx context.Context, providerTxID string, amount int32) error {
	timeStamp := time.Now().Unix()

	// Build refund API call
	postData := url.Values{}
	postData.Set("MerchantID", p.MerchantID)
	postData.Set("RespondType", "JSON")
	postData.Set("TimeStamp", fmt.Sprintf("%d", timeStamp))
	postData.Set("Version", "2.0")
	postData.Set("Amt", fmt.Sprintf("%d", amount))
	postData.Set("TradeNo", providerTxID)

	var baseURL string
	if p.Sandbox {
		baseURL = "https://ccore.newebpay.com/API/CreditCard/Close"
	} else {
		baseURL = "https://core.newebpay.com/API/CreditCard/Close"
	}

	resp, err := p.httpClient.PostForm(baseURL, postData)
	if err != nil {
		return fmt.Errorf("NewebPay refund request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("NewebPay refund parse failed: %w", err)
	}

	return nil
}
