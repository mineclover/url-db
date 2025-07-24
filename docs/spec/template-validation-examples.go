// Package examples demonstrates template validation using santhosh-tekuri/jsonschema
package examples

import (
    "encoding/json"
    "fmt"
    "log"
    "strings"
    
    "github.com/santhosh-tekuri/jsonschema/v6"
)

// TemplateValidator demonstrates template validation patterns
type TemplateValidator struct {
    compiler *jsonschema.Compiler
    schemas  map[string]*jsonschema.Schema
}

// NewTemplateValidator creates a new template validator with embedded schemas
func NewTemplateValidator() (*TemplateValidator, error) {
    compiler := jsonschema.NewCompiler()
    
    // Register custom formats
    compiler.RegisterExtension("custom-formats", map[string]jsonschema.Format{
        "semantic-version": {
            Validate: func(value string) error {
                if !isValidSemanticVersion(value) {
                    return fmt.Errorf("invalid semantic version: %s", value)
                }
                return nil
            },
        },
        "template-name": {
            Validate: func(value string) error {
                if !isValidTemplateName(value) {
                    return fmt.Errorf("invalid template name: %s", value)
                }
                return nil
            },
        },
    })
    
    validator := &TemplateValidator{
        compiler: compiler,
        schemas:  make(map[string]*jsonschema.Schema),
    }
    
    // Load embedded schemas
    if err := validator.loadSchemas(); err != nil {
        return nil, fmt.Errorf("failed to load schemas: %w", err)
    }
    
    return validator, nil
}

// loadSchemas loads all template schemas
func (tv *TemplateValidator) loadSchemas() error {
    schemas := map[string]string{
        "base": baseSchema,
        "layout": layoutSchema,
        "form": formSchema,
        "document": documentSchema,
        "custom": customSchema,
    }
    
    for name, schemaContent := range schemas {
        schemaURL := fmt.Sprintf("https://url-db.internal/schemas/%s.json", name)
        
        if err := tv.compiler.AddResource(schemaURL, schemaContent); err != nil {
            return fmt.Errorf("failed to add schema %s: %w", name, err)
        }
        
        schema, err := tv.compiler.Compile(schemaURL)
        if err != nil {
            return fmt.Errorf("failed to compile schema %s: %w", name, err)
        }
        
        tv.schemas[name] = schema
    }
    
    return nil
}

// ValidateTemplate validates template data against the appropriate schema
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
    
    // Extract template type
    templateType := "base"
    if dataMap, ok := data.(map[string]interface{}); ok {
        if t, exists := dataMap["type"]; exists {
            if typeStr, ok := t.(string); ok {
                templateType = typeStr
            }
        }
    }
    
    // Validate with base schema first
    baseSchema := tv.schemas["base"]
    if err := baseSchema.Validate(data); err != nil {
        return tv.convertError(err), nil
    }
    
    // Validate with specific type schema if available
    if typeSchema, exists := tv.schemas[templateType]; exists {
        if err := typeSchema.Validate(data); err != nil {
            return tv.convertError(err), nil
        }
    }
    
    return &ValidationResult{Valid: true}, nil
}

// convertError converts jsonschema validation errors to our format
func (tv *TemplateValidator) convertError(err error) *ValidationResult {
    var errors []ValidationError
    
    if ve, ok := err.(*jsonschema.ValidationError); ok {
        errors = tv.extractErrors(ve)
    } else {
        errors = []ValidationError{{
            Path:    "$",
            Message: err.Error(),
        }}
    }
    
    return &ValidationResult{
        Valid:  false,
        Errors: errors,
    }
}

// extractErrors recursively extracts validation errors
func (tv *TemplateValidator) extractErrors(ve *jsonschema.ValidationError) []ValidationError {
    var errors []ValidationError
    
    error := ValidationError{
        Path:    ve.InstanceLocation,
        Message: ve.Message,
    }
    
    if ve.InstanceValue != nil {
        error.Value = ve.InstanceValue
    }
    
    errors = append(errors, error)
    
    // Process nested errors
    for _, cause := range ve.Causes {
        if subVe, ok := cause.(*jsonschema.ValidationError); ok {
            errors = append(errors, tv.extractErrors(subVe)...)
        }
    }
    
    return errors
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

// Example usage functions

// ExampleValidateLayoutTemplate demonstrates layout template validation
func ExampleValidateLayoutTemplate() {
    validator, err := NewTemplateValidator()
    if err != nil {
        log.Fatal(err)
    }
    
    // Valid layout template
    validLayout := `{
        "version": "1.0",
        "type": "layout",
        "metadata": {
            "name": "two-column-layout",
            "description": "Responsive two-column layout"
        },
        "content": {
            "structure": {
                "type": "grid",
                "container": {
                    "maxWidth": "1200px",
                    "padding": "20px"
                },
                "areas": [
                    {
                        "name": "sidebar",
                        "width": "300px",
                        "position": "left"
                    },
                    {
                        "name": "main",
                        "width": "auto",
                        "position": "right"
                    }
                ]
            },
            "components": {
                "sidebar": {
                    "type": "navigation"
                },
                "main": {
                    "type": "content"
                }
            }
        },
        "presentation": {
            "theme": "light",
            "responsive": true
        }
    }`
    
    result, err := validator.ValidateTemplate(validLayout)
    if err != nil {
        log.Fatal(err)
    }
    
    if result.Valid {
        fmt.Println("✓ Layout template is valid")
    } else {
        fmt.Println("✗ Layout template validation failed:")
        for _, e := range result.Errors {
            fmt.Printf("  - %s: %s\n", e.Path, e.Message)
        }
    }
}

// ExampleValidateInvalidTemplate demonstrates error handling
func ExampleValidateInvalidTemplate() {
    validator, err := NewTemplateValidator()
    if err != nil {
        log.Fatal(err)
    }
    
    // Invalid template - missing required fields
    invalidTemplate := `{
        "version": "invalid-version",
        "type": "unknown-type",
        "metadata": {
            "name": ""
        }
    }`
    
    result, err := validator.ValidateTemplate(invalidTemplate)
    if err != nil {
        log.Fatal(err)
    }
    
    if !result.Valid {
        fmt.Println("Validation errors found:")
        for _, e := range result.Errors {
            fmt.Printf("  Path: %s\n", e.Path)
            fmt.Printf("  Message: %s\n", e.Message)
            if e.Value != nil {
                fmt.Printf("  Value: %v\n", e.Value)
            }
            fmt.Println()
        }
    }
}

// ExampleValidateFormTemplate demonstrates form template validation
func ExampleValidateFormTemplate() {
    validator, err := NewTemplateValidator()
    if err != nil {
        log.Fatal(err)
    }
    
    formTemplate := `{
        "version": "1.0",
        "type": "form",
        "metadata": {
            "name": "contact-form",
            "description": "Simple contact form"
        },
        "schema": {
            "fields": [
                {
                    "name": "name",
                    "type": "text",
                    "label": "Full Name",
                    "required": true,
                    "placeholder": "Enter your name"
                },
                {
                    "name": "email",
                    "type": "text",
                    "label": "Email Address",
                    "required": true,
                    "placeholder": "your@email.com"
                },
                {
                    "name": "message",
                    "type": "textarea",
                    "label": "Message",
                    "required": true,
                    "placeholder": "Your message here..."
                }
            ],
            "sections": [
                {
                    "title": "Contact Information",
                    "fields": ["name", "email"]
                },
                {
                    "title": "Message",
                    "fields": ["message"]
                }
            ]
        },
        "validation": {
            "rules": {
                "name": {
                    "type": "string",
                    "required": true,
                    "min": 2,
                    "max": 100
                },
                "email": {
                    "type": "string",
                    "required": true,
                    "pattern": "^[^@]+@[^@]+\\.[^@]+$"
                }
            },
            "messages": {
                "name.required": "Name is required",
                "email.pattern": "Please enter a valid email address"
            }
        },
        "presentation": {
            "layout": "vertical",
            "theme": "light"
        }
    }`
    
    result, err := validator.ValidateTemplate(formTemplate)
    if err != nil {
        log.Fatal(err)
    }
    
    if result.Valid {
        fmt.Println("✓ Form template is valid")
    } else {
        fmt.Println("✗ Form template validation failed:")
        for _, e := range result.Errors {
            fmt.Printf("  - %s: %s\n", e.Path, e.Message)
        }
    }
}

// ExampleBatchValidation demonstrates validating multiple templates
func ExampleBatchValidation() {
    validator, err := NewTemplateValidator()
    if err != nil {
        log.Fatal(err)
    }
    
    templates := map[string]string{
        "layout-1": `{"version": "1.0", "type": "layout", "content": {"structure": {"type": "grid", "areas": []}}}`,
        "form-1":   `{"version": "1.0", "type": "form", "schema": {"fields": []}}`,
        "invalid":  `{"version": "invalid", "type": "unknown"}`,
    }
    
    results := make(map[string]*ValidationResult)
    
    for name, template := range templates {
        result, err := validator.ValidateTemplate(template)
        if err != nil {
            log.Printf("Error validating %s: %v", name, err)
            continue
        }
        results[name] = result
    }
    
    // Print results
    for name, result := range results {
        if result.Valid {
            fmt.Printf("✓ %s: Valid\n", name)
        } else {
            fmt.Printf("✗ %s: Invalid (%d errors)\n", name, len(result.Errors))
            for _, e := range result.Errors {
                fmt.Printf("    %s: %s\n", e.Path, e.Message)
            }
        }
    }
}

// Helper functions for custom format validation
func isValidSemanticVersion(version string) bool {
    // Simple semantic version validation
    parts := strings.Split(version, ".")
    if len(parts) < 2 || len(parts) > 3 {
        return false
    }
    
    // Check if each part is numeric
    for _, part := range parts {
        if part == "" {
            return false
        }
        for _, r := range part {
            if r < '0' || r > '9' {
                // Allow pre-release identifiers after dash
                if r == '-' && part == parts[len(parts)-1] {
                    continue
                }
                return false
            }
        }
    }
    
    return true
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

// Schema definitions (in a real implementation, these would be loaded from files)
const baseSchema = `{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "$id": "https://url-db.internal/schemas/base.json",
    "title": "Base Template Schema",
    "type": "object",
    "required": ["version", "type"],
    "properties": {
        "version": {
            "type": "string",
            "format": "semantic-version"
        },
        "type": {
            "type": "string",
            "enum": ["layout", "form", "document", "custom"]
        },
        "metadata": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "format": "template-name",
                    "minLength": 1
                },
                "description": {
                    "type": "string"
                }
            }
        }
    }
}`

const layoutSchema = `{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "$id": "https://url-db.internal/schemas/layout.json",
    "title": "Layout Template Schema",
    "allOf": [
        {"$ref": "base.json"},
        {
            "if": {
                "properties": {"type": {"const": "layout"}}
            },
            "then": {
                "required": ["content"],
                "properties": {
                    "content": {
                        "type": "object",
                        "required": ["structure"],
                        "properties": {
                            "structure": {
                                "type": "object",
                                "required": ["type"],
                                "properties": {
                                    "type": {
                                        "type": "string",
                                        "enum": ["grid", "flex", "card", "list"]
                                    },
                                    "areas": {
                                        "type": "array",
                                        "items": {
                                            "type": "object",
                                            "required": ["name"],
                                            "properties": {
                                                "name": {"type": "string", "minLength": 1},
                                                "width": {"type": "string"},
                                                "position": {"type": "string"}
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    ]
}`

const formSchema = `{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "$id": "https://url-db.internal/schemas/form.json",
    "title": "Form Template Schema",
    "allOf": [
        {"$ref": "base.json"},
        {
            "if": {
                "properties": {"type": {"const": "form"}}
            },
            "then": {
                "required": ["schema"],
                "properties": {
                    "schema": {
                        "type": "object",
                        "required": ["fields"],
                        "properties": {
                            "fields": {
                                "type": "array",
                                "items": {
                                    "type": "object",
                                    "required": ["name", "type", "label"],
                                    "properties": {
                                        "name": {"type": "string", "minLength": 1},
                                        "type": {"type": "string", "enum": ["text", "number", "select", "textarea"]},
                                        "label": {"type": "string", "minLength": 1}
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    ]
}`

const documentSchema = `{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "$id": "https://url-db.internal/schemas/document.json",
    "title": "Document Template Schema",
    "allOf": [{"$ref": "base.json"}]
}`

const customSchema = `{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "$id": "https://url-db.internal/schemas/custom.json",
    "title": "Custom Template Schema",
    "allOf": [{"$ref": "base.json"}]
}`