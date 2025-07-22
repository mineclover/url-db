package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMCPError_Error(t *testing.T) {
	err := &MCPError{
		Code:    "INVALID_REQUEST",
		Message: "Invalid Request",
	}

	assert.Contains(t, err.Error(), "Invalid Request")
}

func TestNewInvalidCompositeKeyError(t *testing.T) {
	err := NewInvalidCompositeKeyError("invalid-key")

	assert.NotNil(t, err)
	assert.Equal(t, "INVALID_COMPOSITE_KEY", err.Code)
	assert.NotEmpty(t, err.Message)
}

func TestNewDomainNotFoundError(t *testing.T) {
	err := NewDomainNotFoundError("test-domain")

	assert.NotNil(t, err)
	assert.Equal(t, "DOMAIN_NOT_FOUND", err.Code)
	assert.NotEmpty(t, err.Message)
}

func TestNewResourceNotFoundError(t *testing.T) {
	err := NewResourceNotFoundError("test-resource")

	assert.NotNil(t, err)
	assert.Equal(t, "RESOURCE_NOT_FOUND", err.Code)
	assert.NotEmpty(t, err.Message)
}

func TestNewNodeNotFoundError(t *testing.T) {
	err := NewNodeNotFoundError("test-domain", "https://example.com")

	assert.NotNil(t, err)
	assert.Equal(t, "NODE_NOT_FOUND", err.Code)
	assert.NotEmpty(t, err.Message)
}

func TestNewAccessDeniedError(t *testing.T) {
	err := NewAccessDeniedError("test operation")

	assert.NotNil(t, err)
	assert.Equal(t, "ACCESS_DENIED", err.Code)
	assert.NotEmpty(t, err.Message)
}

func TestNewBatchPartialFailureError(t *testing.T) {
	failedItems := []string{"item1", "item2"}
	err := NewBatchPartialFailureError(failedItems)

	assert.NotNil(t, err)
	assert.Equal(t, "BATCH_PARTIAL_FAILURE", err.Code)
	assert.NotEmpty(t, err.Message)
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("invalid input")

	assert.NotNil(t, err)
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.NotEmpty(t, err.Message)
}

func TestNewInternalServerError(t *testing.T) {
	err := NewInternalServerError("server failure")

	assert.NotNil(t, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", err.Code)
	assert.NotEmpty(t, err.Message)
}
