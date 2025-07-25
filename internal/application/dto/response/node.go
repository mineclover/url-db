package response

import "time"

// NodeResponse represents the response for node operations
type NodeResponse struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	DomainName  string    `json:"domain_name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NodeListResponse represents the response for node list operations
type NodeListResponse struct {
	Nodes      []NodeResponse `json:"nodes"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	Size       int            `json:"size"`
	TotalPages int            `json:"total_pages"`
}

// NodeWithAttributes represents a node with its attributes for scanning operations
type NodeWithAttributes struct {
	ID          int             `json:"id"`
	Content     string          `json:"content"`
	Title       *string         `json:"title,omitempty"`
	Description *string         `json:"description,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Attributes  []AttributeValue `json:"attributes,omitempty"`
}
