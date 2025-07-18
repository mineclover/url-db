# MCP 서버 API 엔드포인트

## 개요
MCP (Model Context Protocol) 서버로 동작하여 AI 모델이 URL 데이터베이스를 활용할 수 있도록 지원하는 API입니다.
모든 응답에서 합성키(`tool_name:domain_name:id`)를 사용하여 리소스를 식별합니다.

> 합성키 컨벤션: [`../spec/composite-key-conventions.md`](../spec/composite-key-conventions.md)

## MCP 서버 모드

### 1. HTTP/SSE 모드 (기본값)
RESTful API 엔드포인트를 통해 MCP 기능을 제공합니다.

### 2. stdio 모드
표준 입출력을 통한 대화형 명령 인터페이스를 제공합니다.

#### stdio 모드 명령어
```bash
# 서버 실행
./url-db -mcp-mode=stdio DATABASE_URL=file:~/mcp/url-db/url-db.db

# 사용 가능한 명령어
> help                                    # 도움말 표시
> list_domains                           # 모든 도메인 목록
> list_nodes <domain_name>               # 도메인 내 노드 목록
> create_node <domain> <url> [title]     # 새 노드 생성
> get_node <composite_id>                # 노드 상세 정보
> update_node <composite_id> <title>     # 노드 제목 수정
> delete_node <composite_id>             # 노드 삭제
> server_info                            # 서버 정보
> quit                                   # 세션 종료

# 예시
> create_node tech-articles https://example.com/article "Example Article"
> get_node url-db:tech-articles:123
> list_nodes tech-articles
```

## 엔드포인트 목록

### 1. 노드 생성 (MCP)
- **POST** `/api/mcp/nodes`
- **요청 본문**:
```json
{
  "domain_name": "tech-articles",
  "url": "https://example.com/article",
  "title": "Example Article", 
  "description": "This is an example article"
}
```
- **응답 (201)**:
```json
{
  "composite_id": "url-db:tech-articles:123",
  "url": "https://example.com/article",
  "domain_name": "tech-articles",
  "title": "Example Article",
  "description": "This is an example article",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### 2. 노드 조회 (MCP)
- **GET** `/api/mcp/nodes/{composite_id}`
- **예시**: `GET /api/mcp/nodes/url-db:tech-articles:123`
- **응답 (200)**:
```json
{
  "composite_id": "url-db:tech-articles:123",
  "url": "https://example.com/article",
  "domain_name": "tech-articles", 
  "title": "Example Article",
  "description": "This is an example article",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### 3. 노드 목록 조회 (MCP)
- **GET** `/api/mcp/nodes`
- **쿼리 파라미터**:
  - `domain_name` (optional): 도메인 필터링
  - `page` (optional): 페이지 번호 (기본값: 1)
  - `size` (optional): 페이지 크기 (기본값: 20, 최대: 100)
  - `search` (optional): 검색어 (제목, URL에서 검색)
- **응답 (200)**:
```json
{
  "nodes": [
    {
      "composite_id": "url-db:tech-articles:123",
      "url": "https://example.com/article",
      "domain_name": "tech-articles",
      "title": "Example Article", 
      "description": "This is an example article",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "total_count": 1,
  "page": 1,
  "size": 20,
  "total_pages": 1
}
```

### 4. 노드 수정 (MCP)
- **PUT** `/api/mcp/nodes/{composite_id}`
- **예시**: `PUT /api/mcp/nodes/url-db:tech-articles:123`
- **요청 본문**:
```json
{
  "title": "Updated Article Title",
  "description": "Updated description"
}
```
- **응답 (200)**:
```json
{
  "composite_id": "url-db:tech-articles:123",
  "url": "https://example.com/article",
  "domain_name": "tech-articles",
  "title": "Updated Article Title",
  "description": "Updated description", 
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T01:00:00Z"
}
```

### 5. 노드 삭제 (MCP)
- **DELETE** `/api/mcp/nodes/{composite_id}`
- **예시**: `DELETE /api/mcp/nodes/url-db:tech-articles:123`
- **응답 (204)**: 본문 없음

### 6. URL로 노드 찾기 (MCP)
- **POST** `/api/mcp/nodes/find`
- **요청 본문**:
```json
{
  "domain_name": "tech-articles",
  "url": "https://example.com/article"
}
```
- **응답 (200)** - 노드가 존재하는 경우:
```json
{
  "composite_id": "url-db:tech-articles:123",
  "url": "https://example.com/article",
  "domain_name": "tech-articles",
  "title": "Example Article",
  "description": "This is an example article",
  "created_at": "2023-01-01T00:00:00Z", 
  "updated_at": "2023-01-01T00:00:00Z"
}
```
- **응답 (404)** - 노드가 존재하지 않는 경우:
```json
{
  "error": "NODE_NOT_FOUND",
  "message": "노드를 찾을 수 없습니다",
  "domain_name": "tech-articles",
  "url": "https://example.com/article"
}
```

### 7. 배치 조회 (MCP)
- **POST** `/api/mcp/nodes/batch`
- **요청 본문**:
```json
{
  "composite_ids": [
    "url-db:tech-articles:123",
    "url-db:recipes:456",
    "url-db:personal-bookmarks:789"
  ]
}
```
- **응답 (200)**:
```json
{
  "nodes": [
    {
      "composite_id": "url-db:tech-articles:123",
      "url": "https://example.com/article",
      "domain_name": "tech-articles",
      "title": "Example Article",
      "description": "This is an example article",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "not_found": [
    "url-db:recipes:999"
  ]
}
```

## 도메인 관리 (MCP)

### 8. 도메인 목록 조회 (MCP)
- **GET** `/api/mcp/domains`
- **응답 (200)**:
```json
{
  "domains": [
    {
      "name": "tech-articles",
      "description": "기술 관련 아티클",
      "node_count": 150,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

### 9. 도메인 생성 (MCP)
- **POST** `/api/mcp/domains`
- **요청 본문**:
```json
{
  "name": "tech-articles",
  "description": "기술 관련 아티클"
}
```
- **응답 (201)**:
```json
{
  "name": "tech-articles",
  "description": "기술 관련 아티클",
  "node_count": 0,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

## 노드 속성 관리 (MCP)

### 10. 노드 속성 조회 (MCP)
- **GET** `/api/mcp/nodes/{composite_id}/attributes`
- **예시**: `GET /api/mcp/nodes/url-db:tech-articles:123/attributes`
- **응답 (200)**:
```json
{
  "composite_id": "url-db:tech-articles:123",
  "attributes": [
    {
      "name": "category",
      "type": "tag",
      "value": "javascript"
    },
    {
      "name": "rating",
      "type": "number", 
      "value": "4.5"
    }
  ]
}
```

### 11. 노드 속성 설정 (MCP)
- **PUT** `/api/mcp/nodes/{composite_id}/attributes`
- **요청 본문**:
```json
{
  "attributes": [
    {
      "name": "category",
      "value": "javascript"
    },
    {
      "name": "rating",
      "value": "4.5"
    }
  ]
}
```
- **응답 (200)**:
```json
{
  "composite_id": "url-db:tech-articles:123",
  "attributes": [
    {
      "name": "category",
      "type": "tag",
      "value": "javascript"
    },
    {
      "name": "rating", 
      "type": "number",
      "value": "4.5"
    }
  ]
}
```

## 에러 응답

### 합성키 관련 에러
```json
{
  "error": "INVALID_COMPOSITE_KEY",
  "message": "합성키 형식이 올바르지 않습니다",
  "expected_format": "tool_name:domain_name:id",
  "provided": "invalid-key"
}
```

### 도메인 관련 에러
```json
{
  "error": "DOMAIN_NOT_FOUND",
  "message": "지정된 도메인을 찾을 수 없습니다",
  "domain_name": "non-existent-domain"
}
```

### 일반 에러
> 기본 에러 응답 형식: [`../spec/error-codes.md`](../spec/error-codes.md)

## MCP 서버 메타데이터

### 서버 정보
- **GET** `/api/mcp/server/info`
- **응답 (200)**:
```json
{
  "name": "url-db",
  "version": "1.0.0",
  "description": "URL 데이터베이스 MCP 서버",
  "capabilities": [
    "resources",
    "tools",
    "prompts"
  ],
  "composite_key_format": "url-db:domain_name:id"
}
```

## 검증 규칙

### 합성키 검증
- 형식: `tool_name:domain_name:id`
- tool_name: `url-db` (고정)
- domain_name: 영문자, 숫자, 하이픈만 허용
- id: 양의 정수

### 요청 본문 검증
- `url`: 필수, 최대 2048자
- `title`: 선택, 최대 255자
- `description`: 선택, 최대 1000자
- `domain_name`: 필수, 영문자, 숫자, 하이픈만 허용, 최대 50자

## Claude Desktop 통합 예시

### 도메인 및 노드 관리
```text
사용자: "url-db에 'tech-articles' 도메인을 만들어줘"
AI: tech-articles 도메인을 생성했습니다.

사용자: "https://example.com/react-tutorial을 tech-articles에 추가해줘"
AI: URL을 tech-articles 도메인에 추가했습니다. (composite_id: url-db:tech-articles:123)

사용자: "방금 추가한 URL의 제목을 'React Tutorial 2024'로 변경해줘"
AI: 노드의 제목을 업데이트했습니다.
```

### 속성 관리
```text
사용자: "url-db:tech-articles:123에 category 속성을 'frontend'로 설정해줘"
AI: category 속성을 설정했습니다.

사용자: "같은 URL에 priority를 'high'로, rating을 '5'로 설정해줘"
AI: priority와 rating 속성을 설정했습니다.

사용자: "이 URL의 모든 속성을 보여줘"
AI: url-db:tech-articles:123의 속성:
- category: frontend (tag)
- priority: high (string)
- rating: 5 (number)
```

### 검색 및 조회
```text
사용자: "tech-articles 도메인의 모든 URL을 보여줘"
AI: tech-articles 도메인의 URL 목록입니다...

사용자: "https://example.com/react-tutorial이 어느 도메인에 있는지 찾아줘"
AI: 해당 URL은 tech-articles 도메인에 있습니다. (composite_id: url-db:tech-articles:123)

사용자: "url-db:tech-articles:123, url-db:tech-articles:124의 정보를 한번에 가져와줘"
AI: 요청하신 노드들의 정보입니다...
```

## 프로그래밍 언어별 통합 예시

### Python
```python
import requests

# MCP 서버 기본 URL
BASE_URL = "http://localhost:8080/api/mcp"

# 노드 생성
def create_node(domain_name, url, title=None, description=None):
    response = requests.post(f"{BASE_URL}/nodes", json={
        "domain_name": domain_name,
        "url": url,
        "title": title,
        "description": description
    })
    return response.json()

# 노드 조회
def get_node(composite_id):
    response = requests.get(f"{BASE_URL}/nodes/{composite_id}")
    return response.json()

# 속성 설정
def set_node_attributes(composite_id, attributes):
    response = requests.put(
        f"{BASE_URL}/nodes/{composite_id}/attributes",
        json={"attributes": attributes}
    )
    return response.json()

# 사용 예시
node = create_node("tech-articles", "https://example.com/article", "My Article")
print(f"Created node: {node['composite_id']}")

# 속성 추가
set_node_attributes(node['composite_id'], [
    {"name": "category", "value": "javascript"},
    {"name": "priority", "value": "high"}
])
```

### JavaScript/TypeScript
```typescript
// MCP 클라이언트 클래스
class MCPClient {
  constructor(private baseUrl: string = 'http://localhost:8080/api/mcp') {}

  async createNode(params: {
    domain_name: string;
    url: string;
    title?: string;
    description?: string;
  }) {
    const response = await fetch(`${this.baseUrl}/nodes`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(params)
    });
    return response.json();
  }

  async getNode(compositeId: string) {
    const response = await fetch(`${this.baseUrl}/nodes/${compositeId}`);
    return response.json();
  }

  async setNodeAttributes(compositeId: string, attributes: Array<{
    name: string;
    value: string;
  }>) {
    const response = await fetch(
      `${this.baseUrl}/nodes/${compositeId}/attributes`,
      {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ attributes })
      }
    );
    return response.json();
  }
}

// 사용 예시
const client = new MCPClient();

const node = await client.createNode({
  domain_name: 'tech-articles',
  url: 'https://example.com/article',
  title: 'My Article'
});

console.log(`Created node: ${node.composite_id}`);
```

### Go
```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type MCPClient struct {
    BaseURL string
}

type CreateNodeRequest struct {
    DomainName  string `json:"domain_name"`
    URL         string `json:"url"`
    Title       string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
}

type MCPNode struct {
    CompositeID string `json:"composite_id"`
    URL         string `json:"url"`
    DomainName  string `json:"domain_name"`
    Title       string `json:"title"`
}

func (c *MCPClient) CreateNode(req CreateNodeRequest) (*MCPNode, error) {
    data, _ := json.Marshal(req)
    resp, err := http.Post(
        c.BaseURL+"/nodes",
        "application/json",
        bytes.NewBuffer(data),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var node MCPNode
    if err := json.NewDecoder(resp.Body).Decode(&node); err != nil {
        return nil, err
    }
    return &node, nil
}

// 사용 예시
func main() {
    client := &MCPClient{BaseURL: "http://localhost:8080/api/mcp"}
    
    node, err := client.CreateNode(CreateNodeRequest{
        DomainName: "tech-articles",
        URL:        "https://example.com/article",
        Title:      "My Article",
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created node: %s\n", node.CompositeID)
}
```

## 고급 사용 시나리오

### 1. 배치 처리
```python
# 여러 URL을 한 번에 조회
composite_ids = [
    "url-db:tech-articles:123",
    "url-db:tech-articles:124",
    "url-db:recipes:456"
]

response = requests.post(f"{BASE_URL}/nodes/batch", json={
    "composite_ids": composite_ids
})
result = response.json()

print(f"Found {len(result['nodes'])} nodes")
print(f"Not found: {result['not_found']}")
```

### 2. 도메인별 노드 관리
```javascript
// 도메인 생성 후 노드 추가
async function setupDomain(domainName: string, description: string) {
  // 도메인 생성
  await client.createDomain({ name: domainName, description });
  
  // 초기 노드들 추가
  const urls = [
    { url: 'https://example.com/1', title: 'Article 1' },
    { url: 'https://example.com/2', title: 'Article 2' }
  ];
  
  for (const item of urls) {
    await client.createNode({
      domain_name: domainName,
      ...item
    });
  }
}
```

### 3. 속성 기반 워크플로우
```python
# 우선순위가 높은 항목에 속성 추가
def mark_high_priority_items(domain_name):
    # 도메인의 모든 노드 조회
    nodes = get_domain_nodes(domain_name)
    
    for node in nodes:
        if should_be_high_priority(node):
            set_node_attributes(node['composite_id'], [
                {"name": "priority", "value": "high"},
                {"name": "reviewed", "value": "false"}
            ])
```