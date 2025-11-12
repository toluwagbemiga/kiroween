package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Store handles all database operations
type Store struct {
	db *gorm.DB
}

// NewStore creates a new database store
func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// Plan Operations

// CreatePlan creates a new plan in the database
func (s *Store) CreatePlan(ctx context.Context, plan *Plan) error {
	return s.db.WithContext(ctx).Create(plan).Error
}

// GetPlanByID retrieves a plan by ID
func (s *Store) GetPlanByID(ctx context.Context, planID string) (*Plan, error) {
	var plan Plan
	err := s.db.WithContext(ctx).Where("id = ?", planID).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// GetPlanByStripePriceID retrieves a plan by Stripe price ID
func (s *Store) GetPlanByStripePriceID(ctx context.Context, stripePriceID string) (*Plan, error) {
	var plan Plan
	err := s.db.WithContext(ctx).Where("stripe_price_id = ?", stripePriceID).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// ListPlans retrieves all plans, optionally filtering by active status
func (s *Store) ListPlans(ctx context.Context, activeOnly bool) ([]Plan, error) {
	var plans []Plan
	query := s.db.WithContext(ctx)
	
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	
	err := query.Order("price_cents ASC").Find(&plans).Error
	return plans, err
}

// UpdatePlan updates a plan
func (s *Store) UpdatePlan(ctx context.Context, plan *Plan) error {
	return s.db.WithContext(ctx).Save(plan).Error
}

// DeactivatePlan marks a plan as inactive
func (s *Store) DeactivatePlan(ctx context.Context, planID string) error {
	return s.db.WithContext(ctx).Model(&Plan{}).
		Where("id = ?", planID).
		Update("is_active", false).Error
}

// Subscription Operations

// CreateSubscription creates a new subscription
func (s *Store) CreateSubscription(ctx context.Context, subscription *Subscription) error {
	return s.db.WithContext(ctx).Create(subscription).Error
}

// GetSubscriptionByTeamID retrieves a subscription by team ID
func (s *Store) GetSubscriptionByTeamID(ctx context.Context, teamID string) (*Subscription, error) {
	var subscription Subscription
	err := s.db.WithContext(ctx).
		Preload("Plan").
		Where("team_id = ?", teamID).
		First(&subscription).Error
	
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// GetSubscriptionByStripeID retrieves a subscription by Stripe subscription ID
func (s *Store) GetSubscriptionByStripeID(ctx context.Context, stripeSubID string) (*Subscription, error) {
	var subscription Subscription
	err := s.db.WithContext(ctx).
		Preload("Plan").
		Where("stripe_subscription_id = ?", stripeSubID).
		First(&subscription).Error
	
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// UpdateSubscription updates a subscription
func (s *Store) UpdateSubscription(ctx context.Context, subscription *Subscription) error {
	return s.db.WithContext(ctx).Save(subscription).Error
}

// UpdateSubscriptionStatus updates only the status of a subscription
func (s *Store) UpdateSubscriptionStatus(ctx context.Context, subscriptionID, status string) error {
	return s.db.WithContext(ctx).Model(&Subscription{}).
		Where("id = ?", subscriptionID).
		Update("status", status).Error
}

// DeleteSubscription deletes a subscription (soft delete)
func (s *Store) DeleteSubscription(ctx context.Context, subscriptionID string) error {
	return s.db.WithContext(ctx).Delete(&Subscription{}, "id = ?", subscriptionID).Error
}

// Webhook Event Operations (for idempotency)

// CreateWebhookEvent creates a webhook event record
func (s *Store) CreateWebhookEvent(ctx context.Context, event *WebhookEvent) error {
	return s.db.WithContext(ctx).Create(event).Error
}

// GetWebhookEventByStripeID retrieves a webhook event by Stripe event ID
func (s *Store) GetWebhookEventByStripeID(ctx context.Context, stripeEventID string) (*WebhookEvent, error) {
	var event WebhookEvent
	err := s.db.WithContext(ctx).Where("stripe_event_id = ?", stripeEventID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// MarkWebhookEventProcessed marks a webhook event as processed
func (s *Store) MarkWebhookEventProcessed(ctx context.Context, stripeEventID string, processingError *string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"processed":    true,
		"processed_at": &now,
	}
	
	if processingError != nil {
		updates["processing_error"] = processingError
	}
	
	return s.db.WithContext(ctx).Model(&WebhookEvent{}).
		Where("stripe_event_id = ?", stripeEventID).
		Updates(updates).Error
}

// IsWebhookEventProcessed checks if a webhook event has already been processed
func (s *Store) IsWebhookEventProcessed(ctx context.Context, stripeEventID string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&WebhookEvent{}).
		Where("stripe_event_id = ? AND processed = ?", stripeEventID, true).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// Transaction support

// WithTransaction executes a function within a database transaction
func (s *Store) WithTransaction(ctx context.Context, fn func(*Store) error) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txStore := &Store{db: tx}
		return fn(txStore)
	})
}

// Health check

// Ping checks if the database connection is alive
func (s *Store) Ping(ctx context.Context) error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.PingContext(ctx)
}
