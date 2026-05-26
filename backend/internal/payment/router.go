package payment

import (
	"context"
	"fmt"
)

type Router struct {
	providers map[Provider]PaymentProvider
}

func NewRouter() *Router {
	r := &Router{
		providers: make(map[Provider]PaymentProvider),
	}
	r.Register(NewMockProvider())
	return r
}

func (r *Router) Register(p PaymentProvider) {
	r.providers[p.Name()] = p
}

func (r *Router) Pay(ctx context.Context, provider Provider, req *PaymentRequest) (*PaymentResult, error) {
	p, ok := r.providers[provider]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, provider)
	}
	return p.Pay(ctx, req)
}

func (r *Router) VerifyCallback(ctx context.Context, provider Provider, providerTxID, status string) (*PaymentResult, error) {
	p, ok := r.providers[provider]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, provider)
	}
	return p.VerifyCallback(ctx, providerTxID, status)
}

func (r *Router) Refund(ctx context.Context, provider Provider, providerTxID string, amount int32) error {
	p, ok := r.providers[provider]
	if !ok {
		return fmt.Errorf("%w: %s", ErrProviderNotFound, provider)
	}
	return p.Refund(ctx, providerTxID, amount)
}
