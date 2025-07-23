package mcp

import (
	"context"
	"fmt"

	"url-db/internal/models"
)

// Domain operation methods for mcpService

func (s *mcpService) ListDomains(ctx context.Context) (*MCPDomainListResponse, error) {
	response, err := s.domainService.ListDomains(ctx, 1, 1000)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list domains: %v", err))
	}

	mcpDomains := make([]MCPDomain, 0, len(response.Domains))
	for _, domain := range response.Domains {
		nodeCount, err := s.nodeCountService.GetNodeCountByDomain(ctx, domain.ID)
		if err != nil {
			nodeCount = 0
		}

		mcpDomain := s.converter.DomainToMCPDomain(&domain, nodeCount)
		if mcpDomain != nil {
			mcpDomains = append(mcpDomains, *mcpDomain)
		}
	}

	return &MCPDomainListResponse{
		Domains: mcpDomains,
	}, nil
}

func (s *mcpService) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*MCPDomain, error) {
	domain, err := s.domainService.CreateDomain(ctx, req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to create domain: %v", err))
	}

	return s.converter.DomainToMCPDomain(domain, 0), nil
}
