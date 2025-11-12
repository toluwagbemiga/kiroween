package internal

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// MultiProvider sends events to multiple providers in parallel
type MultiProvider struct {
	providers []ExternalProvider
	logger    *zap.Logger
}

// NewMultiProvider creates a new multi-provider
func NewMultiProvider(providers []ExternalProvider, logger *zap.Logger) *MultiProvider {
	return &MultiProvider{
		providers: providers,
		logger:    logger,
	}
}

// SendBatch sends a batch to all providers in parallel
func (p *MultiProvider) SendBatch(ctx context.Context, events []Event) error {
	if len(p.providers) == 0 {
		return fmt.Errorf("no providers configured")
	}

	// Send to all providers in parallel
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(p.providers))
	successCount := 0
	var mu sync.Mutex

	for _, provider := range p.providers {
		wg.Add(1)
		go func(prov ExternalProvider) {
			defer wg.Done()

			if err := prov.SendBatch(ctx, events); err != nil {
				p.logger.Error("provider send failed",
					zap.String("provider", prov.GetName()),
					zap.Int("event_count", len(events)),
					zap.Error(err))
				errorsChan <- fmt.Errorf("%s: %w", prov.GetName(), err)
			} else {
				p.logger.Debug("provider send succeeded",
					zap.String("provider", prov.GetName()),
					zap.Int("event_count", len(events)))
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(provider)
	}

	wg.Wait()
	close(errorsChan)

	// Collect errors
	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	// Return error if all providers failed
	if len(errors) == len(p.providers) {
		return fmt.Errorf("all providers failed: %v", errors)
	}

	// Partial success is OK - log warnings for failed providers
	if len(errors) > 0 {
		p.logger.Warn("some providers failed",
			zap.Int("success_count", successCount),
			zap.Int("failure_count", len(errors)),
			zap.Errors("errors", errors))
	}

	return nil
}

// GetName returns the combined provider names
func (p *MultiProvider) GetName() string {
	if len(p.providers) == 0 {
		return "none"
	}

	names := make([]string, len(p.providers))
	for i, provider := range p.providers {
		names[i] = provider.GetName()
	}
	return strings.Join(names, ",")
}
