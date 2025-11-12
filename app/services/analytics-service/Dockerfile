# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make protobuf-dev

# Copy go mod files
COPY go.mod* go.sum* ./

# Download dependencies first (faster, cacheable)
RUN go mod download || true

# Install protoc-gen-go and protoc-gen-go-grpc (pinned versions for Go 1.21 compatibility)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Copy source code
COPY . .

# Generate proto files if they exist
RUN if [ -f proto/analytics/v1/service.proto ]; then \
        protoc --go_out=. --go_opt=paths=source_relative \
               --go-grpc_out=. --go-grpc_opt=paths=source_relative \
               proto/analytics/v1/*.proto; \
    fi

# Generate go.sum from imports after copying all source
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o analytics-service ./cmd

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Copy binary from builder
COPY --from=builder /build/analytics-service .

# Switch to non-root user
USER appuser

# Expose gRPC port
EXPOSE 50054

# Run the service
CMD ["./analytics-service"]
