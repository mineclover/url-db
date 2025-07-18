package mcp

import (
	"bufio"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"url-db/internal/models"
)

// MockMCPService for testing
type MockMCPService struct {
	mock.Mock
}

func (m *MockMCPService) CreateNode(ctx context.Context, req *models.CreateMCPNodeRequest) (*models.MCPNode, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPNode), args.Error(1)
}

func (m *MockMCPService) GetNode(ctx context.Context, compositeID string) (*models.MCPNode, error) {
	args := m.Called(ctx, compositeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPNode), args.Error(1)
}

func (m *MockMCPService) UpdateNode(ctx context.Context, compositeID string, req *models.UpdateNodeRequest) (*models.MCPNode, error) {
	args := m.Called(ctx, compositeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPNode), args.Error(1)
}

func (m *MockMCPService) DeleteNode(ctx context.Context, compositeID string) error {
	args := m.Called(ctx, compositeID)
	return args.Error(0)
}

func (m *MockMCPService) ListNodes(ctx context.Context, domainName string, page, size int, search string) (*models.MCPNodeListResponse, error) {
	args := m.Called(ctx, domainName, page, size, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPNodeListResponse), args.Error(1)
}

func (m *MockMCPService) FindNodeByURL(ctx context.Context, req *models.FindMCPNodeRequest) (*models.MCPNode, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPNode), args.Error(1)
}

func (m *MockMCPService) BatchGetNodes(ctx context.Context, req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BatchMCPNodeResponse), args.Error(1)
}

func (m *MockMCPService) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*MCPDomain, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MCPDomain), args.Error(1)
}

func (m *MockMCPService) ListDomains(ctx context.Context) (*MCPDomainListResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MCPDomainListResponse), args.Error(1)
}

func (m *MockMCPService) GetNodeAttributes(ctx context.Context, compositeID string) (*MCPNodeAttributeResponse, error) {
	args := m.Called(ctx, compositeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MCPNodeAttributeResponse), args.Error(1)
}

func (m *MockMCPService) SetNodeAttributes(ctx context.Context, compositeID string, req *models.SetMCPNodeAttributesRequest) (*MCPNodeAttributeResponse, error) {
	args := m.Called(ctx, compositeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MCPNodeAttributeResponse), args.Error(1)
}

func (m *MockMCPService) GetServerInfo(ctx context.Context) (*MCPServerInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MCPServerInfo), args.Error(1)
}

func TestNewStdioServer(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)

	assert.NotNil(t, server)
	assert.Equal(t, mockService, server.service)
	assert.NotNil(t, server.reader)
	assert.NotNil(t, server.writer)
}

func TestStdioServer_HandleRequest_ListDomains(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	// Setup mock response
	domains := &MCPDomainListResponse{
		Domains: []MCPDomain{
			{
				Name:        "example.com",
				Description: "Example domain",
				NodeCount:   5,
			},
			{
				Name:        "test.org",
				Description: "Test domain",
				NodeCount:   3,
			},
		},
	}
	
	mockService.On("ListDomains", mock.Anything).Return(domains, nil)
	
	// Create test writer to capture output
	var output strings.Builder
	server.writer = &output
	
	// Test list_domains command
	err := server.handleRequest("list_domains")
	
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Domains (2)")
	assert.Contains(t, output.String(), "example.com")
	assert.Contains(t, output.String(), "test.org")
	mockService.AssertExpectations(t)
}

func TestStdioServer_HandleRequest_CreateNode(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	// Setup mock response
	node := &models.MCPNode{
		CompositeID:  "example.com::https://example.com/page",
		DomainName:   "example.com",
		URL:          "https://example.com/page",
		Title:        "Test Page",
		Description:  "",
	}
	
	mockService.On("CreateNode", mock.Anything, mock.MatchedBy(func(req *models.CreateMCPNodeRequest) bool {
		return req.DomainName == "example.com" && 
			   req.URL == "https://example.com/page" && 
			   req.Title == "Test Page"
	})).Return(node, nil)
	
	// Create test writer to capture output
	var output strings.Builder
	server.writer = &output
	
	// Test create_node command
	err := server.handleRequest("create_node example.com https://example.com/page Test Page")
	
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Created node:")
	assert.Contains(t, output.String(), "example.com::https://example.com/page")
	mockService.AssertExpectations(t)
}

func TestStdioServer_HandleRequest_GetNode(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	// Setup mock response
	node := &models.MCPNode{
		CompositeID:  "example.com::https://example.com/page",
		DomainName:   "example.com",
		URL:          "https://example.com/page",
		Title:        "Test Page",
		Description:  "Test description",
		CreatedAt:    "2024-01-01T00:00:00Z",
		UpdatedAt:    "2024-01-01T00:00:00Z",
	}
	
	mockService.On("GetNode", mock.Anything, "example.com::https://example.com/page").Return(node, nil)
	
	// Create test writer to capture output
	var output strings.Builder
	server.writer = &output
	
	// Test get_node command
	err := server.handleRequest("get_node example.com::https://example.com/page")
	
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Node: example.com::https://example.com/page")
	assert.Contains(t, output.String(), "Domain: example.com")
	assert.Contains(t, output.String(), "URL: https://example.com/page")
	assert.Contains(t, output.String(), "Title: Test Page")
	mockService.AssertExpectations(t)
}

func TestStdioServer_HandleRequest_ListNodes(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	// Setup mock response
	response := &models.MCPNodeListResponse{
		Nodes: []models.MCPNode{
			{
				CompositeID: "example.com::https://example.com/page1",
				URL:         "https://example.com/page1",
				Title:       "Page 1",
			},
			{
				CompositeID: "example.com::https://example.com/page2",
				URL:         "https://example.com/page2",
				Title:       "Page 2",
			},
		},
		TotalCount: 2,
	}
	
	mockService.On("ListNodes", mock.Anything, "example.com", 1, 20, "").Return(response, nil)
	
	// Create test writer to capture output
	var output strings.Builder
	server.writer = &output
	
	// Test list_nodes command
	err := server.handleRequest("list_nodes example.com")
	
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Nodes in domain 'example.com' (2)")
	assert.Contains(t, output.String(), "example.com::https://example.com/page1")
	assert.Contains(t, output.String(), "example.com::https://example.com/page2")
	mockService.AssertExpectations(t)
}

func TestStdioServer_HandleRequest_ServerInfo(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	// Setup mock response
	info := &MCPServerInfo{
		Name:               "URL Database Server",
		Version:            "1.0.0",
		Description:        "MCP-enabled URL management",
		Capabilities:       []string{"domains", "nodes", "attributes"},
		CompositeKeyFormat: "domain_name::url_path",
	}
	
	mockService.On("GetServerInfo", mock.Anything).Return(info, nil)
	
	// Create test writer to capture output
	var output strings.Builder
	server.writer = &output
	
	// Test server_info command
	err := server.handleRequest("server_info")
	
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Server Info:")
	assert.Contains(t, output.String(), "Name: URL Database Server")
	assert.Contains(t, output.String(), "Version: 1.0.0")
	assert.Contains(t, output.String(), "Composite Key Format: domain_name::url_path")
	mockService.AssertExpectations(t)
}

func TestStdioServer_HandleRequest_Help(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	// Create test writer to capture output
	var output strings.Builder
	server.writer = &output
	
	// Test help command
	err := server.handleRequest("help")
	
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Available commands:")
	assert.Contains(t, output.String(), "list_domains")
	assert.Contains(t, output.String(), "create_node")
	assert.Contains(t, output.String(), "get_node")
	assert.Contains(t, output.String(), "quit")
}

func TestStdioServer_HandleRequest_UnknownCommand(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	// Test unknown command
	err := server.handleRequest("unknown_command")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command: unknown_command")
}

func TestStdioServer_HandleRequest_InsufficientArgs(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	tests := []struct {
		command string
		error   string
	}{
		{"list_nodes", "list_nodes requires domain_name argument"},
		{"create_node", "create_node requires domain_name and url arguments"},
		{"create_node example.com", "create_node requires domain_name and url arguments"},
		{"get_node", "get_node requires composite_id argument"},
		{"update_node", "update_node requires composite_id and title arguments"},
		{"delete_node", "delete_node requires composite_id argument"},
	}
	
	for _, test := range tests {
		t.Run(test.command, func(t *testing.T) {
			err := server.handleRequest(test.command)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.error)
		})
	}
}

func TestStdioServer_HandleRequest_UpdateNode(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	// Setup mock response
	node := &models.MCPNode{
		CompositeID: "example.com::https://example.com/page",
		Title:       "Updated Title",
	}
	
	mockService.On("UpdateNode", mock.Anything, "example.com::https://example.com/page", mock.MatchedBy(func(req *models.UpdateNodeRequest) bool {
		return req.Title == "Updated Title"
	})).Return(node, nil)
	
	// Create test writer to capture output
	var output strings.Builder
	server.writer = &output
	
	// Test update_node command
	err := server.handleRequest("update_node example.com::https://example.com/page Updated Title")
	
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Updated node:")
	assert.Contains(t, output.String(), "Title: Updated Title")
	mockService.AssertExpectations(t)
}

func TestStdioServer_HandleRequest_DeleteNode(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	mockService.On("DeleteNode", mock.Anything, "example.com::https://example.com/page").Return(nil)
	
	// Create test writer to capture output
	var output strings.Builder
	server.writer = &output
	
	// Test delete_node command
	err := server.handleRequest("delete_node example.com::https://example.com/page")
	
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Deleted node: example.com::https://example.com/page")
	mockService.AssertExpectations(t)
}

func TestStdioServer_HandleRequest_EmptyCommand(t *testing.T) {
	mockService := &MockMCPService{}
	server := NewStdioServer(mockService)
	
	// Test empty command
	err := server.handleRequest("")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty command")
}