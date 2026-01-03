package configmanager

import (
	"cash-flow-financial/internal/models"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Load() (*models.Config, error) {

	viper.SetDefault("SERVER_PORT", "8080")

	viper.SetDefault("LOG_LEVEL", "info")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "cashflow_user")
	viper.SetDefault("DB_PASSWORD", "cashflow_pass")
	viper.SetDefault("DB_NAME", "cashflow_dev")
	viper.SetDefault("DB_SSL_MODE", "disable")

	viper.SetDefault("RABBITMQ_HOST", "localhost")
	viper.SetDefault("RABBITMQ_PORT", "5672")
	viper.SetDefault("RABBITMQ_USER", "guest")
	viper.SetDefault("RABBITMQ_PASSWORD", "guest")
	viper.SetDefault("RABBITMQ_VHOST", "/")

	viper.SetDefault("API_KEY_HASH_KEY", "cashflow_test_2024_secure_key_123456789")

	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: Could not read config file: %v", err)
		}
	}

	if err := validateConfig(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	config := &models.Config{
		Server: models.ServerConfig{
			Port: getEnvAsString("SERVER_PORT", "8080"),
		},
		Logger: models.LoggerConfig{
			Level: strings.ToLower(getEnvAsString("LOG_LEVEL", "info")),
		},
		Database: models.DatabaseConfig{
			Host:     getEnvAsString("DB_HOST", "localhost"),
			Port:     getEnvAsString("DB_PORT", "5432"),
			User:     getEnvAsString("DB_USER", "cashflow_user"),
			Password: getEnvAsString("DB_PASSWORD", "cashflow_pass"),
			DBName:   getEnvAsString("DB_NAME", "cashflow_dev"),
			SSLMode:  getEnvAsString("DB_SSL_MODE", "disable"),
		},
		RabbitMQ: models.RabbitMQConfig{
			Host:     getEnvAsString("RABBITMQ_HOST", "localhost"),
			Port:     getEnvAsString("RABBITMQ_PORT", "5672"),
			User:     getEnvAsString("RABBITMQ_USER", "guest"),
			Password: getEnvAsString("RABBITMQ_PASSWORD", "guest"),
			VHost:    getEnvAsString("RABBITMQ_VHOST", "/"),
		},
		APIKeyHash: getEnvAsString("API_KEY_HASH_KEY", "cashflow_test_2024_secure_key_123456789"),
	}

	log.Printf("Configuration loaded successfully. Server will run on port %s", config.Server.Port)
	return config, nil
}

func validateConfig() error {
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	logLevel := strings.ToLower(viper.GetString("LOG_LEVEL"))
	if !validLogLevels[logLevel] && logLevel != "" {
		return fmt.Errorf("invalid LOG_LEVEL '%s', must be one of: debug, info, warn, error", logLevel)
	}

	dbPort := viper.GetString("DB_PORT")
	if dbPort != "" {
		if _, err := fmt.Sscanf(dbPort, "%d", new(int)); err != nil {
			return fmt.Errorf("invalid DB_PORT '%s', must be numeric", dbPort)
		}
	}

	serverPort := viper.GetString("SERVER_PORT")
	if serverPort != "" {
		if _, err := fmt.Sscanf(serverPort, "%d", new(int)); err != nil {
			return fmt.Errorf("invalid SERVER_PORT '%s', must be numeric", serverPort)
		}
	}

	return nil
}

func getEnvAsString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
