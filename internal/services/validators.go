package services

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	domainNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	urlRegex        = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func validateDomainName(name string) error {
	if len(name) == 0 {
		return NewValidationError("name", "domain name is required")
	}
	if len(name) > 255 {
		return NewValidationError("name", "domain name cannot exceed 255 characters")
	}
	if !domainNameRegex.MatchString(name) {
		return NewValidationError("name", "domain name can only contain alphanumeric characters and hyphens")
	}
	return nil
}

func validateDescription(description string) error {
	if len(description) > 1000 {
		return NewValidationError("description", "description cannot exceed 1000 characters")
	}
	return nil
}

func validateURL(url string) error {
	if len(url) == 0 {
		return NewValidationError("url", "URL is required")
	}
	if len(url) > 2000 {
		return NewValidationError("url", "URL cannot exceed 2000 characters")
	}
	if !urlRegex.MatchString(url) {
		return NewValidationError("url", "invalid URL format")
	}
	return nil
}

func validateTitle(title string) error {
	if len(title) > 500 {
		return NewValidationError("title", "title cannot exceed 500 characters")
	}
	return nil
}

func validateAttributeName(name string) error {
	if len(name) == 0 {
		return NewValidationError("name", "attribute name is required")
	}
	if len(name) > 255 {
		return NewValidationError("name", "attribute name cannot exceed 255 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
		return NewValidationError("name", "attribute name can only contain alphanumeric characters, underscores, and hyphens")
	}
	return nil
}

func validateAttributeType(attributeType string) error {
	validTypes := []string{"tag", "ordered_tag", "number", "string", "markdown", "image"}
	for _, validType := range validTypes {
		if attributeType == validType {
			return nil
		}
	}
	return NewValidationError("type", "invalid attribute type")
}

func validateAttributeValue(attributeType, value string) error {
	if len(value) == 0 {
		return NewValidationError("value", "attribute value is required")
	}

	switch attributeType {
	case "tag", "ordered_tag":
		if len(value) > 100 {
			return NewValidationError("value", "tag value cannot exceed 100 characters")
		}
		if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
			return NewValidationError("value", "tag value can only contain alphanumeric characters, underscores, and hyphens")
		}
	case "number":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return NewValidationError("value", "number value must be a valid number")
		}
	case "string":
		if len(value) > 1000 {
			return NewValidationError("value", "string value cannot exceed 1000 characters")
		}
	case "markdown":
		if len(value) > 10000 {
			return NewValidationError("value", "markdown value cannot exceed 10000 characters")
		}
	case "image":
		if len(value) > 500 {
			return NewValidationError("value", "image URL cannot exceed 500 characters")
		}
		if !urlRegex.MatchString(value) {
			return NewValidationError("value", "image value must be a valid URL")
		}
	}

	return nil
}

func validatePositiveInteger(value int, fieldName string) error {
	if value <= 0 {
		return NewValidationError(fieldName, fieldName+" must be a positive integer")
	}
	return nil
}

func validatePaginationParams(page, size int) (int, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size, nil
}

func normalizeString(s string) string {
	return strings.TrimSpace(s)
}

func generateTitleFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		domain := parts[2]
		if len(parts) > 3 {
			path := parts[len(parts)-1]
			if path != "" {
				return strings.Title(strings.ReplaceAll(path, "-", " "))
			}
		}
		return strings.Title(domain)
	}
	return "Untitled"
}