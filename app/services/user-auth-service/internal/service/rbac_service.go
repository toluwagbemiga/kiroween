package service

import (
	"context"

	"github.com/haunted-saas/user-auth-service/internal/config"
	"github.com/haunted-saas/user-auth-service/internal/domain"
	"github.com/haunted-saas/user-auth-service/internal/errors"
	"github.com/haunted-saas/user-auth-service/internal/logging"
	"github.com/haunted-saas/user-auth-service/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RBACService handles role-based access control operations
type RBACService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permRepo       repository.PermissionRepository
	permCacheRepo  repository.PermissionCacheRepository
	sessionRepo    repository.SessionRepository
	config         *config.Config
	logger         *logging.Logger
}

// NewRBACService creates a new RBAC service
func NewRBACService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
	permCacheRepo repository.PermissionCacheRepository,
	sessionRepo repository.SessionRepository,
	config *config.Config,
	logger *logging.Logger,
) *RBACService {
	return &RBACService{
		userRepo:      userRepo,
		roleRepo:      roleRepo,
		permRepo:      permRepo,
		permCacheRepo: permCacheRepo,
		sessionRepo:   sessionRepo,
		config:        config,
		logger:        logger,
	}
}

// CreateRole creates a new role
func (s *RBACService) CreateRole(ctx context.Context, name, description string, permissionIDs []string) (*domain.Role, error) {
	// Check if role already exists
	existingRole, err := s.roleRepo.FindByName(ctx, name)
	if err == nil && existingRole != nil {
		return nil, errors.New(errors.ErrCodeInvalidInput, "role already exists")
	}
	
	// Create role
	role := &domain.Role{
		Name:        name,
		Description: description,
		IsSystem:    false,
	}
	
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, errors.Wrap(errors.ErrCodeInternal, "failed to create role", err)
	}
	
	// Assign permissions
	if len(permissionIDs) > 0 {
		if err := s.roleRepo.SetPermissions(ctx, role.ID, permissionIDs); err != nil {
			return nil, errors.Wrap(errors.ErrCodeInternal, "failed to assign permissions", err)
		}
	}
	
	// Reload role with permissions
	role, err = s.roleRepo.FindByID(ctx, role.ID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeInternal, "failed to reload role", err)
	}
	
	s.logger.Info("role created",
		zap.String("role_id", role.ID),
		zap.String("role_name", role.Name),
		zap.Int("permission_count", len(role.Permissions)))
	
	return role, nil
}

// UpdateRole updates a role
func (s *RBACService) UpdateRole(ctx context.Context, roleID, name, description string, permissionIDs []string) (*domain.Role, error) {
	// Get role
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrCodeRoleNotFound, "role not found")
		}
		return nil, errors.Wrap(errors.ErrCodeInternal, "failed to find role", err)
	}
	
	// Check if system role
	if role.IsSystem {
		return nil, errors.New(errors.ErrCodeSystemRoleProtected, "cannot modify system role")
	}
	
	// Update fields
	if name != "" {
		role.Name = name
	}
	if description != "" {
		role.Description = description
	}
	
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, errors.Wrap(errors.ErrCodeInternal, "failed to update role", err)
	}
	
	// Update permissions if provided
	if permissionIDs != nil {
		if err := s.roleRepo.SetPermissions(ctx, role.ID, permissionIDs); err != nil {
			return nil, errors.Wrap(errors.ErrCodeInternal, "failed to update permissions", err)
		}
	}
	
	// Reload role with permissions
	role, err = s.roleRepo.FindByID(ctx, role.ID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeInternal, "failed to reload role", err)
	}
	
	s.logger.Info("role updated",
		zap.String("role_id", role.ID),
		zap.String("role_name", role.Name))
	
	return role, nil
}

// DeleteRole deletes a role
func (s *RBACService) DeleteRole(ctx context.Context, roleID string) error {
	// Get role
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ErrCodeRoleNotFound, "role not found")
		}
		return errors.Wrap(errors.ErrCodeInternal, "failed to find role", err)
	}
	
	// Check if system role
	if role.IsSystem {
		return errors.New(errors.ErrCodeSystemRoleProtected, "cannot delete system role")
	}
	
	// Delete role
	if err := s.roleRepo.Delete(ctx, roleID); err != nil {
		return errors.Wrap(errors.ErrCodeInternal, "failed to delete role", err)
	}
	
	s.logger.Info("role deleted",
		zap.String("role_id", roleID),
		zap.String("role_name", role.Name))
	
	return nil
}

// AssignRoleToUser assigns a role to a user
func (s *RBACService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ErrCodeUserNotFound, "user not found")
		}
		return errors.Wrap(errors.ErrCodeInternal, "failed to find user", err)
	}
	
	// Verify role exists
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ErrCodeRoleNotFound, "role not found")
		}
		return errors.Wrap(errors.ErrCodeInternal, "failed to find role", err)
	}
	
	// Assign role
	if err := s.userRepo.AssignRole(ctx, userID, roleID); err != nil {
		return errors.Wrap(errors.ErrCodeInternal, "failed to assign role", err)
	}
	
	// Invalidate permission cache
	s.permCacheRepo.InvalidateUserPermissions(ctx, userID)
	
	// Invalidate all user sessions (force re-login to get new permissions)
	s.sessionRepo.DeleteAllForUser(ctx, userID)
	
	// Log audit event
	s.logger.LogAuditEvent(&logging.AuditEvent{
		EventType: "user.role.assigned",
		UserID:    userID,
		Email:     user.Email,
		Success:   true,
		Metadata: map[string]interface{}{
			"role_id":   roleID,
			"role_name": role.Name,
		},
	})
	
	return nil
}

// RevokeRoleFromUser revokes a role from a user
func (s *RBACService) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ErrCodeUserNotFound, "user not found")
		}
		return errors.Wrap(errors.ErrCodeInternal, "failed to find user", err)
	}
	
	// Verify role exists
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ErrCodeRoleNotFound, "role not found")
		}
		return errors.Wrap(errors.ErrCodeInternal, "failed to find role", err)
	}
	
	// Revoke role
	if err := s.userRepo.RevokeRole(ctx, userID, roleID); err != nil {
		return errors.Wrap(errors.ErrCodeInternal, "failed to revoke role", err)
	}
	
	// Invalidate permission cache
	s.permCacheRepo.InvalidateUserPermissions(ctx, userID)
	
	// Invalidate all user sessions
	s.sessionRepo.DeleteAllForUser(ctx, userID)
	
	// Log audit event
	s.logger.LogAuditEvent(&logging.AuditEvent{
		EventType: "user.role.revoked",
		UserID:    userID,
		Email:     user.Email,
		Success:   true,
		Metadata: map[string]interface{}{
			"role_id":   roleID,
			"role_name": role.Name,
		},
	})
	
	return nil
}

// CheckPermission checks if a user has a specific permission
func (s *RBACService) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	// Try cache first
	cachedPerms, err := s.permCacheRepo.GetUserPermissions(ctx, userID)
	if err == nil {
		// Check if permission is in cache
		for _, perm := range cachedPerms {
			if perm == permission {
				return true, nil
			}
		}
		return false, nil
	}
	
	// Cache miss - get from database
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, errors.New(errors.ErrCodeUserNotFound, "user not found")
		}
		return false, errors.Wrap(errors.ErrCodeInternal, "failed to find user", err)
	}
	
	// Get all permissions from all roles
	permissions := user.GetPermissions()
	
	// Cache permissions
	s.permCacheRepo.SetUserPermissions(ctx, userID, permissions, s.config.Security.PermissionCacheTTL)
	
	// Check if user has the permission
	for _, perm := range permissions {
		if perm == permission {
			return true, nil
		}
	}
	
	return false, nil
}

// GetUserPermissions gets all permissions for a user
func (s *RBACService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	// Try cache first
	cachedPerms, err := s.permCacheRepo.GetUserPermissions(ctx, userID)
	if err == nil {
		return cachedPerms, nil
	}
	
	// Cache miss - get from database
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrCodeUserNotFound, "user not found")
		}
		return nil, errors.Wrap(errors.ErrCodeInternal, "failed to find user", err)
	}
	
	// Get all permissions from all roles
	permissions := user.GetPermissions()
	
	// Cache permissions
	s.permCacheRepo.SetUserPermissions(ctx, userID, permissions, s.config.Security.PermissionCacheTTL)
	
	return permissions, nil
}
