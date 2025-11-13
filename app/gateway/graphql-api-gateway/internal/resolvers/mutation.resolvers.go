package resolvers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/haunted-saas/graphql-api-gateway/internal/errors"
	"github.com/haunted-saas/graphql-api-gateway/internal/generated"
	"github.com/haunted-saas/graphql-api-gateway/internal/middleware"
	"go.uber.org/zap"

	analyticsv1 "github.com/haunted-saas/analytics-service/proto/analytics/v1"
	billingv1 "github.com/haunted-saas/billing-service/proto/billing/v1"
	llmv1 "github.com/haunted-saas/llm-gateway-service/proto/llm/v1"
	notificationsv1 "github.com/haunted-saas/notifications-service/proto/notifications/v1"
	userauthv1 "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
)

// Mutation resolver
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

// ============================================================================
// AUTHENTICATION MUTATIONS
// ============================================================================

func (r *mutationResolver) Register(ctx context.Context, input generated.RegisterInput) (*generated.AuthPayload, error) {
	resp, err := r.clients.UserAuth.Register(ctx, &userauthv1.RegisterRequest{
		Email:    input.Email,
		Password: input.Password,
		Name:     stringPtrToString(input.Name),
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	// RegisterResponse only returns User, need to login to get token
	loginResp, err := r.clients.UserAuth.Login(ctx, &userauthv1.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return &generated.AuthPayload{
		Token:        loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
		User:         convertUser(resp.User),
		ExpiresAt:    loginResp.ExpiresAt.AsTime(),
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input generated.LoginInput) (*generated.AuthPayload, error) {
	resp, err := r.clients.UserAuth.Login(ctx, &userauthv1.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return &generated.AuthPayload{
		Token:        resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User:         convertUser(resp.User),
		ExpiresAt:    resp.ExpiresAt.AsTime(),
	}, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	token := middleware.GetToken(ctx)

	_, err := r.clients.UserAuth.Logout(ctx, &userauthv1.LogoutRequest{
		SessionToken: token,
		AllDevices:   false,
	})
	if err != nil {
		return false, errors.ConvertGRPCError(err)
	}

	return true, nil
}

func (r *mutationResolver) RequestPasswordReset(ctx context.Context, email string) (bool, error) {
	_, err := r.clients.UserAuth.RequestPasswordReset(ctx, &userauthv1.PasswordResetRequest{
		Email: email,
	})
	if err != nil {
		// Don't expose whether email exists - always return true
		r.logger.Warn("password reset request failed", zap.Error(err))
	}

	return true, nil
}

func (r *mutationResolver) ResetPassword(ctx context.Context, token string, newPassword string) (bool, error) {
	_, err := r.clients.UserAuth.ResetPassword(ctx, &userauthv1.ResetPasswordRequest{
		Token:       token,
		NewPassword: newPassword,
	})
	if err != nil {
		return false, errors.ConvertGRPCError(err)
	}

	return true, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, currentPassword string, newPassword string) (bool, error) {
	// TODO: ChangePassword RPC not implemented in user-auth-service proto yet
	return false, errors.NewBadRequestError("ChangePassword not implemented yet")
}

func (r *mutationResolver) UpdateProfile(ctx context.Context, input generated.UpdateProfileInput) (*generated.User, error) {
	// TODO: UpdateUser RPC not implemented in user-auth-service proto yet
	return nil, errors.NewBadRequestError("UpdateProfile not implemented yet")
}

// ============================================================================
// RBAC MUTATIONS
// ============================================================================

func (r *mutationResolver) AssignRole(ctx context.Context, userID string, roleID string) (*generated.User, error) {
	if err := middleware.RequireRole(ctx, "admin"); err != nil {
		return nil, err
	}

	_, err := r.clients.UserAuth.AssignRoleToUser(ctx, &userauthv1.AssignRoleRequest{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	// TODO: GetUser RPC not implemented - return placeholder
	return &generated.User{
		ID:    userID,
		Email: "user@example.com",
		Name:  stringToPtr("User"),
	}, nil
}

func (r *mutationResolver) RemoveRole(ctx context.Context, userID string, roleID string) (*generated.User, error) {
	if err := middleware.RequireRole(ctx, "admin"); err != nil {
		return nil, err
	}

	_, err := r.clients.UserAuth.RevokeRoleFromUser(ctx, &userauthv1.RevokeRoleRequest{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	// TODO: GetUser RPC not implemented - return placeholder
	return &generated.User{
		ID:    userID,
		Email: "user@example.com",
		Name:  stringToPtr("User"),
	}, nil
}

func (r *mutationResolver) CreateRole(ctx context.Context, input generated.CreateRoleInput) (*generated.Role, error) {
	if err := middleware.RequireRole(ctx, "admin"); err != nil {
		return nil, err
	}

	resp, err := r.clients.UserAuth.CreateRole(ctx, &userauthv1.CreateRoleRequest{
		Name:          input.Name,
		Description:   stringPtrToString(input.Description),
		PermissionIds: input.Permissions,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return convertRole(resp), nil
}



// ============================================================================
// BILLING MUTATIONS
// ============================================================================

func (r *mutationResolver) CreateSubscriptionCheckout(ctx context.Context, planID string) (*generated.CheckoutPayload, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.clients.Billing.CreateCheckoutSession(ctx, &billingv1.CreateCheckoutSessionRequest{
		TeamId:     userID, // Fixed: field is team_id not UserId
		PlanId:     planID,
		SuccessUrl: "http://localhost:3000/success",
		CancelUrl:  "http://localhost:3000/cancel",
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return &generated.CheckoutPayload{
		SessionID: resp.SessionId,
		URL:       resp.CheckoutUrl, // Fixed: field is checkout_url not Url
	}, nil
}

func (r *mutationResolver) CancelSubscription(ctx context.Context) (*generated.Subscription, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.clients.Billing.CancelSubscription(ctx, &billingv1.CancelSubscriptionRequest{
		TeamId:           userID, // Fixed: field is team_id not UserId
		RequestingUserId: userID,
		Immediate:        false,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return convertSubscription(resp.Subscription), nil
}

func (r *mutationResolver) UpdateSubscription(ctx context.Context, planID string) (*generated.Subscription, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.clients.Billing.UpdateSubscription(ctx, &billingv1.UpdateSubscriptionRequest{
		TeamId:           userID, // Fixed: field is team_id not UserId
		NewPlanId:        planID,
		RequestingUserId: userID,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return convertSubscription(resp.Subscription), nil
}

// ============================================================================
// LLM GATEWAY MUTATIONS
// ============================================================================

func (r *mutationResolver) CallPrompt(ctx context.Context, name string, variables map[string]interface{}) (*generated.PromptResponse, error) {
	variablesJSON := "{}"
	if variables != nil {
		jsonBytes, err := json.Marshal(variables)
		if err != nil {
			return nil, errors.NewBadRequestError("invalid variables")
		}
		variablesJSON = string(jsonBytes)
	}

	resp, err := r.clients.LLMGateway.CallPrompt(ctx, &llmv1.CallPromptRequest{
		PromptPath:    name, // Fixed: field is prompt_path not PromptName
		VariablesJson: variablesJSON,
	})
	if err != nil {
		return nil, errors.ConvertGRPCError(err)
	}

	return &generated.PromptResponse{
		Content:      resp.ResponseText, // Fixed: field is response_text
		Model:        resp.ModelUsed,    // Fixed: field is model_used
		TokensUsed:   int(resp.TokenUsage.TotalTokens), // Fixed: nested in token_usage
		Cost:         0.0, // Cost not in proto
		FinishReason: "", // FinishReason not in proto
	}, nil
}

func (r *mutationResolver) CallLlm(ctx context.Context, input generated.LLMCallInput) (*generated.LLMResponse, error) {
	// CallLLM RPC doesn't exist in proto - return not implemented
	return nil, errors.NewBadRequestError("CallLLM not implemented - proto RPC missing")
}

// ============================================================================
// NOTIFICATIONS MUTATIONS
// ============================================================================

func (r *mutationResolver) SendNotification(ctx context.Context, input generated.SendNotificationInput) (bool, error) {
	// SendNotification RPC doesn't exist in proto - use SendToUser instead
	if err := middleware.RequireAuth(ctx); err != nil {
		return false, err
	}

	dataJSON := "{}"
	if input.Data != nil {
		jsonBytes, err := json.Marshal(input.Data)
		if err != nil {
			return false, errors.NewBadRequestError("invalid data")
		}
		dataJSON = string(jsonBytes)
	}

	_, err := r.clients.Notifications.SendToUser(ctx, &notificationsv1.SendToUserRequest{
		UserId:      input.UserID,
		EventType:   input.Type,
		PayloadJson: dataJSON,
	})
	if err != nil {
		return false, errors.ConvertGRPCError(err)
	}

	return true, nil
}

func (r *mutationResolver) UpdateNotificationPreferences(ctx context.Context, input generated.NotificationPreferencesInput) (*generated.NotificationPreferences, error) {
	// UpdatePreferences RPC doesn't exist in proto yet
	return nil, errors.NewBadRequestError("UpdateNotificationPreferences not implemented - proto RPC missing")
}

func (r *mutationResolver) MarkNotificationRead(ctx context.Context, notificationID string) (bool, error) {
	// MarkAsRead RPC doesn't exist in proto yet
	return false, errors.NewBadRequestError("MarkNotificationRead not implemented - proto RPC missing")
}

// ============================================================================
// ANALYTICS MUTATIONS
// ============================================================================

func (r *mutationResolver) TrackEvent(ctx context.Context, input generated.TrackEventInput) (bool, error) {
	userID, _ := middleware.GetUserID(ctx)

	// Convert properties map to proto PropertyValue map
	properties := make(map[string]*analyticsv1.PropertyValue)
	if input.Properties != nil {
		for key, val := range input.Properties {
			pv := &analyticsv1.PropertyValue{}
			switch v := val.(type) {
			case string:
				pv.Value = &analyticsv1.PropertyValue_StringValue{StringValue: v}
			case float64:
				pv.Value = &analyticsv1.PropertyValue_NumberValue{NumberValue: v}
			case bool:
				pv.Value = &analyticsv1.PropertyValue_BoolValue{BoolValue: v}
			}
			properties[key] = pv
		}
	}

	_, err := r.clients.Analytics.TrackEvent(ctx, &analyticsv1.TrackEventRequest{
		EventName:  input.EventName,
		UserId:     userID,
		Properties: properties, // Fixed: use map not JSON
		Timestamp:  time.Now().Unix(), // Fixed: int64 not timestamppb
	})
	if err != nil {
		return false, errors.ConvertGRPCError(err)
	}

	return true, nil
}

func (r *mutationResolver) IdentifyUser(ctx context.Context, properties map[string]interface{}) (bool, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return false, err
	}

	// Convert properties map to proto PropertyValue map
	protoProps := make(map[string]*analyticsv1.PropertyValue)
	if properties != nil {
		for key, val := range properties {
			pv := &analyticsv1.PropertyValue{}
			switch v := val.(type) {
			case string:
				pv.Value = &analyticsv1.PropertyValue_StringValue{StringValue: v}
			case float64:
				pv.Value = &analyticsv1.PropertyValue_NumberValue{NumberValue: v}
			case bool:
				pv.Value = &analyticsv1.PropertyValue_BoolValue{BoolValue: v}
			}
			protoProps[key] = pv
		}
	}

	_, err = r.clients.Analytics.IdentifyUser(ctx, &analyticsv1.IdentifyUserRequest{
		UserId:     userID,
		Properties: protoProps, // Fixed: use map not JSON
	})
	if err != nil {
		return false, errors.ConvertGRPCError(err)
	}

	return true, nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
