package repository

import (
	"context"

	"github.com/haunted-saas/user-auth-service/internal/domain"
	"gorm.io/gorm"
)

// PermissionRepository defines the interface for permission data access
type PermissionRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Permission, error)
	FindByName(ctx context.Context, name string) (*domain.Permission, error)
	FindByIDs(ctx context.Context, ids []string) ([]domain.Permission, error)
	List(ctx context.Context) ([]domain.Permission, error)
}

// permissionRepository implements PermissionRepository
type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// FindByID finds a permission by ID
func (r *permissionRepository) FindByID(ctx context.Context, id string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByName finds a permission by name
func (r *permissionRepository) FindByName(ctx context.Context, name string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByIDs finds permissions by IDs
func (r *permissionRepository) FindByIDs(ctx context.Context, ids []string) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&permissions).Error
	return permissions, err
}

// List lists all permissions
func (r *permissionRepository) List(ctx context.Context) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	return permissions, err
}
