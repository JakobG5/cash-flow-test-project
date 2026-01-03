package models

type Config struct {
	Server      ServerConfig
	Logger      LoggerConfig
	Database    DatabaseConfig
	RabbitMQ    RabbitMQConfig
	APIKeyHash  string
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
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email,max=255"`
}

type CreateMerchantResponse struct {
	Status     bool   `json:"status"`
	MerchantID string `json:"merchant_id,omitempty"`
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	APIKey     string `json:"api_key,omitempty"`
	Message    string `json:"message"`
}

type ErrorResponse struct {
	Status  bool     `json:"status"`
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

type GetMerchantResponse struct {
	Status        bool   `json:"status"`
	MerchantID    string `json:"merchant_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	MerchantStatus string `json:"merchant_status"`
	APIKey        string `json:"api_key"`
	APIKeyStatus  string `json:"api_key_status"`
	CreatedAt     string `json:"created_at"`
	APIKeyCreated string `json:"api_key_created"`
	Message       string `json:"message"`
}
