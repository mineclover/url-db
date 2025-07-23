package testdb

import (
	"database/sql"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed schema.sql
var embeddedSchema embed.FS

// LoadSchema loads the centralized schema.sql file for test databases.
// It first tries to load from the project root, then falls back to embedded schema.
func LoadSchema(t *testing.T, db *sql.DB) {
	schema := getSchemaContent(t)
	_, err := db.Exec(schema)
	require.NoError(t, err, "Failed to execute schema")
}

// getSchemaContent returns the schema content from file or embedded fallback
func getSchemaContent(t *testing.T) string {
	// Try to load from project root first
	if schema := loadSchemaFromFile(); schema != "" {
		return schema
	}

	// Fall back to embedded schema
	t.Log("Using embedded schema fallback")
	return loadEmbeddedSchema(t)
}

// loadSchemaFromFile attempts to load schema.sql from project root
func loadSchemaFromFile() string {
	// Get the project root by finding go.mod
	projectRoot := findProjectRoot()
	if projectRoot == "" {
		return ""
	}

	schemaPath := filepath.Join(projectRoot, "schema.sql")
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return ""
	}

	return string(content)
}

// loadEmbeddedSchema loads schema from embedded file system
func loadEmbeddedSchema(t *testing.T) string {
	content, err := fs.ReadFile(embeddedSchema, "schema.sql")
	require.NoError(t, err, "Failed to read embedded schema")
	return string(content)
}

// findProjectRoot finds the project root by looking for go.mod
func findProjectRoot() string {
	// Start from current file location
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	dir := filepath.Dir(currentFile)
	
	// Walk up the directory tree to find go.mod
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}
	
	return ""
}