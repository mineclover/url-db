# MCP 서버 설정 가이드

## 개요
URL-DB는 MCP (Model Context Protocol) 서버를 지원하여 AI 모델과의 직접 통합이 가능합니다. 이 가이드는 MCP 서버를 설정하고 AI 클라이언트와 연결하는 방법을 설명합니다.

## MCP 서버 모드

URL-DB는 두 가지 MCP 서버 모드를 지원합니다:

### 1. stdio 모드 (권장)
- 표준 입력/출력을 통한 통신
- 로컬 AI 도구와의 직접 연결
- 보안성이 높음

### 2. SSE (Server-Sent Events) 모드
- HTTP 서버를 통한 통신
- 웹 기반 AI 도구와 연결
- 네트워크를 통한 원격 접근 가능

## 설치 및 실행

### 1. 기본 설치
```bash
# 소스 코드 다운로드
git clone https://github.com/yourusername/url-db.git
cd url-db

# 빌드
./build.sh  # Unix/Linux/macOS
# 또는
build.bat   # Windows
```

### 2. MCP stdio 모드 실행
```bash
# Unix/Linux/macOS
./url-db -mcp-mode=stdio

# Windows
url-db.exe -mcp-mode=stdio

# DATABASE_URL을 인자로 전달하는 경우
./url-db -mcp-mode=stdio DATABASE_URL=file:~/mcp/url-db/url-db.db
```

### 3. MCP SSE 모드 실행 (기본값)
```bash
# Unix/Linux/macOS
./bin/url-db -mcp-mode=sse

# Windows
bin\url-db.exe -mcp-mode=sse
```

## Claude Desktop과의 연결

### 1. Claude MCP 명령어 방식 (권장)

가장 간단한 방법은 Claude MCP 명령어를 사용하는 것입니다:

```bash
# 기본 설정
claude mcp add url-db /path/to/url-db/bin/url-db --args="-mcp-mode=stdio"

# 환경변수 포함 설정
claude mcp add url-db /path/to/url-db/bin/url-db \
  --args="-mcp-mode=stdio" \
  --env="DATABASE_URL=file:/path/to/url-db/url-db.db" \
  --env="TOOL_NAME=url-db"

# 현재 프로젝트 경로 예시
claude mcp add url-db /Users/junwoobang/mcp/url-db/bin/url-db \
  --args="-mcp-mode=stdio" \
  --env="DATABASE_URL=file:/Users/junwoobang/mcp/url-db/url-db.db"

# 실행 파일과 인자를 함께 전달하는 방법 (권장)
# 중요: -- 를 사용하여 claude mcp add의 옵션과 프로그램의 인자를 구분
claude mcp add url-db -- ~/mcp/url-db/url-db -mcp-mode=stdio DATABASE_URL=file:~/mcp/url-db/url-db.db
```

### 2. 수동 설정 파일 편집 방식

#### 설정 파일 위치
```bash
# macOS
~/Library/Application Support/Claude/claude_desktop_config.json

# Windows
%APPDATA%\Claude\claude_desktop_config.json

# Linux
~/.config/Claude/claude_desktop_config.json
```

#### stdio 모드 설정
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/path/to/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {
        "DATABASE_URL": "file:/path/to/url-db/url-db.db",
        "TOOL_NAME": "url-db"
      }
    }
  }
}
```

#### SSE 모드 설정
```json
{
  "mcpServers": {
    "url-db": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-fetch", "http://localhost:8080/mcp"],
      "env": {}
    }
  }
}
```

## 고급 설정

### 1. 환경변수 설정
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/path/to/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {
        "DATABASE_URL": "file:/path/to/url-db/url-db.db",
        "TOOL_NAME": "my-url-db",
        "LOG_LEVEL": "debug",
        "PORT": "8080"
      }
    }
  }
}
```

### 2. 다중 인스턴스 설정

#### Claude MCP 명령어 방식
```bash
# 개인용 인스턴스
claude mcp add url-db-personal /path/to/url-db/bin/url-db \
  --args="-mcp-mode=stdio" \
  --env="DATABASE_URL=file:/path/to/personal/url-db.db" \
  --env="TOOL_NAME=personal-db"

# 업무용 인스턴스
claude mcp add url-db-work /path/to/url-db/bin/url-db \
  --args="-mcp-mode=stdio" \
  --env="DATABASE_URL=file:/path/to/work/url-db.db" \
  --env="TOOL_NAME=work-db"
```

#### 수동 설정 파일 방식
```json
{
  "mcpServers": {
    "url-db-personal": {
      "command": "/path/to/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {
        "DATABASE_URL": "file:/path/to/personal/url-db.db",
        "TOOL_NAME": "personal-db"
      }
    },
    "url-db-work": {
      "command": "/path/to/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {
        "DATABASE_URL": "file:/path/to/work/url-db.db",
        "TOOL_NAME": "work-db"
      }
    }
  }
}
```

## 기타 AI 도구와의 연결

### 1. Cline (VS Code Extension)
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/path/to/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"]
    }
  }
}
```

### 2. Continue (VS Code Extension)
```json
{
  "mcpServers": [
    {
      "name": "url-db",
      "command": "/path/to/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"]
    }
  ]
}
```

### 3. 커스텀 MCP 클라이언트
```python
import asyncio
import json
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def main():
    server_params = StdioServerParameters(
        command="/path/to/url-db/bin/url-db",
        args=["-mcp-mode=stdio"],
        env={
            "DATABASE_URL": "file:/path/to/url-db/url-db.db",
            "TOOL_NAME": "url-db"
        }
    )
    
    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            # MCP 서버와 상호작용
            await session.initialize()
            
            # 도구 목록 가져오기
            tools = await session.list_tools()
            print("Available tools:", tools)
            
            # 도구 실행
            result = await session.call_tool(
                "create_mcp_node",
                {
                    "domain_name": "test",
                    "url": "https://example.com",
                    "title": "Test URL"
                }
            )
            print("Result:", result)

if __name__ == "__main__":
    asyncio.run(main())
```

## MCP 서버 기능

### 1. 지원하는 도구 (Tools)
- `create_mcp_node`: 새 노드 생성
- `get_mcp_node`: 노드 조회
- `update_mcp_node`: 노드 업데이트
- `delete_mcp_node`: 노드 삭제
- `list_mcp_nodes`: 노드 목록 조회
- `find_mcp_node_by_url`: URL로 노드 찾기
- `batch_get_mcp_nodes`: 노드 배치 조회
- `get_mcp_node_attributes`: 노드 속성 조회
- `set_mcp_node_attributes`: 노드 속성 설정
- `list_mcp_domains`: 도메인 목록 조회
- `create_mcp_domain`: 도메인 생성
- `get_mcp_server_info`: 서버 정보 조회

### 2. 지원하는 리소스 (Resources)
- `mcp://nodes/{composite_id}`: 개별 노드 리소스
- `mcp://domains/{domain_name}`: 도메인 리소스
- `mcp://domains/{domain_name}/nodes`: 도메인 내 노드 목록

### 3. 합성키 형식
```
{tool_name}:{domain_name}:{node_id}
```

예시:
```
url-db:tech-articles:123
my-db:bookmarks:456
```

## 빠른 시작 가이드

### 1. 자동 설정 (권장) 🎯
```bash
# 설정 스크립트 실행
cd /path/to/url-db
./setup-mcp.sh

# 출력된 명령어 중 하나를 복사해서 실행
```

### 2. 수동 설정 명령어
```bash
# 1단계: URL-DB 빌드 (필요한 경우)
cd /path/to/url-db
./build.sh

# 2단계: Claude MCP에 추가 (다음 두 방법 중 하나 사용)
# 방법 1: 환경변수 사용 (권장)
claude mcp add url-db /path/to/url-db/bin/url-db \
  --args="-mcp-mode=stdio" \
  --env="DATABASE_URL=file:/path/to/url-db/url-db.db"

# 방법 2: CLI 인자로 직접 전달
claude mcp add url-db -- /path/to/url-db/bin/url-db -mcp-mode=stdio DATABASE_URL=file:/path/to/url-db/url-db.db

# 3단계: Claude Desktop 재시작
```

### 실제 예시 (현재 설치 경로)
```bash
# junwoobang 사용자의 경우
claude mcp add url-db "/Users/junwoobang/mcp/url-db/bin/url-db" \
  --args="-mcp-mode=stdio" \
  --env="DATABASE_URL=file:/Users/junwoobang/mcp/url-db/url-db.db"
```

### ✅ 테스트 완료 상태
- **MCP JSON-RPC 2.0 프로토콜**: 완전히 구현됨 (2025-07-22)
- **테스트 점수**: 92% (LLM-as-a-Judge), 100% (통합 테스트)
- **모든 11개 도구**: 정상 작동 확인

### 2. 확인 방법
```bash
# 등록된 MCP 서버 확인
claude mcp list

# 특정 서버 삭제 (필요시)
claude mcp remove url-db
```

## 사용 예시

### 1. 기본 사용법
```bash
# Claude Desktop에서 사용
"url-db에서 새로운 도메인 'tech-articles' 생성해줘"
"https://example.com/article을 tech-articles 도메인에 추가해줘"
"tech-articles 도메인의 모든 URL 목록 보여줘"
```

### 2. 속성 관리
```bash
# 속성 설정
"방금 추가한 URL에 카테고리를 'javascript'로 설정해줘"
"priority를 'high'로 설정해줘"

# 속성 조회
"이 URL의 모든 속성 보여줘"
```

### 3. 고급 검색
```bash
# 배치 조회
"다음 composite ID들의 노드 정보 가져와줘: url-db:tech:1, url-db:tech:2"

# URL로 찾기
"https://example.com/article 이 URL이 어느 도메인에 있는지 찾아줘"
```

## 트러블슈팅

### 1. 연결 문제

#### 문제: Claude Desktop에서 MCP 서버를 찾을 수 없음
```bash
# 해결책 1: 경로 확인
which url-db
# 또는
ls -la /path/to/url-db/bin/url-db

# 해결책 2: 권한 확인
chmod +x /path/to/url-db/bin/url-db

# 해결책 3: 실행 테스트
/path/to/url-db/bin/url-db -mcp-mode=stdio
```

#### 문제: stdio 모드에서 응답 없음
```bash
# 해결책 1: 데이터베이스 파일 확인
ls -la /path/to/url-db/url-db.db

# 해결책 2: 로그 확인
DATABASE_URL=file:/path/to/url-db/url-db.db LOG_LEVEL=debug /path/to/url-db/bin/url-db -mcp-mode=stdio
```

### 2. 성능 문제

#### 문제: 응답 속도 느림
```bash
# 해결책 1: 데이터베이스 최적화
sqlite3 /path/to/url-db/url-db.db "VACUUM; ANALYZE;"

# 해결책 2: 메모리 설정
export GOGC=100
export GOMEMLIMIT=1GiB
```

### 3. 설정 문제

#### 문제: 환경변수 인식 안됨
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/bin/bash",
      "args": ["-c", "cd /path/to/url-db && ./bin/url-db -mcp-mode=stdio"],
      "env": {
        "DATABASE_URL": "file:/path/to/url-db/url-db.db"
      }
    }
  }
}
```

## 보안 고려사항

### 1. 데이터베이스 보안
```bash
# 데이터베이스 파일 권한 설정
chmod 600 /path/to/url-db/url-db.db
chown $USER:$USER /path/to/url-db/url-db.db
```

### 2. 실행 파일 보안
```bash
# 실행 파일 권한 설정
chmod 755 /path/to/url-db/bin/url-db
```

### 3. 네트워크 보안 (SSE 모드)
```bash
# 방화벽 설정
sudo ufw allow from 127.0.0.1 to any port 8080
```

## 고급 기능

### 1. 커스텀 도구 이름 설정
```json
{
  "mcpServers": {
    "my-bookmarks": {
      "command": "/path/to/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {
        "TOOL_NAME": "my-bookmarks",
        "DATABASE_URL": "file:/path/to/bookmarks.db"
      }
    }
  }
}
```

### 2. 로그 레벨 설정
```json
{
  "env": {
    "LOG_LEVEL": "debug"  // trace, debug, info, warn, error
  }
}
```

### 3. 성능 튜닝
```json
{
  "env": {
    "GOGC": "100",
    "GOMEMLIMIT": "1GiB"
  }
}
```

## 개발자 가이드

### 1. MCP 서버 테스트
```bash
# 직접 테스트
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}' | ./bin/url-db -mcp-mode=stdio
```

### 2. 로그 분석
```bash
# 디버그 로그 활성화
LOG_LEVEL=debug ./bin/url-db -mcp-mode=stdio 2>debug.log

# 로그 확인
tail -f debug.log
```

### 3. 프로파일링
```bash
# 메모리 프로파일링
go tool pprof http://localhost:8080/debug/pprof/heap

# CPU 프로파일링
go tool pprof http://localhost:8080/debug/pprof/profile
```

## 참고 자료

### 1. MCP 공식 문서
- [MCP 사양](https://spec.modelcontextprotocol.io/)
- [MCP SDK](https://github.com/modelcontextprotocol/typescript-sdk)

### 2. URL-DB 문서
- [API 문서](api/06-mcp-api.md)
- [설치 가이드](installation-guide.md)
- [메인 문서](README.md)

### 3. 예제 프로젝트
- [MCP 서버 예제](https://github.com/modelcontextprotocol/servers)
- [Claude Desktop 설정 예제](https://docs.anthropic.com/claude/docs/desktop-configuration)

## 지원 및 문의

- GitHub Issues: https://github.com/yourusername/url-db/issues
- MCP 커뮤니티: https://github.com/modelcontextprotocol/specification/discussions
- 문서: https://docs.url-db.com
- 이메일: support@url-db.com