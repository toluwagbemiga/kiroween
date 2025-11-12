package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the gateway configuration
type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	Auth     AuthConfig
	Logging  LoggingConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port int
	Host string
	Env  string
}

// ServicesConfig holds gRPC service addresses
type ServicesConfig struct {
	UserAuthService      string
	BillingService       string
	LLMGatewayService    string
	NotificationsService string
	AnalyticsService     string
	FeatureFlagsService  string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvInt("PORT", 8080),
			Host: getEnv("HOST", "0.0.0.0"),
			Env:  getEnv("ENV", "development"),
		},
		Services: ServicesConfig{
			UserAuthService:      getEnv("USER_AUTH_SERVICE", "localhost:50051"),
			BillingService:       getEnv("BILLING_SERVICE", "localhost:50052"),
			LLMGatewayService:    getEnv("LLM_GATEWAY_SERVICE", "localhost:50053"),
			NotificationsService: getEnv("NOTIFICATIONS_SERVICE", "localhost:50054"),
			AnalyticsService:     getEnv("ANALYTICS_SERVICE", "localhost:50055"),
			FeatureFlagsService:  getEnv("FEATURE_FLAGS_SERVICE", "localhost:50056"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", ""),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Auth.JWTSecret == "" && c.Server.Env == "production" {
		return fmt.Errorf("JWT_SECRET is required in production")
	}

	return nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
