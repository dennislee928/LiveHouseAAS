package payment

import (
	"context"
	"fmt"
)

type CryptoProvider struct {
	// BinancePay: API key, secret
	// EVM: RPC URL, wallet private key for receiving
	Mode string // "binance" or "evm"
}

func NewCryptoProvider(mode string) *CryptoProvider {
	return &CryptoProvider{Mode: mode}
}

func (p *CryptoProvider) Name() Provider {
	return ProviderCrypto
}

func (p *CryptoProvider) Pay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	return nil, fmt.Errorf("Crypto(%s): not implemented yet", p.Mode)
}

func (p *CryptoProvider) VerifyCallback(ctx context.Context, providerTxID string, status string) (*PaymentResult, error) {
	return nil, fmt.Errorf("Crypto(%s): not implemented yet", p.Mode)
}

func (p *CryptoProvider) Refund(ctx context.Context, providerTxID string, amount int32) error {
	return fmt.Errorf("Crypto(%s): not implemented yet", p.Mode)
}
