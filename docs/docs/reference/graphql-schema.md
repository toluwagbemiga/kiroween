---
sidebar_position: 2
title: GraphQL Schema
description: Complete GraphQL schema reference for Haunted SaaS API
---

# GraphQL Schema Reference

This document provides a complete reference for the Haunted SaaS GraphQL API schema.

## Authentication

All authenticated requests require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Core Types

### User

```graphql
type User {
  id: ID!
  email: String!
  name: String
  roles: [Role!]!
  subscription: Subscription
  createdAt: Time!
  updatedAt: Time!
}
```

### Role

```graphql
type Role {
  id: ID!
  name: String!
  permissions: [Permission!]!
}
```

### Subscription

```graphql
type Subscription {
  id: ID!
  userId: ID!
  plan: Plan!
  status: SubscriptionStatus!
  currentPeriodStart: Time!
  currentPeriodEnd: Time!
}

enum SubscriptionStatus {
  ACTIVE
  CANCELED
  PAST_DUE
  TRIALING
}
```

## Queries

### Authentication Queries

```graphql
type Query {
  # Get current authenticated user
  me: User
  
  # Validate a token
  validateToken(token: String!): Boolean!
}
```

### User Queries

```graphql
type Query {
  # Get user by ID (requires admin role)
  user(id: ID!): User
  
  # List all users (requires admin role)
  users(limit: Int, offset: Int): [User!]!
}
```

### Billing Queries

```graphql
type Query {
  # Get available plans
  plans: [Plan!]!
  
  # Get user's subscription
  subscription(userId: ID!): Subscription
}
```

## Mutations

### Authentication Mutations

```graphql
type Mutation {
  # Login with email and password
  login(email: String!, password: String!): AuthPayload!
  
  # Register new user
  register(input: RegisterInput!): AuthPayload!
  
  # Refresh access token
  refreshToken(token: String!): AuthPayload!
  
  # Logout
  logout: Boolean!
}

type AuthPayload {
  token: String!
  user: User!
}

input RegisterInput {
  email: String!
  password: String!
  name: String
}
```

### Billing Mutations

```graphql
type Mutation {
  # Create subscription
  createSubscription(input: CreateSubscriptionInput!): Subscription!
  
  # Cancel subscription
  cancelSubscription(subscriptionId: ID!): Subscription!
  
  # Update payment method
  updatePaymentMethod(input: UpdatePaymentMethodInput!): Boolean!
}
```

### Analytics Mutations

```graphql
type Mutation {
  # Track event
  trackEvent(input: TrackEventInput!): Boolean!
  
  # Identify user
  identifyUser(input: IdentifyUserInput!): Boolean!
}

input TrackEventInput {
  eventName: String!
  properties: JSON
}
```

### Notification Mutations

```graphql
type Mutation {
  # Send notification
  sendNotification(input: SendNotificationInput!): Boolean!
}

input SendNotificationInput {
  userId: ID!
  message: String!
  type: NotificationType!
}

enum NotificationType {
  INFO
  SUCCESS
  WARNING
  ERROR
}
```

### AI Mutations

```graphql
type Mutation {
  # Call AI prompt
  callPrompt(name: String!, variables: JSON!): PromptResponse!
}

type PromptResponse {
  content: String!
  tokensUsed: Int!
  cost: Float!
}
```

## Subscriptions (WebSocket)

```graphql
type Subscription {
  # Subscribe to notifications
  notifications(userId: ID!): Notification!
}

type Notification {
  id: ID!
  userId: ID!
  message: String!
  type: NotificationType!
  createdAt: Time!
}
```

## Scalar Types

```graphql
scalar Time
scalar JSON
```

## Error Handling

All errors follow this structure:

```json
{
  "errors": [
    {
      "message": "Error description",
      "extensions": {
        "code": "ERROR_CODE",
        "field": "fieldName"
      }
    }
  ]
}
```

### Common Error Codes

- `UNAUTHENTICATED`: Missing or invalid authentication
- `FORBIDDEN`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `VALIDATION_ERROR`: Input validation failed
- `INTERNAL_ERROR`: Server error

## Rate Limiting

API requests are rate limited:
- **Authenticated**: 1000 requests per hour
- **Unauthenticated**: 100 requests per hour

Rate limit headers:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640000000
```

---

For interactive exploration, use the GraphQL Playground at `http://localhost:8080/graphql`
