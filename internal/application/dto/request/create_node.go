package request

// CreateNodeRequest represents the request for creating a node
type CreateNodeRequest struct {
	DomainName  string `json:"domain_name" validate:"required"`
	URL         string `json:"url" validate:"required,max=2048"`
	Title       string `json:"title" validate:"max=255"`
	Description string `json:"description" validate:"max=1000"`
}
