package mapper

import (
	"time"
	"url-db/internal/domain/entity"
)

// AttributeDBModel represents the attribute table structure
type AttributeDBModel struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Type        string    `db:"type"`
	Description string    `db:"description"`
	DomainID    int       `db:"domain_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ToAttributeEntity converts a database model to domain entity
func ToAttributeEntity(dbModel *AttributeDBModel) *entity.Attribute {
	if dbModel == nil {
		return nil
	}

	// Create entity using the business logic constructor
	attribute, err := entity.NewAttribute(
		dbModel.Name,
		dbModel.Type,
		dbModel.Description,
		dbModel.DomainID,
	)
	if err != nil {
		return nil
	}

	// Set the ID and timestamps from database
	attribute.SetID(dbModel.ID)
	attribute.SetTimestamps(dbModel.CreatedAt, dbModel.UpdatedAt)

	return attribute
}

// ToAttributeDBModel converts a domain entity to database model
func ToAttributeDBModel(entity *entity.Attribute) *AttributeDBModel {
	return &AttributeDBModel{
		ID:          entity.ID(),
		Name:        entity.Name(),
		Type:        entity.Type(),
		Description: entity.Description(),
		DomainID:    entity.DomainID(),
		CreatedAt:   entity.CreatedAt(),
		UpdatedAt:   entity.UpdatedAt(),
	}
}
