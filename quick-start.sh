#!/bin/bash

# HAUNTED SAAS SKELETON - Quick Start Script
# This script sets up the development environment

set -e

echo "🎃 HAUNTED SAAS SKELETON - Quick Start"
echo "======================================"
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."

command -v docker >/dev/null 2>&1 || { echo "❌ Docker is required but not installed. Aborting." >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose is required but not installed. Aborting." >&2; exit 1; }
command -v openssl >/dev/null 2>&1 || { echo "❌ OpenSSL is required but not installed. Aborting." >&2; exit 1; }

echo "✅ All prerequisites met"
echo ""

# Generate JWT keys
echo "🔐 Generating JWT keys..."
cd keys
if [ ! -f jwt-private.pem ]; then
    ./generate-keys.sh
else
    echo "⚠️  JWT keys already exist, skipping generation"
fi
cd ..
echo ""

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cat > .env << EOF
# Stripe (use test keys)
STRIPE_API_KEY=sk_test_placeholder
STRIPE_WEBHOOK_SECRET=whsec_placeholder

# OpenAI (optional for LLM Gateway)
OPENAI_API_KEY=sk-placeholder

# JWT Secret
JWT_SECRET=dev-secret-key-change-in-production
EOF
    echo "✅ .env file created (update with real keys if needed)"
else
    echo "⚠️  .env file already exists, skipping creation"
fi
echo ""

# Start infrastructure services
echo "🚀 Starting infrastructure services..."
docker-compose up -d postgres redis unleash-db unleash
echo ""

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy (this may take 30-60 seconds)..."
sleep 10

# Check PostgreSQL
echo "   Checking PostgreSQL..."
timeout 60 bash -c 'until docker-compose exec -T postgres pg_isready -U haunted > /dev/null 2>&1; do sleep 2; done' || {
    echo "❌ PostgreSQL failed to start"
    exit 1
}
echo "   ✅ PostgreSQL is ready"

# Check Redis
echo "   Checking Redis..."
timeout 30 bash -c 'until docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; do sleep 2; done' || {
    echo "❌ Redis failed to start"
    exit 1
}
echo "   ✅ Redis is ready"

# Check Unleash
echo "   Checking Unleash..."
timeout 90 bash -c 'until curl -sf http://localhost:4242/health > /dev/null 2>&1; do sleep 3; done' || {
    echo "❌ Unleash failed to start"
    exit 1
}
echo "   ✅ Unleash is ready"

echo ""
echo "✅ Infrastructure is ready!"
echo ""

# Display next steps
echo "📚 Next Steps:"
echo ""
echo "1. Implement services (see IMPLEMENTATION_GUIDE.md)"
echo "   Start with: cd app/services/user-auth-service"
echo ""
echo "2. Or start all services (if already implemented):"
echo "   docker-compose up -d"
echo ""
echo "3. Access services:"
echo "   - Unleash UI:  http://localhost:4242 (admin/unleash4all)"
echo "   - PostgreSQL:  localhost:5432 (haunted/haunted_dev_pass)"
echo "   - Redis:       localhost:6379"
echo ""
echo "4. Generate demo data (after services are running):"
echo "   cd demo && npm install && npm run generate"
echo ""
echo "📖 Documentation:"
echo "   - PROJECT_STATUS.md - Current project status"
echo "   - IMPLEMENTATION_GUIDE.md - Step-by-step implementation"
echo "   - README.md - Project overview"
echo ""
echo "🎃 Happy building!"
