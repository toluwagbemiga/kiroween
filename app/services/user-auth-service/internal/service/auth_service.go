package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/haunted-saas/user-auth-service/internal/auth"
	"github.com/haunted-saas/user-auth-service/internal/config"
	"github.com/haunted-saas/user-auth-service/internal/domain"
	"github.com/haunted-saas/user-auth-service/internal/errors"
	"github.com/haunted-saas/user-auth-service/internal/logging"
	"github.com/haunted-saas/user-auth-service/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo        repository.UserRepository
	roleRepo        repository.RoleRepository
	sessionRepo     repository.SessionRepository
	rateLimiterRepo repository.RateLimiterRepository
	resetRepo       repository.PasswordResetRepository
	tokenManager    *auth.TokenManager
	config          *config.Config
	logger          *logging.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	sessionRepo repository.SessionRepository,
	rateLimiterRepo repository.RateLimiterRepository,
	resetRepo repository.PasswordResetRepository,
	tokenManager *auth.TokenManager,
	config *config.Config,
	logger *logging.Logger,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		roleRepo:        roleRepo,
		sessionRepo:     sessionRepo,
		rateLimiterRepo: rateLimiterRepo,
		resetRepo:       resetRepo,
		tokenManager:    tokenManager,
		config:          config,
		logger:          logger,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, email, password, name string) (*domain.User, error) {
	// Validate input
	if err := auth.ValidateEmail(email); err != nil {
		return nil, errors.New(errors.ErrCodeInvalidEmail, err.Error())
	}
	
	if err := auth.ValidatePassword(password); err != nil {
		return nil, errors.New(errors.ErrCodeWeakPassword, err.Error())
	}
	
	if err := auth.ValidateName(name); err != nil {
		return nil, errors.New(errors.ErrCodeInvalidInput, err.Error())
	}
	
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New(errors.ErrCodeEmailAlreadyExists, "email already registered")
	}
	
	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), s.config.Security.BcryptCost)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeInternal, "failed to hash password", err)
	}
	
	// Create user
	user := &domain.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		Name:         name,
		IsActive:     true,
		IsLocked:     false,
	}
	
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.Wrap(errors.ErrCodeInternal, "failed to create user", err)
	}
	
	// Assign default "member" role
	memberRole, err := s.roleRepo.FindByName(ctx, "member")
	if err != nil {
		s.logger.Error("failed to find member role", 
			zap.Error(err),
			zap.String("user_id", user.ID))
		// Continue even if role assignment fails
	} else {
		if err := s.userRepo.AssignRole(ctx, user.ID, memberRole.ID); err != nil {
			s.logger.Error("failed to assign member role",
				zap.Error(err),
				zap.String("user_id", user.ID))
		}
	}
	
	// Reload user with roles
	user, err = s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeInternal, "failed to reload user", err)
	}
	
	// Log audit event
	s.logger.LogAuditEvent(&logging.AuditEvent{
		EventType: "user.registered",
		UserID:    user.ID,
		Email:     user.Email,
		Success:   true,
	})
	
	return user, nil
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(ctx context.Context, email, password, ipAddress string) (*domain.User, string, time.Time, error) {
	// Check if account is locked
	locked, duration, err := s.rateLimiterRepo.IsLocked(ctx, email)
	if err != nil {
		s.logger.Error("failed to check lock status", zap.Error(err), zap.String("email", email))
	}
	
	if locked {
		s.logger.LogAuditEvent(&logging.AuditEvent{
			EventType:   "user.login.failed",
			Email:       email,
			IPAddress:   ipAddress,
			Success:     false,
			ErrorReason: "account_locked",
			Metadata: map[string]interface{}{
				"locked_duration_remaining": duration.String(),
			},
		})
		return nil, "", time.Time{}, errors.New(errors.ErrCodeAccountLocked, 
			fmt.Sprintf("account locked for %v", duration.Round(time.Second)))
	}
	
	// Find user
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Record failed attempt even for non-existent users (prevent enumeration)
			s.rateLimiterRepo.RecordFailedAttempt(ctx, email)
			
			s.logger.LogAuditEvent(&logging.AuditEvent{
				EventType:   "user.login.failed",
				Email:       email,
				IPAddress:   ipAddress,
				Success:     false,
				ErrorReason: "invalid_credentials",
			})
			
			return nil, "", time.Time{}, errors.New(errors.ErrCodeInvalidCredentials, "invalid email or password")
		}
		return nil, "", time.Time{}, errors.Wrap(errors.ErrCodeInternal, "failed to find user", err)
	}
	
	// Check if account is locked in database
	if user.IsAccountLocked() {
		s.logger.LogAuditEvent(&logging.AuditEvent{
			EventType:   "user.login.failed",
			UserID:      user.ID,
			Email:       email,
			IPAddress:   ipAddress,
			Success:     false,
			ErrorReason: "account_locked",
		})
		return nil, "", time.Time{}, errors.New(errors.ErrCodeAccountLocked, "account is locked")
	}
	
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		// Record failed attempt
		s.rateLimiterRepo.RecordFailedAttempt(ctx, email)
		
		// Check if we should lock the account
		attempts, _ := s.rateLimiterRepo.GetFailedAttempts(ctx, email)
		if attempts >= s.config.Security.MaxLoginAttempts {
			// Lock account
			lockUntil := time.Now().Add(s.config.Security.LockoutDuration)
			user.IsLocked = true
			user.LockedUntil = &lockUntil
			s.userRepo.Update(ctx, user)
			
			// Also set Redis lock
			s.rateLimiterRepo.LockAccount(ctx, email, s.config.Security.LockoutDuration)
			
			s.logger.LogAuditEvent(&logging.AuditEvent{
				EventType:   "user.account.locked",
				UserID:      user.ID,
				Email:       email,
				IPAddress:   ipAddress,
				Success:     false,
				ErrorReason: "max_login_attempts_exceeded",
				Metadata: map[string]interface{}{
					"attempts": attempts,
					"locked_until": lockUntil,
				},
			})
		}
		
		s.logger.LogAuditEvent(&logging.AuditEvent{
			EventType:   "user.login.failed",
			UserID:      user.ID,
			Email:       email,
			IPAddress:   ipAddress,
			Success:     false,
			ErrorReason: "invalid_password",
			Metadata: map[string]interface{}{
				"attempts": attempts,
			},
		})
		
		return nil, "", time.Time{}, errors.New(errors.ErrCodeInvalidCredentials, "invalid email or password")
	}
	
	// Reset failed attempts on successful login
	s.rateLimiterRepo.ResetAttempts(ctx, email)
	
	// Generate session ID
	sessionID := uuid.New().String()
	
	// Generate JWT
	token, err := s.tokenManager.GenerateToken(user, sessionID)
	if err != nil {
		return nil, "", time.Time{}, errors.Wrap(errors.ErrCodeInternal, "failed to generate token", err)
	}
	
	// Extract JTI from token
	claims, _ := s.tokenManager.ExtractClaims(token)
	
	// Create session
	expiresAt := time.Now().Add(s.config.Security.SessionExpiration)
	session := &domain.Session{
		SessionID:    sessionID,
		UserID:       user.ID,
		TokenJTI:     claims.ID,
		IPAddress:    ipAddress,
		CreatedAt:    time.Now(),
		ExpiresAt:    expiresAt,
		LastActivity: time.Now(),
	}
	
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, "", time.Time{}, errors.Wrap(errors.ErrCodeInternal, "failed to create session", err)
	}
	
	// Log audit event
	s.logger.LogAuditEvent(&logging.AuditEvent{
		EventType: "user.login.success",
		UserID:    user.ID,
		Email:     user.Email,
		IPAddress: ipAddress,
		Success:   true,
		Metadata: map[string]interface{}{
			"session_id": sessionID,
		},
	})
	
	return user, token, expiresAt, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	// Validate token signature and expiration
	claims, err := s.tokenManager.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeInvalidToken, "invalid token", err)
	}
	
	// Check if token is revoked
	revoked, err := s.sessionRepo.IsRevoked(ctx, claims.ID)
	if err != nil {
		s.logger.Error("failed to check token revocation", zap.Error(err), zap.String("jti", claims.ID))
	}
	
	if revoked {
		return nil, errors.New(errors.ErrCodeRevokedToken, "token has been revoked")
	}
	
	// Check if session exists
	session, err := s.sessionRepo.Get(ctx, claims.SessionID)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidToken, "session not found")
	}
	
	// Extend session expiration (sliding window)
	if err := s.sessionRepo.ExtendExpiration(ctx, session.SessionID, s.config.Security.SessionExpiration); err != nil {
		s.logger.Error("failed to extend session", zap.Error(err), zap.String("session_id", session.SessionID))
	}
	
	// Get user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeUserNotFound, "user not found", err)
	}
	
	return user, nil
}

// Logout logs out a user
func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	// Extract claims
	claims, err := s.tokenManager.ExtractClaims(tokenString)
	if err != nil {
		return errors.Wrap(errors.ErrCodeInvalidToken, "invalid token", err)
	}
	
	// Delete session
	if err := s.sessionRepo.Delete(ctx, claims.SessionID); err != nil {
		s.logger.Error("failed to delete session", zap.Error(err), zap.String("session_id", claims.SessionID))
	}
	
	// Revoke token
	if err := s.sessionRepo.RevokeToken(ctx, claims.ID, claims.ExpiresAt.Time); err != nil {
		s.logger.Error("failed to revoke token", zap.Error(err), zap.String("jti", claims.ID))
	}
	
	// Log audit event
	s.logger.LogAuditEvent(&logging.AuditEvent{
		EventType: "user.logout",
		UserID:    claims.UserID,
		Email:     claims.Email,
		Success:   true,
		Metadata: map[string]interface{}{
			"session_id": claims.SessionID,
		},
	})
	
	return nil
}

// LogoutAllDevices logs out a user from all devices
func (s *AuthService) LogoutAllDevices(ctx context.Context, userID string) error {
	// Delete all sessions for user
	if err := s.sessionRepo.DeleteAllForUser(ctx, userID); err != nil {
		return errors.Wrap(errors.ErrCodeInternal, "failed to delete sessions", err)
	}
	
	// Log audit event
	s.logger.LogAuditEvent(&logging.AuditEvent{
		EventType: "user.logout.all_devices",
		UserID:    userID,
		Success:   true,
	})
	
	return nil
}

// RequestPasswordReset generates a password reset token
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists
		s.logger.Info("password reset requested for non-existent email", zap.String("email", email))
		return "", nil
	}
	
	// Generate secure token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", errors.Wrap(errors.ErrCodeInternal, "failed to generate token", err)
	}
	token := hex.EncodeToString(tokenBytes)
	
	// Store token in Redis
	resetToken := &repository.PasswordResetToken{
		UserID:    user.ID,
		Email:     user.Email,
		CreatedAt: time.Now(),
	}
	
	_, err = s.resetRepo.CreateResetToken(ctx, resetToken, s.config.Security.PasswordResetTTL)
	if err != nil {
		return "", errors.Wrap(errors.ErrCodeInternal, "failed to store reset token", err)
	}
	
	// Log audit event
	s.logger.LogAuditEvent(&logging.AuditEvent{
		EventType: "user.password_reset.requested",
		UserID:    user.ID,
		Email:     user.Email,
		Success:   true,
	})
	
	// Return the plain token (to be sent to user via email)
	return token, nil
}

// ResetPassword resets a user's password
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate new password
	if err := auth.ValidatePassword(newPassword); err != nil {
		return errors.New(errors.ErrCodeWeakPassword, err.Error())
	}
	
	// Get reset token
	resetToken, err := s.resetRepo.GetResetToken(ctx, token)
	if err != nil {
		return errors.New(errors.ErrCodeInvalidResetToken, "invalid or expired reset token")
	}
	
	// Get user
	user, err := s.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		return errors.Wrap(errors.ErrCodeUserNotFound, "user not found", err)
	}
	
	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.config.Security.BcryptCost)
	if err != nil {
		return errors.Wrap(errors.ErrCodeInternal, "failed to hash password", err)
	}
	
	// Update password
	user.PasswordHash = string(passwordHash)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.Wrap(errors.ErrCodeInternal, "failed to update password", err)
	}
	
	// Delete reset token
	s.resetRepo.DeleteResetToken(ctx, token)
	
	// Invalidate all sessions
	s.sessionRepo.DeleteAllForUser(ctx, user.ID)
	
	// Log audit event
	s.logger.LogAuditEvent(&logging.AuditEvent{
		EventType: "user.password_reset.completed",
		UserID:    user.ID,
		Email:     user.Email,
		Success:   true,
	})
	
	return nil
}
