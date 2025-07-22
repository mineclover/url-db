package mcp

import (
	"context"
	"fmt"

	"url-db/internal/models"
)

// Node operation methods for mcpService

func (s *mcpService) CreateNode(ctx context.Context, req *models.CreateMCPNodeRequest) (*models.MCPNode, error) {
	domain, err := s.domainService.GetDomainByName(ctx, req.DomainName)
	if err != nil {
		return nil, NewDomainNotFoundError(req.DomainName)
	}

	nodeReq := s.converter.CreateMCPNodeRequestToCreateNodeRequest(req)
	node, err := s.nodeService.CreateNode(ctx, domain.ID, nodeReq)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to create node: %v", err))
	}

	return s.converter.NodeToMCPNode(node, domain)
}

func (s *mcpService) GetNode(ctx context.Context, compositeID string) (*models.MCPNode, error) {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domainName, err := s.converter.ExtractDomainNameFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return nil, NewResourceNotFoundError(compositeID)
	}

	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	return s.converter.NodeToMCPNode(node, domain)
}

func (s *mcpService) UpdateNode(ctx context.Context, compositeID string, req *models.UpdateNodeRequest) (*models.MCPNode, error) {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domainName, err := s.converter.ExtractDomainNameFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	node, err := s.nodeService.UpdateNode(ctx, nodeID, req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to update node: %v", err))
	}

	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	return s.converter.NodeToMCPNode(node, domain)
}

func (s *mcpService) DeleteNode(ctx context.Context, compositeID string) error {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return NewInvalidCompositeKeyError(compositeID)
	}

	if err := s.nodeService.DeleteNode(ctx, nodeID); err != nil {
		return NewInternalServerError(fmt.Sprintf("failed to delete node: %v", err))
	}

	return nil
}

func (s *mcpService) ListNodes(ctx context.Context, domainName string, page, size int, search string) (*models.MCPNodeListResponse, error) {
	var domainID int
	var domain *models.Domain
	var err error

	if domainName != "" {
		domain, err = s.domainService.GetDomainByName(ctx, domainName)
		if err != nil {
			return nil, NewDomainNotFoundError(domainName)
		}
		domainID = domain.ID
	}

	response, err := s.nodeService.ListNodes(ctx, domainID, page, size, search)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list nodes: %v", err))
	}

	mcpNodes := make([]models.MCPNode, 0, len(response.Nodes))
	for _, node := range response.Nodes {
		if domain == nil {
			domain, err = s.domainService.GetDomain(ctx, node.DomainID)
			if err != nil {
				continue
			}
		}

		mcpNode, err := s.converter.NodeToMCPNode(&node, domain)
		if err != nil {
			continue
		}
		mcpNodes = append(mcpNodes, *mcpNode)
	}

	return &models.MCPNodeListResponse{
		Nodes:      mcpNodes,
		TotalCount: response.TotalCount,
		Page:       response.Page,
		Size:       response.Size,
		TotalPages: response.TotalPages,
	}, nil
}

func (s *mcpService) FindNodeByURL(ctx context.Context, req *models.FindMCPNodeRequest) (*models.MCPNode, error) {
	domain, err := s.domainService.GetDomainByName(ctx, req.DomainName)
	if err != nil {
		return nil, NewDomainNotFoundError(req.DomainName)
	}

	findReq := &models.FindNodeByURLRequest{URL: req.URL}
	node, err := s.nodeService.FindNodeByURL(ctx, domain.ID, findReq)
	if err != nil {
		return nil, NewNodeNotFoundError(req.DomainName, req.URL)
	}

	return s.converter.NodeToMCPNode(node, domain)
}

func (s *mcpService) BatchGetNodes(ctx context.Context, req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error) {
	nodes := make([]models.MCPNode, 0, len(req.CompositeIDs))
	notFound := make([]string, 0)

	for _, compositeID := range req.CompositeIDs {
		node, err := s.GetNode(ctx, compositeID)
		if err != nil {
			notFound = append(notFound, compositeID)
			continue
		}
		nodes = append(nodes, *node)
	}

	return &models.BatchMCPNodeResponse{
		Nodes:    nodes,
		NotFound: notFound,
	}, nil
}