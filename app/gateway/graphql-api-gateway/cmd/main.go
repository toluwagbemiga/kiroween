package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/haunted-saas/graphql-api-gateway/internal/clients"
	"github.com/haunted-saas/graphql-api-gateway/internal/config"
	"github.com/haunted-saas/graphql-api-gateway/internal/dataloader"
	"github.com/haunted-saas/graphql-api-gateway/internal/generated"
	"github.com/haunted-saas/graphql-api-gateway/internal/middleware"
	"github.com/haunted-saas/graphql-api-gateway/internal/resolvers"
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

	logger.Info("🎃 Starting GraphQL API Gateway",
		zap.Int("port", cfg.Server.Port),
		zap.String("env", cfg.Server.Env))

	// Initialize gRPC clients
	logger.Info("initializing gRPC clients...")
	grpcClients, err := clients.NewGRPCClients(clients.ServicesConfig{
		UserAuthService:      cfg.Services.UserAuthService,
		BillingService:       cfg.Services.BillingService,
		LLMGatewayService:    cfg.Services.LLMGatewayService,
		NotificationsService: cfg.Services.NotificationsService,
		AnalyticsService:     cfg.Services.AnalyticsService,
		FeatureFlagsService:  cfg.Services.FeatureFlagsService,
	}, logger)
	if err != nil {
		logger.Fatal("Failed to initialize gRPC clients", zap.Error(err))
	}
	defer grpcClients.Close()

	logger.Info("✓ all gRPC clients initialized")

	// Initialize dataloaders
	loaders := dataloader.NewLoaders(dataloader.Clients{
		UserAuth: grpcClients.UserAuth,
		Billing:  grpcClients.Billing,
	}, logger)

	// Initialize resolvers
	resolver := resolvers.NewResolver(grpcClients, logger)

	// Create GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(grpcClients.UserAuth, logger)

	// Setup HTTP router
	mux := http.NewServeMux()

	// GraphQL endpoint with auth middleware and dataloaders
	mux.Handle("/graphql", 
		authMiddleware.Middleware(
			dataloader.Middleware(loaders)(srv),
		),
	)

	// GraphQL Playground (only in development)
	if cfg.Server.Env == "development" {
		mux.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
		logger.Info("GraphQL Playground enabled at http://localhost:" + fmt.Sprint(cfg.Server.Port))
	}

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Setup CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Configure this properly in production
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      corsHandler.Handler(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("🚀 GraphQL API Gateway started",
			zap.String("address", httpServer.Addr),
			zap.String("graphql_endpoint", "/graphql"))

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("✓ HTTP server stopped")
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
