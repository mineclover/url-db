package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Schema file path relative to project root
const schemaFilePath = "schema.sql"

type Database struct {
	db     *sql.DB
	sqlxDB *sqlx.DB
	config *Config
}

func New(config *Config) (*Database, error) {
	if config == nil {
		config = DefaultConfig()
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

// loadSchemaFromFile loads schema from external schema.sql file
func (d *Database) loadSchemaFromFile() (string, error) {
	// Find project root by looking for go.mod
	projectRoot, err := findProjectRoot()
	if err != nil {
		return "", fmt.Errorf("could not find project root: %w", err)
	}

	schemaPath := filepath.Join(projectRoot, schemaFilePath)
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return "", fmt.Errorf("could not read schema file %s: %w", schemaPath, err)
	}

	return string(schemaBytes), nil
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
