package compositekey_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"url-db/internal/compositekey"
)

func TestCompositeKeyError_Error(t *testing.T) {
	err := compositekey.CompositeKeyError{
		Code:    "TEST_ERROR",
		Message: "This is a test error",
	}

	expected := "TEST_ERROR: This is a test error"
	assert.Equal(t, expected, err.Error())
}

func TestNewInvalidFormatError(t *testing.T) {
	message := "Invalid format provided"
	err := compositekey.NewInvalidFormatError(message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), compositekey.ErrInvalidFormat)
	assert.Contains(t, err.Error(), message)

	// Test that it's the correct type
	compositeErr, ok := err.(compositekey.CompositeKeyError)
	assert.True(t, ok)
	assert.Equal(t, compositekey.ErrInvalidFormat, compositeErr.Code)
	assert.Equal(t, message, compositeErr.Message)
}

func TestNewInvalidToolNameError(t *testing.T) {
	message := "Invalid tool name provided"
	err := compositekey.NewInvalidToolNameError(message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), compositekey.ErrInvalidToolName)
	assert.Contains(t, err.Error(), message)

	compositeErr, ok := err.(compositekey.CompositeKeyError)
	assert.True(t, ok)
	assert.Equal(t, compositekey.ErrInvalidToolName, compositeErr.Code)
	assert.Equal(t, message, compositeErr.Message)
}

func TestNewInvalidDomainNameError(t *testing.T) {
	message := "Invalid domain name provided"
	err := compositekey.NewInvalidDomainNameError(message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), compositekey.ErrInvalidDomainName)
	assert.Contains(t, err.Error(), message)

	compositeErr, ok := err.(compositekey.CompositeKeyError)
	assert.True(t, ok)
	assert.Equal(t, compositekey.ErrInvalidDomainName, compositeErr.Code)
	assert.Equal(t, message, compositeErr.Message)
}

func TestNewInvalidIDError(t *testing.T) {
	message := "Invalid ID provided"
	err := compositekey.NewInvalidIDError(message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), compositekey.ErrInvalidID)
	assert.Contains(t, err.Error(), message)

	compositeErr, ok := err.(compositekey.CompositeKeyError)
	assert.True(t, ok)
	assert.Equal(t, compositekey.ErrInvalidID, compositeErr.Code)
	assert.Equal(t, message, compositeErr.Message)
}

func TestNewTooLongError(t *testing.T) {
	message := "Composite key too long"
	err := compositekey.NewTooLongError(message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), compositekey.ErrTooLong)
	assert.Contains(t, err.Error(), message)

	compositeErr, ok := err.(compositekey.CompositeKeyError)
	assert.True(t, ok)
	assert.Equal(t, compositekey.ErrTooLong, compositeErr.Code)
	assert.Equal(t, message, compositeErr.Message)
}

func TestErrorConstants(t *testing.T) {
	// Test that error constants are properly defined
	assert.Equal(t, "COMPOSITE_KEY_INVALID_FORMAT", compositekey.ErrInvalidFormat)
	assert.Equal(t, "COMPOSITE_KEY_INVALID_TOOL_NAME", compositekey.ErrInvalidToolName)
	assert.Equal(t, "COMPOSITE_KEY_INVALID_DOMAIN_NAME", compositekey.ErrInvalidDomainName)
	assert.Equal(t, "COMPOSITE_KEY_INVALID_ID", compositekey.ErrInvalidID)
	assert.Equal(t, "COMPOSITE_KEY_TOO_LONG", compositekey.ErrTooLong)
}

func TestCompositeKeyErrorInterface(t *testing.T) {
	// Test that CompositeKeyError implements the error interface
	var err error = compositekey.CompositeKeyError{
		Code:    "TEST_CODE",
		Message: "Test message",
	}

	assert.Equal(t, "TEST_CODE: Test message", err.Error())
}

func TestMultipleErrorTypes(t *testing.T) {
	errors := []error{
		compositekey.NewInvalidFormatError("format error"),
		compositekey.NewInvalidToolNameError("tool name error"),
		compositekey.NewInvalidDomainNameError("domain name error"),
		compositekey.NewInvalidIDError("id error"),
		compositekey.NewTooLongError("too long error"),
	}

	expectedCodes := []string{
		compositekey.ErrInvalidFormat,
		compositekey.ErrInvalidToolName,
		compositekey.ErrInvalidDomainName,
		compositekey.ErrInvalidID,
		compositekey.ErrTooLong,
	}

	for i, err := range errors {
		compositeErr, ok := err.(compositekey.CompositeKeyError)
		assert.True(t, ok, "Error %d should be a CompositeKeyError", i)
		assert.Equal(t, expectedCodes[i], compositeErr.Code, "Error %d should have correct code", i)
		assert.NotEmpty(t, compositeErr.Message, "Error %d should have a message", i)
	}
}