package mcp

import (
	"context"
	"fmt"
	"strings"

	"url-db/internal/models"
)

// Query operation methods for mcpService - enhanced query capabilities

func (s *mcpService) FilterNodesByAttributes(ctx context.Context, domainName string, filters []interface{}, page, size int) (*models.MCPNodeListResponse, error) {
	// Get domain
	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	// Get all nodes in the domain first
	response, err := s.nodeService.ListNodes(ctx, domain.ID, 1, 1000, "")
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list nodes: %v", err))
	}

	// Filter nodes based on attributes
	var filteredNodes []models.MCPNode
	for _, node := range response.Nodes {
		// Get node attributes
		nodeAttrs, err := s.attributeService.GetNodeAttributes(ctx, node.ID)
		if err != nil {
			continue
		}

		// Check if node matches all filters
		matches := true
		for _, filterInterface := range filters {
			filterMap, ok := filterInterface.(map[string]interface{})
			if !ok {
				// Try struct conversion
				filterStruct, ok := filterInterface.(struct {
					Name     string `json:"name"`
					Value    string `json:"value"`
					Operator string `json:"operator"`
				})
				if ok {
					filterMap = map[string]interface{}{
						"name":     filterStruct.Name,
						"value":    filterStruct.Value,
						"operator": filterStruct.Operator,
					}
				} else {
					continue
				}
			}

			filterName, _ := filterMap["name"].(string)
			filterValue, _ := filterMap["value"].(string)
			filterOperator, _ := filterMap["operator"].(string)

			if filterName == "" || filterValue == "" {
				continue
			}

			// Check if node has this attribute with matching value
			attrFound := false
			for _, attr := range nodeAttrs {
				if attr.Name == filterName {
					switch filterOperator {
					case "equals":
						if attr.Value == filterValue {
							attrFound = true
						}
					case "contains":
						if strings.Contains(attr.Value, filterValue) {
							attrFound = true
						}
					case "starts_with":
						if strings.HasPrefix(attr.Value, filterValue) {
							attrFound = true
						}
					case "ends_with":
						if strings.HasSuffix(attr.Value, filterValue) {
							attrFound = true
						}
					}
					break
				}
			}

			if !attrFound {
				matches = false
				break
			}
		}

		if matches {
			mcpNode, err := s.converter.NodeToMCPNode(&node, domain)
			if err == nil {
				filteredNodes = append(filteredNodes, *mcpNode)
			}
		}
	}

	// Apply pagination
	start := (page - 1) * size
	end := start + size
	if start > len(filteredNodes) {
		start = len(filteredNodes)
	}
	if end > len(filteredNodes) {
		end = len(filteredNodes)
	}

	paginatedNodes := filteredNodes[start:end]

	return &models.MCPNodeListResponse{
		Nodes:      paginatedNodes,
		TotalCount: len(filteredNodes),
		Page:       page,
		Size:       size,
		TotalPages: (len(filteredNodes) + size - 1) / size,
	}, nil
}

func (s *mcpService) GetServerInfo(ctx context.Context) (*MCPServerInfo, error) {
	return &MCPServerInfo{
		Name:        "url-db",
		Version:     "1.0.0",
		Description: "URL 데이터베이스 MCP 서버",
		Capabilities: []string{
			"resources",
			"tools",
			"prompts",
		},
		CompositeKeyFormat: "url-db:domain_name:id",
	}, nil
}
