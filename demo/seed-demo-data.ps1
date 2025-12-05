# PowerShell script to seed demo data

Write-Host "=== Seeding Demo Data ===" -ForegroundColor Cyan
Write-Host ""

# Check if docker-compose is running
$containers = docker-compose ps
if ($containers -notmatch "Up") {
    Write-Host "Error: Docker containers are not running. Please start them first with 'docker-compose up -d'" -ForegroundColor Red
    exit 1
}

Write-Host "Building seed script..." -ForegroundColor Yellow
Set-Location demo

# Initialize go module if not exists
if (-not (Test-Path "go.mod")) {
    go mod init demo-seed
}

go get github.com/lib/pq
go get golang.org/x/crypto/bcrypt

Write-Host ""
Write-Host "Running seed script..." -ForegroundColor Yellow

$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "postgres"
$env:DB_NAME = "haunted_saas"

go run seed-users.go

Write-Host ""
Write-Host "=== Demo data seeded successfully! ===" -ForegroundColor Green
