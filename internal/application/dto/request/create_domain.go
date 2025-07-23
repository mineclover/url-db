package request

// CreateDomainRequest represents the request for creating a domain
type CreateDomainRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
}
