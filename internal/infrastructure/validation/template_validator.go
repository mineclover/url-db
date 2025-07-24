package validation

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TemplateValidator provides JSON validation for templates
type TemplateValidator struct {
	// TODO: Implement proper JSON Schema validation with santhosh-tekuri/jsonschema
	// For now, we'll use basic JSON validation
}

// NewTemplateValidator creates a new template validator
func NewTemplateValidator() (*TemplateValidator, error) {
	return &TemplateValidator{}, nil
}

// ValidateTemplate validates template data with basic JSON validation
func (tv *TemplateValidator) ValidateTemplate(templateData string) (*ValidationResult, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(templateData), &data); err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Path:    "$",
				Message: fmt.Sprintf("Invalid JSON: %s", err.Error()),
			}},
		}, nil
	}

	// Basic structure validation
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Path:    "$",
				Message: "Template data must be an object",
			}},
		}, nil
	}

	var errors []ValidationError

	// Check required fields
	if version, exists := dataMap["version"]; !exists {
		errors = append(errors, ValidationError{
			Path:    "$.version",
			Message: "Version field is required",
		})
	} else if versionStr, ok := version.(string); !ok {
		errors = append(errors, ValidationError{
			Path:    "$.version",
			Message: "Version must be a string",
		})
	} else if !isValidSemanticVersion(versionStr) {
		errors = append(errors, ValidationError{
			Path:    "$.version",
			Message: "Invalid semantic version format",
			Value:   versionStr,
		})
	}

	if templateType, exists := dataMap["type"]; !exists {
		errors = append(errors, ValidationError{
			Path:    "$.type",
			Message: "Type field is required",
		})
	} else if typeStr, ok := templateType.(string); !ok {
		errors = append(errors, ValidationError{
			Path:    "$.type",
			Message: "Type must be a string",
		})
	} else if !isValidTemplateType(typeStr) {
		errors = append(errors, ValidationError{
			Path:    "$.type",
			Message: "Invalid template type",
			Value:   typeStr,
		})
	}

	if len(errors) > 0 {
		return &ValidationResult{
			Valid:  false,
			Errors: errors,
		}, nil
	}

	return &ValidationResult{Valid: true}, nil
}

// ValidateWithSchema validates data against a specific schema (placeholder)
func (tv *TemplateValidator) ValidateWithSchema(schemaName, data string) (*ValidationResult, error) {
	// For now, delegate to basic validation
	return tv.ValidateTemplate(data)
}

// GenerateTemplate creates a basic template structure for the given type
func (tv *TemplateValidator) GenerateTemplate(templateType string) (map[string]interface{}, error) {
	templates := map[string]map[string]interface{}{
		"layout": {
			"version": "1.0",
			"type":    "layout",
			"metadata": map[string]interface{}{
				"name":        "",
				"description": "",
				"tags":        []string{},
			},
			"content": map[string]interface{}{
				"structure": map[string]interface{}{
					"type":  "grid",
					"areas": []map[string]interface{}{},
				},
			},
			"presentation": map[string]interface{}{
				"theme":      "light",
				"responsive": true,
			},
		},
		"form": {
			"version": "1.0",
			"type":    "form",
			"metadata": map[string]interface{}{
				"name":        "",
				"description": "",
			},
			"schema": map[string]interface{}{
				"fields":   []map[string]interface{}{},
				"sections": []map[string]interface{}{},
			},
			"validation": map[string]interface{}{
				"rules": map[string]interface{}{},
			},
			"presentation": map[string]interface{}{
				"layout": "vertical",
			},
		},
		"document": {
			"version": "1.0",
			"type":    "document",
			"metadata": map[string]interface{}{
				"name": "",
			},
			"schema": map[string]interface{}{
				"sections": []map[string]interface{}{},
			},
		},
		"custom": {
			"version": "1.0",
			"type":    "custom",
			"metadata": map[string]interface{}{
				"name": "",
			},
			"content": map[string]interface{}{},
		},
	}

	template, exists := templates[templateType]
	if !exists {
		return nil, fmt.Errorf("unknown template type: %s", templateType)
	}

	return template, nil
}

// ValidationResult represents the result of template validation
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Path    string      `json:"path"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// Helper functions for validation
func isValidSemanticVersion(version string) bool {
	if version == "" {
		return false
	}

	// Basic semantic version validation: major.minor.patch[-prerelease]
	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return false
	}

	for i, part := range parts {
		if part == "" {
			return false
		}
		
		// For the last part, allow prerelease suffix
		if i == len(parts)-1 && strings.Contains(part, "-") {
			mainPart := strings.Split(part, "-")[0]
			if mainPart == "" {
				return false
			}
			// Check if main part is numeric
			for _, r := range mainPart {
				if r < '0' || r > '9' {
					return false
				}
			}
		} else {
			// Check if part is numeric
			for _, r := range part {
				if r < '0' || r > '9' {
					return false
				}
			}
		}
	}

	return true
}

func isValidTemplateType(templateType string) bool {
	validTypes := []string{"layout", "form", "document", "custom"}
	for _, validType := range validTypes {
		if templateType == validType {
			return true
		}
	}
	return false
}

func isValidTemplateName(name string) bool {
	if name == "" {
		return false
	}

	for _, r := range name {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_') {
			return false
		}
	}

	return true
}