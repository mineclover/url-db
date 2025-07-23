package attribute

import (
	"context"
	"errors"
	"url-db/internal/application/dto/response"
	"url-db/internal/constants"
	"url-db/internal/domain/repository"
)

type ListAttributesUseCase struct {
	attributeRepo repository.AttributeRepository
	domainRepo    repository.DomainRepository
}

func NewListAttributesUseCase(attributeRepo repository.AttributeRepository, domainRepo repository.DomainRepository) *ListAttributesUseCase {
	return &ListAttributesUseCase{
		attributeRepo: attributeRepo,
		domainRepo:    domainRepo,
	}
}

func (uc *ListAttributesUseCase) Execute(ctx context.Context, domainID int) (*response.AttributeListResponse, error) {
	// Verify domain exists
	domain, err := uc.domainRepo.GetByID(ctx, domainID)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, errors.New(constants.ErrDomainNotFound)
	}

	// Get attributes from repository
	attributes, err := uc.attributeRepo.ListByDomainID(ctx, domainID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	attributeResponses := make([]response.AttributeResponse, len(attributes))
	for i, attr := range attributes {
		attributeResponses[i] = response.AttributeResponse{
			ID:          attr.ID(),
			Name:        attr.Name(),
			Type:        attr.Type(),
			Description: attr.Description(),
			DomainID:    attr.DomainID(),
			CreatedAt:   attr.CreatedAt(),
			UpdatedAt:   attr.UpdatedAt(),
		}
	}

	return &response.AttributeListResponse{
		Attributes: attributeResponses,
		Total:      len(attributeResponses),
	}, nil
}
