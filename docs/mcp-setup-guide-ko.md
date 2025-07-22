# URL-DB MCP 설정 가이드

URL-DB의 MCP (Model Context Protocol) 설정 방법을 설명합니다.

## 개요

URL-DB는 두 가지 MCP 모드를 지원합니다:

1. **stdio 모드** - AI 어시스턴트 연동용 (Claude Desktop, Cursor)
2. **sse 모드** - HTTP 기반 웹 서비스용

## Stdio 모드 (AI 어시스턴트용)

### Claude Desktop 설정

1. 설정 파일 위치:
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`
   - Linux: `~/.config/Claude/claude_desktop_config.json`

2. 설정 내용:
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/Users/your-name/mcp/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {
        "DATABASE_URL": "file:/Users/your-name/mcp/url-db/url-db.sqlite"
      }
    }
  }
}
```

### Cursor 설정

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/path/to/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/path/to/database.sqlite"
      ]
    }
  }
}
```

### 사용 가능한 MCP 도구 (18개)

#### 도메인 관리
- `list_domains` - 모든 도메인 조회
- `create_domain` - 새 도메인 생성

#### 노드(URL) 관리
- `list_nodes` - 노드 목록 조회 (페이징, 검색)
- `create_node` - 새 URL 노드 생성
- `get_node` - 노드 상세 조회
- `update_node` - 노드 정보 수정
- `delete_node` - 노드 삭제
- `find_node_by_url` - URL로 노드 찾기

#### 속성 관리
- `get_node_attributes` - 노드 속성 조회
- `set_node_attributes` - 노드 속성 설정

#### 도메인 스키마
- `list_domain_attributes` - 도메인 속성 정의 목록
- `create_domain_attribute` - 속성 정의 생성
- `get_domain_attribute` - 속성 정의 조회
- `update_domain_attribute` - 속성 설명 수정
- `delete_domain_attribute` - 속성 정의 삭제

#### 고급 쿼리 (신규)
- `get_node_with_attributes` - 노드와 속성을 한번에 조회
- `filter_nodes_by_attributes` - 속성값으로 노드 필터링

#### 서버 정보
- `get_server_info` - 서버 정보 조회

### 사용 예시

Claude에서의 대화 예시:

```
사용자: 기술 관련 북마크를 저장할 도메인을 만들어줘

Claude: create_domain 도구를 사용하여 "tech" 도메인을 생성하겠습니다.

사용자: 방금 만든 도메인에 파이썬 공식 문서 URL을 추가하고 
       "programming", "python", "documentation" 태그를 달아줘

Claude: create_node로 URL을 추가하고 set_node_attributes로 태그를 설정하겠습니다.
```

## SSE 모드 (웹 서비스용)

### 서버 시작

```bash
# SSE 모드로 시작 (기본값)
./bin/url-db -mcp-mode=sse -port=8080

# 또는 단순히 (SSE가 기본 모드)
./bin/url-db
```

### 설정 옵션

```bash
# 포트 변경
./bin/url-db -port=3000

# 데이터베이스 경로 지정
./bin/url-db -db-path=/path/to/database.sqlite

# 도구 이름 변경 (복합 키에 영향)
./bin/url-db -tool-name=my-url-db
```

### REST API 엔드포인트

#### 도메인 관리
```bash
# 도메인 생성
curl -X POST http://localhost:8080/api/domains \
  -H "Content-Type: application/json" \
  -d '{"name":"tech","description":"기술 자료"}'

# 도메인 목록
curl http://localhost:8080/api/domains
```

#### 노드(URL) 관리
```bash
# URL 추가
curl -X POST http://localhost:8080/api/domains/1/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url":"https://docs.python.org",
    "title":"Python Documentation",
    "description":"파이썬 공식 문서"
  }'

# URL 목록 조회
curl http://localhost:8080/api/domains/1/urls
```

#### 속성 관리
```bash
# 속성 추가
curl -X POST http://localhost:8080/api/urls/1/attributes \
  -H "Content-Type: application/json" \
  -d '{
    "attribute_id": 1,
    "value": "programming"
  }'
```

## 주요 차이점

### Stdio 모드
- **용도**: AI 어시스턴트 연동
- **프로토콜**: JSON-RPC 2.0
- **도구**: 18개 MCP 도구
- **고급 기능**: 속성 필터링, 일괄 조회
- **키 형식**: `tool-name:domain:id`

### SSE 모드
- **용도**: 웹 애플리케이션
- **프로토콜**: REST API
- **엔드포인트**: 표준 REST
- **인증**: 직접 구현 필요
- **ID**: 일반 정수 ID

## 도메인 스키마 시스템

URL-DB는 도메인별 스키마를 강제합니다:

1. 도메인 생성 시 속성 타입 정의
2. 노드는 도메인에 정의된 속성만 사용 가능
3. 지원 타입: tag, ordered_tag, number, string, markdown, image

### 스키마 예시

```bash
# 1. 도메인 생성
도메인: "bookmarks"

# 2. 속성 정의
- category (tag): 카테고리
- priority (number): 중요도 (1-5)
- notes (markdown): 상세 메모

# 3. 노드는 위 속성만 사용 가능
```

## 문제 해결

### Stdio 모드
- **응답 없음**: JSON-RPC 형식 확인
- **도구 못 찾음**: 도구 이름 철자 확인
- **복합 키 오류**: `tool:domain:id` 형식 확인

### SSE 모드
- **포트 사용 중**: `-port` 플래그로 변경
- **CORS 오류**: 기본적으로 CORS 헤더 포함됨
- **데이터베이스 잠김**: 단일 인스턴스만 실행

## 고급 사용법

### 속성 필터링 (stdio 모드)

```json
{
  "method": "tools/call",
  "params": {
    "name": "filter_nodes_by_attributes",
    "arguments": {
      "domain_name": "tech",
      "filters": [
        {
          "name": "category",
          "value": "python",
          "operator": "equals"
        }
      ]
    }
  }
}
```

### 일괄 조회 (stdio 모드)

```json
{
  "method": "tools/call",
  "params": {
    "name": "get_node_with_attributes",
    "arguments": {
      "composite_id": "url-db:tech:123"
    }
  }
}
```

## 보안 고려사항

1. **Stdio 모드**: AI 어시스턴트 권한으로 실행
2. **SSE 모드**: 공개 노출 시 인증 구현 필요
3. **데이터베이스**: 적절한 파일 권한 설정
4. **복합 키**: 추측 불가능한 도메인 컨텍스트 포함

## 향후 개선 사항

- SSE 실시간 이벤트 스트리밍
- WebSocket 지원
- SSE 모드 고급 필터링
- 대량 가져오기/내보내기 도구