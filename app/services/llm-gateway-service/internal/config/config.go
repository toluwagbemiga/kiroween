package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the service configuration
type Config struct {
	Server    ServerConfig
	Prompts   PromptsConfig
	LLM       LLMConfig
	Analytics AnalyticsConfig
	Logging   LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	GRPCPort int
	Host     string
}

// PromptsConfig holds prompts configuration
type PromptsConfig struct {
	Directory  string
	WatchMode  bool
}

// LLMConfig holds LLM provider configuration
type LLMConfig struct {
	OpenAIAPIKey       string
	DefaultProvider    string
	DefaultModel       string
	DefaultTimeout     int
	MaxTimeout         int
	TestMode           bool
	MaxRetryAttempts   int
	InitialRetryDelayMs int
	MaxRetryDelayMs    int
}

// AnalyticsConfig holds analytics configuration
type AnalyticsConfig struct {
	ServiceAddr      string
	UsageStoreMaxSize int
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
			GRPCPort: getEnvInt("GRPC_PORT", 50053),
			Host:     getEnv("HOST", "0.0.0.0"),
		},
		Prompts: PromptsConfig{
			Directory: getEnv("PROMPTS_DIR", "/app/prompts"),
			WatchMode: getEnvBool("WATCH_PROMPTS", true),
		},
		LLM: LLMConfig{
			OpenAIAPIKey:       getEnv("OPENAI_API_KEY", ""),
			DefaultProvider:    getEnv("DEFAULT_PROVIDER", "openai"),
			DefaultModel:       getEnv("DEFAULT_MODEL", "gpt-4-turbo-preview"),
			DefaultTimeout:     getEnvInt("DEFAULT_TIMEOUT_SECONDS", 30),
			MaxTimeout:         getEnvInt("MAX_TIMEOUT_SECONDS", 120),
			TestMode:           getEnvBool("TEST_MODE", false),
			MaxRetryAttempts:   getEnvInt("MAX_RETRY_ATTEMPTS", 3),
			InitialRetryDelayMs: getEnvInt("INITIAL_RETRY_DELAY_MS", 1000),
			MaxRetryDelayMs:    getEnvInt("MAX_RETRY_DELAY_MS", 10000),
		},
		Analytics: AnalyticsConfig{
			ServiceAddr:      getEnv("ANALYTICS_SERVICE_ADDR", "analytics-service:50051"),
			UsageStoreMaxSize: getEnvInt("USAGE_STORE_MAX_SIZE", 10000),
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
	// Validate OpenAI API key (unless in test mode)
	if !c.LLM.TestMode && c.LLM.OpenAIAPIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is required (or enable TEST_MODE)")
	}

	// Validate prompts directory
	if c.Prompts.Directory == "" {
		return fmt.Errorf("PROMPTS_DIR is required")
	}

	// Validate timeouts
	if c.LLM.DefaultTimeout < 5 || c.LLM.DefaultTimeout > c.LLM.MaxTimeout {
		return fmt.Errorf("invalid timeout configuration")
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
