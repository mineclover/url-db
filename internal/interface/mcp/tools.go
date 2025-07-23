package mcp

import (
	"context"
	"fmt"

	"url-db/internal/application/dto/request"
	"url-db/internal/interface/setup"
)

// MCPToolHandler handles all MCP tool implementations
type MCPToolHandler struct {
	dependencies *setup.CleanDependencies
}

// NewMCPToolHandler creates a new tool handler
func NewMCPToolHandler(factory *setup.ApplicationFactory) *MCPToolHandler {
	return &MCPToolHandler{
		dependencies: factory.CreateCleanArchitectureDependencies(),
	}
}

// Domain Management Tools

// handleListDomains implements the list_domains tool
func (h *MCPToolHandler) handleListDomains(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Optional pagination parameters
	page := 1
	if p, ok := args["page"].(float64); ok {
		page = int(p)
	}

	size := 20
	if s, ok := args["size"].(float64); ok {
		size = int(s)
	}

	result, err := h.dependencies.ListDomainsUC.Execute(ctx, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	// Convert to MCP response format
	content := []map[string]interface{}{}
	for _, domain := range result.Domains {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("Domain: %s\nDescription: %s\nCreated: %s",
				domain.Name, domain.Description, domain.CreatedAt.Format("2006-01-02 15:04:05")),
		})
	}

	if len(content) == 0 {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": "No domains found",
		})
	}

	return map[string]interface{}{
		"content": content,
		"isError": false,
	}, nil
}

// handleCreateDomain implements the create_domain tool
func (h *MCPToolHandler) handleCreateDomain(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("missing or invalid 'name' parameter")
	}

	description, ok := args["description"].(string)
	if !ok || description == "" {
		return nil, fmt.Errorf("missing or invalid 'description' parameter")
	}

	// Create request DTO
	createReq := &request.CreateDomainRequest{
		Name:        name,
		Description: description,
	}

	// Execute use case
	result, err := h.dependencies.CreateDomainUC.Execute(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully created domain: %s\nDescription: %s\nCreated: %s",
					result.Name, result.Description, result.CreatedAt.Format("2006-01-02 15:04:05")),
			},
		},
		"isError": false,
	}, nil
}

// Node Management Tools

// handleListNodes implements the list_nodes tool
func (h *MCPToolHandler) handleListNodes(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	// Optional parameters with defaults
	page := 1
	if p, ok := args["page"].(float64); ok {
		page = int(p)
	}

	size := 20
	if s, ok := args["size"].(float64); ok {
		size = int(s)
	}

	search := ""
	if s, ok := args["search"].(string); ok {
		search = s
	}
	_ = search // TODO: Implement search functionality

	// Execute use case
	result, err := h.dependencies.ListNodesUC.Execute(ctx, domainName, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	// Convert to MCP response format
	content := []map[string]interface{}{}
	for _, node := range result.Nodes {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("Node ID: %d\nURL: %s\nTitle: %s\nDescription: %s\nCreated: %s",
				node.ID, node.URL, node.Title, node.Description, node.CreatedAt.Format("2006-01-02 15:04:05")),
		})
	}

	if len(content) == 0 {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("No nodes found in domain '%s'", domainName),
		})
	}

	return map[string]interface{}{
		"content": content,
		"isError": false,
	}, nil
}

// handleCreateNode implements the create_node tool
func (h *MCPToolHandler) handleCreateNode(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse required arguments
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	url, ok := args["url"].(string)
	if !ok || url == "" {
		return nil, fmt.Errorf("missing or invalid 'url' parameter")
	}

	// Optional parameters
	title := ""
	if t, ok := args["title"].(string); ok {
		title = t
	}

	description := ""
	if d, ok := args["description"].(string); ok {
		description = d
	}

	// Create request DTO
	createReq := &request.CreateNodeRequest{
		DomainName:  domainName,
		URL:         url,
		Title:       title,
		Description: description,
	}

	// Execute use case
	result, err := h.dependencies.CreateNodeUC.Execute(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully created node in domain '%s'\nURL: %s\nTitle: %s\nDescription: %s\nCreated: %s",
					domainName, result.URL, result.Title, result.Description, result.CreatedAt.Format("2006-01-02 15:04:05")),
			},
		},
		"isError": false,
	}, nil
}