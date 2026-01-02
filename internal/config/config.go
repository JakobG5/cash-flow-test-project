package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Logger LoggerConfig
}

type ServerConfig struct {
	Port string
}

type LoggerConfig struct {
	Level string
}

func Load() (*Config, error) {
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %v", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnvAsString("SERVER_PORT", "8080"),
		},
		Logger: LoggerConfig{
			Level: getEnvAsString("LOG_LEVEL", "info"),
		},
	}

	return config, nil
}

func getEnvAsString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
