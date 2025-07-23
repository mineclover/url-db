package database_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"url-db/internal/database"
)

func TestDefaultConfig(t *testing.T) {
	config := database.DefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "file:./app.db", config.URL)
	assert.Equal(t, 10, config.MaxOpenConns)
	assert.Equal(t, 5, config.MaxIdleConns)
	assert.Equal(t, time.Hour, config.ConnMaxLifetime)
	assert.True(t, config.WALMode)
	assert.True(t, config.ForeignKeys)
	assert.Equal(t, "WAL", config.JournalMode)
	assert.Equal(t, "NORMAL", config.Synchronous)
}

func TestTestConfig(t *testing.T) {
	config := database.TestConfig()

	assert.NotNil(t, config)
	assert.Equal(t, ":memory:", config.URL)
	assert.Equal(t, 1, config.MaxOpenConns)
	assert.Equal(t, 1, config.MaxIdleConns)
	assert.Equal(t, time.Hour, config.ConnMaxLifetime)
	assert.False(t, config.WALMode)
	assert.True(t, config.ForeignKeys)
	assert.Equal(t, "DELETE", config.JournalMode)
	assert.Equal(t, "OFF", config.Synchronous)
}

func TestProductionConfig(t *testing.T) {
	dbPath := "/path/to/production.db"
	config := database.ProductionConfig(dbPath)

	assert.NotNil(t, config)
	assert.Equal(t, "file:"+dbPath, config.URL)
	assert.Equal(t, 100, config.MaxOpenConns)
	assert.Equal(t, 50, config.MaxIdleConns)
	assert.Equal(t, time.Hour, config.ConnMaxLifetime)
	assert.True(t, config.WALMode)
	assert.True(t, config.ForeignKeys)
	assert.Equal(t, "WAL", config.JournalMode)
	assert.Equal(t, "FULL", config.Synchronous)
}

func TestProductionConfig_WithEmptyPath(t *testing.T) {
	config := database.ProductionConfig("")

	assert.NotNil(t, config)
	assert.Equal(t, "file:", config.URL)
}

func TestProductionConfig_WithRelativePath(t *testing.T) {
	config := database.ProductionConfig("data/app.db")

	assert.NotNil(t, config)
	assert.Equal(t, "file:data/app.db", config.URL)
}