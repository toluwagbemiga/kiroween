# CORS and Environment Variable Fix

## Problem

Frontend was getting CORS errors when trying to connect to GraphQL API:
```
Access to fetch at 'http://localhost:8080/graphql' from origin 'http://localhost:3000' 
has been blocked by CORS policy
```

## Root Causes

### 1. Wrong GraphQL URL
- Frontend was connecting to `http://localhost:8080/graphql`
- Should connect to `http://localhost:4000/graphql` (GraphQL Gateway)
- Port 8080 is the billing service webhook endpoint, not GraphQL

### 2. Environment Variables Not Baked Into Static Build
- Next.js static export requires env vars at **build time**
- Docker was only setting runtime env vars (which don't work for static sites)
- Need to pass as **build arguments** in Dockerfile

## Fixes Applied

### 1. Updated `.env.local.example`
Changed default URLs to correct ports:
```bash
NEXT_PUBLIC_GRAPHQL_URL=http://localhost:4000/graphql  # was 8080
NEXT_PUBLIC_SOCKETIO_URL=http://localhost:3002         # was 8085
```

### 2. Updated `Dockerfile`
Added build arguments and environment variables:
```dockerfile
# Build arguments for Next.js public env vars
ARG NEXT_PUBLIC_GRAPHQL_URL=http://localhost:4000/graphql
ARG NEXT_PUBLIC_SOCKETIO_URL=http://localhost:3002
ARG NEXT_PUBLIC_ANALYTICS_ENABLED=true
ARG NEXT_PUBLIC_FEATURE_FLAGS_ENABLED=true

# Set environment variables for build
ENV NEXT_PUBLIC_GRAPHQL_URL=$NEXT_PUBLIC_GRAPHQL_URL
ENV NEXT_PUBLIC_SOCKETIO_URL=$NEXT_PUBLIC_SOCKETIO_URL
# ... etc
```

### 3. Updated `docker-compose.yml`
Pass build args to Docker build:
```yaml
frontend:
  build:
    context: ./app/frontend
    dockerfile: Dockerfile
    args:
      NEXT_PUBLIC_GRAPHQL_URL: http://localhost:4000/graphql
      NEXT_PUBLIC_SOCKETIO_URL: http://localhost:3002
      NEXT_PUBLIC_ANALYTICS_ENABLED: "true"
      NEXT_PUBLIC_FEATURE_FLAGS_ENABLED: "true"
```

## CORS Configuration

The GraphQL Gateway already has CORS properly configured:
```go
corsHandler := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},  // Allows all origins
    AllowedMethods: []string{"GET", "POST", "OPTIONS"},
    AllowedHeaders: []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
})
```

**Note:** In production, change `AllowedOrigins` from `"*"` to specific domains.

## How to Apply

### Rebuild Frontend Container
```bash
# Stop services
docker-compose down

# Rebuild frontend with new build args
docker-compose build frontend

# Start all services
docker-compose up
```

Or rebuild everything:
```bash
docker-compose down
docker-compose build
docker-compose up
```

## Verification

### 1. Check Frontend Logs
```bash
docker-compose logs -f frontend
```

Should see:
```
INFO: Accepting connections at http://localhost:3000
```

### 2. Open Browser Console
Navigate to `http://localhost:3000/login`

Should NOT see CORS errors. GraphQL requests should go to `http://localhost:4000/graphql`

### 3. Test GraphQL Connection
Open browser dev tools → Network tab

Try to login or register. Should see:
- Request URL: `http://localhost:4000/graphql`
- Status: 200 OK (or 400 with GraphQL errors, but not CORS)

## Port Reference

| Service | Port | Purpose |
|---------|------|---------|
| Frontend | 3000 | Next.js static site |
| Socket.IO | 3002 | Real-time notifications |
| GraphQL Gateway | 4000 | Unified GraphQL API |
| Unleash | 4242 | Feature flags UI |
| User Auth | 50051 | gRPC (internal) |
| Billing | 50052 | gRPC (internal) |
| Billing Webhooks | 8080 | Stripe webhooks (HTTP) |
| LLM Gateway | 50053 | gRPC (internal) |
| Notifications | 50054 | gRPC (internal) |
| Analytics | 50055 | gRPC (internal) |
| Feature Flags | 50056 | gRPC (internal) |

## Important Notes

### Static Export Limitations
With `output: 'export'` in Next.js:
- All env vars must be set at **build time**
- Cannot use runtime env vars
- Cannot use server-side features (API routes work differently)
- Must use static file server (`serve`) instead of `next start`

### Environment Variables in Docker
For Next.js static export:
1. ✅ Use `ARG` in Dockerfile
2. ✅ Pass via `build.args` in docker-compose
3. ✅ Set as `ENV` before `npm run build`
4. ❌ Don't use `environment` in docker-compose (runtime only)

### Browser vs Container
- Frontend runs in **browser** (not in container)
- Browser accesses `localhost:4000` from **host machine**
- Container networking doesn't matter for frontend URLs
- Only backend services use container names (e.g., `user-auth-service:50051`)

## Summary

**Fixed:**
1. ✅ Corrected GraphQL URL (4000 not 8080)
2. ✅ Added build-time env vars to Dockerfile
3. ✅ Configured docker-compose build args
4. ✅ Updated .env.local.example with correct defaults

**Result:**
- No more CORS errors
- Frontend connects to correct GraphQL endpoint
- All environment variables properly baked into static build
