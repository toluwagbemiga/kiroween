package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the service configuration
type Config struct {
	Server    ServerConfig
	Analytics AnalyticsConfig
	Logging   LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	GRPCPort int
	Host     string
}

// AnalyticsConfig holds analytics configuration
type AnalyticsConfig struct {
	MixpanelAPIKey    string
	BatchSize         int
	FlushIntervalSec  int
	TestMode          bool
	MaxRetryAttempts  int
	InitialRetryDelay int
	MaxRetryDelay     int
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
			GRPCPort: getEnvInt("GRPC_PORT", 50054),
			Host:     getEnv("HOST", "0.0.0.0"),
		},
		Analytics: AnalyticsConfig{
			MixpanelAPIKey:    getEnv("MIXPANEL_API_KEY", ""),
			BatchSize:         getEnvInt("BATCH_SIZE", 50),
			FlushIntervalSec:  getEnvInt("FLUSH_INTERVAL_SECONDS", 10),
			TestMode:          getEnvBool("TEST_MODE", false),
			MaxRetryAttempts:  getEnvInt("MAX_RETRY_ATTEMPTS", 5),
			InitialRetryDelay: getEnvInt("INITIAL_RETRY_DELAY_MS", 1000),
			MaxRetryDelay:     getEnvInt("MAX_RETRY_DELAY_MS", 30000),
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
	// Validate Mixpanel API key (unless in test mode)
	if !c.Analytics.TestMode && c.Analytics.MixpanelAPIKey == "" {
		return fmt.Errorf("MIXPANEL_API_KEY is required (or enable TEST_MODE)")
	}

	// Validate batch size
	if c.Analytics.BatchSize < 1 || c.Analytics.BatchSize > 1000 {
		return fmt.Errorf("BATCH_SIZE must be between 1 and 1000")
	}

	// Validate flush interval
	if c.Analytics.FlushIntervalSec < 1 || c.Analytics.FlushIntervalSec > 300 {
		return fmt.Errorf("FLUSH_INTERVAL_SECONDS must be between 1 and 300")
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

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
