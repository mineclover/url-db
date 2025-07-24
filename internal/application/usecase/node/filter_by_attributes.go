package node

import (
	"context"
	"url-db/internal/application/dto/response"
	"url-db/internal/domain/repository"
)

// FilterNodesByAttributesUseCase handles filtering nodes by attributes
type FilterNodesByAttributesUseCase struct {
	nodeRepo repository.NodeRepository
}

// NewFilterNodesByAttributesUseCase creates a new instance of FilterNodesByAttributesUseCase
func NewFilterNodesByAttributesUseCase(repo repository.NodeRepository) *FilterNodesByAttributesUseCase {
	return &FilterNodesByAttributesUseCase{nodeRepo: repo}
}

// Execute performs the node filtering use case
func (uc *FilterNodesByAttributesUseCase) Execute(ctx context.Context, domainName string, filters []repository.AttributeFilter, page, size int) (*response.NodeListResponse, error) {
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

	// Get filtered nodes from repository
	nodes, totalCount, err := uc.nodeRepo.FilterByAttributes(ctx, domainName, filters, page, size)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	nodeResponses := make([]response.NodeResponse, len(nodes))
	for i, node := range nodes {
		nodeResponses[i] = response.NodeResponse{
			ID:          node.ID(),
			URL:         node.URL(),
			DomainName:  domainName, // Use domain name from parameter
			Title:       node.Title(),
			Description: node.Description(),
			CreatedAt:   node.CreatedAt(),
			UpdatedAt:   node.UpdatedAt(),
		}
	}

	// Calculate total pages
	totalPages := (totalCount + size - 1) / size

	return &response.NodeListResponse{
		Nodes:      nodeResponses,
		TotalCount: totalCount,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}, nil
}