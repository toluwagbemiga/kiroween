# Docker Startup Fixes Applied

## Issues Fixed

### 1. ✅ Unleash Token Configuration
**Error:** `BadDataError: Client token cannot be scoped to all environments`

**Fix:** Changed token format from `*:*` to proper environment:project format
- Admin token: `default:development.unleash-insecure-admin-token`
- Client token: `development:default.unleash-insecure-client-token`

**Files:** `docker-compose.yml`

### 2. ✅ Missing JWT Keys
**Error:** `failed to read private key: open /app/keys/jwt-private.pem: no such file or directory`

**Fix:** Generated RSA key pair using OpenSSL
```bash
openssl genrsa -out keys/jwt-private.pem 2048
openssl rsa -in keys/jwt-private.pem -pubout -out keys/jwt-public.pem
```

**Files:** 
- `keys/jwt-private.pem` (2048-bit RSA private key)
- `keys/jwt-public.pem` (public key)

### 3. ✅ Analytics Service Missing API Key
**Error:** `MIXPANEL_API_KEY is required (or enable TEST_MODE)`

**Fix:** Added `TEST_MODE: "true"` to analytics-service environment

**Files:** `docker-compose.yml`

## Services Status After Fixes

### ✅ Running Services
- postgres (port 5432)
- redis (port 6379)
- unleash-db (internal)
- unleash (port 4242)
- billing-service (ports 50052, 8080)
- notifications-service (ports 50054, 3002)

### ⚠️ Fixed Services (will start on next run)
- user-auth-service (port 50051) - JWT keys now available
- analytics-service (port 50055) - TEST_MODE enabled
- feature-flags-service (port 50056) - Correct Unleash token

### 🔄 Dependent Services (will start after fixes)
- llm-gateway-service (port 50053) - depends on analytics
- graphql-gateway (port 4000) - depends on all services
- frontend (port 3000) - depends on gateway

## How to Start

### Stop current containers
```bash
docker-compose down
```

### Start all services
```bash
docker-compose up
```

Or start in detached mode:
```bash
docker-compose up -d
```

### Check service logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f user-auth-service
docker-compose logs -f analytics-service
docker-compose logs -f unleash
```

### Verify services are running
```bash
docker-compose ps
```

## Service Dependencies

```
postgres, redis, unleash-db
    ↓
unleash
    ↓
user-auth-service, analytics-service, billing-service, notifications-service
    ↓
feature-flags-service, llm-gateway-service
    ↓
graphql-gateway
    ↓
frontend
```

## Environment Variables

All services now have proper configuration:
- JWT keys mounted from `./keys` directory
- TEST_MODE enabled for analytics (no external API needed)
- Correct Unleash token format
- All database connections configured

## Next Steps

1. Run `docker-compose down` to stop current containers
2. Run `docker-compose up` to start with fixes
3. Wait for all services to be healthy
4. Access frontend at http://localhost:3000
5. Access GraphQL playground at http://localhost:4000/graphql
6. Access Unleash UI at http://localhost:4242
