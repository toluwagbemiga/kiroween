package handler

import (
	"context"

	"github.com/haunted-saas/user-auth-service/internal/errors"
	"github.com/haunted-saas/user-auth-service/internal/service"
	pb "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthHandler handles authentication gRPC requests
type AuthHandler struct {
	pb.UnimplementedUserAuthServiceServer
	authService *service.AuthService
	rbacService *service.RBACService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, rbacService *service.RBACService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		rbacService: rbacService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := h.authService.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return &pb.RegisterResponse{
		User: domainUserToProto(user),
	}, nil
}

// Login handles user login
func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, token, expiresAt, err := h.authService.Login(ctx, req.Email, req.Password, req.IpAddress)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return &pb.LoginResponse{
		AccessToken: token,
		User:        domainUserToProto(user),
		ExpiresAt:   timestamppb.New(expiresAt),
	}, nil
}

// Logout handles user logout
func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.AllDevices {
		// Extract user ID from token
		// For now, we'll use the session token to logout
		if err := h.authService.Logout(ctx, req.SessionToken); err != nil {
			return nil, errors.MapToGRPCError(err)
		}
	} else {
		if err := h.authService.Logout(ctx, req.SessionToken); err != nil {
			return nil, errors.MapToGRPCError(err)
		}
	}
	
	return &pb.LogoutResponse{Success: true}, nil
}

// ValidateToken validates a JWT token
func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	user, err := h.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	permissions, err := h.rbacService.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return &pb.ValidateTokenResponse{
		Valid:       true,
		UserId:      user.ID,
		TeamId:      "", // Add team_id when teams are implemented
		Roles:       user.GetRoleNames(),
		Permissions: permissions,
		User:        domainUserToProto(user),
	}, nil
}

// RefreshSession refreshes a session
func (h *AuthHandler) RefreshSession(ctx context.Context, req *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	// Validate the refresh token
	_, err := h.authService.ValidateToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	// For now, we'll just validate and return the same token
	// In a full implementation, you'd generate a new access token
	return &pb.RefreshSessionResponse{
		AccessToken: req.RefreshToken,
		ExpiresAt:   timestamppb.Now(),
	}, nil
}

// RequestPasswordReset handles password reset requests
func (h *AuthHandler) RequestPasswordReset(ctx context.Context, req *pb.PasswordResetRequest) (*pb.PasswordResetResponse, error) {
	token, err := h.authService.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	// In production, you would send this token via email
	// For now, we return success without exposing the token
	_ = token
	
	return &pb.PasswordResetResponse{
		Success: true,
		Message: "If the email exists, a password reset link has been sent",
	}, nil
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if err := h.authService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return &pb.ResetPasswordResponse{Success: true}, nil
}
