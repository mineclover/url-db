# MCP Debugging Best Practices Implementation

This document demonstrates how URL-DB implements MCP (Model Context Protocol) debugging best practices based on the official documentation at https://modelcontextprotocol.io/docs/tools/debugging.

## Implementation Overview

URL-DB implements structured MCP logging with the following key components:

### 1. Structured Log Notifications (`internal/interface/mcp/types.go`)

```go
type LogLevel string

const (
    LogLevelDebug LogLevel = "debug"
    LogLevelInfo  LogLevel = "info"
    LogLevelWarn  LogLevel = "warning"
    LogLevelError LogLevel = "error"
)

type LogMessage struct {
    Level  LogLevel    `json:"level"`
    Data   interface{} `json:"data"`
    Logger string      `json:"logger,omitempty"`
}

type LogNotification struct {
    JSONRPCVersion string     `json:"jsonrpc"`
    Method         string     `json:"method"`
    Params         LogMessage `json:"params"`
}
```

### 2. MCP-Aware Server Logging (`internal/interface/mcp/server.go`)

The MCPServer provides structured logging capabilities:

```go
// Send structured log messages to MCP clients
func (s *MCPServer) SendLogMessage(level LogLevel, data interface{}, logger string) error

// Convenience methods for different log levels
func (s *MCPServer) LogDebug(data interface{}, logger string) error
func (s *MCPServer) LogInfo(data interface{}, logger string) error
func (s *MCPServer) LogWarn(data interface{}, logger string) error
func (s *MCPServer) LogError(data interface{}, logger string) error

// Enable/disable structured logging
func (s *MCPServer) EnableLogging(enabled bool)
func (s *MCPServer) IsLoggingEnabled() bool
```

### 3. Protocol-Compliant Logger (`internal/interface/mcp/logger.go`)

A comprehensive logger that respects MCP protocol requirements:

```go
type MCPLogger struct {
    server     *MCPServer
    component  string
    fallbackWriter io.Writer
}

// Creates MCP-aware logger for specific component
func NewMCPLogger(server *MCPServer, component string) *MCPLogger

// Structured logging methods
func (l *MCPLogger) Debug(message string)
func (l *MCPLogger) Info(message string)
func (l *MCPLogger) Warn(message string)
func (l *MCPLogger) Error(message string)
func (l *MCPLogger) Fatal(message string)

// Formatted logging methods
func (l *MCPLogger) Debugf(format string, args ...interface{})
func (l *MCPLogger) Infof(format string, args ...interface{})
func (l *MCPLogger) Warnf(format string, args ...interface{})
func (l *MCPLogger) Errorf(format string, args ...interface{})
func (l *MCPLogger) Fatalf(format string, args ...interface{})
```

## MCP Protocol Compliance

### 1. stdio Mode Protection

**Problem**: Local MCP servers should not log messages to stdout as this interferes with JSON-RPC protocol operation.

**Solution**: 
- In stdio mode, all logs are either sent via MCP notifications or discarded
- Fallback logging uses stderr exclusively, never stdout
- Fatal errors exit silently in stdio mode to avoid protocol disruption

```go
// In stdio mode, avoid stdout interference
if s.mode == constants.MCPModeStdio {
    return nil // Don't interfere with stdio JSON-RPC protocol
}
```

### 2. Smart Fallback Logging

When MCP structured logging is not available:
- SSE/HTTP modes: Log to stderr with timestamps and component info
- stdio mode: Silent operation to preserve JSON-RPC protocol integrity

```go
func (l *MCPLogger) fallbackLog(level LogLevel, message string) {
    // Only log to stderr if not in stdio mode
    if l.server != nil && l.server.GetMode() == constants.MCPModeStdio {
        return // Silent in stdio mode
    }
    
    timestamp := time.Now().Format(constants.DateTimeFormat)
    logLine := fmt.Sprintf("[%s] %s [%s] %s\n", timestamp, string(level), l.component, message)
    
    // Always write to stderr, never stdout
    l.fallbackWriter.Write([]byte(logLine))
}
```

### 3. Structured Log Notifications

MCP clients receive structured log messages via JSON-RPC notifications:

```json
{
    "jsonrpc": "2.0",
    "method": "notifications/message",
    "params": {
        "level": "info",
        "data": {
            "message": "MCP server initialized in stdio mode",
            "timestamp": "2025-07-24T10:30:45Z",
            "component": "main"
        },
        "logger": "main"
    }
}
```

## Usage Examples

### 1. Basic Usage in Application Code

```go
// Create MCP-aware logger
mcpLogger := mcp.NewMCPLogger(mcpServer, "database")

// Log various levels
mcpLogger.Info("Database connection established")
mcpLogger.Warnf("High connection count: %d", connectionCount)
mcpLogger.Error("Query failed: timeout occurred")

// Fatal errors (handles stdio mode appropriately)
mcpLogger.Fatal("Critical system failure")
```

### 2. Integration with Standard Libraries

```go
// Create standard logger that respects MCP requirements
standardLogger := mcpLogger.CreateStandardLogger()

// Use with libraries expecting log.Logger interface
database.SetLogger(standardLogger)
```

### 3. Main Application Integration

```go
// cmd/server/main.go
mcpServer, err := mcp.NewMCPServer(factory, *mcpMode)
if err != nil {
    if *mcpMode == constants.MCPModeStdio {
        // stdio mode: stderr output, silent exit
        fmt.Fprintf(os.Stderr, "Failed to create MCP server: %v\n", err)
        os.Exit(1)
    } else {
        log.Fatalf("Failed to create MCP server: %v", err)
    }
}

// Demonstrate structured logging
mcpLogger := mcp.NewMCPLogger(mcpServer, "main")
mcpLogger.Infof("MCP server initialized in %s mode", *mcpMode)
```

## Key Benefits

### 1. Protocol Compliance
- **stdio mode**: Zero stdout interference with JSON-RPC protocol
- **Non-stdio modes**: Rich logging to stderr for debugging
- **All modes**: Structured notifications to MCP clients when possible

### 2. Development Experience
- **Familiar API**: Similar to standard Go logging patterns
- **Component-based**: Each logger instance tied to specific component
- **Fallback safety**: Graceful degradation when MCP logging unavailable

### 3. Debugging Capabilities
- **Structured data**: Rich context in log messages with timestamps, components
- **Level filtering**: Standard debug/info/warn/error levels
- **Real-time**: Live log streaming to MCP clients in supported modes

## Configuration

### Environment Variables

```bash
# Enable/disable MCP structured logging (default: true)
export AUTO_CREATE_ATTRIBUTES=true
```

### Runtime Control

```go
// Enable/disable structured logging
mcpServer.EnableLogging(true)

// Check logging status
if mcpServer.IsLoggingEnabled() {
    mcpLogger.Debug("Debug logging is active")
}

// Custom fallback writer for testing
mcpLogger.SetFallbackWriter(customWriter)
```

## Testing

The implementation includes comprehensive testing support:

```go
// Test with custom writers
mcpLogger := mcp.NewMCPLogger(testServer, "test")
mcpLogger.SetFallbackWriter(testBuffer)

// Verify protocol compliance
assert.Equal(t, "stdio", mcpServer.GetMode())
assert.True(t, mcpServer.IsLoggingEnabled())
```

## Security Considerations

1. **Log Sanitization**: All log data should be sanitized before sending to clients
2. **Rate Limiting**: Consider rate limiting for high-frequency log messages
3. **Sensitive Data**: Never log passwords, tokens, or other sensitive information
4. **Client Trust**: Only send log notifications to trusted MCP clients

## Migration Guide

### From Standard Go Logging

```go
// Before
log.Printf("Server started on port %d", port)
log.Fatal("Database connection failed")

// After
mcpLogger := mcp.NewMCPLogger(mcpServer, "server")
mcpLogger.Infof("Server started on port %d", port)
mcpLogger.Fatal("Database connection failed")
```

### From Custom Logging

```go
// Before
customLogger.Log(LEVEL_INFO, "Operation completed")

// After
mcpLogger := mcp.NewMCPLogger(mcpServer, "operations")
mcpLogger.Info("Operation completed")
```

This implementation follows MCP debugging best practices while maintaining backward compatibility and providing a smooth developer experience.