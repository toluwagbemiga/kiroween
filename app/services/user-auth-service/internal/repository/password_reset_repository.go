package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// PasswordResetRepository defines the interface for password reset tokens
type PasswordResetRepository interface {
	CreateResetToken(ctx context.Context, token *PasswordResetToken, ttl time.Duration) (string, error)
	GetResetToken(ctx context.Context, token string) (*PasswordResetToken, error)
	DeleteResetToken(ctx context.Context, token string) error
}

// passwordResetRepository implements PasswordResetRepository
type passwordResetRepository struct {
	client *redis.Client
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(client *redis.Client) PasswordResetRepository {
	return &passwordResetRepository{client: client}
}

// CreateResetToken creates a password reset token
func (r *passwordResetRepository) CreateResetToken(ctx context.Context, token *PasswordResetToken, ttl time.Duration) (string, error) {
	// Generate a secure token (this will be hashed before storage)
	// The actual token generation happens in the service layer
	// Here we just store it
	
	data, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	
	// Hash the token for storage
	hash := sha256.Sum256([]byte(token.Email + token.UserID + token.CreatedAt.String()))
	hashedToken := hex.EncodeToString(hash[:])
	
	key := fmt.Sprintf("reset:%s", hashedToken)
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return "", err
	}
	
	return hashedToken, nil
}

// GetResetToken retrieves a password reset token
func (r *passwordResetRepository) GetResetToken(ctx context.Context, token string) (*PasswordResetToken, error) {
	key := fmt.Sprintf("reset:%s", token)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("reset token not found or expired")
	}
	if err != nil {
		return nil, err
	}
	
	var resetToken PasswordResetToken
	if err := json.Unmarshal([]byte(data), &resetToken); err != nil {
		return nil, err
	}
	
	return &resetToken, nil
}

// DeleteResetToken deletes a password reset token
func (r *passwordResetRepository) DeleteResetToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("reset:%s", token)
	return r.client.Del(ctx, key).Err()
}
