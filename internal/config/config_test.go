package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 테스트 전 환경변수 백업
	originalAutoCreate := os.Getenv("AUTO_CREATE_ATTRIBUTES")
	defer func() {
		if originalAutoCreate != "" {
			os.Setenv("AUTO_CREATE_ATTRIBUTES", originalAutoCreate)
		} else {
			os.Unsetenv("AUTO_CREATE_ATTRIBUTES")
		}
	}()

	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{
			name:     "기본값 true",
			envValue: "",
			expected: true,
		},
		{
			name:     "true 값",
			envValue: "true",
			expected: true,
		},
		{
			name:     "1 값",
			envValue: "1",
			expected: true,
		},
		{
			name:     "yes 값",
			envValue: "yes",
			expected: true,
		},
		{
			name:     "on 값",
			envValue: "on",
			expected: true,
		},
		{
			name:     "false 값",
			envValue: "false",
			expected: false,
		},
		{
			name:     "0 값",
			envValue: "0",
			expected: false,
		},
		{
			name:     "no 값",
			envValue: "no",
			expected: false,
		},
		{
			name:     "off 값",
			envValue: "off",
			expected: false,
		},
		{
			name:     "대문자 TRUE",
			envValue: "TRUE",
			expected: true,
		},
		{
			name:     "대문자 FALSE",
			envValue: "FALSE",
			expected: false,
		},
		{
			name:     "잘못된 값은 기본값 사용",
			envValue: "invalid",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 환경변수 설정
			if tt.envValue != "" {
				os.Setenv("AUTO_CREATE_ATTRIBUTES", tt.envValue)
			} else {
				os.Unsetenv("AUTO_CREATE_ATTRIBUTES")
			}

			// 설정 로드
			cfg := Load()

			// 결과 검증
			if cfg.AutoCreateAttributes != tt.expected {
				t.Errorf("AutoCreateAttributes = %v, want %v", cfg.AutoCreateAttributes, tt.expected)
			}
		})
	}
}

func TestGetBoolEnv(t *testing.T) {
	// 테스트 전 환경변수 백업
	originalValue := os.Getenv("TEST_BOOL_ENV")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_BOOL_ENV", originalValue)
		} else {
			os.Unsetenv("TEST_BOOL_ENV")
		}
	}()

	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "환경변수 없음, 기본값 true",
			envValue:     "",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "환경변수 없음, 기본값 false",
			envValue:     "",
			defaultValue: false,
			expected:     false,
		},
		{
			name:         "true 값",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "1 값",
			envValue:     "1",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "yes 값",
			envValue:     "yes",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "on 값",
			envValue:     "on",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "false 값",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "0 값",
			envValue:     "0",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "no 값",
			envValue:     "no",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "off 값",
			envValue:     "off",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "잘못된 값은 기본값 사용",
			envValue:     "invalid",
			defaultValue: true,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 환경변수 설정
			if tt.envValue != "" {
				os.Setenv("TEST_BOOL_ENV", tt.envValue)
			} else {
				os.Unsetenv("TEST_BOOL_ENV")
			}

			// 함수 호출
			result := getBoolEnv("TEST_BOOL_ENV", tt.defaultValue)

			// 결과 검증
			if result != tt.expected {
				t.Errorf("getBoolEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}
