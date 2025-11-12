# JWT Keys

This directory contains RSA key pairs for JWT signing and verification.

## Generating Keys

For development:

```bash
./generate-keys.sh
```

For production, use a secure key management service (AWS KMS, HashiCorp Vault, etc.).

## Files

- `jwt-private.pem` - Private key for signing JWTs (keep secure!)
- `jwt-public.pem` - Public key for verifying JWTs (can be shared)

## Security Notes

- **Never commit private keys to version control**
- Use different keys for development and production
- Rotate keys periodically
- Store production keys in a secure key management system
- Set appropriate file permissions (600 for private, 644 for public)
