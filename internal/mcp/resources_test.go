package mcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"url-db/internal/models"
)

func TestResourceRegistry_GetResources(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewResourceRegistry(mockService)

	// Mock domains response
	domainsResponse := &MCPDomainListResponse{
		Domains: []MCPDomain{
			{Name: "test-domain", Description: "Test domain"},
			{Name: "example-domain", Description: "Example domain"},
		},
	}

	mockService.On("ListDomains", mock.Anything).Return(domainsResponse, nil)

	result, err := registry.GetResources(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Should have server info + 2 domains * 2 resources each = 5 total
	expectedResourceCount := 1 + (2 * 2) // server info + domain info + domain nodes for each domain
	assert.Len(t, result.Resources, expectedResourceCount)

	// Check for expected resources
	resourceURIs := make(map[string]bool)
	for _, resource := range result.Resources {
		resourceURIs[resource.URI] = true
		assert.NotEmpty(t, resource.Name)
		assert.NotEmpty(t, resource.Description)
		assert.Equal(t, "application/json", resource.MimeType)
	}

	// Verify expected URIs
	assert.True(t, resourceURIs["mcp://server/info"])
	assert.True(t, resourceURIs["mcp://domains/test-domain"])
	assert.True(t, resourceURIs["mcp://domains/test-domain/nodes"])
	assert.True(t, resourceURIs["mcp://domains/example-domain"])
	assert.True(t, resourceURIs["mcp://domains/example-domain/nodes"])

	mockService.AssertExpectations(t)
}

func TestResourceRegistry_ReadResource_ServerInfo(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewResourceRegistry(mockService)

	serverInfo := &MCPServerInfo{
		Name:               "test-server",
		Version:            "1.0.0",
		Description:        "Test server",
		Capabilities:       []string{"tools", "resources"},
		CompositeKeyFormat: "test:domain:id",
	}

	mockService.On("GetServerInfo", mock.Anything).Return(serverInfo, nil)

	result, err := registry.ReadResource(context.Background(), "mcp://server/info")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "mcp://server/info", content.URI)
	assert.Equal(t, "application/json", content.MimeType)
	assert.Contains(t, content.Text, "test-server")
	assert.Contains(t, content.Text, "1.0.0")

	mockService.AssertExpectations(t)
}

func TestResourceRegistry_ReadResource_DomainInfo(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewResourceRegistry(mockService)

	domainsResponse := &MCPDomainListResponse{
		Domains: []MCPDomain{
			{Name: "test-domain", Description: "Test domain"},
		},
	}

	mockService.On("ListDomains", mock.Anything).Return(domainsResponse, nil)

	result, err := registry.ReadResource(context.Background(), "mcp://domains/test-domain")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "mcp://domains/test-domain", content.URI)
	assert.Equal(t, "application/json", content.MimeType)
	assert.Contains(t, content.Text, "test-domain")

	mockService.AssertExpectations(t)
}

func TestResourceRegistry_ReadResource_DomainNodes(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewResourceRegistry(mockService)

	nodesResponse := &models.MCPNodeListResponse{
		Nodes: []models.MCPNode{
			{CompositeID: "test:test-domain:1", URL: "https://example.com"},
		},
		TotalCount: 1,
	}

	mockService.On("ListNodes", mock.Anything, "test-domain", 1, 100, "").Return(nodesResponse, nil)

	result, err := registry.ReadResource(context.Background(), "mcp://domains/test-domain/nodes")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "mcp://domains/test-domain/nodes", content.URI)
	assert.Equal(t, "application/json", content.MimeType)
	assert.Contains(t, content.Text, "test:test-domain:1")

	mockService.AssertExpectations(t)
}

func TestResourceRegistry_ReadResource_NodeInfo(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewResourceRegistry(mockService)

	node := &models.MCPNode{
		CompositeID: "test:test-domain:1",
		URL:         "https://example.com",
		Title:       "Test Node",
	}

	mockService.On("GetNode", mock.Anything, "test:test-domain:1").Return(node, nil)
	mockService.On("GetNodeAttributes", mock.Anything, "test:test-domain:1").Return(nil, assert.AnError)

	result, err := registry.ReadResource(context.Background(), "mcp://nodes/test:test-domain:1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "mcp://nodes/test:test-domain:1", content.URI)
	assert.Equal(t, "application/json", content.MimeType)
	assert.Contains(t, content.Text, "test:test-domain:1")
	assert.Contains(t, content.Text, "https://example.com")

	mockService.AssertExpectations(t)
}

func TestResourceRegistry_ReadResource_UnknownURI(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewResourceRegistry(mockService)

	result, err := registry.ReadResource(context.Background(), "mcp://unknown/resource")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unknown resource URI")
}

func TestResourceRegistry_ReadResource_DomainNotFound(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewResourceRegistry(mockService)

	domainsResponse := &MCPDomainListResponse{
		Domains: []MCPDomain{},
	}

	mockService.On("ListDomains", mock.Anything).Return(domainsResponse, nil)

	result, err := registry.ReadResource(context.Background(), "mcp://domains/nonexistent-domain")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "domain not found")

	mockService.AssertExpectations(t)
}

func TestResourceRegistry_validateURI(t *testing.T) {
	registry := &ResourceRegistry{}

	// Valid URI
	err := registry.validateURI("mcp://server/info")
	assert.NoError(t, err)

	// Invalid URI scheme
	err = registry.validateURI("http://server/info")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid URI scheme")
}
