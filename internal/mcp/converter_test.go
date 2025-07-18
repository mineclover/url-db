package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"url-db/internal/models"
)

// Mock CompositeKeyAdapter for testing
type MockCompositeKeyAdapter struct {
	mock.Mock
}

func (m *MockCompositeKeyAdapter) NodeToCompositeID(domain *models.Domain, node *models.Node) string {
	args := m.Called(domain, node)
	return args.String(0)
}

func (m *MockCompositeKeyAdapter) ParseCompositeID(compositeID string) (string, int, error) {
	args := m.Called(compositeID)
	return args.String(0), args.Int(1), args.Error(2)
}

func TestNewConverter(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	assert.NotNil(t, converter)
	assert.Equal(t, mockAdapter, converter.compositeKeyAdapter)
}

func TestConverter_DomainToMCP(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	domain := &models.Domain{
		ID:          1,
		Name:        "example.com",
		Description: "Test domain",
		CreatedAt:   "2024-01-01T00:00:00Z",
		UpdatedAt:   "2024-01-01T00:00:00Z",
	}

	nodeCount := 5

	result := converter.DomainToMCP(domain, nodeCount)

	assert.NotNil(t, result)
	assert.Equal(t, domain.Name, result.Name)
	assert.Equal(t, domain.Description, result.Description)
	assert.Equal(t, nodeCount, result.NodeCount)
	assert.Equal(t, domain.CreatedAt, result.CreatedAt)
	assert.Equal(t, domain.UpdatedAt, result.UpdatedAt)
}

func TestConverter_NodeToMCP(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	domain := &models.Domain{
		ID:          1,
		Name:        "example.com",
		Description: "Test domain",
	}

	node := &models.Node{
		ID:          1,
		DomainID:    1,
		URL:         "https://example.com/page",
		Title:       "Test Page",
		Description: "Test description",
		CreatedAt:   "2024-01-01T01:00:00Z",
		UpdatedAt:   "2024-01-01T01:00:00Z",
	}

	expectedCompositeID := "example.com::https://example.com/page"
	mockAdapter.On("NodeToCompositeID", domain, node).Return(expectedCompositeID)

	result := converter.NodeToMCP(domain, node)

	assert.NotNil(t, result)
	assert.Equal(t, expectedCompositeID, result.CompositeID)
	assert.Equal(t, domain.Name, result.DomainName)
	assert.Equal(t, node.URL, result.URL)
	assert.Equal(t, node.Title, result.Title)
	assert.Equal(t, node.Description, result.Description)
	assert.Equal(t, node.CreatedAt, result.CreatedAt)
	assert.Equal(t, node.UpdatedAt, result.UpdatedAt)

	mockAdapter.AssertExpectations(t)
}

func TestConverter_NodeAttributesToMCP(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	compositeID := "example.com::https://example.com/page"
	attributes := []models.NodeAttributeWithInfo{
		{
			ID:          1,
			NodeID:      1,
			AttributeID: 1,
			Name:        "category",
			Type:        models.AttributeTypeTag,
			Value:       "tutorial",
			OrderIndex:  1,
		},
		{
			ID:          2,
			NodeID:      1,
			AttributeID: 2,
			Name:        "rating",
			Type:        models.AttributeTypeNumber,
			Value:       "5",
			OrderIndex:  1,
		},
	}

	result := converter.NodeAttributesToMCP(compositeID, attributes)

	assert.NotNil(t, result)
	assert.Equal(t, compositeID, result.CompositeID)
	assert.Len(t, result.Attributes, 2)

	// Check first attribute
	assert.Equal(t, "category", result.Attributes[0].Name)
	assert.Equal(t, "tag", result.Attributes[0].Type)
	assert.Equal(t, "tutorial", result.Attributes[0].Value)

	// Check second attribute
	assert.Equal(t, "rating", result.Attributes[1].Name)
	assert.Equal(t, "number", result.Attributes[1].Type)
	assert.Equal(t, "5", result.Attributes[1].Value)
}

func TestConverter_NodeAttributesToMCP_EmptyAttributes(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	compositeID := "example.com::https://example.com/page"
	attributes := []models.NodeAttributeWithInfo{}

	result := converter.NodeAttributesToMCP(compositeID, attributes)

	assert.NotNil(t, result)
	assert.Equal(t, compositeID, result.CompositeID)
	assert.Len(t, result.Attributes, 0)
	assert.NotNil(t, result.Attributes) // Should be empty slice, not nil
}

func TestConverter_AttributeTypesToString(t *testing.T) {
	tests := []struct {
		name         string
		attributeType models.AttributeType
		expected     string
	}{
		{"Tag type", models.AttributeTypeTag, "tag"},
		{"Ordered tag type", models.AttributeTypeOrderedTag, "ordered_tag"},
		{"Number type", models.AttributeTypeNumber, "number"},
		{"String type", models.AttributeTypeString, "string"},
		{"Markdown type", models.AttributeTypeMarkdown, "markdown"},
		{"Image type", models.AttributeTypeImage, "image"},
	}

	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test attribute with the specific type
			attributes := []models.NodeAttributeWithInfo{
				{
					Name:  "test",
					Type:  tt.attributeType,
					Value: "test-value",
				},
			}

			result := converter.NodeAttributesToMCP("test::test", attributes)
			assert.Equal(t, tt.expected, result.Attributes[0].Type)
		})
	}
}

func TestConverter_BatchNodeResponse(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	domain := &models.Domain{
		ID:   1,
		Name: "example.com",
	}

	nodes := []models.Node{
		{
			ID:          1,
			DomainID:    1,
			URL:         "https://example.com/page1",
			Title:       "Page 1",
			Description: "First page",
		},
		{
			ID:          2,
			DomainID:    1,
			URL:         "https://example.com/page2",
			Title:       "Page 2",
			Description: "Second page",
		},
	}

	notFound := []string{"example.com::https://example.com/missing"}
	errors := []string{"Error processing composite ID"}

	mockAdapter.On("NodeToCompositeID", domain, &nodes[0]).Return("example.com::https://example.com/page1")
	mockAdapter.On("NodeToCompositeID", domain, &nodes[1]).Return("example.com::https://example.com/page2")

	result := converter.BatchNodeResponseToMCP(domain, nodes, notFound, errors)

	assert.NotNil(t, result)
	assert.Len(t, result.Nodes, 2)
	assert.Equal(t, notFound, result.NotFound)
	assert.Equal(t, errors, result.Errors)

	// Check converted nodes
	assert.Equal(t, "example.com::https://example.com/page1", result.Nodes[0].CompositeID)
	assert.Equal(t, "example.com::https://example.com/page2", result.Nodes[1].CompositeID)

	mockAdapter.AssertExpectations(t)
}

func TestConverter_DomainListToMCP(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	domains := []models.Domain{
		{
			ID:          1,
			Name:        "example.com",
			Description: "Example domain",
			CreatedAt:   "2024-01-01T00:00:00Z",
			UpdatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          2,
			Name:        "test.org",
			Description: "Test domain",
			CreatedAt:   "2024-01-02T00:00:00Z",
			UpdatedAt:   "2024-01-02T00:00:00Z",
		},
	}

	nodeCounts := []int{5, 3}

	result := converter.DomainListToMCP(domains, nodeCounts)

	assert.NotNil(t, result)
	assert.Len(t, result.Domains, 2)

	// Check first domain
	assert.Equal(t, "example.com", result.Domains[0].Name)
	assert.Equal(t, "Example domain", result.Domains[0].Description)
	assert.Equal(t, 5, result.Domains[0].NodeCount)

	// Check second domain
	assert.Equal(t, "test.org", result.Domains[1].Name)
	assert.Equal(t, "Test domain", result.Domains[1].Description)
	assert.Equal(t, 3, result.Domains[1].NodeCount)
}

func TestConverter_DomainListToMCP_MismatchedCounts(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	domains := []models.Domain{
		{ID: 1, Name: "example.com", Description: "Example domain"},
		{ID: 2, Name: "test.org", Description: "Test domain"},
	}

	// Provide fewer node counts than domains
	nodeCounts := []int{5}

	result := converter.DomainListToMCP(domains, nodeCounts)

	assert.NotNil(t, result)
	assert.Len(t, result.Domains, 2)

	// First domain should have the provided count
	assert.Equal(t, 5, result.Domains[0].NodeCount)

	// Second domain should have 0 count (default)
	assert.Equal(t, 0, result.Domains[1].NodeCount)
}

func TestConverter_NodeListToMCP(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	domain := &models.Domain{
		ID:   1,
		Name: "example.com",
	}

	response := &models.NodeListResponse{
		Nodes: []models.Node{
			{
				ID:          1,
				DomainID:    1,
				URL:         "https://example.com/page1",
				Title:       "Page 1",
				Description: "First page",
			},
			{
				ID:          2,
				DomainID:    1,
				URL:         "https://example.com/page2",
				Title:       "Page 2",
				Description: "Second page",
			},
		},
		Page:       1,
		Size:       20,
		TotalCount: 2,
		TotalPages: 1,
	}

	mockAdapter.On("NodeToCompositeID", domain, &response.Nodes[0]).Return("example.com::https://example.com/page1")
	mockAdapter.On("NodeToCompositeID", domain, &response.Nodes[1]).Return("example.com::https://example.com/page2")

	result := converter.NodeListToMCP(domain, response)

	assert.NotNil(t, result)
	assert.Len(t, result.Nodes, 2)
	assert.Equal(t, response.Page, result.Page)
	assert.Equal(t, response.Size, result.Size)
	assert.Equal(t, response.TotalCount, result.TotalCount)
	assert.Equal(t, response.TotalPages, result.TotalPages)

	// Check converted nodes
	assert.Equal(t, "example.com::https://example.com/page1", result.Nodes[0].CompositeID)
	assert.Equal(t, "example.com::https://example.com/page2", result.Nodes[1].CompositeID)

	mockAdapter.AssertExpectations(t)
}

func TestConverter_EdgeCases(t *testing.T) {
	mockAdapter := &MockCompositeKeyAdapter{}
	converter := NewConverter(mockAdapter)

	t.Run("Empty node list", func(t *testing.T) {
		domain := &models.Domain{ID: 1, Name: "example.com"}
		response := &models.NodeListResponse{
			Nodes:      []models.Node{},
			Page:       1,
			Size:       20,
			TotalCount: 0,
			TotalPages: 0,
		}

		result := converter.NodeListToMCP(domain, response)

		assert.NotNil(t, result)
		assert.Len(t, result.Nodes, 0)
		assert.Equal(t, 0, result.TotalCount)
	})

	t.Run("Nil domain", func(t *testing.T) {
		node := &models.Node{
			ID:  1,
			URL: "https://example.com/page",
		}

		// Should not panic even with nil domain
		result := converter.NodeToMCP(nil, node)

		assert.NotNil(t, result)
		assert.Equal(t, "", result.DomainName) // Should handle nil domain gracefully
	})

	t.Run("Empty attribute list", func(t *testing.T) {
		result := converter.NodeAttributesToMCP("test::test", []models.NodeAttributeWithInfo{})

		assert.NotNil(t, result)
		assert.Equal(t, "test::test", result.CompositeID)
		assert.NotNil(t, result.Attributes)
		assert.Len(t, result.Attributes, 0)
	})
}