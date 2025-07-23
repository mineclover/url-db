package entity

import (
	"errors"
	"time"
)

// Domain represents a domain entity in the business domain
type Domain struct {
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

	if len(name) > 255 {
		return nil, errors.New("domain name cannot exceed 255 characters")
	}

	if len(description) > 1000 {
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
func (d *Domain) Name() string         { return d.name }
func (d *Domain) Description() string  { return d.description }
func (d *Domain) CreatedAt() time.Time { return d.createdAt }
func (d *Domain) UpdatedAt() time.Time { return d.updatedAt }

// Business logic methods
func (d *Domain) UpdateDescription(description string) error {
	if len(description) > 1000 {
		return errors.New("domain description cannot exceed 1000 characters")
	}

	d.description = description
	d.updatedAt = time.Now()
	return nil
}

// IsValid checks if the domain is in a valid state
func (d *Domain) IsValid() bool {
	return d.name != "" && len(d.name) <= 255 && len(d.description) <= 1000
}
