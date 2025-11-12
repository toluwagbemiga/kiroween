package internal

import (
	"sync"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

// Connection represents a Socket.IO connection
type Connection struct {
	SocketID    string
	UserID      string
	TeamID      string
	Transport   string // "websocket" or "polling"
	ConnectedAt time.Time
	LastSeen    time.Time
	Conn        socketio.Conn
}

// ConnectionManager manages all active connections
type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*Connection // socket_id -> Connection
	userConns   map[string][]string    // user_id -> []socket_id
	teamConns   map[string][]string    // team_id -> []socket_id
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Connection),
		userConns:   make(map[string][]string),
		teamConns:   make(map[string][]string),
	}
}

// AddConnection adds a connection
func (m *ConnectionManager) AddConnection(conn *Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connections[conn.SocketID] = conn

	// Add to user connections
	m.userConns[conn.UserID] = append(m.userConns[conn.UserID], conn.SocketID)

	// Add to team connections
	if conn.TeamID != "" {
		m.teamConns[conn.TeamID] = append(m.teamConns[conn.TeamID], conn.SocketID)
	}
}

// RemoveConnection removes a connection
func (m *ConnectionManager) RemoveConnection(socketID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, exists := m.connections[socketID]
	if !exists {
		return
	}

	// Remove from user connections
	m.userConns[conn.UserID] = removeFromSlice(m.userConns[conn.UserID], socketID)
	if len(m.userConns[conn.UserID]) == 0 {
		delete(m.userConns, conn.UserID)
	}

	// Remove from team connections
	if conn.TeamID != "" {
		m.teamConns[conn.TeamID] = removeFromSlice(m.teamConns[conn.TeamID], socketID)
		if len(m.teamConns[conn.TeamID]) == 0 {
			delete(m.teamConns, conn.TeamID)
		}
	}

	delete(m.connections, socketID)
}

// GetConnection gets a connection by socket ID
func (m *ConnectionManager) GetConnection(socketID string) (*Connection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, exists := m.connections[socketID]
	return conn, exists
}

// GetUserConnections gets all connections for a user
func (m *ConnectionManager) GetUserConnections(userID string) []*Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	socketIDs := m.userConns[userID]
	conns := make([]*Connection, 0, len(socketIDs))
	for _, socketID := range socketIDs {
		if conn, exists := m.connections[socketID]; exists {
			conns = append(conns, conn)
		}
	}
	return conns
}

// IsUserConnected checks if a user is connected
func (m *ConnectionManager) IsUserConnected(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.userConns[userID]) > 0
}

// GetConnectionCount returns total connection count
func (m *ConnectionManager) GetConnectionCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connections)
}

// GetConnectionsByTransport returns connection counts by transport
func (m *ConnectionManager) GetConnectionsByTransport() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	counts := make(map[string]int)
	for _, conn := range m.connections {
		counts[conn.Transport]++
	}
	return counts
}

// GetConnectionsByTeam returns connection counts by team
func (m *ConnectionManager) GetConnectionsByTeam() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	counts := make(map[string]int)
	for teamID, socketIDs := range m.teamConns {
		counts[teamID] = len(socketIDs)
	}
	return counts
}

// Room represents a Socket.IO room
type Room struct {
	ID        string
	Type      string // "user", "team", "custom"
	Members   map[string]bool
	CreatedAt time.Time
}

// RoomManager manages Socket.IO rooms
type RoomManager struct {
	mu     sync.RWMutex
	rooms  map[string]*Room
	server *socketio.Server
}

// NewRoomManager creates a new room manager
func NewRoomManager(server *socketio.Server) *RoomManager {
	return &RoomManager{
		rooms:  make(map[string]*Room),
		server: server,
	}
}

// JoinRoom adds a socket to a room
func (m *RoomManager) JoinRoom(socketID, roomID, roomType string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create room if it doesn't exist
	if _, exists := m.rooms[roomID]; !exists {
		m.rooms[roomID] = &Room{
			ID:        roomID,
			Type:      roomType,
			Members:   make(map[string]bool),
			CreatedAt: time.Now(),
		}
	}

	// Add member to room
	m.rooms[roomID].Members[socketID] = true

	// Join Socket.IO room
	if m.server != nil {
		// Note: Socket.IO room joining happens at the connection level
		// This is tracked here for management purposes
	}

	return nil
}

// LeaveRoom removes a socket from a room
func (m *RoomManager) LeaveRoom(socketID, roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if room, exists := m.rooms[roomID]; exists {
		delete(room.Members, socketID)

		// Clean up empty rooms
		if len(room.Members) == 0 {
			delete(m.rooms, roomID)
		}
	}
}

// LeaveAllRooms removes a socket from all rooms
func (m *RoomManager) LeaveAllRooms(socketID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for roomID, room := range m.rooms {
		delete(room.Members, socketID)

		// Clean up empty rooms
		if len(room.Members) == 0 {
			delete(m.rooms, roomID)
		}
	}
}

// GetRoomMembers gets all members of a room
func (m *RoomManager) GetRoomMembers(roomID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return []string{}
	}

	members := make([]string, 0, len(room.Members))
	for socketID := range room.Members {
		members = append(members, socketID)
	}
	return members
}

// GetRoomCount returns the number of rooms
func (m *RoomManager) GetRoomCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.rooms)
}

// Helper function to remove an item from a slice
func removeFromSlice(slice []string, item string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
