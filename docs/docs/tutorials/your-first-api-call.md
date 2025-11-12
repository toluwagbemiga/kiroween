---
sidebar_position: 2
title: Your First API Call
description: Learn how to make your first GraphQL API call to Haunted SaaS
---

# Your First API Call

Now that you have Haunted SaaS running, let's make your first API call using the GraphQL API.

## Using GraphQL Playground

The easiest way to explore the API is through the built-in GraphQL Playground.

### Step 1: Open GraphQL Playground

Navigate to: http://localhost:8080/graphql

You'll see an interactive GraphQL IDE where you can write queries and see results in real-time.

### Step 2: Login to Get a Token

First, we need to authenticate. Run this mutation:

```graphql
mutation Login {
  login(email: "user@haunted-saas.com", password: "user123") {
    token
    user {
      id
      email
      name
    }
  }
}
```

**Response:**
```json
{
  "data": {
    "login": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "email": "user@haunted-saas.com",
        "name": "Demo User"
      }
    }
  }
}
```

### Step 3: Set Authorization Header

Copy the token from the response and add it to the HTTP Headers section at the bottom of the playground:

```json
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Step 4: Query Your Profile

Now you can make authenticated requests:

```graphql
query GetMe {
  me {
    id
    email
    name
    roles {
      name
      permissions {
        name
      }
    }
    subscription {
      plan {
        name
        price
      }
      status
    }
  }
}
```

**Response:**
```json
{
  "data": {
    "me": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@haunted-saas.com",
      "name": "Demo User",
      "roles": [
        {
          "name": "user",
          "permissions": [
            { "name": "read:own_profile" },
            { "name": "update:own_profile" }
          ]
        }
      ],
      "subscription": {
        "plan": {
          "name": "Free",
          "price": 0
        },
        "status": "ACTIVE"
      }
    }
  }
}
```

## Using cURL

You can also make API calls using cURL:

```bash
# Login
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { login(email: \"user@haunted-saas.com\", password: \"user123\") { token user { id email } } }"
  }'

# Get profile (replace TOKEN with your actual token)
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "query": "query { me { id email name } }"
  }'
```

## Using JavaScript/TypeScript

```typescript
const GRAPHQL_ENDPOINT = 'http://localhost:8080/graphql';

// Login
async function login(email: string, password: string) {
  const response = await fetch(GRAPHQL_ENDPOINT, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query: `
        mutation Login($email: String!, $password: String!) {
          login(email: $email, password: $password) {
            token
            user {
              id
              email
              name
            }
          }
        }
      `,
      variables: { email, password }
    })
  });

  const data = await response.json();
  return data.data.login;
}

// Get profile
async function getProfile(token: string) {
  const response = await fetch(GRAPHQL_ENDPOINT, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      query: `
        query {
          me {
            id
            email
            name
          }
        }
      `
    })
  });

  const data = await response.json();
  return data.data.me;
}

// Usage
const { token, user } = await login('user@haunted-saas.com', 'user123');
console.log('Logged in as:', user.email);

const profile = await getProfile(token);
console.log('Profile:', profile);
```

## Common Operations

### Track an Analytics Event

```graphql
mutation TrackEvent {
  trackEvent(input: {
    eventName: "button_clicked"
    properties: {
      button_name: "signup"
      page: "homepage"
    }
  })
}
```

### Send a Notification

```graphql
mutation SendNotification {
  sendNotification(input: {
    userId: "123e4567-e89b-12d3-a456-426614174000"
    message: "Welcome to Haunted SaaS!"
    type: INFO
  })
}
```

### Call an AI Prompt

```graphql
mutation CallAI {
  callPrompt(
    name: "v1/support-chatbot"
    variables: {
      user_message: "How do I upgrade my plan?"
    }
  ) {
    content
    tokensUsed
    cost
  }
}
```

### Check Feature Flag

```graphql
query CheckFeature {
  isFeatureEnabled(
    userId: "123e4567-e89b-12d3-a456-426614174000"
    featureKey: "new-dashboard"
  )
}
```

## Error Handling

GraphQL returns errors in a structured format:

```json
{
  "errors": [
    {
      "message": "Invalid credentials",
      "extensions": {
        "code": "UNAUTHENTICATED"
      }
    }
  ]
}
```

Always check for the `errors` field in responses:

```typescript
const response = await fetch(GRAPHQL_ENDPOINT, { /* ... */ });
const data = await response.json();

if (data.errors) {
  console.error('GraphQL errors:', data.errors);
  throw new Error(data.errors[0].message);
}

return data.data;
```

## Next Steps

- 📖 [GraphQL Schema Reference](../reference/graphql-schema) - Complete API documentation
- 🏗️ [Building Your First App](./building-your-first-app) - Create a client application
- 🔐 [Adding Authentication](./adding-authentication) - Implement user auth

---

**Congratulations!** You've made your first API call to Haunted SaaS. You're ready to build amazing applications!
