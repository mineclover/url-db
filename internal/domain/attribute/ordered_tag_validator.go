package attribute

// OrderedTagValidator implements validation for ordered_tag attribute type
type OrderedTagValidator struct{}

// NewOrderedTagValidator creates a new ordered tag validator
func NewOrderedTagValidator() *OrderedTagValidator {
	return &OrderedTagValidator{}
}

// Validate validates an ordered tag attribute value
func (v *OrderedTagValidator) Validate(value string, orderIndex *int) ValidationResult {
	// Check length constraint (max 50 characters)
	if err := validateLength(value, 50); err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: err.Error(),
		}
	}
	
	// Check forbidden characters
	forbiddenChars := []string{",", ";", "|", "\n", "\t"}
	if err := validateForbiddenChars(value, forbiddenChars); err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: err.Error(),
		}
	}
	
	// order_index is required for ordered_tag type
	if orderIndex == nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "order_index is required for ordered_tag type",
		}
	}
	
	// order_index must be non-negative
	if *orderIndex < 0 {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "order_index must be non-negative",
		}
	}
	
	// Normalize to lowercase
	normalizedValue := normalizeCase(value)
	
	return ValidationResult{
		IsValid:         true,
		NormalizedValue: normalizedValue,
	}
}

// GetType returns the attribute type
func (v *OrderedTagValidator) GetType() AttributeType {
	return TypeOrderedTag
}

// GetDescription returns the description of the attribute type
func (v *OrderedTagValidator) GetDescription() string {
	return "순서가 있는 태그. order_index 필수."
}