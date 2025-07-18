package models

import (
	"time"
)

type Domain struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateDomainRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type UpdateDomainRequest struct {
	Description string `json:"description" binding:"max=1000"`
}

type DomainListResponse struct {
	Domains    []Domain `json:"domains"`
	TotalCount int      `json:"total_count"`
	Page       int      `json:"page"`
	Size       int      `json:"size"`
	TotalPages int      `json:"total_pages"`
}

// MCP-specific domain models
type MCPDomain struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	NodeCount   int       `json:"node_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateMCPDomainRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type MCPDomainListResponse struct {
	Domains []MCPDomain `json:"domains"`
}
