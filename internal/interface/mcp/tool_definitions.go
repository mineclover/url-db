package mcp

// ToolDefinition represents an MCP tool definition
type ToolDefinition struct {
	Name        string
	Title       string
	Description string
	InputSchema map[string]interface{}
}

// GetToolDefinitions returns all available MCP tool definitions
func GetToolDefinitions() []ToolDefinition {
	return []ToolDefinition{
		// Server Management
		{
			Name:        "get_server_info",
			Title:       "Server Information",
			Description: "Get server information",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},

		// Domain Management
		{
			Name:        "list_domains",
			Title:       "List Domains",
			Description: "Get all domains",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page": map[string]interface{}{"type": "integer", "default": 1},
					"size": map[string]interface{}{"type": "integer", "default": 20},
				},
			},
		},
		{
			Name:        "create_domain",
			Title:       "Create Domain",
			Description: "Create new domain for organizing URLs",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":        map[string]interface{}{"type": "string", "description": "Domain name"},
					"description": map[string]interface{}{"type": "string", "description": "Domain description"},
				},
				"required": []string{"name", "description"},
			},
		},

		// Node Management
		{
			Name:        "list_nodes",
			Title:       "List URLs",
			Description: "List URLs in domain (requires: domain must exist via create_domain)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "Domain name to list nodes from"},
					"page":        map[string]interface{}{"type": "integer", "default": 1},
					"size":        map[string]interface{}{"type": "integer", "default": 20},
					"search":      map[string]interface{}{"type": "string", "description": "Search query"},
				},
				"required": []string{"domain_name"},
			},
		},
		{
			Name:        "create_node",
			Title:       "Add URL",
			Description: "Add URL to domain (requires: domain must exist via create_domain)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "Domain name"},
					"url":         map[string]interface{}{"type": "string", "description": "URL to store"},
					"title":       map[string]interface{}{"type": "string", "description": "Node title"},
					"description": map[string]interface{}{"type": "string", "description": "Node description"},
				},
				"required": []string{"domain_name", "url"},
			},
		},
		{
			Name:        "get_node",
			Title:       "Get URL Details",
			Description: "Get URL details (requires: node must exist via create_node; returns composite_id from create_node)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "update_node",
			Title:       "Update URL",
			Description: "Update URL title or description (requires: node must exist via create_node; use composite_id from create_node)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
					"title":        map[string]interface{}{"type": "string", "description": "New title"},
					"description":  map[string]interface{}{"type": "string", "description": "New description"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "delete_node",
			Title:       "Remove URL",
			Description: "Remove URL (requires: node must exist via create_node; use composite_id from create_node)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "find_node_by_url",
			Title:       "Find URL",
			Description: "Search by exact URL (requires: domain must exist via create_domain; returns composite_id if found)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "Domain name"},
					"url":         map[string]interface{}{"type": "string", "description": "URL to find"},
				},
				"required": []string{"domain_name", "url"},
			},
		},
		{
			Name:        "scan_all_content",
			Title:       "Scan All Content",
			Description: "Retrieve all URLs and their content from a domain using page-based navigation with token optimization for AI processing",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":         map[string]interface{}{"type": "string", "description": "Domain name to scan"},
					"max_tokens_per_page": map[string]interface{}{"type": "integer", "description": "Maximum tokens per page (recommended: 6000-10000)", "default": 8000},
					"page":                map[string]interface{}{"type": "integer", "description": "Page number (1-based)", "default": 1},
					"include_attributes":  map[string]interface{}{"type": "boolean", "description": "Include node attributes in response", "default": true},
					"compress_attributes": map[string]interface{}{"type": "boolean", "description": "Remove duplicate attribute values for AI context compression", "default": false},
				},
				"required": []string{"domain_name"},
			},
		},

		// Attribute Management
		{
			Name:        "get_node_attributes",
			Title:       "Get URL Attributes",
			Description: "Get URL tags and attributes (requires: node must exist via create_node; attributes defined via create_domain_attribute)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "set_node_attributes",
			Title:       "Set URL Attributes",
			Description: "Add or update URL tags (requires: node must exist via create_node; attributes should be defined via create_domain_attribute unless auto_create_attributes=true)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
					"attributes": map[string]interface{}{
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
					"auto_create_attributes": map[string]interface{}{"type": "boolean", "default": true, "description": "Automatically create attributes if they don't exist"},
				},
				"required": []string{"composite_id", "attributes"},
			},
		},

		// Domain Attribute Schema
		{
			Name:        "list_domain_attributes",
			Title:       "List Domain Attributes",
			Description: "Get available tag types for domain (requires: domain must exist via create_domain)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "The domain to list attributes for"},
				},
				"required": []string{"domain_name"},
			},
		},
		{
			Name:        "create_domain_attribute",
			Title:       "Create Domain Attribute",
			Description: "Define new tag type for domain (requires: domain must exist via create_domain; enables attributes for set_node_attributes)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "The domain to add attribute to"},
					"name":        map[string]interface{}{"type": "string", "description": "Attribute name"},
					"type": map[string]interface{}{
						"type":        "string",
						"description": "One of: tag, ordered_tag, number, string, markdown, image",
						"enum":        []string{"tag", "ordered_tag", "number", "string", "markdown", "image"},
					},
					"description": map[string]interface{}{"type": "string", "description": "Human-readable description"},
				},
				"required": []string{"domain_name", "name", "type"},
			},
		},
		{
			Name:        "get_domain_attribute",
			Title:       "Get Domain Attribute",
			Description: "Get details of a specific domain attribute (requires: attribute must exist via create_domain_attribute)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":    map[string]interface{}{"type": "string", "description": "The domain name"},
					"attribute_name": map[string]interface{}{"type": "string", "description": "The attribute name to get"},
				},
				"required": []string{"domain_name", "attribute_name"},
			},
		},
		{
			Name:        "update_domain_attribute",
			Title:       "Update Domain Attribute",
			Description: "Update domain attribute description (requires: attribute must exist via create_domain_attribute)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":    map[string]interface{}{"type": "string", "description": "The domain name"},
					"attribute_name": map[string]interface{}{"type": "string", "description": "The attribute name to update"},
					"description":    map[string]interface{}{"type": "string", "description": "New description for the attribute"},
				},
				"required": []string{"domain_name", "attribute_name"},
			},
		},
		{
			Name:        "delete_domain_attribute",
			Title:       "Delete Domain Attribute",
			Description: "Remove domain attribute definition (requires: attribute must exist via create_domain_attribute; removes all values from nodes)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":    map[string]interface{}{"type": "string", "description": "The domain name"},
					"attribute_name": map[string]interface{}{"type": "string", "description": "The attribute name to delete"},
				},
				"required": []string{"domain_name", "attribute_name"},
			},
		},

		// Dependency Management
		{
			Name:        "create_dependency",
			Title:       "Create Dependency",
			Description: "Create dependency relationship between nodes (requires: both nodes must exist via create_node; use composite_ids from create_node)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"dependent_node_id":  map[string]interface{}{"type": "string", "description": "Composite ID of the dependent node (format: tool:domain:id)"},
					"dependency_node_id": map[string]interface{}{"type": "string", "description": "Composite ID of the dependency node (format: tool:domain:id)"},
					"dependency_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of dependency",
						"enum":        []string{"hard", "soft", "reference"},
					},
					"cascade_delete": map[string]interface{}{"type": "boolean", "default": false, "description": "Whether to cascade delete"},
					"cascade_update": map[string]interface{}{"type": "boolean", "default": false, "description": "Whether to cascade update"},
					"description":    map[string]interface{}{"type": "string", "description": "Optional description of the dependency"},
				},
				"required": []string{"dependent_node_id", "dependency_node_id", "dependency_type"},
			},
		},
		{
			Name:        "list_node_dependencies",
			Title:       "List Node Dependencies",
			Description: "List what a node depends on (requires: node must exist via create_node; dependencies created via create_dependency)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID of the node (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "list_node_dependents",
			Title:       "List Node Dependents",
			Description: "List what depends on a node (requires: node must exist via create_node; dependencies created via create_dependency)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID of the node (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "delete_dependency",
			Title:       "Delete Dependency",
			Description: "Remove dependency relationship (requires: dependency must exist via create_dependency; use dependency_id from list_node_dependencies)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"dependency_id": map[string]interface{}{"type": "integer", "description": "ID of the dependency relationship to delete"},
				},
				"required": []string{"dependency_id"},
			},
		},

		// Filtering and Queries
		{
			Name:        "filter_nodes_by_attributes",
			Title:       "Filter URLs by Attributes",
			Description: "Filter nodes by attribute values (requires: domain must exist via create_domain; attributes defined via create_domain_attribute)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name": map[string]interface{}{"type": "string", "description": "Domain name to filter nodes from"},
					"filters": map[string]interface{}{
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
					"page": map[string]interface{}{"type": "integer", "default": 1},
					"size": map[string]interface{}{"type": "integer", "default": 20},
				},
				"required": []string{"domain_name", "filters"},
			},
		},
		{
			Name:        "get_node_with_attributes",
			Title:       "Get URL with Attributes",
			Description: "Get URL details with all attributes (requires: node must exist via create_node; combines get_node + get_node_attributes)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:id)"},
				},
				"required": []string{"composite_id"},
			},
		},

		// Template Management
		{
			Name:        "list_templates",
			Title:       "List Templates",
			Description: "List templates in domain (requires: domain must exist via create_domain)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":   map[string]interface{}{"type": "string", "description": "Domain name to list templates from"},
					"page":          map[string]interface{}{"type": "integer", "default": 1},
					"size":          map[string]interface{}{"type": "integer", "default": 20},
					"template_type": map[string]interface{}{"type": "string", "description": "Filter by template type"},
					"active_only":   map[string]interface{}{"type": "boolean", "default": false, "description": "Only return active templates"},
					"search":        map[string]interface{}{"type": "string", "description": "Search query"},
				},
				"required": []string{"domain_name"},
			},
		},
		{
			Name:        "create_template",
			Title:       "Create Template",
			Description: "Create new template in domain (requires: domain must exist via create_domain; use validate_template to check template_data)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"domain_name":   map[string]interface{}{"type": "string", "description": "Domain name"},
					"name":          map[string]interface{}{"type": "string", "description": "Template name"},
					"template_data": map[string]interface{}{"type": "string", "description": "JSON template data"},
					"title":         map[string]interface{}{"type": "string", "description": "Template title"},
					"description":   map[string]interface{}{"type": "string", "description": "Template description"},
				},
				"required": []string{"domain_name", "name", "template_data"},
			},
		},
		{
			Name:        "get_template",
			Title:       "Get Template",
			Description: "Get template details (requires: template must exist via create_template; use composite_id from create_template)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:template:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "update_template",
			Title:       "Update Template",
			Description: "Update template (requires: template must exist via create_template; use validate_template to check new template_data)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id":  map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:template:id)"},
					"template_data": map[string]interface{}{"type": "string", "description": "Updated JSON template data"},
					"title":         map[string]interface{}{"type": "string", "description": "Updated title"},
					"description":   map[string]interface{}{"type": "string", "description": "Updated description"},
					"is_active":     map[string]interface{}{"type": "boolean", "description": "Template active status"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "delete_template",
			Title:       "Delete Template",
			Description: "Delete template (requires: template must exist via create_template)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"composite_id": map[string]interface{}{"type": "string", "description": "Composite ID (format: tool:domain:template:id)"},
				},
				"required": []string{"composite_id"},
			},
		},
		{
			Name:        "clone_template",
			Title:       "Clone Template",
			Description: "Clone existing template (requires: source template must exist via create_template; creates new template with same domain)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"source_composite_id": map[string]interface{}{"type": "string", "description": "Source template composite ID (format: tool:domain:template:id)"},
					"new_name":            map[string]interface{}{"type": "string", "description": "New template name"},
					"new_title":           map[string]interface{}{"type": "string", "description": "New template title"},
					"new_description":     map[string]interface{}{"type": "string", "description": "New template description"},
				},
				"required": []string{"source_composite_id", "new_name"},
			},
		},
		{
			Name:        "generate_template_scaffold",
			Title:       "Generate Template Scaffold",
			Description: "Generate template scaffold for given type (helper: provides starting point for create_template)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"template_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of template to generate",
						"enum":        []string{"layout", "form", "document", "custom"},
					},
				},
				"required": []string{"template_type"},
			},
		},
		{
			Name:        "validate_template",
			Title:       "Validate Template",
			Description: "Validate template data structure (helper: use before create_template or update_template)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"template_data": map[string]interface{}{"type": "string", "description": "JSON template data to validate"},
				},
				"required": []string{"template_data"},
			},
		},
	}
}

// ToMap converts a ToolDefinition to a map for JSON serialization
func (t ToolDefinition) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":        t.Name,
		"title":       t.Title,
		"description": t.Description,
		"inputSchema": t.InputSchema,
	}
}
