# MCP JSON-RPC 구현 작업

## 현재 상황

### 문제점
1. 현재 MCP stdio 서버는 단순한 텍스트 명령어 인터페이스만 지원
2. Claude Desktop과의 통신에 필요한 JSON-RPC 2.0 프로토콜 미지원
3. MCP 표준 스펙을 따르지 않아 Claude와 제대로 통신 불가

### 완료된 작업
1. 빌드 에러 수정
   - `internal/services` 패키지의 인터페이스 포인터 문제 해결
   - int64 -> int 타입 변환 문제 해결
2. 프로젝트 빌드 성공 (`bin/url-db` 생성)
3. Claude MCP에 서버 등록 완료

## 필요한 작업

### 1. JSON-RPC 2.0 프로토콜 구현
- [ ] JSON-RPC 요청/응답 구조체 정의
- [ ] JSON 파싱 및 직렬화 로직 추가
- [ ] 에러 처리 표준화

### 2. MCP 표준 메서드 구현
- [ ] `initialize` - 서버 초기화 및 capability 교환
- [ ] `initialized` - 초기화 완료 알림
- [ ] `tools/list` - 사용 가능한 도구 목록 반환
- [ ] `tools/call` - 도구 실행
- [ ] `resources/list` - 리소스 목록 반환
- [ ] `resources/read` - 리소스 읽기

### 3. 기존 기능 매핑
현재 텍스트 명령어를 MCP 도구로 변환:
- `list_domains` → `list_mcp_domains` tool
- `list_nodes` → `list_mcp_nodes` tool
- `create_node` → `create_mcp_node` tool
- `get_node` → `get_mcp_node` tool
- `update_node` → `update_mcp_node` tool
- `delete_node` → `delete_mcp_node` tool
- 추가 도구들...

### 4. 리소스 구현
- `mcp://nodes/{composite_id}` - 개별 노드 리소스
- `mcp://domains/{domain_name}` - 도메인 리소스
- `mcp://domains/{domain_name}/nodes` - 도메인 내 노드 목록

## 구현 계획

### Phase 1: 기본 JSON-RPC 구조
1. `internal/mcp/jsonrpc.go` 파일 생성
2. Request/Response 구조체 정의
3. JSON 파싱 로직 구현

### Phase 2: MCP 핸들러
1. `internal/mcp/stdio_server.go` 리팩토링
2. JSON-RPC 메시지 처리 로직 추가
3. 표준 입출력 통신 구현

### Phase 3: 도구 및 리소스
1. 도구 정의 및 등록
2. 리소스 정의 및 등록
3. 핸들러 메서드 구현

### Phase 4: 테스트 및 검증
1. 단위 테스트 작성
2. Claude Desktop과의 통합 테스트
3. 에러 케이스 처리

## 예상 파일 변경

- `internal/mcp/stdio_server.go` - 전면 리팩토링
- `internal/mcp/jsonrpc.go` - 신규 생성
- `internal/mcp/tools.go` - 신규 생성
- `internal/mcp/resources.go` - 신규 생성
- `internal/mcp/handlers.go` - 신규 생성

## 참고 자료

- [MCP Specification](https://spec.modelcontextprotocol.io/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [MCP TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk)

## 예상 소요 시간

- Phase 1: 2-3시간
- Phase 2: 3-4시간
- Phase 3: 4-5시간
- Phase 4: 2-3시간

총 예상: 11-15시간