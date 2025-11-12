package internal

import (
	"context"
	"encoding/json"

	pb "github.com/haunted-saas/notifications-service/proto/notifications/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotificationsServer implements the gRPC service
type NotificationsServer struct {
	pb.UnimplementedNotificationsServiceServer
	socketServer *SocketIOServer
	logger       *zap.Logger
}

// NewNotificationsServer creates a new notifications server
func NewNotificationsServer(socketServer *SocketIOServer, logger *zap.Logger) *NotificationsServer {
	return &NotificationsServer{
		socketServer: socketServer,
		logger:       logger,
	}
}

// SendToUser sends a message to a specific user
func (s *NotificationsServer) SendToUser(ctx context.Context, req *pb.SendToUserRequest) (*pb.SendToUserResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.EventType == "" {
		return nil, status.Error(codes.InvalidArgument, "event_type is required")
	}

	// Get user's room
	userRoom := "user_" + req.UserId

	// Parse payload
	var payload interface{}
	if req.PayloadJson != "" {
		if err := json.Unmarshal([]byte(req.PayloadJson), &payload); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid payload_json")
		}
	}

	// Get user connections
	connections := s.socketServer.GetConnectionManager().GetUserConnections(req.UserId)
	connectionCount := len(connections)

	// Emit to user's room
	s.socketServer.GetServer().BroadcastToRoom("/", userRoom, req.EventType, payload)

	s.logger.Info("message sent to user",
		zap.String("user_id", req.UserId),
		zap.String("event_type", req.EventType),
		zap.Int("connection_count", connectionCount),
		zap.String("correlation_id", req.CorrelationId))

	return &pb.SendToUserResponse{
		Delivered:       true,
		ConnectionCount: int32(connectionCount),
	}, nil
}

// SendToUsers sends a message to multiple users
func (s *NotificationsServer) SendToUsers(ctx context.Context, req *pb.SendToUsersRequest) (*pb.SendToUsersResponse, error) {
	// Validate request
	if len(req.UserIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_ids is required")
	}
	if req.EventType == "" {
		return nil, status.Error(codes.InvalidArgument, "event_type is required")
	}

	// Parse payload
	var payload interface{}
	if req.PayloadJson != "" {
		if err := json.Unmarshal([]byte(req.PayloadJson), &payload); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid payload_json")
		}
	}

	// Send to each user
	connectionsByUser := make(map[string]int32)
	deliveredCount := 0

	for _, userID := range req.UserIds {
		userRoom := "user_" + userID
		connections := s.socketServer.GetConnectionManager().GetUserConnections(userID)
		connectionCount := len(connections)

		if connectionCount > 0 {
			s.socketServer.GetServer().BroadcastToRoom("/", userRoom, req.EventType, payload)
			deliveredCount++
		}

		connectionsByUser[userID] = int32(connectionCount)
	}

	s.logger.Info("message sent to multiple users",
		zap.Int("total_users", len(req.UserIds)),
		zap.Int("delivered_count", deliveredCount),
		zap.String("event_type", req.EventType),
		zap.String("correlation_id", req.CorrelationId))

	return &pb.SendToUsersResponse{
		TotalUsers:        int32(len(req.UserIds)),
		DeliveredCount:    int32(deliveredCount),
		ConnectionsByUser: connectionsByUser,
	}, nil
}

// BroadcastToRoom broadcasts a message to a room
func (s *NotificationsServer) BroadcastToRoom(ctx context.Context, req *pb.BroadcastToRoomRequest) (*pb.BroadcastToRoomResponse, error) {
	// Validate request
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}
	if req.EventType == "" {
		return nil, status.Error(codes.InvalidArgument, "event_type is required")
	}

	// Parse payload
	var payload interface{}
	if req.PayloadJson != "" {
		if err := json.Unmarshal([]byte(req.PayloadJson), &payload); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid payload_json")
		}
	}

	// Get room members
	members := s.socketServer.GetRoomManager().GetRoomMembers(req.RoomId)
	recipientCount := len(members)

	// Broadcast to room
	s.socketServer.GetServer().BroadcastToRoom("/", req.RoomId, req.EventType, payload)

	s.logger.Info("message broadcast to room",
		zap.String("room_id", req.RoomId),
		zap.String("event_type", req.EventType),
		zap.Int("recipient_count", recipientCount),
		zap.String("correlation_id", req.CorrelationId))

	return &pb.BroadcastToRoomResponse{
		Delivered:      true,
		RecipientCount: int32(recipientCount),
	}, nil
}

// IsUserConnected checks if a user is connected
func (s *NotificationsServer) IsUserConnected(ctx context.Context, req *pb.IsUserConnectedRequest) (*pb.IsUserConnectedResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	connections := s.socketServer.GetConnectionManager().GetUserConnections(req.UserId)
	socketIDs := make([]string, len(connections))
	for i, conn := range connections {
		socketIDs[i] = conn.SocketID
	}

	return &pb.IsUserConnectedResponse{
		IsConnected:     len(connections) > 0,
		ConnectionCount: int32(len(connections)),
		SocketIds:       socketIDs,
	}, nil
}

// GetConnectionStats returns connection statistics
func (s *NotificationsServer) GetConnectionStats(ctx context.Context, req *pb.GetConnectionStatsRequest) (*pb.GetConnectionStatsResponse, error) {
	connManager := s.socketServer.GetConnectionManager()

	totalConns := connManager.GetConnectionCount()
	transportCounts := connManager.GetConnectionsByTransport()
	teamCounts := connManager.GetConnectionsByTeam()

	// Filter by team if requested
	if req.TeamId != "" {
		filteredTeamCounts := make(map[string]int32)
		if count, exists := teamCounts[req.TeamId]; exists {
			filteredTeamCounts[req.TeamId] = int32(count)
		}
		teamCounts = map[string]int{req.TeamId: int(filteredTeamCounts[req.TeamId])}
	}

	// Convert to proto format
	transportCountsProto := make(map[string]int32)
	for transport, count := range transportCounts {
		transportCountsProto[transport] = int32(count)
	}

	teamCountsProto := make(map[string]int32)
	for teamID, count := range teamCounts {
		teamCountsProto[teamID] = int32(count)
	}

	websocketCount := int32(0)
	pollingCount := int32(0)
	if count, exists := transportCountsProto["websocket"]; exists {
		websocketCount = count
	}
	if count, exists := transportCountsProto["polling"]; exists {
		pollingCount = count
	}

	return &pb.GetConnectionStatsResponse{
		TotalConnections:      int32(totalConns),
		WebsocketConnections:  websocketCount,
		PollingConnections:    pollingCount,
		ConnectionsByTeam:     teamCountsProto,
		ConnectionsByTransport: transportCountsProto,
	}, nil
}

// DisconnectUser disconnects all connections for a user
func (s *NotificationsServer) DisconnectUser(ctx context.Context, req *pb.DisconnectUserRequest) (*pb.DisconnectUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	connections := s.socketServer.GetConnectionManager().GetUserConnections(req.UserId)
	disconnectedCount := 0

	for _, conn := range connections {
		if conn.Conn != nil {
			conn.Conn.Close()
			disconnectedCount++
		}
	}

	s.logger.Info("user disconnected",
		zap.String("user_id", req.UserId),
		zap.String("reason", req.Reason),
		zap.Int("disconnected_count", disconnectedCount))

	return &pb.DisconnectUserResponse{
		DisconnectedCount: int32(disconnectedCount),
	}, nil
}
