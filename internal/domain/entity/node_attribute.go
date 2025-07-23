package entity

import (
	"errors"
	"time"
	"url-db/internal/domain/attribute"
)

// NodeAttribute represents a node's attribute value
type NodeAttribute struct {
	id          int
	nodeID      int
	attributeID int
	value       string
	orderIndex  *int // Used for ordered attributes like ordered_tag
	createdAt   time.Time
}

// NewNodeAttribute creates a new node attribute
func NewNodeAttribute(nodeID int, attributeID int, value string, orderIndex *int) (*NodeAttribute, error) {
	if nodeID <= 0 {
		return nil, errors.New("node ID must be positive")
	}
	
	if attributeID <= 0 {
		return nil, errors.New("attribute ID must be positive")
	}
	
	if value == "" {
		return nil, errors.New("value cannot be empty")
	}
	
	return &NodeAttribute{
		nodeID:      nodeID,
		attributeID: attributeID,
		value:       value,
		orderIndex:  orderIndex,
		createdAt:   time.Now(),
	}, nil
}

// ValidatedNodeAttribute creates a new node attribute with validation
func ValidatedNodeAttribute(nodeID int, attributeID int, attrType attribute.AttributeType, value string, orderIndex *int, registry *attribute.ValidatorRegistry) (*NodeAttribute, error) {
	if nodeID <= 0 {
		return nil, errors.New("node ID must be positive")
	}
	
	if attributeID <= 0 {
		return nil, errors.New("attribute ID must be positive")
	}
	
	// Validate the attribute value using the registry
	result := registry.ValidateAttribute(attrType, value, orderIndex)
	if !result.IsValid {
		return nil, errors.New("attribute validation failed: " + result.ErrorMessage)
	}
	
	return &NodeAttribute{
		nodeID:      nodeID,
		attributeID: attributeID,
		value:       result.NormalizedValue, // Use normalized value
		orderIndex:  orderIndex,
		createdAt:   time.Now(),
	}, nil
}

// Getters
func (na *NodeAttribute) ID() int           { return na.id }
func (na *NodeAttribute) NodeID() int       { return na.nodeID }
func (na *NodeAttribute) AttributeID() int  { return na.attributeID }
func (na *NodeAttribute) Value() string     { return na.value }
func (na *NodeAttribute) OrderIndex() *int  { return na.orderIndex }
func (na *NodeAttribute) CreatedAt() time.Time { return na.createdAt }

// SetID sets the ID (used by repository after insertion)
func (na *NodeAttribute) SetID(id int) {
	na.id = id
}

// UpdateValue updates the attribute value with validation
func (na *NodeAttribute) UpdateValue(value string, orderIndex *int, attrType attribute.AttributeType, registry *attribute.ValidatorRegistry) error {
	// Validate the new value
	result := registry.ValidateAttribute(attrType, value, orderIndex)
	if !result.IsValid {
		return errors.New("attribute validation failed: " + result.ErrorMessage)
	}
	
	na.value = result.NormalizedValue
	na.orderIndex = orderIndex
	return nil
}