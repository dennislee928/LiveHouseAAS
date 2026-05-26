package payment

import (
	"context"
	"fmt"
)

type ECPayProvider struct {
	MerchantID string
	HashKey    string
	HashIV     string
}

func NewECPayProvider(merchantID, hashKey, hashIV string) *ECPayProvider {
	return &ECPayProvider{
		MerchantID: merchantID,
		HashKey:    hashKey,
		HashIV:     hashIV,
	}
}

func (p *ECPayProvider) Name() Provider {
	return ProviderECPay
}

func (p *ECPayProvider) Pay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	return nil, fmt.Errorf("ECPay: not implemented yet")
}

func (p *ECPayProvider) VerifyCallback(ctx context.Context, providerTxID string, status string) (*PaymentResult, error) {
	return nil, fmt.Errorf("ECPay: not implemented yet")
}

func (p *ECPayProvider) Refund(ctx context.Context, providerTxID string, amount int32) error {
	return fmt.Errorf("ECPay: not implemented yet")
}
