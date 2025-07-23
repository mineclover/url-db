package repository

import (
	"context"
	"url-db/internal/domain/entity"
)

// AttributeRepository defines the interface for attribute data access
type AttributeRepository interface {
	Create(ctx context.Context, attribute *entity.Attribute) error
	GetByID(ctx context.Context, id int) (*entity.Attribute, error)
	GetByName(ctx context.Context, domainID int, name string) (*entity.Attribute, error)
	ListByDomainID(ctx context.Context, domainID int) ([]*entity.Attribute, error)
	Update(ctx context.Context, attribute *entity.Attribute) error
	Delete(ctx context.Context, id int) error
}
