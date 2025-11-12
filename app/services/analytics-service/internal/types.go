package internal

import (
	"sync"
	"time"
)

// Event represents an analytics event
type Event struct {
	ID         string
	EventName  string
	UserID     string
	Properties map[string]interface{}
	Timestamp  time.Time
	CreatedAt  time.Time
}

// UserIdentity represents a user's identity and properties
type UserIdentity struct {
	UserID     string
	Properties map[string]interface{}
	FirstSeen  time.Time
	LastSeen   time.Time
	UpdatedAt  time.Time
}

// EventBatch represents a batch of events to be flushed
type EventBatch struct {
	Events    []Event
	CreatedAt time.Time
}

// BatchQueue is a thread-safe queue for batching events
type BatchQueue struct {
	mu         sync.Mutex
	events     []Event
	maxSize    int
	flushTimer *time.Timer
	flushChan  chan struct{}
}

// NewBatchQueue creates a new batch queue
func NewBatchQueue(maxSize int) *BatchQueue {
	return &BatchQueue{
		events:    make([]Event, 0, maxSize),
		maxSize:   maxSize,
		flushChan: make(chan struct{}, 1),
	}
}

// Add adds an event to the queue
func (q *BatchQueue) Add(event Event) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.events = append(q.events, event)

	// Check if we've reached batch size
	if len(q.events) >= q.maxSize {
		// Trigger flush
		select {
		case q.flushChan <- struct{}{}:
		default:
			// Flush already pending
		}
		return true
	}

	return false
}

// GetBatch returns and clears the current batch
func (q *BatchQueue) GetBatch() []Event {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.events) == 0 {
		return nil
	}

	batch := make([]Event, len(q.events))
	copy(batch, q.events)
	q.events = q.events[:0] // Clear the slice

	return batch
}

// Size returns the current queue size
func (q *BatchQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.events)
}

// FlushChannel returns the flush notification channel
func (q *BatchQueue) FlushChannel() <-chan struct{} {
	return q.flushChan
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
		MaxAttempts:   5,
		InitialDelay:  1 * time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
	}
}
