package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"url-db/internal/models"
)

func TestNodeConnection(t *testing.T) {
	now := time.Now()
	
	connection := models.NodeConnection{
		ID:               1,
		SourceNodeID:     100,
		TargetNodeID:     200,
		RelationshipType: models.RelationshipTypeChild,
		Description:      "Parent-child relationship",
		CreatedAt:        now,
	}

	assert.Equal(t, 1, connection.ID)
	assert.Equal(t, 100, connection.SourceNodeID)
	assert.Equal(t, 200, connection.TargetNodeID)
	assert.Equal(t, models.RelationshipTypeChild, connection.RelationshipType)
	assert.Equal(t, "Parent-child relationship", connection.Description)
	assert.Equal(t, now, connection.CreatedAt)
}

func TestCreateNodeConnectionRequest(t *testing.T) {
	req := models.CreateNodeConnectionRequest{
		SourceNodeID:     100,
		TargetNodeID:     200,
		RelationshipType: models.RelationshipTypeLinked,
		Description:      "Linked nodes for reference",
	}

	assert.Equal(t, 100, req.SourceNodeID)
	assert.Equal(t, 200, req.TargetNodeID)
	assert.Equal(t, models.RelationshipTypeLinked, req.RelationshipType)
	assert.Equal(t, "Linked nodes for reference", req.Description)
}

func TestUpdateNodeConnectionRequest(t *testing.T) {
	req := models.UpdateNodeConnectionRequest{
		RelationshipType: models.RelationshipTypeRelated,
		Description:      "Updated relationship description",
	}

	assert.Equal(t, models.RelationshipTypeRelated, req.RelationshipType)
	assert.Equal(t, "Updated relationship description", req.Description)
}

func TestNodeConnectionListResponse(t *testing.T) {
	connections := []models.NodeConnection{
		{ID: 1, SourceNodeID: 100, TargetNodeID: 200, RelationshipType: models.RelationshipTypeChild},
		{ID: 2, SourceNodeID: 200, TargetNodeID: 300, RelationshipType: models.RelationshipTypeNext},
	}
	
	response := models.NodeConnectionListResponse{
		Connections: connections,
		TotalCount:  2,
		Page:        1,
		Size:        10,
		TotalPages:  1,
	}

	assert.Len(t, response.Connections, 2)
	assert.Equal(t, 2, response.TotalCount)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.Size)
	assert.Equal(t, 1, response.TotalPages)
	assert.Equal(t, 1, response.Connections[0].ID)
	assert.Equal(t, models.RelationshipTypeChild, response.Connections[0].RelationshipType)
	assert.Equal(t, 2, response.Connections[1].ID)
	assert.Equal(t, models.RelationshipTypeNext, response.Connections[1].RelationshipType)
}

func TestNodeConnectionWithInfo(t *testing.T) {
	now := time.Now()
	
	connectionWithInfo := models.NodeConnectionWithInfo{
		NodeConnection: models.NodeConnection{
			ID:               1,
			SourceNodeID:     100,
			TargetNodeID:     200,
			RelationshipType: models.RelationshipTypeParent,
			Description:      "Parent relationship",
			CreatedAt:        now,
		},
		SourceNodeURL:   "https://example.com/source",
		TargetNodeURL:   "https://example.com/target",
		SourceNodeTitle: "Source Node Title",
		TargetNodeTitle: "Target Node Title",
	}

	// Test embedded NodeConnection fields
	assert.Equal(t, 1, connectionWithInfo.ID)
	assert.Equal(t, 100, connectionWithInfo.SourceNodeID)
	assert.Equal(t, 200, connectionWithInfo.TargetNodeID)
	assert.Equal(t, models.RelationshipTypeParent, connectionWithInfo.RelationshipType)
	assert.Equal(t, "Parent relationship", connectionWithInfo.Description)
	assert.Equal(t, now, connectionWithInfo.CreatedAt)

	// Test additional info fields
	assert.Equal(t, "https://example.com/source", connectionWithInfo.SourceNodeURL)
	assert.Equal(t, "https://example.com/target", connectionWithInfo.TargetNodeURL)
	assert.Equal(t, "Source Node Title", connectionWithInfo.SourceNodeTitle)
	assert.Equal(t, "Target Node Title", connectionWithInfo.TargetNodeTitle)
}

func TestRelationshipTypeConstants(t *testing.T) {
	// Test all relationship type constants
	assert.Equal(t, "related", models.RelationshipTypeRelated)
	assert.Equal(t, "child", models.RelationshipTypeChild)
	assert.Equal(t, "parent", models.RelationshipTypeParent)
	assert.Equal(t, "next", models.RelationshipTypeNext)
	assert.Equal(t, "previous", models.RelationshipTypePrevious)
	assert.Equal(t, "linked", models.RelationshipTypeLinked)
	assert.Equal(t, "custom", models.RelationshipTypeCustom)
}

func TestNodeConnectionComplexScenarios(t *testing.T) {
	now := time.Now()
	
	// Test bidirectional relationship
	parentToChild := models.NodeConnection{
		ID:               1,
		SourceNodeID:     100,
		TargetNodeID:     200,
		RelationshipType: models.RelationshipTypeChild,
		Description:      "100 has child 200",
		CreatedAt:        now,
	}
	
	childToParent := models.NodeConnection{
		ID:               2,
		SourceNodeID:     200,
		TargetNodeID:     100,
		RelationshipType: models.RelationshipTypeParent,
		Description:      "200 has parent 100",
		CreatedAt:        now,
	}

	assert.Equal(t, parentToChild.SourceNodeID, childToParent.TargetNodeID)
	assert.Equal(t, parentToChild.TargetNodeID, childToParent.SourceNodeID)
	assert.Equal(t, models.RelationshipTypeChild, parentToChild.RelationshipType)
	assert.Equal(t, models.RelationshipTypeParent, childToParent.RelationshipType)

	// Test sequence relationship
	connections := []models.NodeConnection{
		{
			ID:               3,
			SourceNodeID:     100,
			TargetNodeID:     200,
			RelationshipType: models.RelationshipTypeNext,
			Description:      "Step 1 to Step 2",
		},
		{
			ID:               4,
			SourceNodeID:     200,
			TargetNodeID:     300,
			RelationshipType: models.RelationshipTypeNext,
			Description:      "Step 2 to Step 3",
		},
	}

	// Verify sequence chain
	assert.Equal(t, connections[0].TargetNodeID, connections[1].SourceNodeID)
	assert.Equal(t, models.RelationshipTypeNext, connections[0].RelationshipType)
	assert.Equal(t, models.RelationshipTypeNext, connections[1].RelationshipType)
}

func TestNodeConnectionEmptyFields(t *testing.T) {
	// Test connection with minimal required fields
	connection := models.NodeConnection{
		SourceNodeID:     100,
		TargetNodeID:     200,
		RelationshipType: models.RelationshipTypeRelated,
		CreatedAt:        time.Now(),
	}

	assert.Equal(t, 100, connection.SourceNodeID)
	assert.Equal(t, 200, connection.TargetNodeID)
	assert.Equal(t, models.RelationshipTypeRelated, connection.RelationshipType)
	assert.Empty(t, connection.Description) // Optional field should be empty
	assert.Equal(t, 0, connection.ID)       // Default zero value

	// Test requests with empty optional fields
	createReq := models.CreateNodeConnectionRequest{
		SourceNodeID:     100,
		TargetNodeID:     200,
		RelationshipType: models.RelationshipTypeCustom,
		// Description omitted
	}

	assert.Equal(t, 100, createReq.SourceNodeID)
	assert.Equal(t, 200, createReq.TargetNodeID)
	assert.Equal(t, models.RelationshipTypeCustom, createReq.RelationshipType)
	assert.Empty(t, createReq.Description)

	updateReq := models.UpdateNodeConnectionRequest{
		// Both fields omitted - should be empty
	}

	assert.Empty(t, updateReq.RelationshipType)
	assert.Empty(t, updateReq.Description)
}