package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userauthv1 "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
)

// Context keys for user information
type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	TeamIDKey   contextKey = "team_id"
	RolesKey    contextKey = "roles"
	TokenKey    contextKey = "token"
	IsAuthKey   contextKey = "is_authenticated"
)

// AuthMiddleware handles authentication for GraphQL requests
type AuthMiddleware struct {
	userAuthClient userauthv1.UserAuthServiceClient
	logger         *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(userAuthClient userauthv1.UserAuthServiceClient, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		userAuthClient: userAuthClient,
		logger:         logger,
	}
}

// Middleware returns the HTTP middleware function
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No token - mark as unauthenticated and continue
			ctx = context.WithValue(ctx, IsAuthKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Warn("invalid authorization header format")
			ctx = context.WithValue(ctx, IsAuthKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		token := parts[1]

		// Validate token with user-auth-service
		resp, err := m.userAuthClient.ValidateToken(ctx, &userauthv1.ValidateTokenRequest{
			Token: token,
		})

		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unauthenticated {
				m.logger.Debug("invalid token", zap.Error(err))
			} else {
				m.logger.Error("failed to validate token", zap.Error(err))
			}
			ctx = context.WithValue(ctx, IsAuthKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if !resp.Valid {
			m.logger.Debug("token validation failed")
			ctx = context.WithValue(ctx, IsAuthKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Token is valid - inject user information into context
		ctx = context.WithValue(ctx, IsAuthKey, true)
		ctx = context.WithValue(ctx, UserIDKey, resp.UserId)
		ctx = context.WithValue(ctx, TeamIDKey, resp.TeamId)
		ctx = context.WithValue(ctx, RolesKey, resp.Roles)
		ctx = context.WithValue(ctx, TokenKey, token)

		m.logger.Debug("user authenticated",
			zap.String("user_id", resp.UserId),
			zap.String("team_id", resp.TeamId),
			zap.Strings("roles", resp.Roles))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GraphQLAuthDirective enforces authentication on GraphQL operations
func GraphQLAuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	isAuth, ok := ctx.Value(IsAuthKey).(bool)
	if !ok || !isAuth {
		return nil, &gqlerror.Error{
			Message: "Unauthorized: authentication required",
			Extensions: map[string]interface{}{
				"code": "UNAUTHENTICATED",
			},
		}
	}

	return next(ctx)
}

// RequireAuth is a helper that can be called in resolvers to enforce authentication
func RequireAuth(ctx context.Context) error {
	isAuth, ok := ctx.Value(IsAuthKey).(bool)
	if !ok || !isAuth {
		return &gqlerror.Error{
			Message: "Unauthorized: authentication required",
			Extensions: map[string]interface{}{
				"code": "UNAUTHENTICATED",
			},
		}
	}
	return nil
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (string, error) {
	if err := RequireAuth(ctx); err != nil {
		return "", err
	}

	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", &gqlerror.Error{
			Message: "User ID not found in context",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_ERROR",
			},
		}
	}

	return userID, nil
}

// GetTeamID extracts team ID from context
func GetTeamID(ctx context.Context) string {
	teamID, _ := ctx.Value(TeamIDKey).(string)
	return teamID
}

// GetRoles extracts roles from context
func GetRoles(ctx context.Context) []string {
	roles, _ := ctx.Value(RolesKey).([]string)
	return roles
}

// GetToken extracts the JWT token from context
func GetToken(ctx context.Context) string {
	token, _ := ctx.Value(TokenKey).(string)
	return token
}

// IsAuthenticated checks if the request is authenticated
func IsAuthenticated(ctx context.Context) bool {
	isAuth, ok := ctx.Value(IsAuthKey).(bool)
	return ok && isAuth
}

// RequireRole checks if user has a specific role
func RequireRole(ctx context.Context, requiredRole string) error {
	if err := RequireAuth(ctx); err != nil {
		return err
	}

	roles := GetRoles(ctx)
	for _, role := range roles {
		if role == requiredRole {
			return nil
		}
	}

	return &gqlerror.Error{
		Message: "Forbidden: insufficient permissions",
		Extensions: map[string]interface{}{
			"code":          "FORBIDDEN",
			"required_role": requiredRole,
		},
	}
}

// RequireAnyRole checks if user has any of the specified roles
func RequireAnyRole(ctx context.Context, requiredRoles []string) error {
	if err := RequireAuth(ctx); err != nil {
		return err
	}

	roles := GetRoles(ctx)
	roleMap := make(map[string]bool)
	for _, role := range roles {
		roleMap[role] = true
	}

	for _, requiredRole := range requiredRoles {
		if roleMap[requiredRole] {
			return nil
		}
	}

	return &gqlerror.Error{
		Message: "Forbidden: insufficient permissions",
		Extensions: map[string]interface{}{
			"code":           "FORBIDDEN",
			"required_roles": requiredRoles,
		},
	}
}
