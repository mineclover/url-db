# MCP Server LLM-as-a-Judge Testing Scenarios

## Overview

This document defines comprehensive testing scenarios for the URL-DB MCP server using the LLM-as-a-Judge methodology. Each scenario includes detailed criteria, expected behaviors, and evaluation rubrics.

## Testing Context

**Server Under Test**: URL-DB MCP Server (JSON-RPC 2.0)
**Test Client**: Python MCP client (`test_mcp_client.py`)
**MCP Tool Registration**: Assumed to be registered with Claude Desktop or compatible MCP client
**Reference Context**: [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)

## Test Scenarios

### Scenario 1: MCP Protocol Handshake Compliance

**Objective**: Verify complete MCP protocol initialization sequence

**Test Steps**:
1. Send `initialize` request with protocol version "2024-11-05"
2. Verify server responds with proper capabilities and server info
3. Send `initialized` notification
4. Verify server transitions to ready state

**Expected Behavior**:
- Initialize response includes: `protocolVersion`, `capabilities`, `serverInfo`
- Server capabilities include `tools` and `resources`
- Server name is "url-db-mcp-server" with version "1.0.0"
- No errors during handshake sequence

**Evaluation Criteria** (Score 1-10):
- **Protocol Compliance** (3 points): Exact adherence to MCP handshake specification
- **Response Format** (2 points): Proper JSON-RPC 2.0 response structure
- **Capability Declaration** (3 points): Accurate capability advertisement
- **Error Handling** (2 points): Graceful handling of malformed requests

**Pass Threshold**: 8/10

---

### Scenario 2: Tool Discovery and Schema Validation

**Objective**: Verify all 11 MCP tools are properly registered with valid schemas

**Test Steps**:
1. Call `tools/list` after initialization
2. Verify all expected tools are present
3. Validate input schemas for each tool
4. Check tool descriptions and metadata

**Expected Tools**:
1. `list_mcp_domains` - List all domains
2. `create_mcp_domain` - Create new domain
3. `list_mcp_nodes` - List nodes in domain
4. `create_mcp_node` - Create new node/URL
5. `get_mcp_node` - Get node by composite ID
6. `update_mcp_node` - Update node metadata
7. `delete_mcp_node` - Delete node
8. `find_mcp_node_by_url` - Find node by URL
9. `get_mcp_node_attributes` - Get node attributes
10. `set_mcp_node_attributes` - Set node attributes
11. `get_mcp_server_info` - Get server information

**Evaluation Criteria** (Score 1-10):
- **Tool Completeness** (3 points): All 11 tools present and accessible
- **Schema Validity** (3 points): Valid JSON schemas with proper types and constraints
- **Documentation Quality** (2 points): Clear, helpful tool descriptions
- **Parameter Validation** (2 points): Required/optional parameters correctly specified

**Pass Threshold**: 8/10

---

### Scenario 3: Domain Management Workflow

**Objective**: Test complete domain lifecycle management

**Test Steps**:
1. List initial domains (should be empty)
2. Create test domain "test-scenario-domain"
3. Verify domain creation with proper metadata
4. List domains again to confirm presence
5. Attempt to create duplicate domain (should fail gracefully)

**Expected Behavior**:
- Initial domain list returns empty array
- Domain creation returns domain object with timestamps
- Duplicate creation returns appropriate error
- Domain appears in subsequent listings

**Evaluation Criteria** (Score 1-10):
- **CRUD Operations** (4 points): Create, read operations work correctly
- **Data Integrity** (2 points): Proper timestamps, metadata handling
- **Error Handling** (2 points): Duplicate prevention, clear error messages
- **State Consistency** (2 points): Domain persists across operations

**Pass Threshold**: 7/10

---

### Scenario 4: Node/URL Management with Composite Keys

**Objective**: Verify URL storage and retrieval using composite key system

**Test Steps**:
1. Create domain "url-test-domain"
2. Create node with URL "https://example.com/test-page"
3. Verify composite key format "url-db:url-test-domain:N"
4. Retrieve node by composite key
5. Update node title and description
6. Find node by URL search
7. Delete node and verify removal

**Expected Behavior**:
- Node creation returns composite key in correct format
- Node retrieval by composite key returns full metadata
- Updates persist and are reflected in subsequent retrievals
- URL search finds the correct node
- Deletion removes node completely

**Evaluation Criteria** (Score 1-10):
- **Composite Key System** (3 points): Proper format and uniqueness
- **CRUD Completeness** (3 points): All operations work correctly
- **Data Persistence** (2 points): Updates persist across operations
- **Search Functionality** (2 points): URL-based search works accurately

**Pass Threshold**: 8/10

---

### Scenario 5: Resource System Integration

**Objective**: Test MCP resource discovery and access

**Test Steps**:
1. Call `resources/list` to discover available resources
2. Verify expected resource URIs are present
3. Read server info resource: `mcp://server/info`
4. Read domain resource: `mcp://domains/{domain_name}`
5. Read domain nodes resource: `mcp://domains/{domain_name}/nodes`

**Expected Resources**:
- `mcp://server/info` - Server metadata and capabilities
- `mcp://domains/{domain}` - Domain information
- `mcp://domains/{domain}/nodes` - Domain node listings

**Evaluation Criteria** (Score 1-10):
- **Resource Discovery** (3 points): All expected resources are listed
- **Resource Access** (3 points): Resources return proper JSON content
- **URI Format** (2 points): Consistent, predictable URI patterns
- **Content Quality** (2 points): Resources contain useful, accurate data

**Pass Threshold**: 7/10

---

### Scenario 6: Attribute Management System

**Objective**: Test node attribute assignment and retrieval

**Test Steps**:
1. Create domain and node for attribute testing
2. Set multiple attributes on the node
3. Retrieve node attributes
4. Verify attribute types and values
5. Update attribute values
6. Remove attributes

**Test Attributes**:
- `category`: "testing"
- `priority`: "high"
- `description`: "Test node for attribute validation"

**Expected Behavior**:
- Attributes can be set in batch operations
- Attribute retrieval returns all assigned attributes
- Attribute updates modify existing values
- Invalid attribute operations return clear errors

**Evaluation Criteria** (Score 1-10):
- **Attribute Operations** (4 points): Set, get, update operations work
- **Data Types** (2 points): Proper handling of different value types
- **Batch Processing** (2 points): Multiple attributes handled efficiently
- **Error Handling** (2 points): Clear errors for invalid operations

**Pass Threshold**: 7/10

---

### Scenario 7: Error Handling and Edge Cases

**Objective**: Verify robust error handling across all operations

**Test Steps**:
1. Send malformed JSON-RPC requests
2. Call non-existent tools
3. Provide invalid parameters to tools
4. Access non-existent resources
5. Attempt operations on non-existent entities

**Expected Behavior**:
- Malformed requests return proper JSON-RPC error responses
- Unknown tools return "Method not found" errors
- Invalid parameters return validation errors
- Non-existent resources return appropriate 404-style errors
- All errors include helpful error messages

**Evaluation Criteria** (Score 1-10):
- **JSON-RPC Compliance** (3 points): Proper error response format
- **Error Specificity** (3 points): Specific, actionable error messages
- **Edge Case Coverage** (2 points): Handles various edge cases
- **System Stability** (2 points): Errors don't crash or corrupt state

**Pass Threshold**: 8/10

---

### Scenario 8: Performance and Scalability

**Objective**: Evaluate server performance under realistic loads

**Test Steps**:
1. Create 10 domains rapidly
2. Create 100 nodes across domains
3. Perform batch attribute operations
4. Execute concurrent tool calls
5. Measure response times and resource usage

**Expected Behavior**:
- Operations complete within reasonable timeframes
- No memory leaks or resource exhaustion
- Concurrent operations handled correctly
- Database operations remain performant

**Evaluation Criteria** (Score 1-10):
- **Response Times** (3 points): Sub-second response for typical operations
- **Concurrency** (3 points): Handles concurrent requests correctly
- **Resource Usage** (2 points): Reasonable memory and CPU usage
- **Scalability** (2 points): Performance doesn't degrade significantly

**Pass Threshold**: 6/10

---

## Overall Evaluation Rubric

### Scoring System
- **Excellent** (9-10): Exceeds expectations, production-ready
- **Good** (7-8): Meets requirements with minor issues
- **Acceptable** (5-6): Basic functionality works, has limitations
- **Poor** (3-4): Significant issues, not production-ready
- **Failing** (1-2): Major failures, unusable

### Aggregate Scoring
- **Total Possible**: 80 points (8 scenarios × 10 points each)
- **Passing Grade**: 56 points (70%)
- **Production Ready**: 64 points (80%)
- **Excellence**: 72 points (90%)

### Test Report Template

```markdown
# MCP Server Test Report

**Date**: [Test Date]
**Server Version**: url-db-mcp-server v1.0.0
**Test Environment**: [Environment Details]

## Scenario Results

### Scenario 1: Protocol Handshake
- **Score**: X/10
- **Status**: PASS/FAIL
- **Notes**: [Detailed observations]

[Repeat for all scenarios]

## Summary
- **Total Score**: X/80
- **Overall Grade**: [Excellent/Good/Acceptable/Poor/Failing]
- **Production Readiness**: [Assessment]

## Recommendations
[Specific recommendations for improvements]
```

## Test Execution Instructions

1. Ensure URL-DB MCP server is built and ready
2. Run test client: `python3 test_mcp_client.py`
3. Execute each scenario systematically
4. Record detailed observations and scores
5. Calculate aggregate scores and overall assessment
6. Generate recommendations based on results

## Notes for LLM Judge

When evaluating test results:
1. Consider MCP specification compliance as highest priority
2. Weight protocol correctness over performance optimizations
3. Evaluate error messages for clarity and actionability
4. Assess overall user experience from AI assistant perspective
5. Consider production deployment readiness
6. Factor in maintainability and extensibility of implementation