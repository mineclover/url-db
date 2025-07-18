package nodeattributes

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"url-db/internal/models"
)

type Validator interface {
	Validate(attributeType models.AttributeType, value string, orderIndex *int) error
	ValidateValue(attributeType models.AttributeType, value string) error
	ValidateOrderIndex(attributeType models.AttributeType, orderIndex *int) error
}

type validator struct{}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Validate(attributeType models.AttributeType, value string, orderIndex *int) error {
	if err := v.ValidateValue(attributeType, value); err != nil {
		return err
	}

	if err := v.ValidateOrderIndex(attributeType, orderIndex); err != nil {
		return err
	}

	return nil
}

func (v *validator) ValidateValue(attributeType models.AttributeType, value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	switch attributeType {
	case models.AttributeTypeTag:
		return v.validateTag(value)
	case models.AttributeTypeOrderedTag:
		return v.validateOrderedTag(value)
	case models.AttributeTypeNumber:
		return v.validateNumber(value)
	case models.AttributeTypeString:
		return v.validateString(value)
	case models.AttributeTypeMarkdown:
		return v.validateMarkdown(value)
	case models.AttributeTypeImage:
		return v.validateImage(value)
	default:
		return ErrInvalidAttributeType
	}
}

func (v *validator) ValidateOrderIndex(attributeType models.AttributeType, orderIndex *int) error {
	switch attributeType {
	case models.AttributeTypeOrderedTag:
		if orderIndex == nil {
			return ErrOrderIndexRequired
		}
		if *orderIndex < 0 {
			return ErrInvalidOrderIndex
		}
	case models.AttributeTypeTag, models.AttributeTypeNumber, models.AttributeTypeString, models.AttributeTypeMarkdown, models.AttributeTypeImage:
		if orderIndex != nil {
			return ErrOrderIndexNotAllowed
		}
	default:
		return ErrInvalidAttributeType
	}

	return nil
}

func (v *validator) validateTag(value string) error {
	if len(value) > 255 {
		return fmt.Errorf("tag value must be 255 characters or less")
	}

	// Basic tag validation - no special characters except hyphens and underscores
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
		return fmt.Errorf("tag value can only contain letters, numbers, hyphens, and underscores")
	}

	return nil
}

func (v *validator) validateOrderedTag(value string) error {
	if len(value) > 255 {
		return fmt.Errorf("ordered tag value must be 255 characters or less")
	}

	// Same validation as tag
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
		return fmt.Errorf("ordered tag value can only contain letters, numbers, hyphens, and underscores")
	}

	return nil
}

func (v *validator) validateNumber(value string) error {
	// Try to parse as float64 to support both integers and decimals
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("number value must be a valid number")
	}

	return nil
}

func (v *validator) validateString(value string) error {
	if len(value) > 2048 {
		return fmt.Errorf("string value must be 2048 characters or less")
	}

	return nil
}

func (v *validator) validateMarkdown(value string) error {
	if len(value) > 10000 {
		return fmt.Errorf("markdown value must be 10000 characters or less")
	}

	// Basic markdown validation - ensure it's valid text
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("markdown value cannot be empty or only whitespace")
	}

	return nil
}

func (v *validator) validateImage(value string) error {
	// Validate URL format
	parsedURL, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("image value must be a valid URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("image URL must use http or https scheme")
	}

	// Validate image file extensions
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"}
	path := strings.ToLower(parsedURL.Path)

	isValidExtension := false
	for _, ext := range validExtensions {
		if strings.HasSuffix(path, ext) {
			isValidExtension = true
			break
		}
	}

	if !isValidExtension {
		return fmt.Errorf("image URL must have a valid image file extension (.jpg, .jpeg, .png, .gif, .webp, .svg)")
	}

	return nil
}
