package domain

import (
	"time"
)

// User represents a user account in the system
type User struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"not null" json:"-"` // Never serialize password hash
	Name         string     `gorm:"not null" json:"name"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	IsLocked     bool       `gorm:"default:false" json:"is_locked"`
	LockedUntil  *time.Time `gorm:"index" json:"locked_until,omitempty"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	Roles        []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// IsAccountLocked checks if the account is currently locked
func (u *User) IsAccountLocked() bool {
	if !u.IsLocked {
		return false
	}
	if u.LockedUntil == nil {
		return true
	}
	return time.Now().Before(*u.LockedUntil)
}

// GetPermissions aggregates all permissions from all roles
func (u *User) GetPermissions() []string {
	permissionSet := make(map[string]bool)
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			permissionSet[perm.Name] = true
		}
	}
	
	permissions := make([]string, 0, len(permissionSet))
	for perm := range permissionSet {
		permissions = append(permissions, perm)
	}
	return permissions
}

// GetRoleNames returns a list of role names
func (u *User) GetRoleNames() []string {
	roleNames := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		roleNames[i] = role.Name
	}
	return roleNames
}
