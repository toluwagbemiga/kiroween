# Requirements Document: LLM Gateway Service

## Introduction

The LLM Gateway Service acts as the Intelligence Layer of the Haunted SaaS Skeleton platform, serving as the single, secure proxy for all external LLM API calls. It implements a "Prompt-as-Code" strategy by loading prompts from the filesystem and provides a unified interface for other services to execute LLM operations without managing API keys or prompt templates directly.

## Glossary

- **LLM Gateway Service**: The Go-based microservice that proxies all external Large Language Model API requests
- **Prompt-as-Code**: A strategy where prompts are stored as files in a version-controlled directory and loaded dynamically
- **LLM Provider**: An external service that provides language model APIs (e.g., OpenAI, Anthropic)
- **Prompt Template**: A text file containing a prompt with variable placeholders that can be substituted at runtime
- **Token Usage**: The number of tokens consumed by an LLM API request, used for billing and analytics
- **gRPC**: The internal communication protocol used between microservices
- **JSON Variables**: Key-value pairs provided by calling services to populate prompt template placeholders

## Requirements

### Requirement 1

**User Story:** As a developer, I want to store prompts as files in a /prompts directory, so that prompts are version-controlled and can be updated without code changes

#### Acceptance Criteria

1. THE LLM Gateway Service SHALL load all prompt files from the /prompts directory at startup
2. THE LLM Gateway Service SHALL support nested directory structures within /prompts to organize prompts by feature or service
3. WHEN a prompt file is added or modified, THE LLM Gateway Service SHALL reload prompts without requiring a service restart
4. THE LLM Gateway Service SHALL validate that each prompt file contains valid UTF-8 text
5. THE LLM Gateway Service SHALL log the number of prompts loaded and any files that failed to load

### Requirement 2

**User Story:** As a developer, I want to reference prompts by their file path, so that I can organize prompts logically and avoid naming conflicts

#### Acceptance Criteria

1. THE LLM Gateway Service SHALL accept prompt references using relative paths from the /prompts directory
2. THE LLM Gateway Service SHALL support file extensions including .txt, .md, and .prompt
3. WHEN a prompt path includes subdirectories, THE LLM Gateway Service SHALL resolve the full path correctly
4. THE LLM Gateway Service SHALL return an error when a referenced prompt file does not exist
5. THE LLM Gateway Service SHALL treat prompt paths as case-sensitive on all platforms

### Requirement 3

**User Story:** As a calling service, I want to execute a prompt with dynamic variables, so that I can customize prompts for specific contexts without creating multiple files

#### Acceptance Criteria

1. THE LLM Gateway Service SHALL provide a CallPrompt gRPC endpoint that accepts a prompt path and JSON variables
2. WHEN JSON variables are provided, THE LLM Gateway Service SHALL substitute placeholders in the prompt template using the format {{variable_name}}
3. THE LLM Gateway Service SHALL return an error when a required variable placeholder is not provided in the JSON variables
4. THE LLM Gateway Service SHALL preserve variable values exactly as provided without additional formatting
5. THE LLM Gateway Service SHALL support nested JSON objects in variables using dot notation (e.g., {{user.name}})

### Requirement 4

**User Story:** As a security administrator, I want LLM API keys to be managed exclusively by the Gateway Service, so that credentials are not distributed across multiple services

#### Acceptance Criteria

1. THE LLM Gateway Service SHALL load LLM provider API keys from environment variables at startup
2. THE LLM Gateway Service SHALL validate that required API keys are present before accepting requests
3. THE LLM Gateway Service SHALL never expose API keys in gRPC responses, logs, or error messages
4. THE LLM Gateway Service SHALL support multiple provider API keys for different LLM services
5. THE LLM Gateway Service SHALL terminate startup if required API keys are missing or invalid

### Requirement 5

**User Story:** As a calling service, I want to specify which LLM provider and model to use, so that I can choose the appropriate model for different use cases

#### Acceptance Criteria

1. THE LLM Gateway Service SHALL accept an optional provider parameter in the CallPrompt request
2. THE LLM Gateway Service SHALL accept an optional model parameter in the CallPrompt request
3. WHEN no provider is specified, THE LLM Gateway Service SHALL use a default provider configured at startup
4. WHEN no model is specified, THE LLM Gateway Service SHALL use a default model for the selected provider
5. THE LLM Gateway Service SHALL return an error when an unsupported provider or model is requested

### Requirement 6

**User Story:** As a calling service, I want to receive the LLM response as structured data, so that I can process the results programmatically

#### Acceptance Criteria

1. THE LLM Gateway Service SHALL return the LLM response text in the gRPC response
2. THE LLM Gateway Service SHALL include token usage information in the response including prompt tokens and completion tokens
3. THE LLM Gateway Service SHALL include the model used in the response
4. THE LLM Gateway Service SHALL include a unique request ID in the response for tracking and debugging
5. WHEN the LLM API returns an error, THE LLM Gateway Service SHALL return a gRPC error with the provider error message

### Requirement 7

**User Story:** As a system administrator, I want all LLM requests to be logged, so that I can audit usage and troubleshoot issues

#### Acceptance Criteria

1. WHEN a CallPrompt request is received, THE LLM Gateway Service SHALL log the prompt path, calling service, and request ID
2. WHEN an LLM API call completes, THE LLM Gateway Service SHALL log the response time, token usage, and status
3. WHEN an LLM API call fails, THE LLM Gateway Service SHALL log the error message, provider error code, and request details
4. THE LLM Gateway Service SHALL include correlation IDs in all log entries to trace requests across services
5. THE LLM Gateway Service SHALL never log the actual prompt content or LLM response to prevent sensitive data exposure

### Requirement 8

**User Story:** As a data analyst, I want token usage to be tracked for analytics, so that I can monitor costs and usage patterns

#### Acceptance Criteria

1. WHEN an LLM API call completes, THE LLM Gateway Service SHALL emit a usage event to the analytics service
2. THE LLM Gateway Service SHALL include in the usage event the prompt path, model, token counts, and calling service
3. THE LLM Gateway Service SHALL include the request timestamp and response time in the usage event
4. THE LLM Gateway Service SHALL emit usage events asynchronously to avoid blocking the response
5. WHEN the analytics service is unavailable, THE LLM Gateway Service SHALL log the usage data locally and continue processing

### Requirement 9

**User Story:** As a developer, I want to configure LLM parameters like temperature and max tokens, so that I can control response characteristics

#### Acceptance Criteria

1. THE LLM Gateway Service SHALL accept optional parameters in the CallPrompt request including temperature, max_tokens, and top_p
2. WHEN parameters are not provided, THE LLM Gateway Service SHALL use default values configured at startup
3. THE LLM Gateway Service SHALL validate that parameter values are within acceptable ranges for the selected model
4. THE LLM Gateway Service SHALL pass validated parameters to the LLM provider API
5. THE LLM Gateway Service SHALL return an error when invalid parameter values are provided

### Requirement 10

**User Story:** As a calling service, I want requests to timeout appropriately, so that slow LLM responses do not block my operations indefinitely

#### Acceptance Criteria

1. THE LLM Gateway Service SHALL enforce a configurable timeout for LLM API calls with a default of 30 seconds
2. WHEN an LLM API call exceeds the timeout, THE LLM Gateway Service SHALL cancel the request and return a timeout error
3. THE LLM Gateway Service SHALL log timeout events with the prompt path and elapsed time
4. THE LLM Gateway Service SHALL allow calling services to specify a custom timeout in the CallPrompt request
5. THE LLM Gateway Service SHALL validate that custom timeout values are between 5 and 120 seconds

### Requirement 11

**User Story:** As a developer, I want to test prompts without consuming API credits, so that I can develop and debug prompt templates efficiently

#### Acceptance Criteria

1. WHERE a test mode is enabled, THE LLM Gateway Service SHALL return mock responses without calling external LLM APIs
2. WHERE test mode is enabled, THE LLM Gateway Service SHALL return responses with realistic token counts
3. THE LLM Gateway Service SHALL load test mode configuration from an environment variable
4. WHERE test mode is enabled, THE LLM Gateway Service SHALL log that mock responses are being returned
5. THE LLM Gateway Service SHALL support test mode only in non-production environments

### Requirement 12

**User Story:** As a system administrator, I want the service to handle rate limits gracefully, so that temporary API limits do not cause service failures

#### Acceptance Criteria

1. WHEN an LLM provider returns a rate limit error, THE LLM Gateway Service SHALL implement exponential backoff retry logic
2. THE LLM Gateway Service SHALL retry rate-limited requests up to 3 times with increasing delays
3. WHEN retries are exhausted, THE LLM Gateway Service SHALL return a rate limit error to the calling service
4. THE LLM Gateway Service SHALL log rate limit events including the provider and retry attempts
5. THE LLM Gateway Service SHALL include the retry-after duration in the error response when provided by the LLM provider
