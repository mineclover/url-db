package mcp

import (
	"context"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"

	"url-db/internal/application/dto/request"
	"url-db/internal/application/usecase/node"
	"url-db/internal/constants"
	"url-db/internal/interface/setup"
	"url-db/internal/compositekey"
)

// MCPServerV2 implements MCP server using mcp-golang library
type MCPServerV2 struct {
	server  *mcp_golang.Server
	factory *setup.ApplicationFactory
}

// NewMCPServerV2 creates a new MCP server using mcp-golang
func NewMCPServerV2(factory *setup.ApplicationFactory) *MCPServerV2 {
	transport := stdio.NewStdioServerTransport()
	server := mcp_golang.NewServer(transport)

	mcpServer := &MCPServerV2{
		server:  server,
		factory: factory,
	}

	mcpServer.registerTools()
	return mcpServer
}

// Start starts the MCP server
func (s *MCPServerV2) Start(ctx context.Context) error {
	return s.server.Serve()
}

// registerTools registers all MCP tools
func (s *MCPServerV2) registerTools() {
	// Server info
	s.server.RegisterTool("get_server_info", "Get server information", s.handleGetServerInfo)

	// Domain management
	s.server.RegisterTool("list_domains", "Get all domains", s.handleListDomains)
	s.server.RegisterTool("create_domain", "Create new domain for organizing URLs", s.handleCreateDomain)

	// Node management
	s.server.RegisterTool("list_nodes", "List URLs in domain", s.handleListNodes)
	s.server.RegisterTool("create_node", "Add URL to domain", s.handleCreateNode)
	s.server.RegisterTool("get_node", "Get URL details", s.handleGetNode)
	s.server.RegisterTool("update_node", "Update URL title or description", s.handleUpdateNode)
	s.server.RegisterTool("delete_node", "Remove URL", s.handleDeleteNode)
	s.server.RegisterTool("find_node_by_url", "Search by exact URL", s.handleFindNodeByURL)

	// Node attributes
	s.server.RegisterTool("get_node_attributes", "Get URL tags and attributes", s.handleGetNodeAttributes)
	s.server.RegisterTool("set_node_attributes", "Add or update URL tags", s.handleSetNodeAttributes)

	// Domain attributes (schema)
	s.server.RegisterTool("list_domain_attributes", "Get available tag types for domain", s.handleListDomainAttributes)
	s.server.RegisterTool("create_domain_attribute", "Define new tag type for domain", s.handleCreateDomainAttribute)
	s.server.RegisterTool("get_domain_attribute", "Get details of a specific domain attribute", s.handleGetDomainAttribute)
	s.server.RegisterTool("update_domain_attribute", "Update domain attribute description", s.handleUpdateDomainAttribute)
	s.server.RegisterTool("delete_domain_attribute", "Remove domain attribute definition", s.handleDeleteDomainAttribute)

	// Dependencies
	s.server.RegisterTool("create_dependency", "Create dependency relationship between nodes", s.handleCreateDependency)
	s.server.RegisterTool("list_node_dependencies", "List what a node depends on", s.handleListNodeDependencies)
	s.server.RegisterTool("list_node_dependents", "List what depends on a node", s.handleListNodeDependents)
	s.server.RegisterTool("delete_dependency", "Remove dependency relationship", s.handleDeleteDependency)
}

// Tool handlers

func (s *MCPServerV2) handleGetServerInfo(args GetServerInfoArgs) (*mcp_golang.ToolResponse, error) {
	info := fmt.Sprintf("URL Database MCP Server\nVersion: %s\nDescription: %s",
		constants.DefaultServerVersion,
		"URL database management system with Clean Architecture and MCP integration")

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(info)), nil
}

func (s *MCPServerV2) handleListDomains(args ListDomainsArgs) (*mcp_golang.ToolResponse, error) {
	ctx := context.Background()
	useCase := s.factory.CreateListDomainsUseCase()

	if args.Page == 0 {
		args.Page = 1
	}
	if args.Size == 0 {
		args.Size = 20
	}

	response, err := useCase.Execute(ctx, args.Page, args.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	content := "Available domains:\n"
	for _, domain := range response.Domains {
		content += fmt.Sprintf("• %s: %s (created: %s)\n",
			domain.Name, domain.Description, domain.CreatedAt.Format("2006-01-02"))
	}

	if len(response.Domains) == 0 {
		content = "No domains found. Create a domain first using create_domain."
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(content)), nil
}

func (s *MCPServerV2) handleCreateDomain(args CreateDomainArgs) (*mcp_golang.ToolResponse, error) {
	ctx := context.Background()
	useCase := s.factory.CreateCreateDomainUseCase()

	request := &request.CreateDomainRequest{
		Name:        args.Name,
		Description: args.Description,
	}

	response, err := useCase.Execute(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	content := fmt.Sprintf("Successfully created domain: %s\nDescription: %s\nCreated: %s",
		response.Name, response.Description, response.CreatedAt.Format("2006-01-02 15:04:05"))

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(content)), nil
}

func (s *MCPServerV2) handleListNodes(args ListNodesArgs) (*mcp_golang.ToolResponse, error) {
	ctx := context.Background()
	useCase := s.factory.CreateListNodesUseCase()

	if args.Page == 0 {
		args.Page = 1
	}
	if args.Size == 0 {
		args.Size = 20
	}

	response, err := useCase.Execute(ctx, args.DomainName, args.Page, args.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	content := fmt.Sprintf("URLs in domain '%s':\n", args.DomainName)
	for _, node := range response.Nodes {
		compositeID := fmt.Sprintf("%s:%s:%d", constants.DefaultServerName, args.DomainName, node.ID)
		content += fmt.Sprintf("• [%s] %s\n  URL: %s\n  Description: %s\n",
			compositeID, node.Title, node.URL, node.Description)
	}

	if len(response.Nodes) == 0 {
		content = fmt.Sprintf("No URLs found in domain '%s'.", args.DomainName)
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(content)), nil
}

func (s *MCPServerV2) handleCreateNode(args CreateNodeArgs) (*mcp_golang.ToolResponse, error) {
	ctx := context.Background()
	useCase := s.factory.CreateCreateNodeUseCase()

	request := &request.CreateNodeRequest{
		DomainName:  args.DomainName,
		URL:         args.URL,
		Title:       args.Title,
		Description: args.Description,
	}

	response, err := useCase.Execute(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	content := fmt.Sprintf("Successfully created node in domain '%s'\nURL: %s\nTitle: %s\nDescription: %s\nCreated: %s",
		args.DomainName, response.URL, response.Title, response.Description, response.CreatedAt.Format("2006-01-02 15:04:05"))

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(content)), nil
}

func (s *MCPServerV2) handleGetNode(args GetNodeArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Get node feature coming soon")), nil
}

func (s *MCPServerV2) handleUpdateNode(args UpdateNodeArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Update node feature coming soon")), nil
}

func (s *MCPServerV2) handleDeleteNode(args DeleteNodeArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Delete node feature coming soon")), nil
}

func (s *MCPServerV2) handleFindNodeByURL(args FindNodeByURLArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Find node by URL feature coming soon")), nil
}

// Node attributes handlers
func (s *MCPServerV2) handleGetNodeAttributes(args GetNodeAttributesArgs) (*mcp_golang.ToolResponse, error) {
	// TODO: Implement get node attributes - need to create GetNodeAttributesUseCase
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Get node attributes feature coming soon")), nil
}

func (s *MCPServerV2) handleSetNodeAttributes(args SetNodeAttributesArgs) (*mcp_golang.ToolResponse, error) {
	ctx := context.Background()
	useCase := s.factory.CreateSetNodeAttributesUseCase()

	// Parse composite ID to get node ID
	compositeKey, err := compositekey.Parse(args.CompositeID)
	if err != nil {
		return nil, fmt.Errorf("invalid composite ID format: %w", err)
	}

	// Convert MCP attributes to UseCase attributes
	attributes := make([]node.AttributeInput, len(args.Attributes))
	for i, attr := range args.Attributes {
		attributes[i] = node.AttributeInput{
			Name:       attr.Name,
			Value:      attr.Value,
			OrderIndex: attr.OrderIndex,
		}
	}

	err = useCase.Execute(ctx, compositeKey.ID, attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to set node attributes: %w", err)
	}

	content := fmt.Sprintf("Successfully set %d attributes for node: %s", len(args.Attributes), args.CompositeID)
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(content)), nil
}

func (s *MCPServerV2) handleListDomainAttributes(args ListDomainAttributesArgs) (*mcp_golang.ToolResponse, error) {
	ctx := context.Background()
	
	// Get domain by name to get ID
	domain, err := s.factory.GetDomainByName(ctx, args.DomainName)
	if err != nil {
		return nil, fmt.Errorf("failed to find domain '%s': %w", args.DomainName, err)
	}

	useCase := s.factory.CreateListAttributesUseCase()
	response, err := useCase.Execute(ctx, domain.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to list domain attributes: %w", err)
	}

	content := fmt.Sprintf("Attributes for domain '%s':\n", args.DomainName)
	for _, attr := range response.Attributes {
		content += fmt.Sprintf("• %s (%s): %s\n", attr.Name, attr.Type, attr.Description)
	}

	if len(response.Attributes) == 0 {
		content = fmt.Sprintf("No attributes defined for domain '%s'.", args.DomainName)
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(content)), nil
}

func (s *MCPServerV2) handleCreateDomainAttribute(args CreateDomainAttributeArgs) (*mcp_golang.ToolResponse, error) {
	ctx := context.Background()
	
	// Get domain by name to get ID
	domain, err := s.factory.GetDomainByName(ctx, args.DomainName)
	if err != nil {
		return nil, fmt.Errorf("failed to find domain '%s': %w", args.DomainName, err)
	}

	useCase := s.factory.CreateCreateAttributeUseCase()
	request := &request.CreateAttributeRequest{
		DomainID:    domain.ID(),
		Name:        args.Name,
		Type:        args.Type,
		Description: args.Description,
	}

	response, err := useCase.Execute(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain attribute: %w", err)
	}

	content := fmt.Sprintf("Successfully created domain attribute:\nDomain: %s\nName: %s\nType: %s\nDescription: %s\nCreated: %s",
		args.DomainName, response.Name, response.Type, response.Description, response.CreatedAt.Format("2006-01-02 15:04:05"))

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(content)), nil
}

func (s *MCPServerV2) handleGetDomainAttribute(args GetDomainAttributeArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Get domain attribute feature coming soon")), nil
}

func (s *MCPServerV2) handleUpdateDomainAttribute(args UpdateDomainAttributeArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Update domain attribute feature coming soon")), nil
}

func (s *MCPServerV2) handleDeleteDomainAttribute(args DeleteDomainAttributeArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Delete domain attribute feature coming soon")), nil
}

func (s *MCPServerV2) handleCreateDependency(args CreateDependencyArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Create dependency feature coming soon")), nil
}

func (s *MCPServerV2) handleListNodeDependencies(args ListNodeDependenciesArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("List node dependencies feature coming soon")), nil
}

func (s *MCPServerV2) handleListNodeDependents(args ListNodeDependentsArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("List node dependents feature coming soon")), nil
}

func (s *MCPServerV2) handleDeleteDependency(args DeleteDependencyArgs) (*mcp_golang.ToolResponse, error) {
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Delete dependency feature coming soon")), nil
}