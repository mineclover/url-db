package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseJSONRPCRequest(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    *JSONRPCRequest
	}{
		{
			name:  "valid request",
			input: `{"jsonrpc":"2.0","id":1,"method":"test","params":{"key":"value"}}`,
			expected: &JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      float64(1), // JSON unmarshals to float64
				Method:  "test",
				Params:  map[string]interface{}{"key": "value"},
			},
		},
		{
			name:  "request without params",
			input: `{"jsonrpc":"2.0","id":1,"method":"test"}`,
			expected: &JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      float64(1),
				Method:  "test",
				Params:  nil,
			},
		},
		{
			name:        "invalid JSON",
			input:       `{"invalid json}`,
			expectError: true,
		},
		{
			name:        "wrong jsonrpc version",
			input:       `{"jsonrpc":"1.0","id":1,"method":"test"}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSONRPCRequest([]byte(tt.input))
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.JSONRPC, result.JSONRPC)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Method, result.Method)
				assert.Equal(t, tt.expected.Params, result.Params)
			}
		})
	}
}

func TestJSONRPCResponse_ToJSON(t *testing.T) {
	tests := []struct {
		name     string
		response *JSONRPCResponse
		expected string
	}{
		{
			name: "success response",
			response: &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Result:  "success",
			},
			expected: `{"jsonrpc":"2.0","id":1,"result":"success"}`,
		},
		{
			name: "error response",
			response: &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Error: &JSONRPCError{
					Code:    -32600,
					Message: "Invalid Request",
				},
			},
			expected: `{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"Invalid Request"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.response.ToJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestNewJSONRPCRequest(t *testing.T) {
	req := NewJSONRPCRequest(1, "test_method", map[string]string{"key": "value"})
	
	assert.Equal(t, "2.0", req.JSONRPC)
	assert.Equal(t, 1, req.ID)
	assert.Equal(t, "test_method", req.Method)
	assert.Equal(t, map[string]string{"key": "value"}, req.Params)
}

func TestNewJSONRPCResponse(t *testing.T) {
	resp := NewJSONRPCResponse(1, "test_result")
	
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "test_result", resp.Result)
	assert.Nil(t, resp.Error)
}

func TestNewJSONRPCError(t *testing.T) {
	resp := NewJSONRPCError(1, ParseError, "Parse error", "additional data")
	
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Equal(t, 1, resp.ID)
	assert.Nil(t, resp.Result)
	require.NotNil(t, resp.Error)
	assert.Equal(t, ParseError, resp.Error.Code)
	assert.Equal(t, "Parse error", resp.Error.Message)
	assert.Equal(t, "additional data", resp.Error.Data)
}