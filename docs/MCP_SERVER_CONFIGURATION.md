# MCP Server Configuration Guide

Complete guide for configuring URL-DB MCP server with Claude Desktop, including logging and no-logging variants for optimal user experience.

## 🚀 Quick Start

### Prerequisites
- **Claude Desktop**: Download from [claude.ai](https://claude.ai/download)
- **URL-DB Server**: Built and ready (`make build`)
- **Go Runtime**: Required for running the server

### Build the Server
```bash
cd /path/to/url-db
make build
```

## 📋 Configuration Variants

### 🔇 Production Configuration (No Logging - Recommended)

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
      "env": {
        "LOG_LEVEL": "OFF"
      }
    }
  }
}
```

**Features**:
- ✅ Clean Claude Desktop interface
- ✅ No console output interference  
- ✅ Optimal for end users
- ✅ Faster startup

### 🔊 Development Configuration (With Logging)

**Perfect for**: Debugging, development, troubleshooting, learning MCP internals

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/absolute/path/to/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/Users/yourusername/Documents/url-database.db"
      ],
      "env": {
        "LOG_LEVEL": "DEBUG",
        "LOG_FORMAT": "json"
      }
    }
  }
}
```

**Features**:
- 🔍 Detailed request/response logs
- 🐛 Error debugging information
- 📊 Performance metrics
- 🔧 Development insights

## 🎛️ Configuration Options

### Basic Parameters

| Argument | Description | Default | Example |
|----------|-------------|---------|---------|
| `-mcp-mode` | MCP server mode | `stdio` | `-mcp-mode=stdio` |
| `-db-path` | Database file path | `./url-db.sqlite` | `-db-path=/path/to/db.sqlite` |
| `-tool-name` | Composite key prefix | `url-db` | `-tool-name=my-urls` |
| `-port` | HTTP server port | `8080` | `-port=9000` |

### Environment Variables

| Variable | Purpose | Values | Default |
|----------|---------|--------|---------|
| `LOG_LEVEL` | Logging verbosity | `OFF`, `ERROR`, `WARN`, `INFO`, `DEBUG` | `INFO` |
| `LOG_FORMAT` | Log output format | `text`, `json` | `text` |
| `AUTO_CREATE_ATTRIBUTES` | Auto-create missing attributes | `true`, `false` | `true` |

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
        "LOG_LEVEL": "WARN",
        "LOG_FORMAT": "json",
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
        "LOG_LEVEL": "DEBUG",
        "LOG_FORMAT": "text",
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
      "env": {
        "LOG_LEVEL": "ERROR"
      }
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
      "env": { "LOG_LEVEL": "OFF" }
    },
    "personal-urls": {
      "command": "/path/to/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/Users/you/personal-urls.db", 
        "-tool-name=personal"
      ],
      "env": { "LOG_LEVEL": "OFF" }
    }
  }
}
```

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

1. **Enable Debug Logging**
   ```json
   "env": { "LOG_LEVEL": "DEBUG", "LOG_FORMAT": "json" }
   ```

2. **Use Test Database**
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

**Problem**: Too much console output in Claude Desktop
```json
// Change from DEBUG to ERROR
"env": { "LOG_LEVEL": "ERROR" }
```

**Problem**: Can't see MCP protocol details
```json
// Enable detailed logging
"env": { "LOG_LEVEL": "DEBUG", "LOG_FORMAT": "json" }
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

### 3. Logging Verification
**No Logging**: Clean responses with no extra output
**With Logging**: Console shows JSON-RPC requests and responses

## 📈 Performance Optimization

### For Speed (No Logging)
```json
{
  "env": {
    "LOG_LEVEL": "OFF"
  }
}
```

### For Monitoring (Minimal Logging)
```json
{
  "env": {
    "LOG_LEVEL": "WARN",
    "LOG_FORMAT": "json"
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
      "env": { "LOG_LEVEL": "DEBUG" }
    }
  }
}
```

## 📚 Related Documentation

- [MCP Claude Setup Guide](MCP_CLAUDE_SETUP.md) - Basic setup instructions
- [MCP Testing Guide](MCP_TESTING_GUIDE.md) - Testing procedures
- [Tool Specification](../specs/mcp-tools.yaml) - Available MCP tools
- [CLAUDE.md](../CLAUDE.md) - Developer integration guide

---

**Remember**: Always restart Claude Desktop after changing the configuration file!