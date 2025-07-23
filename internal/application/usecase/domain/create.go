package domain

import (
	"context"
	"errors"
	"url-db/internal/application/dto/request"
	"url-db/internal/application/dto/response"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
)

// CreateDomainUseCase handles the creation of a new domain
type CreateDomainUseCase struct {
	domainRepo repository.DomainRepository
}

// NewCreateDomainUseCase creates a new instance of CreateDomainUseCase
func NewCreateDomainUseCase(repo repository.DomainRepository) *CreateDomainUseCase {
	return &CreateDomainUseCase{domainRepo: repo}
}

// Execute performs the domain creation use case
func (uc *CreateDomainUseCase) Execute(ctx context.Context, req *request.CreateDomainRequest) (*response.DomainResponse, error) {
	// Create domain entity
	domain, err := entity.NewDomain(req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	// Check if domain already exists
	exists, err := uc.domainRepo.Exists(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("domain already exists")
	}

	// Save to repository
	if err := uc.domainRepo.Create(ctx, domain); err != nil {
		return nil, err
	}

	// Convert to response
	return &response.DomainResponse{
		Name:        domain.Name(),
		Description: domain.Description(),
		CreatedAt:   domain.CreatedAt(),
		UpdatedAt:   domain.UpdatedAt(),
	}, nil
}
