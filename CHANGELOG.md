# 변경 이력

## [미출시]

### 추가됨
- MCP JSON-RPC 구현 작업 문서 추가
- docs/tasks 디렉토리 구조 생성

### 변경됨
- `internal/services` 패키지의 인터페이스 포인터 문제 수정
  - NodeRepository를 포인터에서 인터페이스로 변경
  - int64를 int로 타입 변환 추가
- 문서 업데이트
  - installation-guide.md에 최신 업데이트 섹션 추가
  - mcp-server-setup-guide.md에 현재 상태 및 주의사항 추가
  - README.md에 개발 작업 섹션 추가

### 수정됨
- 빌드 에러 해결
  - dependency_service.go: NodeRepository 인터페이스 타입 수정
  - subscription_service.go: NodeRepository 인터페이스 타입 수정
  - event_service.go: NodeRepository 인터페이스 타입 수정

### 알려진 문제
- MCP stdio 서버가 JSON-RPC 2.0 프로토콜을 지원하지 않음
- Claude Desktop과의 실제 통신 불가능
- 단순 텍스트 명령어 인터페이스만 구현됨

## [2025.07.18]

### 초기 커밋
- 프로젝트 기본 구조 설정
- MCP 서버 기본 구현
- 외부 종속성 관리 시스템 구현