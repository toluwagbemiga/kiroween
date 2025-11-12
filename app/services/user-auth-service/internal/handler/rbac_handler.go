package handler

import (
	"context"

	"github.com/haunted-saas/user-auth-service/internal/errors"
	pb "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
)

// CreateRole handles role creation
func (h *AuthHandler) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.Role, error) {
	role, err := h.rbacService.CreateRole(ctx, req.Name, req.Description, req.PermissionIds)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return domainRoleToProto(role), nil
}

// UpdateRole handles role updates
func (h *AuthHandler) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.Role, error) {
	role, err := h.rbacService.UpdateRole(ctx, req.RoleId, req.Name, req.Description, req.PermissionIds)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return domainRoleToProto(role), nil
}

// DeleteRole handles role deletion
func (h *AuthHandler) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	if err := h.rbacService.DeleteRole(ctx, req.RoleId); err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return &pb.DeleteRoleResponse{Success: true}, nil
}

// AssignRoleToUser assigns a role to a user
func (h *AuthHandler) AssignRoleToUser(ctx context.Context, req *pb.AssignRoleRequest) (*pb.AssignRoleResponse, error) {
	if err := h.rbacService.AssignRoleToUser(ctx, req.UserId, req.RoleId); err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return &pb.AssignRoleResponse{Success: true}, nil
}

// RevokeRoleFromUser revokes a role from a user
func (h *AuthHandler) RevokeRoleFromUser(ctx context.Context, req *pb.RevokeRoleRequest) (*pb.RevokeRoleResponse, error) {
	if err := h.rbacService.RevokeRoleFromUser(ctx, req.UserId, req.RoleId); err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return &pb.RevokeRoleResponse{Success: true}, nil
}

// CheckPermission checks if a user has a permission
func (h *AuthHandler) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	allowed, err := h.rbacService.CheckPermission(ctx, req.UserId, req.Permission)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	reason := ""
	if !allowed {
		reason = "user does not have the required permission"
	}
	
	return &pb.CheckPermissionResponse{
		Allowed: allowed,
		Reason:  reason,
	}, nil
}

// GetUserPermissions gets all permissions for a user
func (h *AuthHandler) GetUserPermissions(ctx context.Context, req *pb.GetUserPermissionsRequest) (*pb.GetUserPermissionsResponse, error) {
	permissions, err := h.rbacService.GetUserPermissions(ctx, req.UserId)
	if err != nil {
		return nil, errors.MapToGRPCError(err)
	}
	
	return &pb.GetUserPermissionsResponse{
		Permissions: permissions,
	}, nil
}
