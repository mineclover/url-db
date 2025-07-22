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

// DependencyMetadata represents additional metadata for a dependency
type DependencyMetadata struct {
	Relationship string `json:"relationship,omitempty"`
	Description  string `json:"description,omitempty"`
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
	DependencyTypeHard      = "hard"
	DependencyTypeSoft      = "soft"
	DependencyTypeReference = "reference"
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
