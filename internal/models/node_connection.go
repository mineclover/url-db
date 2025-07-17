package models

import "time"

// NodeConnection represents a connection/relationship between two nodes
type NodeConnection struct {
	ID               int       `json:"id" db:"id"`
	SourceNodeID     int       `json:"source_node_id" db:"source_node_id"`
	TargetNodeID     int       `json:"target_node_id" db:"target_node_id"`
	RelationshipType string    `json:"relationship_type" db:"relationship_type"`
	Description      string    `json:"description,omitempty" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// CreateNodeConnectionRequest represents the request to create a node connection
type CreateNodeConnectionRequest struct {
	SourceNodeID     int    `json:"source_node_id" binding:"required"`
	TargetNodeID     int    `json:"target_node_id" binding:"required"`
	RelationshipType string `json:"relationship_type" binding:"required"`
	Description      string `json:"description,omitempty"`
}

// UpdateNodeConnectionRequest represents the request to update a node connection
type UpdateNodeConnectionRequest struct {
	RelationshipType string `json:"relationship_type,omitempty"`
	Description      string `json:"description,omitempty"`
}

// NodeConnectionListResponse represents a paginated list of node connections
type NodeConnectionListResponse struct {
	Connections []NodeConnection `json:"connections"`
	TotalCount  int              `json:"total_count"`
	Page        int              `json:"page"`
	Size        int              `json:"size"`
	TotalPages  int              `json:"total_pages"`
}

// NodeConnectionWithInfo represents a node connection with additional node information
type NodeConnectionWithInfo struct {
	NodeConnection
	SourceNodeURL   string `json:"source_node_url" db:"source_node_url"`
	TargetNodeURL   string `json:"target_node_url" db:"target_node_url"`
	SourceNodeTitle string `json:"source_node_title" db:"source_node_title"`
	TargetNodeTitle string `json:"target_node_title" db:"target_node_title"`
}

// Common relationship types
const (
	RelationshipTypeRelated  = "related"
	RelationshipTypeChild    = "child"
	RelationshipTypeParent   = "parent"
	RelationshipTypeNext     = "next"
	RelationshipTypePrevious = "previous"
	RelationshipTypeLinked   = "linked"
	RelationshipTypeCustom   = "custom"
)