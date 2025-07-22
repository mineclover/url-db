package models

// MCPServerInfo represents server information for MCP
type MCPServerInfo struct {
	Name               string   `json:"name"`
	Version            string   `json:"version"`
	Description        string   `json:"description"`
	Capabilities       []string `json:"capabilities"`
	CompositeKeyFormat string   `json:"composite_key_format"`
}

// MCPAttributeRequest represents an attribute in an MCP request
type MCPAttributeRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
