package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// LLMProvider defines the interface for LLM providers
type LLMProvider interface {
	Call(ctx context.Context, req *LLMRequest) (*LLMResponse, error)
	GetName() string
	ValidateModel(model string) error
}

// OpenAIProvider implements the LLMProvider interface for OpenAI
type OpenAIProvider struct {
	client          *openai.Client
	apiKey          string
	supportedModels map[string]bool
	logger          *zap.Logger
	testMode        bool
}

// Supported OpenAI models
var OpenAIModels = []string{
	"gpt-4-turbo-preview",
	"gpt-4-turbo",
	"gpt-4",
	"gpt-4-32k",
	"gpt-3.5-turbo",
	"gpt-3.5-turbo-16k",
	"gpt-3.5-turbo-1106",
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string, testMode bool, logger *zap.Logger) (*OpenAIProvider, error) {
	if apiKey == "" && !testMode {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	// Build supported models map
	supportedModels := make(map[string]bool)
	for _, model := range OpenAIModels {
		supportedModels[model] = true
	}

	var client *openai.Client
	if !testMode {
		client = openai.NewClient(apiKey)
	}

	return &OpenAIProvider{
		client:          client,
		apiKey:          apiKey,
		supportedModels: supportedModels,
		logger:          logger,
		testMode:        testMode,
	}, nil
}

// Call executes a request to OpenAI
func (p *OpenAIProvider) Call(ctx context.Context, req *LLMRequest) (*LLMResponse, error) {
	startTime := time.Now()

	// Test mode - return mock response
	if p.testMode {
		return p.mockResponse(req, startTime), nil
	}

	// Validate model
	if err := p.ValidateModel(req.Model); err != nil {
		return nil, err
	}

	// Apply timeout
	if req.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, req.Timeout)
		defer cancel()
	}

	// Build OpenAI request
	openaiReq := openai.ChatCompletionRequest{
		Model: req.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req.Prompt,
			},
		},
	}

	// Apply parameters
	if req.Parameters != nil {
		if req.Parameters.Temperature > 0 {
			openaiReq.Temperature = req.Parameters.Temperature
		}
		if req.Parameters.MaxTokens > 0 {
			openaiReq.MaxTokens = int(req.Parameters.MaxTokens)
		}
		if req.Parameters.TopP > 0 {
			openaiReq.TopP = req.Parameters.TopP
		}
		if req.Parameters.FrequencyPenalty != 0 {
			openaiReq.FrequencyPenalty = req.Parameters.FrequencyPenalty
		}
		if req.Parameters.PresencePenalty != 0 {
			openaiReq.PresencePenalty = req.Parameters.PresencePenalty
		}
	}

	p.logger.Debug("calling OpenAI API",
		zap.String("model", req.Model),
		zap.String("request_id", req.RequestID))

	// Call OpenAI API
	resp, err := p.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	// Check if we got a response
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	responseTime := time.Since(startTime)

	// Build response
	llmResp := &LLMResponse{
		Text:  resp.Choices[0].Message.Content,
		Model: resp.Model,
		TokenUsage: &TokenUsage{
			PromptTokens:     int32(resp.Usage.PromptTokens),
			CompletionTokens: int32(resp.Usage.CompletionTokens),
			TotalTokens:      int32(resp.Usage.TotalTokens),
		},
		ResponseTime: responseTime,
	}

	p.logger.Debug("OpenAI API call completed",
		zap.String("model", resp.Model),
		zap.Int("prompt_tokens", resp.Usage.PromptTokens),
		zap.Int("completion_tokens", resp.Usage.CompletionTokens),
		zap.Duration("response_time", responseTime))

	return llmResp, nil
}

// mockResponse returns a mock response for testing
func (p *OpenAIProvider) mockResponse(req *LLMRequest, startTime time.Time) *LLMResponse {
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Generate realistic token counts
	promptTokens := int32(len(req.Prompt) / 4) // Rough estimate: 1 token ≈ 4 chars
	completionTokens := int32(50)              // Mock completion length

	return &LLMResponse{
		Text:  "[TEST MODE] This is a mock response from the LLM Gateway Service. In production, this would be the actual LLM response.",
		Model: req.Model,
		TokenUsage: &TokenUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
		ResponseTime: time.Since(startTime),
	}
}

// GetName returns the provider name
func (p *OpenAIProvider) GetName() string {
	return "openai"
}

// ValidateModel validates if a model is supported
func (p *OpenAIProvider) ValidateModel(model string) error {
	if !p.supportedModels[model] {
		return fmt.Errorf("unsupported model: %s", model)
	}
	return nil
}

// LLMRouter routes requests to the appropriate LLM provider
type LLMRouter struct {
	providers       map[string]LLMProvider
	defaultProvider string
	defaultModels   map[string]string
	retryConfig     *RetryConfig
	logger          *zap.Logger
}

// NewLLMRouter creates a new LLM router
func NewLLMRouter(defaultProvider string, logger *zap.Logger) *LLMRouter {
	return &LLMRouter{
		providers:       make(map[string]LLMProvider),
		defaultProvider: defaultProvider,
		defaultModels: map[string]string{
			"openai": "gpt-4-turbo-preview",
		},
		retryConfig: DefaultRetryConfig(),
		logger:      logger,
	}
}

// RegisterProvider registers an LLM provider
func (r *LLMRouter) RegisterProvider(provider LLMProvider) {
	r.providers[provider.GetName()] = provider
	r.logger.Info("LLM provider registered", zap.String("provider", provider.GetName()))
}

// Route routes a request to the appropriate provider
func (r *LLMRouter) Route(ctx context.Context, req *LLMRequest) (*LLMResponse, error) {
	// Select provider
	providerName := r.defaultProvider
	if req.Model != "" {
		// Try to infer provider from model name
		if strings.HasPrefix(req.Model, "gpt-") {
			providerName = "openai"
		} else if strings.HasPrefix(req.Model, "claude-") {
			providerName = "anthropic"
		}
	}

	provider, ok := r.providers[providerName]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", providerName)
	}

	// Select model
	if req.Model == "" {
		req.Model = r.defaultModels[providerName]
	}

	// Call provider with retry logic
	return r.callWithRetry(ctx, provider, req)
}

// callWithRetry calls a provider with exponential backoff retry
func (r *LLMRouter) callWithRetry(ctx context.Context, provider LLMProvider, req *LLMRequest) (*LLMResponse, error) {
	var lastErr error
	delay := r.retryConfig.InitialDelay

	for attempt := 1; attempt <= r.retryConfig.MaxAttempts; attempt++ {
		resp, err := provider.Call(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Check if error is retryable (rate limit)
		if !isRateLimitError(err) {
			return nil, err
		}

		// Don't retry on last attempt
		if attempt == r.retryConfig.MaxAttempts {
			break
		}

		r.logger.Warn("rate limit hit, retrying",
			zap.String("provider", provider.GetName()),
			zap.Int("attempt", attempt),
			zap.Duration("delay", delay),
			zap.Error(err))

		// Wait before retry
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}

		// Exponential backoff
		delay = time.Duration(float64(delay) * r.retryConfig.BackoffFactor)
		if delay > r.retryConfig.MaxDelay {
			delay = r.retryConfig.MaxDelay
		}
	}

	return nil, fmt.Errorf("max retry attempts exceeded: %w", lastErr)
}

// isRateLimitError checks if an error is a rate limit error
func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "too many requests")
}
