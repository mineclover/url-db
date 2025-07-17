package attributes

import (
	"context"
	"fmt"
	"time"

	"github.com/url-db/internal/models"
)

// AttributeService defines the interface for attribute business logic
type AttributeService interface {
	CreateAttribute(ctx context.Context, domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error)
	GetAttribute(ctx context.Context, id int) (*models.Attribute, error)
	ListAttributes(ctx context.Context, domainID int) (*models.AttributeListResponse, error)
	UpdateAttribute(ctx context.Context, id int, req *models.UpdateAttributeRequest) (*models.Attribute, error)
	DeleteAttribute(ctx context.Context, id int) error
}

// DomainService interface for domain validation
type DomainService interface {
	GetDomain(ctx context.Context, id int) (*models.Domain, error)
}

// attributeService implements AttributeService
type attributeService struct {
	repo          AttributeRepository
	domainService DomainService
}

// NewAttributeService creates a new attribute service
func NewAttributeService(repo AttributeRepository, domainService DomainService) AttributeService {
	return &attributeService{
		repo:          repo,
		domainService: domainService,
	}
}

// validateAttributeName validates attribute name
func (s *attributeService) validateAttributeName(name string) error {
	if len(name) == 0 {
		return ErrAttributeNameRequired
	}
	if len(name) > 255 {
		return ErrAttributeNameTooLong
	}
	return nil
}

// validateAttributeType validates attribute type
func (s *attributeService) validateAttributeType(attrType AttributeType) error {
	if attrType == "" {
		return ErrAttributeTypeRequired
	}
	if !IsValidAttributeType(attrType) {
		return ErrAttributeTypeInvalid
	}
	return nil
}

// validateDescription validates description
func (s *attributeService) validateDescription(description string) error {
	if len(description) > 1000 {
		return ErrDescriptionTooLong
	}
	return nil
}

// CreateAttribute creates a new attribute
func (s *attributeService) CreateAttribute(ctx context.Context, domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error) {
	// Validate input
	if err := s.validateAttributeName(req.Name); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	if err := s.validateAttributeType(req.Type); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	if err := s.validateDescription(req.Description); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if domain exists
	_, err := s.domainService.GetDomain(ctx, domainID)
	if err != nil {
		return nil, ErrDomainNotFound
	}

	// Check if attribute name already exists in domain
	existing, err := s.repo.GetByDomainIDAndName(ctx, domainID, req.Name)
	if err != nil && err != ErrAttributeNotFound {
		return nil, fmt.Errorf("failed to check attribute existence: %w", err)
	}
	if existing != nil {
		return nil, ErrAttributeAlreadyExists
	}

	// Create attribute
	attribute := &models.Attribute{
		DomainID:    domainID,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	err = s.repo.Create(ctx, attribute)
	if err != nil {
		return nil, fmt.Errorf("failed to create attribute: %w", err)
	}

	return attribute, nil
}

// GetAttribute retrieves an attribute by ID
func (s *attributeService) GetAttribute(ctx context.Context, id int) (*models.Attribute, error) {
	attribute, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return attribute, nil
}

// ListAttributes lists all attributes for a domain
func (s *attributeService) ListAttributes(ctx context.Context, domainID int) (*models.AttributeListResponse, error) {
	// Check if domain exists
	_, err := s.domainService.GetDomain(ctx, domainID)
	if err != nil {
		return nil, ErrDomainNotFound
	}

	attributes, err := s.repo.GetByDomainID(ctx, domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to list attributes: %w", err)
	}

	// Convert to response format
	attributeList := make([]models.Attribute, len(attributes))
	for i, attr := range attributes {
		attributeList[i] = *attr
	}

	return &models.AttributeListResponse{
		Attributes: attributeList,
	}, nil
}

// UpdateAttribute updates an attribute
func (s *attributeService) UpdateAttribute(ctx context.Context, id int, req *models.UpdateAttributeRequest) (*models.Attribute, error) {
	// Validate input
	if err := s.validateDescription(req.Description); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get existing attribute
	attribute, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update only description (name and type are immutable)
	attribute.Description = req.Description

	err = s.repo.Update(ctx, attribute)
	if err != nil {
		return nil, fmt.Errorf("failed to update attribute: %w", err)
	}

	return attribute, nil
}

// DeleteAttribute deletes an attribute
func (s *attributeService) DeleteAttribute(ctx context.Context, id int) error {
	// Check if attribute exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if attribute has values
	hasValues, err := s.repo.HasValues(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check attribute values: %w", err)
	}
	if hasValues {
		return ErrAttributeHasValues
	}

	// Delete attribute
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete attribute: %w", err)
	}

	return nil
}