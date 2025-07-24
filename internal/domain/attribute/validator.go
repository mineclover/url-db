package attribute

import (
	"fmt"
	"strconv"
	"strings"
)

// AttributeType represents the type of an attribute
type AttributeType string

const (
	TypeTag        AttributeType = "tag"
	TypeOrderedTag AttributeType = "ordered_tag"
	TypeNumber     AttributeType = "number"
	TypeString     AttributeType = "string"
	TypeMarkdown   AttributeType = "markdown"
	TypeImage      AttributeType = "image"
)

// ValidationResult represents the result of attribute validation
type ValidationResult struct {
	IsValid         bool
	ErrorCode       string
	ErrorMessage    string
	NormalizedValue string
}

// AttributeValidator defines the interface for attribute type validators
type AttributeValidator interface {
	Validate(value string, orderIndex *int) ValidationResult
	GetType() AttributeType
	GetDescription() string
}

// ValidatorRegistry manages all attribute validators
type ValidatorRegistry struct {
	validators map[AttributeType]AttributeValidator
}

// NewValidatorRegistry creates a new validator registry with all built-in validators
func NewValidatorRegistry() *ValidatorRegistry {
	registry := &ValidatorRegistry{
		validators: make(map[AttributeType]AttributeValidator),
	}

	// Register built-in validators
	registry.Register(NewTagValidator())
	registry.Register(NewOrderedTagValidator())
	registry.Register(NewNumberValidator())
	registry.Register(NewStringValidator())
	registry.Register(NewMarkdownValidator())
	registry.Register(NewImageValidator())

	return registry
}

// Register adds a new validator to the registry
func (r *ValidatorRegistry) Register(validator AttributeValidator) {
	r.validators[validator.GetType()] = validator
}

// GetValidator returns the validator for the given type
func (r *ValidatorRegistry) GetValidator(attrType AttributeType) (AttributeValidator, error) {
	validator, exists := r.validators[attrType]
	if !exists {
		return nil, fmt.Errorf("unsupported attribute type: %s", attrType)
	}
	return validator, nil
}

// ValidateAttribute validates an attribute value based on its type
func (r *ValidatorRegistry) ValidateAttribute(attrType AttributeType, value string, orderIndex *int) ValidationResult {
	validator, err := r.GetValidator(attrType)
	if err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "unsupported_type",
			ErrorMessage: err.Error(),
		}
	}

	return validator.Validate(value, orderIndex)
}

// GetSupportedTypes returns all supported attribute types
func (r *ValidatorRegistry) GetSupportedTypes() []AttributeType {
	types := make([]AttributeType, 0, len(r.validators))
	for attrType := range r.validators {
		types = append(types, attrType)
	}
	return types
}

// Common validation helpers

// validateLength checks string length constraints
func validateLength(value string, maxLength int) error {
	if len(value) == 0 {
		return fmt.Errorf("value cannot be empty")
	}
	if len(value) > maxLength {
		return fmt.Errorf("value exceeds maximum length of %d characters", maxLength)
	}
	return nil
}

// validateForbiddenChars checks for forbidden characters
func validateForbiddenChars(value string, forbiddenChars []string) error {
	for _, char := range forbiddenChars {
		if strings.Contains(value, char) {
			return fmt.Errorf("value contains forbidden character: %s", char)
		}
	}
	return nil
}

// normalizeCase converts string to lowercase for case-insensitive attributes
func normalizeCase(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

// validateNumber checks if value is a valid number
func validateNumber(value string) error {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid number format: %s", value)
	}
	return nil
}

