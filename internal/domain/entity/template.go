package entity

import (
	"encoding/json"
	"errors"
	"time"
)

// Template represents a template entity in the business domain
type Template struct {
	id           int
	name         string
	domainID     int
	templateData string // JSON string
	title        string
	description  string
	isActive     bool
	createdAt    time.Time
	updatedAt    time.Time
}

// NewTemplate creates a new template entity with validation
func NewTemplate(name, templateData, title, description string, domainID int) (*Template, error) {
	if name == "" {
		return nil, errors.New("template name cannot be empty")
	}

	if len(name) > 255 {
		return nil, errors.New("template name cannot exceed 255 characters")
	}

	if templateData == "" {
		return nil, errors.New("template data cannot be empty")
	}

	// Validate JSON format
	if !isValidJSON(templateData) {
		return nil, errors.New("template data must be valid JSON")
	}

	if domainID <= 0 {
		return nil, errors.New("domain ID must be positive")
	}

	if len(title) > 255 {
		return nil, errors.New("template title cannot exceed 255 characters")
	}

	if len(description) > 1000 {
		return nil, errors.New("template description cannot exceed 1000 characters")
	}

	now := time.Now()
	return &Template{
		name:         name,
		domainID:     domainID,
		templateData: templateData,
		title:        title,
		description:  description,
		isActive:     true, // Default to active
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// Getters - immutable from outside
func (t *Template) ID() int              { return t.id }
func (t *Template) Name() string         { return t.name }
func (t *Template) DomainID() int        { return t.domainID }
func (t *Template) TemplateData() string { return t.templateData }
func (t *Template) Title() string        { return t.title }
func (t *Template) Description() string  { return t.description }
func (t *Template) IsActive() bool       { return t.isActive }
func (t *Template) CreatedAt() time.Time { return t.createdAt }
func (t *Template) UpdatedAt() time.Time { return t.updatedAt }

// Setters for internal use (e.g., by repository)
func (t *Template) SetID(id int) { t.id = id }
func (t *Template) SetTimestamps(createdAt, updatedAt time.Time) {
	t.createdAt = createdAt
	t.updatedAt = updatedAt
}

// Business logic methods
func (t *Template) UpdateTitle(title string) error {
	if len(title) > 255 {
		return errors.New("template title cannot exceed 255 characters")
	}

	t.title = title
	t.updatedAt = time.Now()
	return nil
}

func (t *Template) UpdateDescription(description string) error {
	if len(description) > 1000 {
		return errors.New("template description cannot exceed 1000 characters")
	}

	t.description = description
	t.updatedAt = time.Now()
	return nil
}

func (t *Template) UpdateTemplateData(templateData string) error {
	if templateData == "" {
		return errors.New("template data cannot be empty")
	}

	if !isValidJSON(templateData) {
		return errors.New("template data must be valid JSON")
	}

	t.templateData = templateData
	t.updatedAt = time.Now()
	return nil
}

func (t *Template) Activate() error {
	if t.isActive {
		return errors.New("template is already active")
	}

	t.isActive = true
	t.updatedAt = time.Now()
	return nil
}

func (t *Template) Deactivate() error {
	if !t.isActive {
		return errors.New("template is already inactive")
	}

	t.isActive = false
	t.updatedAt = time.Now()
	return nil
}

func (t *Template) SetActive(active bool) {
	if t.isActive != active {
		t.isActive = active
		t.updatedAt = time.Now()
	}
}

// UpdateContent updates multiple fields at once
func (t *Template) UpdateContent(title, description, templateData string) error {
	if templateData != "" {
		if err := t.UpdateTemplateData(templateData); err != nil {
			return err
		}
	}

	if title != "" {
		if err := t.UpdateTitle(title); err != nil {
			return err
		}
	}

	if description != "" {
		if err := t.UpdateDescription(description); err != nil {
			return err
		}
	}

	return nil
}

// IsValid checks if the template is in a valid state
func (t *Template) IsValid() bool {
	return t.name != "" &&
		len(t.name) <= 255 &&
		t.domainID > 0 &&
		t.templateData != "" &&
		isValidJSON(t.templateData) &&
		len(t.title) <= 255 &&
		len(t.description) <= 1000
}

// GetTemplateType extracts the template type from JSON data
func (t *Template) GetTemplateType() (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(t.templateData), &data); err != nil {
		return "", errors.New("invalid JSON in template data")
	}

	templateType, ok := data["type"].(string)
	if !ok {
		return "", errors.New("template type not found or not a string")
	}

	return templateType, nil
}

// GetTemplateVersion extracts the template version from JSON data
func (t *Template) GetTemplateVersion() (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(t.templateData), &data); err != nil {
		return "", errors.New("invalid JSON in template data")
	}

	version, ok := data["version"].(string)
	if !ok {
		return "", errors.New("template version not found or not a string")
	}

	return version, nil
}

// GetTemplateMetadata extracts metadata from JSON data
func (t *Template) GetTemplateMetadata() (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(t.templateData), &data); err != nil {
		return nil, errors.New("invalid JSON in template data")
	}

	metadata, ok := data["metadata"].(map[string]interface{})
	if !ok {
		return make(map[string]interface{}), nil // Return empty map if no metadata
	}

	return metadata, nil
}

// CanModify checks if the template can be modified (only active templates)
func (t *Template) CanModify() bool {
	return t.isActive
}

// IsCompatibleWithDomain checks if template can be used with given domain
func (t *Template) IsCompatibleWithDomain(domainID int) bool {
	return t.domainID == domainID && t.isActive
}

// Helper function to validate JSON
func isValidJSON(data string) bool {
	var js interface{}
	return json.Unmarshal([]byte(data), &js) == nil
}

// TemplateType represents the type of template
type TemplateType string

const (
	TemplateTypeLayout   TemplateType = "layout"
	TemplateTypeForm     TemplateType = "form"
	TemplateTypeDocument TemplateType = "document"
	TemplateTypeCustom   TemplateType = "custom"
)

// IsValidTemplateType checks if the given string is a valid template type
func IsValidTemplateType(templateType string) bool {
	switch TemplateType(templateType) {
	case TemplateTypeLayout, TemplateTypeForm, TemplateTypeDocument, TemplateTypeCustom:
		return true
	default:
		return false
	}
}

// GetValidTemplateTypes returns all valid template types
func GetValidTemplateTypes() []string {
	return []string{
		string(TemplateTypeLayout),
		string(TemplateTypeForm),
		string(TemplateTypeDocument),
		string(TemplateTypeCustom),
	}
}
