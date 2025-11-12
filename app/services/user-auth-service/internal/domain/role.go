package domain

import (
	"time"
)

// Role represents a role that can be assigned to users
type Role struct {
	ID          string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string       `gorm:"uniqueIndex;not null" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
	IsSystem    bool         `gorm:"default:false" json:"is_system"` // System roles cannot be deleted
	CreatedAt   time.Time    `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"not null;default:now()" json:"updated_at"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// TableName specifies the table name for GORM
func (Role) TableName() string {
	return "roles"
}

// CanDelete checks if the role can be deleted (not a system role)
func (r *Role) CanDelete() bool {
	return !r.IsSystem
}
