package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultEndpoint = "http://localhost:8080/mcp"
	DefaultTimeout  = 30
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// RPCError represents a JSON-RPC 2.0 error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Bridge handles stdio to SSE conversion
type Bridge struct {
	endpoint string
	timeout  time.Duration
	debug    bool
	client   *http.Client
}

// NewBridge creates a new bridge instance
func NewBridge(endpoint string, timeout time.Duration, debug bool) *Bridge {
	return &Bridge{
		endpoint: endpoint,
		timeout:  timeout,
		debug:    debug,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// logDebug logs debug messages to stderr if debug mode is enabled
func (b *Bridge) logDebug(format string, args ...interface{}) {
	if b.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

// sendToSSE sends a request to the SSE server and parses the response
func (b *Bridge) sendToSSE(ctx context.Context, data *JSONRPCRequest) *JSONRPCResponse {
	b.logDebug("Sending request: %s", toJSON(data))

	// Marshal request
	jsonData, err := json.Marshal(data)
	if err != nil {
		b.logDebug("JSON marshal error: %v", err)
		return b.createErrorResponse(data.ID, -32700, fmt.Sprintf("JSON marshal error: %v", err))
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", b.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		b.logDebug("Request creation error: %v", err)
		return b.createErrorResponse(data.ID, -32603, fmt.Sprintf("Request creation error: %v", err))
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := b.client.Do(req)
	if err != nil {
		b.logDebug("HTTP request error: %v", err)
		if strings.Contains(err.Error(), "timeout") {
			return b.createErrorResponse(data.ID, -32603, fmt.Sprintf("Request timeout after %v", b.timeout))
		}
		if strings.Contains(err.Error(), "connection refused") {
			return b.createErrorResponse(data.ID, -32603, fmt.Sprintf("Cannot connect to SSE server at %s", b.endpoint))
		}
		return b.createErrorResponse(data.ID, -32603, fmt.Sprintf("HTTP request error: %v", err))
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		b.logDebug("Response read error: %v", err)
		return b.createErrorResponse(data.ID, -32603, fmt.Sprintf("Response read error: %v", err))
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		b.logDebug("HTTP error status: %d", resp.StatusCode)
		return b.createErrorResponse(data.ID, -32603, fmt.Sprintf("HTTP error: %d %s", resp.StatusCode, resp.Status))
	}

	// Parse SSE format
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			jsonStr := line[6:] // Remove "data: " prefix
			var result JSONRPCResponse
			if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
				b.logDebug("JSON decode error: %v", err)
				return b.createErrorResponse(data.ID, -32700, fmt.Sprintf("Invalid JSON in response: %v", err))
			}
			b.logDebug("Received response: %s", jsonStr)
			return &result
		}
	}

	// No SSE data found
	b.logDebug("No SSE data found in response")
	return b.createErrorResponse(data.ID, -32603, "No data found in SSE response")
}

// createErrorResponse creates an error response
func (b *Bridge) createErrorResponse(id interface{}, code int, message string) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
}

// run starts the bridge main loop
func (b *Bridge) run(ctx context.Context) error {
	b.logDebug("Bridge started, SSE endpoint: %s", b.endpoint)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse JSON-RPC request
		var request JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			b.logDebug("Invalid JSON from stdin: %v", err)
			errorResponse := b.createErrorResponse(nil, -32700, fmt.Sprintf("Parse error: %v", err))
			fmt.Println(toJSON(errorResponse))
			continue
		}

		// Forward to SSE server
		response := b.sendToSSE(ctx, &request)
		if response != nil {
			fmt.Println(toJSON(response))
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stdin scanner error: %v", err)
	}

	return nil
}

// toJSON converts a struct to JSON string
func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":null,"error":{"code":-32603,"message":"JSON marshal error: %v"}}`, err)
	}
	return string(data)
}

func main() {
	var (
		endpoint = flag.String("endpoint", getEnv("SSE_ENDPOINT", DefaultEndpoint), "SSE server endpoint")
		timeout  = flag.Int("timeout", getEnvInt("TIMEOUT", DefaultTimeout), "Request timeout in seconds")
		debug    = flag.Bool("debug", getEnv("DEBUG", "") != "", "Enable debug logging")
		help     = flag.Bool("help", false, "Show help message")
		version  = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *help {
		fmt.Println("MCP SSE Bridge - Converts stdio MCP protocol to HTTP SSE requests")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  mcp-bridge [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -endpoint string   SSE server endpoint (default: http://localhost:8080/mcp)")
		fmt.Println("  -timeout int       Request timeout in seconds (default: 30)")
		fmt.Println("  -debug            Enable debug logging")
		fmt.Println("  -help             Show help message")
		fmt.Println("  -version          Show version information")
		fmt.Println()
		fmt.Println("Environment Variables:")
		fmt.Println("  SSE_ENDPOINT      SSE server endpoint")
		fmt.Println("  TIMEOUT           Request timeout in seconds")
		fmt.Println("  DEBUG             Enable debug logging (any non-empty value)")
		os.Exit(0)
	}

	if *version {
		fmt.Println("MCP SSE Bridge v1.0.0")
		fmt.Println("Converts stdio MCP protocol to HTTP SSE requests")
		os.Exit(0)
	}

	// Create bridge
	bridge := NewBridge(*endpoint, time.Duration(*timeout)*time.Second, *debug)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		// Wait for interrupt signal (Ctrl+C)
		// In a real implementation, you'd handle os.Signal here
		// For simplicity, we'll just run until stdin closes
	}()

	// Run bridge
	if err := bridge.run(ctx); err != nil {
		bridge.logDebug("Bridge error: %v", err)
		os.Exit(1)
	}

	bridge.logDebug("Bridge stopped")
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets environment variable as integer with default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}