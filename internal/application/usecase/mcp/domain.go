package mcp

import (
	"context"
	"fmt"
	"url-db/internal/application/dto/request"
	"url-db/internal/application/dto/response"
	"url-db/internal/application/usecase/domain"
)

// DomainUseCases groups domain-related use cases for MCP
type DomainUseCases struct {
	Create *domain.CreateDomainUseCase
	List   *domain.ListDomainsUseCase
}

// NewDomainUseCases creates domain use cases for MCP
func NewDomainUseCases(createUC *domain.CreateDomainUseCase, listUC *domain.ListDomainsUseCase) *DomainUseCases {
	return &DomainUseCases{
		Create: createUC,
		List:   listUC,
	}
}

// MCPListDomainsUseCase wraps the list domains use case for MCP
type MCPListDomainsUseCase struct {
	useCase *domain.ListDomainsUseCase
}

// NewMCPListDomainsUseCase creates a new MCP list domains use case
func NewMCPListDomainsUseCase(useCase *domain.ListDomainsUseCase) *MCPListDomainsUseCase {
	return &MCPListDomainsUseCase{useCase: useCase}
}

// Execute executes the list domains use case with MCP-specific formatting
func (uc *MCPListDomainsUseCase) Execute(ctx context.Context, page, size int) (interface{}, error) {
	result, err := uc.useCase.Execute(ctx, page, size)
	if err != nil {
		return nil, err
	}

	// Convert to MCP response format
	content := []map[string]interface{}{}
	for _, domain := range result.Domains {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": formatDomainForMCP(&domain),
		})
	}

	if len(content) == 0 {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": "No domains found",
		})
	}

	return map[string]interface{}{
		"content": content,
		"isError": false,
	}, nil
}

// MCPCreateDomainUseCase wraps the create domain use case for MCP
type MCPCreateDomainUseCase struct {
	useCase *domain.CreateDomainUseCase
}

// NewMCPCreateDomainUseCase creates a new MCP create domain use case
func NewMCPCreateDomainUseCase(useCase *domain.CreateDomainUseCase) *MCPCreateDomainUseCase {
	return &MCPCreateDomainUseCase{useCase: useCase}
}

// Execute executes the create domain use case with MCP-specific formatting
func (uc *MCPCreateDomainUseCase) Execute(ctx context.Context, name, description string) (interface{}, error) {
	req := &request.CreateDomainRequest{
		Name:        name,
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
				"text": formatCreateDomainResultForMCP(result),
			},
		},
		"isError": false,
	}, nil
}

// Helper functions for MCP formatting
func formatDomainForMCP(domain *response.DomainResponse) string {
	return fmt.Sprintf("Domain: %s\nDescription: %s\nCreated: %s",
		domain.Name, domain.Description, domain.CreatedAt.Format("2006-01-02 15:04:05"))
}

func formatCreateDomainResultForMCP(result *response.DomainResponse) string {
	return fmt.Sprintf("Successfully created domain: %s\nDescription: %s\nCreated: %s",
		result.Name, result.Description, result.CreatedAt.Format("2006-01-02 15:04:05"))
}
