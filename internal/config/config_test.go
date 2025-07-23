package config_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"url-db/internal/config"
	"url-db/internal/constants"
)

func TestLoad(t *testing.T) {
	// Save original environment
	originalPort := os.Getenv("PORT")
	originalDatabaseURL := os.Getenv("DATABASE_URL")
	originalToolName := os.Getenv("TOOL_NAME")

	// Clean up after test
	defer func() {
		os.Setenv("PORT", originalPort)
		os.Setenv("DATABASE_URL", originalDatabaseURL)
		os.Setenv("TOOL_NAME", originalToolName)
	}()

	tests := []struct {
		name        string
		setupEnv    func()
		expectedCfg *config.Config
	}{
		{
			name: "default values",
			setupEnv: func() {
				os.Unsetenv("PORT")
				os.Unsetenv("DATABASE_URL")
				os.Unsetenv("TOOL_NAME")
			},
			expectedCfg: &config.Config{
				Port:        strconv.Itoa(constants.DefaultPort),
				DatabaseURL: "file:./" + constants.DefaultDBPath,
				ToolName:    constants.DefaultServerName,
			},
		},
		{
			name: "custom environment values",
			setupEnv: func() {
				os.Setenv("PORT", "9090")
				os.Setenv("DATABASE_URL", "file:./custom.db")
				os.Setenv("TOOL_NAME", "custom-tool")
			},
			expectedCfg: &config.Config{
				Port:        "9090",
				DatabaseURL: "file:./custom.db",
				ToolName:    "custom-tool",
			},
		},
		{
			name: "partial environment values",
			setupEnv: func() {
				os.Setenv("PORT", "8000")
				os.Unsetenv("DATABASE_URL")
				os.Setenv("TOOL_NAME", "partial-tool")
			},
			expectedCfg: &config.Config{
				Port:        "8000",
				DatabaseURL: "file:./" + constants.DefaultDBPath,
				ToolName:    "partial-tool",
			},
		},
		{
			name: "empty environment values should use defaults",
			setupEnv: func() {
				os.Setenv("PORT", "")
				os.Setenv("DATABASE_URL", "")
				os.Setenv("TOOL_NAME", "")
			},
			expectedCfg: &config.Config{
				Port:        strconv.Itoa(constants.DefaultPort),
				DatabaseURL: "file:./" + constants.DefaultDBPath,
				ToolName:    constants.DefaultServerName,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			
			cfg := config.Load()
			
			assert.Equal(t, tt.expectedCfg.Port, cfg.Port)
			assert.Equal(t, tt.expectedCfg.DatabaseURL, cfg.DatabaseURL)
			assert.Equal(t, tt.expectedCfg.ToolName, cfg.ToolName)
		})
	}
}

// TestGetEnv를 제거하고 Load()를 통해 간접적으로 테스트
// getEnv는 private 함수이므로 별도 패키지에서 직접 테스트할 수 없음

func TestConfig_Struct(t *testing.T) {
	cfg := &config.Config{
		Port:        "8080",
		DatabaseURL: "file:./test.db",
		ToolName:    "test-tool",
	}

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "file:./test.db", cfg.DatabaseURL)
	assert.Equal(t, "test-tool", cfg.ToolName)
}

func TestLoad_Integration(t *testing.T) {
	// This test ensures Load() creates a valid Config struct
	cfg := config.Load()

	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.DatabaseURL)
	assert.NotEmpty(t, cfg.ToolName)
	
	// Port should be numeric string or valid port
	if port, err := strconv.Atoi(cfg.Port); err == nil {
		assert.True(t, port > 0 && port <= 65535, "Port should be in valid range")
	}
	
	// DatabaseURL should have some basic structure
	assert.Contains(t, cfg.DatabaseURL, "file:", "DatabaseURL should contain file: prefix")
}