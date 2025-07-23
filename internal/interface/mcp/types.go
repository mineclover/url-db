package mcp

// Tool argument types with JSON schema annotations for mcp-golang

// Domain management
type CreateDomainArgs struct {
	Name        string `json:"name" jsonschema:"required,title=Domain Name,description=Domain name"`
	Description string `json:"description" jsonschema:"required,title=Description,description=Domain description"`
}

type ListDomainsArgs struct {
	Page int `json:"page,omitempty" jsonschema:"title=Page,description=Page number,default=1"`
	Size int `json:"size,omitempty" jsonschema:"title=Size,description=Page size,default=20"`
}

// Node management
type CreateNodeArgs struct {
	DomainName  string `json:"domain_name" jsonschema:"required,title=Domain Name,description=Domain name"`
	URL         string `json:"url" jsonschema:"required,title=URL,description=URL to store"`
	Title       string `json:"title,omitempty" jsonschema:"title=Title,description=Node title"`
	Description string `json:"description,omitempty" jsonschema:"title=Description,description=Node description"`
}

type ListNodesArgs struct {
	DomainName string `json:"domain_name" jsonschema:"required,title=Domain Name,description=Domain name to list nodes from"`
	Page       int    `json:"page,omitempty" jsonschema:"title=Page,description=Page number,default=1"`
	Size       int    `json:"size,omitempty" jsonschema:"title=Size,description=Page size,default=20"`
	Search     string `json:"search,omitempty" jsonschema:"title=Search,description=Search query"`
}

type GetNodeArgs struct {
	CompositeID string `json:"composite_id" jsonschema:"required,title=Composite ID,description=Composite ID (format: tool:domain:id)"`
}

type UpdateNodeArgs struct {
	CompositeID string `json:"composite_id" jsonschema:"required,title=Composite ID,description=Composite ID (format: tool:domain:id)"`
	Title       string `json:"title,omitempty" jsonschema:"title=Title,description=New title"`
	Description string `json:"description,omitempty" jsonschema:"title=Description,description=New description"`
}

type DeleteNodeArgs struct {
	CompositeID string `json:"composite_id" jsonschema:"required,title=Composite ID,description=Composite ID (format: tool:domain:id)"`
}

type FindNodeByURLArgs struct {
	DomainName string `json:"domain_name" jsonschema:"required,title=Domain Name,description=Domain name"`
	URL        string `json:"url" jsonschema:"required,title=URL,description=URL to find"`
}

// Node attributes
type AttributeValue struct {
	Name       string `json:"name" jsonschema:"required,title=Name,description=Attribute name"`
	Value      string `json:"value" jsonschema:"required,title=Value,description=Attribute value"`
	OrderIndex *int   `json:"order_index,omitempty" jsonschema:"title=Order Index,description=Order index (for ordered_tag type)"`
}

type SetNodeAttributesArgs struct {
	CompositeID           string           `json:"composite_id" jsonschema:"required,title=Composite ID,description=Composite ID (format: tool:domain:id)"`
	Attributes            []AttributeValue `json:"attributes" jsonschema:"required,title=Attributes,description=Array of attributes to set"`
	AutoCreateAttributes  bool             `json:"auto_create_attributes,omitempty" jsonschema:"title=Auto Create,description=Automatically create attributes if they don't exist,default=true"`
}

type GetNodeAttributesArgs struct {
	CompositeID string `json:"composite_id" jsonschema:"required,title=Composite ID,description=Composite ID (format: tool:domain:id)"`
}

// Domain attributes (schema)
type ListDomainAttributesArgs struct {
	DomainName string `json:"domain_name" jsonschema:"required,title=Domain Name,description=The domain to list attributes for"`
}

type CreateDomainAttributeArgs struct {
	DomainName  string `json:"domain_name" jsonschema:"required,title=Domain Name,description=The domain to add attribute to"`
	Name        string `json:"name" jsonschema:"required,title=Name,description=Attribute name"`
	Type        string `json:"type" jsonschema:"required,title=Type,description=One of: tag, ordered_tag, number, string, markdown, image,enum=tag,enum=ordered_tag,enum=number,enum=string,enum=markdown,enum=image"`
	Description string `json:"description,omitempty" jsonschema:"title=Description,description=Human-readable description"`
}

type GetDomainAttributeArgs struct {
	DomainName    string `json:"domain_name" jsonschema:"required,title=Domain Name,description=The domain name"`
	AttributeName string `json:"attribute_name" jsonschema:"required,title=Attribute Name,description=The attribute name to get"`
}

type UpdateDomainAttributeArgs struct {
	DomainName    string `json:"domain_name" jsonschema:"required,title=Domain Name,description=The domain name"`
	AttributeName string `json:"attribute_name" jsonschema:"required,title=Attribute Name,description=The attribute name to update"`
	Description   string `json:"description" jsonschema:"required,title=Description,description=New description for the attribute"`
}

type DeleteDomainAttributeArgs struct {
	DomainName    string `json:"domain_name" jsonschema:"required,title=Domain Name,description=The domain name"`
	AttributeName string `json:"attribute_name" jsonschema:"required,title=Attribute Name,description=The attribute name to delete"`
}

// Dependencies
type CreateDependencyArgs struct {
	DependentNodeID   string `json:"dependent_node_id" jsonschema:"required,title=Dependent Node ID,description=Composite ID of the dependent node (format: tool:domain:id)"`
	DependencyNodeID  string `json:"dependency_node_id" jsonschema:"required,title=Dependency Node ID,description=Composite ID of the dependency node (format: tool:domain:id)"`
	DependencyType    string `json:"dependency_type" jsonschema:"required,title=Dependency Type,description=Type of dependency,enum=hard,enum=soft,enum=reference"`
	Description       string `json:"description,omitempty" jsonschema:"title=Description,description=Optional description of the dependency"`
	CascadeDelete     bool   `json:"cascade_delete,omitempty" jsonschema:"title=Cascade Delete,description=Whether to cascade delete,default=false"`
	CascadeUpdate     bool   `json:"cascade_update,omitempty" jsonschema:"title=Cascade Update,description=Whether to cascade update,default=false"`
}

type ListNodeDependenciesArgs struct {
	CompositeID string `json:"composite_id" jsonschema:"required,title=Composite ID,description=Composite ID of the node (format: tool:domain:id)"`
}

type ListNodeDependentsArgs struct {
	CompositeID string `json:"composite_id" jsonschema:"required,title=Composite ID,description=Composite ID of the node (format: tool:domain:id)"`
}

type DeleteDependencyArgs struct {
	DependencyID int `json:"dependency_id" jsonschema:"required,title=Dependency ID,description=ID of the dependency relationship to delete"`
}

// Server info
type GetServerInfoArgs struct {
	// No arguments needed
}