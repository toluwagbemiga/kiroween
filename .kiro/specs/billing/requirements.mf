# Requirements Document

## Introduction

The User Authentication Service (`user-auth-service`) is a core microservice in the HAUNTED SAAS SKELETON architecture that provides secure authentication, authorization, and granular Role-Based Access Control (RBAC) capabilities. This service will be implemented in Go and communicate via gRPC with other services, while exposing functionality through the GraphQL API Gateway for external clients.

## Glossary

- **User Auth Service**: The Go-based microservice responsible for authentication, authorization, and RBAC
- **RBAC**: Role-Based Access Control - a method of regulating access based on roles assigned to users
- **JWT**: JSON Web Token - a compact, URL-safe means of representing claims to be transferred between parties
- **gRPC**: Google Remote Procedure Call - the internal communication protocol between services
- **GraphQL Gateway**: The single external API entry point that translates GraphQL queries to internal gRPC calls
- **Session Store**: Redis-based storage for active user sessions
- **User Database**: PostgreSQL database containing user accounts, roles, and permissions
- **Permission**: A granular capability that can be assigned to roles (e.g., "users:read", "billing:write")
- **Role**: A collection of permissions that can be assigned to users (e.g., "admin", "member", "viewer")

## Requirements

### Requirement 1

**User Story:** As a new user, I want to register an account with email and password, so that I can access the platform

#### Acceptance Criteria

1. WHEN a registration request is received with valid email and password, THE User Auth Service SHALL create a new user account in the User Database
2. WHEN a registration request is received with an email that already exists, THE User Auth Service SHALL return an error indicating the email is already registered
3. THE User Auth Service SHALL hash passwords using bcrypt with a minimum cost factor of 12 before storing in the User Database
4. WHEN a user account is successfully created, THE User Auth Service SHALL assign the default "member" role to the new user
5. THE User Auth Service SHALL validate that email addresses conform to RFC 5322 standard format before account creation

### Requirement 2

**User Story:** As a registered user, I want to authenticate with my credentials, so that I can access protected resources

#### Acceptance Criteria

1. WHEN a login request is received with valid credentials, THE User Auth Service SHALL generate a JWT containing user ID, roles, and permissions
2. WHEN a login request is received with invalid credentials, THE User Auth Service SHALL return an authentication error without revealing whether the email or password was incorrect
3. WHEN a JWT is successfully generated, THE User Auth Service SHALL create a session record in the Session Store with a 24-hour expiration
4. THE User Auth Service SHALL include the following claims in the JWT: user_id, email, roles, permissions, issued_at, and expires_at
5. WHEN a user has multiple failed login attempts exceeding 5 within 15 minutes, THE User Auth Service SHALL temporarily lock the account for 30 minutes

### Requirement 3

**User Story:** As an authenticated user, I want my session to be validated on each request, so that only authorized users can access protected resources

#### Acceptance Criteria

1. WHEN a request includes a valid JWT, THE User Auth Service SHALL verify the token signature and expiration
2. WHEN a request includes an expired JWT, THE User Auth Service SHALL return an authentication error requiring re-login
3. WHEN a request includes a JWT for a session that has been revoked, THE User Auth Service SHALL reject the request and return an authentication error
4. THE User Auth Service SHALL validate JWT signatures using the RS256 algorithm with public key verification
5. WHEN a JWT is validated successfully, THE User Auth Service SHALL extend the session expiration in the Session Store by 24 hours from the current time

### Requirement 4

**User Story:** As a system administrator, I want to assign roles and permissions to users, so that I can control access to different parts of the system

#### Acceptance Criteria

1. WHEN an administrator assigns a role to a user, THE User Auth Service SHALL update the user's role assignments in the User Database
2. WHEN an administrator creates a custom role, THE User Auth Service SHALL store the role definition with its associated permissions in the User Database
3. THE User Auth Service SHALL support multiple role assignments per user
4. WHEN a user's roles are modified, THE User Auth Service SHALL invalidate all active sessions for that user in the Session Store
5. THE User Auth Service SHALL provide gRPC endpoints for role management operations: CreateRole, UpdateRole, DeleteRole, AssignRoleToUser, and RevokeRoleFromUser

### Requirement 5

**User Story:** As a service developer, I want to check if a user has specific permissions, so that I can enforce authorization rules in my service

#### Acceptance Criteria

1. WHEN a permission check request is received via gRPC, THE User Auth Service SHALL evaluate whether the user has the requested permission through any of their assigned roles
2. WHEN evaluating permissions, THE User Auth Service SHALL aggregate permissions from all roles assigned to the user
3. THE User Auth Service SHALL cache permission lookups in Redis for 5 minutes to optimize performance
4. WHEN a permission check is performed for an invalid or expired session, THE User Auth Service SHALL return an authorization error
5. THE User Auth Service SHALL provide a gRPC endpoint CheckPermission that accepts user_id and permission_name and returns a boolean result

### Requirement 6

**User Story:** As a user, I want to log out of my account, so that my session is terminated and my access is revoked

#### Acceptance Criteria

1. WHEN a logout request is received with a valid JWT, THE User Auth Service SHALL remove the session from the Session Store
2. WHEN a logout request is received, THE User Auth Service SHALL add the JWT to a revocation list in Redis with expiration matching the token's original expiration
3. THE User Auth Service SHALL provide a gRPC endpoint for logout that accepts a session token
4. WHEN a user logs out, THE User Auth Service SHALL return a success confirmation
5. WHERE a user requests to log out from all devices, THE User Auth Service SHALL remove all sessions for that user from the Session Store

### Requirement 7

**User Story:** As a security administrator, I want password reset functionality, so that users can regain access to their accounts securely

#### Acceptance Criteria

1. WHEN a password reset is requested, THE User Auth Service SHALL generate a cryptographically secure reset token and store it in Redis with a 1-hour expiration
2. WHEN a password reset token is used, THE User Auth Service SHALL validate the token exists in Redis and has not expired
3. WHEN a password is successfully reset, THE User Auth Service SHALL invalidate all active sessions for that user in the Session Store
4. THE User Auth Service SHALL hash the reset token before storing in Redis
5. WHEN a password reset is completed, THE User Auth Service SHALL remove the reset token from Redis

### Requirement 8

**User Story:** As a platform operator, I want comprehensive audit logging of authentication events, so that I can monitor security and troubleshoot issues

#### Acceptance Criteria

1. WHEN any authentication event occurs, THE User Auth Service SHALL emit an event to the analytics service with event type, user ID, timestamp, and IP address
2. THE User Auth Service SHALL log the following events: registration, login_success, login_failure, logout, password_reset_requested, password_reset_completed, role_assigned, role_revoked
3. WHEN a security-relevant event occurs, THE User Auth Service SHALL write structured logs to stdout in JSON format
4. THE User Auth Service SHALL include correlation IDs in all log entries to enable request tracing
5. WHEN logging authentication failures, THE User Auth Service SHALL not include sensitive information such as passwords or tokens in log messages
