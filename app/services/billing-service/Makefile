.PHONY: proto test build run clean

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/billing/v1/service.proto

test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/billing-service ./cmd

run:
	go run ./cmd/main.go

clean:
	rm -rf bin/ coverage.out coverage.html

lint:
	golangci-lint run

docker-build:
	docker build -t haunted-billing-service:latest .

.DEFAULT_GOAL := build
