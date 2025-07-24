# MCP 서버 로그 출력 문제 해결 문서

## 📋 개요

MCP (Model Context Protocol) 서버가 stdio 모드에서 실행될 때 일반 텍스트 로그가 JSON-RPC 프로토콜과 섞여서 클라이언트 파싱 오류를 발생시키는 문제를 해결한 과정을 문서화합니다.

## 🔍 문제 분석

### 1. 문제 현상
```
2025-07-24T06:35:24.722Z [url-db-local] [error] Unexpected token 'S', "Starting M"... is not valid JSON
```

### 2. 원인 분석
- **데이터베이스 초기화 로그**: `database.go`의 `loadSchemaFromFile()` 함수에서 stderr로 출력
- **MCP 서버 시작 로그**: `server.go`의 `Start()` 함수에서 stdout으로 출력
- **프로토콜 간섭**: 일반 텍스트 로그가 JSON-RPC 스트림에 섞여서 전송

### 3. 영향 범위
- MCP 클라이언트의 JSON 파싱 실패
- 서버-클라이언트 간 통신 중단
- 템플릿 검증 로직 실행 불가

## 🛠️ 해결 방법

### 1. 데이터베이스 로그 억제 (`internal/database/database.go`)

#### 수정 전
```go
fmt.Fprintf(os.Stderr, "[INFO] Schema loaded relative to executable: %s\n", schemaPath)
```

#### 수정 후
```go
// isMCPServerMode checks if the application is running in MCP server mode
func isMCPServerMode() bool {
    return os.Getenv("MCP_MODE") == "stdio" || 
           strings.Contains(strings.Join(os.Args, " "), "-mcp-mode=stdio")
}

// logInfo logs info message only if not in MCP stdio mode
func logInfo(format string, args ...interface{}) {
    if !isMCPServerMode() {
        fmt.Fprintf(os.Stderr, format, args...)
    }
}

// 사용 예시
logInfo("[INFO] Schema loaded relative to executable: %s\n", schemaPath)
```

### 2. MCP 서버 로그 억제 (`internal/interface/mcp/server.go`)

#### 수정 전
```go
fmt.Printf("Starting MCP server in %s mode\n", s.mode)
```

#### 수정 후
```go
// Don't log in stdio mode as it interferes with JSON-RPC communication
if s.mode != "stdio" {
    fmt.Printf("Starting MCP server in %s mode\n", s.mode)
}
```

## ✅ 검증 결과

### 1. 수정 전 테스트
```bash
echo '{"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}},"jsonrpc":"2.0","id":1}' | ./bin/url-db -mcp-mode=stdio -db-path=./url-db.db
```

**출력**:
```
Starting MCP server in stdio mode
[INFO] Schema loaded relative to executable: /Users/junwoobang/mcp/url-db/schema.sql
{"jsonrpc":"2.0","id":1,"result":{...}}
```

### 2. 수정 후 테스트
```bash
echo '{"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}},"jsonrpc":"2.0","id":1}' | ./bin/url-db -mcp-mode=stdio -db-path=./url-db.db
```

**출력**:
```
{"jsonrpc":"2.0","id":1,"result":{"capabilities":{"resources":{"listChanged":true,"subscribe":true},"tools":{"listChanged":true}},"protocolVersion":"2025-06-18","serverInfo":{"name":"url-db-mcp-server","version":"1.0.0"}}}
```

## 🎯 해결 효과

### 1. 즉시 효과
- ✅ JSON 파싱 오류 해결
- ✅ 순수한 JSON-RPC 통신
- ✅ 클라이언트-서버 간 안정적인 통신

### 2. 장기적 효과
- ✅ MCP를 통한 태그 속성 생성 가능
- ✅ 템플릿 생성 및 검증 로직 정상 실행
- ✅ 전체 MCP 도구 체인 정상 작동

## 📝 구현 세부사항

### 1. MCP 모드 감지 로직
```go
func isMCPServerMode() bool {
    return os.Getenv("MCP_MODE") == "stdio" || 
           strings.Contains(strings.Join(os.Args, " "), "-mcp-mode=stdio")
}
```

**감지 방법**:
- 환경 변수 `MCP_MODE=stdio` 확인
- 명령행 인수에 `-mcp-mode=stdio` 포함 여부 확인

### 2. 조건부 로그 출력
```go
func logInfo(format string, args ...interface{}) {
    if !isMCPServerMode() {
        fmt.Fprintf(os.Stderr, format, args...)
    }
}
```

**동작 원리**:
- MCP stdio 모드가 아닐 때만 로그 출력
- stdio 모드일 때는 로그 출력 억제

### 3. 적용된 파일들
- `internal/database/database.go`: 데이터베이스 초기화 로그 억제
- `internal/interface/mcp/server.go`: MCP 서버 시작/종료 로그 억제

## 🔧 유지보수 가이드

### 1. 새로운 로그 추가 시 주의사항
```go
// ❌ 잘못된 방법
fmt.Printf("로그 메시지\n")

// ✅ 올바른 방법
if !isMCPServerMode() {
    fmt.Printf("로그 메시지\n")
}
```

### 2. MCP 모드 감지 확장
```go
// 추가 환경 변수나 명령행 옵션을 감지하려면
func isMCPServerMode() bool {
    return os.Getenv("MCP_MODE") == "stdio" || 
           os.Getenv("MCP_STDIO") == "true" ||
           strings.Contains(strings.Join(os.Args, " "), "-mcp-mode=stdio")
}
```

### 3. 디버깅 모드 추가
```go
func logInfo(format string, args ...interface{}) {
    if !isMCPServerMode() || os.Getenv("MCP_DEBUG") == "true" {
        fmt.Fprintf(os.Stderr, format, args...)
    }
}
```

## 📚 관련 문서

- [MCP 서버 설정 가이드](./MCP_SERVER_CONFIGURATION.md)
- [MCP 테스팅 가이드](./MCP_TESTING_GUIDE.md)
- [템플릿 검증 로직 분석](./template_validation_flow_analysis.md)

## 🏷️ 태그

- `#MCP` `#logging` `#stdio` `#JSON-RPC` `#protocol` `#debugging` 