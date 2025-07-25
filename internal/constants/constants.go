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
	
	// MCP notification methods
	MCPLogNotificationMethod = "notifications/message"

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

// Content scanning constants
const (
	DefaultMaxTokensPerPage = 3000  // Default tokens per page
	MaxTokensPerPage        = 5000  // Maximum tokens per page
	MinTokensPerNode        = 20    // Minimum tokens per node
	AvgTokensPerNode        = 100   // Average tokens per node
	ScanBatchSize           = 100   // Batch size for scanning
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

// Date/Time formatting
const (
	DateTimeFormat    = "2006-01-02 15:04:05"
	ISODateTimeFormat = "2006-01-02T15:04:05Z"
)

// Pagination and search limits
const (
	DefaultSearchLimit         = 10
	LargeFetchLimit           = 1000
	DefaultRecentlyModified   = 10
	DefaultPaginationOffset   = 0
)

// Attribute validation limits
const (
	MaxStringLength    = 500
	MaxTagLength       = 50
	MaxMarkdownLength  = 10000
	MaxImageSize       = 10 * 1024 * 1024 // 10MB
	MBInBytes         = 1024 * 1024
	MinOrderIndex     = 0
	MaxTemplateNameLength = 255
)

// Database configuration
const (
	DefaultMaxOpenConns    = 10
	DefaultMaxIdleConns    = 5
	ProductionMaxOpenConns = 100
	ProductionMaxIdleConns = 50
	TestMaxConns          = 1
	DirectoryPermissions  = 0755
	
	// Database journal modes
	JournalModeWAL    = "WAL"
	JournalModeDelete = "DELETE"
	
	// Database synchronous modes  
	SyncModeNormal = "NORMAL"
	SyncModeOff    = "OFF"
	SyncModeFull   = "FULL"
	
	// Special database URLs
	InMemoryDB = ":memory:"
)

// Template validation error codes
const (
	ErrTemplateValueNotAllowed     = "template_value_not_allowed"
	ErrTemplateRequiredButMissing  = "template_required_but_missing"
	ErrTemplateValueFormatMismatch = "template_value_format_mismatch"
)

// Common validation error messages
const (
	ValidationErrorCode           = "validation_error"
	ErrOrderIndexNotAllowed       = "order_index not allowed for %s type"
	ErrOrderIndexRequired         = "order_index is required for ordered_tag type"
	ErrOrderIndexNonNegative      = "order_index must be non-negative"
	ErrInvalidMarkdownSyntax      = "invalid markdown syntax: unbalanced brackets or parentheses"
	ErrUnsupportedImageType       = "unsupported image type: %s. Supported types: jpeg, png, gif, webp"
	ErrInvalidBase64Encoding      = "invalid base64 encoding"
	ErrImageSizeExceeded          = "image size exceeds maximum limit of 10MB (actual: %.2fMB)"
	ErrInvalidURLFormat           = "invalid URL format"
	ErrURLMustUseHTTPS            = "URL must use http or https scheme"
	ErrURLMustHaveHost            = "URL must have a valid host"
)

// Template service error messages
const (
	ErrTemplateDataValidationFailed = "Template data validation failed"
	ErrInactiveTemplateModification = "inactive templates cannot be modified"
	ErrTemplateNameEmpty           = "template name cannot be empty"
	ErrTemplateNameTooLong         = "template name cannot exceed 255 characters"
	ErrTemplateNameInvalidChars    = "template name can only contain letters, numbers, hyphens, and underscores"
	ErrTemplateNameInvalidStartEnd = "template name cannot start or end with hyphen or underscore"
	ErrTemplateNotFound           = "template not found"
	ErrTemplateTypeNotFound       = "template type not found or not a string"
	ErrTemplateVersionNotFound    = "template version not found or not a string"
	ErrInvalidJSON                = "invalid JSON"
)

// Image validation constants
const (
	DataImagePrefix = "data:image/"
	Base64Separator = ";base64,"
	Base64Encoding  = "base64"
	
	// Image MIME types
	ImageJPEG = "data:image/jpeg"
	ImagePNG  = "data:image/png"
	ImageGIF  = "data:image/gif"
	ImageWEBP = "data:image/webp"
)

// Template validation method types
const (
	ValidationMethodAllowedValues = "allowed_values"
	ValidationMethodEnum         = "enum"
	ValidationMethodPattern      = "pattern"
	ValidationMethodRange        = "range"
	ValidationMethodSingleValue  = "single_value"
	ValidationMethodUnknown      = "unknown"
	ValidationMethodNoConstraints = "no_template_constraints"
)

// Slice constants
var (
	// Tag forbidden characters
	TagForbiddenChars = []string{",", ";", "|", "\n", "\t"}
	
	// Supported image MIME types
	SupportedImageTypes = []string{
		ImageJPEG,
		ImagePNG,
		ImageGIF,
		ImageWEBP,
	}
)
