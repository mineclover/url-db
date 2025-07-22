# MCP (Model Context Protocol) Setup Guide

This guide explains how to configure and use URL-DB with MCP in both stdio and SSE modes.

## Overview

URL-DB supports MCP (Model Context Protocol) in two modes:

1. **stdio mode** - For AI assistant integration (Claude Desktop, Cursor)
2. **SSE mode** - For HTTP-based integration with Server-Sent Events

## Stdio Mode Setup

### What is stdio mode?

Stdio mode uses standard input/output for JSON-RPC 2.0 communication. This is the primary mode for AI assistant integration.

### Configuration

#### For Claude Desktop

Add the following to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/path/to/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {
        "DATABASE_URL": "file:/path/to/your/database.sqlite"
      }
    }
  }
}
```

#### For Cursor

Add to your Cursor MCP configuration:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/path/to/url-db",
      "args": ["-mcp-mode=stdio", "-db-path=/path/to/database.sqlite"]
    }
  }
}
```

### Available MCP Tools (18 total)

All tools are defined in `/specs/mcp-tools.yaml` with auto-generated constants. When using stdio mode, the following tools are available:

#### Domain Management
- `list_domains` - List all domains
- `create_domain` - Create a new domain

#### Node Operations
- `list_nodes` - List nodes with pagination and search
- `create_node` - Create a new URL node
- `get_node` - Get node by composite ID
- `update_node` - Update node title/description
- `delete_node` - Delete a node
- `find_node_by_url` - Find node by URL

#### Node Attributes
- `get_node_attributes` - Get all attributes for a node
- `set_node_attributes` - Set/update node attributes

#### Domain Schema
- `list_domain_attributes` - List attribute definitions
- `create_domain_attribute` - Create attribute definition
- `get_domain_attribute` - Get attribute definition
- `update_domain_attribute` - Update attribute description
- `delete_domain_attribute` - Delete unused attribute

#### Enhanced Queries
- `get_node_with_attributes` - Get node with attributes in one call
- `filter_nodes_by_attributes` - Filter nodes by attribute values

#### Server Info
- `get_server_info` - Get server capabilities

### Testing stdio Mode

```bash
# Start in stdio mode
./bin/url-db -mcp-mode=stdio -db-path=./test.db

# Send JSON-RPC request (example)
{"jsonrpc":"2.0","id":1,"method":"tools/list"}
```

### Example Usage in Claude

```
User: Create a new domain for my bookmarks

Claude will use: create_domain tool with appropriate parameters

User: Add a URL to the tech domain with tags

Claude will use: create_node and set_node_attributes tools
```

## SSE Mode Setup

### What is SSE mode?

SSE (Server-Sent Events) mode runs URL-DB as an HTTP server with REST API endpoints. This mode is for web-based integrations.

### Starting the Server

```bash
# Start in SSE mode (default)
./bin/url-db -mcp-mode=sse -port=8080

# Or simply (SSE is default)
./bin/url-db
```

### Configuration Options

```bash
# Custom port
./bin/url-db -mcp-mode=sse -port=3000

# Custom database
./bin/url-db -mcp-mode=sse -db-path=/path/to/database.sqlite

# Custom tool name (affects composite keys)
./bin/url-db -mcp-mode=sse -tool-name=my-url-db
```

### REST API Endpoints

SSE mode provides standard REST API endpoints:

#### Health Check
- `GET /health` - Server health status

#### Domains
- `POST /api/domains` - Create domain
- `GET /api/domains` - List domains
- `GET /api/domains/:id` - Get domain
- `PUT /api/domains/:id` - Update domain
- `DELETE /api/domains/:id` - Delete domain

#### Domain Attributes
- `POST /api/domains/:id/attributes` - Create attribute
- `GET /api/domains/:id/attributes` - List attributes

#### Nodes/URLs
- `POST /api/domains/:id/urls` - Create node
- `GET /api/domains/:id/urls` - List nodes
- `GET /api/urls/:id` - Get node
- `PUT /api/urls/:id` - Update node
- `DELETE /api/urls/:id` - Delete node

#### Node Attributes
- `POST /api/urls/:id/attributes` - Set attribute
- `GET /api/urls/:id/attributes` - Get attributes
- `DELETE /api/urls/:id/attributes/:attr_id` - Delete attribute

### Testing SSE Mode

```bash
# Health check
curl http://localhost:8080/health

# Create domain
curl -X POST http://localhost:8080/api/domains \
  -H "Content-Type: application/json" \
  -d '{"name":"tech","description":"Technology links"}'

# Create node
curl -X POST http://localhost:8080/api/domains/1/urls \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","title":"Example"}'
```

## Key Differences Between Modes

### Stdio Mode
- **Protocol**: JSON-RPC 2.0 over stdin/stdout
- **Use Case**: AI assistant integration
- **Tools**: 18 MCP tools available
- **Features**: Enhanced queries, batch operations
- **Composite Keys**: `tool-name:domain:id` format

### SSE Mode
- **Protocol**: REST API over HTTP
- **Use Case**: Web applications, services
- **Endpoints**: Standard REST endpoints
- **Features**: Real-time updates via SSE (planned)
- **IDs**: Standard integer IDs

## Environment Variables

Both modes support:

- `DATABASE_URL` - Database connection string (default: file:./url-db.sqlite)
- `PORT` - Server port (SSE mode only, default: 8080)
- `TOOL_NAME` - Tool name for composite keys (default: url-db)

All default values are defined in `/internal/constants/constants.go` for consistency.

## Composite Key Format

In stdio mode, resources use composite keys:

```
Format: tool-name:domain:id
Example: url-db:tech:123
```

## Domain Schema System

URL-DB enforces domain schemas:

1. Define attributes at domain level
2. Nodes can only have attributes defined in their domain
3. Attribute types: tag, ordered_tag, number, string, markdown, image

## Best Practices

### For AI Assistants (stdio mode)
1. Use descriptive domain names
2. Define clear attribute schemas
3. Leverage enhanced query tools
4. Use composite keys consistently

### For Web Apps (SSE mode)
1. Implement proper error handling
2. Use pagination for large datasets
3. Cache frequently accessed data
4. Monitor server health endpoint

## Troubleshooting

### Stdio Mode Issues

1. **No response**: Check JSON-RPC format
2. **Tool not found**: Verify tool name spelling
3. **Composite key errors**: Check format `tool:domain:id`

### SSE Mode Issues

1. **Port in use**: Change port with `-port` flag
2. **CORS errors**: Server includes CORS headers by default
3. **Database locked**: Ensure single instance running

## Migration Between Modes

The same database works with both modes:

```bash
# Create data in SSE mode
./bin/url-db -mcp-mode=sse

# Access same data in stdio mode
./bin/url-db -mcp-mode=stdio
```

## Security Considerations

1. **stdio mode**: Runs with AI assistant permissions
2. **SSE mode**: Implement authentication if exposed publicly
3. **Database**: Use appropriate file permissions
4. **Composite keys**: Not guessable, include domain context

## Future Enhancements

- Real-time SSE events for data changes
- WebSocket support
- Enhanced filtering in SSE mode
- Bulk import/export tools