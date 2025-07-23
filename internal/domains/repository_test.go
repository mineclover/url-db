package domains_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"url-db/internal/domains"
	"url-db/internal/models"
	"url-db/internal/shared/testdb"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Load centralized schema
	testdb.LoadSchema(t, db)

	return db
}

func TestDomainRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	domain := &models.Domain{
		Name:        "test-domain",
		Description: "Test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, domain)
	assert.NoError(t, err)
	assert.NotZero(t, domain.ID)
	assert.NotZero(t, domain.CreatedAt)
	assert.NotZero(t, domain.UpdatedAt)
}

func TestDomainRepository_Create_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	domain1 := &models.Domain{
		Name:        "test-domain",
		Description: "First domain",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	domain2 := &models.Domain{
		Name:        "test-domain",
		Description: "Second domain",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, domain1)
	assert.NoError(t, err)

	err = repo.Create(ctx, domain2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}

func TestDomainRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	// Create a domain first
	domain := &models.Domain{
		Name:        "test-domain",
		Description: "Test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, domain)
	require.NoError(t, err)

	// Get by ID
	retrieved, err := repo.GetByID(ctx, domain.ID)
	assert.NoError(t, err)
	assert.Equal(t, domain.ID, retrieved.ID)
	assert.Equal(t, domain.Name, retrieved.Name)
	assert.Equal(t, domain.Description, retrieved.Description)
}

func TestDomainRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 999)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestDomainRepository_GetByName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	// Create a domain first
	domain := &models.Domain{
		Name:        "test-domain",
		Description: "Test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, domain)
	require.NoError(t, err)

	// Get by name
	retrieved, err := repo.GetByName(ctx, "test-domain")
	assert.NoError(t, err)
	assert.Equal(t, domain.ID, retrieved.ID)
	assert.Equal(t, domain.Name, retrieved.Name)
	assert.Equal(t, domain.Description, retrieved.Description)
}

func TestDomainRepository_GetByName_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	_, err := repo.GetByName(ctx, "nonexistent")
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestDomainRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	// Create multiple domains
	domains := []*models.Domain{
		{Name: "domain1", Description: "Description 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "domain2", Description: "Description 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "domain3", Description: "Description 3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, domain := range domains {
		err := repo.Create(ctx, domain)
		require.NoError(t, err)
	}

	// Test pagination
	result, totalCount, err := repo.List(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 3, totalCount)

	// Test second page
	result, totalCount, err = repo.List(ctx, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 3, totalCount)

	// Test all items
	result, totalCount, err = repo.List(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, 3, totalCount)
}

func TestDomainRepository_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	result, totalCount, err := repo.List(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 0)
	assert.Equal(t, 0, totalCount)
}

func TestDomainRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	// Create a domain first
	domain := &models.Domain{
		Name:        "test-domain",
		Description: "Original description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, domain)
	require.NoError(t, err)

	// Update the domain
	domain.Description = "Updated description"
	
	err = repo.Update(ctx, domain)
	assert.NoError(t, err)
	// UpdatedAt should be updated by the database
	assert.NotZero(t, domain.UpdatedAt)

	// Verify update
	updated, err := repo.GetByID(ctx, domain.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated description", updated.Description)
}

func TestDomainRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	// Create a domain first
	domain := &models.Domain{
		Name:        "test-domain",
		Description: "Test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, domain)
	require.NoError(t, err)

	// Delete the domain
	err = repo.Delete(ctx, domain.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, domain.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestDomainRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestDomainRepository_ExistsByName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	// Check non-existent domain
	exists, err := repo.ExistsByName(ctx, "nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Create a domain
	domain := &models.Domain{
		Name:        "test-domain",
		Description: "Test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = repo.Create(ctx, domain)
	require.NoError(t, err)

	// Check existing domain
	exists, err = repo.ExistsByName(ctx, "test-domain")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDomainRepository_DatabaseError(t *testing.T) {
	db := setupTestDB(t)
	db.Close() // Close database to trigger errors

	repo := domains.NewDomainRepository(db)
	ctx := context.Background()

	domain := &models.Domain{
		Name:        "test-domain",
		Description: "Test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test Create error
	err := repo.Create(ctx, domain)
	assert.Error(t, err)

	// Test GetByID error
	_, err = repo.GetByID(ctx, 1)
	assert.Error(t, err)

	// Test GetByName error
	_, err = repo.GetByName(ctx, "test")
	assert.Error(t, err)

	// Test List error
	_, _, err = repo.List(ctx, 1, 10)
	assert.Error(t, err)

	// Test Update error
	err = repo.Update(ctx, domain)
	assert.Error(t, err)

	// Test Delete error
	err = repo.Delete(ctx, 1)
	assert.Error(t, err)

	// Test ExistsByName error
	_, err = repo.ExistsByName(ctx, "test")
	assert.Error(t, err)
}