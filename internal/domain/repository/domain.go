package repository

import (
	"context"
	"url-db/internal/domain/entity"
)

// DomainRepository defines the interface for domain persistence operations
type DomainRepository interface {
	// Create creates a new domain
	Create(ctx context.Context, domain *entity.Domain) error

	// GetByID retrieves a domain by its ID
	GetByID(ctx context.Context, id int) (*entity.Domain, error)

	// GetByName retrieves a domain by its name
	GetByName(ctx context.Context, name string) (*entity.Domain, error)

	// List retrieves all domains with optional pagination
	List(ctx context.Context, page, size int) ([]*entity.Domain, int, error)

	// Update updates an existing domain
	Update(ctx context.Context, domain *entity.Domain) error

	// Delete deletes a domain by its name
	Delete(ctx context.Context, name string) error

	// Exists checks if a domain exists by name
	Exists(ctx context.Context, name string) (bool, error)
}
