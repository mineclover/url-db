package compositekey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeToolName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:     "기본 정규화",
			input:    "URL Database",
			expected: "url-database",
		},
		{
			name:     "특수문자 변환",
			input:    "My@Special#Tool!",
			expected: "my-special-tool",
		},
		{
			name:     "연속 구분자 처리",
			input:    "my--tool___name",
			expected: "my-tool-name",
		},
		{
			name:     "앞뒤 공백 제거",
			input:    "  my tool  ",
			expected: "my-tool",
		},
		{
			name:     "카멜케이스 변환",
			input:    "TechArticles",
			expected: "techarticles",
		},
		{
			name:     "이미 정규화된 문자열",
			input:    "url-db",
			expected: "url-db",
		},
		{
			name:        "빈 문자열",
			input:       "",
			expectError: true,
		},
		{
			name:        "공백만 있는 문자열",
			input:       "   ",
			expectError: true,
		},
		{
			name:        "너무 긴 문자열",
			input:       "this-is-a-very-long-tool-name-that-exceeds-the-maximum-length-limit-of-fifty-characters",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeToolName(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeDomainName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:     "기본 정규화",
			input:    "Tech Articles",
			expected: "tech-articles",
		},
		{
			name:     "특수문자 변환",
			input:    "My@Domain#Name!",
			expected: "my-domain-name",
		},
		{
			name:     "연속 구분자 처리",
			input:    "my--domain___name",
			expected: "my-domain-name",
		},
		{
			name:     "앞뒤 공백 제거",
			input:    "  my domain  ",
			expected: "my-domain",
		},
		{
			name:     "카멜케이스 변환",
			input:    "PersonalBookmarks",
			expected: "personalbookmarks",
		},
		{
			name:     "이미 정규화된 문자열",
			input:    "tech-articles",
			expected: "tech-articles",
		},
		{
			name:        "빈 문자열",
			input:       "",
			expectError: true,
		},
		{
			name:        "공백만 있는 문자열",
			input:       "   ",
			expectError: true,
		},
		{
			name:        "너무 긴 문자열",
			input:       "this-is-a-very-long-domain-name-that-exceeds-the-maximum-length-limit-of-fifty-characters",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeDomainName(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateNormalized(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		domainName  string
		id          int
		expected    CompositeKey
		expectError bool
	}{
		{
			name:       "기본 정규화된 생성",
			toolName:   "URL Database",
			domainName: "Tech Articles",
			id:         123,
			expected: CompositeKey{
				ToolName:   "url-database",
				DomainName: "tech-articles",
				ID:         123,
			},
		},
		{
			name:       "특수문자 포함",
			toolName:   "My@Tool!",
			domainName: "Special#Domain",
			id:         456,
			expected: CompositeKey{
				ToolName:   "my-tool",
				DomainName: "special-domain",
				ID:         456,
			},
		},
		{
			name:        "잘못된 도구명",
			toolName:    "",
			domainName:  "valid-domain",
			id:          123,
			expectError: true,
		},
		{
			name:        "잘못된 도메인명",
			toolName:    "valid-tool",
			domainName:  "",
			id:          123,
			expectError: true,
		},
		{
			name:        "잘못된 ID",
			toolName:    "valid-tool",
			domainName:  "valid-domain",
			id:          0,
			expectError: true,
		},
		{
			name:        "음수 ID",
			toolName:    "valid-tool",
			domainName:  "valid-domain",
			id:          -123,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CreateNormalized(tt.toolName, tt.domainName, tt.id)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
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
			name:     "기본 정규화",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "특수문자 변환",
			input:    "hello@world#test!",
			expected: "hello-world-test",
		},
		{
			name:     "연속 구분자",
			input:    "hello--world___test",
			expected: "hello-world-test",
		},
		{
			name:     "앞뒤 공백과 하이픈",
			input:    "  -hello-world-  ",
			expected: "hello-world",
		},
		{
			name:     "빈 문자열",
			input:    "",
			expected: "",
		},
		{
			name:     "공백만",
			input:    "   ",
			expected: "",
		},
		{
			name:     "하이픈만",
			input:    "---",
			expected: "",
		},
		{
			name:     "대소문자 변환",
			input:    "HelloWorld",
			expected: "helloworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}