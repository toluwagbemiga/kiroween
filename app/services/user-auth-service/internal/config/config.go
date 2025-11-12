package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Security SecurityConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	GRPCPort int
	Host     string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL             string
	MaxConnections  int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	Expiration     time.Duration
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	BcryptCost           int
	MaxLoginAttempts     int
	LockoutDuration      time.Duration
	PermissionCacheTTL   time.Duration
	SessionExpiration    time.Duration
	PasswordResetTTL     time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	viper.AutomaticEnv()

	config := &Config{
		Server: ServerConfig{
			GRPCPort: getEnvAsInt("GRPC_PORT", 50051),
			Host:     getEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", ""),
			MaxConnections:  getEnvAsInt("DB_MAX_CONNECTIONS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 5)) * time.Minute,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			PrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "/app/keys/jwt-private.pem"),
			PublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "/app/keys/jwt-public.pem"),
			Expiration:     time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
		},
		Security: SecurityConfig{
			BcryptCost:           getEnvAsInt("BCRYPT_COST", 12),
			MaxLoginAttempts:     getEnvAsInt("MAX_LOGIN_ATTEMPTS", 5),
			LockoutDuration:      time.Duration(getEnvAsInt("LOCKOUT_DURATION_MINUTES", 30)) * time.Minute,
			PermissionCacheTTL:   time.Duration(getEnvAsInt("PERMISSION_CACHE_TTL_MINUTES", 5)) * time.Minute,
			SessionExpiration:    time.Duration(getEnvAsInt("SESSION_EXPIRATION_HOURS", 24)) * time.Hour,
			PasswordResetTTL:     time.Duration(getEnvAsInt("PASSWORD_RESET_TTL_MINUTES", 60)) * time.Minute,
		},
	}

	// Validate required configuration
	if config.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
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
