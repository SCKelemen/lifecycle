package lifecycle

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

// Producer provides a high-level API for emitting structured lifecycle events
// It replaces standard loggers to prevent direct logging
// These are OBSERVABILITY events for engineers, NOT domain events
type Producer struct {
	service     string
	host        string
	logger      *slog.Logger
	output      io.Writer
	piiDetector *PIIDetector
	redactor    *Redactor
	otel        *OTelIntegration
}

// ProducerOption configures the Producer
type ProducerOption func(*Producer)

// WithLogger sets a custom logger (for internal logging only)
func WithLogger(logger *slog.Logger) ProducerOption {
	return func(p *Producer) {
		p.logger = logger
	}
}

// WithOutput sets a custom output writer (default: os.Stdout)
func WithOutput(output io.Writer) ProducerOption {
	return func(p *Producer) {
		p.output = output
	}
}

// WithPIIDetector sets a custom PII detector
func WithPIIDetector(detector *PIIDetector) ProducerOption {
	return func(p *Producer) {
		p.piiDetector = detector
	}
}

// WithRedactor sets a custom redactor
func WithRedactor(redactor *Redactor) ProducerOption {
	return func(p *Producer) {
		p.redactor = redactor
	}
}

// WithOTelIntegration sets OpenTelemetry integration
func WithOTelIntegration(otel *OTelIntegration) ProducerOption {
	return func(p *Producer) {
		p.otel = otel
	}
}

// NewProducer creates a new lifecycle event producer
// This replaces standard loggers - developers should use this instead of log.Printf, etc.
// These are OBSERVABILITY events for engineers, NOT domain events
func NewProducer(service, host string, opts ...ProducerOption) *Producer {
	p := &Producer{
		service:     service,
		host:        host,
		logger:      slog.Default(),
		output:      os.Stdout,
		piiDetector: NewPIIDetector(),
		redactor:    NewRedactor(),
		otel:        NewOTelIntegration(service),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// createBaseEvent creates a base event with common fields
func (p *Producer) createBaseEvent(eventType, correlationID string, metadata map[string]interface{}) *BaseEvent {
	base := &BaseEvent{
		EventType:     eventType,
		Timestamp:     time.Now(),
		Service:       p.service,
		Host:          p.host,
		CorrelationID: correlationID,
		Metadata:      metadata,
	}

	return base
}

// redactData redacts PII from data based on schema annotations
func (p *Producer) redactData(data map[string]interface{}, schemaAnnotations map[string]FieldAnnotations) map[string]interface{} {
	if data == nil {
		return nil
	}

	redacted := make(map[string]interface{})
	for key, value := range data {
		// Check if field has PII annotations
		annotations, hasAnnotations := schemaAnnotations[key]
		if hasAnnotations && (annotations.Encrypted || annotations.Redactable) {
			// Redact PII fields
			redacted[key] = p.redactor.Redact(value)
		} else {
			// Recursively check nested structures
			if nestedMap, ok := value.(map[string]interface{}); ok {
				redacted[key] = p.redactData(nestedMap, schemaAnnotations)
			} else {
				redacted[key] = value
			}
		}
	}

	return redacted
}

// emitEvent writes the event to the configured output as JSON
// Also creates OpenTelemetry spans and records metrics
func (p *Producer) emitEvent(ctx context.Context, event Event, duration time.Duration) error {
	// Redact PII before serialization
	if eventWithData, ok := event.(EventWithData); ok {
		eventWithData.RedactPII(p.piiDetector, p.redactor)
	}

	// Create OpenTelemetry span
	if p.otel != nil {
		attrs := EventAttributes(event)
		spanCtx, span := p.otel.StartSpan(ctx, event.GetEventType(), attrs...)
		defer span.End()

		// Record metrics
		p.otel.RecordMetric(spanCtx, event.GetEventType(), duration, attrs...)
	}

	// Emit structured log (JSON)
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if _, err := fmt.Fprintln(p.output, string(jsonData)); err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}

	return nil
}

// Service Lifecycle Events

// EmitServiceStarted emits a service.started event
func (p *Producer) EmitServiceStarted(ctx context.Context, version string, pid int32) error {
	event := &ServiceStartedEvent{
		Base:    p.createBaseEvent("service.started", "", nil),
		Version: version,
		PID:     pid,
	}
	return p.emitEvent(ctx, event, 0)
}

// EmitServiceHealthy emits a service.healthy event
func (p *Producer) EmitServiceHealthy(ctx context.Context, healthChecks []string) error {
	event := &ServiceHealthyEvent{
		Base:         p.createBaseEvent("service.healthy", "", nil),
		HealthChecks: healthChecks,
	}
	return p.emitEvent(ctx, event, 0)
}

// EmitServiceShutdown emits a service.shutdown event
func (p *Producer) EmitServiceShutdown(ctx context.Context, reason string, exitCode int32) error {
	event := &ServiceShutdownEvent{
		Base:     p.createBaseEvent("service.shutdown", "", nil),
		Reason:   reason,
		ExitCode: exitCode,
	}
	return p.emitEvent(ctx, event, 0)
}

// EmitServiceCrashed emits a service.crashed event
func (p *Producer) EmitServiceCrashed(ctx context.Context, reason, stackTrace string, exitCode int32) error {
	event := &ServiceCrashedEvent{
		Base:       p.createBaseEvent("service.crashed", "", nil),
		Reason:     reason,
		StackTrace: stackTrace,
		ExitCode:   exitCode,
	}
	return p.emitEvent(ctx, event, 0)
}

// API Events

// EmitRequestReceived emits an api.request.received event
func (p *Producer) EmitRequestReceived(ctx context.Context, correlationID, method, path string, metadata map[string]interface{}) error {
	event := &RequestReceivedEvent{
		Base:    p.createBaseEvent("api.request.received", correlationID, metadata),
		Method:  method,
		Path:    path,
		UserAgent: extractUserAgent(ctx),
		RemoteAddr: extractRemoteAddr(ctx),
	}
	return p.emitEvent(ctx, event, 0)
}

// EmitRequestHandled emits an api.request.handled event
func (p *Producer) EmitRequestHandled(ctx context.Context, correlationID string, actor *Actor, resource *Resource,
	statusCode int32, durationMs int64, responseSizeBytes int64) error {
	event := &RequestHandledEvent{
		Base:              p.createBaseEvent("api.request.handled", correlationID, nil),
		Actor:             actor,
		Resource:          resource,
		Status:            StatusSuccess,
		DurationMs:        durationMs,
		StatusCode:        statusCode,
		ResponseSizeBytes: responseSizeBytes,
	}
	return p.emitEvent(ctx, event, time.Duration(durationMs)*time.Millisecond)
}

// EmitRequestErrored emits an api.request.errored event
func (p *Producer) EmitRequestErrored(ctx context.Context, correlationID, errorMessage, errorCode string,
	statusCode int32, durationMs int64) error {
	event := &RequestErroredEvent{
		Base:         p.createBaseEvent("api.request.errored", correlationID, nil),
		Status:       StatusError,
		ErrorMessage: errorMessage,
		ErrorCode:    errorCode,
		StatusCode:   statusCode,
		DurationMs:   durationMs,
	}
	return p.emitEvent(ctx, event, time.Duration(durationMs)*time.Millisecond)
}

// EmitRequestRetried emits an api.request.retried event
func (p *Producer) EmitRequestRetried(ctx context.Context, correlationID string, retryCount int32,
	delayMs int64, retryReason string) error {
	event := &RequestRetriedEvent{
		Base:        p.createBaseEvent("api.request.retried", correlationID, nil),
		RetryCount:  retryCount,
		DelayMs:     delayMs,
		RetryReason: retryReason,
	}
	return p.emitEvent(ctx, event, time.Duration(delayMs)*time.Millisecond)
}

// Database Tracing Events

// EmitQueryStarted emits a db.query.started event
func (p *Producer) EmitQueryStarted(ctx context.Context, queryID, query string, params []interface{}) error {
	// Redact PII from query parameters
	redactedParams := p.redactor.RedactParams(params)
	
	event := &QueryStartedEvent{
		Base:   p.createBaseEvent("db.query.started", extractCorrelationID(ctx), nil),
		QueryID: queryID,
		Query:   query,
		Params:  redactedParams,
	}
	return p.emitEvent(ctx, event, 0)
}

// EmitQueryCompleted emits a db.query.completed event
func (p *Producer) EmitQueryCompleted(ctx context.Context, queryID string, durationMs int64, rowsAffected int64) error {
	event := &QueryCompletedEvent{
		Base:          p.createBaseEvent("db.query.completed", extractCorrelationID(ctx), nil),
		QueryID:       queryID,
		DurationMs:    durationMs,
		RowsAffected:  rowsAffected,
	}
	return p.emitEvent(ctx, event, time.Duration(durationMs)*time.Millisecond)
}

// EmitQueryErrored emits a db.query.errored event
func (p *Producer) EmitQueryErrored(ctx context.Context, queryID, errorMessage, errorCode string, durationMs int64) error {
	event := &QueryErroredEvent{
		Base:         p.createBaseEvent("db.query.errored", extractCorrelationID(ctx), nil),
		QueryID:      queryID,
		ErrorMessage: errorMessage,
		ErrorCode:    errorCode,
		DurationMs:   durationMs,
	}
	return p.emitEvent(ctx, event, time.Duration(durationMs)*time.Millisecond)
}

// EmitTransactionStarted emits a db.transaction.started event
func (p *Producer) EmitTransactionStarted(ctx context.Context, transactionID string) error {
	event := &TransactionStartedEvent{
		Base:          p.createBaseEvent("db.transaction.started", extractCorrelationID(ctx), nil),
		TransactionID: transactionID,
	}
	return p.emitEvent(ctx, event, 0)
}

// EmitTransactionCommitted emits a db.transaction.committed event
func (p *Producer) EmitTransactionCommitted(ctx context.Context, transactionID string, durationMs int64) error {
	event := &TransactionCommittedEvent{
		Base:          p.createBaseEvent("db.transaction.committed", extractCorrelationID(ctx), nil),
		TransactionID: transactionID,
		DurationMs:   durationMs,
	}
	return p.emitEvent(ctx, event, time.Duration(durationMs)*time.Millisecond)
}

// EmitTransactionRolledBack emits a db.transaction.rolled_back event
func (p *Producer) EmitTransactionRolledBack(ctx context.Context, transactionID, reason string, durationMs int64) error {
	event := &TransactionRolledBackEvent{
		Base:          p.createBaseEvent("db.transaction.rolled_back", extractCorrelationID(ctx), nil),
		TransactionID: transactionID,
		Reason:        reason,
		DurationMs:    durationMs,
	}
	return p.emitEvent(ctx, event, time.Duration(durationMs)*time.Millisecond)
}

// Resource Events

// EmitResourceCreated emits a resource.created event
func (p *Producer) EmitResourceCreated(ctx context.Context, correlationID string, actor *Actor,
	resource *Resource, resourceData map[string]interface{}, schemaAnnotations map[string]FieldAnnotations) error {
	// Redact PII from resource data
	redactedData := p.redactData(resourceData, schemaAnnotations)
	
	event := &ResourceCreatedEvent{
		Base:         p.createBaseEvent("resource.created", correlationID, nil),
		Actor:        actor,
		Resource:     resource,
		ResourceData: redactedData,
	}
	return p.emitEvent(ctx, event, 0)
}

// EmitResourceUpdated emits a resource.updated event
func (p *Producer) EmitResourceUpdated(ctx context.Context, correlationID string, actor *Actor,
	resource *Resource, previousData, newData map[string]interface{}, updatedFields []string, schemaAnnotations map[string]FieldAnnotations) error {
	// Redact PII from both previous and new data
	redactedPrevious := p.redactData(previousData, schemaAnnotations)
	redactedNew := p.redactData(newData, schemaAnnotations)
	
	event := &ResourceUpdatedEvent{
		Base:          p.createBaseEvent("resource.updated", correlationID, nil),
		Actor:         actor,
		Resource:      resource,
		PreviousData:  redactedPrevious,
		NewData:       redactedNew,
		UpdatedFields: updatedFields,
	}
	return p.emitEvent(ctx, event, 0)
}

// EmitResourceDeleted emits a resource.deleted event
func (p *Producer) EmitResourceDeleted(ctx context.Context, correlationID string, actor *Actor,
	resource *Resource, softDelete bool, finalData map[string]interface{}, schemaAnnotations map[string]FieldAnnotations) error {
	// Redact PII from final data
	redactedData := p.redactData(finalData, schemaAnnotations)
	
	event := &ResourceDeletedEvent{
		Base:       p.createBaseEvent("resource.deleted", correlationID, nil),
		Actor:      actor,
		Resource:   resource,
		SoftDelete: softDelete,
		FinalData:  redactedData,
	}
	return p.emitEvent(ctx, event, 0)
}

// Helper functions

// extractCorrelationID extracts correlation ID from context
func extractCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value("correlation_id").(string); ok {
		return id
	}
	return ""
}

// extractUserAgent extracts user agent from context
func extractUserAgent(ctx context.Context) string {
	if ua, ok := ctx.Value("user_agent").(string); ok {
		return ua
	}
	return ""
}

// extractRemoteAddr extracts remote address from context
func extractRemoteAddr(ctx context.Context) string {
	if addr, ok := ctx.Value("remote_addr").(string); ok {
		return addr
	}
	return ""
}

