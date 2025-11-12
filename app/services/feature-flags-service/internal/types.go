package internal

import (
	"encoding/json"
	"time"
)

// FeatureContext represents the context for feature flag evaluation
type FeatureContext struct {
	UserID     string
	TeamID     string
	Properties map[string]interface{}
	RemoteAddr string
	UserAgent  string
	SessionID  string
}

// FeatureFlag represents a feature flag configuration
type FeatureFlag struct {
	Name        string
	Enabled     bool
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// FeatureFlagResult represents the result of a feature flag evaluation
type FeatureFlagResult struct {
	FeatureName string
	Enabled     bool
	Variant     string
	Context     *FeatureContext
	EvaluatedAt time.Time
}

// UnleashConfig holds Unleash SDK configuration
type UnleashConfig struct {
	ServerURL       string
	APIToken        string
	AppName         string
	InstanceID      string
	RefreshInterval time.Duration
	MetricsInterval time.Duration
	DisableMetrics  bool
}

// ContextProperty represents a property in the feature flag context
type ContextProperty struct {
	Key   string
	Value interface{}
}

// ParsePropertiesJSON parses a JSON string into a map of properties
func ParsePropertiesJSON(propertiesJSON string) (map[string]interface{}, error) {
	if propertiesJSON == "" {
		return make(map[string]interface{}), nil
	}

	var properties map[string]interface{}
	err := json.Unmarshal([]byte(propertiesJSON), &properties)
	if err != nil {
		return nil, err
	}

	return properties, nil
}

// ToJSON converts properties map to JSON string
func ToJSON(properties map[string]interface{}) (string, error) {
	if len(properties) == 0 {
		return "{}", nil
	}

	jsonBytes, err := json.Marshal(properties)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
