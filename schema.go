package lifecycle

// SchemaFieldAnnotations represents field annotations from the API schema system
// This type matches the FieldFlags from github.com/SCKelemen/api/internal/schema
// It's used to integrate with the API generator's annotation system
type SchemaFieldAnnotations struct {
	// PII flags
	PII        bool `json:"pii"`        // Contains personally identifiable information
	Encrypted  bool `json:"encrypted"`  // Field-level encryption required
	Redactable bool `json:"redactable"` // Can be redacted for GDPR Article 17
	Sensitive  bool `json:"sensitive"`  // Sensitive data (general)

	// Other field flags
	Immutable  bool `json:"immutable,omitempty"`
	OutputOnly bool `json:"output_only,omitempty"`
	InputOnly  bool `json:"input_only,omitempty"`
	Required   bool `json:"required,omitempty"`
}

// ConvertFromSchemaFieldFlags converts API schema FieldFlags to lifecycle FieldAnnotations
// This allows the lifecycle library to work with annotations from the API generator
func ConvertFromSchemaFieldFlags(schemaFlags map[string]interface{}) map[string]FieldAnnotations {
	if schemaFlags == nil {
		return nil
	}

	result := make(map[string]FieldAnnotations)
	for fieldName, flags := range schemaFlags {
		if flagsMap, ok := flags.(map[string]interface{}); ok {
			annotations := FieldAnnotations{}
			
			if pii, ok := flagsMap["pii"].(bool); ok {
				annotations.PII = pii
			}
			if encrypted, ok := flagsMap["encrypted"].(bool); ok {
				annotations.Encrypted = encrypted
			}
			if redactable, ok := flagsMap["redactable"].(bool); ok {
				annotations.Redactable = redactable
			}
			if sensitive, ok := flagsMap["sensitive"].(bool); ok {
				annotations.Sensitive = sensitive
			}
			if immutable, ok := flagsMap["immutable"].(bool); ok {
				annotations.Immutable = immutable
			}
			
			result[fieldName] = annotations
		}
	}
	
	return result
}

// ShouldRedact checks if a field should be redacted based on schema annotations
func ShouldRedact(annotations FieldAnnotations) bool {
	return annotations.PII || annotations.Redactable || annotations.Encrypted || annotations.Sensitive
}

// GetPIIFields extracts all PII field names from schema annotations
func GetPIIFields(schemaAnnotations map[string]FieldAnnotations) []string {
	var piiFields []string
	for fieldName, annotations := range schemaAnnotations {
		if ShouldRedact(annotations) {
			piiFields = append(piiFields, fieldName)
		}
	}
	return piiFields
}


