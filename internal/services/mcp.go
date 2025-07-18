package services

import (
	"context"
	"fmt"
	"log"

	"url-db/internal/models"
)

type mcpService struct {
	nodeService          NodeService
	domainService        DomainService
	attributeService     AttributeService
	nodeAttributeService NodeAttributeService
	compositeKeyService  CompositeKeyService
	toolName             string
	version              string
	logger               *log.Logger
}

func NewMCPService(
	nodeService NodeService,
	domainService DomainService,
	attributeService AttributeService,
	nodeAttributeService NodeAttributeService,
	compositeKeyService CompositeKeyService,
	toolName, version string,
	logger *log.Logger,
) MCPService {
	return &mcpService{
		nodeService:          nodeService,
		domainService:        domainService,
		attributeService:     attributeService,
		nodeAttributeService: nodeAttributeService,
		compositeKeyService:  compositeKeyService,
		toolName:             toolName,
		version:              version,
		logger:               logger,
	}
}

func (s *mcpService) CreateNode(ctx context.Context, req *models.CreateMCPNodeRequest) (*models.MCPNode, error) {
	req.DomainName = normalizeString(req.DomainName)
	req.URL = normalizeString(req.URL)
	req.Title = normalizeString(req.Title)
	req.Description = normalizeString(req.Description)

	if err := validateDomainName(req.DomainName); err != nil {
		return nil, err
	}

	domain, err := s.domainService.GetDomainByName(ctx, req.DomainName)
	if err != nil {
		if serviceErr, ok := err.(*ServiceError); ok && serviceErr.Code == "DOMAIN_NOT_FOUND" {
			createDomainReq := &models.CreateDomainRequest{
				Name:        req.DomainName,
				Description: fmt.Sprintf("Auto-created domain for %s", req.DomainName),
			}
			domain, err = s.domainService.CreateDomain(ctx, createDomainReq)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	createNodeReq := &models.CreateNodeRequest{
		URL:         req.URL,
		Title:       req.Title,
		Description: req.Description,
	}

	node, err := s.nodeService.CreateNode(ctx, domain.ID, createNodeReq)
	if err != nil {
		return nil, err
	}

	return s.convertToMCPNode(node, domain), nil
}

func (s *mcpService) GetNode(ctx context.Context, compositeID string) (*models.MCPNode, error) {
	ck, err := s.compositeKeyService.Parse(compositeID)
	if err != nil {
		return nil, err
	}

	node, err := s.nodeService.GetNode(ctx, ck.ID)
	if err != nil {
		return nil, err
	}

	domain, err := s.domainService.GetDomain(ctx, node.DomainID)
	if err != nil {
		return nil, err
	}

	return s.convertToMCPNode(node, domain), nil
}

func (s *mcpService) ListNodes(ctx context.Context, domainName string, page, size int, search string) (*models.MCPNodeListResponse, error) {
	domainName = normalizeString(domainName)

	if err := validateDomainName(domainName); err != nil {
		return nil, err
	}

	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, err
	}

	nodeListResp, err := s.nodeService.ListNodesByDomain(ctx, domain.ID, page, size, search)
	if err != nil {
		return nil, err
	}

	mcpNodes := make([]models.MCPNode, len(nodeListResp.Nodes))
	for i, node := range nodeListResp.Nodes {
		mcpNodes[i] = *s.convertToMCPNode(&node, domain)
	}

	return &models.MCPNodeListResponse{
		Nodes:      mcpNodes,
		TotalCount: nodeListResp.TotalCount,
		Page:       nodeListResp.Page,
		Size:       nodeListResp.Size,
		TotalPages: nodeListResp.TotalPages,
	}, nil
}

func (s *mcpService) UpdateNode(ctx context.Context, compositeID string, req *models.UpdateMCPNodeRequest) (*models.MCPNode, error) {
	ck, err := s.compositeKeyService.Parse(compositeID)
	if err != nil {
		return nil, err
	}

	updateReq := &models.UpdateNodeRequest{
		Title:       req.Title,
		Description: req.Description,
	}

	node, err := s.nodeService.UpdateNode(ctx, ck.ID, updateReq)
	if err != nil {
		return nil, err
	}

	domain, err := s.domainService.GetDomain(ctx, node.DomainID)
	if err != nil {
		return nil, err
	}

	return s.convertToMCPNode(node, domain), nil
}

func (s *mcpService) DeleteNode(ctx context.Context, compositeID string) error {
	ck, err := s.compositeKeyService.Parse(compositeID)
	if err != nil {
		return err
	}

	return s.nodeService.DeleteNode(ctx, ck.ID)
}

func (s *mcpService) FindNodeByURL(ctx context.Context, req *models.FindMCPNodeRequest) (*models.MCPNode, error) {
	req.DomainName = normalizeString(req.DomainName)

	if err := validateDomainName(req.DomainName); err != nil {
		return nil, err
	}

	domain, err := s.domainService.GetDomainByName(ctx, req.DomainName)
	if err != nil {
		return nil, err
	}

	findReq := &models.FindNodeByURLRequest{
		URL: req.URL,
	}

	node, err := s.nodeService.FindNodeByURL(ctx, domain.ID, findReq)
	if err != nil {
		return nil, err
	}

	return s.convertToMCPNode(node, domain), nil
}

func (s *mcpService) BatchGetNodes(ctx context.Context, req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error) {
	results := make([]models.MCPNode, 0, len(req.CompositeIDs))
	errors := make([]string, 0)

	for _, compositeID := range req.CompositeIDs {
		node, err := s.GetNode(ctx, compositeID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to get node %s: %v", compositeID, err))
			continue
		}
		results = append(results, *node)
	}

	return &models.BatchMCPNodeResponse{
		Nodes:  results,
		Errors: errors,
	}, nil
}

func (s *mcpService) ListDomains(ctx context.Context) (*models.MCPDomainListResponse, error) {
	domainsResp, err := s.domainService.ListDomains(ctx, 1, 100)
	if err != nil {
		return nil, err
	}

	mcpDomains := make([]models.MCPDomain, len(domainsResp.Domains))
	for i, domain := range domainsResp.Domains {
		mcpDomains[i] = models.MCPDomain{
			Name:        domain.Name,
			Description: domain.Description,
			CreatedAt:   domain.CreatedAt,
			UpdatedAt:   domain.UpdatedAt,
		}
	}

	return &models.MCPDomainListResponse{
		Domains: mcpDomains,
	}, nil
}

func (s *mcpService) CreateDomain(ctx context.Context, req *models.CreateMCPDomainRequest) (*models.MCPDomain, error) {
	createReq := &models.CreateDomainRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	domain, err := s.domainService.CreateDomain(ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &models.MCPDomain{
		Name:        domain.Name,
		Description: domain.Description,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}, nil
}

func (s *mcpService) GetNodeAttributes(ctx context.Context, compositeID string) (*models.MCPNodeAttributeResponse, error) {
	ck, err := s.compositeKeyService.Parse(compositeID)
	if err != nil {
		return nil, err
	}

	nodeAttrsResp, err := s.nodeAttributeService.ListNodeAttributesByNode(ctx, ck.ID)
	if err != nil {
		return nil, err
	}

	attributes := make([]models.MCPAttribute, len(nodeAttrsResp.NodeAttributes))
	for i, nodeAttr := range nodeAttrsResp.NodeAttributes {
		attributes[i] = models.MCPAttribute{
			Name:  fmt.Sprintf("attr_%d", nodeAttr.AttributeID),
			Type:  string(nodeAttr.Type),
			Value: nodeAttr.Value,
		}
	}

	return &models.MCPNodeAttributeResponse{
		CompositeID: compositeID,
		Attributes:  attributes,
	}, nil
}

func (s *mcpService) SetNodeAttributes(ctx context.Context, compositeID string, req *models.SetMCPNodeAttributesRequest) error {
	ck, err := s.compositeKeyService.Parse(compositeID)
	if err != nil {
		return err
	}

	for _, attr := range req.Attributes {
		// TODO: Need to implement logic to find or create attribute by name
		// For now, this is a placeholder implementation
		s.logger.Printf("Setting attribute %s=%s for node %d", attr.Name, attr.Value, ck.ID)
		// The actual implementation would need to:
		// 1. Find the attribute by name in the domain
		// 2. Create a NodeAttribute with the found attribute ID
	}

	return nil
}

func (s *mcpService) GetServerInfo(ctx context.Context) (*models.MCPServerInfo, error) {
	return &models.MCPServerInfo{
		Name:        s.toolName,
		Version:     s.version,
		Description: "URL Database Management System",
	}, nil
}

func (s *mcpService) convertToMCPNode(node *models.Node, domain *models.Domain) *models.MCPNode {
	compositeID := s.compositeKeyService.Create(domain.Name, node.ID)

	return &models.MCPNode{
		CompositeID: compositeID,
		URL:         node.Content,
		DomainName:  domain.Name,
		Title:       node.Title,
		Description: node.Description,
		CreatedAt:   node.CreatedAt,
		UpdatedAt:   node.UpdatedAt,
	}
}
