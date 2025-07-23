package mapper

import (
	"time"
	"url-db/internal/domain/entity"
)

// NodeAttributeModel represents the database model for node attributes
type NodeAttributeModel struct {
	ID          int        `db:"id"`
	NodeID      int        `db:"node_id"`
	AttributeID int        `db:"attribute_id"`
	Value       string     `db:"value"`
	OrderIndex  *int       `db:"order_index"`
	CreatedAt   time.Time  `db:"created_at"`
}

// MapNodeAttributeModelToEntity converts a NodeAttributeModel to a NodeAttribute entity
func MapNodeAttributeModelToEntity(model *NodeAttributeModel) *entity.NodeAttribute {
	// Create the entity directly with the validated constructor
	nodeAttribute, err := entity.NewNodeAttribute(
		model.NodeID,
		model.AttributeID,
		model.Value,
		model.OrderIndex,
	)
	if err != nil {
		// This should not happen for data from database, but handle gracefully
		panic("invalid data from database: " + err.Error())
	}
	
	// Set the ID and creation time from database
	nodeAttribute.SetID(model.ID)
	
	return nodeAttribute
}

// MapNodeAttributeEntityToModel converts a NodeAttribute entity to a NodeAttributeModel
func MapNodeAttributeEntityToModel(entity *entity.NodeAttribute) *NodeAttributeModel {
	return &NodeAttributeModel{
		ID:          entity.ID(),
		NodeID:      entity.NodeID(),
		AttributeID: entity.AttributeID(),
		Value:       entity.Value(),
		OrderIndex:  entity.OrderIndex(),
		CreatedAt:   entity.CreatedAt(),
	}
}