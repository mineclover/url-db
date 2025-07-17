package compositekey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	service := NewService("url-db")
	assert.NotNil(t, service)
	assert.Equal(t, "url-db", service.defaultToolName)
}

func TestServiceCreate(t *testing.T) {
	service := NewService("url-db")
	
	tests := []struct {
		name        string
		domainName  string
		id          int
		expected    string
		expectError bool
	}{
		{
			name:       "기본 생성",
			domainName: "Tech Articles",
			id:         123,
			expected:   "url-db:tech-articles:123",
		},
		{
			name:       "특수문자 포함",
			domainName: "My@Special#Domain!",
			id:         456,
			expected:   "url-db:my-special-domain:456",
		},
		{
			name:        "잘못된 도메인명",
			domainName:  "",
			id:          123,
			expectError: true,
		},
		{
			name:        "잘못된 ID",
			domainName:  "valid-domain",
			id:          0,
			expectError: true,
		},
		{
			name:        "음수 ID",
			domainName:  "valid-domain",
			id:          -123,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Create(tt.domainName, tt.id)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServiceCreateWithTool(t *testing.T) {
	service := NewService("url-db")
	
	tests := []struct {
		name        string
		toolName    string
		domainName  string
		id          int
		expected    string
		expectError bool
	}{
		{
			name:       "기본 생성",
			toolName:   "Custom Tool",
			domainName: "Tech Articles",
			id:         123,
			expected:   "custom-tool:tech-articles:123",
		},
		{
			name:       "특수문자 포함",
			toolName:   "My@Tool!",
			domainName: "My@Domain!",
			id:         456,
			expected:   "my-tool:my-domain:456",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CreateWithTool(tt.toolName, tt.domainName, tt.id)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServiceParse(t *testing.T) {
	service := NewService("url-db")
	
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
		},
		{
			name:         "다른 유효한 합성키",
			compositeKey: "bookmark-manager:personal-links:456",
			expectedResult: CompositeKey{
				ToolName:   "bookmark-manager",
				DomainName: "personal-links",
				ID:         456,
			},
		},
		{
			name:         "잘못된 형식",
			compositeKey: "url-db:tech-articles",
			expectError:  true,
		},
		{
			name:         "잘못된 ID",
			compositeKey: "url-db:tech-articles:abc",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Parse(tt.compositeKey)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestServiceValidate(t *testing.T) {
	service := NewService("url-db")
	
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
			result := service.Validate(tt.compositeKey)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServiceParseComponents(t *testing.T) {
	service := NewService("url-db")
	
	tests := []struct {
		name               string
		compositeKey       string
		expectedToolName   string
		expectedDomainName string
		expectedID         int
		expectError        bool
	}{
		{
			name:               "유효한 합성키",
			compositeKey:       "url-db:tech-articles:123",
			expectedToolName:   "url-db",
			expectedDomainName: "tech-articles",
			expectedID:         123,
		},
		{
			name:               "다른 유효한 합성키",
			compositeKey:       "bookmark-manager:personal-links:456",
			expectedToolName:   "bookmark-manager",
			expectedDomainName: "personal-links",
			expectedID:         456,
		},
		{
			name:         "잘못된 형식",
			compositeKey: "url-db:tech-articles",
			expectError:  true,
		},
		{
			name:         "잘못된 ID",
			compositeKey: "url-db:tech-articles:abc",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolName, domainName, id, err := service.ParseComponents(tt.compositeKey)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expectedToolName, toolName)
			assert.Equal(t, tt.expectedDomainName, domainName)
			assert.Equal(t, tt.expectedID, id)
		})
	}
}

func TestServiceGetToolName(t *testing.T) {
	service := NewService("url-db")
	
	tests := []struct {
		name         string
		compositeKey string
		expected     string
		expectError  bool
	}{
		{
			name:         "유효한 합성키",
			compositeKey: "url-db:tech-articles:123",
			expected:     "url-db",
		},
		{
			name:         "다른 유효한 합성키",
			compositeKey: "bookmark-manager:personal-links:456",
			expected:     "bookmark-manager",
		},
		{
			name:         "잘못된 형식",
			compositeKey: "url-db:tech-articles",
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
			result, err := service.GetToolName(tt.compositeKey)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServiceGetDomainName(t *testing.T) {
	service := NewService("url-db")
	
	tests := []struct {
		name         string
		compositeKey string
		expected     string
		expectError  bool
	}{
		{
			name:         "유효한 합성키",
			compositeKey: "url-db:tech-articles:123",
			expected:     "tech-articles",
		},
		{
			name:         "다른 유효한 합성키",
			compositeKey: "bookmark-manager:personal-links:456",
			expected:     "personal-links",
		},
		{
			name:         "잘못된 형식",
			compositeKey: "url-db:tech-articles",
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
			result, err := service.GetDomainName(tt.compositeKey)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServiceGetID(t *testing.T) {
	service := NewService("url-db")
	
	tests := []struct {
		name         string
		compositeKey string
		expected     int
		expectError  bool
	}{
		{
			name:         "유효한 합성키",
			compositeKey: "url-db:tech-articles:123",
			expected:     123,
		},
		{
			name:         "다른 유효한 합성키",
			compositeKey: "bookmark-manager:personal-links:456",
			expected:     456,
		},
		{
			name:         "잘못된 형식",
			compositeKey: "url-db:tech-articles",
			expectError:  true,
		},
		{
			name:         "잘못된 ID",
			compositeKey: "url-db:tech-articles:abc",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetID(tt.compositeKey)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServiceIsValidFormat(t *testing.T) {
	service := NewService("url-db")
	
	tests := []struct {
		name         string
		compositeKey string
		expected     bool
	}{
		{
			name:         "유효한 형식",
			compositeKey: "url-db:tech-articles:123",
			expected:     true,
		},
		{
			name:         "잘못된 형식",
			compositeKey: "url-db:tech-articles",
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
			result := service.IsValidFormat(tt.compositeKey)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServiceNormalizeComponents(t *testing.T) {
	service := NewService("url-db")
	
	tests := []struct {
		name                   string
		toolName               string
		domainName             string
		expectedToolName       string
		expectedDomainName     string
		expectError            bool
	}{
		{
			name:               "기본 정규화",
			toolName:           "URL Database",
			domainName:         "Tech Articles",
			expectedToolName:   "url-database",
			expectedDomainName: "tech-articles",
		},
		{
			name:               "특수문자 처리",
			toolName:           "My@Tool!",
			domainName:         "My@Domain!",
			expectedToolName:   "my-tool",
			expectedDomainName: "my-domain",
		},
		{
			name:        "잘못된 도구명",
			toolName:    "",
			domainName:  "valid-domain",
			expectError: true,
		},
		{
			name:        "잘못된 도메인명",
			toolName:    "valid-tool",
			domainName:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolName, domainName, err := service.NormalizeComponents(tt.toolName, tt.domainName)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expectedToolName, toolName)
			assert.Equal(t, tt.expectedDomainName, domainName)
		})
	}
}

func TestServiceValidateComponents(t *testing.T) {
	service := NewService("url-db")
	
	tests := []struct {
		name        string
		toolName    string
		domainName  string
		id          int
		expectError bool
	}{
		{
			name:       "유효한 구성 요소",
			toolName:   "url-db",
			domainName: "tech-articles",
			id:         123,
		},
		{
			name:        "잘못된 도구명",
			toolName:    "url@db",
			domainName:  "tech-articles",
			id:          123,
			expectError: true,
		},
		{
			name:        "잘못된 도메인명",
			toolName:    "url-db",
			domainName:  "tech@articles",
			id:          123,
			expectError: true,
		},
		{
			name:        "잘못된 ID",
			toolName:    "url-db",
			domainName:  "tech-articles",
			id:          0,
			expectError: true,
		},
		{
			name:        "음수 ID",
			toolName:    "url-db",
			domainName:  "tech-articles",
			id:          -123,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateComponents(tt.toolName, tt.domainName, tt.id)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}