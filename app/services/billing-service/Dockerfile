FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make protobuf-dev

# Copy go mod files
COPY go.mod* go.sum* ./

# Download dependencies first (faster, cacheable)
RUN go mod download || true

# Copy source code
COPY . .

# Install protoc-gen-go and protoc-gen-go-grpc (pinned versions for Go 1.21 compatibility)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Generate proto files
RUN mkdir -p proto/billing/v1 && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/billing/v1/*.proto || true

# Generate go.sum from imports
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o billing-service ./cmd

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/billing-service .

EXPOSE 50052 8080

CMD ["./billing-service"]
