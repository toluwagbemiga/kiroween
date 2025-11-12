package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/haunted-saas/user-auth-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	Get(ctx context.Context, sessionID string) (*domain.Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteAllForUser(ctx context.Context, userID string) error
	ExtendExpiration(ctx context.Context, sessionID string, duration time.Duration) error
	IsRevoked(ctx context.Context, tokenJTI string) (bool, error)
	RevokeToken(ctx context.Context, tokenJTI string, expiresAt time.Time) error
}

// sessionRepository implements SessionRepository
type sessionRepository struct {
	client *redis.Client
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(client *redis.Client) SessionRepository {
	return &sessionRepository{client: client}
}

// Create creates a new session
func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) error {
	key := fmt.Sprintf("session:%s", session.SessionID)
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	
	ttl := time.Until(session.ExpiresAt)
	return r.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a session by ID
func (r *sessionRepository) Get(ctx context.Context, sessionID string) (*domain.Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, err
	}
	
	var session domain.Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}
	
	return &session, nil
}

// Delete deletes a session
func (r *sessionRepository) Delete(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return r.client.Del(ctx, key).Err()
}

// DeleteAllForUser deletes all sessions for a user
func (r *sessionRepository) DeleteAllForUser(ctx context.Context, userID string) error {
	// Scan for all session keys
	pattern := "session:*"
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	
	for iter.Next(ctx) {
		key := iter.Val()
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		
		var session domain.Session
		if err := json.Unmarshal([]byte(data), &session); err != nil {
			continue
		}
		
		if session.UserID == userID {
			r.client.Del(ctx, key)
		}
	}
	
	return iter.Err()
}

// ExtendExpiration extends the expiration of a session (sliding window)
func (r *sessionRepository) ExtendExpiration(ctx context.Context, sessionID string, duration time.Duration) error {
	// Get current session
	session, err := r.Get(ctx, sessionID)
	if err != nil {
		return err
	}
	
	// Update expiration and last activity
	session.ExpiresAt = time.Now().Add(duration)
	session.LastActivity = time.Now()
	
	// Save back to Redis
	return r.Create(ctx, session)
}

// IsRevoked checks if a token is revoked
func (r *sessionRepository) IsRevoked(ctx context.Context, tokenJTI string) (bool, error) {
	key := fmt.Sprintf("revoked:%s", tokenJTI)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// RevokeToken adds a token to the revocation list
func (r *sessionRepository) RevokeToken(ctx context.Context, tokenJTI string, expiresAt time.Time) error {
	key := fmt.Sprintf("revoked:%s", tokenJTI)
	ttl := time.Until(expiresAt)
	if ttl < 0 {
		ttl = 0
	}
	return r.client.Set(ctx, key, "1", ttl).Err()
}
