package services

import (
	"context"

	"github.com/url-db/internal/models"
)

type DomainService interface {
	CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*models.Domain, error)
	GetDomain(ctx context.Context, id int) (*models.Domain, error)
	GetDomainByName(ctx context.Context, name string) (*models.Domain, error)
	ListDomains(ctx context.Context, page, size int) (*models.DomainListResponse, error)
	UpdateDomain(ctx context.Context, id int, req *models.UpdateDomainRequest) (*models.Domain, error)
	DeleteDomain(ctx context.Context, id int) error
}

type NodeService interface {
	CreateNode(ctx context.Context, domainID int, req *models.CreateNodeRequest) (*models.Node, error)
	GetNode(ctx context.Context, id int) (*models.Node, error)
	GetNodeByDomainAndURL(ctx context.Context, domainID int, url string) (*models.Node, error)
	ListNodesByDomain(ctx context.Context, domainID int, page, size int, search string) (*models.NodeListResponse, error)
	UpdateNode(ctx context.Context, id int, req *models.UpdateNodeRequest) (*models.Node, error)
	DeleteNode(ctx context.Context, id int) error
	FindNodeByURL(ctx context.Context, domainID int, req *models.FindNodeByURLRequest) (*models.Node, error)
}

type AttributeService interface {
	CreateAttribute(ctx context.Context, domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error)
	GetAttribute(ctx context.Context, id int) (*models.Attribute, error)
	ListAttributesByDomain(ctx context.Context, domainID int) (*models.AttributeListResponse, error)
	UpdateAttribute(ctx context.Context, id int, req *models.UpdateAttributeRequest) (*models.Attribute, error)
	DeleteAttribute(ctx context.Context, id int) error
	ValidateAttributeValue(ctx context.Context, attributeID int, value string) error
}

type NodeAttributeService interface {
	CreateNodeAttribute(ctx context.Context, nodeID int, req *models.CreateNodeAttributeRequest) (*models.NodeAttribute, error)
	GetNodeAttribute(ctx context.Context, id int) (*models.NodeAttribute, error)
	ListNodeAttributesByNode(ctx context.Context, nodeID int) (*models.NodeAttributeListResponse, error)
	UpdateNodeAttribute(ctx context.Context, id int, req *models.UpdateNodeAttributeRequest) (*models.NodeAttribute, error)
	DeleteNodeAttribute(ctx context.Context, id int) error
	ValidateNodeAttributeValue(ctx context.Context, nodeID, attributeID int, value string) error
}

type CompositeKeyService interface {
	Create(domainName string, id int) string
	Parse(compositeKey string) (*models.CompositeKey, error)
	Validate(compositeKey string) error
	GetToolName() string
}

type MCPService interface {
	CreateNode(ctx context.Context, req *models.CreateMCPNodeRequest) (*models.MCPNode, error)
	GetNode(ctx context.Context, compositeID string) (*models.MCPNode, error)
	ListNodes(ctx context.Context, domainName string, page, size int, search string) (*models.MCPNodeListResponse, error)
	UpdateNode(ctx context.Context, compositeID string, req *models.UpdateMCPNodeRequest) (*models.MCPNode, error)
	DeleteNode(ctx context.Context, compositeID string) error
	FindNodeByURL(ctx context.Context, req *models.FindMCPNodeRequest) (*models.MCPNode, error)
	BatchGetNodes(ctx context.Context, req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error)
	
	// 도메인 관리
	ListDomains(ctx context.Context) (*models.MCPDomainListResponse, error)
	CreateDomain(ctx context.Context, req *models.CreateMCPDomainRequest) (*models.MCPDomain, error)
	
	// 속성 관리
	GetNodeAttributes(ctx context.Context, compositeID string) (*models.MCPNodeAttributeResponse, error)
	SetNodeAttributes(ctx context.Context, compositeID string, req *models.SetMCPNodeAttributesRequest) error
	
	// 서버 정보
	GetServerInfo(ctx context.Context) (*models.MCPServerInfo, error)
}