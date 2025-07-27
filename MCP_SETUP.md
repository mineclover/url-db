# URL-DB MCP 서버 설정 가이드

## Claude Desktop 설정

1. Claude Desktop 설정 파일 위치를 찾으세요:
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`

2. 설정 파일에 다음 내용을 추가하세요:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/Users/username/mcp/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "cwd": "/Users/username/mcp/url-db",
      "env": {
        "DATABASE_URL": "/Users/username/mcp/url-db/url-db.sqlite"
      }
    }
  }
}
```

3. Claude Desktop을 재시작하세요.

## Cursor 설정

1. Cursor 설정 파일: `.cursor/mcp.json`

```json
{
  "mcpServers": {
    "url-db-local": {
      "command": "/Users/username/mcp/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "cwd": "/Users/username/mcp/url-db",
      "env": {
        "DATABASE_URL": "/Users/username/mcp/url-db/url-db.sqlite"
      }
    }
  }
}
```

2. Cursor를 재시작하세요.

## SSE 모드 설정 (URL만 필요)

일부 MCP 클라이언트는 SSE (Server-Sent Events) 방식을 지원합니다. 이 경우 브릿지 없이 직접 URL만 설정하면 됩니다.

### 1. SSE 서버 시작
```bash
# 로컬에서 직접 실행
./bin/url-db -mcp-mode=sse -port=8080

# Docker로 실행
docker run -d -p 8080:8080 -v $(pwd)/data:/data --name url-db-sse asfdassdssa/url-db:latest -mcp-mode=sse
```

### 2. 클라이언트 설정 (예시)
```json
{
  "mcpServers": {
    "url-db-sse": {
      "url": "http://localhost:8080/mcp",
      "transport": "sse"
    }
  }
}
```

**참고**: 
- **stdio 모드**: Claude Desktop, Cursor 등 대부분의 MCP 클라이언트가 지원
- **SSE 모드**: Server-Sent Events를 통한 실시간 HTTP 통신, 웹 기반 클라이언트용
- **HTTP 모드**: 일반적인 HTTP JSON-RPC 2.0 통신, REST API처럼 사용 가능

## 사용 가능한 도구들

URL-DB MCP 서버는 30개의 도구를 제공합니다:

### 도메인 관리
- `list_domains` - 도메인 목록 조회
- `create_domain` - 새 도메인 생성
- `get_server_info` - 서버 정보 확인

### URL 관리
- `list_nodes` - URL 목록 조회
- `create_node` - URL 추가
- `get_node` - URL 상세 정보
- `update_node` - URL 정보 수정
- `delete_node` - URL 삭제
- `find_node_by_url` - URL로 검색

### 속성 관리
- `get_node_attributes` - URL 속성 조회
- `set_node_attributes` - URL 속성 설정
- `list_domain_attributes` - 도메인 속성 타입 목록
- `create_domain_attribute` - 새 속성 타입 정의
- `get_domain_attribute` - 속성 타입 상세 정보
- `update_domain_attribute` - 속성 타입 수정
- `delete_domain_attribute` - 속성 타입 삭제

### 템플릿 관리 (새로 추가!)
- `list_templates` - 템플릿 목록 조회
- `create_template` - 새 템플릿 생성
- `get_template` - 템플릿 상세 정보
- `update_template` - 템플릿 수정
- `delete_template` - 템플릿 삭제
- `clone_template` - 템플릿 복제
- `generate_template_scaffold` - 템플릿 스캐폴드 생성
- `validate_template` - 템플릿 검증

### 고급 기능
- `filter_nodes_by_attributes` - 속성 기반 URL 필터링
- `get_node_with_attributes` - 속성 포함 URL 상세 정보
- `create_dependency` - URL 간 의존성 생성
- `list_node_dependencies` - URL 의존성 목록
- `list_node_dependents` - URL 의존관계 목록
- `delete_dependency` - 의존성 삭제

## 문제 해결

### 로그 확인
Claude Desktop 로그: `/Users/username/Library/Logs/Claude/mcp-server-url-db*.log`

### 일반적인 문제들
1. **"Server disconnected" 에러**: 
   - `cwd` 설정이 올바른지 확인
   - 바이너리 경로가 정확한지 확인
   - 권한 문제가 없는지 확인

2. **데이터베이스 에러**:
   - `DATABASE_URL` 환경 변수 확인
   - 디렉토리에 쓰기 권한이 있는지 확인

3. **도구가 보이지 않음**:
   - 클라이언트를 완전히 재시작
   - 설정 파일 형식(JSON) 확인

## 버전 정보

- MCP 프로토콜 버전: 2025-06-18
- 서버 버전: 1.0.0
- 지원 도구: 30개
- 템플릿 시스템: 완전 구현