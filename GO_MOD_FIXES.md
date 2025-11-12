# Go Module Version Fixes

## Issues Found and Fixed

### Problem
Docker builds were failing with "unknown revision" errors for Go dependencies.

### Root Causes
Invalid version numbers in go.mod files:

1. **golang.org/x/net v1.20.0** → Should be **v0.20.0**
   - The `golang.org/x/net` package uses v0.x versioning, not v1.x

2. **google.golang.org/genproto/googleapis/rpc v1.0.0** → Should be **v0.0.0-20231212172506-995d672761c0**
   - This package uses pseudo-versions (commit-based), not semantic versioning

## Files Fixed

### golang.org/x/net version (v1.20.0 → v0.20.0)
- ✅ `app/services/analytics-service/go.mod`
- ✅ `app/services/notifications-service/go.mod`
- ✅ `app/services/feature-flags-service/go.mod`
- ✅ `app/gateway/graphql-api-gateway/go.mod`

### golang.org/x/sys version (v1.16.0 → v0.16.0)
- ✅ `app/services/notifications-service/go.mod`
- ✅ `app/services/feature-flags-service/go.mod`
- ✅ `app/gateway/graphql-api-gateway/go.mod`

### golang.org/x/text version (v1.14.0 → v0.14.0)
- ✅ `app/services/notifications-service/go.mod`
- ✅ `app/services/feature-flags-service/go.mod`
- ✅ `app/gateway/graphql-api-gateway/go.mod`

### google.golang.org/genproto version (v1.0.0 → v0.0.0-20231212172506-995d672761c0)
- ✅ `app/services/analytics-service/go.mod`
- ✅ `app/services/notifications-service/go.mod`
- ✅ `app/services/feature-flags-service/go.mod`
- ✅ `app/gateway/graphql-api-gateway/go.mod`

## Why This Happened

These invalid versions likely came from:
1. Manual editing of go.mod files
2. Copy-paste errors
3. Incorrect version assumptions

## How to Prevent

### Use `go get` to add dependencies:
```bash
# Correct way
go get golang.org/x/net@latest

# This will add the correct version automatically
```

### Use `go mod tidy` to clean up:
```bash
cd app/services/analytics-service
go mod tidy
# This validates and fixes version numbers
```

### Check versions before committing:
```bash
# Verify a specific package version exists
go list -m -versions golang.org/x/net
```

## Testing the Fix

```bash
# Build all services
docker-compose build

# Or build specific service
docker-compose build analytics-service
```

## Expected Behavior

After these fixes:
- ✅ `go mod download` should complete successfully
- ✅ All dependencies should resolve correctly
- ✅ Docker builds should progress past the download stage
- ✅ Services should compile successfully

## Common Go Version Patterns

### Semantic Versioning (v1.x.x, v2.x.x)
```
github.com/google/uuid v1.5.0
google.golang.org/grpc v1.60.1
```

### v0.x Versioning (pre-1.0 packages)
```
golang.org/x/net v0.20.0
golang.org/x/crypto v0.17.0
```

### Pseudo-versions (commit-based)
```
google.golang.org/genproto/googleapis/rpc v0.0.0-20231212172506-995d672761c0
```

Format: `v0.0.0-YYYYMMDDHHMMSS-commithash`

## Verification

To verify all go.mod files are valid:

```bash
# Check each service
for dir in app/services/*/; do
  echo "Checking $dir"
  cd "$dir"
  go mod verify
  cd -
done
```

## Status

✅ All invalid versions fixed
✅ Ready to rebuild Docker images
✅ Services should now build successfully

---

**Next Step**: Run `docker-compose up --build` to rebuild with fixed dependencies
