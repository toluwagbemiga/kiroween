package internal

import (
	"context"
	"time"

	pb "github.com/haunted-saas/feature-flags-service/proto/featureflags/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// FeatureFlagsServer implements the gRPC service
type FeatureFlagsServer struct {
	pb.UnimplementedFeatureFlagsServiceServer
	unleashClient *UnleashClient
	logger        *zap.Logger
}

// NewFeatureFlagsServer creates a new feature flags server
func NewFeatureFlagsServer(unleashClient *UnleashClient, logger *zap.Logger) *FeatureFlagsServer {
	return &FeatureFlagsServer{
		unleashClient: unleashClient,
		logger:        logger,
	}
}

// IsFeatureEnabled checks if a feature is enabled for the given context
// This is the core proxy function - extremely simple and fast
func (s *FeatureFlagsServer) IsFeatureEnabled(ctx context.Context, req *pb.IsFeatureEnabledRequest) (*pb.IsFeatureEnabledResponse, error) {
	// Validate request
	if req.FeatureName == "" {
		return nil, status.Error(codes.InvalidArgument, "feature_name is required")
	}

	// Parse properties JSON
	properties, err := ParsePropertiesJSON(req.PropertiesJson)
	if err != nil {
		s.logger.Warn("invalid properties JSON",
			zap.String("feature_name", req.FeatureName),
			zap.String("user_id", req.UserId),
			zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid properties_json")
	}

	// Extract metadata from gRPC context
	remoteAddr, userAgent, sessionID := s.extractMetadata(ctx)

	// Build feature context
	featureContext := &FeatureContext{
		UserID:     req.UserId,
		TeamID:     req.TeamId,
		Properties: properties,
		RemoteAddr: remoteAddr,
		UserAgent:  userAgent,
		SessionID:  sessionID,
	}

	// Call Unleash SDK (in-memory cache lookup - no network call!)
	enabled := s.unleashClient.IsFeatureEnabled(req.FeatureName, featureContext)

	s.logger.Debug("feature flag evaluated",
		zap.String("feature_name", req.FeatureName),
		zap.String("user_id", req.UserId),
		zap.String("team_id", req.TeamId),
		zap.Bool("enabled", enabled))

	// Return result immediately
	return &pb.IsFeatureEnabledResponse{
		Enabled: enabled,
	}, nil
}

// GetFeatureVariant gets the variant for a feature flag
func (s *FeatureFlagsServer) GetFeatureVariant(ctx context.Context, req *pb.GetFeatureVariantRequest) (*pb.GetFeatureVariantResponse, error) {
	// Validate request
	if req.FeatureName == "" {
		return nil, status.Error(codes.InvalidArgument, "feature_name is required")
	}

	// Parse properties JSON
	properties, err := ParsePropertiesJSON(req.PropertiesJson)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid properties_json")
	}

	// Extract metadata
	remoteAddr, userAgent, sessionID := s.extractMetadata(ctx)

	// Build feature context
	featureContext := &FeatureContext{
		UserID:     req.UserId,
		TeamID:     req.TeamId,
		Properties: properties,
		RemoteAddr: remoteAddr,
		UserAgent:  userAgent,
		SessionID:  sessionID,
	}

	// Get variant from Unleash SDK (in-memory cache)
	variant := s.unleashClient.GetVariant(req.FeatureName, featureContext)

	s.logger.Debug("feature variant evaluated",
		zap.String("feature_name", req.FeatureName),
		zap.String("user_id", req.UserId),
		zap.String("variant_name", variant.Name),
		zap.Bool("enabled", variant.Enabled))

	// Convert variant payload to JSON
	payloadJSON := "{}"
	if variant.Payload.Value != "" {
		payloadJSON = variant.Payload.Value
	}

	return &pb.GetFeatureVariantResponse{
		Enabled:     variant.Enabled,
		VariantName: variant.Name,
		PayloadJson: payloadJSON,
	}, nil
}

// ListFeatures lists all available features (for debugging/admin)
func (s *FeatureFlagsServer) ListFeatures(ctx context.Context, req *pb.ListFeaturesRequest) (*pb.ListFeaturesResponse, error) {
	// Get features from Unleash SDK
	features := s.unleashClient.GetFeatureToggles()

	// Convert to proto format
	protoFeatures := make([]*pb.Feature, len(features))
	for i, feature := range features {
		protoFeatures[i] = &pb.Feature{
			Name:        feature.Name,
			Description: feature.Description,
			Enabled:     feature.Enabled,
			CreatedAt:   feature.CreatedAt.Format(time.RFC3339),
		}
	}

	s.logger.Debug("features listed", zap.Int("count", len(features)))

	return &pb.ListFeaturesResponse{
		Features: protoFeatures,
	}, nil
}

// GetServiceHealth returns the health status of the service
func (s *FeatureFlagsServer) GetServiceHealth(ctx context.Context, req *pb.GetServiceHealthRequest) (*pb.GetServiceHealthResponse, error) {
	isReady := s.unleashClient.IsReady()

	status := "healthy"
	if !isReady {
		status = "not_ready"
	}

	return &pb.GetServiceHealthResponse{
		Status:  status,
		IsReady: isReady,
	}, nil
}

// extractMetadata extracts useful metadata from gRPC context
func (s *FeatureFlagsServer) extractMetadata(ctx context.Context) (remoteAddr, userAgent, sessionID string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// Extract remote address
		if addrs := md.Get("x-forwarded-for"); len(addrs) > 0 {
			remoteAddr = addrs[0]
		} else if addrs := md.Get("x-real-ip"); len(addrs) > 0 {
			remoteAddr = addrs[0]
		}

		// Extract user agent
		if agents := md.Get("user-agent"); len(agents) > 0 {
			userAgent = agents[0]
		}

		// Extract session ID
		if sessions := md.Get("x-session-id"); len(sessions) > 0 {
			sessionID = sessions[0]
		}
	}

	return remoteAddr, userAgent, sessionID
}


// GetUserFeatures gets all enabled features for a user
// This is useful for frontend to get all flags at once
func (s *FeatureFlagsServer) GetUserFeatures(ctx context.Context, req *pb.GetUserFeaturesRequest) (*pb.GetUserFeaturesResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Parse properties JSON
	properties, err := ParsePropertiesJSON(req.PropertiesJson)
	if err != nil {
		s.logger.Warn("invalid properties JSON",
			zap.String("user_id", req.UserId),
			zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid properties_json")
	}

	// Build feature context
	featureContext := make(map[string]interface{})
	featureContext["userId"] = req.UserId
	featureContext["teamId"] = req.TeamId
	for k, v := range properties {
		featureContext[k] = v
	}

	// STUB: Return empty list for now
	// TODO: Implement full Unleash integration
	enabledFeatures := []string{}

	s.logger.Debug("user features evaluated",
		zap.String("user_id", req.UserId),
		zap.Int("enabled_features", len(enabledFeatures)))

	return &pb.GetUserFeaturesResponse{
		EnabledFeatures: enabledFeatures,
	}, nil
}
