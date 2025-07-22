# URL-DB MCP Server

A powerful URL database management system with MCP (Model Context Protocol) integration, supporting unlimited attribute tagging and domain-based organization.

## 🚀 Quick Start

### Install URL-DB as MCP Server in Claude Desktop

```bash
# Run the setup script
./setup-mcp.sh

# Copy and run the displayed command, for example:
claude mcp add url-db "/Users/junwoobang/mcp/url-db/bin/url-db" \
  --args="-mcp-mode=stdio" \
  --env="DATABASE_URL=file:/Users/junwoobang/mcp/url-db/url-db.db"

# Restart Claude Desktop
```

## ✅ Features

- **11 MCP Tools**: Complete URL and domain management
- **Unlimited Tagging**: 6 attribute types (tag, ordered_tag, number, string, markdown, image)
- **Domain Organization**: Group URLs by domains/namespaces
- **Composite Keys**: `url-db:domain:id` format for unique identification
- **Resource System**: Access data via `mcp://` URIs
- **Batch Operations**: Efficient bulk processing
- **JSON-RPC 2.0**: Full protocol compliance (92% LLM-as-a-Judge score)

## 📋 Available MCP Tools

1. `list_mcp_domains` - List all domains
2. `create_mcp_domain` - Create new domain
3. `list_mcp_nodes` - List URLs in domain
4. `create_mcp_node` - Add new URL
5. `get_mcp_node` - Get URL by composite ID
6. `update_mcp_node` - Update URL metadata
7. `delete_mcp_node` - Delete URL
8. `find_mcp_node_by_url` - Find URL in domain
9. `get_mcp_node_attributes` - Get URL attributes
10. `set_mcp_node_attributes` - Set URL attributes
11. `get_mcp_server_info` - Get server information

## 🛠️ Manual Installation

### 1. Build from Source

```bash
# Clone repository
git clone https://github.com/yourusername/url-db.git
cd url-db

# Build
./build.sh

# Binary will be at: bin/url-db
```

### 2. Configure Claude Desktop

Option 1 - Using Claude CLI:
```bash
claude mcp add url-db /path/to/url-db/bin/url-db \
  --args="-mcp-mode=stdio" \
  --env="DATABASE_URL=file:/path/to/url-db/url-db.db"
```

Option 2 - Manual config (`~/Library/Application Support/Claude/claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/path/to/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {
        "DATABASE_URL": "file:/path/to/url-db/url-db.db"
      }
    }
  }
}
```

## 📖 Usage Examples

### Basic URL Management
```
"Create a new domain called 'tech-articles'"
"Add https://example.com/article to tech-articles domain"
"List all URLs in tech-articles domain"
"Find the URL https://example.com/article"
```

### Attribute Management
```
"Set category to 'javascript' for this URL"
"Add tags: tutorial, beginner, react"
"Set priority to high"
"Show all attributes for this URL"
```

### Advanced Operations
```
"Get nodes with IDs: url-db:tech:1, url-db:tech:2"
"Update the title of url-db:tech:1 to 'New Title'"
"Delete url-db:tech:5"
```

## 🧪 Test Results

- **MCP Protocol Compliance**: 92% (Excellent)
- **Integration Tests**: 100% success rate
- **Test Coverage**: 15.8% (MCP package)
- **All 11 tools**: Verified and working

See [Test Results](docs/testing/mcp-test-results.md) for details.

## 📚 Documentation

- [MCP Server Setup Guide](docs/mcp-server-setup-guide.md)
- [API Documentation](docs/api/)
- [Database Schema](docs/database-schema.md)
- [Test Scenarios](docs/testing/mcp-llm-judge-scenarios.md)

## 🔧 Configuration

Environment variables:
- `DATABASE_URL`: SQLite database path (default: `file:./url-db.db`)
- `TOOL_NAME`: Tool name prefix for composite keys (default: `url-db`)
- `PORT`: HTTP server port for SSE mode (default: `8080`)
- `LOG_LEVEL`: Logging level (default: `info`)

## 🐛 Troubleshooting

### MCP server not found in Claude
```bash
# Check if binary exists and is executable
ls -la /path/to/url-db/bin/url-db
chmod +x /path/to/url-db/bin/url-db

# Test directly
/path/to/url-db/bin/url-db -mcp-mode=stdio
```

### Database issues
```bash
# Check database file
ls -la /path/to/url-db/url-db.db

# Create if missing
touch /path/to/url-db/url-db.db
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with the [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- Tested with [Claude Desktop](https://claude.ai)
- Powered by SQLite and Go

---
🤖 Generated with [Claude Code](https://claude.ai/code)