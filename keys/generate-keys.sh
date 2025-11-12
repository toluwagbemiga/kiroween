#!/bin/bash

# Generate RSA key pair for JWT signing
# Run this script before starting the services

set -e

echo "🔐 Generating JWT RSA key pair..."

# Generate private key
openssl genrsa -out jwt-private.pem 2048

# Generate public key from private key
openssl rsa -in jwt-private.pem -pubout -out jwt-public.pem

# Set appropriate permissions
chmod 600 jwt-private.pem
chmod 644 jwt-public.pem

echo "✅ JWT keys generated successfully!"
echo "   Private key: jwt-private.pem"
echo "   Public key:  jwt-public.pem"
echo ""
echo "⚠️  Keep jwt-private.pem secure and never commit it to version control!"
