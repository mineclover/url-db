package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

// CompositeKey represents a composite key value object
type CompositeKey struct {
	toolName   string
	domainName string
	id         string
}

// NewCompositeKey creates a new composite key with validation
func NewCompositeKey(toolName, domainName, id string) (*CompositeKey, error) {
	if toolName == "" {
		return nil, errors.New("tool name cannot be empty")
	}

	if domainName == "" {
		return nil, errors.New("domain name cannot be empty")
	}

	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	return &CompositeKey{
		toolName:   toolName,
		domainName: domainName,
		id:         id,
	}, nil
}

// ParseCompositeKey parses a composite key string
func ParseCompositeKey(compositeKey string) (*CompositeKey, error) {
	parts := strings.Split(compositeKey, ":")
	if len(parts) != 3 {
		return nil, errors.New("invalid composite key format, expected tool:domain:id")
	}

	return NewCompositeKey(parts[0], parts[1], parts[2])
}

// Getters
func (ck *CompositeKey) ToolName() string   { return ck.toolName }
func (ck *CompositeKey) DomainName() string { return ck.domainName }
func (ck *CompositeKey) ID() string         { return ck.id }

// String returns the string representation of the composite key
func (ck *CompositeKey) String() string {
	return fmt.Sprintf("%s:%s:%s", ck.toolName, ck.domainName, ck.id)
}

// IsValid checks if the composite key is in a valid state
func (ck *CompositeKey) IsValid() bool {
	return ck.toolName != "" && ck.domainName != "" && ck.id != ""
}
