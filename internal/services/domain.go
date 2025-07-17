package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/url-db/internal/models"
)

type DomainRepository interface {
	Create(ctx context.Context, domain *models.Domain) error
	GetByID(ctx context.Context, id int) (*models.Domain, error)
	GetByName(ctx context.Context, name string) (*models.Domain, error)
	List(ctx context.Context, page, size int) ([]*models.Domain, int, error)
	Update(ctx context.Context, domain *models.Domain) error
	Delete(ctx context.Context, id int) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type domainService struct {
	domainRepo DomainRepository
	logger     *log.Logger
}

func NewDomainService(domainRepo DomainRepository, logger *log.Logger) DomainService {
	return &domainService{
		domainRepo: domainRepo,
		logger:     logger,
	}
}

func (s *domainService) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*models.Domain, error) {
	req.Name = normalizeString(req.Name)
	req.Description = normalizeString(req.Description)

	if err := validateDomainName(req.Name); err != nil {
		return nil, err
	}
	
	if err := validateDescription(req.Description); err != nil {
		return nil, err
	}
	
	exists, err := s.domainRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		s.logger.Printf("Failed to check domain existence: %v", err)
		return nil, err
	}
	if exists {
		return nil, NewDomainAlreadyExistsError(req.Name)
	}
	
	domain := &models.Domain{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	if err := s.domainRepo.Create(ctx, domain); err != nil {
		s.logger.Printf("Failed to create domain: %v", err)
		return nil, err
	}
	
	s.logger.Printf("Created domain: %s (ID: %d)", domain.Name, domain.ID)
	return domain, nil
}

func (s *domainService) GetDomain(ctx context.Context, id int) (*models.Domain, error) {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return nil, err
	}

	domain, err := s.domainRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewDomainNotFoundError(id)
		}
		s.logger.Printf("Failed to get domain: %v", err)
		return nil, err
	}
	
	return domain, nil
}

func (s *domainService) GetDomainByName(ctx context.Context, name string) (*models.Domain, error) {
	name = normalizeString(name)
	
	if err := validateDomainName(name); err != nil {
		return nil, err
	}

	domain, err := s.domainRepo.GetByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewDomainNotFoundError(0)
		}
		s.logger.Printf("Failed to get domain by name: %v", err)
		return nil, err
	}
	
	return domain, nil
}

func (s *domainService) ListDomains(ctx context.Context, page, size int) (*models.DomainListResponse, error) {
	page, size, err := validatePaginationParams(page, size)
	if err != nil {
		return nil, err
	}
	
	domains, totalCount, err := s.domainRepo.List(ctx, page, size)
	if err != nil {
		s.logger.Printf("Failed to list domains: %v", err)
		return nil, err
	}
	
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
	if err := validatePositiveInteger(id, "id"); err != nil {
		return nil, err
	}

	req.Description = normalizeString(req.Description)
	
	if err := validateDescription(req.Description); err != nil {
		return nil, err
	}
	
	domain, err := s.domainRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewDomainNotFoundError(id)
		}
		s.logger.Printf("Failed to get domain: %v", err)
		return nil, err
	}
	
	domain.Description = req.Description
	domain.UpdatedAt = time.Now()
	
	if err := s.domainRepo.Update(ctx, domain); err != nil {
		s.logger.Printf("Failed to update domain: %v", err)
		return nil, err
	}
	
	s.logger.Printf("Updated domain: %s (ID: %d)", domain.Name, domain.ID)
	return domain, nil
}

func (s *domainService) DeleteDomain(ctx context.Context, id int) error {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return err
	}

	err := s.domainRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewDomainNotFoundError(id)
		}
		s.logger.Printf("Failed to delete domain: %v", err)
		return err
	}
	
	s.logger.Printf("Deleted domain with ID: %d", id)
	return nil
}