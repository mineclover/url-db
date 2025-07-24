package repository

import (
	"context"
	"url-db/internal/domain/entity"
)

// TemplateRepository defines the interface for template persistence operations
type TemplateRepository interface {
	// Create creates a new template
	Create(ctx context.Context, template *entity.Template) error

	// GetByID retrieves a template by its ID
	GetByID(ctx context.Context, id int) (*entity.Template, error)

	// GetByName retrieves a template by its name and domain
	GetByName(ctx context.Context, name, domainName string) (*entity.Template, error)

	// List retrieves templates by domain with optional pagination
	List(ctx context.Context, domainName string, page, size int) ([]*entity.Template, int, error)

	// ListActive retrieves only active templates by domain
	ListActive(ctx context.Context, domainName string, page, size int) ([]*entity.Template, int, error)

	// ListByType retrieves templates by type and domain
	ListByType(ctx context.Context, domainName, templateType string, page, size int) ([]*entity.Template, int, error)

	// Update updates an existing template
	Update(ctx context.Context, template *entity.Template) error

	// Delete deletes a template by its ID
	Delete(ctx context.Context, id int) error

	// Exists checks if a template exists by name and domain
	Exists(ctx context.Context, name, domainName string) (bool, error)

	// GetBatch retrieves multiple templates by their IDs
	GetBatch(ctx context.Context, ids []int) ([]*entity.Template, error)

	// GetDomainByTemplateID retrieves the domain for a given template ID
	GetDomainByTemplateID(ctx context.Context, templateID int) (*entity.Domain, error)

	// FilterByAttributes retrieves templates by domain with attribute filters
	FilterByAttributes(ctx context.Context, domainName string, filters []AttributeFilter, page, size int) ([]*entity.Template, int, error)

	// Clone creates a copy of an existing template with a new name
	Clone(ctx context.Context, sourceID int, newName, newTitle, newDescription string) (*entity.Template, error)

	// SetActive updates the active status of a template
	SetActive(ctx context.Context, id int, active bool) error

	// GetTemplatesByDomainID retrieves templates by domain ID
	GetTemplatesByDomainID(ctx context.Context, domainID int, page, size int) ([]*entity.Template, int, error)

	// CountByType returns the count of templates by type in a domain
	CountByType(ctx context.Context, domainName string) (map[string]int, error)

	// GetRecentlyModified retrieves recently modified templates
	GetRecentlyModified(ctx context.Context, domainName string, limit int) ([]*entity.Template, error)

	// Search searches templates by name, title, or description
	Search(ctx context.Context, domainName, query string, page, size int) ([]*entity.Template, int, error)
}

// TemplateAttributeRepository defines the interface for template attribute operations
type TemplateAttributeRepository interface {
	// Basic CRUD operations
	CreateTemplateAttribute(ctx context.Context, templateAttribute *entity.TemplateAttribute) error
	GetTemplateAttributeByID(ctx context.Context, id int) (*entity.TemplateAttribute, error)
	UpdateTemplateAttribute(ctx context.Context, templateAttribute *entity.TemplateAttribute) error
	DeleteTemplateAttributeByID(ctx context.Context, id int) error

	// Template-specific queries
	GetTemplateAttributes(ctx context.Context, templateID int) ([]*entity.TemplateAttribute, error)
	GetTemplateAttributesWithDetails(ctx context.Context, templateID int) ([]*entity.TemplateAttributeWithDetails, error)
	DeleteAllTemplateAttributes(ctx context.Context, templateID int) error

	// SetTemplateAttribute sets a single attribute for a template
	SetTemplateAttribute(ctx context.Context, templateID, attributeID int, value string, orderIndex *int) error

	// SetTemplateAttributes sets multiple attributes for a template (replaces all)
	SetTemplateAttributes(ctx context.Context, templateID int, attributes []TemplateAttributeValue) error

	// DeleteTemplateAttribute deletes a specific attribute from a template
	DeleteTemplateAttribute(ctx context.Context, templateID, attributeID int) error

	// GetTemplatesByAttribute retrieves templates that have a specific attribute value
	GetTemplatesByAttribute(ctx context.Context, domainName, attributeName, attributeValue string) ([]*entity.Template, error)

	// GetTemplateAttributesByName retrieves specific attribute values for a template
	GetTemplateAttributesByName(ctx context.Context, templateID int, attributeNames []string) ([]*entity.TemplateAttribute, error)

	// Batch operations
	CreateTemplateAttributesBatch(ctx context.Context, templateAttributes []*entity.TemplateAttribute) error
	UpdateTemplateAttributesBatch(ctx context.Context, templateAttributes []*entity.TemplateAttribute) error

	// Search and filter operations
	FindTemplateAttributesByValue(ctx context.Context, value string) ([]*entity.TemplateAttribute, error)
	GetTemplateAttributeUsageStats(ctx context.Context, domainName string) (map[string]int, error)
}

// TemplateAttributeValue represents an attribute value to be set on a template
type TemplateAttributeValue struct {
	AttributeName string // Name of the attribute
	Value         string // Attribute value
	OrderIndex    *int   // Optional order index for ordered attributes
}

// TemplateFilter represents a filter condition for templates
type TemplateFilter struct {
	Name        *string // Filter by name (partial match)
	Type        *string // Filter by template type
	IsActive    *bool   // Filter by active status
	HasMetadata *bool   // Filter templates that have metadata
	CreatedFrom *string // Filter by creation date (ISO format)
	CreatedTo   *string // Filter by creation date (ISO format)
}

// TemplateSearchOptions represents search options for templates
type TemplateSearchOptions struct {
	Query           string         // Search query
	DomainName      string         // Domain to search in
	Filters         TemplateFilter // Additional filters
	SortBy          string         // Sort field: "name", "created_at", "updated_at"
	SortDirection   string         // Sort direction: "asc", "desc"
	IncludeInactive bool           // Include inactive templates
	Page            int            // Page number (1-based)
	Size            int            // Page size
}

// TemplateStats represents statistics about templates in a domain
type TemplateStats struct {
	TotalCount      int            `json:"total_count"`
	ActiveCount     int            `json:"active_count"`
	InactiveCount   int            `json:"inactive_count"`
	TypeCounts      map[string]int `json:"type_counts"`
	RecentlyAdded   int            `json:"recently_added"`   // Count of templates added in last 30 days
	RecentlyUpdated int            `json:"recently_updated"` // Count of templates updated in last 7 days
}

// TemplateRepositoryStats defines interface for template statistics
type TemplateRepositoryStats interface {
	// GetStats retrieves statistics for templates in a domain
	GetStats(ctx context.Context, domainName string) (*TemplateStats, error)

	// GetUsageStats retrieves usage statistics (if templates are referenced elsewhere)
	GetUsageStats(ctx context.Context, templateID int) (map[string]int, error)
}
