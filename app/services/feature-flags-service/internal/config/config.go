package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the service configuration
type Config struct {
	Server  ServerConfig
	Unleash UnleashConfig
	Logging LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	GRPCPort int
	Host     string
}

// UnleashConfig holds Unleash SDK configuration
type UnleashConfig struct {
	ServerURL       string
	APIToken        string
	AppName         string
	InstanceID      string
	RefreshInterval time.Duration
	MetricsInterval time.Duration
	DisableMetrics  bool
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
			GRPCPort: getEnvInt("GRPC_PORT", 50056),
			Host:     getEnv("HOST", "0.0.0.0"),
		},
		Unleash: UnleashConfig{
			ServerURL:       getEnv("UNLEASH_SERVER_URL", ""),
			APIToken:        getEnv("UNLEASH_API_TOKEN", ""),
			AppName:         getEnv("UNLEASH_APP_NAME", "feature-flags-service"),
			InstanceID:      getEnv("UNLEASH_INSTANCE_ID", "feature-flags-service-1"),
			RefreshInterval: time.Duration(getEnvInt("UNLEASH_REFRESH_INTERVAL_SECONDS", 10)) * time.Second,
			MetricsInterval: time.Duration(getEnvInt("UNLEASH_METRICS_INTERVAL_SECONDS", 60)) * time.Second,
			DisableMetrics:  getEnvBool("UNLEASH_DISABLE_METRICS", false),
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
	// Validate Unleash server URL
	if c.Unleash.ServerURL == "" {
		return fmt.Errorf("UNLEASH_SERVER_URL is required")
	}

	// Validate Unleash API token
	if c.Unleash.APIToken == "" {
		return fmt.Errorf("UNLEASH_API_TOKEN is required")
	}

	// Validate refresh interval
	if c.Unleash.RefreshInterval < 1*time.Second {
		return fmt.Errorf("UNLEASH_REFRESH_INTERVAL_SECONDS must be at least 1")
	}

	// Validate metrics interval
	if c.Unleash.MetricsInterval < 10*time.Second {
		return fmt.Errorf("UNLEASH_METRICS_INTERVAL_SECONDS must be at least 10")
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
