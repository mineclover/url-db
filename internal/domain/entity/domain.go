package entity

import (
	"errors"
	"time"
	"url-db/internal/constants"
)

// Domain represents a domain entity in the business domain
type Domain struct {
	id          int
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewDomain creates a new domain entity with validation
func NewDomain(name, description string) (*Domain, error) {
	if name == "" {
		return nil, errors.New("domain name cannot be empty")
	}

	if len(name) > constants.MaxDomainNameLength {
		return nil, errors.New("domain name cannot exceed 255 characters")
	}

	if len(description) > constants.MaxDescriptionLength {
		return nil, errors.New("domain description cannot exceed 1000 characters")
	}

	now := time.Now()
	return &Domain{
		name:        name,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Getters - immutable from outside
func (d *Domain) ID() int              { return d.id }
func (d *Domain) Name() string         { return d.name }
func (d *Domain) Description() string  { return d.description }
func (d *Domain) CreatedAt() time.Time { return d.createdAt }
func (d *Domain) UpdatedAt() time.Time { return d.updatedAt }

// Business logic methods
func (d *Domain) UpdateDescription(description string) error {
	if len(description) > constants.MaxDescriptionLength {
		return errors.New("domain description cannot exceed 1000 characters")
	}

	d.description = description
	d.updatedAt = time.Now()
	return nil
}

// IsValid checks if the domain is in a valid state
func (d *Domain) IsValid() bool {
	return d.name != "" && len(d.name) <= constants.MaxDomainNameLength && len(d.description) <= constants.MaxDescriptionLength
}

// SetID sets the domain ID (for repository usage)
func (d *Domain) SetID(id int) {
	d.id = id
}

// SetTimestamps sets creation and update timestamps (for repository usage)
func (d *Domain) SetTimestamps(createdAt, updatedAt time.Time) {
	d.createdAt = createdAt
	d.updatedAt = updatedAt
}
