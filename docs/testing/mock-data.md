# Mock Data for Testing

## Overview

This document provides standardized mock data for use across all unit tests in the URL Database system. Using consistent test data ensures reliable and predictable test outcomes.

## Test Data Categories

### 1. Domain Test Data

#### Standard Domains
```go
var TestDomains = []models.Domain{
    {
        ID:          1,
        Name:        "example.com",
        Description: "Example domain for testing",
        CreatedAt:   "2024-01-01T00:00:00Z",
        UpdatedAt:   "2024-01-01T00:00:00Z",
    },
    {
        ID:          2,
        Name:        "test-site.org",
        Description: "Test organization website",
        CreatedAt:   "2024-01-02T00:00:00Z",
        UpdatedAt:   "2024-01-02T00:00:00Z",
    },
    {
        ID:          3,
        Name:        "api.service.com",
        Description: "API service documentation",
        CreatedAt:   "2024-01-03T00:00:00Z",
        UpdatedAt:   "2024-01-03T00:00:00Z",
    },
}
```

#### Edge Case Domains
```go
var EdgeCaseDomains = []models.Domain{
    {
        ID:          100,
        Name:        "unicode-测试.com",
        Description: "Unicode domain name test",
        CreatedAt:   "2024-01-10T00:00:00Z",
        UpdatedAt:   "2024-01-10T00:00:00Z",
    },
    {
        ID:          101,
        Name:        "very-long-domain-name.with.multiple.subdomains.example.com",
        Description: "Very long domain name for testing limits",
        CreatedAt:   "2024-01-11T00:00:00Z",
        UpdatedAt:   "2024-01-11T00:00:00Z",
    },
    {
        ID:          102,
        Name:        "192.168.1.1",
        Description: "IP address as domain",
        CreatedAt:   "2024-01-12T00:00:00Z",
        UpdatedAt:   "2024-01-12T00:00:00Z",
    },
}
```

### 2. URL/Node Test Data

#### Standard URLs
```go
var TestNodes = []models.Node{
    {
        ID:          1,
        DomainID:    1,
        URL:         "https://example.com/page1",
        Title:       "Example Page 1",
        Description: "First example page for testing",
        Content:     "https://example.com/page1",
        CreatedAt:   "2024-01-01T01:00:00Z",
        UpdatedAt:   "2024-01-01T01:00:00Z",
    },
    {
        ID:          2,
        DomainID:    1,
        URL:         "https://example.com/page2",
        Title:       "Example Page 2",
        Description: "Second example page for testing",
        Content:     "https://example.com/page2",
        CreatedAt:   "2024-01-01T02:00:00Z",
        UpdatedAt:   "2024-01-01T02:00:00Z",
    },
    {
        ID:          3,
        DomainID:    2,
        URL:         "https://test-site.org/article",
        Title:       "Test Article",
        Description: "Article for testing search functionality",
        Content:     "https://test-site.org/article",
        CreatedAt:   "2024-01-02T01:00:00Z",
        UpdatedAt:   "2024-01-02T01:00:00Z",
    },
}
```

#### Edge Case URLs
```go
var EdgeCaseNodes = []models.Node{
    {
        ID:          100,
        DomainID:    1,
        URL:         "https://example.com/very-long-url-path/with/multiple/segments/and/query?param1=value1&param2=value2&param3=very-long-value#fragment",
        Title:       "Very Long URL",
        Description: "URL with very long path and query parameters",
        Content:     "https://example.com/very-long-url-path/with/multiple/segments/and/query?param1=value1&param2=value2&param3=very-long-value#fragment",
        CreatedAt:   "2024-01-10T01:00:00Z",
        UpdatedAt:   "2024-01-10T01:00:00Z",
    },
    {
        ID:          101,
        DomainID:    100,
        URL:         "https://unicode-测试.com/页面",
        Title:       "Unicode URL",
        Description: "URL with unicode characters",
        Content:     "https://unicode-测试.com/页面",
        CreatedAt:   "2024-01-11T01:00:00Z",
        UpdatedAt:   "2024-01-11T01:00:00Z",
    },
    {
        ID:          102,
        DomainID:    1,
        URL:         "https://example.com/special-chars!@#$%^&*()_+",
        Title:       "Special Characters URL",
        Description: "URL with special characters",
        Content:     "https://example.com/special-chars!@#$%^&*()_+",
        CreatedAt:   "2024-01-12T01:00:00Z",
        UpdatedAt:   "2024-01-12T01:00:00Z",
    },
}
```

### 3. Attribute Test Data

#### Standard Attributes
```go
var TestAttributes = []models.Attribute{
    {
        ID:          1,
        DomainID:    1,
        Name:        "category",
        Type:        models.AttributeTypeTag,
        Description: "Content category",
        CreatedAt:   "2024-01-01T00:30:00Z",
    },
    {
        ID:          2,
        DomainID:    1,
        Name:        "rating",
        Type:        models.AttributeTypeNumber,
        Description: "User rating (1-5)",
        CreatedAt:   "2024-01-01T00:31:00Z",
    },
    {
        ID:          3,
        DomainID:    1,
        Name:        "notes",
        Type:        models.AttributeTypeString,
        Description: "Additional notes",
        CreatedAt:   "2024-01-01T00:32:00Z",
    },
    {
        ID:          4,
        DomainID:    1,
        Name:        "priority",
        Type:        models.AttributeTypeOrderedTag,
        Description: "Priority level with ordering",
        CreatedAt:   "2024-01-01T00:33:00Z",
    },
    {
        ID:          5,
        DomainID:    1,
        Name:        "content",
        Type:        models.AttributeTypeMarkdown,
        Description: "Rich content in markdown",
        CreatedAt:   "2024-01-01T00:34:00Z",
    },
    {
        ID:          6,
        DomainID:    1,
        Name:        "thumbnail",
        Type:        models.AttributeTypeImage,
        Description: "Thumbnail image URL",
        CreatedAt:   "2024-01-01T00:35:00Z",
    },
}
```

### 4. Node Attribute Test Data

#### Standard Node Attributes
```go
var TestNodeAttributes = []models.NodeAttribute{
    {
        ID:          1,
        NodeID:      1,
        AttributeID: 1,
        Value:       "tutorial",
        OrderIndex:  1,
        CreatedAt:   "2024-01-01T01:30:00Z",
    },
    {
        ID:          2,
        NodeID:      1,
        AttributeID: 2,
        Value:       "5",
        OrderIndex:  1,
        CreatedAt:   "2024-01-01T01:31:00Z",
    },
    {
        ID:          3,
        NodeID:      1,
        AttributeID: 3,
        Value:       "This is a great tutorial page",
        OrderIndex:  1,
        CreatedAt:   "2024-01-01T01:32:00Z",
    },
    {
        ID:          4,
        NodeID:      2,
        AttributeID: 1,
        Value:       "documentation",
        OrderIndex:  1,
        CreatedAt:   "2024-01-01T02:30:00Z",
    },
    {
        ID:          5,
        NodeID:      2,
        AttributeID: 4,
        Value:       "high",
        OrderIndex:  1,
        CreatedAt:   "2024-01-01T02:31:00Z",
    },
}
```

### 5. MCP Test Data

#### MCP Domains
```go
var TestMCPDomains = []models.MCPDomain{
    {
        Name:        "example.com",
        Description: "Example domain for MCP testing",
        NodeCount:   5,
        CreatedAt:   "2024-01-01T00:00:00Z",
        UpdatedAt:   "2024-01-01T00:00:00Z",
    },
    {
        Name:        "test-site.org",
        Description: "Test site for MCP integration",
        NodeCount:   3,
        CreatedAt:   "2024-01-02T00:00:00Z",
        UpdatedAt:   "2024-01-02T00:00:00Z",
    },
}
```

#### MCP Nodes
```go
var TestMCPNodes = []models.MCPNode{
    {
        CompositeID:  "example.com::https://example.com/page1",
        DomainName:   "example.com",
        URL:          "https://example.com/page1",
        Title:        "Example Page 1",
        Description:  "First example page",
        CreatedAt:    "2024-01-01T01:00:00Z",
        UpdatedAt:    "2024-01-01T01:00:00Z",
    },
    {
        CompositeID:  "example.com::https://example.com/page2",
        DomainName:   "example.com",
        URL:          "https://example.com/page2",
        Title:        "Example Page 2",
        Description:  "Second example page",
        CreatedAt:    "2024-01-01T02:00:00Z",
        UpdatedAt:    "2024-01-01T02:00:00Z",
    },
}
```

#### MCP Attributes
```go
var TestMCPAttributes = []models.MCPAttribute{
    {
        Name:  "category",
        Type:  "tag",
        Value: "tutorial",
    },
    {
        Name:  "rating",
        Type:  "number",
        Value: "5",
    },
    {
        Name:  "notes",
        Type:  "string",
        Value: "Excellent tutorial content",
    },
}
```

### 6. Request/Response Test Data

#### Create Domain Requests
```go
var TestCreateDomainRequests = []models.CreateDomainRequest{
    {
        Name:        "new-domain.com",
        Description: "New domain for testing",
    },
    {
        Name:        "api.example.com",
        Description: "API documentation domain",
    },
}
```

#### Create Node Requests
```go
var TestCreateNodeRequests = []models.CreateNodeRequest{
    {
        URL:         "https://example.com/new-page",
        Title:       "New Page",
        Description: "Newly created page",
    },
    {
        URL:         "https://example.com/api/docs",
        Title:       "API Documentation",
        Description: "Complete API documentation",
    },
}
```

#### Create MCP Node Requests
```go
var TestCreateMCPNodeRequests = []models.CreateMCPNodeRequest{
    {
        DomainName:  "example.com",
        URL:         "https://example.com/mcp-page",
        Title:       "MCP Test Page",
        Description: "Page created via MCP",
    },
}
```

### 7. Error Test Data

#### Invalid Domain Data
```go
var InvalidDomainData = []struct {
    Name        string
    Request     models.CreateDomainRequest
    ExpectedErr string
}{
    {
        Name: "Empty domain name",
        Request: models.CreateDomainRequest{
            Name:        "",
            Description: "Valid description",
        },
        ExpectedErr: "domain name is required",
    },
    {
        Name: "Domain name too long",
        Request: models.CreateDomainRequest{
            Name:        strings.Repeat("a", 256),
            Description: "Valid description",
        },
        ExpectedErr: "domain name too long",
    },
}
```

#### Invalid URL Data
```go
var InvalidURLData = []struct {
    Name        string
    Request     models.CreateNodeRequest
    ExpectedErr string
}{
    {
        Name: "Empty URL",
        Request: models.CreateNodeRequest{
            URL:         "",
            Title:       "Valid title",
            Description: "Valid description",
        },
        ExpectedErr: "URL is required",
    },
    {
        Name: "Invalid URL format",
        Request: models.CreateNodeRequest{
            URL:         "not-a-valid-url",
            Title:       "Valid title",
            Description: "Valid description",
        },
        ExpectedErr: "invalid URL format",
    },
}
```

## Test Data Factories

### Domain Factory
```go
func CreateTestDomain(overrides ...func(*models.Domain)) *models.Domain {
    domain := &models.Domain{
        ID:          1,
        Name:        "example.com",
        Description: "Test domain",
        CreatedAt:   time.Now().Format(time.RFC3339),
        UpdatedAt:   time.Now().Format(time.RFC3339),
    }
    
    for _, override := range overrides {
        override(domain)
    }
    
    return domain
}
```

### Node Factory
```go
func CreateTestNode(overrides ...func(*models.Node)) *models.Node {
    node := &models.Node{
        ID:          1,
        DomainID:    1,
        URL:         "https://example.com/page",
        Title:       "Test Page",
        Description: "Test description",
        Content:     "https://example.com/page",
        CreatedAt:   time.Now().Format(time.RFC3339),
        UpdatedAt:   time.Now().Format(time.RFC3339),
    }
    
    for _, override := range overrides {
        override(node)
    }
    
    return node
}
```

### Attribute Factory
```go
func CreateTestAttribute(overrides ...func(*models.Attribute)) *models.Attribute {
    attribute := &models.Attribute{
        ID:          1,
        DomainID:    1,
        Name:        "category",
        Type:        models.AttributeTypeTag,
        Description: "Test attribute",
        CreatedAt:   time.Now().Format(time.RFC3339),
    }
    
    for _, override := range overrides {
        override(attribute)
    }
    
    return attribute
}
```

## Usage Examples

### Using Standard Test Data
```go
func TestSomething(t *testing.T) {
    // Use predefined test domain
    domain := TestDomains[0]
    
    // Use predefined test node
    node := TestNodes[0]
    
    // Test logic here
}
```

### Using Factories with Overrides
```go
func TestCustomScenario(t *testing.T) {
    // Create custom domain
    domain := CreateTestDomain(func(d *models.Domain) {
        d.Name = "custom.com"
        d.Description = "Custom test domain"
    })
    
    // Create custom node
    node := CreateTestNode(func(n *models.Node) {
        n.DomainID = domain.ID
        n.URL = "https://custom.com/page"
    })
    
    // Test logic here
}
```

### Database Setup with Test Data
```go
func setupTestData(db *sql.DB) error {
    // Insert test domains
    for _, domain := range TestDomains {
        _, err := db.Exec(
            "INSERT INTO domains (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
            domain.ID, domain.Name, domain.Description, domain.CreatedAt, domain.UpdatedAt,
        )
        if err != nil {
            return err
        }
    }
    
    // Insert test nodes
    for _, node := range TestNodes {
        _, err := db.Exec(
            "INSERT INTO nodes (id, domain_id, url, title, description, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
            node.ID, node.DomainID, node.URL, node.Title, node.Description, node.Content, node.CreatedAt, node.UpdatedAt,
        )
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

## Best Practices

1. **Consistency**: Use the same test data across related tests
2. **Isolation**: Each test should set up its own data or use factories
3. **Cleanup**: Clean up test data after each test
4. **Realistic Data**: Use realistic URLs, domain names, and content
5. **Edge Cases**: Include edge case data for thorough testing
6. **Maintainability**: Keep test data definitions centralized and versioned