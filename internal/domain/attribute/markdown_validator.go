package attribute

import (
	"fmt"
	"strings"
	"url-db/internal/constants"
)

// MarkdownValidator implements validation for markdown attribute type
type MarkdownValidator struct{}

// NewMarkdownValidator creates a new markdown validator
func NewMarkdownValidator() *MarkdownValidator {
	return &MarkdownValidator{}
}

// Validate validates a markdown attribute value
func (v *MarkdownValidator) Validate(value string, orderIndex *int) ValidationResult {
	// Check length constraint (max 10000 characters for markdown)
	if err := validateLength(value, constants.MaxMarkdownLength); err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    constants.ValidationErrorCode,
			ErrorMessage: err.Error(),
		}
	}

	// order_index should not be used for markdown type
	if orderIndex != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    constants.ValidationErrorCode,
			ErrorMessage: fmt.Sprintf(constants.ErrOrderIndexNotAllowed, "markdown"),
		}
	}

	// Basic markdown validation - check for balanced brackets/parentheses
	if !v.validateMarkdownSyntax(value) {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    constants.ValidationErrorCode,
			ErrorMessage: constants.ErrInvalidMarkdownSyntax,
		}
	}

	// Trim whitespace but preserve formatting
	normalizedValue := strings.TrimSpace(value)

	return ValidationResult{
		IsValid:         true,
		NormalizedValue: normalizedValue,
	}
}

// validateMarkdownSyntax performs basic markdown syntax validation
func (v *MarkdownValidator) validateMarkdownSyntax(value string) bool {
	// Check for balanced square brackets []
	squareBrackets := 0
	// Check for balanced parentheses ()
	parentheses := 0

	for _, char := range value {
		switch char {
		case '[':
			squareBrackets++
		case ']':
			squareBrackets--
			if squareBrackets < 0 {
				return false
			}
		case '(':
			parentheses++
		case ')':
			parentheses--
			if parentheses < 0 {
				return false
			}
		}
	}

	return squareBrackets == 0 && parentheses == 0
}

// GetType returns the attribute type
func (v *MarkdownValidator) GetType() AttributeType {
	return TypeMarkdown
}

// GetDescription returns the description of the attribute type
func (v *MarkdownValidator) GetDescription() string {
	return "마크다운 형식 텍스트. 최대 10,000자."
}
