# GraphQL API Gateway - Quick Start Guide

Get the gateway running in 5 minutes.

## Prerequisites

- Go 1.21+
- All 6 backend services running (or at least user-auth-service)
- Docker (optional)

## Step 1: Environment Setup

```bash
cd app/gateway/graphql-api-gateway
cp .env.example .env
```

Edit `.env`:

```bash
# Minimum required for local development
USER_AUTH_SERVICE=localhost:50051
BILLING_SERVICE=localhost:50052
LLM_GATEWAY_SERVICE=localhost:50053
NOTIFICATIONS_SERVICE=localhost:50054
ANALYTICS_SERVICE=localhost:50055
FEATURE_FLAGS_SERVICE=localhost:50056

JWT_SECRET=dev-secret-change-in-production
```

## Step 2: Generate GraphQL Code

```bash
make generate
```

This creates:
- `internal/generated/exec.go`
- `internal/generated/models.go`
- Resolver stubs

## Step 3: Install Dependencies

```bash
go mod download
```

## Step 4: Run the Gateway

```bash
make run
```

You should see:

```
🎃 Starting GraphQL API Gateway
✓ connected to user-auth-service
✓ connected to billing-service
✓ connected to llm-gateway-service
✓ connected to notifications-service
✓ connected to analytics-service
✓ connected to feature-flags-service
✓ all gRPC clients initialized
🚀 GraphQL API Gateway started address=0.0.0.0:8080
```

## Step 5: Test with GraphQL Playground

Open http://localhost:8080 in your browser.

### Test Query 1: Register

```graphql
mutation Register {
  register(input: {
    email: "test@example.com"
    password: "SecurePass123!"
    name: "Test User"
  }) {
    token
    user {
      id
      email
      name
    }
  }
}
```

### Test Query 2: Login

```graphql
mutation Login {
  login(input: {
    email: "test@example.com"
    password: "SecurePass123!"
  }) {
    token
    user {
      id
      email
      roles {
        name
        permissions
      }
    }
  }
}
```

### Test Query 3: Get Current User (Authenticated)

Add the token to HTTP Headers:

```json
{
  "Authorization": "Bearer <your-token-here>"
}
```

Then query:

```graphql
query Me {
  me {
    id
    email
    name
    roles {
      name
      permissions
    }
  }
}
```

## Common Issues

### Issue: "Failed to connect to user-auth-service"

**Solution**: Make sure user-auth-service is running on port 50051

```bash
cd app/services/user-auth-service
make run
```

### Issue: "JWT_SECRET is required in production"

**Solution**: Set JWT_SECRET in .env file

```bash
JWT_SECRET=your-secret-here
```

### Issue: "CORS error in browser"

**Solution**: CORS is configured to allow all origins in development. For production, update `cmd/main.go`:

```go
corsHandler := cors.New(cors.Options{
    AllowedOrigins: []string{"https://yourdomain.com"},
    // ...
})
```

## Docker Quick Start

```bash
# Build
docker build -t graphql-gateway .

# Run
docker run -p 8080:8080 \
  -e USER_AUTH_SERVICE=host.docker.internal:50051 \
  -e BILLING_SERVICE=host.docker.internal:50052 \
  -e LLM_GATEWAY_SERVICE=host.docker.internal:50053 \
  -e NOTIFICATIONS_SERVICE=host.docker.internal:50054 \
  -e ANALYTICS_SERVICE=host.docker.internal:50055 \
  -e FEATURE_FLAGS_SERVICE=host.docker.internal:50056 \
  -e JWT_SECRET=dev-secret \
  graphql-gateway
```

## Next Steps

1. **Explore the Schema**: Use GraphQL Playground's "Docs" tab
2. **Test All Queries**: Try billing, feature flags, LLM calls
3. **Integrate Frontend**: See README.md for React/Vue examples
4. **Add Custom Resolvers**: Extend `internal/resolvers/`
5. **Deploy**: Use Docker or Kubernetes

## Useful Commands

```bash
# Generate GraphQL code
make generate

# Run tests
make test

# Build binary
make build

# Run with auto-reload (requires air)
make dev

# Lint code
make lint

# Build Docker image
make docker-build
```

## Health Check

```bash
curl http://localhost:8080/health
# {"status":"healthy"}
```

## Example Frontend Integration

```typescript
// Apollo Client setup
import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client';

const client = new ApolloClient({
  uri: 'http://localhost:8080/graphql',
  cache: new InMemoryCache(),
});

// Use in React
import { useQuery, gql } from '@apollo/client';

const GET_ME = gql`
  query GetMe {
    me {
      id
      email
      name
    }
  }
`;

function Profile() {
  const { loading, error, data } = useQuery(GET_ME);
  
  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;
  
  return <h1>Welcome, {data.me.name}!</h1>;
}
```

## Support

- **Documentation**: See README.md
- **Implementation Details**: See IMPLEMENTATION_COMPLETE.md
- **Schema Reference**: Use GraphQL Playground introspection

---

**You're ready to go!** 🚀

The gateway is now running and ready to handle requests from your frontend application.
