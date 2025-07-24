package mcp

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"url-db/internal/constants"
	"url-db/internal/interface/setup"
)

// MCPServer represents the refactored MCP JSON-RPC 2.0 server with transport abstraction
// This replaces the original MCPServer with a more modular, extensible architecture
type MCPServer struct {
	factory          *setup.ApplicationFactory
	protocolHandler  *MCPProtocolHandler
	transport        Transport
	transportFactory *TransportFactory
	mode             string
	port             string
	logWriter        io.Writer // For sending log notifications to client
	logEnabled       bool      // Whether to send log notifications
}

// NewMCPServer creates a new MCP server instance with transport abstraction
func NewMCPServer(factory *setup.ApplicationFactory, mode string) (*MCPServer, error) {
	transportFactory := NewTransportFactory()

	// Validate the mode
	if err := transportFactory.ValidateMode(mode); err != nil {
		return nil, err
	}

	server := &MCPServer{
		factory:          factory,
		protocolHandler:  NewMCPProtocolHandler(factory, mode),
		transportFactory: transportFactory,
		mode:             mode,
		port:             strconv.Itoa(constants.DefaultPort),
		logEnabled:       true, // Enable structured logging by default
	}

	// Create transport based on mode
	if err := server.initializeTransport(); err != nil {
		return nil, fmt.Errorf("failed to initialize transport: %w", err)
	}

	return server, nil
}

// SetPort sets the port for network-based transports
func (s *MCPServer) SetPort(port string) {
	s.port = port
	if s.transport != nil {
		s.transport.SetPort(port)
	}
}

// SetIOStreams sets custom input/output streams (useful for testing)
func (s *MCPServer) SetIOStreams(reader io.Reader, writer io.Writer) error {
	if s.mode != constants.MCPModeStdio {
		return fmt.Errorf("custom IO streams only supported for stdio mode")
	}

	// Recreate stdio transport with custom streams
	config := &TransportConfig{
		Mode:   s.mode,
		Port:   s.port,
		Reader: reader,
		Writer: writer,
	}

	transport, err := s.transportFactory.CreateTransport(config)
	if err != nil {
		return fmt.Errorf("failed to create transport with custom streams: %w", err)
	}

	s.transport = transport
	s.transport.SetRequestHandler(s.protocolHandler.HandleRequest)
	return nil
}

// Start begins the MCP server operation
func (s *MCPServer) Start(ctx context.Context) error {
	if s.transport == nil {
		return fmt.Errorf("transport not initialized")
	}

	// Don't log in stdio mode as it interferes with JSON-RPC communication
	if s.mode != "stdio" {
		fmt.Printf("Starting MCP server in %s mode\n", s.mode)
	}
	return s.transport.Start(ctx)
}

// Stop gracefully shuts down the MCP server
func (s *MCPServer) Stop() error {
	if s.transport == nil {
		return nil
	}

	// Don't log in stdio mode as it interferes with JSON-RPC communication
	if s.mode != "stdio" {
		fmt.Printf("Stopping MCP server (%s mode)\n", s.mode)
	}
	return s.transport.Stop()
}

// GetMode returns the current transport mode
func (s *MCPServer) GetMode() string {
	return s.mode
}

// GetSupportedModes returns all supported transport modes
func (s *MCPServer) GetSupportedModes() []string {
	return s.transportFactory.GetSupportedModes()
}

// SwitchMode dynamically switches to a different transport mode
func (s *MCPServer) SwitchMode(newMode string) error {
	// Validate the new mode
	if err := s.transportFactory.ValidateMode(newMode); err != nil {
		return err
	}

	// Stop current transport if running
	if s.transport != nil {
		if err := s.transport.Stop(); err != nil {
			return fmt.Errorf("failed to stop current transport: %w", err)
		}
	}

	// Update mode and recreate transport
	s.mode = newMode
	s.protocolHandler = NewMCPProtocolHandler(s.factory, newMode)

	if err := s.initializeTransport(); err != nil {
		return fmt.Errorf("failed to initialize new transport: %w", err)
	}

	fmt.Printf("Switched to %s mode\n", newMode)
	return nil
}

// initializeTransport creates and configures the transport based on current mode
func (s *MCPServer) initializeTransport() error {
	config := &TransportConfig{
		Mode:   s.mode,
		Port:   s.port,
		Reader: os.Stdin,  // Default for stdio
		Writer: os.Stdout, // Default for stdio
	}

	transport, err := s.transportFactory.CreateTransport(config)
	if err != nil {
		return err
	}

	// Set the request handler
	transport.SetRequestHandler(s.protocolHandler.HandleRequest)
	transport.SetPort(s.port)

	s.transport = transport
	return nil
}

// GetTransportInfo returns information about the current transport
func (s *MCPServer) GetTransportInfo() map[string]interface{} {
	info := map[string]interface{}{
		"mode": s.mode,
		"port": s.port,
	}

	if s.transport != nil {
		info["transport_name"] = s.transport.GetName()
		info["supported_modes"] = s.GetSupportedModes()
	}

	return info
}

// SendLogMessage sends a structured log message to the MCP client via notifications
func (s *MCPServer) SendLogMessage(level LogLevel, data interface{}, logger string) error {
	if !s.logEnabled || s.transport == nil {
		return nil // Silently ignore if logging disabled or transport not available
	}

	// Only send log notifications in modes that support it (not stdio)
	if s.mode == constants.MCPModeStdio {
		return nil // Don't interfere with stdio JSON-RPC protocol
	}

	// Create log notification
	logNotification := LogNotification{
		JSONRPCVersion: constants.JSONRPCVersion,
		Method:         constants.MCPLogNotificationMethod,
		Params: LogMessage{
			Level:  level,
			Data:   data,
			Logger: logger,
		},
	}

	// Send notification via transport (for SSE/HTTP modes)
	if responseWriter, ok := s.getResponseWriter(); ok {
		return s.sendNotification(responseWriter, &logNotification)
	}

	return nil
}

// EnableLogging enables or disables structured log notifications
func (s *MCPServer) EnableLogging(enabled bool) {
	s.logEnabled = enabled
}

// IsLoggingEnabled returns whether structured logging is enabled
func (s *MCPServer) IsLoggingEnabled() bool {
	return s.logEnabled
}

// LogDebug sends a debug level log message
func (s *MCPServer) LogDebug(data interface{}, logger string) error {
	return s.SendLogMessage(LogLevelDebug, data, logger)
}

// LogInfo sends an info level log message
func (s *MCPServer) LogInfo(data interface{}, logger string) error {
	return s.SendLogMessage(LogLevelInfo, data, logger)
}

// LogWarn sends a warning level log message
func (s *MCPServer) LogWarn(data interface{}, logger string) error {
	return s.SendLogMessage(LogLevelWarn, data, logger)
}

// LogError sends an error level log message
func (s *MCPServer) LogError(data interface{}, logger string) error {
	return s.SendLogMessage(LogLevelError, data, logger)
}

// getResponseWriter gets the response writer from the current transport
func (s *MCPServer) getResponseWriter() (ResponseWriter, bool) {
	// This is a simplified implementation - in a full implementation,
	// you would need to extract the ResponseWriter from each transport type
	// For now, return nil to indicate that notification sending is not implemented
	// for the current transport configuration
	return nil, false
}

// sendNotification sends a notification via the response writer
func (s *MCPServer) sendNotification(writer ResponseWriter, notification *LogNotification) error {
	// Convert notification to JSON-RPC response format
	// This would need to be implemented based on the specific transport
	// For now, this is a placeholder implementation
	return fmt.Errorf("notification sending not yet implemented")
}
