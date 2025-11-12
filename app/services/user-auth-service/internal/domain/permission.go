package domain

import (
	"fmt"
	"time"
)

// Permission represents a granular capability in the system
type Permission struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"` // e.g., "users:read", "billing:write"
	Resource    string    `gorm:"not null" json:"resource"`         // e.g., "users", "billing"
	Action      string    `gorm:"not null" json:"action"`           // e.g., "read", "write", "delete"
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// TableName specifies the table name for GORM
func (Permission) TableName() string {
	return "permissions"
}

// NewPermission creates a new permission with resource:action format
func NewPermission(resource, action, description string) *Permission {
	return &Permission{
		Name:        fmt.Sprintf("%s:%s", resource, action),
		Resource:    resource,
		Action:      action,
		Description: description,
	}
}
