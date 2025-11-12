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

// AmplitudeProvider implements the ExternalProvider interface for Amplitude
type AmplitudeProvider struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
	logger     *zap.Logger
	testMode   bool
}

// NewAmplitudeProvider creates a new Amplitude provider
func NewAmplitudeProvider(apiKey string, testMode bool, logger *zap.Logger) *AmplitudeProvider {
	return &AmplitudeProvider{
		apiKey: apiKey,
		apiURL: "https://api2.amplitude.com/2/httpapi",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger:   logger,
		testMode: testMode,
	}
}

// SendBatch sends a batch of events to Amplitude
func (p *AmplitudeProvider) SendBatch(ctx context.Context, events []Event) error {
	if p.testMode {
		p.logger.Info("TEST MODE: would send batch to Amplitude",
			zap.Int("event_count", len(events)))
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	// Convert events to Amplitude format
	amplitudeEvents := make([]map[string]interface{}, len(events))
	for i, event := range events {
		amplitudeEvent := map[string]interface{}{
			"user_id":    event.UserID,
			"event_type": event.EventName,
			"time":       event.Timestamp.UnixMilli(),
			"insert_id":  event.ID, // For deduplication
		}

		// Add event properties
		if event.Properties != nil {
			amplitudeEvent["event_properties"] = event.Properties
		}

		// Handle identify events
		if event.EventName == "$identify" {
			amplitudeEvent["user_properties"] = event.Properties
		}

		amplitudeEvents[i] = amplitudeEvent
	}

	// Prepare request
	payload := map[string]interface{}{
		"api_key": p.apiKey,
		"events":  amplitudeEvents,
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
		return fmt.Errorf("amplitude API returned status %d", resp.StatusCode)
	}

	p.logger.Debug("batch sent to Amplitude",
		zap.Int("event_count", len(events)),
		zap.Int("status_code", resp.StatusCode))

	return nil
}

// GetName returns the provider name
func (p *AmplitudeProvider) GetName() string {
	return "amplitude"
}
