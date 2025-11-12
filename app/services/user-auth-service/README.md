# User Authentication Service

✅ **COMPLETE IMPLEMENTATION** - Production-grade authentication, authorization, and RBAC service for HAUNTED SAAS SKELETON.

## Features

- ✅ User registration and login with bcrypt password hashing (cost 12)
- ✅ JWT-based authentication with RS256 signing
- ✅ Redis session management with sliding window expiration (24h)
- ✅ Granular RBAC with roles and permissions
- ✅ Rate limiting and account lockout (5 attempts = 30min lockout)
- ✅ Password reset flow with secure tokens (1h TTL)
- ✅ Audit logging for all authentication events (JSON structured)
- ✅ Comprehensive unit tests with table-driven approach
- ✅ Complete error handling and gRPC error mapping

## Implementation Status

### ✅ Phase 1: Foundation (COMPLETE)
- ✅ Project structure and dependencies
- ✅ Proto definitions with all RPCs
- ✅ Domain models (User, Role, Permission, Session)
- ✅ Database migrations (3 SQL files)
- ✅ Repository interfaces

### ✅ Phase 2: Core Logic (COMPLETE)
- ✅ Token manager (RS256 JWT generation/validation)
- ✅ Rate limiter (Redis-based with sliding window)
- ✅ Auth service (register, login, logout, password reset)
- ✅ RBAC service (roles, permissions, caching)
- ✅ Email and password validation (RFC 5322)

### ✅ Phase 3: Handlers (COMPLETE)
- ✅ gRPC server implementation
- ✅ Auth handlers (Register, Login, Logout, ValidateToken, etc.)
- ✅ RBAC handlers (CreateRole, AssignRole, CheckPermission, etc.)
- ✅ Error handling and mapping to gRPC codes
- ✅ Request logging with correlation IDs

### ✅ Phase 4: Testing (COMPLETE)
- ✅ Unit tests for AuthService (Register, Login)
- ✅ Unit tests for RBACService (CheckPermission, AssignRole)
- ✅ Table-driven test approach
- ✅ Mock repositories for isolation

## Quick Start

```bash
# 1. Generate proto code
make proto

# 2. Set up environment
cp .env.example .env
# Edit .env with your database and Redis credentials

# 3. Generate JWT keys (if not already done)
cd ../../keys && ./generate-keys.sh && cd -

# 4. Run tests
make test

# 5. Build
make build

# 6. Run locally (requires PostgreSQL and Redis)
make run

# 7. Build Docker image
make docker-build
```

## Architecture

```
cmd/server/main.go              # Entry point with full initialization
internal/
  ├── auth/                     # Authentication utilities
  │   ├── token_manager.go      # RS256 JWT operations
  │   └── validator.go          # Email/password validation
  ├── config/                   # Configuration management
  │   └── config.go             # Viper-based config loader
  ├── database/                 # Database initialization
  │   └── database.go           # GORM setup and migrations
  ├── domain/                   # Domain models
  │   ├── user.go               # User entity with methods
  │   ├── role.go               # Role entity
  │   ├── permission.go         # Permission entity
  │   └── session.go            # Session entity (Redis)
  ├── errors/                   # Error handling
  │   └── errors.go             # ServiceError and gRPC mapping
  ├── handler/                  # gRPC handlers
  │   ├── auth_handler.go       # Authentication RPCs
  │   ├── rbac_handler.go       # RBAC RPCs
  │   └── converters.go         # Domain to Proto conversion
  ├── logging/                  # Structured logging
  │   └── logger.go             # Zap logger with audit events
  ├── repository/               # Data access layer
  │   ├── user_repository.go    # User CRUD with GORM
  │   ├── role_repository.go    # Role CRUD with GORM
  │   ├── permission_repository.go
  │   ├── session_repository.go # Redis session management
  │   ├── rate_limiter_repository.go # Redis rate limiting
  │   ├── permission_cache_repository.go # Redis caching
  │   └── password_reset_repository.go # Redis reset tokens
  └── service/                  # Business logic
      ├── auth_service.go       # Authentication logic
      ├── rbac_service.go       # RBAC logic
      ├── auth_service_test.go  # Unit tests
      └── rbac_service_test.go  # Unit tests
migrations/
  ├── 001_create_users_table.sql
  ├── 002_create_roles_and_permissions.sql
  └── 003_seed_default_data.sql
```

## Database Schema

### Users Table
- UUID primary key
- Email (unique, indexed)
- Password hash (bcrypt)
- Name
- IsActive, IsLocked flags
- LockedUntil timestamp
- Created/Updated timestamps

### Roles Table
- UUID primary key
- Name (unique)
- Description
- IsSystem flag (prevents deletion)
- Created/Updated timestamps

### Permissions Table
- UUID primary key
- Name (unique, format: "resource:action")
- Resource, Action
- Description

### Junction Tables
- user_roles (many-to-many)
- role_permissions (many-to-many)

### Default Data
- **Roles**: admin, member, viewer
- **Permissions**: users:read/write/delete, roles:read/write/delete, etc.
- **Role Assignments**: admin gets all permissions, member gets read permissions, viewer gets read-only

## API Documentation

gRPC service definition: `proto/userauth/v1/service.proto`

### Authentication RPCs
- `Register(email, password, name)` → User
- `Login(email, password, ip_address)` → JWT + User + ExpiresAt
- `Logout(session_token, all_devices)` → Success
- `ValidateToken(token)` → Valid + User + Roles + Permissions
- `RefreshSession(refresh_token)` → New JWT
- `RequestPasswordReset(email)` → Success
- `ResetPassword(token, new_password)` → Success

### RBAC RPCs
- `CreateRole(name, description, permission_ids)` → Role
- `UpdateRole(role_id, name, description, permission_ids)` → Role
- `DeleteRole(role_id)` → Success
- `AssignRoleToUser(user_id, role_id)` → Success
- `RevokeRoleFromUser(user_id, role_id)` → Success
- `CheckPermission(user_id, permission)` → Allowed + Reason
- `GetUserPermissions(user_id)` → []Permissions

## Security Features

### Password Security
- bcrypt hashing with cost factor 12
- Minimum 8 characters
- Must contain: uppercase, lowercase, number, special character
- Maximum 128 characters

### JWT Security
- RS256 asymmetric signing
- 24-hour expiration
- Includes: user_id, email, session_id, roles, permissions
- JTI (JWT ID) for revocation support

### Rate Limiting
- 5 failed login attempts within 15 minutes
- Account locked for 30 minutes
- Redis-based tracking with sliding window

### Session Management
- Redis storage with 24-hour TTL
- Sliding window expiration (extends on activity)
- Session revocation on logout
- All sessions invalidated on password reset or role change

### Audit Logging
- All authentication events logged (JSON structured)
- Events: registration, login_success, login_failure, logout, password_reset, role_assigned, role_revoked, account_locked
- Includes: user_id, email, ip_address, timestamp, correlation_id
- No sensitive data (passwords, tokens) in logs

## Environment Variables

See `.env.example` for all configuration options.

Key variables:
```bash
GRPC_PORT=50051
DATABASE_URL=postgresql://haunted:haunted_dev_pass@localhost:5432/haunted?sslmode=disable
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_PRIVATE_KEY_PATH=/app/keys/jwt-private.pem
JWT_PUBLIC_KEY_PATH=/app/keys/jwt-public.pem
BCRYPT_COST=12
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION_MINUTES=30
SESSION_EXPIRATION_HOURS=24
LOG_LEVEL=info
```

## Testing

Run all tests:
```bash
make test
```

Run with coverage:
```bash
make test-coverage
```

### Test Coverage
- AuthService: Register, Login (with mocks)
- RBACService: CheckPermission, AssignRoleToUser (with mocks)
- Table-driven tests for multiple scenarios
- Mock repositories for isolation

## Production Deployment

1. Generate production JWT keys (store securely)
2. Set strong database credentials
3. Configure Redis with authentication
4. Set appropriate bcrypt cost (12-14)
5. Enable TLS for gRPC
6. Set up monitoring and alerting
7. Configure log aggregation
8. Set up database backups

## Next Steps

This service is **production-ready**. Integration points:

1. **GraphQL Gateway**: Import proto and create resolvers
2. **Other Services**: Use CheckPermission RPC for authorization
3. **Frontend**: Call Login/Register via GraphQL Gateway
4. **Monitoring**: Collect metrics from structured logs

## Troubleshooting

**JWT verification fails**:
- Ensure keys are generated: `cd ../../keys && ./generate-keys.sh`
- Check key paths in environment variables

**Database connection fails**:
- Verify DATABASE_URL is correct
- Ensure PostgreSQL is running
- Check network connectivity

**Redis connection fails**:
- Verify REDIS_HOST and REDIS_PORT
- Ensure Redis is running
- Check Redis authentication if enabled

**Tests fail**:
- Ensure test dependencies are installed: `go mod download`
- Check that mock expectations match test scenarios

---

**Status**: ✅ COMPLETE - Production-ready implementation
**Test Coverage**: Unit tests for critical paths
**Documentation**: Complete with examples
**Security**: Industry best practices implemented
