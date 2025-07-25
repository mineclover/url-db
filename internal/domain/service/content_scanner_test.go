package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"url-db/internal/application/dto/response"
	"url-db/internal/constants"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
	"url-db/internal/domain/service"
)

// Mock repositories for testing
type mockNodeRepository struct {
	nodes []*entity.Node
}

func (m *mockNodeRepository) CountByDomain(ctx context.Context, domainID int) (int, error) {
	count := 0
	for _, node := range m.nodes {
		if node.DomainID() == domainID {
			count++
		}
	}
	return count, nil
}

func (m *mockNodeRepository) GetByDomainFromCursor(ctx context.Context, domainID int, lastNodeID int, limit int) ([]*entity.Node, error) {
	var result []*entity.Node
	for _, node := range m.nodes {
		if node.DomainID() == domainID && node.ID() > lastNodeID {
			result = append(result, node)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

// Implement other required methods (stub implementations)
func (m *mockNodeRepository) Create(ctx context.Context, node *entity.Node) error { return nil }
func (m *mockNodeRepository) GetByID(ctx context.Context, id int) (*entity.Node, error) { return nil, nil }
func (m *mockNodeRepository) GetByURL(ctx context.Context, url, domainName string) (*entity.Node, error) { return nil, nil }
func (m *mockNodeRepository) List(ctx context.Context, domainName string, page, size int) ([]*entity.Node, int, error) { return nil, 0, nil }
func (m *mockNodeRepository) Update(ctx context.Context, node *entity.Node) error { return nil }
func (m *mockNodeRepository) Delete(ctx context.Context, id int) error { return nil }
func (m *mockNodeRepository) Exists(ctx context.Context, url, domainName string) (bool, error) { return false, nil }
func (m *mockNodeRepository) GetBatch(ctx context.Context, ids []int) ([]*entity.Node, error) { return nil, nil }
func (m *mockNodeRepository) GetDomainByNodeID(ctx context.Context, nodeID int) (*entity.Domain, error) { return nil, nil }
func (m *mockNodeRepository) FilterByAttributes(ctx context.Context, domainName string, filters []repository.AttributeFilter, page, size int) ([]*entity.Node, int, error) { return nil, 0, nil }

type mockNodeAttributeRepository struct {
	attributes map[int][]*entity.NodeAttribute
}

func (m *mockNodeAttributeRepository) GetByNodeID(ctx context.Context, nodeID int) ([]*entity.NodeAttribute, error) {
	return m.attributes[nodeID], nil
}

// Implement other required methods (stub implementations)
func (m *mockNodeAttributeRepository) Create(ctx context.Context, nodeAttribute *entity.NodeAttribute) error { return nil }
func (m *mockNodeAttributeRepository) GetByNodeAndAttribute(ctx context.Context, nodeID int, attributeID int) (*entity.NodeAttribute, error) { return nil, nil }
func (m *mockNodeAttributeRepository) Update(ctx context.Context, nodeAttribute *entity.NodeAttribute) error { return nil }
func (m *mockNodeAttributeRepository) Delete(ctx context.Context, nodeID int, attributeID int) error { return nil }
func (m *mockNodeAttributeRepository) DeleteAllByNode(ctx context.Context, nodeID int) error { return nil }
func (m *mockNodeAttributeRepository) SetNodeAttributes(ctx context.Context, nodeID int, attributes []*entity.NodeAttribute) error { return nil }
func (m *mockNodeAttributeRepository) GetNodesWithAttribute(ctx context.Context, attributeID int, value *string) ([]int, error) { return nil, nil }

type mockDomainRepository struct {
	domain *entity.Domain
}

func (m *mockDomainRepository) GetByName(ctx context.Context, name string) (*entity.Domain, error) {
	return m.domain, nil
}

// Implement other required methods (stub implementations)
func (m *mockDomainRepository) Create(ctx context.Context, domain *entity.Domain) error { return nil }
func (m *mockDomainRepository) GetByID(ctx context.Context, id int) (*entity.Domain, error) { return nil, nil }
func (m *mockDomainRepository) List(ctx context.Context, page, size int) ([]*entity.Domain, int, error) { return nil, 0, nil }
func (m *mockDomainRepository) Update(ctx context.Context, domain *entity.Domain) error { return nil }
func (m *mockDomainRepository) Delete(ctx context.Context, name string) error { return nil }
func (m *mockDomainRepository) Exists(ctx context.Context, name string) (bool, error) { return false, nil }

func TestContentScanner_ScanAllContent(t *testing.T) {
	// Create test domain
	domain, _ := entity.NewDomain("test", "Test domain")
	domain.SetID(1)

	// Create test nodes
	node1, _ := entity.NewNode("https://example.com/1", "Title 1", "Description 1", 1)
	node1.SetID(1)
	node1.SetTimestamps(time.Now(), time.Now())

	node2, _ := entity.NewNode("https://example.com/2", "Title 2", "Description 2", 1)
	node2.SetID(2)
	node2.SetTimestamps(time.Now(), time.Now())

	// Create mock repositories
	nodeRepo := &mockNodeRepository{
		nodes: []*entity.Node{node1, node2},
	}

	nodeAttrRepo := &mockNodeAttributeRepository{
		attributes: make(map[int][]*entity.NodeAttribute),
	}

	domainRepo := &mockDomainRepository{
		domain: domain,
	}

	// Create content scanner
	scanner := service.NewContentScanner(nodeRepo, nodeAttrRepo, domainRepo)

	// Test scan request (first page)
	req := service.ScanRequest{
		DomainName:        "test",
		MaxTokensPerPage:  constants.DefaultMaxTokensPerPage,
		Page:              1, // Test page-based navigation
		IncludeAttributes: false,
	}

	// Execute scan
	result, err := scanner.ScanAllContent(context.Background(), req)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result.Items))
	}

	if result.Metadata.TotalNodes != 2 {
		t.Errorf("Expected total nodes 2, got %d", result.Metadata.TotalNodes)
	}

	if result.Pagination.HasMore {
		t.Error("Expected no more pages for small dataset")
	}

	if result.Pagination.CurrentPage != 1 {
		t.Errorf("Expected current page 1, got %d", result.Pagination.CurrentPage)
	}

	if result.Pagination.HasPrevious {
		t.Error("Expected no previous page for first page")
	}

	// Check first item
	firstItem := result.Items[0]
	if firstItem.Content != "https://example.com/1" {
		t.Errorf("Expected first item content to be 'https://example.com/1', got '%s'", firstItem.Content)
	}

	if firstItem.Title == nil || *firstItem.Title != "Title 1" {
		t.Error("Expected first item to have correct title")
	}
}

func TestSmartChunker_EstimateNodeTokens(t *testing.T) {
	chunker := service.NewSmartChunker(8000, false)

	// Create test node
	node := response.NodeWithAttributes{
		ID:      1,
		Content: "https://example.com/test-url",
		Title:   stringPtr("Test Title"),
		Description: stringPtr("This is a test description"),
	}

	tokens := chunker.EstimateNodeTokens(node)

	// Should be at least minimum tokens
	if tokens < constants.MinTokensPerNode {
		t.Errorf("Expected at least %d tokens, got %d", constants.MinTokensPerNode, tokens)
	}

	// Should be reasonable estimate (not too high)
	if tokens > 1000 {
		t.Errorf("Token estimate seems too high: %d", tokens)
	}
}

func TestSmartChunker_CanAddNode(t *testing.T) {
	chunker := service.NewSmartChunker(200, false) // Small limit for testing

	// Create a large node that should not fit
	largeNode := response.NodeWithAttributes{
		ID:      1,
		Content: "https://example.com/" + strings.Repeat("a", 500), // Large URL
		Title:   stringPtr("Large Title " + strings.Repeat("b", 200)),
	}

	// Should not be able to add large node to empty chunker with small limit
	canAdd := chunker.CanAddNode(largeNode)
	if canAdd {
		t.Error("Expected large node to not fit in small chunker")
	}

	// Create small chunker with higher limit
	bigChunker := service.NewSmartChunker(8000, false)
	canAddToBig := bigChunker.CanAddNode(largeNode)
	if !canAddToBig {
		t.Error("Expected large node to fit in big chunker")
	}
}

func TestContentScanner_ScanAllContent_WithCompression(t *testing.T) {
	// Create test domain
	domain, _ := entity.NewDomain("test", "Test domain")
	domain.SetID(1)

	// Create test nodes
	node1, _ := entity.NewNode("https://example.com/1", "Title 1", "Description 1", 1)
	node1.SetID(1)
	node1.SetTimestamps(time.Now(), time.Now())

	node2, _ := entity.NewNode("https://example.com/2", "Title 2", "Description 2", 1)
	node2.SetID(2)
	node2.SetTimestamps(time.Now(), time.Now())

	// Create mock attributes with duplicates for testing compression
	attr1, _ := entity.NewNodeAttribute(1, 1, "tech", nil)
	attr1.SetName("category")
	attr1.SetAttributeType(stringPtr("tag"))
	
	attr2, _ := entity.NewNodeAttribute(2, 1, "tech", nil) // duplicate value
	attr2.SetName("category")
	attr2.SetAttributeType(stringPtr("tag"))
	
	attr3, _ := entity.NewNodeAttribute(1, 2, "high", nil)
	attr3.SetName("priority")
	attr3.SetAttributeType(stringPtr("tag"))

	nodeAttrRepo := &mockNodeAttributeRepository{
		attributes: map[int][]*entity.NodeAttribute{
			1: {attr1, attr3},
			2: {attr2},
		},
	}

	nodeRepo := &mockNodeRepository{
		nodes: []*entity.Node{node1, node2},
	}

	domainRepo := &mockDomainRepository{
		domain: domain,
	}

	scanner := service.NewContentScanner(nodeRepo, nodeAttrRepo, domainRepo)

	// Test with compression enabled
	req := service.ScanRequest{
		DomainName:         "test",
		MaxTokensPerPage:   constants.DefaultMaxTokensPerPage,
		Page:               1,
		IncludeAttributes:  true,
		CompressAttributes: true, // Enable compression
	}

	result, err := scanner.ScanAllContent(context.Background(), req)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.Metadata.CompressedOutput {
		t.Error("Expected compressed output flag to be true")
	}

	if result.Metadata.AttributeSummary == nil {
		t.Error("Expected attribute summary with compression enabled")
	}

	// Check that attributes are included
	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result.Items))
	}

	// Verify first item has attributes
	firstItem := result.Items[0]
	if firstItem.Attributes == nil || len(firstItem.Attributes) == 0 {
		t.Error("Expected first item to have attributes")
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}