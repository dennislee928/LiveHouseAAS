package payment

import (
	"context"
	"fmt"
)

type NewebPayProvider struct {
	MerchantID  string
	HashKey     string
	HashIV      string
}

func NewNewebPayProvider(merchantID, hashKey, hashIV string) *NewebPayProvider {
	return &NewebPayProvider{
		MerchantID: merchantID,
		HashKey:    hashKey,
		HashIV:     hashIV,
	}
}

func (p *NewebPayProvider) Name() Provider {
	return ProviderNewebPay
}

func (p *NewebPayProvider) Pay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	return nil, fmt.Errorf("NewebPay: not implemented yet")
}

func (p *NewebPayProvider) VerifyCallback(ctx context.Context, providerTxID string, status string) (*PaymentResult, error) {
	return nil, fmt.Errorf("NewebPay: not implemented yet")
}

func (p *NewebPayProvider) Refund(ctx context.Context, providerTxID string, amount int32) error {
	return fmt.Errorf("NewebPay: not implemented yet")
}
