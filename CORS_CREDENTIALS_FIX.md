# CORS Credentials Fix

## Problem

After fixing the GraphQL URL, got a new CORS error:
```
The value of the 'Access-Control-Allow-Origin' header in the response must not be 
the wildcard '*' when the request's credentials mode is 'include'.
```

## Root Cause

Apollo Client was configured with `credentials: 'include'`, which tells the browser to send cookies with requests. However, when using credentials, CORS doesn't allow wildcard (`*`) origins - you must specify exact origins.

## Solution Options

### Option 1: Remove credentials (Chosen for Development)
Simpler for development. JWT tokens are sent via Authorization header, not cookies.

**Changed in `app/frontend/src/lib/apollo-client.ts`:**
```typescript
const httpLink = new HttpLink({
  uri: process.env.NEXT_PUBLIC_GRAPHQL_URL || 'http://localhost:4000/graphql',
  // Removed: credentials: 'include',
});
```

### Option 2: Specify Exact Origins (For Production)
If you need cookie-based auth, update gateway CORS:

**In `app/gateway/graphql-api-gateway/cmd/main.go`:**
```go
corsHandler := cors.New(cors.Options{
    AllowedOrigins: []string{
        "http://localhost:3000",
        "https://yourdomain.com",
    },
    AllowedMethods: []string{"GET", "POST", "OPTIONS"},
    AllowedHeaders: []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
})
```

## Why This Works

Our auth system uses **JWT tokens in Authorization headers**, not cookies:

```typescript
// Auth middleware adds token to headers
const authMiddleware = new ApolloLink((operation, forward) => {
  const token = getToken();  // From localStorage
  
  operation.setContext(({ headers = {} }) => ({
    headers: {
      ...headers,
      ...(token ? { authorization: `Bearer ${token}` } : {}),
    },
  }));

  return forward(operation);
});
```

So we don't need `credentials: 'include'` at all!

## How to Apply

### Rebuild Frontend
```bash
docker-compose down
docker-compose build frontend
docker-compose up
```

Or just restart if using local development:
```bash
cd app/frontend
npm run build
npm run dev
```

## Verification

1. Open browser console at `http://localhost:3000/login`
2. Try to login/register
3. Should see successful GraphQL requests to `http://localhost:4000/graphql`
4. No CORS errors

## Authentication Flow

With this fix, authentication works as follows:

1. **Login/Register** → GraphQL mutation
2. **Receive JWT token** → Store in localStorage
3. **Subsequent requests** → Token sent in `Authorization: Bearer <token>` header
4. **Gateway validates** → Calls user-auth-service to validate token
5. **Response** → Data returned to frontend

No cookies needed, so no `credentials: 'include'` required!

## Summary

**Fixed:** Removed `credentials: 'include'` from Apollo Client
**Why:** JWT auth uses headers, not cookies
**Result:** CORS works with wildcard origins
**Trade-off:** Can't use cookie-based auth (but we don't need it)
