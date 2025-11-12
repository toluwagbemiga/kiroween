package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haunted-saas/notifications-service/internal"
	"github.com/haunted-saas/notifications-service/internal/config"
	pb "github.com/haunted-saas/notifications-service/proto/notifications/v1"
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

	logger.Info("🎃 Starting Notifications Service",
		zap.Int("grpc_port", cfg.Server.GRPCPort),
		zap.Int("socketio_port", cfg.Server.SocketIOPort),
		zap.Int("max_connections", cfg.SocketIO.MaxConnections),
		zap.Strings("allowed_origins", cfg.SocketIO.AllowedOrigins))

	// Initialize JWT authentication middleware
	authMW := internal.NewAuthMiddleware(cfg.Authentication.JWTSecret, logger)
	logger.Info("✓ JWT authentication middleware initialized")

	// Initialize Socket.IO server
	socketServer, err := internal.NewSocketIOServer(
		authMW,
		cfg.SocketIO.MaxConnections,
		cfg.SocketIO.AllowedOrigins,
		logger,
	)
	if err != nil {
		logger.Fatal("Failed to create Socket.IO server", zap.Error(err))
	}
	logger.Info("✓ Socket.IO server initialized")

	// Start Socket.IO server
	go func() {
		logger.Info("🚀 Socket.IO server starting",
			zap.String("address", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.SocketIOPort)))

		if err := http.ListenAndServe(
			fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.SocketIOPort),
			socketServer.GetServer(),
		); err != nil {
			logger.Fatal("Failed to serve Socket.IO", zap.Error(err))
		}
	}()

	// Initialize gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	// Register notifications service
	notificationsService := internal.NewNotificationsServer(socketServer, logger)
	pb.RegisterNotificationsServiceServer(grpcServer, notificationsService)

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
		logger.Info("🚀 gRPC server started", zap.String("address", lis.Addr().String()))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down servers...")

	// Graceful shutdown
	grpcServer.GracefulStop()
	logger.Info("✓ gRPC server stopped")

	// Note: Socket.IO server shutdown would require additional implementation
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
