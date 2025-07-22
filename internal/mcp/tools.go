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
	case ListDomainsTool:
		return tr.callListDomains(ctx, arguments)
	case CreateDomainTool:
		return tr.callCreateDomain(ctx, arguments)
	case ListNodesTool:
		return tr.callListNodes(ctx, arguments)
	case CreateNodeTool:
		return tr.callCreateNode(ctx, arguments)
	case GetNodeTool:
		return tr.callGetNode(ctx, arguments)
	case UpdateNodeTool:
		return tr.callUpdateNode(ctx, arguments)
	case DeleteNodeTool:
		return tr.callDeleteNode(ctx, arguments)
	case FindNodeByUrlTool:
		return tr.callFindNodeByURL(ctx, arguments)
	case GetNodeAttributesTool:
		return tr.callGetNodeAttributes(ctx, arguments)
	case SetNodeAttributesTool:
		return tr.callSetNodeAttributes(ctx, arguments)
	case ListDomainAttributesTool:
		return tr.callListDomainAttributes(ctx, arguments)
	case CreateDomainAttributeTool:
		return tr.callCreateDomainAttribute(ctx, arguments)
	case GetDomainAttributeTool:
		return tr.callGetDomainAttribute(ctx, arguments)
	case UpdateDomainAttributeTool:
		return tr.callUpdateDomainAttribute(ctx, arguments)
	case DeleteDomainAttributeTool:
		return tr.callDeleteDomainAttribute(ctx, arguments)
	case GetNodeWithAttributesTool:
		return tr.callGetNodeWithAttributes(ctx, arguments)
	case FilterNodesByAttributesTool:
		return tr.callFilterNodesByAttributes(ctx, arguments)
	case GetServerInfoTool:
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
			Name: ListDomainsTool,
			Description: "Get all domains",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name: CreateDomainTool,
			Description: "Create new domain for organizing URLs",
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
			Name: ListNodesTool,
			Description: "List URLs in domain",
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
			Name: CreateNodeTool,
			Description: "Add URL to domain",
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
			Name: GetNodeTool,
			Description: "Get URL details",
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
			Name: UpdateNodeTool,
			Description: "Update URL title or description",
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
			Name: DeleteNodeTool,
			Description: "Remove URL",
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
			Name: FindNodeByUrlTool,
			Description: "Search by exact URL",
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
			Name: GetNodeAttributesTool,
			Description: "Get URL tags and attributes",
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
			Name: SetNodeAttributesTool,
			Description: "Add or update URL tags",
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
								"order_index": map[string]interface{}{
									"type":        "integer",
									"description": "Order index (required for ordered_tag type)",
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
			Name: GetServerInfoTool,
			Description: "Get server info",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name: ListDomainAttributesTool,
			Description: "Get available tag types for domain",
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
			Name: CreateDomainAttributeTool,
			Description: "Define new tag type for domain",
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
			Name: GetDomainAttributeTool,
			Description: "Get tag type details",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool-name:domain:attr-{id})",
					},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name: UpdateDomainAttributeTool,
			Description: "Update tag type description",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool-name:domain:attr-{id})",
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
			Name: DeleteDomainAttributeTool,
			Description: "Remove unused tag type",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool-name:domain:attr-{id})",
					},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name: GetNodeWithAttributesTool,
			Description: "Get URL with all tags at once",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{
						"type":        "string",
						"description": "Composite ID (format: tool-name:domain:id)",
					},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name: FilterNodesByAttributesTool,
			Description: "Search URLs by tag values",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{
						"type":        "string",
						"description": "Domain name to search in",
					},
					"filters": map[string]interface{}{
						"type":        "array",
						"description": "Array of attribute filters",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name": map[string]interface{}{
									"type":        "string",
									"description": "Attribute name",
								},
								"value": map[string]interface{}{
									"type":        "string",
									"description": "Attribute value to match",
								},
								"operator": map[string]interface{}{
									"type":        "string",
									"description": "Comparison operator: equals, contains, starts_with, ends_with",
									"enum":        []string{"equals", "contains", "starts_with", "ends_with"},
									"default":     "equals",
								},
							},
							"required": []string{"name", "value"},
						},
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
				},
				"required": []string{"domain_name", "filters"},
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
		Name       string `json:"name" binding:"required"`
		Value      string `json:"value" binding:"required"`
		OrderIndex *int   `json:"order_index"`
	}
	for _, attrRaw := range attributesArray {
		attrMap, ok := attrRaw.(map[string]interface{})
		if !ok {
			continue
		}

		var orderIndex *int
		if orderIndexRaw, exists := attrMap["order_index"]; exists {
			if orderIndexFloat, ok := orderIndexRaw.(float64); ok {
				orderIndexInt := int(orderIndexFloat)
				orderIndex = &orderIndexInt
			}
		}

		attributes = append(attributes, struct {
			Name       string `json:"name" binding:"required"`
			Value      string `json:"value" binding:"required"`
			OrderIndex *int   `json:"order_index"`
		}{
			Name:       attrMap["name"].(string),
			Value:      attrMap["value"].(string),
			OrderIndex: orderIndex,
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

func (tr *ToolRegistry) callGetNodeWithAttributes(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
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

	// Get node information
	node, err := tr.service.GetNode(ctx, compositeID)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error getting node: %v", err)}},
			IsError: true,
		}, nil
	}

	// Get node attributes
	attributes, err := tr.service.GetNodeAttributes(ctx, compositeID)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error getting node attributes: %v", err)}},
			IsError: true,
		}, nil
	}

	// Combine node and attributes into a single response
	response := map[string]interface{}{
		"node":       node,
		"attributes": attributes.Attributes,
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}

func (tr *ToolRegistry) callFilterNodesByAttributes(ctx context.Context, arguments interface{}) (*CallToolResult, error) {
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

	// Parse filters
	filtersRaw, exists := argsMap["filters"]
	if !exists {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Missing filters parameter"}},
			IsError: true,
		}, nil
	}

	filtersArray, ok := filtersRaw.([]interface{})
	if !ok {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: "Invalid filters format"}},
			IsError: true,
		}, nil
	}

	type attributeFilter struct {
		Name     string `json:"name"`
		Value    string `json:"value"`
		Operator string `json:"operator"`
	}

	var filters []attributeFilter
	for _, filterRaw := range filtersArray {
		filterMap, ok := filterRaw.(map[string]interface{})
		if !ok {
			continue
		}

		operator := getStringArg(filterMap, "operator")
		if operator == "" {
			operator = "equals"
		}

		filter := attributeFilter{
			Name:     getStringArg(filterMap, "name"),
			Value:    getStringArg(filterMap, "value"),
			Operator: operator,
		}

		filters = append(filters, filter)
	}

	// Parse pagination
	page := 1
	size := 20
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

	// Convert filters to interface slice as maps
	var filterInterfaces []interface{}
	for _, f := range filters {
		filterMap := map[string]interface{}{
			"name":     f.Name,
			"value":    f.Value,
			"operator": f.Operator,
		}
		filterInterfaces = append(filterInterfaces, filterMap)
	}

	// Call service method to filter nodes
	response, err := tr.service.FilterNodesByAttributes(ctx, domainName, filterInterfaces, page, size)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error filtering nodes: %v", err)}},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return &CallToolResult{
		Content: []Content{{Type: "text", Text: string(result)}},
	}, nil
}
