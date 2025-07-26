# MCP 클라이언트에서 SSE 서버 연결 가이드

URL-DB SSE 서버에 다양한 MCP 클라이언트를 연결하는 방법을 설명합니다.

## 목차
- [개요](#개요)
- [서버 준비](#서버-준비)
- [Claude Desktop 연결](#claude-desktop-연결)
- [Cursor 연결](#cursor-연결)
- [Continue 연결](#continue-연결)
- [MCP Python SDK 연결](#mcp-python-sdk-연결)
- [MCP TypeScript SDK 연결](#mcp-typescript-sdk-연결)
- [사용자 정의 클라이언트](#사용자-정의-클라이언트)
- [문제 해결](#문제-해결)

## 개요

URL-DB는 두 가지 MCP 연결 방식을 지원합니다:

| 모드 | 용도 | 클라이언트 예시 |
|------|------|----------------|
| **stdio** | AI 어시스턴트 | Claude Desktop, Cursor, Continue |
| **SSE** | HTTP 기반 클라이언트 | 웹 앱, 사용자 정의 도구, 원격 연결 |

이 가이드는 **SSE 모드** 연결에 대해 다룹니다.

## 서버 준비

### 1. SSE 서버 시작

```bash
# Docker로 SSE 서버 시작
docker run -d \
  --name url-db-sse \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse

# 서버 상태 확인
curl http://localhost:8080/health
```

### 2. 서버 정보 확인

```bash
# 서버 정보 조회
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

## Claude Desktop 연결

Claude Desktop은 기본적으로 stdio 모드를 사용하지만, HTTP 프록시를 통해 SSE 서버에 연결할 수 있습니다.

### HTTP 프록시 사용

1. **HTTP-to-stdio 브리지 스크립트 생성**:

```bash
# bridge-script.py 생성
cat > bridge-script.py << 'EOF'
#!/usr/bin/env python3
import sys
import json
import requests

SSE_ENDPOINT = "http://localhost:8080/mcp"

def send_to_sse(data):
    try:
        response = requests.post(SSE_ENDPOINT, 
                               headers={'Content-Type': 'application/json'},
                               json=data)
        
        # Parse SSE response
        lines = response.text.split('\n')
        for line in lines:
            if line.startswith('data: '):
                return json.loads(line[6:])
        return None
    except Exception as e:
        return {"error": str(e)}

# Read from stdin and forward to SSE
for line in sys.stdin:
    try:
        request = json.loads(line.strip())
        response = send_to_sse(request)
        if response:
            print(json.dumps(response))
            sys.stdout.flush()
    except Exception as e:
        error_response = {
            "jsonrpc": "2.0", 
            "id": None,
            "error": {"code": -32603, "message": str(e)}
        }
        print(json.dumps(error_response))
        sys.stdout.flush()
EOF

chmod +x bridge-script.py
```

2. **Claude Desktop 설정**:

```json
{
  "mcpServers": {
    "url-db-sse": {
      "command": "python3",
      "args": ["/path/to/bridge-script.py"]
    }
  }
}
```

## Cursor 연결

Cursor도 마찬가지로 HTTP 브리지를 사용합니다.

### Cursor 설정

1. **settings.json 열기**: `Ctrl/Cmd + Shift + P` → "Preferences: Open Settings (JSON)"

2. **MCP 서버 설정 추가**:

```json
{
  "cursor.experimental.mcpServers": {
    "url-db-sse": {
      "command": "python3",
      "args": ["/path/to/bridge-script.py"]
    }
  }
}
```

## Continue 연결

Continue의 경우 직접 HTTP 연결을 지원할 수 있습니다.

### Continue 설정

1. **config.json 수정**:

```json
{
  "mcpServers": [
    {
      "name": "url-db-sse",
      "serverUrl": "http://localhost:8080/mcp",
      "protocol": "http"
    }
  ]
}
```

## MCP Python SDK 연결

Python으로 SSE 서버에 연결하는 예제입니다.

### Python 클라이언트 예제

```python
import asyncio
import json
import aiohttp
from typing import Dict, Any, Optional

class MCPSSEClient:
    def __init__(self, endpoint: str = "http://localhost:8080/mcp"):
        self.endpoint = endpoint
        self.session: Optional[aiohttp.ClientSession] = None
        self.request_id = 0
        
    async def __aenter__(self):
        self.session = aiohttp.ClientSession()
        return self
        
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self.session:
            await self.session.close()
    
    async def send_request(self, method: str, params: Dict[str, Any] = None) -> Dict[str, Any]:
        self.request_id += 1
        request = {
            "jsonrpc": "2.0",
            "method": method,
            "params": params or {},
            "id": self.request_id
        }
        
        async with self.session.post(
            self.endpoint,
            headers={"Content-Type": "application/json"},
            json=request
        ) as response:
            text = await response.text()
            
            # Parse SSE format
            for line in text.split('\n'):
                if line.startswith('data: '):
                    return json.loads(line[6:])
            
            raise Exception("No data found in SSE response")
    
    async def initialize(self) -> Dict[str, Any]:
        return await self.send_request("initialize", {
            "protocolVersion": "2025-06-18",
            "capabilities": {"roots": {"listChanged": True}},
            "clientInfo": {"name": "python-mcp-client", "version": "1.0.0"}
        })
    
    async def list_tools(self) -> Dict[str, Any]:
        return await self.send_request("tools/list")
    
    async def call_tool(self, name: str, arguments: Dict[str, Any] = None) -> Dict[str, Any]:
        return await self.send_request("tools/call", {
            "name": name,
            "arguments": arguments or {}
        })

# 사용 예제
async def main():
    async with MCPSSEClient() as client:
        # 초기화
        init_result = await client.initialize()
        print("Initialized:", init_result)
        
        # 도구 목록 조회
        tools = await client.list_tools()
        print(f"Available tools: {len(tools['result']['tools'])}")
        
        # 도메인 생성
        domain = await client.call_tool("create_domain", {
            "name": "python-test",
            "description": "Python에서 생성한 도메인"
        })
        print("Domain created:", domain)
        
        # URL 추가
        node = await client.call_tool("create_node", {
            "domain_name": "python-test",
            "url": "https://python.org",
            "title": "Python 공식 사이트"
        })
        print("Node created:", node)

if __name__ == "__main__":
    asyncio.run(main())
```

## MCP TypeScript SDK 연결

TypeScript/JavaScript로 SSE 서버에 연결하는 예제입니다.

### TypeScript 클라이언트 예제

```typescript
interface JSONRPCRequest {
  jsonrpc: string;
  method: string;
  params?: any;
  id: number | string;
}

interface JSONRPCResponse {
  jsonrpc: string;
  result?: any;
  error?: any;
  id: number | string;
}

class MCPSSEClient {
  private endpoint: string;
  private requestId: number = 0;

  constructor(endpoint: string = 'http://localhost:8080/mcp') {
    this.endpoint = endpoint;
  }

  private async sendRequest(method: string, params?: any): Promise<any> {
    const request: JSONRPCRequest = {
      jsonrpc: '2.0',
      method,
      params: params || {},
      id: ++this.requestId
    };

    const response = await fetch(this.endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request)
    });

    const text = await response.text();
    
    // Parse SSE format
    const lines = text.split('\n');
    for (const line of lines) {
      if (line.startsWith('data: ')) {
        return JSON.parse(line.substring(6));
      }
    }
    
    throw new Error('No data found in SSE response');
  }

  async initialize(): Promise<JSONRPCResponse> {
    return this.sendRequest('initialize', {
      protocolVersion: '2025-06-18',
      capabilities: {
        roots: { listChanged: true }
      },
      clientInfo: {
        name: 'typescript-mcp-client',
        version: '1.0.0'
      }
    });
  }

  async listTools(): Promise<JSONRPCResponse> {
    return this.sendRequest('tools/list');
  }

  async callTool(name: string, arguments?: any): Promise<JSONRPCResponse> {
    return this.sendRequest('tools/call', {
      name,
      arguments: arguments || {}
    });
  }
}

// 사용 예제
async function main() {
  const client = new MCPSSEClient();
  
  try {
    // 초기화
    const initResult = await client.initialize();
    console.log('Initialized:', initResult);
    
    // 도구 목록 조회
    const tools = await client.listTools();
    console.log('Available tools:', tools.result.tools.length);
    
    // 도메인 생성
    const domain = await client.callTool('create_domain', {
      name: 'ts-test',
      description: 'TypeScript에서 생성한 도메인'
    });
    console.log('Domain created:', domain);
    
  } catch (error) {
    console.error('Error:', error);
  }
}

main();
```

## 사용자 정의 클라이언트

### 기본 HTTP 클라이언트 구현

어떤 언어든 HTTP POST 요청을 보낼 수 있다면 SSE 서버에 연결할 수 있습니다.

#### Go 예제

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
)

type JSONRPCRequest struct {
    JSONRPC string      `json:"jsonrpc"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params,omitempty"`
    ID      int         `json:"id"`
}

type JSONRPCResponse struct {
    JSONRPC string      `json:"jsonrpc"`
    Result  interface{} `json:"result,omitempty"`
    Error   interface{} `json:"error,omitempty"`
    ID      int         `json:"id"`
}

type MCPClient struct {
    endpoint  string
    requestID int
}

func NewMCPClient(endpoint string) *MCPClient {
    return &MCPClient{
        endpoint: endpoint,
        requestID: 0,
    }
}

func (c *MCPClient) SendRequest(method string, params interface{}) (*JSONRPCResponse, error) {
    c.requestID++
    
    req := JSONRPCRequest{
        JSONRPC: "2.0",
        Method:  method,
        Params:  params,
        ID:      c.requestID,
    }
    
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.Post(c.endpoint, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    // Parse SSE format
    lines := strings.Split(string(body), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "data: ") {
            var response JSONRPCResponse
            err := json.Unmarshal([]byte(line[6:]), &response)
            return &response, err
        }
    }
    
    return nil, fmt.Errorf("no data found in SSE response")
}

func main() {
    client := NewMCPClient("http://localhost:8080/mcp")
    
    // 초기화
    initResp, err := client.SendRequest("initialize", map[string]interface{}{
        "protocolVersion": "2025-06-18",
    })
    if err != nil {
        panic(err)
    }
    fmt.Println("Initialized:", initResp)
    
    // 도메인 생성
    domainResp, err := client.SendRequest("tools/call", map[string]interface{}{
        "name": "create_domain",
        "arguments": map[string]interface{}{
            "name":        "go-test",
            "description": "Go에서 생성한 도메인",
        },
    })
    if err != nil {
        panic(err)
    }
    fmt.Println("Domain created:", domainResp)
}
```

## 문제 해결

### 일반적인 연결 문제

#### 1. 연결 거부 (Connection Refused)

```bash
# 서버가 실행 중인지 확인
docker ps | grep url-db-sse

# 포트가 열려있는지 확인
curl http://localhost:8080/health

# 방화벽 확인 (필요시)
sudo ufw status
```

#### 2. SSE 응답 파싱 오류

```javascript
// 올바른 SSE 파싱 방법
const text = await response.text();
const lines = text.split('\n');
for (const line of lines) {
    if (line.startsWith('data: ')) {
        const data = JSON.parse(line.substring(6));
        return data;
    }
}
```

#### 3. CORS 오류 (브라우저 클라이언트)

```bash
# 개발 환경에서는 CORS가 허용되어 있음
# 필요시 프록시 사용
npx http-proxy-middleware --target http://localhost:8080
```

#### 4. 타임아웃 오류

```javascript
// 타임아웃 설정 증가
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), 30000); // 30초

const response = await fetch(endpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
    signal: controller.signal
});

clearTimeout(timeoutId);
```

### 디버깅 팁

#### 1. 서버 로그 확인

```bash
# Docker 로그 확인
docker logs -f url-db-sse

# 디버그 모드로 실행
docker run -d \
  --name url-db-sse-debug \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  -e LOG_LEVEL=debug \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

#### 2. 네트워크 테스트

```bash
# 기본 연결 테스트
curl -v http://localhost:8080/health

# MCP 프로토콜 테스트
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}' \
  -v
```

#### 3. 응답 형식 확인

```bash
# 응답 형식 상세 확인
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2025-06-18"},"id":1}' \
  | cat -A  # 숨겨진 문자 표시
```

## 성능 최적화

### 연결 풀링

```python
# aiohttp 연결 풀 설정
connector = aiohttp.TCPConnector(
    limit=100,  # 총 연결 수 제한
    limit_per_host=30,  # 호스트당 연결 수 제한
    keepalive_timeout=30,
    enable_cleanup_closed=True
)

session = aiohttp.ClientSession(connector=connector)
```

### 배치 요청

```javascript
// 여러 요청을 동시에 처리
const requests = [
    client.callTool('list_domains'),
    client.callTool('get_server_info'),
    client.callTool('list_tools')
];

const results = await Promise.all(requests);
console.log('Batch results:', results);
```

## 보안 고려사항

### 1. 네트워크 보안

```bash
# 로컬호스트만 허용
docker run -d \
  --name url-db-sse \
  -p 127.0.0.1:8080:8080 \
  -v $(pwd)/data:/data \
  asfdassdssa/url-db:latest \
  -mcp-mode=sse
```

### 2. 데이터 암호화

프로덕션 환경에서는 HTTPS 프록시 사용을 권장합니다:

```bash
# nginx 프록시 예제
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location /mcp {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

이 가이드를 통해 다양한 MCP 클라이언트에서 URL-DB SSE 서버에 성공적으로 연결할 수 있습니다!