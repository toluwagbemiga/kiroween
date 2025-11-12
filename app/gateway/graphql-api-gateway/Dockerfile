# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate GraphQL code
RUN go run github.com/99designs/gqlgen generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o graphql-api-gateway ./cmd

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Copy binary from builder
COPY --from=builder /build/graphql-api-gateway .

# Switch to non-root user
USER appuser

# Expose HTTP port
EXPOSE 8080

# Run the gateway
CMD ["./graphql-api-gateway"]
