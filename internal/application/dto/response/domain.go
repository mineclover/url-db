package response

import "time"

// DomainResponse represents the response for domain operations
type DomainResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DomainListResponse represents the response for domain list operations
type DomainListResponse struct {
	Domains    []DomainResponse `json:"domains"`
	TotalCount int              `json:"total_count"`
	Page       int              `json:"page"`
	Size       int              `json:"size"`
	TotalPages int              `json:"total_pages"`
}
