package mapper

import (
	"url-db/internal/domain/entity"
	"url-db/internal/models"
)

// ToDomainEntity converts a database model to a domain entity
func ToDomainEntity(dbModel *models.Domain) *entity.Domain {
	if dbModel == nil {
		return nil
	}

	domain, _ := entity.NewDomain(dbModel.Name, dbModel.Description)
	if domain != nil {
		// Note: In a real implementation, you might need to set internal fields
		// This is a simplified version for demonstration
	}

	return domain
}

// ToDBModel converts a domain entity to a database model
func ToDBModel(entity *entity.Domain) *models.Domain {
	if entity == nil {
		return nil
	}

	return &models.Domain{
		Name:        entity.Name(),
		Description: entity.Description(),
		CreatedAt:   entity.CreatedAt(),
		UpdatedAt:   entity.UpdatedAt(),
	}
}
