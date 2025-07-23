package compositekey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name         string
		compositeKey string
		expectError  bool
	}{
		{
			name:         "유효한 형식",
			compositeKey: "url-db:tech-articles:123",
			expectError:  false,
		},
		{
			name:         "빈 문자열",
			compositeKey: "",
			expectError:  true,
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
			name:         "단일 구성 요소",
			compositeKey: "url-db",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFormat(tt.compositeKey)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateToolName(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		expectError bool
	}{
		{
			name:        "유효한 도구명",
			toolName:    "url-db",
			expectError: false,
		},
		{
			name:        "언더스코어 포함",
			toolName:    "url_db",
			expectError: false,
		},
		{
			name:        "숫자 포함",
			toolName:    "url-db-v2",
			expectError: false,
		},
		{
			name:        "빈 문자열",
			toolName:    "",
			expectError: true,
		},
		{
			name:        "너무 긴 도구명",
			toolName:    "this-is-a-very-long-tool-name-that-exceeds-the-maximum-length-limit-of-fifty-characters",
			expectError: true,
		},
		{
			name:        "특수문자 포함",
			toolName:    "url@db",
			expectError: true,
		},
		{
			name:        "공백 포함",
			toolName:    "url db",
			expectError: true,
		},
		{
			name:        "하이픈으로 시작",
			toolName:    "-url-db",
			expectError: true,
		},
		{
			name:        "하이픈으로 끝남",
			toolName:    "url-db-",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToolName(tt.toolName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDomainName(t *testing.T) {
	tests := []struct {
		name        string
		domainName  string
		expectError bool
	}{
		{
			name:        "유효한 도메인명",
			domainName:  "tech-articles",
			expectError: false,
		},
		{
			name:        "언더스코어 포함",
			domainName:  "tech_articles",
			expectError: false,
		},
		{
			name:        "숫자 포함",
			domainName:  "tech-articles-2024",
			expectError: false,
		},
		{
			name:        "빈 문자열",
			domainName:  "",
			expectError: true,
		},
		{
			name:        "너무 긴 도메인명",
			domainName:  "this-is-a-very-long-domain-name-that-exceeds-the-maximum-length-limit-of-fifty-characters",
			expectError: true,
		},
		{
			name:        "특수문자 포함",
			domainName:  "tech@articles",
			expectError: true,
		},
		{
			name:        "공백 포함",
			domainName:  "tech articles",
			expectError: true,
		},
		{
			name:        "하이픈으로 시작",
			domainName:  "-tech-articles",
			expectError: true,
		},
		{
			name:        "하이픈으로 끝남",
			domainName:  "tech-articles-",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDomainName(tt.domainName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name        string
		idStr       string
		expectError bool
	}{
		{
			name:        "유효한 ID",
			idStr:       "123",
			expectError: false,
		},
		{
			name:        "큰 ID",
			idStr:       "999999",
			expectError: false,
		},
		{
			name:        "빈 문자열",
			idStr:       "",
			expectError: true,
		},
		{
			name:        "너무 긴 ID",
			idStr:       "123456789012345678901",
			expectError: true,
		},
		{
			name:        "문자 포함",
			idStr:       "abc",
			expectError: true,
		},
		{
			name:        "음수",
			idStr:       "-123",
			expectError: true,
		},
		{
			name:        "0",
			idStr:       "0",
			expectError: true,
		},
		{
			name:        "부동소수점",
			idStr:       "123.45",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.idStr)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCompositeKey(t *testing.T) {
	tests := []struct {
		name         string
		compositeKey string
		expectError  bool
	}{
		{
			name:         "유효한 합성키",
			compositeKey: "url-db:tech-articles:123",
			expectError:  false,
		},
		{
			name:         "다른 유효한 합성키",
			compositeKey: "bookmark_manager:personal_links:456",
			expectError:  false,
		},
		{
			name:         "잘못된 형식",
			compositeKey: "url-db:tech-articles",
			expectError:  true,
		},
		{
			name:         "잘못된 도구명",
			compositeKey: "url@db:tech-articles:123",
			expectError:  true,
		},
		{
			name:         "잘못된 도메인명",
			compositeKey: "url-db:tech@articles:123",
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
			err := ValidateCompositeKey(tt.compositeKey)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCompositeKeyStruct(t *testing.T) {
	tests := []struct {
		name        string
		ck          CompositeKey
		expectError bool
	}{
		{
			name: "유효한 구조체",
			ck: CompositeKey{
				ToolName:   "url-db",
				DomainName: "tech-articles",
				ID:         123,
			},
			expectError: false,
		},
		{
			name: "잘못된 도구명",
			ck: CompositeKey{
				ToolName:   "url@db",
				DomainName: "tech-articles",
				ID:         123,
			},
			expectError: true,
		},
		{
			name: "잘못된 도메인명",
			ck: CompositeKey{
				ToolName:   "url-db",
				DomainName: "tech@articles",
				ID:         123,
			},
			expectError: true,
		},
		{
			name: "잘못된 ID",
			ck: CompositeKey{
				ToolName:   "url-db",
				DomainName: "tech-articles",
				ID:         0,
			},
			expectError: true,
		},
		{
			name: "음수 ID",
			ck: CompositeKey{
				ToolName:   "url-db",
				DomainName: "tech-articles",
				ID:         -123,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCompositeKeyStruct(tt.ck)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
