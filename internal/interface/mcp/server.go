package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"url-db/internal/constants"
	"url-db/internal/interface/setup"
)

// MCPServer represents the MCP JSON-RPC 2.0 server
type MCPServer struct {
	factory     *setup.ApplicationFactory
	toolHandler *MCPToolHandler
	mode        string
	reader      io.Reader
	writer      io.Writer
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(factory *setup.ApplicationFactory, mode string) *MCPServer {
	return &MCPServer{
		factory:     factory,
		toolHandler: NewMCPToolHandler(factory),
		mode:        mode,
		reader:      os.Stdin,
		writer:      os.Stdout,
	}
}

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC 2.0 error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC 2.0 error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// Start begins the MCP server operation
func (s *MCPServer) Start(ctx context.Context) error {
	switch s.mode {
	case constants.MCPModeStdio:
		return s.startStdioMode(ctx)
	case constants.MCPModeSSE:
		return s.startSSEMode(ctx)
	case constants.MCPModeHTTP:
		return s.startHTTPMode(ctx)
	default:
		return fmt.Errorf("unsupported MCP mode: %s", s.mode)
	}
}

// startStdioMode handles stdin/stdout communication
func (s *MCPServer) startStdioMode(ctx context.Context) error {
	decoder := json.NewDecoder(s.reader)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var req JSONRPCRequest
			if err := decoder.Decode(&req); err != nil {
				if err == io.EOF {
					return nil
				}
				s.sendError(nil, ParseError, "Parse error", err.Error())
				continue
			}

			s.handleRequest(ctx, &req)
		}
	}
}

// startSSEMode handles Server-Sent Events mode (placeholder)
func (s *MCPServer) startSSEMode(ctx context.Context) error {
	return fmt.Errorf("SSE mode not implemented yet")
}

// startHTTPMode handles HTTP mode (placeholder)
func (s *MCPServer) startHTTPMode(ctx context.Context) error {
	return fmt.Errorf("HTTP mode not implemented yet")
}

// handleRequest processes a JSON-RPC request
func (s *MCPServer) handleRequest(ctx context.Context, req *JSONRPCRequest) {
	if req.JSONRPC != constants.JSONRPCVersion {
		s.sendError(req.ID, InvalidRequest, "Invalid JSON-RPC version", nil)
		return
	}

	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolCall(ctx, req)
	case "resources/list":
		s.handleResourcesList(req)
	case "resources/read":
		s.handleResourceRead(req)
	case "notifications/initialized":
		// Client notification that initialization is complete
		// No response needed for notifications
		return
	default:
		s.sendError(req.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", req.Method), nil)
	}
}

// handleInitialize handles MCP initialization
func (s *MCPServer) handleInitialize(req *JSONRPCRequest) {
	result := map[string]interface{}{
		"protocolVersion": constants.MCPProtocolVersion,
		"capabilities": map[string]interface{}{
			"tools":     map[string]interface{}{},
			"resources": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    constants.MCPServerName,
			"version": constants.DefaultServerVersion,
		},
	}

	s.sendResult(req.ID, result)
}

// handleToolsList returns available MCP tools
func (s *MCPServer) handleToolsList(req *JSONRPCRequest) {
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
			"description": "List URLs in domain",
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
			"description": "Add URL to domain",
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
			"description": "Get URL details",
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
			"description": "Update URL title or description",
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
			"description": "Remove URL",
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
			"description": "Search by exact URL",
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
			"name":        "get_node_attributes",
			"description": "Get URL tags and attributes",
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
			"description": "Add or update URL tags",
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
			"description": "Get available tag types for domain",
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
			"description": "Define new tag type for domain",
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
			"description": "Get details of a specific domain attribute",
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
			"description": "Update domain attribute description",
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
			"description": "Remove domain attribute definition",
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
			"description": "Create dependency relationship between nodes",
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
			"description": "List what a node depends on",
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
			"description": "List what depends on a node",
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
			"description": "Remove dependency relationship",
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
			"description": "Filter nodes by attribute values",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "Domain name to filter nodes from"},
					"filters": map[string]interface{}{
						"type": "array",
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
			"description": "Get URL details with all attributes",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
	}

	result := map[string]interface{}{
		"tools": tools,
	}

	s.sendResult(req.ID, result)
}

// handleToolCall executes a tool call
func (s *MCPServer) handleToolCall(ctx context.Context, req *JSONRPCRequest) {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, InvalidParams, "Invalid tool call parameters", err.Error())
		return
	}

	var result interface{}
	var err error

	switch params.Name {
	case "get_server_info":
		s.handleGetServerInfo(req)
		return
	case "list_domains":
		result, err = s.toolHandler.handleListDomains(ctx, params.Arguments)
	case "create_domain":
		result, err = s.toolHandler.handleCreateDomain(ctx, params.Arguments)
	case "list_nodes":
		result, err = s.toolHandler.handleListNodes(ctx, params.Arguments)
	case "create_node":
		result, err = s.toolHandler.handleCreateNode(ctx, params.Arguments)
	case "get_node":
		result, err = s.toolHandler.handleGetNode(ctx, params.Arguments)
	case "update_node":
		result, err = s.toolHandler.handleUpdateNode(ctx, params.Arguments)
	case "delete_node":
		result, err = s.toolHandler.handleDeleteNode(ctx, params.Arguments)
	case "find_node_by_url":
		result, err = s.toolHandler.handleFindNodeByURL(ctx, params.Arguments)
	case "get_node_attributes":
		result, err = s.toolHandler.handleGetNodeAttributes(ctx, params.Arguments)
	case "set_node_attributes":
		result, err = s.toolHandler.handleSetNodeAttributes(ctx, params.Arguments)
	case "list_domain_attributes":
		result, err = s.toolHandler.handleListDomainAttributes(ctx, params.Arguments)
	case "create_domain_attribute":
		result, err = s.toolHandler.handleCreateDomainAttribute(ctx, params.Arguments)
	case "get_domain_attribute":
		result, err = s.toolHandler.handleGetDomainAttribute(ctx, params.Arguments)
	case "update_domain_attribute":
		result, err = s.toolHandler.handleUpdateDomainAttribute(ctx, params.Arguments)
	case "delete_domain_attribute":
		result, err = s.toolHandler.handleDeleteDomainAttribute(ctx, params.Arguments)
	case "create_dependency":
		result, err = s.toolHandler.handleCreateDependency(ctx, params.Arguments)
	case "list_node_dependencies":
		result, err = s.toolHandler.handleListNodeDependencies(ctx, params.Arguments)
	case "list_node_dependents":
		result, err = s.toolHandler.handleListNodeDependents(ctx, params.Arguments)
	case "delete_dependency":
		result, err = s.toolHandler.handleDeleteDependency(ctx, params.Arguments)
	case "filter_nodes_by_attributes":
		result, err = s.toolHandler.handleFilterNodesByAttributes(ctx, params.Arguments)
	case "get_node_with_attributes":
		result, err = s.toolHandler.handleGetNodeWithAttributes(ctx, params.Arguments)
	default:
		s.sendError(req.ID, MethodNotFound, fmt.Sprintf("Tool not found: %s", params.Name), nil)
		return
	}

	// Handle the response
	if err != nil {
		s.sendError(req.ID, InternalError, "Tool execution failed", err.Error())
		return
	}

	s.sendResult(req.ID, result)
}

// handleGetServerInfo returns server information
func (s *MCPServer) handleGetServerInfo(req *JSONRPCRequest) {
	result := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Server: %s v%s\nMode: %s\nProtocol: MCP %s",
					constants.MCPServerName,
					constants.DefaultServerVersion,
					s.mode,
					constants.MCPProtocolVersion,
				),
			},
		},
	}

	s.sendResult(req.ID, result)
}

// handleResourcesList returns available resources (placeholder)
func (s *MCPServer) handleResourcesList(req *JSONRPCRequest) {
	result := map[string]interface{}{
		"resources": []interface{}{},
	}

	s.sendResult(req.ID, result)
}

// handleResourceRead reads a resource (placeholder)
func (s *MCPServer) handleResourceRead(req *JSONRPCRequest) {
	s.sendError(req.ID, MethodNotFound, "Resource reading not implemented", nil)
}

// sendResult sends a successful JSON-RPC response
func (s *MCPServer) sendResult(id interface{}, result interface{}) {
	response := JSONRPCResponse{
		JSONRPC: constants.JSONRPCVersion,
		ID:      id,
		Result:  result,
	}

	s.sendResponse(&response)
}

// sendError sends an error JSON-RPC response
func (s *MCPServer) sendError(id interface{}, code int, message string, data interface{}) {
	response := JSONRPCResponse{
		JSONRPC: constants.JSONRPCVersion,
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	s.sendResponse(&response)
}

// sendResponse writes a JSON-RPC response to the output
func (s *MCPServer) sendResponse(response *JSONRPCResponse) {
	encoder := json.NewEncoder(s.writer)
	if err := encoder.Encode(response); err != nil {
		// Log to stderr to avoid corrupting JSON-RPC protocol
		fmt.Fprintf(os.Stderr, "Failed to send response: %v\n", err)
	}
}
