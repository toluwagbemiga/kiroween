package dataloader

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"go.uber.org/zap"

	billingv1 "github.com/haunted-saas/billing-service/proto/billing/v1"
	userauthv1 "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
)

// Loaders holds all dataloaders
type Loaders struct {
	UserByID         *dataloader.Loader[string, *userauthv1.User]
	SubscriptionByID *dataloader.Loader[string, *billingv1.Subscription]
	PlanByID         *dataloader.Loader[string, *billingv1.Plan]
}

// contextKey is the type for context keys
type contextKey string

const loadersKey contextKey = "dataloaders"

// Clients holds gRPC clients needed for dataloaders
type Clients struct {
	UserAuth userauthv1.UserAuthServiceClient
	Billing  billingv1.BillingServiceClient
}

// NewLoaders creates a new set of dataloaders
func NewLoaders(clients Clients, logger *zap.Logger) *Loaders {
	return &Loaders{
		UserByID: dataloader.NewBatchedLoader(
			userBatchFunc(clients.UserAuth, logger),
			dataloader.WithWait[string, *userauthv1.User](10*time.Millisecond),
			dataloader.WithBatchCapacity[string, *userauthv1.User](100),
		),
		SubscriptionByID: dataloader.NewBatchedLoader(
			subscriptionBatchFunc(clients.Billing, logger),
			dataloader.WithWait[string, *billingv1.Subscription](10*time.Millisecond),
			dataloader.WithBatchCapacity[string, *billingv1.Subscription](100),
		),
		PlanByID: dataloader.NewBatchedLoader(
			planBatchFunc(clients.Billing, logger),
			dataloader.WithWait[string, *billingv1.Plan](10*time.Millisecond),
			dataloader.WithBatchCapacity[string, *billingv1.Plan](100),
		),
	}
}

// Middleware injects dataloaders into the request context
func Middleware(loaders *Loaders) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// For returns the dataloaders from context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// userBatchFunc batches user lookups
// Note: GetUser RPC doesn't exist in proto yet, so this is a placeholder
// Users will be fetched directly in resolvers for now
func userBatchFunc(client userauthv1.UserAuthServiceClient, logger *zap.Logger) dataloader.BatchFunc[string, *userauthv1.User] {
	return func(ctx context.Context, userIDs []string) []*dataloader.Result[*userauthv1.User] {
		logger.Debug("batching user lookups", zap.Int("count", len(userIDs)))

		results := make([]*dataloader.Result[*userauthv1.User], len(userIDs))

		// TODO: Implement batch user lookup when GetUser RPC is added to proto
		// For now, return errors - users will be fetched via ValidateToken in resolvers
		for i, userID := range userIDs {
			logger.Warn("user dataloader not implemented - GetUser RPC missing from proto", zap.String("user_id", userID))
			results[i] = &dataloader.Result[*userauthv1.User]{
				Error: errors.New("GetUser RPC not implemented"),
			}
		}

		return results
	}
}

// subscriptionBatchFunc batches subscription lookups
func subscriptionBatchFunc(client billingv1.BillingServiceClient, logger *zap.Logger) dataloader.BatchFunc[string, *billingv1.Subscription] {
	return func(ctx context.Context, subscriptionIDs []string) []*dataloader.Result[*billingv1.Subscription] {
		logger.Debug("batching subscription lookups", zap.Int("count", len(subscriptionIDs)))

		results := make([]*dataloader.Result[*billingv1.Subscription], len(subscriptionIDs))

		// GetSubscription expects team_id, not subscription_id
		// This dataloader needs to be redesigned or removed
		for i, subID := range subscriptionIDs {
			resp, err := client.GetSubscription(ctx, &billingv1.GetSubscriptionRequest{
				TeamId: subID, // Using subID as teamID for now
			})

			if err != nil {
				logger.Error("failed to fetch subscription", zap.String("team_id", subID), zap.Error(err))
				results[i] = &dataloader.Result[*billingv1.Subscription]{Error: err}
			} else {
				results[i] = &dataloader.Result[*billingv1.Subscription]{Data: resp.Subscription}
			}
		}

		return results
	}
}

// planBatchFunc batches plan lookups
func planBatchFunc(client billingv1.BillingServiceClient, logger *zap.Logger) dataloader.BatchFunc[string, *billingv1.Plan] {
	return func(ctx context.Context, planIDs []string) []*dataloader.Result[*billingv1.Plan] {
		logger.Debug("batching plan lookups", zap.Int("count", len(planIDs)))

		results := make([]*dataloader.Result[*billingv1.Plan], len(planIDs))

		// Fetch all plans once
		plansResp, err := client.ListPlans(ctx, &billingv1.ListPlansRequest{})
		if err != nil {
			logger.Error("failed to fetch plans", zap.Error(err))
			for i := range results {
				results[i] = &dataloader.Result[*billingv1.Plan]{Error: err}
			}
			return results
		}

		// Create a map for quick lookup
		planMap := make(map[string]*billingv1.Plan)
		for _, plan := range plansResp.Plans {
			planMap[plan.Id] = plan
		}

		// Match plans to requested IDs
		for i, planID := range planIDs {
			if plan, ok := planMap[planID]; ok {
				results[i] = &dataloader.Result[*billingv1.Plan]{Data: plan}
			} else {
				results[i] = &dataloader.Result[*billingv1.Plan]{
					Error: errors.New("plan not found"),
				}
			}
		}

		return results
	}
}
