package attributes

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagValidator_Validate(t *testing.T) {
	validator := &TagValidator{}

	tests := []struct {
		name       string
		value      string
		orderIndex *int
		wantErr    error
	}{
		{
			name:    "valid tag",
			value:   "test-tag",
			wantErr: nil,
		},
		{
			name:    "empty value",
			value:   "",
			wantErr: ErrValueRequired,
		},
		{
			name:    "too long value",
			value:   strings.Repeat("a", 256),
			wantErr: ErrValueTooLong,
		},
		{
			name:    "max length value",
			value:   strings.Repeat("a", 255),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value, tt.orderIndex)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderedTagValidator_Validate(t *testing.T) {
	validator := &OrderedTagValidator{}

	tests := []struct {
		name       string
		value      string
		orderIndex *int
		wantErr    error
	}{
		{
			name:       "valid ordered tag",
			value:      "test-tag",
			orderIndex: intPtr(1),
			wantErr:    nil,
		},
		{
			name:    "empty value",
			value:   "",
			wantErr: ErrValueRequired,
		},
		{
			name:    "missing order index",
			value:   "test-tag",
			wantErr: ErrOrderIndexRequired,
		},
		{
			name:    "too long value",
			value:   strings.Repeat("a", 256),
			wantErr: ErrValueTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value, tt.orderIndex)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNumberValidator_Validate(t *testing.T) {
	validator := &NumberValidator{}

	tests := []struct {
		name       string
		value      string
		orderIndex *int
		wantErr    error
	}{
		{
			name:    "valid integer",
			value:   "123",
			wantErr: nil,
		},
		{
			name:    "valid float",
			value:   "123.45",
			wantErr: nil,
		},
		{
			name:    "valid negative number",
			value:   "-123.45",
			wantErr: nil,
		},
		{
			name:    "empty value",
			value:   "",
			wantErr: ErrValueRequired,
		},
		{
			name:    "invalid number",
			value:   "not-a-number",
			wantErr: ErrInvalidNumber,
		},
		{
			name:    "scientific notation",
			value:   "1.23e10",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value, tt.orderIndex)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStringValidator_Validate(t *testing.T) {
	validator := &StringValidator{}

	tests := []struct {
		name       string
		value      string
		orderIndex *int
		wantErr    error
	}{
		{
			name:    "valid string",
			value:   "test string",
			wantErr: nil,
		},
		{
			name:    "empty value",
			value:   "",
			wantErr: ErrValueRequired,
		},
		{
			name:    "too long value",
			value:   strings.Repeat("a", 2049),
			wantErr: ErrValueTooLong,
		},
		{
			name:    "max length value",
			value:   strings.Repeat("a", 2048),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value, tt.orderIndex)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMarkdownValidator_Validate(t *testing.T) {
	validator := &MarkdownValidator{}

	tests := []struct {
		name       string
		value      string
		orderIndex *int
		wantErr    error
	}{
		{
			name:    "valid markdown",
			value:   "# Header\n\nThis is **bold** text and *italic* text.\n\n[Link](https://example.com)",
			wantErr: nil,
		},
		{
			name:    "empty value",
			value:   "",
			wantErr: ErrValueRequired,
		},
		{
			name:    "too long value",
			value:   strings.Repeat("a", 10001),
			wantErr: ErrValueTooLong,
		},
		{
			name:    "valid with code block",
			value:   "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```",
			wantErr: nil,
		},
		{
			name:    "valid with inline code",
			value:   "Use `fmt.Println()` to print.",
			wantErr: nil,
		},
		{
			name:    "valid with relative link",
			value:   "[Link](./relative/path)",
			wantErr: nil,
		},
		{
			name:    "valid with anchor",
			value:   "[Link](#anchor)",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value, tt.orderIndex)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestImageValidator_Validate(t *testing.T) {
	validator := &ImageValidator{}

	tests := []struct {
		name       string
		value      string
		orderIndex *int
		wantErr    error
	}{
		{
			name:    "valid image URL with extension",
			value:   "https://example.com/image.jpg",
			wantErr: nil,
		},
		{
			name:    "valid imgur URL",
			value:   "https://i.imgur.com/abc123.png",
			wantErr: nil,
		},
		{
			name:    "valid unsplash URL",
			value:   "https://images.unsplash.com/photo-123",
			wantErr: nil,
		},
		{
			name:    "valid GitHub avatar",
			value:   "https://avatars.githubusercontent.com/u/123",
			wantErr: nil,
		},
		{
			name:    "empty value",
			value:   "",
			wantErr: ErrValueRequired,
		},
		{
			name:    "too long value",
			value:   "https://example.com/" + strings.Repeat("a", 2030) + ".jpg",
			wantErr: ErrValueTooLong,
		},
		{
			name:    "invalid URL",
			value:   "not-a-url",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "relative URL",
			value:   "./image.jpg",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "non-HTTP scheme",
			value:   "ftp://example.com/image.jpg",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "valid HTTPS image",
			value:   "https://example.com/image.png",
			wantErr: nil,
		},
		{
			name:    "valid HTTP image",
			value:   "http://example.com/image.gif",
			wantErr: nil,
		},
		{
			name:    "valid webp image",
			value:   "https://example.com/image.webp",
			wantErr: nil,
		},
		{
			name:    "valid svg image",
			value:   "https://example.com/image.svg",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value, tt.orderIndex)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatorFactory_GetValidator(t *testing.T) {
	factory := NewValidatorFactory()

	tests := []struct {
		name     string
		attrType AttributeType
		wantType interface{}
		wantErr  bool
	}{
		{
			name:     "tag validator",
			attrType: AttributeTypeTag,
			wantType: &TagValidator{},
			wantErr:  false,
		},
		{
			name:     "ordered tag validator",
			attrType: AttributeTypeOrderedTag,
			wantType: &OrderedTagValidator{},
			wantErr:  false,
		},
		{
			name:     "number validator",
			attrType: AttributeTypeNumber,
			wantType: &NumberValidator{},
			wantErr:  false,
		},
		{
			name:     "string validator",
			attrType: AttributeTypeString,
			wantType: &StringValidator{},
			wantErr:  false,
		},
		{
			name:     "markdown validator",
			attrType: AttributeTypeMarkdown,
			wantType: &MarkdownValidator{},
			wantErr:  false,
		},
		{
			name:     "image validator",
			attrType: AttributeTypeImage,
			wantType: &ImageValidator{},
			wantErr:  false,
		},
		{
			name:     "unsupported type",
			attrType: "unsupported",
			wantType: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := factory.GetValidator(tt.attrType)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, validator)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, tt.wantType, validator)
			}
		})
	}
}

func TestValidateAttributeValue(t *testing.T) {
	tests := []struct {
		name       string
		attrType   AttributeType
		value      string
		orderIndex *int
		wantErr    error
	}{
		{
			name:     "valid tag",
			attrType: AttributeTypeTag,
			value:    "test-tag",
			wantErr:  nil,
		},
		{
			name:       "valid ordered tag",
			attrType:   AttributeTypeOrderedTag,
			value:      "test-tag",
			orderIndex: intPtr(1),
			wantErr:    nil,
		},
		{
			name:     "valid number",
			attrType: AttributeTypeNumber,
			value:    "123.45",
			wantErr:  nil,
		},
		{
			name:     "invalid number",
			attrType: AttributeTypeNumber,
			value:    "not-a-number",
			wantErr:  ErrInvalidNumber,
		},
		{
			name:     "valid string",
			attrType: AttributeTypeString,
			value:    "test string",
			wantErr:  nil,
		},
		{
			name:     "valid markdown",
			attrType: AttributeTypeMarkdown,
			value:    "# Header\n\nParagraph",
			wantErr:  nil,
		},
		{
			name:     "valid image",
			attrType: AttributeTypeImage,
			value:    "https://example.com/image.jpg",
			wantErr:  nil,
		},
		{
			name:     "unsupported type",
			attrType: "unsupported",
			value:    "value",
			wantErr:  nil, // Should return the factory error, not nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAttributeValue(tt.attrType, tt.value, tt.orderIndex)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else if tt.attrType == "unsupported" {
				assert.Error(t, err) // Should error for unsupported type
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
