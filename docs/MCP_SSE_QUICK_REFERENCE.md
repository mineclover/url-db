# MCP SSE 연결 빠른 참조

## 🚀 서버 시작

```bash
# Docker로 SSE 서버 시작
docker run -d -p 8080:8080 -v $(pwd)/data:/data --name url-db-sse asfdassdssa/url-db:latest -mcp-mode=sse

# 상태 확인
curl http://localhost:8080/health
```

## 📋 클라이언트별 설정

### Claude Desktop

**Note**: Claude Desktop currently supports stdio mode. For SSE usage, you would need a custom bridge implementation:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/path/to/url-db", 
      "args": ["-mcp-mode=stdio"]
    }
  }
}
```

### Cursor

For Cursor, use stdio mode configuration:

```json
{
  "cursor.experimental.mcpServers": {
    "url-db": {
      "command": "/path/to/url-db",
      "args": ["-mcp-mode=stdio"]
    }
  }
}
```

### Continue

```json
{
  "mcpServers": [{
    "name": "url-db-sse",
    "serverUrl": "http://localhost:8080/mcp",
    "protocol": "http"
  }]
}
```

## 🔧 직접 SSE 서버 사용법

```bash
# 빌드
make build

# SSE 모드로 서버 시작
./bin/url-db -mcp-mode=sse

# 포트 지정
./bin/url-db -mcp-mode=sse -port=8081

# 디버그 모드
LOG_LEVEL=debug ./bin/url-db -mcp-mode=sse

# Docker로 SSE 서버 시작
docker run -d -p 8080:8080 -v $(pwd)/data:/data --name url-db-sse asfdassdssa/url-db:latest -mcp-mode=sse
```

## 🌐 HTTP 클라이언트 예제

### cURL

```bash
# 초기화
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2025-06-18"},"id":1}'

# 도구 목록
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":2}'

# 도메인 생성
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_domain","arguments":{"name":"test","description":"테스트 도메인"}},"id":3}'
```

### JavaScript

```javascript
async function callMCP(method, params = {}) {
  const response = await fetch('http://localhost:8080/mcp', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      jsonrpc: '2.0', method, params, id: Date.now()
    })
  });
  
  const text = await response.text();
  const data = text.split('\n').find(l => l.startsWith('data: '));
  return JSON.parse(data.substring(6));
}

// 사용 예제
const init = await callMCP('initialize', {protocolVersion: '2025-06-18'});
const tools = await callMCP('tools/list');
const domain = await callMCP('tools/call', {
  name: 'create_domain',
  arguments: {name: 'js-test', description: 'JavaScript 테스트'}
});
```

### Python

```python
import requests, json

def call_mcp(method, params=None):
    response = requests.post('http://localhost:8080/mcp', json={
        'jsonrpc': '2.0', 'method': method, 'params': params or {}, 'id': 1
    })
    for line in response.text.split('\n'):
        if line.startswith('data: '):
            return json.loads(line[6:])

# 사용 예제
init = call_mcp('initialize', {'protocolVersion': '2025-06-18'})
tools = call_mcp('tools/list')
domain = call_mcp('tools/call', {
    'name': 'create_domain',
    'arguments': {'name': 'py-test', 'description': 'Python 테스트'}
})
```

## 🎯 자주 사용하는 도구 호출

### 도메인 관리

```bash
# 도메인 목록
{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_domains","arguments":{}},"id":1}

# 도메인 생성
{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_domain","arguments":{"name":"bookmarks","description":"북마크 모음"}},"id":2}
```

### URL 관리

```bash
# URL 추가
{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_node","arguments":{"domain_name":"bookmarks","url":"https://example.com","title":"예시 사이트"}},"id":3}

# URL 목록
{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_nodes","arguments":{"domain_name":"bookmarks"}},"id":4}

# 전체 컨텐츠 스캔
{"jsonrpc":"2.0","method":"tools/call","params":{"name":"scan_all_content","arguments":{"domain_name":"bookmarks","max_tokens_per_page":3000}},"id":5}
```

## 🐛 문제 해결

### 연결 확인

```bash
# 서버 상태
curl http://localhost:8080/health

# 포트 사용량 확인
lsof -i :8080

# Docker 컨테이너 상태
docker ps | grep url-db-sse

# 로그 확인
docker logs -f url-db-sse
```

### 일반적인 오류

| 오류 | 해결 방법 |
|------|-----------|
| Connection Refused | 서버가 시작되었는지 확인 |
| Invalid JSON | 요청 형식 확인 |
| Method not found | 메서드 이름 확인 |
| CORS Error | 브라우저에서 직접 호출 시 프록시 사용 |

### 디버그 모드

```bash
# 디버그 로그와 함께 실행
docker run -d -p 8080:8080 -v $(pwd)/data:/data -e LOG_LEVEL=debug --name url-db-sse asfdassdssa/url-db:latest -mcp-mode=sse
```

## 💡 팁

- **응답 파싱**: 항상 `data: ` 접두사 제거 후 JSON 파싱
- **요청 ID**: 각 요청마다 고유한 ID 사용
- **에러 처리**: JSON-RPC 2.0 에러 형식 확인
- **연결 풀링**: 여러 요청 시 연결 재사용
- **타임아웃**: 긴 작업의 경우 타임아웃 설정 증가