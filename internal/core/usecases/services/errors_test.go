package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceError_Error(t *testing.T) {
	err := &ServiceError{
		Code:    "TEST_ERROR",
		Message: "This is a test error",
	}

	expected := "TEST_ERROR: This is a test error"
	assert.Equal(t, expected, err.Error())
}

func TestServiceError_ErrorWithoutMessage(t *testing.T) {
	err := &ServiceError{
		Code: "TEST_ERROR",
	}

	expected := "TEST_ERROR"
	assert.Equal(t, expected, err.Error())
}

func TestNewValidationError(t *testing.T) {
	field := "username"
	message := "Username is required"

	err := NewValidationError(field, message)

	assert.NotNil(t, err)
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Contains(t, err.Message, field)
	assert.Contains(t, err.Message, message)
}

func TestNewDomainNotFoundError(t *testing.T) {
	domainID := 123

	err := NewDomainNotFoundError(domainID)

	assert.NotNil(t, err)
	assert.Equal(t, "DOMAIN_NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "123")
	assert.Contains(t, err.Message, "not found")
}

func TestNewDomainAlreadyExistsError(t *testing.T) {
	domainName := "test-domain"

	err := NewDomainAlreadyExistsError(domainName)

	assert.NotNil(t, err)
	assert.Equal(t, "DOMAIN_ALREADY_EXISTS", err.Code)
	assert.Contains(t, err.Message, domainName)
	assert.Contains(t, err.Message, "already exists")
}

func TestNewBusinessLogicError(t *testing.T) {
	message := "Operation not allowed"

	err := NewBusinessLogicError(message)

	assert.NotNil(t, err)
	assert.Equal(t, "BUSINESS_LOGIC_ERROR", err.Code)
	assert.Equal(t, message, err.Message)
}