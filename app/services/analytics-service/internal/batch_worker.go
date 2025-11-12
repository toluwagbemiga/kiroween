package internal

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// ExternalProvider defines the interface for external analytics providers
type ExternalProvider interface {
	SendBatch(ctx context.Context, events []Event) error
	GetName() string
}

// BatchWorker processes events in batches
type BatchWorker struct {
	queue        *BatchQueue
	provider     ExternalProvider
	retryConfig  *RetryConfig
	flushTimer   *time.Timer
	flushInterval time.Duration
	logger       *zap.Logger
	stopChan     chan struct{}
	doneChan     chan struct{}
}

// NewBatchWorker creates a new batch worker
func NewBatchWorker(
	queue *BatchQueue,
	provider ExternalProvider,
	flushInterval time.Duration,
	retryConfig *RetryConfig,
	logger *zap.Logger,
) *BatchWorker {
	return &BatchWorker{
		queue:         queue,
		provider:      provider,
		retryConfig:   retryConfig,
		flushInterval: flushInterval,
		logger:        logger,
		stopChan:      make(chan struct{}),
		doneChan:      make(chan struct{}),
	}
}

// Start starts the batch worker
func (w *BatchWorker) Start() {
	w.logger.Info("batch worker started",
		zap.Duration("flush_interval", w.flushInterval),
		zap.Int("max_retry_attempts", w.retryConfig.MaxAttempts))

	// Start flush timer
	w.flushTimer = time.NewTimer(w.flushInterval)

	go w.run()
}

// run is the main worker loop
func (w *BatchWorker) run() {
	defer close(w.doneChan)

	for {
		select {
		case <-w.queue.FlushChannel():
			// Batch size reached - flush immediately
			w.logger.Debug("batch size reached, flushing")
			w.flush()
			w.resetTimer()

		case <-w.flushTimer.C:
			// Timer expired - flush if we have events
			w.logger.Debug("flush timer expired")
			w.flush()
			w.resetTimer()

		case <-w.stopChan:
			// Shutdown requested - final flush
			w.logger.Info("batch worker stopping, performing final flush")
			w.flush()
			return
		}
	}
}

// flush processes the current batch
func (w *BatchWorker) flush() {
	batch := w.queue.GetBatch()
	if len(batch) == 0 {
		w.logger.Debug("no events to flush")
		return
	}

	w.logger.Info("flushing batch",
		zap.Int("event_count", len(batch)),
		zap.String("provider", w.provider.GetName()))

	// Send batch with retry logic
	if err := w.sendBatchWithRetry(context.Background(), batch); err != nil {
		w.logger.Error("failed to flush batch after retries",
			zap.Int("event_count", len(batch)),
			zap.Error(err))
		// TODO: Consider dead letter queue for failed batches
	} else {
		w.logger.Info("batch flushed successfully",
			zap.Int("event_count", len(batch)))
	}
}

// sendBatchWithRetry sends a batch with exponential backoff retry
func (w *BatchWorker) sendBatchWithRetry(ctx context.Context, batch []Event) error {
	var lastErr error
	delay := w.retryConfig.InitialDelay

	for attempt := 1; attempt <= w.retryConfig.MaxAttempts; attempt++ {
		err := w.provider.SendBatch(ctx, batch)
		if err == nil {
			if attempt > 1 {
				w.logger.Info("batch sent successfully after retry",
					zap.Int("attempt", attempt))
			}
			return nil
		}

		lastErr = err

		// Don't retry on last attempt
		if attempt == w.retryConfig.MaxAttempts {
			break
		}

		w.logger.Warn("batch send failed, retrying",
			zap.Int("attempt", attempt),
			zap.Duration("delay", delay),
			zap.Error(err))

		// Wait before retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		// Exponential backoff
		delay = time.Duration(float64(delay) * w.retryConfig.BackoffFactor)
		if delay > w.retryConfig.MaxDelay {
			delay = w.retryConfig.MaxDelay
		}
	}

	return fmt.Errorf("max retry attempts exceeded: %w", lastErr)
}

// resetTimer resets the flush timer
func (w *BatchWorker) resetTimer() {
	if !w.flushTimer.Stop() {
		select {
		case <-w.flushTimer.C:
		default:
		}
	}
	w.flushTimer.Reset(w.flushInterval)
}

// Stop stops the batch worker gracefully
func (w *BatchWorker) Stop() {
	w.logger.Info("stopping batch worker")
	close(w.stopChan)
	<-w.doneChan
	w.logger.Info("batch worker stopped")
}
