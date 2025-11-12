package internal

import (
	"time"

	"go.uber.org/zap"
)

// UnleashClient wraps the Unleash SDK client - STUB IMPLEMENTATION
// TODO: Implement full Unleash integration when SDK version is compatible
type UnleashClient struct {
	config *UnleashConfig
	logger *zap.Logger
}

// Feature represents a feature toggle
type Feature struct {
	Name        string
	Description string
	Enabled     bool
	CreatedAt   time.Time
}

// Variant represents a feature variant
type Variant struct {
	Name    string
	Enabled bool
	Payload VariantPayload
}

// VariantPayload represents variant payload
type VariantPayload struct {
	Type  string
	Value string
}

// NewUnleashClient creates a new Unleash client - STUB
func NewUnleashClient(config *UnleashConfig, logger *zap.Logger) (*UnleashClient, error) {
	logger.Warn("Using stub Unleash client - feature flags will return default values",
		zap.String("server_url", config.ServerURL),
		zap.String("app_name", config.AppName))
	
	return &UnleashClient{
		config: config,
		logger: logger,
	}, nil
}

// IsFeatureEnabled checks if a feature is enabled - STUB returns false
func (c *UnleashClient) IsFeatureEnabled(featureKey string, context *FeatureContext) bool {
	c.logger.Debug("feature flag check (stub)",
		zap.String("feature_key", featureKey),
		zap.Bool("enabled", false))
	return false
}

// IsEnabled checks if a feature is enabled with map context - STUB returns false
func (c *UnleashClient) IsEnabled(featureKey string, context map[string]interface{}) bool {
	c.logger.Debug("feature flag check (stub)",
		zap.String("feature_key", featureKey),
		zap.Bool("enabled", false))
	return false
}

// GetVariant gets a feature variant - STUB returns empty variant
func (c *UnleashClient) GetVariant(featureKey string, context *FeatureContext) Variant {
	c.logger.Debug("feature variant check (stub)",
		zap.String("feature_key", featureKey))
	return Variant{
		Name:    "",
		Enabled: false,
		Payload: VariantPayload{
			Type:  "string",
			Value: "{}",
		},
	}
}

// GetFeatureToggles returns all feature toggles - STUB returns empty list
func (c *UnleashClient) GetFeatureToggles() []Feature {
	c.logger.Debug("get feature toggles (stub)")
	return []Feature{}
}

// IsReady checks if the client is ready - STUB returns true
func (c *UnleashClient) IsReady() bool {
	return true
}

// WaitForReady waits for the client to be ready - STUB returns immediately
func (c *UnleashClient) WaitForReady(ctx interface{}) error {
	c.logger.Info("Unleash client ready (stub)")
	return nil
}

// Close closes the Unleash client
func (c *UnleashClient) Close() error {
	return nil
}
