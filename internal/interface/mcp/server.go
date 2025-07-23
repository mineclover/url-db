package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"url-db/internal/constants"
	"url-db/internal/interface/setup"
)

// MCPServer represents the MCP JSON-RPC 2.0 server
type MCPServer struct {
	factory *setup.ApplicationFactory
	mode    string
	reader  io.Reader
	writer  io.Writer
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(factory *setup.ApplicationFactory, mode string) *MCPServer {
	return &MCPServer{
		factory: factory,
		mode:    mode,
		reader:  os.Stdin,
		writer:  os.Stdout,
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
		// TODO: Add more tools from mcp-tools.yaml specification
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

	switch params.Name {
	case "get_server_info":
		s.handleGetServerInfo(req)
	default:
		s.sendError(req.ID, MethodNotFound, fmt.Sprintf("Tool not found: %s", params.Name), nil)
	}
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
		log.Printf("Failed to send response: %v", err)
	}
}
