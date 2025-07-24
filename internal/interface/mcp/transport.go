package mcp

import (
	"context"
	"io"
)

// Transport represents different communication transports for MCP server
type Transport interface {
	// Start begins the transport operation
	Start(ctx context.Context) error
	// Stop gracefully shuts down the transport
	Stop() error
	// SetRequestHandler sets the request handler for processing incoming requests
	SetRequestHandler(handler RequestHandler)
	// SetPort configures the port for network-based transports
	SetPort(port string)
	// GetName returns the transport mode name
	GetName() string
}

// RequestHandler processes JSON-RPC requests and returns responses
type RequestHandler func(ctx context.Context, req *JSONRPCRequest) *JSONRPCResponse

// ResponseWriter provides a unified interface for writing responses across different transports
type ResponseWriter interface {
	// WriteResponse writes a JSON-RPC response
	WriteResponse(response *JSONRPCResponse) error
	// WriteError writes an error response with the specified parameters
	WriteError(id interface{}, code int, message string, data interface{}) error
	// GetWriter returns the underlying io.Writer (for backward compatibility)
	GetWriter() io.Writer
}

// TransportConfig holds configuration for transport initialization
type TransportConfig struct {
	Mode   string
	Port   string
	Reader io.Reader
	Writer io.Writer
}
