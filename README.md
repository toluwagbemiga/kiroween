# HAUNTED SAAS SKELETON

A production-grade, microservices-based SaaS platform built with Go, Next.js, and modern DevOps practices.

## Architecture

This is a monorepo containing:

- **Backend Services** (Go microservices with gRPC)
  - `user-auth-service` - Authentication, authorization, and RBAC
  - `billing-service` - Stripe integration for subscriptions
  - `llm-gateway-service` - LLM proxy with Prompt-as-Code
  - `notifications-service` - Real-time Socket.IO communication
  - `analytics-service` - Event tracking and analytics
  - `feature-flags-service` - Feature flag management with Unleash

- **Frontend** (Next.js 14+ with TypeScript)
  - Modern React application with Design System
  - GraphQL API Gateway integration
  - WCAG 2.1 Level AA accessibility

- **Documentation** (Docusaurus)
  - API reference auto-generated from OpenAPI specs
  - Tutorials, how-to guides, and explanations

- **Demo Sandbox** (Docker Compose)
  - Fully functional local environment
  - Pre-populated with realistic test data

## Quick Start

```bash
# Start the complete demo environment
docker-compose up -d

# Access services:
# - Frontend: http://localhost:3000
# - Docs: http://localhost:3001
# - Unleash UI: http://localhost:4242
# - GraphQL Playground: http://localhost:4000/graphql
```

## Development

See individual service READMEs in `/app/services/` for development instructions.

## CI/CD

GitHub Actions workflows handle:
- Service builds and tests
- Docker image publishing
- Documentation deployment
- Demo sandbox validation

## License

MIT
