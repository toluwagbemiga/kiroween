package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config holds all configuration for the billing service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Stripe   StripeConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	GRPCPort int
	HTTPPort int
	Host     string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL            string
	MaxConnections int
}

// StripeConfig holds Stripe configuration
type StripeConfig struct {
	APIKey        string
	WebhookSecret string
	APIVersion    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	viper.AutomaticEnv()

	config := &Config{
		Server: ServerConfig{
			GRPCPort: getEnvAsInt("GRPC_PORT", 50052),
			HTTPPort: getEnvAsInt("HTTP_PORT", 8080),
			Host:     getEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			URL:            getEnv("DATABASE_URL", ""),
			MaxConnections: getEnvAsInt("DB_MAX_CONNECTIONS", 25),
		},
		Stripe: StripeConfig{
			APIKey:        getEnv("STRIPE_API_KEY", ""),
			WebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
			APIVersion:    getEnv("STRIPE_API_VERSION", "2023-10-16"),
		},
	}

	// Validate required configuration
	if config.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if config.Stripe.APIKey == "" {
		return nil, fmt.Errorf("STRIPE_API_KEY is required")
	}

	if config.Stripe.WebhookSecret == "" {
		return nil, fmt.Errorf("STRIPE_WEBHOOK_SECRET is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
