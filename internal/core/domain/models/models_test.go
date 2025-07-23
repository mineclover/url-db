package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"url-db/internal/core/domain/models"
)

func TestCompositeKey(t *testing.T) {
	key := &models.CompositeKey{
		ToolName:   "url-db",
		DomainName: "test-domain",
		ID:         123,
	}

	assert.Equal(t, "url-db", key.ToolName)
	assert.Equal(t, "test-domain", key.DomainName)
	assert.Equal(t, 123, key.ID)
}

func TestDomain(t *testing.T) {
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
		Name:        "new-domain",
		Description: "New domain description",
	}

	assert.Equal(t, "new-domain", request.Name)
	assert.Equal(t, "New domain description", request.Description)
}

func TestNode(t *testing.T) {
	now := time.Now()
	node := models.Node{
		ID:          1,
		Content:     "https://example.com",
		Title:       "Example Site",
		Description: "Example description",
		DomainID:    1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, 1, node.ID)
	assert.Equal(t, "https://example.com", node.Content)
	assert.Equal(t, "Example Site", node.Title)
	assert.Equal(t, "Example description", node.Description)
	assert.Equal(t, 1, node.DomainID)
	assert.Equal(t, now, node.CreatedAt)
	assert.Equal(t, now, node.UpdatedAt)
}

func TestAttribute(t *testing.T) {
	now := time.Now()
	attribute := models.Attribute{
		ID:          1,
		Name:        "category",
		Type:        models.AttributeTypeString,
		Description: "Category attribute",
		DomainID:    1,
		CreatedAt:   now,
	}

	assert.Equal(t, 1, attribute.ID)
	assert.Equal(t, "category", attribute.Name)
	assert.Equal(t, models.AttributeTypeString, attribute.Type)
	assert.Equal(t, "Category attribute", attribute.Description)
	assert.Equal(t, 1, attribute.DomainID)
	assert.Equal(t, now, attribute.CreatedAt)
}

func TestAttributeTypeConstants(t *testing.T) {
	assert.Equal(t, models.AttributeType("tag"), models.AttributeTypeTag)
	assert.Equal(t, models.AttributeType("ordered_tag"), models.AttributeTypeOrderedTag)
	assert.Equal(t, models.AttributeType("number"), models.AttributeTypeNumber)
	assert.Equal(t, models.AttributeType("string"), models.AttributeTypeString)
	assert.Equal(t, models.AttributeType("markdown"), models.AttributeTypeMarkdown)
	assert.Equal(t, models.AttributeType("image"), models.AttributeTypeImage)
}

func TestServerInfo(t *testing.T) {
	serverInfo := models.MCPServerInfo{
		Name:               "test-server",
		Version:            "1.0.0",
		Description:        "Test server",
		Capabilities:       []string{"nodes", "domains"},
		CompositeKeyFormat: "tool:domain:id",
	}

	assert.Equal(t, "test-server", serverInfo.Name)
	assert.Equal(t, "1.0.0", serverInfo.Version)
	assert.Equal(t, "Test server", serverInfo.Description)
	assert.Len(t, serverInfo.Capabilities, 2)
	assert.Contains(t, serverInfo.Capabilities, "nodes")
	assert.Contains(t, serverInfo.Capabilities, "domains")
	assert.Equal(t, "tool:domain:id", serverInfo.CompositeKeyFormat)
}

func TestSubscriptionModels(t *testing.T) {
	now := time.Now()
	endpoint := "https://webhook.example.com"
	
	subscription := models.NodeSubscription{
		ID:                 1,
		SubscriberService:  "test-service",
		SubscriberEndpoint: &endpoint,
		SubscribedNodeID:   100,
		EventTypes:         models.EventTypeList{"created", "updated"},
		IsActive:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	assert.Equal(t, int64(1), subscription.ID)
	assert.Equal(t, "test-service", subscription.SubscriberService)
	assert.NotNil(t, subscription.SubscriberEndpoint)
	assert.Equal(t, "https://webhook.example.com", *subscription.SubscriberEndpoint)
	assert.Equal(t, int64(100), subscription.SubscribedNodeID)
	assert.Len(t, subscription.EventTypes, 2)
	assert.True(t, subscription.IsActive)
}

func TestDependencyModels(t *testing.T) {
	now := time.Now()
	metadata := &models.DependencyMetadata{
		Relationship: "parent-child",
		Description:  "Child depends on parent",
	}
	
	dependency := models.NodeDependency{
		ID:               1,
		DependentNodeID:  100,
		DependencyNodeID: 200,
		DependencyType:   models.DependencyTypeHard,
		CascadeDelete:    true,
		CascadeUpdate:    false,
		Metadata:         metadata,
		CreatedAt:        now,
	}

	assert.Equal(t, int64(1), dependency.ID)
	assert.Equal(t, int64(100), dependency.DependentNodeID)
	assert.Equal(t, int64(200), dependency.DependencyNodeID)
	assert.Equal(t, models.DependencyTypeHard, dependency.DependencyType)
	assert.True(t, dependency.CascadeDelete)
	assert.False(t, dependency.CascadeUpdate)
	assert.NotNil(t, dependency.Metadata)
	assert.Equal(t, "parent-child", dependency.Metadata.Relationship)
}

func TestNodeConnection(t *testing.T) {
	now := time.Now()
	
	connection := models.NodeConnection{
		ID:               1,
		SourceNodeID:     100,
		TargetNodeID:     200,
		RelationshipType: models.RelationshipTypeChild,
		Description:      "Parent-child connection",
		CreatedAt:        now,
	}

	assert.Equal(t, 1, connection.ID)
	assert.Equal(t, 100, connection.SourceNodeID)
	assert.Equal(t, 200, connection.TargetNodeID)
	assert.Equal(t, models.RelationshipTypeChild, connection.RelationshipType)
	assert.Equal(t, "Parent-child connection", connection.Description)
	assert.Equal(t, now, connection.CreatedAt)
}

func TestEventTypeListScanValue(t *testing.T) {
	// Test Value method
	eventTypes := models.EventTypeList{"created", "updated"}
	value, err := eventTypes.Value()
	assert.NoError(t, err)
	assert.Equal(t, `["created","updated"]`, value)

	// Test empty list
	emptyTypes := models.EventTypeList{}
	value, err = emptyTypes.Value()
	assert.NoError(t, err)
	assert.Equal(t, "[]", value)

	// Test Scan method
	var scannedTypes models.EventTypeList
	err = scannedTypes.Scan(`["deleted","attribute_changed"]`)
	assert.NoError(t, err)
	assert.Len(t, scannedTypes, 2)
	assert.Contains(t, scannedTypes, "deleted")
	assert.Contains(t, scannedTypes, "attribute_changed")

	// Test Scan with nil
	err = scannedTypes.Scan(nil)
	assert.NoError(t, err)
	assert.Len(t, scannedTypes, 0)
}

func TestDependencyMetadataScanValue(t *testing.T) {
	// Test Value method
	metadata := &models.DependencyMetadata{
		Relationship: "test-rel",
		Description:  "test desc",
	}
	
	value, err := metadata.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	// Test nil Value
	var nilMetadata *models.DependencyMetadata
	value, err = nilMetadata.Value()
	assert.NoError(t, err)
	assert.Nil(t, value)

	// Test Scan method
	var scannedMetadata models.DependencyMetadata
	jsonData := `{"relationship":"parent","description":"parent dependency"}`
	err = scannedMetadata.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, "parent", scannedMetadata.Relationship)
	assert.Equal(t, "parent dependency", scannedMetadata.Description)

	// Test Scan with nil
	err = scannedMetadata.Scan(nil)
	assert.NoError(t, err)
}