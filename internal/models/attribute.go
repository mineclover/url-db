package models

import (
	"time"
)

type AttributeType string

const (
	AttributeTypeTag        AttributeType = "tag"
	AttributeTypeOrderedTag AttributeType = "ordered_tag"
	AttributeTypeNumber     AttributeType = "number"
	AttributeTypeString     AttributeType = "string"
	AttributeTypeMarkdown   AttributeType = "markdown"
	AttributeTypeImage      AttributeType = "image"
)

type Attribute struct {
	ID          int           `json:"id" db:"id"`
	DomainID    int           `json:"domain_id" db:"domain_id"`
	Name        string        `json:"name" db:"name"`
	Type        AttributeType `json:"type" db:"type"`
	Description string        `json:"description" db:"description"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
}

type CreateAttributeRequest struct {
	Name        string        `json:"name" binding:"required,max=255"`
	Type        AttributeType `json:"type" binding:"required"`
	Description string        `json:"description" binding:"max=1000"`
}

type UpdateAttributeRequest struct {
	Description string `json:"description" binding:"max=1000"`
}

type AttributeListResponse struct {
	Attributes []Attribute `json:"attributes"`
}

type NodeAttribute struct {
	ID          int       `json:"id" db:"id"`
	NodeID      int       `json:"node_id" db:"node_id"`
	AttributeID int       `json:"attribute_id" db:"attribute_id"`
	Value       string    `json:"value" db:"value"`
	OrderIndex  *int      `json:"order_index" db:"order_index"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateNodeAttributeRequest struct {
	AttributeID int    `json:"attribute_id" binding:"required"`
	Value       string `json:"value" binding:"required,max=2048"`
	OrderIndex  *int   `json:"order_index"`
}

type UpdateNodeAttributeRequest struct {
	Value      string `json:"value" binding:"required,max=2048"`
	OrderIndex *int   `json:"order_index"`
}

type NodeAttributeWithInfo struct {
	ID          int           `json:"id"`
	NodeID      int           `json:"node_id"`
	AttributeID int           `json:"attribute_id"`
	Name        string        `json:"name"`
	Type        AttributeType `json:"type"`
	Value       string        `json:"value"`
	OrderIndex  *int          `json:"order_index"`
	CreatedAt   time.Time     `json:"created_at"`
}

type NodeAttributeListResponse struct {
	Attributes []NodeAttributeWithInfo `json:"attributes"`
}

// MCP-specific models
type MCPAttribute struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type MCPNodeAttributeResponse struct {
	CompositeID string         `json:"composite_id"`
	Attributes  []MCPAttribute `json:"attributes"`
}

type SetMCPNodeAttributesRequest struct {
	Attributes []struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value" binding:"required"`
	} `json:"attributes" binding:"required"`
}