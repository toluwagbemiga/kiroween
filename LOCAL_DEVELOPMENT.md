# Local Development Guide (Without Docker)

## Why Run Locally?

Docker builds are hitting issues with:
- Missing proto-generated files
- Complex build dependencies
- Network timeouts

Running services locally is faster for development and easier to debug.

## Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- Protocol Buffers compiler (protoc)

## Quick Setup

### 1. Install Dependencies

**Windows (using Chocolatey)**:
```powershell
choco install golang nodejs postgresql redis protoc
```

**Or download manually**:
- Go: https://go.dev/dl/
- Node.js: https://nodejs.org/
- PostgreSQL: https://www.postgresql.org/download/windows/
- Redis: https://github.com/microsoftarchive/redis/releases
- Protoc: https://github.com/protocolbuffers/protobuf/releases

### 2. Start Infrastructure

**PostgreSQL**:
```powershell
# Start PostgreSQL service
net start postgresql-x64-15

# Create database
psql -U postgres
CREATE DATABASE haunted;
CREATE USER haunted WITH PASSWORD 'haunted_dev_pass';
GRANT ALL PRIVILEGES ON DATABASE haunted TO haunted;
\q
```

**Redis**:
```powershell
# Start Redis
redis-server
```

**Unleash (Docker)**:
```powershell
docker-compose up -d unleash-db unleash
```

### 3. Setup Environment Variables

Create `.env` files for each service (copy from `.env.example`):

```powershell
# Copy all .env.example files
Get-ChildItem -Path app\services -Recurse -Filter ".env.example" | ForEach-Object {
    Copy-Item $_.FullName ($_.FullName -replace '\.example$','')
}
```

### 4. Generate Proto Files

For each service with proto files:

```powershell
# Analytics Service
cd app\services\analytics-service
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       proto/analytics/v1/service.proto

# Repeat for other services...
```

Or use the Makefiles:
```powershell
cd app\services\analytics-service
make proto
```

### 5. Run Database Migrations

```powershell
# User Auth Service
cd app\services\user-auth-service
# Run migrations (implement migration runner or use psql)
psql -U haunted -d haunted -f migrations/001_create_users_table.sql
psql -U haunted -d haunted -f migrations/002_create_roles_and_permissions.sql
psql -U haunted -d haunted -f migrations/003_seed_default_data.sql

# Billing Service
cd app\services\billing-service
psql -U haunted -d haunted -f migrations/001_create_plans_table.sql
psql -U haunted -d haunted -f migrations/002_create_subscriptions_table.sql
```

### 6. Start Services

Open multiple terminal windows:

**Terminal 1 - User Auth Service**:
```powershell
cd app\services\user-auth-service
go run cmd\server\main.go
```

**Terminal 2 - Billing Service**:
```powershell
cd app\services\billing-service
go run cmd\main.go
```

**Terminal 3 - Analytics Service**:
```powershell
cd app\services\analytics-service
go run cmd\main.go
```

**Terminal 4 - Notifications Service**:
```powershell
cd app\services\notifications-service
go run cmd\main.go
```

**Terminal 5 - LLM Gateway Service**:
```powershell
cd app\services\llm-gateway-service
go run cmd\main.go
```

**Terminal 6 - Feature Flags Service**:
```powershell
cd app\services\feature-flags-service
go run cmd\main.go
```

**Terminal 7 - GraphQL Gateway**:
```powershell
cd app\gateway\graphql-api-gateway
go run cmd\main.go
```

**Terminal 8 - Frontend**:
```powershell
cd app\frontend
npm install
npm run dev
```

**Terminal 9 - Documentation**:
```powershell
cd docs
npm install
npm start
```

## Service Ports

| Service | Port | URL |
|---------|------|-----|
| User Auth | 50051 | gRPC only |
| Billing | 50052, 8080 | gRPC + HTTP webhooks |
| Analytics | 50055 | gRPC only |
| Notifications | 50054, 3002 | gRPC + Socket.IO |
| LLM Gateway | 50053 | gRPC only |
| Feature Flags | 50056 | gRPC only |
| GraphQL Gateway | 4000 | http://localhost:4000/graphql |
| Frontend | 3000 | http://localhost:3000 |
| Docs | 3001 | http://localhost:3001 |

## Simplified Startup Script

Create `start-local.ps1`:

```powershell
# Start infrastructure
Write-Host "Starting infrastructure..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-Command", "redis-server"
docker-compose up -d unleash-db unleash

# Wait for databases
Start-Sleep -Seconds 5

# Start services
Write-Host "Starting services..." -ForegroundColor Green

$services = @(
    "app\services\user-auth-service",
    "app\services\billing-service",
    "app\services\analytics-service",
    "app\services\notifications-service",
    "app\services\llm-gateway-service",
    "app\services\feature-flags-service"
)

foreach ($service in $services) {
    $name = Split-Path $service -Leaf
    Write-Host "Starting $name..." -ForegroundColor Cyan
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd $service; go run cmd\main.go"
}

# Start gateway
Start-Sleep -Seconds 3
Write-Host "Starting GraphQL Gateway..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd app\gateway\graphql-api-gateway; go run cmd\main.go"

# Start frontend
Start-Sleep -Seconds 2
Write-Host "Starting Frontend..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd app\frontend; npm run dev"

Write-Host "`nAll services starting!" -ForegroundColor Green
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Yellow
Write-Host "GraphQL: http://localhost:4000/graphql" -ForegroundColor Yellow
```

## Advantages of Local Development

✅ **Faster iteration** - No Docker rebuild needed
✅ **Better debugging** - Direct access to logs and debuggers
✅ **Hot reload** - Changes reflect immediately
✅ **Easier troubleshooting** - Clear error messages
✅ **Resource efficient** - Lower memory/CPU usage

## When to Use Docker

- ✅ Production deployments
- ✅ CI/CD pipelines
- ✅ Testing full system integration
- ✅ Sharing with team members
- ✅ Consistent environments

## Troubleshooting

### Port Already in Use
```powershell
# Find process using port
netstat -ano | findstr :3000

# Kill process
taskkill /PID <PID> /F
```

### Database Connection Failed
```powershell
# Check PostgreSQL is running
Get-Service postgresql*

# Start if stopped
Start-Service postgresql-x64-15
```

### Proto Files Not Found
```powershell
# Generate for all services
cd app\services\analytics-service
make proto

# Or manually
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       proto/analytics/v1/service.proto
```

## Next Steps

Once all services are running:
1. Access frontend at http://localhost:3000
2. Login with demo credentials
3. Explore the GraphQL API at http://localhost:4000/graphql
4. View documentation at http://localhost:3001

---

**Recommendation**: Use local development for now, fix Docker builds later when you have more time.
