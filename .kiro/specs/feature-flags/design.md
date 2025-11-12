# Feature Flag Service Design Document

## Overview

The Feature Flag Service is a Go-based microservice that provides a self-hosted abstraction layer over Unleash for managing feature flags within the HAUNTED SAAS SKELETON platform. It exposes a simple gRPC interface that allows other microservices to check if features are enabled for specific users or contexts. The service handles feature flag evaluation, caching, and integration with the Analytics Service for tracking flag usage.

The service follows a layered architecture with the Unleash SDK at its core, wrapped by business logic that adds caching, analytics, and permission checks. It uses Redis for caching evaluations and integrates with the User Auth Service for RBAC enforcement.

## Architecture

### High-Level Architecture

```mermaid
graph TB
    Client[Next.js Client] -->|GraphQL| Gateway[GraphQL Gateway]
    Services[Other Microservices] -->|gRPC| FeatureFlags[Feature Flag Service]
    Gateway -->|gRPC| FeatureFlags
    FeatureFlags -->|Evaluate| Unleash[Unleash Server]
    FeatureFlags -->|Cache| Redis[(Redis)]
    FeatureFlags -->|Track| Analytics[Analytics Service]
    FeatureFlags -->|Verify| Auth[User Auth Service]
    Unleash -->|Store| UnleashDB[(Unleash PostgreSQL)]
    Admin[Admin User] -->|Manage| UnleashUI[Unleash Admin UI]
    UnleashUI --> Unleash
```

### Service Structure

```
feature-flag-service/
├── cmd/
│   └── server/
│       └── main.go                 # Service entry point
├── internal/
│   ├── api/
│   │   └── grpc/
│   │       ├── handler.go          # gRPC handler implementation
│   │       └── middleware.go       # Logging, tracing
│   ├── service/
│   │   ├── featureflags.go         # Business logic
│   │   ├── unleash.go              # Unleash SDK wrapper
│   │   └── permissions.go          # RBAC integration
│   ├── repository/
│   │   └── redis/
│   │       └── cache.go            # Caching layer
│   ├── models/
│   │   ├── context.go              # Evaluation context model
│   │   └── flag.go                 # Feature flag model
│   ├── analytics/
│   │   └── tracker.go              # Analytics integration
│   └── config/
│       └── config.go               # Configuration management
├── proto/
│   └── featureflags/
│       └── v1/
│           └── featureflags.proto  # gRPC service definition
├── Dockerfile
└── go.mod
```

## Components and Interfaces

### gRPC Service Definition

```protobuf
syntax = "proto3";

package featureflags.v1;

option go_package = "github.com/haunted-saas/feature-flag-service/proto/featureflags/v1";

service FeatureFlagService {
  rpc IsFeatureEnabled(IsFeatureEnabledRequest) returns (IsFeatureEnabledResponse);
  rpc GetAllFeatureFlags(GetAllFeatureFlagsRequest) returns (GetAllFeatureFlagsResponse);
  rpc GetFeatureFlagDetails(GetFeatureFlagDetailsRequest) returns (GetFeatureFlagDetailsResponse);
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}

message IsFeatureEnabledRequest {
  string feature_name = 1;
  EvaluationContext context = 2;
}

message EvaluationContext {
  string user_id = 1;
  string session_id = 2;
  string environment = 3;
  map<string, PropertyValue> properties = 4;
}

message PropertyValue {
  oneof value {
    string string_value = 1;
    double number_value = 2;
    bool bool_value = 3;
  }
}

message IsFeatureEnabledResponse {
  bool enabled = 1;
  string variant = 2;  // Optional variant name
}

message GetAllFeatureFlagsRequest {
  EvaluationContext context = 1;
}

message GetAllFeatureFlagsResponse {
  map<string, bool> flags = 1;
}

message GetFeatureFlagDetailsRequest {
  string feature_name = 1;
}

message GetFeatureFlagDetailsResponse {
  string name = 1;
  string description = 2;
  bool enabled = 3;
  repeated string strategies = 4;
}

message HealthCheckRequest {}

message HealthCheckResponse {
  string status = 1;
  bool unleash_connected = 2;
}
```

### Core Interfaces

```go
// Service layer interface
type FeatureFlagService interface {
    IsFeatureEnabled(ctx context.Context, featureName string, evalCtx *EvaluationContext) (bool, error)
    GetAllFeatureFlags(ctx context.Context, evalCtx *EvaluationContext) (map[string]bool, error)
    GetFeatureFlagDetails(ctx context.Context, featureName string) (*FeatureFlagDetails, error)
}

// Unleash wrapper interface
type UnleashClient interface {
    IsEnabled(featureName string, context *unleash.Context) bool
    GetVariant(featureName string, context *unleash.Context) *unleash.Variant
    ListFeatures() []string
    WaitForReady() error
}

// Cache interface
type CacheRepository interface {
    Get(ctx context.Context, key string) (bool, bool, error)  // value, found, error
    Set(ctx context.Context, key string, value bool, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

// Analytics tracker interface
type AnalyticsTracker interface {
    TrackFlagEvaluation(ctx context.Context, event *FlagEvaluationEvent) error
}

// Permission checker interface
type PermissionChecker interface {
    HasFeaturePermission(ctx context.Context, userID, featureName string) (bool, error)
}
```

## Data Models

### Go Domain Models

```go
type EvaluationContext struct {
    UserID      string
    SessionID   string
    Environment string
    Properties  map[string]interface{}
}

type FeatureFlagDetails struct {
    Name        string
    Description string
    Enabled     bool
    Strategies  []string
}

type FlagEvaluationEvent struct {
    FeatureName string
    UserID      string
    Enabled     bool
    Variant     string
    Strategy    string
    Timestamp   time.Time
}

// Unleash context conversion
func (ec *EvaluationContext) ToUnleashContext() *unleash.Context {
    return &unleash.Context{
        UserId:      ec.UserID,
        SessionId:   ec.SessionID,
        Environment: ec.Environment,
        Properties:  ec.Properties,
    }
}
```

### Redis Cache Keys

```
# Cache key format
feature_flag:{feature_name}:{user_id}:{context_hash} -> bool

# Example
feature_flag:new_dashboard:user123:abc123 -> true

# TTL: 30 seconds
```

## Error Handling

### Error Types

```go
var (
    ErrInvalidFeatureName  = errors.New("invalid feature name")
    ErrInvalidContext      = errors.New("invalid evaluation context")
    ErrUnleashUnavailable  = errors.New("unleash service unavailable")
    ErrPermissionDenied    = errors.New("user lacks permission for feature")
    ErrCacheUnavailable    = errors.New("cache unavailable")
)
```

### Error Handling Strategy

1. **Validation Errors**: Return immediately with descriptive error messages
2. **Unleash Errors**: Log error, return false as safe default
3. **Cache Errors**: Log warning, continue with Unleash evaluation
4. **Permission Errors**: Return false and log denial
5. **Panic Recovery**: Middleware catches panics, logs stack trace, returns false

### Fallback Behavior

```go
func (s *FeatureFlagService) IsFeatureEnabled(ctx context.Context, featureName string, evalCtx *EvaluationContext) (bool, error) {
    // Try cache first
    if cached, found, err := s.cache.Get(ctx, cacheKey); err == nil && found {
        return cached, nil
    }
    
    // Evaluate with Unleash
    enabled := s.unleash.IsEnabled(featureName, evalCtx.ToUnleashContext())
    
    // Cache result (ignore cache errors)
    _ = s.cache.Set(ctx, cacheKey, enabled, 30*time.Second)
    
    // Track evaluation asynchronously
    go s.trackEvaluation(featureName, evalCtx, enabled)
    
    return enabled, nil
}
```

## Integration with Unleash

### Unleash SDK Configuration

```go
func NewUnleashClient(cfg *config.Config) (*unleash.Client, error) {
    return unleash.NewClient(
        unleash.WithUrl(cfg.UnleashURL),
        unleash.WithAppName("feature-flag-service"),
        unleash.WithInstanceId(cfg.ServiceInstanceID),
        unleash.WithCustomHeaders(http.Header{
            "Authorization": []string{cfg.UnleashAPIToken},
        }),
        unleash.WithRefreshInterval(15*time.Second),
        unleash.WithMetricsInterval(60*time.Second),
        unleash.WithListener(&UnleashListener{
            logger: cfg.Logger,
        }),
    )
}

type UnleashListener struct {
    logger *zap.Logger
}

func (l *UnleashListener) OnReady() {
    l.logger.Info("Unleash client ready")
}

func (l *UnleashListener) OnError(err error) {
    l.logger.Error("Unleash client error", zap.Error(err))
}

func (l *UnleashListener) OnCount(name string, enabled bool) {
    // Metrics tracking
}
```

### Unleash Strategies Support

The service supports all standard Unleash strategies:

1. **Default**: Always on/off
2. **Gradual Rollout**: Percentage-based targeting
3. **User IDs**: Whitelist specific users
4. **IP Addresses**: Target by IP range
5. **Custom**: User-defined strategies

## Testing Strategy

### Unit Tests

- **Service Layer**: Mock Unleash client, test evaluation logic
- **Cache Layer**: Test cache hit/miss scenarios
- **Permission Checks**: Mock auth service, test RBAC integration
- **Analytics Tracking**: Mock analytics service, verify events

### Integration Tests

- **gRPC Handlers**: Test full request/response cycle
- **Unleash Integration**: Use Unleash test server
- **Redis Integration**: Use testcontainers for Redis
- **End-to-End**: Test complete flow with all dependencies

### Test Structure

```go
func TestFeatureFlagService_IsFeatureEnabled(t *testing.T) {
    tests := []struct {
        name           string
        featureName    string
        context        *EvaluationContext
        unleashEnabled bool
        cacheHit       bool
        cachedValue    bool
        want           bool
        wantErr        bool
    }{
        {
            name:        "cache hit returns cached value",
            featureName: "new_feature",
            context: &EvaluationContext{
                UserID: "user123",
            },
            cacheHit:    true,
            cachedValue: true,
            want:        true,
            wantErr:     false,
        },
        {
            name:        "cache miss evaluates with unleash",
            featureName: "new_feature",
            context: &EvaluationContext{
                UserID: "user123",
            },
            cacheHit:       false,
            unleashEnabled: true,
            want:           true,
            wantErr:        false,
        },
        {
            name:        "unknown feature returns false",
            featureName: "unknown_feature",
            context: &EvaluationContext{
                UserID: "user123",
            },
            cacheHit:       false,
            unleashEnabled: false,
            want:           false,
            wantErr:        false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Configuration

### Environment Variables

```
# Service
SERVICE_PORT=50052
SERVICE_NAME=feature-flag-service
SERVICE_INSTANCE_ID=instance-1
LOG_LEVEL=info

# Unleash
UNLEASH_URL=http://localhost:4242/api
UNLEASH_API_TOKEN=*:*.unleash-insecure-api-token
UNLEASH_REFRESH_INTERVAL=15s
UNLEASH_METRICS_INTERVAL=60s

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=1
REDIS_CACHE_TTL=30

# Analytics Integration
ANALYTICS_SERVICE_URL=localhost:50051
ANALYTICS_BATCH_SIZE=100
ANALYTICS_FLUSH_INTERVAL=10s

# Auth Integration
AUTH_SERVICE_URL=localhost:50053
RBAC_ENABLED=true
PERMISSION_CACHE_TTL=60

# Observability
PROMETHEUS_PORT=9091
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
```

## Docker Compose Integration

### Unleash Service Configuration

```yaml
services:
  unleash-db:
    image: postgres:15
    environment:
      POSTGRES_DB: unleash
      POSTGRES_USER: unleash
      POSTGRES_PASSWORD: unleash
    volumes:
      - unleash-db-data:/var/lib/postgresql/data

  unleash:
    image: unleashorg/unleash-server:latest
    ports:
      - "4242:4242"
    environment:
      DATABASE_URL: postgres://unleash:unleash@unleash-db:5432/unleash
      DATABASE_SSL: "false"
      INIT_ADMIN_API_TOKENS: "*:*.unleash-insecure-api-token"
      INIT_CLIENT_API_TOKENS: "*:*.unleash-insecure-api-token"
    depends_on:
      - unleash-db

  feature-flag-service:
    build: ./feature-flag-service
    ports:
      - "50052:50052"
    environment:
      UNLEASH_URL: http://unleash:4242/api
      UNLEASH_API_TOKEN: "*:*.unleash-insecure-api-token"
      REDIS_HOST: redis
      ANALYTICS_SERVICE_URL: analytics-service:50051
      AUTH_SERVICE_URL: user-auth-service:50053
    depends_on:
      - unleash
      - redis
```

## Security Considerations

1. **API Token Security**: Store Unleash API token in environment variables
2. **RBAC Integration**: Verify user permissions before returning flag state
3. **Rate Limiting**: Prevent abuse of evaluation endpoints
4. **Audit Logging**: Log all flag evaluations with user context
5. **mTLS**: Secure gRPC communication between services

## Performance Optimizations

1. **Redis Caching**: Cache evaluations for 30 seconds
2. **Async Analytics**: Track evaluations asynchronously to avoid blocking
3. **Connection Pooling**: Reuse HTTP connections to Unleash
4. **Batch Processing**: Batch analytics events before sending
5. **Context Hashing**: Hash evaluation context for efficient cache keys

## Observability

### Prometheus Metrics

```go
var (
    flagEvaluations = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "feature_flag_evaluations_total",
            Help: "Total number of feature flag evaluations",
        },
        []string{"feature_name", "enabled"},
    )
    
    cacheHitRate = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "feature_flag_cache_hit_rate",
            Help: "Cache hit rate for feature flag evaluations",
        },
        []string{"feature_name"},
    )
    
    unleashSyncStatus = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "unleash_sync_status",
            Help: "Unleash synchronization status (1=connected, 0=disconnected)",
        },
    )
)
```

### Structured Logging

```go
logger.Info("feature flag evaluated",
    zap.String("feature_name", featureName),
    zap.String("user_id", evalCtx.UserID),
    zap.Bool("enabled", enabled),
    zap.String("variant", variant),
    zap.Duration("latency", latency),
)
```
