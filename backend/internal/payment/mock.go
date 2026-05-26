package payment

import (
	"context"
	"fmt"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) Name() Provider {
	return ProviderMock
}

func (p *MockProvider) Pay(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	return &PaymentResult{
		Provider:     ProviderMock,
		Status:       "completed",
		ProviderTxID: fmt.Sprintf("mock_tx_%s", req.OrderID),
		RedirectURL:  req.ReturnURL,
		PaymentURL:   req.ReturnURL,
	}, nil
}

func (p *MockProvider) VerifyCallback(ctx context.Context, providerTxID string, status string) (*PaymentResult, error) {
	return &PaymentResult{
		Provider:     ProviderMock,
		Status:       "completed",
		ProviderTxID: providerTxID,
	}, nil
}

func (p *MockProvider) Refund(ctx context.Context, providerTxID string, amount int32) error {
	return nil
}
