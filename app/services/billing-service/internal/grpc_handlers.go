package internal

import (
	"context"
	"time"

	"github.com/haunted-saas/billing-service/internal/db"
	pb "github.com/haunted-saas/billing-service/proto/billing/v1"
	"github.com/stripe/stripe-go/v76"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// BillingServiceServer implements the gRPC billing service
type BillingServiceServer struct {
	pb.UnimplementedBillingServiceServer
	stripeClient *StripeClient
	store        *db.Store
	logger       *zap.Logger
}

// NewBillingServiceServer creates a new billing service server
func NewBillingServiceServer(stripeClient *StripeClient, store *db.Store, logger *zap.Logger) *BillingServiceServer {
	return &BillingServiceServer{
		stripeClient: stripeClient,
		store:        store,
		logger:       logger,
	}
}

// Plan Management

// CreatePlan creates a new subscription plan
func (s *BillingServiceServer) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest) (*pb.CreatePlanResponse, error) {
	s.logger.Info("creating plan",
		zap.String("name", req.Name),
		zap.Int64("price_cents", req.PriceCents),
		zap.String("interval", req.BillingInterval))
	
	// Validate input
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "plan name is required")
	}
	if req.PriceCents < 0 {
		return nil, status.Error(codes.InvalidArgument, "price cannot be negative")
	}
	if req.BillingInterval != "month" && req.BillingInterval != "year" {
		return nil, status.Error(codes.InvalidArgument, "billing interval must be 'month' or 'year'")
	}
	
	currency := req.Currency
	if currency == "" {
		currency = "usd"
	}
	
	// Create Stripe product
	stripeProduct, err := s.stripeClient.CreateProduct(req.Name, map[string]string{
		"created_by": req.CreatedByUserId,
	})
	if err != nil {
		s.logger.Error("failed to create Stripe product", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}
	
	// Create Stripe price
	stripePrice, err := s.stripeClient.CreatePrice(
		stripeProduct.ID,
		req.PriceCents,
		currency,
		req.BillingInterval,
	)
	if err != nil {
		s.logger.Error("failed to create Stripe price", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create price: %v", err)
	}
	
	// Create plan in database
	plan := &db.Plan{
		Name:            req.Name,
		PriceCents:      req.PriceCents,
		Currency:        currency,
		BillingInterval: req.BillingInterval,
		Features:        req.Features,
		IsActive:        true,
		StripePriceID:   stripePrice.ID,
		StripeProductID: stripeProduct.ID,
		TrialDays:       req.TrialDays,
	}
	
	if req.CreatedByUserId != "" {
		plan.CreatedByUserID = &req.CreatedByUserId
	}
	
	if err := s.store.CreatePlan(ctx, plan); err != nil {
		s.logger.Error("failed to create plan in database", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to save plan: %v", err)
	}
	
	s.logger.Info("plan created successfully",
		zap.String("plan_id", plan.ID),
		zap.String("stripe_price_id", plan.StripePriceID))
	
	return &pb.CreatePlanResponse{
		Plan: dbPlanToProto(plan),
	}, nil
}

// GetPlan retrieves a plan by ID
func (s *BillingServiceServer) GetPlan(ctx context.Context, req *pb.GetPlanRequest) (*pb.GetPlanResponse, error) {
	if req.PlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "plan_id is required")
	}
	
	plan, err := s.store.GetPlanByID(ctx, req.PlanId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		s.logger.Error("failed to get plan", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get plan: %v", err)
	}
	
	return &pb.GetPlanResponse{
		Plan: dbPlanToProto(plan),
	}, nil
}

// ListPlans lists all plans
func (s *BillingServiceServer) ListPlans(ctx context.Context, req *pb.ListPlansRequest) (*pb.ListPlansResponse, error) {
	plans, err := s.store.ListPlans(ctx, req.ActiveOnly)
	if err != nil {
		s.logger.Error("failed to list plans", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list plans: %v", err)
	}
	
	pbPlans := make([]*pb.Plan, len(plans))
	for i, plan := range plans {
		pbPlans[i] = dbPlanToProto(&plan)
	}
	
	return &pb.ListPlansResponse{
		Plans: pbPlans,
	}, nil
}

// UpdatePlan updates a plan
func (s *BillingServiceServer) UpdatePlan(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.UpdatePlanResponse, error) {
	if req.PlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "plan_id is required")
	}
	
	plan, err := s.store.GetPlanByID(ctx, req.PlanId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get plan: %v", err)
	}
	
	// Update fields
	if req.Name != "" {
		plan.Name = req.Name
		// Update Stripe product name
		if _, err := s.stripeClient.UpdateProduct(plan.StripeProductID, req.Name, nil); err != nil {
			s.logger.Error("failed to update Stripe product", zap.Error(err))
		}
	}
	
	if req.Features != nil {
		plan.Features = req.Features
	}
	
	if err := s.store.UpdatePlan(ctx, plan); err != nil {
		s.logger.Error("failed to update plan", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update plan: %v", err)
	}
	
	return &pb.UpdatePlanResponse{
		Plan: dbPlanToProto(plan),
	}, nil
}

// DeactivatePlan deactivates a plan
func (s *BillingServiceServer) DeactivatePlan(ctx context.Context, req *pb.DeactivatePlanRequest) (*pb.DeactivatePlanResponse, error) {
	if req.PlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "plan_id is required")
	}
	
	if err := s.store.DeactivatePlan(ctx, req.PlanId); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		s.logger.Error("failed to deactivate plan", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to deactivate plan: %v", err)
	}
	
	s.logger.Info("plan deactivated", zap.String("plan_id", req.PlanId))
	
	return &pb.DeactivatePlanResponse{
		Success: true,
	}, nil
}

// Subscription Management

// CreateCheckoutSession creates a Stripe Checkout session
func (s *BillingServiceServer) CreateCheckoutSession(ctx context.Context, req *pb.CreateCheckoutSessionRequest) (*pb.CreateCheckoutSessionResponse, error) {
	s.logger.Info("creating checkout session",
		zap.String("team_id", req.TeamId),
		zap.String("plan_id", req.PlanId))
	
	// Validate input
	if req.TeamId == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}
	if req.PlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "plan_id is required")
	}
	if req.SuccessUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "success_url is required")
	}
	if req.CancelUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "cancel_url is required")
	}
	
	// Get plan
	plan, err := s.store.GetPlanByID(ctx, req.PlanId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get plan: %v", err)
	}
	
	if !plan.IsActive {
		return nil, status.Error(codes.FailedPrecondition, "plan is not active")
	}
	
	// Check if team already has a subscription
	existingSub, err := s.store.GetSubscriptionByTeamID(ctx, req.TeamId)
	if err == nil && existingSub != nil && existingSub.IsActive() {
		return nil, status.Error(codes.AlreadyExists, "team already has an active subscription")
	}
	
	// Create or get Stripe customer
	var customerID string
	if existingSub != nil {
		customerID = existingSub.StripeCustomerID
	} else if req.CustomerEmail != "" {
		customer, err := s.stripeClient.CreateCustomer(req.CustomerEmail, req.TeamId, nil)
		if err != nil {
			s.logger.Error("failed to create Stripe customer", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "failed to create customer: %v", err)
		}
		customerID = customer.ID
	}
	
	// Create checkout session
	session, err := s.stripeClient.CreateCheckoutSession(
		plan.StripePriceID,
		customerID,
		req.SuccessUrl,
		req.CancelUrl,
		map[string]string{
			"team_id": req.TeamId,
			"plan_id": req.PlanId,
		},
		plan.TrialDays,
	)
	if err != nil {
		s.logger.Error("failed to create checkout session", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create checkout session: %v", err)
	}
	
	s.logger.Info("checkout session created",
		zap.String("session_id", session.ID),
		zap.String("team_id", req.TeamId))
	
	return &pb.CreateCheckoutSessionResponse{
		CheckoutUrl: session.URL,
		SessionId:   session.ID,
	}, nil
}

// GetSubscription retrieves a subscription by team ID
func (s *BillingServiceServer) GetSubscription(ctx context.Context, req *pb.GetSubscriptionRequest) (*pb.GetSubscriptionResponse, error) {
	if req.TeamId == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}
	
	subscription, err := s.store.GetSubscriptionByTeamID(ctx, req.TeamId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "subscription not found")
		}
		s.logger.Error("failed to get subscription", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get subscription: %v", err)
	}
	
	return &pb.GetSubscriptionResponse{
		Subscription: dbSubscriptionToProto(subscription),
	}, nil
}

// CancelSubscription cancels a subscription
func (s *BillingServiceServer) CancelSubscription(ctx context.Context, req *pb.CancelSubscriptionRequest) (*pb.CancelSubscriptionResponse, error) {
	s.logger.Info("canceling subscription",
		zap.String("team_id", req.TeamId),
		zap.Bool("immediate", req.Immediate))
	
	if req.TeamId == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}
	
	// Get subscription
	subscription, err := s.store.GetSubscriptionByTeamID(ctx, req.TeamId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "subscription not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get subscription: %v", err)
	}
	
	// Cancel in Stripe
	cancelAtPeriodEnd := !req.Immediate
	stripeSub, err := s.stripeClient.CancelSubscription(subscription.StripeSubscriptionID, cancelAtPeriodEnd)
	if err != nil {
		s.logger.Error("failed to cancel Stripe subscription", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to cancel subscription: %v", err)
	}
	
	// Update in database
	subscription.Status = string(stripeSub.Status)
	if stripeSub.CancelAt > 0 {
		cancelAt := time.Unix(stripeSub.CancelAt, 0)
		subscription.CancelAt = &cancelAt
	}
	if stripeSub.CanceledAt > 0 {
		canceledAt := time.Unix(stripeSub.CanceledAt, 0)
		subscription.CanceledAt = &canceledAt
	}
	
	if err := s.store.UpdateSubscription(ctx, subscription); err != nil {
		s.logger.Error("failed to update subscription", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update subscription: %v", err)
	}
	
	cancellationDate := subscription.CurrentPeriodEnd.Format("2006-01-02")
	if req.Immediate {
		cancellationDate = time.Now().Format("2006-01-02")
	}
	
	s.logger.Info("subscription canceled",
		zap.String("team_id", req.TeamId),
		zap.String("cancellation_date", cancellationDate))
	
	return &pb.CancelSubscriptionResponse{
		CancellationDate: cancellationDate,
		Subscription:     dbSubscriptionToProto(subscription),
	}, nil
}

// UpdateSubscription updates a subscription (change plan)
func (s *BillingServiceServer) UpdateSubscription(ctx context.Context, req *pb.UpdateSubscriptionRequest) (*pb.UpdateSubscriptionResponse, error) {
	s.logger.Info("updating subscription",
		zap.String("team_id", req.TeamId),
		zap.String("new_plan_id", req.NewPlanId))
	
	if req.TeamId == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}
	if req.NewPlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "new_plan_id is required")
	}
	
	// Get current subscription
	subscription, err := s.store.GetSubscriptionByTeamID(ctx, req.TeamId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "subscription not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get subscription: %v", err)
	}
	
	// Get new plan
	newPlan, err := s.store.GetPlanByID(ctx, req.NewPlanId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get plan: %v", err)
	}
	
	if !newPlan.IsActive {
		return nil, status.Error(codes.FailedPrecondition, "plan is not active")
	}
	
	// Update subscription in Stripe with proration
	stripeSub, err := s.stripeClient.UpdateSubscription(
		subscription.StripeSubscriptionID,
		newPlan.StripePriceID,
		"create_prorations",
	)
	if err != nil {
		s.logger.Error("failed to update Stripe subscription", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update subscription: %v", err)
	}
	
	// Update in database
	subscription.PlanID = newPlan.ID
	subscription.Status = string(stripeSub.Status)
	subscription.CurrentPeriodStart = time.Unix(stripeSub.CurrentPeriodStart, 0)
	subscription.CurrentPeriodEnd = time.Unix(stripeSub.CurrentPeriodEnd, 0)
	
	if err := s.store.UpdateSubscription(ctx, subscription); err != nil {
		s.logger.Error("failed to update subscription in database", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update subscription: %v", err)
	}
	
	// Get upcoming invoice for proration amount
	var prorationAmount int64
	upcomingInvoice, err := s.stripeClient.GetUpcomingInvoice(subscription.StripeCustomerID)
	if err == nil {
		prorationAmount = upcomingInvoice.AmountDue
	}
	
	s.logger.Info("subscription updated",
		zap.String("team_id", req.TeamId),
		zap.String("new_plan_id", newPlan.ID))
	
	return &pb.UpdateSubscriptionResponse{
		Subscription:          dbSubscriptionToProto(subscription),
		NextBillingAmountCents: newPlan.PriceCents,
		ProrationAmountCents:   prorationAmount,
	}, nil
}

// CreateCustomerPortalSession creates a Stripe Customer Portal session
func (s *BillingServiceServer) CreateCustomerPortalSession(ctx context.Context, req *pb.CreateCustomerPortalSessionRequest) (*pb.CreateCustomerPortalSessionResponse, error) {
	if req.TeamId == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}
	if req.ReturnUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "return_url is required")
	}
	
	// Get subscription to get customer ID
	subscription, err := s.store.GetSubscriptionByTeamID(ctx, req.TeamId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "subscription not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get subscription: %v", err)
	}
	
	// Create portal session
	portalSession, err := s.stripeClient.CreateCustomerPortalSession(subscription.StripeCustomerID, req.ReturnUrl)
	if err != nil {
		s.logger.Error("failed to create customer portal session", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create portal session: %v", err)
	}
	
	return &pb.CreateCustomerPortalSessionResponse{
		PortalUrl: portalSession.URL,
	}, nil
}

// GetUpcomingInvoice retrieves the upcoming invoice for a team
func (s *BillingServiceServer) GetUpcomingInvoice(ctx context.Context, req *pb.GetUpcomingInvoiceRequest) (*pb.GetUpcomingInvoiceResponse, error) {
	if req.TeamId == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}
	
	subscription, err := s.store.GetSubscriptionByTeamID(ctx, req.TeamId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "subscription not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get subscription: %v", err)
	}
	
	invoice, err := s.stripeClient.GetUpcomingInvoice(subscription.StripeCustomerID)
	if err != nil {
		s.logger.Error("failed to get upcoming invoice", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get invoice: %v", err)
	}
	
	return &pb.GetUpcomingInvoiceResponse{
		Invoice: stripeInvoiceToProto(invoice),
	}, nil
}

// ListInvoices lists invoices for a team
func (s *BillingServiceServer) ListInvoices(ctx context.Context, req *pb.ListInvoicesRequest) (*pb.ListInvoicesResponse, error) {
	if req.TeamId == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}
	
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	
	subscription, err := s.store.GetSubscriptionByTeamID(ctx, req.TeamId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "subscription not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get subscription: %v", err)
	}
	
	invoices, err := s.stripeClient.ListInvoices(subscription.StripeCustomerID, int64(limit))
	if err != nil {
		s.logger.Error("failed to list invoices", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list invoices: %v", err)
	}
	
	pbInvoices := make([]*pb.Invoice, len(invoices))
	for i, inv := range invoices {
		pbInvoices[i] = stripeInvoiceToProto(inv)
	}
	
	return &pb.ListInvoicesResponse{
		Invoices: pbInvoices,
	}, nil
}

// Helper functions to convert between database and proto models

func dbPlanToProto(plan *db.Plan) *pb.Plan {
	return &pb.Plan{
		Id:              plan.ID,
		Name:            plan.Name,
		PriceCents:      plan.PriceCents,
		Currency:        plan.Currency,
		BillingInterval: plan.BillingInterval,
		Features:        plan.Features,
		IsActive:        plan.IsActive,
		StripePriceId:   plan.StripePriceID,
		StripeProductId: plan.StripeProductID,
		TrialDays:       plan.TrialDays,
		CreatedAt:       timestamppb.New(plan.CreatedAt),
		UpdatedAt:       timestamppb.New(plan.UpdatedAt),
	}
}

func dbSubscriptionToProto(sub *db.Subscription) *pb.Subscription {
	pbSub := &pb.Subscription{
		Id:                   sub.ID,
		TeamId:               sub.TeamID,
		PlanId:               sub.PlanID,
		Status:               sub.Status,
		StripeSubscriptionId: sub.StripeSubscriptionID,
		StripeCustomerId:     sub.StripeCustomerID,
		CurrentPeriodStart:   timestamppb.New(sub.CurrentPeriodStart),
		CurrentPeriodEnd:     timestamppb.New(sub.CurrentPeriodEnd),
		CreatedAt:            timestamppb.New(sub.CreatedAt),
		UpdatedAt:            timestamppb.New(sub.UpdatedAt),
	}
	
	if sub.CancelAt != nil {
		pbSub.CancelAt = timestamppb.New(*sub.CancelAt)
	}
	if sub.CanceledAt != nil {
		pbSub.CanceledAt = timestamppb.New(*sub.CanceledAt)
	}
	if sub.TrialEnd != nil {
		pbSub.TrialEnd = timestamppb.New(*sub.TrialEnd)
	}
	if sub.Plan.ID != "" {
		pbSub.Plan = dbPlanToProto(&sub.Plan)
	}
	
	return pbSub
}

func stripeInvoiceToProto(inv *stripe.Invoice) *pb.Invoice {
	pbInv := &pb.Invoice{
		Id:         inv.ID,
		AmountDue:  inv.AmountDue,
		AmountPaid: inv.AmountPaid,
		Currency:   string(inv.Currency),
		Status:     string(inv.Status),
		Created:    timestamppb.New(time.Unix(inv.Created, 0)),
	}
	
	if inv.DueDate > 0 {
		pbInv.DueDate = timestamppb.New(time.Unix(inv.DueDate, 0))
	}
	if inv.InvoicePDF != "" {
		pbInv.InvoicePdf = inv.InvoicePDF
	}
	if inv.HostedInvoiceURL != "" {
		pbInv.HostedInvoiceUrl = inv.HostedInvoiceURL
	}
	
	return pbInv
}
