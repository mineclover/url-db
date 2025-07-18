package domains

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"url-db/internal/models"
)

type DomainService interface {
	CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*models.Domain, error)
	GetDomain(ctx context.Context, id int) (*models.Domain, error)
	ListDomains(ctx context.Context, page, size int) (*models.DomainListResponse, error)
	UpdateDomain(ctx context.Context, id int, req *models.UpdateDomainRequest) (*models.Domain, error)
	DeleteDomain(ctx context.Context, id int) error
}

type domainService struct {
	repo DomainRepository
}

func NewDomainService(repo DomainRepository) DomainService {
	return &domainService{repo: repo}
}

var domainNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

func (s *domainService) validateDomainName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("domain name is required")
	}
	if len(name) > 255 {
		return fmt.Errorf("domain name cannot exceed 255 characters")
	}
	if !domainNameRegex.MatchString(name) {
		return fmt.Errorf("domain name can only contain alphanumeric characters and hyphens")
	}
	return nil
}

func (s *domainService) validateDescription(description string) error {
	if len(description) > 1000 {
		return fmt.Errorf("description cannot exceed 1000 characters")
	}
	return nil
}

func (s *domainService) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*models.Domain, error) {
	if err := s.validateDomainName(req.Name); err != nil {
		return nil, fmt.Errorf("validation_error: %w", err)
	}

	if err := s.validateDescription(req.Description); err != nil {
		return nil, fmt.Errorf("validation_error: %w", err)
	}

	// Check if domain already exists
	exists, err := s.repo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check domain existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("conflict: domain with name '%s' already exists", req.Name)
	}

	domain := &models.Domain{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.repo.Create(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	return domain, nil
}

func (s *domainService) GetDomain(ctx context.Context, id int) (*models.Domain, error) {
	domain, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not_found: domain with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return domain, nil
}

func (s *domainService) ListDomains(ctx context.Context, page, size int) (*models.DomainListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	domains, totalCount, err := s.repo.List(ctx, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	// Convert to response format
	domainList := make([]models.Domain, len(domains))
	for i, domain := range domains {
		domainList[i] = *domain
	}

	totalPages := (totalCount + size - 1) / size

	return &models.DomainListResponse{
		Domains:    domainList,
		TotalCount: totalCount,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}, nil
}

func (s *domainService) UpdateDomain(ctx context.Context, id int, req *models.UpdateDomainRequest) (*models.Domain, error) {
	if err := s.validateDescription(req.Description); err != nil {
		return nil, fmt.Errorf("validation_error: %w", err)
	}

	// Check if domain exists
	domain, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not_found: domain with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// Update only description
	domain.Description = req.Description
	domain.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to update domain: %w", err)
	}

	return domain, nil
}

func (s *domainService) DeleteDomain(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("not_found: domain with id %d not found", id)
		}
		return fmt.Errorf("failed to delete domain: %w", err)
	}

	return nil
}
