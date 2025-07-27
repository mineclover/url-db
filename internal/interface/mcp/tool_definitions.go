package mcp

// Helper functions for creating pointers
func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

// ToolDefinition represents an MCP tool definition according to TypeScript schema 2025-06-18
type ToolDefinition struct {
	Name         string                 `json:"name"`
	Description  *string                `json:"description,omitempty"`
	InputSchema  InputSchema            `json:"inputSchema"`
	OutputSchema *OutputSchema          `json:"outputSchema,omitempty"`
	Annotations  *ToolAnnotations       `json:"annotations,omitempty"`
	Meta         map[string]interface{} `json:"_meta,omitempty"`
}

// InputSchema represents the input schema for tools (must be object type)
type InputSchema struct {
	Type       string                            `json:"type"`
	Properties map[string]map[string]interface{} `json:"properties,omitempty"`
	Required   []string                          `json:"required,omitempty"`
}

// OutputSchema represents the output schema for tools (must be object type)  
type OutputSchema struct {
	Type       string                            `json:"type"`
	Properties map[string]map[string]interface{} `json:"properties,omitempty"`
	Required   []string                          `json:"required,omitempty"`
}

// ToolAnnotations represents optional tool annotations as per MCP 2025-06-18 schema
type ToolAnnotations struct {
	Title           *string `json:"title,omitempty"`
	ReadOnlyHint    *bool   `json:"readOnlyHint,omitempty"`
	DestructiveHint *bool   `json:"destructiveHint,omitempty"`
	IdempotentHint  *bool   `json:"idempotentHint,omitempty"`
	OpenWorldHint   *bool   `json:"openWorldHint,omitempty"`
}

// GetToolDefinitions returns all available MCP tool definitions
func GetToolDefinitions() []ToolDefinition {
	return []ToolDefinition{
		// Server Management
		{
			Name:        "get_server_info",
			Description: stringPtr("Get server information"),
			InputSchema: InputSchema{
				Type:       "object",
				Properties: map[string]map[string]interface{}{},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		// Domain Management
		{
			Name:        "list_domains",
			Description: stringPtr("Get all domains"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"page": {"type": "integer", "default": 1},
					"size": {"type": "integer", "default": 20},
				},
			},
			OutputSchema: &OutputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domains": {
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name":        map[string]interface{}{"type": "string"},
								"description": map[string]interface{}{"type": "string"},
								"created_at":  map[string]interface{}{"type": "string", "format": "date-time"},
							},
						},
					},
					"total_count": {"type": "integer"},
					"page":        {"type": "integer"},
					"total_pages": {"type": "integer"},
				},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "create_domain",
			Description: stringPtr("Create new domain for organizing URLs"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"name":        {"type": "string", "description": "Domain name"},
					"description": {"type": "string", "description": "Domain description"},
				},
				Required: []string{"name", "description"},
			},
			OutputSchema: &OutputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"name":        {"type": "string"},
					"description": {"type": "string"},
					"created_at":  {"type": "string", "format": "date-time"},
				},
				Required: []string{"name", "description", "created_at"},
			},
		},

		// Node Management
		{
			Name:        "list_nodes",
			Description: stringPtr("List URLs in domain (requires: domain must exist via create_domain)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name": {"type": "string", "description": "Domain name to list nodes from"},
					"page":        {"type": "integer", "default": 1},
					"size":        {"type": "integer", "default": 20},
					"search":      {"type": "string", "description": "Search query"},
				},
				Required: []string{"domain_name"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "create_node",
			Description: stringPtr("Add URL to domain (requires: domain must exist via create_domain)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name": {"type": "string", "description": "Domain name"},
					"url":         {"type": "string", "description": "URL to store"},
					"title":       {"type": "string", "description": "Node title"},
					"description": {"type": "string", "description": "Node description"},
				},
				Required: []string{"domain_name", "url"},
			},
		},

		{
			Name:        "get_node",
			Description: stringPtr("Get URL details (requires: node must exist via create_node; returns composite_id from create_node)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				Required: []string{"composite_id"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "update_node",
			Description: stringPtr("Update URL title or description (requires: node must exist via create_node; use composite_id from create_node)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID (format: tool:domain:id)"},
					"title":        {"type": "string", "description": "New title"},
					"description":  {"type": "string", "description": "New description"},
				},
				Required: []string{"composite_id"},
			},
		},

		{
			Name:        "delete_node",
			Description: stringPtr("Remove URL (requires: node must exist via create_node; use composite_id from create_node)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				Required: []string{"composite_id"},
			},
			OutputSchema: &OutputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"deleted":      {"type": "boolean"},
					"composite_id": {"type": "string"},
					"message":      {"type": "string"},
				},
				Required: []string{"deleted", "composite_id"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(true),
				OpenWorldHint:   boolPtr(false),
			},
		},

		{
			Name:        "find_node_by_url",
			Description: stringPtr("Search by exact URL (requires: domain must exist via create_domain; returns composite_id if found)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name": {"type": "string", "description": "Domain name"},
					"url":         {"type": "string", "description": "URL to find"},
				},
				Required: []string{"domain_name", "url"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "scan_all_content",
			Description: stringPtr("Retrieve all URLs and their content from a domain using page-based navigation with token optimization for AI processing"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name":         {"type": "string", "description": "Domain name to scan"},
					"max_tokens_per_page": {"type": "integer", "description": "Maximum tokens per page (recommended: 6000-10000)", "default": 8000},
					"page":                {"type": "integer", "description": "Page number (1-based)", "default": 1},
					"include_attributes":  {"type": "boolean", "description": "Include node attributes in response", "default": true},
					"compress_attributes": {"type": "boolean", "description": "Remove duplicate attribute values for AI context compression", "default": false},
				},
				Required: []string{"domain_name"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		// Attribute Management
		{
			Name:        "get_node_attributes",
			Description: stringPtr("Get URL tags and attributes (requires: node must exist via create_node; attributes defined via create_domain_attribute)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				Required: []string{"composite_id"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "set_node_attributes",
			Description: stringPtr("Add or update URL tags (requires: node must exist via create_node; attributes should be defined via create_domain_attribute unless auto_create_attributes=true)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID (format: tool:domain:id)"},
					"attributes": {
						"type":        "array",
						"description": "Array of attributes to set",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name":        map[string]interface{}{"type": "string", "description": "Attribute name"},
								"value":       map[string]interface{}{"type": "string", "description": "Attribute value"},
								"order_index": map[string]interface{}{"type": "integer", "description": "Order index (for ordered_tag type)"},
							},
							"required": []string{"name", "value"},
						},
					},
					"auto_create_attributes": {"type": "boolean", "default": true, "description": "Automatically create attributes if they don't exist"},
				},
				Required: []string{"composite_id", "attributes"},
			},
		},

		// Domain Attribute Schema
		{
			Name:        "list_domain_attributes",
			Description: stringPtr("Get available tag types for domain (requires: domain must exist via create_domain)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name": {"type": "string", "description": "The domain to list attributes for"},
				},
				Required: []string{"domain_name"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "create_domain_attribute",
			Description: stringPtr("Define new tag type for domain (requires: domain must exist via create_domain; enables attributes for set_node_attributes)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name": {"type": "string", "description": "The domain to add attribute to"},
					"name":        {"type": "string", "description": "Attribute name"},
					"type": {
						"type":        "string",
						"description": "One of: tag, ordered_tag, number, string, markdown, image",
						"enum":        []string{"tag", "ordered_tag", "number", "string", "markdown", "image"},
					},
					"description": {"type": "string", "description": "Human-readable description"},
				},
				Required: []string{"domain_name", "name", "type"},
			},
		},

		{
			Name:        "get_domain_attribute",
			Description: stringPtr("Get details of a specific domain attribute (requires: attribute must exist via create_domain_attribute)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name":    {"type": "string", "description": "The domain name"},
					"attribute_name": {"type": "string", "description": "The attribute name to get"},
				},
				Required: []string{"domain_name", "attribute_name"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "update_domain_attribute",
			Description: stringPtr("Update domain attribute description (requires: attribute must exist via create_domain_attribute)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name":    {"type": "string", "description": "The domain name"},
					"attribute_name": {"type": "string", "description": "The attribute name to update"},
					"description":    {"type": "string", "description": "New description for the attribute"},
				},
				Required: []string{"domain_name", "attribute_name"},
			},
		},

		{
			Name:        "delete_domain_attribute",
			Description: stringPtr("Remove domain attribute definition (requires: attribute must exist via create_domain_attribute; removes all values from nodes)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name":    {"type": "string", "description": "The domain name"},
					"attribute_name": {"type": "string", "description": "The attribute name to delete"},
				},
				Required: []string{"domain_name", "attribute_name"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(true),
				OpenWorldHint:   boolPtr(false),
			},
		},

		// Dependency Management
		{
			Name:        "create_dependency",
			Description: stringPtr("Create dependency relationship between nodes (requires: both nodes must exist via create_node; use composite_ids from create_node)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"dependent_node_id":  {"type": "string", "description": "Composite ID of the dependent node (format: tool:domain:id)"},
					"dependency_node_id": {"type": "string", "description": "Composite ID of the dependency node (format: tool:domain:id)"},
					"dependency_type": {
						"type":        "string",
						"description": "Type of dependency",
						"enum":        []string{"hard", "soft", "reference"},
					},
					"cascade_delete": {"type": "boolean", "default": false, "description": "Whether to cascade delete"},
					"cascade_update": {"type": "boolean", "default": false, "description": "Whether to cascade update"},
					"description":    {"type": "string", "description": "Optional description of the dependency"},
				},
				Required: []string{"dependent_node_id", "dependency_node_id", "dependency_type"},
			},
		},

		{
			Name:        "list_node_dependencies",
			Description: stringPtr("List what a node depends on (requires: node must exist via create_node; dependencies created via create_dependency)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID of the node (format: tool:domain:id)"},
				},
				Required: []string{"composite_id"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "list_node_dependents",
			Description: stringPtr("List what depends on a node (requires: node must exist via create_node; dependencies created via create_dependency)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID of the node (format: tool:domain:id)"},
				},
				Required: []string{"composite_id"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "delete_dependency",
			Description: stringPtr("Remove dependency relationship (requires: dependency must exist via create_dependency; use dependency_id from list_node_dependencies)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"dependency_id": {"type": "integer", "description": "ID of the dependency relationship to delete"},
				},
				Required: []string{"dependency_id"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(true),
				OpenWorldHint:   boolPtr(false),
			},
		},

		// Filtering and Queries
		{
			Name:        "filter_nodes_by_attributes",
			Description: stringPtr("Filter nodes by attribute values (requires: domain must exist via create_domain; attributes defined via create_domain_attribute)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name": {"type": "string", "description": "Domain name to filter nodes from"},
					"filters": {
						"type":        "array",
						"description": "Array of attribute filters",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name":     map[string]interface{}{"type": "string", "description": "Attribute name"},
								"value":    map[string]interface{}{"type": "string", "description": "Attribute value"},
								"operator": map[string]interface{}{"type": "string", "description": "Comparison operator", "enum": []string{"equals", "contains", "starts_with", "ends_with"}, "default": "equals"},
							},
							"required": []string{"name", "value"},
						},
					},
					"page": {"type": "integer", "default": 1},
					"size": {"type": "integer", "default": 20},
				},
				Required: []string{"domain_name", "filters"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "get_node_with_attributes",
			Description: stringPtr("Get URL details with all attributes (requires: node must exist via create_node; combines get_node + get_node_attributes)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				Required: []string{"composite_id"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		// Template Management
		{
			Name:        "list_templates",
			Description: stringPtr("List templates in domain (requires: domain must exist via create_domain)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name":   {"type": "string", "description": "Domain name to list templates from"},
					"page":          {"type": "integer", "default": 1},
					"size":          {"type": "integer", "default": 20},
					"template_type": {"type": "string", "description": "Filter by template type"},
					"active_only":   {"type": "boolean", "default": false, "description": "Only return active templates"},
					"search":        {"type": "string", "description": "Search query"},
				},
				Required: []string{"domain_name"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "create_template",
			Description: stringPtr("Create new template in domain (requires: domain must exist via create_domain; use validate_template to check template_data)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"domain_name":   {"type": "string", "description": "Domain name"},
					"name":          {"type": "string", "description": "Template name"},
					"template_data": {"type": "string", "description": "JSON template data"},
					"title":         {"type": "string", "description": "Template title"},
					"description":   {"type": "string", "description": "Template description"},
				},
				Required: []string{"domain_name", "name", "template_data"},
			},
		},

		{
			Name:        "get_template",
			Description: stringPtr("Get template details (requires: template must exist via create_template; use composite_id from create_template)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID (format: tool:domain:template:id)"},
				},
				Required: []string{"composite_id"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "update_template",
			Description: stringPtr("Update template (requires: template must exist via create_template; use validate_template to check new template_data)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id":  {"type": "string", "description": "Composite ID (format: tool:domain:template:id)"},
					"template_data": {"type": "string", "description": "Updated JSON template data"},
					"title":         {"type": "string", "description": "Updated title"},
					"description":   {"type": "string", "description": "Updated description"},
					"is_active":     {"type": "boolean", "description": "Template active status"},
				},
				Required: []string{"composite_id"},
			},
		},

		{
			Name:        "delete_template",
			Description: stringPtr("Delete template (requires: template must exist via create_template)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"composite_id": {"type": "string", "description": "Composite ID (format: tool:domain:template:id)"},
				},
				Required: []string{"composite_id"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(true),
				OpenWorldHint:   boolPtr(false),
			},
		},

		{
			Name:        "clone_template",
			Description: stringPtr("Clone existing template (requires: source template must exist via create_template; creates new template with same domain)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"source_composite_id": {"type": "string", "description": "Source template composite ID (format: tool:domain:template:id)"},
					"new_name":            {"type": "string", "description": "New template name"},
					"new_title":           {"type": "string", "description": "New template title"},
					"new_description":     {"type": "string", "description": "New template description"},
				},
				Required: []string{"source_composite_id", "new_name"},
			},
		},

		{
			Name:        "generate_template_scaffold",
			Description: stringPtr("Generate template scaffold for given type (helper: provides starting point for create_template)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"template_type": {
						"type":        "string",
						"description": "Type of template to generate",
						"enum":        []string{"layout", "form", "document", "custom"},
					},
				},
				Required: []string{"template_type"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},

		{
			Name:        "validate_template",
			Description: stringPtr("Validate template data structure (helper: use before create_template or update_template)"),
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"template_data": {"type": "string", "description": "JSON template data to validate"},
				},
				Required: []string{"template_data"},
			},
			Annotations: &ToolAnnotations{
				ReadOnlyHint:  boolPtr(true),
				OpenWorldHint: boolPtr(false),
			},
		},
	}
}

// ToMap converts a ToolDefinition to a map for JSON serialization according to TypeScript schema 2025-06-18
func (t ToolDefinition) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"name":        t.Name,
		"inputSchema": t.InputSchema,
	}
	
	// Add optional fields only if present
	if t.Description != nil {
		result["description"] = *t.Description
	}
	
	if t.OutputSchema != nil {
		result["outputSchema"] = *t.OutputSchema
	}
	
	// Add annotations if present
	if t.Annotations != nil {
		annotations := make(map[string]interface{})
		
		if t.Annotations.Title != nil {
			annotations["title"] = *t.Annotations.Title
		}
		if t.Annotations.ReadOnlyHint != nil {
			annotations["readOnlyHint"] = *t.Annotations.ReadOnlyHint
		}
		if t.Annotations.DestructiveHint != nil {
			annotations["destructiveHint"] = *t.Annotations.DestructiveHint
		}
		if t.Annotations.IdempotentHint != nil {
			annotations["idempotentHint"] = *t.Annotations.IdempotentHint
		}
		if t.Annotations.OpenWorldHint != nil {
			annotations["openWorldHint"] = *t.Annotations.OpenWorldHint
		}
		
		if len(annotations) > 0 {
			result["annotations"] = annotations
		}
	}
	
	// Add meta if present
	if t.Meta != nil && len(t.Meta) > 0 {
		result["_meta"] = t.Meta
	}
	
	return result
}