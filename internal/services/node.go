package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"url-db/internal/models"
)

type NodeRepository interface {
	Create(ctx context.Context, node *models.Node) error
	GetByID(ctx context.Context, id int) (*models.Node, error)
	GetByDomainAndContent(ctx context.Context, domainID int, content string) (*models.Node, error)
	ListByDomain(ctx context.Context, domainID int, page, size int, search string) ([]*models.Node, int, error)
	Update(ctx context.Context, node *models.Node) error
	Delete(ctx context.Context, id int) error
	ExistsByDomainAndContent(ctx context.Context, domainID int, content string) (bool, error)
}

type nodeService struct {
	nodeRepo   NodeRepository
	domainRepo DomainRepository
	logger     *log.Logger
}

func NewNodeService(nodeRepo NodeRepository, domainRepo DomainRepository, logger *log.Logger) NodeService {
	return &nodeService{
		nodeRepo:   nodeRepo,
		domainRepo: domainRepo,
		logger:     logger,
	}
}

func (s *nodeService) CreateNode(ctx context.Context, domainID int, req *models.CreateNodeRequest) (*models.Node, error) {
	if err := validatePositiveInteger(domainID, "domainID"); err != nil {
		return nil, err
	}

	req.URL = normalizeString(req.URL)
	req.Title = normalizeString(req.Title)
	req.Description = normalizeString(req.Description)

	if err := validateURL(req.URL); err != nil {
		return nil, err
	}

	if err := validateTitle(req.Title); err != nil {
		return nil, err
	}

	if err := validateDescription(req.Description); err != nil {
		return nil, err
	}

	domain, err := s.domainRepo.GetByID(ctx, domainID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewDomainNotFoundError(domainID)
		}
		s.logger.Printf("Failed to get domain: %v", err)
		return nil, err
	}

	exists, err := s.nodeRepo.ExistsByDomainAndContent(ctx, domainID, req.URL)
	if err != nil {
		s.logger.Printf("Failed to check node existence: %v", err)
		return nil, err
	}
	if exists {
		return nil, NewNodeAlreadyExistsError(req.URL)
	}

	title := req.Title
	if title == "" {
		title = generateTitleFromURL(req.URL)
	}

	node := &models.Node{
		Content:     req.URL,
		DomainID:    domainID,
		Title:       title,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.nodeRepo.Create(ctx, node); err != nil {
		s.logger.Printf("Failed to create node: %v", err)
		return nil, err
	}

	s.logger.Printf("Created node: %s in domain %s (ID: %d)", node.Content, domain.Name, node.ID)
	return node, nil
}

func (s *nodeService) GetNode(ctx context.Context, id int) (*models.Node, error) {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return nil, err
	}

	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNodeNotFoundError(id)
		}
		s.logger.Printf("Failed to get node: %v", err)
		return nil, err
	}

	return node, nil
}

func (s *nodeService) GetNodeByDomainAndURL(ctx context.Context, domainID int, url string) (*models.Node, error) {
	if err := validatePositiveInteger(domainID, "domainID"); err != nil {
		return nil, err
	}

	url = normalizeString(url)
	if err := validateURL(url); err != nil {
		return nil, err
	}

	node, err := s.nodeRepo.GetByDomainAndContent(ctx, domainID, url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNodeNotFoundError(0)
		}
		s.logger.Printf("Failed to get node by domain and URL: %v", err)
		return nil, err
	}

	return node, nil
}

func (s *nodeService) ListNodesByDomain(ctx context.Context, domainID int, page, size int, search string) (*models.NodeListResponse, error) {
	if err := validatePositiveInteger(domainID, "domainID"); err != nil {
		return nil, err
	}

	page, size, err := validatePaginationParams(page, size)
	if err != nil {
		return nil, err
	}

	search = normalizeString(search)

	_, err = s.domainRepo.GetByID(ctx, domainID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewDomainNotFoundError(domainID)
		}
		s.logger.Printf("Failed to get domain: %v", err)
		return nil, err
	}

	nodes, totalCount, err := s.nodeRepo.ListByDomain(ctx, domainID, page, size, search)
	if err != nil {
		s.logger.Printf("Failed to list nodes: %v", err)
		return nil, err
	}

	nodeList := make([]models.Node, len(nodes))
	for i, node := range nodes {
		nodeList[i] = *node
	}

	totalPages := (totalCount + size - 1) / size

	return &models.NodeListResponse{
		Nodes:      nodeList,
		TotalCount: totalCount,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}, nil
}

func (s *nodeService) UpdateNode(ctx context.Context, id int, req *models.UpdateNodeRequest) (*models.Node, error) {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return nil, err
	}

	req.Title = normalizeString(req.Title)
	req.Description = normalizeString(req.Description)

	if err := validateTitle(req.Title); err != nil {
		return nil, err
	}

	if err := validateDescription(req.Description); err != nil {
		return nil, err
	}

	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNodeNotFoundError(id)
		}
		s.logger.Printf("Failed to get node: %v", err)
		return nil, err
	}

	node.Title = req.Title
	node.Description = req.Description
	node.UpdatedAt = time.Now()

	if err := s.nodeRepo.Update(ctx, node); err != nil {
		s.logger.Printf("Failed to update node: %v", err)
		return nil, err
	}

	s.logger.Printf("Updated node: %s (ID: %d)", node.Content, node.ID)
	return node, nil
}

func (s *nodeService) DeleteNode(ctx context.Context, id int) error {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return err
	}

	err := s.nodeRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewNodeNotFoundError(id)
		}
		s.logger.Printf("Failed to delete node: %v", err)
		return err
	}

	s.logger.Printf("Deleted node with ID: %d", id)
	return nil
}

func (s *nodeService) FindNodeByURL(ctx context.Context, domainID int, req *models.FindNodeByURLRequest) (*models.Node, error) {
	if err := validatePositiveInteger(domainID, "domainID"); err != nil {
		return nil, err
	}

	req.URL = normalizeString(req.URL)
	if err := validateURL(req.URL); err != nil {
		return nil, err
	}

	return s.GetNodeByDomainAndURL(ctx, domainID, req.URL)
}
