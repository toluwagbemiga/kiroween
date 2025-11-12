# Real-Time Notifications Integration Guide

## Overview

The notifications-service has been successfully integrated into the frontend application using Socket.IO for real-time, bidirectional communication. This enables:

- **Real-time Notifications**: Instant push notifications to users
- **Toast Notifications**: Automatic display of notifications as toast messages
- **Persistent Connection**: Maintains WebSocket connection with automatic reconnection
- **JWT Authentication**: Secure connection using user's authentication token
- **Fallback Support**: Gracefully falls back to HTTP long-polling if WebSocket fails

## Architecture

### Backend (Notifications Service)

**Socket.IO Server:**
- Location: `app/services/notifications-service`
- Port: 8085 (default)
- Supports WebSocket and HTTP long-polling
- JWT authentication via Socket.IO auth middleware
- Events: `new_notification`, `broadcast`, `notification_read`

### Frontend Integration

#### 1. SocketProvider (`src/lib/SocketProvider.tsx`)

Manages Socket.IO connection and handles real-time events:

```typescript
const { socket, isConnected } = useSocket();
```

**Features:**
- Automatic connection when user is authenticated
- JWT token passed in Socket.IO auth header
- Automatic reconnection with exponential backoff
- Event listeners for notifications
- Automatic toast display for incoming notifications
- Connection status tracking

#### 2. Connection Flow

1. User logs in → JWT token stored
2. SocketProvider detects authentication
3. Creates Socket.IO connection with JWT
4. Backend validates JWT and establishes connection
5. Listens for `new_notification` events
6. Displays toast notifications automatically

#### 3. Notification Types

The system supports four notification types:
- `info`: Blue toast (default)
- `success`: Green toast
- `warning`: Yellow toast
- `error`: Red toast

## Usage Examples

### Automatic Notifications

Notifications are automatically displayed when received from the server. No additional code needed in components!

```typescript
// Backend sends notification
socket.emit('new_notification', {
  id: 'notif-123',
  title: 'New Message',
  message: 'You have a new message from John',
  type: 'info',
  timestamp: new Date().toISOString()
});

// Frontend automatically displays toast
// ✓ Toast appears with title and message
```

### Check Connection Status

```typescript
import { useSocket } from '@/lib/SocketProvider';

function MyComponent() {
  const { isConnected } = useSocket();

  return (
    <div>
      {isConnected ? (
        <span>🟢 Connected</span>
      ) : (
        <span>🔴 Disconnected</span>
      )}
    </div>
  );
}
```

### Send Notifications Programmatically

```typescript
import { useSendNotification } from '@/lib/SocketProvider';

function AdminPanel() {
  const { sendNotification, isConnected } = useSendNotification();

  const notifyUser = () => {
    sendNotification(
      'user-123',
      'System Update',
      'The system will be updated in 5 minutes',
      'warning',
      { updateId: 'update-456' }
    );
  };

  return (
    <button onClick={notifyUser} disabled={!isConnected}>
      Send Notification
    </button>
  );
}
```

### Access Socket Directly

```typescript
import { useSocket } from '@/lib/SocketProvider';

function AdvancedComponent() {
  const { socket } = useSocket();

  useEffect(() => {
    if (!socket) return;

    // Listen for custom events
    socket.on('custom_event', (data) => {
      console.log('Custom event received:', data);
    });

    // Emit custom events
    socket.emit('custom_action', { action: 'ping' });

    return () => {
      socket.off('custom_event');
    };
  }, [socket]);

  return <div>Advanced Socket.IO Usage</div>;
}
```

## Notification Payload Structure

```typescript
interface NotificationPayload {
  id: string;                    // Unique notification ID
  title?: string;                // Optional title
  message: string;               // Notification message (required)
  type: 'info' | 'success' | 'warning' | 'error';
  data?: Record<string, any>;    // Optional additional data
  timestamp: string;             // ISO 8601 timestamp
}
```

## Backend Configuration

### Environment Variables

```bash
# Notifications Service
NOTIFICATIONS_PORT=8085
NOTIFICATIONS_SOCKET_PORT=8085

# JWT Secret (must match user-auth-service)
JWT_SECRET=your-secret-key

# CORS Origins
CORS_ORIGINS=http://localhost:3000,https://your-domain.com
```

### Sending Notifications from Backend

#### Via gRPC (Recommended)

```go
import notificationsv1 "github.com/haunted-saas/notifications-service/proto/notifications/v1"

// Send notification via gRPC
_, err := notificationsClient.SendNotification(ctx, &notificationsv1.SendNotificationRequest{
    UserId:  "user-123",
    Title:   "New Message",
    Message: "You have a new message",
    Type:    "info",
    DataJson: `{"messageId": "msg-456"}`,
})
```

#### Via Socket.IO (Direct)

```go
// In notifications-service
server.BroadcastToUser(userID, "new_notification", NotificationPayload{
    ID:        uuid.New().String(),
    Title:     "New Message",
    Message:   "You have a new message",
    Type:      "info",
    Timestamp: time.Now().Format(time.RFC3339),
})
```

## Frontend Configuration

### Environment Variables

Create `.env.local` file:

```bash
# Notifications Service URL
NEXT_PUBLIC_NOTIFICATIONS_URL=http://localhost:8085

# For production
# NEXT_PUBLIC_NOTIFICATIONS_URL=https://notifications.your-domain.com
```

### Connection Options

The Socket.IO client is configured with:

```typescript
{
  auth: {
    token: jwtToken,  // JWT from authentication
  },
  transports: ['websocket', 'polling'],  // Try WebSocket first
  reconnection: true,
  reconnectionDelay: 1000,
  reconnectionDelayMax: 5000,
  reconnectionAttempts: 5,
}
```

## Components

### NotificationStatus

Shows connection status in the UI:

```typescript
import { NotificationStatus } from '@/components/NotificationStatus';

<NotificationStatus />
// Shows: 🟢 Live (when connected)
```

### DetailedNotificationStatus

Shows detailed connection info (for debugging):

```typescript
import { DetailedNotificationStatus } from '@/components/NotificationStatus';

<DetailedNotificationStatus />
// Shows: Connection status, Socket ID, Transport type
```

## Events

### Received Events (Frontend Listens)

| Event | Description | Payload |
|-------|-------------|---------|
| `connect` | Socket connected | - |
| `disconnect` | Socket disconnected | `reason: string` |
| `new_notification` | New notification received | `NotificationPayload` |
| `broadcast` | Broadcast message to all users | `NotificationPayload` |
| `notification_read` | Notification marked as read | `{ notificationId: string }` |

### Emitted Events (Frontend Sends)

| Event | Description | Payload |
|-------|-------------|---------|
| `send_notification` | Send notification to user | `NotificationPayload` |
| `mark_as_read` | Mark notification as read | `{ notificationId: string }` |

## Best Practices

### 1. Handle Connection States

```typescript
const { isConnected } = useSocket();

// Disable features when disconnected
<button disabled={!isConnected}>
  Send Notification
</button>

// Show warning when disconnected
{!isConnected && (
  <Alert variant="warning">
    Real-time notifications are currently unavailable
  </Alert>
)}
```

### 2. Clean Up Event Listeners

```typescript
useEffect(() => {
  if (!socket) return;

  const handleEvent = (data) => {
    // Handle event
  };

  socket.on('custom_event', handleEvent);

  return () => {
    socket.off('custom_event', handleEvent);
  };
}, [socket]);
```

### 3. Don't Block on Socket Operations

```typescript
// Good - non-blocking
const sendNotification = () => {
  if (socket?.connected) {
    socket.emit('send_notification', data);
  }
};

// Bad - blocking
const sendNotification = async () => {
  await new Promise((resolve) => {
    socket.emit('send_notification', data, resolve);
  });
};
```

### 4. Use Notification Types Appropriately

```typescript
// Success - for completed actions
{ type: 'success', message: 'Profile updated successfully' }

// Error - for failures
{ type: 'error', message: 'Failed to save changes' }

// Warning - for important info
{ type: 'warning', message: 'Your session will expire in 5 minutes' }

// Info - for general updates
{ type: 'info', message: 'New feature available' }
```

## Troubleshooting

### Connection Issues

**Problem:** Socket not connecting

**Solutions:**
1. Check if notifications-service is running on port 8085
2. Verify `NEXT_PUBLIC_NOTIFICATIONS_URL` environment variable
3. Check browser console for connection errors
4. Verify JWT token is valid
5. Check CORS configuration on backend

**Debug:**
```typescript
// Add to SocketProvider for debugging
console.log('Socket URL:', socketUrl);
console.log('Token:', token ? 'Present' : 'Missing');
console.log('Socket state:', socket?.connected);
```

### Notifications Not Appearing

**Problem:** Connected but no toast notifications

**Solutions:**
1. Check browser console for event logs
2. Verify notification payload structure
3. Check if ToastProvider is in component tree
4. Verify event name is `new_notification`

**Debug:**
```typescript
socket.on('new_notification', (payload) => {
  console.log('Notification received:', payload);
  // Should see this in console when notification arrives
});
```

### Reconnection Issues

**Problem:** Socket keeps disconnecting

**Solutions:**
1. Check network stability
2. Verify backend is not crashing
3. Check JWT token expiration
4. Increase reconnection attempts
5. Check backend logs for errors

### CORS Errors

**Problem:** CORS policy blocking connection

**Solution:** Update backend CORS configuration:
```go
// In notifications-service
server := socketio.NewServer(&engineio.Options{
    Transports: []transport.Transport{
        &polling.Transport{
            CheckOrigin: func(r *http.Request) bool {
                origin := r.Header.Get("Origin")
                return origin == "http://localhost:3000" || 
                       origin == "https://your-domain.com"
            },
        },
        &websocket.Transport{
            CheckOrigin: func(r *http.Request) bool {
                origin := r.Header.Get("Origin")
                return origin == "http://localhost:3000" || 
                       origin == "https://your-domain.com"
            },
        },
    },
})
```

## Performance Considerations

### Connection Overhead

- Socket.IO maintains a persistent connection
- Minimal overhead after initial connection
- Automatic heartbeat keeps connection alive
- Falls back to polling if WebSocket unavailable

### Scalability

For production with multiple servers:

1. **Use Redis Adapter:**
   ```go
   // In notifications-service
   adapter := socketio.NewRedisAdapter(redisClient)
   server.SetAdapter(adapter)
   ```

2. **Sticky Sessions:**
   Configure load balancer for sticky sessions or use Redis adapter

3. **Horizontal Scaling:**
   Multiple notification service instances can share state via Redis

## Security

### Authentication

- JWT token required for connection
- Token validated on connection
- Invalid tokens rejected immediately
- Token expiration handled automatically

### Authorization

- Users only receive their own notifications
- Broadcast messages sent to all connected users
- Admin-only events can be restricted

### Best Practices

1. **Never expose sensitive data** in notification payloads
2. **Validate all data** on backend before sending
3. **Rate limit** notification sending
4. **Log all notification events** for audit trail
5. **Use HTTPS/WSS** in production

## Testing

### Manual Testing

1. **Start services:**
   ```bash
   # Terminal 1: Notifications service
   cd app/services/notifications-service
   go run cmd/main.go

   # Terminal 2: Frontend
   cd app/frontend
   npm run dev
   ```

2. **Login to application**
3. **Check connection status** in header (should show "Live")
4. **Send test notification** via backend or admin panel
5. **Verify toast appears**

### Automated Testing

```typescript
// Mock Socket.IO for testing
jest.mock('socket.io-client', () => ({
  io: jest.fn(() => ({
    on: jest.fn(),
    emit: jest.fn(),
    disconnect: jest.fn(),
  })),
}));
```

## Next Steps

Consider adding:
- **Notification Center**: List of all notifications
- **Notification Preferences**: Let users configure notification types
- **Sound Alerts**: Play sound when notification arrives
- **Browser Notifications**: Use Notification API for desktop notifications
- **Notification History**: Store and display past notifications
- **Read/Unread Status**: Track which notifications user has seen
- **Notification Actions**: Add buttons to notifications (e.g., "View", "Dismiss")
