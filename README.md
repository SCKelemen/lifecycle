# Lifecycle Events Library

A structured observability library that prevents direct logging and enforces lifecycle events for engineers. Compatible with OpenTelemetry, structured logging, and metrics.

## Overview

This library replaces ad-hoc log statements with typed, structured observability events. These are **NOT domain events** - they are **observability events for engineers** to monitor, debug, and understand system behavior.

### Architecture Support

The library supports complex service architectures:
- **Multiple APIs per service**: A single service can host multiple APIs (e.g., `user-service` hosting both `examples.User` and `examples.Order` APIs)
- **APIs across services**: A single API can span multiple services (e.g., `examples.User` API implemented across `user-service`, `user-cache-service`, and `user-search-service`)

Events include both `service` (service instance) and `api` (API identifier) fields to enable filtering and aggregation at both levels.

It enforces structured observability to:

1. **Prevent direct logging** - Developers cannot use standard loggers directly
2. **Protect sensitive data** - Automatic PII detection and redaction based on schema annotations
3. **Enable tooling** - Structured events enable powerful developer tooling
4. **Ensure compliance** - GDPR-compliant logging with automatic PII handling
5. **OpenTelemetry integration** - Native support for traces, spans, and metrics
6. **Structured logging** - Compatible with log aggregation systems
7. **Metrics** - Built-in metrics support (counters, gauges, histograms)

## Key Features

- ✅ **No Direct Logging** - Wraps standard loggers to prevent unstructured logs
- ✅ **PII Detection & Redaction** - Automatic detection and redaction based on field annotations
- ✅ **OpenTelemetry Compatible** - Native integration with OTel traces, spans, and metrics
- ✅ **Structured Logging** - JSON logs compatible with log aggregation systems
- ✅ **Metrics Support** - Built-in counters, gauges, and histograms
- ✅ **Multiple Event Categories** - API events, service lifecycle, DB tracing, etc.
- ✅ **Schema Integration** - Integrates with API schema system for field-level PII detection

## Event Categories

### Service Lifecycle
- `service.started` - Service startup
- `service.healthy` - Health check passed
- `service.shutdown` - Graceful shutdown
- `service.crashed` - Unexpected crash

### API Events
- `api.request.received` - HTTP/gRPC request received
- `api.request.handled` - Request handled successfully
- `api.request.errored` - Request failed
- `api.request.retried` - Request retried

### Database Tracing
- `db.query.started` - Database query started
- `db.query.completed` - Query completed successfully
- `db.query.errored` - Query failed
- `db.transaction.started` - Transaction started
- `db.transaction.committed` - Transaction committed
- `db.transaction.rolled_back` - Transaction rolled back

### Resource Events
- `resource.created` - Resource created
- `resource.updated` - Resource updated
- `resource.deleted` - Resource deleted

## OpenTelemetry Integration

Events automatically create OpenTelemetry spans and metrics:

```go
producer := lifecycle.NewProducer("my-service", "pod-123")

// Creates an OTel span and logs structured event
producer.EmitRequestReceived(ctx, "req-123", "GET", "/api/users", nil)

// Automatically records metrics
// - api.request.count (counter)
// - api.request.duration (histogram)
// - api.request.size (histogram)
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/SCKelemen/lifecycle"
)

func main() {
    // Create lifecycle producer (replaces standard logger)
    // Service: "user-service-pod-123" (service instance)
    // API: "examples.User" (optional - can be set per-event or via WithAPI)
    producer := lifecycle.NewProducer("user-service-pod-123", "pod-123",
        lifecycle.WithAPI("examples.User"), // Optional: set default API
    )
    
    // Service lifecycle (creates OTel span) - no API field
    producer.EmitServiceStarted(context.Background(), "v1.0.0", 12345)
    
    // API events (creates OTel span + metrics) - API inferred from producer or resource
    producer.EmitRequestReceived(context.Background(), 
        "req-123", "GET", "/api/users", nil, "examples.User") // API can be specified per-event
    
    // Multiple APIs in same service
    producer.EmitRequestReceived(context.Background(), 
        "req-124", "GET", "/api/orders", nil, "examples.Order") // Different API
    
    // DB tracing (creates OTel span)
    producer.EmitQueryStarted(context.Background(), 
        "query-456", "SELECT * FROM users WHERE id = ?", []interface{}{123})
}
```

## PII Handling

The library automatically detects and redacts PII based on schema annotations from the API generator:

```go
// Fields marked with `field: { pii: true, encrypted: true, redactable: true }` 
// in the API schema are automatically redacted in events
schemaAnnotations := map[string]lifecycle.FieldAnnotations{
    "email": {PII: true, Encrypted: true, Redactable: true},
    "phone": {PII: true, Encrypted: true, Redactable: true},
    "name":  {PII: true, Redactable: true},
}

producer.EmitResourceCreated(ctx, correlationID, actor, resource, 
    map[string]interface{}{
        "email": "user@example.com",  // Will be redacted (PII=true)
        "phone": "+1234567890",        // Will be redacted (PII=true)
        "name": "John Doe",            // Will be redacted (PII=true)
        "id": "user-123",              // Not redacted (no PII annotation)
    }, schemaAnnotations)
```

### Schema Integration

The library integrates with the API schema system's `FieldFlags` annotations:
- **`pii: true`** - Field contains personally identifiable information
- **`encrypted: true`** - Field requires encryption (also triggers redaction)
- **`redactable: true`** - Field can be redacted for GDPR Article 17
- **`sensitive: true`** - General sensitive data flag

Fields are redacted if they have **any** of these flags set. The library also falls back to pattern-based detection if schema annotations are not provided.

## OpenTelemetry Setup

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "github.com/SCKelemen/lifecycle"
)

// Initialize OTel
exporter, _ := otlptracegrpc.New(ctx)
tp := trace.NewTracerProvider(trace.WithBatcher(exporter))
otel.SetTracerProvider(tp)

// Create producer with OTel integration
producer := lifecycle.NewProducer("my-service", "pod-123",
    lifecycle.WithTracerProvider(tp),
    lifecycle.WithMeterProvider(otel.GetMeterProvider()),
)
```

## Structured Logging

Events are automatically logged as structured JSON:

```json
{
  "event_type": "api.request.received",
  "timestamp": "2025-01-15T10:30:00Z",
  "service": "my-service",
  "host": "pod-123",
  "correlation_id": "req-123",
  "method": "GET",
  "path": "/api/users"
}
```

## Metrics

Automatic metrics are recorded for all events:

- **Counters**: `api.request.count`, `db.query.count`, etc.
- **Histograms**: `api.request.duration`, `db.query.duration`, etc.
- **Gauges**: `service.health.status`, etc.

## Integration with Generated Services

The library integrates with generated services from the API schema tool:

```go
// Generated services automatically use lifecycle events
// Direct logging is prevented - all logs go through lifecycle events
```

## License

MIT
