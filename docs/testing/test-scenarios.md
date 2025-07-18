# Test Scenarios - URL Database System

## Test Scenario Categories

### 1. Domain Management Tests

#### Happy Path Scenarios
- **Create Domain**: Successfully create a new domain with valid name and description
- **List Domains**: Retrieve all domains with pagination
- **Get Domain**: Retrieve specific domain by ID
- **Update Domain**: Modify domain description
- **Delete Domain**: Remove empty domain

#### Error Scenarios
- **Duplicate Domain**: Attempt to create domain with existing name
- **Invalid Name**: Create domain with empty or too long name
- **Domain Not Found**: Access non-existent domain
- **Delete Non-Empty Domain**: Attempt to delete domain with URLs

#### Edge Cases
- **Special Characters**: Domain names with special characters
- **Unicode Names**: International domain names
- **Maximum Length**: Names at character limits
- **Null Values**: Handling null/empty descriptions

### 2. URL (Node) Management Tests

#### Happy Path Scenarios
- **Save URL**: Add new URL with title and description
- **List URLs**: Retrieve URLs with pagination and filtering
- **Search URLs**: Find URLs by content search
- **Update URL**: Modify title and description
- **Delete URL**: Remove URL and its attributes

#### Error Scenarios
- **Invalid URL**: Malformed URL strings
- **Duplicate URL**: Same URL in same domain
- **URL Too Long**: URLs exceeding length limits
- **Domain Not Found**: Add URL to non-existent domain

#### Edge Cases
- **Special URLs**: Data URLs, file URLs, custom schemes
- **Unicode URLs**: International URLs with non-ASCII characters
- **Query Parameters**: URLs with complex query strings
- **Fragments**: URLs with hash fragments

### 3. Attribute System Tests

#### Happy Path Scenarios
- **Create Attribute**: Define new attribute type for domain
- **List Attributes**: Get all attributes for domain
- **Assign Attribute**: Add attribute value to URL
- **Update Attribute Value**: Modify existing attribute value
- **Delete Attribute**: Remove attribute definition

#### Error Scenarios
- **Invalid Type**: Unsupported attribute types
- **Value Validation**: Invalid values for attribute types
- **Attribute Not Found**: Reference non-existent attributes
- **Type Mismatch**: Wrong value type for attribute

#### Attribute Type Tests
- **Tag Attributes**: Simple string labels
- **Number Attributes**: Numeric values with validation
- **String Attributes**: Text content with length limits
- **Markdown Attributes**: Rich text content
- **Image Attributes**: Image URLs and validation
- **Ordered Tag Attributes**: Tags with ordering

### 4. MCP Integration Tests

#### Happy Path Scenarios
- **Server Info**: Get MCP server capabilities
- **List MCP Domains**: Retrieve domains in MCP format
- **Create MCP Node**: Add URL using MCP protocol
- **Get MCP Node**: Retrieve URL by composite ID
- **Update MCP Node**: Modify URL metadata
- **Batch Operations**: Bulk URL retrieval

#### Error Scenarios
- **Invalid Composite ID**: Malformed composite keys
- **MCP Protocol Errors**: Invalid request formats
- **Permission Denied**: Unauthorized operations
- **Resource Not Found**: Missing resources

#### Stdio Mode Tests
- **Command Parsing**: Parse command-line inputs
- **Interactive Session**: Multi-command sessions
- **Help System**: Built-in help commands
- **Error Handling**: Invalid commands and parameters
- **Session Termination**: Graceful exit

### 5. Search and Filtering Tests

#### Search Functionality
- **Content Search**: Search in titles and descriptions
- **Partial Matching**: Substring searches
- **Case Insensitive**: Case-insensitive search
- **Special Characters**: Search with special characters
- **Empty Results**: Searches with no matches

#### Filtering Options
- **Domain Filtering**: Filter by specific domains
- **Attribute Filtering**: Filter by attribute values
- **Date Ranges**: Filter by creation/update dates
- **Combined Filters**: Multiple filter criteria

#### Pagination Tests
- **Standard Pagination**: Page and size parameters
- **Edge Cases**: First/last pages, empty pages
- **Large Datasets**: Performance with many results
- **Invalid Parameters**: Negative or zero values

### 6. Validation Tests

#### Input Validation
- **Required Fields**: Missing mandatory fields
- **Field Lengths**: Maximum length constraints
- **Data Types**: Type validation and conversion
- **Format Validation**: URL, email, date formats

#### Business Rule Validation
- **Uniqueness Constraints**: Duplicate prevention
- **Reference Integrity**: Foreign key validation
- **State Validation**: Valid state transitions
- **Authorization Rules**: Permission checks

### 7. Error Handling Tests

#### HTTP Error Responses
- **400 Bad Request**: Invalid input data
- **404 Not Found**: Missing resources
- **409 Conflict**: Constraint violations
- **500 Internal Server Error**: System errors

#### Error Message Quality
- **Descriptive Messages**: Clear error descriptions
- **Error Codes**: Consistent error coding
- **Localization**: Multi-language error messages
- **Debug Information**: Development vs production

### 8. Performance Tests

#### Response Time Tests
- **Database Operations**: Query performance
- **API Endpoints**: Response time measurement
- **Search Operations**: Search performance
- **Batch Operations**: Bulk operation efficiency

#### Load Tests
- **Concurrent Requests**: Multiple simultaneous operations
- **Large Datasets**: Performance with many records
- **Memory Usage**: Memory consumption patterns
- **Resource Cleanup**: Proper resource management

### 9. Integration Scenarios

#### Cross-Component Tests
- **Domain-URL Integration**: URLs within domains
- **URL-Attribute Integration**: Attributes on URLs
- **Search Integration**: Search across all components
- **MCP Integration**: End-to-end MCP workflows

#### Database Integration
- **Transaction Handling**: Multi-table operations
- **Constraint Enforcement**: Database-level validation
- **Migration Testing**: Schema change compatibility
- **Backup/Recovery**: Data persistence tests

### 10. Security Tests

#### Input Sanitization
- **SQL Injection**: Database injection attempts
- **XSS Prevention**: Script injection in content
- **Command Injection**: System command injection
- **Path Traversal**: File system access attempts

#### Data Protection
- **Sensitive Data**: Handling of sensitive information
- **Access Control**: Permission-based access
- **Audit Logging**: Operation tracking
- **Data Validation**: Comprehensive input validation

## Test Data Requirements

### Domain Test Data
```yaml
domains:
  - name: "example.com"
    description: "Example domain for testing"
  - name: "test-site.org"
    description: "Test organization site"
  - name: "unicode-测试.com"
    description: "Unicode domain name test"
```

### URL Test Data
```yaml
urls:
  - url: "https://example.com/page1"
    title: "Example Page 1"
    description: "First example page"
  - url: "https://example.com/page2"
    title: "Example Page 2"
    description: "Second example page"
  - url: "https://test-site.org/article"
    title: "Test Article"
    description: "Article for testing"
```

### Attribute Test Data
```yaml
attributes:
  - name: "category"
    type: "tag"
    description: "Content category"
  - name: "rating"
    type: "number"
    description: "User rating (1-5)"
  - name: "notes"
    type: "string"
    description: "Additional notes"
  - name: "priority"
    type: "ordered_tag"
    description: "Priority level"
```

## Expected Test Coverage

### Coverage Targets
- **Repository Layer**: 95% - Critical data operations
- **Service Layer**: 90% - Business logic coverage
- **Handler Layer**: 85% - HTTP endpoint coverage
- **Utility Functions**: 90% - Helper function coverage
- **MCP Integration**: 85% - Protocol implementation

### Critical Paths
- User authentication and authorization
- Data validation and sanitization
- Error handling and recovery
- Database transaction management
- MCP protocol compliance