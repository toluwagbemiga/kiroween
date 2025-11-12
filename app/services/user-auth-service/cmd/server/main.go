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

	"github.com/haunted-saas/user-auth-service/internal/auth"
	"github.com/haunted-saas/user-auth-service/internal/config"
	"github.com/haunted-saas/user-auth-service/internal/database"
	"github.com/haunted-saas/user-auth-service/internal/handler"
	"github.com/haunted-saas/user-auth-service/internal/logging"
	"github.com/haunted-saas/user-auth-service/internal/repository"
	"github.com/haunted-saas/user-auth-service/internal/service"
	pb "github.com/haunted-saas/user-auth-service/proto/userauth/v1"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
	logger, err := logging.NewLogger(logging.GetLogLevel())
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("🎃 Starting User Auth Service",
		zap.Int("port", cfg.Server.GRPCPort),
		zap.Int("bcrypt_cost", cfg.Security.BcryptCost))

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	logger.Info("✓ Database connected")

	// Run migrations
	if err := database.RunMigrations(db, "migrations"); err != nil {
		logger.Warn("Failed to run migrations", zap.Error(err))
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logger.Info("✓ Redis connected")

	// Initialize token manager
	tokenManager, err := auth.NewTokenManager(
		cfg.JWT.PrivateKeyPath,
		cfg.JWT.PublicKeyPath,
		cfg.JWT.Expiration,
	)
	if err != nil {
		logger.Fatal("Failed to initialize token manager", zap.Error(err))
	}

	logger.Info("✓ Token manager initialized")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	sessionRepo := repository.NewSessionRepository(redisClient)
	rateLimiterRepo := repository.NewRateLimiterRepository(redisClient)
	permCacheRepo := repository.NewPermissionCacheRepository(redisClient)
	resetRepo := repository.NewPasswordResetRepository(redisClient)

	// Initialize services
	authService := service.NewAuthService(
		userRepo,
		roleRepo,
		sessionRepo,
		rateLimiterRepo,
		resetRepo,
		tokenManager,
		cfg,
		logger,
	)

	rbacService := service.NewRBACService(
		userRepo,
		roleRepo,
		permRepo,
		permCacheRepo,
		sessionRepo,
		cfg,
		logger,
	)

	// Initialize handler
	authHandler := handler.NewAuthHandler(authService, rbacService)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	// Register services
	pb.RegisterUserAuthServiceServer(grpcServer, authHandler)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for development
	reflection.Register(grpcServer)

	// Start server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GRPCPort))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	logger.Info("🚀 User Auth Service started",
		zap.String("address", lis.Addr().String()))

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	grpcServer.GracefulStop()
	logger.Info("Server stopped")
}

func loggingInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
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
			logger.Info("gRPC call succeeded",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration))
		}

		return resp, err
	}
}
