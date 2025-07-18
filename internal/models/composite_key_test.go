package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompositeKey_Valid(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		domainName  string
		identifier  string
		expectedKey string
	}{
		{
			name:        "Standard composite key",
			toolName:    "url-db",
			domainName:  "example.com",
			identifier:  "123",
			expectedKey: "url-db:example.com:123",
		},
		{
			name:        "Long domain name",
			toolName:    "url-db",
			domainName:  "very-long-domain-name.example.com",
			identifier:  "456",
			expectedKey: "url-db:very-long-domain-name.example.com:456",
		},
		{
			name:        "Numeric identifier",
			toolName:    "url-db",
			domainName:  "test.org",
			identifier:  "789",
			expectedKey: "url-db:test.org:789",
		},
		{
			name:        "UUID identifier",
			toolName:    "url-db",
			domainName:  "api.service.com",
			identifier:  "550e8400-e29b-41d4-a716-446655440000",
			expectedKey: "url-db:api.service.com:550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := NewCompositeKey(tt.toolName, tt.domainName, tt.identifier)
			assert.Equal(t, tt.expectedKey, key.String())
			assert.Equal(t, tt.toolName, key.ToolName)
			assert.Equal(t, tt.domainName, key.DomainName)
			assert.Equal(t, tt.identifier, key.Identifier)
		})
	}
}

func TestCompositeKey_Parse(t *testing.T) {
	tests := []struct {
		name           string
		keyString      string
		expectedKey    *CompositeKey
		expectedError  bool
		errorMessage   string
	}{
		{
			name:      "Valid composite key",
			keyString: "url-db:example.com:123",
			expectedKey: &CompositeKey{
				ToolName:    "url-db",
				DomainName:  "example.com",
				Identifier:  "123",
			},
			expectedError: false,
		},
		{
			name:      "Valid key with complex identifier",
			keyString: "url-db:api.service.com:550e8400-e29b-41d4-a716-446655440000",
			expectedKey: &CompositeKey{
				ToolName:    "url-db",
				DomainName:  "api.service.com",
				Identifier:  "550e8400-e29b-41d4-a716-446655440000",
			},
			expectedError: false,
		},
		{
			name:          "Invalid format - too few parts",
			keyString:     "url-db:example.com",
			expectedKey:   nil,
			expectedError: true,
			errorMessage:  "invalid composite key format",
		},
		{
			name:          "Invalid format - too many parts",
			keyString:     "url-db:example.com:123:extra",
			expectedKey:   nil,
			expectedError: true,
			errorMessage:  "invalid composite key format",
		},
		{
			name:          "Empty key string",
			keyString:     "",
			expectedKey:   nil,
			expectedError: true,
			errorMessage:  "invalid composite key format",
		},
		{
			name:          "Missing tool name",
			keyString:     ":example.com:123",
			expectedKey:   nil,
			expectedError: true,
			errorMessage:  "tool name cannot be empty",
		},
		{
			name:          "Missing domain name",
			keyString:     "url-db::123",
			expectedKey:   nil,
			expectedError: true,
			errorMessage:  "domain name cannot be empty",
		},
		{
			name:          "Missing identifier",
			keyString:     "url-db:example.com:",
			expectedKey:   nil,
			expectedError: true,
			errorMessage:  "identifier cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := ParseCompositeKey(tt.keyString)
			
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, key)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, key)
				assert.Equal(t, tt.expectedKey.ToolName, key.ToolName)
				assert.Equal(t, tt.expectedKey.DomainName, key.DomainName)
				assert.Equal(t, tt.expectedKey.Identifier, key.Identifier)
			}
		})
	}
}

func TestCompositeKey_Validate(t *testing.T) {
	tests := []struct {
		name          string
		key           *CompositeKey
		expectedError bool
		errorMessage  string
	}{
		{
			name: "Valid composite key",
			key: &CompositeKey{
				ToolName:    "url-db",
				DomainName:  "example.com",
				Identifier:  "123",
			},
			expectedError: false,
		},
		{
			name: "Empty tool name",
			key: &CompositeKey{
				ToolName:    "",
				DomainName:  "example.com",
				Identifier:  "123",
			},
			expectedError: true,
			errorMessage:  "tool name cannot be empty",
		},
		{
			name: "Empty domain name",
			key: &CompositeKey{
				ToolName:    "url-db",
				DomainName:  "",
				Identifier:  "123",
			},
			expectedError: true,
			errorMessage:  "domain name cannot be empty",
		},
		{
			name: "Empty identifier",
			key: &CompositeKey{
				ToolName:    "url-db",
				DomainName:  "example.com",
				Identifier:  "",
			},
			expectedError: true,
			errorMessage:  "identifier cannot be empty",
		},
		{
			name: "Tool name with colon",
			key: &CompositeKey{
				ToolName:    "url:db",
				DomainName:  "example.com",
				Identifier:  "123",
			},
			expectedError: true,
			errorMessage:  "tool name cannot contain colon",
		},
		{
			name: "Domain name with colon",
			key: &CompositeKey{
				ToolName:    "url-db",
				DomainName:  "example:com",
				Identifier:  "123",
			},
			expectedError: true,
			errorMessage:  "domain name cannot contain colon",
		},
		{
			name: "Identifier with colon",
			key: &CompositeKey{
				ToolName:    "url-db",
				DomainName:  "example.com",
				Identifier:  "12:3",
			},
			expectedError: true,
			errorMessage:  "identifier cannot contain colon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.key.Validate()
			
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCompositeKey_Equal(t *testing.T) {
	key1 := &CompositeKey{
		ToolName:    "url-db",
		DomainName:  "example.com",
		Identifier:  "123",
	}
	
	key2 := &CompositeKey{
		ToolName:    "url-db",
		DomainName:  "example.com",
		Identifier:  "123",
	}
	
	key3 := &CompositeKey{
		ToolName:    "url-db",
		DomainName:  "example.com",
		Identifier:  "456",
	}
	
	key4 := &CompositeKey{
		ToolName:    "url-db",
		DomainName:  "different.com",
		Identifier:  "123",
	}
	
	key5 := &CompositeKey{
		ToolName:    "different-tool",
		DomainName:  "example.com",
		Identifier:  "123",
	}

	assert.True(t, key1.Equal(key2))
	assert.True(t, key2.Equal(key1))
	assert.False(t, key1.Equal(key3))
	assert.False(t, key1.Equal(key4))
	assert.False(t, key1.Equal(key5))
	assert.False(t, key1.Equal(nil))
}

func TestCompositeKey_RoundTrip(t *testing.T) {
	original := &CompositeKey{
		ToolName:    "url-db",
		DomainName:  "example.com",
		Identifier:  "123",
	}
	
	// Convert to string and back
	keyString := original.String()
	parsed, err := ParseCompositeKey(keyString)
	
	assert.NoError(t, err)
	assert.True(t, original.Equal(parsed))
}

func TestCompositeKey_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		domainName  string
		identifier  string
		shouldPanic bool
	}{
		{
			name:        "Very long parts",
			toolName:    "very-long-tool-name-that-might-cause-issues",
			domainName:  "very-long-domain-name.with.multiple.subdomains.example.com",
			identifier:  "very-long-identifier-with-lots-of-characters-and-dashes-and-numbers-123456789",
			shouldPanic: false,
		},
		{
			name:        "Special characters in identifier",
			toolName:    "url-db",
			domainName:  "example.com",
			identifier:  "id-with-special_chars.and@symbols#123",
			shouldPanic: false,
		},
		{
			name:        "Numeric domain",
			toolName:    "url-db",
			domainName:  "192.168.1.1",
			identifier:  "123",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					NewCompositeKey(tt.toolName, tt.domainName, tt.identifier)
				})
			} else {
				key := NewCompositeKey(tt.toolName, tt.domainName, tt.identifier)
				assert.NotNil(t, key)
				
				// Test round trip
				keyString := key.String()
				parsed, err := ParseCompositeKey(keyString)
				assert.NoError(t, err)
				assert.True(t, key.Equal(parsed))
			}
		})
	}
}