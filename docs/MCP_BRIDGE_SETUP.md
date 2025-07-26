# MCP Bridge 설정 가이드

Go로 작성된 MCP Bridge를 사용하여 stdio 기반 MCP 클라이언트를 SSE 서버에 연결하는 방법입니다.

## 🚀 빠른 시작

### 1. 빌드

```bash
# 서버와 브리지 동시 빌드
make build

# 브리지만 빌드하려면
go build -o bin/mcp-bridge cmd/bridge/main.go
```

### 2. SSE 서버 시작

```bash
# Docker 사용
docker run -d -p 8080:8080 -v $(pwd)/data:/data --name url-db-sse asfdassdssa/url-db:latest -mcp-mode=sse

# 로컬 빌드 사용
./bin/url-db -mcp-mode=sse -port=8080
```

### 3. 브리지 테스트

```bash
# 기본 설정으로 테스트
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05"},"id":1}' | ./bin/mcp-bridge
```

## 📋 클라이언트별 설정

### Claude Desktop

`~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/path/to/url-db/bin/mcp-bridge",
      "args": ["-endpoint", "http://localhost:8080/mcp"],
      "env": {
        "DEBUG": "1"
      }
    }
  }
}
```

### Cursor

`.cursor/config.json`:

```json
{
  "cursor.experimental.mcpServers": {
    "url-db": {
      "command": "/path/to/url-db/bin/mcp-bridge",
      "args": ["-endpoint", "http://localhost:8080/mcp", "-debug"],
      "env": {
        "TIMEOUT": "60"
      }
    }
  }
}
```

### VS Code (MCP Extension)

`.vscode/mcp.json`:

```json
{
  "mcpServers": [
    {
      "name": "url-db",
      "command": "/path/to/url-db/bin/mcp-bridge",
      "args": ["-endpoint", "http://localhost:8080/mcp"],
      "env": {
        "DEBUG": "1",
        "TIMEOUT": "45"
      }
    }
  ]
}
```

## 🔧 브리지 옵션

### 명령행 플래그

```bash
./bin/mcp-bridge [옵션]

옵션:
  -endpoint string   SSE 서버 엔드포인트 (기본값: http://localhost:8080/mcp)
  -timeout int       요청 타임아웃 (초) (기본값: 30)
  -debug            디버그 로깅 활성화
  -help             도움말 표시
  -version          버전 정보 표시
```

### 환경 변수

```bash
export SSE_ENDPOINT=http://localhost:8080/mcp  # 서버 엔드포인트
export TIMEOUT=30                              # 타임아웃 (초)
export DEBUG=1                                 # 디버그 모드 (빈 값이 아닌 경우)
```

## 🌐 Docker 환경에서 사용

### Docker Compose로 전체 스택 실행

```yaml
version: '3.8'
services:
  url-db-sse:
    image: asfdassdssa/url-db:latest
    command: ["-mcp-mode=sse", "-port=8080"]
    ports:
      - "8080:8080"
    volumes:
      - url-db-data:/data
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  url-db-data:
```

### 브리지를 Docker 컨테이너로 실행

```bash
# Docker 이미지에 포함된 브리지 사용
docker run -it --rm --network host asfdassdssa/url-db:latest ./mcp-bridge -endpoint http://localhost:8080/mcp -debug
```

## 🐛 문제 해결

### 연결 문제

```bash
# 서버 상태 확인
curl http://localhost:8080/health

# 브리지 디버그 모드로 실행
./bin/mcp-bridge -debug -endpoint http://localhost:8080/mcp

# 포트 사용 확인
lsof -i :8080
```

### 일반적인 오류

| 오류 메시지 | 원인 | 해결 방법 |
|-------------|------|-----------|
| `connection refused` | SSE 서버가 실행되지 않음 | 서버 시작 확인 |
| `timeout` | 응답 시간 초과 | `-timeout` 값 증가 |
| `invalid JSON` | 잘못된 요청 형식 | JSON-RPC 2.0 형식 확인 |
| `no data found` | SSE 응답 파싱 실패 | 서버 응답 형식 확인 |

### 로그 분석

디버그 모드에서 브리지는 stderr로 다음 정보를 출력합니다:

```
[DEBUG] Bridge started, SSE endpoint: http://localhost:8080/mcp
[DEBUG] Sending request: {"jsonrpc":"2.0","method":"initialize",...}
[DEBUG] Received response: {"jsonrpc":"2.0","result":{...},"id":1}
```

## 📊 성능 최적화

### HTTP 클라이언트 튜닝

브리지는 내부적으로 다음 설정을 사용합니다:

- **Connection Pooling**: 연결 재사용으로 성능 향상
- **Timeout**: 기본 30초, 환경에 따라 조정 가능
- **Keep-Alive**: HTTP 연결 유지로 지연 시간 감소

### 대용량 응답 처리

```bash
# 타임아웃 증가 (대용량 데이터 처리시)
./bin/mcp-bridge -timeout 120 -endpoint http://localhost:8080/mcp
```

## 🔒 보안 고려사항

### HTTPS 사용

```bash
# HTTPS 엔드포인트 사용
./bin/mcp-bridge -endpoint https://your-domain.com/mcp
```

### 인증 헤더 (향후 지원 예정)

현재 버전은 기본 인증을 지원하지 않습니다. 필요시 다음과 같이 확장 가능합니다:

```go
// 향후 지원 예정
req.Header.Set("Authorization", "Bearer "+token)
```

## 📈 모니터링

### 상태 확인

```bash
# 브리지 상태 확인 (stdin으로 테스트)
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | ./bin/mcp-bridge

# 서버 상태 확인
curl http://localhost:8080/health
```

### 성능 메트릭

디버그 모드에서 요청/응답 시간을 모니터링할 수 있습니다:

```bash
# 시간 측정과 함께 실행
time echo '{"jsonrpc":"2.0","method":"list_domains","id":1}' | ./bin/mcp-bridge
```

## 🚀 고급 사용법

### 여러 서버 연결

각기 다른 엔드포인트에 대해 별도의 브리지 인스턴스를 실행:

```bash
# 개발 서버
./bin/mcp-bridge -endpoint http://dev-server:8080/mcp &

# 프로덕션 서버  
./bin/mcp-bridge -endpoint https://prod-server/mcp &
```

### 스크립트 자동화

```bash
#!/bin/bash
# start-bridge.sh

# 서버 시작 대기
until curl -s http://localhost:8080/health > /dev/null; do 
  echo "Waiting for server..."
  sleep 1
done

# 브리지 시작
exec ./bin/mcp-bridge -endpoint http://localhost:8080/mcp -debug
```

## 💡 팁

- **디버그 모드**: 개발 시 항상 `-debug` 플래그 사용
- **타임아웃 조정**: 네트워크 환경에 따라 타임아웃 값 조정
- **환경 변수**: 반복 사용시 환경 변수 활용
- **로그 파일**: 운영 환경에서는 로그를 파일로 리다이렉션
- **자동 재시작**: systemd 또는 supervisor로 자동 재시작 설정