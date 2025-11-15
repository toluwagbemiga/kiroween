# Auth Service Fix - Migrations Missing

## Problem

The user-auth-service Docker container doesn't have the migrations folder, causing registration to fail with "internal error".

## Root Cause

The Dockerfile wasn't copying the `migrations/` folder into the final container image.

## Fix Applied

Updated `app/services/user-auth-service/Dockerfile`:

```dockerfile
# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app  # Changed from /root/

# Copy the binary from builder
COPY --from=builder /app/user-auth-service .

# Copy migrations folder (NEW!)
COPY --from=builder /app/migrations ./migrations

# Copy JWT keys directory (will be mounted as volume)
RUN mkdir -p /app/keys

EXPOSE 50051

CMD ["./user-auth-service"]
```

## How to Apply

### Rebuild and Restart
```bash
# Stop services
docker-compose down

# Rebuild user-auth-service
docker-compose build user-auth-service

# Start all services
docker-compose up
```

## What Will Happen

When the service starts, it will:
1. Connect to PostgreSQL
2. Find the migrations folder at `/app/migrations`
3. Run all 3 migration files:
   - `001_create_users_table.sql`
   - `002_create_roles_and_permissions.sql`
   - `003_seed_default_data.sql`
4. Create all necessary tables
5. Seed default roles (admin, user, etc.)

## Verification

After restart, check the logs:
```bash
docker-compose logs -f user-auth-service
```

Should see:
```
✓ Database connected
✓ Migrations completed
✓ Redis connected
✓ Token manager initialized
🚀 User Auth Service started
```

Then check the database:
```bash
docker-compose exec postgres psql -U haunted -d haunted -c "\dt"
```

Should see tables:
- users
- roles
- permissions
- role_permissions
- user_roles
- sessions
- password_resets
- rate_limits

## Test Registration

After fix, you can register with:
- **Email**: test@example.com
- **Password**: Test123! (uppercase, lowercase, number, special char)
- **Name**: Test User

The registration should succeed and you'll be logged in automatically!
