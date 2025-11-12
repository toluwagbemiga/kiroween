package service

import (
	"context"
	"testing"

	"github.com/haunted-saas/user-auth-service/internal/config"
	"github.com/haunted-saas/user-auth-service/internal/domain"
	"github.com/haunted-saas/user-auth-service/internal/errors"
	"github.com/haunted-saas/user-auth-service/internal/logging"
	"github.com/haunted-saas/user-auth-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) FindByID(ctx context.Context, id string) (*domain.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permission), args.Error(1)
}

func (m *MockPermissionRepository) FindByName(ctx context.Context, name string) (*domain.Permission, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permission), args.Error(1)
}

func (m *MockPermissionRepository) FindByIDs(ctx context.Context, ids []string) ([]domain.Permission, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]domain.Permission), args.Error(1)
}

func (m *MockPermissionRepository) List(ctx context.Context) ([]domain.Permission, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Permission), args.Error(1)
}

type MockPermissionCacheRepository struct {
	mock.Mock
}

func (m *MockPermissionCacheRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPermissionCacheRepository) SetUserPermissions(ctx context.Context, userID string, permissions []string, ttl interface{}) error {
	args := m.Called(ctx, userID, permissions, ttl)
	return args.Error(0)
}

func (m *MockPermissionCacheRepository) InvalidateUserPermissions(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Test CheckPermission
func TestRBACService_CheckPermission(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		permission     string
		setupMocks     func(*MockUserRepository, *MockPermissionCacheRepository)
		expectedResult bool
		expectedError  error
	}{
		{
			name:       "permission found in cache",
			userID:     "user-123",
			permission: "users:read",
			setupMocks: func(userRepo *MockUserRepository, cacheRepo *MockPermissionCacheRepository) {
				cacheRepo.On("GetUserPermissions", mock.Anything, "user-123").Return([]string{"users:read", "users:write"}, nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:       "permission not in cache - found in database",
			userID:     "user-123",
			permission: "users:read",
			setupMocks: func(userRepo *MockUserRepository, cacheRepo *MockPermissionCacheRepository) {
				cacheRepo.On("GetUserPermissions", mock.Anything, "user-123").Return(nil, repository.ErrNotFound)
				userRepo.On("FindByID", mock.Anything, "user-123").Return(&domain.User{
					ID:    "user-123",
					Email: "test@example.com",
					Roles: []domain.Role{
						{
							Name: "member",
							Permissions: []domain.Permission{
								{Name: "users:read"},
								{Name: "users:write"},
							},
						},
					},
				}, nil)
				cacheRepo.On("SetUserPermissions", mock.Anything, "user-123", mock.Anything, mock.Anything).Return(nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:       "permission not found",
			userID:     "user-123",
			permission: "admin:delete",
			setupMocks: func(userRepo *MockUserRepository, cacheRepo *MockPermissionCacheRepository) {
				cacheRepo.On("GetUserPermissions", mock.Anything, "user-123").Return([]string{"users:read"}, nil)
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:       "user not found",
			userID:     "nonexistent",
			permission: "users:read",
			setupMocks: func(userRepo *MockUserRepository, cacheRepo *MockPermissionCacheRepository) {
				cacheRepo.On("GetUserPermissions", mock.Anything, "nonexistent").Return(nil, repository.ErrNotFound)
				userRepo.On("FindByID", mock.Anything, "nonexistent").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedResult: false,
			expectedError:  errors.New(errors.ErrCodeUserNotFound, ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			cacheRepo := new(MockPermissionCacheRepository)
			sessionRepo := new(MockSessionRepository)

			tt.setupMocks(userRepo, cacheRepo)

			// Create service
			logger, _ := logging.NewLogger("error")
			cfg := &config.Config{}

			service := NewRBACService(
				userRepo,
				nil,
				nil,
				cacheRepo,
				sessionRepo,
				cfg,
				logger,
			)

			// Execute
			result, err := service.CheckPermission(context.Background(), tt.userID, tt.permission)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				serviceErr, ok := err.(*errors.ServiceError)
				if ok {
					expectedErr := tt.expectedError.(*errors.ServiceError)
					assert.Equal(t, expectedErr.Code, serviceErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			userRepo.AssertExpectations(t)
			cacheRepo.AssertExpectations(t)
		})
	}
}

// Test AssignRoleToUser
func TestRBACService_AssignRoleToUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		roleID        string
		setupMocks    func(*MockUserRepository, *MockRoleRepository, *MockPermissionCacheRepository, *MockSessionRepository)
		expectedError error
	}{
		{
			name:   "successful role assignment",
			userID: "user-123",
			roleID: "role-456",
			setupMocks: func(userRepo *MockUserRepository, roleRepo *MockRoleRepository, cacheRepo *MockPermissionCacheRepository, sessionRepo *MockSessionRepository) {
				userRepo.On("FindByID", mock.Anything, "user-123").Return(&domain.User{
					ID:    "user-123",
					Email: "test@example.com",
				}, nil)
				roleRepo.On("FindByID", mock.Anything, "role-456").Return(&domain.Role{
					ID:   "role-456",
					Name: "admin",
				}, nil)
				userRepo.On("AssignRole", mock.Anything, "user-123", "role-456").Return(nil)
				cacheRepo.On("InvalidateUserPermissions", mock.Anything, "user-123").Return(nil)
				sessionRepo.On("DeleteAllForUser", mock.Anything, "user-123").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			roleID: "role-456",
			setupMocks: func(userRepo *MockUserRepository, roleRepo *MockRoleRepository, cacheRepo *MockPermissionCacheRepository, sessionRepo *MockSessionRepository) {
				userRepo.On("FindByID", mock.Anything, "nonexistent").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: errors.New(errors.ErrCodeUserNotFound, ""),
		},
		{
			name:   "role not found",
			userID: "user-123",
			roleID: "nonexistent",
			setupMocks: func(userRepo *MockUserRepository, roleRepo *MockRoleRepository, cacheRepo *MockPermissionCacheRepository, sessionRepo *MockSessionRepository) {
				userRepo.On("FindByID", mock.Anything, "user-123").Return(&domain.User{
					ID:    "user-123",
					Email: "test@example.com",
				}, nil)
				roleRepo.On("FindByID", mock.Anything, "nonexistent").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: errors.New(errors.ErrCodeRoleNotFound, ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			roleRepo := new(MockRoleRepository)
			cacheRepo := new(MockPermissionCacheRepository)
			sessionRepo := new(MockSessionRepository)

			tt.setupMocks(userRepo, roleRepo, cacheRepo, sessionRepo)

			// Create service
			logger, _ := logging.NewLogger("error")
			cfg := &config.Config{}

			service := NewRBACService(
				userRepo,
				roleRepo,
				nil,
				cacheRepo,
				sessionRepo,
				cfg,
				logger,
			)

			// Execute
			err := service.AssignRoleToUser(context.Background(), tt.userID, tt.roleID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				serviceErr, ok := err.(*errors.ServiceError)
				if ok {
					expectedErr := tt.expectedError.(*errors.ServiceError)
					assert.Equal(t, expectedErr.Code, serviceErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
			roleRepo.AssertExpectations(t)
			cacheRepo.AssertExpectations(t)
			sessionRepo.AssertExpectations(t)
		})
	}
}

// Define ErrNotFound for tests
var ErrNotFound = repository.ErrNotFound
