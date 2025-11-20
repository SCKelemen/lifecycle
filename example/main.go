package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SCKelemen/lifecycle"
)

func main() {
	// Create lifecycle producer (replaces standard logger)
	producer := lifecycle.NewProducer("example-service", "pod-123")

	ctx := context.Background()

	// Service lifecycle events
	fmt.Println("=== Service Lifecycle Events ===")
	producer.EmitServiceStarted(ctx, "v1.0.0", 12345)
	producer.EmitServiceHealthy(ctx, []string{"database", "cache"})

	// API events
	fmt.Println("\n=== API Events ===")
	correlationID := "req-123"
	actor := lifecycle.NewHumanActor("user-456")
	resource := lifecycle.NewResource("User", "user-789")

	producer.EmitRequestReceived(ctx, correlationID, "GET", "/api/users/user-789", nil)
	time.Sleep(10 * time.Millisecond) // Simulate processing
	producer.EmitRequestHandled(ctx, correlationID, actor, resource, 200, 10, 1024)

	// Database tracing events
	fmt.Println("\n=== Database Tracing Events ===")
	queryID := "query-001"
	producer.EmitQueryStarted(ctx, queryID, "SELECT * FROM users WHERE email = ?", []interface{}{"user@example.com"})
	time.Sleep(5 * time.Millisecond)
	producer.EmitQueryCompleted(ctx, queryID, 5, 1)

	// Resource events with PII redaction
	fmt.Println("\n=== Resource Events (with PII redaction) ===")
	resourceData := map[string]interface{}{
		"id":    "user-789",
		"name":  "John Doe",
		"email": "john.doe@example.com", // This will be redacted
		"phone": "+1234567890",          // This will be redacted
	}

	// Schema annotations indicate which fields are PII
	schemaAnnotations := map[string]lifecycle.FieldAnnotations{
		"email": {Encrypted: true, Redactable: true},
		"phone": {Encrypted: true, Redactable: true},
	}

	producer.EmitResourceCreated(ctx, correlationID, actor, resource, resourceData, schemaAnnotations)

	// Service shutdown
	fmt.Println("\n=== Service Shutdown ===")
	producer.EmitServiceShutdown(ctx, "graceful", 0)
}

