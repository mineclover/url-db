package nodeattributes

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"url-db/internal/models"
)

type Repository interface {
	Create(nodeAttribute *models.NodeAttribute) error
	GetByID(id int) (*models.NodeAttributeWithInfo, error)
	GetByNodeID(nodeID int) ([]models.NodeAttributeWithInfo, error)
	Update(id int, req *models.UpdateNodeAttributeRequest) (*models.NodeAttribute, error)
	Delete(id int) error
	DeleteByNodeIDAndAttributeID(nodeID, attributeID int) error
	GetByNodeIDAndAttributeID(nodeID, attributeID int) (*models.NodeAttributeWithInfo, error)
	GetMaxOrderIndex(nodeID, attributeID int) (int, error)
	ReorderAfterIndex(nodeID, attributeID, afterIndex int) error
	ValidateNodeAndAttributeDomain(nodeID, attributeID int) error
	GetAttributeType(attributeID int) (models.AttributeType, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(nodeAttribute *models.NodeAttribute) error {
	query := `
		INSERT INTO node_attributes (node_id, attribute_id, value, order_index)
		VALUES (:node_id, :attribute_id, :value, :order_index)
	`

	result, err := r.db.NamedExec(query, nodeAttribute)
	if err != nil {
		return fmt.Errorf("failed to create node attribute: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	nodeAttribute.ID = int(id)
	return nil
}

func (r *repository) GetByID(id int) (*models.NodeAttributeWithInfo, error) {
	query := `
		SELECT 
			na.id,
			na.node_id,
			na.attribute_id,
			a.name,
			a.type,
			na.value,
			na.order_index,
			na.created_at
		FROM node_attributes na
		JOIN attributes a ON na.attribute_id = a.id
		WHERE na.id = ?
	`

	var nodeAttr models.NodeAttributeWithInfo
	err := r.db.Get(&nodeAttr, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get node attribute by id: %w", err)
	}

	return &nodeAttr, nil
}

func (r *repository) GetByNodeID(nodeID int) ([]models.NodeAttributeWithInfo, error) {
	query := `
		SELECT 
			na.id,
			na.node_id,
			na.attribute_id,
			a.name,
			a.type,
			na.value,
			na.order_index,
			na.created_at
		FROM node_attributes na
		JOIN attributes a ON na.attribute_id = a.id
		WHERE na.node_id = ?
		ORDER BY a.name, na.order_index, na.created_at
	`

	var nodeAttrs []models.NodeAttributeWithInfo
	err := r.db.Select(&nodeAttrs, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node attributes by node id: %w", err)
	}

	return nodeAttrs, nil
}

func (r *repository) Update(id int, req *models.UpdateNodeAttributeRequest) (*models.NodeAttribute, error) {
	query := `
		UPDATE node_attributes 
		SET value = ?, order_index = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, req.Value, req.OrderIndex, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update node attribute: %w", err)
	}

	// Return updated node attribute
	selectQuery := `
		SELECT id, node_id, attribute_id, value, order_index, created_at
		FROM node_attributes
		WHERE id = ?
	`

	var nodeAttr models.NodeAttribute
	err = r.db.Get(&nodeAttr, selectQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated node attribute: %w", err)
	}

	return &nodeAttr, nil
}

func (r *repository) Delete(id int) error {
	query := `DELETE FROM node_attributes WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete node attribute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("node attribute not found")
	}

	return nil
}

func (r *repository) DeleteByNodeIDAndAttributeID(nodeID, attributeID int) error {
	query := `DELETE FROM node_attributes WHERE node_id = ? AND attribute_id = ?`

	result, err := r.db.Exec(query, nodeID, attributeID)
	if err != nil {
		return fmt.Errorf("failed to delete node attribute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("node attribute not found")
	}

	return nil
}

func (r *repository) GetByNodeIDAndAttributeID(nodeID, attributeID int) (*models.NodeAttributeWithInfo, error) {
	query := `
		SELECT 
			na.id,
			na.node_id,
			na.attribute_id,
			a.name,
			a.type,
			na.value,
			na.order_index,
			na.created_at
		FROM node_attributes na
		JOIN attributes a ON na.attribute_id = a.id
		WHERE na.node_id = ? AND na.attribute_id = ?
	`

	var nodeAttr models.NodeAttributeWithInfo
	err := r.db.Get(&nodeAttr, query, nodeID, attributeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get node attribute: %w", err)
	}

	return &nodeAttr, nil
}

func (r *repository) GetMaxOrderIndex(nodeID, attributeID int) (int, error) {
	query := `
		SELECT COALESCE(MAX(order_index), 0)
		FROM node_attributes
		WHERE node_id = ? AND attribute_id = ?
	`

	var maxIndex int
	err := r.db.Get(&maxIndex, query, nodeID, attributeID)
	if err != nil {
		return 0, fmt.Errorf("failed to get max order index: %w", err)
	}

	return maxIndex, nil
}

func (r *repository) ReorderAfterIndex(nodeID, attributeID, afterIndex int) error {
	query := `
		UPDATE node_attributes
		SET order_index = order_index + 1
		WHERE node_id = ? AND attribute_id = ? AND order_index > ?
	`

	_, err := r.db.Exec(query, nodeID, attributeID, afterIndex)
	if err != nil {
		return fmt.Errorf("failed to reorder node attributes: %w", err)
	}

	return nil
}

func (r *repository) ValidateNodeAndAttributeDomain(nodeID, attributeID int) error {
	query := `
		SELECT COUNT(*)
		FROM nodes n
		JOIN attributes a ON n.domain_id = a.domain_id
		WHERE n.id = ? AND a.id = ?
	`

	var count int
	err := r.db.Get(&count, query, nodeID, attributeID)
	if err != nil {
		return fmt.Errorf("failed to validate node and attribute domain: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("node and attribute must belong to the same domain")
	}

	return nil
}

func (r *repository) GetAttributeType(attributeID int) (models.AttributeType, error) {
	query := `SELECT type FROM attributes WHERE id = ?`

	var attributeType models.AttributeType
	err := r.db.Get(&attributeType, query, attributeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("attribute not found")
		}
		return "", fmt.Errorf("failed to get attribute type: %w", err)
	}

	return attributeType, nil
}
