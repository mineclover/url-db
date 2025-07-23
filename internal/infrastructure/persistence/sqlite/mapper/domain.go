package mapper

import (
	"time"
	"url-db/internal/domain/entity"
)

// DatabaseDomain represents the domain as stored in database (raw SQL row)
type DatabaseDomain struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ToDomainEntity converts a database row to a domain entity
func ToDomainEntity(dbRow *DatabaseDomain) *entity.Domain {
	if dbRow == nil {
		return nil
	}

	domain, err := entity.NewDomain(dbRow.Name, dbRow.Description)
	if err != nil {
		return nil
	}

	// Set database-specific fields
	domain.SetID(dbRow.ID)
	domain.SetTimestamps(dbRow.CreatedAt, dbRow.UpdatedAt)

	return domain
}

// FromDomainEntity converts a domain entity to database row format
func FromDomainEntity(domain *entity.Domain) *DatabaseDomain {
	if domain == nil {
		return nil
	}

	return &DatabaseDomain{
		ID:          domain.ID(),
		Name:        domain.Name(),
		Description: domain.Description(),
		CreatedAt:   domain.CreatedAt(),
		UpdatedAt:   domain.UpdatedAt(),
	}
}
