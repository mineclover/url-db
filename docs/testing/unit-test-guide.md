# Unit Testing Guide - URL Database System

## Overview

This document provides comprehensive guidance for unit testing the URL Database system. All tests are designed to be fast, isolated, and focused on individual components.

## Test Architecture

### Directory Structure
```
url-db/
├── internal/
│   ├── attributes/
│   │   ├── handler_test.go
│   │   ├── repository_test.go
│   │   ├── service_test.go
│   │   └── validators_test.go
│   ├── compositekey/
│   │   ├── service_test.go
│   │   ├── validator_test.go
│   │   └── normalizer_test.go
│   ├── domains/
│   │   ├── handler_test.go
│   │   ├── repository_test.go
│   │   └── service_test.go
│   ├── mcp/
│   │   ├── service_test.go
│   │   ├── converter_test.go
│   │   └── stdio_server_test.go
│   ├── nodes/
│   │   ├── handler_test.go
│   │   ├── repository_test.go
│   │   └── service_test.go
│   └── nodeattributes/
│       ├── handler_test.go
│       ├── repository_test.go
│       └── service_test.go
├── cmd/server/
│   └── main_test.go
└── docs/testing/
    ├── unit-test-guide.md
    ├── test-scenarios.md
    └── mock-data.md
```

## Testing Principles

### 1. Test Isolation
- Each test runs independently
- No shared state between tests
- Use in-memory SQLite for database tests
- Mock external dependencies

### 2. Test Coverage
- **Repositories**: Database operations, error handling
- **Services**: Business logic, validation, error scenarios
- **Handlers**: HTTP request/response, validation, status codes
- **Utils**: Helper functions, edge cases
- **MCP**: Protocol implementation, stdio operations

### 3. Test Categories

#### Repository Tests
- CRUD operations
- Query filtering and pagination
- Constraint violations
- Transaction handling
- Error scenarios

#### Service Tests  
- Business logic validation
- Error handling and propagation
- Integration between components
- Data transformation
- Authorization logic

#### Handler Tests
- HTTP request parsing
- Response formatting
- Status code validation
- Error response structure
- Input validation

#### MCP Tests
- Protocol compliance
- Composite key handling
- Command parsing
- Error responses
- Interactive flows

## Test Utilities

### Mock Framework
Using `github.com/stretchr/testify/mock` for mocking dependencies.

### Test Database
Using in-memory SQLite with migrations for repository tests.

### Test Data
Centralized test data factories for consistent test scenarios.

## Running Tests

### All Tests
```bash
make test
```

### With Coverage
```bash
make test-coverage
```

### Specific Package
```bash
go test ./internal/domains/... -v
go test ./internal/mcp/... -v
```

### Integration Tests
```bash
go test -tags=integration ./...
```

## Test Scenarios

### Core Domain Operations
1. **Domain Management**
   - Create, read, update, delete domains
   - Duplicate domain handling
   - Domain validation

2. **URL Management**
   - Save URLs with metadata
   - Search and filtering
   - URL validation and normalization
   - Duplicate detection

3. **Attribute System**
   - Create attribute definitions
   - Assign attributes to URLs
   - Validate attribute types
   - Order management

4. **MCP Integration**
   - Composite key generation
   - Command parsing and execution
   - Error handling
   - Interactive sessions

### Error Scenarios
1. **Validation Errors**
   - Invalid input data
   - Missing required fields
   - Data type mismatches

2. **Database Errors**
   - Connection failures
   - Constraint violations
   - Transaction rollbacks

3. **Business Logic Errors**
   - Domain not found
   - Duplicate resources
   - Permission denied

## Test Data Management

### Factories
Consistent test data creation using factory functions:

```go
func CreateTestDomain() *models.Domain
func CreateTestNode() *models.Node
func CreateTestAttribute() *models.Attribute
```

### Fixtures
Pre-defined test datasets for complex scenarios:

```go
var TestDomains = []models.Domain{...}
var TestNodes = []models.Node{...}
```

### Cleanup
Automatic cleanup between tests using setup/teardown functions.

## Performance Testing

### Benchmarks
Performance testing for critical operations:

```bash
go test -bench=. ./internal/...
```

### Memory Testing
Memory usage validation:

```bash
go test -memprofile=mem.prof ./internal/...
```

## Continuous Integration

### GitHub Actions
Automated testing on:
- Pull requests
- Main branch commits
- Release tags

### Coverage Requirements
- Minimum 80% code coverage
- All new code must include tests
- Critical paths require 95% coverage

## Best Practices

### Test Naming
```go
func TestServiceName_MethodName_Scenario(t *testing.T)
func TestDomainService_CreateDomain_Success(t *testing.T)
func TestDomainService_CreateDomain_DuplicateName(t *testing.T)
```

### Test Structure
```go
func TestExample(t *testing.T) {
    // Arrange
    // Setup test data and mocks
    
    // Act
    // Execute the function under test
    
    // Assert
    // Verify results and side effects
}
```

### Error Testing
```go
func TestExample_Error(t *testing.T) {
    // Test both error conditions and error messages
    err := service.Method(invalidInput)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "expected error message")
}
```

### Mock Usage
```go
func TestWithMock(t *testing.T) {
    mockRepo := &MockRepository{}
    mockRepo.On("Method", mock.Anything).Return(expectedResult, nil)
    
    service := NewService(mockRepo)
    result, err := service.Method(input)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
    mockRepo.AssertExpectations(t)
}
```

## Testing Tools

### Required Dependencies
```go
github.com/stretchr/testify/assert
github.com/stretchr/testify/require  
github.com/stretchr/testify/mock
github.com/stretchr/testify/suite
```

### Optional Tools
- `go-sqlmock` for advanced database mocking
- `httptest` for HTTP handler testing
- `goleak` for goroutine leak detection

## Maintenance

### Test Updates
- Update tests when APIs change
- Add tests for new features
- Remove tests for deprecated functionality

### Performance Monitoring
- Track test execution time
- Identify slow tests
- Optimize test performance

### Documentation
- Keep test documentation current
- Document complex test scenarios
- Maintain test data examples