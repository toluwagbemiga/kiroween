package resolvers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/haunted-saas/graphql-api-gateway/internal/errors"
	"github.com/haunted-saas/graphql-api-gateway/internal/generated"
	"github.com/haunted-saas/graphql-api-gateway/internal/middleware"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	analyticsv1 "github.com/haunted-saas/analytics-service/proto/analytics/v1"
	billingv1 "github.com/haunted-saas/billing-service/proto/billing/v1"
	featureflagsv1 "github.com/haunted-saas/feature-flags-service/proto/featureflags/v1"
	llmv1 "github.com/haunted-saas/llm-gateway-service/proto/llm/v1"
	notificationsv1 "github.com/haunted-saas/notifications-service/proto/notifications/v1"
	userauthv1 "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
)

// Query resolver
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

// ============================================================================
// USER & AUTH QUERIES
// ============================================================================

func (r *queryResolver) Me(ctx context.Context) (*generated.User, error) {
	// GetUser RPC doesn't exist - use ValidateToken as workaround
	token := middleware.GetToken(ctx)
	resp, err := r.clients.UserAuth.ValidateToken(ctx, &userauthv1.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		r.logger.Error("failed to validate token", zap.Error(err))
		return nil, errors.ConvertGRPCError(err)
	}

	return convertUser(resp.User), nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*generated.User, error) {
	// GetUser RPC doesn't exist in proto yet
	return nil, errors.NewBadRequestError("User query not implemented - GetUser RPC missing")
}

func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int) (*generated.UserConnection, error) {
	// ListUsers RPC doesn't exist in proto yet
	return nil, errors.NewBadRequestError("Users query not implemented - ListUsers RPC missing")
}

func (r *queryResolver) MyPermissions(ctx context.Context) ([]string, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.clients.UserAuth.GetUserPermissions(ctx, &userauthv1.GetUserPermissionsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return resp.Permissions, nil
}

func (r *queryResolver) Roles(ctx context.Context) ([]*generated.Role, error) {
	// ListRoles RPC doesn't exist in proto yet
	return nil, errors.NewBadRequestError("Roles query not implemented - ListRoles RPC missing")
}

func (r *queryResolver) Role(ctx context.Context, id string) (*generated.Role, error) {
	// GetRole RPC doesn't exist in proto yet
	return nil, errors.NewBadRequestError("Role query not implemented - GetRole RPC missing")
}

// ============================================================================
// BILLING QUERIES
// ============================================================================

func (r *queryResolver) Plans(ctx context.Context) ([]*generated.Plan, error) {
	resp, err := r.clients.Billing.ListPlans(ctx, &billingv1.ListPlansRequest{})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	plans := make([]*generated.Plan, len(resp.Plans))
	for i, plan := range resp.Plans {
		plans[i] = convertPlan(plan)
	}

	return plans, nil
}

func (r *queryResolver) MySubscription(ctx context.Context) (*generated.Subscription, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.clients.Billing.GetUserSubscription(ctx, &billingv1.GetUserSubscriptionRequest{
		UserId: userID,
	})
	if err != nil {
		// User might not have a subscription - return nil instead of error
		return nil, nil
	}

	return convertSubscription(resp.Subscription), nil
}

func (r *queryResolver) Subscription(ctx context.Context, id string) (*generated.Subscription, error) {
	if err := middleware.RequireAuth(ctx); err != nil {
		return nil, err
	}

	resp, err := r.clients.Billing.GetSubscription(ctx, &billingv1.GetSubscriptionRequest{
		SubscriptionId: id,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return convertSubscription(resp.Subscription), nil
}

func (r *queryResolver) BillingPortalURL(ctx context.Context) (string, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return "", err
	}

	resp, err := r.clients.Billing.CreatePortalSession(ctx, &billingv1.CreatePortalSessionRequest{
		UserId: userID,
	})
	if err != nil {
		return "", errors.ConvertGRPCError(err)
	}

	return resp.Url, nil
}

// ============================================================================
// FEATURE FLAGS QUERIES
// ============================================================================

func (r *queryResolver) IsFeatureEnabled(ctx context.Context, featureName string, properties map[string]interface{}) (bool, error) {
	userID, _ := middleware.GetUserID(ctx)
	teamID := middleware.GetTeamID(ctx)

	propertiesJSON := "{}"
	if properties != nil {
		jsonBytes, err := json.Marshal(properties)
		if err != nil {
			return false, errors.NewBadRequestError("invalid properties")
		}
		propertiesJSON = string(jsonBytes)
	}

	resp, err := r.clients.FeatureFlags.IsFeatureEnabled(ctx, &featureflagsv1.IsFeatureEnabledRequest{
		FeatureName:    featureName,
		UserId:         userID,
		TeamId:         teamID,
		PropertiesJson: propertiesJSON,
	})
	if err != nil {
		return false, errors.ConvertGRPCError(err)
	}

	return resp.Enabled, nil
}

func (r *queryResolver) FeatureVariant(ctx context.Context, featureName string, properties map[string]interface{}) (*generated.FeatureVariant, error) {
	userID, _ := middleware.GetUserID(ctx)
	teamID := middleware.GetTeamID(ctx)

	propertiesJSON := "{}"
	if properties != nil {
		jsonBytes, err := json.Marshal(properties)
		if err != nil {
			return nil, errors.NewBadRequestError("invalid properties")
		}
		propertiesJSON = string(jsonBytes)
	}

	resp, err := r.clients.FeatureFlags.GetFeatureVariant(ctx, &featureflagsv1.GetFeatureVariantRequest{
		FeatureName:    featureName,
		UserId:         userID,
		TeamId:         teamID,
		PropertiesJson: propertiesJSON,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	var payload map[string]interface{}
	if resp.PayloadJson != "" && resp.PayloadJson != "{}" {
		if err := json.Unmarshal([]byte(resp.PayloadJson), &payload); err != nil {
			r.logger.Warn("failed to parse variant payload", zap.Error(err))
		}
	}

	return &generated.FeatureVariant{
		Enabled:     resp.Enabled,
		VariantName: resp.VariantName,
		Payload:     payload,
	}, nil
}

func (r *queryResolver) AvailableFeatures(ctx context.Context) ([]*generated.Feature, error) {
	if err := middleware.RequireRole(ctx, "admin"); err != nil {
		return nil, err
	}

	resp, err := r.clients.FeatureFlags.ListFeatures(ctx, &featureflagsv1.ListFeaturesRequest{})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	features := make([]*generated.Feature, len(resp.Features))
	for i, f := range resp.Features {
		features[i] = convertFeature(f)
	}

	return features, nil
}

// ============================================================================
// LLM GATEWAY QUERIES
// ============================================================================

func (r *queryResolver) AvailablePrompts(ctx context.Context) ([]*generated.PromptMetadata, error) {
	if err := middleware.RequireAuth(ctx); err != nil {
		return nil, err
	}

	resp, err := r.clients.LLMGateway.ListPrompts(ctx, &llmv1.ListPromptsRequest{})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	prompts := make([]*generated.PromptMetadata, len(resp.Prompts))
	for i, p := range resp.Prompts {
		prompts[i] = convertPromptMetadata(p)
	}

	return prompts, nil
}

func (r *queryResolver) PromptDetails(ctx context.Context, name string) (*generated.PromptMetadata, error) {
	if err := middleware.RequireAuth(ctx); err != nil {
		return nil, err
	}

	resp, err := r.clients.LLMGateway.GetPromptMetadata(ctx, &llmv1.GetPromptMetadataRequest{
		PromptName: name,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return convertPromptMetadata(resp.Metadata), nil
}

func (r *queryResolver) MyLLMUsage(ctx context.Context) (*generated.LLMUsageStats, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.clients.LLMGateway.GetUsageStats(ctx, &llmv1.GetUsageStatsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	var callsByModel map[string]interface{}
	if resp.CallsByModelJson != "" {
		if err := json.Unmarshal([]byte(resp.CallsByModelJson), &callsByModel); err != nil {
			r.logger.Warn("failed to parse calls by model", zap.Error(err))
		}
	}

	return &generated.LLMUsageStats{
		TotalCalls:  int(resp.TotalCalls),
		TotalTokens: int(resp.TotalTokens),
		TotalCost:   resp.TotalCost,
		CallsByModel: callsByModel,
	}, nil
}

// ============================================================================
// NOTIFICATIONS QUERIES
// ============================================================================

func (r *queryResolver) NotificationToken(ctx context.Context) (*generated.NotificationToken, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.clients.Notifications.GenerateConnectionToken(ctx, &notificationsv1.GenerateConnectionTokenRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return &generated.NotificationToken{
		Token:     resp.Token,
		SocketURL: resp.SocketUrl,
		ExpiresAt: resp.ExpiresAt.AsTime(),
	}, nil
}

func (r *queryResolver) MyNotificationPreferences(ctx context.Context) (*generated.NotificationPreferences, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.clients.Notifications.GetPreferences(ctx, &notificationsv1.GetPreferencesRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	var channels map[string]interface{}
	if resp.ChannelsJson != "" {
		if err := json.Unmarshal([]byte(resp.ChannelsJson), &channels); err != nil {
			r.logger.Warn("failed to parse channels", zap.Error(err))
		}
	}

	return &generated.NotificationPreferences{
		EmailEnabled: resp.EmailEnabled,
		PushEnabled:  resp.PushEnabled,
		InAppEnabled: resp.InAppEnabled,
		Channels:     channels,
	}, nil
}

// ============================================================================
// ANALYTICS QUERIES
// ============================================================================

func (r *queryResolver) MyAnalytics(ctx context.Context, startDate *time.Time, endDate *time.Time) (*generated.AnalyticsSummary, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	req := &analyticsv1.GetAnalyticsSummaryRequest{
		UserId: userID,
	}

	if startDate != nil {
		req.StartDate = timestamppb.New(*startDate)
	}
	if endDate != nil {
		req.EndDate = timestamppb.New(*endDate)
	}

	resp, err := r.clients.Analytics.GetAnalyticsSummary(ctx, req)
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	var eventsByType map[string]interface{}
	if resp.EventsByTypeJson != "" {
		if err := json.Unmarshal([]byte(resp.EventsByTypeJson), &eventsByType); err != nil {
			r.logger.Warn("failed to parse events by type", zap.Error(err))
		}
	}

	topEvents := make([]*generated.EventCount, len(resp.TopEvents))
	for i, e := range resp.TopEvents {
		topEvents[i] = &generated.EventCount{
			EventName: e.EventName,
			Count:     int(e.Count),
		}
	}

	return &generated.AnalyticsSummary{
		TotalEvents:  int(resp.TotalEvents),
		UniqueUsers:  int(resp.UniqueUsers),
		EventsByType: eventsByType,
		TopEvents:    topEvents,
	}, nil
}
