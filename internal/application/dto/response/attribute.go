package response

import "time"

type AttributeResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	DomainID    int       `json:"domain_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AttributeListResponse struct {
	Attributes []AttributeResponse `json:"attributes"`
	Total      int                 `json:"total"`
}
