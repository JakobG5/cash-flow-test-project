package models

type Config struct {
	Server   ServerConfig
	Logger   LoggerConfig
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
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
