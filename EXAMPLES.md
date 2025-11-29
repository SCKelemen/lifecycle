# Running Examples

## Lifecycle Example - Styled Output with Colors

Demonstrates styled terminal output with colors from API generator annotations.

### Quick Start

```bash
cd /Users/samuel.kelemen/Code/github.com/SCKelemen/lifecycle
go run ./example
```

### What It Shows

- ✅ **Service colors** - Service names colored (blue: `#3B82F6`)
- ✅ **API colors** - API names colored (green: `#10B981` for User, orange: `#F59E0B` for Order)
- ✅ **Event colors** - Event types colored (various colors)
- ✅ **Status colors** - HTTP status codes and status values automatically colored
- ✅ **Styled output** - Beautiful terminal formatting using `charmbracelet/log`

### Expected Output

You'll see colored output like:
- `INFO service.started` (blue event type, blue service name)
- `INFO api.request.received` (indigo event type, green API name)
- `INFO api.request.handled` (green event type, green status code 200)
- `ERRO api.request.errored` (red event type, red status code 500)

Colors are applied to:
- Event type names
- Service names
- API names
- Status codes (2xx=green, 4xx=orange, 5xx=red)
- Status values (success=green, error=red)

