package domain

import (
	"time"
)

// Team represents a team/organization in the system
type Team struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string     `gorm:"not null" json:"name"`
	Slug          string     `gorm:"uniqueIndex;not null" json:"slug"`
	OwnerID       *string    `gorm:"type:uuid" json:"owner_id,omitempty"`
	IsMultiTenant bool       `gorm:"default:false" json:"is_multi_tenant"`
	TenantMode    string     `gorm:"default:'single'" json:"tenant_mode"` // 'single' or 'multi'
	MaxMembers    int        `gorm:"default:10" json:"max_members"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	Owner         *User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members       []TeamMember `gorm:"foreignKey:TeamID" json:"members,omitempty"`
}

// TableName specifies the table name for GORM
func (Team) TableName() string {
	return "teams"
}

// TeamMember represents a user's membership in a team
type TeamMember struct {
	ID        string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TeamID    string     `gorm:"type:uuid;not null" json:"team_id"`
	UserID    string     `gorm:"type:uuid;not null" json:"user_id"`
	Role      string     `gorm:"default:'member'" json:"role"` // 'owner', 'admin', 'member'
	Status    string     `gorm:"default:'pending'" json:"status"` // 'pending', 'active', 'suspended'
	InvitedBy *string    `gorm:"type:uuid" json:"invited_by,omitempty"`
	JoinedAt  *time.Time `json:"joined_at,omitempty"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"created_at"`
	Team      *Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for GORM
func (TeamMember) TableName() string {
	return "team_members"
}

// TeamInvitation represents an invitation to join a team
type TeamInvitation struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TeamID    string    `gorm:"type:uuid;not null" json:"team_id"`
	Email     string    `gorm:"not null" json:"email"`
	InvitedBy string    `gorm:"type:uuid;not null" json:"invited_by"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	Status    string    `gorm:"default:'pending'" json:"status"` // 'pending', 'accepted', 'rejected', 'expired'
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	Team      *Team     `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

// TableName specifies the table name for GORM
func (TeamInvitation) TableName() string {
	return "team_invitations"
}

// IsExpired checks if the invitation has expired
func (ti *TeamInvitation) IsExpired() bool {
	return time.Now().After(ti.ExpiresAt)
}
