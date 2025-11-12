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

	"github.com/haunted-saas/llm-gateway-service/internal"
	"github.com/haunted-saas/llm-gateway-service/internal/config"
	pb "github.com/haunted-saas/llm-gateway-service/proto/llm/v1"
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

	logger.Info("🎃 Starting LLM Gateway Service",
		zap.Int("grpc_port", cfg.Server.GRPCPort),
		zap.String("prompts_dir", cfg.Prompts.Directory),
		zap.Bool("test_mode", cfg.LLM.TestMode))

	// Initialize prompt loader
	promptLoader, err := internal.NewPromptLoader(cfg.Prompts.Directory, cfg.Prompts.WatchMode, logger)
	if err != nil {
		logger.Fatal("Failed to create prompt loader", zap.Error(err))
	}
	defer promptLoader.Close()

	// Load all prompts
	if err := promptLoader.LoadAllPrompts(); err != nil {
		logger.Fatal("Failed to load prompts", zap.Error(err))
	}

	// Start file watcher
	if cfg.Prompts.WatchMode {
		if err := promptLoader.WatchForChanges(); err != nil {
			logger.Fatal("Failed to start file watcher", zap.Error(err))
		}
	}

	// Initialize LLM providers
	router := internal.NewLLMRouter(cfg.LLM.DefaultProvider, logger)

	// Register OpenAI provider
	openaiProvider, err := internal.NewOpenAIProvider(cfg.LLM.OpenAIAPIKey, cfg.LLM.TestMode, logger)
	if err != nil {
		logger.Fatal("Failed to create OpenAI provider", zap.Error(err))
	}
	router.RegisterProvider(openaiProvider)

	if cfg.LLM.TestMode {
		logger.Warn("⚠️  TEST MODE ENABLED - Using mock LLM responses")
	}

	// Initialize usage tracker
	usageTracker := internal.NewUsageTracker(cfg.Analytics.UsageStoreMaxSize, logger)

	// Initialize gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	// Register LLM gateway service
	llmService := internal.NewLLMGatewayServer(promptLoader, router, usageTracker, logger)
	pb.RegisterLLMGatewayServiceServer(grpcServer, llmService)

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

	logger.Info("Shutting down server...")

	// Graceful shutdown
	grpcServer.GracefulStop()
	logger.Info("Server stopped")
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
