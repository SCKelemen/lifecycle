package lifecycle

import "time"

// Event is the base interface for all lifecycle events
type Event interface {
	GetEventType() string
	GetTimestamp() time.Time
	GetService() string
	GetHost() string
	GetCorrelationID() string
}

// EventWithData is an event that contains data that may need PII redaction
type EventWithData interface {
	Event
	RedactPII(detector *PIIDetector, redactor *Redactor)
}

// BaseEvent contains common fields for all events
type BaseEvent struct {
	EventType     string                 `json:"event_type"`
	Timestamp     time.Time              `json:"timestamp"`
	Service       string                 `json:"service"`
	Host          string                 `json:"host"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

func (e *BaseEvent) GetEventType() string     { return e.EventType }
func (e *BaseEvent) GetTimestamp() time.Time  { return e.Timestamp }
func (e *BaseEvent) GetService() string       { return e.Service }
func (e *BaseEvent) GetHost() string          { return e.Host }
func (e *BaseEvent) GetCorrelationID() string { return e.CorrelationID }

// Actor represents the actor performing an action
type Actor struct {
	UserID    string    `json:"user_id"`
	ActorType ActorType `json:"actor_type"`
}

// ActorType represents the type of actor
type ActorType string

const (
	ActorTypeHuman     ActorType = "human"
	ActorTypeSystem    ActorType = "system"
	ActorTypeSynthetic ActorType = "synthetic"
)

// Resource represents a resource being acted upon
type Resource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Status represents the status of an operation
type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

// Service Lifecycle Events

// ServiceStartedEvent represents a service.started event
type ServiceStartedEvent struct {
	Base    *BaseEvent `json:"base"`
	Version string     `json:"version"`
	PID     int32      `json:"pid"`
}

func (e *ServiceStartedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *ServiceStartedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *ServiceStartedEvent) GetService() string       { return e.Base.GetService() }
func (e *ServiceStartedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *ServiceStartedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// ServiceHealthyEvent represents a service.healthy event
type ServiceHealthyEvent struct {
	Base         *BaseEvent `json:"base"`
	HealthChecks []string   `json:"health_checks"`
}

func (e *ServiceHealthyEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *ServiceHealthyEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *ServiceHealthyEvent) GetService() string       { return e.Base.GetService() }
func (e *ServiceHealthyEvent) GetHost() string          { return e.Base.GetHost() }
func (e *ServiceHealthyEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// ServiceShutdownEvent represents a service.shutdown event
type ServiceShutdownEvent struct {
	Base     *BaseEvent `json:"base"`
	Reason   string     `json:"reason"`
	ExitCode int32      `json:"exit_code"`
}

func (e *ServiceShutdownEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *ServiceShutdownEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *ServiceShutdownEvent) GetService() string       { return e.Base.GetService() }
func (e *ServiceShutdownEvent) GetHost() string          { return e.Base.GetHost() }
func (e *ServiceShutdownEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// ServiceCrashedEvent represents a service.crashed event
type ServiceCrashedEvent struct {
	Base       *BaseEvent `json:"base"`
	Reason     string     `json:"reason"`
	StackTrace string     `json:"stack_trace"`
	ExitCode   int32      `json:"exit_code"`
}

func (e *ServiceCrashedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *ServiceCrashedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *ServiceCrashedEvent) GetService() string       { return e.Base.GetService() }
func (e *ServiceCrashedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *ServiceCrashedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// API Events

// RequestReceivedEvent represents an api.request.received event
type RequestReceivedEvent struct {
	Base       *BaseEvent `json:"base"`
	Method     string     `json:"method"`
	Path       string     `json:"path"`
	UserAgent  string     `json:"user_agent,omitempty"`
	RemoteAddr string     `json:"remote_addr,omitempty"`
}

func (e *RequestReceivedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *RequestReceivedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *RequestReceivedEvent) GetService() string       { return e.Base.GetService() }
func (e *RequestReceivedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *RequestReceivedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// RequestHandledEvent represents an api.request.handled event
type RequestHandledEvent struct {
	Base              *BaseEvent `json:"base"`
	Actor             *Actor     `json:"actor,omitempty"`
	Resource          *Resource  `json:"resource,omitempty"`
	Status            Status     `json:"status"`
	DurationMs        int64      `json:"duration_ms"`
	StatusCode        int32      `json:"status_code"`
	ResponseSizeBytes int64      `json:"response_size_bytes,omitempty"`
}

func (e *RequestHandledEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *RequestHandledEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *RequestHandledEvent) GetService() string       { return e.Base.GetService() }
func (e *RequestHandledEvent) GetHost() string          { return e.Base.GetHost() }
func (e *RequestHandledEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// RequestErroredEvent represents an api.request.errored event
type RequestErroredEvent struct {
	Base         *BaseEvent `json:"base"`
	Status       Status     `json:"status"`
	ErrorMessage string     `json:"error_message"`
	ErrorCode    string     `json:"error_code,omitempty"`
	StatusCode   int32      `json:"status_code"`
	DurationMs   int64      `json:"duration_ms"`
}

func (e *RequestErroredEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *RequestErroredEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *RequestErroredEvent) GetService() string       { return e.Base.GetService() }
func (e *RequestErroredEvent) GetHost() string          { return e.Base.GetHost() }
func (e *RequestErroredEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// RequestRetriedEvent represents an api.request.retried event
type RequestRetriedEvent struct {
	Base        *BaseEvent `json:"base"`
	RetryCount  int32      `json:"retry_count"`
	DelayMs     int64      `json:"delay_ms"`
	RetryReason string     `json:"retry_reason,omitempty"`
}

func (e *RequestRetriedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *RequestRetriedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *RequestRetriedEvent) GetService() string       { return e.Base.GetService() }
func (e *RequestRetriedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *RequestRetriedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// Database Tracing Events

// QueryStartedEvent represents a db.query.started event
type QueryStartedEvent struct {
	Base    *BaseEvent    `json:"base"`
	QueryID string        `json:"query_id"`
	Query   string        `json:"query"`
	Params  []interface{} `json:"params,omitempty"`
}

func (e *QueryStartedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *QueryStartedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *QueryStartedEvent) GetService() string       { return e.Base.GetService() }
func (e *QueryStartedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *QueryStartedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// QueryCompletedEvent represents a db.query.completed event
type QueryCompletedEvent struct {
	Base         *BaseEvent `json:"base"`
	QueryID      string     `json:"query_id"`
	DurationMs   int64      `json:"duration_ms"`
	RowsAffected int64      `json:"rows_affected,omitempty"`
}

func (e *QueryCompletedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *QueryCompletedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *QueryCompletedEvent) GetService() string       { return e.Base.GetService() }
func (e *QueryCompletedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *QueryCompletedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// QueryErroredEvent represents a db.query.errored event
type QueryErroredEvent struct {
	Base         *BaseEvent `json:"base"`
	QueryID      string     `json:"query_id"`
	ErrorMessage string     `json:"error_message"`
	ErrorCode    string     `json:"error_code,omitempty"`
	DurationMs   int64      `json:"duration_ms"`
}

func (e *QueryErroredEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *QueryErroredEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *QueryErroredEvent) GetService() string       { return e.Base.GetService() }
func (e *QueryErroredEvent) GetHost() string          { return e.Base.GetHost() }
func (e *QueryErroredEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// TransactionStartedEvent represents a db.transaction.started event
type TransactionStartedEvent struct {
	Base          *BaseEvent `json:"base"`
	TransactionID string     `json:"transaction_id"`
}

func (e *TransactionStartedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *TransactionStartedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *TransactionStartedEvent) GetService() string       { return e.Base.GetService() }
func (e *TransactionStartedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *TransactionStartedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// TransactionCommittedEvent represents a db.transaction.committed event
type TransactionCommittedEvent struct {
	Base          *BaseEvent `json:"base"`
	TransactionID string     `json:"transaction_id"`
	DurationMs    int64      `json:"duration_ms"`
}

func (e *TransactionCommittedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *TransactionCommittedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *TransactionCommittedEvent) GetService() string       { return e.Base.GetService() }
func (e *TransactionCommittedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *TransactionCommittedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// TransactionRolledBackEvent represents a db.transaction.rolled_back event
type TransactionRolledBackEvent struct {
	Base          *BaseEvent `json:"base"`
	TransactionID string     `json:"transaction_id"`
	Reason        string     `json:"reason,omitempty"`
	DurationMs    int64      `json:"duration_ms"`
}

func (e *TransactionRolledBackEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *TransactionRolledBackEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *TransactionRolledBackEvent) GetService() string       { return e.Base.GetService() }
func (e *TransactionRolledBackEvent) GetHost() string          { return e.Base.GetHost() }
func (e *TransactionRolledBackEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

// Resource Events

// ResourceCreatedEvent represents a resource.created event
type ResourceCreatedEvent struct {
	Base         *BaseEvent             `json:"base"`
	Actor        *Actor                 `json:"actor,omitempty"`
	Resource     *Resource              `json:"resource"`
	ResourceData map[string]interface{} `json:"resource_data,omitempty"`
}

func (e *ResourceCreatedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *ResourceCreatedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *ResourceCreatedEvent) GetService() string       { return e.Base.GetService() }
func (e *ResourceCreatedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *ResourceCreatedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

func (e *ResourceCreatedEvent) RedactPII(detector *PIIDetector, redactor *Redactor) {
	if e.ResourceData != nil {
		e.ResourceData = redactor.RedactMap(e.ResourceData, detector)
	}
}

// ResourceUpdatedEvent represents a resource.updated event
type ResourceUpdatedEvent struct {
	Base          *BaseEvent             `json:"base"`
	Actor         *Actor                 `json:"actor,omitempty"`
	Resource      *Resource              `json:"resource"`
	PreviousData  map[string]interface{} `json:"previous_data,omitempty"`
	NewData       map[string]interface{} `json:"new_data,omitempty"`
	UpdatedFields []string               `json:"updated_fields,omitempty"`
}

func (e *ResourceUpdatedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *ResourceUpdatedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *ResourceUpdatedEvent) GetService() string       { return e.Base.GetService() }
func (e *ResourceUpdatedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *ResourceUpdatedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

func (e *ResourceUpdatedEvent) RedactPII(detector *PIIDetector, redactor *Redactor) {
	if e.PreviousData != nil {
		e.PreviousData = redactor.RedactMap(e.PreviousData, detector)
	}
	if e.NewData != nil {
		e.NewData = redactor.RedactMap(e.NewData, detector)
	}
}

// ResourceDeletedEvent represents a resource.deleted event
type ResourceDeletedEvent struct {
	Base       *BaseEvent             `json:"base"`
	Actor      *Actor                 `json:"actor,omitempty"`
	Resource   *Resource              `json:"resource"`
	SoftDelete bool                   `json:"soft_delete"`
	FinalData  map[string]interface{} `json:"final_data,omitempty"`
}

func (e *ResourceDeletedEvent) GetEventType() string     { return e.Base.GetEventType() }
func (e *ResourceDeletedEvent) GetTimestamp() time.Time  { return e.Base.GetTimestamp() }
func (e *ResourceDeletedEvent) GetService() string       { return e.Base.GetService() }
func (e *ResourceDeletedEvent) GetHost() string          { return e.Base.GetHost() }
func (e *ResourceDeletedEvent) GetCorrelationID() string { return e.Base.GetCorrelationID() }

func (e *ResourceDeletedEvent) RedactPII(detector *PIIDetector, redactor *Redactor) {
	if e.FinalData != nil {
		e.FinalData = redactor.RedactMap(e.FinalData, detector)
	}
}

// FieldAnnotations represents field-level annotations from the schema system
type FieldAnnotations struct {
	Encrypted  bool `json:"encrypted"`
	Redactable bool `json:"redactable"`
}
