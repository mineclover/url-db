# 노드 관리 API 엔드포인트

## 개요
노드(URL 컨텐츠) 생성, 조회, 수정, 삭제 기능을 제공하는 REST API  
API에서는 URL로 처리하지만 내부적으로는 nodeid 기반으로 관리

## 엔드포인트 목록

### 1. 노드 생성
- **POST** `/api/domains/{domain_id}/urls`
- **요청 본문**:
```json
{
  "url": "https://example.com/article",
  "title": "Example Article",
  "description": "This is an example article"
}
```
- **응답 (201)**:
```json
{
  "id": 1,
  "content": "https://example.com/article",
  "domain_id": 1,
  "title": "Example Article",
  "description": "This is an example article",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### 2. 도메인별 노드 목록 조회
- **GET** `/api/domains/{domain_id}/urls`
- **쿼리 파라미터**:
  - `page` (optional): 페이지 번호 (기본값: 1)
  - `size` (optional): 페이지 크기 (기본값: 20, 최대: 100)
  - `search` (optional): 검색어 (제목, 컨텐츠에서 검색)
- **응답 (200)**:
```json
{
  "nodes": [
    {
      "id": 1,
      "content": "https://example.com/article",
      "domain_id": 1,
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

### 3. 노드 조회
- **GET** `/api/urls/{id}`
- **응답 (200)**:
```json
{
  "id": 1,
  "content": "https://example.com/article",
  "domain_id": 1,
  "title": "Example Article",
  "description": "This is an example article",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### 4. 노드 수정
- **PUT** `/api/urls/{id}`
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
  "id": 1,
  "content": "https://example.com/article",
  "domain_id": 1,
  "title": "Updated Article Title",
  "description": "Updated description",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T01:00:00Z"
}
```

### 5. 노드 삭제
- **DELETE** `/api/urls/{id}`
- **응답 (204)**: 본문 없음

### 6. URL로 노드 ID 조회
- **POST** `/api/domains/{domain_id}/urls/find`
- **요청 본문**:
```json
{
  "url": "https://example.com/article"
}
```
- **응답 (200)** - 노드가 존재하는 경우:
```json
{
  "id": 1,
  "content": "https://example.com/article",
  "domain_id": 1,
  "title": "Example Article",
  "description": "This is an example article",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "composite_id": "url-db:tech-articles:1"
}
```
- **응답 (404)** - 노드가 존재하지 않는 경우:
```json
{
  "error": "NODE_NOT_FOUND",
  "message": "노드를 찾을 수 없습니다"
}
```

## MCP 서버 지원

이 API는 MCP (Model Context Protocol) 서버로도 동작합니다. MCP 전용 엔드포인트는 다음 문서를 참조하세요:

> MCP API 문서: [`06-mcp-api.md`](06-mcp-api.md)  
> 합성키 컨벤션: [`../spec/composite-key-conventions.md`](../spec/composite-key-conventions.md)

### 기존 API와 MCP API 차이점
- **식별자**: 기존 API는 내부 ID 사용, MCP API는 합성키 사용
- **응답 형식**: MCP API는 `composite_id` 필드 포함  
- **도메인 참조**: MCP API는 `domain_name` 사용, 기존 API는 `domain_id` 사용

## 컨텐츠 처리

### 원본 보존
- 입력된 URL을 그대로 저장
- 정규화 처리 없음
- 사용자가 입력한 형태 그대로 유지

## 에러 응답

> 에러 응답 형식: [`../spec/error-codes.md`](../spec/error-codes.md)  
> 노드 관련 에러: [`../spec/node-errors.md`](../spec/node-errors.md)

## 검증 규칙
- `url`: 필수, 최대 2048자 (형식 검증 없음)
- `title`: 선택, 최대 255자 (비어있으면 컨텐츠에서 자동 생성)
- `description`: 선택, 최대 1000자
- 컨텐츠와 domain_id는 수정 불가