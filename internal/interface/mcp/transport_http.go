package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"url-db/internal/constants"
)

// HTTPTransport implements Transport for HTTP communication
type HTTPTransport struct {
	port           string
	server         *http.Server
	requestHandler RequestHandler
}

// NewHTTPTransport creates a new HTTP transport
func NewHTTPTransport(config *TransportConfig) *HTTPTransport {
	port := config.Port
	if port == "" {
		port = strconv.Itoa(constants.DefaultPort)
	}

	return &HTTPTransport{
		port: port,
	}
}

// Start begins HTTP server operation
func (t *HTTPTransport) Start(ctx context.Context) error {
	if t.requestHandler == nil {
		return fmt.Errorf("request handler not set")
	}

	mux := http.NewServeMux()

	// MCP endpoint for JSON-RPC communication
	mux.HandleFunc("/mcp", t.handleHTTPEndpoint)

	// Health check endpoint
	mux.HandleFunc("/health", t.handleHealthCheck)

	t.server = &http.Server{
		Addr:    ":" + t.port,
		Handler: mux,
	}

	fmt.Printf("Starting MCP HTTP server on port %s\n", t.port)
	fmt.Printf("MCP endpoint: http://localhost:%s/mcp\n", t.port)
	fmt.Printf("Health check: http://localhost:%s/health\n", t.port)

	return t.server.ListenAndServe()
}

// Stop gracefully shuts down the transport
func (t *HTTPTransport) Stop() error {
	if t.server != nil {
		return t.server.Shutdown(context.Background())
	}
	return nil
}

// SetRequestHandler sets the request handler
func (t *HTTPTransport) SetRequestHandler(handler RequestHandler) {
	t.requestHandler = handler
}

// SetPort sets the port for the HTTP server
func (t *HTTPTransport) SetPort(port string) {
	t.port = port
}

// GetName returns the transport name
func (t *HTTPTransport) GetName() string {
	return constants.MCPModeHTTP
}

// handleHTTPEndpoint handles HTTP requests for MCP communication
func (t *HTTPTransport) handleHTTPEndpoint(w http.ResponseWriter, r *http.Request) {
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		t.setCORSHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set response headers
	t.setCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")

	// Read and parse the JSON-RPC request
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWriter := NewHTTPResponseWriter(w)
		responseWriter.WriteError(nil, ParseError, "Parse error", err.Error())
		return
	}

	// Create response writer and handle the request
	responseWriter := NewHTTPResponseWriter(w)
	response := t.requestHandler(r.Context(), &req)
	
	if response != nil {
		if err := responseWriter.WriteResponse(response); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
}

// handleHealthCheck handles health check requests
func (t *HTTPTransport) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"mode":   "http",
		"server": constants.MCPServerName,
	})
}

// setCORSHeaders sets Cross-Origin Resource Sharing headers
func (t *HTTPTransport) setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// HTTPResponseWriter implements ResponseWriter for HTTP
type HTTPResponseWriter struct {
	responseWriter http.ResponseWriter
}

// NewHTTPResponseWriter creates a new HTTP response writer
func NewHTTPResponseWriter(w http.ResponseWriter) *HTTPResponseWriter {
	return &HTTPResponseWriter{
		responseWriter: w,
	}
}

// WriteResponse writes a JSON-RPC response to HTTP response
func (w *HTTPResponseWriter) WriteResponse(response *JSONRPCResponse) error {
	w.responseWriter.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w.responseWriter).Encode(response)
}

// WriteError writes an error response to HTTP response
func (w *HTTPResponseWriter) WriteError(id interface{}, code int, message string, data interface{}) error {
	response := &JSONRPCResponse{
		JSONRPC: constants.JSONRPCVersion,
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	return w.WriteResponse(response)
}

// GetWriter returns the underlying http.ResponseWriter as io.Writer
func (w *HTTPResponseWriter) GetWriter() io.Writer {
	return w.responseWriter
}