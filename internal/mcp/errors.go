package mcp

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCompositeKey = errors.New("invalid composite key format")
	ErrDomainNotFound      = errors.New("domain not found")
	ErrResourceNotFound    = errors.New("resource not found")
	ErrAccessDenied        = errors.New("access denied")
	ErrBatchPartialFailure = errors.New("batch partial failure")
)

type MCPError struct {
	Code     string      `json:"error"`
	Message  string      `json:"message"`
	Details  interface{} `json:"details,omitempty"`
	HTTPCode int         `json:"-"`
}

func (e MCPError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewInvalidCompositeKeyError(provided string) *MCPError {
	return &MCPError{
		Code:     "INVALID_COMPOSITE_KEY",
		Message:  "합성키 형식이 올바르지 않습니다",
		HTTPCode: 400,
		Details: map[string]interface{}{
			"expected_format": "tool_name:domain_name:id",
			"provided":        provided,
		},
	}
}

func NewDomainNotFoundError(domainName string) *MCPError {
	return &MCPError{
		Code:     "DOMAIN_NOT_FOUND",
		Message:  "지정된 도메인을 찾을 수 없습니다",
		HTTPCode: 404,
		Details: map[string]interface{}{
			"domain_name": domainName,
		},
	}
}

func NewResourceNotFoundError(compositeID string) *MCPError {
	return &MCPError{
		Code:     "RESOURCE_NOT_FOUND",
		Message:  "리소스를 찾을 수 없습니다",
		HTTPCode: 404,
		Details: map[string]interface{}{
			"composite_id": compositeID,
		},
	}
}

func NewNodeNotFoundError(domainName, url string) *MCPError {
	return &MCPError{
		Code:     "NODE_NOT_FOUND",
		Message:  "노드를 찾을 수 없습니다",
		HTTPCode: 404,
		Details: map[string]interface{}{
			"domain_name": domainName,
			"url":         url,
		},
	}
}

func NewAccessDeniedError(resource string) *MCPError {
	return &MCPError{
		Code:     "ACCESS_DENIED",
		Message:  "해당 리소스에 대한 접근 권한이 없습니다",
		HTTPCode: 403,
		Details: map[string]interface{}{
			"resource": resource,
		},
	}
}

func NewBatchPartialFailureError(failed []string) *MCPError {
	return &MCPError{
		Code:     "BATCH_PARTIAL_FAILURE",
		Message:  "배치 처리 중 일부 항목에서 오류가 발생했습니다",
		HTTPCode: 207,
		Details: map[string]interface{}{
			"failed_items": failed,
		},
	}
}

func NewValidationError(message string) *MCPError {
	return &MCPError{
		Code:     "VALIDATION_ERROR",
		Message:  message,
		HTTPCode: 400,
	}
}

func NewInternalServerError(message string) *MCPError {
	return &MCPError{
		Code:     "INTERNAL_SERVER_ERROR",
		Message:  message,
		HTTPCode: 500,
	}
}