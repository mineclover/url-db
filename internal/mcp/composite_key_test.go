package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"url-db/internal/compositekey"
)

func TestCompositeKeyAdapter_Create(t *testing.T) {
	service := compositekey.NewService("url-db")
	adapter := NewCompositeKeyAdapter(service)

	result := adapter.Create("test-domain", 123)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "url-db")
	assert.Contains(t, result, "test-domain")
	assert.Contains(t, result, "123")
}

func TestCompositeKeyAdapter_Parse(t *testing.T) {
	service := compositekey.NewService("url-db")
	adapter := NewCompositeKeyAdapter(service)

	compositeID := "url-db:test-domain:123"
	result, err := adapter.Parse(compositeID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-domain", result.DomainName)
	assert.Equal(t, 123, result.ID)
}

func TestCompositeKeyAdapter_Parse_InvalidFormat(t *testing.T) {
	service := compositekey.NewService("url-db")
	adapter := NewCompositeKeyAdapter(service)

	compositeID := "invalid-format"
	result, err := adapter.Parse(compositeID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCompositeKeyAdapter_Validate(t *testing.T) {
	service := compositekey.NewService("url-db")
	adapter := NewCompositeKeyAdapter(service)

	tests := []struct {
		name        string
		compositeID string
		expectError bool
	}{
		{
			name:        "valid composite ID",
			compositeID: "url-db:test-domain:123",
			expectError: false,
		},
		{
			name:        "invalid format",
			compositeID: "invalid-format",
			expectError: true,
		},
		{
			name:        "empty composite ID",
			compositeID: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.Validate(tt.compositeID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCompositeKeyAdapter_ExtractDomainName(t *testing.T) {
	service := compositekey.NewService("url-db")
	adapter := NewCompositeKeyAdapter(service)

	compositeID := "url-db:test-domain:123"
	result, err := adapter.ExtractDomainName(compositeID)

	assert.NoError(t, err)
	assert.Equal(t, "test-domain", result)
}

func TestCompositeKeyAdapter_ExtractNodeID(t *testing.T) {
	service := compositekey.NewService("url-db")
	adapter := NewCompositeKeyAdapter(service)

	compositeID := "url-db:test-domain:123"
	result, err := adapter.ExtractNodeID(compositeID)

	assert.NoError(t, err)
	assert.Equal(t, 123, result)
}

func TestCompositeKeyAdapter_ExtractToolName(t *testing.T) {
	service := compositekey.NewService("url-db")
	adapter := NewCompositeKeyAdapter(service)

	compositeID := "url-db:test-domain:123"
	result, err := adapter.ExtractToolName(compositeID)

	assert.NoError(t, err)
	assert.Equal(t, "url-db", result)
}