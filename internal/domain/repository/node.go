package repository

import (
	"context"
	"url-db/internal/domain/entity"
)

// NodeRepository defines the interface for node persistence operations
type NodeRepository interface {
	// Create creates a new node
	Create(ctx context.Context, node *entity.Node) error

	// GetByID retrieves a node by its ID
	GetByID(ctx context.Context, id int) (*entity.Node, error)

	// GetByURL retrieves a node by its URL and domain
	GetByURL(ctx context.Context, url, domainName string) (*entity.Node, error)

	// List retrieves nodes by domain with optional pagination
	List(ctx context.Context, domainName string, page, size int) ([]*entity.Node, int, error)

	// Update updates an existing node
	Update(ctx context.Context, node *entity.Node) error

	// Delete deletes a node by its ID
	Delete(ctx context.Context, id int) error

	// Exists checks if a node exists by URL and domain
	Exists(ctx context.Context, url, domainName string) (bool, error)

	// GetBatch retrieves multiple nodes by their IDs
	GetBatch(ctx context.Context, ids []int) ([]*entity.Node, error)

	// GetDomainByNodeID retrieves the domain for a given node ID
	GetDomainByNodeID(ctx context.Context, nodeID int) (*entity.Domain, error)

	// FilterByAttributes retrieves nodes by domain with attribute filters
	FilterByAttributes(ctx context.Context, domainName string, filters []AttributeFilter, page, size int) ([]*entity.Node, int, error)

	// CountByDomain counts nodes in a domain
	CountByDomain(ctx context.Context, domainID int) (int, error)

	// GetByDomainFromCursor retrieves nodes starting from a cursor position
	GetByDomainFromCursor(ctx context.Context, domainID int, lastNodeID int, limit int) ([]*entity.Node, error)
}

// AttributeFilter represents a filter condition for node attributes
type AttributeFilter struct {
	Name     string // Attribute name
	Value    string // Attribute value
	Operator string // Comparison operator: "equals", "contains", "starts_with", "ends_with"
}
