# MCP SQLite Server Pattern Analysis for URL-DB Implementation

## Executive Summary

After analyzing the MCP SQLite server implementation and comparing it with the URL-DB MCP server, I've identified several patterns and best practices that could enhance the URL-DB implementation. The SQLite server follows a simpler, more focused architecture that could inform improvements to URL-DB's more complex implementation.

## Key Patterns from SQLite MCP Server

### 1. **Simplified Architecture**

#### SQLite Pattern:
```python
# Single database class with clear responsibilities
class SqliteDatabase:
    def __init__(self, db_path: str):
        self.db_path = str(Path(db_path).expanduser())
        Path(self.db_path).parent.mkdir(parents=True, exist_ok=True)
        self._init_database()
        self.insights: list[str] = []
```

#### URL-DB Current:
- Complex multi-layer architecture with separate repositories, services, and handlers
- Multiple adapter layers to bridge interfaces

**Recommendation**: Consider consolidating some layers where appropriate, particularly for MCP-specific operations.

### 2. **Direct Handler Registration**

#### SQLite Pattern:
```python
async def main(db_path: str):
    db = SqliteDatabase(db_path)
    server = Server("sqlite-manager")
    
    @server.list_resources()
    async def handle_list_resources() -> list[types.Resource]:
        # Direct implementation
    
    @server.read_resource()
    async def handle_read_resource(uri: AnyUrl) -> str:
        # Direct implementation
```

#### URL-DB Current:
- Handlers are registered through a more complex routing mechanism
- Multiple layers of abstraction between MCP server and actual handlers

**Recommendation**: Consider using decorator-based handler registration for cleaner MCP integration.

### 3. **Unified Entry Point**

#### SQLite Pattern:
```python
def main():
    parser = argparse.ArgumentParser(description='SQLite MCP Server')
    parser.add_argument('--db-path', default="./sqlite_mcp_server.db",
                       help='Path to SQLite database file')
    args = parser.parse_args()
    asyncio.run(server.main(args.db_path))
```

#### URL-DB Current:
- Complex main function with multiple modes and extensive configuration
- Mixed HTTP and MCP server initialization

**Recommendation**: Consider separating MCP-specific entry point for cleaner separation of concerns.

### 4. **Resource Management**

#### SQLite Pattern:
- Simple resource URIs (e.g., `memo://insights`)
- Direct resource reading without complex routing

#### URL-DB Current:
- Complex composite key system
- Multiple resource types with intricate relationships

**Recommendation**: Simplify resource URI patterns where possible while maintaining functionality.

### 5. **Error Handling**

#### SQLite Pattern:
```python
try:
    # Operation
except Exception as e:
    return [types.TextContent(
        type="text",
        text=f"Error: {str(e)}"
    )]
```

#### URL-DB Current:
- Custom error types for each domain
- Complex error propagation through layers

**Recommendation**: Standardize MCP-specific error handling separate from HTTP API errors.

## Best Practices to Apply

### 1. **Cleaner Separation of MCP and HTTP Concerns**

Create a dedicated MCP package structure:
```
internal/mcp/
├── server.go        # MCP server implementation
├── handlers.go      # MCP-specific handlers
├── resources.go     # Resource management
├── tools.go         # Tool implementations
└── standalone/      # Standalone MCP server
    └── main.go      # Dedicated MCP entry point
```

### 2. **Simplified Tool Registration**

Instead of complex registry patterns, consider:
```go
func RegisterTools(server *StdioServer, service MCPService) {
    server.RegisterTool("create_url", "Create a new URL entry", 
        func(args map[string]interface{}) (interface{}, error) {
            // Direct implementation
        })
}
```

### 3. **Streamlined Configuration**

For MCP mode:
```go
type MCPConfig struct {
    DatabasePath string
    ToolName     string
    ServerName   string
    Version      string
}
```

### 4. **Direct Database Access for MCP**

Consider a simplified service layer specifically for MCP operations:
```go
type MCPDirectService struct {
    db *sql.DB
    compositeKeyService *CompositeKeyService
}
```

### 5. **Improved stdio Mode**

Enhancements to stdio server:
```go
func (s *StdioServer) Start() error {
    // Add proper signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    // Add graceful shutdown
    go func() {
        <-sigChan
        s.Shutdown()
    }()
    
    // Process requests
    return s.processRequests()
}
```

## Specific Improvements for URL-DB

### 1. **Separate MCP Binary**

Create `cmd/mcp-server/main.go`:
```go
package main

import (
    "flag"
    "os"
    "url-db/internal/mcp/standalone"
)

func main() {
    var (
        dbPath   = flag.String("db", "", "Database path")
        toolName = flag.String("tool", "url-db", "Tool name")
    )
    flag.Parse()
    
    server := standalone.NewServer(*dbPath, *toolName)
    if err := server.Run(); err != nil {
        os.Exit(1)
    }
}
```

### 2. **Simplified Resource URIs**

Current: `url-db://domain:example.com/url:123`
Proposed: `url://example.com/123` or `domain://example.com`

### 3. **Batch Operations**

Following SQLite's pattern of business logic:
```go
type BatchOperations struct {
    service MCPService
}

func (b *BatchOperations) ImportURLs(ctx context.Context, domain string, urls []string) error {
    // Implement efficient batch import
}
```

### 4. **MCP-Specific Logging**

```go
type MCPLogger struct {
    enabled bool
    output  io.Writer
}

func (l *MCPLogger) Log(format string, args ...interface{}) {
    if l.enabled && l.output != nil {
        fmt.Fprintf(l.output, format+"\n", args...)
    }
}
```

### 5. **Tool Simplification**

Current tools are comprehensive but complex. Consider:
- Grouping related operations
- Providing both simple and advanced versions
- Better parameter validation with clear error messages

## Implementation Priority

1. **High Priority**:
   - Separate MCP entry point
   - Simplified error handling for MCP
   - Cleaner stdio server implementation

2. **Medium Priority**:
   - Resource URI simplification
   - Direct service layer for MCP
   - Improved logging for debugging

3. **Low Priority**:
   - Batch operations
   - Advanced tool grouping
   - Performance optimizations

## Conclusion

The SQLite MCP server demonstrates that a simpler, more focused architecture can be effective for MCP implementations. While URL-DB has more complex requirements, adopting some of these patterns—particularly around separation of concerns, simplified handler registration, and cleaner error handling—could improve maintainability and user experience.

The key insight is that MCP servers benefit from being treated as a distinct concern from HTTP APIs, with their own optimized patterns and simplified architectures where appropriate.