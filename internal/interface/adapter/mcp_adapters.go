package adapter

import (
	"context"
	"fmt"

	"url-db/internal/attributes"
	"url-db/internal/domains"
	"url-db/internal/interfaces/mcp"
	"url-db/internal/models"
	"url-db/internal/nodeattributes"
	"url-db/internal/nodes"
	"url-db/internal/repositories"
)

// MCPDomainServiceAdapter adapts the legacy domain service for MCP usage
type MCPDomainServiceAdapter struct {
	domainService domains.DomainService
	domainRepo    domains.DomainRepository
}

// NewMCPDomainServiceAdapter creates a new MCPDomainServiceAdapter
func NewMCPDomainServiceAdapter(domainService domains.DomainService, domainRepo domains.DomainRepository) *MCPDomainServiceAdapter {
	return &MCPDomainServiceAdapter{
		domainService: domainService,
		domainRepo:    domainRepo,
	}
}

func (a *MCPDomainServiceAdapter) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*models.Domain, error) {
	return a.domainService.CreateDomain(ctx, req)
}

func (a *MCPDomainServiceAdapter) GetDomain(ctx context.Context, domainID int) (*models.Domain, error) {
	return a.domainService.GetDomain(ctx, domainID)
}

func (a *MCPDomainServiceAdapter) GetDomainByName(ctx context.Context, name string) (*models.Domain, error) {
	return a.domainRepo.GetByName(ctx, name)
}

func (a *MCPDomainServiceAdapter) ListDomains(ctx context.Context, page, size int) (*models.DomainListResponse, error) {
	return a.domainService.ListDomains(ctx, page, size)
}

func (a *MCPDomainServiceAdapter) UpdateDomain(ctx context.Context, domainID int, req *models.UpdateDomainRequest) (*models.Domain, error) {
	return a.domainService.UpdateDomain(ctx, domainID, req)
}

func (a *MCPDomainServiceAdapter) DeleteDomain(ctx context.Context, domainID int) error {
	return a.domainService.DeleteDomain(ctx, domainID)
}

// MCPNodeServiceAdapter adapts the legacy node service for MCP usage
type MCPNodeServiceAdapter struct {
	nodeService nodes.NodeService
}

// NewMCPNodeServiceAdapter creates a new MCPNodeServiceAdapter
func NewMCPNodeServiceAdapter(nodeService nodes.NodeService) *MCPNodeServiceAdapter {
	return &MCPNodeServiceAdapter{
		nodeService: nodeService,
	}
}

func (a *MCPNodeServiceAdapter) CreateNode(ctx context.Context, domainID int, req *models.CreateNodeRequest) (*models.Node, error) {
	return a.nodeService.CreateNode(domainID, req)
}

func (a *MCPNodeServiceAdapter) GetNode(ctx context.Context, nodeID int) (*models.Node, error) {
	return a.nodeService.GetNodeByID(nodeID)
}

func (a *MCPNodeServiceAdapter) UpdateNode(ctx context.Context, nodeID int, req *models.UpdateNodeRequest) (*models.Node, error) {
	return a.nodeService.UpdateNode(nodeID, req)
}

func (a *MCPNodeServiceAdapter) DeleteNode(ctx context.Context, nodeID int) error {
	return a.nodeService.DeleteNode(nodeID)
}

func (a *MCPNodeServiceAdapter) ListNodes(ctx context.Context, domainID int, page, size int, search string) (*models.NodeListResponse, error) {
	if search != "" {
		return a.nodeService.SearchNodes(domainID, search, page, size)
	}
	return a.nodeService.GetNodesByDomainID(domainID, page, size)
}

func (a *MCPNodeServiceAdapter) FindNodeByURL(ctx context.Context, domainID int, req *models.FindNodeByURLRequest) (*models.Node, error) {
	return a.nodeService.FindNodeByURL(domainID, req)
}

// MCPAttributeServiceAdapter adapts the legacy attribute services for MCP usage
type MCPAttributeServiceAdapter struct {
	nodeAttributeService nodeattributes.Service
	attributeService     attributes.AttributeService
}

// NewMCPAttributeServiceAdapter creates a new MCPAttributeServiceAdapter
func NewMCPAttributeServiceAdapter(nodeAttributeService nodeattributes.Service, attributeService attributes.AttributeService) *MCPAttributeServiceAdapter {
	return &MCPAttributeServiceAdapter{
		nodeAttributeService: nodeAttributeService,
		attributeService:     attributeService,
	}
}

func (a *MCPAttributeServiceAdapter) GetNodeAttributes(ctx context.Context, nodeID int) ([]models.NodeAttributeWithInfo, error) {
	response, err := a.nodeAttributeService.GetNodeAttributesByNodeID(nodeID)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (a *MCPAttributeServiceAdapter) SetNodeAttribute(ctx context.Context, nodeID int, req *models.CreateNodeAttributeRequest) (*models.NodeAttributeWithInfo, error) {
	attr, err := a.nodeAttributeService.CreateNodeAttribute(nodeID, req)
	if err != nil {
		return nil, err
	}

	// Convert NodeAttribute to NodeAttributeWithInfo - we need to get additional info
	return &models.NodeAttributeWithInfo{
		ID:          attr.ID,
		NodeID:      attr.NodeID,
		AttributeID: attr.AttributeID,
		Value:       attr.Value,
		OrderIndex:  attr.OrderIndex,
		CreatedAt:   attr.CreatedAt,
	}, nil
}

func (a *MCPAttributeServiceAdapter) GetAttributeByName(ctx context.Context, domainID int, name string) (*models.Attribute, error) {
	// Get all attributes for the domain and find by name
	response, err := a.attributeService.ListAttributes(ctx, domainID)
	if err != nil {
		return nil, err
	}

	for _, attr := range response.Attributes {
		if attr.Name == name {
			return &attr, nil
		}
	}

	return nil, fmt.Errorf("attribute '%s' not found in domain %d", name, domainID)
}

func (a *MCPAttributeServiceAdapter) DeleteNodeAttribute(ctx context.Context, nodeID, attributeID int) error {
	return a.nodeAttributeService.DeleteNodeAttributeByNodeIDAndAttributeID(nodeID, attributeID)
}

// Domain attribute management methods
func (a *MCPAttributeServiceAdapter) ListAttributes(ctx context.Context, domainID int) (*models.AttributeListResponse, error) {
	return a.attributeService.ListAttributes(ctx, domainID)
}

func (a *MCPAttributeServiceAdapter) CreateAttribute(ctx context.Context, domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error) {
	return a.attributeService.CreateAttribute(ctx, domainID, req)
}

func (a *MCPAttributeServiceAdapter) GetAttribute(ctx context.Context, id int) (*models.Attribute, error) {
	return a.attributeService.GetAttribute(ctx, id)
}

func (a *MCPAttributeServiceAdapter) UpdateAttribute(ctx context.Context, id int, req *models.UpdateAttributeRequest) (*models.Attribute, error) {
	return a.attributeService.UpdateAttribute(ctx, id, req)
}

func (a *MCPAttributeServiceAdapter) DeleteAttribute(ctx context.Context, id int) error {
	return a.attributeService.DeleteAttribute(ctx, id)
}

// NodeCountServiceAdapter provides node count functionality
type NodeCountServiceAdapter struct {
	nodeRepo repositories.NodeRepository
}

// NewNodeCountServiceAdapter creates a new NodeCountServiceAdapter
func NewNodeCountServiceAdapter(nodeRepo repositories.NodeRepository) *NodeCountServiceAdapter {
	return &NodeCountServiceAdapter{
		nodeRepo: nodeRepo,
	}
}

func (a *NodeCountServiceAdapter) GetNodeCountByDomain(ctx context.Context, domainID int) (int, error) {
	return a.nodeRepo.CountNodesByDomain(ctx, domainID)
}

// MCPHandlerServiceAdapter adapts the MCP service for HTTP handler usage
type MCPHandlerServiceAdapter struct {
	mcpService mcp.MCPService
}

// NewMCPHandlerServiceAdapter creates a new MCPHandlerServiceAdapter
func NewMCPHandlerServiceAdapter(mcpService mcp.MCPService) *MCPHandlerServiceAdapter {
	return &MCPHandlerServiceAdapter{
		mcpService: mcpService,
	}
}

func (a *MCPHandlerServiceAdapter) CreateMCPNode(req *models.CreateMCPNodeRequest) (*models.MCPNode, error) {
	return a.mcpService.CreateNode(context.Background(), req)
}

func (a *MCPHandlerServiceAdapter) GetMCPNodeByCompositeID(compositeID string) (*models.MCPNode, error) {
	return a.mcpService.GetNode(context.Background(), compositeID)
}

func (a *MCPHandlerServiceAdapter) GetMCPNodes(domainName string, page, size int, search string) (*models.MCPNodeListResponse, error) {
	return a.mcpService.ListNodes(context.Background(), domainName, page, size, search)
}

func (a *MCPHandlerServiceAdapter) UpdateMCPNode(compositeID string, req *models.UpdateMCPNodeRequest) (*models.MCPNode, error) {
	// Convert UpdateMCPNodeRequest to UpdateNodeRequest
	nodeReq := &models.UpdateNodeRequest{
		Title:       req.Title,
		Description: req.Description,
	}
	return a.mcpService.UpdateNode(context.Background(), compositeID, nodeReq)
}

func (a *MCPHandlerServiceAdapter) DeleteMCPNode(compositeID string) error {
	return a.mcpService.DeleteNode(context.Background(), compositeID)
}

func (a *MCPHandlerServiceAdapter) FindMCPNodeByURL(req *models.FindMCPNodeRequest) (*models.MCPNode, error) {
	return a.mcpService.FindNodeByURL(context.Background(), req)
}

func (a *MCPHandlerServiceAdapter) BatchGetMCPNodes(req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error) {
	return a.mcpService.BatchGetNodes(context.Background(), req)
}

func (a *MCPHandlerServiceAdapter) GetMCPDomains() ([]models.MCPDomain, error) {
	response, err := a.mcpService.ListDomains(context.Background())
	if err != nil {
		return nil, err
	}

	// Convert mcp.MCPDomain to models.MCPDomain
	result := make([]models.MCPDomain, len(response.Domains))
	for i, domain := range response.Domains {
		result[i] = models.MCPDomain{
			Name:        domain.Name,
			Description: domain.Description,
			NodeCount:   domain.NodeCount,
			CreatedAt:   domain.CreatedAt,
			UpdatedAt:   domain.UpdatedAt,
		}
	}
	return result, nil
}

func (a *MCPHandlerServiceAdapter) CreateMCPDomain(req *models.CreateMCPDomainRequest) (*models.MCPDomain, error) {
	// Convert CreateMCPDomainRequest to CreateDomainRequest
	domainReq := &models.CreateDomainRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	mcpDomain, err := a.mcpService.CreateDomain(context.Background(), domainReq)
	if err != nil {
		return nil, err
	}

	// Convert mcp.MCPDomain to models.MCPDomain
	return &models.MCPDomain{
		Name:        mcpDomain.Name,
		Description: mcpDomain.Description,
		NodeCount:   mcpDomain.NodeCount,
		CreatedAt:   mcpDomain.CreatedAt,
		UpdatedAt:   mcpDomain.UpdatedAt,
	}, nil
}

func (a *MCPHandlerServiceAdapter) GetMCPNodeAttributes(compositeID string) (*models.MCPNodeAttributesResponse, error) {
	response, err := a.mcpService.GetNodeAttributes(context.Background(), compositeID)
	if err != nil {
		return nil, err
	}

	// Convert MCPNodeAttributeResponse to MCPNodeAttributesResponse
	return &models.MCPNodeAttributesResponse{
		CompositeID: response.CompositeID,
		Attributes:  response.Attributes,
	}, nil
}

func (a *MCPHandlerServiceAdapter) SetMCPNodeAttributes(compositeID string, req *models.SetMCPNodeAttributesRequest) (*models.MCPNodeAttributesResponse, error) {
	response, err := a.mcpService.SetNodeAttributes(context.Background(), compositeID, req)
	if err != nil {
		return nil, err
	}

	// Convert MCPNodeAttributeResponse to MCPNodeAttributesResponse
	return &models.MCPNodeAttributesResponse{
		CompositeID: response.CompositeID,
		Attributes:  response.Attributes,
	}, nil
}

func (a *MCPHandlerServiceAdapter) GetMCPServerInfo() (*models.MCPServerInfo, error) {
	serverInfo, err := a.mcpService.GetServerInfo(context.Background())
	if err != nil {
		return nil, err
	}

	// Convert mcp.MCPServerInfo to models.MCPServerInfo
	return &models.MCPServerInfo{
		Name:               serverInfo.Name,
		Version:            serverInfo.Version,
		Description:        serverInfo.Description,
		Capabilities:       serverInfo.Capabilities,
		CompositeKeyFormat: serverInfo.CompositeKeyFormat,
	}, nil
}

// Domain attribute management methods
func (a *MCPHandlerServiceAdapter) ListDomainAttributes(domainName string) (*models.MCPDomainAttributeListResponse, error) {
	return a.mcpService.ListDomainAttributes(context.Background(), domainName)
}

func (a *MCPHandlerServiceAdapter) CreateDomainAttribute(domainName string, req *models.CreateAttributeRequest) (*models.MCPDomainAttribute, error) {
	return a.mcpService.CreateDomainAttribute(context.Background(), domainName, req)
}

func (a *MCPHandlerServiceAdapter) GetDomainAttribute(compositeID string) (*models.MCPDomainAttribute, error) {
	return a.mcpService.GetDomainAttribute(context.Background(), compositeID)
}

func (a *MCPHandlerServiceAdapter) UpdateDomainAttribute(compositeID string, req *models.UpdateAttributeRequest) (*models.MCPDomainAttribute, error) {
	return a.mcpService.UpdateDomainAttribute(context.Background(), compositeID, req)
}

func (a *MCPHandlerServiceAdapter) DeleteDomainAttribute(compositeID string) error {
	return a.mcpService.DeleteDomainAttribute(context.Background(), compositeID)
}