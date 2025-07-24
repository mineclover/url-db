package mcp

import (
	"encoding/json"
)

// LogLevel represents MCP log message levels
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warning"
	LogLevelError LogLevel = "error"
)

// LogMessage represents an MCP log message notification
type LogMessage struct {
	Level  LogLevel    `json:"level"`
	Data   interface{} `json:"data"`
	Logger string      `json:"logger,omitempty"`
}

// LogNotification represents a log message notification to the client
type LogNotification struct {
	JSONRPCVersion string     `json:"jsonrpc"`
	Method         string     `json:"method"`
	Params         LogMessage `json:"params"`
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
