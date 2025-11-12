package service

import (
	"context"
	"testing"
	"time"

	"github.com/haunted-saas/user-auth-service/internal/auth"
	"github.com/haunted-saas/user-auth-service/internal/config"
	"github.com/haunted-saas/user-auth-service/internal/domain"
	"github.com/haunted-saas/user-auth-service/internal/errors"
	"github.com/haunted-saas/user-auth-service/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Mock repositories
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Role), args.Error(1)
}

func (m *MockUserRepository) AssignRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockUserRepository) RevokeRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	return nil
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id string) (*domain.Role, error) {
	return nil, nil
}

func (m *MockRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	return nil
}

func (m *MockRoleRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockRoleRepository) GetRolePermissions(ctx context.Context, roleID string) ([]domain.Permission, error) {
	return nil, nil
}

func (m *MockRoleRepository) AssignPermission(ctx context.Context, roleID, permissionID string) error {
	return nil
}

func (m *MockRoleRepository) RevokePermission(ctx context.Context, roleID, permissionID string) error {
	return nil
}

func (m *MockRoleRepository) SetPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return nil
}

type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) Get(ctx context.Context, sessionID string) (*domain.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) Delete(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteAllForUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockSessionRepository) ExtendExpiration(ctx context.Context, sessionID string, duration time.Duration) error {
	args := m.Called(ctx, sessionID, duration)
	return args.Error(0)
}

func (m *MockSessionRepository) IsRevoked(ctx context.Context, tokenJTI string) (bool, error) {
	args := m.Called(ctx, tokenJTI)
	return args.Bool(0), args.Error(1)
}

func (m *MockSessionRepository) RevokeToken(ctx context.Context, tokenJTI string, expiresAt time.Time) error {
	args := m.Called(ctx, tokenJTI, expiresAt)
	return args.Error(0)
}

type MockRateLimiterRepository struct {
	mock.Mock
}

func (m *MockRateLimiterRepository) RecordFailedAttempt(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockRateLimiterRepository) GetFailedAttempts(ctx context.Context, email string) (int, error) {
	args := m.Called(ctx, email)
	return args.Int(0), args.Error(1)
}

func (m *MockRateLimiterRepository) ResetAttempts(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockRateLimiterRepository) IsLocked(ctx context.Context, email string) (bool, time.Duration, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Get(1).(time.Duration), args.Error(2)
}

func (m *MockRateLimiterRepository) LockAccount(ctx context.Context, email string, duration time.Duration) error {
	args := m.Called(ctx, email, duration)
	return args.Error(0)
}

// Test Register
func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		userName      string
		setupMocks    func(*MockUserRepository, *MockRoleRepository)
		expectedError error
	}{
		{
			name:     "successful registration",
			email:    "test@example.com",
			password: "ValidPass123!",
			userName: "Test User",
			setupMocks: func(userRepo *MockUserRepository, roleRepo *MockRoleRepository) {
				userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, gorm.ErrRecordNotFound)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
				userRepo.On("FindByID", mock.Anything, mock.Anything).Return(&domain.User{
					ID:    "user-123",
					Email: "test@example.com",
					Name:  "Test User",
					Roles: []domain.Role{{Name: "member"}},
				}, nil)
				roleRepo.On("FindByName", mock.Anything, "member").Return(&domain.Role{
					ID:   "role-123",
					Name: "member",
				}, nil)
				userRepo.On("AssignRole", mock.Anything, mock.Anything, "role-123").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "invalid email",
			email:    "invalid-email",
			password: "ValidPass123!",
			userName: "Test User",
			setupMocks: func(userRepo *MockUserRepository, roleRepo *MockRoleRepository) {
				// No mocks needed - validation fails before DB access
			},
			expectedError: errors.New(errors.ErrCodeInvalidEmail, ""),
		},
		{
			name:     "weak password",
			email:    "test@example.com",
			password: "weak",
			userName: "Test User",
			setupMocks: func(userRepo *MockUserRepository, roleRepo *MockRoleRepository) {
				// No mocks needed - validation fails before DB access
			},
			expectedError: errors.New(errors.ErrCodeWeakPassword, ""),
		},
		{
			name:     "email already exists",
			email:    "existing@example.com",
			password: "ValidPass123!",
			userName: "Test User",
			setupMocks: func(userRepo *MockUserRepository, roleRepo *MockRoleRepository) {
				userRepo.On("FindByEmail", mock.Anything, "existing@example.com").Return(&domain.User{
					Email: "existing@example.com",
				}, nil)
			},
			expectedError: errors.New(errors.ErrCodeEmailAlreadyExists, ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			roleRepo := new(MockRoleRepository)
			sessionRepo := new(MockSessionRepository)
			rateLimiterRepo := new(MockRateLimiterRepository)
			
			tt.setupMocks(userRepo, roleRepo)

			// Create service
			logger, _ := logging.NewLogger("error")
			cfg := &config.Config{
				Security: config.SecurityConfig{
					BcryptCost: 12,
				},
			}
			
			service := NewAuthService(
				userRepo,
				roleRepo,
				sessionRepo,
				rateLimiterRepo,
				nil,
				nil,
				cfg,
				logger,
			)

			// Execute
			user, err := service.Register(context.Background(), tt.email, tt.password, tt.userName)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
				serviceErr, ok := err.(*errors.ServiceError)
				if ok {
					expectedErr := tt.expectedError.(*errors.ServiceError)
					assert.Equal(t, expectedErr.Code, serviceErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}

			userRepo.AssertExpectations(t)
			roleRepo.AssertExpectations(t)
		})
	}
}

// Test Login
func TestAuthService_Login(t *testing.T) {
	// Create a valid password hash
	validPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("ValidPass123!"), 12)

	tests := []struct {
		name          string
		email         string
		password      string
		ipAddress     string
		setupMocks    func(*MockUserRepository, *MockRateLimiterRepository, *MockSessionRepository)
		expectedError error
	}{
		{
			name:      "successful login",
			email:     "test@example.com",
			password:  "ValidPass123!",
			ipAddress: "192.168.1.1",
			setupMocks: func(userRepo *MockUserRepository, rateLimiter *MockRateLimiterRepository, sessionRepo *MockSessionRepository) {
				rateLimiter.On("IsLocked", mock.Anything, "test@example.com").Return(false, time.Duration(0), nil)
				userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(&domain.User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: string(validPasswordHash),
					IsActive:     true,
					IsLocked:     false,
					Roles: []domain.Role{
						{Name: "member", Permissions: []domain.Permission{{Name: "users:read"}}},
					},
				}, nil)
				rateLimiter.On("ResetAttempts", mock.Anything, "test@example.com").Return(nil)
				sessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Session")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "invalid password",
			email:     "test@example.com",
			password:  "WrongPassword123!",
			ipAddress: "192.168.1.1",
			setupMocks: func(userRepo *MockUserRepository, rateLimiter *MockRateLimiterRepository, sessionRepo *MockSessionRepository) {
				rateLimiter.On("IsLocked", mock.Anything, "test@example.com").Return(false, time.Duration(0), nil)
				userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(&domain.User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: string(validPasswordHash),
					IsActive:     true,
					IsLocked:     false,
				}, nil)
				rateLimiter.On("RecordFailedAttempt", mock.Anything, "test@example.com").Return(nil)
				rateLimiter.On("GetFailedAttempts", mock.Anything, "test@example.com").Return(1, nil)
			},
			expectedError: errors.New(errors.ErrCodeInvalidCredentials, ""),
		},
		{
			name:      "account locked",
			email:     "locked@example.com",
			password:  "ValidPass123!",
			ipAddress: "192.168.1.1",
			setupMocks: func(userRepo *MockUserRepository, rateLimiter *MockRateLimiterRepository, sessionRepo *MockSessionRepository) {
				rateLimiter.On("IsLocked", mock.Anything, "locked@example.com").Return(true, 30*time.Minute, nil)
			},
			expectedError: errors.New(errors.ErrCodeAccountLocked, ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			rateLimiterRepo := new(MockRateLimiterRepository)
			sessionRepo := new(MockSessionRepository)

			tt.setupMocks(userRepo, rateLimiterRepo, sessionRepo)

			// Create service (without token manager for this test)
			logger, _ := logging.NewLogger("error")
			cfg := &config.Config{
				Security: config.SecurityConfig{
					BcryptCost:          12,
					MaxLoginAttempts:    5,
					LockoutDuration:     30 * time.Minute,
					SessionExpiration:   24 * time.Hour,
				},
			}

			// Create a mock token manager
			tokenManager, _ := auth.NewTokenManager(
				"../../keys/jwt-private.pem",
				"../../keys/jwt-public.pem",
				24*time.Hour,
			)

			service := NewAuthService(
				userRepo,
				nil,
				sessionRepo,
				rateLimiterRepo,
				nil,
				tokenManager,
				cfg,
				logger,
			)

			// Execute
			user, token, expiresAt, err := service.Login(context.Background(), tt.email, tt.password, tt.ipAddress)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Empty(t, token)
				serviceErr, ok := err.(*errors.ServiceError)
				if ok {
					expectedErr := tt.expectedError.(*errors.ServiceError)
					assert.Equal(t, expectedErr.Code, serviceErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotEmpty(t, token)
				assert.False(t, expiresAt.IsZero())
			}

			userRepo.AssertExpectations(t)
			rateLimiterRepo.AssertExpectations(t)
			sessionRepo.AssertExpectations(t)
		})
	}
}
