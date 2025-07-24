package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Schema file path relative to project root
const schemaFilePath = "schema.sql"

// isMCPServerMode checks if the application is running in MCP server mode
func isMCPServerMode() bool {
	return os.Getenv("MCP_MODE") == "stdio" || 
		   strings.Contains(strings.Join(os.Args, " "), "-mcp-mode=stdio")
}

// logInfo logs info message only if not in MCP stdio mode
func logInfo(format string, args ...interface{}) {
	if !isMCPServerMode() {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

type Database struct {
	db     *sql.DB
	sqlxDB *sqlx.DB
	config *Config
}

func New(config *Config) (*Database, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Ensure database file and directory exist
	if err := ensureDatabaseExists(config.URL); err != nil {
		return nil, fmt.Errorf("failed to ensure database exists: %w", err)
	}

	db, err := sql.Open("sqlite3", config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := configureDatabase(db, config); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to configure database: %w", err)
	}

	// Create sqlx wrapper
	sqlxDB := sqlx.NewDb(db, "sqlite3")

	database := &Database{
		db:     db,
		sqlxDB: sqlxDB,
		config: config,
	}

	if err := database.createSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return database, nil
}

func configureDatabase(db *sql.DB, config *Config) error {
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	pragmas := []string{
		fmt.Sprintf("PRAGMA journal_mode = %s", config.JournalMode),
		fmt.Sprintf("PRAGMA synchronous = %s", config.Synchronous),
	}

	if config.ForeignKeys {
		pragmas = append(pragmas, "PRAGMA foreign_keys = ON")
	}

	if config.WALMode {
		pragmas = append(pragmas, "PRAGMA journal_mode = WAL")
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma %s: %w", pragma, err)
		}
	}

	return nil
}

func (d *Database) createSchema() error {
	// Load schema from external file (single source of truth)
	schema, err := d.loadSchemaFromFile()
	if err != nil {
		return fmt.Errorf("failed to load schema from %s: %w", schemaFilePath, err)
	}

	if _, err := d.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// loadSchemaFromFile loads schema with multiple fallback strategies
func (d *Database) loadSchemaFromFile() (string, error) {
	var lastErr error
	
	// Strategy 1: Try to find project root by looking for go.mod
	if projectRoot, err := findProjectRoot(); err == nil {
		schemaPath := filepath.Join(projectRoot, schemaFilePath)
		if schemaBytes, err := os.ReadFile(schemaPath); err == nil {
			logInfo("[INFO] Schema loaded from project root: %s\n", schemaPath)
			return string(schemaBytes), nil
		} else {
			lastErr = err
		}
	} else {
		lastErr = err
	}

	// Strategy 2: Try relative to executable
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		// Try ../schema.sql (if executable is in bin/)
		schemaPath := filepath.Join(execDir, "..", schemaFilePath)
		if schemaBytes, err := os.ReadFile(schemaPath); err == nil {
			logInfo("[INFO] Schema loaded relative to executable: %s\n", schemaPath)
			return string(schemaBytes), nil
		}
		// Try ./schema.sql (same directory as executable)
		schemaPath = filepath.Join(execDir, schemaFilePath)
		if schemaBytes, err := os.ReadFile(schemaPath); err == nil {
			logInfo("[INFO] Schema loaded from executable directory: %s\n", schemaPath)
			return string(schemaBytes), nil
		}
	}

	// Strategy 3: Try current working directory
	if cwd, err := os.Getwd(); err == nil {
		schemaPath := filepath.Join(cwd, schemaFilePath)
		if schemaBytes, err := os.ReadFile(schemaPath); err == nil {
			logInfo("[INFO] Schema loaded from working directory: %s\n", schemaPath)
			return string(schemaBytes), nil
		}
	}

	// Strategy 4: Use minimal fallback schema (always succeeds)
	logInfo("[INFO] Using minimal fallback schema (last error: %v)\n", lastErr)
	return getFallbackSchema(), nil
}

// findProjectRoot finds the project root by looking for go.mod file
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree looking for go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root directory
		}
		dir = parent
	}

	return "", fmt.Errorf("go.mod not found")
}

// getFallbackSchema returns a minimal embedded schema as fallback
func getFallbackSchema() string {
	return `
-- Fallback minimal schema for URL-DB MCP Server
-- Domains table
CREATE TABLE IF NOT EXISTS domains (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Nodes table  
CREATE TABLE IF NOT EXISTS nodes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	domain_id INTEGER NOT NULL,
	url TEXT NOT NULL,
	title TEXT,
	description TEXT,
	content TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
	UNIQUE(domain_id, url)
);

-- Attributes table
CREATE TABLE IF NOT EXISTS attributes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	domain_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL DEFAULT 'tag',
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
	UNIQUE(domain_id, name)
);

-- Node attributes table
CREATE TABLE IF NOT EXISTS node_attributes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	node_id INTEGER NOT NULL,
	attribute_id INTEGER NOT NULL,
	value TEXT NOT NULL,
	order_index INTEGER,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE,
	FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE
);

-- Templates table
CREATE TABLE IF NOT EXISTS templates (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	domain_id INTEGER NOT NULL,
	template_data TEXT NOT NULL,
	title TEXT,
	description TEXT,
	is_active BOOLEAN DEFAULT TRUE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
	UNIQUE(name, domain_id)
);

-- Template attributes table
CREATE TABLE IF NOT EXISTS template_attributes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	template_id INTEGER NOT NULL,
	attribute_id INTEGER NOT NULL,
	value TEXT NOT NULL,
	order_index INTEGER,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE CASCADE,
	FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE
);

-- Basic indexes
CREATE INDEX IF NOT EXISTS idx_nodes_domain ON nodes(domain_id);
CREATE INDEX IF NOT EXISTS idx_attributes_domain ON attributes(domain_id);
CREATE INDEX IF NOT EXISTS idx_node_attributes_node ON node_attributes(node_id);
CREATE INDEX IF NOT EXISTS idx_node_attributes_attribute ON node_attributes(attribute_id);
CREATE INDEX IF NOT EXISTS idx_templates_domain ON templates(domain_id);
CREATE INDEX IF NOT EXISTS idx_template_attributes_template ON template_attributes(template_id);
CREATE INDEX IF NOT EXISTS idx_template_attributes_attribute ON template_attributes(attribute_id);

-- Update triggers
CREATE TRIGGER IF NOT EXISTS nodes_updated_at 
	AFTER UPDATE ON nodes 
	FOR EACH ROW 
	BEGIN 
		UPDATE nodes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

CREATE TRIGGER IF NOT EXISTS domains_updated_at 
	AFTER UPDATE ON domains 
	FOR EACH ROW 
	BEGIN 
		UPDATE domains SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

CREATE TRIGGER IF NOT EXISTS templates_updated_at 
	AFTER UPDATE ON templates 
	FOR EACH ROW 
	BEGIN 
		UPDATE templates SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;
`
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *Database) Ping() error {
	return d.db.Ping()
}

func (d *Database) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// SQLXDB returns the sqlx database instance
func (d *Database) SQLXDB() *sqlx.DB {
	return d.sqlxDB
}

// InitDB initializes the database with the given URL
func InitDB(url string) (*Database, error) {
	config := DefaultConfig()
	config.URL = url
	return New(config)
}

// ensureDatabaseExists creates the database file and directory if they don't exist
func ensureDatabaseExists(url string) error {
	// Parse the database URL to extract the file path
	dbPath, err := parseDatabasePath(url)
	if err != nil {
		return fmt.Errorf("failed to parse database path: %w", err)
	}

	// Skip if it's an in-memory database
	if dbPath == ":memory:" || dbPath == "" {
		return nil
	}

	// Check if database file already exists
	if _, err := os.Stat(dbPath); err == nil {
		// File exists, nothing to do
		return nil
	} else if !os.IsNotExist(err) {
		// Some other error occurred
		return fmt.Errorf("failed to check database file: %w", err)
	}

	// Create directory structure if it doesn't exist
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory %s: %w", dir, err)
		}
		logInfo("[INFO] Created database directory: %s\n", dir)
	}

	// Create empty database file
	file, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database file %s: %w", dbPath, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close database file: %w", err)
	}
	
	logInfo("[INFO] Created database file: %s\n", dbPath)
	return nil
}

// parseDatabasePath extracts the file path from SQLite database URL
func parseDatabasePath(url string) (string, error) {
	// Handle common SQLite URL formats:
	// - file:path/to/db.sqlite
	// - file://path/to/db.sqlite  
	// - path/to/db.sqlite (direct path)
	// - :memory: (in-memory)

	if url == ":memory:" {
		return ":memory:", nil
	}

	// Remove file: prefix if present
	if strings.HasPrefix(url, "file:") {
		url = strings.TrimPrefix(url, "file:")
		// Handle file:// format
		if strings.HasPrefix(url, "//") {
			url = strings.TrimPrefix(url, "//")
		}
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(url) {
		abs, err := filepath.Abs(url)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path: %w", err)
		}
		url = abs
	}

	return url, nil
}
