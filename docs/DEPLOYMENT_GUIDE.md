# URL-DB 배포 가이드

**완전한 3-in-1 Docker 이미지**: stdio, SSE 서버, SSE 클라이언트(브리지) 모두 포함

## 🚀 빠른 배포

### Docker Hub에서 바로 사용

```bash
# stdio 모드 (Claude Desktop/Cursor 등)
docker run -it --rm -v url-db-data:/data asfdassdssa/url-db:latest

# SSE 서버 모드 (HTTP API 제공)
docker run -d -p 8080:8080 -v url-db-data:/data --name url-db-sse asfdassdssa/url-db:latest -mcp-mode=sse

# HTTP 서버 모드 (REST API)
docker run -d -p 8080:8080 -v url-db-data:/data --name url-db-http asfdassdssa/url-db:latest -mcp-mode=http
```

## 📦 3가지 배포 시나리오

### 1. stdio 모드 - AI 어시스턴트 직접 연결

**용도**: Claude Desktop, Cursor 등에서 직접 MCP 서버로 사용

```bash
# Docker 컨테이너 실행
docker run -it --rm -v url-db-data:/data asfdassdssa/url-db:latest
```

**Claude Desktop 설정**:
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

### 2. SSE 서버 모드 - 중앙 집중식 서버

**용도**: 여러 클라이언트가 하나의 서버를 공유

```bash
# SSE 서버 시작
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v url-db-data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse

# 상태 확인
curl http://localhost:8080/health
```

### 3. SSE 클라이언트(브리지) - stdio를 SSE로 변환

**용도**: stdio 기반 클라이언트를 SSE 서버에 연결

```bash
# 브리지 실행 (서버가 localhost:8080에서 실행 중이어야 함)
docker run -it --rm --network host asfdassdssa/url-db:latest ./mcp-bridge -endpoint http://localhost:8080/mcp
```

**Claude Desktop 설정 (브리지 사용)**:
```json
{
  "mcpServers": {
    "url-db-sse": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm", "--network", "host",
        "asfdassdssa/url-db:latest", "./mcp-bridge",
        "-endpoint", "http://localhost:8080/mcp"
      ]
    }
  }
}
```

## 🏗️ 운영 환경 배포

### Docker Compose로 전체 스택

```yaml
# docker-compose.yml
version: '3.8'

services:
  # SSE 서버
  url-db-sse:
    image: asfdassdssa/url-db:latest
    command: ["-mcp-mode=sse", "-port=8080"]
    ports:
      - "8080:8080"
    volumes:
      - url-db-data:/data
    environment:
      - GIN_MODE=release
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    restart: unless-stopped

  # HTTP API 서버
  url-db-http:
    image: asfdassdssa/url-db:latest
    command: ["-mcp-mode=http", "-port=8081"]
    ports:
      - "8081:8081"
    volumes:
      - url-db-data:/data
    environment:
      - GIN_MODE=release
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    restart: unless-stopped

volumes:
  url-db-data:
```

```bash
# 전체 스택 시작
docker-compose up -d

# 로그 확인
docker-compose logs -f

# 스택 중지
docker-compose down
```

### Kubernetes 배포

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: url-db-sse
spec:
  replicas: 2
  selector:
    matchLabels:
      app: url-db-sse
  template:
    metadata:
      labels:
        app: url-db-sse
    spec:
      containers:
      - name: url-db
        image: asfdassdssa/url-db:latest
        args: ["-mcp-mode=sse", "-port=8080"]
        ports:
        - containerPort: 8080
        env:
        - name: GIN_MODE
          value: "release"
        volumeMounts:
        - name: data
          mountPath: /data
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: url-db-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: url-db-service
spec:
  selector:
    app: url-db-sse
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: url-db-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

```bash
kubectl apply -f k8s-deployment.yaml
```

## 🔧 설정 옵션

### 환경 변수

```bash
# 데이터베이스 경로
DATABASE_URL="file:/data/url-db.sqlite"

# 도구 이름 (composite key에 사용)
TOOL_NAME="url-db"

# Gin 모드 (release/debug)
GIN_MODE="release"

# 로그 레벨
LOG_LEVEL="info"
```

### 명령행 옵션

```bash
# 서버 옵션
-mcp-mode string    # MCP 서버 모드 (stdio, sse, http)
-port string        # 포트 (기본값: 8080)
-db-path string     # 데이터베이스 파일 경로
-tool-name string   # composite key에 사용할 도구 이름

# 브리지 옵션  
-endpoint string    # SSE 서버 엔드포인트
-timeout int        # 요청 타임아웃 (초)
-debug             # 디버그 로깅 활성화
```

## 📊 모니터링

### 헬스 체크

```bash
# 서버 상태 확인
curl http://localhost:8080/health

# 응답: {"status":"ok","timestamp":"2025-01-26T10:30:00Z"}
```

### 로그 모니터링

```bash
# Docker 로그
docker logs -f url-db-sse

# Docker Compose 로그
docker-compose logs -f url-db-sse

# Kubernetes 로그
kubectl logs -f deployment/url-db-sse
```

### 메트릭 수집 (Prometheus 호환)

```bash
# 메트릭 엔드포인트 (향후 지원 예정)
curl http://localhost:8080/metrics
```

## 🔒 보안 설정

### HTTPS 지원 (Reverse Proxy)

```yaml
# nginx.conf
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 네트워크 보안

```bash
# 내부 네트워크만 허용
docker run -d \
  --name url-db-sse \
  --network internal \
  -p 127.0.0.1:8080:8080 \
  -v url-db-data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

## 🚨 문제 해결

### 일반적인 이슈

| 문제 | 원인 | 해결 방법 |
|------|------|-----------|
| Container won't start | 포트 충돌 | `lsof -i :8080`로 포트 확인 |
| Database locked | 여러 인스턴스 실행 | 단일 인스턴스만 실행 |
| Permission denied | 볼륨 권한 | `docker volume rm url-db-data` 후 재생성 |
| Bridge connection failed | SSE 서버 미실행 | SSE 서버 먼저 시작 |

### 디버그 모드

```bash
# 디버그 모드로 실행
docker run -it --rm \
  -v url-db-data:/data \
  -e LOG_LEVEL=debug \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse

# 브리지 디버그
docker run -it --rm --network host \
  asfdassdssa/url-db:latest \
  ./mcp-bridge -debug -endpoint http://localhost:8080/mcp
```

## 📈 성능 튜닝

### 컨테이너 리소스 제한

```bash
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v url-db-data:/data \
  --memory=512m \
  --cpus=1.0 \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

### 데이터베이스 최적화

```bash
# SSD 스토리지 사용 권장
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v /fast-ssd/url-db:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

## 📚 추가 리소스

- **GitHub**: https://github.com/your-username/url-db
- **Docker Hub**: https://hub.docker.com/r/asfdassdssa/url-db
- **문서**: `/docs` 디렉터리의 추가 가이드들
- **이슈 트래킹**: GitHub Issues 페이지