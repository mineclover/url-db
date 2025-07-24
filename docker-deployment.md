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

### 1. Docker Volume (기본 설정)

기본적인 Docker 볼륨을 사용하는 설정:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "url-db-data:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 2. 로컬 디렉토리 마운트 (권장)

호스트의 특정 폴더에 SQLite 파일을 저장:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "/Users/username/url-db-data:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 3. 현재 사용자 홈 디렉토리

홈 디렉토리 하위에 저장 (macOS/Linux):

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "~/Documents/url-db:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 4. 프로젝트별 데이터베이스

특정 프로젝트 폴더에 저장:

```json
{
  "mcpServers": {
    "url-db-project": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "/Users/username/projects/my-project/db:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 5. Windows 사용자용 설정

Windows 경로를 사용하는 경우:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "C:/Users/username/url-db-data:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 6. 데스크톱에 저장

데스크톱 폴더에 데이터베이스 저장:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "/Users/username/Desktop/url-db:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 7. 여러 환경 설정

개발용과 운영용 분리:

```json
{
  "mcpServers": {
    "url-db-dev": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "/Users/username/dev/url-db:/data",
        "asfdassdssa/url-db:latest"
      ]
    },
    "url-db-prod": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "/Users/username/prod/url-db:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 8. 권한 문제 해결 (Linux/macOS)

사용자 권한을 매핑하여 권한 문제 해결:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "/Users/username/url-db-data:/data",
        "-u", "1000:1000",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 💡 **경로 설정 팁**

1. **절대 경로 사용**: 상대 경로보다 절대 경로 권장
2. **폴더 미리 생성**: Docker 실행 전에 호스트 폴더 생성
3. **권한 확인**: 폴더에 읽기/쓰기 권한이 있는지 확인
4. **백슬래시 주의**: Windows에서는 경로 구분자 주의

### 📂 **SQLite 파일 위치 확인**

설정 후 다음 경로에서 SQLite 파일을 확인할 수 있습니다:
- 로컬 경로: `{your-path}/url-db.sqlite`
- 직접 접근: `sqlite3 {your-path}/url-db.sqlite`

### 🎯 **실제 사용 예시**

현재 사용자의 구체적인 설정 예시:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "/Users/username/url-db-data:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

이 설정을 사용하면:
- SQLite 파일: `/Users/username/url-db-data/url-db.sqlite`
- 직접 접근: `sqlite3 /Users/username/url-db-data/url-db.sqlite`
- Finder에서: `open /Users/username/url-db-data/`

### 🔧 **폴더 생성 및 테스트**

설정 전에 폴더를 미리 생성하고 테스트:

```bash
# 1. 폴더 생성
mkdir -p /Users/username/url-db-data

# 2. 권한 확인
ls -la /Users/username/url-db-data

# 3. 테스트 실행
docker run -it --rm \
  -v /Users/username/url-db-data:/data \
  asfdassdssa/url-db:latest

# 4. 데이터베이스 파일 확인
ls -la /Users/username/url-db-data/
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