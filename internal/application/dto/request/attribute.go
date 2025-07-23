package request

type CreateAttributeRequest struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description"`
	DomainID    int    `json:"domain_id"`
}

type UpdateAttributeRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}