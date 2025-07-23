package constants_test

import (
	"regexp"
	"strings"
	"testing"

	"url-db/internal/shared/constants"

	"github.com/stretchr/testify/assert"
)

func TestServerConstants(t *testing.T) {
	t.Run("Server metadata constants", func(t *testing.T) {
		assert.Equal(t, "url-db", constants.DefaultServerName)
		assert.Equal(t, "1.0.0", constants.DefaultServerVersion)
		assert.Equal(t, "url-db-mcp-server", constants.MCPServerName)
		assert.Equal(t, "URL 데이터베이스 MCP 서버", constants.ServerDescription)
	})

	t.Run("Network and protocol constants", func(t *testing.T) {
		assert.Equal(t, 8080, constants.DefaultPort)
		assert.Equal(t, "stdio", constants.DefaultMCPMode)
		assert.Equal(t, "application/json", constants.HTTPContentTypeJSON)
		assert.Equal(t, "http", constants.HTTPScheme)
		assert.Equal(t, "https", constants.HTTPSScheme)
	})

	t.Run("Database constants", func(t *testing.T) {
		assert.Equal(t, "url-db.sqlite", constants.DefaultDBPath)
		assert.Equal(t, "sqlite3", constants.DefaultDBDriver)
		assert.Equal(t, "test_", constants.TestDBPrefix)
	})
}

func TestLimitsAndValidation(t *testing.T) {
	t.Run("Length limits are positive", func(t *testing.T) {
		assert.Greater(t, constants.MaxDomainNameLength, 0)
		assert.Greater(t, constants.MaxTitleLength, 0)
		assert.Greater(t, constants.MaxDescriptionLength, 0)
		assert.Greater(t, constants.MaxURLLength, 0)
		assert.Greater(t, constants.MaxAttributeValueLength, 0)
	})

	t.Run("Batch and page sizes are reasonable", func(t *testing.T) {
		assert.Greater(t, constants.MaxBatchSize, 0)
		assert.Greater(t, constants.MaxPageSize, 0)
		assert.Greater(t, constants.DefaultPageSize, 0)
		
		// DefaultPageSize should not exceed MaxPageSize
		assert.LessOrEqual(t, constants.DefaultPageSize, constants.MaxPageSize)
	})

	t.Run("Specific limit values", func(t *testing.T) {
		assert.Equal(t, 50, constants.MaxDomainNameLength)
		assert.Equal(t, 255, constants.MaxTitleLength)
		assert.Equal(t, 1000, constants.MaxDescriptionLength)
		assert.Equal(t, 2048, constants.MaxURLLength)
		assert.Equal(t, 2048, constants.MaxAttributeValueLength)
		assert.Equal(t, 100, constants.MaxBatchSize)
		assert.Equal(t, 100, constants.MaxPageSize)
		assert.Equal(t, 20, constants.DefaultPageSize)
	})
}

func TestCompositeKeyConstants(t *testing.T) {
	t.Run("Composite key format and separator", func(t *testing.T) {
		assert.Equal(t, "url-db:domain:id", constants.CompositeKeyFormat)
		assert.Equal(t, ":", constants.CompositeKeySeparator)
		
		// Ensure format contains separator
		assert.Contains(t, constants.CompositeKeyFormat, constants.CompositeKeySeparator)
	})
}

func TestMCPProtocolConstants(t *testing.T) {
	t.Run("Protocol versions", func(t *testing.T) {
		assert.Equal(t, "2024-11-05", constants.MCPProtocolVersion)
		assert.Equal(t, "2.0", constants.JSONRPCVersion)
	})

	t.Run("Version formats are valid", func(t *testing.T) {
		// MCP protocol version should be in date format
		datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		assert.True(t, datePattern.MatchString(constants.MCPProtocolVersion), "MCP protocol version should match date format")
		
		// JSON-RPC version should be semantic version-like
		versionPattern := regexp.MustCompile(`^\d+\.\d+$`)
		assert.True(t, versionPattern.MatchString(constants.JSONRPCVersion), "JSON-RPC version should match version format")
	})
}

func TestFileExtensions(t *testing.T) {
	t.Run("Extension formats", func(t *testing.T) {
		extensions := []string{
			constants.SQLiteExtension,
			constants.YAMLExtension,
			constants.JSONExtension,
			constants.GoExtension,
			constants.PythonExtension,
		}

		for _, ext := range extensions {
			// All extensions should start with a dot
			assert.True(t, strings.HasPrefix(ext, "."), "Extension %s should start with a dot", ext)
			// All extensions should be lowercase
			assert.Equal(t, strings.ToLower(ext), ext, "Extension %s should be lowercase", ext)
		}
	})

	t.Run("Specific extension values", func(t *testing.T) {
		assert.Equal(t, ".sqlite", constants.SQLiteExtension)
		assert.Equal(t, ".yaml", constants.YAMLExtension)
		assert.Equal(t, ".json", constants.JSONExtension)
		assert.Equal(t, ".go", constants.GoExtension)
		assert.Equal(t, ".py", constants.PythonExtension)
	})
}

func TestErrorMessageConstants(t *testing.T) {
	t.Run("Basic error messages", func(t *testing.T) {
		errorMessages := []string{
			constants.ErrDomainNotFound,
			constants.ErrNodeNotFound,
			constants.ErrAttributeNotFound,
			constants.ErrInvalidCompositeID,
			constants.ErrDuplicateDomain,
			constants.ErrDuplicateAttribute,
			constants.ErrInvalidURL,
			constants.ErrInvalidParameters,
			constants.ErrDatabaseError,
			constants.ErrServerNotInitialized,
			constants.ErrToolNotFound,
			constants.ErrResourceNotFound,
		}

		for _, msg := range errorMessages {
			assert.NotEmpty(t, msg, "Error message should not be empty")
			assert.True(t, len(msg) > 5, "Error message should be descriptive: %s", msg)
		}
	})

	t.Run("MCP protocol error messages", func(t *testing.T) {
		mcpErrors := []string{
			constants.ErrParseError,
			constants.ErrInvalidInitParams,
			constants.ErrInvalidToolCallParams,
			constants.ErrInvalidResourceParams,
			constants.ErrToolExecutionFailed,
			constants.ErrFailedToGetResources,
			constants.ErrFailedToReadResource,
			constants.ErrMethodNotFound,
		}

		for _, msg := range mcpErrors {
			assert.NotEmpty(t, msg, "MCP error message should not be empty")
		}
	})

	t.Run("Error message with format placeholder", func(t *testing.T) {
		assert.Contains(t, constants.ErrMethodNotFound, "%s", "ErrMethodNotFound should contain format placeholder")
	})
}

func TestHTTPStatusCodes(t *testing.T) {
	t.Run("HTTP status code values", func(t *testing.T) {
		assert.Equal(t, 200, constants.StatusOK)
		assert.Equal(t, 201, constants.StatusCreated)
		assert.Equal(t, 400, constants.StatusBadRequest)
		assert.Equal(t, 404, constants.StatusNotFound)
		assert.Equal(t, 500, constants.StatusInternalServerError)
	})

	t.Run("HTTP status codes are valid", func(t *testing.T) {
		statusCodes := []int{
			constants.StatusOK,
			constants.StatusCreated,
			constants.StatusBadRequest,
			constants.StatusNotFound,
			constants.StatusInternalServerError,
		}

		for _, code := range statusCodes {
			assert.GreaterOrEqual(t, code, 100, "HTTP status code should be >= 100")
			assert.LessOrEqual(t, code, 599, "HTTP status code should be <= 599")
		}
	})
}

func TestLogConstants(t *testing.T) {
	t.Run("Log levels", func(t *testing.T) {
		logLevels := []string{
			constants.LogLevelDebug,
			constants.LogLevelInfo,
			constants.LogLevelWarn,
			constants.LogLevelError,
		}

		for _, level := range logLevels {
			assert.NotEmpty(t, level, "Log level should not be empty")
			assert.Equal(t, strings.ToLower(level), level, "Log level should be lowercase")
		}
	})

	t.Run("Log categories", func(t *testing.T) {
		logCategories := []string{
			constants.LogCategoryMCP,
			constants.LogCategoryHTTP,
			constants.LogCategoryDatabase,
			constants.LogCategoryService,
		}

		for _, category := range logCategories {
			assert.NotEmpty(t, category, "Log category should not be empty")
			assert.Equal(t, strings.ToLower(category), category, "Log category should be lowercase")
		}
	})
}

func TestEnvironmentVariables(t *testing.T) {
	t.Run("Environment variable names", func(t *testing.T) {
		envVars := []string{
			constants.EnvDatabaseURL,
			constants.EnvPort,
			constants.EnvLogLevel,
			constants.EnvMCPMode,
		}

		for _, envVar := range envVars {
			assert.NotEmpty(t, envVar, "Environment variable should not be empty")
			// Environment variables should be uppercase by convention
			assert.Equal(t, strings.ToUpper(envVar), envVar, "Environment variable should be uppercase")
		}
	})

	t.Run("Specific environment variable values", func(t *testing.T) {
		assert.Equal(t, "DATABASE_URL", constants.EnvDatabaseURL)
		assert.Equal(t, "PORT", constants.EnvPort)
		assert.Equal(t, "LOG_LEVEL", constants.EnvLogLevel)
		assert.Equal(t, "MCP_MODE", constants.EnvMCPMode)
	})
}

func TestResourceURISchemes(t *testing.T) {
	t.Run("URI schemes", func(t *testing.T) {
		schemes := []string{
			constants.MCPResourceScheme,
			constants.FileResourceScheme,
			constants.HTTPResourceScheme,
		}

		for _, scheme := range schemes {
			assert.NotEmpty(t, scheme, "URI scheme should not be empty")
			assert.Equal(t, strings.ToLower(scheme), scheme, "URI scheme should be lowercase")
		}
	})

	t.Run("Specific scheme values", func(t *testing.T) {
		assert.Equal(t, "mcp", constants.MCPResourceScheme)
		assert.Equal(t, "file", constants.FileResourceScheme)
		assert.Equal(t, "http", constants.HTTPResourceScheme)
	})
}

func TestValidationPatterns(t *testing.T) {
	t.Run("Regex patterns compile successfully", func(t *testing.T) {
		patterns := map[string]string{
			"DomainNamePattern": constants.DomainNamePattern,
			"URLPattern":        constants.URLPattern,
			"EmailPattern":      constants.EmailPattern,
		}

		for name, pattern := range patterns {
			_, err := regexp.Compile(pattern)
			assert.NoError(t, err, "Pattern %s should compile successfully: %s", name, pattern)
		}
	})

	t.Run("Domain name pattern validation", func(t *testing.T) {
		re, _ := regexp.Compile(constants.DomainNamePattern)
		
		validNames := []string{"test", "test-domain", "test_domain", "test123", "123test"}
		invalidNames := []string{"test.domain", "test domain", "test@domain", "test/domain", ""}

		for _, name := range validNames {
			assert.True(t, re.MatchString(name), "Domain name '%s' should be valid", name)
		}

		for _, name := range invalidNames {
			assert.False(t, re.MatchString(name), "Domain name '%s' should be invalid", name)
		}
	})

	t.Run("URL pattern validation", func(t *testing.T) {
		re, _ := regexp.Compile(constants.URLPattern)
		
		validURLs := []string{
			"http://example.com",
			"https://example.com",
			"http://example.com/path",
			"https://example.com/path?query=value",
		}
		invalidURLs := []string{
			"ftp://example.com",
			"example.com",
			"mailto:test@example.com",
		}

		for _, url := range validURLs {
			assert.True(t, re.MatchString(url), "URL '%s' should be valid", url)
		}

		for _, url := range invalidURLs {
			assert.False(t, re.MatchString(url), "URL '%s' should be invalid", url)
		}
	})

	t.Run("Email pattern validation", func(t *testing.T) {
		re, _ := regexp.Compile(constants.EmailPattern)
		
		validEmails := []string{
			"test@example.com",
			"user.name@example.co.uk",
			"test123@test-domain.com",
		}
		invalidEmails := []string{
			"invalid-email",
			"@example.com",
			"test@",
			"test.example.com",
		}

		for _, email := range validEmails {
			assert.True(t, re.MatchString(email), "Email '%s' should be valid", email)
		}

		for _, email := range invalidEmails {
			assert.False(t, re.MatchString(email), "Email '%s' should be invalid", email)
		}
	})
}

func TestConstantsConsistency(t *testing.T) {
	t.Run("Default values consistency", func(t *testing.T) {
		// DefaultPageSize should be reasonable compared to MaxPageSize
		assert.LessOrEqual(t, constants.DefaultPageSize, constants.MaxPageSize)
		
		// URL length should be generous for modern URLs
		assert.GreaterOrEqual(t, constants.MaxURLLength, 2000)
		
		// Server name consistency
		assert.Contains(t, constants.MCPServerName, constants.DefaultServerName)
	})

	t.Run("String constants are not empty", func(t *testing.T) {
		stringConstants := []string{
			constants.DefaultServerName,
			constants.DefaultServerVersion,
			constants.MCPServerName,
			constants.DefaultMCPMode,
			constants.HTTPContentTypeJSON,
			constants.DefaultDBPath,
			constants.DefaultDBDriver,
			constants.CompositeKeyFormat,
			constants.CompositeKeySeparator,
			constants.MCPProtocolVersion,
			constants.JSONRPCVersion,
		}

		for _, constant := range stringConstants {
			assert.NotEmpty(t, constant, "String constant should not be empty")
		}
	})
}