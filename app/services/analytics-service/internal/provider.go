package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// MixpanelProvider implements the ExternalProvider interface for Mixpanel
type MixpanelProvider struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
	logger     *zap.Logger
	testMode   bool
}

// NewMixpanelProvider creates a new Mixpanel provider
func NewMixpanelProvider(apiKey string, testMode bool, logger *zap.Logger) *MixpanelProvider {
	return &MixpanelProvider{
		apiKey: apiKey,
		apiURL: "https://api.mixpanel.com/track",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger:   logger,
		testMode: testMode,
	}
}

// SendBatch sends a batch of events to Mixpanel
func (p *MixpanelProvider) SendBatch(ctx context.Context, events []Event) error {
	if p.testMode {
		p.logger.Info("TEST MODE: would send batch to Mixpanel",
			zap.Int("event_count", len(events)))
		// Simulate processing time
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	// Convert events to Mixpanel format
	mixpanelEvents := make([]map[string]interface{}, len(events))
	for i, event := range events {
		mixpanelEvents[i] = map[string]interface{}{
			"event": event.EventName,
			"properties": map[string]interface{}{
				"distinct_id": event.UserID,
				"time":        event.Timestamp.Unix(),
				"$insert_id":  event.ID, // For deduplication
			},
		}

		// Merge custom properties
		if event.Properties != nil {
			for k, v := range event.Properties {
				mixpanelEvents[i]["properties"].(map[string]interface{})[k] = v
			}
		}
	}

	// Prepare request
	payload := map[string]interface{}{
		"api_key": p.apiKey,
		"events":  mixpanelEvents,
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

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mixpanel API returned status %d", resp.StatusCode)
	}

	p.logger.Debug("batch sent to Mixpanel",
		zap.Int("event_count", len(events)),
		zap.Int("status_code", resp.StatusCode))

	return nil
}

// GetName returns the provider name
func (p *MixpanelProvider) GetName() string {
	return "mixpanel"
}
