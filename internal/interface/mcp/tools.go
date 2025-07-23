package mcp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

// Additional Node Management Tools

// handleGetNode implements the get_node tool
func (h *MCPToolHandler) handleGetNode(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse composite_id argument
	compositeID, ok := args["composite_id"].(string)
	if !ok || compositeID == "" {
		return nil, fmt.Errorf("missing or invalid 'composite_id' parameter")
	}

	// Parse composite ID to extract node ID
	// composite_id format: "tool-name:domain:id"
	parts := strings.Split(compositeID, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid composite_id format, expected 'tool-name:domain:id'")
	}

	nodeIDStr := parts[2]
	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID in composite_id: %v", err)
	}

	// Get node from repository
	node, err := h.dependencies.NodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Node ID: %d\nURL: %s\nTitle: %s\nDescription: %s\nCreated: %s\nUpdated: %s",
					node.ID(), node.URL(), node.Title(), node.Description(),
					node.CreatedAt().Format("2006-01-02 15:04:05"),
					node.UpdatedAt().Format("2006-01-02 15:04:05")),
			},
		},
		"isError": false,
	}, nil
}

// handleUpdateNode implements the update_node tool
func (h *MCPToolHandler) handleUpdateNode(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse composite_id argument
	compositeID, ok := args["composite_id"].(string)
	if !ok || compositeID == "" {
		return nil, fmt.Errorf("missing or invalid 'composite_id' parameter")
	}

	// Parse composite ID to extract node ID
	parts := strings.Split(compositeID, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid composite_id format, expected 'tool-name:domain:id'")
	}

	nodeIDStr := parts[2]
	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID in composite_id: %v", err)
	}

	// Get existing node
	node, err := h.dependencies.NodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// Update fields if provided
	updated := false
	if title, ok := args["title"].(string); ok {
		if err := node.UpdateTitle(title); err != nil {
			return nil, fmt.Errorf("failed to update title: %w", err)
		}
		updated = true
	}

	if description, ok := args["description"].(string); ok {
		if err := node.UpdateDescription(description); err != nil {
			return nil, fmt.Errorf("failed to update description: %w", err)
		}
		updated = true
	}

	if !updated {
		return nil, fmt.Errorf("at least one field (title or description) must be provided for update")
	}

	// Save updated node
	if err := h.dependencies.NodeRepo.Update(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully updated node:\nID: %d\nURL: %s\nTitle: %s\nDescription: %s\nUpdated: %s",
					node.ID(), node.URL(), node.Title(), node.Description(),
					node.UpdatedAt().Format("2006-01-02 15:04:05")),
			},
		},
		"isError": false,
	}, nil
}

// handleDeleteNode implements the delete_node tool
func (h *MCPToolHandler) handleDeleteNode(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse composite_id argument
	compositeID, ok := args["composite_id"].(string)
	if !ok || compositeID == "" {
		return nil, fmt.Errorf("missing or invalid 'composite_id' parameter")
	}

	// Parse composite ID to extract node ID
	parts := strings.Split(compositeID, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid composite_id format, expected 'tool-name:domain:id'")
	}

	nodeIDStr := parts[2]
	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID in composite_id: %v", err)
	}

	// Get node before deleting (for confirmation message)
	node, err := h.dependencies.NodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// Delete the node
	if err := h.dependencies.NodeRepo.Delete(ctx, nodeID); err != nil {
		return nil, fmt.Errorf("failed to delete node: %w", err)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully deleted node:\nID: %d\nURL: %s\nTitle: %s",
					node.ID(), node.URL(), node.Title()),
			},
		},
		"isError": false,
	}, nil
}

// handleFindNodeByURL implements the find_node_by_url tool
func (h *MCPToolHandler) handleFindNodeByURL(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	url, ok := args["url"].(string)
	if !ok || url == "" {
		return nil, fmt.Errorf("missing or invalid 'url' parameter")
	}

	// Find node by URL
	node, err := h.dependencies.NodeRepo.GetByURL(ctx, url, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to find node: %w", err)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Found node:\nID: %d\nURL: %s\nTitle: %s\nDescription: %s\nCreated: %s",
					node.ID(), node.URL(), node.Title(), node.Description(),
					node.CreatedAt().Format("2006-01-02 15:04:05")),
			},
		},
		"isError": false,
	}, nil
}

// Attribute Management Tools

// handleGetNodeAttributes implements the get_node_attributes tool
func (h *MCPToolHandler) handleGetNodeAttributes(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse composite_id argument
	compositeID, ok := args["composite_id"].(string)
	if !ok || compositeID == "" {
		return nil, fmt.Errorf("missing or invalid 'composite_id' parameter")
	}

	// Parse composite ID to extract node ID
	parts := strings.Split(compositeID, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid composite_id format, expected 'tool-name:domain:id'")
	}

	nodeIDStr := parts[2]
	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID in composite_id: %v", err)
	}

	// Get node to ensure it exists
	node, err := h.dependencies.NodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// TODO: Get node attributes from database
	// For now, return a placeholder response
	content := []map[string]interface{}{
		{
			"type": "text",
			"text": fmt.Sprintf("Node attributes for: %s\nURL: %s\n\nNote: Attribute functionality is not fully implemented yet",
				node.Title(), node.URL()),
		},
	}

	return map[string]interface{}{
		"content": content,
		"isError": false,
	}, nil
}

// handleSetNodeAttributes implements the set_node_attributes tool
func (h *MCPToolHandler) handleSetNodeAttributes(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse composite_id argument
	compositeID, ok := args["composite_id"].(string)
	if !ok || compositeID == "" {
		return nil, fmt.Errorf("missing or invalid 'composite_id' parameter")
	}

	// Parse attributes argument
	attributesRaw, ok := args["attributes"]
	if !ok {
		return nil, fmt.Errorf("missing 'attributes' parameter")
	}

	attributes, ok := attributesRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid 'attributes' parameter, expected array")
	}

	// Parse composite ID to extract node ID
	parts := strings.Split(compositeID, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid composite_id format, expected 'tool-name:domain:id'")
	}

	nodeIDStr := parts[2]
	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID in composite_id: %v", err)
	}

	// Get node to ensure it exists
	node, err := h.dependencies.NodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// TODO: Implement attribute setting logic
	// For now, return a placeholder response
	attributeCount := len(attributes)

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Would set %d attributes for node: %s\nURL: %s\n\nNote: Attribute setting functionality is not fully implemented yet",
					attributeCount, node.Title(), node.URL()),
			},
		},
		"isError": false,
	}, nil
}

// Domain Schema Management Tools

// handleListDomainAttributes implements the list_domain_attributes tool
func (h *MCPToolHandler) handleListDomainAttributes(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse domain_name argument
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	// Get domain first to get domain ID
	domain, err := h.dependencies.DomainRepo.GetByName(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// Get attributes for this domain
	attributes, err := h.dependencies.AttributeRepo.ListByDomainID(ctx, domain.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to list domain attributes: %w", err)
	}

	// Convert to MCP response format
	content := []map[string]interface{}{}
	for _, attr := range attributes {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("Attribute: %s\nType: %s\nDescription: %s\nCreated: %s",
				attr.Name(), attr.Type(), attr.Description(),
				attr.CreatedAt().Format("2006-01-02 15:04:05")),
		})
	}

	if len(content) == 0 {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("No attributes defined for domain '%s'", domainName),
		})
	}

	return map[string]interface{}{
		"content": content,
		"isError": false,
	}, nil
}

// handleCreateDomainAttribute implements the create_domain_attribute tool
func (h *MCPToolHandler) handleCreateDomainAttribute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("missing or invalid 'name' parameter")
	}

	attrType, ok := args["type"].(string)
	if !ok || attrType == "" {
		return nil, fmt.Errorf("missing or invalid 'type' parameter")
	}

	description := ""
	if d, ok := args["description"].(string); ok {
		description = d
	}

	// Get domain first to get domain ID
	domain, err := h.dependencies.DomainRepo.GetByName(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// Create attribute request DTO  
	createReq := &request.CreateAttributeRequest{
		DomainID:    domain.ID(),
		Name:        name,
		Type:        attrType,
		Description: description,
	}

	// Execute use case
	result, err := h.dependencies.CreateAttributeUC.Execute(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain attribute: %w", err)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully created domain attribute:\nDomain: %s\nName: %s\nType: %s\nDescription: %s\nCreated: %s",
					domainName, result.Name, result.Type, result.Description,
					result.CreatedAt.Format("2006-01-02 15:04:05")),
			},
		},
		"isError": false,
	}, nil
}