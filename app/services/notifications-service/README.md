# Notifications Service

✅ **COMPLETE IMPLEMENTATION** - Real-time bidirectional communication with Socket.IO for HAUNTED SAAS SKELETON.

## Features

- ✅ **Socket.IO Integration**: WebSocket + HTTP long-polling fallback
- ✅ **JWT Authentication**: Validates all connections with JWT middleware
- ✅ **Room Auto-Subscription**: Automatic user and team room joining
- ✅ **gRPC-to-Socket.IO**: SendToUser and BroadcastToRoom handlers
- ✅ **Connection Management**: Thread-safe tracking of all connections
- ✅ **Room Management**: Efficient room membership tracking
- ✅ **CORS Support**: Configurable allowed origins
- ✅ **Connection Limits**: Configurable max concurrent connections
- ✅ **Transport Fallback**: Automatic WebSocket → Polling fallback

## Architecture

```
cmd/main.go                     # Dual server (Socket.IO + gRPC)
internal/
  ├── types.go                  # Connection & Room data structures
  ├── auth_middleware.go        # JWT authentication (HIGH PRIORITY)
  ├── socketio_server.go        # Socket.IO with auto-subscription
  ├── grpc_handlers.go          # gRPC-to-Socket.IO bridge
  └── config/config.go          # Configuration
proto/notifications/v1/service.proto # gRPC definitions
```

## Critical Implementation Details

### 1. Socket.IO with HTTP Long-Polling ✅

```go
server := socketio.NewServer(&engineio.Options{
    Transports: []transport.Transport{
        &websocket.Transport{...},      // Primary
        &polling.Transport{...},         // Fallback
    },
    PingTimeout:  60 * time.Second,
    PingInterval: 25 * time.Second,
})
```

**Automatic fallback for restrictive networks!**

### 2. JWT Authentication Middleware (HIGH PRIORITY) ✅

```go
func (s *SocketIOServer) handleConnect(conn socketio.Conn) error {
    // Authenticate connection (REQUIRED)
    claims, err := s.authMW.Authenticate(conn)
    if err != nil {
        return fmt.Errorf("authentication failed: %w", err)
    }
    
    // Connection rejected if JWT invalid
    // ...
}
```

**JWT Validation:**
- Extracts token from query params or headers
- Validates signature using shared secret
- Checks expiration
- Extracts user_id and team_id claims
- **Rejects connection on failure**

### 3. Room Auto-Subscription ✅

```go
// After successful authentication
userRoom := fmt.Sprintf("user_%s", claims.UserID)
teamRoom := fmt.Sprintf("team_%s", claims.TeamID)

// Join user room (private)
conn.Join(userRoom)
s.roomManager.JoinRoom(socketID, userRoom, "user")

// Join team room (shared)
conn.Join(teamRoom)
s.roomManager.JoinRoom(socketID, teamRoom, "team")

// Emit connection_ready event
conn.Emit("connection_ready", map[string]interface{}{
    "user_id": claims.UserID,
    "rooms":   []string{userRoom, teamRoom},
})
```

**Automatic room subscription on connect!**

### 4. gRPC-to-Socket.IO Logic ✅

**SendToUser Handler:**
```go
func (s *NotificationsServer) SendToUser(ctx context.Context, req *pb.SendToUserRequest) (*pb.SendToUserResponse, error) {
    // Get user's private room
    userRoom := "user_" + req.UserId
    
    // Broadcast to user's room
    s.socketServer.GetServer().BroadcastToRoom("/", userRoom, req.EventType, payload)
    
    return &pb.SendToUserResponse{Delivered: true}, nil
}
```

**BroadcastToRoom Handler:**
```go
func (s *NotificationsServer) BroadcastToRoom(ctx context.Context, req *pb.BroadcastToRoomRequest) (*pb.BroadcastToRoomResponse, error) {
    // Broadcast to specified room (e.g., team_xyz-789)
    s.socketServer.GetServer().BroadcastToRoom("/", req.RoomId, req.EventType, payload)
    
    return &pb.BroadcastToRoomResponse{Delivered: true}, nil
}
```

## Environment Variables

```bash
# Required
JWT_SECRET=your-jwt-secret-key-here

# Server Ports
SOCKETIO_PORT=3000              # Socket.IO HTTP server
GRPC_PORT=50055                 # gRPC server

# CORS (Important!)
ALLOWED_ORIGINS=http://localhost:3000,https://app.example.com

# Connection Limits
MAX_CONNECTIONS=10000

# Socket.IO Configuration
PING_TIMEOUT_SECONDS=60
PING_INTERVAL_SECONDS=25
ENABLE_WEBSOCKET=true
ENABLE_POLLING=true

# Logging
LOG_LEVEL=info
```

## Quick Start

```bash
# 1. Set up environment
cd app/services/notifications-service
cp .env.example .env
# Edit .env with your JWT secret

# 2. Generate proto code
make proto

# 3. Build
make build

# 4. Run
make run

# Service listens on:
# - Socket.IO: http://localhost:3000
# - gRPC: localhost:50055
```

## Client Connection (Frontend)

### JavaScript/TypeScript Client

```javascript
import io from 'socket.io-client';

// Get JWT token from your auth system
const token = 'your-jwt-token-here';

// Connect with JWT authentication
const socket = io('http://localhost:3000', {
  auth: {
    token: token  // JWT passed in auth object
  },
  transports: ['websocket', 'polling'],  // WebSocket first, polling fallback
  reconnection: true,
  reconnectionDelay: 1000,
  reconnectionAttempts: 5
});

// Connection established
socket.on('connect', () => {
  console.log('Connected:', socket.id);
});

// Connection ready (after authentication)
socket.on('connection_ready', (data) => {
  console.log('Ready:', data);
  // data = { user_id: "...", rooms: ["user_...", "team_..."], socket_id: "..." }
});

// Listen for notifications
socket.on('notification', (data) => {
  console.log('Notification received:', data);
});

// Listen for alerts
socket.on('alert', (data) => {
  console.log('Alert received:', data);
});

// Disconnection
socket.on('disconnect', (reason) => {
  console.log('Disconnected:', reason);
});

// Authentication error
socket.on('connect_error', (error) => {
  console.error('Connection error:', error.message);
});
```

## Backend Usage (gRPC)

### Send to Specific User

```go
import pb "github.com/haunted-saas/notifications-service/proto/notifications/v1"

conn, _ := grpc.Dial("notifications-service:50055", grpc.WithInsecure())
client := pb.NewNotificationsServiceClient(conn)

// Send notification to user
resp, err := client.SendToUser(ctx, &pb.SendToUserRequest{
    UserId:    "user_123",
    EventType: "notification",
    PayloadJson: `{
        "title": "New Message",
        "body": "You have a new message from Alice",
        "priority": "high"
    }`,
})

fmt.Printf("Delivered to %d connections\n", resp.ConnectionCount)
```

### Broadcast to Team

```go
// Broadcast to entire team
resp, err := client.BroadcastToRoom(ctx, &pb.BroadcastToRoomRequest{
    RoomId:    "team_xyz-789",
    EventType: "alert",
    PayloadJson: `{
        "type": "system",
        "message": "System maintenance in 10 minutes"
    }`,
})

fmt.Printf("Delivered to %d recipients\n", resp.RecipientCount)
```

### Check if User is Connected

```go
resp, err := client.IsUserConnected(ctx, &pb.IsUserConnectedRequest{
    UserId: "user_123",
})

if resp.IsConnected {
    fmt.Printf("User has %d active connections\n", resp.ConnectionCount)
} else {
    fmt.Println("User is offline")
}
```

## How It Works

### Connection Flow

```
1. Client connects to Socket.IO server
2. JWT extracted from auth object or headers
3. JWT validated (signature, expiration, claims)
4. If invalid → Connection rejected
5. If valid → Extract user_id and team_id
6. Auto-subscribe to:
   - user_{user_id} (private room)
   - team_{team_id} (team room)
7. Emit connection_ready event to client
8. Client ready to receive messages
```

### Message Delivery Flow

```
Backend Service → gRPC SendToUser
                ↓
Notifications Service
                ↓
Lookup user's room (user_{user_id})
                ↓
Socket.IO BroadcastToRoom
                ↓
All user's connections receive message
```

### Transport Fallback

```
Client attempts WebSocket
       ↓
WebSocket blocked by firewall?
       ↓
Automatic fallback to HTTP long-polling
       ↓
Connection established
```

## Security

### JWT Authentication

**Token Requirements:**
- Must contain `user_id` claim
- Must contain `team_id` claim (optional)
- Must be signed with shared secret
- Must not be expired

**Validation:**
- Signature verified using HMAC
- Expiration checked
- Required claims validated
- **Connection rejected on any failure**

### CORS Protection

```bash
ALLOWED_ORIGINS=https://app.example.com,https://admin.example.com
```

Only specified origins can connect.

### Connection Limits

```bash
MAX_CONNECTIONS=10000
```

Prevents resource exhaustion.

## Monitoring

### Connection Stats

```go
resp, err := client.GetConnectionStats(ctx, &pb.GetConnectionStatsRequest{})

fmt.Printf("Total connections: %d\n", resp.TotalConnections)
fmt.Printf("WebSocket: %d\n", resp.WebsocketConnections)
fmt.Printf("Polling: %d\n", resp.PollingConnections)
```

### Logs to Watch

```
INFO  connection attempt socket_id=abc123 transport=websocket
INFO  authentication successful user_id=user_123 team_id=team_456
INFO  connection established user_room=user_user_123 team_room=team_team_456
INFO  message sent to user user_id=user_123 connection_count=2
INFO  connection closed socket_id=abc123 duration=5m30s
```

## Production Checklist

- [ ] Set production JWT secret (same as user-auth-service)
- [ ] Configure allowed origins (frontend domains)
- [ ] Set appropriate connection limits
- [ ] Enable both WebSocket and polling
- [ ] Configure ping timeout/interval
- [ ] Set up monitoring for connection count
- [ ] Test JWT authentication
- [ ] Test transport fallback
- [ ] Verify room auto-subscription
- [ ] Test message delivery

## Integration with Other Services

### User Auth Service

```go
// Use same JWT secret for token validation
JWT_SECRET=same-secret-as-user-auth-service
```

### Any Backend Service

```go
// Send notification from any service
notificationsClient.SendToUser(ctx, &pb.SendToUserRequest{
    UserId:    userID,
    EventType: "notification",
    PayloadJson: notificationJSON,
})
```

## Troubleshooting

**Connection rejected:**
- Check JWT token is valid
- Verify JWT_SECRET matches user-auth-service
- Check token expiration
- Verify user_id claim exists

**WebSocket not working:**
- Check firewall/proxy settings
- Verify ENABLE_WEBSOCKET=true
- Service should fallback to polling automatically

**Messages not received:**
- Check user is connected (IsUserConnected)
- Verify room names match (user_{user_id})
- Check event type matches client listener
- Verify payload JSON is valid

---

**Status**: ✅ COMPLETE - Production-ready with Socket.IO  
**Security**: JWT authentication on all connections  
**Reliability**: Automatic WebSocket → Polling fallback  
**Ready for**: Real-time notifications, chat, live updates
