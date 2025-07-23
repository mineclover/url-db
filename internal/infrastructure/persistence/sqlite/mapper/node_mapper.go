package mapper

import (
	"url-db/internal/domain/entity"
	"url-db/internal/models"
)

// ToNodeEntity converts a database model to a node entity
func ToNodeEntity(dbModel *models.Node) *entity.Node {
	if dbModel == nil {
		return nil
	}

	node, _ := entity.NewNode(dbModel.Content, "", dbModel.Title, dbModel.Description)
	if node != nil {
		node.SetID(dbModel.ID)
		// Note: Domain name would need to be set from domain lookup
	}

	return node
}

// ToNodeDBModel converts a node entity to a database model
func ToNodeDBModel(entity *entity.Node) *models.Node {
	if entity == nil {
		return nil
	}

	return &models.Node{
		ID:          entity.ID(),
		Content:     entity.URL(),
		Title:       entity.Title(),
		Description: entity.Description(),
		CreatedAt:   entity.CreatedAt(),
		UpdatedAt:   entity.UpdatedAt(),
	}
}
