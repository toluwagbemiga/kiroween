package repository

import (
	"context"

	"github.com/haunted-saas/user-auth-service/internal/domain"
	"gorm.io/gorm"
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	Create(ctx context.Context, role *domain.Role) error
	FindByID(ctx context.Context, id string) (*domain.Role, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]domain.Permission, error)
	AssignPermission(ctx context.Context, roleID, permissionID string) error
	RevokePermission(ctx context.Context, roleID, permissionID string) error
	SetPermissions(ctx context.Context, roleID string, permissionIDs []string) error
}

// roleRepository implements RoleRepository
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// Create creates a new role
func (r *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// FindByID finds a role by ID with permissions preloaded
func (r *roleRepository) FindByID(ctx context.Context, id string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("id = ?", id).
		First(&role).Error
	
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName finds a role by name with permissions preloaded
func (r *roleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("name = ?", name).
		First(&role).Error
	
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Update updates a role
func (r *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// Delete deletes a role
func (r *roleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Role{}, "id = ?", id).Error
}

// GetRolePermissions gets all permissions for a role
func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID string) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	
	return permissions, err
}

// AssignPermission assigns a permission to a role
func (r *roleRepository) AssignPermission(ctx context.Context, roleID, permissionID string) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		roleID, permissionID,
	).Error
}

// RevokePermission revokes a permission from a role
func (r *roleRepository) RevokePermission(ctx context.Context, roleID, permissionID string) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?",
		roleID, permissionID,
	).Error
}

// SetPermissions sets all permissions for a role (replaces existing)
func (r *roleRepository) SetPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing permissions
		if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
			return err
		}
		
		// Insert new permissions
		for _, permID := range permissionIDs {
			if err := tx.Exec(
				"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
				roleID, permID,
			).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}
