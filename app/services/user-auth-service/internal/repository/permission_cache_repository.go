package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// PermissionCacheRepository defines the interface for permission caching
type PermissionCacheRepository interface {
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	SetUserPermissions(ctx context.Context, userID string, permissions []string, ttl time.Duration) error
	InvalidateUserPermissions(ctx context.Context, userID string) error
}

// permissionCacheRepository implements PermissionCacheRepository
type permissionCacheRepository struct {
	client *redis.Client
}

// NewPermissionCacheRepository creates a new permission cache repository
func NewPermissionCacheRepository(client *redis.Client) PermissionCacheRepository {
	return &permissionCacheRepository{client: client}
}

// GetUserPermissions retrieves cached permissions for a user
func (r *permissionCacheRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	key := fmt.Sprintf("permissions:%s", userID)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("permissions not cached")
	}
	if err != nil {
		return nil, err
	}
	
	var permissions []string
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, err
	}
	
	return permissions, nil
}

// SetUserPermissions caches permissions for a user
func (r *permissionCacheRepository) SetUserPermissions(ctx context.Context, userID string, permissions []string, ttl time.Duration) error {
	key := fmt.Sprintf("permissions:%s", userID)
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, data, ttl).Err()
}

// InvalidateUserPermissions invalidates cached permissions for a user
func (r *permissionCacheRepository) InvalidateUserPermissions(ctx context.Context, userID string) error {
	key := fmt.Sprintf("permissions:%s", userID)
	return r.client.Del(ctx, key).Err()
}
