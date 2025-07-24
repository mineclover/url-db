package attribute

import "strings"

// StringValidator implements validation for string attribute type
type StringValidator struct{}

// NewStringValidator creates a new string validator
func NewStringValidator() *StringValidator {
	return &StringValidator{}
}

// Validate validates a string attribute value
func (v *StringValidator) Validate(value string, orderIndex *int) ValidationResult {
	// Check length constraint (max 500 characters)
	if err := validateLength(value, 500); err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: err.Error(),
		}
	}

	// order_index should not be used for string type
	if orderIndex != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "order_index not allowed for string type",
		}
	}

	// Trim whitespace but preserve case
	normalizedValue := strings.TrimSpace(value)

	return ValidationResult{
		IsValid:         true,
		NormalizedValue: normalizedValue,
	}
}

// GetType returns the attribute type
func (v *StringValidator) GetType() AttributeType {
	return TypeString
}

// GetDescription returns the description of the attribute type
func (v *StringValidator) GetDescription() string {
	return "일반 문자열. 최대 500자."
}
