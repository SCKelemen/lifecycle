# Lifecycle Example - Styled Output with Colors

This example demonstrates the lifecycle library's styled output with color integration from API generator annotations.

## Running the Example

From the `lifecycle` repository root:

```bash
go run ./example
```

Or build and run:

```bash
go build -o example/main ./example
./example/main
```

## What You'll See

The example demonstrates:

1. **Service Lifecycle Events** - Service names are colored (blue)
2. **API Events** - API names are colored (green/orange), event types are colored
3. **Status Codes** - HTTP status codes are automatically colored (green for 2xx, red for 5xx)
4. **Database Tracing** - Query events are colored (purple)
5. **Resource Events** - Status values are colored (green for created, etc.)

## Color Registration

Colors are registered from API generator type/event annotations:

- **Service colors**: `#3B82F6` (blue)
- **API colors**: `#10B981` (green) for `examples.User`, `#F59E0B` (orange) for `examples.Order`
- **Event colors**: Various colors for different event types
- **Status colors**: Automatic coloring based on status (success=green, error=red)

## Output Format

The styled output uses `charmbracelet/log` for beautiful terminal formatting with:
- Colored event types
- Colored service names
- Colored API names
- Colored status codes and status values
- Structured key-value pairs

