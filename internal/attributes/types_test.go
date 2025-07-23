package attributes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"url-db/internal/attributes"
	"url-db/internal/models"
)

func TestIsValidAttributeType(t *testing.T) {
	tests := []struct {
		name     string
		attrType attributes.AttributeType
		expected bool
	}{
		{
			name:     "valid tag type",
			attrType: attributes.AttributeTypeTag,
			expected: true,
		},
		{
			name:     "valid ordered tag type",
			attrType: attributes.AttributeTypeOrderedTag,
			expected: true,
		},
		{
			name:     "valid number type",
			attrType: attributes.AttributeTypeNumber,
			expected: true,
		},
		{
			name:     "valid string type",
			attrType: attributes.AttributeTypeString,
			expected: true,
		},
		{
			name:     "valid markdown type",
			attrType: attributes.AttributeTypeMarkdown,
			expected: true,
		},
		{
			name:     "valid image type",
			attrType: attributes.AttributeTypeImage,
			expected: true,
		},
		{
			name:     "invalid type",
			attrType: "invalid_type",
			expected: false,
		},
		{
			name:     "empty type",
			attrType: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := attributes.IsValidAttributeType(tt.attrType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetSupportedAttributeTypes(t *testing.T) {
	supportedTypes := attributes.GetSupportedAttributeTypes()

	assert.NotNil(t, supportedTypes)
	assert.Len(t, supportedTypes, 6)

	expectedTypes := []attributes.AttributeType{
		attributes.AttributeTypeTag,
		attributes.AttributeTypeOrderedTag,
		attributes.AttributeTypeNumber,
		attributes.AttributeTypeString,
		attributes.AttributeTypeMarkdown,
		attributes.AttributeTypeImage,
	}

	for _, expectedType := range expectedTypes {
		assert.Contains(t, supportedTypes, expectedType)
	}
}

func TestAttributeTypeConstants(t *testing.T) {
	// Test that our constants match the models constants
	assert.Equal(t, models.AttributeTypeTag, attributes.AttributeTypeTag)
	assert.Equal(t, models.AttributeTypeOrderedTag, attributes.AttributeTypeOrderedTag)
	assert.Equal(t, models.AttributeTypeNumber, attributes.AttributeTypeNumber)
	assert.Equal(t, models.AttributeTypeString, attributes.AttributeTypeString)
	assert.Equal(t, models.AttributeTypeMarkdown, attributes.AttributeTypeMarkdown)
	assert.Equal(t, models.AttributeTypeImage, attributes.AttributeTypeImage)
}

func TestAttributeTypeAlias(t *testing.T) {
	// Test that AttributeType is properly aliased to models.AttributeType
	var attrType attributes.AttributeType = "test"
	var modelType models.AttributeType = "test"
	
	assert.Equal(t, modelType, models.AttributeType(attrType))
}