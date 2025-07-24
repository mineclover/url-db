# Docker Deployment Guide for URL-DB MCP Server

## Overview

This guide provides instructions for deploying the URL-DB MCP server using Docker, optimized for AI assistant integration.

## Quick Start

### Build the Docker image

```bash
make docker-build
```

### Run in MCP stdio mode (for AI assistants)

```bash
make docker-run
```

### Run all services with Docker Compose

```bash
make docker-compose-up
```

## Docker Configuration

### Dockerfile Features

- **Multi-stage build**: Optimized image size (~50MB)
- **Alpine Linux base**: Minimal attack surface
- **Non-root user**: Enhanced security
- **Volume support**: Persistent database storage
- **Multiple modes**: Support for stdio, HTTP, SSE, and MCP-HTTP

### Docker Compose Services

The `docker-compose.yml` defines four services:

1. **url-db-mcp-stdio**: MCP stdio mode for AI assistants (Claude Desktop, Cursor)
2. **url-db-http**: REST API on port 8080
3. **url-db-mcp-sse**: Server-Sent Events on port 8081
4. **url-db-mcp-http**: HTTP-based MCP on port 8082

## Usage Examples

### 1. Run MCP server for AI assistants

```bash
# Using make command
make docker-run

# Using docker directly
docker run -it --rm \
  --name url-db-mcp \
  -v url-db-data:/data \
  url-db:latest
```

### 2. Start all services

```bash
# Start all services in background
make docker-compose-up

# View logs
make docker-logs

# Stop all services
make docker-compose-down
```

### 3. Use specific MCP mode

```bash
# HTTP mode
docker run -it --rm \
  -p 8080:8080 \
  -v url-db-data:/data \
  url-db:latest \
  -port=8080 -db-path=/data/url-db.sqlite

# SSE mode
docker run -it --rm \
  -p 8081:8081 \
  -v url-db-data:/data \
  url-db:latest \
  -mcp-mode=sse -port=8081 -db-path=/data/url-db.sqlite
```

## Claude Desktop Integration

To use with Claude Desktop, add to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "url-db-data:/data",
        "url-db:latest"
      ]
    }
  }
}
```

## Deployment to Container Registry

### Push to Docker Hub

```bash
# Tag and push to Docker Hub
docker tag url-db:latest your-dockerhub-username/url-db:latest
docker push your-dockerhub-username/url-db:latest

# Or use make command
make docker-push DOCKER_REGISTRY=your-dockerhub-username
```

### Push to private registry

```bash
# Tag and push to private registry
make docker-push DOCKER_REGISTRY=registry.example.com DOCKER_TAG=v1.0.0
```

## Environment Variables

- `DATABASE_URL`: Database connection string (default: `file:/data/url-db.sqlite`)
- `TOOL_NAME`: Tool name for composite keys (default: `url-db`)

## Volume Management

The Docker setup uses a named volume `url-db-data` for persistent storage:

```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect url-db-data

# Backup database
docker run --rm -v url-db-data:/data -v $(pwd):/backup alpine \
  cp /data/url-db.sqlite /backup/url-db-backup.sqlite

# Restore database
docker run --rm -v url-db-data:/data -v $(pwd):/backup alpine \
  cp /backup/url-db-backup.sqlite /data/url-db.sqlite
```

## Troubleshooting

### View container logs

```bash
docker logs url-db-mcp
```

### Access container shell

```bash
docker exec -it url-db-mcp sh
```

### Clean up resources

```bash
make docker-clean
```

## Security Considerations

1. **Non-root user**: Container runs as user `urldb` (UID 1000)
2. **Read-only filesystem**: Consider adding `--read-only` flag for production
3. **Resource limits**: Add `--memory` and `--cpus` limits as needed
4. **Network isolation**: Use custom networks for multi-container deployments

## Production Deployment

For production deployment, consider:

1. Using a reverse proxy (nginx, traefik) for HTTPS
2. Setting resource limits in docker-compose.yml
3. Implementing health checks
4. Using secrets management for sensitive configuration
5. Regular backups of the database volume