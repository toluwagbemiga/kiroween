package internal

import (
	"context"
	"testing"
	"text/template"

	pb "github.com/haunted-saas/llm-gateway-service/proto/llm/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLLMGatewayServer_SubstituteVariables(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	server := &LLMGatewayServer{logger: logger}

	tests := []struct {
		name          string
		promptContent string
		variablesJSON string
		requiredVars  []string
		expectError   bool
		expectedText  string
	}{
		{
			name:          "simple substitution",
			promptContent: "Hello {{.name}}!",
			variablesJSON: `{"name": "Alice"}`,
			requiredVars:  []string{"name"},
			expectError:   false,
			expectedText:  "Hello Alice!",
		},
		{
			name:          "multiple variables",
			promptContent: "Hello {{.name}}, you are {{.age}} years old.",
			variablesJSON: `{"name": "Bob", "age": 30}`,
			requiredVars:  []string{"name", "age"},
			expectError:   false,
			expectedText:  "Hello Bob, you are 30 years old.",
		},
		{
			name:          "nested variables",
			promptContent: "User: {{.user.name}}, Email: {{.user.email}}",
			variablesJSON: `{"user": {"name": "Charlie", "email": "charlie@example.com"}}`,
			requiredVars:  []string{"user"},
			expectError:   false,
			expectedText:  "User: Charlie, Email: charlie@example.com",
		},
		{
			name:          "missing required variable",
			promptContent: "Hello {{.name}}!",
			variablesJSON: `{}`,
			requiredVars:  []string{"name"},
			expectError:   true,
		},
		{
			name:          "invalid JSON",
			promptContent: "Hello {{.name}}!",
			variablesJSON: `{invalid json}`,
			requiredVars:  []string{"name"},
			expectError:   true,
		},
		{
			name:          "empty variables",
			promptContent: "Static prompt",
			variablesJSON: `{}`,
			requiredVars:  []string{},
			expectError:   false,
			expectedText:  "Static prompt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create prompt
			prompt := &Prompt{
				Content:      tt.promptContent,
				RequiredVars: tt.requiredVars,
			}

			// Parse template
			tmpl, err := template.New("test").Parse(tt.promptContent)
			require.NoError(t, err)
			prompt.Template = tmpl

			// Substitute variables
			result, err := server.substituteVariables(prompt, tt.variablesJSON)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedText, result)
			}
		})
	}
}

func TestLLMGatewayServer_ValidateParameters(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	server := &LLMGatewayServer{logger: logger}

	tests := []struct {
		name        string
		params      *pb.LLMParameters
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid parameters",
			params: &pb.LLMParameters{
				Temperature:      0.7,
				MaxTokens:        1000,
				TopP:             0.9,
				FrequencyPenalty: 0.5,
				PresencePenalty:  0.5,
			},
			expectError: false,
		},
		{
			name: "temperature too high",
			params: &pb.LLMParameters{
				Temperature: 2.5,
			},
			expectError: true,
			errorMsg:    "temperature",
		},
		{
			name: "temperature negative",
			params: &pb.LLMParameters{
				Temperature: -0.1,
			},
			expectError: true,
			errorMsg:    "temperature",
		},
		{
			name: "max_tokens too high",
			params: &pb.LLMParameters{
				MaxTokens: 50000,
			},
			expectError: true,
			errorMsg:    "max_tokens",
		},
		{
			name: "top_p out of range",
			params: &pb.LLMParameters{
				TopP: 1.5,
			},
			expectError: true,
			errorMsg:    "top_p",
		},
		{
			name: "frequency_penalty out of range",
			params: &pb.LLMParameters{
				FrequencyPenalty: 3.0,
			},
			expectError: true,
			errorMsg:    "frequency_penalty",
		},
		{
			name:        "nil parameters",
			params:      nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := server.validateParameters(tt.params, &Prompt{})

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestLLMGatewayServer_CallPrompt_Validation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	// Create minimal server setup
	cache := NewPromptCache()
	promptLoader := &PromptLoader{
		cache:  cache,
		logger: logger,
	}
	
	router := NewLLMRouter("openai", logger)
	usageTracker := NewUsageTracker(1000, logger)
	
	server := NewLLMGatewayServer(promptLoader, router, usageTracker, logger)

	tests := []struct {
		name         string
		request      *pb.CallPromptRequest
		expectedCode codes.Code
	}{
		{
			name: "missing prompt_path",
			request: &pb.CallPromptRequest{
				PromptPath: "",
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "invalid prompt_path with directory traversal",
			request: &pb.CallPromptRequest{
				PromptPath: "../../../etc/passwd",
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "prompt not found",
			request: &pb.CallPromptRequest{
				PromptPath: "nonexistent.txt",
			},
			expectedCode: codes.NotFound,
		},
		{
			name: "invalid timeout",
			request: &pb.CallPromptRequest{
				PromptPath:     "test.txt",
				TimeoutSeconds: 200, // Too high
			},
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := server.CallPrompt(context.Background(), tt.request)
			
			assert.Error(t, err)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, st.Code())
		})
	}
}
