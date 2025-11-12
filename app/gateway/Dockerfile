# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make protobuf-dev

# Install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Copy all source code (gateway and services for proto files)
COPY ./gateway/graphql-api-gateway/ ./
COPY ./services/ ../services/

# Generate proto files for all services
RUN cd ../services/user-auth-service && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/userauth/v1/*.proto && \
    cd /app

RUN cd ../services/billing-service && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/billing/v1/*.proto && \
    cd /app

RUN cd ../services/llm-gateway-service && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/llm/v1/*.proto && \
    cd /app

RUN cd ../services/notifications-service && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/notifications/v1/*.proto && \
    cd /app

RUN cd ../services/analytics-service && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/analytics/v1/*.proto && \
    cd /app

RUN cd ../services/feature-flags-service && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/featureflags/v1/*.proto && \
    cd /app

# Create stubs for generated code and models to satisfy imports temporarily
RUN mkdir -p internal/generated internal/models && \
    echo 'package generated' > internal/generated/stub.go && \
    echo 'package models' > internal/models/stub.go

# Download dependencies and tidy
RUN go mod download && go mod tidy

# Generate GraphQL code (this will replace the stub)
# We ignore errors because schema.resolvers.go will have compilation errors
RUN go run github.com/99designs/gqlgen generate || echo "gqlgen generated with expected errors"

# Remove the auto-generated schema.resolvers.go since we have manual resolvers
# This file conflicts with our manually written mutation.resolvers.go and query.resolvers.go
RUN rm -f internal/resolvers/schema.resolvers.go

# Run tidy again after removing conflicting file
RUN go mod tidy

# Build the application (now without the conflicting file)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Production stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/main .

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S appuser -u 1001 -G appuser

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:4000/health || exit 1

# Run the application
CMD ["./main"]
