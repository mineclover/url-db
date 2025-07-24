package repository

import (
	"context"
	"url-db/internal/domain/entity"
)

// NodeAttributeRepository defines the contract for node attribute persistence
type NodeAttributeRepository interface {
	// Create creates a new node attribute
	Create(ctx context.Context, nodeAttribute *entity.NodeAttribute) error

	// GetByNodeID retrieves all attributes for a specific node
	GetByNodeID(ctx context.Context, nodeID int) ([]*entity.NodeAttribute, error)

	// GetByNodeAndAttribute retrieves a specific attribute for a node
	GetByNodeAndAttribute(ctx context.Context, nodeID int, attributeID int) (*entity.NodeAttribute, error)

	// Update updates an existing node attribute
	Update(ctx context.Context, nodeAttribute *entity.NodeAttribute) error

	// Delete deletes a node attribute
	Delete(ctx context.Context, nodeID int, attributeID int) error

	// DeleteAllByNode deletes all attributes for a node
	DeleteAllByNode(ctx context.Context, nodeID int) error

	// SetNodeAttributes sets multiple attributes for a node (replaces existing ones)
	SetNodeAttributes(ctx context.Context, nodeID int, attributes []*entity.NodeAttribute) error

	// GetNodesWithAttribute retrieves nodes that have a specific attribute with optional value filter
	GetNodesWithAttribute(ctx context.Context, attributeID int, value *string) ([]int, error)
}
