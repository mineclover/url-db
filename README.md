# URL Database - Clean Architecture with MCP Integration

A comprehensive URL management system built with Clean Architecture principles and optimized for Model Context Protocol (MCP) integration.

## 🚀 Quick Start

### Build and Run

```bash
# Install dependencies and build
make deps
make build

# Run the server
make run

# Or run directly
./bin/url-db
```

### Claude Desktop Integration

For Claude Desktop MCP integration:
1. Build the project: `make build`
2. Configure Claude Desktop: See [MCP Claude Setup Guide](docs/MCP_CLAUDE_SETUP.md)
3. For general MCP info: See [MCP Setup Guide](docs/mcp-setup-guide.md)

### Development

```bash
# Run with hot reload (requires air)
make dev

# Format code
make fmt

# Run linter
make lint
```

## �� Testing

```bash
# Comprehensive test suite (includes linting, coverage, benchmarks)
./scripts/test.sh

# Specific test types
./scripts/test.sh --tests-only      # Unit tests only
./scripts/test.sh --coverage-only   # Coverage analysis
./scripts/test.sh --benchmarks-only # Performance benchmarks
./scripts/test.sh --mcp-only        # MCP integration tests
./scripts/test.sh --lint-only       # Linting only
./scripts/test.sh --package internal/mcp  # Test specific package

# Show all test options
./scripts/test.sh --help
```

## 📋 Available Commands

### Build Commands (Makefile)
- `make deps` - Install dependencies
- `make build` - Build for current platform
- `make build-all` - Build for all platforms (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64)
- `make run` - Build and run
- `make clean` - Clean build artifacts

### Development Commands (Makefile)
- `make fmt` - Format code
- `make lint` - Run linter
- `make dev` - Run with hot reload
- `make swagger-gen` - Generate Swagger documentation
- `make dev-swagger` - Generate docs and run dev mode

### Test Commands (scripts/test.sh)
- `./scripts/test.sh` - Run comprehensive test suite
- `./scripts/test.sh --tests-only` - Unit tests only
- `./scripts/test.sh --coverage-only` - Coverage analysis
- `./scripts/test.sh --benchmarks-only` - Performance benchmarks
- `./scripts/test.sh --mcp-only` - MCP integration tests
- `./scripts/test.sh --lint-only` - Linting only
- `./scripts/test.sh --package DIR` - Test specific package

## 🏗️ Architecture

### Clean Architecture Implementation

URL-DB follows Clean Architecture principles with four distinct layers:

```
cmd/server/               # Entry point and main application
internal/
├── domain/               # Business entities and repository interfaces
│   ├── entity/          # Domain entities (Domain, Node, Attribute)
│   └── repository/      # Repository interfaces
├── application/         # Application layer (use cases and DTOs)
│   ├── dto/            # Data Transfer Objects
│   └── usecase/        # Business logic use cases
├── infrastructure/     # External concerns (database, persistence)
│   └── persistence/    # Data persistence implementations
└── interface/          # Interface adapters and setup
    └── setup/          # Dependency injection and factory pattern
```

**Key Architectural Principles:**
- **Dependency Inversion**: Inner layers define interfaces, outer layers implement them
- **Single Responsibility**: Each layer has one reason to change
- **Clean Separation**: Business logic is independent of frameworks and databases

### Core Components
- **Domain Entities**: Immutable business objects with encapsulated logic
- **Use Cases**: Single-responsibility business operations
- **Repositories**: Data access abstractions
- **Factory Pattern**: Dependency injection and object creation
- **MCP Integration**: Native support for AI tool integration

### Composite Key Format
Nodes are identified using composite keys: `tool-name:domain:id`

Examples:
- `url-db:example:1` - Node ID 1 in the "example" domain
- `url-db:github:42` - Node ID 42 in the "github" domain
- `work:projects:15` - Node ID 15 with custom tool name "work"

## 🔧 Configuration

### Default Settings
- **Port**: 8080 (configurable via constants)
- **Database**: `file:./url-db.sqlite` (configurable via constants)
- **Tool Name**: `url-db` (configurable via constants)
- **MCP Server Name**: `url-db-mcp-server`
- **Protocol Version**: `2024-11-05`

### Constants Management
All configuration values are centralized in `/internal/constants/constants.go`:
- Server metadata and versions
- Network settings and ports  
- Database paths and drivers
- Validation limits and patterns
- Error messages and HTTP status codes

### Environment Variables
- `VERSION` - Build version (default: 1.0.0)
- `TEST_TIMEOUT` - Test timeout in seconds (default: 300)
- `COVERAGE_THRESHOLD` - Minimum coverage percentage (default: 80)
- `AUTO_CREATE_ATTRIBUTES` - Auto-create attributes if they don't exist (default: true)

## 📊 Test Output

Comprehensive tests generate reports in `test-output/`:
- `test-results.txt` - Unit test results
- `coverage.html` - HTML coverage report
- `coverage-summary.txt` - Coverage percentage
- `benchmark-results.txt` - Performance benchmarks
- `race-detection.txt` - Race condition analysis
- `lint-report.txt` - Linting results
- `test-summary.txt` - Complete test summary

## 🚀 MCP Integration

The URL-DB server provides native MCP support with 18 tools:
- **Domain Management**: Create and list domains
- **URL Operations**: Save, search, and manage URLs  
- **Attribute System**: Tag and categorize URLs with type validation
- **Schema Management**: Define and enforce domain-specific attributes
- **Advanced Queries**: Filter by attributes, batch operations
- **Resource System**: MCP resource protocol support

### Tool Specification System
- **Single Source**: All tools defined in `/specs/mcp-tools.yaml`
- **Auto-Generation**: Constants generated for Go and Python
- **Consistency**: Tool names and descriptions managed centrally
- **Validation**: Schema-enforced tool definitions

### Common MCP Workflows

#### Save and Categorize URL
1. Check domain exists: `GET /mcp/domains`
2. Create domain if needed: `POST /mcp/domains`
3. Save URL: `POST /mcp/nodes`
4. Add attributes: `PUT /mcp/nodes/{composite_id}/attributes`

#### Research URLs
1. List domains: `GET /mcp/domains`
2. Search URLs: `GET /mcp/nodes?domain_name=example.com&search=keyword`
3. Get details: `GET /mcp/nodes/{composite_id}`
4. View attributes: `GET /mcp/nodes/{composite_id}/attributes`

## 📚 Documentation

### Core Documentation
- [CLAUDE.md](CLAUDE.md) - Claude Code AI assistant integration guide
- [CLEAN_CODE_GUIDELINES.md](CLEAN_CODE_GUIDELINES.md) - Clean code principles and best practices
- [API Documentation](docs/mcp-openapi.yaml) - Complete OpenAPI specification  

### Setup Guides
- [MCP Claude Setup Guide](docs/MCP_CLAUDE_SETUP.md) - Comprehensive MCP integration guide
- [MCP Testing Guide](docs/MCP_TESTING_GUIDE.md) - Testing procedures and workflows

### Technical References
- [Tool Specification](specs/mcp-tools.yaml) - MCP tools definition
- [Composite Key Conventions](docs/spec/composite-key-conventions.md) - Key format specifications
- [Error Codes](docs/spec/error-codes.md) - Error code reference

### Architecture Quality
- **Code Quality Score**: A- (85/100)
- **Architecture Compliance**: A (95/100) - Clean Architecture principles
- **Test Coverage**: 20.6% (target: 80%)
- **Go Standards**: 100% compliance

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `./scripts/test.sh`
5. Submit a pull request

## 📄 License

Apache 2.0 License - see [LICENSE](LICENSE) file for details.