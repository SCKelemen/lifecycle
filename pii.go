package lifecycle

import (
	"fmt"
	"regexp"
	"strings"
)

// PIIDetector detects PII in data based on field names and patterns
type PIIDetector struct {
	// PII field patterns (field names that indicate PII)
	piiFieldPatterns []*regexp.Regexp
	// PII value patterns (values that match PII patterns)
	piiValuePatterns []*regexp.Regexp
}

// NewPIIDetector creates a new PII detector with default patterns
func NewPIIDetector() *PIIDetector {
	return &PIIDetector{
		piiFieldPatterns: []*regexp.Regexp{
			// Common PII field names
			regexp.MustCompile(`(?i)(email|e-mail)`),
			regexp.MustCompile(`(?i)(phone|telephone|mobile)`),
			regexp.MustCompile(`(?i)(ssn|social.security)`),
			regexp.MustCompile(`(?i)(credit.card|card.number)`),
			regexp.MustCompile(`(?i)(password|passwd|pwd)`),
			regexp.MustCompile(`(?i)(secret|token|key)`),
			regexp.MustCompile(`(?i)(address|street|city|zip|postal)`),
			regexp.MustCompile(`(?i)(name|firstname|lastname|fullname)`),
			regexp.MustCompile(`(?i)(dob|date.of.birth|birthdate)`),
			regexp.MustCompile(`(?i)(ip.address|ip_addr)`),
		},
		piiValuePatterns: []*regexp.Regexp{
			// Email pattern
			regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
			// Phone pattern (E.164 or common formats)
			regexp.MustCompile(`^\+?[1-9]\d{1,14}$`),
			regexp.MustCompile(`^\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}$`),
			// Credit card pattern (basic)
			regexp.MustCompile(`^\d{4}[\s\-]?\d{4}[\s\-]?\d{4}[\s\-]?\d{4}$`),
			// SSN pattern
			regexp.MustCompile(`^\d{3}-?\d{2}-?\d{4}$`),
		},
	}
}

// IsPIIField checks if a field name indicates PII
func (d *PIIDetector) IsPIIField(fieldName string) bool {
	for _, pattern := range d.piiFieldPatterns {
		if pattern.MatchString(fieldName) {
			return true
		}
	}
	return false
}

// IsPIIValue checks if a value matches PII patterns
func (d *PIIDetector) IsPIIValue(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	for _, pattern := range d.piiValuePatterns {
		if pattern.MatchString(str) {
			return true
		}
	}
	return false
}

// Redactor redacts PII from data
type Redactor struct {
	redactionString string
}

// NewRedactor creates a new redactor
func NewRedactor() *Redactor {
	return &Redactor{
		redactionString: "[REDACTED]",
	}
}

// WithRedactionString sets a custom redaction string
func (r *Redactor) WithRedactionString(s string) *Redactor {
	r.redactionString = s
	return r
}

// Redact redacts a value if it's PII
func (r *Redactor) Redact(value interface{}) interface{} {
	if value == nil {
		return value
	}

	// Check if it's a string that looks like PII
	if str, ok := value.(string); ok {
		detector := NewPIIDetector()
		if detector.IsPIIValue(str) {
			return r.redactionString
		}
	}

	return value
}

// RedactMap redacts PII from a map based on field names and values
func (r *Redactor) RedactMap(data map[string]interface{}, detector *PIIDetector) map[string]interface{} {
	if data == nil {
		return nil
	}

	redacted := make(map[string]interface{})
	for key, value := range data {
		// Check if field name indicates PII
		if detector.IsPIIField(key) {
			redacted[key] = r.redactionString
			continue
		}

		// Check if value matches PII patterns
		if detector.IsPIIValue(value) {
			redacted[key] = r.redactionString
			continue
		}

		// Recursively handle nested maps
		if nestedMap, ok := value.(map[string]interface{}); ok {
			redacted[key] = r.RedactMap(nestedMap, detector)
		} else if nestedSlice, ok := value.([]interface{}); ok {
			redacted[key] = r.RedactSlice(nestedSlice, detector)
		} else {
			redacted[key] = value
		}
	}

	return redacted
}

// RedactSlice redacts PII from a slice
func (r *Redactor) RedactSlice(slice []interface{}, detector *PIIDetector) []interface{} {
	if slice == nil {
		return nil
	}

	redacted := make([]interface{}, len(slice))
	for i, value := range slice {
		if detector.IsPIIValue(value) {
			redacted[i] = r.redactionString
		} else if nestedMap, ok := value.(map[string]interface{}); ok {
			redacted[i] = r.RedactMap(nestedMap, detector)
		} else {
			redacted[i] = value
		}
	}

	return redacted
}

// RedactParams redacts PII from query parameters
func (r *Redactor) RedactParams(params []interface{}) []interface{} {
	if params == nil {
		return nil
	}

	detector := NewPIIDetector()
	redacted := make([]interface{}, len(params))
	for i, param := range params {
		if detector.IsPIIValue(param) {
			redacted[i] = r.redactionString
		} else {
			redacted[i] = param
		}
	}

	return redacted
}

// RedactString redacts PII from a string value
func (r *Redactor) RedactString(value string) string {
	detector := NewPIIDetector()
	if detector.IsPIIValue(value) {
		return r.redactionString
	}
	return value
}

// FormatRedacted formats a redacted value for display
func (r *Redactor) FormatRedacted(fieldName string, value interface{}) string {
	detector := NewPIIDetector()
	
	// Check field name
	if detector.IsPIIField(fieldName) {
		return fmt.Sprintf("%s=%s", fieldName, r.redactionString)
	}

	// Check value
	if detector.IsPIIValue(value) {
		return fmt.Sprintf("%s=%s", fieldName, r.redactionString)
	}

	// Return original
	return fmt.Sprintf("%s=%v", fieldName, value)
}

// MaskEmail masks an email address (e.g., "user@example.com" -> "u***@example.com")
func (r *Redactor) MaskEmail(email string) string {
	if email == "" {
		return email
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return r.redactionString
	}

	local := parts[0]
	domain := parts[1]

	if len(local) == 0 {
		return r.redactionString
	}

	// Mask local part (keep first character)
	maskedLocal := string(local[0]) + strings.Repeat("*", len(local)-1)
	return maskedLocal + "@" + domain
}

// MaskPhone masks a phone number (e.g., "+1234567890" -> "+1*******90")
func (r *Redactor) MaskPhone(phone string) string {
	if phone == "" {
		return phone
	}

	// Keep first 2 and last 2 characters
	if len(phone) <= 4 {
		return strings.Repeat("*", len(phone))
	}

	return phone[:2] + strings.Repeat("*", len(phone)-4) + phone[len(phone)-2:]
}


