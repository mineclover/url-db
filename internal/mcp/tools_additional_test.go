package mcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"url-db/internal/models"
)

func TestToolRegistry_CallTool_ListNodes(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)

	expectedResponse := &models.MCPNodeListResponse{
		Nodes: []models.MCPNode{
			{CompositeID: "url-db:test-domain:1", URL: "https://example.com"},
		},
		TotalCount: 1,
	}

	mockService.On("ListNodes", mock.Anything, "test-domain", 1, 10, "").Return(expectedResponse, nil)

	arguments := map[string]interface{}{
		"domain_name": "test-domain",
		"page":        float64(1),
		"size":        float64(10),
		"search":      "",
	}

	result, err := registry.CallTool(context.Background(), "list_nodes", arguments)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "url-db:test-domain:1")

	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_CreateNode(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)

	expectedNode := &models.MCPNode{
		CompositeID: "url-db:test-domain:1",
		URL:         "https://example.com",
		Title:       "Test Node",
	}

	mockService.On("CreateNode", mock.Anything, mock.AnythingOfType("*models.CreateMCPNodeRequest")).Return(expectedNode, nil)

	arguments := map[string]interface{}{
		"domain_name": "test-domain",
		"url":         "https://example.com",
		"title":       "Test Node",
	}

	result, err := registry.CallTool(context.Background(), "create_node", arguments)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "url-db:test-domain:1")

	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_GetNode(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)

	expectedNode := &models.MCPNode{
		CompositeID: "url-db:test-domain:1",
		URL:         "https://example.com",
		Title:       "Test Node",
	}

	mockService.On("GetNode", mock.Anything, "url-db:test-domain:1").Return(expectedNode, nil)

	arguments := map[string]interface{}{
		"composite_id": "url-db:test-domain:1",
	}

	result, err := registry.CallTool(context.Background(), "get_mcp_node", arguments)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "url-db:test-domain:1")

	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_UpdateNode(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)

	expectedNode := &models.MCPNode{
		CompositeID: "url-db:test-domain:1",
		URL:         "https://example.com",
		Title:       "Updated Node",
	}

	mockService.On("UpdateNode", mock.Anything, "url-db:test-domain:1", mock.AnythingOfType("*models.UpdateNodeRequest")).Return(expectedNode, nil)

	arguments := map[string]interface{}{
		"composite_id": "url-db:test-domain:1",
		"title":        "Updated Node",
	}

	result, err := registry.CallTool(context.Background(), "update_mcp_node", arguments)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "Updated Node")

	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_DeleteNode(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)

	mockService.On("DeleteNode", mock.Anything, "url-db:test-domain:1").Return(nil)

	arguments := map[string]interface{}{
		"composite_id": "url-db:test-domain:1",
	}

	result, err := registry.CallTool(context.Background(), "delete_mcp_node", arguments)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "deleted successfully")

	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_FindNodeByURL(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)

	expectedNode := &models.MCPNode{
		CompositeID: "url-db:test-domain:1",
		URL:         "https://example.com",
		Title:       "Found Node",
	}

	mockService.On("FindNodeByURL", mock.Anything, mock.AnythingOfType("*models.FindMCPNodeRequest")).Return(expectedNode, nil)

	arguments := map[string]interface{}{
		"domain_name": "test-domain",
		"url":         "https://example.com",
	}

	result, err := registry.CallTool(context.Background(), "find_mcp_node_by_url", arguments)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "Found Node")

	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_GetNodeAttributes(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)

	expectedResponse := &models.MCPNodeAttributeResponse{
		CompositeID: "url-db:test-domain:1",
		Attributes: []models.MCPAttribute{
			{Name: "category", Value: "technology"},
		},
	}

	mockService.On("GetNodeAttributes", mock.Anything, "url-db:test-domain:1").Return(expectedResponse, nil)

	arguments := map[string]interface{}{
		"composite_id": "url-db:test-domain:1",
	}

	result, err := registry.CallTool(context.Background(), "get_mcp_node_attributes", arguments)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "category")

	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_SetNodeAttributes(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)

	expectedResponse := &models.MCPNodeAttributeResponse{
		CompositeID: "url-db:test-domain:1",
		Attributes: []models.MCPAttribute{
			{Name: "category", Value: "technology"},
		},
	}

	mockService.On("SetNodeAttributes", mock.Anything, "url-db:test-domain:1", mock.AnythingOfType("*models.SetMCPNodeAttributesRequest")).Return(expectedResponse, nil)

	arguments := map[string]interface{}{
		"composite_id": "url-db:test-domain:1",
		"attributes": []interface{}{
			map[string]interface{}{
				"name":  "category",
				"value": "technology",
			},
		},
	}

	result, err := registry.CallTool(context.Background(), "set_mcp_node_attributes", arguments)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "category")

	mockService.AssertExpectations(t)
}

func TestToolRegistry_CallTool_GetServerInfo(t *testing.T) {
	mockService := new(MockMCPService)
	registry := NewToolRegistry(mockService)

	expectedInfo := &MCPServerInfo{
		Name:        "url-db",
		Version:     "1.0.0",
		Description: "URL Database MCP Server",
	}

	mockService.On("GetServerInfo", mock.Anything).Return(expectedInfo, nil)

	result, err := registry.CallTool(context.Background(), "get_mcp_server_info", map[string]interface{}{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Contains(t, result.Content[0].Text, "url-db")

	mockService.AssertExpectations(t)
}
