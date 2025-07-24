package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
	"url-db/internal/infrastructure/validation"
)

// TemplateService represents template business logic
type TemplateService interface {
	// Template CRUD operations
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*entity.Template, error)
	GetTemplate(ctx context.Context, id int) (*entity.Template, error)
	GetTemplateByName(ctx context.Context, domainName, name string) (*entity.Template, error)
	UpdateTemplate(ctx context.Context, id int, req *UpdateTemplateRequest) (*entity.Template, error)
	DeleteTemplate(ctx context.Context, id int) error
	
	// Template listing and search
	ListTemplates(ctx context.Context, domainName string, page, size int) ([]*entity.Template, int, error)
	ListActiveTemplates(ctx context.Context, domainName string, page, size int) ([]*entity.Template, int, error)
	ListTemplatesByType(ctx context.Context, domainName, templateType string, page, size int) ([]*entity.Template, int, error)
	SearchTemplates(ctx context.Context, domainName, query string, page, size int) ([]*entity.Template, int, error)
	
	// Template management operations
	ActivateTemplate(ctx context.Context, id int) error
	DeactivateTemplate(ctx context.Context, id int) error
	CloneTemplate(ctx context.Context, sourceID int, newName, newTitle, newDescription string) (*entity.Template, error)
	
	// Template validation and generation
	ValidateTemplateData(templateData string) (*validation.ValidationResult, error)
	GenerateTemplateScaffold(templateType string) (string, error)
	GetValidTemplateTypes() []string
	
	// Template statistics
	GetTemplateStats(ctx context.Context, domainName string) (*repository.TemplateStats, error)
	GetRecentlyModified(ctx context.Context, domainName string, limit int) ([]*entity.Template, error)
	
	// Template utilities
	ExtractTemplateType(templateData string) (string, error)
	ExtractTemplateVersion(templateData string) (string, error)
	ValidateTemplateName(name string) error
}

type templateService struct {
	templateRepo repository.TemplateRepository
	domainRepo   repository.DomainRepository
	validator    *validation.TemplateValidator
}

// NewTemplateService creates a new template service
func NewTemplateService(templateRepo repository.TemplateRepository, domainRepo repository.DomainRepository) (TemplateService, error) {
	validator, err := validation.NewTemplateValidator()
	if err != nil {
		return nil, fmt.Errorf("failed to create template validator: %w", err)
	}

	return &templateService{
		templateRepo: templateRepo,
		domainRepo:   domainRepo,
		validator:    validator,
	}, nil
}

// CreateTemplateRequest represents a request to create a new template
type CreateTemplateRequest struct {
	Name         string `json:"name"`
	DomainName   string `json:"domain_name"`
	TemplateData string `json:"template_data"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

// UpdateTemplateRequest represents a request to update a template
type UpdateTemplateRequest struct {
	TemplateData *string `json:"template_data,omitempty"`
	Title        *string `json:"title,omitempty"`
	Description  *string `json:"description,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

func (s *templateService) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*entity.Template, error) {
	// Validate input
	if err := s.ValidateTemplateName(req.Name); err != nil {
		return nil, fmt.Errorf("invalid template name: %w", err)
	}

	// Validate template data
	result, err := s.ValidateTemplateData(req.TemplateData)
	if err != nil {
		return nil, fmt.Errorf("template validation error: %w", err)
	}

	if !result.Valid {
		return nil, &ValidationError{
			Message: "Template data validation failed",
			Errors:  result.Errors,
		}
	}

	// Get domain
	domain, err := s.domainRepo.GetByName(ctx, req.DomainName)
	if err != nil {
		return nil, fmt.Errorf("domain not found: %w", err)
	}

	// Check if template name already exists in domain
	exists, err := s.templateRepo.Exists(ctx, req.Name, req.DomainName)
	if err != nil {
		return nil, fmt.Errorf("failed to check template existence: %w", err)
	}
	if exists {
		return nil, repository.ErrDuplicateKey
	}

	// Create template entity
	template, err := entity.NewTemplate(
		req.Name,
		req.TemplateData,
		req.Title,
		req.Description,
		domain.ID(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create template entity: %w", err)
	}

	// Save to repository
	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return template, nil
}

func (s *templateService) GetTemplate(ctx context.Context, id int) (*entity.Template, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return template, nil
}

func (s *templateService) GetTemplateByName(ctx context.Context, domainName, name string) (*entity.Template, error) {
	template, err := s.templateRepo.GetByName(ctx, name, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get template by name: %w", err)
	}
	return template, nil
}

func (s *templateService) UpdateTemplate(ctx context.Context, id int, req *UpdateTemplateRequest) (*entity.Template, error) {
	// Get existing template
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	// Check if template can be modified
	if !template.CanModify() {
		return nil, errors.New("inactive templates cannot be modified")
	}

	// Update template data if provided
	if req.TemplateData != nil {
		// Validate new template data
		result, err := s.ValidateTemplateData(*req.TemplateData)
		if err != nil {
			return nil, fmt.Errorf("template validation error: %w", err)
		}

		if !result.Valid {
			return nil, &ValidationError{
				Message: "Template data validation failed",
				Errors:  result.Errors,
			}
		}

		if err := template.UpdateTemplateData(*req.TemplateData); err != nil {
			return nil, fmt.Errorf("failed to update template data: %w", err)
		}
	}

	// Update other fields
	if req.Title != nil {
		if err := template.UpdateTitle(*req.Title); err != nil {
			return nil, fmt.Errorf("failed to update title: %w", err)
		}
	}

	if req.Description != nil {
		if err := template.UpdateDescription(*req.Description); err != nil {
			return nil, fmt.Errorf("failed to update description: %w", err)
		}
	}

	if req.IsActive != nil {
		template.SetActive(*req.IsActive)
	}

	// Save changes
	if err := s.templateRepo.Update(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	return template, nil
}

func (s *templateService) DeleteTemplate(ctx context.Context, id int) error {
	// Check if template exists
	_, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	// Delete template
	if err := s.templateRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	return nil
}

func (s *templateService) ListTemplates(ctx context.Context, domainName string, page, size int) ([]*entity.Template, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20 // Default page size
	}

	templates, total, err := s.templateRepo.List(ctx, domainName, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list templates: %w", err)
	}

	return templates, total, nil
}

func (s *templateService) ListActiveTemplates(ctx context.Context, domainName string, page, size int) ([]*entity.Template, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	templates, total, err := s.templateRepo.ListActive(ctx, domainName, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list active templates: %w", err)
	}

	return templates, total, nil
}

func (s *templateService) ListTemplatesByType(ctx context.Context, domainName, templateType string, page, size int) ([]*entity.Template, int, error) {
	if !entity.IsValidTemplateType(templateType) {
		return nil, 0, fmt.Errorf("invalid template type: %s", templateType)
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	templates, total, err := s.templateRepo.ListByType(ctx, domainName, templateType, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list templates by type: %w", err)
	}

	return templates, total, nil
}

func (s *templateService) SearchTemplates(ctx context.Context, domainName, query string, page, size int) ([]*entity.Template, int, error) {
	if query == "" {
		return s.ListTemplates(ctx, domainName, page, size)
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	templates, total, err := s.templateRepo.Search(ctx, domainName, query, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search templates: %w", err)
	}

	return templates, total, nil
}

func (s *templateService) ActivateTemplate(ctx context.Context, id int) error {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	if err := template.Activate(); err != nil {
		return err
	}

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return fmt.Errorf("failed to activate template: %w", err)
	}

	return nil
}

func (s *templateService) DeactivateTemplate(ctx context.Context, id int) error {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	if err := template.Deactivate(); err != nil {
		return err
	}

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return fmt.Errorf("failed to deactivate template: %w", err)
	}

	return nil
}

func (s *templateService) CloneTemplate(ctx context.Context, sourceID int, newName, newTitle, newDescription string) (*entity.Template, error) {
	// Validate new template name
	if err := s.ValidateTemplateName(newName); err != nil {
		return nil, fmt.Errorf("invalid template name: %w", err)
	}

	// Check if source template exists
	_, err := s.templateRepo.GetByID(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("source template not found: %w", err)
	}

	// Get domain name for existence check
	domain, err := s.templateRepo.GetDomainByTemplateID(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// Check if new name already exists
	exists, err := s.templateRepo.Exists(ctx, newName, domain.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to check template existence: %w", err)
	}
	if exists {
		return nil, repository.ErrDuplicateKey
	}

	// Clone template
	clonedTemplate, err := s.templateRepo.Clone(ctx, sourceID, newName, newTitle, newDescription)
	if err != nil {
		return nil, fmt.Errorf("failed to clone template: %w", err)
	}

	return clonedTemplate, nil
}

func (s *templateService) ValidateTemplateData(templateData string) (*validation.ValidationResult, error) {
	return s.validator.ValidateTemplate(templateData)
}

func (s *templateService) GenerateTemplateScaffold(templateType string) (string, error) {
	if !entity.IsValidTemplateType(templateType) {
		return "", fmt.Errorf("invalid template type: %s", templateType)
	}

	template, err := s.validator.GenerateTemplate(templateType)
	if err != nil {
		return "", fmt.Errorf("failed to generate template: %w", err)
	}

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %w", err)
	}

	return string(data), nil
}

func (s *templateService) GetValidTemplateTypes() []string {
	return entity.GetValidTemplateTypes()
}

func (s *templateService) GetTemplateStats(ctx context.Context, domainName string) (*repository.TemplateStats, error) {
	if statsRepo, ok := s.templateRepo.(repository.TemplateRepositoryStats); ok {
		return statsRepo.GetStats(ctx, domainName)
	}
	
	// Fallback implementation
	templates, total, err := s.templateRepo.List(ctx, domainName, 1, 1000) // Get a large number
	if err != nil {
		return nil, fmt.Errorf("failed to get templates for stats: %w", err)
	}
	
	stats := &repository.TemplateStats{
		TotalCount:   total,
		ActiveCount:  0,
		InactiveCount: 0,
		TypeCounts:   make(map[string]int),
	}
	
	for _, template := range templates {
		if template.IsActive() {
			stats.ActiveCount++
		} else {
			stats.InactiveCount++
		}
		
		templateType, err := template.GetTemplateType()
		if err == nil {
			stats.TypeCounts[templateType]++
		}
	}
	
	return stats, nil
}

func (s *templateService) GetRecentlyModified(ctx context.Context, domainName string, limit int) ([]*entity.Template, error) {
	if limit <= 0 || limit > 100 {
		limit = 10 // Default limit
	}

	templates, err := s.templateRepo.GetRecentlyModified(ctx, domainName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recently modified templates: %w", err)
	}

	return templates, nil
}

func (s *templateService) ExtractTemplateType(templateData string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(templateData), &data); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	templateType, ok := data["type"].(string)
	if !ok {
		return "", errors.New("template type not found or not a string")
	}

	return templateType, nil
}

func (s *templateService) ExtractTemplateVersion(templateData string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(templateData), &data); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	version, ok := data["version"].(string)
	if !ok {
		return "", errors.New("template version not found or not a string")
	}

	return version, nil
}

func (s *templateService) ValidateTemplateName(name string) error {
	if name == "" {
		return errors.New("template name cannot be empty")
	}

	if len(name) > 255 {
		return errors.New("template name cannot exceed 255 characters")
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_') {
			return errors.New("template name can only contain letters, numbers, hyphens, and underscores")
		}
	}

	// Name cannot start or end with hyphen or underscore
	if strings.HasPrefix(name, "-") || strings.HasPrefix(name, "_") ||
		strings.HasSuffix(name, "-") || strings.HasSuffix(name, "_") {
		return errors.New("template name cannot start or end with hyphen or underscore")
	}

	return nil
}

// ValidationError represents a template validation error
type ValidationError struct {
	Message string                         `json:"message"`
	Errors  []validation.ValidationError   `json:"errors"`
}

func (e *ValidationError) Error() string {
	return e.Message
}