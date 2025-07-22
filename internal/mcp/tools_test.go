package mcp

import (
	"context"
	"testing"

	"url-db/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock MCP Service for testing
type MockMCPService struct {
	mock.Mock
}

func (m *MockMCPService) ListDomains(ctx context.Context) (*MCPDomainListResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MCPDomainListResponse), args.Error(1)
}

func (m *MockMCPService) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*MCPDomain, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MCPDomain), args.Error(1)
}

func (m *MockMCPService) ListNodes(ctx context.Context, domainName string, page, size int, search string) (*models.MCPNodeListResponse, error) {
	args := m.Called(ctx, domainName, page, size, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPNodeListResponse), args.Error(1)
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

func (m *MockMCPService) FindNodeByURL(ctx context.Context, req *models.FindMCPNodeRequest) (*models.MCPNode, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPNode), args.Error(1)
}

func (m *MockMCPService) GetNodeAttributes(ctx context.Context, compositeID string) (*models.MCPNodeAttributeResponse, error) {
	args := m.Called(ctx, compositeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPNodeAttributeResponse), args.Error(1)
}

func (m *MockMCPService) SetNodeAttributes(ctx context.Context, compositeID string, req *models.SetMCPNodeAttributesRequest) (*models.MCPNodeAttributeResponse, error) {
	args := m.Called(ctx, compositeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPNodeAttributeResponse), args.Error(1)
}

func (m *MockMCPService) GetServerInfo(ctx context.Context) (*MCPServerInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MCPServerInfo), args.Error(1)
}

// Domain attribute management methods
func (m *MockMCPService) ListDomainAttributes(ctx context.Context, domainName string) (*models.MCPDomainAttributeListResponse, error) {
	args := m.Called(ctx, domainName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPDomainAttributeListResponse), args.Error(1)
}

func (m *MockMCPService) CreateDomainAttribute(ctx context.Context, domainName string, req *models.CreateAttributeRequest) (*models.MCPDomainAttribute, error) {
	args := m.Called(ctx, domainName, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPDomainAttribute), args.Error(1)
}

func (m *MockMCPService) GetDomainAttribute(ctx context.Context, compositeID string) (*models.MCPDomainAttribute, error) {
	args := m.Called(ctx, compositeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPDomainAttribute), args.Error(1)
}

func (m *MockMCPService) UpdateDomainAttribute(ctx context.Context, compositeID string, req *models.UpdateAttributeRequest) (*models.MCPDomainAttribute, error) {
	args := m.Called(ctx, compositeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MCPDomainAttribute), args.Error(1)
}

func (m *MockMCPService) DeleteDomainAttribute(ctx context.Context, compositeID string) error {
	args := m.Called(ctx, compositeID)
	return args.Error(0)
}

func (m *MockMCPService) BatchGetNodes(ctx context.Context, req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BatchMCPNodeResponse), args.Error(1)
}

func TestToolRegistry_GetTools(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)
	
	tools := registry.GetTools()
	
	// Check that all expected tools are present
	expectedTools := []string{
		"list_mcp_domains", "create_mcp_domain", "list_mcp_nodes",
		"create_mcp_node", "get_mcp_node", "update_mcp_node",
		"delete_mcp_node", "find_mcp_node_by_url", "get_mcp_node_attributes",
		"set_mcp_node_attributes", "get_mcp_server_info",
		"list_mcp_domain_attributes", "create_mcp_domain_attribute",
		"get_mcp_domain_attribute", "update_mcp_domain_attribute",
		"delete_mcp_domain_attribute",
	}
	
	assert.Len(t, tools, len(expectedTools))
	
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
		// Check that each tool has required fields
		assert.NotEmpty(t, tool.Name)
		assert.NotEmpty(t, tool.Description)
		assert.NotNil(t, tool.InputSchema)
	}
	
	// Verify all expected tools are present
	for _, expectedTool := range expectedTools {
		assert.True(t, toolNames[expectedTool], "Missing tool: %s", expectedTool)
	}
}

func TestToolRegistry_CallTool_ListDomains(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)
	
	expectedResponse := &MCPDomainListResponse{
		Domains: []MCPDomain{
			{Name: "test-domain", Description: "Test domain"},
		},
	}
	
	mockService.On("ListDomains", mock.Anything).Return(expectedResponse, nil)
	
	result, err := registry.CallTool(context.Background(), "list_mcp_domains", map[string]interface{}{})
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Equal(t, "text", result.Content[0].Type)
	assert.Contains(t, result.Content[0].Text, "test-domain")
	
	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_CreateDomain(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)
	
	expectedDomain := &MCPDomain{
		Name:        "new-domain",
		Description: "New test domain",
	}
	
	mockService.On("CreateDomain", mock.Anything, mock.AnythingOfType("*models.CreateDomainRequest")).Return(expectedDomain, nil)
	
	arguments := map[string]interface{}{
		"name":        "new-domain",
		"description": "New test domain",
	}
	
	result, err := registry.CallTool(context.Background(), "create_mcp_domain", arguments)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "new-domain")
	
	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_UnknownTool(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)
	
	result, err := registry.CallTool(context.Background(), "unknown_tool", map[string]interface{}{})
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsError)
	assert.Contains(t, result.Content[0].Text, "Unknown tool")
}

func TestGetStringArg(t *testing.T) {
	args := map[string]interface{}{
		"test": "value",
		"num":  123,
	}
	
	assert.Equal(t, "value", getStringArg(args, "test"))
	assert.Equal(t, "", getStringArg(args, "missing"))
	assert.Equal(t, "", getStringArg(args, "num"))
}

// Domain attribute management tests
func TestToolRegistry_CallTool_ListDomainAttributes(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)
	
	expectedResponse := &models.MCPDomainAttributeListResponse{
		DomainName: "test-domain",
		Attributes: []models.MCPDomainAttribute{
			{
				CompositeID: "url-db:test-domain:1",
				Name:        "status",
				Type:        "tag",
				Description: "Project status",
			},
		},
		TotalCount: 1,
	}
	
	mockService.On("ListDomainAttributes", mock.Anything, "test-domain").Return(expectedResponse, nil)
	
	arguments := map[string]interface{}{
		"domain_name": "test-domain",
	}
	
	result, err := registry.CallTool(context.Background(), "list_mcp_domain_attributes", arguments)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "test-domain")
	
	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_CreateDomainAttribute(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)
	
	expectedAttribute := &models.MCPDomainAttribute{
		CompositeID: "url-db:test-domain:1",
		Name:        "priority",
		Type:        "ordered_tag",
		Description: "Project priority",
	}
	
	mockService.On("CreateDomainAttribute", mock.Anything, "test-domain", mock.AnythingOfType("*models.CreateAttributeRequest")).Return(expectedAttribute, nil)
	
	arguments := map[string]interface{}{
		"domain_name": "test-domain",
		"name":        "priority",
		"type":        "ordered_tag",
		"description": "Project priority",
	}
	
	result, err := registry.CallTool(context.Background(), "create_mcp_domain_attribute", arguments)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "url-db:test-domain:1")
	
	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_GetDomainAttribute(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)
	
	expectedAttribute := &models.MCPDomainAttribute{
		CompositeID: "url-db:test-domain:1",
		Name:        "status",
		Type:        "tag",
		Description: "Project status",
	}
	
	mockService.On("GetDomainAttribute", mock.Anything, "url-db:test-domain:1").Return(expectedAttribute, nil)
	
	arguments := map[string]interface{}{
		"composite_id": "url-db:test-domain:1",
	}
	
	result, err := registry.CallTool(context.Background(), "get_mcp_domain_attribute", arguments)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "status")
	
	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_UpdateDomainAttribute(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)
	
	expectedAttribute := &models.MCPDomainAttribute{
		CompositeID: "url-db:test-domain:1",
		Name:        "status",
		Type:        "tag",
		Description: "Updated project status",
	}
	
	mockService.On("UpdateDomainAttribute", mock.Anything, "url-db:test-domain:1", mock.AnythingOfType("*models.UpdateAttributeRequest")).Return(expectedAttribute, nil)
	
	arguments := map[string]interface{}{
		"composite_id": "url-db:test-domain:1",
		"description":  "Updated project status",
	}
	
	result, err := registry.CallTool(context.Background(), "update_mcp_domain_attribute", arguments)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "Updated project status")
	
	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_DeleteDomainAttribute(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)
	
	mockService.On("DeleteDomainAttribute", mock.Anything, "url-db:test-domain:1").Return(nil)
	
	arguments := map[string]interface{}{
		"composite_id": "url-db:test-domain:1",
	}
	
	result, err := registry.CallTool(context.Background(), "delete_mcp_domain_attribute", arguments)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "deleted successfully")
	
	mockService.AssertExpectations(t)
}