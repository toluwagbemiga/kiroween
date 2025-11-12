package internal

import (
	"context"
	"time"

	"github.com/google/uuid"
	pb "github.com/haunted-saas/analytics-service/proto/analytics/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AnalyticsServer implements the gRPC service
type AnalyticsServer struct {
	pb.UnimplementedAnalyticsServiceServer
	queue  *BatchQueue
	logger *zap.Logger
}

// NewAnalyticsServer creates a new analytics server
func NewAnalyticsServer(queue *BatchQueue, logger *zap.Logger) *AnalyticsServer {
	return &AnalyticsServer{
		queue:  queue,
		logger: logger,
	}
}

// TrackEvent tracks an analytics event (NON-BLOCKING)
func (s *AnalyticsServer) TrackEvent(ctx context.Context, req *pb.TrackEventRequest) (*pb.TrackEventResponse, error) {
	// Validate request
	if req.EventName == "" {
		return nil, status.Error(codes.InvalidArgument, "event_name is required")
	}

	// Generate event ID
	eventID := uuid.New().String()

	// Convert properties
	properties := make(map[string]interface{})
	for key, propValue := range req.Properties {
		properties[key] = convertPropertyValue(propValue)
	}

	// Create event
	event := Event{
		ID:         eventID,
		EventName:  req.EventName,
		UserID:     req.UserId,
		Properties: properties,
		Timestamp:  time.Now(),
		CreatedAt:  time.Now(),
	}

	// Override timestamp if provided
	if req.Timestamp > 0 {
		event.Timestamp = time.Unix(req.Timestamp, 0)
	}

	// Add to queue (NON-BLOCKING - just adds to in-memory queue)
	shouldFlush := s.queue.Add(event)

	s.logger.Debug("event queued",
		zap.String("event_id", eventID),
		zap.String("event_name", req.EventName),
		zap.String("user_id", req.UserId),
		zap.Bool("should_flush", shouldFlush),
		zap.Int("queue_size", s.queue.Size()))

	// Return immediately (non-blocking)
	return &pb.TrackEventResponse{
		Success: true,
		EventId: eventID,
	}, nil
}

// IdentifyUser identifies a user (NON-BLOCKING)
func (s *AnalyticsServer) IdentifyUser(ctx context.Context, req *pb.IdentifyUserRequest) (*pb.IdentifyUserResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Convert properties
	properties := make(map[string]interface{})
	for key, propValue := range req.Properties {
		properties[key] = convertPropertyValue(propValue)
	}

	// Create identify event (special event type)
	event := Event{
		ID:         uuid.New().String(),
		EventName:  "$identify",
		UserID:     req.UserId,
		Properties: properties,
		Timestamp:  time.Now(),
		CreatedAt:  time.Now(),
	}

	// Add to queue (NON-BLOCKING)
	s.queue.Add(event)

	s.logger.Debug("user identified",
		zap.String("user_id", req.UserId),
		zap.Int("property_count", len(properties)))

	// Return immediately (non-blocking)
	return &pb.IdentifyUserResponse{
		Success: true,
	}, nil
}

// GetEventCount returns event counts (placeholder for future implementation)
func (s *AnalyticsServer) GetEventCount(ctx context.Context, req *pb.GetEventCountRequest) (*pb.GetEventCountResponse, error) {
	// TODO: Implement querying from external provider or local database
	s.logger.Warn("GetEventCount not yet implemented")
	return &pb.GetEventCountResponse{
		Counts: []*pb.EventCount{},
	}, nil
}

// GetUserCount returns unique user count (placeholder for future implementation)
func (s *AnalyticsServer) GetUserCount(ctx context.Context, req *pb.GetUserCountRequest) (*pb.GetUserCountResponse, error) {
	// TODO: Implement querying from external provider or local database
	s.logger.Warn("GetUserCount not yet implemented")
	return &pb.GetUserCountResponse{
		UniqueUsers: 0,
	}, nil
}

// HealthCheck returns service health status
func (s *AnalyticsServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: "healthy",
	}, nil
}

// convertPropertyValue converts a proto PropertyValue to interface{}
func convertPropertyValue(pv *pb.PropertyValue) interface{} {
	if pv == nil {
		return nil
	}

	switch v := pv.Value.(type) {
	case *pb.PropertyValue_StringValue:
		return v.StringValue
	case *pb.PropertyValue_NumberValue:
		return v.NumberValue
	case *pb.PropertyValue_BoolValue:
		return v.BoolValue
	default:
		return nil
	}
}
