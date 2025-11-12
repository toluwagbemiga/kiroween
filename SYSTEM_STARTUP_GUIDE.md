# System Startup Guide

Complete guide to running the entire HAUNTED SAAS SKELETON platform locally.

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (or use Docker)
- Redis (or use Docker)
- Node.js 18+ (for frontend)
- Protocol Buffers compiler (`protoc`)

## Quick Start (Docker Compose)

The fastest way to run everything:

```bash
# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

## Manual Startup (Development)

For development, you'll want to run services individually.

### Step 1: Start Infrastructure

```bash
# Start PostgreSQL
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=haunted_saas \
  -p 5432:5432 \
  postgres:15

# Start Redis
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:7-alpine

# Start Unleash (for feature flags)
docker run -d \
  --name unleash \
  -e DATABASE_URL=postgresql://postgres:postgres@host.docker.internal:5432/unleash \
  -p 4242:4242 \
  unleashorg/unleash-server:latest
```

### Step 2: Start Backend Services

Start each service in a separate terminal:

#### Terminal 1: User Auth Service

```bash
cd app/services/user-auth-service

# Set up environment
cp .env.example .env
# Edit .env with your database credentials

# Run migrations
make migrate-up

# Generate proto files
make proto

# Start service
make run

# Service will start on port 50051
```

#### Terminal 2: Billing Service

```bash
cd app/services/billing-service

# Set up environment
cp .env.example .env
# Add your Stripe API keys

# Run migrations
make migrate-up

# Generate proto files
make proto

# Start service
make run

# Service will start on port 50052
```

#### Terminal 3: LLM Gateway Service

```bash
cd app/services/llm-gateway-service

# Set up environment
cp .env.example .env
# Add your OpenAI API key

# Generate proto files
make proto

# Start service
make run

# Service will start on port 50053
```

#### Terminal 4: Notifications Service

```bash
cd app/services/notifications-service

# Set up environment
cp .env.example .env

# Generate proto files
make proto

# Start service
make run

# Service will start on:
# - gRPC: 50054
# - Socket.IO: 3001
```

#### Terminal 5: Analytics Service

```bash
cd app/services/analytics-service

# Set up environment
cp .env.example .env
# Add Mixpanel/Amplitude credentials (optional)

# Generate proto files
make proto

# Start service
make run

# Service will start on port 50055
```

#### Terminal 6: Feature Flags Service

```bash
cd app/services/feature-flags-service

# Set up environment
cp .env.example .env
# Add Unleash server URL and API token

# Generate proto files
make proto

# Start service
make run

# Service will start on port 50056
```

### Step 3: Start GraphQL API Gateway

```bash
cd app/gateway/graphql-api-gateway

# Set up environment
cp .env.example .env
# Verify all service addresses are correct

# Generate GraphQL code
make generate

# Start gateway
make run

# Gateway will start on port 8080
# GraphQL Playground: http://localhost:8080
```

### Step 4: Start Frontend (Optional)

```bash
cd app/frontend

# Install dependencies
npm install

# Set up environment
cp .env.example .env.local
# Set NEXT_PUBLIC_GRAPHQL_URL=http://localhost:8080/graphql

# Start development server
npm run dev

# Frontend will start on port 3000
```

## Verification

### Check All Services Are Running

```bash
# User Auth Service
grpcurl -plaintext localhost:50051 list

# Billing Service
grpcurl -plaintext localhost:50052 list

# LLM Gateway Service
grpcurl -plaintext localhost:50053 list

# Notifications Service
grpcurl -plaintext localhost:50054 list

# Analytics Service
grpcurl -plaintext localhost:50055 list

# Feature Flags Service
grpcurl -plaintext localhost:50056 list

# GraphQL Gateway
curl http://localhost:8080/health
```

### Test End-to-End Flow

#### 1. Register a User

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { register(input: { email: \"test@example.com\", password: \"SecurePass123!\", name: \"Test User\" }) { token user { id email name } } }"
  }'
```

#### 2. Login

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { login(input: { email: \"test@example.com\", password: \"SecurePass123!\" }) { token user { id email } } }"
  }'
```

Save the token from the response.

#### 3. Get Current User

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "query { me { id email name roles { name permissions } } }"
  }'
```

#### 4. Check Feature Flag

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "query { isFeatureEnabled(featureName: \"new_dashboard\") }"
  }'
```

#### 5. Call LLM Prompt

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "mutation { callPrompt(name: \"welcome-email\", variables: { user_name: \"Test\" }) { content tokensUsed cost } }"
  }'
```

## Service Ports Reference

| Service | gRPC Port | HTTP Port | Purpose |
|---------|-----------|-----------|---------|
| User Auth | 50051 | - | Authentication & RBAC |
| Billing | 50052 | 8081 | Stripe webhooks |
| LLM Gateway | 50053 | - | OpenAI proxy |
| Notifications | 50054 | 3001 | Socket.IO server |
| Analytics | 50055 | - | Event tracking |
| Feature Flags | 50056 | - | Unleash proxy |
| GraphQL Gateway | - | 8080 | Public API |
| Frontend | - | 3000 | Next.js app |

## Environment Variables Checklist

### Required for All Services

```bash
# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### User Auth Service

```bash
GRPC_PORT=50051
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/haunted_saas
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION_HOURS=24
```

### Billing Service

```bash
GRPC_PORT=50052
HTTP_PORT=8081
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/haunted_saas
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
```

### LLM Gateway Service

```bash
GRPC_PORT=50053
OPENAI_API_KEY=sk-...
PROMPTS_DIR=../../../prompts
```

### Notifications Service

```bash
GRPC_PORT=50054
SOCKETIO_PORT=3001
JWT_SECRET=your-secret-key-here
CORS_ORIGINS=http://localhost:3000
```

### Analytics Service

```bash
GRPC_PORT=50055
PROVIDER=mixpanel
MIXPANEL_TOKEN=your-token
BATCH_SIZE=100
FLUSH_INTERVAL_SECONDS=10
```

### Feature Flags Service

```bash
GRPC_PORT=50056
UNLEASH_SERVER_URL=http://localhost:4242/api
UNLEASH_API_TOKEN=*:*.your-token-here
UNLEASH_REFRESH_INTERVAL_SECONDS=10
```

### GraphQL Gateway

```bash
PORT=8080
USER_AUTH_SERVICE=localhost:50051
BILLING_SERVICE=localhost:50052
LLM_GATEWAY_SERVICE=localhost:50053
NOTIFICATIONS_SERVICE=localhost:50054
ANALYTICS_SERVICE=localhost:50055
FEATURE_FLAGS_SERVICE=localhost:50056
JWT_SECRET=your-secret-key-here
```

## Common Issues & Solutions

### Issue: "connection refused" on service startup

**Cause**: Service dependencies not running

**Solution**:
```bash
# Check PostgreSQL
docker ps | grep postgres

# Check Redis
docker ps | grep redis

# Check service logs
docker logs service-name
```

### Issue: "proto file not found"

**Cause**: Proto files not generated

**Solution**:
```bash
cd app/services/service-name
make proto
```

### Issue: "database does not exist"

**Cause**: Database not created

**Solution**:
```bash
# Create database
psql -U postgres -c "CREATE DATABASE haunted_saas;"

# Run migrations
cd app/services/user-auth-service
make migrate-up
```

### Issue: GraphQL Gateway can't connect to services

**Cause**: Service addresses incorrect

**Solution**:
```bash
# Check .env file
cat app/gateway/graphql-api-gateway/.env

# Verify services are running
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50052 list
# ... etc
```

### Issue: "Unleash client not ready"

**Cause**: Unleash server not accessible

**Solution**:
```bash
# Check Unleash is running
curl http://localhost:4242/health

# Check API token is correct
# Go to http://localhost:4242 → API Access → Create token
```

### Issue: Frontend can't connect to GraphQL

**Cause**: CORS or wrong URL

**Solution**:
```bash
# Check GraphQL gateway is running
curl http://localhost:8080/health

# Check CORS settings in gateway
# Update app/gateway/graphql-api-gateway/cmd/main.go
# AllowedOrigins: []string{"http://localhost:3000"}

# Check frontend .env.local
cat app/frontend/.env.local
# NEXT_PUBLIC_GRAPHQL_URL=http://localhost:8080/graphql
```

## Development Workflow

### Making Changes to a Service

```bash
# 1. Make code changes
vim app/services/user-auth-service/internal/service/auth_service.go

# 2. If proto changed, regenerate
make proto

# 3. Run tests
make test

# 4. Restart service
# Ctrl+C in terminal, then
make run
```

### Adding a New gRPC Method

```bash
# 1. Update proto file
vim app/services/user-auth-service/proto/userauth/v1/service.proto

# 2. Regenerate proto
make proto

# 3. Implement handler
vim app/services/user-auth-service/internal/handler/auth_handler.go

# 4. Update GraphQL schema (if exposing to frontend)
vim app/gateway/graphql-api-gateway/schema.graphqls

# 5. Regenerate GraphQL
cd app/gateway/graphql-api-gateway
make generate

# 6. Implement resolver
vim app/gateway/graphql-api-gateway/internal/resolvers/query.resolvers.go
```

### Database Migrations

```bash
# Create new migration
cd app/services/user-auth-service
migrate create -ext sql -dir migrations -seq add_new_table

# Edit migration files
vim migrations/000004_add_new_table.up.sql
vim migrations/000004_add_new_table.down.sql

# Run migration
make migrate-up

# Rollback if needed
make migrate-down
```

## Production Deployment

### Docker Build

```bash
# Build all services
docker-compose build

# Or build individually
cd app/services/user-auth-service
docker build -t user-auth-service:latest .
```

### Kubernetes Deployment

```bash
# Apply configurations
kubectl apply -f k8s/

# Check status
kubectl get pods
kubectl get services

# View logs
kubectl logs -f deployment/user-auth-service
```

### Environment Variables in Production

```bash
# Use Kubernetes secrets
kubectl create secret generic app-secrets \
  --from-literal=jwt-secret=your-secret \
  --from-literal=stripe-key=sk_live_... \
  --from-literal=openai-key=sk-...

# Reference in deployment
env:
  - name: JWT_SECRET
    valueFrom:
      secretKeyRef:
        name: app-secrets
        key: jwt-secret
```

## Monitoring

### Health Checks

```bash
# Check all services
for port in 50051 50052 50053 50054 50055 50056; do
  echo "Checking port $port..."
  grpcurl -plaintext localhost:$port list || echo "FAILED"
done

# Check GraphQL gateway
curl http://localhost:8080/health
```

### Logs

```bash
# View logs for all services
docker-compose logs -f

# View logs for specific service
docker-compose logs -f user-auth-service

# Follow logs in real-time
tail -f app/services/user-auth-service/logs/app.log
```

### Metrics (Future)

```bash
# Prometheus metrics endpoint (to be implemented)
curl http://localhost:8080/metrics

# Grafana dashboard
open http://localhost:3000
```

## Testing

### Unit Tests

```bash
# Test all services
for service in user-auth-service billing-service llm-gateway-service notifications-service analytics-service feature-flags-service; do
  echo "Testing $service..."
  cd app/services/$service
  make test
  cd ../../..
done
```

### Integration Tests

```bash
# Start all services
docker-compose up -d

# Run integration tests
go test ./tests/integration/...

# Stop services
docker-compose down
```

### Load Testing

```bash
# Install k6
brew install k6

# Run load test
k6 run tests/load/graphql-load-test.js
```

## Summary

### Startup Order

1. **Infrastructure**: PostgreSQL, Redis, Unleash
2. **Core Services**: User Auth, Feature Flags
3. **Business Services**: Billing, LLM Gateway, Notifications, Analytics
4. **Gateway**: GraphQL API Gateway
5. **Frontend**: Next.js app

### Verification Checklist

- [ ] All 6 backend services responding on their ports
- [ ] GraphQL gateway health check passes
- [ ] Can register and login a user
- [ ] Can query user data with JWT
- [ ] Feature flags return values
- [ ] LLM prompts execute successfully
- [ ] Frontend connects to GraphQL gateway

### Next Steps

1. Configure external services (Stripe, OpenAI, Unleash)
2. Set up production databases
3. Configure monitoring and logging
4. Set up CI/CD pipelines
5. Deploy to staging environment

---

**Status**: Complete startup guide  
**Last Updated**: 2024  
**Tested On**: macOS, Linux, Windows (WSL2)
