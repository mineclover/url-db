package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
	"url-db/internal/infrastructure/persistence/sqlite/mapper"
)

// sqliteNodeAttributeRepository implements the NodeAttributeRepository interface
type sqliteNodeAttributeRepository struct {
	db *sqlx.DB
}

// NewSQLiteNodeAttributeRepository creates a new SQLite node attribute repository
func NewSQLiteNodeAttributeRepository(db *sqlx.DB) repository.NodeAttributeRepository {
	return &sqliteNodeAttributeRepository{db: db}
}

// Create creates a new node attribute
func (r *sqliteNodeAttributeRepository) Create(ctx context.Context, nodeAttribute *entity.NodeAttribute) error {
	query := `
		INSERT INTO node_attributes (node_id, attribute_id, value, order_index, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		nodeAttribute.NodeID(),
		nodeAttribute.AttributeID(),
		nodeAttribute.Value(),
		nodeAttribute.OrderIndex(),
		nodeAttribute.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to create node attribute: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	
	nodeAttribute.SetID(int(id))
	return nil
}

// GetByNodeID retrieves all attributes for a specific node
func (r *sqliteNodeAttributeRepository) GetByNodeID(ctx context.Context, nodeID int) ([]*entity.NodeAttribute, error) {
	query := `
		SELECT id, node_id, attribute_id, value, order_index, created_at
		FROM node_attributes
		WHERE node_id = ?
		ORDER BY attribute_id, order_index
	`
	
	rows, err := r.db.QueryContext(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query node attributes: %w", err)
	}
	defer rows.Close()
	
	var attributes []*entity.NodeAttribute
	for rows.Next() {
		model := &mapper.NodeAttributeModel{}
		err := rows.Scan(
			&model.ID,
			&model.NodeID,
			&model.AttributeID,
			&model.Value,
			&model.OrderIndex,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node attribute: %w", err)
		}
		
		attribute := mapper.MapNodeAttributeModelToEntity(model)
		attributes = append(attributes, attribute)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate node attributes: %w", err)
	}
	
	return attributes, nil
}

// GetByNodeAndAttribute retrieves a specific attribute for a node
func (r *sqliteNodeAttributeRepository) GetByNodeAndAttribute(ctx context.Context, nodeID int, attributeID int) (*entity.NodeAttribute, error) {
	query := `
		SELECT id, node_id, attribute_id, value, order_index, created_at
		FROM node_attributes
		WHERE node_id = ? AND attribute_id = ?
	`
	
	model := &mapper.NodeAttributeModel{}
	err := r.db.GetContext(ctx, model, query, nodeID, attributeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get node attribute: %w", err)
	}
	
	return mapper.MapNodeAttributeModelToEntity(model), nil
}

// Update updates an existing node attribute
func (r *sqliteNodeAttributeRepository) Update(ctx context.Context, nodeAttribute *entity.NodeAttribute) error {
	query := `
		UPDATE node_attributes
		SET value = ?, order_index = ?
		WHERE node_id = ? AND attribute_id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		nodeAttribute.Value(),
		nodeAttribute.OrderIndex(),
		nodeAttribute.NodeID(),
		nodeAttribute.AttributeID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update node attribute: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("node attribute not found for update")
	}
	
	return nil
}

// Delete deletes a node attribute
func (r *sqliteNodeAttributeRepository) Delete(ctx context.Context, nodeID int, attributeID int) error {
	query := `DELETE FROM node_attributes WHERE node_id = ? AND attribute_id = ?`
	
	result, err := r.db.ExecContext(ctx, query, nodeID, attributeID)
	if err != nil {
		return fmt.Errorf("failed to delete node attribute: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("node attribute not found for deletion")
	}
	
	return nil
}

// DeleteAllByNode deletes all attributes for a node
func (r *sqliteNodeAttributeRepository) DeleteAllByNode(ctx context.Context, nodeID int) error {
	query := `DELETE FROM node_attributes WHERE node_id = ?`
	
	_, err := r.db.ExecContext(ctx, query, nodeID)
	if err != nil {
		return fmt.Errorf("failed to delete node attributes: %w", err)
	}
	
	return nil
}

// SetNodeAttributes sets multiple attributes for a node (replaces existing ones)
func (r *sqliteNodeAttributeRepository) SetNodeAttributes(ctx context.Context, nodeID int, attributes []*entity.NodeAttribute) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Delete existing attributes for the node
	_, err = tx.ExecContext(ctx, "DELETE FROM node_attributes WHERE node_id = ?", nodeID)
	if err != nil {
		return fmt.Errorf("failed to delete existing attributes: %w", err)
	}
	
	// Insert new attributes
	insertQuery := `
		INSERT INTO node_attributes (node_id, attribute_id, value, order_index, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	for _, attr := range attributes {
		_, err = tx.ExecContext(ctx, insertQuery,
			attr.NodeID(),
			attr.AttributeID(),
			attr.Value(),
			attr.OrderIndex(),
			attr.CreatedAt(),
		)
		if err != nil {
			return fmt.Errorf("failed to insert node attribute: %w", err)
		}
	}
	
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// GetNodesWithAttribute retrieves nodes that have a specific attribute with optional value filter
func (r *sqliteNodeAttributeRepository) GetNodesWithAttribute(ctx context.Context, attributeID int, value *string) ([]int, error) {
	var query string
	var args []interface{}
	
	if value != nil {
		query = `
			SELECT DISTINCT node_id
			FROM node_attributes
			WHERE attribute_id = ? AND value = ?
			ORDER BY node_id
		`
		args = []interface{}{attributeID, *value}
	} else {
		query = `
			SELECT DISTINCT node_id
			FROM node_attributes
			WHERE attribute_id = ?
			ORDER BY node_id
		`
		args = []interface{}{attributeID}
	}
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query nodes with attribute: %w", err)
	}
	defer rows.Close()
	
	var nodeIDs []int
	for rows.Next() {
		var nodeID int
		if err := rows.Scan(&nodeID); err != nil {
			return nil, fmt.Errorf("failed to scan node ID: %w", err)
		}
		nodeIDs = append(nodeIDs, nodeID)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate node IDs: %w", err)
	}
	
	return nodeIDs, nil
}