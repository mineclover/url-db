package nodes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTitleFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Simple domain",
			url:      "https://example.com",
			expected: "example.com",
		},
		{
			name:     "Domain with www",
			url:      "https://www.example.com",
			expected: "example.com",
		},
		{
			name:     "URL with path",
			url:      "https://example.com/article",
			expected: "Article - example.com",
		},
		{
			name:     "URL with dashed path",
			url:      "https://example.com/my-awesome-article",
			expected: "My Awesome Article - example.com",
		},
		{
			name:     "URL with underscored path",
			url:      "https://example.com/my_awesome_article",
			expected: "My Awesome Article - example.com",
		},
		{
			name:     "URL with file extension",
			url:      "https://example.com/article.html",
			expected: "Article - example.com",
		},
		{
			name:     "URL with nested path",
			url:      "https://example.com/blog/2023/my-article",
			expected: "My Article - example.com",
		},
		{
			name:     "URL with trailing slash",
			url:      "https://example.com/article/",
			expected: "example.com",
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: "",
		},
		{
			name:     "Invalid URL",
			url:      "not-a-url",
			expected: "not-a-url",
		},
		{
			name:     "URL without protocol",
			url:      "example.com/article",
			expected: "example.com/article",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateTitleFromURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractTitleFromRawURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "URL with path",
			url:      "https://example.com/article",
			expected: "article - example.com",
		},
		{
			name:     "URL with www",
			url:      "https://www.example.com/article",
			expected: "article - example.com",
		},
		{
			name:     "URL without path",
			url:      "https://example.com",
			expected: "example.com",
		},
		{
			name:     "Invalid format",
			url:      "not-a-url",
			expected: "not-a-url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTitleFromRawURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expectErr bool
	}{
		{
			name:      "Valid HTTP URL",
			url:       "https://example.com",
			expectErr: false,
		},
		{
			name:      "Valid HTTPS URL",
			url:       "http://example.com",
			expectErr: false,
		},
		{
			name:      "Protocol-relative URL",
			url:       "//example.com",
			expectErr: false,
		},
		{
			name:      "Empty URL",
			url:       "",
			expectErr: true,
		},
		{
			name:      "URL without protocol",
			url:       "example.com",
			expectErr: true,
		},
		{
			name:      "Very long URL",
			url:       "https://example.com/" + string(make([]byte, 2050)),
			expectErr: true,
		},
		{
			name:      "URL with custom protocol",
			url:       "ftp://example.com/file.txt",
			expectErr: false,
		},
		{
			name:      "URL with port",
			url:       "https://example.com:8080/path",
			expectErr: false,
		},
		{
			name:      "URL with query params",
			url:       "https://example.com/search?q=test&page=1",
			expectErr: false,
		},
		{
			name:      "URL with fragment",
			url:       "https://example.com/page#section",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, ErrNodeURLInvalid, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
