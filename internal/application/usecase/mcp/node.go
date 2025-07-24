package mcp

import (
	"context"
	"fmt"
	"url-db/internal/application/dto/request"
	"url-db/internal/application/dto/response"
	"url-db/internal/application/usecase/node"
)

// NodeUseCases groups node-related use cases for MCP
type NodeUseCases struct {
	Create *node.CreateNodeUseCase
	List   *node.ListNodesUseCase
}

// NewNodeUseCases creates node use cases for MCP
func NewNodeUseCases(createUC *node.CreateNodeUseCase, listUC *node.ListNodesUseCase) *NodeUseCases {
	return &NodeUseCases{
		Create: createUC,
		List:   listUC,
	}
}

// MCPListNodesUseCase wraps the list nodes use case for MCP
type MCPListNodesUseCase struct {
	useCase *node.ListNodesUseCase
}

// NewMCPListNodesUseCase creates a new MCP list nodes use case
func NewMCPListNodesUseCase(useCase *node.ListNodesUseCase) *MCPListNodesUseCase {
	return &MCPListNodesUseCase{useCase: useCase}
}

// Execute executes the list nodes use case with MCP-specific formatting
func (uc *MCPListNodesUseCase) Execute(ctx context.Context, domainName string, page, size int) (interface{}, error) {
	result, err := uc.useCase.Execute(ctx, domainName, page, size)
	if err != nil {
		return nil, err
	}

	// Convert to MCP response format
	content := []map[string]interface{}{}
	for _, node := range result.Nodes {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": formatNodeForMCP(&node),
		})
	}

	if len(content) == 0 {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": "No nodes found",
		})
	}

	return map[string]interface{}{
		"content": content,
		"isError": false,
	}, nil
}

// MCPCreateNodeUseCase wraps the create node use case for MCP
type MCPCreateNodeUseCase struct {
	useCase *node.CreateNodeUseCase
}

// NewMCPCreateNodeUseCase creates a new MCP create node use case
func NewMCPCreateNodeUseCase(useCase *node.CreateNodeUseCase) *MCPCreateNodeUseCase {
	return &MCPCreateNodeUseCase{useCase: useCase}
}

// Execute executes the create node use case with MCP-specific formatting
func (uc *MCPCreateNodeUseCase) Execute(ctx context.Context, domainName, url, title, description string) (interface{}, error) {
	req := &request.CreateNodeRequest{
		DomainName:  domainName,
		URL:         url,
		Title:       title,
		Description: description,
	}

	result, err := uc.useCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": formatCreateNodeResultForMCP(result),
			},
		},
		"isError": false,
	}, nil
}

// Helper functions for MCP formatting
func formatNodeForMCP(node *response.NodeResponse) string {
	return fmt.Sprintf("Node ID: %d\nURL: %s\nTitle: %s\nDescription: %s\nCreated: %s",
		node.ID, node.URL, node.Title, node.Description, node.CreatedAt.Format("2006-01-02 15:04:05"))
}

func formatCreateNodeResultForMCP(result *response.NodeResponse) string {
	return fmt.Sprintf("Successfully created node in domain '%s'\nURL: %s\nTitle: %s\nDescription: %s\nCreated: %s",
		result.DomainName, result.URL, result.Title, result.Description, result.CreatedAt.Format("2006-01-02 15:04:05"))
}
