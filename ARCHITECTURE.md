# HAUNTED SAAS SKELETON - Architecture

## System Overview

```mermaid
graph TB
    subgraph "External Services"
        Stripe[Stripe API]
        OpenAI[OpenAI API]
        Unleash[Unleash Server]
    end
    
    subgraph "Client Layer"
        Browser[Web Browser]
        Mobile[Mobile App]
    end
    
    subgraph "Frontend"
        NextJS[Next.js Frontend<br/>Port 3000]
    end
    
    subgraph "API Gateway"
        GraphQL[GraphQL Gateway<br/>Port 4000<br/>JWT Auth]
    end
    
    subgraph "Microservices"
        Auth[User Auth Service<br/>Port 50051<br/>JWT, RBAC, Sessions]
        Billing[Billing Service<br/>Port 50052<br/>Stripe Integration]
        LLM[LLM Gateway Service<br/>Port 50053<br/>Prompt-as-Code]
        Notifications[Notifications Service<br/>Port 50054<br/>Socket.IO]
        Analytics[Analytics Service<br/>Port 50055<br/>Event Tracking]
        FeatureFlags[Feature Flags Service<br/>Port 50056<br/>Unleash Proxy]
    end
    
    subgraph "Data Layer"
        Postgres[(PostgreSQL<br/>Port 5432)]
        Redis[(Redis<br/>Port 6379)]
    end
    
    subgraph "Documentation"
        Docs[Docusaurus<br/>Port 3001]
    end
    
    Browser -->|HTTPS| NextJS
    Mobile -->|HTTPS| NextJS
    NextJS -->|GraphQL| GraphQL
    NextJS -->|Socket.IO| Notifications
    
    GraphQL -->|gRPC| Auth
    GraphQL -->|gRPC| Billing
    GraphQL -->|gRPC| LLM
    GraphQL -->|gRPC| Notifications
    GraphQL -->|gRPC| Analytics
    GraphQL -->|gRPC| FeatureFlags
    
    Auth -->|SQL| Postgres
    Auth -->|Cache| Redis
    Billing -->|SQL| Postgres
    Billing -->|Webhooks| Stripe
    LLM -->|API| OpenAI
    LLM -->|gRPC| Analytics
    Notifications -->|Cache| Redis
    Analytics -->|SQL| Postgres
    Analytics -->|Cache| Redis
    FeatureFlags -->|API| Unleash
    FeatureFlags -->|Cache| Redis
    FeatureFlags -->|gRPC| Auth
    FeatureFlags -->|gRPC| Analytics
    
    Unleash -->|SQL| Postgres
```

## Service Communication Matrix

| Service | Calls | Called By | Protocol |
|---------|-------|-----------|----------|
| **User Auth** | PostgreSQL, Redis | GraphQL Gateway, Feature Flags | gRPC |
| **Billing** | PostgreSQL, Stripe | GraphQL Gateway | gRPC, HTTP (webhooks) |
| **LLM Gateway** | OpenAI, Analytics | GraphQL Gateway | gRPC |
| **Notifications** | Redis | GraphQL Gateway, All Services | gRPC, Socket.IO |
| **Analytics** | PostgreSQL, Redis | All Services | gRPC |
| **Feature Flags** | Unleash, Redis, Auth, Analytics | GraphQL Gateway, All Services | gRPC |

## Data Flow Diagrams

### Authentication Flow

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant Auth
    participant Postgres
    participant Redis
    
    Client->>Gateway: POST /graphql (login mutation)
    Gateway->>Auth: Login(email, password) [gRPC]
    Auth->>Postgres: Query user by email
    Postgres-->>Auth: User record
    Auth->>Auth: Verify password (bcrypt)
    Auth->>Auth: Generate JWT (RS256)
    Auth->>Redis: Create session
    Redis-->>Auth: Session created
    Auth-->>Gateway: JWT + user data
    Gateway-->>Client: AuthPayload with token
    
    Note over Client: Store JWT in localStorage
    
    Client->>Gateway: Query with Authorization header
    Gateway->>Gateway: Verify JWT signature
    Gateway->>Auth: ValidateToken(jwt) [gRPC]
    Auth->>Redis: Check session
    Redis-->>Auth: Session valid
    Auth-->>Gateway: User + permissions
    Gateway->>Gateway: Execute query with user context
```

### Real-time Notification Flow

```mermaid
sequenceDiagram
    participant Client
    participant Notifications
    participant Redis
    participant OtherService
    
    Client->>Notifications: Connect Socket.IO (with JWT)
    Notifications->>Notifications: Verify JWT
    Notifications->>Redis: Create connection record
    Notifications->>Notifications: Join user + team rooms
    Notifications-->>Client: connection_ready event
    
    Note over OtherService: Event occurs (e.g., new message)
    
    OtherService->>Notifications: SendToUser(user_id, event) [gRPC]
    Notifications->>Notifications: Find user's connections
    Notifications->>Client: Emit event via Socket.IO
    Client->>Client: Display notification
```

### Feature Flag Evaluation Flow

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant FeatureFlags
    participant Redis
    participant Unleash
    participant Auth
    participant Analytics
    
    Client->>Gateway: Query isFeatureEnabled
    Gateway->>FeatureFlags: IsFeatureEnabled(name, context) [gRPC]
    
    FeatureFlags->>Redis: Check cache
    alt Cache Hit
        Redis-->>FeatureFlags: Cached result
    else Cache Miss
        FeatureFlags->>Auth: CheckPermission (if RBAC enabled)
        Auth-->>FeatureFlags: Permission result
        FeatureFlags->>Unleash: Evaluate flag
        Unleash-->>FeatureFlags: Evaluation result
        FeatureFlags->>Redis: Cache result (30s TTL)
    end
    
    FeatureFlags->>Analytics: TrackFlagEvaluation (async)
    FeatureFlags-->>Gateway: Boolean result
    Gateway-->>Client: Feature enabled/disabled
```

### Billing Subscription Flow

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant Billing
    participant Postgres
    participant Stripe
    
    Client->>Gateway: Mutation createCheckoutSession
    Gateway->>Billing: CreateCheckoutSession(plan_id) [gRPC]
    Billing->>Postgres: Get plan details
    Postgres-->>Billing: Plan data
    Billing->>Stripe: Create Checkout Session
    Stripe-->>Billing: Session URL
    Billing-->>Gateway: Checkout URL
    Gateway-->>Client: Redirect to Stripe
    
    Note over Client: User completes payment on Stripe
    
    Stripe->>Billing: Webhook: checkout.session.completed
    Billing->>Billing: Verify webhook signature
    Billing->>Postgres: Create subscription record
    Billing->>Stripe: Retrieve subscription details
    Stripe-->>Billing: Subscription data
    Billing->>Postgres: Update subscription
    Billing-->>Stripe: 200 OK
```

### LLM Prompt Execution Flow

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant LLM
    participant FileSystem
    participant OpenAI
    participant Analytics
    
    Client->>Gateway: Mutation callPrompt
    Gateway->>LLM: CallPrompt(path, variables) [gRPC]
    LLM->>FileSystem: Load prompt template
    FileSystem-->>LLM: Prompt content
    LLM->>LLM: Substitute variables
    LLM->>OpenAI: ChatCompletion API
    OpenAI-->>LLM: Response + token usage
    LLM->>Analytics: TrackUsage (async)
    LLM-->>Gateway: Response text + metadata
    Gateway-->>Client: LLM response
```

## Technology Stack

### Backend Services (Go)
- **Language**: Go 1.21+
- **gRPC**: google.golang.org/grpc
- **Database**: PostgreSQL 15 with GORM
- **Cache**: Redis 7
- **JWT**: github.com/golang-jwt/jwt/v5
- **Password**: golang.org/x/crypto/bcrypt
- **Logging**: go.uber.org/zap
- **Config**: github.com/spf13/viper

### Frontend (TypeScript)
- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **GraphQL**: Apollo Client or urql
- **Real-time**: Socket.IO client
- **Forms**: React Hook Form + Zod
- **Testing**: Jest + React Testing Library

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus + Grafana (planned)
- **Logging**: Structured JSON logs
- **Tracing**: OpenTelemetry (planned)

### External Services
- **Payments**: Stripe
- **LLM**: OpenAI
- **Feature Flags**: Unleash (self-hosted)

## Security Architecture

### Authentication & Authorization

```mermaid
graph LR
    A[Client Request] --> B{Has JWT?}
    B -->|No| C[Return 401]
    B -->|Yes| D[Verify Signature]
    D -->|Invalid| C
    D -->|Valid| E[Check Expiration]
    E -->|Expired| C
    E -->|Valid| F[Check Revocation]
    F -->|Revoked| C
    F -->|Active| G[Extract Claims]
    G --> H[Check Permissions]
    H -->|Denied| I[Return 403]
    H -->|Allowed| J[Process Request]
```

### Security Layers

1. **Transport Security**
   - TLS/HTTPS for all external communication
   - mTLS for service-to-service (optional)

2. **Authentication**
   - JWT with RS256 (asymmetric signing)
   - 24-hour token expiration
   - Session management in Redis
   - Token revocation list

3. **Authorization**
   - Role-Based Access Control (RBAC)
   - Granular permissions (resource:action)
   - Permission caching (5 min TTL)

4. **Data Protection**
   - bcrypt password hashing (cost 12)
   - PII encryption at rest
   - Sensitive data sanitization in logs

5. **Rate Limiting**
   - Account lockout (5 attempts = 30 min)
   - API rate limiting per user
   - Connection limits (10,000 concurrent)

6. **Input Validation**
   - Request validation at gateway
   - SQL injection prevention (parameterized queries)
   - XSS prevention (output encoding)

## Scalability Patterns

### Horizontal Scaling

All services are stateless and can scale horizontally:

```
┌─────────────┐
│   Load      │
│  Balancer   │
└──────┬──────┘
       │
   ┌───┴───┬───────┬───────┐
   │       │       │       │
┌──▼──┐ ┌──▼──┐ ┌──▼──┐ ┌──▼──┐
│Svc 1│ │Svc 2│ │Svc 3│ │Svc N│
└──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘
   │       │       │       │
   └───┬───┴───────┴───────┘
       │
   ┌───▼────┐
   │ Redis  │
   │Cluster │
   └────────┘
```

### Caching Strategy

1. **Session Cache** (Redis)
   - TTL: 24 hours (sliding window)
   - Invalidation: On logout or role change

2. **Permission Cache** (Redis)
   - TTL: 5 minutes
   - Invalidation: On role/permission change

3. **Feature Flag Cache** (Redis)
   - TTL: 30 seconds
   - Invalidation: On flag update

4. **Analytics Aggregation Cache** (Redis)
   - TTL: 5 minutes
   - Invalidation: Time-based

### Database Optimization

1. **Connection Pooling**
   - Max 25 connections per service
   - Connection reuse

2. **Indexes**
   - User email (unique)
   - Session user_id
   - Event timestamp
   - Subscription team_id

3. **Query Optimization**
   - Prepared statements
   - Batch inserts
   - Pagination

## Deployment Architecture

### Development
```
Docker Compose
├── PostgreSQL
├── Redis
├── Unleash
└── All Services (local build)
```

### Production (Kubernetes)
```
┌─────────────────────────────────────┐
│           Ingress Controller         │
│         (TLS Termination)            │
└────────────┬────────────────────────┘
             │
    ┌────────┴────────┐
    │                 │
┌───▼────┐      ┌────▼─────┐
│Frontend│      │ Gateway  │
│  Pod   │      │   Pod    │
└───┬────┘      └────┬─────┘
    │                │
    └────────┬───────┘
             │
    ┌────────┴────────────────────┐
    │                             │
┌───▼────┐  ┌────────┐  ┌────────▼┐
│ Auth   │  │Billing │  │  Other  │
│  Pod   │  │  Pod   │  │Services │
└───┬────┘  └───┬────┘  └────┬────┘
    │           │            │
    └───────┬───┴────────────┘
            │
    ┌───────┴────────┐
    │                │
┌───▼────┐    ┌─────▼──┐
│Postgres│    │ Redis  │
│Cluster │    │Cluster │
└────────┘    └────────┘
```

## Monitoring & Observability

### Metrics (Prometheus)
- Request rate, latency, errors
- Connection counts
- Cache hit rates
- Database query performance
- Token generation/validation rates

### Logging (Structured JSON)
- Request/response logs
- Error logs with stack traces
- Audit logs (auth events)
- Performance logs (slow queries)

### Tracing (OpenTelemetry)
- Distributed request tracing
- Service dependency mapping
- Performance bottleneck identification

### Alerting
- Service health
- Error rate thresholds
- Database connection issues
- Cache failures
- External API failures

## Disaster Recovery

### Backup Strategy
- **Database**: Daily automated backups
- **Redis**: AOF persistence + snapshots
- **Configuration**: Version controlled
- **Keys**: Secure key management service

### Recovery Procedures
1. Database restore from backup
2. Redis restore from snapshot
3. Service redeployment
4. Configuration restoration
5. Verification testing

---

**Architecture Version**: 1.0
**Last Updated**: 2025-01-11
**Status**: Foundation complete, ready for implementation
