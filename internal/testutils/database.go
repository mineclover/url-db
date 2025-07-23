package testutils

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"url-db/internal/infrastructure/database"
)

// SetupTestDB creates an in-memory database using the single source of truth schema
func SetupTestDB(t *testing.T) *sql.DB {
	// Create an in-memory database
	config := database.DefaultConfig()
	config.URL = ":memory:"
	
	db, err := database.New(config)
	require.NoError(t, err, "Failed to create test database")
	
	// Create common test data
	CreateTestData(t, db.DB())
	
	// Return the underlying sql.DB for compatibility with existing tests
	return db.DB()
}

// SetupTestDatabase creates a test database instance with proper configuration
func SetupTestDatabase(t *testing.T) *database.Database {
	config := database.DefaultConfig()
	config.URL = ":memory:"
	
	db, err := database.New(config)
	require.NoError(t, err, "Failed to create test database")
	
	return db
}

// CleanupTestDB properly closes the test database
func CleanupTestDB(t *testing.T, db *sql.DB) {
	if db != nil {
		err := db.Close()
		require.NoError(t, err, "Failed to close test database")
	}
}

// CleanupTestDatabase properly closes the test database instance
func CleanupTestDatabase(t *testing.T, db *database.Database) {
	if db != nil {
		err := db.Close()
		require.NoError(t, err, "Failed to close test database")
	}
}

// CreateTestData inserts common test data for tests
func CreateTestData(t *testing.T, db *sql.DB) {
	// Insert test domain
	_, err := db.Exec(`
		INSERT OR IGNORE INTO domains (id, name, description, created_at, updated_at) 
		VALUES (1, 'test-domain', 'Test domain for testing', datetime('now'), datetime('now'))
	`)
	require.NoError(t, err, "Failed to insert test domain")
	
	// Insert test nodes  
	_, err = db.Exec(`
		INSERT OR IGNORE INTO nodes (id, content, domain_id, title, description, created_at, updated_at) 
		VALUES 
		(1, 'https://example.com/test1', 1, 'Test Node 1', 'First test node', datetime('now'), datetime('now')),
		(2, 'https://example.com/test2', 1, 'Test Node 2', 'Second test node', datetime('now'), datetime('now'))
	`)
	require.NoError(t, err, "Failed to insert test nodes")
	
	// Insert test attributes
	_, err = db.Exec(`
		INSERT OR IGNORE INTO attributes (id, domain_id, name, type, description, created_at) 
		VALUES 
		(1, 1, 'test-tag', 'tag', 'Test tag attribute', datetime('now')),
		(2, 1, 'test-string', 'string', 'Test string attribute', datetime('now'))
	`)
	require.NoError(t, err, "Failed to insert test attributes")
}