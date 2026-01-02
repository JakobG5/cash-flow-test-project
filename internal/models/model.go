package models

// Config holds all application configuration settings
type Config struct {
	Server   ServerConfig   // HTTP server configuration
	Logger   LoggerConfig   // Logging configuration
	Database DatabaseConfig // PostgreSQL database configuration
	RabbitMQ RabbitMQConfig // RabbitMQ message queue configuration
}

// ServerConfig contains HTTP server related settings
type ServerConfig struct {
	Port string // Port for the HTTP server to listen on
}

// LoggerConfig contains logging related settings
type LoggerConfig struct {
	Level string // Log level (debug, info, warn, error)
}

// DatabaseConfig contains PostgreSQL database connection settings
type DatabaseConfig struct {
	Host     string // Database host address
	Port     string // Database port
	User     string // Database username
	Password string // Database password
	DBName   string // Database name
	SSLMode  string // SSL mode for connection
}

// RabbitMQConfig contains RabbitMQ message queue settings
type RabbitMQConfig struct {
	Host     string // RabbitMQ host address
	Port     string // RabbitMQ port
	User     string // RabbitMQ username
	Password string // RabbitMQ password
	VHost    string // RabbitMQ virtual host
}
