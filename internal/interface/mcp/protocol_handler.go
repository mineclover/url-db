package mcp

import (
	"context"
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
		// Check if this might be a direct tool call attempt
		toolNames := []string{"get_server_info", "list_domains", "create_domain", "list_nodes", "create_node", 
			"get_node", "update_node", "delete_node", "find_node_by_url", "scan_all_content",
			"get_node_attributes", "set_node_attributes", "list_domain_attributes", 
			"create_domain_attribute", "get_domain_attribute", "update_domain_attribute",
			"delete_domain_attribute"}
		
		for _, toolName := range toolNames {
			if req.Method == toolName {
				return h.createErrorResponse(req.ID, MethodNotFound, 
					fmt.Sprintf("Direct tool calls are not supported. Use 'tools/call' method with parameters: {\"name\":\"%s\",\"arguments\":{}}", req.Method), 
					map[string]interface{}{
						"hint": "Example: {\"method\":\"tools/call\",\"params\":{\"name\":\"" + req.Method + "\",\"arguments\":{}}}",
						"available_methods": []string{"initialize", "tools/list", "tools/call", "resources/list", "resources/read"},
					})
			}
		}
		
		return h.createErrorResponse(req.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", req.Method), 
			map[string]interface{}{
				"available_methods": []string{"initialize", "tools/list", "tools/call", "resources/list", "resources/read"},
			})
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

// handleToolsList returns available MCP tools with standard format
func (h *MCPProtocolHandler) handleToolsList(req *JSONRPCRequest) *JSONRPCResponse {
	toolDefs := GetToolDefinitions()
	tools := make([]map[string]interface{}, len(toolDefs))
	
	for i, def := range toolDefs {
		tools[i] = def.ToMap()
	}

	result := map[string]interface{}{
		"tools": tools,
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
