package service

import (
	"context"
	"errors"
	"regexp"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
)

// DomainService represents domain business logic
type DomainService interface {
	CreateDomain(ctx context.Context, name, description string) (*entity.Domain, error)
	ValidateDomainName(name string) error
	ValidateDescription(description string) error
}

type domainService struct {
	domainRepo repository.DomainRepository
}

// NewDomainService creates a new domain service
func NewDomainService(domainRepo repository.DomainRepository) DomainService {
	return &domainService{
		domainRepo: domainRepo,
	}
}

var domainNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

// ValidateDomainName validates domain name according to business rules
func (s *domainService) ValidateDomainName(name string) error {
	if len(name) == 0 {
		return errors.New("domain name is required")
	}
	if len(name) > 255 {
		return errors.New("domain name cannot exceed 255 characters")
	}
	if !domainNameRegex.MatchString(name) {
		return errors.New("domain name can only contain alphanumeric characters and hyphens")
	}
	return nil
}

// ValidateDescription validates domain description according to business rules
func (s *domainService) ValidateDescription(description string) error {
	if len(description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}
	return nil
}

// CreateDomain creates a new domain with business validation
func (s *domainService) CreateDomain(ctx context.Context, name, description string) (*entity.Domain, error) {
	// Validate input
	if err := s.ValidateDomainName(name); err != nil {
		return nil, err
	}

	if err := s.ValidateDescription(description); err != nil {
		return nil, err
	}

	// Check if domain already exists
	exists, err := s.domainRepo.Exists(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("domain already exists")
	}

	// Create domain entity
	domain, err := entity.NewDomain(name, description)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.domainRepo.Create(ctx, domain); err != nil {
		return nil, err
	}

	return domain, nil
}
