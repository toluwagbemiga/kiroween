package internal

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// UsageTracker tracks LLM usage for analytics
type UsageTracker struct {
	store  *UsageStore
	logger *zap.Logger
	// analyticsClient would go here for real analytics service integration
}

// NewUsageTracker creates a new usage tracker
func NewUsageTracker(maxStoreSize int, logger *zap.Logger) *UsageTracker {
	return &UsageTracker{
		store:  NewUsageStore(maxStoreSize),
		logger: logger,
	}
}

// TrackUsage tracks a usage event
func (t *UsageTracker) TrackUsage(ctx context.Context, event *UsageEvent) error {
	// Store locally
	t.store.Add(*event)

	// Log the event
	t.logger.Info("usage tracked",
		zap.String("request_id", event.RequestID),
		zap.String("prompt_path", event.PromptPath),
		zap.String("calling_service", event.CallingService),
		zap.String("model", event.Model),
		zap.Int32("total_tokens", event.TotalTokens),
		zap.Bool("success", event.Success))

	// Send to analytics service asynchronously (fire and forget)
	go func() {
		// Note: In production, you'd initialize this client once and reuse it
		// For now, this demonstrates the integration pattern
		// The actual client should be passed to UsageTracker constructor
		t.logger.Debug("usage event logged locally",
			zap.String("calling_service", event.CallingService),
			zap.String("model", event.Model),
			zap.Int32("total_tokens", event.TotalTokens))
	}()

	return nil
}

// GetStats returns usage statistics
func (t *UsageTracker) GetStats(timeRange string, serviceFilter string) (*UsageStats, error) {
	// Parse time range
	var since time.Time
	switch timeRange {
	case "hour":
		since = time.Now().Add(-1 * time.Hour)
	case "day":
		since = time.Now().Add(-24 * time.Hour)
	case "week":
		since = time.Now().Add(-7 * 24 * time.Hour)
	default:
		since = time.Now().Add(-24 * time.Hour) // Default to day
	}

	return t.store.GetStats(since, serviceFilter), nil
}

// UsageStore stores usage events in memory
type UsageStore struct {
	mu      sync.RWMutex
	events  []UsageEvent
	maxSize int
}

// NewUsageStore creates a new usage store
func NewUsageStore(maxSize int) *UsageStore {
	return &UsageStore{
		events:  make([]UsageEvent, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add adds a usage event
func (s *UsageStore) Add(event UsageEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, event)

	// Prune old events if we exceed max size
	if len(s.events) > s.maxSize {
		// Keep only the most recent events
		s.events = s.events[len(s.events)-s.maxSize:]
	}
}

// GetRecent returns recent events
func (s *UsageStore) GetRecent(limit int) []UsageEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit > len(s.events) {
		limit = len(s.events)
	}

	// Return most recent events
	start := len(s.events) - limit
	result := make([]UsageEvent, limit)
	copy(result, s.events[start:])
	return result
}

// GetStats returns aggregated statistics
func (s *UsageStore) GetStats(since time.Time, serviceFilter string) *UsageStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &UsageStats{
		RequestsByService: make(map[string]int64),
		TokensByModel:     make(map[string]int64),
	}

	var totalResponseTime int64
	var responseCount int64

	for _, event := range s.events {
		// Filter by time
		if event.Timestamp.Before(since) {
			continue
		}

		// Filter by service if specified
		if serviceFilter != "" && event.CallingService != serviceFilter {
			continue
		}

		// Aggregate stats
		stats.TotalRequests++
		stats.TotalTokens += int64(event.TotalTokens)
		stats.RequestsByService[event.CallingService]++
		stats.TokensByModel[event.Model] += int64(event.TotalTokens)

		if event.Success {
			totalResponseTime += event.ResponseTimeMs
			responseCount++
		}
	}

	// Calculate average response time
	if responseCount > 0 {
		stats.AverageResponseMs = totalResponseTime / responseCount
	}

	return stats
}

// Prune removes events older than the specified time
func (s *UsageStore) Prune(olderThan time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find first event to keep
	keepIdx := 0
	for i, event := range s.events {
		if event.Timestamp.After(olderThan) {
			keepIdx = i
			break
		}
	}

	// Remove old events
	if keepIdx > 0 {
		s.events = s.events[keepIdx:]
	}
}
