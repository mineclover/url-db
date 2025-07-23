package usecases_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-db/internal/core/usecases"
	"url-db/internal/models"
	"url-db/internal/testutils"
)

func setupTestDB(t *testing.T) *sql.DB {
	return testutils.SetupTestDB(t)
}

func TestSQLiteAttributeRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	attribute := &models.Attribute{
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Test description",
		CreatedAt:   time.Now(),
	}

	err := repo.Create(ctx, attribute)
	assert.NoError(t, err)
	assert.NotZero(t, attribute.ID)
}

func TestSQLiteAttributeRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	// Create attribute
	attribute := &models.Attribute{
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Test description",
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, attribute)
	require.NoError(t, err)

	// Get by ID
	retrieved, err := repo.GetByID(ctx, attribute.ID)
	assert.NoError(t, err)
	assert.Equal(t, attribute.ID, retrieved.ID)
	assert.Equal(t, attribute.Name, retrieved.Name)
	assert.Equal(t, attribute.Type, retrieved.Type)
	assert.Equal(t, attribute.Description, retrieved.Description)
}

func TestSQLiteAttributeRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 999)
	assert.ErrorIs(t, err, usecases.ErrAttributeNotFound)
}

func TestSQLiteAttributeRepository_GetByDomainID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	// Create multiple attributes
	attributes := []*models.Attribute{
		{
			DomainID:    1,
			Name:        "attr1",
			Type:        models.AttributeTypeTag,
			Description: "Description 1",
			CreatedAt:   time.Now(),
		},
		{
			DomainID:    1,
			Name:        "attr2",
			Type:        models.AttributeTypeString,
			Description: "Description 2",
			CreatedAt:   time.Now(),
		},
	}

	for _, attr := range attributes {
		err := repo.Create(ctx, attr)
		require.NoError(t, err)
	}

	// Get by domain ID
	retrieved, err := repo.GetByDomainID(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 2)
	assert.Equal(t, "attr1", retrieved[0].Name) // Should be sorted by name
	assert.Equal(t, "attr2", retrieved[1].Name)
}

func TestSQLiteAttributeRepository_GetByDomainIDAndName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	// Create attribute
	attribute := &models.Attribute{
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Test description",
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, attribute)
	require.NoError(t, err)

	// Get by domain ID and name
	retrieved, err := repo.GetByDomainIDAndName(ctx, 1, "test-attribute")
	assert.NoError(t, err)
	assert.Equal(t, attribute.ID, retrieved.ID)
	assert.Equal(t, attribute.Name, retrieved.Name)
}

func TestSQLiteAttributeRepository_GetByDomainIDAndName_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	_, err := repo.GetByDomainIDAndName(ctx, 1, "nonexistent")
	assert.ErrorIs(t, err, usecases.ErrAttributeNotFound)
}

func TestSQLiteAttributeRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	// Create attribute
	attribute := &models.Attribute{
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Original description",
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, attribute)
	require.NoError(t, err)

	// Update attribute
	attribute.Description = "Updated description"
	err = repo.Update(ctx, attribute)
	assert.NoError(t, err)

	// Verify update
	retrieved, err := repo.GetByID(ctx, attribute.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated description", retrieved.Description)
}

func TestSQLiteAttributeRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	attribute := &models.Attribute{
		ID:          999,
		Description: "Updated description",
	}

	err := repo.Update(ctx, attribute)
	assert.ErrorIs(t, err, usecases.ErrAttributeNotFound)
}

func TestSQLiteAttributeRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	// Create attribute
	attribute := &models.Attribute{
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Test description",
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, attribute)
	require.NoError(t, err)

	// Delete attribute
	err = repo.Delete(ctx, attribute.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, attribute.ID)
	assert.ErrorIs(t, err, usecases.ErrAttributeNotFound)
}

func TestSQLiteAttributeRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	assert.ErrorIs(t, err, usecases.ErrAttributeNotFound)
}

func TestSQLiteAttributeRepository_HasValues(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := usecases.NewSQLiteAttributeRepository(db)
	ctx := context.Background()

	// Create attribute
	attribute := &models.Attribute{
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Test description",
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, attribute)
	require.NoError(t, err)

	// Check if attribute has values (should be false)
	hasValues, err := repo.HasValues(ctx, attribute.ID)
	assert.NoError(t, err)
	assert.False(t, hasValues)

	// Create node and node attribute
	_, err = db.Exec(`INSERT INTO nodes (id, domain_id, url, created_at) VALUES (1, 1, 'http://example.com', '2023-01-01 00:00:00')`)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO node_attributes (node_id, attribute_id, value, created_at) VALUES (1, ?, 'test-value', '2023-01-01 00:00:00')`, attribute.ID)
	require.NoError(t, err)

	// Check if attribute has values (should be true)
	hasValues, err = repo.HasValues(ctx, attribute.ID)
	assert.NoError(t, err)
	assert.True(t, hasValues)
}
