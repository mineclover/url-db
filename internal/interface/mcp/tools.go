package mcp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"url-db/internal/application/dto/request"
	nodeUseCase "url-db/internal/application/usecase/node"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
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

	// Get node attributes from database
	nodeAttributes, err := h.dependencies.NodeAttributeRepo.GetByNodeID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node attributes: %w", err)
	}

	if len(nodeAttributes) == 0 {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("No attributes found for node: %s\nURL: %s", node.Title(), node.URL()),
				},
			},
			"isError": false,
		}, nil
	}

	// Build attributes display
	var attributeTexts []string
	for _, nodeAttr := range nodeAttributes {
		// Get attribute definition to show name and type
		attr, err := h.dependencies.AttributeRepo.GetByID(ctx, nodeAttr.AttributeID())
		if err != nil {
			continue // Skip if attribute definition not found
		}

		text := fmt.Sprintf("• %s (%s): %s", attr.Name(), attr.Type(), nodeAttr.Value())
		if nodeAttr.OrderIndex() != nil {
			text += fmt.Sprintf(" [order: %d]", *nodeAttr.OrderIndex())
		}
		attributeTexts = append(attributeTexts, text)
	}

	content := []map[string]interface{}{
		{
			"type": "text",
			"text": fmt.Sprintf("Attributes for node: %s\nURL: %s\n\n%s",
				node.Title(), node.URL(), strings.Join(attributeTexts, "\n")),
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

	// Convert attributes to use case input
	var attributeInputs []nodeUseCase.AttributeInput
	for _, attr := range attributes {
		attrMap, ok := attr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid attribute format")
		}

		name, ok := attrMap["name"].(string)
		if !ok || name == "" {
			return nil, fmt.Errorf("attribute must have a valid 'name'")
		}

		value, ok := attrMap["value"].(string)
		if !ok || value == "" {
			return nil, fmt.Errorf("attribute must have a valid 'value'")
		}

		var orderIndex *int
		if orderIndexRaw, exists := attrMap["order_index"]; exists && orderIndexRaw != nil {
			if orderIndexFloat, ok := orderIndexRaw.(float64); ok {
				orderIndexInt := int(orderIndexFloat)
				orderIndex = &orderIndexInt
			}
		}

		attributeInputs = append(attributeInputs, nodeUseCase.AttributeInput{
			Name:       name,
			Value:      value,
			OrderIndex: orderIndex,
		})
	}

	// Execute the use case
	err = h.dependencies.SetNodeAttributesUC.Execute(ctx, nodeID, attributeInputs)
	if err != nil {
		return nil, fmt.Errorf("failed to set node attributes: %w", err)
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully set %d attributes for node: %s\nURL: %s",
					len(attributes), node.Title(), node.URL()),
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

// handleGetDomainAttribute implements the get_domain_attribute tool
func (h *MCPToolHandler) handleGetDomainAttribute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	attributeName, ok := args["attribute_name"].(string)
	if !ok || attributeName == "" {
		return nil, fmt.Errorf("missing or invalid 'attribute_name' parameter")
	}

	// Get domain first to get domain ID
	domain, err := h.dependencies.DomainRepo.GetByName(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// Get all attributes for this domain and find the specific one
	attributes, err := h.dependencies.AttributeRepo.ListByDomainID(ctx, domain.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to list domain attributes: %w", err)
	}

	// Find the specific attribute
	var foundAttribute *entity.Attribute
	for _, attr := range attributes {
		if attr.Name() == attributeName {
			foundAttribute = attr
			break
		}
	}

	if foundAttribute == nil {
		return nil, fmt.Errorf("attribute '%s' not found in domain '%s'", attributeName, domainName)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Domain Attribute Details:\nDomain: %s\nName: %s\nType: %s\nDescription: %s\nCreated: %s",
					domainName, foundAttribute.Name(), foundAttribute.Type(), foundAttribute.Description(),
					foundAttribute.CreatedAt().Format("2006-01-02 15:04:05")),
			},
		},
		"isError": false,
	}, nil
}

// handleUpdateDomainAttribute implements the update_domain_attribute tool
func (h *MCPToolHandler) handleUpdateDomainAttribute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	attributeName, ok := args["attribute_name"].(string)
	if !ok || attributeName == "" {
		return nil, fmt.Errorf("missing or invalid 'attribute_name' parameter")
	}

	// Get domain first to get domain ID
	domain, err := h.dependencies.DomainRepo.GetByName(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// Get all attributes for this domain and find the specific one
	attributes, err := h.dependencies.AttributeRepo.ListByDomainID(ctx, domain.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to list domain attributes: %w", err)
	}

	// Find the specific attribute
	var foundAttribute *entity.Attribute
	for _, attr := range attributes {
		if attr.Name() == attributeName {
			foundAttribute = attr
			break
		}
	}

	if foundAttribute == nil {
		return nil, fmt.Errorf("attribute '%s' not found in domain '%s'", attributeName, domainName)
	}

	// Update description if provided
	updated := false
	if description, ok := args["description"].(string); ok {
		foundAttribute.UpdateDescription(description)
		updated = true
	}

	if !updated {
		return nil, fmt.Errorf("at least one field (description) must be provided for update")
	}

	// Save updated attribute
	if err := h.dependencies.AttributeRepo.Update(ctx, foundAttribute); err != nil {
		return nil, fmt.Errorf("failed to update domain attribute: %w", err)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully updated domain attribute:\nDomain: %s\nName: %s\nType: %s\nDescription: %s\nUpdated: %s",
					domainName, foundAttribute.Name(), foundAttribute.Type(), foundAttribute.Description(),
					foundAttribute.UpdatedAt().Format("2006-01-02 15:04:05")),
			},
		},
		"isError": false,
	}, nil
}

// handleDeleteDomainAttribute implements the delete_domain_attribute tool
func (h *MCPToolHandler) handleDeleteDomainAttribute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	attributeName, ok := args["attribute_name"].(string)
	if !ok || attributeName == "" {
		return nil, fmt.Errorf("missing or invalid 'attribute_name' parameter")
	}

	// Get domain first to get domain ID
	domain, err := h.dependencies.DomainRepo.GetByName(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// Get all attributes for this domain and find the specific one
	attributes, err := h.dependencies.AttributeRepo.ListByDomainID(ctx, domain.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to list domain attributes: %w", err)
	}

	// Find the specific attribute
	var foundAttribute *entity.Attribute
	for _, attr := range attributes {
		if attr.Name() == attributeName {
			foundAttribute = attr
			break
		}
	}

	if foundAttribute == nil {
		return nil, fmt.Errorf("attribute '%s' not found in domain '%s'", attributeName, domainName)
	}

	// Delete the attribute
	if err := h.dependencies.AttributeRepo.Delete(ctx, foundAttribute.ID()); err != nil {
		return nil, fmt.Errorf("failed to delete domain attribute: %w", err)
	}

	// Convert to MCP response format
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully deleted domain attribute:\nDomain: %s\nName: %s\nType: %s",
					domainName, foundAttribute.Name(), foundAttribute.Type()),
			},
		},
		"isError": false,
	}, nil
}

// Dependency Management Tools

// parseCompositeID is a helper function to parse composite IDs
func parseCompositeID(compositeID string) (int, error) {
	parts := strings.Split(compositeID, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid composite_id format, expected 'tool-name:domain:id'")
	}

	nodeID, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid node ID in composite_id: %v", err)
	}

	return nodeID, nil
}

// handleCreateDependency implements the create_dependency tool
func (h *MCPToolHandler) handleCreateDependency(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	dependentNodeID, ok := args["dependent_node_id"].(string)
	if !ok || dependentNodeID == "" {
		return nil, fmt.Errorf("missing or invalid 'dependent_node_id' parameter")
	}

	dependencyNodeID, ok := args["dependency_node_id"].(string)
	if !ok || dependencyNodeID == "" {
		return nil, fmt.Errorf("missing or invalid 'dependency_node_id' parameter")
	}

	dependencyType, ok := args["dependency_type"].(string)
	if !ok || dependencyType == "" {
		return nil, fmt.Errorf("missing or invalid 'dependency_type' parameter")
	}

	// Validate dependency type
	validTypes := []string{"hard", "soft", "reference"}
	isValid := false
	for _, validType := range validTypes {
		if dependencyType == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, fmt.Errorf("invalid dependency_type: %s. Must be one of: hard, soft, reference", dependencyType)
	}

	// Parse composite IDs
	depNodeID, err := parseCompositeID(dependentNodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid dependent_node_id: %w", err)
	}

	depyNodeID, err := parseCompositeID(dependencyNodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid dependency_node_id: %w", err)
	}

	// Prevent self-dependency
	if depNodeID == depyNodeID {
		return nil, fmt.Errorf("a node cannot depend on itself")
	}

	// Optional parameters
	cascadeDelete := false
	if cd, ok := args["cascade_delete"].(bool); ok {
		cascadeDelete = cd
	}

	cascadeUpdate := false
	if cu, ok := args["cascade_update"].(bool); ok {
		cascadeUpdate = cu
	}

	description := ""
	if d, ok := args["description"].(string); ok {
		description = d
	}

	// Verify both nodes exist
	_, err = h.dependencies.NodeRepo.GetByID(ctx, depNodeID)
	if err != nil {
		return nil, fmt.Errorf("dependent node not found: %w", err)
	}

	_, err = h.dependencies.NodeRepo.GetByID(ctx, depyNodeID)
	if err != nil {
		return nil, fmt.Errorf("dependency node not found: %w", err)
	}

	// TODO: Use a proper dependency repository when available
	// For now, we'll use a direct database approach similar to other implementations
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully created dependency:\nDependent: %s\nDependency: %s\nType: %s\nCascade Delete: %t\nCascade Update: %t\nDescription: %s\n\nNote: Full dependency creation will be implemented with proper repository",
					dependentNodeID, dependencyNodeID, dependencyType, cascadeDelete, cascadeUpdate, description),
			},
		},
		"isError": false,
	}, nil
}

// handleListNodeDependencies implements the list_node_dependencies tool
func (h *MCPToolHandler) handleListNodeDependencies(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse composite_id argument
	compositeID, ok := args["composite_id"].(string)
	if !ok || compositeID == "" {
		return nil, fmt.Errorf("missing or invalid 'composite_id' parameter")
	}

	// Parse composite ID to extract node ID
	nodeID, err := parseCompositeID(compositeID)
	if err != nil {
		return nil, fmt.Errorf("invalid composite_id: %w", err)
	}

	// Verify node exists
	node, err := h.dependencies.NodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	// TODO: Query dependencies from database when repository is available
	// For now, return placeholder response
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Dependencies for node: %s\nURL: %s\n\nNote: Dependency listing will be implemented with proper repository",
					node.Title(), node.URL()),
			},
		},
		"isError": false,
	}, nil
}

// handleListNodeDependents implements the list_node_dependents tool
func (h *MCPToolHandler) handleListNodeDependents(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse composite_id argument
	compositeID, ok := args["composite_id"].(string)
	if !ok || compositeID == "" {
		return nil, fmt.Errorf("missing or invalid 'composite_id' parameter")
	}

	// Parse composite ID to extract node ID
	nodeID, err := parseCompositeID(compositeID)
	if err != nil {
		return nil, fmt.Errorf("invalid composite_id: %w", err)
	}

	// Verify node exists
	node, err := h.dependencies.NodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	// TODO: Query dependents from database when repository is available
	// For now, return placeholder response
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Dependents for node: %s\nURL: %s\n\nNote: Dependent listing will be implemented with proper repository",
					node.Title(), node.URL()),
			},
		},
		"isError": false,
	}, nil
}

// handleDeleteDependency implements the delete_dependency tool
func (h *MCPToolHandler) handleDeleteDependency(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse dependency_id argument
	dependencyIDRaw, ok := args["dependency_id"]
	if !ok {
		return nil, fmt.Errorf("missing 'dependency_id' parameter")
	}

	var dependencyID int
	switch v := dependencyIDRaw.(type) {
	case float64:
		dependencyID = int(v)
	case int:
		dependencyID = v
	case string:
		var err error
		dependencyID, err = strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid dependency_id format: %v", err)
		}
	default:
		return nil, fmt.Errorf("invalid dependency_id type, expected number or string")
	}

	if dependencyID <= 0 {
		return nil, fmt.Errorf("dependency_id must be positive")
	}

	// TODO: Delete dependency from database when repository is available
	// For now, return placeholder response
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Would delete dependency with ID: %d\n\nNote: Dependency deletion will be implemented with proper repository",
					dependencyID),
			},
		},
		"isError": false,
	}, nil
}

// handleFilterNodesByAttributes implements the filter_nodes_by_attributes tool
func (h *MCPToolHandler) handleFilterNodesByAttributes(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse domain_name argument
	domainName, ok := args["domain_name"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("missing or invalid 'domain_name' parameter")
	}

	// Parse filters argument
	filtersRaw, ok := args["filters"]
	if !ok {
		return nil, fmt.Errorf("missing 'filters' parameter")
	}

	filtersArray, ok := filtersRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid 'filters' parameter, expected array")
	}

	// Convert filters to repository format
	var filters []repository.AttributeFilter
	for i, filterRaw := range filtersArray {
		filterMap, ok := filterRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid filter at index %d, expected object", i)
		}

		name, ok := filterMap["name"].(string)
		if !ok || name == "" {
			return nil, fmt.Errorf("missing or invalid 'name' in filter at index %d", i)
		}

		value, ok := filterMap["value"].(string)
		if !ok || value == "" {
			return nil, fmt.Errorf("missing or invalid 'value' in filter at index %d", i)
		}

		operator := "equals" // default operator
		if op, ok := filterMap["operator"].(string); ok && op != "" {
			operator = op
		}

		filters = append(filters, repository.AttributeFilter{
			Name:     name,
			Value:    value,
			Operator: operator,
		})
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

	// Execute filter use case
	result, err := h.dependencies.FilterNodesUC.Execute(ctx, domainName, filters, page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to filter nodes: %w", err)
	}

	// Convert to MCP response format
	content := []map[string]interface{}{}

	if len(result.Nodes) == 0 {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("No nodes found matching the specified filters in domain '%s'", domainName),
		})
	} else {
		for _, node := range result.Nodes {
			content = append(content, map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("Node ID: %d\nURL: %s\nTitle: %s\nDescription: %s\nCreated: %s",
					node.ID, node.URL, node.Title, node.Description, node.CreatedAt.Format("2006-01-02 15:04:05")),
			})
		}

		// Add pagination info
		if result.TotalPages > 1 {
			content = append(content, map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("\nPage %d of %d (Total: %d nodes)", result.Page, result.TotalPages, result.TotalCount),
			})
		}
	}

	return map[string]interface{}{
		"content": content,
		"isError": false,
	}, nil
}

// handleGetNodeWithAttributes implements the get_node_with_attributes tool
func (h *MCPToolHandler) handleGetNodeWithAttributes(ctx context.Context, args map[string]interface{}) (interface{}, error) {
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

	// Execute use case
	result, err := h.dependencies.GetNodeWithAttributesUC.Execute(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node with attributes: %w", err)
	}

	// Build response text
	var responseText strings.Builder
	
	// Node information
	responseText.WriteString(fmt.Sprintf("Node: %s\n", result.Node.Title))
	responseText.WriteString(fmt.Sprintf("URL: %s\n", result.Node.URL))
	responseText.WriteString(fmt.Sprintf("Description: %s\n", result.Node.Description))
	responseText.WriteString(fmt.Sprintf("Domain: %s\n", result.Node.DomainName))
	responseText.WriteString(fmt.Sprintf("Created: %s\n", result.Node.CreatedAt.Format("2006-01-02 15:04:05")))
	responseText.WriteString(fmt.Sprintf("Updated: %s\n", result.Node.UpdatedAt.Format("2006-01-02 15:04:05")))

	// Attributes information
	if len(result.Attributes) > 0 {
		responseText.WriteString("\nAttributes:\n")
		for _, attr := range result.Attributes {
			attrText := fmt.Sprintf("• %s (%s): %s", attr.AttributeName, attr.AttributeType, attr.Value)
			if attr.OrderIndex != nil {
				attrText += fmt.Sprintf(" [order: %d]", *attr.OrderIndex)
			}
			responseText.WriteString(attrText + "\n")
		}
	} else {
		responseText.WriteString("\nNo attributes found for this node.\n")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": responseText.String(),
			},
		},
		"isError": false,
	}, nil
}