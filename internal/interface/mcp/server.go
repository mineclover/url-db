package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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
	port        string
	server      *http.Server
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(factory *setup.ApplicationFactory, mode string) *MCPServer {
	return &MCPServer{
		factory:     factory,
		toolHandler: NewMCPToolHandler(factory),
		mode:        mode,
		reader:      os.Stdin,
		writer:      os.Stdout,
		port:        strconv.Itoa(constants.DefaultPort),
	}
}

// SetPort sets the port for HTTP/SSE mode
func (s *MCPServer) SetPort(port string) {
	s.port = port
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

// startSSEMode handles Server-Sent Events mode
func (s *MCPServer) startSSEMode(ctx context.Context) error {
	mux := http.NewServeMux()

	// SSE endpoint for MCP communication
	mux.HandleFunc("/mcp", s.handleSSEEndpoint)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
			"mode":   "sse",
			"server": constants.MCPServerName,
		})
	})

	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: mux,
	}

	fmt.Printf("Starting MCP SSE server on port %s\n", s.port)
	fmt.Printf("SSE endpoint: http://localhost:%s/mcp\n", s.port)
	fmt.Printf("Health check: http://localhost:%s/health\n", s.port)

	return s.server.ListenAndServe()
}

// handleSSEEndpoint handles SSE connections for MCP communication
func (s *MCPServer) handleSSEEndpoint(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests for initial setup
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the initial JSON-RPC request
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create a custom writer that sends to SSE
	sseWriter := &SSEWriter{responseWriter: w}
	s.writer = sseWriter

	// Handle the request and send response via SSE
	s.handleRequest(r.Context(), &req)
}

// SSEWriter implements io.Writer for SSE message sending
type SSEWriter struct {
	responseWriter http.ResponseWriter
}

func (w *SSEWriter) Write(p []byte) (n int, err error) {
	// Send SSE message
	fmt.Fprintf(w.responseWriter, "data: %s\n\n", p)
	if f, ok := w.responseWriter.(http.Flusher); ok {
		f.Flush()
	}
	return len(p), nil
}

// startHTTPMode handles HTTP mode
func (s *MCPServer) startHTTPMode(ctx context.Context) error {
	mux := http.NewServeMux()

	// MCP endpoint for JSON-RPC communication
	mux.HandleFunc("/mcp", s.handleHTTPEndpoint)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
			"mode":   "http",
			"server": constants.MCPServerName,
		})
	})

	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: mux,
	}

	fmt.Printf("Starting MCP HTTP server on port %s\n", s.port)
	fmt.Printf("MCP endpoint: http://localhost:%s/mcp\n", s.port)
	fmt.Printf("Health check: http://localhost:%s/health\n", s.port)

	return s.server.ListenAndServe()
}

// handleHTTPEndpoint handles HTTP requests for MCP communication
func (s *MCPServer) handleHTTPEndpoint(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Read and parse the JSON-RPC request
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendHTTPError(w, nil, ParseError, "Parse error", err.Error())
		return
	}

	// Create a custom writer that writes to HTTP response
	httpWriter := &HTTPWriter{responseWriter: w}
	s.writer = httpWriter

	// Handle the request
	s.handleRequest(r.Context(), &req)
}

// HTTPWriter implements io.Writer for HTTP response writing
type HTTPWriter struct {
	responseWriter http.ResponseWriter
}

func (w *HTTPWriter) Write(p []byte) (n int, err error) {
	return w.responseWriter.Write(p)
}

// sendHTTPError sends an error response via HTTP
func (s *MCPServer) sendHTTPError(w http.ResponseWriter, id interface{}, code int, message string, data interface{}) {
	response := JSONRPCResponse{
		JSONRPC: constants.JSONRPCVersion,
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
