# GraphQL API Gateway

Single external API entry point that translates GraphQL to internal gRPC calls.

## Features

- ✅ GraphQL schema stitching all services
- ✅ JWT authentication middleware
- ✅ gRPC client connections to all services
- ✅ Error handling and mapping
- ✅ Request logging and tracing
- ✅ CORS configuration

## Architecture

```
cmd/server/main.go
internal/
  ├── schema/
  │   ├── schema.graphql        # Main schema
  │   ├── auth.graphql          # Auth types
  │   ├── billing.graphql       # Billing types
  │   └── ...
  ├── resolvers/
  │   ├── auth_resolver.go
  │   ├── billing_resolver.go
  │   └── ...
  ├── clients/                  # gRPC clients
  │   ├── auth_client.go
  │   ├── billing_client.go
  │   └── ...
  ├── middleware/
  │   ├── auth.go               # JWT verification
  │   └── logging.go
  └── config/
```

## Technology Stack

- **Framework**: gqlgen (Go GraphQL server)
- **gRPC**: Client connections to all services
- **JWT**: Token verification with public key
- **HTTP**: Chi router for middleware

## GraphQL Schema Example

```graphql
type Mutation {
  # Auth
  register(input: RegisterInput!): AuthPayload!
  login(input: LoginInput!): AuthPayload!
  logout: Boolean!
  
  # Billing
  createCheckoutSession(planId: ID!): CheckoutSession!
  cancelSubscription: Subscription!
  
  # Analytics
  trackEvent(input: TrackEventInput!): Boolean!
  identifyUser(input: IdentifyUserInput!): Boolean!
}

type Query {
  me: User!
  subscription: Subscription
  plans: [Plan!]!
  isFeatureEnabled(featureName: String!): Boolean!
}
```

## Environment Variables

```bash
PORT=4000
USER_AUTH_SERVICE=user-auth-service:50051
BILLING_SERVICE=billing-service:50052
LLM_GATEWAY_SERVICE=llm-gateway-service:50053
NOTIFICATIONS_SERVICE=notifications-service:50054
ANALYTICS_SERVICE=analytics-service:50055
FEATURE_FLAGS_SERVICE=feature-flags-service:50056
JWT_PUBLIC_KEY_PATH=/app/keys/jwt-public.pem
ALLOWED_ORIGINS=http://localhost:3000
```

## Implementation Steps

1. Define GraphQL schema for all services
2. Generate resolvers with gqlgen
3. Implement gRPC clients for each service
4. Add JWT authentication middleware
5. Implement resolvers (call gRPC services)
6. Add error handling and logging
7. Test with GraphQL Playground

## Authentication Flow

1. Client sends JWT in Authorization header
2. Middleware verifies JWT with public key
3. Extract user ID and permissions from claims
4. Pass user context to resolvers
5. Resolvers include user ID in gRPC calls

## Error Handling

Map gRPC errors to GraphQL errors:
- `codes.Unauthenticated` → `UNAUTHENTICATED`
- `codes.PermissionDenied` → `FORBIDDEN`
- `codes.NotFound` → `NOT_FOUND`
- `codes.InvalidArgument` → `BAD_USER_INPUT`

## Next Steps

1. Install gqlgen and generate boilerplate
2. Define complete GraphQL schema
3. Implement gRPC client connections
4. Add authentication middleware
5. Implement all resolvers
6. Test with Playground
