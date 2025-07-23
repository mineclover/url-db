# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

URL-DB is a Go-based URL database management system with MCP (Model Context Protocol) integration. It provides domain-based URL organization with unlimited attribute tagging, designed for AI assistant integration (Claude Desktop, Cursor).

## Common Development Commands

```bash
# Build the project
make build                    # Build using Makefile
go build ./cmd/server         # Direct Go build

# Run tests  
go test -v ./...             # Run all tests
go test -v ./internal/mcp/... # Run specific package tests
./scripts/test.sh            # Comprehensive test suite with coverage

# Lint and format
make lint                    # Run golangci-lint (install: brew install golangci-lint)
make fmt                     # Format all Go files

# Development mode
make dev                     # Hot reload (requires: go install github.com/cosmtrek/air@latest)

# Run the server
./bin/url-db                 # HTTP mode (default port 8080) 
./bin/url-db -mcp-mode=stdio # MCP stdio mode for AI assistants

# Constants management
python scripts/generate-tool-constants.py  # Generate tool constants from YAML spec
```

## Architecture

The codebase follows a clean layered architecture:

1. **Database Layer** (`/internal/database/` & `/schema.sql`)
   - SQLite with sqlx for enhanced operations
   - **Single Source of Truth**: Schema managed via `/schema.sql` file
   - Enhanced dependency system with 8 built-in types and advanced features
   - Automatic schema loading with fallback strategy in `/internal/database/database.go`
   - Tables: domains, nodes, attributes, node_attributes, node_connections, node_subscriptions, node_dependencies, node_dependencies_v2, dependency_types, dependency_history, dependency_graph_cache, dependency_rules, dependency_impact_analysis, node_events

2. **Repository Layer** (`/internal/repositories/`)
   - Data access patterns with transaction support
   - Key files: `domain.go`, `node.go`, `attribute.go`

3. **Service Layer** (`/internal/services/`)
   - Business logic and validation
   - Cross-domain operations
   - Domain schema enforcement

4. **Handler Layer** (`/internal/handlers/`)
   - REST API endpoints (40+ endpoints)
   - Error handling and response formatting

5. **MCP Layer** (`/internal/mcp/`)
   - JSON-RPC 2.0 implementation
   - 18 tools with distinctive names (without 'mcp' prefix)
   - Resource system in `/internal/mcp/resources.go`
   - Composite key format: `tool-name:domain:id`

6. **Constants Layer** (`/internal/constants/`)
   - Centralized configuration constants
   - Server metadata, network settings, database paths
   - Error messages and validation patterns
   - Single source of truth for all hardcoded values

7. **Advanced Services Layer** (`/internal/services/advanced/`)
   - Enterprise-grade dependency management services
   - Circular dependency detection using Tarjan's algorithm
   - Impact analysis with scoring and recommendations
   - Graph caching and performance optimization

## Domain Schema System

URL-DB implements a powerful domain schema system that ensures data consistency:

1. **Schema Definition**: Each domain defines its allowed attributes with specific types
   - Attribute types: `tag`, `ordered_tag`, `number`, `string`, `markdown`, `image`
   - Attributes are defined at the domain level using `create_domain_attribute`
   
2. **Schema Enforcement**: Nodes can only have attributes defined in their domain's schema
   - Enforced at database level through foreign key constraints
   - Validated in service layer before operations
   - Invalid attributes are rejected with clear error messages

3. **Benefits**:
   - Type safety: Attribute values are validated against their defined types
   - Data consistency: All nodes in a domain follow the same schema
   - Extensibility: New attributes can be added to domains as needed
   - Clear organization: Each domain has its own attribute namespace

Example workflow:
```bash
# 1. Create a domain
create_domain(name="products", description="Product catalog")

# 2. Define the domain schema
create_domain_attribute(domain_name="products", name="price", type="number")
create_domain_attribute(domain_name="products", name="category", type="tag")
create_domain_attribute(domain_name="products", name="description", type="markdown")

# 3. Create nodes that follow the schema
create_node(domain_name="products", url="https://example.com/product1")
set_node_attributes(composite_id="url-db:products:1", attributes=[
  {name: "price", value: "29.99"},
  {name: "category", value: "electronics"}
])

# 4. Invalid attributes are rejected
set_node_attributes(composite_id="url-db:products:1", attributes=[
  {name: "invalid_attr", value: "fail"}  # Error: attribute not in schema
])
```

## Enhanced Dependency System

URL-DB features an enterprise-grade dependency management system with advanced capabilities:

### Database Schema Management
- **Single Source of Truth**: All schema managed in `/schema.sql`
- **Automatic Loading**: `database.go` loads schema from external file with fallback
- **Project Root Detection**: Automatic `go.mod` search for flexible deployment

### Dependency Types (8 Built-in Types)
```go
// Structural Dependencies
- hard:      Strong coupling with cascading operations
- soft:      Loose coupling without cascading  
- reference: Informational reference link only

// Behavioral Dependencies  
- runtime:   Required at runtime execution
- compile:   Required at build/compile time
- optional:  Optional enhancement dependency

// Data Dependencies
- sync:      Synchronous data dependency
- async:     Asynchronous data dependency
```

### Advanced Features
- **Circular Dependency Detection**: Tarjan's algorithm for cycle detection
- **Impact Analysis**: Comprehensive analysis with scoring (0-100) and recommendations
- **Version Constraints**: Semantic versioning support with compatibility checking
- **History Tracking**: Complete audit trail of dependency changes
- **Graph Caching**: Performance optimization for large dependency networks
- **Validation Rules**: Domain-specific dependency rules and constraints

### Core Tables
```sql
dependency_types           -- Type registry with 8 built-in types
node_dependencies_v2       -- Enhanced dependency management
dependency_history         -- Complete change tracking
dependency_graph_cache     -- Performance optimization cache
dependency_rules           -- Validation and constraint rules
dependency_impact_analysis -- Impact analysis results
node_events               -- Event logging system
```

### Performance Optimization
- 25+ specialized indexes for optimal query performance
- Graph traversal caching with automatic expiration
- Batch operations for bulk dependency management
- Memory-efficient algorithms for large networks

## Key Patterns and Conventions

1. **Composite Keys**: Always use format `tool-name:domain:id` (e.g., `url-db:tech:123`)
2. **Error Handling**: Use MCP error codes (-32602 for invalid params, -32603 for internal errors)
3. **Validation**: Validate at service layer before database operations
4. **Transactions**: Use repository transaction methods for multi-step operations
5. **Testing**: Use testify for assertions, create test databases for integration tests
6. **Constants Usage**: Import and use constants from `/internal/constants/` package
7. **Tool Definitions**: Managed through `/specs/mcp-tools.yaml` with auto-generated constants

## MCP Integration

The MCP server supports two modes:
- **stdio**: For AI assistants (Claude Desktop, Cursor)
- **sse**: For HTTP-based integration

MCP provides 18 tools following strict JSON-RPC 2.0 protocol:
- Domain management: `list_domains`, `create_domain`
- Node operations: `list_nodes`, `create_node`, `get_node`, `update_node`, `delete_node`, `find_node_by_url`
- Attribute management: `get_node_attributes`, `set_node_attributes`
- Enhanced queries: `get_node_with_attributes`, `filter_nodes_by_attributes`
- Domain schema: `list_domain_attributes`, `create_domain_attribute`, `get_domain_attribute`, `update_domain_attribute`, `delete_domain_attribute`
- Server info: `get_server_info`

All tools are defined in `/internal/mcp/tools.go` with built-in validation. Tool names and descriptions are generated from `/specs/mcp-tools.yaml` for consistency.

## Testing Approach

### Test Organization and Separation

**Test-Business Logic Separation Principle**: All tests are organized using Go's `package_test` pattern to ensure clean separation between test code and business logic implementation.

```bash
# Run specific test  
go test -v -run TestCreateNode ./internal/repositories/

# Run MCP tests
go test -v ./internal/mcp/...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Integration tests
go test -v ./internal/database/... -tags=integration
```

### Test Structure and Coverage (Current: 20.6%)

#### High Coverage Packages (60%+)
- `internal/compositekey`: **100.0%** - Complete composite key functionality
- `internal/config`: **100.0%** - Configuration management  
- `internal/shared/utils`: **99.3%** - Utility functions
- `internal/database`: **70.3%** - Database connections and operations
- `internal/attributes`: **65.2%** - Attribute type system

#### Test File Organization
Tests follow the `package_test` pattern for clear separation:
- **Repository tests**: `*_test.go` in separate `*_test` packages
- **MCP protocol tests**: `/internal/mcp/*_test.go` 
- **Integration tests**: `/internal/database/database_test.go`
- **Unit tests**: Component-specific test files with `package_test` naming

#### Key Test Categories
1. **Unit Tests**: Individual function/method testing with mocks
2. **Integration Tests**: Database operations with in-memory SQLite
3. **MCP Protocol Tests**: JSON-RPC 2.0 compliance and tool execution
4. **Validation Tests**: Input validation, error handling, edge cases
5. **Utility Tests**: String processing, normalization, composite keys

### Test Database Strategy
- **Isolation**: Each test uses independent in-memory SQLite databases
- **Cleanup**: Automatic cleanup after each test case
- **Fixtures**: Standardized test data setup via helper functions
- **Transactions**: Repository tests wrap operations in transactions

## Important Implementation Details

1. **Attribute System**: 6 types (tag, ordered_tag, number, string, markdown, image)
2. **Database Path**: Use `-db-path` flag or `DATABASE_URL` env var
3. **Tool Name**: Customizable via `-tool-name` flag (affects composite keys)
4. **Resource URIs**: Format `mcp://resource-type/path` for MCP resource system
5. **Batch Operations**: Use `SetNodeAttributes` for efficient bulk updates
6. **Constants Management**: All configuration values centralized in `/internal/constants/`
7. **Tool Specification**: Single source of truth in `/specs/mcp-tools.yaml`
8. **Code Generation**: Use `scripts/generate-tool-constants.py` to update constants

## Development Principles

### Test-Business Logic Separation
**Core Principle**: Maintain strict separation between test code and business logic to ensure clean, maintainable code architecture.

#### Package Test Pattern (`package_test`)
- **Separation**: All test files use `package_test` naming (e.g., `package compositekey_test`)
- **Isolation**: Tests access only public APIs of the packages they test
- **Independence**: Business logic remains unaware of test implementation details
- **Maintainability**: Changes to business logic don't break test organization

#### Test Development Guidelines
1. **Test Package Naming**: Always use `package [packagename]_test` 
2. **Import Strategy**: Import the package being tested explicitly
3. **Public API Focus**: Test only exported functions and types
4. **Mock Dependencies**: Use interfaces and dependency injection for testability
5. **Test Data Isolation**: Each test should create its own test data
6. **Cleanup**: Ensure proper cleanup after each test case

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
4. **Write Tests with Separation**: Use `package_test` pattern for all test files
5. **Add REST API (if needed)**: REST endpoints are secondary to MCP tools
6. **Verify Coverage**: Ensure adequate test coverage for new components
7. **Update Documentation**: Document MCP usage patterns and examples

#### Test-First Development Process
1. **Define Test Package**: Create `*_test.go` files with `package [name]_test`
2. **Write Interface Tests**: Test public APIs and expected behaviors
3. **Implement Business Logic**: Write minimal code to pass tests
4. **Refactor with Confidence**: Tests ensure functionality remains intact
5. **Monitor Coverage**: Use `go test -cover` to track test completeness

### MCP Tool Design Guidelines

1. **Tool Naming**: Use clear, action-oriented names (e.g., `create_mcp_domain`, not `domain_new`)
2. **Parameter Design**: 
   - Use descriptive parameter names
   - Required parameters should be minimal
   - Optional parameters for extended functionality
3. **Return Values**: Always return useful information for chaining operations
4. **Error Messages**: Provide clear, actionable error messages
5. **Composite Keys**: Return composite IDs for created/updated resources

### Potential Future MCP Features

Features that could be added for enhanced functionality:
- Bulk operations for efficiency
- Advanced search/filter capabilities
- Export/import functionality
- Node connections and relationships management
- Subscription and dependency management