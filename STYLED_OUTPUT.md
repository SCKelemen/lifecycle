# Styled Output Integration

This document describes how to use beautiful terminal styling with the lifecycle events library, integrating ideas from [charmbracelet/log](https://github.com/charmbracelet/log).

## Overview

The lifecycle library now supports beautiful terminal styling while maintaining structured JSON output for log aggregation. This gives you:

- **Beautiful terminal output** - Colorful, formatted logs for development and debugging
- **Structured JSON** - Machine-readable logs for log aggregation systems
- **PII protection** - Automatic redaction still works with styled output
- **OpenTelemetry integration** - Traces and metrics still work

## Basic Usage

### Enable Styled Output

```go
package main

import (
	"context"
	"os"
	
	"github.com/SCKelemen/lifecycle"
)

func main() {
	// Create styled output for beautiful terminal logs
	styled := lifecycle.NewStyledOutput(os.Stdout)
	
	// Create producer with styled output
	producer := lifecycle.NewProducer(
		"user-service",
		"pod-123",
		lifecycle.WithStyledOutput(styled),
	)
	
	// Events will now be displayed with beautiful styling
	ctx := context.Background()
	producer.EmitServiceStarted(ctx, "v1.0.0", 12345)
	producer.EmitRequestReceived(ctx, "corr-123", "GET", "/api/users", nil)
}
```

### Dual Output: Styled Terminal + JSON File

For production, you often want styled output for developers but JSON for log aggregation:

```go
// Open JSON file for log aggregation
jsonFile, _ := os.OpenFile("logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
defer jsonFile.Close()

// Create styled output that also writes JSON
styled := lifecycle.NewStyledOutput(
	os.Stdout, // Styled output to terminal
	lifecycle.WithJSONOutput(jsonFile), // JSON to file
)

producer := lifecycle.NewProducer(
	"user-service",
	"pod-123",
	lifecycle.WithStyledOutput(styled),
)
```

### JSON-Only Mode

For log aggregation systems that don't need terminal styling:

```go
styled := lifecycle.NewStyledOutput(
	os.Stdout,
	lifecycle.WithJSONOnly(), // Only output JSON, no styling
)

producer := lifecycle.NewProducer(
	"user-service",
	"pod-123",
	lifecycle.WithStyledOutput(styled),
)
```

## Integration with Debug Library

The styled output works seamlessly with the debug library for conditional debug logging:

```go
import (
	"github.com/SCKelemen/debug"
	"github.com/SCKelemen/lifecycle"
)

func main() {
	// Create debug manager
	dm := debug.NewDebugManager()
	dm.SetFlags("api.request|db.query")
	
	// Create styled output
	styled := lifecycle.NewStyledOutput(os.Stdout)
	
	// Create producer
	producer := lifecycle.NewProducer(
		"user-service",
		"pod-123",
		lifecycle.WithStyledOutput(styled),
	)
	
	// Use debug flags to conditionally enable detailed logging
	if dm.IsEnabled(APIV1Request) {
		producer.EmitRequestReceived(ctx, correlationID, method, path, detailedMetadata)
	}
}
```

## Features from charmbracelet/log

The styled output integrates the following features from charmbracelet/log:

1. **Colorful Output** - Different colors for different log levels (info, warn, error)
2. **Structured Logging** - Key-value pairs displayed beautifully
3. **Level Indicators** - Visual indicators for log levels
4. **Timestamp Formatting** - Human-readable timestamps
5. **Field Highlighting** - Important fields are highlighted

## Example Output

### Styled Output (Terminal)

```
INFO  api.request.received  service=user-service api=examples.User correlation_id=corr-123 method=GET path=/api/users
INFO  api.request.handled   service=user-service api=examples.User correlation_id=corr-123 status_code=200 duration_ms=45
ERROR api.request.errored   service=user-service api=examples.User correlation_id=corr-124 error="database connection failed" status_code=500
```

### JSON Output (Log Aggregation)

```json
{"event_type":"api.request.received","timestamp":"2025-01-25T10:00:00Z","service":"user-service","api":"examples.User","correlation_id":"corr-123","method":"GET","path":"/api/users"}
{"event_type":"api.request.handled","timestamp":"2025-01-25T10:00:00Z","service":"user-service","api":"examples.User","correlation_id":"corr-123","status_code":200,"duration_ms":45}
{"event_type":"api.request.errored","timestamp":"2025-01-25T10:00:00Z","service":"user-service","api":"examples.User","correlation_id":"corr-124","error":"database connection failed","status_code":500}
```

## Configuration Options

### Custom Logger

```go
import "github.com/charmbracelet/log"

customLogger := log.New(os.Stdout)
customLogger.SetLevel(log.DebugLevel)

styled := lifecycle.NewStyledOutput(
	os.Stdout,
	lifecycle.WithStyledLogger(customLogger),
)
```

### Environment-Based Configuration

```go
func createProducer() *lifecycle.Producer {
	var styled *lifecycle.StyledOutput
	
	if os.Getenv("LOG_FORMAT") == "json" {
		// Production: JSON only
		styled = lifecycle.NewStyledOutput(
			os.Stdout,
			lifecycle.WithJSONOnly(),
		)
	} else {
		// Development: Styled output
		styled = lifecycle.NewStyledOutput(os.Stdout)
		
		// Also write JSON to file if configured
		if jsonPath := os.Getenv("LOG_JSON_FILE"); jsonPath != "" {
			jsonFile, _ := os.OpenFile(jsonPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			styled = lifecycle.NewStyledOutput(
				os.Stdout,
				lifecycle.WithJSONOutput(jsonFile),
			)
		}
	}
	
	return lifecycle.NewProducer(
		"user-service",
		"pod-123",
		lifecycle.WithStyledOutput(styled),
	)
}
```

## Benefits

1. **Developer Experience** - Beautiful, readable logs during development
2. **Production Ready** - JSON output for log aggregation systems
3. **No Breaking Changes** - Existing code continues to work
4. **PII Protection** - Automatic redaction still works
5. **OpenTelemetry** - Traces and metrics still work
6. **Debug Integration** - Works seamlessly with debug flags

## Migration Guide

### Before (Plain JSON)

```go
producer := lifecycle.NewProducer("user-service", "pod-123")
```

### After (Styled Output)

```go
styled := lifecycle.NewStyledOutput(os.Stdout)
producer := lifecycle.NewProducer(
	"user-service",
	"pod-123",
	lifecycle.WithStyledOutput(styled),
)
```

That's it! No other code changes needed.

