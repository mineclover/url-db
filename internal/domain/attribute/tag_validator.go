package attribute

// TagValidator implements validation for tag attribute type
type TagValidator struct{}

// NewTagValidator creates a new tag validator
func NewTagValidator() *TagValidator {
	return &TagValidator{}
}

// Validate validates a tag attribute value
func (v *TagValidator) Validate(value string, orderIndex *int) ValidationResult {
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

	// order_index should not be used for tag type
	if orderIndex != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "order_index not allowed for tag type",
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
func (v *TagValidator) GetType() AttributeType {
	return TypeTag
}

// GetDescription returns the description of the attribute type
func (v *TagValidator) GetDescription() string {
	return "순서 없는 일반 태그. 중복 값 허용하지 않음."
}
