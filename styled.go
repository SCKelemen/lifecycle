package lifecycle

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/log"
)

// StyledOutput provides beautiful terminal styling for lifecycle events
// while maintaining structured JSON output for log aggregation
type StyledOutput struct {
	logger        *log.Logger
	jsonOutput    io.Writer      // Separate JSON output for log aggregation
	jsonOnly      bool           // If true, only output JSON (no styling)
	colorRegistry *ColorRegistry // Color registry for services, APIs, events, statuses
}

// StyledOutputOption configures the styled output
type StyledOutputOption func(*StyledOutput)

// WithJSONOutput sets a separate writer for JSON output (for log aggregation)
// When set, styled output goes to terminal, JSON goes to this writer
func WithJSONOutput(w io.Writer) StyledOutputOption {
	return func(s *StyledOutput) {
		s.jsonOutput = w
	}
}

// WithJSONOnly disables styling and only outputs JSON
func WithJSONOnly() StyledOutputOption {
	return func(s *StyledOutput) {
		s.jsonOnly = true
	}
}

// WithStyledLogger sets a custom charmbracelet/log logger
func WithStyledLogger(logger *log.Logger) StyledOutputOption {
	return func(s *StyledOutput) {
		s.logger = logger
	}
}

// WithStyledColorRegistry sets a color registry for styled output
func WithStyledColorRegistry(registry *ColorRegistry) StyledOutputOption {
	return func(s *StyledOutput) {
		s.colorRegistry = registry
	}
}

// NewStyledOutput creates a new styled output handler
func NewStyledOutput(w io.Writer, opts ...StyledOutputOption) *StyledOutput {
	s := &StyledOutput{
		logger:        log.New(w),
		jsonOutput:    nil, // No separate JSON output by default
		jsonOnly:      false,
		colorRegistry: NewColorRegistry(), // Default color registry
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WriteEvent writes a lifecycle event with beautiful styling
// Also writes JSON to jsonOutput if configured
func (s *StyledOutput) WriteEvent(event Event) error {
	// Always write JSON if jsonOutput is configured (for log aggregation)
	if s.jsonOutput != nil {
		jsonData, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}
		if _, err := fmt.Fprintln(s.jsonOutput, string(jsonData)); err != nil {
			return fmt.Errorf("failed to write JSON event: %w", err)
		}
	}

	// If JSON-only mode, skip styling
	if s.jsonOnly {
		return nil
	}

	// Write styled output to terminal
	return s.writeStyledEvent(event)
}

// writeStyledEvent writes a beautifully styled version of the event
func (s *StyledOutput) writeStyledEvent(event Event) error {
	eventType := event.GetEventType()

	// Determine log level from event type
	level := s.eventTypeToLevel(eventType)

	// Get event color from registry
	eventColor := ""
	if s.colorRegistry != nil {
		eventColor = s.colorRegistry.GetEventColor(eventType)
	}

	// Build key-value pairs for structured logging
	fields := s.buildFields(event, eventColor)

	// Format event type with color if available
	styledEventType := eventType
	if eventColor != "" {
		styledEventType = FormatWithColor(eventType, eventColor)
	}

	// Use charmbracelet/log's structured logging
	switch level {
	case log.DebugLevel:
		s.logger.Debug(styledEventType, fields...)
	case log.InfoLevel:
		s.logger.Info(styledEventType, fields...)
	case log.WarnLevel:
		s.logger.Warn(styledEventType, fields...)
	case log.ErrorLevel:
		s.logger.Error(styledEventType, fields...)
	case log.FatalLevel:
		s.logger.Fatal(styledEventType, fields...)
	default:
		s.logger.Info(styledEventType, fields...)
	}

	return nil
}

// eventTypeToLevel maps event types to log levels
func (s *StyledOutput) eventTypeToLevel(eventType string) log.Level {
	switch {
	case contains(eventType, "error", "errored", "failed", "crashed"):
		return log.ErrorLevel
	case contains(eventType, "warn", "warning"):
		return log.WarnLevel
	case contains(eventType, "debug", "trace"):
		return log.DebugLevel
	case contains(eventType, "started", "completed", "handled", "created", "updated", "deleted"):
		return log.InfoLevel
	default:
		return log.InfoLevel
	}
}

// buildFields extracts key-value pairs from the event for structured logging
// Colors are applied to service, API, and status fields
func (s *StyledOutput) buildFields(event Event, eventColor string) []interface{} {
	fields := []interface{}{}

	// Add base event fields from Event interface with colors
	if service := event.GetService(); service != "" {
		serviceColor := ""
		if s.colorRegistry != nil {
			serviceColor = s.colorRegistry.GetServiceColor(service)
		}
		if serviceColor != "" {
			fields = append(fields, "service", FormatWithColor(service, serviceColor))
		} else {
			fields = append(fields, "service", service)
		}
	}
	if api := event.GetAPI(); api != "" {
		apiColor := ""
		if s.colorRegistry != nil {
			apiColor = s.colorRegistry.GetAPIColor(api)
		}
		if apiColor != "" {
			fields = append(fields, "api", FormatWithColor(api, apiColor))
		} else {
			fields = append(fields, "api", api)
		}
	}
	if host := event.GetHost(); host != "" {
		fields = append(fields, "host", host)
	}
	if correlationID := event.GetCorrelationID(); correlationID != "" {
		fields = append(fields, "correlation_id", correlationID)
	}
	if !event.GetTimestamp().IsZero() {
		fields = append(fields, "timestamp", event.GetTimestamp().Format(time.RFC3339))
	}

	// Add event-specific fields based on event type (with status colors)
	s.addEventSpecificFields(event, &fields)

	return fields
}

// addEventSpecificFields extracts fields from specific event types
func (s *StyledOutput) addEventSpecificFields(event Event, fields *[]interface{}) {
	switch e := event.(type) {
	case *ServiceStartedEvent:
		if e != nil && e.Base != nil {
			if e.Version != "" {
				*fields = append(*fields, "version", e.Version)
			}
			if e.PID > 0 {
				*fields = append(*fields, "pid", e.PID)
			}
		}

	case *ServiceShutdownEvent:
		if e != nil && e.Base != nil {
			if e.Reason != "" {
				*fields = append(*fields, "reason", e.Reason)
			}
			if e.ExitCode != 0 {
				*fields = append(*fields, "exit_code", e.ExitCode)
			}
		}

	case *ServiceCrashedEvent:
		if e != nil && e.Base != nil {
			if e.Reason != "" {
				*fields = append(*fields, "reason", e.Reason)
			}
			if e.StackTrace != "" {
				*fields = append(*fields, "stack_trace", e.StackTrace)
			}
			if e.ExitCode != 0 {
				*fields = append(*fields, "exit_code", e.ExitCode)
			}
		}

	case *RequestReceivedEvent:
		if e != nil && e.Base != nil {
			if e.Method != "" {
				*fields = append(*fields, "method", e.Method)
			}
			if e.Path != "" {
				*fields = append(*fields, "path", e.Path)
			}
			if e.UserAgent != "" {
				*fields = append(*fields, "user_agent", e.UserAgent)
			}
			if e.RemoteAddr != "" {
				*fields = append(*fields, "remote_addr", e.RemoteAddr)
			}
		}

	case *RequestHandledEvent:
		if e != nil && e.Base != nil {
			if e.StatusCode > 0 {
				statusStr := fmt.Sprintf("%d", e.StatusCode)
				// Color status code based on HTTP status
				statusColor := s.getStatusCodeColor(e.StatusCode)
				if statusColor != "" {
					*fields = append(*fields, "status_code", FormatWithColor(statusStr, statusColor))
				} else {
					*fields = append(*fields, "status_code", e.StatusCode)
				}
			}
			if e.DurationMs > 0 {
				*fields = append(*fields, "duration_ms", e.DurationMs)
			}
			if e.ResponseSizeBytes > 0 {
				*fields = append(*fields, "response_size_bytes", e.ResponseSizeBytes)
			}
			if e.Actor != nil && e.Actor.UserID != "" {
				*fields = append(*fields, "actor", e.Actor.UserID)
			}
			if e.Resource != nil && e.Resource.ID != "" {
				*fields = append(*fields, "resource", e.Resource.ID)
			}
			// Add status with color
			if e.Status != "" {
				statusColor := ""
				if s.colorRegistry != nil {
					statusColor = s.colorRegistry.GetStatusColor(string(e.Status))
				}
				if statusColor != "" {
					*fields = append(*fields, "status", FormatWithColor(string(e.Status), statusColor))
				} else {
					*fields = append(*fields, "status", string(e.Status))
				}
			}
		}

	case *RequestErroredEvent:
		if e != nil && e.Base != nil {
			if e.StatusCode > 0 {
				statusStr := fmt.Sprintf("%d", e.StatusCode)
				statusColor := s.getStatusCodeColor(e.StatusCode)
				if statusColor != "" {
					*fields = append(*fields, "status_code", FormatWithColor(statusStr, statusColor))
				} else {
					*fields = append(*fields, "status_code", e.StatusCode)
				}
			}
			if e.DurationMs > 0 {
				*fields = append(*fields, "duration_ms", e.DurationMs)
			}
			if e.ErrorMessage != "" {
				*fields = append(*fields, "error", e.ErrorMessage)
			}
			if e.ErrorCode != "" {
				*fields = append(*fields, "error_code", e.ErrorCode)
			}
			// Add status with color (error status)
			if e.Status != "" {
				statusColor := ""
				if s.colorRegistry != nil {
					statusColor = s.colorRegistry.GetStatusColor(string(e.Status))
				}
				if statusColor != "" {
					*fields = append(*fields, "status", FormatWithColor(string(e.Status), statusColor))
				} else {
					*fields = append(*fields, "status", string(e.Status))
				}
			}
		}

	case *QueryStartedEvent:
		if e != nil && e.Base != nil {
			if e.QueryID != "" {
				*fields = append(*fields, "query_id", e.QueryID)
			}
			if e.Query != "" {
				*fields = append(*fields, "query", e.Query)
			}
		}

	case *QueryCompletedEvent:
		if e != nil && e.Base != nil {
			if e.QueryID != "" {
				*fields = append(*fields, "query_id", e.QueryID)
			}
			if e.DurationMs > 0 {
				*fields = append(*fields, "duration_ms", e.DurationMs)
			}
			if e.RowsAffected > 0 {
				*fields = append(*fields, "rows_affected", e.RowsAffected)
			}
		}

	case *QueryErroredEvent:
		if e != nil && e.Base != nil {
			if e.QueryID != "" {
				*fields = append(*fields, "query_id", e.QueryID)
			}
			if e.DurationMs > 0 {
				*fields = append(*fields, "duration_ms", e.DurationMs)
			}
			if e.ErrorMessage != "" {
				*fields = append(*fields, "error", e.ErrorMessage)
			}
			if e.ErrorCode != "" {
				*fields = append(*fields, "error_code", e.ErrorCode)
			}
		}

	case *ResourceCreatedEvent:
		if e != nil && e.Base != nil {
			if e.Resource != nil && e.Resource.ID != "" {
				*fields = append(*fields, "resource", e.Resource.ID)
			}
			// Status is "created"
			statusColor := ""
			if s.colorRegistry != nil {
				statusColor = s.colorRegistry.GetStatusColor("created")
			}
			if statusColor != "" {
				*fields = append(*fields, "status", FormatWithColor("created", statusColor))
			} else {
				*fields = append(*fields, "status", "created")
			}
		}
	case *ResourceUpdatedEvent:
		if e != nil && e.Base != nil {
			if e.Resource != nil && e.Resource.ID != "" {
				*fields = append(*fields, "resource", e.Resource.ID)
			}
			// Status is "updated"
			statusColor := ""
			if s.colorRegistry != nil {
				statusColor = s.colorRegistry.GetStatusColor("updated")
			}
			if statusColor != "" {
				*fields = append(*fields, "status", FormatWithColor("updated", statusColor))
			} else {
				*fields = append(*fields, "status", "updated")
			}
		}
	case *ResourceDeletedEvent:
		if e != nil && e.Base != nil {
			if e.Resource != nil && e.Resource.ID != "" {
				*fields = append(*fields, "resource", e.Resource.ID)
			}
			// Status is "deleted"
			statusColor := ""
			if s.colorRegistry != nil {
				statusColor = s.colorRegistry.GetStatusColor("deleted")
			}
			if statusColor != "" {
				*fields = append(*fields, "status", FormatWithColor("deleted", statusColor))
			} else {
				*fields = append(*fields, "status", "deleted")
			}
		}
	}
}

// getStatusCodeColor returns a color for HTTP status codes
func (s *StyledOutput) getStatusCodeColor(statusCode int32) string {
	if s.colorRegistry == nil {
		return ""
	}

	switch {
	case statusCode >= 200 && statusCode < 300:
		return s.colorRegistry.GetStatusColor("success")
	case statusCode >= 300 && statusCode < 400:
		return s.colorRegistry.GetStatusColor("info")
	case statusCode >= 400 && statusCode < 500:
		return s.colorRegistry.GetStatusColor("warning")
	case statusCode >= 500:
		return s.colorRegistry.GetStatusColor("error")
	default:
		return ""
	}
}

// contains checks if any of the substrings are contained in the string
func contains(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}
