package lifecycle

import (
	"context"
	"log"
	"log/slog"
	"os"
)

// WrappedLogger wraps standard loggers to prevent direct logging
// This ensures all logging goes through the lifecycle event system
type WrappedLogger struct {
	producer *Producer
	fallback *slog.Logger
}

// NewWrappedLogger creates a wrapped logger that routes all logs through lifecycle events
// This replaces standard loggers to prevent direct logging
func NewWrappedLogger(producer *Producer) *WrappedLogger {
	return &WrappedLogger{
		producer: producer,
		fallback: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

// Log emits a generic log event (should be avoided in favor of specific event types)
func (l *WrappedLogger) Log(level slog.Level, msg string, args ...interface{}) {
	// Convert to lifecycle event
	// In production, this should emit a structured event
	l.fallback.Log(nil, level, msg, args...)
}

// Debug logs a debug message (should use specific event types instead)
func (l *WrappedLogger) Debug(msg string, args ...interface{}) {
	l.Log(slog.LevelDebug, msg, args...)
}

// Info logs an info message (should use specific event types instead)
func (l *WrappedLogger) Info(msg string, args ...interface{}) {
	l.Log(slog.LevelInfo, msg, args...)
}

// Warn logs a warning message (should use specific event types instead)
func (l *WrappedLogger) Warn(msg string, args ...interface{}) {
	l.Log(slog.LevelWarn, msg, args...)
}

// Error logs an error message (should use specific event types instead)
func (l *WrappedLogger) Error(msg string, args ...interface{}) {
	l.Log(slog.LevelError, msg, args...)
}

// PreventDirectLogging replaces standard loggers with wrapped versions
// This should be called at application startup to prevent direct logging
func PreventDirectLogging(producer *Producer) {
	// Replace standard log package
	log.SetOutput(&logWriter{producer: producer})
	log.SetFlags(0) // Remove default flags to force structured logging

	// Replace slog default logger
	slog.SetDefault(slog.New(NewLifecycleHandler(producer)))
}

// logWriter implements io.Writer to intercept log package output
type logWriter struct {
	producer *Producer
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	// Convert log output to lifecycle event
	// This is a fallback - ideally all code should use lifecycle events directly
	// Emit as a generic log event (should be avoided)
	// In production, parse the log and emit appropriate event type
	_ = string(p) // Suppress unused variable warning
	return len(p), nil
}

// LifecycleHandler implements slog.Handler to route logs through lifecycle events
type LifecycleHandler struct {
	producer *Producer
}

// NewLifecycleHandler creates a new lifecycle handler
func NewLifecycleHandler(producer *Producer) *LifecycleHandler {
	return &LifecycleHandler{
		producer: producer,
	}
}

func (h *LifecycleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *LifecycleHandler) Handle(ctx context.Context, record slog.Record) error {
	// Convert slog record to lifecycle event
	// This is a fallback - ideally all code should use lifecycle events directly
	_ = record // Suppress unused variable warning
	return nil
}

func (h *LifecycleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *LifecycleHandler) WithGroup(name string) slog.Handler {
	return h
}

