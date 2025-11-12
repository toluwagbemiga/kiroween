# Docker Troubleshooting Guide

## Common Issues and Solutions

### 1. "failed to read dockerfile: no such file or directory"

**Problem**: Docker can't find the Dockerfile for a service.

**Solution**:
```bash
# Verify all Dockerfiles exist
ls -la app/services/*/Dockerfile
ls -la app/frontend/Dockerfile
ls -la app/gateway/Dockerfile
ls -la docs/Dockerfile

# If any are missing, they need to be created
```

**Services that need Dockerfiles**:
- ✅ `app/services/user-auth-service/Dockerfile`
- ✅ `app/services/billing-service/Dockerfile`
- ✅ `app/services/llm-gateway-service/Dockerfile`
- ✅ `app/services/notifications-service/Dockerfile`
- ✅ `app/services/analytics-service/Dockerfile`
- ✅ `app/services/feature-flags-service/Dockerfile`
- ✅ `app/gateway/Dockerfile`
- ✅ `app/frontend/Dockerfile`
- ✅ `docs/Dockerfile`

### 2. Build Context Issues

**Problem**: Docker can't find files referenced in Dockerfile.

**Solution**:
```bash
# Check build context in docker-compose.yml
# Context should point to the directory containing the Dockerfile

# Example:
services:
  docs:
    build:
      context: ./docs  # This directory must contain Dockerfile
      dockerfile: Dockerfile
```

### 3. Port Conflicts

**Problem**: Port already in use.

**Solution**:
```bash
# Check what's using the port (Windows)
netstat -ano | findstr :3000

# Kill the process
taskkill /PID <PID> /F

# Or change the port in docker-compose.yml
```

### 4. Build Fails Due to Missing Dependencies

**Problem**: npm install or go mod download fails.

**Solution**:
```bash
# For Node.js services
cd app/frontend
npm install  # Test locally first

# For Go services
cd app/services/user-auth-service
go mod download
go mod tidy
```

### 5. Services Won't Start

**Problem**: Container exits immediately.

**Solution**:
```bash
# Check logs
docker-compose logs <service-name>

# Example
docker-compose logs user-auth-service

# Check if dependencies are healthy
docker-compose ps
```

### 6. Database Connection Issues

**Problem**: Services can't connect to PostgreSQL or Redis.

**Solution**:
```bash
# Wait for databases to be healthy
docker-compose up -d postgres redis
docker-compose ps

# Check health status
docker-compose exec postgres pg_isready -U haunted
docker-compose exec redis redis-cli ping

# Restart dependent services
docker-compose restart user-auth-service
```

### 7. Build Cache Issues

**Problem**: Changes not reflected in container.

**Solution**:
```bash
# Rebuild without cache
docker-compose build --no-cache <service-name>

# Or rebuild all
docker-compose build --no-cache

# Then restart
docker-compose up -d
```

### 8. Volume Permission Issues

**Problem**: Permission denied errors.

**Solution**:
```bash
# On Windows, ensure Docker has access to the drive
# Docker Desktop -> Settings -> Resources -> File Sharing

# Check volume mounts in docker-compose.yml
volumes:
  - ./prompts:/app/prompts:ro  # :ro = read-only
```

### 9. Network Issues

**Problem**: Services can't communicate.

**Solution**:
```bash
# Check network
docker network ls
docker network inspect haunted-saas_default

# Restart networking
docker-compose down
docker-compose up -d
```

### 10. Out of Disk Space

**Problem**: No space left on device.

**Solution**:
```bash
# Clean up Docker
docker system prune -a --volumes

# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune
```

## Quick Fixes

### Start Fresh

```bash
# Stop everything
docker-compose down -v

# Remove all containers and volumes
docker-compose rm -f -v

# Rebuild and start
docker-compose up --build -d
```

### Check Service Health

```bash
# View all services
docker-compose ps

# Check specific service logs
docker-compose logs -f user-auth-service

# Check last 100 lines
docker-compose logs --tail=100 user-auth-service
```

### Restart Single Service

```bash
# Rebuild and restart one service
docker-compose up -d --build user-auth-service

# Just restart (no rebuild)
docker-compose restart user-auth-service
```

### Access Service Shell

```bash
# Access running container
docker-compose exec user-auth-service sh

# Or for services with bash
docker-compose exec postgres bash
```

## Build Order

Services should be built in this order due to dependencies:

1. **Infrastructure** (no dependencies):
   - postgres
   - redis
   - unleash-db
   - unleash

2. **Core Services** (depend on infrastructure):
   - user-auth-service
   - billing-service
   - analytics-service
   - notifications-service
   - llm-gateway-service
   - feature-flags-service

3. **Gateway** (depends on core services):
   - graphql-gateway

4. **Frontend** (depends on gateway):
   - frontend
   - docs

## Recommended Startup

```bash
# 1. Start infrastructure
docker-compose up -d postgres redis unleash-db unleash

# 2. Wait for health checks
docker-compose ps

# 3. Start backend services
docker-compose up -d user-auth-service billing-service analytics-service \
  notifications-service llm-gateway-service feature-flags-service

# 4. Start gateway
docker-compose up -d graphql-gateway

# 5. Start frontend
docker-compose up -d frontend docs

# Or just start everything (Docker handles dependencies)
docker-compose up -d
```

## Monitoring

```bash
# Watch all logs
docker-compose logs -f

# Watch specific service
docker-compose logs -f user-auth-service

# Check resource usage
docker stats

# Check container details
docker-compose ps -a
```

## Environment Variables

```bash
# Check environment variables in container
docker-compose exec user-auth-service env

# Override environment variables
docker-compose up -d -e DATABASE_URL=custom_url
```

## Common Error Messages

### "dial tcp: lookup postgres: no such host"
- **Cause**: Service started before postgres
- **Fix**: Wait for postgres health check, then restart service

### "connection refused"
- **Cause**: Target service not running or wrong port
- **Fix**: Check service is running with `docker-compose ps`

### "permission denied"
- **Cause**: File permissions or volume mount issues
- **Fix**: Check file permissions and Docker file sharing settings

### "port is already allocated"
- **Cause**: Port conflict with another process
- **Fix**: Stop the conflicting process or change port in docker-compose.yml

### "no space left on device"
- **Cause**: Docker disk space full
- **Fix**: Run `docker system prune -a --volumes`

## Getting Help

1. **Check logs**: `docker-compose logs <service>`
2. **Check health**: `docker-compose ps`
3. **Verify config**: `docker-compose config`
4. **Test locally**: Run service outside Docker first
5. **Check documentation**: Service-specific README files

## Useful Commands

```bash
# View compose file with resolved variables
docker-compose config

# Validate compose file
docker-compose config --quiet

# List all containers
docker-compose ps -a

# Remove stopped containers
docker-compose rm

# Stop and remove everything
docker-compose down -v

# Follow logs for multiple services
docker-compose logs -f user-auth-service billing-service

# Execute command in running container
docker-compose exec user-auth-service ls -la

# Run one-off command
docker-compose run --rm user-auth-service sh
```

## Performance Tips

1. **Use .dockerignore**: Exclude unnecessary files from build context
2. **Multi-stage builds**: Reduce final image size
3. **Layer caching**: Order Dockerfile commands for better caching
4. **Prune regularly**: Clean up unused images and volumes
5. **Limit resources**: Set memory and CPU limits in docker-compose.yml

---

**Need more help?** Check service-specific README files or logs for detailed error messages.
