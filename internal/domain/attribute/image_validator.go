package attribute

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// ImageValidator implements validation for image attribute type
type ImageValidator struct{}

// NewImageValidator creates a new image validator
func NewImageValidator() *ImageValidator {
	return &ImageValidator{}
}

// Validate validates an image attribute value
func (v *ImageValidator) Validate(value string, orderIndex *int) ValidationResult {
	// order_index should not be used for image type
	if orderIndex != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "order_index not allowed for image type",
		}
	}
	
	// Check if it's a data URL or HTTP(S) URL
	if strings.HasPrefix(value, "data:image/") {
		return v.validateDataURL(value)
	} else if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return v.validateHTTPURL(value)
	}
	
	return ValidationResult{
		IsValid:      false,
		ErrorCode:    "validation_error",
		ErrorMessage: "image must be either data URL (data:image/...) or HTTP(S) URL",
	}
}

// validateDataURL validates a data URL format
func (v *ImageValidator) validateDataURL(value string) ValidationResult {
	// Parse data URL format: data:image/{type};base64,{data}
	if !strings.Contains(value, ";base64,") {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "data URL must use base64 encoding",
		}
	}
	
	parts := strings.SplitN(value, ";base64,", 2)
	if len(parts) != 2 {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "invalid data URL format",
		}
	}
	
	// Validate MIME type
	mimeType := parts[0]
	supportedTypes := []string{
		"data:image/jpeg",
		"data:image/png", 
		"data:image/gif",
		"data:image/webp",
	}
	
	isSupported := false
	for _, supportedType := range supportedTypes {
		if mimeType == supportedType {
			isSupported = true
			break
		}
	}
	
	if !isSupported {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: fmt.Sprintf("unsupported image type: %s. Supported types: jpeg, png, gif, webp", 
				strings.TrimPrefix(mimeType, "data:image/")),
		}
	}
	
	// Validate base64 data
	base64Data := parts[1]
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "invalid base64 encoding",
		}
	}
	
	// Check size limit (10MB)
	const maxSize = 10 * 1024 * 1024 // 10MB
	if len(decodedData) > maxSize {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: fmt.Sprintf("image size exceeds maximum limit of 10MB (actual: %.2fMB)", 
				float64(len(decodedData))/1024/1024),
		}
	}
	
	return ValidationResult{
		IsValid:         true,
		NormalizedValue: value, // Keep data URL as-is
	}
}

// validateHTTPURL validates an HTTP(S) URL
func (v *ImageValidator) validateHTTPURL(value string) ValidationResult {
	// Parse URL
	parsedURL, err := url.Parse(value)
	if err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "invalid URL format",
		}
	}
	
	// Ensure it's HTTP or HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "URL must use http or https scheme",
		}
	}
	
	// Ensure host is present
	if parsedURL.Host == "" {
		return ValidationResult{
			IsValid:      false,
			ErrorCode:    "validation_error",
			ErrorMessage: "URL must have a valid host",
		}
	}
	
	return ValidationResult{
		IsValid:         true,
		NormalizedValue: value, // Keep URL as-is
	}
}

// GetType returns the attribute type
func (v *ImageValidator) GetType() AttributeType {
	return TypeImage
}

// GetDescription returns the description of the attribute type
func (v *ImageValidator) GetDescription() string {
	return "이미지 데이터. Base64 또는 URL 형식. 최대 10MB."
}