package compositekey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name       string
		toolName   string
		domainName string
		id         int
		expected   string
	}{
		{
			name:       "기본 생성",
			toolName:   "url-db",
			domainName: "tech-articles",
			id:         123,
			expected:   "url-db:tech-articles:123",
		},
		{
			name:       "다른 도구명",
			toolName:   "bookmark-manager",
			domainName: "personal-links",
			id:         456,
			expected:   "bookmark-manager:personal-links:456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ck := Create(tt.toolName, tt.domainName, tt.id)
			assert.Equal(t, tt.expected, ck.String())
			assert.Equal(t, tt.toolName, ck.ToolName)
			assert.Equal(t, tt.domainName, ck.DomainName)
			assert.Equal(t, tt.id, ck.ID)
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name           string
		compositeKey   string
		expectedResult CompositeKey
		expectError    bool
	}{
		{
			name:         "유효한 합성키",
			compositeKey: "url-db:tech-articles:123",
			expectedResult: CompositeKey{
				ToolName:   "url-db",
				DomainName: "tech-articles",
				ID:         123,
			},
			expectError: false,
		},
		{
			name:         "다른 유효한 합성키",
			compositeKey: "bookmark-manager:personal-links:456",
			expectedResult: CompositeKey{
				ToolName:   "bookmark-manager",
				DomainName: "personal-links",
				ID:         456,
			},
			expectError: false,
		},
		{
			name:         "구성 요소 부족",
			compositeKey: "url-db:tech-articles",
			expectError:  true,
		},
		{
			name:         "구성 요소 초과",
			compositeKey: "url-db:tech-articles:123:extra",
			expectError:  true,
		},
		{
			name:         "잘못된 ID",
			compositeKey: "url-db:tech-articles:abc",
			expectError:  true,
		},
		{
			name:         "음수 ID",
			compositeKey: "url-db:tech-articles:-123",
			expectError:  true,
		},
		{
			name:         "0 ID",
			compositeKey: "url-db:tech-articles:0",
			expectError:  true,
		},
		{
			name:         "빈 문자열",
			compositeKey: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.compositeKey)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name         string
		compositeKey string
		expected     bool
	}{
		{
			name:         "유효한 합성키",
			compositeKey: "url-db:tech-articles:123",
			expected:     true,
		},
		{
			name:         "잘못된 형식",
			compositeKey: "url-db:tech-articles",
			expected:     false,
		},
		{
			name:         "잘못된 ID",
			compositeKey: "url-db:tech-articles:abc",
			expected:     false,
		},
		{
			name:         "빈 문자열",
			compositeKey: "",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValid(tt.compositeKey)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCompositeKeyString(t *testing.T) {
	ck := CompositeKey{
		ToolName:   "url-db",
		DomainName: "tech-articles",
		ID:         123,
	}

	expected := "url-db:tech-articles:123"
	assert.Equal(t, expected, ck.String())
}
