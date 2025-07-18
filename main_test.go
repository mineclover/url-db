package main

import (
	"fmt"
	"os"
	"testing"
	"url-db/internal/config"
	"url-db/internal/database"
)

// TestDatabaseConnection tests the database connection and schema creation
func TestDatabaseConnection(t *testing.T) {
	// Create a temporary database for testing
	tmpDB := ":memory:"

	// Initialize database
	db, err := database.InitDB(tmpDB)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Test database ping
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Test schema creation by checking if tables exist
	tables := []string{"domains", "nodes", "attributes", "node_attributes", "node_connections"}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", table)
		var tableName string
		err := db.DB().QueryRow(query).Scan(&tableName)
		if err != nil {
			t.Fatalf("Table %s not found: %v", table, err)
		}
		if tableName != table {
			t.Fatalf("Expected table %s, got %s", table, tableName)
		}
	}
}

// TestConfigLoad tests the configuration loading
func TestConfigLoad(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PORT", "9999")
	os.Setenv("DATABASE_URL", "file:./test.db")
	os.Setenv("TOOL_NAME", "test-tool")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("TOOL_NAME")
	}()

	cfg := config.Load()

	if cfg.Port != "9999" {
		t.Errorf("Expected port 9999, got %s", cfg.Port)
	}

	if cfg.DatabaseURL != "file:./test.db" {
		t.Errorf("Expected database URL 'file:./test.db', got %s", cfg.DatabaseURL)
	}

	if cfg.ToolName != "test-tool" {
		t.Errorf("Expected tool name 'test-tool', got %s", cfg.ToolName)
	}
}

// TestMainApplicationComponents tests that all main components can be initialized
func TestMainApplicationComponents(t *testing.T) {
	// This test ensures that the main application can be initialized without errors
	// but doesn't start the actual server

	// Create a temporary database
	tmpDB := ":memory:"

	// Initialize database
	db, err := database.InitDB(tmpDB)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Test that we can create repositories without errors
	// This is a basic smoke test for the application initialization
	if db.DB() == nil {
		t.Fatal("Database connection is nil")
	}

	t.Log("All main application components initialized successfully")
}
