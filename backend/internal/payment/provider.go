package payment

import (
	"context"
	"errors"
)

var (
	ErrProviderNotFound = errors.New("payment provider not found")
	ErrPaymentFailed    = errors.New("payment failed")
)

type Provider string

const (
	ProviderMock    Provider = "mock"
	ProviderECPay   Provider = "ecpay"
	ProviderNewebPay Provider = "newebpay"
	ProviderCrypto  Provider = "crypto"
)

type PaymentRequest struct {
	OrderID     string `json:"order_id"`
	Amount      int32  `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	ReturnURL   string `json:"return_url"`
	CallbackURL string `json:"callback_url"`
}

type PaymentResult struct {
	Provider      Provider `json:"provider"`
	Status        string   `json:"status"`
	ProviderTxID  string   `json:"provider_tx_id"`
	RedirectURL   string   `json:"redirect_url,omitempty"`
	PaymentURL    string   `json:"payment_url,omitempty"`
	ErrorMessage  string   `json:"error_message,omitempty"`
}

type PaymentProvider interface {
	Name() Provider
	Pay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error)
	VerifyCallback(ctx context.Context, providerTxID string, status string) (*PaymentResult, error)
	Refund(ctx context.Context, providerTxID string, amount int32) error
}
