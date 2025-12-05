#!/bin/bash

echo "=== Seeding Demo Data ==="
echo ""

# Check if docker-compose is running
if ! docker-compose ps | grep -q "Up"; then
    echo "Error: Docker containers are not running. Please start them first with 'docker-compose up -d'"
    exit 1
fi

echo "Building seed script..."
cd demo
go mod init demo-seed 2>/dev/null || true
go get github.com/lib/pq
go get golang.org/x/crypto/bcrypt

echo ""
echo "Running seed script..."
DB_HOST=localhost \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=haunted_saas \
go run seed-users.go

echo ""
echo "=== Demo data seeded successfully! ==="
