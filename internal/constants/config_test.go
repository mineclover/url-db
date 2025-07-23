package constants_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"url-db/internal/constants"
)

func TestConfigConstants(t *testing.T) {
	// Test database constants
	assert.Equal(t, "url-db.sqlite", constants.DefaultDBPath)
	assert.Equal(t, "sqlite3", constants.DefaultDBDriver)
	assert.Equal(t, "url-db", constants.DefaultServerName)

	// Test server constants
	assert.Equal(t, 8080, constants.DefaultPort)
	assert.Equal(t, "stdio", constants.DefaultMCPMode)
}

func TestMCPConstants(t *testing.T) {
	// Test MCP protocol constants
	assert.Equal(t, "2024-11-05", constants.MCPProtocolVersion)
	assert.Equal(t, "url-db-mcp-server", constants.MCPServerName)
	assert.Equal(t, "1.0.0", constants.DefaultServerVersion)
	assert.Equal(t, "2.0", constants.JSONRPCVersion)
}

func TestErrorConstants(t *testing.T) {
	// Test error message constants
	assert.Equal(t, "domain not found", constants.ErrDomainNotFound)
	assert.Equal(t, "node not found", constants.ErrNodeNotFound)
	assert.Equal(t, "attribute not found", constants.ErrAttributeNotFound)
	assert.Equal(t, "invalid composite ID format", constants.ErrInvalidCompositeID)
	assert.Equal(t, "domain already exists", constants.ErrDuplicateDomain)
}

func TestValidationConstants(t *testing.T) {
	// Test validation limits
	assert.Equal(t, 50, constants.MaxDomainNameLength)
	assert.Equal(t, 255, constants.MaxTitleLength)
	assert.Equal(t, 1000, constants.MaxDescriptionLength)
	assert.Equal(t, 2048, constants.MaxURLLength)
	assert.Equal(t, 2048, constants.MaxAttributeValueLength)
}

func TestCompositeKeyConstants(t *testing.T) {
	// Test composite key constants
	assert.Equal(t, ":", constants.CompositeKeySeparator)
	assert.Equal(t, "url-db:domain:id", constants.CompositeKeyFormat)
}

func TestResourceConstants(t *testing.T) {
	// Test MCP resource constants
	assert.Equal(t, "mcp", constants.MCPResourceScheme)
	assert.Equal(t, "file", constants.FileResourceScheme)
	assert.Equal(t, "http", constants.HTTPResourceScheme)
}

func TestHTTPConstants(t *testing.T) {
	// Test HTTP constants
	assert.Equal(t, 200, constants.StatusOK)
	assert.Equal(t, 201, constants.StatusCreated)
	assert.Equal(t, 400, constants.StatusBadRequest)
	assert.Equal(t, 404, constants.StatusNotFound)
	assert.Equal(t, 500, constants.StatusInternalServerError)
}

func TestValidationPatterns(t *testing.T) {
	// Test validation patterns
	assert.Equal(t, "^[a-zA-Z0-9_-]+$", constants.DomainNamePattern)
	assert.Equal(t, "^https?://.*", constants.URLPattern)
	assert.Contains(t, constants.EmailPattern, "@")
}