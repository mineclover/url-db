# MCP Server Configuration Guide

Complete guide for configuring URL-DB MCP server with Claude Desktop, including production and development variants for optimal user experience.

## 🚀 Quick Start

### Prerequisites
- **Claude Desktop**: Download from [claude.ai](https://claude.ai/download)
- **Claude Code CLI**: For command-line MCP management
- **URL-DB Server**: Built and ready (`make build`)
- **Go Runtime**: Required for running the server

### Build the Server
```bash
cd /path/to/url-db
make build
```

### Method 1: Claude Code CLI (Recommended)

**For local project scope:**
```bash
# Using -- separator to properly pass arguments to the binary
claude mcp add url-db "/absolute/path/to/url-db/bin/url-db" -- -mcp-mode=stdio -db-path=/path/to/your/database.db
```

**For user scope (across all projects):**
```bash
claude mcp add url-db "/absolute/path/to/url-db/bin/url-db" -s user -- -mcp-mode=stdio -db-path=/path/to/your/database.db
```

**Verify the connection:**
```bash
claude mcp list
# Should show: url-db: ✓ Connected

claude mcp get url-db
# Shows detailed configuration
```

**Example with full paths:**
```bash
claude mcp add url-db "/Users/yourname/mcp/url-db/bin/url-db" -- -mcp-mode=stdio -db-path=/Users/yourname/Documents/url-database.db
```

**⚠️ Important**: Always use the `--` separator before your binary's arguments to prevent "unknown option" errors.

### Method 2: Manual Configuration Files

## 📋 Configuration Variants

### 🔇 Production Configuration (Recommended)

**Perfect for**: Daily use, clean Claude Desktop experience, production environments

**Location**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/absolute/path/to/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/Users/yourusername/Documents/url-database.db"
      ],
      "env": {}
    }
  }
}
```

**Features**:
- ✅ Clean Claude Desktop interface
- ✅ Minimal console output
- ✅ Optimal for end users
- ✅ Faster startup

### 🔊 Development Configuration

**Perfect for**: Debugging, development, troubleshooting, testing with separate database

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/absolute/path/to/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/Users/yourusername/Documents/url-database.db"
      ],
      "env": {}
    }
  }
}
```

**Features**:
- 🔧 Same functionality as production
- 📝 Manual debugging via console output
- 🧪 Test database for safe experimentation

## 🎛️ Configuration Options

### Basic Parameters

| Argument | Description | Default | Example |
|----------|-------------|---------|---------|
| `-mcp-mode` | MCP server mode | `stdio` | `-mcp-mode=stdio` |
| `-db-path` | Database file path | `./url-db.sqlite` | `-db-path=/path/to/db.sqlite` |
| `-tool-name` | Composite key prefix | `url-db` | `-tool-name=my-urls` |
| `-port` | HTTP server port | `8080` | `-port=9000` |

### MCP Server Modes

| Mode | Description | Use Case | Endpoint |
|------|-------------|----------|----------|
| `stdio` | Standard input/output | AI assistants (Claude Desktop, Cursor) | stdin/stdout |
| `http` | HTTP JSON-RPC | Web applications, REST clients | `http://localhost:port/mcp` |
| `sse` | Server-Sent Events | Real-time applications | `http://localhost:port/mcp` |

### Environment Variables

| Variable | Purpose | Values | Default |
|----------|---------|--------|---------|
| `AUTO_CREATE_ATTRIBUTES` | Auto-create missing attributes | `true`, `false` | `true` |

**Note**: Logging is currently handled through standard Go logging without environment variable control.

## 📊 Configuration Templates

### 🏢 Enterprise Setup

```json
{
  "mcpServers": {
    "corporate-urls": {
      "command": "/opt/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/var/lib/url-db/corporate.db",
        "-tool-name=corp-links"
      ],
      "env": {
        "AUTO_CREATE_ATTRIBUTES": "false"
      }
    }
  }
}
```

### 🎓 Development/Learning Setup

```json
{
  "mcpServers": {
    "url-db-dev": {
      "command": "/Users/dev/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/tmp/url-db-dev.sqlite"
      ],
      "env": {
        "AUTO_CREATE_ATTRIBUTES": "true"
      }
    }
  }
}
```

### 🏠 Personal Use Setup

```json
{
  "mcpServers": {
    "my-bookmarks": {
      "command": "/Users/yourname/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/Users/yourname/Documents/my-bookmarks.db",
        "-tool-name=bookmarks"
      ],
      "env": {}
    }
  }
}
```

### 🔄 Multiple Database Setup

```json
{
  "mcpServers": {
    "work-urls": {
      "command": "/path/to/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio", 
        "-db-path=/Users/you/work-urls.db",
        "-tool-name=work"
      ],
      "env": {}
    },
    "personal-urls": {
      "command": "/path/to/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/Users/you/personal-urls.db", 
        "-tool-name=personal"
      ],
      "env": {}
    }
  }
}
```

### 🌐 HTTP Mode Setup

```json
{
  "mcpServers": {
    "url-db-http": {
      "command": "/path/to/url-db/bin/url-db",
      "args": [
        "-mcp-mode=http",
        "-port=8080",
        "-db-path=/path/to/database.db"
      ],
      "env": {}
    }
  }
}
```

**HTTP Mode Features**:
- ✅ RESTful API endpoints
- ✅ JSON-RPC 2.0 protocol
- ✅ CORS support
- ✅ Health check endpoint
- ✅ Easy integration with web applications

## 🔧 Configuration Best Practices

### ✅ Production Recommendations

1. **Use Absolute Paths**
   ```json
   "command": "/full/path/to/url-db/bin/url-db"
   ```

2. **Disable Logging for Clean Experience**
   ```json
   "env": { "LOG_LEVEL": "OFF" }
   ```

3. **Secure Database Location**
   ```json
   "-db-path=/Users/yourname/Documents/secure-folder/urls.db"
   ```

4. **Meaningful Tool Names**
   ```json
   "-tool-name=my-research-links"
   ```

### ⚠️ Development Guidelines

1. **Use Test Database**
   ```json
   "-db-path=/tmp/url-db-test.sqlite"
   ```

3. **Allow Auto-Creation**
   ```json
   "env": { "AUTO_CREATE_ATTRIBUTES": "true" }
   ```

## 🔍 Troubleshooting by Configuration

### No Logging Configuration Issues

**Problem**: Server not responding, no error output
```bash
# Enable temporary logging to diagnose
./bin/url-db -mcp-mode=stdio -db-path=test.db
# Check for error messages
```

**Solution**: Switch to development configuration temporarily:
```json
"env": { "LOG_LEVEL": "ERROR" }
```

### Logging Configuration Issues

**Problem**: Want to see server startup messages
```bash
# Run server manually to see console output
./bin/url-db -mcp-mode=stdio -db-path=test.db
```

## 🧪 Testing Your Configuration

### 1. Manual Server Test
```bash
# Test without Claude Desktop
./bin/url-db -mcp-mode=stdio -db-path=test.db
```

### 2. Claude Desktop Integration Test
Ask Claude: 
```
"What MCP servers are available?"
"Can you list domains in my URL database?"
```

### 3. Console Output Verification
**Normal Operation**: Minimal startup messages, clean responses
**Manual Testing**: Console shows server activity when run directly

### 4. HTTP Mode Testing
```bash
# Start HTTP server
./bin/url-db -mcp-mode=http -port=8080 -db-path=test.db

# Test health endpoint
curl http://localhost:8080/health

# Test MCP endpoint
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}}'
```

## 📈 Performance Optimization

### For Speed
```json
{
  "env": {}
}
```

### For Development
```json
{
  "env": {
    "AUTO_CREATE_ATTRIBUTES": "true"
  }
}
```

## 🔐 Security Considerations

1. **Database Permissions**
   ```bash
   chmod 600 /path/to/your/database.db
   ```

2. **Directory Access**
   ```bash
   mkdir -p ~/Documents/url-db
   chmod 755 ~/Documents/url-db
   ```

3. **Path Validation**
   - Always use absolute paths
   - Avoid paths with spaces or special characters
   - Test paths before configuring

## 🆘 Quick Fix Commands

### Reset to Minimal Configuration
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/absolute/path/to/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {}
    }
  }
}
```

### Emergency Debug Mode
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/absolute/path/to/url-db/bin/url-db", 
      "args": ["-mcp-mode=stdio", "-db-path=/tmp/debug.db"],
      "env": {}
    }
  }
}
```

**For debugging**: Run the server manually in terminal to see console output:
```bash
./bin/url-db -mcp-mode=stdio -db-path=/tmp/debug.db
```

## 📚 Related Documentation

- [MCP Specification](https://spec.modelcontextprotocol.io/)
- [Claude Desktop Setup](https://claude.ai/download)
- [Cursor MCP Integration](https://cursor.sh/docs/mcp)
- [URL-DB API Documentation](./api/)

## 🎯 Success Metrics

### ✅ Working Configuration Checklist

- [ ] Server starts without errors
- [ ] Claude Desktop recognizes the MCP server
- [ ] Tools are available in Claude's interface
- [ ] Database operations work correctly
- [ ] No console errors during normal operation
- [ ] Health check endpoint responds (HTTP mode)
- [ ] MCP endpoint responds to JSON-RPC requests

### 🚀 Performance Benchmarks

- **Startup Time**: < 2 seconds
- **Tool Response Time**: < 500ms
- **Memory Usage**: < 50MB
- **Database Operations**: < 100ms for simple queries

## 🔄 Migration Guide

### From Legacy Configuration

If you have an older configuration:

1. **Backup your current config**
2. **Update to new format**
3. **Test with new server**
4. **Remove old configuration**

### Version Compatibility

| URL-DB Version | MCP Protocol | Claude Desktop | Cursor |
|----------------|---------------|----------------|--------|
| 1.0.0+ | 2024-11-05 | ✅ | ✅ |
| 0.9.x | 2024-11-05 | ✅ | ✅ |
| 0.8.x | 2024-11-05 | ⚠️ | ⚠️ |

## 📞 Support

For issues with MCP configuration:

1. **Check the logs**: Run server manually to see error messages
2. **Verify paths**: Ensure all paths are absolute and accessible
3. **Test connectivity**: Use curl to test HTTP endpoints
4. **Check permissions**: Ensure binary and database are readable
5. **Restart applications**: Restart Claude Desktop/Cursor after config changes

---

**Last Updated**: 2024-07-24
**Version**: 1.0.0
**Status**: ✅ Production Ready