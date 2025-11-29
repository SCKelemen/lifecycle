package lifecycle

// ColorLoader provides utilities to load colors from API generator type definitions
// This allows services to automatically use colors from their type/event annotations

// LoadColorsFromTypeDefinitions loads colors from API generator type definitions
// This would typically be called at service startup with the type definitions
//
// Example usage:
//
//	colors := LoadColorsFromTypeDefinitions(typeFiles)
//	registry := NewColorRegistry()
//	for api, color := range colors.APIs {
//	    registry.RegisterAPIColor(api, color)
//	}
//	for event, color := range colors.Events {
//	    registry.RegisterEventColor(event, color)
//	}
type ColorDefinitions struct {
	APIs     map[string]string // API type -> color (e.g., "examples.User" -> "#3B82F6")
	Events   map[string]string // Event type -> color (e.g., "examples.OrderCreated" -> "#10B981")
	Services map[string]string // Service name -> color (optional, can be set via config)
}

// LoadColorsFromTypeDefinitions extracts colors from type definitions
// This function would be implemented by integrating with the API generator's schema loader
// For now, this is a placeholder that shows the expected interface
//
// In practice, this would:
// 1. Load type files using the API generator's schema loader
// 2. Extract color annotations from TypeSpec.Annotations
// 3. Map type names to colors
// 4. Return a ColorDefinitions struct
func LoadColorsFromTypeDefinitions(typeFiles interface{}) *ColorDefinitions {
	// This is a placeholder - actual implementation would:
	// 1. Iterate through typeFiles
	// 2. Extract color from annotations: typeFile.Spec.Annotations (look for color annotation)
	// 3. Map type name to color: typeFile.Spec.Type -> color
	// 4. Determine if it's an API (Kind: "Type") or Event (Kind: "Event")

	return &ColorDefinitions{
		APIs:     make(map[string]string),
		Events:   make(map[string]string),
		Services: make(map[string]string),
	}
}

// ExtractColorFromAnnotations extracts color value from annotations
// This matches the logic from the API generator's CLI
func ExtractColorFromAnnotations(annotations interface{}) string {
	// This would need to match the annotation structure from the API generator
	// For now, this is a placeholder

	// Expected structure:
	// annotations: []Annotation
	// Annotation.Color can be:
	//   - string: "#RRGGBB"
	//   - map[string]interface{}: {"value": "#RRGGBB"} or {"color": "#RRGGBB"}

	return ""
}
