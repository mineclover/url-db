# Python Integration Tests

This directory contains Python integration tests for the URL-DB MCP server. These tests validate the MCP protocol implementation from an external client perspective.

## Overview

The Python tests serve as **integration tests** that complement the Go unit tests by:
- Testing the MCP protocol from a client perspective
- Validating cross-language compatibility
- Ensuring real-world usage scenarios work correctly
- Testing the complete request/response cycle

## Test Structure

### Test Categories

#### 🚀 Demo Tests
Basic demonstration scripts and simple client tests:
- `server_info_demo.py` - Server information retrieval demo
- `list_domains.py` - Basic MCP client demo

#### 🔧 Basic Tests  
Core functionality tests:
- `test_mcp_client.py` - MCP protocol handshake and basic operations
- `test_node_attributes.py` - Node creation with typed attributes

#### 🛠️ Comprehensive Tests
Complete MCP tool testing:
- `test_all_mcp_tools.py` - Tests all 16 MCP tools
- `test_mcp_tools.py` - Tests 11 core MCP tools
- `test_mcp_final.py` - Safe integration test with error handling

#### 📋 Domain Attribute Tests
Domain schema and attribute management:
- `test_domain_attributes.py` - Domain attribute CRUD (Korean comments)
- `test_mcp_domain_attributes.py` - Domain attribute management
- `test_final.py` - Domain attribute operations with validation
- `test_mcp_persistent.py` - Persistent connection tests

#### 🎯 Advanced Tests
Production-readiness testing:
- `test_scenarios.py` - LLM-as-a-Judge comprehensive testing

## Running Tests

### Quick Start

Run all tests from project root:
```bash
./run_python_tests.sh
```

### Using the Test Runner

From the `tests/python/` directory:

```bash
# Run all tests
python3 run_tests.py

# Run specific category
python3 run_tests.py --category basic
python3 run_tests.py --category comprehensive

# List available tests
python3 run_tests.py --list

# Verbose output
python3 run_tests.py --verbose
```

### Manual Test Execution

Run individual tests:
```bash
cd tests/python/
python3 test_mcp_client.py
python3 test_all_mcp_tools.py
```

## Prerequisites

1. **Server Binary**: The URL-DB server must be built first:
   ```bash
   make build  # Creates bin/url-db
   ```

2. **Python 3**: Tests use Python standard library only, no external dependencies required.

## Test Architecture

### MCP Protocol Testing
All tests use JSON-RPC 2.0 over stdio communication:
- Server started in `stdio` mode
- JSON messages sent/received via stdin/stdout
- Protocol handshake validation
- Tool discovery and execution

### Key Validation Points
- ✅ MCP protocol compliance (JSON-RPC 2.0)
- ✅ Tool registration and discovery
- ✅ Composite key format (`url-db:domain:id`)
- ✅ Domain schema enforcement
- ✅ Error handling and response format
- ✅ Resource management

### Test Data
Tests use:
- In-memory SQLite database (`:memory:`)
- Temporary test data that doesn't persist
- Clean state for each test run

## Development Guidelines

### Adding New Tests
1. Follow existing naming conventions (`test_*.py`)
2. Use standard library only (no external dependencies)
3. Include proper error handling
4. Test both success and failure scenarios
5. Update the test runner categories if needed

### Test Quality Standards
- Clear test descriptions and comments
- Proper setup/teardown procedures
- Comprehensive error checking
- Meaningful assertions
- Good logging for debugging

### Language Standards
- Use English for all comments and documentation
- Korean comments should be translated to English
- Follow Python PEP 8 style guidelines

## Integration with CI/CD

These tests should be run as part of the CI/CD pipeline:

```yaml
# Example GitHub Actions step
- name: Run Python Integration Tests
  run: |
    make build
    ./run_python_tests.sh
```

## Troubleshooting

### Common Issues

1. **Server binary not found**:
   ```bash
   make build
   ```

2. **Permission denied**:
   ```bash
   chmod +x run_python_tests.sh
   chmod +x tests/python/run_tests.py
   ```

3. **Tests timeout**: 
   - Check if server is starting properly
   - Look for port conflicts
   - Verify database permissions

### Debug Mode
Run tests with verbose output:
```bash
python3 run_tests.py --verbose
```

## Contributing

When contributing Python tests:
1. Ensure tests are deterministic and repeatable
2. Include both positive and negative test cases
3. Document any special requirements
4. Update this README for new test categories
5. Maintain compatibility with the Go codebase

## Architecture Benefits

This dual-language testing approach provides:
- **Go Tests**: Fast unit tests, internal component validation
- **Python Tests**: Integration testing, external client validation
- **Complete Coverage**: Both internal and external perspectives
- **Cross-Language Validation**: Ensures protocol works with any client