.PHONY: generate test build run clean dev

# Generate GraphQL code
generate:
	go run github.com/99designs/gqlgen generate

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Build the gateway
build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/graphql-api-gateway ./cmd

# Run the gateway
run:
	go run ./cmd/main.go

# Development mode with auto-reload
dev:
	air

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Lint code
lint:
	golangci-lint run

# Docker build
docker-build:
	docker build -t haunted-graphql-gateway:latest .

# Format code
fmt:
	go fmt ./...
	gofmt -s -w .

.DEFAULT_GOAL := build
