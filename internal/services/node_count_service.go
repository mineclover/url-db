package services

import (
	"context"
	"fmt"
)

// NodeCountService provides methods to count nodes by domain
type NodeCountService interface {
	GetNodeCountByDomain(ctx context.Context, domainID int) (int, error)
}

// NodeCountRepository interface for data access
type NodeCountRepository interface {
	CountNodesByDomain(ctx context.Context, domainID int) (int, error)
}

// nodeCountService implements NodeCountService
type nodeCountService struct {
	nodeRepo NodeCountRepository
}

// NewNodeCountService creates a new node count service
func NewNodeCountService(nodeRepo NodeCountRepository) NodeCountService {
	return &nodeCountService{
		nodeRepo: nodeRepo,
	}
}

// GetNodeCountByDomain returns the count of nodes for a specific domain
func (s *nodeCountService) GetNodeCountByDomain(ctx context.Context, domainID int) (int, error) {
	if domainID <= 0 {
		return 0, fmt.Errorf("domain ID must be positive")
	}

	count, err := s.nodeRepo.CountNodesByDomain(ctx, domainID)
	if err != nil {
		return 0, fmt.Errorf("failed to count nodes by domain: %w", err)
	}

	return count, nil
}