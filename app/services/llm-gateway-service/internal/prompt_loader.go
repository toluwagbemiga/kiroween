package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// PromptLoader loads and manages prompt templates
type PromptLoader struct {
	promptsDir string
	cache      *PromptCache
	watcher    *fsnotify.Watcher
	logger     *zap.Logger
	watchMode  bool
}

// NewPromptLoader creates a new prompt loader
func NewPromptLoader(promptsDir string, watchMode bool, logger *zap.Logger) (*PromptLoader, error) {
	if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("prompts directory does not exist: %s", promptsDir)
	}

	loader := &PromptLoader{
		promptsDir: promptsDir,
		cache:      NewPromptCache(),
		logger:     logger,
		watchMode:  watchMode,
	}

	return loader, nil
}

// LoadAllPrompts loads all prompts from the prompts directory
func (l *PromptLoader) LoadAllPrompts() error {
	l.logger.Info("loading prompts", zap.String("directory", l.promptsDir))
	
	loadedCount := 0
	failedCount := 0

	err := filepath.Walk(l.promptsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			l.logger.Warn("error accessing path", zap.String("path", path), zap.Error(err))
			return nil // Continue walking
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file has valid extension
		if !l.isValidPromptFile(path) {
			return nil
		}

		// Get relative path from prompts directory
		relPath, err := filepath.Rel(l.promptsDir, path)
		if err != nil {
			l.logger.Error("failed to get relative path", zap.String("path", path), zap.Error(err))
			failedCount++
			return nil
		}

		// Normalize path separators to forward slashes
		relPath = filepath.ToSlash(relPath)

		// Load the prompt
		if err := l.loadPrompt(relPath, path); err != nil {
			l.logger.Error("failed to load prompt",
				zap.String("path", relPath),
				zap.Error(err))
			failedCount++
		} else {
			loadedCount++
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk prompts directory: %w", err)
	}

	l.logger.Info("prompts loaded",
		zap.Int("loaded", loadedCount),
		zap.Int("failed", failedCount),
		zap.Int("total", l.cache.Count()))

	return nil
}

// isValidPromptFile checks if a file has a valid prompt extension
func (l *PromptLoader) isValidPromptFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".txt" || ext == ".md" || ext == ".prompt"
}

// loadPrompt loads a single prompt file
func (l *PromptLoader) loadPrompt(relPath, absPath string) error {
	// Read file content
	content, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Get file info
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Parse frontmatter and content
	metadata, promptContent := l.parseFrontmatter(content)

	// Extract required variables from template
	requiredVars := l.extractRequiredVariables(promptContent)

	// Merge metadata required vars with extracted vars
	if metadata != nil && len(metadata.RequiredVars) > 0 {
		requiredVars = l.mergeUniqueVars(requiredVars, metadata.RequiredVars)
	}

	// Parse as Go template
	tmpl, err := template.New(relPath).Parse(promptContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create prompt object
	prompt := &Prompt{
		Path:          relPath,
		Content:       promptContent,
		Template:      tmpl,
		RequiredVars:  requiredVars,
		LastModified:  info.ModTime(),
		FileSizeBytes: info.Size(),
		Metadata:      metadata,
	}

	// Store in cache
	l.cache.Set(relPath, prompt)

	l.logger.Debug("prompt loaded",
		zap.String("path", relPath),
		zap.Int("required_vars", len(requiredVars)),
		zap.Int64("size_bytes", info.Size()))

	return nil
}

// parseFrontmatter parses YAML frontmatter from prompt content
func (l *PromptLoader) parseFrontmatter(content []byte) (*PromptMetadata, string) {
	// Check if content starts with ---
	if !bytes.HasPrefix(content, []byte("---\n")) && !bytes.HasPrefix(content, []byte("---\r\n")) {
		return nil, string(content)
	}

	// Find the closing ---
	lines := bytes.Split(content, []byte("\n"))
	endIdx := -1
	for i := 1; i < len(lines); i++ {
		line := bytes.TrimSpace(lines[i])
		if bytes.Equal(line, []byte("---")) {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		// No closing ---, treat as regular content
		return nil, string(content)
	}

	// Extract frontmatter
	frontmatterBytes := bytes.Join(lines[1:endIdx], []byte("\n"))
	
	// Parse YAML
	var metadata PromptMetadata
	if err := yaml.Unmarshal(frontmatterBytes, &metadata); err != nil {
		l.logger.Warn("failed to parse frontmatter, ignoring", zap.Error(err))
		return nil, string(content)
	}

	// Extract content after frontmatter
	contentBytes := bytes.Join(lines[endIdx+1:], []byte("\n"))
	promptContent := strings.TrimSpace(string(contentBytes))

	return &metadata, promptContent
}

// extractRequiredVariables extracts variable placeholders from template content
func (l *PromptLoader) extractRequiredVariables(content string) []string {
	// Match {{variable_name}} or {{object.property}}
	re := regexp.MustCompile(`\{\{\.?([a-zA-Z0-9_\.]+)\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	vars := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			varName := match[1]
			// Extract root variable name (before first dot)
			if idx := strings.Index(varName, "."); idx != -1 {
				varName = varName[:idx]
			}
			vars[varName] = true
		}
	}

	// Convert to slice
	result := make([]string, 0, len(vars))
	for varName := range vars {
		result = append(result, varName)
	}

	return result
}

// mergeUniqueVars merges two variable lists, removing duplicates
func (l *PromptLoader) mergeUniqueVars(vars1, vars2 []string) []string {
	varMap := make(map[string]bool)
	for _, v := range vars1 {
		varMap[v] = true
	}
	for _, v := range vars2 {
		varMap[v] = true
	}

	result := make([]string, 0, len(varMap))
	for v := range varMap {
		result = append(result, v)
	}
	return result
}

// GetPrompt retrieves a prompt from the cache
func (l *PromptLoader) GetPrompt(path string) (*Prompt, error) {
	prompt, ok := l.cache.Get(path)
	if !ok {
		return nil, fmt.Errorf("prompt not found: %s", path)
	}
	return prompt, nil
}

// ListPrompts returns all loaded prompts
func (l *PromptLoader) ListPrompts(directoryFilter string) []*Prompt {
	allPrompts := l.cache.GetAll()
	
	if directoryFilter == "" {
		// Return all prompts
		result := make([]*Prompt, 0, len(allPrompts))
		for _, prompt := range allPrompts {
			result = append(result, prompt)
		}
		return result
	}

	// Filter by directory
	result := make([]*Prompt, 0)
	for _, prompt := range allPrompts {
		if strings.HasPrefix(prompt.Path, directoryFilter) {
			result = append(result, prompt)
		}
	}
	return result
}

// WatchForChanges starts watching the prompts directory for changes
func (l *PromptLoader) WatchForChanges() error {
	if !l.watchMode {
		l.logger.Info("prompt watching disabled")
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	l.watcher = watcher

	// Add prompts directory and all subdirectories
	if err := l.addWatchRecursive(l.promptsDir); err != nil {
		return fmt.Errorf("failed to add watch paths: %w", err)
	}

	// Start watching in goroutine
	go l.watchLoop()

	l.logger.Info("prompt file watching started", zap.String("directory", l.promptsDir))
	return nil
}

// addWatchRecursive adds a directory and all subdirectories to the watcher
func (l *PromptLoader) addWatchRecursive(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := l.watcher.Add(path); err != nil {
				l.logger.Warn("failed to watch directory", zap.String("path", path), zap.Error(err))
			}
		}
		return nil
	})
}

// watchLoop processes file system events
func (l *PromptLoader) watchLoop() {
	for {
		select {
		case event, ok := <-l.watcher.Events:
			if !ok {
				return
			}
			l.handleFileEvent(event)

		case err, ok := <-l.watcher.Errors:
			if !ok {
				return
			}
			l.logger.Error("file watcher error", zap.Error(err))
		}
	}
}

// handleFileEvent handles a file system event
func (l *PromptLoader) handleFileEvent(event fsnotify.Event) {
	// Check if it's a prompt file
	if !l.isValidPromptFile(event.Name) {
		return
	}

	// Get relative path
	relPath, err := filepath.Rel(l.promptsDir, event.Name)
	if err != nil {
		l.logger.Error("failed to get relative path", zap.String("path", event.Name), zap.Error(err))
		return
	}

	// Normalize path separators
	relPath = filepath.ToSlash(relPath)

	switch {
	case event.Op&fsnotify.Write == fsnotify.Write:
		l.logger.Info("prompt file modified, reloading", zap.String("path", relPath))
		if err := l.loadPrompt(relPath, event.Name); err != nil {
			l.logger.Error("failed to reload prompt", zap.String("path", relPath), zap.Error(err))
		}

	case event.Op&fsnotify.Create == fsnotify.Create:
		l.logger.Info("prompt file created, loading", zap.String("path", relPath))
		if err := l.loadPrompt(relPath, event.Name); err != nil {
			l.logger.Error("failed to load new prompt", zap.String("path", relPath), zap.Error(err))
		}

	case event.Op&fsnotify.Remove == fsnotify.Remove:
		l.logger.Info("prompt file removed, deleting from cache", zap.String("path", relPath))
		l.cache.Delete(relPath)

	case event.Op&fsnotify.Rename == fsnotify.Rename:
		l.logger.Info("prompt file renamed, deleting from cache", zap.String("path", relPath))
		l.cache.Delete(relPath)
	}
}

// Close closes the file watcher
func (l *PromptLoader) Close() error {
	if l.watcher != nil {
		return l.watcher.Close()
	}
	return nil
}
