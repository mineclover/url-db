package domains_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"url-db/internal/domains"
)

func TestDomainError_Error(t *testing.T) {
	err := &domains.DomainError{
		Code:    "TEST_CODE",
		Message: "Test error message",
		Err:     errors.New("underlying error"),
	}

	assert.Equal(t, "Test error message", err.Error())
}

func TestDomainError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := &domains.DomainError{
		Code:    "TEST_CODE",
		Message: "Test error message",
		Err:     underlyingErr,
	}

	assert.Equal(t, underlyingErr, err.Unwrap())
}

func TestNewDomainError(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := domains.NewDomainError("TEST_CODE", "Test message", underlyingErr)

	assert.NotNil(t, err)
	assert.Equal(t, "TEST_CODE", err.Code)
	assert.Equal(t, "Test message", err.Message)
	assert.Equal(t, underlyingErr, err.Err)
}

func TestNewValidationError(t *testing.T) {
	message := "Validation failed"
	err := domains.NewValidationError(message)

	assert.NotNil(t, err)
	assert.Equal(t, domains.ErrorCodeValidation, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, domains.ErrValidation, err.Err)
}

func TestNewDomainNotFoundError(t *testing.T) {
	id := 123
	err := domains.NewDomainNotFoundError(id)

	assert.NotNil(t, err)
	assert.Equal(t, domains.ErrorCodeDomainNotFound, err.Code)
	assert.Equal(t, "Domain not found", err.Message)
	assert.Equal(t, domains.ErrDomainNotFound, err.Err)
}

func TestNewDomainAlreadyExistsError(t *testing.T) {
	name := "test-domain"
	err := domains.NewDomainAlreadyExistsError(name)

	assert.NotNil(t, err)
	assert.Equal(t, domains.ErrorCodeDomainAlreadyExists, err.Code)
	assert.Equal(t, "Domain already exists", err.Message)
	assert.Equal(t, domains.ErrDomainAlreadyExists, err.Err)
}

func TestErrorConstants(t *testing.T) {
	// Test that error constants are properly defined
	assert.Equal(t, "DOMAIN_NOT_FOUND", domains.ErrorCodeDomainNotFound)
	assert.Equal(t, "DOMAIN_ALREADY_EXISTS", domains.ErrorCodeDomainAlreadyExists)
	assert.Equal(t, "DOMAIN_NAME_INVALID", domains.ErrorCodeDomainNameInvalid)
	assert.Equal(t, "DOMAIN_HAS_DEPENDENCIES", domains.ErrorCodeDomainHasDependencies)
	assert.Equal(t, "VALIDATION_ERROR", domains.ErrorCodeValidation)
}

func TestPreDefinedErrors(t *testing.T) {
	// Test that predefined errors are not nil
	assert.NotNil(t, domains.ErrDomainNotFound)
	assert.NotNil(t, domains.ErrDomainAlreadyExists)
	assert.NotNil(t, domains.ErrDomainNameInvalid)
	assert.NotNil(t, domains.ErrDomainHasDependencies)
	assert.NotNil(t, domains.ErrValidation)

	// Test error messages
	assert.Equal(t, "domain not found", domains.ErrDomainNotFound.Error())
	assert.Equal(t, "domain already exists", domains.ErrDomainAlreadyExists.Error())
	assert.Equal(t, "domain name is invalid", domains.ErrDomainNameInvalid.Error())
	assert.Equal(t, "domain has dependencies and cannot be deleted", domains.ErrDomainHasDependencies.Error())
	assert.Equal(t, "validation error", domains.ErrValidation.Error())
}

func TestDomainErrorImplementsError(t *testing.T) {
	// Test that DomainError implements the error interface
	var err error = &domains.DomainError{
		Code:    "TEST_CODE",
		Message: "Test message",
		Err:     errors.New("underlying"),
	}

	assert.Equal(t, "Test message", err.Error())
}

func TestDomainErrorWrapping(t *testing.T) {
	// Test error wrapping functionality
	originalErr := errors.New("original error")
	domainErr := domains.NewDomainError("TEST_CODE", "Domain error", originalErr)

	// Test that errors.Is works
	assert.True(t, errors.Is(domainErr, originalErr))

	// Test that errors.Unwrap works
	unwrapped := errors.Unwrap(domainErr)
	assert.Equal(t, originalErr, unwrapped)
}

func TestNilUnderlyingError(t *testing.T) {
	err := domains.NewDomainError("TEST_CODE", "Test message", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "TEST_CODE", err.Code)
	assert.Equal(t, "Test message", err.Message)
	assert.Nil(t, err.Err)
	assert.Nil(t, err.Unwrap())
}