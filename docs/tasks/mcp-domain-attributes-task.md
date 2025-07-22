# Task: Implement Domain Attribute Management in MCP

## Problem Statement

The MCP layer currently lacks domain attribute definition management capabilities, despite these features being fully implemented in the database and API layers. Users cannot create or manage attribute definitions through MCP, creating a dependency on the REST API for basic setup tasks.

## Current State

### What Works
- Database schema supports full attribute management
- REST API provides complete CRUD operations for attributes
- MCP can get/set node attribute values (but only for pre-existing attributes)

### What's Missing
- No way to list available attributes for a domain via MCP
- No way to create new attribute definitions via MCP
- No way to view attribute metadata (type, description) via MCP
- No way to update or delete attribute definitions via MCP

## Proposed Solution

Add 5 new MCP tools for domain attribute management:

### 1. `list_mcp_domain_attributes`
- **Purpose**: List all attribute definitions for a domain
- **Parameters**: 
  - `domain_name` (required): The domain to list attributes for
- **Returns**: Array of attribute definitions with id, name, type, description

### 2. `create_mcp_domain_attribute`
- **Purpose**: Create a new attribute definition for a domain
- **Parameters**:
  - `domain_name` (required): The domain to add attribute to
  - `name` (required): Attribute name
  - `type` (required): One of: tag, ordered_tag, number, string, markdown, image
  - `description` (optional): Human-readable description
- **Returns**: Created attribute with composite ID

### 3. `get_mcp_domain_attribute`
- **Purpose**: Get a specific attribute definition
- **Parameters**:
  - `composite_id` (required): Format `tool-name:domain:attribute-id`
- **Returns**: Attribute definition details

### 4. `update_mcp_domain_attribute`
- **Purpose**: Update attribute description
- **Parameters**:
  - `composite_id` (required): Format `tool-name:domain:attribute-id`
  - `description` (required): New description
- **Returns**: Updated attribute

### 5. `delete_mcp_domain_attribute`
- **Purpose**: Delete an attribute definition (if no values exist)
- **Parameters**:
  - `composite_id` (required): Format `tool-name:domain:attribute-id`
- **Returns**: Success confirmation

## Implementation Plan

### Phase 1: Core Implementation
1. Add attribute composite key support in `/internal/mcp/composite_key.go`
2. Implement the 5 new tools in `/internal/mcp/tools.go`
3. Add attribute validators in `/internal/mcp/validators.go`
4. Update tool registration in MCP handler

### Phase 2: Testing
1. Add unit tests for composite key parsing with attributes
2. Add integration tests for each new tool
3. Update MCP protocol compliance tests

### Phase 3: Documentation
1. Update README.md with new tool descriptions
2. Add usage examples for attribute management workflow
3. Update CLAUDE.md with attribute management patterns

## Technical Considerations

### Composite Key Format
- Current: `tool-name:domain:node-id`
- New: `tool-name:domain:attribute-id` for attribute operations
- Need to differentiate between node and attribute IDs in parsing

### Validation Requirements
- Attribute names must be unique within a domain
- Attribute types must be valid (enum validation)
- Cannot delete attributes that have values assigned

### Error Handling
- -32602: Invalid parameters (bad type, missing required fields)
- -32603: Internal errors (database failures)
- -32001: Attribute already exists
- -32002: Attribute not found
- -32003: Cannot delete attribute with existing values

## User Workflow Example

```bash
# 1. Create a domain
"Create a new domain called 'projects'"

# 2. Define attributes for the domain
"Create a 'status' attribute of type 'tag' for projects domain"
"Create a 'priority' attribute of type 'ordered_tag' for projects domain"
"Create a 'description' attribute of type 'markdown' for projects domain"

# 3. List available attributes
"Show me all attributes for the projects domain"

# 4. Add URLs with attributes
"Add https://github.com/myproject to projects"
"Set status to 'active' and priority to 'high' for this URL"
```

## Success Criteria

1. All 5 tools implemented and working
2. Full test coverage for new functionality
3. MCP protocol compliance maintained at 90%+
4. Seamless integration with existing node attribute operations
5. Clear documentation and examples

## Estimated Effort

- Implementation: 4-6 hours
- Testing: 2-3 hours
- Documentation: 1-2 hours
- Total: 7-11 hours