# Design Document: LLM Gateway Service

## Overview

The LLM Gateway Service acts as the Intelligence Layer, providing a secure, centralized proxy for all external LLM API calls. It implements a "Prompt-as-Code" strategy where prompts are stored as files in a `/prompts` directory and loaded dynamically. The service manages API keys, tracks token usage, and provides a unified gRPC interface for other services to execute LLM operations without direct access to provider APIs.

## Architecture

### High-Level Architecture

```mermaid
graph TB
    subgraph "External LLM Providers"
        OpenAI[OpenAI API]
        Anthropic[Anthropic API]
        OtherLLM[Other LLM APIs]
    end
    
    subgraph "LLM Gateway Service"
        gRPCServer[gRPC Server]
        PromptLoader[Prompt Loader]
        PromptCache[Prompt Cache]
        LLMRouter[LLM Provider Router]
        OpenAIClient[OpenAI Client]
        AnthropicClient[Anthropic Client]
        UsageTracker[Usage Tracker]
        FileWatcher[File Watcher]
    end
    
    subgraph "Internal Services"
        CallingServices[Other Services]
    end
    
    subgraph "File System"
        PromptsDir[/prompts Directory]
    end
    
    CallingServices -->|gRPC: CallPrompt| gRPCServer
    gRPCServer --> LLMRouter
    LLMRouter --> PromptCache
    PromptCache --> PromptLoader
    PromptLoader --> PromptsDir
    FileWatcher --> PromptsDir
    FileWatcher -->|Reload| PromptCache
    LLMRouter --> OpenAIClient
    LLMRouter --> AnthropicClient
    OpenAIClient -->|HTTPS| OpenAI
    AnthropicClient -->|HTTPS| Anthropic
    LLMRouter --> UsageTracker
    UsageTracker -->|gRPC| AnalyticsService[Analytics Service]
```

### Technology Stack

- **Language**: Go 1.21+
- **gRPC Framework**: google.golang.org/grpc
- **OpenAI SDK**: github.com/sashabaranov/go-openai
- **File Watching**: github.com/fsnotify/fsnotify
- **Template Engine**: text/template (standard library)
- **Configuration**: Environment variables via godotenv
- **Logging**: Structured logging with zerolog

## Components and Interfaces

### 1. gRPC Service Definition

```protobuf
syntax = "proto3";

package llm.v1;

service LLMGatewayService {
  rpc CallPrompt(CallPromptRequest) returns (CallPromptResponse);
  rpc GetPromptMetadata(GetPromptMetadataRequest) returns (GetPromptMetadataResponse);
  rpc ListPrompts(ListPromptsRequest) returns (ListPromptsResponse);
  rpc GetUsageStats(GetUsageStatsRequest) returns (GetUsageStatsResponse);
}

message CallPromptRequest {
  string prompt_path = 1;
  string variables_json = 2;
  string provider = 3; // Optional: "openai", "anthropic"
  string model = 4; // Optional: "gpt-4", "claude-3-opus"
  LLMParameters parameters = 5; // Optional
  int32 timeout_seconds = 6; // Optional
  string calling_service = 7;
  string correlation_id = 8;
}

message LLMParameters {
  float temperature = 1;
  int32 max_tokens = 2;
  float top_p = 3;
  float frequency_penalty = 4;
  float presence_penalty = 5;
}

message CallPromptResponse {
  string response_text = 1;
  TokenUsage token_usage = 2;
  string model_used = 3;
  string request_id = 4;
  int64 response_time_ms = 5;
}

message TokenUsage {
  int32 prompt_tokens = 1;
  int32 completion_tokens = 2;
  int32 total_tokens = 3;
}

message GetPromptMetadataRequest {
  string prompt_path = 1;
}

message GetPromptMetadataResponse {
  string prompt_path = 1;
  int64 file_size_bytes = 2;
  string last_modified = 3;
  repeated string required_variables = 4;
}

message ListPromptsRequest {
  string directory_filter = 1; // Optional: filter by subdirectory
}

message ListPromptsResponse {
  repeated PromptInfo prompts = 1;
}

message PromptInfo {
  string path = 1;
  int64 size_bytes = 2;
  string last_modified = 3;
}

message GetUsageStatsRequest {
  string time_range = 1; // "hour", "day", "week"
  string calling_service = 2; // Optional filter
}

message GetUsageStatsResponse {
  int64 total_requests = 1;
  int64 total_tokens = 2;
  map<string, int64> requests_by_service = 3;
  map<string, int64> tokens_by_model = 4;
}
```

### 2. Prompt Loader and Cache

```go
type PromptLoader struct {
    promptsDir string
    cache      *PromptCache
    watcher    *fsnotify.Watcher
    logger     *zerolog.Logger
}

type PromptCache struct {
    mu      sync.RWMutex
    prompts map[string]*Prompt
}

type Prompt struct {
    Path             string
    Content          string
    Template         *template.Template
    RequiredVars     []string
    LastModified     time.Time
    FileSizeBytes    int64
}

func NewPromptLoader(promptsDir string) (*PromptLoader, error)
func (l *PromptLoader) LoadAllPrompts() error
func (l *PromptLoader) GetPrompt(path string) (*Prompt, error)
func (l *PromptLoader) WatchForChanges() error
func (l *PromptLoader) ReloadPrompt(path string) error

// Template variable extraction
func (l *PromptLoader) extractRequiredVariables(content string) []string
```

### 3. LLM Provider Router

```go
type LLMRouter struct {
    providers      map[string]LLMProvider
    defaultProvider string
    defaultModels   map[string]string
    usageTracker   *UsageTracker
    logger         *zerolog.Logger
}

type LLMProvider interface {
    Call(ctx context.Context, req *LLMRequest) (*LLMResponse, error)
    GetName() string
    ValidateModel(model string) error
}

type LLMRequest struct {
    Prompt      string
    Model       string
    Parameters  *LLMParameters
    Timeout     time.Duration
    RequestID   string
}

type LLMResponse struct {
    Text         string
    TokenUsage   *TokenUsage
    Model        string
    ResponseTime time.Duration
}

func (r *LLMRouter) Route(ctx context.Context, req *CallPromptRequest, prompt *Prompt) (*LLMResponse, error)
func (r *LLMRouter) selectProvider(providerName string) (LLMProvider, error)
func (r *LLMRouter) selectModel(provider LLMProvider, modelName string) (string, error)
```

### 4. OpenAI Provider Implementation

```go
type OpenAIProvider struct {
    client       *openai.Client
    apiKey       string
    supportedModels []string
    logger       *zerolog.Logger
}

func NewOpenAIProvider(apiKey string) (*OpenAIProvider, error)

func (p *OpenAIProvider) Call(ctx context.Context, req *LLMRequest) (*LLMResponse, error) {
    // 1. Build OpenAI request
    // 2. Set timeout context
    // 3. Call OpenAI API
    // 4. Parse response
    // 5. Extract token usage
    // 6. Return structured response
}

func (p *OpenAIProvider) ValidateModel(model string) error
func (p *OpenAIProvider) GetName() string

// Supported models
var OpenAIModels = []string{
    "gpt-4-turbo-preview",
    "gpt-4",
    "gpt-3.5-turbo",
    "gpt-3.5-turbo-16k",
}
```

### 5. Usage Tracker

```go
type UsageTracker struct {
    analyticsClient analytics.AnalyticsServiceClient
    localStore      *UsageStore
    logger          *zerolog.Logger
}

type UsageEvent struct {
    RequestID       string
    PromptPath      string
    CallingService  string
    Provider        string
    Model           string
    PromptTokens    int32
    CompletionTokens int32
    TotalTokens     int32
    ResponseTimeMs  int64
    Timestamp       time.Time
    Success         bool
    ErrorMessage    string
}

func (t *UsageTracker) TrackUsage(ctx context.Context, event *UsageEvent) error
func (t *UsageTracker) GetStats(timeRange string, serviceFilter string) (*UsageStats, error)
func (t *UsageTracker) sendToAnalytics(event *UsageEvent) error
```

### 6. Template Variable Substitution

```go
type VariableSubstitutor struct {
    logger *zerolog.Logger
}

func (s *VariableSubstitutor) Substitute(prompt *Prompt, variablesJSON string) (string, error) {
    // 1. Parse JSON variables
    // 2. Validate required variables are present
    // 3. Execute template with variables
    // 4. Return rendered prompt
}

func (s *VariableSubstitutor) parseVariables(jsonStr string) (map[string]interface{}, error)
func (s *VariableSubstitutor) validateRequiredVars(prompt *Prompt, vars map[string]interface{}) error
```

## Data Models

### Prompt File Format

Prompts are stored as plain text files with optional frontmatter for metadata:

```markdown
---
description: Generate a welcome email for new users
required_vars: [user_name, user_email, team_name]
default_model: gpt-4
---

You are a friendly customer success manager writing a welcome email.

Write a personalized welcome email for {{user_name}} ({{user_email}}) who just joined the {{team_name}} team.

The email should:
- Welcome them warmly
- Explain the next steps
- Provide helpful resources
- Be professional but friendly

Email:
```

### Directory Structure

```
/prompts/
├── onboarding/
│   ├── welcome-email.md
│   ├── setup-guide.txt
│   └── feature-tour.md
├── notifications/
│   ├── alert-message.txt
│   └── reminder-message.txt
├── analytics/
│   ├── insight-summary.md
│   └── report-generation.md
└── support/
    ├── ticket-response.md
    └── faq-answer.txt
```

### In-Memory Usage Store

```go
type UsageStore struct {
    mu     sync.RWMutex
    events []UsageEvent
    maxSize int
}

func (s *UsageStore) Add(event UsageEvent)
func (s *UsageStore) GetRecent(limit int) []UsageEvent
func (s *UsageStore) GetStats(since time.Time) *UsageStats
func (s *UsageStore) Prune(olderThan time.Time)
```

## Error Handling

### Error Types

```go
type LLMGatewayError struct {
    Code    string
    Message string
    Cause   error
}

const (
    ErrCodePromptNotFound      = "PROMPT_NOT_FOUND"
    ErrCodeInvalidVariables    = "INVALID_VARIABLES"
    ErrCodeMissingVariable     = "MISSING_VARIABLE"
    ErrCodeProviderNotFound    = "PROVIDER_NOT_FOUND"
    ErrCodeModelNotSupported   = "MODEL_NOT_SUPPORTED"
    ErrCodeProviderAPIError    = "PROVIDER_API_ERROR"
    ErrCodeTimeout             = "TIMEOUT"
    ErrCodeRateLimited         = "RATE_LIMITED"
    ErrCodeInvalidParameters   = "INVALID_PARAMETERS"
)
```

### Error Handling Strategy

1. **Prompt Loading Errors**
   - Log file read errors but continue loading other prompts
   - Return specific error when requested prompt not found
   - Validate template syntax on load

2. **Variable Substitution Errors**
   - Return clear error messages indicating missing variables
   - Validate JSON format before parsing
   - Provide helpful error context

3. **Provider API Errors**
   - Implement exponential backoff for rate limits
   - Wrap provider errors with context
   - Log full error details for debugging
   - Return user-friendly error messages

4. **Timeout Handling**
   - Enforce configurable timeouts per request
   - Cancel provider API calls on timeout
   - Log timeout events with request details

### Retry Logic

```go
type RetryConfig struct {
    MaxAttempts     int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    BackoffFactor   float64
}

func (r *LLMRouter) callWithRetry(ctx context.Context, provider LLMProvider, req *LLMRequest) (*LLMResponse, error) {
    // Retry only on rate limit errors
    // Use exponential backoff
    // Respect context cancellation
}
```

## Testing Strategy

### Unit Tests

1. **Prompt Loader Tests**
   - Test loading prompts from directory
   - Test template variable extraction
   - Test file watching and reload
   - Test error handling for invalid files

2. **Variable Substitution Tests**
   - Test simple variable replacement
   - Test nested object access
   - Test missing variable detection
   - Test invalid JSON handling

3. **Provider Tests**
   - Mock provider API responses
   - Test request construction
   - Test response parsing
   - Test error handling

4. **Router Tests**
   - Test provider selection
   - Test model selection
   - Test default fallbacks
   - Test retry logic

### Integration Tests

1. **End-to-End Prompt Execution**
   - Test full CallPrompt flow
   - Test with real prompt files
   - Test variable substitution
   - Mock LLM provider responses

2. **File Watching Tests**
   - Test prompt reload on file change
   - Test handling of file deletion
   - Test handling of new files

3. **Usage Tracking Tests**
   - Test event emission
   - Test stats aggregation
   - Mock analytics service

### Test Data

```go
// Test prompts
const (
    TestPromptSimple = "Hello {{name}}!"
    TestPromptComplex = `
You are a {{role}}.
User: {{user.name}} ({{user.email}})
Task: {{task}}
`
)

// Test variables
var TestVariables = map[string]interface{}{
    "name": "Alice",
    "role": "assistant",
    "user": map[string]interface{}{
        "name":  "Bob",
        "email": "bob@example.com",
    },
    "task": "Write a summary",
}
```

## Configuration

### Environment Variables

```bash
# Server Configuration
GRPC_PORT=50052

# Prompts
PROMPTS_DIR=/app/prompts
WATCH_PROMPTS=true

# LLM Providers
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
DEFAULT_PROVIDER=openai
DEFAULT_MODEL=gpt-4-turbo-preview

# Timeouts
DEFAULT_TIMEOUT_SECONDS=30
MAX_TIMEOUT_SECONDS=120

# Usage Tracking
ANALYTICS_SERVICE_ADDR=analytics-service:50051
USAGE_STORE_MAX_SIZE=10000

# Test Mode
TEST_MODE=false
MOCK_RESPONSES=false

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Retry Configuration
MAX_RETRY_ATTEMPTS=3
INITIAL_RETRY_DELAY_MS=1000
MAX_RETRY_DELAY_MS=10000
```

### Startup Sequence

1. Load and validate environment variables
2. Validate LLM provider API keys
3. Initialize prompt loader and load all prompts
4. Start file watcher for prompt directory
5. Initialize LLM provider clients
6. Initialize usage tracker
7. Start gRPC server
8. Register health check endpoint
9. Log startup completion with prompt count

## Security Considerations

1. **API Key Management**
   - Load provider keys from environment only
   - Never log or expose keys in responses
   - Validate keys at startup
   - Use separate keys for test/production

2. **Prompt Access Control**
   - Validate prompt paths to prevent directory traversal
   - Restrict access to prompts directory
   - Log all prompt access attempts

3. **Variable Injection**
   - Sanitize user-provided variables
   - Prevent template injection attacks
   - Validate variable types and formats

4. **Rate Limiting**
   - Implement per-service rate limits
   - Track usage by calling service
   - Prevent abuse of expensive models

5. **Logging Security**
   - Never log full prompts (may contain sensitive data)
   - Never log LLM responses (may contain PII)
   - Log only metadata (paths, token counts, errors)

## Performance Considerations

1. **Prompt Caching**
   - Cache loaded prompts in memory
   - Reload only on file changes
   - Use read-write locks for concurrent access

2. **Concurrent Request Handling**
   - Process multiple LLM requests in parallel
   - Use goroutines for async operations
   - Implement connection pooling for provider APIs

3. **Usage Tracking**
   - Emit usage events asynchronously
   - Buffer events before sending to analytics
   - Prune old events from local store

4. **Template Rendering**
   - Pre-compile templates on load
   - Cache compiled templates
   - Reuse template instances

## Monitoring and Observability

### Metrics to Track

```go
// Request metrics
- llm_requests_total (counter, by provider, model, status)
- llm_request_duration_seconds (histogram, by provider, model)
- llm_tokens_total (counter, by provider, model, type)

// Prompt metrics
- prompts_loaded_total (gauge)
- prompt_reload_total (counter, by status)

// Error metrics
- llm_errors_total (counter, by error_code, provider)
- llm_timeouts_total (counter, by provider)
- llm_rate_limits_total (counter, by provider)

// Usage metrics
- llm_cost_estimate_total (counter, by provider, model)
```

### Health Checks

```go
type HealthChecker struct {
    promptLoader *PromptLoader
    providers    map[string]LLMProvider
}

func (h *HealthChecker) Check() *HealthStatus {
    // Check prompt loader status
    // Check provider connectivity (optional)
    // Check usage tracker status
    // Return overall health
}
```
