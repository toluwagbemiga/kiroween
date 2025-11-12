# Design Document: Billing Service

## Overview

The Billing Service is a Go-based microservice that manages all payment processing, subscription plans, and billing operations through Stripe integration. It exposes gRPC endpoints for internal service communication and provides HTTP endpoints for Stripe webhook callbacks. The service maintains local subscription state synchronized with Stripe through webhook events.

## Architecture

### High-Level Architecture

```mermaid
graph TB
    subgraph "External"
        Stripe[Stripe API]
        StripeWebhook[Stripe Webhooks]
    end
    
    subgraph "Billing Service"
        gRPCServer[gRPC Server]
        HTTPServer[HTTP Server - Webhooks]
        BillingLogic[Billing Logic Layer]
        StripeClient[Stripe Client]
        DB[(PostgreSQL)]
    end
    
    subgraph "Internal Services"
        AuthService[Auth Service]
        OtherServices[Other Services]
    end
    
    OtherServices -->|gRPC| gRPCServer
    AuthService -->|gRPC| gRPCServer
    gRPCServer --> BillingLogic
    HTTPServer --> BillingLogic
    BillingLogic --> StripeClient
    BillingLogic --> DB
    StripeClient <-->|HTTPS| Stripe
    StripeWebhook -->|HTTPS| HTTPServer
```

### Technology Stack

- **Language**: Go 1.21+
- **gRPC Framework**: google.golang.org/grpc
- **HTTP Framework**: net/http (standard library)
- **Database**: PostgreSQL with pgx driver and GORM
- **Stripe SDK**: github.com/stripe/stripe-go/v76
- **Configuration**: Environment variables via godotenv
- **Logging**: Structured logging with zerolog

## Components and Interfaces

### 1. gRPC Service Definition

```protobuf
syntax = "proto3";

package billing.v1;

service BillingService {
  // Plan Management
  rpc CreatePlan(CreatePlanRequest) returns (CreatePlanResponse);
  rpc GetPlan(GetPlanRequest) returns (GetPlanResponse);
  rpc ListPlans(ListPlansRequest) returns (ListPlansResponse);
  rpc UpdatePlan(UpdatePlanRequest) returns (UpdatePlanResponse);
  rpc DeactivatePlan(DeactivatePlanRequest) returns (DeactivatePlanResponse);
  
  // Subscription Management
  rpc CreateCheckoutSession(CreateCheckoutSessionRequest) returns (CreateCheckoutSessionResponse);
  rpc GetSubscription(GetSubscriptionRequest) returns (GetSubscriptionResponse);
  rpc CancelSubscription(CancelSubscriptionRequest) returns (CancelSubscriptionResponse);
  rpc UpdateSubscription(UpdateSubscriptionRequest) returns (UpdateSubscriptionResponse);
}

message CreatePlanRequest {
  string name = 1;
  int64 price_cents = 2;
  string currency = 3;
  string billing_interval = 4; // "month" or "year"
  map<string, string> features = 5;
  string created_by_user_id = 6;
}

message CreatePlanResponse {
  Plan plan = 1;
}

message CreateCheckoutSessionRequest {
  string team_id = 1;
  string plan_id = 2;
  string success_url = 3;
  string cancel_url = 4;
}

message CreateCheckoutSessionResponse {
  string checkout_url = 1;
  string session_id = 2;
}

message GetSubscriptionRequest {
  string team_id = 1;
}

message GetSubscriptionResponse {
  Subscription subscription = 1;
}

message CancelSubscriptionRequest {
  string team_id = 1;
  string requesting_user_id = 2;
}

message CancelSubscriptionResponse {
  string cancellation_date = 1;
  Subscription subscription = 2;
}

message UpdateSubscriptionRequest {
  string team_id = 1;
  string new_plan_id = 2;
  string requesting_user_id = 3;
}

message UpdateSubscriptionResponse {
  Subscription subscription = 1;
  int64 next_billing_amount_cents = 2;
}

message Plan {
  string id = 1;
  string name = 2;
  int64 price_cents = 3;
  string currency = 4;
  string billing_interval = 5;
  map<string, string> features = 6;
  bool is_active = 7;
  string stripe_price_id = 8;
  string stripe_product_id = 9;
  string created_at = 10;
  string updated_at = 11;
}

message Subscription {
  string id = 1;
  string team_id = 2;
  string plan_id = 3;
  string status = 4; // "active", "canceled", "past_due", "incomplete"
  string stripe_subscription_id = 5;
  string stripe_customer_id = 6;
  string current_period_start = 7;
  string current_period_end = 8;
  string cancel_at = 9;
  string created_at = 10;
  string updated_at = 11;
}
```

### 2. HTTP Webhook Handler

```go
type WebhookHandler struct {
    billingService *BillingService
    webhookSecret  string
}

// POST /webhooks/stripe
func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
    // 1. Read request body
    // 2. Verify Stripe signature
    // 3. Parse event
    // 4. Route to appropriate handler
    // 5. Return 200 immediately
}

// Supported webhook events:
// - checkout.session.completed
// - customer.subscription.created
// - customer.subscription.updated
// - customer.subscription.deleted
// - invoice.payment_succeeded
// - invoice.payment_failed
```

### 3. Core Business Logic Layer

```go
type BillingService struct {
    db            *gorm.DB
    stripeClient  *stripe.Client
    planRepo      *PlanRepository
    subRepo       *SubscriptionRepository
    webhookQueue  *WebhookQueue
}

// Plan operations
func (s *BillingService) CreatePlan(ctx context.Context, req *CreatePlanRequest) (*Plan, error)
func (s *BillingService) GetPlan(ctx context.Context, planID string) (*Plan, error)
func (s *BillingService) ListPlans(ctx context.Context, activeOnly bool) ([]*Plan, error)
func (s *BillingService) UpdatePlan(ctx context.Context, planID string, updates *PlanUpdates) (*Plan, error)
func (s *BillingService) DeactivatePlan(ctx context.Context, planID string) error

// Subscription operations
func (s *BillingService) CreateCheckoutSession(ctx context.Context, req *CheckoutRequest) (*CheckoutSession, error)
func (s *BillingService) GetSubscription(ctx context.Context, teamID string) (*Subscription, error)
func (s *BillingService) CancelSubscription(ctx context.Context, teamID string) (*Subscription, error)
func (s *BillingService) UpdateSubscription(ctx context.Context, teamID, newPlanID string) (*Subscription, error)

// Webhook processing
func (s *BillingService) ProcessWebhook(ctx context.Context, event *stripe.Event) error
```

### 4. Stripe Client Wrapper

```go
type StripeClientWrapper struct {
    client *client.API
}

// Product and Price management
func (c *StripeClientWrapper) CreateProduct(name string, metadata map[string]string) (*stripe.Product, error)
func (c *StripeClientWrapper) CreatePrice(productID string, amount int64, currency, interval string) (*stripe.Price, error)
func (c *StripeClientWrapper) UpdateProduct(productID string, updates *ProductUpdates) (*stripe.Product, error)

// Checkout Session
func (c *StripeClientWrapper) CreateCheckoutSession(params *CheckoutSessionParams) (*stripe.CheckoutSession, error)

// Subscription management
func (c *StripeClientWrapper) GetSubscription(subscriptionID string) (*stripe.Subscription, error)
func (c *StripeClientWrapper) CancelSubscription(subscriptionID string, cancelAtPeriodEnd bool) (*stripe.Subscription, error)
func (c *StripeClientWrapper) UpdateSubscription(subscriptionID string, newPriceID string) (*stripe.Subscription, error)

// Customer management
func (c *StripeClientWrapper) CreateCustomer(email, teamID string) (*stripe.Customer, error)
func (c *StripeClientWrapper) GetCustomer(customerID string) (*stripe.Customer, error)
```

## Data Models

### Database Schema

```sql
-- Plans table
CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    price_cents BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'usd',
    billing_interval VARCHAR(20) NOT NULL, -- 'month' or 'year'
    features JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    stripe_price_id VARCHAR(255) NOT NULL UNIQUE,
    stripe_product_id VARCHAR(255) NOT NULL,
    created_by_user_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_interval CHECK (billing_interval IN ('month', 'year')),
    CONSTRAINT positive_price CHECK (price_cents >= 0)
);

CREATE INDEX idx_plans_active ON plans(is_active) WHERE is_active = true;
CREATE INDEX idx_plans_stripe_price ON plans(stripe_price_id);

-- Subscriptions table
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL UNIQUE, -- One subscription per team
    plan_id UUID NOT NULL REFERENCES plans(id),
    status VARCHAR(50) NOT NULL,
    stripe_subscription_id VARCHAR(255) NOT NULL UNIQUE,
    stripe_customer_id VARCHAR(255) NOT NULL,
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    cancel_at TIMESTAMP,
    canceled_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_status CHECK (status IN ('active', 'canceled', 'past_due', 'incomplete', 'trialing', 'unpaid'))
);

CREATE INDEX idx_subscriptions_team ON subscriptions(team_id);
CREATE INDEX idx_subscriptions_stripe_sub ON subscriptions(stripe_subscription_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);

-- Webhook events table (for idempotency and dead letter queue)
CREATE TABLE webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_event_id VARCHAR(255) NOT NULL UNIQUE,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    processed BOOLEAN NOT NULL DEFAULT false,
    processing_error TEXT,
    received_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP
);

CREATE INDEX idx_webhook_events_stripe_id ON webhook_events(stripe_event_id);
CREATE INDEX idx_webhook_events_processed ON webhook_events(processed);
CREATE INDEX idx_webhook_events_type ON webhook_events(event_type);
```

### Go Domain Models

```go
type Plan struct {
    ID              string            `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Name            string            `gorm:"not null"`
    PriceCents      int64             `gorm:"not null"`
    Currency        string            `gorm:"not null;default:usd"`
    BillingInterval string            `gorm:"not null"`
    Features        map[string]string `gorm:"type:jsonb;not null;default:'{}'"`
    IsActive        bool              `gorm:"not null;default:true"`
    StripePriceID   string            `gorm:"not null;unique"`
    StripeProductID string            `gorm:"not null"`
    CreatedByUserID *string           `gorm:"type:uuid"`
    CreatedAt       time.Time         `gorm:"not null;default:now()"`
    UpdatedAt       time.Time         `gorm:"not null;default:now()"`
}

type Subscription struct {
    ID                   string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    TeamID               string    `gorm:"type:uuid;not null;unique"`
    PlanID               string    `gorm:"type:uuid;not null"`
    Status               string    `gorm:"not null"`
    StripeSubscriptionID string    `gorm:"not null;unique"`
    StripeCustomerID     string    `gorm:"not null"`
    CurrentPeriodStart   time.Time `gorm:"not null"`
    CurrentPeriodEnd     time.Time `gorm:"not null"`
    CancelAt             *time.Time
    CanceledAt           *time.Time
    CreatedAt            time.Time `gorm:"not null;default:now()"`
    UpdatedAt            time.Time `gorm:"not null;default:now()"`
    
    Plan Plan `gorm:"foreignKey:PlanID"`
}

type WebhookEvent struct {
    ID              string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    StripeEventID   string          `gorm:"not null;unique"`
    EventType       string          `gorm:"not null"`
    Payload         json.RawMessage `gorm:"type:jsonb;not null"`
    Processed       bool            `gorm:"not null;default:false"`
    ProcessingError *string
    ReceivedAt      time.Time `gorm:"not null;default:now()"`
    ProcessedAt     *time.Time
}
```

## Error Handling

### Error Types

```go
type BillingError struct {
    Code    string
    Message string
    Cause   error
}

const (
    ErrCodePlanNotFound         = "PLAN_NOT_FOUND"
    ErrCodePlanInactive         = "PLAN_INACTIVE"
    ErrCodeSubscriptionNotFound = "SUBSCRIPTION_NOT_FOUND"
    ErrCodeStripeAPIError       = "STRIPE_API_ERROR"
    ErrCodeInvalidRequest       = "INVALID_REQUEST"
    ErrCodeUnauthorized         = "UNAUTHORIZED"
    ErrCodeDuplicateSubscription = "DUPLICATE_SUBSCRIPTION"
    ErrCodeWebhookVerification  = "WEBHOOK_VERIFICATION_FAILED"
)
```

### Error Handling Strategy

1. **gRPC Errors**: Map internal errors to appropriate gRPC status codes
   - `PLAN_NOT_FOUND` → `codes.NotFound`
   - `UNAUTHORIZED` → `codes.PermissionDenied`
   - `INVALID_REQUEST` → `codes.InvalidArgument`
   - `STRIPE_API_ERROR` → `codes.Unavailable`

2. **Webhook Errors**: Always return HTTP 200 to prevent Stripe retries, log errors to dead letter queue

3. **Stripe API Errors**: Wrap Stripe errors with context, log full error details, return user-friendly messages

4. **Database Errors**: Retry transient errors (connection issues), log and return internal error for others

### Logging Strategy

```go
type LogEntry struct {
    Level         string
    Message       string
    CorrelationID string
    UserID        string
    TeamID        string
    PlanID        string
    StripeEventID string
    Error         error
    Timestamp     time.Time
}

// Log levels:
// - INFO: Successful operations (plan created, subscription updated)
// - WARN: Recoverable issues (webhook already processed)
// - ERROR: Failures requiring attention (Stripe API errors, database errors)
// - DEBUG: Detailed flow information (development only)
```

## Testing Strategy

### Unit Tests

1. **Business Logic Tests**
   - Test plan creation with valid/invalid data
   - Test subscription lifecycle (create, update, cancel)
   - Test webhook event processing with various event types
   - Mock Stripe client and database

2. **Repository Tests**
   - Test CRUD operations with test database
   - Test query methods (find by team, find active plans)
   - Test transaction handling

3. **Stripe Client Wrapper Tests**
   - Mock Stripe API responses
   - Test error handling for API failures
   - Test request parameter construction

### Integration Tests

1. **gRPC Endpoint Tests**
   - Test full request/response cycle
   - Test authentication and authorization
   - Test error responses

2. **Webhook Handler Tests**
   - Test signature verification
   - Test event routing
   - Test idempotency

3. **Database Integration Tests**
   - Test with real PostgreSQL instance (testcontainers)
   - Test concurrent operations
   - Test transaction rollback scenarios

### Test Data

```go
// Test fixtures
var (
    TestPlan = &Plan{
        ID:              "test-plan-id",
        Name:            "Pro Plan",
        PriceCents:      2999,
        Currency:        "usd",
        BillingInterval: "month",
        Features:        map[string]string{"users": "10", "storage": "100GB"},
        IsActive:        true,
        StripePriceID:   "price_test123",
        StripeProductID: "prod_test123",
    }
    
    TestSubscription = &Subscription{
        ID:                   "test-sub-id",
        TeamID:               "test-team-id",
        PlanID:               "test-plan-id",
        Status:               "active",
        StripeSubscriptionID: "sub_test123",
        StripeCustomerID:     "cus_test123",
        CurrentPeriodStart:   time.Now(),
        CurrentPeriodEnd:     time.Now().AddDate(0, 1, 0),
    }
)
```

## Configuration

### Environment Variables

```bash
# Server Configuration
GRPC_PORT=50051
HTTP_PORT=8080

# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/billing?sslmode=disable

# Stripe
STRIPE_API_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_API_VERSION=2023-10-16

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Service Discovery (optional)
AUTH_SERVICE_ADDR=auth-service:50051
```

### Startup Sequence

1. Load and validate environment variables
2. Initialize database connection and run migrations
3. Initialize Stripe client with API key
4. Start gRPC server on configured port
5. Start HTTP server for webhooks on configured port
6. Register health check endpoints
7. Log startup completion with service version

## Security Considerations

1. **API Key Management**
   - Load Stripe keys from environment only
   - Never log or expose keys in responses
   - Use separate keys for test/production

2. **Webhook Verification**
   - Always verify Stripe signature before processing
   - Use constant-time comparison for signature validation
   - Reject requests with invalid signatures immediately

3. **Authorization**
   - Validate user permissions before plan management operations
   - Ensure team admins can only manage their own team's subscription
   - Super admin role required for plan creation/modification

4. **Data Protection**
   - Store minimal customer data (no credit card info)
   - Use Stripe customer IDs for reference
   - Encrypt sensitive data at rest (database encryption)

5. **Rate Limiting**
   - Implement rate limiting on gRPC endpoints
   - Separate limits for plan management vs subscription operations
   - Higher limits for webhook endpoint
