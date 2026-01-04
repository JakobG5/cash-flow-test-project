package worker

import (
	"context"
)

type PaymentIntentMessage struct {
	PaymentIntentID string `json:"payment_intent_id"`
	MerchantID      string `json:"merchant_id"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Timestamp       string `json:"timestamp"`
}

type IWorker interface {
	Start(ctx context.Context) error
	Stop() error
	ProcessPaymentIntent(ctx context.Context, message PaymentIntentMessage) error
}
