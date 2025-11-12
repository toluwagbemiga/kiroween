package internal

import (
	"sync"
	"text/template"
	"time"
)

// Prompt represents a loaded prompt template
type Prompt struct {
	Path          string
	Content       string
	Template      *template.Template
	RequiredVars  []string
	LastModified  time.Time
	FileSizeBytes int64
	Metadata      *PromptMetadata
}

// PromptMetadata contains optional frontmatter metadata
type PromptMetadata struct {
	Description  string   `yaml:"description"`
	RequiredVars []string `yaml:"required_vars"`
	DefaultModel string   `yaml:"default_model"`
	Temperature  *float32 `yaml:"temperature"`
	MaxTokens    *int32   `yaml:"max_tokens"`
}

// PromptCache is a thread-safe cache for loaded prompts
type PromptCache struct {
	mu      sync.RWMutex
	prompts map[string]*Prompt
}

// NewPromptCache creates a new prompt cache
func NewPromptCache() *PromptCache {
	return &PromptCache{
		prompts: make(map[string]*Prompt),
	}
}

// Get retrieves a prompt from the cache
func (c *PromptCache) Get(path string) (*Prompt, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	prompt, ok := c.prompts[path]
	return prompt, ok
}

// Set stores a prompt in the cache
func (c *PromptCache) Set(path string, prompt *Prompt) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prompts[path] = prompt
}

// Delete removes a prompt from the cache
func (c *PromptCache) Delete(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.prompts, path)
}

// List returns all prompt paths
func (c *PromptCache) List() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	paths := make([]string, 0, len(c.prompts))
	for path := range c.prompts {
		paths = append(paths, path)
	}
	return paths
}

// GetAll returns all prompts
func (c *PromptCache) GetAll() map[string]*Prompt {
	c.mu.RLock()
	defer c.mu.RUnlock()
	prompts := make(map[string]*Prompt, len(c.prompts))
	for path, prompt := range c.prompts {
		prompts[path] = prompt
	}
	return prompts
}

// Count returns the number of cached prompts
func (c *PromptCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.prompts)
}

// LLMRequest represents a request to an LLM provider
type LLMRequest struct {
	Prompt     string
	Model      string
	Parameters *LLMParameters
	Timeout    time.Duration
	RequestID  string
}

// LLMParameters contains LLM generation parameters
type LLMParameters struct {
	Temperature      float32
	MaxTokens        int32
	TopP             float32
	FrequencyPenalty float32
	PresencePenalty  float32
}

// LLMResponse represents a response from an LLM provider
type LLMResponse struct {
	Text         string
	TokenUsage   *TokenUsage
	Model        string
	ResponseTime time.Duration
}

// TokenUsage contains token usage information
type TokenUsage struct {
	PromptTokens     int32
	CompletionTokens int32
	TotalTokens      int32
}

// UsageEvent represents a usage tracking event
type UsageEvent struct {
	RequestID        string
	PromptPath       string
	CallingService   string
	Provider         string
	Model            string
	PromptTokens     int32
	CompletionTokens int32
	TotalTokens      int32
	ResponseTimeMs   int64
	Timestamp        time.Time
	Success          bool
	ErrorMessage     string
}

// UsageStats contains aggregated usage statistics
type UsageStats struct {
	TotalRequests      int64
	TotalTokens        int64
	RequestsByService  map[string]int64
	TokensByModel      map[string]int64
	AverageResponseMs  int64
}

// RetryConfig contains retry configuration
type RetryConfig struct {
	MaxAttempts   int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  1 * time.Second,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
	}
}
