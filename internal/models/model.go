package models

import "time"

type Config struct {
	Server     ServerConfig
	Logger     LoggerConfig
	Database   DatabaseConfig
	RabbitMQ   RabbitMQConfig
	APIKeyHash string
}

type ServerConfig struct {
	Port string
}

type LoggerConfig struct {
	Level string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
}

type CreateMerchantRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
	Email string `json:"email" validate:"required,email,max=255" example:"john.doe@example.com"`
}

type CreateMerchantResponse struct {
	Status     bool   `json:"status" example:"true"`
	MerchantID string `json:"merchant_id,omitempty" example:"CASM-ABC123"`
	Name       string `json:"name,omitempty" example:"John Doe"`
	Email      string `json:"email,omitempty" example:"john.doe@example.com"`
	APIKey     string `json:"api_key,omitempty" example:"cash_test_abc123def456"`
	Message    string `json:"message" example:"Merchant created successfully"`
}

type ErrorResponse struct {
	Status  bool     `json:"status"`
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

type MerchantBalance struct {
	Currency              string `json:"currency"`
	AvailableBalance      string `json:"available_balance"`
	TotalDeposit          string `json:"total_deposit"`
	TotalTransactionCount int32  `json:"total_transaction_count"`
	LastUpdated           string `json:"last_updated"`
}

type MerchantTransaction struct {
	ID                  string `json:"id"`
	PaymentIntentID     string `json:"payment_intent_id"`
	MerchantID          string `json:"merchant_id"`
	Amount              string `json:"amount"`
	Currency            string `json:"currency"`
	Status              string `json:"status"`
	ThirdPartyReference string `json:"third_party_reference,omitempty"`
	PaymentMethod       string `json:"payment_method,omitempty"`
	FeeAmount           string `json:"fee_amount,omitempty"`
	AccountNumber       string `json:"account_number,omitempty"`
	ProcessedAt         string `json:"processed_at,omitempty"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type GetMerchantResponse struct {
	Status         bool                  `json:"status"`
	MerchantID     string                `json:"merchant_id"`
	Name           string                `json:"name"`
	Email          string                `json:"email"`
	MerchantStatus string                `json:"merchant_status"`
	APIKey         string                `json:"api_key"`
	APIKeyStatus   string                `json:"api_key_status"`
	CreatedAt      string                `json:"created_at"`
	APIKeyCreated  string                `json:"api_key_created"`
	Balances       []MerchantBalance     `json:"balances,omitempty"`
	Transactions   []MerchantTransaction `json:"transactions,omitempty"`
	Message        string                `json:"message"`
}

type CreatePaymentIntentRequest struct {
	Amount      float64                `json:"amount" validate:"required,gt=0,lte=100000" example:"100.50"`
	Currency    string                 `json:"currency" validate:"required,len=3,oneof=ETB USD" example:"ETB"`
	Description string                 `json:"description,omitempty" validate:"max=500" example:"Payment for order #123"`
	CallbackURL string                 `json:"callback_url" validate:"required,url" example:"https://example.com/callback"`
	Nonce       string                 `json:"nonce" validate:"required,min=16,max=64" example:"unique_nonce_123456789"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type CreatePaymentIntentResponse struct {
	Status          bool      `json:"status" example:"true"`
	PaymentIntentID string    `json:"payment_intent_id" example:"PI-ABC123"`
	Amount          float64   `json:"amount" example:"100.5"`
	Currency        string    `json:"currency" example:"ETB"`
	PaymentStatus   string    `json:"payment_status" example:"pending"`
	Description     string    `json:"description,omitempty" example:"Payment for order #123"`
	CreatedAt       time.Time `json:"created_at" example:"2024-01-05T10:30:00Z"`
	ExpiresAt       time.Time `json:"expires_at" example:"2024-01-05T10:45:00Z"`
	Message         string    `json:"message" example:"Payment intent created successfully"`
}
