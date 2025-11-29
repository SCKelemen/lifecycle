# Color Integration Guide

This document describes how to wire color annotations from the API generator through all tools and services.

## Overview

Colors are defined in type/event annotations in the API generator:
- **API Types**: `annotations: - color: "#RRGGBB"` on resource types
- **Event Types**: `annotations: - color: "#RRGGBB"` on event types
- **Services**: Can be configured via environment or configuration
- **Statuses**: Default colors for common statuses (success, error, etc.)

These colors flow through:
1. **Lifecycle Events** - Service, API, event, and status colors in styled output
2. **Debug Library** - Path colors for debug flags
3. **All Tools** - CLI, TUI, logs, metrics dashboards

## Loading Colors from API Generator

### Step 1: Extract Colors from Type Definitions

```go
package main

import (
	"github.com/SCKelemen/api/internal/schema"
	"github.com/SCKelemen/lifecycle"
)

func loadColorsFromTypes(typesPath string) (*lifecycle.ColorRegistry, error) {
	// Load all type files
	typeFiles, err := schema.LoadAllTypeFiles(typesPath)
	if err != nil {
		return nil, err
	}

	registry := lifecycle.NewColorRegistry()

	// Extract colors from type definitions
	for _, tf := range typeFiles {
		// Get color from annotations
		color := extractColorFromAnnotations(tf.Spec.Annotations)
		if color == "" && tf.Spec.Color != "" {
			color = tf.Spec.Color // Fallback to deprecated field
		}

		if color != "" {
			if tf.Kind == "Event" {
				// Register event color
				registry.RegisterEventColor(tf.Spec.Type, color)
			} else {
				// Register API color
				registry.RegisterAPIColor(tf.Spec.Type, color)
			}
		}
	}

	return registry, nil
}

func extractColorFromAnnotations(annotations []schema.Annotation) string {
	for _, ann := range annotations {
		if ann.Color != nil {
			switch v := ann.Color.(type) {
			case string:
				return v
			case map[string]interface{}:
				if val, ok := v["value"].(string); ok {
					return val
				}
				if val, ok := v["color"].(string); ok {
					return val
				}
			}
		}
	}
	return ""
}
```

### Step 2: Configure Lifecycle Producer with Colors

```go
package main

import (
	"github.com/SCKelemen/lifecycle"
)

func setupProducer(service, host string) (*lifecycle.Producer, error) {
	// Load colors from type definitions
	colorRegistry, err := loadColorsFromTypes("./types")
	if err != nil {
		return nil, err
	}

	// Register service color (can come from config/env)
	colorRegistry.RegisterServiceColor(service, "#3B82F6") // Blue

	// Create styled output with color registry
	styled := lifecycle.NewStyledOutput(
		os.Stdout,
		lifecycle.WithColorRegistry(colorRegistry),
		lifecycle.WithJSONOutput(jsonFile), // Optional: also write JSON
	)

	// Create producer with styled output
	producer := lifecycle.NewProducer(
		service,
		host,
		lifecycle.WithStyledOutput(styled),
		lifecycle.WithColorRegistry(colorRegistry),
	)

	return producer, nil
}
```

### Step 3: Configure Debug Manager with Colors

```go
package main

import (
	"github.com/SCKelemen/debug"
	"github.com/SCKelemen/api/internal/schema"
)

func setupDebugManager(typesPath string) (*debug.DebugManager, error) {
	// Load type files to extract colors
	typeFiles, err := schema.LoadAllTypeFiles(typesPath)
	if err != nil {
		return nil, err
	}

	// Create debug manager
	dm := debug.NewDebugManager(debug.NewV2Parser())

	// Register flags with colors
	definitions := []debug.FlagDefinition{}
	
	// Register API flags with colors
	for _, tf := range typeFiles {
		if tf.Kind == "Type" {
			color := extractColorFromAnnotations(tf.Spec.Annotations)
			definitions = append(definitions, debug.FlagDefinition{
				Name:  fmt.Sprintf("api.%s", tf.Spec.Type),
				Flag:  debug.DebugFlag(1 << len(definitions)),
				Path:  fmt.Sprintf("api.%s", tf.Spec.Type),
				Color: color,
			})
		}
	}

	// Register event flags with colors
	for _, tf := range typeFiles {
		if tf.Kind == "Event" {
			color := extractColorFromAnnotations(tf.Spec.Annotations)
			definitions = append(definitions, debug.FlagDefinition{
				Name:  fmt.Sprintf("event.%s", tf.Spec.Type),
				Flag:  debug.DebugFlag(1 << len(definitions)),
				Path:  fmt.Sprintf("event.%s", tf.Spec.Type),
				Color: color,
			})
		}
	}

	dm.RegisterFlags(definitions)
	return dm, nil
}
```

## Color Usage in Styled Output

When colors are registered, they automatically appear in lifecycle event output:

```
INFO  api.request.received  service=user-service api=examples.User correlation_id=corr-123 method=GET path=/api/users
      ^^^^^^^^^^^^^^^^^^^^^  ^^^^^^^^^^^^^^^^^^^^ ^^^^^^^^^^^^^^^^^
      (colored by event)     (colored by service) (colored by API)
```

## Color Usage in Debug Logs

Debug logs show colors for paths:

```
DEBUG [api.examples.User] Creating user with email: user@example.com
      ^^^^^^^^^^^^^^^^^^^
      (colored by API type)
```

## Service Configuration

### Environment Variables

```bash
# Service color
export SERVICE_COLOR="#3B82F6"

# Load colors from types directory
export TYPES_PATH="./types"
```

### Configuration File

```yaml
# config.yaml
service:
  name: "user-service"
  color: "#3B82F6"  # Optional: override default

colors:
  # Override specific API colors
  apis:
    "examples.User": "#10B981"
    "examples.Order": "#F59E0B"
  
  # Override specific event colors
  events:
    "examples.OrderCreated": "#10B981"
    "examples.OrderCancelled": "#EF4444"
```

## Integration in Generated Services

Generated services should automatically:

1. **Load colors** from type definitions at startup
2. **Register colors** with lifecycle and debug libraries
3. **Use colors** in all log output

Example generated service initialization:

```go
// generated/main.go
func main() {
	// Load type definitions
	typeFiles, _ := schema.LoadAllTypeFiles("./types")
	
	// Extract and register colors
	colorRegistry := extractColorsFromTypes(typeFiles)
	
	// Setup lifecycle producer with colors
	styled := lifecycle.NewStyledOutput(
		os.Stdout,
		lifecycle.WithColorRegistry(colorRegistry),
	)
	producer := lifecycle.NewProducer(
		"user-service",
		"pod-123",
		lifecycle.WithStyledOutput(styled),
		lifecycle.WithColorRegistry(colorRegistry),
	)
	
	// Setup debug manager with colors
	dm := setupDebugManagerWithColors(typeFiles)
	
	// Use producer and debug manager throughout service
	// Colors will automatically appear in all logs
}
```

## Color Format

Colors should be in CSS-compatible hex format: `#RRGGBB`

Examples:
- `#3B82F6` - Blue
- `#10B981` - Green
- `#F59E0B` - Orange
- `#EF4444` - Red
- `#8B5CF6` - Purple

## Status Colors

Default status colors (can be overridden):

- `success` / `completed` → Green (`#00FF00`)
- `error` / `failed` → Red (`#FF0000`)
- `warning` → Orange (`#FFA500`)
- `info` → Blue (`#00BFFF`)
- `pending` → Yellow (`#FFFF00`)
- `in_progress` → Purple (`#9370DB`)
- `created` → Blue (`#00BFFF`)
- `updated` → Orange (`#FFA500`)
- `deleted` → Red (`#FF0000`)

HTTP status codes are automatically colored:
- 2xx → Green (success)
- 3xx → Blue (info)
- 4xx → Orange (warning)
- 5xx → Red (error)

## Benefits

1. **Visual Consistency** - Same colors across all tools and services
2. **Easy Identification** - Quickly spot services, APIs, and events in logs
3. **Better UX** - Colorful logs are easier to read and parse
4. **Automatic** - Colors flow from type definitions automatically
5. **Configurable** - Can override colors via config if needed

