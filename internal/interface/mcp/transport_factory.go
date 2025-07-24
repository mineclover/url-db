package mcp

import (
	"fmt"

	"url-db/internal/constants"
)

// TransportFactory creates Transport instances based on mode configuration
type TransportFactory struct{}

// NewTransportFactory creates a new transport factory
func NewTransportFactory() *TransportFactory {
	return &TransportFactory{}
}

// CreateTransport creates a Transport instance based on the provided configuration
func (f *TransportFactory) CreateTransport(config *TransportConfig) (Transport, error) {
	switch config.Mode {
	case constants.MCPModeStdio:
		return NewStdioTransport(config), nil
	case constants.MCPModeHTTP:
		return NewHTTPTransport(config), nil
	case constants.MCPModeSSE:
		return NewSSETransport(config), nil
	default:
		return nil, fmt.Errorf("unsupported transport mode: %s", config.Mode)
	}
}

// GetSupportedModes returns a list of supported transport modes
func (f *TransportFactory) GetSupportedModes() []string {
	return []string{
		constants.MCPModeStdio,
		constants.MCPModeHTTP,
		constants.MCPModeSSE,
	}
}

// ValidateMode checks if the provided mode is supported
func (f *TransportFactory) ValidateMode(mode string) error {
	supportedModes := f.GetSupportedModes()
	for _, supportedMode := range supportedModes {
		if mode == supportedMode {
			return nil
		}
	}
	return fmt.Errorf("unsupported mode '%s'. Supported modes: %v", mode, supportedModes)
}
