package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"url-db/internal/constants"
	"url-db/internal/interface/setup"
)

// MCPProtocolHandler handles MCP JSON-RPC 2.0 protocol logic
type MCPProtocolHandler struct {
	factory     *setup.ApplicationFactory
	toolHandler *MCPToolHandler
	mode        string
}

// NewMCPProtocolHandler creates a new protocol handler
func NewMCPProtocolHandler(factory *setup.ApplicationFactory, mode string) *MCPProtocolHandler {
	return &MCPProtocolHandler{
		factory:     factory,
		toolHandler: NewMCPToolHandler(factory),
		mode:        mode,
	}
}

// HandleRequest processes a JSON-RPC request and returns a response
func (h *MCPProtocolHandler) HandleRequest(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	// Validate JSON-RPC version
	if req.JSONRPC != constants.JSONRPCVersion {
		return h.createErrorResponse(req.ID, InvalidRequest, "Invalid JSON-RPC version", nil)
	}

	// Route the request based on method
	switch req.Method {
	case "initialize":
		return h.handleInitialize(req)
	case "tools/list":
		return h.handleToolsList(req)
	case "tools/call":
		return h.handleToolCall(ctx, req)
	case "resources/list":
		return h.handleResourcesList(req)
	case "resources/read":
		return h.handleResourceRead(req)
	case "notifications/initialized":
		// Client notification that initialization is complete
		// No response needed for notifications
		return nil
	default:
		return h.createErrorResponse(req.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", req.Method), nil)
	}
}

// handleInitialize handles MCP initialization
func (h *MCPProtocolHandler) handleInitialize(req *JSONRPCRequest) *JSONRPCResponse {
	result := map[string]interface{}{
		"protocolVersion": constants.MCPProtocolVersion,
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true,
			},
			"resources": map[string]interface{}{
				"subscribe":   true,
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    constants.MCPServerName,
			"version": constants.DefaultServerVersion,
		},
	}

	return h.createSuccessResponse(req.ID, result)
}

// handleToolsList returns available MCP tools
func (h *MCPProtocolHandler) handleToolsList(req *JSONRPCRequest) *JSONRPCResponse {
	tools := []map[string]interface{}{
		{
			"name":        "get_server_info",
			"description": "Get server information",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "list_domains",
			"description": "Get all domains",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page": map[string]interface{}{"type": "integer", "default": 1},
					"size": map[string]interface{}{"type": "integer", "default": 20},
				},
			},
		},
		{
			"name":        "create_domain",
			"description": "Create new domain for organizing URLs",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":        map[string]interface{}{"type": "string", "description": "Domain name"},
					"description": map[string]interface{}{"type": "string", "description": "Domain description"},
				},
				"required": []string{"name", "description"},
			},
		},
		{
			"name":        "list_nodes",
			"description": "List URLs in domain (requires: domain must exist via create_domain)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "Domain name to list nodes from"},
					"page":        map[string]interface{}{"type": "integer", "default": 1},
					"size":        map[string]interface{}{"type": "integer", "default": 20},
					"search":      map[string]interface{}{"type": "string", "description": "Search query"},
				},
				"required": []string{"domain_name"},
			},
		},
		{
			"name":        "create_node",
			"description": "Add URL to domain (requires: domain must exist via create_domain)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "Domain name"},
					"url":         map[string]interface{}{"type": "string", "description": "URL to store"},
					"title":       map[string]interface{}{"type": "string", "description": "Node title"},
					"description": map[string]interface{}{"type": "string", "description": "Node description"},
				},
				"required": []string{"domain_name", "url"},
			},
		},
		{
			"name":        "get_node",
			"description": "Get URL details (requires: node must exist via create_node; returns composite_id from create_node)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "update_node",
			"description": "Update URL title or description (requires: node must exist via create_node; use composite_id from create_node)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
					"title":        map[string]interface{}{"type": "string", "description": "New title"},
					"description":  map[string]interface{}{"type": "string", "description": "New description"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "delete_node",
			"description": "Remove URL (requires: node must exist via create_node; use composite_id from create_node)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "find_node_by_url",
			"description": "Search by exact URL (requires: domain must exist via create_domain; returns composite_id if found)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "Domain name"},
					"url":         map[string]interface{}{"type": "string", "description": "URL to find"},
				},
				"required": []string{"domain_name", "url"},
			},
		},
		{
			"name":        "scan_all_content",
			"description": "Retrieve all URLs and their content from a domain using page-based navigation with token optimization for AI processing (requires: domain must exist via create_domain; includes attributes if domain has create_domain_attribute)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":         map[string]interface{}{"type": "string", "description": "Domain name to scan"},
					"max_tokens_per_page": map[string]interface{}{"type": "integer", "description": "Maximum tokens per page (recommended: 6000-10000)", "default": 8000},
					"page":                map[string]interface{}{"type": "integer", "description": "Page number (1-based)", "default": 1},
					"include_attributes":  map[string]interface{}{"type": "boolean", "description": "Include node attributes in response", "default": true},
					"compress_attributes": map[string]interface{}{"type": "boolean", "description": "Remove duplicate attribute values for AI context compression", "default": false},
				},
				"required": []string{"domain_name"},
			},
		},
		{
			"name":        "get_node_attributes",
			"description": "Get URL tags and attributes (requires: node must exist via create_node; attributes defined via create_domain_attribute)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "set_node_attributes",
			"description": "Add or update URL tags (requires: node must exist via create_node; attributes should be defined via create_domain_attribute unless auto_create_attributes=true)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
					"attributes": map[string]interface{}{
						"type":        "array",
						"description": "Array of attributes to set",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name":        map[string]interface{}{"type": "string", "description": "Attribute name"},
								"value":       map[string]interface{}{"type": "string", "description": "Attribute value"},
								"order_index": map[string]interface{}{"type": "integer", "description": "Order index (for ordered_tag type)"},
							},
							"required": []string{"name", "value"},
						},
					},
					"auto_create_attributes": map[string]interface{}{"type": "boolean", "default": true, "description": "Automatically create attributes if they don't exist"},
				},
				"required": []string{"composite_id", "attributes"},
			},
		},
		{
			"name":        "list_domain_attributes",
			"description": "Get available tag types for domain (requires: domain must exist via create_domain)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "The domain to list attributes for"},
				},
				"required": []string{"domain_name"},
			},
		},
		{
			"name":        "create_domain_attribute",
			"description": "Define new tag type for domain (requires: domain must exist via create_domain; enables attributes for set_node_attributes)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "The domain to add attribute to"},
					"name":        map[string]interface{}{"type": "string", "description": "Attribute name"},
					"type": map[string]interface{}{
						"type":        "string",
						"description": "One of: tag, ordered_tag, number, string, markdown, image",
						"enum":        []string{"tag", "ordered_tag", "number", "string", "markdown", "image"},
					},
					"description": map[string]interface{}{"type": "string", "description": "Human-readable description"},
				},
				"required": []string{"domain_name", "name", "type"},
			},
		},
		{
			"name":        "get_domain_attribute",
			"description": "Get details of a specific domain attribute (requires: attribute must exist via create_domain_attribute)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":    map[string]interface{}{"type": "string", "description": "The domain name"},
					"attribute_name": map[string]interface{}{"type": "string", "description": "The attribute name to get"},
				},
				"required": []string{"domain_name", "attribute_name"},
			},
		},
		{
			"name":        "update_domain_attribute",
			"description": "Update domain attribute description (requires: attribute must exist via create_domain_attribute)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":    map[string]interface{}{"type": "string", "description": "The domain name"},
					"attribute_name": map[string]interface{}{"type": "string", "description": "The attribute name to update"},
					"description":    map[string]interface{}{"type": "string", "description": "New description for the attribute"},
				},
				"required": []string{"domain_name", "attribute_name"},
			},
		},
		{
			"name":        "delete_domain_attribute",
			"description": "Remove domain attribute definition (requires: attribute must exist via create_domain_attribute; removes all values from nodes)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":    map[string]interface{}{"type": "string", "description": "The domain name"},
					"attribute_name": map[string]interface{}{"type": "string", "description": "The attribute name to delete"},
				},
				"required": []string{"domain_name", "attribute_name"},
			},
		},
		{
			"name":        "create_dependency",
			"description": "Create dependency relationship between nodes (requires: both nodes must exist via create_node; use composite_ids from create_node)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"dependent_node_id":  map[string]interface{}{"type": "string", "description": "Composite ID of the dependent node (format: tool:domain:id)"},
					"dependency_node_id": map[string]interface{}{"type": "string", "description": "Composite ID of the dependency node (format: tool:domain:id)"},
					"dependency_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of dependency",
						"enum":        []string{"hard", "soft", "reference"},
					},
					"cascade_delete": map[string]interface{}{"type": "boolean", "default": false, "description": "Whether to cascade delete"},
					"cascade_update": map[string]interface{}{"type": "boolean", "default": false, "description": "Whether to cascade update"},
					"description":    map[string]interface{}{"type": "string", "description": "Optional description of the dependency"},
				},
				"required": []string{"dependent_node_id", "dependency_node_id", "dependency_type"},
			},
		},
		{
			"name":        "list_node_dependencies",
			"description": "List what a node depends on (requires: node must exist via create_node; dependencies created via create_dependency)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID of the node (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "list_node_dependents",
			"description": "List what depends on a node (requires: node must exist via create_node; dependencies created via create_dependency)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID of the node (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "delete_dependency",
			"description": "Remove dependency relationship (requires: dependency must exist via create_dependency; use dependency_id from list_node_dependencies)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"dependency_id": map[string]interface{}{"type": "integer", "description": "ID of the dependency relationship to delete"},
				},
				"required": []string{"dependency_id"},
			},
		},
		{
			"name":        "filter_nodes_by_attributes",
			"description": "Filter nodes by attribute values (requires: domain must exist via create_domain; attributes defined via create_domain_attribute)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "Domain name to filter nodes from"},
					"filters": map[string]interface{}{
						"type":        "array",
						"description": "Array of attribute filters",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name":     map[string]interface{}{"type": "string", "description": "Attribute name"},
								"value":    map[string]interface{}{"type": "string", "description": "Attribute value"},
								"operator": map[string]interface{}{"type": "string", "description": "Comparison operator", "enum": []string{"equals", "contains", "starts_with", "ends_with"}, "default": "equals"},
							},
							"required": []string{"name", "value"},
						},
					},
					"page": map[string]interface{}{"type": "integer", "default": 1},
					"size": map[string]interface{}{"type": "integer", "default": 20},
				},
				"required": []string{"domain_name", "filters"},
			},
		},
		{
			"name":        "get_node_with_attributes",
			"description": "Get URL details with all attributes (requires: node must exist via create_node; combines get_node + get_node_attributes)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "list_templates",
			"description": "List templates in domain (requires: domain must exist via create_domain)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":   map[string]interface{}{"type": "string", "description": "Domain name to list templates from"},
					"page":          map[string]interface{}{"type": "integer", "default": 1},
					"size":          map[string]interface{}{"type": "integer", "default": 20},
					"template_type": map[string]interface{}{"type": "string", "description": "Filter by template type"},
					"active_only":   map[string]interface{}{"type": "boolean", "default": false, "description": "Only return active templates"},
					"search":        map[string]interface{}{"type": "string", "description": "Search query"},
				},
				"required": []string{"domain_name"},
			},
		},
		{
			"name":        "create_template",
			"description": "Create new template in domain (requires: domain must exist via create_domain; use validate_template to check template_data)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":   map[string]interface{}{"type": "string", "description": "Domain name"},
					"name":          map[string]interface{}{"type": "string", "description": "Template name"},
					"template_data": map[string]interface{}{"type": "string", "description": "JSON template data"},
					"title":         map[string]interface{}{"type": "string", "description": "Template title"},
					"description":   map[string]interface{}{"type": "string", "description": "Template description"},
				},
				"required": []string{"domain_name", "name", "template_data"},
			},
		},
		{
			"name":        "get_template",
			"description": "Get template details (requires: template must exist via create_template; use composite_id from create_template)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:template:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "update_template",
			"description": "Update template (requires: template must exist via create_template; use validate_template to check new template_data)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id":  map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:template:id)"},
					"template_data": map[string]interface{}{"type": "string", "description": "Updated JSON template data"},
					"title":         map[string]interface{}{"type": "string", "description": "Updated title"},
					"description":   map[string]interface{}{"type": "string", "description": "Updated description"},
					"is_active":     map[string]interface{}{"type": "boolean", "description": "Template active status"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "delete_template",
			"description": "Delete template (requires: template must exist via create_template)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:template:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			"name":        "clone_template",
			"description": "Clone existing template (requires: source template must exist via create_template; creates new template with same domain)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"source_composite_id": map[string]interface{}{"type": "string", "description": "Source template composite ID (format: tool:domain:template:id)"},
					"new_name":            map[string]interface{}{"type": "string", "description": "New template name"},
					"new_title":           map[string]interface{}{"type": "string", "description": "New template title"},
					"new_description":     map[string]interface{}{"type": "string", "description": "New template description"},
				},
				"required": []string{"source_composite_id", "new_name"},
			},
		},
		{
			"name":        "generate_template_scaffold",
			"description": "Generate template scaffold for given type (helper: provides starting point for create_template)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"template_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of template to generate",
						"enum":        []string{"layout", "form", "document", "custom"},
					},
				},
				"required": []string{"template_type"},
			},
		},
		{
			"name":        "validate_template",
			"description": "Validate template data structure (helper: use before create_template or update_template)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"template_data": map[string]interface{}{"type": "string", "description": "JSON template data to validate"},
				},
				"required": []string{"template_data"},
			},
		},
	}

	result := map[string]interface{}{
		"tools": tools,
	}

	return h.createSuccessResponse(req.ID, result)
}

// handleToolCall executes a tool call
func (h *MCPProtocolHandler) handleToolCall(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return h.createErrorResponse(req.ID, InvalidParams, "Invalid tool call parameters", err.Error())
	}

	var result interface{}
	var err error

	switch params.Name {
	case "get_server_info":
		return h.handleGetServerInfo(req)
	case "list_domains":
		result, err = h.toolHandler.handleListDomains(ctx, params.Arguments)
	case "create_domain":
		result, err = h.toolHandler.handleCreateDomain(ctx, params.Arguments)
	case "list_nodes":
		result, err = h.toolHandler.handleListNodes(ctx, params.Arguments)
	case "create_node":
		result, err = h.toolHandler.handleCreateNode(ctx, params.Arguments)
	case "get_node":
		result, err = h.toolHandler.handleGetNode(ctx, params.Arguments)
	case "update_node":
		result, err = h.toolHandler.handleUpdateNode(ctx, params.Arguments)
	case "delete_node":
		result, err = h.toolHandler.handleDeleteNode(ctx, params.Arguments)
	case "find_node_by_url":
		result, err = h.toolHandler.handleFindNodeByURL(ctx, params.Arguments)
	case "scan_all_content":
		result, err = h.toolHandler.handleScanAllContent(ctx, params.Arguments)
	case "get_node_attributes":
		result, err = h.toolHandler.handleGetNodeAttributes(ctx, params.Arguments)
	case "set_node_attributes":
		result, err = h.toolHandler.handleSetNodeAttributes(ctx, params.Arguments)
	case "list_domain_attributes":
		result, err = h.toolHandler.handleListDomainAttributes(ctx, params.Arguments)
	case "create_domain_attribute":
		result, err = h.toolHandler.handleCreateDomainAttribute(ctx, params.Arguments)
	case "get_domain_attribute":
		result, err = h.toolHandler.handleGetDomainAttribute(ctx, params.Arguments)
	case "update_domain_attribute":
		result, err = h.toolHandler.handleUpdateDomainAttribute(ctx, params.Arguments)
	case "delete_domain_attribute":
		result, err = h.toolHandler.handleDeleteDomainAttribute(ctx, params.Arguments)
	case "create_dependency":
		result, err = h.toolHandler.handleCreateDependency(ctx, params.Arguments)
	case "list_node_dependencies":
		result, err = h.toolHandler.handleListNodeDependencies(ctx, params.Arguments)
	case "list_node_dependents":
		result, err = h.toolHandler.handleListNodeDependents(ctx, params.Arguments)
	case "delete_dependency":
		result, err = h.toolHandler.handleDeleteDependency(ctx, params.Arguments)
	case "filter_nodes_by_attributes":
		result, err = h.toolHandler.handleFilterNodesByAttributes(ctx, params.Arguments)
	case "get_node_with_attributes":
		result, err = h.toolHandler.handleGetNodeWithAttributes(ctx, params.Arguments)
	case "list_templates":
		result, err = h.toolHandler.handleListTemplates(ctx, params.Arguments)
	case "create_template":
		result, err = h.toolHandler.handleCreateTemplate(ctx, params.Arguments)
	case "get_template":
		result, err = h.toolHandler.handleGetTemplate(ctx, params.Arguments)
	case "update_template":
		result, err = h.toolHandler.handleUpdateTemplate(ctx, params.Arguments)
	case "delete_template":
		result, err = h.toolHandler.handleDeleteTemplate(ctx, params.Arguments)
	case "clone_template":
		result, err = h.toolHandler.handleCloneTemplate(ctx, params.Arguments)
	case "generate_template_scaffold":
		result, err = h.toolHandler.handleGenerateTemplateScaffold(ctx, params.Arguments)
	case "validate_template":
		result, err = h.toolHandler.handleValidateTemplate(ctx, params.Arguments)
	default:
		return h.createErrorResponse(req.ID, MethodNotFound, fmt.Sprintf("Tool not found: %s", params.Name), nil)
	}

	// Handle the response
	if err != nil {
		return h.createErrorResponse(req.ID, InternalError, "Tool execution failed", err.Error())
	}

	return h.createSuccessResponse(req.ID, result)
}

// handleGetServerInfo returns server information
func (h *MCPProtocolHandler) handleGetServerInfo(req *JSONRPCRequest) *JSONRPCResponse {
	result := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Server: %s v%s\nMode: %s\nProtocol: MCP %s",
					constants.MCPServerName,
					constants.DefaultServerVersion,
					h.mode,
					constants.MCPProtocolVersion,
				),
			},
		},
	}

	return h.createSuccessResponse(req.ID, result)
}

// handleResourcesList returns available resources (placeholder)
func (h *MCPProtocolHandler) handleResourcesList(req *JSONRPCRequest) *JSONRPCResponse {
	result := map[string]interface{}{
		"resources": []interface{}{},
	}

	return h.createSuccessResponse(req.ID, result)
}

// handleResourceRead reads a resource (placeholder)
func (h *MCPProtocolHandler) handleResourceRead(req *JSONRPCRequest) *JSONRPCResponse {
	return h.createErrorResponse(req.ID, MethodNotFound, "Resource reading not implemented", nil)
}

// createSuccessResponse creates a successful JSON-RPC response
func (h *MCPProtocolHandler) createSuccessResponse(id interface{}, result interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: constants.JSONRPCVersion,
		ID:      id,
		Result:  result,
	}
}

// createErrorResponse creates an error JSON-RPC response
func (h *MCPProtocolHandler) createErrorResponse(id interface{}, code int, message string, data interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: constants.JSONRPCVersion,
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}
