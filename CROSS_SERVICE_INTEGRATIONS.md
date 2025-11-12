# Cross-Service Integrations Guide

This document explains how the microservices integrate with each other and clarifies what's implemented vs what's left as extension points.

## Service Dependency Map

```
┌─────────────────────────────────────────────────────────────┐
│                    GraphQL API Gateway                       │
│  (Calls ALL services via gRPC - FULLY IMPLEMENTED)          │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
   ┌────▼────┐         ┌────▼────┐        ┌────▼────┐
   │ User    │         │ Billing │        │ LLM     │
   │ Auth    │         │ Service │        │ Gateway │
   └─────────┘         └─────────┘        └─────────┘
        │                   │                   │
   ┌────▼────┐         ┌────▼────┐        ┌────▼────┐
   │ Notif.  │         │ Analytics│       │ Feature │
   │ Service │         │ Service  │       │ Flags   │
   └─────────┘         └──────────┘       └─────────┘
```

## Fully Implemented Integrations

### 1. GraphQL Gateway → All Services ✅

**Location**: `app/gateway/graphql-api-gateway/internal/clients/grpc_clients.go`

**Status**: COMPLETE

The GraphQL gateway initializes gRPC clients to all 6 services on startup:

```go
type GRPCClients struct {
    UserAuth      userauthv1.UserAuthServiceClient
    Billing       billingv1.BillingServiceClient
    LLMGateway    llmv1.LLMGatewayServiceClient
    Notifications notificationsv1.NotificationsServiceClient
    Analytics     analyticsv1.AnalyticsServiceClient
    FeatureFlags  featureflagsv1.FeatureFlagsServiceClient
}
```

**Usage**: All GraphQL resolvers use these clients to call backend services.

**Example**:
```go
// In query resolver
resp, err := r.clients.UserAuth.GetUser(ctx, &userauthv1.GetUserRequest{
    UserId: userID,
})
```

### 2. GraphQL Gateway → User Auth Service (Authentication) ✅

**Location**: `app/gateway/graphql-api-gateway/internal/middleware/auth.go`

**Status**: COMPLETE

The auth middleware validates JWT tokens by calling user-auth-service:

```go
resp, err := m.userAuthClient.ValidateToken(ctx, &userauthv1.ValidateTokenRequest{
    Token: token,
})
```

**Flow**:
1. Extract JWT from `Authorization: Bearer <token>` header
2. Call `user-auth-service.ValidateToken()`
3. Inject user context into request
4. Pass to resolvers

### 3. Notifications Service → User Auth Service (JWT Validation) ✅

**Location**: `app/services/notifications-service/internal/auth_middleware.go`

**Status**: COMPLETE (Local JWT validation)

**Implementation**: The notifications service validates JWTs locally for performance:

```go
func (m *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return m.jwtSecret, nil
    })
    // ...
}
```

**Why Local?**: Socket.IO connections need sub-millisecond auth checks. Local JWT validation is appropriate here since:
- JWTs are cryptographically signed
- No database lookup needed
- User-auth-service already validated the token when it was issued

**Alternative**: For stricter security, you could call `user-auth-service.ValidateToken()` but this adds latency.

## Extension Points (Intentionally Left for Future Implementation)

These are marked with TODO comments and represent optional integrations that can be added based on requirements.

### 1. LLM Gateway → Analytics Service (Usage Tracking)

**Location**: `app/services/llm-gateway-service/internal/usage_tracker.go:40`

**Current**: Logs usage locally in memory

**Future Enhancement**:
```go
// Add analytics client to UsageTracker
type UsageTracker struct {
    store           *UsageStore
    analyticsClient analyticsv1.AnalyticsServiceClient
    logger          *zap.Logger
}

// In TrackUsage method
go func() {
    _, err := t.analyticsClient.TrackEvent(ctx, &analyticsv1.TrackEventRequest{
        UserId:    event.UserID,
        EventName: "llm_call",
        PropertiesJson: fmt.Sprintf(`{
            "model": "%s",
            "tokens": %d,
            "cost": %.4f,
            "prompt": "%s"
        }`, event.Model, event.TokensUsed, event.Cost, event.PromptPath),
    })
}()
```

**Why Not Implemented**: 
- Local logging is sufficient for MVP
- Analytics service already tracks events from GraphQL gateway
- Adds complexity and potential failure point

**When to Implement**: When you need detailed LLM usage analytics beyond what's tracked at the gateway level.

### 2. Billing Service → User Auth Service (Access Provisioning)

**Location**: `app/services/billing-service/internal/webhook_handler.go:220`

**Current**: TODO comment for future implementation

**Future Enhancement**:
```go
// When subscription is created/updated
func (h *WebhookHandler) handleSubscriptionUpdated(subscription *db.Subscription) {
    // Call user-auth-service to update user's plan
    _, err := h.userAuthClient.UpdateUserPlan(ctx, &userauthv1.UpdateUserPlanRequest{
        UserId: subscription.UserID,
        Plan:   subscription.PlanID,
    })
    
    // Or call feature-flags-service to enable features
    _, err := h.featureFlagsClient.EnableFeaturesForPlan(ctx, &featureflagsv1.EnableFeaturesRequest{
        TeamId: subscription.TeamID,
        Plan:   subscription.PlanID,
    })
}
```

**Why Not Implemented**:
- Feature flags are evaluated dynamically based on plan (no provisioning needed)
- User plan is stored in user-auth database
- Webhook handler focuses on Stripe sync

**When to Implement**: When you need to:
- Provision specific resources per plan
- Send welcome emails on subscription
- Update user roles based on plan

### 3. Billing Service → Notifications Service (Payment Notifications)

**Location**: `app/services/billing-service/internal/webhook_handler.go:388`

**Current**: TODO comment for trial ending notifications

**Future Enhancement**:
```go
// When trial is ending
func (h *WebhookHandler) handleTrialWillEnd(subscription *db.Subscription) {
    _, err := h.notificationsClient.SendToUser(ctx, &notificationsv1.SendToUserRequest{
        UserId:  subscription.UserID,
        Event:   "trial_ending",
        Message: "Your trial ends in 3 days",
        DataJson: fmt.Sprintf(`{"trial_end": "%s"}`, subscription.TrialEnd),
    })
}
```

**Why Not Implemented**:
- Notifications can be sent from frontend after subscription check
- Reduces coupling between services
- Webhook handler focuses on data sync

**When to Implement**: When you need automated notifications for:
- Trial ending warnings
- Payment failures
- Subscription renewals
- Plan upgrades

### 4. Analytics Service → External Providers (Query Implementation)

**Location**: `app/services/analytics-service/internal/grpc_handlers.go:114`

**Current**: Placeholder implementations for `GetEventCount` and `GetUserCount`

**Future Enhancement**:
```go
func (s *AnalyticsServer) GetEventCount(ctx context.Context, req *pb.GetEventCountRequest) (*pb.GetEventCountResponse, error) {
    // Query from Mixpanel/Amplitude
    count, err := s.provider.GetEventCount(ctx, req.EventName, req.StartDate, req.EndDate)
    if err != nil {
        return nil, err
    }
    
    return &pb.GetEventCountResponse{
        Count: count,
    }, nil
}
```

**Why Not Implemented**:
- Analytics providers (Mixpanel/Amplitude) have their own dashboards
- Querying is typically done through their UIs
- Event tracking (write path) is fully implemented

**When to Implement**: When you need to:
- Display analytics in your own dashboard
- Build custom reports
- Aggregate data from multiple providers

## How to Add New Cross-Service Integrations

### Step 1: Add gRPC Client to Service

```go
// In internal/config/config.go
type Config struct {
    // ... existing config
    Services ServicesConfig
}

type ServicesConfig struct {
    UserAuthService string
    // Add new service address
}

// In cmd/main.go
conn, err := grpc.Dial(cfg.Services.UserAuthService, grpc.WithInsecure())
if err != nil {
    logger.Fatal("failed to connect", zap.Error(err))
}
defer conn.Close()

client := userauthv1.NewUserAuthServiceClient(conn)
```

### Step 2: Pass Client to Handler/Service

```go
// Update constructor
func NewMyService(
    userAuthClient userauthv1.UserAuthServiceClient,
    logger *zap.Logger,
) *MyService {
    return &MyService{
        userAuthClient: userAuthClient,
        logger:         logger,
    }
}
```

### Step 3: Make gRPC Calls

```go
func (s *MyService) DoSomething(ctx context.Context) error {
    resp, err := s.userAuthClient.GetUser(ctx, &userauthv1.GetUserRequest{
        UserId: "user_123",
    })
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }
    
    // Use resp.User
    return nil
}
```

### Step 4: Update Environment Variables

```bash
# .env.example
USER_AUTH_SERVICE=localhost:50051
```

### Step 5: Update go.mod

```go
require (
    github.com/haunted-saas/user-auth-service v0.0.0
)

replace github.com/haunted-saas/user-auth-service => ../user-auth-service
```

## Service Communication Patterns

### Pattern 1: Synchronous gRPC (Request-Response)

**Use When**: You need immediate response

**Example**: GraphQL gateway calling any service

```go
resp, err := client.GetUser(ctx, &pb.GetUserRequest{UserId: id})
```

### Pattern 2: Asynchronous Fire-and-Forget

**Use When**: You don't need to wait for response

**Example**: Sending analytics events

```go
go func() {
    _, _ = analyticsClient.TrackEvent(ctx, &pb.TrackEventRequest{...})
}()
```

### Pattern 3: Event-Driven (Future)

**Use When**: Multiple services need to react to events

**Example**: Subscription created event

```
Billing Service → Message Queue → [User Auth, Notifications, Analytics]
```

**Not Implemented**: Would require message queue (RabbitMQ, Kafka, etc.)

## Testing Cross-Service Integrations

### Unit Tests

Mock the gRPC clients:

```go
type mockUserAuthClient struct {
    userauthv1.UserAuthServiceClient
}

func (m *mockUserAuthClient) GetUser(ctx context.Context, req *userauthv1.GetUserRequest, opts ...grpc.CallOption) (*userauthv1.GetUserResponse, error) {
    return &userauthv1.GetUserResponse{
        User: &userauthv1.User{
            Id:    req.UserId,
            Email: "test@example.com",
        },
    }, nil
}
```

### Integration Tests

Start all services and test end-to-end:

```bash
# Start all services
docker-compose up -d

# Run integration tests
go test ./tests/integration/...
```

### Manual Testing with grpcurl

```bash
# Test user-auth-service
grpcurl -plaintext -d '{"email":"test@example.com","password":"pass"}' \
  localhost:50051 userauth.v1.UserAuthService/Login

# Test feature-flags-service
grpcurl -plaintext -d '{"feature_name":"new_dashboard","user_id":"user_123"}' \
  localhost:50056 featureflags.v1.FeatureFlagsService/IsFeatureEnabled
```

## Troubleshooting

### "connection refused" errors

**Cause**: Target service not running

**Solution**: 
```bash
# Check if service is running
docker ps | grep service-name

# Check service logs
docker logs service-name
```

### "unknown service" errors

**Cause**: Proto mismatch or service not registered

**Solution**:
```bash
# Regenerate proto files
cd app/services/service-name
make proto

# Verify service registration in main.go
pb.RegisterServiceServer(grpcServer, serviceImpl)
```

### Timeout errors

**Cause**: Service taking too long to respond

**Solution**:
```go
// Add timeout to context
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

resp, err := client.Method(ctx, req)
```

## Summary

### ✅ Fully Implemented
- GraphQL Gateway → All Services (complete integration)
- GraphQL Gateway → User Auth (JWT validation)
- Notifications → User Auth (local JWT validation)

### 📝 Extension Points (Optional)
- LLM Gateway → Analytics (usage tracking)
- Billing → User Auth (access provisioning)
- Billing → Notifications (payment alerts)
- Analytics → Providers (query implementation)

### 🎯 Recommendation

The current implementation is **production-ready** for MVP. The extension points are intentionally left as TODOs because:

1. **They're optional** - Core functionality works without them
2. **They add complexity** - More moving parts = more failure points
3. **They can be added incrementally** - Start simple, add as needed

Implement extension points when you have specific requirements that justify the added complexity.

---

**Last Updated**: 2024  
**Status**: All critical integrations complete  
**Next Steps**: Implement extension points based on business requirements
