package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "default config",
			config:  nil,
			wantErr: false,
		},
		{
			name:    "test config",
			config:  TestConfig(),
			wantErr: false,
		},
		{
			name: "custom config",
			config: &Config{
				URL:             ":memory:",
				MaxOpenConns:    5,
				MaxIdleConns:    2,
				ConnMaxLifetime: time.Minute * 30,
				WALMode:         false,
				ForeignKeys:     true,
				JournalMode:     "DELETE",
				Synchronous:     "OFF",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				defer db.Close()

				err = db.Ping()
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatabase_createSchema(t *testing.T) {
	db, err := New(TestConfig())
	require.NoError(t, err)
	defer db.Close()

	tables := []string{
		"domains",
		"nodes",
		"attributes",
		"node_attributes",
		"node_connections",
	}

	for _, table := range tables {
		var name string
		err := db.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		assert.NoError(t, err)
		assert.Equal(t, table, name)
	}

	indexes := []string{
		"idx_nodes_domain",
		"idx_nodes_content",
		"idx_attributes_domain",
		"idx_node_attributes_node",
		"idx_node_attributes_attribute",
		"idx_node_connections_source",
		"idx_node_connections_target",
	}

	for _, index := range indexes {
		var name string
		err := db.db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name=?", index).Scan(&name)
		assert.NoError(t, err)
		assert.Equal(t, index, name)
	}
}

func TestDatabase_WithTransaction(t *testing.T) {
	db, err := New(TestConfig())
	require.NoError(t, err)
	defer db.Close()

	t.Run("successful transaction", func(t *testing.T) {
		err := db.WithTransaction(func(tx *sql.Tx) error {
			_, err := tx.Exec("INSERT INTO domains (name, description) VALUES (?, ?)", "test", "test domain")
			return err
		})
		assert.NoError(t, err)

		var count int
		err = db.db.QueryRow("SELECT COUNT(*) FROM domains WHERE name = ?", "test").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("failed transaction rollback", func(t *testing.T) {
		err := db.WithTransaction(func(tx *sql.Tx) error {
			_, err := tx.Exec("INSERT INTO domains (name, description) VALUES (?, ?)", "test2", "test domain 2")
			if err != nil {
				return err
			}
			return assert.AnError
		})
		assert.Error(t, err)

		var count int
		err = db.db.QueryRow("SELECT COUNT(*) FROM domains WHERE name = ?", "test2").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestDatabase_Close(t *testing.T) {
	db, err := New(TestConfig())
	require.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)

	err = db.Ping()
	assert.Error(t, err)
}

func TestDatabase_Ping(t *testing.T) {
	db, err := New(TestConfig())
	require.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	assert.NoError(t, err)
}

func TestDatabase_DB(t *testing.T) {
	db, err := New(TestConfig())
	require.NoError(t, err)
	defer db.Close()

	sqlDB := db.DB()
	assert.NotNil(t, sqlDB)
	assert.Equal(t, db.db, sqlDB)
}

func TestConfigureDatabase(t *testing.T) {
	config := TestConfig()
	db, err := New(config)
	require.NoError(t, err)
	defer db.Close()

	var foreignKeys bool
	err = db.db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys)
	assert.NoError(t, err)
	assert.True(t, foreignKeys)

	var journalMode string
	err = db.db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	assert.NoError(t, err)
	assert.Equal(t, config.JournalMode, journalMode)

	var synchronous string
	err = db.db.QueryRow("PRAGMA synchronous").Scan(&synchronous)
	assert.NoError(t, err)
	assert.Equal(t, "0", synchronous)
}
