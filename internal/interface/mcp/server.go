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
