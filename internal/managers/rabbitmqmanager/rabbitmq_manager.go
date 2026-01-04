package rabbitmqmanager

type PaymentMessage struct {
	PaymentIntentID string `json:"payment_intent_id"`
	TransactionID   string `json:"transaction_id"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Timestamp       string `json:"timestamp"`
}

type IRabbitMQManager interface {
	Close() error
	HealthCheck() error
	PublishPaymentIntent(message PaymentMessage) error
}
