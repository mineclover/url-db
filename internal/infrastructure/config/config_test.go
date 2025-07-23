package config_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"url-db/internal/constants"
	"url-db/internal/infrastructure/config"
)

func TestConfig_Load_DefaultValues(t *testing.T) {
	// Clear environment variables to ensure we get defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("TOOL_NAME")

	cfg := config.Load()

	expectedPort := strconv.Itoa(constants.DefaultPort)
	expectedDatabaseURL := "file:./" + constants.DefaultDBPath
	expectedToolName := constants.DefaultServerName

	assert.Equal(t, expectedPort, cfg.Port)
	assert.Equal(t, expectedDatabaseURL, cfg.DatabaseURL)
	assert.Equal(t, expectedToolName, cfg.ToolName)
}

func TestConfig_Load_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	expectedPort := "9090"
	expectedDatabaseURL := "file:/custom/path/database.db"
	expectedToolName := "custom-tool"

	os.Setenv("PORT", expectedPort)
	os.Setenv("DATABASE_URL", expectedDatabaseURL)
	os.Setenv("TOOL_NAME", expectedToolName)
	
	// Clean up after test
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("TOOL_NAME")
	}()

	cfg := config.Load()

	assert.Equal(t, expectedPort, cfg.Port)
	assert.Equal(t, expectedDatabaseURL, cfg.DatabaseURL)
	assert.Equal(t, expectedToolName, cfg.ToolName)
}

func TestConfig_Load_PartialEnvironmentVariables(t *testing.T) {
	// Set only some environment variables
	os.Setenv("PORT", "3000")
	os.Unsetenv("DATABASE_URL")
	os.Setenv("TOOL_NAME", "partial-tool")
	
	// Clean up after test
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("TOOL_NAME")
	}()

	cfg := config.Load()

	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, "file:./"+constants.DefaultDBPath, cfg.DatabaseURL)
	assert.Equal(t, "partial-tool", cfg.ToolName)
}

func TestConfig_Load_EmptyEnvironmentVariables(t *testing.T) {
	// Set environment variables to empty strings
	os.Setenv("PORT", "")
	os.Setenv("DATABASE_URL", "")
	os.Setenv("TOOL_NAME", "")
	
	// Clean up after test
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("TOOL_NAME")
	}()

	cfg := config.Load()

	// Empty strings should fall back to defaults
	expectedPort := strconv.Itoa(constants.DefaultPort)
	expectedDatabaseURL := "file:./" + constants.DefaultDBPath
	expectedToolName := constants.DefaultServerName

	assert.Equal(t, expectedPort, cfg.Port)
	assert.Equal(t, expectedDatabaseURL, cfg.DatabaseURL)
	assert.Equal(t, expectedToolName, cfg.ToolName)
}

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

func TestConfig_Load_DifferentEnvironmentCombinations(t *testing.T) {
	testCases := []struct {
		name        string
		envVars     map[string]string
		expectedCfg config.Config
	}{
		{
			name: "production-like config",
			envVars: map[string]string{
				"PORT":         "80",
				"DATABASE_URL": "postgresql://prod-db:5432/urldb",
				"TOOL_NAME":    "url-db-prod",
			},
			expectedCfg: config.Config{
				Port:        "80",
				DatabaseURL: "postgresql://prod-db:5432/urldb",
				ToolName:    "url-db-prod",
			},
		},
		{
			name: "development config",
			envVars: map[string]string{
				"PORT":         "8080",
				"DATABASE_URL": "file:./dev.db",
				"TOOL_NAME":    "url-db-dev",
			},
			expectedCfg: config.Config{
				Port:        "8080",
				DatabaseURL: "file:./dev.db",
				ToolName:    "url-db-dev",
			},
		},
		{
			name: "only port override",
			envVars: map[string]string{
				"PORT": "4000",
			},
			expectedCfg: config.Config{
				Port:        "4000",
				DatabaseURL: "file:./" + constants.DefaultDBPath,
				ToolName:    constants.DefaultServerName,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up environment first
			os.Unsetenv("PORT")
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("TOOL_NAME")

			// Set test environment variables
			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}

			// Clean up after test
			defer func() {
				for key := range tc.envVars {
					os.Unsetenv(key)
				}
			}()

			cfg := config.Load()

			assert.Equal(t, tc.expectedCfg.Port, cfg.Port)
			assert.Equal(t, tc.expectedCfg.DatabaseURL, cfg.DatabaseURL)
			assert.Equal(t, tc.expectedCfg.ToolName, cfg.ToolName)
		})
	}
}

func TestConfig_Load_SpecialCharacters(t *testing.T) {
	// Test with special characters in environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("DATABASE_URL", "file:./test with spaces.db")
	os.Setenv("TOOL_NAME", "tool-with-dashes_and_underscores")
	
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("TOOL_NAME")
	}()

	cfg := config.Load()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "file:./test with spaces.db", cfg.DatabaseURL)
	assert.Equal(t, "tool-with-dashes_and_underscores", cfg.ToolName)
}

func TestConfig_Load_UnicodeCharacters(t *testing.T) {
	// Test with unicode characters
	os.Setenv("TOOL_NAME", "url-db-测试")
	os.Setenv("DATABASE_URL", "file:./データベース.db")
	
	defer func() {
		os.Unsetenv("TOOL_NAME")
		os.Unsetenv("DATABASE_URL")
	}()

	cfg := config.Load()

	assert.Equal(t, strconv.Itoa(constants.DefaultPort), cfg.Port)
	assert.Equal(t, "file:./データベース.db", cfg.DatabaseURL)
	assert.Equal(t, "url-db-测试", cfg.ToolName)
}

func TestConfig_Load_MultipleRuns(t *testing.T) {
	// Test that multiple Load() calls return consistent results
	os.Setenv("PORT", "5000")
	os.Setenv("DATABASE_URL", "file:./multi.db")
	
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
	}()

	cfg1 := config.Load()
	cfg2 := config.Load()

	assert.Equal(t, cfg1.Port, cfg2.Port)
	assert.Equal(t, cfg1.DatabaseURL, cfg2.DatabaseURL)
	assert.Equal(t, cfg1.ToolName, cfg2.ToolName)

	// Values should match expected
	assert.Equal(t, "5000", cfg1.Port)
	assert.Equal(t, "file:./multi.db", cfg1.DatabaseURL)
	assert.Equal(t, constants.DefaultServerName, cfg1.ToolName)
}