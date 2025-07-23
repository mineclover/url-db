package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// MCP Resource 정의 및 구현

// ResourceRegistry manages all available MCP resources
type ResourceRegistry struct {
	service MCPService
}

// NewResourceRegistry creates a new resource registry
func NewResourceRegistry(service MCPService) *ResourceRegistry {
	return &ResourceRegistry{
		service: service,
	}
}

// GetResources returns all available resources
func (rr *ResourceRegistry) GetResources(ctx context.Context) (*ResourcesListResult, error) {
	// Get all domains to generate resource URIs
	domainsResponse, err := rr.service.ListDomains(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	var resources []Resource

	// Add server info resource
	resources = append(resources, Resource{
		URI:         "mcp://server/info",
		Name:        "Server Information",
		Description: "URL Database server information and capabilities",
		MimeType:    "application/json",
	})

	// Add domain resources
	for _, domain := range domainsResponse.Domains {
		// Domain resource
		resources = append(resources, Resource{
			URI:         fmt.Sprintf("mcp://domains/%s", domain.Name),
			Name:        fmt.Sprintf("Domain: %s", domain.Name),
			Description: fmt.Sprintf("Domain information for %s", domain.Name),
			MimeType:    "application/json",
		})

		// Domain nodes resource
		resources = append(resources, Resource{
			URI:         fmt.Sprintf("mcp://domains/%s/nodes", domain.Name),
			Name:        fmt.Sprintf("Nodes in %s", domain.Name),
			Description: fmt.Sprintf("List of all nodes in domain %s", domain.Name),
			MimeType:    "application/json",
		})
	}

	return &ResourcesListResult{
		Resources: resources,
	}, nil
}

// ReadResource reads a specific resource by URI
func (rr *ResourceRegistry) ReadResource(ctx context.Context, uri string) (*ReadResourceResult, error) {
	// Parse the URI to determine resource type
	if uri == "mcp://server/info" {
		return rr.readServerInfo(ctx)
	}

	// Match domain resource pattern: mcp://domains/{domain_name}
	domainPattern := regexp.MustCompile(`^mcp://domains/([^/]+)$`)
	if matches := domainPattern.FindStringSubmatch(uri); len(matches) == 2 {
		domainName := matches[1]
		return rr.readDomainInfo(ctx, domainName)
	}

	// Match domain nodes pattern: mcp://domains/{domain_name}/nodes
	domainNodesPattern := regexp.MustCompile(`^mcp://domains/([^/]+)/nodes$`)
	if matches := domainNodesPattern.FindStringSubmatch(uri); len(matches) == 2 {
		domainName := matches[1]
		return rr.readDomainNodes(ctx, domainName)
	}

	// Match individual node pattern: mcp://nodes/{composite_id}
	nodePattern := regexp.MustCompile(`^mcp://nodes/(.+)$`)
	if matches := nodePattern.FindStringSubmatch(uri); len(matches) == 2 {
		compositeID := matches[1]
		return rr.readNodeInfo(ctx, compositeID)
	}

	return nil, fmt.Errorf("unknown resource URI: %s", uri)
}

// readServerInfo reads server information
func (rr *ResourceRegistry) readServerInfo(ctx context.Context) (*ReadResourceResult, error) {
	info, err := rr.service.GetServerInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server info: %w", err)
	}

	content, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal server info: %w", err)
	}

	return &ReadResourceResult{
		Contents: []ResourceContent{{
			URI:      "mcp://server/info",
			MimeType: "application/json",
			Text:     string(content),
		}},
	}, nil
}

// readDomainInfo reads domain information
func (rr *ResourceRegistry) readDomainInfo(ctx context.Context, domainName string) (*ReadResourceResult, error) {
	// Get domain info by listing all domains and finding the one we want
	domainsResponse, err := rr.service.ListDomains(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	var targetDomain *MCPDomain
	for _, domain := range domainsResponse.Domains {
		if domain.Name == domainName {
			targetDomain = &domain
			break
		}
	}

	if targetDomain == nil {
		return nil, fmt.Errorf("domain not found: %s", domainName)
	}

	content, err := json.MarshalIndent(targetDomain, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal domain info: %w", err)
	}

	return &ReadResourceResult{
		Contents: []ResourceContent{{
			URI:      fmt.Sprintf("mcp://domains/%s", domainName),
			MimeType: "application/json",
			Text:     string(content),
		}},
	}, nil
}

// readDomainNodes reads all nodes in a domain
func (rr *ResourceRegistry) readDomainNodes(ctx context.Context, domainName string) (*ReadResourceResult, error) {
	// Get first page of nodes (we might want to implement pagination later)
	nodesResponse, err := rr.service.ListNodes(ctx, domainName, 1, 100, "")
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes for domain %s: %w", domainName, err)
	}

	content, err := json.MarshalIndent(nodesResponse, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal nodes: %w", err)
	}

	return &ReadResourceResult{
		Contents: []ResourceContent{{
			URI:      fmt.Sprintf("mcp://domains/%s/nodes", domainName),
			MimeType: "application/json",
			Text:     string(content),
		}},
	}, nil
}

// readNodeInfo reads individual node information
func (rr *ResourceRegistry) readNodeInfo(ctx context.Context, compositeID string) (*ReadResourceResult, error) {
	// Get node info
	node, err := rr.service.GetNode(ctx, compositeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s: %w", compositeID, err)
	}

	// Also get node attributes
	attributes, err := rr.service.GetNodeAttributes(ctx, compositeID)
	if err != nil {
		// Don't fail if attributes can't be fetched, just log it
		attributes = nil
	}

	// Combine node info and attributes
	result := map[string]interface{}{
		"node":       node,
		"attributes": attributes,
	}

	content, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal node info: %w", err)
	}

	return &ReadResourceResult{
		Contents: []ResourceContent{{
			URI:      fmt.Sprintf("mcp://nodes/%s", compositeID),
			MimeType: "application/json",
			Text:     string(content),
		}},
	}, nil
}

// Helper function to validate resource URI format
func (rr *ResourceRegistry) validateURI(uri string) error {
	if !strings.HasPrefix(uri, "mcp://") {
		return fmt.Errorf("invalid URI scheme, must start with 'mcp://'")
	}
	return nil
}
