---
inclusion: always
---

# HAUNTED SAAS SKELETON - Project Standards

## Architecture

**Monorepo Structure:**
- `/app` - Backend services (Go) and frontend (Next.js)
- `/docs` - Docusaurus documentation site
- `/demo` - Docker sandbox and demo data scripts
- `/prompts` - Prompt-as-Code templates for LLM interactions

**Service Communication:**
- Internal: gRPC for all service-to-service calls
- External: Single GraphQL API Gateway using `gqlgen`
- Real-time: Socket.IO (not plain WebSockets) for HTTP long-polling fallback

## Technology Stack

**Backend:**
- Go 1.21+ for all services
- PostgreSQL with `pgx` driver and `GORM` ORM
- Redis for caching and session management
- gRPC for internal APIs

**Frontend:**
- Next.js 14+ with TypeScript
- React components from shared Design System
- WCAG 2.1 Level AA accessibility compliance required

**Documentation:**
- Docusaurus with MDX
- Follow Diátaxis framework: tutorials, how-to guides, reference, explanation

## Code Standards

**Go Services:**
- Use table-driven tests for all unit tests
- Target >80% test coverage
- Generate OpenAPI 3.0 specs from gRPC definitions and code comments
- Follow standard Go project layout: `/cmd`, `/internal`, `/pkg`

**TypeScript/React:**
- Build all UI components as part of reusable Design System
- Components must work in both Next.js app and Docusaurus MDX
- Ensure accessibility compliance in all generated code

**API Documentation:**
- Auto-generate `openapi.yaml` from gRPC specs
- Consume OpenAPI specs in Docusaurus for API reference docs

## Core Services

**llm-gateway-service:**
- Single interface to external LLMs (OpenAI, etc.)
- Load prompts from `/prompts` directory as Prompt-as-Code
- Implement prompt versioning and template management

**user-auth-service:**
- Granular RBAC with roles and permissions in core user model
- JWT-based authentication
- Permission checks via gRPC

**feature-flag-service:**
- Self-hosted abstraction over Unleash
- Provide `IsFeatureEnabled(userId, featureKey)` gRPC endpoint
- Support gradual rollouts and A/B testing

**notifications-service:**
- Manage all real-time communication
- Use Socket.IO for WebSocket with long-polling fallback
- Support multiple notification channels (in-app, email, push)

**analytics-service:**
- Event-based analytics abstraction layer
- Provide `TrackEvent` and `IdentifyUser` endpoints
- Compatible with Mixpanel/Amplitude patterns

**recommendation-service:**
- Personalized recommendation engine
- Abstract interface for pluggable recommendation algorithms

## DevOps & CI/CD

**Containerization:**
- Every service must have a `Dockerfile`
- Provide `docker-compose.yml` for local development and demo sandbox

**GitHub Actions Workflows:**
- `app-ci-cd.yml` - Build, test, deploy Next.js and Go services
- `docs-deploy.yml` - Build and deploy Docusaurus to Vercel/Read the Docs
- `demo-sandbox-build.yml` - Spin up Docker sandbox with demo data

**Demo Environment:**
- Generate `data-generation-script.js` for realistic test data
- Populate all services with demo users, events, and content
- Ensure demo is fully functional without external dependencies


