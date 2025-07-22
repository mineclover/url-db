package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NodeDependency represents a dependency relationship between nodes
type NodeDependency struct {
	ID               int64               `db:"id" json:"id"`
	DependentNodeID  int64               `db:"dependent_node_id" json:"dependent_node_id"`
	DependencyNodeID int64               `db:"dependency_node_id" json:"dependency_node_id"`
	DependencyType   string              `db:"dependency_type" json:"dependency_type"`
	CascadeDelete    bool                `db:"cascade_delete" json:"cascade_delete"`
	CascadeUpdate    bool                `db:"cascade_update" json:"cascade_update"`
	Metadata         *DependencyMetadata `db:"metadata" json:"metadata,omitempty"`
	CreatedAt        time.Time           `db:"created_at" json:"created_at"`
}

// NodeDependencyV2 represents an enhanced dependency relationship with advanced features
type NodeDependencyV2 struct {
	ID                int64                 `db:"id" json:"id"`
	DependentNodeID   int64                 `db:"dependent_node_id" json:"dependent_node_id"`
	DependencyNodeID  int64                 `db:"dependency_node_id" json:"dependency_node_id"`
	DependencyType    string                `db:"dependency_type" json:"dependency_type"`
	Category          string                `db:"category" json:"category"`
	Strength          int                   `db:"strength" json:"strength"` // 0-100
	Priority          int                   `db:"priority" json:"priority"` // 0-100
	CascadeDelete     bool                  `db:"cascade_delete" json:"cascade_delete"`
	CascadeUpdate     bool                  `db:"cascade_update" json:"cascade_update"`
	Metadata          *DependencyMetadataV2 `db:"metadata" json:"metadata,omitempty"`
	VersionConstraint *string               `db:"version_constraint" json:"version_constraint,omitempty"`
	IsRequired        bool                  `db:"is_required" json:"is_required"`
	IsActive          bool                  `db:"is_active" json:"is_active"`
	ValidFrom         time.Time             `db:"valid_from" json:"valid_from"`
	ValidUntil        *time.Time            `db:"valid_until" json:"valid_until,omitempty"`
	CreatedAt         time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time             `db:"updated_at" json:"updated_at"`
	CreatedBy         *string               `db:"created_by" json:"created_by,omitempty"`
}

// DependencyMetadata represents additional metadata for a dependency
type DependencyMetadata struct {
	Relationship string `json:"relationship,omitempty"`
	Description  string `json:"description,omitempty"`
}

// DependencyMetadataV2 represents enhanced metadata with type-specific fields
type DependencyMetadataV2 struct {
	Relationship       string                 `json:"relationship,omitempty"`
	Description        string                 `json:"description,omitempty"`
	HealthCheckURL     string                 `json:"health_check_url,omitempty"`
	SyncFrequency      string                 `json:"sync_frequency,omitempty"`
	RetryPolicy        map[string]interface{} `json:"retry_policy,omitempty"`
	StartupOrder       int                    `json:"startup_order,omitempty"`
	FallbackBehavior   string                 `json:"fallback_behavior,omitempty"`
	ConflictResolution string                 `json:"conflict_resolution,omitempty"`
	CustomFields       map[string]interface{} `json:"custom_fields,omitempty"`
}

// Scan implements sql.Scanner interface
func (m *DependencyMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), m)
	case []byte:
		return json.Unmarshal(v, m)
	default:
		return nil
	}
}

// Value implements driver.Valuer interface
func (m *DependencyMetadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	data, err := json.Marshal(m)
	return string(data), err
}

// DependencyType constants
const (
	// Structural dependency types
	DependencyTypeHard      = "hard"
	DependencyTypeSoft      = "soft"
	DependencyTypeReference = "reference"

	// Behavioral dependency types
	DependencyTypeRuntime  = "runtime"
	DependencyTypeCompile  = "compile"
	DependencyTypeOptional = "optional"

	// Data dependency types
	DependencyTypeSync  = "sync"
	DependencyTypeAsync = "async"

	// Dependency categories
	CategoryStructural = "structural"
	CategoryBehavioral = "behavioral"
	CategoryData       = "data"

	// Impact levels
	ImpactLevelCritical = "critical"
	ImpactLevelHigh     = "high"
	ImpactLevelMedium   = "medium"
	ImpactLevelLow      = "low"
)

// CreateNodeDependencyRequest represents a request to create a dependency
type CreateNodeDependencyRequest struct {
	DependencyNodeID int64               `json:"dependency_node_id" validate:"required"`
	DependencyType   string              `json:"dependency_type" validate:"required,oneof=hard soft reference"`
	CascadeDelete    bool                `json:"cascade_delete"`
	CascadeUpdate    bool                `json:"cascade_update"`
	Metadata         *DependencyMetadata `json:"metadata,omitempty"`
}

// NodeEvent represents an event that occurred on a node
type NodeEvent struct {
	ID          int64      `db:"id" json:"id"`
	NodeID      int64      `db:"node_id" json:"node_id"`
	EventType   string     `db:"event_type" json:"event_type"`
	EventData   *EventData `db:"event_data" json:"event_data,omitempty"`
	OccurredAt  time.Time  `db:"occurred_at" json:"occurred_at"`
	ProcessedAt *time.Time `db:"processed_at" json:"processed_at,omitempty"`
}

// EventData represents the data associated with an event
type EventData struct {
	EventID   string                 `json:"event_id"`
	NodeID    int64                  `json:"node_id"`
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Changes   *EventChanges          `json:"changes,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// EventChanges represents before/after state in an event
type EventChanges struct {
	Before map[string]interface{} `json:"before,omitempty"`
	After  map[string]interface{} `json:"after,omitempty"`
}

// Scan implements sql.Scanner interface
func (e *EventData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), e)
	case []byte:
		return json.Unmarshal(v, e)
	default:
		return nil
	}
}

// Value implements driver.Valuer interface
func (e *EventData) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	data, err := json.Marshal(e)
	return string(data), err
}

// Event type constants
const (
	EventTypeCreated           = "created"
	EventTypeUpdated           = "updated"
	EventTypeDeleted           = "deleted"
	EventTypeAttributeChanged  = "attribute_changed"
	EventTypeConnectionChanged = "connection_changed"
)

// DependencyGraph represents a node's dependency relationships
type DependencyGraph struct {
	NodeID       int64                  `json:"node_id"`
	Dependencies []DependencyNode       `json:"dependencies"`
	Dependents   []DependencyNode       `json:"dependents"`
	Depth        int                    `json:"depth"`
	HasCircular  bool                   `json:"has_circular"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// DependencyNode represents a node in the dependency graph
type DependencyNode struct {
	NodeID            int64                  `json:"node_id"`
	CompositeID       string                 `json:"composite_id"`
	Title             string                 `json:"title"`
	DependencyType    string                 `json:"dependency_type"`
	Category          string                 `json:"category"`
	Strength          int                    `json:"strength"`
	Priority          int                    `json:"priority"`
	IsRequired        bool                   `json:"is_required"`
	VersionConstraint *string                `json:"version_constraint,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	Children          []DependencyNode       `json:"children,omitempty"`
}

// ImpactAnalysisResult contains the results of dependency impact analysis
type ImpactAnalysisResult struct {
	SourceNodeID    int64          `json:"source_node_id"`
	SourceComposite string         `json:"source_composite_id"`
	ImpactType      string         `json:"impact_type"`
	AffectedNodes   []AffectedNode `json:"affected_nodes"`
	ImpactScore     int            `json:"impact_score"` // 0-100
	CascadeDepth    int            `json:"cascade_depth"`
	EstimatedTime   string         `json:"estimated_time"`
	Warnings        []string       `json:"warnings,omitempty"`
	Recommendations []string       `json:"recommendations,omitempty"`
}

// AffectedNode represents a node affected by a dependency change
type AffectedNode struct {
	NodeID       int64   `json:"node_id"`
	CompositeID  string  `json:"composite_id"`
	Title        string  `json:"title"`
	ImpactLevel  string  `json:"impact_level"` // 'critical', 'high', 'medium', 'low'
	Reason       string  `json:"reason"`
	ActionNeeded string  `json:"action_needed"`
	Path         []int64 `json:"path,omitempty"` // Path from source to this node
}

// CircularDependency represents a circular dependency path
type CircularDependency struct {
	Path        []int64  `json:"path"`
	NodeDetails []string `json:"node_details"`
	Strength    int      `json:"strength"` // Weakest link in the cycle
}

// DependencyValidationResult contains validation results
type DependencyValidationResult struct {
	IsValid  bool                 `json:"is_valid"`
	Errors   []string             `json:"errors,omitempty"`
	Warnings []string             `json:"warnings,omitempty"`
	Cycles   []CircularDependency `json:"cycles,omitempty"`
}

// DependencyTypeConfig represents configuration for a dependency type
type DependencyTypeConfig struct {
	TypeName           string                 `json:"type_name"`
	Category           string                 `json:"category"`
	CascadeDelete      bool                   `json:"cascade_delete"`
	CascadeUpdate      bool                   `json:"cascade_update"`
	ValidationRequired bool                   `json:"validation_required"`
	DefaultStrength    int                    `json:"default_strength"`
	DefaultPriority    int                    `json:"default_priority"`
	MetadataSchema     map[string]interface{} `json:"metadata_schema,omitempty"`
	Description        string                 `json:"description"`
}

// DependencyRule represents a validation rule for dependencies
type DependencyRule struct {
	ID         int64                  `db:"id" json:"id"`
	DomainID   *int64                 `db:"domain_id" json:"domain_id,omitempty"`
	RuleName   string                 `db:"rule_name" json:"rule_name"`
	RuleType   string                 `db:"rule_type" json:"rule_type"`
	RuleConfig map[string]interface{} `db:"rule_config" json:"rule_config"`
	IsActive   bool                   `db:"is_active" json:"is_active"`
	CreatedAt  time.Time              `db:"created_at" json:"created_at"`
}
