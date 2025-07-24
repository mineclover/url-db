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

# Run tests (actual working commands)
go test -v ./...                                    # Run all tests (basic)
go test -v ./internal/mcp/...                      # Run specific package tests
go test -coverprofile=coverage.out ./...          # Generate coverage
go tool cover -html=coverage.out -o coverage.html # Generate HTML coverage report
go tool cover -func=coverage.out                  # Show coverage by function

# Lint and format
make lint                    # Run golangci-lint (install: brew install golangci-lint)
make fmt                     # Format all Go files

# Development mode
make dev                     # Hot reload (requires: go install github.com/cosmtrek/air@latest)

# Run the server
./bin/url-db                 # HTTP mode (default port 8080) 
./bin/url-db -mcp-mode=stdio # MCP stdio mode for AI assistants
./bin/url-db -mcp-mode=sse   # Server-Sent Events mode

# Docker commands (added in latest version)
make docker-build            # Build Docker image
make docker-run              # Run container in MCP stdio mode
make docker-compose-up       # Start all services with Docker Compose
make docker-compose-down     # Stop all services
make docker-logs             # Show Docker logs
make docker-push DOCKER_REGISTRY=your-registry # Push to registry
make docker-clean            # Clean Docker resources

# Build commands
make deps                    # Install dependencies
make build-all              # Build for all platforms
make run                    # Build and run
make clean                  # Clean build artifacts

# Documentation
make swagger-gen            # Generate Swagger documentation
make dev-swagger           # Generate docs and run dev mode

# MCP Tool specification
# Single source: /specs/mcp-tools.yaml - contains all tool definitions, descriptions, and usage info
python scripts/generate-tool-constants.py  # Generate Go constants (only if needed for compile-time constants)
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
node_dependencies          -- Enhanced dependency management
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

The MCP server supports three modes:
- **stdio**: For AI assistants (Claude Desktop, Cursor)
- **http**: For HTTP-based integration
- **sse**: For Server-Sent Events (experimental)

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

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

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

1. **Tool Naming**: Use clear, action-oriented names (e.g., `create_domain`, not `domain_new`)
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

## Docker Deployment

The URL-DB MCP server is available as a Docker image for easy deployment:

### Quick Start with Docker

```bash
# Build Docker image locally
make docker-build

# Run in MCP stdio mode (for AI assistants)
make docker-run

# Run all services with Docker Compose
make docker-compose-up
```

### Docker Hub Image

The project is available on Docker Hub as `asfdassdssa/url-db:latest`:

```bash
# Pull and run from Docker Hub
docker run -it --rm asfdassdssa/url-db:latest

# With persistent data storage
docker run -it --rm -v url-db-data:/data asfdassdssa/url-db:latest

# With host directory mounting
docker run -it --rm -v $(pwd)/database:/data asfdassdssa/url-db:latest
```

### Claude Desktop Integration

Add to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "url-db-data:/data", 
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### Docker Features

- **Multi-stage build**: Optimized 58.8MB Alpine Linux image
- **Non-root user**: Enhanced security with user `urldb`
- **Volume support**: Persistent data storage with `/data` volume
- **Multiple modes**: Supports stdio, HTTP, SSE, and MCP-HTTP modes
- **Host database**: SQLite can be stored on host filesystem for direct access

### Docker Configuration Files

- `Dockerfile`: Multi-stage build configuration
- `docker-compose.yml`: Multi-service deployment
- `docker-compose-host-db.yml`: Host database storage example
- `.dockerignore`: Optimized build context
- `docker-deployment.md`: Comprehensive deployment guide
- `sqlite-host-storage-guide.md`: Host database storage guide

## Configuration and Environment

### Default Settings (from /internal/constants/)
- **Port**: 8080 (configurable via constants)
- **Database**: `file:./url-db.sqlite` (configurable via constants)
- **Tool Name**: `url-db` (configurable via constants)
- **MCP Server Name**: `url-db-mcp-server`
- **Protocol Version**: `2024-11-05`

### Environment Variables
- `VERSION` - Build version (default: 1.0.0)
- `TEST_TIMEOUT` - Test timeout in seconds (default: 300)
- `COVERAGE_THRESHOLD` - Minimum coverage percentage (default: 80)
- `AUTO_CREATE_ATTRIBUTES` - Auto-create attributes if they don't exist (default: true)

### Build Configuration (Makefile)
The Korean-commented Makefile provides comprehensive build automation with color-coded output:
- Multi-platform builds (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64)
- Hot reload development mode with air
- Swagger documentation generation
- Comprehensive linting with golangci-lint
- Clean build artifact management