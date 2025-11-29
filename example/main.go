package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SCKelemen/lifecycle"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	// Create color registry and register colors for services, APIs, and events
	// In a real application, these would come from API generator type definitions
	colorRegistry := lifecycle.NewColorRegistry()

	// Register service color
	colorRegistry.RegisterServiceColor("example-service", "#3B82F6") // Blue

	// Register API colors (from type annotations)
	colorRegistry.RegisterAPIColor("examples.User", "#10B981")  // Green
	colorRegistry.RegisterAPIColor("examples.Order", "#F59E0B") // Orange

	// Register event colors (from event type annotations)
	colorRegistry.RegisterEventColor("api.request.received", "#6366F1") // Indigo
	colorRegistry.RegisterEventColor("api.request.handled", "#10B981")  // Green
	colorRegistry.RegisterEventColor("api.request.errored", "#EF4444")  // Red
	colorRegistry.RegisterEventColor("service.started", "#3B82F6")      // Blue
	colorRegistry.RegisterEventColor("service.healthy", "#10B981")      // Green
	colorRegistry.RegisterEventColor("service.shutdown", "#6B7280")     // Gray
	colorRegistry.RegisterEventColor("db.query.started", "#8B5CF6")     // Purple
	colorRegistry.RegisterEventColor("db.query.completed", "#10B981")   // Green
	colorRegistry.RegisterEventColor("resource.created", "#10B981")     // Green

	// Create styled output with color registry
	// This will show beautiful colored terminal output
	styled := lifecycle.NewStyledOutput(
		os.Stdout,
		lifecycle.WithStyledColorRegistry(colorRegistry),
	)

	// Create lifecycle producer with styled output and color registry
	producer := lifecycle.NewProducer(
		"example-service",
		"pod-123",
		lifecycle.WithStyledOutput(styled),
		lifecycle.WithColorRegistry(colorRegistry),
		lifecycle.WithAPI("examples.User"), // Set default API
	)

	ctx := context.Background()

	fmt.Println("=== Service Lifecycle Events (with colors) ===")
	fmt.Println()

	// Service lifecycle events - service name will be colored
	producer.EmitServiceStarted(ctx, "v1.0.0", 12345)
	time.Sleep(100 * time.Millisecond)
	producer.EmitServiceHealthy(ctx, []string{"database", "cache"})

	// API events - API name and event type will be colored
	fmt.Println()
	fmt.Println("=== API Events (with colors) ===")
	fmt.Println()

	correlationID := "req-123"
	actor := lifecycle.NewHumanActor("user-456")
	resource := lifecycle.NewResource("User", "user-789")

	// Request received - event type and API will be colored
	producer.EmitRequestReceived(ctx, correlationID, "GET", "/api/users/user-789", nil, "examples.User")
	time.Sleep(10 * time.Millisecond)

	// Request handled - status code will be colored (green for 200)
	producer.EmitRequestHandled(ctx, correlationID, actor, resource, 200, 10, 1024, "examples.User")

	// Request errored - status code will be colored (red for 500)
	fmt.Println()
	producer.EmitRequestErrored(ctx, "req-456", "database connection failed", "DB_CONN_ERROR", 500, 50, "examples.User")

	// Database tracing events - event type will be colored
	fmt.Println()
	fmt.Println("=== Database Tracing Events (with colors) ===")
	fmt.Println()

	queryID := "query-001"
	producer.EmitQueryStarted(ctx, queryID, "SELECT * FROM users WHERE email = ?", []interface{}{"user@example.com"})
	time.Sleep(5 * time.Millisecond)
	producer.EmitQueryCompleted(ctx, queryID, 5, 1)

	// Resource events with PII redaction - status will be colored
	fmt.Println()
	fmt.Println("=== Resource Events (with colors and PII redaction) ===")
	fmt.Println()

	resourceData := map[string]interface{}{
		"id":    "user-789",
		"name":  "John Doe",
		"email": "john.doe@example.com", // This will be redacted
		"phone": "+1234567890",          // This will be redacted
	}

	// Schema annotations from API generator indicate which fields are PII
	schemaAnnotations := map[string]lifecycle.FieldAnnotations{
		"email": {PII: true, Encrypted: true, Redactable: true},
		"phone": {PII: true, Encrypted: true, Redactable: true},
		"name":  {PII: true, Redactable: true},
	}

	producer.EmitResourceCreated(ctx, correlationID, actor, resource, resourceData, schemaAnnotations, "examples.User")

	// Service shutdown - event type will be colored
	fmt.Println()
	fmt.Println("=== Service Shutdown (with colors) ===")
	fmt.Println()
	producer.EmitServiceShutdown(ctx, "graceful", 0)

	fmt.Println()
	fmt.Println("=== Summary ===")
	fmt.Println("Notice how:")
	fmt.Println("  - Service names are colored (blue)")
	fmt.Println("  - API names are colored (green/orange)")
	fmt.Println("  - Event types are colored (various colors)")
	fmt.Println("  - Status codes are colored (green for success, red for errors)")
	fmt.Println("  - Status values are colored (green for success, red for errors)")
	fmt.Println()
	fmt.Println("=== Color Support Check ===")

	// Test if colors are working
	testStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6"))
	fmt.Printf("If you see this text in BLUE, colors are working: %s\n", testStyle.Render("COLORED TEXT"))
	fmt.Println("If the text above is NOT blue, your terminal may not support colors.")
	fmt.Println("See example/TROUBLESHOOTING.md for help.")
	fmt.Println()
	fmt.Printf("Terminal type: %s\n", os.Getenv("TERM"))
	fmt.Printf("Color profile: %s\n", lipgloss.ColorProfile())
}
