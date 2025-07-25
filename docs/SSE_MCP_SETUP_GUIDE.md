# SSE Mode Setup Guide for URL-DB MCP Server

This guide explains how to set up and configure the URL-DB MCP server in SSE (Server-Sent Events) mode for HTTP-based MCP communication.

## Table of Contents
- [Overview](#overview)
- [Local Setup](#local-setup)
- [Docker Setup](#docker-setup)
- [MCP Client Configuration](#mcp-client-configuration)
- [Testing the Connection](#testing-the-connection)
- [Troubleshooting](#troubleshooting)

## Overview

SSE mode allows the MCP server to communicate over HTTP using Server-Sent Events format. This is useful for:
- HTTP-based MCP clients
- Network-accessible MCP servers
- Cross-platform compatibility

**Note**: The current SSE implementation uses request-response pattern with SSE formatting, not persistent streaming connections.

## Local Setup

### 1. Build the Server

```bash
# Clone the repository
git clone https://github.com/your-username/url-db.git
cd url-db

# Build the binary
make build
# or
go build -o bin/url-db ./cmd/server
```

### 2. Start in SSE Mode

```bash
# Basic SSE mode (default port 8080)
./bin/url-db -mcp-mode=sse

# Custom port
./bin/url-db -mcp-mode=sse -port=9090

# With custom database path
./bin/url-db -mcp-mode=sse -port=8080 -db-path=/path/to/database.sqlite

# With custom tool name
./bin/url-db -mcp-mode=sse -tool-name=my-url-db
```

### 3. Verify Server is Running

```bash
# Check health endpoint
curl http://localhost:8080/health

# Expected response:
# {"mode":"sse","server":"url-db-mcp-server","status":"ok"}
```

## Docker Setup

### 1. Using Pre-built Docker Image

```bash
# Pull the image
docker pull asfdassdssa/url-db:latest

# Run in SSE mode
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v url-db-data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

### 2. Docker Compose Setup

Create `docker-compose-sse.yml`:

```yaml
version: '3.8'

services:
  url-db-sse:
    image: asfdassdssa/url-db:latest
    container_name: url-db-sse-server
    ports:
      - "${SSE_PORT:-8080}:8080"
    volumes:
      - url-db-data:/data
    environment:
      - LOG_LEVEL=${LOG_LEVEL:-info}
      - AUTO_CREATE_ATTRIBUTES=${AUTO_CREATE_ATTRIBUTES:-true}
    command: ["-mcp-mode=sse", "-port=8080"]
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    networks:
      - mcp-network

volumes:
  url-db-data:
    driver: local

networks:
  mcp-network:
    driver: bridge
```

### 3. Start with Docker Compose

```bash
# Start the service
docker-compose -f docker-compose-sse.yml up -d

# View logs
docker-compose -f docker-compose-sse.yml logs -f

# Stop service
docker-compose -f docker-compose-sse.yml down
```

## MCP Client Configuration

### 1. Generic HTTP Client Configuration

```json
{
  "servers": {
    "url-db-sse": {
      "type": "http",
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Content-Type": "application/json"
      }
    }
  }
}
```

### 2. JavaScript/TypeScript Client Example

```javascript
class MCPSSEClient {
  constructor(endpoint = 'http://localhost:8080/mcp') {
    this.endpoint = endpoint;
    this.requestId = 0;
  }

  async callTool(toolName, args = {}) {
    const request = {
      jsonrpc: '2.0',
      method: 'tools/call',
      params: {
        name: toolName,
        arguments: args
      },
      id: ++this.requestId
    };

    const response = await fetch(this.endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request)
    });

    const text = await response.text();
    // Parse SSE format
    const data = text.split('\n')
      .find(line => line.startsWith('data: '))
      ?.substring(6);
    
    return JSON.parse(data || '{}');
  }

  async initialize() {
    const request = {
      jsonrpc: '2.0',
      method: 'initialize',
      params: {
        protocolVersion: '2025-06-18',
        capabilities: {
          roots: { listChanged: true }
        },
        clientInfo: {
          name: 'mcp-sse-client',
          version: '1.0.0'
        }
      },
      id: ++this.requestId
    };

    return this.sendRequest(request);
  }

  async sendRequest(request) {
    const response = await fetch(this.endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request)
    });

    const text = await response.text();
    const data = text.split('\n')
      .find(line => line.startsWith('data: '))
      ?.substring(6);
    
    return JSON.parse(data || '{}');
  }
}

// Usage example
async function main() {
  const client = new MCPSSEClient();
  
  // Initialize
  const initResult = await client.initialize();
  console.log('Initialized:', initResult);
  
  // Create domain
  const domain = await client.callTool('create_domain', {
    name: 'test-domain',
    description: 'Test domain for SSE'
  });
  console.log('Domain created:', domain);
  
  // List domains
  const domains = await client.callTool('list_domains', {});
  console.log('Domains:', domains);
}
```

### 3. Python Client Example

```python
import requests
import json

class MCPSSEClient:
    def __init__(self, endpoint='http://localhost:8080/mcp'):
        self.endpoint = endpoint
        self.request_id = 0
    
    def call_tool(self, tool_name, args=None):
        self.request_id += 1
        request = {
            'jsonrpc': '2.0',
            'method': 'tools/call',
            'params': {
                'name': tool_name,
                'arguments': args or {}
            },
            'id': self.request_id
        }
        
        response = requests.post(self.endpoint, json=request)
        
        # Parse SSE format
        lines = response.text.split('\n')
        for line in lines:
            if line.startswith('data: '):
                return json.loads(line[6:])
        
        return None
    
    def initialize(self):
        self.request_id += 1
        request = {
            'jsonrpc': '2.0',
            'method': 'initialize',
            'params': {
                'protocolVersion': '2025-06-18',
                'capabilities': {
                    'roots': {'listChanged': True}
                },
                'clientInfo': {
                    'name': 'python-mcp-client',
                    'version': '1.0.0'
                }
            },
            'id': self.request_id
        }
        
        response = requests.post(self.endpoint, json=request)
        lines = response.text.split('\n')
        for line in lines:
            if line.startswith('data: '):
                return json.loads(line[6:])
        return None

# Usage
if __name__ == '__main__':
    client = MCPSSEClient()
    
    # Initialize
    init_result = client.initialize()
    print('Initialized:', init_result)
    
    # Create domain
    domain = client.call_tool('create_domain', {
        'name': 'python-test',
        'description': 'Created from Python'
    })
    print('Domain:', domain)
```

## Testing the Connection

### 1. Basic Connection Test

```bash
# Test with curl
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {
      "protocolVersion": "2025-06-18"
    },
    "id": 1
  }'
```

### 2. Tool Call Test

```bash
# List available tools
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "id": 2
  }'

# Create a domain
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "create_domain",
      "arguments": {
        "name": "test",
        "description": "Test domain"
      }
    },
    "id": 3
  }'
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   ```bash
   # Check if server is running
   ps aux | grep url-db
   
   # Check port availability
   lsof -i :8080
   ```

2. **Docker Container Not Starting**
   ```bash
   # Check container logs
   docker logs url-db-sse
   
   # Verify volume permissions
   docker exec url-db-sse ls -la /data
   ```

3. **SSE Format Issues**
   - Ensure client properly parses SSE format (data: prefix)
   - Check Content-Type header is text/event-stream

4. **Database Access Issues**
   ```bash
   # For Docker, ensure volume is properly mounted
   docker volume inspect url-db-data
   
   # Check database file permissions
   docker exec url-db-sse ls -la /data/url-db.sqlite
   ```

### Debug Mode

```bash
# Run with verbose logging
./bin/url-db -mcp-mode=sse -debug

# Docker debug mode
docker run -it --rm \
  -p 8080:8080 \
  -e LOG_LEVEL=debug \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

## Performance Tuning

### 1. Connection Pooling

For high-traffic scenarios, configure connection pooling:

```bash
# Set database connection limits
./bin/url-db -mcp-mode=sse \
  -max-open-conns=100 \
  -max-idle-conns=50
```

### 2. Resource Limits

```yaml
# Docker resource limits
services:
  url-db-sse:
    # ... other config
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 256M
```

## Monitoring

### Health Checks

```bash
# Simple health check script
#!/bin/bash
while true; do
  if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "Health check failed!"
    # Send alert
  fi
  sleep 30
done
```

## Conclusion

The SSE mode provides HTTP-based access to the MCP server, suitable for HTTP clients and network deployments. While the current implementation uses request-response pattern, it provides a foundation for HTTP-based MCP communication.