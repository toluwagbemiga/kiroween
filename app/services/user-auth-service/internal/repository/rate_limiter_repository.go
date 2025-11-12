package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiterRepository defines the interface for rate limiting
type RateLimiterRepository interface {
	RecordFailedAttempt(ctx context.Context, email string) error
	GetFailedAttempts(ctx context.Context, email string) (int, error)
	ResetAttempts(ctx context.Context, email string) error
	IsLocked(ctx context.Context, email string) (bool, time.Duration, error)
	LockAccount(ctx context.Context, email string, duration time.Duration) error
}

// rateLimiterRepository implements RateLimiterRepository
type rateLimiterRepository struct {
	client *redis.Client
}

// NewRateLimiterRepository creates a new rate limiter repository
func NewRateLimiterRepository(client *redis.Client) RateLimiterRepository {
	return &rateLimiterRepository{client: client}
}

// RecordFailedAttempt records a failed login attempt
func (r *rateLimiterRepository) RecordFailedAttempt(ctx context.Context, email string) error {
	key := fmt.Sprintf("ratelimit:login:%s", email)
	
	// Increment counter
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	
	// Set expiration on first attempt (15 minutes)
	if count == 1 {
		r.client.Expire(ctx, key, 15*time.Minute)
	}
	
	return nil
}

// GetFailedAttempts gets the number of failed attempts
func (r *rateLimiterRepository) GetFailedAttempts(ctx context.Context, email string) (int, error) {
	key := fmt.Sprintf("ratelimit:login:%s", email)
	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}

// ResetAttempts resets the failed attempt counter
func (r *rateLimiterRepository) ResetAttempts(ctx context.Context, email string) error {
	key := fmt.Sprintf("ratelimit:login:%s", email)
	return r.client.Del(ctx, key).Err()
}

// IsLocked checks if an account is locked
func (r *rateLimiterRepository) IsLocked(ctx context.Context, email string) (bool, time.Duration, error) {
	key := fmt.Sprintf("ratelimit:locked:%s", email)
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}
	
	if ttl > 0 {
		return true, ttl, nil
	}
	
	return false, 0, nil
}

// LockAccount locks an account for a duration
func (r *rateLimiterRepository) LockAccount(ctx context.Context, email string, duration time.Duration) error {
	key := fmt.Sprintf("ratelimit:locked:%s", email)
	return r.client.Set(ctx, key, "1", duration).Err()
}
