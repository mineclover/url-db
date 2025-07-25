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

type NodeAttributeResponse struct {
	AttributeName string `json:"attribute_name"`
	AttributeType string `json:"attribute_type"`
	Value         string `json:"value"`
	OrderIndex    *int   `json:"order_index,omitempty"`
}

// AttributeValue represents an attribute value for scanning operations
type AttributeValue struct {
	Name          string  `json:"name"`
	Value         string  `json:"value"`
	AttributeType *string `json:"attribute_type,omitempty"`
	OrderIndex    *int    `json:"order_index,omitempty"`
}
