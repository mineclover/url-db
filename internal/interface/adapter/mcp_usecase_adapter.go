package adapter

import (
	"context"
	"fmt"
	
	"url-db/internal/application/usecase/domain"
	"url-db/internal/application/usecase/node"
	"url-db/internal/application/usecase/attribute"
	"url-db/internal/application/dto/request"
	"url-db/internal/interfaces/mcp"
	"url-db/internal/models"
)

// MCPUseCaseAdapter directly adapts use cases to MCP interface
// This simplifies the adapter chain by removing intermediate layers
type MCPUseCaseAdapter struct {
	// Domain use cases
	createDomainUC *domain.CreateDomainUseCase
	listDomainsUC  *domain.ListDomainsUseCase
	
	// Node use cases
	createNodeUC *node.CreateNodeUseCase
	listNodesUC  *node.ListNodesUseCase
	
	// Attribute use cases
	createAttributeUC *attribute.CreateAttributeUseCase
	listAttributesUC  *attribute.ListAttributesUseCase
	
	// MCP converter for composite keys
	converter *mcp.Converter
}

// NewMCPUseCaseAdapter creates a new simplified MCP adapter
func NewMCPUseCaseAdapter(
	createDomainUC *domain.CreateDomainUseCase,
	listDomainsUC *domain.ListDomainsUseCase,
	createNodeUC *node.CreateNodeUseCase,
	listNodesUC *node.ListNodesUseCase,
	createAttributeUC *attribute.CreateAttributeUseCase,
	listAttributesUC *attribute.ListAttributesUseCase,
	converter *mcp.Converter,
) *MCPUseCaseAdapter {
	return &MCPUseCaseAdapter{
		createDomainUC:    createDomainUC,
		listDomainsUC:     listDomainsUC,
		createNodeUC:      createNodeUC,
		listNodesUC:       listNodesUC,
		createAttributeUC: createAttributeUC,
		listAttributesUC:  listAttributesUC,
		converter:         converter,
	}
}

// Domain operations
func (a *MCPUseCaseAdapter) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*models.MCPDomain, error) {
	// Convert to use case request
	ucReq := &request.CreateDomainRequest{
		Name:        req.Name,
		Description: req.Description,
	}
	
	// Execute use case
	result, err := a.createDomainUC.Execute(ctx, ucReq)
	if err != nil {
		return nil, err
	}
	
	// Convert to MCP response
	return &models.MCPDomain{
		Name:        result.Name,
		Description: result.Description,
		NodeCount:   0, // Will be calculated separately if needed
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.CreatedAt, // For newly created domain
	}, nil
}

func (a *MCPUseCaseAdapter) ListDomains(ctx context.Context) (*models.MCPDomainListResponse, error) {
	// Execute use case with reasonable defaults
	result, err := a.listDomainsUC.Execute(ctx, 1, 100)
	if err != nil {
		return nil, err
	}
	
	// Convert to MCP response
	mcpDomains := make([]models.MCPDomain, len(result.Domains))
	for i, d := range result.Domains {
		mcpDomains[i] = models.MCPDomain{
			Name:        d.Name,
			Description: d.Description,
			NodeCount:   0, // Will be calculated separately if needed
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		}
	}
	
	return &models.MCPDomainListResponse{
		Domains: mcpDomains,
	}, nil
}

// Node operations
func (a *MCPUseCaseAdapter) CreateNode(ctx context.Context, req *models.CreateMCPNodeRequest) (*models.MCPNode, error) {
	// For now, we need to get domain ID from name
	// This is a simplified implementation
	_ = 1 // domainID would be looked up from domainName
	domainName := req.DomainName
	
	// Convert to use case request
	ucReq := &request.CreateNodeRequest{
		URL:         req.URL,
		Title:       req.Title,
		Description: req.Description,
	}
	
	// Execute use case
	result, err := a.createNodeUC.Execute(ctx, ucReq)
	if err != nil {
		return nil, err
	}
	
	// Convert to MCP response
	// For now using simple format - this should use proper composite key format
	compositeID := fmt.Sprintf("url-db:%s:%d", domainName, result.ID)
	
	return &models.MCPNode{
		CompositeID: compositeID,
		URL:         result.URL,
		Title:       result.Title,
		Description: result.Description,
		DomainName:  domainName,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

// Attribute operations
func (a *MCPUseCaseAdapter) CreateDomainAttribute(ctx context.Context, domainName string, req *models.CreateAttributeRequest) (*models.MCPDomainAttribute, error) {
	// For now, we need to get domain ID from name
	// This is a simplified implementation
	domainID := 1 // This should be looked up properly
	
	// Convert to use case request
	ucReq := &request.CreateAttributeRequest{
		Name:        req.Name,
		Type:        string(req.Type),
		Description: req.Description,
		DomainID:    domainID,
	}
	
	// Execute use case
	result, err := a.createAttributeUC.Execute(ctx, ucReq)
	if err != nil {
		return nil, err
	}
	
	// Convert to MCP response
	// For now using simple format - this should use proper composite key format
	compositeID := fmt.Sprintf("url-db:%s:%d", domainName, result.ID)
	
	return &models.MCPDomainAttribute{
		CompositeID: compositeID,
		Name:        result.Name,
		Type:        models.AttributeType(result.Type),
		Description: result.Description,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (a *MCPUseCaseAdapter) ListDomainAttributes(ctx context.Context, domainName string) (*models.MCPDomainAttributeListResponse, error) {
	// For now, we need to get domain ID from name
	// This is a simplified implementation
	domainID := 1 // This should be looked up properly
	
	// Execute use case
	result, err := a.listAttributesUC.Execute(ctx, domainID)
	if err != nil {
		return nil, err
	}
	
	// Convert to MCP response
	mcpAttributes := make([]models.MCPDomainAttribute, len(result.Attributes))
	for i, attr := range result.Attributes {
		// For now using simple format - this should use proper composite key format
		compositeID := fmt.Sprintf("url-db:%s:%d", domainName, attr.ID)
		mcpAttributes[i] = models.MCPDomainAttribute{
			CompositeID: compositeID,
			Name:        attr.Name,
			Type:        models.AttributeType(attr.Type),
			Description: attr.Description,
			CreatedAt:   attr.CreatedAt,
			UpdatedAt:   attr.UpdatedAt,
		}
	}
	
	return &models.MCPDomainAttributeListResponse{
		DomainName: domainName,
		Attributes: mcpAttributes,
		TotalCount: result.Total,
	}, nil
}

// Server info
func (a *MCPUseCaseAdapter) GetServerInfo(ctx context.Context) (*models.MCPServerInfo, error) {
	return &models.MCPServerInfo{
		Name:               "url-db",
		Version:            "1.0.0",
		Description:        "URL Database with MCP support",
		Capabilities:       []string{"domains", "nodes", "attributes"},
		CompositeKeyFormat: "tool-name:domain:id",
	}, nil
}