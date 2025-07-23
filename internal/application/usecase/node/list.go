package node

import (
	"context"
	"url-db/internal/application/dto/response"
	"url-db/internal/domain/repository"
)

// ListNodesUseCase handles the listing of nodes
type ListNodesUseCase struct {
	nodeRepo repository.NodeRepository
}

// NewListNodesUseCase creates a new instance of ListNodesUseCase
func NewListNodesUseCase(repo repository.NodeRepository) *ListNodesUseCase {
	return &ListNodesUseCase{nodeRepo: repo}
}

// Execute performs the node listing use case
func (uc *ListNodesUseCase) Execute(ctx context.Context, domainName string, page, size int) (*response.NodeListResponse, error) {
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

	// Get nodes from repository
	nodes, totalCount, err := uc.nodeRepo.List(ctx, domainName, page, size)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	nodeResponses := make([]response.NodeResponse, len(nodes))
	for i, node := range nodes {
		nodeResponses[i] = response.NodeResponse{
			ID:          node.ID(),
			URL:         node.URL(),
			DomainName:  node.DomainName(),
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
