# SSE Mode Quick Start Guide

## 🚀 30-Second Setup

### Using Docker (Recommended)

```bash
# 1. Start SSE server
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse

# 2. Test connection
curl http://localhost:8080/health
```

### Using Docker Compose

```bash
# 1. Get docker-compose file
curl -o docker-compose-sse.yml https://raw.githubusercontent.com/your-repo/url-db/main/docker-compose-sse.yml

# 2. Start service
docker-compose -f docker-compose-sse.yml up -d

# 3. Test connection
curl http://localhost:8080/health
```

## 📡 Quick API Test

### Test with cURL

```bash
# Initialize MCP connection
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {"protocolVersion": "2025-06-18"},
    "id": 1
  }'

# Create a domain
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "create_domain",
      "arguments": {"name": "test", "description": "Test domain"}
    },
    "id": 2
  }'

# Add a URL
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "create_node",
      "arguments": {
        "domain_name": "test",
        "url": "https://example.com",
        "title": "Example Site"
      }
    },
    "id": 3
  }'
```

## 💻 Client Examples

### JavaScript

```javascript
async function testSSE() {
  const endpoint = 'http://localhost:8080/mcp';
  
  // Helper function to call MCP
  async function callMCP(method, params = {}) {
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        jsonrpc: '2.0',
        method,
        params,
        id: Date.now()
      })
    });
    
    const text = await response.text();
    const data = text.split('\n').find(l => l.startsWith('data: '));
    return JSON.parse(data.substring(6));
  }
  
  // Initialize
  const init = await callMCP('initialize', {
    protocolVersion: '2025-06-18'
  });
  console.log('Initialized:', init);
  
  // List tools
  const tools = await callMCP('tools/list');
  console.log('Available tools:', tools.result.tools.length);
  
  // Create domain
  const domain = await callMCP('tools/call', {
    name: 'create_domain',
    arguments: { name: 'bookmarks', description: 'My bookmarks' }
  });
  console.log('Domain created:', domain);
}

testSSE();
```

### Python

```python
import requests
import json

def call_mcp(method, params=None):
    response = requests.post('http://localhost:8080/mcp', json={
        'jsonrpc': '2.0',
        'method': method,
        'params': params or {},
        'id': 1
    })
    
    # Parse SSE format
    for line in response.text.split('\n'):
        if line.startswith('data: '):
            return json.loads(line[6:])

# Initialize
init_result = call_mcp('initialize', {'protocolVersion': '2025-06-18'})
print('Initialized:', init_result)

# Create domain
domain = call_mcp('tools/call', {
    'name': 'create_domain',
    'arguments': {'name': 'python-test', 'description': 'Created from Python'}
})
print('Domain:', domain)

# List domains
domains = call_mcp('tools/call', {'name': 'list_domains', 'arguments': {}})
print('Domains:', domains)
```

## 🎯 Common Operations

### Domain Management

```bash
# List all domains
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_domains","arguments":{}},"id":1}'

# Create domain
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_domain","arguments":{"name":"work","description":"Work URLs"}},"id":2}'
```

### URL Management

```bash
# Add URL to domain
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_node","arguments":{"domain_name":"work","url":"https://github.com","title":"GitHub"}},"id":3}'

# List URLs in domain
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_nodes","arguments":{"domain_name":"work"}},"id":4}'

# Scan all content (new feature)
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"work","max_tokens_per_page":3000}},"id":5}'
```

## 🔧 Management Commands

### Start/Stop Server

```bash
# Stop server
docker stop url-db-sse

# Start server
docker start url-db-sse

# View logs
docker logs -f url-db-sse

# Remove server (data persists in volume)
docker rm url-db-sse
```

### Different Ports

```bash
# Run on port 9090
docker run -d \
  --name url-db-sse-9090 \
  -p 9090:8080 \
  -v $(pwd)/data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse

# Test on different port
curl http://localhost:9090/health
```

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Check what's using port 8080
lsof -i :8080

# Use different port
docker run -d \
  --name url-db-sse \
  -p 8081:8080 \
  -v $(pwd)/data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

### Container Won't Start

```bash
# Check container status
docker ps -a | grep url-db-sse

# View container logs
docker logs url-db-sse

# Remove and recreate
docker rm url-db-sse
docker run -d --name url-db-sse -p 8080:8080 -v $(pwd)/data:/data asfdassdssa/url-db:latest -mcp-mode=sse
```

### Data Persistence

```bash
# Check if data directory exists
ls -la data/

# Check database file
ls -la data/url-db.sqlite

# Access container to check database
docker exec -it url-db-sse ls -la /data/
```

## 📚 Next Steps

1. Read the full [SSE Deployment Guide](./SSE_DEPLOYMENT.md)
2. Check [MCP Testing Guide](./MCP_TESTING_GUIDE.md)
3. See [API Documentation](../README.md)

## 💡 Tips

- Use named volumes for production: `-v url-db-data:/data`
- Set resource limits for production deployments
- SSE mode is for HTTP clients, use stdio mode for AI assistants
- Response format is SSE: `data: {json}\n\n`