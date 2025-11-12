# Billing Service

Stripe integration for subscription management and payment processing.

## Features

- ✅ Subscription plan management (CRUD)
- ✅ Stripe Checkout Session creation
- ✅ Subscription lifecycle (create, update, cancel)
- ✅ Webhook handling for Stripe events
- ✅ Prorated plan changes
- ✅ Secure API key management

## Architecture

```
cmd/server/main.go
internal/
  ├── domain/                   # Plan, Subscription, WebhookEvent models
  ├── repository/               # PostgreSQL repositories
  ├── service/                  # Business logic
  │   ├── billing_service.go
  │   └── stripe_client.go
  ├── handler/                  # gRPC + HTTP webhook handlers
  └── config/
```

## Implementation Checklist

See `.kiro/specs/billing/tasks.md` for full task list.

### Key Components
- [ ] Proto definitions (CreatePlan, CreateCheckoutSession, etc.)
- [ ] Stripe client wrapper
- [ ] Plan repository (GORM)
- [ ] Subscription repository
- [ ] Webhook signature verification
- [ ] Webhook event handlers (checkout.session.completed, etc.)
- [ ] Idempotency tracking

## Environment Variables

```bash
GRPC_PORT=50052
HTTP_PORT=8080
DATABASE_URL=postgresql://...
STRIPE_API_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
```

## Endpoints

**gRPC:**
- CreatePlan, GetPlan, ListPlans, UpdatePlan, DeactivatePlan
- CreateCheckoutSession, GetSubscription, CancelSubscription, UpdateSubscription

**HTTP:**
- POST /webhooks/stripe - Stripe webhook endpoint

## Stripe Integration

Uses `github.com/stripe/stripe-go/v76` for:
- Product and Price creation
- Checkout Session management
- Subscription lifecycle
- Customer management

## Security

- Webhook signature verification (HMAC-SHA256)
- API keys from environment only
- Separate test/production keys
- Never log sensitive data

## Next Steps

1. Implement proto definitions
2. Create Stripe client wrapper
3. Implement plan management
4. Implement subscription flows
5. Add webhook handlers with signature verification
6. Test with Stripe test mode
