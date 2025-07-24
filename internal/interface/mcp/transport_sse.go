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

// SSETransport implements Transport for Server-Sent Events communication
type SSETransport struct {
	port           string
	server         *http.Server
	requestHandler RequestHandler
}

// NewSSETransport creates a new SSE transport
func NewSSETransport(config *TransportConfig) *SSETransport {
	port := config.Port
	if port == "" {
		port = strconv.Itoa(constants.DefaultPort)
	}

	return &SSETransport{
		port: port,
	}
}

// Start begins SSE server operation
func (t *SSETransport) Start(ctx context.Context) error {
	if t.requestHandler == nil {
		return fmt.Errorf("request handler not set")
	}

	mux := http.NewServeMux()

	// SSE endpoint for MCP communication
	mux.HandleFunc("/mcp", t.handleSSEEndpoint)

	// Health check endpoint
	mux.HandleFunc("/health", t.handleHealthCheck)

	t.server = &http.Server{
		Addr:    ":" + t.port,
		Handler: mux,
	}

	fmt.Printf("Starting MCP SSE server on port %s\n", t.port)
	fmt.Printf("SSE endpoint: http://localhost:%s/mcp\n", t.port)
	fmt.Printf("Health check: http://localhost:%s/health\n", t.port)

	return t.server.ListenAndServe()
}

// Stop gracefully shuts down the transport
func (t *SSETransport) Stop() error {
	if t.server != nil {
		return t.server.Shutdown(context.Background())
	}
	return nil
}

// SetRequestHandler sets the request handler
func (t *SSETransport) SetRequestHandler(handler RequestHandler) {
	t.requestHandler = handler
}

// SetPort sets the port for the SSE server
func (t *SSETransport) SetPort(port string) {
	t.port = port
}

// GetName returns the transport name
func (t *SSETransport) GetName() string {
	return constants.MCPModeSSE
}

// handleSSEEndpoint handles SSE connections for MCP communication
func (t *SSETransport) handleSSEEndpoint(w http.ResponseWriter, r *http.Request) {
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		t.setSSEHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests for initial setup
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set SSE headers
	t.setSSEHeaders(w)

	// Read the initial JSON-RPC request
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create SSE response writer and handle the request
	responseWriter := NewSSEResponseWriter(w)
	response := t.requestHandler(r.Context(), &req)
	
	if response != nil {
		if err := responseWriter.WriteResponse(response); err != nil {
			fmt.Printf("Failed to send SSE response: %v\n", err)
		}
	}
}

// handleHealthCheck handles health check requests
func (t *SSETransport) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"mode":   "sse",
		"server": constants.MCPServerName,
	})
}

// setSSEHeaders sets Server-Sent Events headers
func (t *SSETransport) setSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
}

// SSEResponseWriter implements ResponseWriter for Server-Sent Events
type SSEResponseWriter struct {
	responseWriter http.ResponseWriter
}

// NewSSEResponseWriter creates a new SSE response writer
func NewSSEResponseWriter(w http.ResponseWriter) *SSEResponseWriter {
	return &SSEResponseWriter{
		responseWriter: w,
	}
}

// WriteResponse writes a JSON-RPC response via SSE
func (w *SSEResponseWriter) WriteResponse(response *JSONRPCResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	// Send SSE message
	fmt.Fprintf(w.responseWriter, "data: %s\n\n", data)
	if f, ok := w.responseWriter.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// WriteError writes an error response via SSE
func (w *SSEResponseWriter) WriteError(id interface{}, code int, message string, data interface{}) error {
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

// GetWriter returns a custom writer that formats data for SSE
func (w *SSEResponseWriter) GetWriter() io.Writer {
	return &sseWriter{responseWriter: w.responseWriter}
}

// sseWriter implements io.Writer for SSE format
type sseWriter struct {
	responseWriter http.ResponseWriter
}

func (w *sseWriter) Write(p []byte) (n int, err error) {
	// Send SSE message
	fmt.Fprintf(w.responseWriter, "data: %s\n\n", p)
	if f, ok := w.responseWriter.(http.Flusher); ok {
		f.Flush()
	}
	return len(p), nil
}