# HAUNTED SAAS SKELETON - Implementation Guide

This guide provides step-by-step instructions for completing the implementation of all services.

## 🎯 Implementation Philosophy

Each service follows this pattern:

1. **Proto First**: Define gRPC contracts
2. **Domain Models**: Create data structures
3. **Repository Layer**: Database/cache access
4. **Service Layer**: Business logic
5. **Handler Layer**: gRPC endpoints
6. **Tests**: Unit + integration tests

## 📚 Service Implementation Order

### 1. User Auth Service (PRIORITY 1)

**Why First**: All other services depend on authentication.

**Task List**: `.kiro/specs/user-auth/tasks.md`

**Key Steps**:

```bash
cd app/services/user-auth-service

# 1. Generate proto code
make proto

# 2. Create database migrations
mkdir -p migrations
# Write SQL files for users, roles, permissions tables

# 3. Implement repositories
# internal/repository/user_repository.go
# internal/repository/role_repository.go
# internal/repository/session_repository.go (Redis)
# internal/repository/permission_cache_repository.go (Redis)

# 4. Implement token manager
# internal/auth/token_manager.go
# - Load RSA keys
# - Generate JWT with RS256
# - Validate JWT

# 5. Implement rate limiter
# internal/auth/rate_limiter.go
# - Redis counter
# - 5 attempts = 30min lockout

# 6. Implement auth service
# internal/service/auth_service.go
# - Register (bcrypt hash, assign default role)
# - Login (verify, generate JWT, create session)
# - ValidateToken (verify JWT, check revocation)
# - Logout (remove session, revoke token)
# - Password reset flow

# 7. Implement RBAC service
# internal/service/rbac_service.go
# - CreateRole, UpdateRole, DeleteRole
# - AssignRoleToUser, RevokeRoleFromUser
# - CheckPermission (with caching)
# - GetUserPermissions

# 8. Implement gRPC handlers
# internal/handler/auth_handler.go
# internal/handler/rbac_handler.go
# - Wire up services
# - Error mapping

# 9. Wire everything in main.go
# - Load config
# - Connect to PostgreSQL and Redis
# - Initialize all components
# - Start gRPC server

# 10. Write tests
make test
```

**Testing Checklist**:
- [ ] User registration with valid/invalid data
- [ ] Login with correct/incorrect credentials
- [ ] Account lockout after 5 failed attempts
- [ ] JWT generation and validation
- [ ] Session management (create, extend, revoke)
- [ ] Role assignment and permission checks
- [ ] Password reset flow

### 2. Feature Flags Service

**Task List**: `.kiro/specs/feature-flags/tasks.md`

**Key Steps**:

```bash
cd app/services/feature-flags-service

# 1. Proto definitions
# proto/featureflags/v1/service.proto

# 2. Unleash client wrapper
# internal/service/unleash_client.go
# - Initialize with API token
# - WaitForReady
# - IsEnabled, GetVariant

# 3. Redis caching
# internal/repository/redis/cache.go
# - Cache key: feature_flag:{name}:{user_id}:{hash}
# - TTL: 30 seconds

# 4. Feature flags service
# internal/service/feature_flags_service.go
# - IsFeatureEnabled (cache → Unleash)
# - GetAllFeatureFlags
# - Track evaluation to Analytics

# 5. Permission checker
# internal/service/permissions.go
# - Call Auth Service to verify permissions
# - Cache results (60s TTL)

# 6. Analytics tracker
# internal/analytics/tracker.go
# - Async event emission
# - Batch processing

# 7. gRPC handlers
# internal/handler/handler.go

# 8. Test with Unleash UI
# Create flags at http://localhost:4242
```

### 3. Analytics Service

**Task List**: `.kiro/specs/analytics/tasks.md`

**Key Steps**:

```bash
cd app/services/analytics-service

# 1. Proto definitions
# proto/analytics/v1/service.proto

# 2. Database schema
# migrations/001_create_events_table.sql
# migrations/002_create_users_table.sql

# 3. Repositories
# internal/repository/postgres/events.go
# internal/repository/postgres/users.go
# internal/repository/redis/cache.go

# 4. Validation
# internal/service/validation.go
# - Event name, payload size, time ranges

# 5. Analytics service
# internal/service/analytics_service.go
# - TrackEvent (validate → store)
# - IdentifyUser (upsert properties)
# - GetEventCount (cache → query)
# - GetUserCount

# 6. gRPC handlers
# internal/handler/handler.go

# 7. Test with sample events
```

### 4. LLM Gateway Service

**Task List**: `.kiro/specs/llm-gateway/tasks.md`

**Key Steps**:

```bash
cd app/services/llm-gateway-service

# 1. Proto definitions
# proto/llm/v1/service.proto

# 2. Prompt loader
# internal/service/prompt_loader.go
# - Scan /prompts directory
# - Parse frontmatter
# - Extract variables
# - Cache prompts

# 3. File watcher
# - Use fsnotify
# - Reload on changes

# 4. Variable substitution
# internal/service/variable_substitutor.go
# - Parse JSON variables
# - Execute Go template
# - Validate required vars

# 5. OpenAI provider
# internal/service/openai_provider.go
# - Initialize client
# - Call ChatCompletion API
# - Parse response and token usage

# 6. LLM router
# internal/service/llm_router.go
# - Select provider
# - Select model
# - Retry on rate limits

# 7. Usage tracker
# internal/analytics/tracker.go
# - Emit to Analytics Service

# 8. gRPC handlers
# internal/handler/handler.go

# 9. Test with sample prompts
# Use prompts in /prompts directory
```

### 5. Notifications Service

**Task List**: `.kiro/specs/notifications/tasks.md`

**Key Steps**:

```bash
cd app/services/notifications-service

# 1. Proto definitions
# proto/notifications/v1/service.proto

# 2. Socket.IO server
# internal/socketio/server.go
# - Configure transports (WebSocket + polling)
# - Set ping intervals

# 3. JWT auth middleware
# internal/socketio/auth_middleware.go
# - Extract token from handshake
# - Verify JWT
# - Extract user ID and team ID

# 4. Connection manager
# internal/socketio/connection_manager.go
# - Thread-safe maps
# - Track connections by user
# - Update last seen

# 5. Room manager
# internal/socketio/room_manager.go
# - Auto-join user and team rooms
# - EmitToRoom, EmitToRoomExcept

# 6. Message router
# internal/service/message_router.go
# - SendToUser
# - SendToUsers (parallel)
# - BroadcastToRoom

# 7. Stats collector
# internal/service/stats_collector.go
# - Track connections
# - Track messages

# 8. gRPC handlers
# internal/handler/handler.go

# 9. Test with Socket.IO client
# See client example in README
```

### 6. Billing Service

**Task List**: `.kiro/specs/billing/tasks.md`

**Key Steps**:

```bash
cd app/services/billing-service

# 1. Proto definitions
# proto/billing/v1/service.proto

# 2. Database schema
# migrations/001_create_plans_table.sql
# migrations/002_create_subscriptions_table.sql
# migrations/003_create_webhook_events_table.sql

# 3. Stripe client wrapper
# internal/service/stripe_client.go
# - CreateProduct, CreatePrice
# - CreateCheckoutSession
# - GetSubscription, CancelSubscription, UpdateSubscription
# - CreateCustomer

# 4. Repositories
# internal/repository/plan_repository.go
# internal/repository/subscription_repository.go
# internal/repository/webhook_event_repository.go

# 5. Billing service
# internal/service/billing_service.go
# - CreatePlan (Stripe + DB)
# - CreateCheckoutSession
# - GetSubscription
# - CancelSubscription (at period end)
# - UpdateSubscription (with proration)

# 6. Webhook handler
# internal/handler/webhook_handler.go
# - Verify signature
# - Route events
# - Handle: checkout.session.completed, customer.subscription.*
# - Idempotency checks

# 7. gRPC handlers
# internal/handler/grpc_handler.go

# 8. Test with Stripe test mode
# Use Stripe CLI for webhook testing
```

## 🌐 GraphQL Gateway Implementation

**Location**: `app/gateway/`

**Steps**:

```bash
cd app/gateway

# 1. Initialize gqlgen
go run github.com/99designs/gqlgen init

# 2. Define schema
# schema/schema.graphql
# - Combine all service types
# - Define mutations and queries

# 3. Generate resolvers
go run github.com/99designs/gqlgen generate

# 4. Implement gRPC clients
# internal/clients/auth_client.go
# internal/clients/billing_client.go
# ... (one for each service)

# 5. Implement resolvers
# internal/resolvers/auth_resolver.go
# - Call gRPC clients
# - Map responses to GraphQL types

# 6. Add JWT middleware
# internal/middleware/auth.go
# - Verify JWT with public key
# - Extract user context

# 7. Wire up in main.go
# - Initialize all gRPC clients
# - Create GraphQL server
# - Add middleware

# 8. Test with GraphQL Playground
# http://localhost:4000/graphql
```

## 🎨 Frontend Implementation

**Location**: `app/frontend/`

**Steps**:

```bash
cd app/frontend

# 1. Initialize Next.js
npx create-next-app@latest . --typescript --tailwind --app

# 2. Install dependencies
npm install @apollo/client graphql socket.io-client
npm install react-hook-form zod @hookform/resolvers
npm install @radix-ui/react-* # For accessible components

# 3. Set up Apollo Client
# lib/graphql/client.ts

# 4. Create Design System
# components/ui/Button.tsx
# components/ui/Card.tsx
# components/ui/Input.tsx
# ... (all UI primitives)

# 5. Implement authentication
# app/(auth)/login/page.tsx
# app/(auth)/register/page.tsx
# lib/auth.ts

# 6. Create dashboard layout
# app/(dashboard)/layout.tsx
# components/dashboard/BentoGrid.tsx

# 7. Add Socket.IO
# lib/socket.ts
# components/notifications/NotificationToast.tsx

# 8. Implement pages
# app/(dashboard)/page.tsx - Dashboard
# app/(dashboard)/settings/page.tsx
# app/(dashboard)/billing/page.tsx

# 9. Test accessibility
# Use axe DevTools
# Test keyboard navigation
# Test screen readers
```

## 📚 Documentation Implementation

**Location**: `docs/`

**Steps**:

```bash
cd docs

# 1. Initialize Docusaurus
npx create-docusaurus@latest . classic --typescript

# 2. Configure structure
# docs/tutorial/
# docs/how-to/
# docs/reference/
# docs/explanation/

# 3. Write tutorial content
# docs/tutorial/getting-started.md
# docs/tutorial/authentication.md
# docs/tutorial/first-feature.md

# 4. Generate OpenAPI specs
# Use protoc-gen-openapi or grpc-gateway
# Output to static/openapi/

# 5. Configure OpenAPI plugin
# docusaurus.config.js

# 6. Write how-to guides
# docs/how-to/manage-users.md
# docs/how-to/configure-billing.md

# 7. Add code examples
# Use MDX for interactive examples

# 8. Configure search
# Algolia DocSearch

# 9. Deploy
# Vercel deployment via GitHub Actions
```

## 🧪 Testing Strategy

### Unit Tests
```bash
# Each service
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Target: >80% coverage
```

### Integration Tests
```bash
# Use testcontainers for PostgreSQL and Redis
# Test full gRPC flows
# Test webhook handling
```

### End-to-End Tests
```bash
# Start all services
docker-compose up -d

# Run demo data script
cd demo && npm run generate

# Test complete user flows
# - Registration → Login → Dashboard
# - Create subscription → Webhook → Update
# - Feature flag evaluation
# - Real-time notifications
```

## 🚀 Deployment Checklist

### Pre-deployment
- [ ] All tests passing
- [ ] Code coverage >80%
- [ ] Security audit (gosec, npm audit)
- [ ] Environment variables documented
- [ ] Database migrations tested
- [ ] Backup strategy defined

### Production Setup
- [ ] Generate production JWT keys (secure storage)
- [ ] Configure production Stripe keys
- [ ] Set up production database (managed PostgreSQL)
- [ ] Set up production Redis (managed Redis)
- [ ] Configure production Unleash
- [ ] Set up monitoring (Prometheus + Grafana)
- [ ] Configure logging (ELK or similar)
- [ ] Set up alerts
- [ ] Configure CDN for frontend
- [ ] Set up SSL certificates

### Deployment
- [ ] Deploy services to Kubernetes/ECS
- [ ] Run database migrations
- [ ] Verify health checks
- [ ] Test critical paths
- [ ] Monitor error rates
- [ ] Set up rollback plan

## 📊 Progress Tracking

Use this checklist to track overall progress:

- [ ] User Auth Service (100% complete)
- [ ] Feature Flags Service (100% complete)
- [ ] Analytics Service (100% complete)
- [ ] LLM Gateway Service (100% complete)
- [ ] Notifications Service (100% complete)
- [ ] Billing Service (100% complete)
- [ ] GraphQL Gateway (100% complete)
- [ ] Frontend (100% complete)
- [ ] Documentation (100% complete)
- [ ] Demo Data Script (100% complete)
- [ ] End-to-End Tests (100% complete)
- [ ] Production Deployment (100% complete)

## 🆘 Troubleshooting

### Common Issues

**Proto generation fails**:
```bash
# Install protoc and plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**Database connection fails**:
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Test connection
psql postgresql://haunted:haunted_dev_pass@localhost:5432/haunted
```

**Redis connection fails**:
```bash
# Check Redis is running
docker-compose ps redis

# Test connection
redis-cli -h localhost -p 6379 ping
```

**JWT verification fails**:
```bash
# Ensure keys are generated
cd keys && ./generate-keys.sh

# Check key permissions
ls -la keys/
```

## 📞 Support

- Refer to service READMEs for specific guidance
- Check spec files for requirements
- Review task lists for implementation steps
- Use GitHub Issues for tracking

---

**Happy Building! 🎃**
