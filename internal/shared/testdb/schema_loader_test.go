package testdb_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"url-db/internal/shared/testdb"
)

func TestLoadSchema_Success(t *testing.T) {
	// Create in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Load schema
	testdb.LoadSchema(t, db)

	// Verify that main tables exist by running simple queries
	tables := []string{
		"domains",
		"nodes", 
		"attributes",
		"node_attributes",
		"node_connections",
		"node_subscriptions",
		"dependency_types",
		"node_dependencies",
		"node_dependencies_v2",
		"dependency_history",
		"dependency_graph_cache",
		"dependency_rules",
		"dependency_impact_analysis",
		"node_events",
	}

	for _, table := range tables {
		query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
		var name string
		err := db.QueryRow(query, table).Scan(&name)
		assert.NoError(t, err, "Table %s should exist", table)
		assert.Equal(t, table, name, "Table name should match")
	}
}

func TestLoadSchema_WithData(t *testing.T) {
	// Create in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Load schema
	testdb.LoadSchema(t, db)

	// Test that we can insert and query data
	_, err = db.Exec("INSERT INTO domains (name, description) VALUES (?, ?)", "test-domain", "Test description")
	assert.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM domains WHERE name = ?", "test-domain").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Test dependency types are populated
	err = db.QueryRow("SELECT COUNT(*) FROM dependency_types").Scan(&count)
	assert.NoError(t, err)
	assert.Greater(t, count, 0, "Dependency types should be pre-populated")
}

func TestLoadSchema_WithForeignKeys(t *testing.T) {
	// Create in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)

	// Load schema
	testdb.LoadSchema(t, db)

	// Test foreign key constraint by inserting a domain
	_, err = db.Exec("INSERT INTO domains (name, description) VALUES (?, ?)", "test-domain", "Test description")
	require.NoError(t, err)

	// Test that foreign key constraint works for nodes
	_, err = db.Exec("INSERT INTO nodes (content, domain_id, title, description) VALUES (?, ?, ?, ?)",
		"https://example.com", 1, "Test Node", "Test description")
	assert.NoError(t, err, "Should be able to insert node with valid domain_id")

	// Test foreign key constraint violation
	_, err = db.Exec("INSERT INTO nodes (content, domain_id, title, description) VALUES (?, ?, ?, ?)",
		"https://example2.com", 999, "Test Node 2", "Test description 2")
	assert.Error(t, err, "Should fail with invalid domain_id due to foreign key constraint")
}