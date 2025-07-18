package mcp

import (
	"strconv"
	"strings"

	"url-db/internal/compositekey"
)

// CompositeKeyAdapter adapts the compositekey.Service to the CompositeKeyService interface
type CompositeKeyAdapter struct {
	service *compositekey.Service
}

// NewCompositeKeyAdapter creates a new adapter
func NewCompositeKeyAdapter(service *compositekey.Service) *CompositeKeyAdapter {
	return &CompositeKeyAdapter{
		service: service,
	}
}

// Create creates a composite key from domain name and node ID
func (a *CompositeKeyAdapter) Create(domainName string, nodeID int) string {
	compositeKey, err := a.service.Create(domainName, nodeID)
	if err != nil {
		// Return empty string if creation fails
		return ""
	}
	return compositeKey
}

// Parse parses a composite key string and returns a CompositeKey
func (a *CompositeKeyAdapter) Parse(compositeID string) (*CompositeKey, error) {
	toolName, domainName, id, err := a.service.ParseComponents(compositeID)
	if err != nil {
		return nil, err
	}

	return &CompositeKey{
		ToolName:   toolName,
		DomainName: domainName,
		ID:         id,
	}, nil
}

// Validate validates a composite key
func (a *CompositeKeyAdapter) Validate(compositeID string) error {
	if !a.service.Validate(compositeID) {
		return NewInvalidCompositeKeyError(compositeID)
	}
	return nil
}

// Additional helper methods for compatibility
func (a *CompositeKeyAdapter) ExtractDomainName(compositeID string) (string, error) {
	return a.service.GetDomainName(compositeID)
}

func (a *CompositeKeyAdapter) ExtractNodeID(compositeID string) (int, error) {
	return a.service.GetID(compositeID)
}

func (a *CompositeKeyAdapter) ExtractToolName(compositeID string) (string, error) {
	return a.service.GetToolName(compositeID)
}

// Simple parsing fallback for basic format validation
func parseCompositeID(compositeID string) (string, string, int, error) {
	parts := strings.Split(compositeID, ":")
	if len(parts) != 3 {
		return "", "", 0, NewInvalidCompositeKeyError(compositeID)
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", 0, NewInvalidCompositeKeyError(compositeID)
	}

	return parts[0], parts[1], id, nil
}
