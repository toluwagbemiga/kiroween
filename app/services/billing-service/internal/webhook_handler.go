package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/haunted-saas/billing-service/internal/db"
	"github.com/stripe/stripe-go/v76"
	"go.uber.org/zap"
)

// WebhookHandler handles Stripe webhook events
type WebhookHandler struct {
	stripeClient  *StripeClient
	store         *db.Store
	webhookSecret string
	logger        *zap.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(stripeClient *StripeClient, store *db.Store, webhookSecret string, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{
		stripeClient:  stripeClient,
		store:         store,
		webhookSecret: webhookSecret,
		logger:        logger,
	}
}

// HandleWebhook handles incoming Stripe webhook requests
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Read the request body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read webhook payload", zap.Error(err))
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	
	// Get the Stripe signature from headers
	signature := r.Header.Get("Stripe-Signature")
	if signature == "" {
		h.logger.Warn("webhook request missing Stripe-Signature header")
		http.Error(w, "missing Stripe-Signature header", http.StatusBadRequest)
		return
	}
	
	// Verify the webhook signature
	event, err := h.stripeClient.ConstructEvent(payload, signature, h.webhookSecret)
	if err != nil {
		h.logger.Error("webhook signature verification failed",
			zap.Error(err),
			zap.String("signature", signature))
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}
	
	h.logger.Info("webhook received",
		zap.String("event_id", event.ID),
		zap.String("event_type", string(event.Type)))
	
	// Check idempotency - has this event already been processed?
	processed, err := h.store.IsWebhookEventProcessed(ctx, event.ID)
	if err != nil {
		h.logger.Error("failed to check webhook idempotency",
			zap.Error(err),
			zap.String("event_id", event.ID))
		// Continue processing - better to risk duplicate than miss an event
	}
	
	if processed {
		h.logger.Info("webhook event already processed (idempotent)",
			zap.String("event_id", event.ID))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "already_processed"})
		return
	}
	
	// Record the webhook event
	webhookEvent := &db.WebhookEvent{
		StripeEventID: event.ID,
		EventType:     string(event.Type),
		Processed:     false,
		ReceivedAt:    time.Now(),
	}
	
	if err := h.store.CreateWebhookEvent(ctx, webhookEvent); err != nil {
		h.logger.Error("failed to create webhook event record",
			zap.Error(err),
			zap.String("event_id", event.ID))
		// Continue processing even if we can't record it
	}
	
	// Process the event based on type
	var processingError *string
	if err := h.processEvent(ctx, event); err != nil {
		h.logger.Error("failed to process webhook event",
			zap.Error(err),
			zap.String("event_id", event.ID),
			zap.String("event_type", string(event.Type)))
		errMsg := err.Error()
		processingError = &errMsg
	}
	
	// Mark the event as processed
	if err := h.store.MarkWebhookEventProcessed(ctx, event.ID, processingError); err != nil {
		h.logger.Error("failed to mark webhook event as processed",
			zap.Error(err),
			zap.String("event_id", event.ID))
	}
	
	// Always return 200 OK to Stripe to prevent retries
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}

// processEvent processes a Stripe event
func (h *WebhookHandler) processEvent(ctx context.Context, event stripe.Event) error {
	switch event.Type {
	case "checkout.session.completed":
		return h.handleCheckoutSessionCompleted(ctx, event)
	
	case "customer.subscription.created":
		return h.handleSubscriptionCreated(ctx, event)
	
	case "customer.subscription.updated":
		return h.handleSubscriptionUpdated(ctx, event)
	
	case "customer.subscription.deleted":
		return h.handleSubscriptionDeleted(ctx, event)
	
	case "invoice.payment_succeeded":
		return h.handleInvoicePaymentSucceeded(ctx, event)
	
	case "invoice.payment_failed":
		return h.handleInvoicePaymentFailed(ctx, event)
	
	case "customer.subscription.trial_will_end":
		return h.handleTrialWillEnd(ctx, event)
	
	default:
		h.logger.Info("unhandled webhook event type",
			zap.String("event_type", string(event.Type)))
		return nil
	}
}

// handleCheckoutSessionCompleted handles checkout.session.completed events
func (h *WebhookHandler) handleCheckoutSessionCompleted(ctx context.Context, event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return fmt.Errorf("failed to unmarshal checkout session: %w", err)
	}
	
	// Extract team_id from metadata
	teamID, ok := session.Subscription.Metadata["team_id"]
	if !ok || teamID == "" {
		return fmt.Errorf("team_id not found in subscription metadata")
	}
	
	h.logger.Info("processing checkout session completed",
		zap.String("session_id", session.ID),
		zap.String("team_id", teamID),
		zap.String("subscription_id", session.Subscription.ID))
	
	// Get the subscription details from Stripe
	stripeSub, err := h.stripeClient.GetSubscription(session.Subscription.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription from Stripe: %w", err)
	}
	
	// Get the plan from our database using the Stripe price ID
	plan, err := h.store.GetPlanByStripePriceID(ctx, stripeSub.Items.Data[0].Price.ID)
	if err != nil {
		return fmt.Errorf("failed to get plan: %w", err)
	}
	
	// Create or update subscription in our database
	subscription := &db.Subscription{
		TeamID:               teamID,
		PlanID:               plan.ID,
		Status:               string(stripeSub.Status),
		StripeSubscriptionID: stripeSub.ID,
		StripeCustomerID:     stripeSub.Customer.ID,
		CurrentPeriodStart:   time.Unix(stripeSub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(stripeSub.CurrentPeriodEnd, 0),
	}
	
	if stripeSub.TrialEnd > 0 {
		trialEnd := time.Unix(stripeSub.TrialEnd, 0)
		subscription.TrialEnd = &trialEnd
	}
	
	// Check if subscription already exists
	existing, err := h.store.GetSubscriptionByTeamID(ctx, teamID)
	if err == nil && existing != nil {
		// Update existing subscription
		subscription.ID = existing.ID
		if err := h.store.UpdateSubscription(ctx, subscription); err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}
	} else {
		// Create new subscription
		if err := h.store.CreateSubscription(ctx, subscription); err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}
	}
	
	h.logger.Info("subscription provisioned",
		zap.String("team_id", teamID),
		zap.String("plan_id", plan.ID),
		zap.String("status", subscription.Status))
	
	// TODO: Call feature-flags-service or user-auth-service to provision access
	// This would be a gRPC call to enable features for the team
	
	return nil
}

// handleSubscriptionCreated handles customer.subscription.created events
func (h *WebhookHandler) handleSubscriptionCreated(ctx context.Context, event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}
	
	h.logger.Info("subscription created",
		zap.String("subscription_id", stripeSub.ID),
		zap.String("status", string(stripeSub.Status)))
	
	// This is usually handled by checkout.session.completed
	// But we log it for monitoring
	return nil
}

// handleSubscriptionUpdated handles customer.subscription.updated events
func (h *WebhookHandler) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}
	
	h.logger.Info("subscription updated",
		zap.String("subscription_id", stripeSub.ID),
		zap.String("status", string(stripeSub.Status)))
	
	// Get subscription from our database
	subscription, err := h.store.GetSubscriptionByStripeID(ctx, stripeSub.ID)
	if err != nil {
		h.logger.Warn("subscription not found in database",
			zap.String("subscription_id", stripeSub.ID))
		return nil // Not an error - might be a subscription we don't track
	}
	
	// Update subscription details
	subscription.Status = string(stripeSub.Status)
	subscription.CurrentPeriodStart = time.Unix(stripeSub.CurrentPeriodStart, 0)
	subscription.CurrentPeriodEnd = time.Unix(stripeSub.CurrentPeriodEnd, 0)
	
	if stripeSub.CancelAt > 0 {
		cancelAt := time.Unix(stripeSub.CancelAt, 0)
		subscription.CancelAt = &cancelAt
	}
	
	if stripeSub.CanceledAt > 0 {
		canceledAt := time.Unix(stripeSub.CanceledAt, 0)
		subscription.CanceledAt = &canceledAt
	}
	
	if err := h.store.UpdateSubscription(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	
	return nil
}

// handleSubscriptionDeleted handles customer.subscription.deleted events
func (h *WebhookHandler) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}
	
	h.logger.Info("subscription deleted",
		zap.String("subscription_id", stripeSub.ID))
	
	// Get subscription from our database
	subscription, err := h.store.GetSubscriptionByStripeID(ctx, stripeSub.ID)
	if err != nil {
		h.logger.Warn("subscription not found in database",
			zap.String("subscription_id", stripeSub.ID))
		return nil
	}
	
	// Update status to canceled
	subscription.Status = "canceled"
	now := time.Now()
	subscription.CanceledAt = &now
	
	if err := h.store.UpdateSubscription(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	
	h.logger.Info("subscription access revoked",
		zap.String("team_id", subscription.TeamID))
	
	// TODO: Call feature-flags-service or user-auth-service to revoke access
	// This would be a gRPC call to disable features for the team
	
	return nil
}

// handleInvoicePaymentSucceeded handles invoice.payment_succeeded events
func (h *WebhookHandler) handleInvoicePaymentSucceeded(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %w", err)
	}
	
	h.logger.Info("invoice payment succeeded",
		zap.String("invoice_id", invoice.ID),
		zap.String("subscription_id", invoice.Subscription.ID),
		zap.Int64("amount_paid", invoice.AmountPaid))
	
	// Payment succeeded - ensure subscription is active
	if invoice.Subscription != nil {
		subscription, err := h.store.GetSubscriptionByStripeID(ctx, invoice.Subscription.ID)
		if err == nil && subscription != nil {
			if subscription.Status != "active" {
				subscription.Status = "active"
				h.store.UpdateSubscription(ctx, subscription)
			}
		}
	}
	
	return nil
}

// handleInvoicePaymentFailed handles invoice.payment_failed events
func (h *WebhookHandler) handleInvoicePaymentFailed(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %w", err)
	}
	
	h.logger.Warn("invoice payment failed",
		zap.String("invoice_id", invoice.ID),
		zap.String("subscription_id", invoice.Subscription.ID),
		zap.Int64("amount_due", invoice.AmountDue))
	
	// Payment failed - update subscription status
	if invoice.Subscription != nil {
		subscription, err := h.store.GetSubscriptionByStripeID(ctx, invoice.Subscription.ID)
		if err == nil && subscription != nil {
			subscription.Status = "past_due"
			if err := h.store.UpdateSubscription(ctx, subscription); err != nil {
				return fmt.Errorf("failed to update subscription: %w", err)
			}
			
			h.logger.Info("subscription marked as past_due",
				zap.String("team_id", subscription.TeamID))
			
			// TODO: Notify team about payment failure
			// TODO: Consider revoking access after grace period
		}
	}
	
	return nil
}

// handleTrialWillEnd handles customer.subscription.trial_will_end events
func (h *WebhookHandler) handleTrialWillEnd(ctx context.Context, event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}
	
	h.logger.Info("trial will end soon",
		zap.String("subscription_id", stripeSub.ID),
		zap.Int64("trial_end", stripeSub.TrialEnd))
	
	// TODO: Send notification to team about trial ending
	// This is typically sent 3 days before trial ends
	
	return nil
}
