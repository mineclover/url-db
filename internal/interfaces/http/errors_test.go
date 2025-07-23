package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleError_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		err := NewValidationError("Invalid input", map[string]string{"field": "required"})
		HandleError(c, err)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
	assert.Equal(t, "Invalid input", response.Message)
	assert.NotNil(t, response.Details)
}

func TestHandleError_NotFoundError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		err := NewNotFoundError("Resource not found")
		HandleError(c, err)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", response.Error)
	assert.Equal(t, "Resource not found", response.Message)
}

func TestHandleError_ConflictError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		err := NewConflictError("Resource already exists")
		HandleError(c, err)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "conflict", response.Error)
	assert.Equal(t, "Resource already exists", response.Message)
}

func TestHandleError_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		err := NewInternalError("Internal server error")
		HandleError(c, err)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
	assert.Equal(t, "Internal server error", response.Message)
}

func TestHandleError_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		err := assert.AnError // Use a standard Go error
		HandleError(c, err)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
	assert.Equal(t, "An unexpected error occurred", response.Message)
}

func TestErrorTypes_ErrorMethod(t *testing.T) {
	validationErr := NewValidationError("validation failed", nil)
	assert.Equal(t, "validation failed", validationErr.Error())

	notFoundErr := NewNotFoundError("not found")
	assert.Equal(t, "not found", notFoundErr.Error())

	conflictErr := NewConflictError("conflict")
	assert.Equal(t, "conflict", conflictErr.Error())

	internalErr := NewInternalError("internal error")
	assert.Equal(t, "internal error", internalErr.Error())
}
