package models

// CompositeKey represents a parsed composite key
type CompositeKey struct {
	ToolName   string `json:"tool_name"`
	DomainName string `json:"domain_name"`
	ID         int    `json:"id"`
}
