package internal

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// SegmentProvider implements the ExternalProvider interface for Segment
type SegmentProvider struct {
	writeKey   string
	apiURL     string
	httpClient *http.Client
	logger     *zap.Logger
	testMode   bool
}

// NewSegmentProvider creates a new Segment provider
func NewSegmentProvider(writeKey string, testMode bool, logger *zap.Logger) *SegmentProvider {
	return &SegmentProvider{
		writeKey: writeKey,
		apiURL:   "https://api.segment.io/v1/batch",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger:   logger,
		testMode: testMode,
	}
}

// SendBatch sends a batch of events to Segment
func (p *SegmentProvider) SendBatch(ctx context.Context, events []Event) error {
	if p.testMode {
		p.logger.Info("TEST MODE: would send batch to Segment",
			zap.Int("event_count", len(events)))
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	// Convert events to Segment format
	segmentBatch := make([]map[string]interface{}, len(events))
	for i, event := range events {
		// Determine event type
		eventType := "track"
		if event.EventName == "$identify" {
			eventType = "identify"
		}

		segmentEvent := map[string]interface{}{
			"type":      eventType,
			"userId":    event.UserID,
			"messageId": event.ID,
			"timestamp": event.Timestamp.Format(time.RFC3339),
		}

		if eventType == "track" {
			segmentEvent["event"] = event.EventName
			segmentEvent["properties"] = event.Properties
		} else {
			// Identify event
			segmentEvent["traits"] = event.Properties
		}

		segmentBatch[i] = segmentEvent
	}

	// Prepare request
	payload := map[string]interface{}{
		"batch": segmentBatch,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	// Send HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", p.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Segment uses Basic Auth with write key as username
	auth := base64.StdEncoding.EncodeToString([]byte(p.writeKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("segment API returned status %d", resp.StatusCode)
	}

	p.logger.Debug("batch sent to Segment",
		zap.Int("event_count", len(events)),
		zap.Int("status_code", resp.StatusCode))

	return nil
}

// GetName returns the provider name
func (p *SegmentProvider) GetName() string {
	return "segment"
}
