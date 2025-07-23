package domain

import (
	"context"
	"url-db/internal/application/dto/response"
	"url-db/internal/domain/repository"
)

// ListDomainsUseCase handles the listing of domains
type ListDomainsUseCase struct {
	domainRepo repository.DomainRepository
}

// NewListDomainsUseCase creates a new instance of ListDomainsUseCase
func NewListDomainsUseCase(repo repository.DomainRepository) *ListDomainsUseCase {
	return &ListDomainsUseCase{domainRepo: repo}
}

// Execute performs the domain listing use case
func (uc *ListDomainsUseCase) Execute(ctx context.Context, page, size int) (*response.DomainListResponse, error) {
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

	// Get domains from repository
	domains, totalCount, err := uc.domainRepo.List(ctx, page, size)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	domainResponses := make([]response.DomainResponse, len(domains))
	for i, domain := range domains {
		domainResponses[i] = response.DomainResponse{
			Name:        domain.Name(),
			Description: domain.Description(),
			CreatedAt:   domain.CreatedAt(),
			UpdatedAt:   domain.UpdatedAt(),
		}
	}

	// Calculate total pages
	totalPages := (totalCount + size - 1) / size

	return &response.DomainListResponse{
		Domains:    domainResponses,
		TotalCount: totalCount,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}, nil
}
