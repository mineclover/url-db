package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// StdioServer implements MCP protocol over stdin/stdout with JSON-RPC 2.0
type StdioServer struct {
	service          MCPService
	reader           *bufio.Reader
	writer           io.Writer
	toolRegistry     *ToolRegistry
	resourceRegistry *ResourceRegistry
	initialized      bool
	shutdownOnce     sync.Once
	shutdown         chan struct{}
}

// NewStdioServer creates a new MCP stdio server
func NewStdioServer(service MCPService) *StdioServer {
	return &StdioServer{
		service:          service,
		reader:           bufio.NewReader(os.Stdin),
		writer:           os.Stdout,
		toolRegistry:     NewToolRegistry(service),
		resourceRegistry: NewResourceRegistry(service),
		initialized:      false,
		shutdown:         make(chan struct{}),
	}
}

// Start begins the stdio MCP session with JSON-RPC 2.0 protocol
func (s *StdioServer) Start() error {
	// log.Println("Starting MCP JSON-RPC 2.0 stdio server...") // Disabled for clean JSON output

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// log.Println("EOF received, ending session") // Disabled for clean JSON output
				return nil
			}
			return fmt.Errorf("error reading input: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse JSON-RPC request
		req, err := ParseJSONRPCRequest([]byte(line))
		if err != nil {
			// Send parse error response
			errorResp := NewJSONRPCError(nil, ParseError, "Parse error", err.Error())
			s.sendResponse(errorResp)
			continue
		}

		// Handle the request
		resp := s.handleJSONRPCRequest(req)
		if resp != nil {
			s.sendResponse(resp)
		}
	}
}

// handleJSONRPCRequest processes a JSON-RPC 2.0 request
func (s *StdioServer) handleJSONRPCRequest(req *JSONRPCRequest) *JSONRPCResponse {
	ctx := context.Background()

	switch req.Method {
	case "initialize":
		return s.handleInitialize(ctx, req)
	case "notifications/initialized":
		return s.handleInitialized(ctx, req)
	case "tools/list":
		return s.handleToolsList(ctx, req)
	case "tools/call":
		return s.handleToolsCall(ctx, req)
	case "resources/list":
		return s.handleResourcesList(ctx, req)
	case "resources/read":
		return s.handleResourcesRead(ctx, req)
	default:
		return NewJSONRPCError(req.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", req.Method), nil)
	}
}

// Shutdown gracefully shuts down the server
func (s *StdioServer) Shutdown() {
	s.shutdownOnce.Do(func() {
		close(s.shutdown)
	})
}

// sendResponse sends a JSON-RPC response
func (s *StdioServer) sendResponse(resp *JSONRPCResponse) {
	data, err := resp.ToJSON()
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	fmt.Fprintf(s.writer, "%s\n", string(data))
}

// handleInitialize handles the MCP initialize request
func (s *StdioServer) handleInitialize(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	// Parse initialize request
	var initReq InitializeRequest
	if req.Params != nil {
		paramsData, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(paramsData, &initReq); err != nil {
			return NewJSONRPCError(req.ID, InvalidParams, "Invalid initialize parameters", err.Error())
		}
	}

	// Create response
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Resources: &ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
		},
		ServerInfo: ServerInfo{
			Name:    "url-db-mcp-server",
			Version: "1.0.0",
		},
	}

	return NewJSONRPCResponse(req.ID, result)
}

// handleInitialized handles the MCP initialized notification
func (s *StdioServer) handleInitialized(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	s.initialized = true
	log.Println("MCP server initialized successfully")

	// Notification - no response needed for initialized
	if req.ID == nil {
		return nil
	}

	return NewJSONRPCResponse(req.ID, nil)
}

// handleToolsList handles the tools/list request
func (s *StdioServer) handleToolsList(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	if !s.initialized {
		return NewJSONRPCError(req.ID, InvalidRequest, "Server not initialized", nil)
	}

	tools := s.toolRegistry.GetTools()
	result := ToolsListResult{
		Tools: tools,
	}

	return NewJSONRPCResponse(req.ID, result)
}

// handleToolsCall handles the tools/call request
func (s *StdioServer) handleToolsCall(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	if !s.initialized {
		return NewJSONRPCError(req.ID, InvalidRequest, "Server not initialized", nil)
	}

	// Parse tool call request
	var callReq CallToolRequest
	if req.Params != nil {
		paramsData, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(paramsData, &callReq); err != nil {
			return NewJSONRPCError(req.ID, InvalidParams, "Invalid tool call parameters", err.Error())
		}
	}

	// Call the tool
	result, err := s.toolRegistry.CallTool(ctx, callReq.Name, callReq.Arguments)
	if err != nil {
		return NewJSONRPCError(req.ID, InternalError, "Tool execution failed", err.Error())
	}

	return NewJSONRPCResponse(req.ID, result)
}

// handleResourcesList handles the resources/list request
func (s *StdioServer) handleResourcesList(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	if !s.initialized {
		return NewJSONRPCError(req.ID, InvalidRequest, "Server not initialized", nil)
	}

	result, err := s.resourceRegistry.GetResources(ctx)
	if err != nil {
		return NewJSONRPCError(req.ID, InternalError, "Failed to get resources", err.Error())
	}

	return NewJSONRPCResponse(req.ID, result)
}

// handleResourcesRead handles the resources/read request
func (s *StdioServer) handleResourcesRead(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse {
	if !s.initialized {
		return NewJSONRPCError(req.ID, InvalidRequest, "Server not initialized", nil)
	}

	// Parse resource read request
	var readReq ReadResourceRequest
	if req.Params != nil {
		paramsData, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(paramsData, &readReq); err != nil {
			return NewJSONRPCError(req.ID, InvalidParams, "Invalid resource read parameters", err.Error())
		}
	}

	// Read the resource
	result, err := s.resourceRegistry.ReadResource(ctx, readReq.URI)
	if err != nil {
		return NewJSONRPCError(req.ID, InternalError, "Failed to read resource", err.Error())
	}

	return NewJSONRPCResponse(req.ID, result)
}
