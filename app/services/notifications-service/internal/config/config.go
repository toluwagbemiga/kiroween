package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds the service configuration
type Config struct {
	Server         ServerConfig
	SocketIO       SocketIOConfig
	Authentication AuthConfig
	Logging        LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	GRPCPort     int
	SocketIOPort int
	Host         string
}

// SocketIOConfig holds Socket.IO configuration
type SocketIOConfig struct {
	AllowedOrigins     []string
	MaxConnections     int
	PingTimeoutSec     int
	PingIntervalSec    int
	EnableWebSocket    bool
	EnablePolling      bool
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
			GRPCPort:     getEnvInt("GRPC_PORT", 50055),
			SocketIOPort: getEnvInt("SOCKETIO_PORT", 3000),
			Host:         getEnv("HOST", "0.0.0.0"),
		},
		SocketIO: SocketIOConfig{
			AllowedOrigins:  parseAllowedOrigins(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
			MaxConnections:  getEnvInt("MAX_CONNECTIONS", 10000),
			PingTimeoutSec:  getEnvInt("PING_TIMEOUT_SECONDS", 60),
			PingIntervalSec: getEnvInt("PING_INTERVAL_SECONDS", 25),
			EnableWebSocket: getEnvBool("ENABLE_WEBSOCKET", true),
			EnablePolling:   getEnvBool("ENABLE_POLLING", true),
		},
		Authentication: AuthConfig{
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
	// Validate JWT secret
	if c.Authentication.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	// Validate at least one transport is enabled
	if !c.SocketIO.EnableWebSocket && !c.SocketIO.EnablePolling {
		return fmt.Errorf("at least one transport (WebSocket or Polling) must be enabled")
	}

	// Validate max connections
	if c.SocketIO.MaxConnections < 1 {
		return fmt.Errorf("MAX_CONNECTIONS must be at least 1")
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

func parseAllowedOrigins(originsStr string) []string {
	if originsStr == "" {
		return []string{}
	}
	origins := strings.Split(originsStr, ",")
	result := make([]string, 0, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
