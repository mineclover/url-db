package attribute

import (
	"context"
	"errors"
	"url-db/internal/application/dto/request"
	"url-db/internal/application/dto/response"
	"url-db/internal/constants"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
)

type CreateAttributeUseCase struct {
	attributeRepo repository.AttributeRepository
	domainRepo    repository.DomainRepository
}

func NewCreateAttributeUseCase(attributeRepo repository.AttributeRepository, domainRepo repository.DomainRepository) *CreateAttributeUseCase {
	return &CreateAttributeUseCase{
		attributeRepo: attributeRepo,
		domainRepo:    domainRepo,
	}
}

func (uc *CreateAttributeUseCase) Execute(ctx context.Context, req *request.CreateAttributeRequest) (*response.AttributeResponse, error) {
	// Verify domain exists
	domain, err := uc.domainRepo.GetByID(ctx, req.DomainID)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, errors.New(constants.ErrDomainNotFound)
	}

	// Create attribute entity
	attribute, err := entity.NewAttribute(req.Name, req.Type, req.Description, req.DomainID)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.attributeRepo.Create(ctx, attribute); err != nil {
		return nil, err
	}

	// Convert to response
	return &response.AttributeResponse{
		ID:          attribute.ID(),
		Name:        attribute.Name(),
		Type:        attribute.Type(),
		Description: attribute.Description(),
		DomainID:    attribute.DomainID(),
		CreatedAt:   attribute.CreatedAt(),
		UpdatedAt:   attribute.UpdatedAt(),
	}, nil
}
