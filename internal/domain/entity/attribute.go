package entity

import (
	"errors"
	"time"
)

// Attribute represents a domain attribute that can be assigned to nodes
type Attribute struct {
	id            int
	name          string
	attributeType string
	description   string
	domainID      int
	createdAt     time.Time
	updatedAt     time.Time
}

// NewAttribute creates a new attribute entity with validation
func NewAttribute(name, attributeType, description string, domainID int) (*Attribute, error) {
	if name == "" {
		return nil, errors.New("attribute name cannot be empty")
	}

	if attributeType == "" {
		return nil, errors.New("attribute type cannot be empty")
	}

	if domainID <= 0 {
		return nil, errors.New("domain ID must be positive")
	}

	// Validate attribute type
	validTypes := map[string]bool{
		"tag":         true,
		"ordered_tag": true,
		"number":      true,
		"string":      true,
		"markdown":    true,
		"image":       true,
	}

	if !validTypes[attributeType] {
		return nil, errors.New("invalid attribute type")
	}

	now := time.Now()
	return &Attribute{
		name:          name,
		attributeType: attributeType,
		description:   description,
		domainID:      domainID,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

// Getters - ensuring immutability from outside
func (a *Attribute) ID() int              { return a.id }
func (a *Attribute) Name() string         { return a.name }
func (a *Attribute) Type() string         { return a.attributeType }
func (a *Attribute) Description() string  { return a.description }
func (a *Attribute) DomainID() int        { return a.domainID }
func (a *Attribute) CreatedAt() time.Time { return a.createdAt }
func (a *Attribute) UpdatedAt() time.Time { return a.updatedAt }

// Business logic methods
func (a *Attribute) UpdateDescription(description string) {
	a.description = description
	a.updatedAt = time.Now()
}

// SetID is used by infrastructure layer after persistence
func (a *Attribute) SetID(id int) {
	if a.id == 0 { // Only allow setting ID once
		a.id = id
	}
}

// SetTimestamps sets creation and update timestamps (for repository usage)
func (a *Attribute) SetTimestamps(createdAt, updatedAt time.Time) {
	a.createdAt = createdAt
	a.updatedAt = updatedAt
}

// IsValid checks if the attribute is in a valid state
func (a *Attribute) IsValid() bool {
	validTypes := map[string]bool{
		"tag":         true,
		"ordered_tag": true,
		"number":      true,
		"string":      true,
		"markdown":    true,
		"image":       true,
	}

	return a.name != "" &&
		a.domainID > 0 &&
		validTypes[a.attributeType]
}
