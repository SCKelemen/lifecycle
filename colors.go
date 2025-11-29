package lifecycle

import (
	"github.com/charmbracelet/lipgloss"
)

// ColorRegistry manages color mappings for services, APIs, events, and statuses
// Colors come from type/event annotations in the API generator
type ColorRegistry struct {
	serviceColors map[string]string // service name -> color
	apiColors     map[string]string // API type (e.g., "examples.User") -> color
	eventColors   map[string]string // event type (e.g., "examples.OrderCreated") -> color
	statusColors  map[string]string // status -> color (e.g., "success" -> green, "error" -> red)
}

// NewColorRegistry creates a new color registry
func NewColorRegistry() *ColorRegistry {
	return &ColorRegistry{
		serviceColors: make(map[string]string),
		apiColors:     make(map[string]string),
		eventColors:   make(map[string]string),
		statusColors:  defaultStatusColors(),
	}
}

// defaultStatusColors returns default colors for common statuses
func defaultStatusColors() map[string]string {
	return map[string]string{
		"success":     "#00FF00", // Green
		"error":       "#FF0000", // Red
		"warning":     "#FFA500", // Orange
		"info":        "#00BFFF", // Blue
		"pending":     "#FFFF00", // Yellow
		"in_progress": "#9370DB", // Purple
		"completed":   "#00FF00", // Green
		"failed":      "#FF0000", // Red
		"cancelled":   "#808080", // Gray
		"created":     "#00BFFF", // Blue
		"updated":     "#FFA500", // Orange
		"deleted":     "#FF0000", // Red
	}
}

// RegisterServiceColor registers a color for a service
func (r *ColorRegistry) RegisterServiceColor(service, color string) {
	r.serviceColors[service] = color
}

// RegisterAPIColor registers a color for an API type
func (r *ColorRegistry) RegisterAPIColor(api, color string) {
	r.apiColors[api] = color
}

// RegisterEventColor registers a color for an event type
func (r *ColorRegistry) RegisterEventColor(eventType, color string) {
	r.eventColors[eventType] = color
}

// RegisterStatusColor registers a color for a status
func (r *ColorRegistry) RegisterStatusColor(status, color string) {
	r.statusColors[status] = color
}

// GetServiceColor returns the color for a service, or empty string if not found
func (r *ColorRegistry) GetServiceColor(service string) string {
	return r.serviceColors[service]
}

// GetAPIColor returns the color for an API, or empty string if not found
func (r *ColorRegistry) GetAPIColor(api string) string {
	return r.apiColors[api]
}

// GetEventColor returns the color for an event type, or empty string if not found
func (r *ColorRegistry) GetEventColor(eventType string) string {
	return r.eventColors[eventType]
}

// GetStatusColor returns the color for a status, or default if not found
func (r *ColorRegistry) GetStatusColor(status string) string {
	if color, ok := r.statusColors[status]; ok {
		return color
	}
	// Default to gray for unknown statuses
	return "#808080"
}

// GetColorStyle returns a lipgloss style with the given color
// Handles hex colors (#RRGGBB) and named colors
func GetColorStyle(color string) lipgloss.Style {
	if color == "" {
		return lipgloss.NewStyle()
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

// FormatWithColor formats text with the given color
func FormatWithColor(text, color string) string {
	if color == "" {
		return text
	}
	return GetColorStyle(color).Render(text)
}
