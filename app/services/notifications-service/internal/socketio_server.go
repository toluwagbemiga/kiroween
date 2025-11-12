package internal

import (
	"fmt"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"go.uber.org/zap"
)

// SocketIOServer manages the Socket.IO server
type SocketIOServer struct {
	server      *socketio.Server
	connManager *ConnectionManager
	roomManager *RoomManager
	authMW      *AuthMiddleware
	logger      *zap.Logger
	maxConns    int
}

// NewSocketIOServer creates a new Socket.IO server
func NewSocketIOServer(
	authMW *AuthMiddleware,
	maxConns int,
	allowedOrigins []string,
	logger *zap.Logger,
) (*SocketIOServer, error) {
	// Create Socket.IO server with WebSocket and polling transports
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return checkOrigin(r, allowedOrigins)
				},
			},
			&polling.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return checkOrigin(r, allowedOrigins)
				},
			},
		},
		PingTimeout:  60 * time.Second,
		PingInterval: 25 * time.Second,
	})

	connManager := NewConnectionManager()
	roomManager := NewRoomManager(server)

	s := &SocketIOServer{
		server:      server,
		connManager: connManager,
		roomManager: roomManager,
		authMW:      authMW,
		logger:      logger,
		maxConns:    maxConns,
	}

	// Register event handlers
	s.registerHandlers()

	return s, nil
}

// registerHandlers registers Socket.IO event handlers
func (s *SocketIOServer) registerHandlers() {
	// Connection handler
	s.server.OnConnect("/", func(conn socketio.Conn) error {
		return s.handleConnect(conn)
	})

	// Disconnection handler
	s.server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		s.handleDisconnect(conn, reason)
	})

	// Error handler
	s.server.OnError("/", func(conn socketio.Conn, err error) {
		s.handleError(conn, err)
	})
}

// handleConnect handles new connections
func (s *SocketIOServer) handleConnect(conn socketio.Conn) error {
	socketID := conn.ID()

	s.logger.Info("connection attempt",
		zap.String("socket_id", socketID))

	// Check connection limit
	if s.connManager.GetConnectionCount() >= s.maxConns {
		s.logger.Warn("connection limit reached",
			zap.Int("current", s.connManager.GetConnectionCount()),
			zap.Int("max", s.maxConns))
		return fmt.Errorf("connection limit reached")
	}

	// Authenticate connection (HIGH PRIORITY)
	claims, err := s.authMW.Authenticate(conn)
	if err != nil {
		s.logger.Warn("authentication failed",
			zap.String("socket_id", socketID),
			zap.Error(err))
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create connection object
	connection := &Connection{
		SocketID:    socketID,
		UserID:      claims.UserID,
		TeamID:      claims.TeamID,
		Transport:   "websocket", // Default transport
		ConnectedAt: time.Now(),
		LastSeen:    time.Now(),
		Conn:        conn,
	}

	// Add to connection manager
	s.connManager.AddConnection(connection)

	// Auto-subscribe to rooms (HIGH PRIORITY)
	userRoom := fmt.Sprintf("user_%s", claims.UserID)
	teamRoom := fmt.Sprintf("team_%s", claims.TeamID)

	// Join user room
	conn.Join(userRoom)
	s.roomManager.JoinRoom(socketID, userRoom, "user")

	// Join team room (if team ID exists)
	if claims.TeamID != "" {
		conn.Join(teamRoom)
		s.roomManager.JoinRoom(socketID, teamRoom, "team")
	}

	s.logger.Info("connection established",
		zap.String("socket_id", socketID),
		zap.String("user_id", claims.UserID),
		zap.String("team_id", claims.TeamID),
		zap.String("user_room", userRoom),
		zap.String("team_room", teamRoom))

	// Emit connection_ready event to client
	conn.Emit("connection_ready", map[string]interface{}{
		"user_id": claims.UserID,
		"rooms":   []string{userRoom, teamRoom},
		"socket_id": socketID,
	})

	return nil
}

// handleDisconnect handles disconnections
func (s *SocketIOServer) handleDisconnect(conn socketio.Conn, reason string) {
	socketID := conn.ID()

	connection, exists := s.connManager.GetConnection(socketID)
	if !exists {
		s.logger.Warn("disconnect for unknown connection", zap.String("socket_id", socketID))
		return
	}

	duration := time.Since(connection.ConnectedAt)

	s.logger.Info("connection closed",
		zap.String("socket_id", socketID),
		zap.String("user_id", connection.UserID),
		zap.String("reason", reason),
		zap.Duration("duration", duration))

	// Leave all rooms
	s.roomManager.LeaveAllRooms(socketID)

	// Remove from connection manager
	s.connManager.RemoveConnection(socketID)
}

// handleError handles Socket.IO errors
func (s *SocketIOServer) handleError(conn socketio.Conn, err error) {
	socketID := conn.ID()

	s.logger.Error("socket error",
		zap.String("socket_id", socketID),
		zap.Error(err))
}

// GetServer returns the underlying Socket.IO server
func (s *SocketIOServer) GetServer() *socketio.Server {
	return s.server
}

// GetConnectionManager returns the connection manager
func (s *SocketIOServer) GetConnectionManager() *ConnectionManager {
	return s.connManager
}

// GetRoomManager returns the room manager
func (s *SocketIOServer) GetRoomManager() *RoomManager {
	return s.roomManager
}

// checkOrigin checks if the origin is allowed
func checkOrigin(r *http.Request, allowedOrigins []string) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true // Allow requests without origin (same-origin)
	}

	for _, allowed := range allowedOrigins {
		if origin == allowed || allowed == "*" {
			return true
		}
	}

	return false
}
