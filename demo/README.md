# Demo Sandbox

This directory contains scripts and configuration for the HAUNTED SAAS SKELETON demo environment.

## Quick Start

```bash
# From project root
docker-compose up -d

# Wait for services to be healthy (about 30 seconds)
docker-compose ps

# Generate demo data
cd demo
npm install
npm run generate
```

## What's Included

- **Demo Users**: Pre-configured users with different roles
- **Sample Plans**: Free, Pro, and Enterprise subscription tiers
- **Feature Flags**: Example flags with various strategies
- **Analytics Data**: 1000+ sample events for testing
- **LLM Prompts**: Ready-to-use prompt templates

## Demo Credentials

- **Admin**: admin@haunted.dev / Admin123!
- **User**: user@haunted.dev / User123!
- **Viewer**: viewer@haunted.dev / Viewer123!

## Service URLs

- Frontend: http://localhost:3000
- Documentation: http://localhost:3001
- GraphQL Playground: http://localhost:4000/graphql
- Unleash UI: http://localhost:4242 (admin/unleash4all)

## Resetting the Demo

```bash
# Stop and remove all data
docker-compose down -v

# Start fresh
docker-compose up -d
cd demo && npm run generate
```
