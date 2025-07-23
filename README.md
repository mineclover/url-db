# URL Database - MCP Enhanced

A comprehensive URL management system optimized for Model Context Protocol (MCP) integration.

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

### Core Components
- **Domains**: Organize URLs by website or topic
- **Nodes**: Individual URLs with rich metadata
- **Attributes**: Flexible tagging system (tags, numbers, text, markdown, images)
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

- [CLAUDE.md](CLAUDE.md) - Claude Code AI assistant integration guide
- [API Documentation](docs/mcp-openapi.yaml) - Complete OpenAPI specification  
- [MCP Setup Guide (Korean)](docs/mcp-setup-guide-ko.md) - MCP integration guide
- [MCP Setup Guide (English)](docs/mcp-setup-guide.md) - MCP integration guide
- [Scripts README](scripts/README.md) - Maintenance and utility scripts guide
- [Tool Specification](specs/mcp-tools.yaml) - MCP tools definition

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `./scripts/test.sh`
5. Submit a pull request

## 📄 License

Apache 2.0 License - see [LICENSE](LICENSE) file for details.