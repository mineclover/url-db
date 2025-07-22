package mcp

import (
	"context"

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
	FilterNodesByAttributes(ctx context.Context, domainName string, filters []interface{}, page, size int) (*models.MCPNodeListResponse, error)

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

// Node operations are in node_operations.go
// Domain operations are in domain_operations.go
// Attribute operations are in attribute_operations.go
// Query operations are in query_operations.go
