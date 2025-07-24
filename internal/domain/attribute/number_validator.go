package attribute

import (
	"fmt"
	"url-db/internal/constants"
)

// NumberValidator implements validation for number attribute type
type NumberValidator struct{}

// NewNumberValidator creates a new number validator
func NewNumberValidator() *NumberValidator {
	return &NumberValidator{}
}

// Validate validates a number attribute value
func (v *NumberValidator) Validate(value string, orderIndex *int) ValidationResult {
	// Check if it's a valid number
	if err := validateNumber(value); err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    constants.ValidationErrorCode,
			ErrorMessage: err.Error(),
		}
	}

	// order_index should not be used for number type
	if orderIndex != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    constants.ValidationErrorCode,
			ErrorMessage: fmt.Sprintf(constants.ErrOrderIndexNotAllowed, "number"),
		}
	}

	// No normalization needed for numbers, return as-is
	return ValidationResult{
		IsValid:         true,
		NormalizedValue: value,
	}
}

// GetType returns the attribute type
func (v *NumberValidator) GetType() AttributeType {
	return TypeNumber
}

// GetDescription returns the description of the attribute type
func (v *NumberValidator) GetDescription() string {
	return "숫자 값. 정수 또는 실수 허용."
}
