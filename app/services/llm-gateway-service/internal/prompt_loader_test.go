package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPromptLoader_LoadAllPrompts(t *testing.T) {
	// Create temporary prompts directory
	tmpDir := t.TempDir()
	logger, _ := zap.NewDevelopment()

	// Create test prompt files
	testPrompts := map[string]string{
		"simple.txt":           "Hello {{name}}!",
		"nested/complex.md":    "User: {{user.name}}\nEmail: {{user.email}}",
		"with-frontmatter.md": "---\ndescription: Test prompt\nrequired_vars: [name, age]\n---\nHello {{name}}, you are {{age}} years old.",
	}

	for path, content := range testPrompts {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create prompt loader
	loader, err := NewPromptLoader(tmpDir, false, logger)
	require.NoError(t, err)

	// Load all prompts
	err = loader.LoadAllPrompts()
	require.NoError(t, err)

	// Verify prompts were loaded
	assert.Equal(t, 3, loader.cache.Count())

	// Test simple prompt
	prompt, err := loader.GetPrompt("simple.txt")
	require.NoError(t, err)
	assert.Equal(t, "simple.txt", prompt.Path)
	assert.Contains(t, prompt.RequiredVars, "name")

	// Test nested prompt
	prompt, err = loader.GetPrompt("nested/complex.md")
	require.NoError(t, err)
	assert.Equal(t, "nested/complex.md", prompt.Path)
	assert.Contains(t, prompt.RequiredVars, "user")

	// Test prompt with frontmatter
	prompt, err = loader.GetPrompt("with-frontmatter.md")
	require.NoError(t, err)
	assert.NotNil(t, prompt.Metadata)
	assert.Equal(t, "Test prompt", prompt.Metadata.Description)
	assert.Contains(t, prompt.RequiredVars, "name")
	assert.Contains(t, prompt.RequiredVars, "age")
}

func TestPromptLoader_ExtractRequiredVariables(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	loader := &PromptLoader{logger: logger}

	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "simple variable",
			content:  "Hello {{name}}!",
			expected: []string{"name"},
		},
		{
			name:     "multiple variables",
			content:  "Hello {{name}}, you are {{age}} years old.",
			expected: []string{"name", "age"},
		},
		{
			name:     "nested variable",
			content:  "User: {{user.name}}, Email: {{user.email}}",
			expected: []string{"user"},
		},
		{
			name:     "duplicate variables",
			content:  "{{name}} {{name}} {{name}}",
			expected: []string{"name"},
		},
		{
			name:     "no variables",
			content:  "This is a static prompt.",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := loader.extractRequiredVariables(tt.content)
			
			// Check that all expected variables are present
			for _, expectedVar := range tt.expected {
				assert.Contains(t, vars, expectedVar)
			}
			
			// Check count matches (accounting for uniqueness)
			assert.Equal(t, len(tt.expected), len(vars))
		})
	}
}

func TestPromptLoader_ParseFrontmatter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	loader := &PromptLoader{logger: logger}

	tests := []struct {
		name            string
		content         string
		expectMetadata  bool
		expectedContent string
	}{
		{
			name: "with frontmatter",
			content: `---
description: Test prompt
required_vars: [name, age]
default_model: gpt-4
---
Hello {{name}}!`,
			expectMetadata:  true,
			expectedContent: "Hello {{name}}!",
		},
		{
			name:            "without frontmatter",
			content:         "Hello {{name}}!",
			expectMetadata:  false,
			expectedContent: "Hello {{name}}!",
		},
		{
			name: "invalid frontmatter",
			content: `---
invalid yaml: [
---
Hello {{name}}!`,
			expectMetadata:  false,
			expectedContent: `---
invalid yaml: [
---
Hello {{name}}!`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, content := loader.parseFrontmatter([]byte(tt.content))
			
			if tt.expectMetadata {
				assert.NotNil(t, metadata)
			} else {
				assert.Nil(t, metadata)
			}
			
			assert.Equal(t, tt.expectedContent, content)
		})
	}
}

func TestPromptLoader_GetPrompt_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	logger, _ := zap.NewDevelopment()

	loader, err := NewPromptLoader(tmpDir, false, logger)
	require.NoError(t, err)

	_, err = loader.GetPrompt("nonexistent.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestPromptLoader_ListPrompts(t *testing.T) {
	tmpDir := t.TempDir()
	logger, _ := zap.NewDevelopment()

	// Create test prompts
	testPrompts := map[string]string{
		"root.txt":           "Root prompt",
		"feature1/test.txt":  "Feature 1 prompt",
		"feature2/test.txt":  "Feature 2 prompt",
	}

	for path, content := range testPrompts {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	loader, err := NewPromptLoader(tmpDir, false, logger)
	require.NoError(t, err)
	err = loader.LoadAllPrompts()
	require.NoError(t, err)

	// Test list all
	prompts := loader.ListPrompts("")
	assert.Equal(t, 3, len(prompts))

	// Test list with filter
	prompts = loader.ListPrompts("feature1")
	assert.Equal(t, 1, len(prompts))
	assert.Equal(t, "feature1/test.txt", prompts[0].Path)
}
