package mapper

import (
	"time"
	"url-db/internal/domain/entity"
)

// DatabaseTemplate represents the template as stored in database (raw SQL row)
type DatabaseTemplate struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	DomainID     int       `db:"domain_id"`
	TemplateData string    `db:"template_data"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ToTemplateEntity converts a database row to a template entity
func ToTemplateEntity(dbRow *DatabaseTemplate) *entity.Template {
	if dbRow == nil {
		return nil
	}

	template, err := entity.NewTemplate(
		dbRow.Name,
		dbRow.TemplateData,
		dbRow.Title,
		dbRow.Description,
		dbRow.DomainID,
	)
	if err != nil {
		return nil
	}

	// Set database-specific fields
	template.SetID(dbRow.ID)
	template.SetTimestamps(dbRow.CreatedAt, dbRow.UpdatedAt)
	template.SetActive(dbRow.IsActive)

	return template
}

// FromTemplateEntity converts a template entity to database row format
func FromTemplateEntity(template *entity.Template) *DatabaseTemplate {
	if template == nil {
		return nil
	}

	return &DatabaseTemplate{
		ID:           template.ID(),
		Name:         template.Name(),
		DomainID:     template.DomainID(),
		TemplateData: template.TemplateData(),
		Title:        template.Title(),
		Description:  template.Description(),
		IsActive:     template.IsActive(),
		CreatedAt:    template.CreatedAt(),
		UpdatedAt:    template.UpdatedAt(),
	}
}

// DatabaseTemplateWithDomain represents a template with domain information
type DatabaseTemplateWithDomain struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	DomainID     int       `db:"domain_id"`
	DomainName   string    `db:"domain_name"`
	TemplateData string    `db:"template_data"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ToTemplateEntityWithDomain converts a database row with domain to template entity
func ToTemplateEntityWithDomain(dbRow *DatabaseTemplateWithDomain) *entity.Template {
	if dbRow == nil {
		return nil
	}

	template, err := entity.NewTemplate(
		dbRow.Name,
		dbRow.TemplateData,
		dbRow.Title,
		dbRow.Description,
		dbRow.DomainID,
	)
	if err != nil {
		return nil
	}

	// Set database-specific fields
	template.SetID(dbRow.ID)
	template.SetTimestamps(dbRow.CreatedAt, dbRow.UpdatedAt)
	template.SetActive(dbRow.IsActive)

	return template
}
