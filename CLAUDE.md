# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

URL-DB is a Go-based URL database management system with MCP (Model Context Protocol) integration. It provides domain-based URL organization with unlimited attribute tagging, designed for AI assistant integration (Claude Desktop, Cursor).

## Common Development Commands

```bash
# Build the project
./build.sh                    # Unix/Mac - builds and runs tests
make build                    # Alternative using Makefile

# Run tests
go test -v ./...             # Run all tests
go test -v ./internal/mcp/...  # Run specific package tests
make test-coverage           # Run tests with coverage report

# Lint and format
make lint                    # Run golangci-lint (install: brew install golangci-lint)
make fmt                     # Format all Go files

# Development mode
make dev                     # Hot reload (requires: go install github.com/cosmtrek/air@latest)

# Run the server
./bin/url-db                 # HTTP mode (default port 8080)
./bin/url-db -mcp-mode=stdio # MCP stdio mode for AI assistants
```

## Architecture

The codebase follows a clean layered architecture:

1. **Database Layer** (`/internal/database/`)
   - SQLite with sqlx for enhanced operations
   - Schema: domains, nodes, attributes, node_attributes, node_connections, node_subscriptions, node_dependencies
   - Migration system in `/internal/database/migrations.go`

2. **Repository Layer** (`/internal/repositories/`)
   - Data access patterns with transaction support
   - Key files: `domain.go`, `node.go`, `attribute.go`

3. **Service Layer** (`/internal/services/`)
   - Business logic and validation
   - Cross-domain operations

4. **Handler Layer** (`/internal/handlers/`)
   - REST API endpoints (40+ endpoints)
   - Error handling and response formatting

5. **MCP Layer** (`/internal/mcp/`)
   - JSON-RPC 2.0 implementation
   - 11 tools defined in `/internal/mcp/tools.go`
   - Resource system in `/internal/mcp/resources.go`
   - Composite key format: `tool-name:domain:id`

## Key Patterns and Conventions

1. **Composite Keys**: Always use format `tool-name:domain:id` (e.g., `url-db:tech:123`)
2. **Error Handling**: Use MCP error codes (-32602 for invalid params, -32603 for internal errors)
3. **Validation**: Validate at service layer before database operations
4. **Transactions**: Use repository transaction methods for multi-step operations
5. **Testing**: Use testify for assertions, create test databases for integration tests

## MCP Integration

The MCP server supports two modes:
- **stdio**: For AI assistants (Claude Desktop, Cursor)
- **sse**: For HTTP-based integration

MCP tools follow strict JSON-RPC 2.0 protocol with composite key requirements. All tools are defined in `/internal/mcp/tools.go` with validation in `/internal/mcp/validators.go`.

## Testing Approach

```bash
# Run specific test
go test -v -run TestCreateNode ./internal/repositories/

# Run MCP tests
go test -v ./internal/mcp/...

# Integration tests
go test -v ./internal/database/... -tags=integration
```

Tests use in-memory SQLite databases for isolation. Key test files:
- Repository tests: `*_test.go` in `/internal/repositories/`
- MCP protocol tests: `/internal/mcp/*_test.go`
- Integration tests: `/internal/database/database_test.go`

## Important Implementation Details

1. **Attribute System**: 6 types (tag, ordered_tag, number, string, markdown, image)
2. **Database Path**: Use `-db-path` flag or `DATABASE_URL` env var
3. **Tool Name**: Customizable via `-tool-name` flag (affects composite keys)
4. **Resource URIs**: Format `mcp://resource-type/path` for MCP resource system
5. **Batch Operations**: Use `SetNodeAttributes` for efficient bulk updates

## Development Principles

### MCP-First Development
This project prioritizes MCP functionality as the primary interface. Follow these principles:

1. **MCP is the Primary Interface**: Every major feature should be accessible through MCP tools
2. **Feature Parity**: If a feature exists in the REST API, it should have an MCP equivalent
3. **Composite Key Consistency**: Always use the `tool-name:domain:id` format for identification
4. **AI-Friendly Design**: Design tools with natural language interaction in mind

### Development Workflow

When implementing new features:
1. **Start with MCP Tool Design**: Define the MCP tool interface first
2. **Implement Backend Support**: Add repository/service layers as needed
3. **Create MCP Tool Implementation**: Implement the tool in `/internal/mcp/tools.go`
4. **Add REST API (if needed)**: REST endpoints are secondary to MCP tools
5. **Write Tests**: Focus on MCP integration tests first
6. **Update Documentation**: Document MCP usage patterns and examples

### MCP Tool Design Guidelines

1. **Tool Naming**: Use clear, action-oriented names (e.g., `create_mcp_domain`, not `domain_new`)
2. **Parameter Design**: 
   - Use descriptive parameter names
   - Required parameters should be minimal
   - Optional parameters for extended functionality
3. **Return Values**: Always return useful information for chaining operations
4. **Error Messages**: Provide clear, actionable error messages
5. **Composite Keys**: Return composite IDs for created/updated resources

### Missing MCP Features to Implement

Currently missing MCP tools that should be added:
- Domain attribute management (see `/docs/tasks/mcp-domain-attributes-task.md`)
- Bulk operations for efficiency
- Search/filter capabilities
- Export/import functionality