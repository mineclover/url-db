package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NodeSubscription represents a subscription to node events
type NodeSubscription struct {
	ID                 int64            `db:"id" json:"id"`
	SubscriberService  string           `db:"subscriber_service" json:"subscriber_service"`
	SubscriberEndpoint *string          `db:"subscriber_endpoint" json:"subscriber_endpoint,omitempty"`
	SubscribedNodeID   int64            `db:"subscribed_node_id" json:"subscribed_node_id"`
	EventTypes         EventTypeList    `db:"event_types" json:"event_types"`
	FilterConditions   *FilterCondition `db:"filter_conditions" json:"filter_conditions,omitempty"`
	IsActive           bool             `db:"is_active" json:"is_active"`
	CreatedAt          time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time        `db:"updated_at" json:"updated_at"`
}

// EventTypeList represents a list of event types
type EventTypeList []string

// Scan implements sql.Scanner interface
func (e *EventTypeList) Scan(value interface{}) error {
	if value == nil {
		*e = EventTypeList{}
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
func (e EventTypeList) Value() (driver.Value, error) {
	if len(e) == 0 {
		return "[]", nil
	}
	data, err := json.Marshal(e)
	return string(data), err
}

// FilterCondition represents subscription filter conditions
type FilterCondition struct {
	AttributeFilters []AttributeFilter `json:"attribute_filters,omitempty"`
	ChangeFilters    *ChangeFilter     `json:"change_filters,omitempty"`
}

// AttributeFilter represents a filter on node attributes
type AttributeFilter struct {
	AttributeName string      `json:"attribute_name"`
	Operator      string      `json:"operator"`
	Value         interface{} `json:"value"`
}

// ChangeFilter represents filters on what changes to track
type ChangeFilter struct {
	Fields       []string `json:"fields,omitempty"`
	IgnoreFields []string `json:"ignore_fields,omitempty"`
}

// Scan implements sql.Scanner interface
func (f *FilterCondition) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), f)
	case []byte:
		return json.Unmarshal(v, f)
	default:
		return nil
	}
}

// Value implements driver.Valuer interface
func (f *FilterCondition) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
	}
	data, err := json.Marshal(f)
	return string(data), err
}

// CreateNodeSubscriptionRequest represents a request to create a subscription
type CreateNodeSubscriptionRequest struct {
	SubscriberService  string           `json:"subscriber_service" validate:"required"`
	SubscriberEndpoint *string          `json:"subscriber_endpoint,omitempty"`
	EventTypes         []string         `json:"event_types" validate:"required,min=1"`
	FilterConditions   *FilterCondition `json:"filter_conditions,omitempty"`
}

// UpdateNodeSubscriptionRequest represents a request to update a subscription
type UpdateNodeSubscriptionRequest struct {
	SubscriberEndpoint *string          `json:"subscriber_endpoint,omitempty"`
	EventTypes         []string         `json:"event_types,omitempty"`
	FilterConditions   *FilterCondition `json:"filter_conditions,omitempty"`
	IsActive           *bool            `json:"is_active,omitempty"`
}