package nodeattributes_test

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"url-db/internal/models"
	"url-db/internal/nodeattributes"
	"url-db/internal/testutils"
)

// Helper function to create *int pointers
func intPtr(i int) *int {
	return &i
}

func setupTestDB(t *testing.T) *sqlx.DB {
	// Get the sql.DB from testutils and wrap it with sqlx
	sqlDB := testutils.SetupTestDB(t)
	db := sqlx.NewDb(sqlDB, "sqlite3")
	
	return db
}

func TestNewRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)
	assert.NotNil(t, repo)
}

func TestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	nodeAttr := &models.NodeAttribute{
		NodeID:      1,
		AttributeID: 1,
		Value:       "test-value",
		OrderIndex:  intPtr(1),
		CreatedAt:   time.Now(),
	}

	err := repo.Create(nodeAttr)
	assert.NoError(t, err)
	assert.NotZero(t, nodeAttr.ID)
}

func TestRepository_Create_InvalidNodeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	nodeAttr := &models.NodeAttribute{
		NodeID:      999, // Non-existent node
		AttributeID: 1,
		Value:       "test-value",
		OrderIndex:  intPtr(1),
		CreatedAt:   time.Now(),
	}

	err := repo.Create(nodeAttr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FOREIGN KEY constraint failed")
}

func TestRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	// Create a node attribute first
	nodeAttr := &models.NodeAttribute{
		NodeID:      1,
		AttributeID: 1,
		Value:       "test-value",
		OrderIndex:  intPtr(1),
		CreatedAt:   time.Now(),
	}
	
	err := repo.Create(nodeAttr)
	require.NoError(t, err)

	// Retrieve by ID
	retrieved, err := repo.GetByID(nodeAttr.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, nodeAttr.ID, retrieved.ID)
	assert.Equal(t, nodeAttr.NodeID, retrieved.NodeID)
	assert.Equal(t, nodeAttr.AttributeID, retrieved.AttributeID)
	assert.Equal(t, nodeAttr.Value, retrieved.Value)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	retrieved, err := repo.GetByID(999)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestRepository_GetByNodeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	// Create multiple node attributes
	attrs := []*models.NodeAttribute{
		{NodeID: 1, AttributeID: 1, Value: "value1", OrderIndex: intPtr(1), CreatedAt: time.Now()},
	}

	for _, attr := range attrs {
		err := repo.Create(attr)
		require.NoError(t, err)
	}

	// Retrieve by node ID
	retrieved, err := repo.GetByNodeID(1)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 1)
	assert.Equal(t, "value1", retrieved[0].Value)
}

func TestRepository_GetByNodeID_EmptyResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	retrieved, err := repo.GetByNodeID(999) // Non-existent node
	assert.NoError(t, err)
	assert.Len(t, retrieved, 0)
}

func TestRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	// Create a node attribute first
	nodeAttr := &models.NodeAttribute{
		NodeID:      1,
		AttributeID: 1,
		Value:       "original-value",
		OrderIndex:  intPtr(1),
		CreatedAt:   time.Now(),
	}
	
	err := repo.Create(nodeAttr)
	require.NoError(t, err)

	// Update the attribute
	updateReq := &models.UpdateNodeAttributeRequest{
		Value:      "updated-value",
		OrderIndex: intPtr(2),
	}

	updated, err := repo.Update(nodeAttr.ID, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "updated-value", updated.Value)
	assert.NotNil(t, updated.OrderIndex)
	assert.Equal(t, 2, *updated.OrderIndex)
}

func TestRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	updateReq := &models.UpdateNodeAttributeRequest{
		Value: "updated-value",
	}

	updated, err := repo.Update(999, updateReq)
	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "failed to get updated node attribute")
}

func TestRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	// Create a node attribute first
	nodeAttr := &models.NodeAttribute{
		NodeID:      1,
		AttributeID: 1,
		Value:       "test-value",
		OrderIndex:  intPtr(1),
		CreatedAt:   time.Now(),
	}
	
	err := repo.Create(nodeAttr)
	require.NoError(t, err)

	// Delete the attribute
	err = repo.Delete(nodeAttr.ID)
	assert.NoError(t, err)

	// Verify deletion
	deleted, err := repo.GetByID(nodeAttr.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}

func TestRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	err := repo.Delete(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRepository_DeleteByNodeIDAndAttributeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	// Create a node attribute first
	nodeAttr := &models.NodeAttribute{
		NodeID:      1,
		AttributeID: 1,
		Value:       "test-value",
		OrderIndex:  intPtr(1),
		CreatedAt:   time.Now(),
	}
	
	err := repo.Create(nodeAttr)
	require.NoError(t, err)

	// Delete by node ID and attribute ID
	err = repo.DeleteByNodeIDAndAttributeID(1, 1)
	assert.NoError(t, err)

	// Verify deletion
	deleted, err := repo.GetByID(nodeAttr.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}

func TestRepository_DeleteByNodeIDAndAttributeID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	err := repo.DeleteByNodeIDAndAttributeID(999, 999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRepository_GetByNodeIDAndAttributeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	// Create a node attribute first
	nodeAttr := &models.NodeAttribute{
		NodeID:      1,
		AttributeID: 1,
		Value:       "test-value",
		OrderIndex:  intPtr(1),
		CreatedAt:   time.Now(),
	}
	
	err := repo.Create(nodeAttr)
	require.NoError(t, err)

	// Retrieve by node ID and attribute ID
	retrieved, err := repo.GetByNodeIDAndAttributeID(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, nodeAttr.Value, retrieved.Value)
}

func TestRepository_GetByNodeIDAndAttributeID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	retrieved, err := repo.GetByNodeIDAndAttributeID(999, 999)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestRepository_GetMaxOrderIndex(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	// Create multiple node attributes with different order indices
	attrs := []*models.NodeAttribute{
		{NodeID: 1, AttributeID: 1, Value: "value1", OrderIndex: intPtr(1), CreatedAt: time.Now()},
		{NodeID: 1, AttributeID: 1, Value: "value2", OrderIndex: intPtr(3), CreatedAt: time.Now()},
		{NodeID: 1, AttributeID: 1, Value: "value3", OrderIndex: intPtr(2), CreatedAt: time.Now()},
	}

	for _, attr := range attrs {
		err := repo.Create(attr)
		require.NoError(t, err)
	}

	// Get max order index
	maxIndex, err := repo.GetMaxOrderIndex(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 3, maxIndex)
}

func TestRepository_GetMaxOrderIndex_NoAttributes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	maxIndex, err := repo.GetMaxOrderIndex(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, maxIndex) // Should return 0 when no attributes exist
}

func TestRepository_ValidateNodeAndAttributeDomain(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	// Valid case - node and attribute are in the same domain
	err := repo.ValidateNodeAndAttributeDomain(1, 1)
	assert.NoError(t, err)
}

func TestRepository_ValidateNodeAndAttributeDomain_DomainMismatch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Add another domain and attribute for testing
	_, err := db.Exec(`INSERT INTO domains (id, name, created_at) VALUES (2, 'domain2', '2023-01-01 00:00:00')`)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO attributes (id, domain_id, name, type, created_at) VALUES (2, 2, 'attr2', 'string', '2023-01-01 00:00:00')`)
	require.NoError(t, err)

	repo := nodeattributes.NewRepository(db)

	// Invalid case - node is in domain 1, attribute is in domain 2
	err = repo.ValidateNodeAndAttributeDomain(1, 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node and attribute must belong to the same domain")
}

func TestRepository_GetAttributeType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	attrType, err := repo.GetAttributeType(1)
	assert.NoError(t, err)
	assert.Equal(t, models.AttributeType("string"), attrType)
}

func TestRepository_GetAttributeType_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := nodeattributes.NewRepository(db)

	_, err := repo.GetAttributeType(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "attribute not found")
}