package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"url-db/internal/constants"
)

// StdioTransport implements Transport for stdin/stdout communication
type StdioTransport struct {
	reader         io.Reader
	writer         ResponseWriter
	requestHandler RequestHandler
}

// NewStdioTransport creates a new stdio transport
func NewStdioTransport(config *TransportConfig) *StdioTransport {
	reader := config.Reader
	if reader == nil {
		reader = os.Stdin
	}

	writer := config.Writer
	if writer == nil {
		writer = os.Stdout
	}

	return &StdioTransport{
		reader: reader,
		writer: NewStdioResponseWriter(writer),
	}
}

// Start begins stdio communication
func (t *StdioTransport) Start(ctx context.Context) error {
	if t.requestHandler == nil {
		return fmt.Errorf("request handler not set")
	}

	decoder := json.NewDecoder(t.reader)

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
				t.writer.WriteError(nil, ParseError, "Parse error", err.Error())
				continue
			}

			response := t.requestHandler(ctx, &req)
			if response != nil {
				if err := t.writer.WriteResponse(response); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to send response: %v\n", err)
				}
			}
		}
	}
}

// Stop gracefully shuts down the transport
func (t *StdioTransport) Stop() error {
	// No cleanup needed for stdio
	return nil
}

// SetRequestHandler sets the request handler
func (t *StdioTransport) SetRequestHandler(handler RequestHandler) {
	t.requestHandler = handler
}

// SetPort is not applicable for stdio transport (no-op)
func (t *StdioTransport) SetPort(port string) {
	// No-op for stdio transport
}

// GetName returns the transport name
func (t *StdioTransport) GetName() string {
	return constants.MCPModeStdio
}

// StdioResponseWriter implements ResponseWriter for stdio
type StdioResponseWriter struct {
	writer io.Writer
}

// NewStdioResponseWriter creates a new stdio response writer
func NewStdioResponseWriter(writer io.Writer) *StdioResponseWriter {
	return &StdioResponseWriter{
		writer: writer,
	}
}

// WriteResponse writes a JSON-RPC response to stdout
func (w *StdioResponseWriter) WriteResponse(response *JSONRPCResponse) error {
	encoder := json.NewEncoder(w.writer)
	return encoder.Encode(response)
}

// WriteError writes an error response to stdout
func (w *StdioResponseWriter) WriteError(id interface{}, code int, message string, data interface{}) error {
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

// GetWriter returns the underlying io.Writer
func (w *StdioResponseWriter) GetWriter() io.Writer {
	return w.writer
}
