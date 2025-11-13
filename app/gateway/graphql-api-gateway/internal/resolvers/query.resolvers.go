package resolvers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/haunted-saas/graphql-api-gateway/internal/errors"
	"github.com/haunted-saas/graphql-api-gateway/internal/generated"
	"github.com/haunted-saas/graphql-api-gateway/internal/middleware"
	"go.uber.org/zap"

	billingv1 "github.com/haunted-saas/billing-service/proto/billing/v1"
	featureflagsv1 "github.com/haunted-saas/feature-flags-service/proto/featureflags/v1"
	llmv1 "github.com/haunted-saas/llm-gateway-service/proto/llm/v1"
	userauthv1 "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
)

// Query resolver
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// Subscription resolver (for GraphQL subscriptions/real-time updates)
// TODO: Implement GraphQL subscriptions when needed
func (r *Resolver) Subscription() generated.SubscriptionResolver {
	return nil
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

	resp, err := r.clients.Billing.GetSubscription(ctx, &billingv1.GetSubscriptionRequest{
		TeamId: userID, // Fixed: use GetSubscription with team_id
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
		TeamId: id, // Fixed: field is team_id not SubscriptionId
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

	resp, err := r.clients.Billing.CreateCustomerPortalSession(ctx, &billingv1.CreateCustomerPortalSessionRequest{
		TeamId:    userID, // Fixed: field is team_id and method is CreateCustomerPortalSession
		ReturnUrl: "http://localhost:3000/dashboard",
	})
	if err != nil {
		return "", errors.ConvertGRPCError(err)
	}

	return resp.PortalUrl, nil // Fixed: field is portal_url
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
		prompts[i] = convertPromptInfo(p) // Fixed: use convertPromptInfo for PromptInfo type
	}

	return prompts, nil
}

func (r *queryResolver) PromptDetails(ctx context.Context, name string) (*generated.PromptMetadata, error) {
	if err := middleware.RequireAuth(ctx); err != nil {
		return nil, err
	}

	resp, err := r.clients.LLMGateway.GetPromptMetadata(ctx, &llmv1.GetPromptMetadataRequest{
		PromptPath: name, // Fixed: field is prompt_path
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return convertPromptMetadata(resp), nil // Fixed: no Metadata field, pass resp directly
}

func (r *queryResolver) MyLLMUsage(ctx context.Context) (*generated.LLMUsageStats, error) {
	_, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.clients.LLMGateway.GetUsageStats(ctx, &llmv1.GetUsageStatsRequest{
		TimeRange: "week", // Fixed: UserId field doesn't exist
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	// CallsByModelJson doesn't exist in proto
	var callsByModel map[string]interface{}

	return &generated.LLMUsageStats{
		TotalCalls:   int(resp.TotalRequests), // Fixed: field is total_requests
		TotalTokens:  int(resp.TotalTokens),
		TotalCost:    0.0, // TotalCost doesn't exist in proto
		CallsByModel: callsByModel,
	}, nil
}

// ============================================================================
// NOTIFICATIONS QUERIES
// ============================================================================

func (r *queryResolver) NotificationToken(ctx context.Context) (*generated.NotificationToken, error) {
	// GenerateConnectionToken RPC doesn't exist in proto yet
	return nil, errors.NewBadRequestError("NotificationToken not implemented - GenerateConnectionToken RPC missing")
}

func (r *queryResolver) MyNotificationPreferences(ctx context.Context) (*generated.NotificationPreferences, error) {
	// GetPreferences RPC doesn't exist in proto yet
	return nil, errors.NewBadRequestError("MyNotificationPreferences not implemented - GetPreferences RPC missing")
}

// ============================================================================
// ANALYTICS QUERIES
// ============================================================================

func (r *queryResolver) MyAnalytics(ctx context.Context, startDate *time.Time, endDate *time.Time) (*generated.AnalyticsSummary, error) {
	// GetAnalyticsSummary RPC doesn't exist in proto yet
	return nil, errors.NewBadRequestError("MyAnalytics not implemented - GetAnalyticsSummary RPC missing")
}


