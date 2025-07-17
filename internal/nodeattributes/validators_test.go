package nodeattributes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"internal/models"
)

func TestValidator_ValidateValue(t *testing.T) {
	v := NewValidator()

	t.Run("tag validation", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			wantErr bool
		}{
			{"valid tag", "programming", false},
			{"valid tag with hyphens", "web-development", false},
			{"valid tag with underscores", "machine_learning", false},
			{"valid tag with numbers", "python3", false},
			{"empty tag", "", true},
			{"tag too long", string(make([]byte, 256)), true},
			{"tag with spaces", "web development", true},
			{"tag with special chars", "web@development", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidateValue(models.AttributeTypeTag, tt.value)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("ordered_tag validation", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			wantErr bool
		}{
			{"valid ordered tag", "priority-high", false},
			{"empty ordered tag", "", true},
			{"ordered tag too long", string(make([]byte, 256)), true},
			{"ordered tag with spaces", "priority high", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidateValue(models.AttributeTypeOrderedTag, tt.value)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("number validation", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			wantErr bool
		}{
			{"valid integer", "123", false},
			{"valid float", "123.45", false},
			{"valid negative", "-123", false},
			{"valid zero", "0", false},
			{"empty number", "", true},
			{"invalid number", "abc", true},
			{"number with spaces", "123 ", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidateValue(models.AttributeTypeNumber, tt.value)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("string validation", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			wantErr bool
		}{
			{"valid string", "This is a valid string", false},
			{"empty string", "", true},
			{"string too long", string(make([]byte, 2049)), true},
			{"string at max length", string(make([]byte, 2048)), false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidateValue(models.AttributeTypeString, tt.value)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("markdown validation", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			wantErr bool
		}{
			{"valid markdown", "# Title\n\nThis is **bold** text", false},
			{"empty markdown", "", true},
			{"whitespace only", "   ", true},
			{"markdown too long", string(make([]byte, 10001)), true},
			{"markdown at max length", string(make([]byte, 10000)), false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidateValue(models.AttributeTypeMarkdown, tt.value)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("image validation", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			wantErr bool
		}{
			{"valid jpg", "https://example.com/image.jpg", false},
			{"valid png", "https://example.com/image.png", false},
			{"valid gif", "https://example.com/image.gif", false},
			{"valid webp", "https://example.com/image.webp", false},
			{"valid svg", "https://example.com/image.svg", false},
			{"http url", "http://example.com/image.jpg", false},
			{"empty image", "", true},
			{"invalid url", "not-a-url", true},
			{"invalid extension", "https://example.com/image.txt", true},
			{"no extension", "https://example.com/image", true},
			{"ftp scheme", "ftp://example.com/image.jpg", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidateValue(models.AttributeTypeImage, tt.value)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestValidator_ValidateOrderIndex(t *testing.T) {
	v := NewValidator()

	t.Run("ordered_tag requires order index", func(t *testing.T) {
		err := v.ValidateOrderIndex(models.AttributeTypeOrderedTag, nil)
		assert.Error(t, err)
		assert.Equal(t, ErrOrderIndexRequired, err)

		validIndex := 1
		err = v.ValidateOrderIndex(models.AttributeTypeOrderedTag, &validIndex)
		assert.NoError(t, err)

		invalidIndex := -1
		err = v.ValidateOrderIndex(models.AttributeTypeOrderedTag, &invalidIndex)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidOrderIndex, err)
	})

	t.Run("other types don't allow order index", func(t *testing.T) {
		orderIndex := 1
		
		types := []models.AttributeType{
			models.AttributeTypeTag,
			models.AttributeTypeNumber,
			models.AttributeTypeString,
			models.AttributeTypeMarkdown,
			models.AttributeTypeImage,
		}

		for _, attrType := range types {
			t.Run(string(attrType), func(t *testing.T) {
				err := v.ValidateOrderIndex(attrType, &orderIndex)
				assert.Error(t, err)
				assert.Equal(t, ErrOrderIndexNotAllowed, err)

				err = v.ValidateOrderIndex(attrType, nil)
				assert.NoError(t, err)
			})
		}
	})
}

func TestValidator_Validate(t *testing.T) {
	v := NewValidator()

	t.Run("valid tag without order index", func(t *testing.T) {
		err := v.Validate(models.AttributeTypeTag, "programming", nil)
		assert.NoError(t, err)
	})

	t.Run("valid ordered tag with order index", func(t *testing.T) {
		orderIndex := 1
		err := v.Validate(models.AttributeTypeOrderedTag, "priority-high", &orderIndex)
		assert.NoError(t, err)
	})

	t.Run("invalid tag with order index", func(t *testing.T) {
		orderIndex := 1
		err := v.Validate(models.AttributeTypeTag, "programming", &orderIndex)
		assert.Error(t, err)
		assert.Equal(t, ErrOrderIndexNotAllowed, err)
	})

	t.Run("invalid ordered tag without order index", func(t *testing.T) {
		err := v.Validate(models.AttributeTypeOrderedTag, "priority-high", nil)
		assert.Error(t, err)
		assert.Equal(t, ErrOrderIndexRequired, err)
	})
}