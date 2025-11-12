# User Auth Service - Quick Reference

## File Count: 35+ files created

## Core Components

### 1. Domain Models (4 files)
- `internal/domain/user.go` - User entity with GetPermissions(), GetRoleNames()
- `internal/domain/role.go` - Role entity with CanDelete()
- `internal/domain/permission.go` - Permission entity with NewPermission()
- `internal/domain/session.go` - Session entity for Redis

### 2. Repositories (8 files)
- `internal/repository/user_repository.go` - User CRUD + role assignment
- `internal/repository/role_repository.go` - Role CRUD + permission management
- `internal/repository/permission_repository.go` - Permission queries
- `internal/repository/session_repository.go` - Session management (Redis)
- `internal/repository/rate_limiter_repository.go` - Rate limiting (Redis)
- `internal/repository/permission_cache_repository.go` - Permission caching (Redis)
- `internal/repository/password_reset_repository.go` - Reset tokens (Redis)
- `internal/repository/errors.go` - Repository error types

### 3. Services (4 files)
- `internal/service/auth_service.go` - Register, Login, Logout, Password Reset
- `internal/service/rbac_service.go` - Role/Permission management
- `internal/service/auth_service_test.go` - Auth unit tests
- `internal/service/rbac_service_test.go` - RBAC unit tests

### 4. Handlers (3 files)
- `internal/handler/auth_handler.go` - Auth gRPC handlers
- `internal/handler/rbac_handler.go` - RBAC gRPC handlers
- `internal/handler/converters.go` - Domain to Proto conversion

### 5. Infrastructure (8 files)
- `internal/auth/token_manager.go` - RS256 JWT operations
- `internal/auth/validator.go` - Email/password validation
- `internal/config/config.go` - Configuration management
- `internal/database/database.go` - Database initialization
- `internal/errors/errors.go` - Error types and gRPC mapping
- `internal/logging/logger.go` - Structured logging
- `cmd/server/main.go` - Server entry point
- `proto/userauth/v1/service.proto` - gRPC service definition

### 6. Database (3 files)
- `migrations/001_create_users_table.sql`
- `migrations/002_create_roles_and_permissions.sql`
- `migrations/003_seed_default_data.sql`

### 7. Configuration (5 files)
- `go.mod` - Go dependencies
- `Makefile` - Build automation
- `Dockerfile` - Container image
- `.env.example` - Environment template
- `README.md` - Complete documentation

## Key Functions

### AuthService
```go
Register(email, password, name) (*User, error)
Login(email, password, ipAddress) (*User, string, time.Time, error)
ValidateToken(tokenString) (*User, error)
Logout(tokenString) error
LogoutAllDevices(userID) error
RequestPasswordReset(email) (string, error)
ResetPassword(token, newPassword) error
```

### RBACService
```go
CreateRole(name, description, permissionIDs) (*Role, error)
UpdateRole(roleID, name, description, permissionIDs) (*Role, error)
DeleteRole(roleID) error
AssignRoleToUser(userID, roleID) error
RevokeRoleFromUser(userID, roleID) error
CheckPermission(userID, permission) (bool, error)
GetUserPermissions(userID) ([]string, error)
```

### TokenManager
```go
GenerateToken(user, sessionID) (string, error)
ValidateToken(tokenString) (*TokenClaims, error)
ExtractClaims(tokenString) (*TokenClaims, error)
```

## gRPC Endpoints

### Authentication
- `Register` - Create user account
- `Login` - Authenticate and get JWT
- `Logout` - End session
- `ValidateToken` - Verify JWT
- `RefreshSession` - Renew token
- `RequestPasswordReset` - Request reset
- `ResetPassword` - Complete reset

### RBAC
- `CreateRole` - Create new role
- `UpdateRole` - Modify role
- `DeleteRole` - Remove role
- `AssignRoleToUser` - Grant role
- `RevokeRoleFromUser` - Remove role
- `CheckPermission` - Verify permission
- `GetUserPermissions` - List permissions

## Error Codes
- `INVALID_CREDENTIALS` - Wrong email/password
- `USER_NOT_FOUND` - User doesn't exist
- `EMAIL_ALREADY_EXISTS` - Duplicate email
- `INVALID_TOKEN` - Bad JWT
- `EXPIRED_TOKEN` - JWT expired
- `REVOKED_TOKEN` - JWT revoked
- `ACCOUNT_LOCKED` - Too many failed attempts
- `PERMISSION_DENIED` - Insufficient permissions
- `INVALID_INPUT` - Bad request data
- `ROLE_NOT_FOUND` - Role doesn't exist
- `SYSTEM_ROLE_PROTECTED` - Can't modify system role
- `INVALID_EMAIL` - Bad email format
- `WEAK_PASSWORD` - Password too weak
- `INVALID_RESET_TOKEN` - Bad reset token

## Audit Events
1. `user.registered`
2. `user.login.success`
3. `user.login.failed`
4. `user.logout`
5. `user.password_reset.requested`
6. `user.password_reset.completed`
7. `user.role.assigned`
8. `user.role.revoked`
9. `user.account.locked`
10. `user.logout.all_devices`

## Default Roles & Permissions

### Admin Role
- All permissions (users:*, roles:*, permissions:*, billing:*, analytics:*, features:*)

### Member Role
- users:read
- billing:read
- analytics:read
- features:read

### Viewer Role
- users:read
- roles:read
- permissions:read
- analytics:read
- features:read

## Configuration Keys

### Server
- `GRPC_PORT` - gRPC server port (default: 50051)
- `HOST` - Bind address (default: 0.0.0.0)

### Database
- `DATABASE_URL` - PostgreSQL connection string
- `DB_MAX_CONNECTIONS` - Max connections (default: 25)
- `DB_MAX_IDLE_CONNS` - Max idle (default: 5)
- `DB_CONN_MAX_LIFETIME_MINUTES` - Connection lifetime (default: 5)

### Redis
- `REDIS_HOST` - Redis host (default: localhost)
- `REDIS_PORT` - Redis port (default: 6379)
- `REDIS_PASSWORD` - Redis password (optional)
- `REDIS_DB` - Redis database (default: 0)

### JWT
- `JWT_PRIVATE_KEY_PATH` - Private key path
- `JWT_PUBLIC_KEY_PATH` - Public key path
- `JWT_EXPIRATION_HOURS` - Token lifetime (default: 24)

### Security
- `BCRYPT_COST` - Password hash cost (default: 12)
- `MAX_LOGIN_ATTEMPTS` - Failed attempts limit (default: 5)
- `LOCKOUT_DURATION_MINUTES` - Lockout time (default: 30)
- `PERMISSION_CACHE_TTL_MINUTES` - Cache TTL (default: 5)
- `SESSION_EXPIRATION_HOURS` - Session lifetime (default: 24)
- `PASSWORD_RESET_TTL_MINUTES` - Reset token TTL (default: 60)

### Logging
- `LOG_LEVEL` - Log level (debug, info, warn, error)

## Commands

```bash
# Generate proto code
make proto

# Run tests
make test

# Run with coverage
make test-coverage

# Build binary
make build

# Run service
make run

# Build Docker image
make docker-build

# Clean build artifacts
make clean

# Lint code
make lint
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test -run TestAuthService_Register ./internal/service

# Run with verbose output
go test -v ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Docker

```bash
# Build image
docker build -t user-auth-service .

# Run container
docker run -p 50051:50051 \
  -e DATABASE_URL=postgresql://... \
  -e REDIS_HOST=redis \
  -v $(pwd)/keys:/app/keys \
  user-auth-service
```

## Common Issues

**Proto generation fails**:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**JWT keys missing**:
```bash
cd ../../keys && ./generate-keys.sh
```

**Database connection fails**:
- Check DATABASE_URL format
- Verify PostgreSQL is running
- Test connection: `psql $DATABASE_URL`

**Redis connection fails**:
- Verify REDIS_HOST and REDIS_PORT
- Test connection: `redis-cli -h $REDIS_HOST ping`

## Performance Tips

1. **Connection Pooling**: Adjust `DB_MAX_CONNECTIONS` based on load
2. **Cache TTL**: Tune `PERMISSION_CACHE_TTL_MINUTES` for your use case
3. **bcrypt Cost**: Higher = more secure but slower (12-14 recommended)
4. **Session Expiration**: Balance security vs user experience

## Security Checklist

- ✅ Use strong DATABASE_URL password
- ✅ Generate unique JWT keys for production
- ✅ Set REDIS_PASSWORD in production
- ✅ Use TLS for gRPC in production
- ✅ Set appropriate BCRYPT_COST (12-14)
- ✅ Monitor failed login attempts
- ✅ Regularly rotate JWT keys
- ✅ Keep dependencies updated

---

**Quick Start**: `make proto && make test && make run`
**Documentation**: See README.md for complete guide
**Status**: ✅ Production-ready
