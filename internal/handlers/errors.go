package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type ValidationError struct {
	Message string
	Details interface{}
}

func (e *ValidationError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return e.Message
}

func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *ValidationError:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: e.Message,
			Details: e.Details,
		})
	case *NotFoundError:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: e.Message,
		})
	case *ConflictError:
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "conflict",
			Message: e.Message,
		})
	case *InternalError:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: e.Message,
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "An unexpected error occurred",
		})
	}
}

func NewValidationError(message string, details interface{}) *ValidationError {
	return &ValidationError{
		Message: message,
		Details: details,
	}
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{
		Message: message,
	}
}

func NewInternalError(message string) *InternalError {
	return &InternalError{
		Message: message,
	}
}