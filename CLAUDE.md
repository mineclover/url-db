# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

URL-DB is a Go-based URL database management system built with Clean Architecture principles and MCP (Model Context Protocol) integration. It provides domain-based URL organization with unlimited attribute tagging, designed for AI assistant integration (Claude Desktop, Cursor).

**Architecture**: Clean Architecture with 4-layer separation (Domain, Application, Infrastructure, Interface)
**Code Quality**: A- (85/100) with excellent SOLID principles implementation
**Current Status**: Production-ready with comprehensive dependency injection and use case patterns

## Common Development Commands

```bash
# Build the project
make build                    # Build using Makefile
go build ./cmd/server         # Direct Go build

# Run tests  
./scripts/test_runner.sh                    # Comprehensive test runner with options
./scripts/test_runner.sh -m coverage       # Run with detailed coverage analysis
./scripts/test_runner.sh -p internal/mcp   # Test specific package
./scripts/test_runner.sh -m unit -v        # Unit tests with verbose output

# Coverage analysis
./scripts/coverage_analysis.sh             # Detailed coverage analysis and recommendations
go test -coverprofile=coverage.out ./...   # Basic coverage
go tool cover -html=coverage.out -o coverage.html  # Generate HTML coverage report

# Legacy test commands
go test -v ./...             # Run all tests (basic)
go test -v ./internal/mcp/... # Run specific package tests

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

The codebase follows Clean Architecture principles with strict layer separation:

### Clean Architecture Layers

1. **Domain Layer** (`/internal/domain/`)
   - **Entities**: Core business objects (Domain, Node, Attribute) with encapsulated logic
   - **Repository Interfaces**: Data access contracts defined by business needs
   - **Business Rules**: Validation and domain logic within entities

2. **Application Layer** (`/internal/application/`)
   - **Use Cases**: Single-responsibility business operations
   - **DTOs**: Data Transfer Objects for request/response
   - **Business Logic**: Coordinated operations across domain entities

3. **Infrastructure Layer** (`/internal/infrastructure/`)
   - **Persistence**: Repository implementations using SQLite with sqlx
   - **Database**: Schema managed via `/schema.sql` file
   - **Mappers**: Translation between domain entities and database models

4. **Interface Layer** (`/internal/interface/`)
   - **Setup**: Dependency injection using factory pattern
   - **Adapters**: Interface between external world and application
   - **Router**: HTTP routing and endpoint configuration

### Additional Supporting Layers

5. **Database Schema** (`/schema.sql`)
   - Single source of truth for database structure
   - Enhanced dependency system with 8 built-in types
   - Tables: domains, nodes, attributes, node_attributes, dependency management

6. **MCP Integration** (`/internal/mcp/`)
   - JSON-RPC 2.0 implementation
   - 18 tools with distinctive names (without 'mcp' prefix)
   - Resource system for MCP protocol support
   - Composite key format: `tool-name:domain:id`

7. **Constants Layer** (`/internal/constants/`)
   - Centralized configuration constants
   - Server metadata, network settings, database paths
   - Error messages and validation patterns

### Key Architecture Benefits
- **Dependency Inversion**: Business logic independent of frameworks
- **Testability**: Clean separation enables comprehensive testing
- **Maintainability**: Single responsibility at each layer
- **Extensibility**: Easy to add new features without breaking existing code

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

### Test Automation Scripts

#### Coverage Analysis Script (`./scripts/coverage_analysis.sh`)
Comprehensive coverage analysis with actionable insights:
```bash
./scripts/coverage_analysis.sh    # Full analysis with color-coded output
```

**Features:**
- **Overall Statistics**: Total coverage with status indicators
- **Package-level Analysis**: High/Medium/Low coverage categorization
- **Critical Functions**: Lists 0% coverage functions requiring immediate attention
- **High Potential**: Functions with 75-95% coverage that can easily reach 100%
- **Missing Tests**: Identifies packages without test files
- **HTML Report**: Generates detailed coverage.html report
- **Improvement Suggestions**: Prioritized recommendations for coverage improvement

#### Test Runner Script (`./scripts/test_runner.sh`)
Flexible test execution with multiple modes and options:
```bash
./scripts/test_runner.sh -h                    # Show all options
./scripts/test_runner.sh                       # Run all tests
./scripts/test_runner.sh -m coverage           # Run with coverage analysis
./scripts/test_runner.sh -p internal/mcp -v    # Test specific package with verbose
./scripts/test_runner.sh -m unit --timeout 15m # Unit tests with custom timeout
```

**Modes:**
- `all`: Run all tests (default)
- `unit`: Unit tests only (with -short flag)
- `integration`: Integration tests only (with -tags=integration)
- `mcp`: MCP-specific tests only
- `coverage`: Tests with detailed coverage analysis

**Options:**
- `-v, --verbose`: Detailed test output
- `-p, --package PKG`: Target specific package
- `-t, --timeout TIME`: Custom timeout (default: 10m)
- `--clean`: Clean coverage files before running

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

### Coverage Improvement Strategy
Use the coverage analysis script to identify improvement opportunities:

1. **Immediate Impact** (0% coverage packages):
   - `internal/services/advanced` - 0 test files, high business value
   - `cmd/server/main.go` - MCP adapter functions (lines 355-522)

2. **Quick Wins** (75-95% functions):
   - Functions already mostly tested, small gaps to close
   - Use analysis script to identify specific functions

3. **Medium-term Goals**:
   - Add missing test files for untested packages
   - Improve MCP protocol test coverage (currently 14.4%)
   - Enhance HTTP interface tests (currently 1.5%)

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

### Development Workflow (Clean Architecture)

When implementing new features:
1. **Design Domain Entity**: Define business rules and validation in domain layer
2. **Create Repository Interface**: Define data access contracts in domain layer
3. **Implement Use Case**: Business logic orchestration in application layer
4. **Create Infrastructure**: Repository implementation in infrastructure layer
5. **Wire Dependencies**: Factory pattern for dependency injection in interface layer
6. **Add MCP Tools**: Expose functionality via MCP protocol
7. **Write Tests**: Follow `package_test` separation pattern
8. **Update Documentation**: Document Clean Architecture patterns and examples

#### Clean Architecture Development Process
1. **Domain Modeling**: Identify entities, value objects, business rules
2. **Use Case Design**: Single-responsibility business operations
3. **Interface Definition**: Repository contracts from domain perspective
4. **Implementation**: Infrastructure layer implementations
5. **Integration**: Factory pattern dependency wiring
6. **Testing**: Comprehensive test coverage with `package_test` separation

### MCP Tool Design Guidelines

1. **Tool Naming**: Use clear, action-oriented names (e.g., `create_mcp_domain`, not `domain_new`)
2. **Parameter Design**: 
   - Use descriptive parameter names
   - Required parameters should be minimal
   - Optional parameters for extended functionality
3. **Return Values**: Always return useful information for chaining operations
4. **Error Messages**: Provide clear, actionable error messages
5. **Composite Keys**: Return composite IDs for created/updated resources

### Code Quality Standards

**Current Quality Metrics:**
- **Code Quality Score**: A- (85/100)
- **Architecture Compliance**: A (95/100) - Clean Architecture principles
- **SOLID Principles**: Fully implemented
- **Go Standards**: 100% compliance
- **Function Size**: Average 15 lines (target: <20 lines)
- **Test Coverage**: 20.6% (target: 80%)

**Quality Guidelines:**
1. **Meaningful Names**: Intention-revealing, pronounceable, searchable
2. **Small Functions**: Single responsibility, <20 lines
3. **DRY Principle**: Eliminate code duplication
4. **Error Handling**: Consistent Go idiom patterns
5. **Immutability**: Domain entities with getter methods
6. **Clean Architecture**: Strict layer separation and dependency inversion

### Potential Future Features

Features that could be added for enhanced functionality:
- Comprehensive test suite completion (target: 80% coverage)
- Advanced search/filter capabilities
- Export/import functionality
- Node connections and relationships management
- Subscription and dependency management
- Architecture tests to enforce dependency rules