package nodes

import (
	"fmt"
	"math"
	"strings"

	"url-db/internal/models"
)

type NodeService interface {
	CreateNode(domainID int, req *models.CreateNodeRequest) (*models.Node, error)
	GetNodeByID(id int) (*models.Node, error)
	GetNodesByDomainID(domainID, page, size int) (*models.NodeListResponse, error)
	FindNodeByURL(domainID int, req *models.FindNodeByURLRequest) (*models.Node, error)
	UpdateNode(id int, req *models.UpdateNodeRequest) (*models.Node, error)
	DeleteNode(id int) error
	SearchNodes(domainID int, query string, page, size int) (*models.NodeListResponse, error)
}

type NodeServiceImpl struct {
	repo NodeRepository
}

func NewNodeService(repo NodeRepository) NodeService {
	return &NodeServiceImpl{repo: repo}
}

func (s *NodeServiceImpl) CreateNode(domainID int, req *models.CreateNodeRequest) (*models.Node, error) {
	// Validate domain exists
	exists, err := s.repo.CheckDomainExists(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to check domain existence: %w", err)
	}
	if !exists {
		return nil, ErrNodeDomainNotFound
	}
	
	// Validate URL
	if err := ValidateURL(req.URL); err != nil {
		return nil, err
	}
	
	// Check if node already exists
	existingNode, err := s.repo.GetByURL(domainID, req.URL)
	if err != nil && err != ErrNodeNotFound {
		return nil, fmt.Errorf("failed to check existing node: %w", err)
	}
	if existingNode != nil {
		return nil, ErrNodeAlreadyExists
	}
	
	// Generate title if not provided
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = GenerateTitleFromURL(req.URL)
	}
	
	// Create node
	node := &models.Node{
		Content:     req.URL,
		DomainID:    domainID,
		Title:       title,
		Description: strings.TrimSpace(req.Description),
	}
	
	if err := s.repo.Create(node); err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}
	
	return node, nil
}

func (s *NodeServiceImpl) GetNodeByID(id int) (*models.Node, error) {
	node, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get node by id: %w", err)
	}
	
	return node, nil
}

func (s *NodeServiceImpl) GetNodesByDomainID(domainID, page, size int) (*models.NodeListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	
	// Check if domain exists
	exists, err := s.repo.CheckDomainExists(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to check domain existence: %w", err)
	}
	if !exists {
		return nil, ErrNodeDomainNotFound
	}
	
	nodes, totalCount, err := s.repo.GetByDomainID(domainID, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes by domain id: %w", err)
	}
	
	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))
	
	return &models.NodeListResponse{
		Nodes:      nodes,
		TotalCount: totalCount,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}, nil
}

func (s *NodeServiceImpl) FindNodeByURL(domainID int, req *models.FindNodeByURLRequest) (*models.Node, error) {
	// Validate domain exists
	exists, err := s.repo.CheckDomainExists(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to check domain existence: %w", err)
	}
	if !exists {
		return nil, ErrNodeDomainNotFound
	}
	
	// Validate URL
	if err := ValidateURL(req.URL); err != nil {
		return nil, err
	}
	
	node, err := s.repo.GetByURL(domainID, req.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to find node by url: %w", err)
	}
	
	return node, nil
}

func (s *NodeServiceImpl) UpdateNode(id int, req *models.UpdateNodeRequest) (*models.Node, error) {
	// Get existing node
	node, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get node for update: %w", err)
	}
	
	// Update fields
	node.Title = strings.TrimSpace(req.Title)
	node.Description = strings.TrimSpace(req.Description)
	
	// Generate title if empty
	if node.Title == "" {
		node.Title = GenerateTitleFromURL(node.Content)
	}
	
	if err := s.repo.Update(node); err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}
	
	return node, nil
}

func (s *NodeServiceImpl) DeleteNode(id int) error {
	// Check if node exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get node for deletion: %w", err)
	}
	
	// Delete node (cascade will handle attributes)
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}
	
	return nil
}

func (s *NodeServiceImpl) SearchNodes(domainID int, query string, page, size int) (*models.NodeListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	
	// Check if domain exists
	exists, err := s.repo.CheckDomainExists(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to check domain existence: %w", err)
	}
	if !exists {
		return nil, ErrNodeDomainNotFound
	}
	
	// Validate search query
	query = strings.TrimSpace(query)
	if query == "" {
		return &models.NodeListResponse{
			Nodes:      []models.Node{},
			TotalCount: 0,
			Page:       page,
			Size:       size,
			TotalPages: 0,
		}, nil
	}
	
	nodes, totalCount, err := s.repo.Search(domainID, query, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to search nodes: %w", err)
	}
	
	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))
	
	return &models.NodeListResponse{
		Nodes:      nodes,
		TotalCount: totalCount,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}, nil
}