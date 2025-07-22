package mcp

import (
	"context"
	"fmt"

	"url-db/internal/models"
)

type MCPService interface {
	CreateNode(ctx context.Context, req *models.CreateMCPNodeRequest) (*models.MCPNode, error)
	GetNode(ctx context.Context, compositeID string) (*models.MCPNode, error)
	UpdateNode(ctx context.Context, compositeID string, req *models.UpdateNodeRequest) (*models.MCPNode, error)
	DeleteNode(ctx context.Context, compositeID string) error
	ListNodes(ctx context.Context, domainName string, page, size int, search string) (*models.MCPNodeListResponse, error)
	FindNodeByURL(ctx context.Context, req *models.FindMCPNodeRequest) (*models.MCPNode, error)
	BatchGetNodes(ctx context.Context, req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error)

	ListDomains(ctx context.Context) (*MCPDomainListResponse, error)
	CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*MCPDomain, error)

	GetNodeAttributes(ctx context.Context, compositeID string) (*models.MCPNodeAttributeResponse, error)
	SetNodeAttributes(ctx context.Context, compositeID string, req *models.SetMCPNodeAttributesRequest) (*models.MCPNodeAttributeResponse, error)

	// Domain attribute management methods
	ListDomainAttributes(ctx context.Context, domainName string) (*models.MCPDomainAttributeListResponse, error)
	CreateDomainAttribute(ctx context.Context, domainName string, req *models.CreateAttributeRequest) (*models.MCPDomainAttribute, error)
	GetDomainAttribute(ctx context.Context, compositeID string) (*models.MCPDomainAttribute, error)
	UpdateDomainAttribute(ctx context.Context, compositeID string, req *models.UpdateAttributeRequest) (*models.MCPDomainAttribute, error)
	DeleteDomainAttribute(ctx context.Context, compositeID string) error

	GetServerInfo(ctx context.Context) (*MCPServerInfo, error)
}

type NodeService interface {
	CreateNode(ctx context.Context, domainID int, req *models.CreateNodeRequest) (*models.Node, error)
	GetNode(ctx context.Context, nodeID int) (*models.Node, error)
	UpdateNode(ctx context.Context, nodeID int, req *models.UpdateNodeRequest) (*models.Node, error)
	DeleteNode(ctx context.Context, nodeID int) error
	ListNodes(ctx context.Context, domainID int, page, size int, search string) (*models.NodeListResponse, error)
	FindNodeByURL(ctx context.Context, domainID int, req *models.FindNodeByURLRequest) (*models.Node, error)
}

type DomainService interface {
	CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*models.Domain, error)
	GetDomain(ctx context.Context, domainID int) (*models.Domain, error)
	GetDomainByName(ctx context.Context, name string) (*models.Domain, error)
	ListDomains(ctx context.Context, page, size int) (*models.DomainListResponse, error)
	UpdateDomain(ctx context.Context, domainID int, req *models.UpdateDomainRequest) (*models.Domain, error)
	DeleteDomain(ctx context.Context, domainID int) error
}

type AttributeService interface {
	GetNodeAttributes(ctx context.Context, nodeID int) ([]models.NodeAttributeWithInfo, error)
	SetNodeAttribute(ctx context.Context, nodeID int, req *models.CreateNodeAttributeRequest) (*models.NodeAttributeWithInfo, error)
	GetAttributeByName(ctx context.Context, domainID int, name string) (*models.Attribute, error)
	DeleteNodeAttribute(ctx context.Context, nodeID, attributeID int) error
	
	// Domain attribute management methods
	ListAttributes(ctx context.Context, domainID int) (*models.AttributeListResponse, error)
	CreateAttribute(ctx context.Context, domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error)
	GetAttribute(ctx context.Context, id int) (*models.Attribute, error)
	UpdateAttribute(ctx context.Context, id int, req *models.UpdateAttributeRequest) (*models.Attribute, error)
	DeleteAttribute(ctx context.Context, id int) error
}

type NodeCountService interface {
	GetNodeCountByDomain(ctx context.Context, domainID int) (int, error)
}

type mcpService struct {
	nodeService      NodeService
	domainService    DomainService
	attributeService AttributeService
	nodeCountService NodeCountService
	converter        *Converter
}

func NewMCPService(
	nodeService NodeService,
	domainService DomainService,
	attributeService AttributeService,
	nodeCountService NodeCountService,
	converter *Converter,
) MCPService {
	return &mcpService{
		nodeService:      nodeService,
		domainService:    domainService,
		attributeService: attributeService,
		nodeCountService: nodeCountService,
		converter:        converter,
	}
}

func (s *mcpService) CreateNode(ctx context.Context, req *models.CreateMCPNodeRequest) (*models.MCPNode, error) {
	domain, err := s.domainService.GetDomainByName(ctx, req.DomainName)
	if err != nil {
		return nil, NewDomainNotFoundError(req.DomainName)
	}

	nodeReq := s.converter.CreateMCPNodeRequestToCreateNodeRequest(req)
	node, err := s.nodeService.CreateNode(ctx, domain.ID, nodeReq)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to create node: %v", err))
	}

	return s.converter.NodeToMCPNode(node, domain)
}

func (s *mcpService) GetNode(ctx context.Context, compositeID string) (*models.MCPNode, error) {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domainName, err := s.converter.ExtractDomainNameFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return nil, NewResourceNotFoundError(compositeID)
	}

	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	return s.converter.NodeToMCPNode(node, domain)
}

func (s *mcpService) UpdateNode(ctx context.Context, compositeID string, req *models.UpdateNodeRequest) (*models.MCPNode, error) {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domainName, err := s.converter.ExtractDomainNameFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	node, err := s.nodeService.UpdateNode(ctx, nodeID, req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to update node: %v", err))
	}

	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	return s.converter.NodeToMCPNode(node, domain)
}

func (s *mcpService) DeleteNode(ctx context.Context, compositeID string) error {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return NewInvalidCompositeKeyError(compositeID)
	}

	if err := s.nodeService.DeleteNode(ctx, nodeID); err != nil {
		return NewInternalServerError(fmt.Sprintf("failed to delete node: %v", err))
	}

	return nil
}

func (s *mcpService) ListNodes(ctx context.Context, domainName string, page, size int, search string) (*models.MCPNodeListResponse, error) {
	var domainID int
	var domain *models.Domain
	var err error

	if domainName != "" {
		domain, err = s.domainService.GetDomainByName(ctx, domainName)
		if err != nil {
			return nil, NewDomainNotFoundError(domainName)
		}
		domainID = domain.ID
	}

	response, err := s.nodeService.ListNodes(ctx, domainID, page, size, search)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list nodes: %v", err))
	}

	mcpNodes := make([]models.MCPNode, 0, len(response.Nodes))
	for _, node := range response.Nodes {
		if domain == nil {
			domain, err = s.domainService.GetDomain(ctx, node.DomainID)
			if err != nil {
				continue
			}
		}

		mcpNode, err := s.converter.NodeToMCPNode(&node, domain)
		if err != nil {
			continue
		}
		mcpNodes = append(mcpNodes, *mcpNode)
	}

	return &models.MCPNodeListResponse{
		Nodes:      mcpNodes,
		TotalCount: response.TotalCount,
		Page:       response.Page,
		Size:       response.Size,
		TotalPages: response.TotalPages,
	}, nil
}

func (s *mcpService) FindNodeByURL(ctx context.Context, req *models.FindMCPNodeRequest) (*models.MCPNode, error) {
	domain, err := s.domainService.GetDomainByName(ctx, req.DomainName)
	if err != nil {
		return nil, NewDomainNotFoundError(req.DomainName)
	}

	findReq := &models.FindNodeByURLRequest{URL: req.URL}
	node, err := s.nodeService.FindNodeByURL(ctx, domain.ID, findReq)
	if err != nil {
		return nil, NewNodeNotFoundError(req.DomainName, req.URL)
	}

	return s.converter.NodeToMCPNode(node, domain)
}

func (s *mcpService) BatchGetNodes(ctx context.Context, req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error) {
	nodes := make([]models.MCPNode, 0, len(req.CompositeIDs))
	notFound := make([]string, 0)

	for _, compositeID := range req.CompositeIDs {
		node, err := s.GetNode(ctx, compositeID)
		if err != nil {
			notFound = append(notFound, compositeID)
			continue
		}
		nodes = append(nodes, *node)
	}

	return &models.BatchMCPNodeResponse{
		Nodes:    nodes,
		NotFound: notFound,
	}, nil
}

func (s *mcpService) ListDomains(ctx context.Context) (*MCPDomainListResponse, error) {
	response, err := s.domainService.ListDomains(ctx, 1, 1000)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list domains: %v", err))
	}

	mcpDomains := make([]MCPDomain, 0, len(response.Domains))
	for _, domain := range response.Domains {
		nodeCount, err := s.nodeCountService.GetNodeCountByDomain(ctx, domain.ID)
		if err != nil {
			nodeCount = 0
		}

		mcpDomain := s.converter.DomainToMCPDomain(&domain, nodeCount)
		if mcpDomain != nil {
			mcpDomains = append(mcpDomains, *mcpDomain)
		}
	}

	return &MCPDomainListResponse{
		Domains: mcpDomains,
	}, nil
}

func (s *mcpService) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*MCPDomain, error) {
	domain, err := s.domainService.CreateDomain(ctx, req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to create domain: %v", err))
	}

	return s.converter.DomainToMCPDomain(domain, 0), nil
}

func (s *mcpService) GetNodeAttributes(ctx context.Context, compositeID string) (*models.MCPNodeAttributeResponse, error) {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attributes, err := s.attributeService.GetNodeAttributes(ctx, nodeID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get node attributes: %v", err))
	}

	mcpAttributes := s.converter.NodeAttributesToMCPAttributes(attributes)

	return &models.MCPNodeAttributeResponse{
		CompositeID: compositeID,
		Attributes:  mcpAttributes,
	}, nil
}

func (s *mcpService) SetNodeAttributes(ctx context.Context, compositeID string, req *models.SetMCPNodeAttributesRequest) (*models.MCPNodeAttributeResponse, error) {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domainName, err := s.converter.ExtractDomainNameFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	existingAttributes, err := s.attributeService.GetNodeAttributes(ctx, nodeID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get existing attributes: %v", err))
	}

	existingAttrMap := make(map[string]int)
	for _, attr := range existingAttributes {
		existingAttrMap[attr.Name] = attr.AttributeID
	}

	for _, attrReq := range req.Attributes {
		attribute, err := s.attributeService.GetAttributeByName(ctx, domain.ID, attrReq.Name)
		if err != nil {
			return nil, NewValidationError(fmt.Sprintf("attribute '%s' not found", attrReq.Name))
		}

		if existingAttrID, exists := existingAttrMap[attrReq.Name]; exists {
			if err := s.attributeService.DeleteNodeAttribute(ctx, nodeID, existingAttrID); err != nil {
				return nil, NewInternalServerError(fmt.Sprintf("failed to delete existing attribute: %v", err))
			}
		}

		createReq := &models.CreateNodeAttributeRequest{
			AttributeID: attribute.ID,
			Value:       attrReq.Value,
			OrderIndex:  attrReq.OrderIndex,
		}

		if _, err := s.attributeService.SetNodeAttribute(ctx, nodeID, createReq); err != nil {
			return nil, NewInternalServerError(fmt.Sprintf("failed to set attribute: %v", err))
		}
	}

	return s.GetNodeAttributes(ctx, compositeID)
}

func (s *mcpService) GetServerInfo(ctx context.Context) (*MCPServerInfo, error) {
	return &MCPServerInfo{
		Name:        "url-db",
		Version:     "1.0.0",
		Description: "URL 데이터베이스 MCP 서버",
		Capabilities: []string{
			"resources",
			"tools",
			"prompts",
		},
		CompositeKeyFormat: "url-db:domain_name:id",
	}, nil
}

// Domain attribute management methods
func (s *mcpService) ListDomainAttributes(ctx context.Context, domainName string) (*models.MCPDomainAttributeListResponse, error) {
	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	response, err := s.attributeService.ListAttributes(ctx, domain.ID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list attributes: %v", err))
	}

	mcpAttributes := make([]models.MCPDomainAttribute, len(response.Attributes))
	for i, attr := range response.Attributes {
		compositeID := s.converter.CreateAttributeCompositeID(domain.Name, attr.ID)
		mcpAttributes[i] = models.MCPDomainAttribute{
			CompositeID: compositeID,
			Name:        attr.Name,
			Type:        attr.Type,
			Description: attr.Description,
			CreatedAt:   attr.CreatedAt,
			UpdatedAt:   attr.CreatedAt, // Assuming no updated_at field in Attribute model
		}
	}

	return &models.MCPDomainAttributeListResponse{
		DomainName:  domainName,
		Attributes:  mcpAttributes,
		TotalCount:  len(mcpAttributes),
	}, nil
}

func (s *mcpService) CreateDomainAttribute(ctx context.Context, domainName string, req *models.CreateAttributeRequest) (*models.MCPDomainAttribute, error) {
	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	attribute, err := s.attributeService.CreateAttribute(ctx, domain.ID, req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to create attribute: %v", err))
	}

	compositeID := s.converter.CreateAttributeCompositeID(domain.Name, attribute.ID)
	return &models.MCPDomainAttribute{
		CompositeID: compositeID,
		Name:        attribute.Name,
		Type:        attribute.Type,
		Description: attribute.Description,
		CreatedAt:   attribute.CreatedAt,
		UpdatedAt:   attribute.CreatedAt,
	}, nil
}

func (s *mcpService) GetDomainAttribute(ctx context.Context, compositeID string) (*models.MCPDomainAttribute, error) {
	if err := s.converter.ValidateAttributeCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attributeID, err := s.converter.ExtractAttributeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attribute, err := s.attributeService.GetAttribute(ctx, attributeID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get attribute: %v", err))
	}

	return &models.MCPDomainAttribute{
		CompositeID: compositeID,
		Name:        attribute.Name,
		Type:        attribute.Type,
		Description: attribute.Description,
		CreatedAt:   attribute.CreatedAt,
		UpdatedAt:   attribute.CreatedAt,
	}, nil
}

func (s *mcpService) UpdateDomainAttribute(ctx context.Context, compositeID string, req *models.UpdateAttributeRequest) (*models.MCPDomainAttribute, error) {
	if err := s.converter.ValidateAttributeCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attributeID, err := s.converter.ExtractAttributeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attribute, err := s.attributeService.UpdateAttribute(ctx, attributeID, req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to update attribute: %v", err))
	}

	return &models.MCPDomainAttribute{
		CompositeID: compositeID,
		Name:        attribute.Name,
		Type:        attribute.Type,
		Description: attribute.Description,
		CreatedAt:   attribute.CreatedAt,
		UpdatedAt:   attribute.CreatedAt,
	}, nil
}

func (s *mcpService) DeleteDomainAttribute(ctx context.Context, compositeID string) error {
	if err := s.converter.ValidateAttributeCompositeID(compositeID); err != nil {
		return NewInvalidCompositeKeyError(compositeID)
	}

	attributeID, err := s.converter.ExtractAttributeIDFromCompositeID(compositeID)
	if err != nil {
		return NewInvalidCompositeKeyError(compositeID)
	}

	if err := s.attributeService.DeleteAttribute(ctx, attributeID); err != nil {
		return NewInternalServerError(fmt.Sprintf("failed to delete attribute: %v", err))
	}

	return nil
}
