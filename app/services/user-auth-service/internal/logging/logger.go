package logging

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger
func NewLogger(level string) (*Logger, error) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}
	
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	
	return &Logger{Logger: logger}, nil
}

// WithCorrelationID adds a correlation ID to the logger
func (l *Logger) WithCorrelationID(ctx context.Context, correlationID string) *zap.Logger {
	return l.With(zap.String("correlation_id", correlationID))
}

// AuditEvent represents an audit log event
type AuditEvent struct {
	EventType     string
	UserID        string
	Email         string
	IPAddress     string
	Success       bool
	ErrorReason   string
	CorrelationID string
	Metadata      map[string]interface{}
}

// LogAuditEvent logs an audit event
func (l *Logger) LogAuditEvent(event *AuditEvent) {
	fields := []zap.Field{
		zap.String("event_type", event.EventType),
		zap.String("user_id", event.UserID),
		zap.String("email", event.Email),
		zap.String("ip_address", event.IPAddress),
		zap.Bool("success", event.Success),
		zap.String("correlation_id", event.CorrelationID),
	}
	
	if event.ErrorReason != "" {
		fields = append(fields, zap.String("error_reason", event.ErrorReason))
	}
	
	if event.Metadata != nil {
		for key, value := range event.Metadata {
			fields = append(fields, zap.Any(key, value))
		}
	}
	
	if event.Success {
		l.Info("audit_event", fields...)
	} else {
		l.Warn("audit_event", fields...)
	}
}

// GetLogLevel returns the current log level from environment
func GetLogLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		return "info"
	}
	return level
}
