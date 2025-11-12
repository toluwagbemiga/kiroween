package domain

import (
	"time"
)

// Session represents an active user session stored in Redis
type Session struct {
	SessionID    string    `json:"session_id"`
	UserID       string    `json:"user_id"`
	TokenJTI     string    `json:"token_jti"` // JWT ID for revocation
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastActivity time.Time `json:"last_activity"`
}
