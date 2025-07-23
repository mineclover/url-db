package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"url-db/internal/models"
)

func TestDomainStruct(t *testing.T) {
	now := time.Now()
	domain := models.Domain{
		ID:          1,
		Name:        "test-domain",
		Description: "Test description",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, 1, domain.ID)
	assert.Equal(t, "test-domain", domain.Name)
	assert.Equal(t, "Test description", domain.Description)
	assert.Equal(t, now, domain.CreatedAt)
	assert.Equal(t, now, domain.UpdatedAt)
}

func TestCreateDomainRequest(t *testing.T) {
	request := models.CreateDomainRequest{
		Name:        "test-domain",
		Description: "Test description",
	}

	assert.Equal(t, "test-domain", request.Name)
	assert.Equal(t, "Test description", request.Description)
}

func TestNodeStruct(t *testing.T) {
	now := time.Now()
	node := models.Node{
		ID:          1,
		Content:     "https://example.com",
		Title:       "Example",
		Description: "Test node",
		DomainID:    1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, 1, node.ID)
	assert.Equal(t, "https://example.com", node.Content)
	assert.Equal(t, "Example", node.Title)
	assert.Equal(t, "Test node", node.Description)
	assert.Equal(t, 1, node.DomainID)
	assert.Equal(t, now, node.CreatedAt)
	assert.Equal(t, now, node.UpdatedAt)
}

func TestAttributeStruct(t *testing.T) {
	now := time.Now()
	attribute := models.Attribute{
		ID:          1,
		Name:        "test-attr",
		Type:        models.AttributeTypeString,
		Description: "Test attribute",
		DomainID:    1,
		CreatedAt:   now,
	}

	assert.Equal(t, 1, attribute.ID)
	assert.Equal(t, "test-attr", attribute.Name)
	assert.Equal(t, models.AttributeTypeString, attribute.Type)
	assert.Equal(t, "Test attribute", attribute.Description)
	assert.Equal(t, 1, attribute.DomainID)
	assert.Equal(t, now, attribute.CreatedAt)
}

func TestAttributeTypeConstants(t *testing.T) {
	// Test that all attribute type constants are defined
	assert.Equal(t, models.AttributeType("tag"), models.AttributeTypeTag)
	assert.Equal(t, models.AttributeType("ordered_tag"), models.AttributeTypeOrderedTag)
	assert.Equal(t, models.AttributeType("number"), models.AttributeTypeNumber)
	assert.Equal(t, models.AttributeType("string"), models.AttributeTypeString)
	assert.Equal(t, models.AttributeType("markdown"), models.AttributeTypeMarkdown)
	assert.Equal(t, models.AttributeType("image"), models.AttributeTypeImage)
}

func TestCompositeKeyStruct(t *testing.T) {
	key := &models.CompositeKey{
		ToolName:   "url-db",
		DomainName: "test",
		ID:         123,
	}

	assert.Equal(t, "url-db", key.ToolName)
	assert.Equal(t, "test", key.DomainName)
	assert.Equal(t, 123, key.ID)
}

// Server model tests
func TestMCPServerInfo(t *testing.T) {
	serverInfo := models.MCPServerInfo{
		Name:               "url-db",
		Version:            "1.0.0",
		Description:        "URL Database Server",
		Capabilities:       []string{"nodes", "domains"},
		CompositeKeyFormat: "tool-name:domain:id",
	}

	assert.Equal(t, "url-db", serverInfo.Name)
	assert.Equal(t, "1.0.0", serverInfo.Version)
	assert.Equal(t, "URL Database Server", serverInfo.Description)
	assert.Len(t, serverInfo.Capabilities, 2)
	assert.Contains(t, serverInfo.Capabilities, "nodes")
	assert.Contains(t, serverInfo.Capabilities, "domains")
	assert.Equal(t, "tool-name:domain:id", serverInfo.CompositeKeyFormat)
}

func TestMCPAttributeRequest(t *testing.T) {
	req := models.MCPAttributeRequest{
		Name:  "category",
		Value: "tech",
	}

	assert.Equal(t, "category", req.Name)
	assert.Equal(t, "tech", req.Value)
}