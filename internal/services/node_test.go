package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"url-db/internal/models"
)

func TestNodeService_CreateNode(t *testing.T) {
	service, _, mockDomainRepo := CreateTestNodeService(t)

	ctx := CreateTestContext()
	
	// Create domain first
	testDomain := CreateTestDomain("test-domain", "Test domain")
	mockDomainRepo.Create(ctx, testDomain)

	req := &models.CreateNodeRequest{
		URL:         "https://example.com",
		Title:       "Example",
		Description: "Example website",
	}

	node, err := service.CreateNode(ctx, testDomain.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, "https://example.com", node.Content)
	assert.Equal(t, "Example", node.Title)
}

func TestNodeService_GetNode(t *testing.T) {
	service, mockNodeRepo, mockDomainRepo := CreateTestNodeService(t)

	ctx := CreateTestContext()
	testDomain := CreateTestDomain("test", "Test")
	mockDomainRepo.Create(ctx, testDomain)

	testNode := CreateTestNode(testDomain.ID, "https://example.com", "Example", "Test node")
	mockNodeRepo.Create(ctx, testNode)

	node, err := service.GetNode(ctx, testNode.ID)

	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, testNode.ID, node.ID)
	assert.Equal(t, testNode.Content, node.Content)
}

func TestNodeService_ListNodesByDomain(t *testing.T) {
	service, mockNodeRepo, mockDomainRepo := CreateTestNodeService(t)

	ctx := CreateTestContext()
	testDomain := CreateTestDomain("test", "Test")
	mockDomainRepo.Create(ctx, testDomain)

	node1 := CreateTestNode(testDomain.ID, "https://example1.com", "Example1", "Node 1")
	node2 := CreateTestNode(testDomain.ID, "https://example2.com", "Example2", "Node 2")
	
	mockNodeRepo.Create(ctx, node1)
	mockNodeRepo.Create(ctx, node2)

	response, err := service.ListNodesByDomain(ctx, testDomain.ID, 1, 20, "")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, response.TotalCount)
	assert.Len(t, response.Nodes, 2)
}

func TestNodeService_UpdateNode(t *testing.T) {
	service, mockNodeRepo, mockDomainRepo := CreateTestNodeService(t)

	ctx := CreateTestContext()
	testDomain := CreateTestDomain("test", "Test")
	mockDomainRepo.Create(ctx, testDomain)

	testNode := CreateTestNode(testDomain.ID, "https://example.com", "Original", "Original description")
	mockNodeRepo.Create(ctx, testNode)

	req := &models.UpdateNodeRequest{
		Title:       "Updated Title",
		Description: "Updated description",
	}

	node, err := service.UpdateNode(ctx, testNode.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, "Updated Title", node.Title)
	assert.Equal(t, "Updated description", node.Description)
}

func TestNodeService_DeleteNode(t *testing.T) {
	service, mockNodeRepo, mockDomainRepo := CreateTestNodeService(t)

	ctx := CreateTestContext()
	testDomain := CreateTestDomain("test", "Test")
	mockDomainRepo.Create(ctx, testDomain)

	testNode := CreateTestNode(testDomain.ID, "https://example.com", "Test", "Test node")
	mockNodeRepo.Create(ctx, testNode)

	err := service.DeleteNode(ctx, testNode.ID)

	assert.NoError(t, err)

	// Verify it's deleted
	_, err = service.GetNode(ctx, testNode.ID)
	assert.Error(t, err)
}