package attribute

import (
	"fmt"
	"strings"
	"url-db/internal/constants"
)

// StringValidator implements validation for string attribute type
type StringValidator struct{}

// NewStringValidator creates a new string validator
func NewStringValidator() *StringValidator {
	return &StringValidator{}
}

// Validate validates a string attribute value
func (v *StringValidator) Validate(value string, orderIndex *int) ValidationResult {
	// Check length constraint (max 500 characters)
	if err := validateLength(value, constants.MaxStringLength); err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    constants.ValidationErrorCode,
			ErrorMessage: err.Error(),
		}
	}

	// order_index should not be used for string type
	if orderIndex != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    constants.ValidationErrorCode,
			ErrorMessage: fmt.Sprintf(constants.ErrOrderIndexNotAllowed, "string"),
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
