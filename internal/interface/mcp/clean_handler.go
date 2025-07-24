package mcp

import (
	"context"
	"fmt"
	"url-db/internal/application/usecase/mcp"
	"url-db/internal/interface/setup"
)

// CleanMCPToolHandler implements MCP tools following Clean Architecture principles
type CleanMCPToolHandler struct {
	domainUseCases *mcp.DomainUseCases
	nodeUseCases   *mcp.NodeUseCases
}

// NewCleanMCPToolHandler creates a new clean MCP tool handler
func NewCleanMCPToolHandler(factory *setup.ApplicationFactory) *CleanMCPToolHandler {
	// Create domain use cases
	createDomainUC, listDomainsUC := factory.CreateDomainUseCases(factory.CreateDomainRepository())
	domainUseCases := mcp.NewDomainUseCases(createDomainUC, listDomainsUC)

	// Create node use cases
	createNodeUC, listNodesUC := factory.CreateNodeUseCases(factory.CreateNodeRepository(), factory.CreateDomainRepository())
	nodeUseCases := mcp.NewNodeUseCases(createNodeUC, listNodesUC)

	return &CleanMCPToolHandler{
		domainUseCases: domainUseCases,
		nodeUseCases:   nodeUseCases,
	}
}

// Domain Management Tools

// HandleListDomains implements the list_domains tool with clean architecture
func (h *CleanMCPToolHandler) HandleListDomains(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Optional pagination parameters
	page := 1
	if p, ok := args["page"].(float64); ok {
		page = int(p)
	}

	size := 20
	if s, ok := args["size"].(float64); ok {
		size = int(s)
	}

	// Use MCP-specific use case
	mcpUseCase := mcp.NewMCPListDomainsUseCase(h.domainUseCases.List)
	return mcpUseCase.Execute(ctx, page, size)
}

// HandleCreateDomain implements the create_domain tool with clean architecture
func (h *CleanMCPToolHandler) HandleCreateDomain(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("missing or invalid 'name' parameter")
	}

	description, ok := args["description"].(string)
	if !ok || description == "" {
		return nil, fmt.Errorf("missing or invalid 'description' parameter")
	}

	// Use MCP-specific use case
	mcpUseCase := mcp.NewMCPCreateDomainUseCase(h.domainUseCases.Create)
	return mcpUseCase.Execute(ctx, name, description)
}

// Node Management Tools

// HandleListNodes implements the list_nodes tool with clean architecture
func (h *CleanMCPToolHandler) HandleListNodes(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse required arguments
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	// Optional pagination parameters
	page := 1
	if p, ok := args["page"].(float64); ok {
		page = int(p)
	}

	size := 20
	if s, ok := args["size"].(float64); ok {
		size = int(s)
	}

	// Use MCP-specific use case
	mcpUseCase := mcp.NewMCPListNodesUseCase(h.nodeUseCases.List)
	return mcpUseCase.Execute(ctx, domainName, page, size)
}

// HandleCreateNode implements the create_node tool with clean architecture
func (h *CleanMCPToolHandler) HandleCreateNode(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse required arguments
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	url, ok := args["url"].(string)
	if !ok || url == "" {
		return nil, fmt.Errorf("missing or invalid 'url' parameter")
	}

	// Optional arguments
	title := ""
	if t, ok := args["title"].(string); ok {
		title = t
	}

	description := ""
	if d, ok := args["description"].(string); ok {
		description = d
	}

	// Use MCP-specific use case
	mcpUseCase := mcp.NewMCPCreateNodeUseCase(h.nodeUseCases.Create)
	return mcpUseCase.Execute(ctx, domainName, url, title, description)
}

// GetToolHandler returns the appropriate tool handler based on tool name
func (h *CleanMCPToolHandler) GetToolHandler(toolName string) func(context.Context, map[string]interface{}) (interface{}, error) {
	handlers := map[string]func(context.Context, map[string]interface{}) (interface{}, error){
		"list_domains":  h.HandleListDomains,
		"create_domain": h.HandleCreateDomain,
		"list_nodes":    h.HandleListNodes,
		"create_node":   h.HandleCreateNode,
	}

	return handlers[toolName]
}
