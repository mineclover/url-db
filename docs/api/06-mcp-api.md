# MCP 서버 API 엔드포인트

## 개요
MCP (Model Context Protocol) 서버로 동작하여 AI 모델이 URL 데이터베이스를 활용할 수 있도록 지원하는 API입니다.
모든 응답에서 합성키(`tool_name:domain_name:id`)를 사용하여 리소스를 식별합니다.

> 합성키 컨벤션: [`../spec/composite-key-conventions.md`](../spec/composite-key-conventions.md)

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