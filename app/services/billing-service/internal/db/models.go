package db

import (
	"time"
)

// Plan represents a subscription plan
type Plan struct {
	ID              string            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name            string            `gorm:"not null" json:"name"`
	PriceCents      int64             `gorm:"not null" json:"price_cents"`
	Currency        string            `gorm:"not null;default:usd" json:"currency"`
	BillingInterval string            `gorm:"not null" json:"billing_interval"` // "month" or "year"
	Features        map[string]string `gorm:"type:jsonb;not null;default:'{}'" json:"features"`
	IsActive        bool              `gorm:"not null;default:true" json:"is_active"`
	StripePriceID   string            `gorm:"not null;unique" json:"stripe_price_id"`
	StripeProductID string            `gorm:"not null" json:"stripe_product_id"`
	TrialDays       int32             `gorm:"default:0" json:"trial_days"`
	CreatedByUserID *string           `gorm:"type:uuid" json:"created_by_user_id,omitempty"`
	CreatedAt       time.Time         `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (Plan) TableName() string {
	return "plans"
}

// Subscription represents a team's subscription
type Subscription struct {
	ID                   string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TeamID               string     `gorm:"type:uuid;not null;unique" json:"team_id"` // One subscription per team
	PlanID               string     `gorm:"type:uuid;not null" json:"plan_id"`
	Status               string     `gorm:"not null" json:"status"`
	StripeSubscriptionID string     `gorm:"not null;unique" json:"stripe_subscription_id"`
	StripeCustomerID     string     `gorm:"not null" json:"stripe_customer_id"`
	CurrentPeriodStart   time.Time  `gorm:"not null" json:"current_period_start"`
	CurrentPeriodEnd     time.Time  `gorm:"not null" json:"current_period_end"`
	CancelAt             *time.Time `json:"cancel_at,omitempty"`
	CanceledAt           *time.Time `json:"canceled_at,omitempty"`
	TrialEnd             *time.Time `json:"trial_end,omitempty"`
	CreatedAt            time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	
	Plan Plan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
}

// TableName specifies the table name for GORM
func (Subscription) TableName() string {
	return "subscriptions"
}

// IsActive checks if the subscription is currently active
func (s *Subscription) IsActive() bool {
	return s.Status == "active" || s.Status == "trialing"
}

// IsCanceled checks if the subscription is canceled
func (s *Subscription) IsCanceled() bool {
	return s.Status == "canceled"
}

// WebhookEvent represents a processed webhook event for idempotency
type WebhookEvent struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StripeEventID   string    `gorm:"not null;unique" json:"stripe_event_id"`
	EventType       string    `gorm:"not null" json:"event_type"`
	Processed       bool      `gorm:"not null;default:false" json:"processed"`
	ProcessingError *string   `json:"processing_error,omitempty"`
	ReceivedAt      time.Time `gorm:"not null;default:now()" json:"received_at"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
}

// TableName specifies the table name for GORM
func (WebhookEvent) TableName() string {
	return "webhook_events"
}
