# Analytics Service - Integration Examples

## Table of Contents
1. [Backend Service Integration](#backend-service-integration)
2. [GraphQL Gateway Integration](#graphql-gateway-integration)
3. [Frontend Integration](#frontend-integration)
4. [Custom Event Properties](#custom-event-properties)
5. [Multi-Provider Setup](#multi-provider-setup)

---

## Backend Service Integration

### From User Auth Service

Track user registration and login events:

```go
// In user-auth-service/internal/service/auth_service.go

import (
    pb "github.com/haunted-saas/analytics-service/proto/analytics/v1"
    "google.golang.org/grpc"
)

type AuthService struct {
    analyticsClient pb.AnalyticsServiceClient
    // ... other fields
}

func NewAuthService(analyticsAddr string) (*AuthService, error) {
    // Connect to analytics service
    conn, err := grpc.Dial(analyticsAddr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    return &AuthService{
        analyticsClient: pb.NewAnalyticsServiceClient(conn),
    }, nil
}

// Track user registration
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
    // ... create user logic ...
    
    // Track registration event
    go func() {
        s.analyticsClient.TrackEvent(context.Background(), &pb.TrackEventRequest{
            EventName: "user_registered",
            UserId:    user.ID,
            Properties: map[string]*pb.PropertyValue{
                "email": {Value: &pb.PropertyValue_StringValue{StringValue: user.Email}},
                "signup_method": {Value: &pb.PropertyValue_StringValue{StringValue: "email"}},
                "plan": {Value: &pb.PropertyValue_StringValue{StringValue: "free"}},
            },
        })
    }()
    
    return user, nil
}

// Track user login
func (s *AuthService) Login(ctx context.Context, email, password string) (*Session, error) {
    // ... authentication logic ...
    
    // Track login event
    go func() {
        s.analyticsClient.TrackEvent(context.Background(), &pb.TrackEventRequest{
            EventName: "user_logged_in",
            UserId:    user.ID,
            Properties: map[string]*pb.PropertyValue{
                "login_method": {Value: &pb.PropertyValue_StringValue{StringValue: "password"}},
                "ip_address": {Value: &pb.PropertyValue_StringValue{StringValue: ipAddr}},
            },
        })
    }()
    
    return session, nil
}

// Identify user with properties
func (s *AuthService) UpdateUserProfile(ctx context.Context, userID string, profile *Profile) error {
    // ... update profile logic ...
    
    // Identify user with updated properties
    go func() {
        s.analyticsClient.IdentifyUser(context.Background(), &pb.IdentifyUserRequest{
            UserId: userID,
            Properties: map[string]*pb.PropertyValue{
                "name": {Value: &pb.PropertyValue_StringValue{StringValue: profile.Name}},
                "email": {Value: &pb.PropertyValue_StringValue{StringValue: profile.Email}},
                "company": {Value: &pb.PropertyValue_StringValue{StringValue: profile.Company}},
                "plan": {Value: &pb.PropertyValue_StringValue{StringValue: profile.Plan}},
            },
        })
    }()
    
    return nil
}
```

### From Billing Service

Track subscription events:

```go
// In billing-service/internal/webhook_handler.go

func (h *WebhookHandler) handleCheckoutSessionCompleted(ctx context.Context, event stripe.Event) error {
    // ... subscription provisioning logic ...
    
    // Track subscription created event
    go func() {
        h.analyticsClient.TrackEvent(context.Background(), &pb.TrackEventRequest{
            EventName: "subscription_created",
            UserId:    userID,
            Properties: map[string]*pb.PropertyValue{
                "plan_id": {Value: &pb.PropertyValue_StringValue{StringValue: planID}},
                "plan_name": {Value: &pb.PropertyValue_StringValue{StringValue: planName}},
                "amount": {Value: &pb.PropertyValue_NumberValue{NumberValue: float64(amount)}},
                "currency": {Value: &pb.PropertyValue_StringValue{StringValue: "usd"}},
                "trial_days": {Value: &pb.PropertyValue_NumberValue{NumberValue: float64(trialDays)}},
            },
        })
    }()
    
    return nil
}

func (h *WebhookHandler) handleSubscriptionCanceled(ctx context.Context, event stripe.Event) error {
    // ... cancellation logic ...
    
    // Track subscription canceled event
    go func() {
        h.analyticsClient.TrackEvent(context.Background(), &pb.TrackEventRequest{
            EventName: "subscription_canceled",
            UserId:    userID,
            Properties: map[string]*pb.PropertyValue{
                "plan_id": {Value: &pb.PropertyValue_StringValue{StringValue: planID}},
                "cancellation_reason": {Value: &pb.PropertyValue_StringValue{StringValue: reason}},
                "days_active": {Value: &pb.PropertyValue_NumberValue{NumberValue: float64(daysActive)}},
            },
        })
    }()
    
    return nil
}
```

### From LLM Gateway Service

Track AI usage:

```go
// In llm-gateway-service/internal/grpc_handlers.go

func (s *LLMGatewayServer) CallPrompt(ctx context.Context, req *pb.CallPromptRequest) (*pb.CallPromptResponse, error) {
    // ... LLM call logic ...
    
    // Track LLM usage
    go func() {
        s.analyticsClient.TrackEvent(context.Background(), &pb.TrackEventRequest{
            EventName: "llm_prompt_executed",
            UserId:    req.CallingService, // Or extract from context
            Properties: map[string]*pb.PropertyValue{
                "prompt_path": {Value: &pb.PropertyValue_StringValue{StringValue: req.PromptPath}},
                "model": {Value: &pb.PropertyValue_StringValue{StringValue: llmResp.Model}},
                "prompt_tokens": {Value: &pb.PropertyValue_NumberValue{NumberValue: float64(llmResp.TokenUsage.PromptTokens)}},
                "completion_tokens": {Value: &pb.PropertyValue_NumberValue{NumberValue: float64(llmResp.TokenUsage.CompletionTokens)}},
                "total_tokens": {Value: &pb.PropertyValue_NumberValue{NumberValue: float64(llmResp.TokenUsage.TotalTokens)}},
                "response_time_ms": {Value: &pb.PropertyValue_NumberValue{NumberValue: float64(responseTime.Milliseconds())}},
            },
        })
    }()
    
    return response, nil
}
```

---

## GraphQL Gateway Integration

### Setup Analytics Client

```go
// In gateway/internal/analytics/client.go

package analytics

import (
    pb "github.com/haunted-saas/analytics-service/proto/analytics/v1"
    "google.golang.org/grpc"
)

type Client struct {
    client pb.AnalyticsServiceClient
}

func NewClient(addr string) (*Client, error) {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    return &Client{
        client: pb.NewAnalyticsServiceClient(conn),
    }, nil
}

func (c *Client) TrackEvent(ctx context.Context, userID, eventName string, properties map[string]interface{}) error {
    // Convert properties to proto format
    protoProps := make(map[string]*pb.PropertyValue)
    for key, value := range properties {
        protoProps[key] = convertToPropertyValue(value)
    }
    
    _, err := c.client.TrackEvent(ctx, &pb.TrackEventRequest{
        EventName:  eventName,
        UserId:     userID,
        Properties: protoProps,
    })
    
    return err
}

func convertToPropertyValue(value interface{}) *pb.PropertyValue {
    switch v := value.(type) {
    case string:
        return &pb.PropertyValue{Value: &pb.PropertyValue_StringValue{StringValue: v}}
    case float64:
        return &pb.PropertyValue{Value: &pb.PropertyValue_NumberValue{NumberValue: v}}
    case int:
        return &pb.PropertyValue{Value: &pb.PropertyValue_NumberValue{NumberValue: float64(v)}}
    case bool:
        return &pb.PropertyValue{Value: &pb.PropertyValue_BoolValue{BoolValue: v}}
    default:
        return &pb.PropertyValue{Value: &pb.PropertyValue_StringValue{StringValue: fmt.Sprintf("%v", v)}}
    }
}
```

### GraphQL Mutations

```go
// In gateway/internal/resolvers/analytics.go

type AnalyticsResolver struct {
    analyticsClient *analytics.Client
}

// Track event mutation
func (r *AnalyticsResolver) TrackEvent(ctx context.Context, args struct {
    EventName  string
    Properties map[string]interface{}
}) (*bool, error) {
    // Get user ID from context
    userID := getUserIDFromContext(ctx)
    
    // Track event
    err := r.analyticsClient.TrackEvent(ctx, userID, args.EventName, args.Properties)
    if err != nil {
        return nil, err
    }
    
    success := true
    return &success, nil
}

// Identify user mutation
func (r *AnalyticsResolver) IdentifyUser(ctx context.Context, args struct {
    Properties map[string]interface{}
}) (*bool, error) {
    userID := getUserIDFromContext(ctx)
    
    // Convert properties
    protoProps := make(map[string]*pb.PropertyValue)
    for key, value := range args.Properties {
        protoProps[key] = convertToPropertyValue(value)
    }
    
    _, err := r.analyticsClient.client.IdentifyUser(ctx, &pb.IdentifyUserRequest{
        UserId:     userID,
        Properties: protoProps,
    })
    
    if err != nil {
        return nil, err
    }
    
    success := true
    return &success, nil
}
```

### GraphQL Schema

```graphql
# In gateway/schema.graphql

type Mutation {
  # Track an analytics event
  trackEvent(
    eventName: String!
    properties: JSON
  ): Boolean!
  
  # Identify user with properties
  identifyUser(
    properties: JSON!
  ): Boolean!
}

# Example usage:
# mutation {
#   trackEvent(
#     eventName: "button_clicked"
#     properties: {
#       button_name: "signup"
#       page: "/landing"
#     }
#   )
# }
```

---

## Frontend Integration

### React/Next.js Integration

```typescript
// lib/analytics.ts

import { gql } from '@apollo/client';
import { apolloClient } from './apollo-client';

const TRACK_EVENT = gql`
  mutation TrackEvent($eventName: String!, $properties: JSON) {
    trackEvent(eventName: $eventName, properties: $properties)
  }
`;

const IDENTIFY_USER = gql`
  mutation IdentifyUser($properties: JSON!) {
    identifyUser(properties: $properties)
  }
`;

export const analytics = {
  // Track an event
  track: async (eventName: string, properties?: Record<string, any>) => {
    try {
      await apolloClient.mutate({
        mutation: TRACK_EVENT,
        variables: { eventName, properties },
      });
    } catch (error) {
      console.error('Analytics track error:', error);
    }
  },
  
  // Identify user
  identify: async (properties: Record<string, any>) => {
    try {
      await apolloClient.mutate({
        mutation: IDENTIFY_USER,
        variables: { properties },
      });
    } catch (error) {
      console.error('Analytics identify error:', error);
    }
  },
  
  // Track page view
  page: async (pageName: string, properties?: Record<string, any>) => {
    await analytics.track('page_viewed', {
      page_name: pageName,
      ...properties,
    });
  },
};
```

### Usage in Components

```typescript
// components/SignupForm.tsx

import { analytics } from '@/lib/analytics';

export function SignupForm() {
  const handleSubmit = async (data: SignupData) => {
    try {
      // Create account
      const user = await createAccount(data);
      
      // Track signup
      await analytics.track('user_signed_up', {
        signup_method: 'email',
        plan: 'free',
      });
      
      // Identify user
      await analytics.identify({
        email: user.email,
        name: user.name,
        created_at: new Date().toISOString(),
      });
      
      router.push('/dashboard');
    } catch (error) {
      // Track error
      await analytics.track('signup_failed', {
        error: error.message,
      });
    }
  };
  
  return <form onSubmit={handleSubmit}>...</form>;
}
```

### Track User Actions

```typescript
// components/Dashboard.tsx

import { useEffect } from 'react';
import { analytics } from '@/lib/analytics';

export function Dashboard() {
  useEffect(() => {
    // Track page view
    analytics.page('dashboard');
  }, []);
  
  const handleFeatureClick = (featureName: string) => {
    // Track feature usage
    analytics.track('feature_clicked', {
      feature_name: featureName,
      page: 'dashboard',
    });
  };
  
  return (
    <div>
      <button onClick={() => handleFeatureClick('export_data')}>
        Export Data
      </button>
    </div>
  );
}
```

---

## Custom Event Properties

### Adding Custom Properties to Events

```go
// Extend the Event struct with custom fields
// In internal/types.go (optional - properties map is already flexible)

// In internal/grpc_handlers.go - add custom properties
func (s *AnalyticsServer) TrackEvent(ctx context.Context, req *pb.TrackEventRequest) (*pb.TrackEventResponse, error) {
    // Parse custom properties
    properties := make(map[string]interface{})
    for key, propValue := range req.Properties {
        properties[key] = convertPropertyValue(propValue)
    }
    
    // Add automatic properties
    properties["$timestamp"] = time.Now().Unix()
    properties["$service"] = "analytics-service"
    properties["$version"] = "1.0.0"
    
    // Add IP address if available (from gRPC metadata)
    if md, ok := metadata.FromIncomingContext(ctx); ok {
        if ips := md.Get("x-forwarded-for"); len(ips) > 0 {
            properties["$ip"] = ips[0]
        }
    }
    
    // Create event with enriched properties
    event := Event{
        ID:         eventID,
        EventName:  req.EventName,
        UserID:     req.UserId,
        Properties: properties,
        Timestamp:  time.Now(),
    }
    
    s.queue.Add(event)
    return &pb.TrackEventResponse{Success: true, EventId: eventID}, nil
}
```

### Custom Property Examples

```go
// E-commerce tracking
analytics.TrackEvent(ctx, &pb.TrackEventRequest{
    EventName: "product_purchased",
    UserId:    userID,
    Properties: map[string]*pb.PropertyValue{
        "product_id": {Value: &pb.PropertyValue_StringValue{StringValue: "prod_123"}},
        "product_name": {Value: &pb.PropertyValue_StringValue{StringValue: "Pro Plan"}},
        "price": {Value: &pb.PropertyValue_NumberValue{NumberValue: 29.99}},
        "currency": {Value: &pb.PropertyValue_StringValue{StringValue: "USD"}},
        "quantity": {Value: &pb.PropertyValue_NumberValue{NumberValue: 1}},
        "category": {Value: &pb.PropertyValue_StringValue{StringValue: "subscription"}},
    },
})

// Feature usage tracking
analytics.TrackEvent(ctx, &pb.TrackEventRequest{
    EventName: "feature_used",
    UserId:    userID,
    Properties: map[string]*pb.PropertyValue{
        "feature_name": {Value: &pb.PropertyValue_StringValue{StringValue: "ai_assistant"}},
        "usage_count": {Value: &pb.PropertyValue_NumberValue{NumberValue: 5}},
        "session_duration": {Value: &pb.PropertyValue_NumberValue{NumberValue: 120}},
        "success": {Value: &pb.PropertyValue_BoolValue{BoolValue: true}},
    },
})

// Error tracking
analytics.TrackEvent(ctx, &pb.TrackEventRequest{
    EventName: "error_occurred",
    UserId:    userID,
    Properties: map[string]*pb.PropertyValue{
        "error_type": {Value: &pb.PropertyValue_StringValue{StringValue: "validation_error"}},
        "error_message": {Value: &pb.PropertyValue_StringValue{StringValue: "Invalid email format"}},
        "page": {Value: &pb.PropertyValue_StringValue{StringValue: "/signup"}},
        "severity": {Value: &pb.PropertyValue_StringValue{StringValue: "low"}},
    },
})
```

---

## Multi-Provider Setup

### Configure Multiple Providers

```bash
# In .env
ANALYTICS_PROVIDER=mixpanel,amplitude,segment
MIXPANEL_API_KEY=your-mixpanel-key
AMPLITUDE_API_KEY=your-amplitude-key
SEGMENT_WRITE_KEY=your-segment-key
```

### Update Configuration

```go
// In internal/config/config.go

type AnalyticsConfig struct {
    Providers           []string // "mixpanel", "amplitude", "segment"
    MixpanelAPIKey     string
    AmplitudeAPIKey    string
    SegmentWriteKey    string
    BatchSize          int
    FlushIntervalSec   int
    TestMode           bool
}

func Load() (*Config, error) {
    cfg := &Config{
        Analytics: AnalyticsConfig{
            Providers:        parseProviders(getEnv("ANALYTICS_PROVIDER", "mixpanel")),
            MixpanelAPIKey:   getEnv("MIXPANEL_API_KEY", ""),
            AmplitudeAPIKey:  getEnv("AMPLITUDE_API_KEY", ""),
            SegmentWriteKey:  getEnv("SEGMENT_WRITE_KEY", ""),
            // ...
        },
    }
    return cfg, nil
}

func parseProviders(providersStr string) []string {
    providers := strings.Split(providersStr, ",")
    result := make([]string, 0, len(providers))
    for _, p := range providers {
        trimmed := strings.TrimSpace(p)
        if trimmed != "" {
            result = append(result, trimmed)
        }
    }
    return result
}
```

### Multi-Provider Worker

```go
// In cmd/main.go

// Initialize all configured providers
var providers []internal.ExternalProvider

for _, providerName := range cfg.Analytics.Providers {
    switch providerName {
    case "mixpanel":
        if cfg.Analytics.MixpanelAPIKey != "" {
            provider := internal.NewMixpanelProvider(cfg.Analytics.MixpanelAPIKey, cfg.Analytics.TestMode, logger)
            providers = append(providers, provider)
            logger.Info("✓ Mixpanel provider initialized")
        }
    case "amplitude":
        if cfg.Analytics.AmplitudeAPIKey != "" {
            provider := internal.NewAmplitudeProvider(cfg.Analytics.AmplitudeAPIKey, cfg.Analytics.TestMode, logger)
            providers = append(providers, provider)
            logger.Info("✓ Amplitude provider initialized")
        }
    case "segment":
        if cfg.Analytics.SegmentWriteKey != "" {
            provider := internal.NewSegmentProvider(cfg.Analytics.SegmentWriteKey, cfg.Analytics.TestMode, logger)
            providers = append(providers, provider)
            logger.Info("✓ Segment provider initialized")
        }
    }
}

// Create multi-provider wrapper
multiProvider := internal.NewMultiProvider(providers, logger)

// Initialize batch worker with multi-provider
worker := internal.NewBatchWorker(queue, multiProvider, flushInterval, retryConfig, logger)
```

### Multi-Provider Implementation

```go
// In internal/multi_provider.go

type MultiProvider struct {
    providers []ExternalProvider
    logger    *zap.Logger
}

func NewMultiProvider(providers []ExternalProvider, logger *zap.Logger) *MultiProvider {
    return &MultiProvider{
        providers: providers,
        logger:    logger,
    }
}

func (p *MultiProvider) SendBatch(ctx context.Context, events []Event) error {
    var errors []error
    
    // Send to all providers in parallel
    var wg sync.WaitGroup
    errorsChan := make(chan error, len(p.providers))
    
    for _, provider := range p.providers {
        wg.Add(1)
        go func(prov ExternalProvider) {
            defer wg.Done()
            if err := prov.SendBatch(ctx, events); err != nil {
                p.logger.Error("provider send failed",
                    zap.String("provider", prov.GetName()),
                    zap.Error(err))
                errorsChan <- fmt.Errorf("%s: %w", prov.GetName(), err)
            } else {
                p.logger.Debug("provider send succeeded",
                    zap.String("provider", prov.GetName()))
            }
        }(provider)
    }
    
    wg.Wait()
    close(errorsChan)
    
    // Collect errors
    for err := range errorsChan {
        errors = append(errors, err)
    }
    
    // Return error if all providers failed
    if len(errors) == len(p.providers) {
        return fmt.Errorf("all providers failed: %v", errors)
    }
    
    // Partial success is OK
    return nil
}

func (p *MultiProvider) GetName() string {
    names := make([]string, len(p.providers))
    for i, provider := range p.providers {
        names[i] = provider.GetName()
    }
    return strings.Join(names, ",")
}
```

---

## Summary

The analytics service now supports:

✅ **Multiple Providers**: Mixpanel, Amplitude, Segment  
✅ **Backend Integration**: Easy gRPC client setup  
✅ **GraphQL Gateway**: Mutations for frontend tracking  
✅ **Frontend SDK**: React/Next.js integration  
✅ **Custom Properties**: Flexible event enrichment  
✅ **Multi-Provider**: Send to multiple services simultaneously  

All integrations are non-blocking and production-ready! 🎃
