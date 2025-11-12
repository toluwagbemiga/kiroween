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

	"github.com/haunted-saas/billing-service/internal"
	"github.com/haunted-saas/billing-service/internal/config"
	"github.com/haunted-saas/billing-service/internal/db"
	pb "github.com/haunted-saas/billing-service/proto/billing/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("🎃 Starting Billing Service",
		zap.Int("grpc_port", cfg.Server.GRPCPort),
		zap.Int("http_port", cfg.Server.HTTPPort))

	// Initialize database
	gormDB, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		zapLogger.Fatal("Failed to get database instance", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	zapLogger.Info("✓ Database connected")

	// Run migrations
	if err := runMigrations(gormDB); err != nil {
		zapLogger.Warn("Failed to run migrations", zap.Error(err))
	}

	// Initialize store
	store := db.NewStore(gormDB)

	// Initialize Stripe client
	stripeClient := internal.NewStripeClient(cfg.Stripe.APIKey)
	zapLogger.Info("✓ Stripe client initialized")

	// Initialize gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(zapLogger)),
	)

	// Register billing service
	billingService := internal.NewBillingServiceServer(stripeClient, store, zapLogger)
	pb.RegisterBillingServiceServer(grpcServer, billingService)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for development
	reflection.Register(grpcServer)

	// Start gRPC server
	grpcLis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GRPCPort))
	if err != nil {
		zapLogger.Fatal("Failed to listen on gRPC port", zap.Error(err))
	}

	go func() {
		zapLogger.Info("🚀 gRPC server started", zap.String("address", grpcLis.Addr().String()))
		if err := grpcServer.Serve(grpcLis); err != nil {
			zapLogger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// Initialize webhook handler
	webhookHandler := internal.NewWebhookHandler(stripeClient, store, cfg.Stripe.WebhookSecret, zapLogger)

	// Start HTTP server for webhooks
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/webhooks/stripe", webhookHandler.HandleWebhook)
	httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort),
		Handler:      httpMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		zapLogger.Info("🚀 HTTP server started (webhooks)", zap.String("address", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Failed to serve HTTP", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down servers...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		zapLogger.Error("HTTP server shutdown error", zap.Error(err))
	}

	grpcServer.GracefulStop()
	zapLogger.Info("Servers stopped")
}

func runMigrations(database *gorm.DB) error {
	return database.AutoMigrate(
		&db.Plan{},
		&db.Subscription{},
		&db.WebhookEvent{},
	)
}

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
			logger.Info("gRPC call succeeded",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration))
		}

		return resp, err
	}
}
