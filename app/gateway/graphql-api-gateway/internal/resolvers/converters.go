package resolvers

import (
	"time"

	"github.com/haunted-saas/graphql-api-gateway/internal/generated"

	billingv1 "github.com/haunted-saas/billing-service/proto/billing/v1"
	featureflagsv1 "github.com/haunted-saas/feature-flags-service/proto/featureflags/v1"
	llmv1 "github.com/haunted-saas/llm-gateway-service/proto/llm/v1"
	userauthv1 "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
)

// ============================================================================
// USER & AUTH CONVERTERS
// ============================================================================

func convertUser(u *userauthv1.User) *generated.User {
	if u == nil {
		return nil
	}

	roles := make([]*generated.Role, len(u.Roles))
	for i, role := range u.Roles {
		roles[i] = convertRole(role)
	}

	// Extract permission names from roles
	permissionSet := make(map[string]bool)
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			permissionSet[perm.Name] = true
		}
	}
	permissions := make([]string, 0, len(permissionSet))
	for perm := range permissionSet {
		permissions = append(permissions, perm)
	}

	return &generated.User{
		ID:          u.Id,
		Email:       u.Email,
		Name:        stringToPtr(u.Name),
		TeamID:      nil, // TeamId field doesn't exist in proto yet
		Roles:       roles,
		Permissions: permissions,
		CreatedAt:   u.CreatedAt.AsTime(),
		UpdatedAt:   u.UpdatedAt.AsTime(),
	}
}

func convertRole(r *userauthv1.Role) *generated.Role {
	if r == nil {
		return nil
	}

	// Convert Permission objects to permission name strings
	permissions := make([]string, len(r.Permissions))
	for i, perm := range r.Permissions {
		permissions[i] = perm.Name
	}

	return &generated.Role{
		ID:          r.Id,
		Name:        r.Name,
		Description: stringToPtr(r.Description),
		Permissions: permissions,
		IsSystem:    r.IsSystem,
		CreatedAt:   r.CreatedAt.AsTime(),
	}
}

// ============================================================================
// BILLING CONVERTERS
// ============================================================================

func convertPlan(p *billingv1.Plan) *generated.Plan {
	if p == nil {
		return nil
	}

	// Convert features map to array of strings
	features := make([]string, 0, len(p.Features))
	for key, value := range p.Features {
		features = append(features, key+": "+value)
	}

	return &generated.Plan{
		ID:            p.Id,
		Name:          p.Name,
		Description:   nil, // Description field doesn't exist in proto
		Price:         float64(p.PriceCents) / 100.0, // Convert cents to dollars
		Currency:      p.Currency,
		Interval:      p.BillingInterval,
		Features:      features,
		StripePriceID: p.StripePriceId,
		IsActive:      p.IsActive,
	}
}

func convertSubscription(s *billingv1.Subscription) *generated.Subscription {
	if s == nil {
		return nil
	}

	// TODO: Fix this once gqlgen regenerates with correct schema
	// The generated.Subscription type doesn't match the GraphQL schema yet
	return nil
	
	/* Commented out until gqlgen regenerates properly
	cancelAtPeriodEnd := false
	if s.CancelAt != nil {
		cancelAtPeriodEnd = true
	}

	return &generated.Subscription{
		ID:                    s.Id,
		UserID:                s.TeamId,
		PlanID:                s.PlanId,
		Status:                s.Status,
		CurrentPeriodStart:    s.CurrentPeriodStart.AsTime(),
		CurrentPeriodEnd:      s.CurrentPeriodEnd.AsTime(),
		CancelAtPeriodEnd:     cancelAtPeriodEnd,
		StripeSubscriptionID:  s.StripeSubscriptionId,
		CreatedAt:             s.CreatedAt.AsTime(),
		UpdatedAt:             s.UpdatedAt.AsTime(),
	}
	*/
}

// ============================================================================
// FEATURE FLAGS CONVERTERS
// ============================================================================

func convertFeature(f *featureflagsv1.Feature) *generated.Feature {
	if f == nil {
		return nil
	}

	createdAt, _ := time.Parse(time.RFC3339, f.CreatedAt)

	return &generated.Feature{
		Name:        f.Name,
		Description: f.Description,
		Enabled:     f.Enabled,
		CreatedAt:   createdAt,
	}
}

// ============================================================================
// LLM GATEWAY CONVERTERS
// ============================================================================

func convertPromptInfo(p *llmv1.PromptInfo) *generated.PromptMetadata {
	if p == nil {
		return nil
	}

	return &generated.PromptMetadata{
		Name:        p.Path,
		Description: "",
		Version:     "",
		Variables:   []string{},
		Model:       "",
		Temperature: 0.7,
	}
}

func convertPromptMetadata(resp *llmv1.GetPromptMetadataResponse) *generated.PromptMetadata {
	if resp == nil {
		return nil
	}

	return &generated.PromptMetadata{
		Name:        resp.PromptPath,
		Description: "",
		Version:     resp.LastModified,
		Variables:   resp.RequiredVariables,
		Model:       "",
		Temperature: 0.7,
	}
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
