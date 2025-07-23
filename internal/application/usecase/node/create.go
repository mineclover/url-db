package node

import (
	"context"
	"errors"
	"url-db/internal/application/dto/request"
	"url-db/internal/application/dto/response"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
)

// CreateNodeUseCase handles the creation of a new node
type CreateNodeUseCase struct {
	nodeRepo   repository.NodeRepository
	domainRepo repository.DomainRepository
}

// NewCreateNodeUseCase creates a new instance of CreateNodeUseCase
func NewCreateNodeUseCase(nodeRepo repository.NodeRepository, domainRepo repository.DomainRepository) *CreateNodeUseCase {
	return &CreateNodeUseCase{
		nodeRepo:   nodeRepo,
		domainRepo: domainRepo,
	}
}

// Execute performs the node creation use case
func (uc *CreateNodeUseCase) Execute(ctx context.Context, req *request.CreateNodeRequest) (*response.NodeResponse, error) {
	// Check if domain exists
	domain, err := uc.domainRepo.GetByName(ctx, req.DomainName)
	if err != nil {
		return nil, err
	}

	if domain == nil {
		return nil, errors.New("domain not found")
	}

	// Create node entity
	node, err := entity.NewNode(req.URL, req.DomainName, req.Title, req.Description)
	if err != nil {
		return nil, err
	}

	// Check if node already exists
	exists, err := uc.nodeRepo.Exists(ctx, req.URL, req.DomainName)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("node already exists in this domain")
	}

	// Save to repository
	if err := uc.nodeRepo.Create(ctx, node); err != nil {
		return nil, err
	}

	// Convert to response
	return &response.NodeResponse{
		ID:          node.ID(),
		URL:         node.URL(),
		DomainName:  node.DomainName(),
		Title:       node.Title(),
		Description: node.Description(),
		CreatedAt:   node.CreatedAt(),
		UpdatedAt:   node.UpdatedAt(),
	}, nil
}
