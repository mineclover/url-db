package mcp

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"url-db/internal/constants"
)

// MCPLogger provides structured logging that respects MCP protocol requirements
type MCPLogger struct {
	server     *MCPServer
	component  string
	fallbackWriter io.Writer
}

// NewMCPLogger creates a new MCP-aware logger
func NewMCPLogger(server *MCPServer, component string) *MCPLogger {
	return &MCPLogger{
		server:     server,
		component:  component,
		fallbackWriter: os.Stderr, // Always use stderr for fallback to avoid stdio interference
	}
}

// Debug logs a debug message using MCP structured logging if available
func (l *MCPLogger) Debug(message string) {
	l.log(LogLevelDebug, message)
}

// Debugf logs a formatted debug message
func (l *MCPLogger) Debugf(format string, args ...interface{}) {
	l.log(LogLevelDebug, fmt.Sprintf(format, args...))
}

// Info logs an info message using MCP structured logging if available
func (l *MCPLogger) Info(message string) {
	l.log(LogLevelInfo, message)
}

// Infof logs a formatted info message
func (l *MCPLogger) Infof(format string, args ...interface{}) {
	l.log(LogLevelInfo, fmt.Sprintf(format, args...))
}

// Warn logs a warning message using MCP structured logging if available
func (l *MCPLogger) Warn(message string) {
	l.log(LogLevelWarn, message)
}

// Warnf logs a formatted warning message
func (l *MCPLogger) Warnf(format string, args ...interface{}) {
	l.log(LogLevelWarn, fmt.Sprintf(format, args...))
}

// Error logs an error message using MCP structured logging if available
func (l *MCPLogger) Error(message string) {
	l.log(LogLevelError, message)
}

// Errorf logs a formatted error message
func (l *MCPLogger) Errorf(format string, args ...interface{}) {
	l.log(LogLevelError, fmt.Sprintf(format, args...))
}

// Fatal logs a fatal error and exits the program appropriately for MCP mode
func (l *MCPLogger) Fatal(message string) {
	l.log(LogLevelError, message)
	l.handleFatal()
}

// Fatalf logs a formatted fatal error and exits the program
func (l *MCPLogger) Fatalf(format string, args ...interface{}) {
	l.log(LogLevelError, fmt.Sprintf(format, args...))
	l.handleFatal()
}

// log is the internal logging method that handles MCP vs fallback logging
func (l *MCPLogger) log(level LogLevel, message string) {
	// Try to send via MCP structured logging first
	if l.server != nil && l.server.IsLoggingEnabled() {
		logData := map[string]interface{}{
			"message":   message,
			"timestamp": time.Now().Format(constants.ISODateTimeFormat),
			"component": l.component,
		}
		
		if err := l.server.SendLogMessage(level, logData, l.component); err == nil {
			return // Successfully sent via MCP
		}
	}

	// Fallback to stderr logging (never use stdout to avoid MCP protocol interference)
	l.fallbackLog(level, message)
}

// fallbackLog provides stderr-based logging when MCP logging is not available
func (l *MCPLogger) fallbackLog(level LogLevel, message string) {
	// Only log to stderr if not in stdio mode to avoid JSON-RPC interference
	if l.server != nil && l.server.GetMode() == constants.MCPModeStdio {
		return // Silent in stdio mode
	}

	timestamp := time.Now().Format(constants.DateTimeFormat)
	logLine := fmt.Sprintf("[%s] %s [%s] %s\n", timestamp, string(level), l.component, message)
	
	// Always write to stderr, never stdout
	if _, err := l.fallbackWriter.Write([]byte(logLine)); err != nil {
		// If we can't even write to stderr, there's nothing more we can do
		return
	}
}

// handleFatal manages program termination in MCP-appropriate way
func (l *MCPLogger) handleFatal() {
	if l.server != nil && l.server.GetMode() == constants.MCPModeStdio {
		// In stdio mode, exit silently to avoid disrupting JSON-RPC protocol
		os.Exit(1)
	} else {
		// In other modes, use standard Go log.Fatal behavior
		log.Fatal("Fatal error occurred")
	}
}

// SetFallbackWriter allows customizing the fallback writer (useful for testing)
func (l *MCPLogger) SetFallbackWriter(writer io.Writer) {
	l.fallbackWriter = writer
}

// CreateStandardLogger creates a standard Go logger that respects MCP requirements
// This is useful for libraries that expect a standard log.Logger interface
func (l *MCPLogger) CreateStandardLogger() *log.Logger {
	// Create a logger that writes to stderr in non-stdio modes, nowhere in stdio mode
	var writer io.Writer
	if l.server != nil && l.server.GetMode() == constants.MCPModeStdio {
		writer = io.Discard // Discard logs in stdio mode
	} else {
		writer = l.fallbackWriter
	}
	
	return log.New(writer, fmt.Sprintf("[%s] ", l.component), log.LstdFlags)
}