package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid HTTPS URL",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "valid HTTP URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid URL - no protocol",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "very long URL",
			url:     "https://example.com/" + string(make([]byte, 2100)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "valid title",
			title:   "Test Title",
			wantErr: false,
		},
		{
			name:    "empty title - should be valid",
			title:   "",
			wantErr: false,
		},
		{
			name:    "title too long",
			title:   string(make([]byte, 501)), // > 500 chars
			wantErr: true,
		},
		{
			name:    "normal length title",
			title:   string(make([]byte, 500)), // exactly 500 chars
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTitle(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAttributeName(t *testing.T) {
	tests := []struct {
		name      string
		attrName  string
		wantErr   bool
	}{
		{
			name:     "valid attribute name",
			attrName: "category",
			wantErr:  false,
		},
		{
			name:     "attribute name with underscore",
			attrName: "sub_category",
			wantErr:  false,
		},
		{
			name:     "attribute name with hyphen",
			attrName: "sub-category",
			wantErr:  false,
		},
		{
			name:     "empty attribute name",
			attrName: "",
			wantErr:  true,
		},
		{
			name:     "attribute name too long",
			attrName: string(make([]byte, 101)), // > 100 chars
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAttributeName(tt.attrName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAttributeType(t *testing.T) {
	tests := []struct {
		name     string
		attrType string
		wantErr  bool
	}{
		{
			name:     "valid type - tag",
			attrType: "tag",
			wantErr:  false,
		},
		{
			name:     "valid type - ordered_tag",
			attrType: "ordered_tag",
			wantErr:  false,
		},
		{
			name:     "valid type - number",
			attrType: "number",
			wantErr:  false,
		},
		{
			name:     "valid type - string",
			attrType: "string",
			wantErr:  false,
		},
		{
			name:     "valid type - markdown",
			attrType: "markdown",
			wantErr:  false,
		},
		{
			name:     "valid type - image",
			attrType: "image",
			wantErr:  false,
		},
		{
			name:     "invalid type",
			attrType: "invalid_type",
			wantErr:  true,
		},
		{
			name:     "empty type",
			attrType: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAttributeType(tt.attrType)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePositiveInteger(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		fieldName string
		wantErr   bool
	}{
		{
			name:      "positive integer",
			value:     10,
			fieldName: "test",
			wantErr:   false,
		},
		{
			name:      "zero",
			value:     0,
			fieldName: "test",
			wantErr:   true,
		},
		{
			name:      "negative integer",
			value:     -5,
			fieldName: "test",
			wantErr:   true,
		},
		{
			name:      "large positive integer",
			value:     1000000,
			fieldName: "test",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePositiveInteger(tt.value, tt.fieldName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePaginationParams(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		size     int
		wantErr  bool
		wantPage int
		wantSize int
	}{
		{
			name:     "valid pagination",
			page:     1,
			size:     10,
			wantErr:  false,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "zero page - should default to 1",
			page:     0,
			size:     10,
			wantErr:  false,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "zero size - should default to 20",
			page:     1,
			size:     0,
			wantErr:  false,
			wantPage: 1,
			wantSize: 20,
		},
		{
			name:     "size too large - should cap at 100",
			page:     1,
			size:     200,
			wantErr:  false,
			wantPage: 1,
			wantSize: 100,
		},
		{
			name:     "negative values - should use defaults",
			page:     -1,
			size:     -5,
			wantErr:  false,
			wantPage: 1,
			wantSize: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, size, err := validatePaginationParams(tt.page, tt.size)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPage, page)
				assert.Equal(t, tt.wantSize, size)
			}
		})
	}
}

func TestNormalizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "test string",
			expected: "test string",
		},
		{
			name:     "string with leading/trailing spaces",
			input:    "  test string  ",
			expected: "test string",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "string with tabs and newlines",
			input:    "\t test string \n",
			expected: "test string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateTitleFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "simple domain",
			url:      "https://example.com",
			expected: "Example.Com",
		},
		{
			name:     "domain with www",
			url:      "https://www.example.com",
			expected: "Www.Example.Com",
		},
		{
			name:     "URL with path",
			url:      "https://example.com/article",
			expected: "Article",
		},
		{
			name:     "URL with hyphenated path",
			url:      "https://example.com/my-article",
			expected: "My Article",
		},
		{
			name:     "empty URL",
			url:      "",
			expected: "Untitled",
		},
		{
			name:     "invalid URL",
			url:      "not-a-url",
			expected: "Untitled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateTitleFromURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}