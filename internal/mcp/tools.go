package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"url-db/internal/models"
)

// MCP Tool 정의 및 구현

// ToolRegistry manages all available MCP tools
type ToolRegistry struct {
	service MCPService
	tools   []Tool
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry(service MCPService) *ToolRegistry {
	registry := &ToolRegistry{
		service: service,
	}
	registry.registerTools()
	return registry
}

// GetTools returns all available tools
func (tr *ToolRegistry) GetTools() []Tool {
	return tr.tools
}

// CallTool executes a tool by name
func (tr *ToolRegistry) CallTool(ctx context.Context, name string, arguments interface{}) (*CallToolResult, error) {
	switch name {
	case "list_mcp_domains":
		return tr.callListDomains(ctx, arguments)
	case "create_mcp_domain":
		return tr.callCreateDomain(ctx, arguments)
	case "list_mcp_nodes":
		return tr.callListNodes(ctx, arguments)
	case "create_mcp_node":
		return tr.callCreateNode(ctx, arguments)
	case "get_mcp_node":
		return tr.callGetNode(ctx, arguments)
	case "update_mcp_node":
		return tr.callUpdateNode(ctx, arguments)
	case "delete_mcp_node":
		return tr.callDeleteNode(ctx, arguments)
	case "find_mcp_node_by_url":
		return tr.callFindNodeByURL(ctx, arguments)
	case "get_mcp_node_attributes":
		return tr.callGetNodeAttributes(ctx, arguments)
	case "set_mcp_node_attributes":
		return tr.callSetNodeAttributes(ctx, arguments)
	case "list_mcp_domain_attributes":
		return tr.callListDomainAttributes(ctx, arguments)
	case "create_mcp_domain_attribute":
		return tr.callCreateDomainAttribute(ctx, arguments)
	case "get_mcp_domain_attribute":
		return tr.callGetDomainAttribute(ctx, arguments)
	case "update_mcp_domain_attribute":
		return tr.callUpdateDomainAttribute(ctx, arguments)
	case "delete_mcp_domain_attribute":
		return tr.callDeleteDomainAttribute(ctx, arguments)
	case "get_mcp_server_info":
		return tr.callGetServerInfo(ctx, arguments)
	default:
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Unknown tool: %s", name)}},
			IsError: true,
		}, nil
	}
}

// registerTools registers all available tools
func (tr *ToolRegistry) registerTools() {
	tr.tools = []Tool{
		{
			Name:        "list_mcp_domains",
			Description: "List all domains in the URL database",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "create_mcp_domain",
			Description: "Create a new domain",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Domain name",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Domain description",
					},
				},
				"required": []string{"name", "description"},
			},
		},
		{
			Name:        "list_mcp_nodes",
			Description: "List nodes in a specific domain",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{
						"type":        "string",
						"description": "Domain name to list nodes from",
					},
					"page": map[string]interface{}{
						"type":        "integer",
						"description": "Page number (default: 1)",
						"default":     1,
					},
					"size": map[string]interface{}{
						"type":        "integer",
						"description": "Page size (default: 20)",
						"default":     20,
					},
					"search": map[string]interface{}{
						"type":        "string",
						"description": "Search query (optional)",
					},
				},
				"required": []string{"domain_name"},
			},
		},
		{
			Name:        "create_mcp_node",
			Description: "Create a new node (URL) in a domain",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{
						"type":        "string",
						"description": "Domain name",
					},
					"url": map[string]interface{}{
						"type":        "string",
						"description": "URL to store",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Node title (optional)",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Node description (optional)",
					},
				},
				"required": []string{"domain_name", "url"},
			},
		},
		{
			Name:        "get_mcp_node",
			Description: "Get a node by composite ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool:domain:id)",
					},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "update_mcp_node",
			Description: "Update a node's title and description",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool:domain:id)",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "New title (optional)",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "New description (optional)",
					},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "delete_mcp_node",
			Description: "Delete a node by composite ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool:domain:id)",
					},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "find_mcp_node_by_url",
			Description: "Find a node by URL in a domain",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{
						"type":        "string",
						"description": "Domain name",
					},
					"url": map[string]interface{}{
						"type":        "string",
						"description": "URL to find",
					},
				},
				"required": []string{"domain_name", "url"},
			},
		},
		{
			Name:        "get_mcp_node_attributes",
			Description: "Get all attributes for a node",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool:domain:id)",
					},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "set_mcp_node_attributes",
			Description: "Set attributes for a node",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool:domain:id)",
					},
					"attributes": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name": map[string]interface{}{
									"type":        "string",
									"description": "Attribute name",
								},
								"value": map[string]interface{}{
									"type":        "string",
									"description": "Attribute value",
								},
							},
							"required": []string{"name", "value"},
						},
						"description": "Array of attributes to set",
					},
				},
				"required": []string{"composite_id", "attributes"},
			},
		},
		{
			Name:        "get_mcp_server_info",
			Description: "Get server information and capabilities",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "list_mcp_domain_attributes",
			Description: "List all attribute definitions for a domain",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{
						"type":        "string",
						"description": "The domain to list attributes for",
					},
				},
				"required": []string{"domain_name"},
			},
		},
		{
			Name:        "create_mcp_domain_attribute",
			Description: "Create a new attribute definition for a domain",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{
						"type":        "string",
						"description": "The domain to add attribute to",
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Attribute name",
					},
					"type": map[string]interface{}{
						"type":        "string",
						"description": "One of: tag, ordered_tag, number, string, markdown, image",
						"enum":        []string{"tag", "ordered_tag", "number", "string", "markdown", "image"},
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Human-readable description (optional)",
					},
				},
				"required": []string{"domain_name", "name", "type"},
			},
		},
		{
			Name:        "get_mcp_domain_attribute",
			Description: "Get a specific attribute definition",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool-name:domain:attribute-id)",
					},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "update_mcp_domain_attribute",
			Description: "Update attribute description",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool-name:domain:attribute-id)",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "New description",
					},
				},
				"required": []string{"composite_id", "description"},
			},
		},
		{
			Name:        "delete_mcp_domain_attribute",
			Description: "Delete an attribute definition (if no values exist)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool-name:domain:attribute-id)",
					},
				},
				"required": []string{"composite_id"},
			},
		},
	}
}

// Tool implementation methods

func (tr *ToolRegistry) callListDomains(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	response, err := tr.service.ListDomains(ctx)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error listing domains: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callCreateDomain(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	req := &models.CreateDomainRequest{
		Name:        argsMap["name"].(string),
		Description: argsMap["description"].(string),
	}

	domain, err := tr.service.CreateDomain(ctx, req)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error creating domain: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(domain, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callListNodes(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	domainName := argsMap["domain_name"].(string)
	page := 1
	size := 20
	search := ""

	if p, exists := argsMap["page"]; exists {
		if pFloat, ok := p.(float64); ok {
			page = int(pFloat)
		}
	}
	if s, exists := argsMap["size"]; exists {
		if sFloat, ok := s.(float64); ok {
			size = int(sFloat)
		}
	}
	if s, exists := argsMap["search"]; exists {
		search = s.(string)
	}

	response, err := tr.service.ListNodes(ctx, domainName, page, size, search)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error listing nodes: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callCreateNode(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	req := &models.CreateMCPNodeRequest{
		DomainName:  argsMap["domain_name"].(string),
		URL:         argsMap["url"].(string),
		Title:       getStringArg(argsMap, "title"),
		Description: getStringArg(argsMap, "description"),
	}

	node, err := tr.service.CreateNode(ctx, req)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error creating node: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(node, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callGetNode(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	compositeID := argsMap["composite_id"].(string)
	node, err := tr.service.GetNode(ctx, compositeID)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error getting node: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(node, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callUpdateNode(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	compositeID := argsMap["composite_id"].(string)
	req := &models.UpdateNodeRequest{
		Title:       getStringArg(argsMap, "title"),
		Description: getStringArg(argsMap, "description"),
	}

	node, err := tr.service.UpdateNode(ctx, compositeID, req)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error updating node: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(node, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callDeleteNode(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	compositeID := argsMap["composite_id"].(string)
	err := tr.service.DeleteNode(ctx, compositeID)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error deleting node: %v", err)}},
			IsError: true,
		}, nil
	}

	return &CallToolResult{
		Content: []Content{{Type: "text", Text: fmt.Sprintf("Node %s deleted successfully", compositeID)}},
	}, nil
}

func (tr *ToolRegistry) callFindNodeByURL(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	req := &models.FindMCPNodeRequest{
		DomainName: argsMap["domain_name"].(string),
		URL:        argsMap["url"].(string),
	}

	node, err := tr.service.FindNodeByURL(ctx, req)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error finding node: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(node, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callGetNodeAttributes(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	compositeID := argsMap["composite_id"].(string)
	response, err := tr.service.GetNodeAttributes(ctx, compositeID)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error getting node attributes: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callSetNodeAttributes(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	compositeID := argsMap["composite_id"].(string)
	
	// Parse attributes array
	attributesRaw, exists := argsMap["attributes"]
	if !exists {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Missing attributes parameter"}},
			IsError: true,
		}, nil
	}

	attributesArray, ok := attributesRaw.([]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid attributes format"}},
			IsError: true,
		}, nil
	}

	var attributes []struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	for _, attrRaw := range attributesArray {
		attrMap, ok := attrRaw.(map[string]interface{})
		if !ok {
			continue
		}
		
		attributes = append(attributes, struct {
			Name  string `json:"name" binding:"required"`
			Value string `json:"value" binding:"required"`
		}{
			Name:  attrMap["name"].(string),
			Value: attrMap["value"].(string),
		})
	}

	req := &models.SetMCPNodeAttributesRequest{
		Attributes: attributes,
	}

	response, err := tr.service.SetNodeAttributes(ctx, compositeID, req)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error setting node attributes: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callGetServerInfo(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	info, err := tr.service.GetServerInfo(ctx)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error getting server info: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(info, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

// Helper function to safely get string arguments
func getStringArg(args map[string]interface{}, key string) string {
	if val, exists := args[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// Domain attribute management methods
func (tr *ToolRegistry) callListDomainAttributes(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	domainName := getStringArg(argsMap, "domain_name")
	if domainName == "" {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Missing domain_name parameter"}},
			IsError: true,
		}, nil
	}

	response, err := tr.service.ListDomainAttributes(ctx, domainName)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error listing domain attributes: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callCreateDomainAttribute(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	domainName := getStringArg(argsMap, "domain_name")
	name := getStringArg(argsMap, "name")
	typeStr := getStringArg(argsMap, "type")
	description := getStringArg(argsMap, "description")

	if domainName == "" || name == "" || typeStr == "" {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Missing required parameters: domain_name, name, type"}},
			IsError: true,
		}, nil
	}

	req := &models.CreateAttributeRequest{
		Name:        name,
		Type:        models.AttributeType(typeStr),
		Description: description,
	}

	response, err := tr.service.CreateDomainAttribute(ctx, domainName, req)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error creating domain attribute: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callGetDomainAttribute(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	compositeID := getStringArg(argsMap, "composite_id")
	if compositeID == "" {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Missing composite_id parameter"}},
			IsError: true,
		}, nil
	}

	response, err := tr.service.GetDomainAttribute(ctx, compositeID)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error getting domain attribute: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callUpdateDomainAttribute(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	compositeID := getStringArg(argsMap, "composite_id")
	description := getStringArg(argsMap, "description")

	if compositeID == "" || description == "" {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Missing required parameters: composite_id, description"}},
			IsError: true,
		}, nil
	}

	req := &models.UpdateAttributeRequest{
		Description: description,
	}

	response, err := tr.service.UpdateDomainAttribute(ctx, compositeID, req)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error updating domain attribute: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callDeleteDomainAttribute(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
	argsMap, ok := arguments.(map[string]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid arguments format"}},
			IsError: true,
		}, nil
	}

	compositeID := getStringArg(argsMap, "composite_id")
	if compositeID == "" {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Missing composite_id parameter"}},
			IsError: true,
		}, nil
	}

	err := tr.service.DeleteDomainAttribute(ctx, compositeID)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error deleting domain attribute: %v", err)}},
			IsError: true,
		}, nil
	}

	return &CallToolResult{
		Content: []Content{{Type: "text", Text: "Domain attribute deleted successfully"}},
	}, nil
}