# Environment Variables Guide

## TL;DR: Do I Need to Set Env Vars?

**No, the system will start without any .env file.** All services have sensible defaults.

## What Works Without Env Vars?

### ✅ Fully Functional (No Env Vars Needed)
- **User Auth Service** - Uses generated JWT keys from `./keys/`
- **Analytics Service** - Runs in TEST_MODE (no external API)
- **Notifications Service** - Uses default JWT secret
- **Feature Flags Service** - Uses Unleash (configured in docker-compose)
- **GraphQL Gateway** - Connects to all services
- **Frontend** - Connects to gateway
- **Postgres, Redis, Unleash** - All configured

### ⚠️ Limited Functionality (Defaults Used)
- **Billing Service** - Uses placeholder Stripe keys
  - Service starts ✅
  - Real payments won't work ❌
  - Need real keys for: subscriptions, webhooks, checkout
  
- **LLM Gateway Service** - Uses placeholder OpenAI key
  - Service starts ✅
  - LLM calls will fail ❌
  - Need real key for: AI chat, prompt execution

## When Do You Need Real Env Vars?

### For Development/Testing
**You don't need any.** The system is fully functional for:
- User registration/login
- RBAC and permissions
- Feature flags
- Real-time notifications
- Analytics tracking (test mode)
- GraphQL API exploration

### For Production Features
**You need real keys for:**

1. **Billing/Payments** → Set `STRIPE_API_KEY` and `STRIPE_WEBHOOK_SECRET`
2. **AI Features** → Set `OPENAI_API_KEY`
3. **Production Analytics** → Set `MIXPANEL_API_KEY` or `AMPLITUDE_API_KEY`
4. **Security** → Set strong `JWT_SECRET`

## How to Set Env Vars

### Option 1: Create .env File (Recommended)
```bash
# Copy the example
cp .env.example .env

# Edit with your values
nano .env  # or use your editor
```

Docker Compose automatically loads `.env` file.

### Option 2: Export in Shell
```bash
export STRIPE_API_KEY=sk_test_your_key
export OPENAI_API_KEY=sk-your_key
docker-compose up
```

### Option 3: Inline with Command
```bash
STRIPE_API_KEY=sk_test_key OPENAI_API_KEY=sk-key docker-compose up
```

## Current Configuration Status

### ✅ Already Configured
- JWT keys generated in `./keys/`
- Unleash tokens fixed in docker-compose
- Analytics TEST_MODE enabled
- All database connections configured
- All service-to-service URLs configured

### 📝 Optional Configuration
- Stripe keys (for billing)
- OpenAI key (for LLM)
- Production JWT secret
- External analytics providers

## Testing Without Real Keys

### Test Billing Flow
The billing service starts and accepts gRPC calls, but Stripe operations will fail gracefully. You can:
- Test the GraphQL API structure
- See the database schema
- Mock the responses in your frontend

### Test LLM Flow
The LLM gateway starts and accepts gRPC calls, but OpenAI calls will fail. You can:
- Test prompt loading from `/prompts`
- See the API structure
- Mock responses for frontend development

## Summary

**Your three startup errors were NOT caused by missing env vars:**
1. Unleash token format → Fixed in docker-compose
2. Missing JWT keys → Generated with OpenSSL
3. Analytics API key → Enabled TEST_MODE

**The system will start completely without a .env file.** You only need real API keys when you want to test actual billing or AI features.
