# Docker Hub 배포 가이드

## Docker Hub 계정: asfdassdssa

### 1. Docker Hub 로그인

먼저 Docker Hub에 로그인해야 합니다:

```bash
docker login
```

Username과 Password를 입력하세요:
- Username: `asfdassdssa`
- Password: Docker Hub 비밀번호

### 2. 이미지 푸시

로그인 후 다음 명령어로 이미지를 푸시합니다:

```bash
# 이미 태그된 이미지 푸시
docker push asfdassdssa/url-db:latest

# 또는 Makefile 사용
make docker-push DOCKER_REGISTRY=asfdassdssa
```

### 3. 푸시된 이미지 사용하기

#### Claude Desktop 설정

Claude Desktop의 설정 파일에 다음을 추가합니다:

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

#### 직접 실행

```bash
# MCP stdio 모드 (AI 어시스턴트용)
docker run -it --rm \
  -v url-db-data:/data \
  asfdassdssa/url-db:latest

# HTTP 서버 모드
docker run -d \
  -p 8080:8080 \
  -v url-db-data:/data \
  asfdassdssa/url-db:latest \
  -port=8080 -db-path=/data/url-db.sqlite

# MCP SSE 모드
docker run -d \
  -p 8081:8081 \
  -v url-db-data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse -port=8081 -db-path=/data/url-db.sqlite
```

### 4. Docker Compose 사용

docker-compose.yml 파일을 수정하여 Docker Hub 이미지를 사용할 수 있습니다:

```yaml
version: '3.8'

services:
  url-db-mcp-stdio:
    image: asfdassdssa/url-db:latest
    container_name: url-db-mcp-stdio
    volumes:
      - url-db-data:/data
    command: ["-mcp-mode=stdio", "-db-path=/data/url-db.sqlite"]
    stdin_open: true
    tty: true
    restart: unless-stopped

volumes:
  url-db-data:
    driver: local
```

### 5. 이미지 업데이트

새 버전을 배포할 때:

```bash
# 새 버전 빌드
make build
make docker-build

# 버전 태그 추가
docker tag url-db:latest asfdassdssa/url-db:v1.0.1
docker tag url-db:latest asfdassdssa/url-db:latest

# 모든 태그 푸시
docker push asfdassdssa/url-db:v1.0.1
docker push asfdassdssa/url-db:latest
```

### 6. 이미지 정보 확인

Docker Hub에서 이미지를 확인할 수 있습니다:
- URL: https://hub.docker.com/r/asfdassdssa/url-db

### 7. 다른 사용자가 사용하기

다른 사용자는 다음과 같이 간단히 사용할 수 있습니다:

```bash
# 이미지 다운로드 및 실행
docker run -it --rm asfdassdssa/url-db:latest

# 또는 docker-compose.yml에서 직접 사용
image: asfdassdssa/url-db:latest
```

## 주의사항

1. Docker Hub 무료 계정은 public 리포지토리만 지원합니다
2. 이미지가 public이므로 민감한 정보를 포함하지 마세요
3. 정기적으로 보안 업데이트를 확인하고 새 버전을 배포하세요