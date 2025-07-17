package models

import (
	"time"
)

type Node struct {
	ID          int       `json:"id" db:"id"`
	Content     string    `json:"content" db:"content"`
	DomainID    int       `json:"domain_id" db:"domain_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateNodeRequest struct {
	URL         string `json:"url" binding:"required,max=2048"`
	Title       string `json:"title" binding:"max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type UpdateNodeRequest struct {
	Title       string `json:"title" binding:"max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type FindNodeByURLRequest struct {
	URL string `json:"url" binding:"required,max=2048"`
}

type NodeListResponse struct {
	Nodes      []Node `json:"nodes"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	TotalPages int    `json:"total_pages"`
}

// MCP-specific models
type MCPNode struct {
	CompositeID string    `json:"composite_id"`
	URL         string    `json:"url"`
	DomainName  string    `json:"domain_name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateMCPNodeRequest struct {
	DomainName  string `json:"domain_name" binding:"required"`
	URL         string `json:"url" binding:"required,max=2048"`
	Title       string `json:"title" binding:"max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type FindMCPNodeRequest struct {
	DomainName string `json:"domain_name" binding:"required"`
	URL        string `json:"url" binding:"required,max=2048"`
}

type BatchMCPNodeRequest struct {
	CompositeIDs []string `json:"composite_ids" binding:"required"`
}

type BatchMCPNodeResponse struct {
	Nodes    []MCPNode `json:"nodes"`
	NotFound []string  `json:"not_found"`
}

type MCPNodeListResponse struct {
	Nodes      []MCPNode `json:"nodes"`
	TotalCount int       `json:"total_count"`
	Page       int       `json:"page"`
	Size       int       `json:"size"`
	TotalPages int       `json:"total_pages"`
}