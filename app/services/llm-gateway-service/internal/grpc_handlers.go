package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pb "github.com/haunted-saas/llm-gateway-service/proto/llm/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LLMGatewayServer implements the gRPC service
type LLMGatewayServer struct {
	pb.UnimplementedLLMGatewayServiceServer
	promptLoader   *PromptLoader
	router         *LLMRouter
	usageTracker   *UsageTracker
	logger         *zap.Logger
	defaultTimeout time.Duration
	maxTimeout     time.Duration
}

// NewLLMGatewayServer creates a new LLM gateway server
func NewLLMGatewayServer(
	promptLoader *PromptLoader,
	router *LLMRouter,
	usageTracker *UsageTracker,
	logger *zap.Logger,
) *LLMGatewayServer {
	return &LLMGatewayServer{
		promptLoader:   promptLoader,
		router:         router,
		usageTracker:   usageTracker,
		logger:         logger,
		defaultTimeout: 30 * time.Second,
		maxTimeout:     120 * time.Second,
	}
}

// CallPrompt executes a prompt with variables
func (s *LLMGatewayServer) CallPrompt(ctx context.Context, req *pb.CallPromptRequest) (*pb.CallPromptResponse, error) {
	startTime := time.Now()
	requestID := generateRequestID()

	s.logger.Info("CallPrompt request received",
		zap.String("prompt_path", req.PromptPath),
		zap.String("calling_service", req.CallingService),
		zap.String("request_id", requestID),
		zap.String("correlation_id", req.CorrelationId))

	// Validate request
	if req.PromptPath == "" {
		return nil, status.Error(codes.InvalidArgument, "prompt_path is required")
	}

	// Validate prompt path (prevent directory traversal)
	if strings.Contains(req.PromptPath, "..") {
		return nil, status.Error(codes.InvalidArgument, "invalid prompt path")
	}

	// Load prompt from cache
	prompt, err := s.promptLoader.GetPrompt(req.PromptPath)
	if err != nil {
		s.logger.Warn("prompt not found",
			zap.String("prompt_path", req.PromptPath),
			zap.Error(err))
		return nil, status.Error(codes.NotFound, fmt.Sprintf("prompt not found: %s", req.PromptPath))
	}

	// Substitute variables
	renderedPrompt, err := s.substituteVariables(prompt, req.VariablesJson)
	if err != nil {
		s.logger.Error("variable substitution failed",
			zap.String("prompt_path", req.PromptPath),
			zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("variable substitution failed: %v", err))
	}

	// Validate and apply parameters
	params, err := s.validateParameters(req.Parameters, prompt)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid parameters: %v", err))
	}

	// Determine timeout
	timeout := s.defaultTimeout
	if req.TimeoutSeconds > 0 {
		if req.TimeoutSeconds < 5 || req.TimeoutSeconds > 120 {
			return nil, status.Error(codes.InvalidArgument, "timeout must be between 5 and 120 seconds")
		}
		timeout = time.Duration(req.TimeoutSeconds) * time.Second
	}

	// Build LLM request
	llmReq := &LLMRequest{
		Prompt:     renderedPrompt,
		Model:      req.Model,
		Parameters: params,
		Timeout:    timeout,
		RequestID:  requestID,
	}

	// Apply metadata defaults if available
	if prompt.Metadata != nil {
		if llmReq.Model == "" && prompt.Metadata.DefaultModel != "" {
			llmReq.Model = prompt.Metadata.DefaultModel
		}
		if params.Temperature == 0 && prompt.Metadata.Temperature != nil {
			params.Temperature = *prompt.Metadata.Temperature
		}
		if params.MaxTokens == 0 && prompt.Metadata.MaxTokens != nil {
			params.MaxTokens = *prompt.Metadata.MaxTokens
		}
	}

	// Route to LLM provider
	llmResp, err := s.router.Route(ctx, llmReq)
	if err != nil {
		s.logger.Error("LLM call failed",
			zap.String("prompt_path", req.PromptPath),
			zap.String("request_id", requestID),
			zap.Error(err))

		// Track failed usage
		s.trackUsageAsync(&UsageEvent{
			RequestID:      requestID,
			PromptPath:     req.PromptPath,
			CallingService: req.CallingService,
			Provider:       req.Provider,
			Model:          req.Model,
			Timestamp:      time.Now(),
			Success:        false,
			ErrorMessage:   err.Error(),
		})

		// Map error to gRPC code
		if isRateLimitError(err) {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, status.Error(codes.DeadlineExceeded, "request timeout")
		}
		return nil, status.Error(codes.Internal, "LLM provider error")
	}

	responseTime := time.Since(startTime)

	// Track successful usage
	s.trackUsageAsync(&UsageEvent{
		RequestID:        requestID,
		PromptPath:       req.PromptPath,
		CallingService:   req.CallingService,
		Provider:         req.Provider,
		Model:            llmResp.Model,
		PromptTokens:     llmResp.TokenUsage.PromptTokens,
		CompletionTokens: llmResp.TokenUsage.CompletionTokens,
		TotalTokens:      llmResp.TokenUsage.TotalTokens,
		ResponseTimeMs:   responseTime.Milliseconds(),
		Timestamp:        time.Now(),
		Success:          true,
	})

	s.logger.Info("CallPrompt completed",
		zap.String("prompt_path", req.PromptPath),
		zap.String("request_id", requestID),
		zap.String("model", llmResp.Model),
		zap.Int32("total_tokens", llmResp.TokenUsage.TotalTokens),
		zap.Duration("response_time", responseTime))

	// Build response
	return &pb.CallPromptResponse{
		ResponseText: llmResp.Text,
		TokenUsage: &pb.TokenUsage{
			PromptTokens:     llmResp.TokenUsage.PromptTokens,
			CompletionTokens: llmResp.TokenUsage.CompletionTokens,
			TotalTokens:      llmResp.TokenUsage.TotalTokens,
		},
		ModelUsed:      llmResp.Model,
		RequestId:      requestID,
		ResponseTimeMs: responseTime.Milliseconds(),
	}, nil
}

// substituteVariables substitutes variables in a prompt template
func (s *LLMGatewayServer) substituteVariables(prompt *Prompt, variablesJSON string) (string, error) {
	// Parse variables JSON
	var variables map[string]interface{}
	if variablesJSON != "" {
		if err := json.Unmarshal([]byte(variablesJSON), &variables); err != nil {
			return "", fmt.Errorf("invalid JSON: %w", err)
		}
	} else {
		variables = make(map[string]interface{})
	}

	// Check required variables
	for _, requiredVar := range prompt.RequiredVars {
		if _, ok := variables[requiredVar]; !ok {
			return "", fmt.Errorf("missing required variable: %s", requiredVar)
		}
	}

	// Execute template
	var buf bytes.Buffer
	if err := prompt.Template.Execute(&buf, variables); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}

// validateParameters validates and applies default parameters
func (s *LLMGatewayServer) validateParameters(params *pb.LLMParameters, prompt *Prompt) (*LLMParameters, error) {
	result := &LLMParameters{}

	if params == nil {
		return result, nil
	}

	// Validate temperature (0.0 - 2.0)
	if params.Temperature < 0 || params.Temperature > 2.0 {
		return nil, fmt.Errorf("temperature must be between 0.0 and 2.0")
	}
	result.Temperature = params.Temperature

	// Validate max_tokens (1 - 32000, model dependent)
	if params.MaxTokens < 0 || params.MaxTokens > 32000 {
		return nil, fmt.Errorf("max_tokens must be between 1 and 32000")
	}
	result.MaxTokens = params.MaxTokens

	// Validate top_p (0.0 - 1.0)
	if params.TopP < 0 || params.TopP > 1.0 {
		return nil, fmt.Errorf("top_p must be between 0.0 and 1.0")
	}
	result.TopP = params.TopP

	// Validate frequency_penalty (-2.0 - 2.0)
	if params.FrequencyPenalty < -2.0 || params.FrequencyPenalty > 2.0 {
		return nil, fmt.Errorf("frequency_penalty must be between -2.0 and 2.0")
	}
	result.FrequencyPenalty = params.FrequencyPenalty

	// Validate presence_penalty (-2.0 - 2.0)
	if params.PresencePenalty < -2.0 || params.PresencePenalty > 2.0 {
		return nil, fmt.Errorf("presence_penalty must be between -2.0 and 2.0")
	}
	result.PresencePenalty = params.PresencePenalty

	return result, nil
}

// trackUsageAsync tracks usage asynchronously
func (s *LLMGatewayServer) trackUsageAsync(event *UsageEvent) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.usageTracker.TrackUsage(ctx, event); err != nil {
			s.logger.Error("failed to track usage",
				zap.String("request_id", event.RequestID),
				zap.Error(err))
		}
	}()
}

// GetPromptMetadata returns metadata for a prompt
func (s *LLMGatewayServer) GetPromptMetadata(ctx context.Context, req *pb.GetPromptMetadataRequest) (*pb.GetPromptMetadataResponse, error) {
	if req.PromptPath == "" {
		return nil, status.Error(codes.InvalidArgument, "prompt_path is required")
	}

	prompt, err := s.promptLoader.GetPrompt(req.PromptPath)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("prompt not found: %s", req.PromptPath))
	}

	return &pb.GetPromptMetadataResponse{
		PromptPath:      prompt.Path,
		FileSizeBytes:   prompt.FileSizeBytes,
		LastModified:    prompt.LastModified.Format(time.RFC3339),
		RequiredVariables: prompt.RequiredVars,
	}, nil
}

// ListPrompts lists all available prompts
func (s *LLMGatewayServer) ListPrompts(ctx context.Context, req *pb.ListPromptsRequest) (*pb.ListPromptsResponse, error) {
	prompts := s.promptLoader.ListPrompts(req.DirectoryFilter)

	promptInfos := make([]*pb.PromptInfo, len(prompts))
	for i, prompt := range prompts {
		promptInfos[i] = &pb.PromptInfo{
			Path:         prompt.Path,
			SizeBytes:    prompt.FileSizeBytes,
			LastModified: prompt.LastModified.Format(time.RFC3339),
		}
	}

	return &pb.ListPromptsResponse{
		Prompts: promptInfos,
	}, nil
}

// GetUsageStats returns usage statistics
func (s *LLMGatewayServer) GetUsageStats(ctx context.Context, req *pb.GetUsageStatsRequest) (*pb.GetUsageStatsResponse, error) {
	stats, err := s.usageTracker.GetStats(req.TimeRange, req.CallingService)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get stats: %v", err))
	}

	return &pb.GetUsageStatsResponse{
		TotalRequests:     stats.TotalRequests,
		TotalTokens:       stats.TotalTokens,
		RequestsByService: stats.RequestsByService,
		TokensByModel:     stats.TokensByModel,
	}, nil
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
