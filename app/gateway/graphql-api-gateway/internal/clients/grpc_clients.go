package clients

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	analyticsv1 "github.com/haunted-saas/analytics-service/proto/analytics/v1"
	billingv1 "github.com/haunted-saas/billing-service/proto/billing/v1"
	featureflagsv1 "github.com/haunted-saas/feature-flags-service/proto/featureflags/v1"
	llmv1 "github.com/haunted-saas/llm-gateway-service/proto/llm/v1"
	notificationsv1 "github.com/haunted-saas/notifications-service/proto/notifications/v1"
	userauthv1 "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
)

// GRPCClients holds all gRPC client connections
type GRPCClients struct {
	UserAuth      userauthv1.UserAuthServiceClient
	Billing       billingv1.BillingServiceClient
	LLMGateway    llmv1.LLMGatewayServiceClient
	Notifications notificationsv1.NotificationsServiceClient
	Analytics     analyticsv1.AnalyticsServiceClient
	FeatureFlags  featureflagsv1.FeatureFlagsServiceClient

	// Store connections for cleanup
	conns []*grpc.ClientConn
}

// ServicesConfig holds service addresses
type ServicesConfig struct {
	UserAuthService      string
	BillingService       string
	LLMGatewayService    string
	NotificationsService string
	AnalyticsService     string
	FeatureFlagsService  string
}

// NewGRPCClients initializes all gRPC clients
func NewGRPCClients(config ServicesConfig, logger *zap.Logger) (*GRPCClients, error) {
	clients := &GRPCClients{
		conns: make([]*grpc.ClientConn, 0, 6),
	}

	// Initialize User Auth Service client
	logger.Info("connecting to user-auth-service", zap.String("address", config.UserAuthService))
	userAuthConn, err := grpc.Dial(
		config.UserAuthService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user-auth-service: %w", err)
	}
	clients.conns = append(clients.conns, userAuthConn)
	clients.UserAuth = userauthv1.NewUserAuthServiceClient(userAuthConn)
	logger.Info("✓ connected to user-auth-service")

	// Initialize Billing Service client
	logger.Info("connecting to billing-service", zap.String("address", config.BillingService))
	billingConn, err := grpc.Dial(
		config.BillingService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to billing-service: %w", err)
	}
	clients.conns = append(clients.conns, billingConn)
	clients.Billing = billingv1.NewBillingServiceClient(billingConn)
	logger.Info("✓ connected to billing-service")

	// Initialize LLM Gateway Service client
	logger.Info("connecting to llm-gateway-service", zap.String("address", config.LLMGatewayService))
	llmConn, err := grpc.Dial(
		config.LLMGatewayService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to llm-gateway-service: %w", err)
	}
	clients.conns = append(clients.conns, llmConn)
	clients.LLMGateway = llmv1.NewLLMGatewayServiceClient(llmConn)
	logger.Info("✓ connected to llm-gateway-service")

	// Initialize Notifications Service client
	logger.Info("connecting to notifications-service", zap.String("address", config.NotificationsService))
	notificationsConn, err := grpc.Dial(
		config.NotificationsService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notifications-service: %w", err)
	}
	clients.conns = append(clients.conns, notificationsConn)
	clients.Notifications = notificationsv1.NewNotificationsServiceClient(notificationsConn)
	logger.Info("✓ connected to notifications-service")

	// Initialize Analytics Service client
	logger.Info("connecting to analytics-service", zap.String("address", config.AnalyticsService))
	analyticsConn, err := grpc.Dial(
		config.AnalyticsService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to analytics-service: %w", err)
	}
	clients.conns = append(clients.conns, analyticsConn)
	clients.Analytics = analyticsv1.NewAnalyticsServiceClient(analyticsConn)
	logger.Info("✓ connected to analytics-service")

	// Initialize Feature Flags Service client
	logger.Info("connecting to feature-flags-service", zap.String("address", config.FeatureFlagsService))
	featureFlagsConn, err := grpc.Dial(
		config.FeatureFlagsService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to feature-flags-service: %w", err)
	}
	clients.conns = append(clients.conns, featureFlagsConn)
	clients.FeatureFlags = featureflagsv1.NewFeatureFlagsServiceClient(featureFlagsConn)
	logger.Info("✓ connected to feature-flags-service")

	logger.Info("✓ all gRPC clients initialized successfully")

	return clients, nil
}

// Close closes all gRPC connections
func (c *GRPCClients) Close() error {
	for _, conn := range c.conns {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
