# SSE Mode Docker Deployment Guide

## 🚀 Quick Deployment Commands

### 1. Single Docker Command (Simplest)

```bash
# Run SSE server on port 8080
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse

# Test the connection
curl http://localhost:8080/health
```

### 2. Docker Compose (Recommended)

```bash
# Download docker-compose-sse.yml (or use existing one)
curl -o docker-compose-sse.yml https://raw.githubusercontent.com/your-repo/url-db/main/docker-compose-sse.yml

# Start the service
docker-compose -f docker-compose-sse.yml up -d

# View logs
docker-compose -f docker-compose-sse.yml logs -f

# Stop the service
docker-compose -f docker-compose-sse.yml down
```

### 3. Different Port

```bash
# Run on port 9090 instead of 8080
docker run -d \
  --name url-db-sse \
  -p 9090:8080 \
  -v $(pwd)/data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

### 4. Named Volume (Persistent Data)

```bash
# Create a named volume for persistent data
docker volume create url-db-data

# Run with named volume
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v url-db-data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

## 📋 Management Commands

### Start/Stop/Restart

```bash
# Stop the container
docker stop url-db-sse

# Start the container
docker start url-db-sse

# Restart the container
docker restart url-db-sse

# Remove the container
docker rm url-db-sse
```

### View Logs

```bash
# View recent logs
docker logs url-db-sse

# Follow logs in real-time
docker logs -f url-db-sse

# View last 100 lines
docker logs --tail 100 url-db-sse
```

### Health Check

```bash
# Basic health check
curl http://localhost:8080/health

# Expected response:
# {"mode":"sse","server":"url-db-mcp-server","status":"ok"}

# Container health status
docker ps --filter name=url-db-sse
```

## 🔧 Configuration Options

### Environment Variables

```bash
# Run with debug logging
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  -e LOG_LEVEL=debug \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

### Custom Database Path

```bash
# Use custom database file
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v $(pwd)/my-db:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse -db-path=/data/my-urls.sqlite
```

### Custom Tool Name

```bash
# Use custom tool name for composite keys
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse -tool-name=my-url-db
```

## 📝 Docker Compose Examples

### Basic docker-compose-sse.yml

```yaml
version: '3.8'

services:
  url-db-sse:
    image: asfdassdssa/url-db:latest
    container_name: url-db-sse
    ports:
      - "8080:8080"
    volumes:
      - url-db-data:/data
    command: ["-mcp-mode=sse", "-port=8080"]
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  url-db-data:
```

### With Custom Port

```yaml
version: '3.8'

services:
  url-db-sse:
    image: asfdassdssa/url-db:latest
    container_name: url-db-sse
    ports:
      - "9090:8080"  # External port 9090
    volumes:
      - ./data:/data  # Host directory
    command: ["-mcp-mode=sse", "-port=8080"]
    restart: unless-stopped
```

## 🧪 Testing the Deployment

### 1. Health Check

```bash
curl http://localhost:8080/health
```

### 2. Initialize MCP

```bash
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

### 3. List Tools

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "id": 2
  }'
```

### 4. Create Test Domain

```bash
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

## 🐛 Troubleshooting

### Container Won't Start

```bash
# Check if port is already in use
lsof -i :8080

# Check container logs
docker logs url-db-sse

# Verify image is available
docker images | grep url-db
```

### Permission Issues

```bash
# Fix data directory permissions
sudo chown -R 1000:1000 ./data

# Or use different user
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --user $(id -u):$(id -g) \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

### Database Issues

```bash
# Check database file
docker exec url-db-sse ls -la /data/

# Access container shell
docker exec -it url-db-sse sh

# Check database contents
docker exec url-db-sse sqlite3 /data/url-db.sqlite ".tables"
```

## 🔄 Updates

### Update to Latest Version

```bash
# Pull latest image
docker pull asfdassdssa/url-db:latest

# Stop and remove old container
docker stop url-db-sse
docker rm url-db-sse

# Start new container (data is preserved in volume)
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v url-db-data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

### With Docker Compose

```bash
# Pull latest and restart
docker-compose -f docker-compose-sse.yml pull
docker-compose -f docker-compose-sse.yml up -d
```

## 🌐 Production Deployment

For production use, consider:

1. **Reverse Proxy**: Use nginx or traefik for SSL/TLS
2. **Monitoring**: Add health check endpoints to monitoring
3. **Backup**: Regular database backups
4. **Security**: Network isolation, access controls
5. **Performance**: Resource limits and scaling

### Example with Resource Limits

```yaml
version: '3.8'

services:
  url-db-sse:
    image: asfdassdssa/url-db:latest
    container_name: url-db-sse
    ports:
      - "8080:8080"
    volumes:
      - url-db-data:/data
    command: ["-mcp-mode=sse", "-port=8080"]
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  url-db-data:
```