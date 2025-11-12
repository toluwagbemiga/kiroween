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
	// Try to get token from auth object (Socket.IO client sends it here)
	url := conn.URL()
	if auth := url.Query().Get("token"); auth != "" {
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

	return "", fmt.Errorf("no token found in connection")
}

// validateToken validates a JWT token and returns claims
func (m *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate required claims
	if claims.UserID == "" {
		return nil, fmt.Errorf("missing user_id in token")
	}

	return claims, nil
}
