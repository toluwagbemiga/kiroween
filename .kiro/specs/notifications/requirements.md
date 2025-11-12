# Requirements Document: Notifications Service

## Introduction

The Notifications Service manages all real-time communication for the Haunted SaaS Skeleton platform, serving as the Liveness Layer. It uses Socket.IO to provide reliable real-time messaging with automatic fallback to HTTP long-polling for network-restricted environments. The service authenticates connections using JWT and provides gRPC endpoints for other backend services to push messages to users and teams.

## Glossary

- **Notifications Service**: The Go-based microservice responsible for managing real-time communication
- **Socket.IO**: A real-time communication library that provides WebSocket connections with HTTP long-polling fallback
- **HTTP Long-Polling**: A fallback communication method where clients make repeated HTTP requests to receive updates
- **JWT**: JSON Web Token used for authenticating socket connections
- **Room**: A Socket.IO concept representing a group of connected clients (e.g., all users in a team)
- **gRPC**: The internal communication protocol used between microservices
- **Connection**: An active Socket.IO connection between a client and the server

## Requirements

### Requirement 1

**User Story:** As a frontend client, I want to establish a real-time connection using Socket.IO, so that I can receive instant notifications

#### Acceptance Criteria

1. THE Notifications Service SHALL accept Socket.IO connections on a configurable port
2. THE Notifications Service SHALL support WebSocket transport as the primary connection method
3. WHEN WebSocket connections fail, THE Notifications Service SHALL automatically fall back to HTTP long-polling
4. THE Notifications Service SHALL configure Socket.IO with CORS settings to allow connections from the frontend domain
5. THE Notifications Service SHALL log each new connection attempt with the client IP and transport method

### Requirement 2

**User Story:** As a frontend client, I want to authenticate my socket connection with a JWT, so that only authorized users can receive notifications

#### Acceptance Criteria

1. WHEN a client attempts to connect, THE Notifications Service SHALL require a JWT in the connection handshake
2. THE Notifications Service SHALL validate the JWT signature using the shared secret key
3. WHEN the JWT is invalid or expired, THE Notifications Service SHALL reject the connection with an authentication error
4. WHEN the JWT is valid, THE Notifications Service SHALL extract the user ID and team ID from the token claims
5. THE Notifications Service SHALL associate the user ID and team ID with the socket connection for routing messages

### Requirement 3

**User Story:** As a frontend client, I want to automatically join rooms based on my team membership, so that I receive team-wide notifications

#### Acceptance Criteria

1. WHEN a client connection is authenticated, THE Notifications Service SHALL automatically join the client to a room named after their team ID
2. THE Notifications Service SHALL join the client to a user-specific room named after their user ID
3. THE Notifications Service SHALL log room join events with the user ID, team ID, and socket ID
4. WHEN a client disconnects, THE Notifications Service SHALL automatically remove the client from all rooms
5. THE Notifications Service SHALL support clients being members of multiple rooms simultaneously

### Requirement 4

**User Story:** As a backend service, I want to send a notification to a specific user, so that I can deliver personalized messages

#### Acceptance Criteria

1. THE Notifications Service SHALL provide a SendToUser gRPC endpoint that accepts user ID, event type, and message payload
2. WHEN SendToUser is called, THE Notifications Service SHALL emit the message to all active connections in the user's room
3. THE Notifications Service SHALL return success when the message is emitted regardless of whether the user is connected
4. THE Notifications Service SHALL include the event type and payload in the Socket.IO emit
5. THE Notifications Service SHALL log the message delivery with user ID, event type, and number of active connections

### Requirement 5

**User Story:** As a backend service, I want to broadcast a notification to all users in a team, so that I can send team-wide announcements

#### Acceptance Criteria

1. THE Notifications Service SHALL provide a BroadcastToRoom gRPC endpoint that accepts team ID, event type, and message payload
2. WHEN BroadcastToRoom is called, THE Notifications Service SHALL emit the message to all active connections in the team's room
3. THE Notifications Service SHALL return success when the message is emitted regardless of the number of connected users
4. THE Notifications Service SHALL include the event type and payload in the Socket.IO emit
5. THE Notifications Service SHALL log the broadcast with team ID, event type, and number of recipients

### Requirement 6

**User Story:** As a backend service, I want to send notifications to multiple users at once, so that I can efficiently deliver messages to a group

#### Acceptance Criteria

1. THE Notifications Service SHALL provide a SendToUsers gRPC endpoint that accepts a list of user IDs, event type, and message payload
2. WHEN SendToUsers is called, THE Notifications Service SHALL emit the message to all specified users' rooms
3. THE Notifications Service SHALL process the user list and emit messages in parallel
4. THE Notifications Service SHALL return success when all messages are emitted
5. THE Notifications Service SHALL log the batch delivery with the number of target users and event type

### Requirement 7

**User Story:** As a system administrator, I want to monitor active connections, so that I can understand real-time usage and troubleshoot issues

#### Acceptance Criteria

1. THE Notifications Service SHALL maintain a count of active Socket.IO connections
2. THE Notifications Service SHALL provide a GetConnectionStats gRPC endpoint that returns connection metrics
3. THE Notifications Service SHALL include in the stats the total connections, connections by transport type, and connections by team
4. WHEN a connection is established or terminated, THE Notifications Service SHALL update the connection metrics
5. THE Notifications Service SHALL expose connection metrics in a format suitable for monitoring systems

### Requirement 8

**User Story:** As a frontend client, I want to receive acknowledgment when my connection is established, so that I know I'm ready to receive notifications

#### Acceptance Criteria

1. WHEN a client connection is authenticated and rooms are joined, THE Notifications Service SHALL emit a connection_ready event to the client
2. THE Notifications Service SHALL include in the connection_ready event the user ID and list of joined rooms
3. THE Notifications Service SHALL emit the connection_ready event before processing any other messages
4. WHEN a client reconnects, THE Notifications Service SHALL emit a new connection_ready event
5. THE Notifications Service SHALL log the connection_ready event emission with the socket ID

### Requirement 9

**User Story:** As a frontend client, I want to handle disconnections gracefully, so that I can reconnect automatically without losing notifications

#### Acceptance Criteria

1. THE Notifications Service SHALL configure Socket.IO to allow automatic reconnection attempts
2. WHEN a client disconnects, THE Notifications Service SHALL log the disconnection with user ID and reason
3. WHEN a client reconnects, THE Notifications Service SHALL re-authenticate the JWT and rejoin rooms
4. THE Notifications Service SHALL maintain connection state for 30 seconds after disconnection to support quick reconnects
5. THE Notifications Service SHALL clean up connection resources after the reconnection window expires

### Requirement 10

**User Story:** As a system administrator, I want connection errors to be logged, so that I can diagnose connectivity issues

#### Acceptance Criteria

1. WHEN a connection attempt fails authentication, THE Notifications Service SHALL log the failure with the client IP and error reason
2. WHEN a Socket.IO error occurs, THE Notifications Service SHALL log the error with the socket ID and error message
3. WHEN a message delivery fails, THE Notifications Service SHALL log the failure with the target user or room and error details
4. THE Notifications Service SHALL include correlation IDs in all log entries to trace requests across services
5. THE Notifications Service SHALL log at appropriate levels (INFO for connections, ERROR for failures)

### Requirement 11

**User Story:** As a security administrator, I want JWT secrets to be securely managed, so that authentication credentials are not exposed

#### Acceptance Criteria

1. THE Notifications Service SHALL load the JWT secret from an environment variable at startup
2. THE Notifications Service SHALL validate that the JWT secret is present before accepting connections
3. THE Notifications Service SHALL never log or expose the JWT secret in responses or error messages
4. THE Notifications Service SHALL use the same JWT secret as other services for consistent authentication
5. THE Notifications Service SHALL terminate startup if the JWT secret is missing or invalid

### Requirement 12

**User Story:** As a backend service, I want to check if a user is currently connected, so that I can decide whether to send real-time or offline notifications

#### Acceptance Criteria

1. THE Notifications Service SHALL provide a IsUserConnected gRPC endpoint that accepts a user ID
2. WHEN IsUserConnected is called, THE Notifications Service SHALL check if the user has any active socket connections
3. THE Notifications Service SHALL return a boolean indicating connection status and the number of active connections
4. THE Notifications Service SHALL respond to IsUserConnected requests within 100 milliseconds
5. THE Notifications Service SHALL log IsUserConnected requests at DEBUG level to avoid log noise

### Requirement 13

**User Story:** As a developer, I want to send different types of notifications, so that clients can handle messages appropriately

#### Acceptance Criteria

1. THE Notifications Service SHALL support arbitrary event types specified by calling services
2. THE Notifications Service SHALL pass the event type to Socket.IO emit without modification
3. THE Notifications Service SHALL support common event types including notification, alert, update, and message
4. THE Notifications Service SHALL validate that event types are non-empty strings
5. THE Notifications Service SHALL log the event type with each message delivery for debugging

### Requirement 14

**User Story:** As a system administrator, I want to configure connection limits, so that the service can handle load appropriately

#### Acceptance Criteria

1. THE Notifications Service SHALL enforce a configurable maximum number of concurrent connections
2. WHEN the connection limit is reached, THE Notifications Service SHALL reject new connections with a capacity error
3. THE Notifications Service SHALL log when the connection limit is reached
4. THE Notifications Service SHALL allow the connection limit to be configured via environment variable
5. THE Notifications Service SHALL set a default connection limit of 10000 when not configured
