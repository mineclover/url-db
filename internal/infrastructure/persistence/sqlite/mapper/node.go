package mapper

import (
	"time"
	"url-db/internal/domain/entity"
)

// DatabaseNode represents the node as stored in database (raw SQL row)
type DatabaseNode struct {
	ID          int       `db:"id"`
	Content     string    `db:"content"`  // This is the URL field
	DomainID    int       `db:"domain_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ToNodeEntity converts a database row to a node entity
func ToNodeEntity(dbRow *DatabaseNode) *entity.Node {
	if dbRow == nil {
		return nil
	}

	node, err := entity.NewNode(dbRow.Content, dbRow.Title, dbRow.Description, dbRow.DomainID)
	if err != nil {
		return nil
	}

	// Set database-specific fields
	node.SetID(dbRow.ID)
	node.SetTimestamps(dbRow.CreatedAt, dbRow.UpdatedAt)

	return node
}

// FromNodeEntity converts a node entity to database row format
func FromNodeEntity(node *entity.Node) *DatabaseNode {
	if node == nil {
		return nil
	}

	return &DatabaseNode{
		ID:          node.ID(),
		Content:     node.Content(),
		DomainID:    node.DomainID(),
		Title:       node.Title(),
		Description: node.Description(),
		CreatedAt:   node.CreatedAt(),
		UpdatedAt:   node.UpdatedAt(),
	}
}