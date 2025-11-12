package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haunted-saas/feature-flags-service/internal"
	"github.com/haunted-saas/feature-flags-service/internal/config"
	pb "github.com/haunted-saas/feature-flags-service/proto/featureflags/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := initLogger(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("🎃 Starting Feature Flags Service",
		zap.Int("grpc_port", cfg.Server.GRPCPort),
		zap.String("unleash_server", cfg.Unleash.ServerURL),
		zap.String("app_name", cfg.Unleash.AppName),
		zap.Duration("refresh_interval", cfg.Unleash.RefreshInterval))

	// Initialize Unleash client (HIGH PRIORITY - SDK initialization)
	unleashConfig := &internal.UnleashConfig{
		ServerURL:       cfg.Unleash.ServerURL,
		APIToken:        cfg.Unleash.APIToken,
		AppName:         cfg.Unleash.AppName,
		InstanceID:      cfg.Unleash.InstanceID,
		RefreshInterval: cfg.Unleash.RefreshInterval,
		MetricsInterval: cfg.Unleash.MetricsInterval,
		DisableMetrics:  cfg.Unleash.DisableMetrics,
	}

	unleashClient, err := internal.NewUnleashClient(unleashConfig, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Unleash client", zap.Error(err))
	}
	defer unleashClient.Close()

	// Wait for Unleash client to be ready (load feature flags into memory)
	logger.Info("Waiting for Unleash client to load feature flags...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := unleashClient.WaitForReady(ctx); err != nil {
		logger.Fatal("Unleash client failed to become ready", zap.Error(err))
	}

	logger.Info("✓ Unleash client ready - feature flags loaded into in-memory cache")

	// Initialize gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	// Register feature flags service
	featureFlagsService := internal.NewFeatureFlagsServer(unleashClient, logger)
	pb.RegisterFeatureFlagsServiceServer(grpcServer, featureFlagsService)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for development
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GRPCPort))
	if err != nil {
		logger.Fatal("Failed to listen on gRPC port", zap.Error(err))
	}

	go func() {
		logger.Info("🚀 gRPC server started (high-speed proxy ready)",
			zap.String("address", lis.Addr().String()))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	grpcServer.GracefulStop()
	logger.Info("✓ gRPC server stopped")

	logger.Info("Shutdown complete")
}

// initLogger initializes the logger
func initLogger(level, format string) (*zap.Logger, error) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	if format == "console" {
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return config.Build()
}

// loggingInterceptor logs gRPC requests
func loggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			logger.Error("gRPC call failed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.Error(err))
		} else {
			logger.Debug("gRPC call succeeded",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration))
		}

		return resp, err
	}
}
