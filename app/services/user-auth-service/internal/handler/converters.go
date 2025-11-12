package handler

import (
	"github.com/haunted-saas/user-auth-service/internal/domain"
	pb "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// domainUserToProto converts a domain user to proto
func domainUserToProto(user *domain.User) *pb.User {
	if user == nil {
		return nil
	}
	
	pbUser := &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		IsActive:  user.IsActive,
		IsLocked:  user.IsLocked,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
	
	// Convert roles
	if len(user.Roles) > 0 {
		pbUser.Roles = make([]*pb.Role, len(user.Roles))
		for i, role := range user.Roles {
			pbUser.Roles[i] = domainRoleToProto(&role)
		}
	}
	
	return pbUser
}

// domainRoleToProto converts a domain role to proto
func domainRoleToProto(role *domain.Role) *pb.Role {
	if role == nil {
		return nil
	}
	
	pbRole := &pb.Role{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   timestamppb.New(role.CreatedAt),
		UpdatedAt:   timestamppb.New(role.UpdatedAt),
	}
	
	// Convert permissions
	if len(role.Permissions) > 0 {
		pbRole.Permissions = make([]*pb.Permission, len(role.Permissions))
		for i, perm := range role.Permissions {
			pbRole.Permissions[i] = domainPermissionToProto(&perm)
		}
	}
	
	return pbRole
}

// domainPermissionToProto converts a domain permission to proto
func domainPermissionToProto(perm *domain.Permission) *pb.Permission {
	if perm == nil {
		return nil
	}
	
	return &pb.Permission{
		Id:          perm.ID,
		Name:        perm.Name,
		Resource:    perm.Resource,
		Action:      perm.Action,
		Description: perm.Description,
		CreatedAt:   timestamppb.New(perm.CreatedAt),
	}
}
