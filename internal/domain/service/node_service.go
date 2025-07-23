package service

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
)

// NodeService represents node business logic
type NodeService interface {
	CreateNode(ctx context.Context, domainName, urlStr, title, description string) (*entity.Node, error)
	ValidateURL(urlStr string) error
	ValidateTitle(title string) error
	ValidateDescription(description string) error
	GenerateTitleFromURL(urlStr string) string
	NormalizeString(str string) string
}

type nodeService struct {
	nodeRepo   repository.NodeRepository
	domainRepo repository.DomainRepository
}

// NewNodeService creates a new node service
func NewNodeService(nodeRepo repository.NodeRepository, domainRepo repository.DomainRepository) NodeService {
	return &nodeService{
		nodeRepo:   nodeRepo,
		domainRepo: domainRepo,
	}
}

// ValidateURL validates URL format and content
func (s *nodeService) ValidateURL(urlStr string) error {
	if urlStr == "" {
		return errors.New("URL is required")
	}

	if len(urlStr) > 2048 {
		return errors.New("URL cannot exceed 2048 characters")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if parsedURL.Scheme == "" {
		return errors.New("URL must have a scheme (http:// or https://)")
	}

	if parsedURL.Host == "" {
		return errors.New("URL must have a host")
	}

	return nil
}

// ValidateTitle validates node title
func (s *nodeService) ValidateTitle(title string) error {
	if len(title) > 255 {
		return errors.New("title cannot exceed 255 characters")
	}
	return nil
}

// ValidateDescription validates node description
func (s *nodeService) ValidateDescription(description string) error {
	if len(description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}
	return nil
}

// NormalizeString normalizes string input
func (s *nodeService) NormalizeString(str string) string {
	return strings.TrimSpace(str)
}

// GenerateTitleFromURL generates a title from URL if not provided
func (s *nodeService) GenerateTitleFromURL(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	host := parsedURL.Host
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}

	return host
}

// CreateNode creates a new node with business validation
func (s *nodeService) CreateNode(ctx context.Context, domainName, urlStr, title, description string) (*entity.Node, error) {
	// Normalize inputs
	urlStr = s.NormalizeString(urlStr)
	title = s.NormalizeString(title)
	description = s.NormalizeString(description)

	// Validate inputs
	if err := s.ValidateURL(urlStr); err != nil {
		return nil, err
	}

	if err := s.ValidateTitle(title); err != nil {
		return nil, err
	}

	if err := s.ValidateDescription(description); err != nil {
		return nil, err
	}

	// Check if domain exists
	domain, err := s.domainRepo.GetByName(ctx, domainName)
	if err != nil {
		return nil, err
	}

	if domain == nil {
		return nil, errors.New("domain not found")
	}

	// Check if node already exists
	exists, err := s.nodeRepo.Exists(ctx, urlStr, domainName)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("node already exists in this domain")
	}

	// Generate title if not provided
	if title == "" {
		title = s.GenerateTitleFromURL(urlStr)
	}

	// Create node entity
	node, err := entity.NewNode(urlStr, title, description, domain.ID())
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.nodeRepo.Create(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}
