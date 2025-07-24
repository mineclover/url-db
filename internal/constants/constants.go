package constants

// Server configuration constants
const (
	// Server metadata
	DefaultServerName    = "url-db"
	DefaultServerVersion = "1.0.0"
	MCPServerName        = "url-db-mcp-server"
	ServerDescription    = "URL 데이터베이스 MCP 서버"

	// Network and protocol
	DefaultPort         = 8080
	DefaultMCPMode      = "stdio"
	HTTPContentTypeJSON = "application/json"
	HTTPScheme          = "http"
	HTTPSScheme         = "https"

	// MCP operating modes
	MCPModeStdio = "stdio"
	MCPModeSSE   = "sse"
	MCPModeHTTP  = "http"

	// Database
	DefaultDBPath   = "url-db.sqlite"
	DefaultDBDriver = "sqlite3"
	TestDBPrefix    = "test_"

	// Limits and validation
	MaxDomainNameLength     = 50
	MaxToolNameLength       = 50
	MaxIDLength             = 20
	MaxTitleLength          = 255
	MaxDescriptionLength    = 1000
	MaxURLLength            = 2048
	MaxAttributeValueLength = 2048
	MaxBatchSize            = 100
	MaxPageSize             = 100
	DefaultPageSize         = 20

	// Composite key format
	CompositeKeyFormat    = "url-db:domain:id"
	CompositeKeySeparator = ":"

	// MCP protocol
	MCPProtocolVersion = "2025-06-18"
	JSONRPCVersion     = "2.0"

	// File extensions and types
	SQLiteExtension = ".sqlite"
	YAMLExtension   = ".yaml"
	JSONExtension   = ".json"
	GoExtension     = ".go"
	PythonExtension = ".py"
)

// Error message constants
const (
	ErrDomainNotFound       = "domain not found"
	ErrNodeNotFound         = "node not found"
	ErrAttributeNotFound    = "attribute not found"
	ErrInvalidCompositeID   = "invalid composite ID format"
	ErrDuplicateDomain      = "domain already exists"
	ErrDuplicateNode        = "node already exists in this domain"
	ErrDuplicateAttribute   = "attribute already exists"
	ErrInvalidURL           = "invalid URL format"
	ErrInvalidParameters    = "invalid parameters"
	ErrDatabaseError        = "database error"
	ErrServerNotInitialized = "server not initialized"
	ErrToolNotFound         = "tool not found"
	ErrResourceNotFound     = "resource not found"

	// MCP Protocol error messages
	ErrParseError            = "Parse error"
	ErrInvalidInitParams     = "Invalid initialize parameters"
	ErrInvalidToolCallParams = "Invalid tool call parameters"
	ErrInvalidResourceParams = "Invalid resource read parameters"
	ErrToolExecutionFailed   = "Tool execution failed"
	ErrFailedToGetResources  = "Failed to get resources"
	ErrFailedToReadResource  = "Failed to read resource"
	ErrMethodNotFound        = "Method not found: %s"
)

// HTTP status codes
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// Log levels and categories
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"

	LogCategoryMCP      = "mcp"
	LogCategoryHTTP     = "http"
	LogCategoryDatabase = "database"
	LogCategoryService  = "service"
)

// Environment variables
const (
	EnvDatabaseURL          = "DATABASE_URL"
	EnvPort                 = "PORT"
	EnvLogLevel             = "LOG_LEVEL"
	EnvMCPMode              = "MCP_MODE"
	EnvAutoCreateAttributes = "AUTO_CREATE_ATTRIBUTES"
)

// Resource URI schemes
const (
	MCPResourceScheme  = "mcp"
	FileResourceScheme = "file"
	HTTPResourceScheme = "http"
)

// Validation patterns
const (
	DomainNamePattern = `^[a-zA-Z0-9_-]+$`
	URLPattern        = `^https?://.*`
	EmailPattern      = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)
