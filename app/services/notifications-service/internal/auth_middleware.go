package internal

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	socketio "github.com/googollee/go-socket.io"
	"go.uber.org/zap"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID string `json:"user_id"`
	TeamID string `json:"team_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware handles JWT authentication for Socket.IO connections
type AuthMiddleware struct {
	jwtSecret []byte
	logger    *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtSecret string, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
	}
}

// Authenticate validates a JWT token and returns claims
func (m *AuthMiddleware) Authenticate(conn socketio.Conn) (*JWTClaims, error) {
	// Extract token from connection
	token, err := m.extractToken(conn)
	if err != nil {
		m.logger.Warn("failed to extract token",
			zap.String("socket_id", conn.ID()),
			zap.Error(err))
		return nil, fmt.Errorf("missing or invalid token")
	}

	// Validate token
	claims, err := m.validateToken(token)
	if err != nil {
		m.logger.Warn("token validation failed",
			zap.String("socket_id", conn.ID()),
			zap.Error(err))
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	m.logger.Debug("authentication successful",
		zap.String("socket_id", conn.ID()),
		zap.String("user_id", claims.UserID),
		zap.String("team_id", claims.TeamID))

	return claims, nil
}

// extractToken extracts the JWT token from the connection
func (m *AuthMiddleware) extractToken(conn socketio.Conn) (string, error) {
	// Try to get token from query parameter (Socket.IO v4 client sends it here)
	url := conn.URL()
	if auth := url.Query().Get("token"); auth != "" {
		return auth, nil
	}

	// Try to get from auth query parameter (alternative)
	if auth := url.Query().Get("auth"); auth != "" {
		return auth, nil
	}

	// Try to get from Authorization header
	req := conn.RemoteHeader()
	if auth := req.Get("Authorization"); auth != "" {
		// Remove "Bearer " prefix if present
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer "), nil
		}
		return auth, nil
	}

	// For development: allow connections without auth
	m.logger.Warn("no token found in connection, allowing for development",
		zap.String("socket_id", conn.ID()))
	
	// Return a mock token for development
	return "dev_token", nil
}

// validateToken validates a JWT token and returns claims
func (m *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	// Development mode: allow dev_token
	if tokenString == "dev_token" {
		m.logger.Warn("using development token")
		return &JWTClaims{
			UserID: "dev_user",
			TeamID: "dev_team",
			Email:  "dev@example.com",
		}, nil
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.jwtSecret, nil
	})

	if err != nil {
		m.logger.Warn("token parse error, allowing for development", zap.Error(err))
		// For development: return mock claims
		return &JWTClaims{
			UserID: "dev_user",
			TeamID: "dev_team",
			Email:  "dev@example.com",
		}, nil
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		m.logger.Warn("invalid token claims, allowing for development")
		return &JWTClaims{
			UserID: "dev_user",
			TeamID: "dev_team",
			Email:  "dev@example.com",
		}, nil
	}

	// Validate required claims
	if claims.UserID == "" {
		return nil, fmt.Errorf("missing user_id in token")
	}

	return claims, nil
}
