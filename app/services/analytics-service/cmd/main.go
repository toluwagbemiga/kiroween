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

	"github.com/haunted-saas/analytics-service/internal"
	"github.com/haunted-saas/analytics-service/internal/config"
	pb "github.com/haunted-saas/analytics-service/proto/analytics/v1"
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

	logger.Info("🎃 Starting Analytics Service",
		zap.Int("grpc_port", cfg.Server.GRPCPort),
		zap.Int("batch_size", cfg.Analytics.BatchSize),
		zap.Int("flush_interval_sec", cfg.Analytics.FlushIntervalSec),
		zap.Bool("test_mode", cfg.Analytics.TestMode))

	// Initialize batch queue
	queue := internal.NewBatchQueue(cfg.Analytics.BatchSize)
	logger.Info("✓ Batch queue initialized", zap.Int("max_size", cfg.Analytics.BatchSize))

	// Initialize external provider (Mixpanel)
	provider := internal.NewMixpanelProvider(cfg.Analytics.MixpanelAPIKey, cfg.Analytics.TestMode, logger)
	logger.Info("✓ Analytics provider initialized", zap.String("provider", provider.GetName()))

	if cfg.Analytics.TestMode {
		logger.Warn("⚠️  TEST MODE ENABLED - Events will not be sent to external provider")
	}

	// Initialize retry config
	retryConfig := &internal.RetryConfig{
		MaxAttempts:   cfg.Analytics.MaxRetryAttempts,
		InitialDelay:  time.Duration(cfg.Analytics.InitialRetryDelay) * time.Millisecond,
		MaxDelay:      time.Duration(cfg.Analytics.MaxRetryDelay) * time.Millisecond,
		BackoffFactor: 2.0,
	}

	// Initialize batch worker
	flushInterval := time.Duration(cfg.Analytics.FlushIntervalSec) * time.Second
	worker := internal.NewBatchWorker(queue, provider, flushInterval, retryConfig, logger)
	
	// Start batch worker (concurrent goroutine)
	worker.Start()
	logger.Info("✓ Batch worker started")

	// Initialize gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	// Register analytics service
	analyticsService := internal.NewAnalyticsServer(queue, logger)
	pb.RegisterAnalyticsServiceServer(grpcServer, analyticsService)

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

	// Wait for interrupt signal (SIGTERM or SIGINT)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("Shutdown signal received", zap.String("signal", sig.String()))

	// Graceful shutdown
	logger.Info("Shutting down gracefully...")

	// Stop accepting new gRPC requests
	grpcServer.GracefulStop()
	logger.Info("✓ gRPC server stopped")

	// Stop batch worker (triggers final flush)
	worker.Stop()
	logger.Info("✓ Batch worker stopped (final flush completed)")

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
