# MCP Claude Desktop 설정 가이드 (상세)

이 가이드는 URL-DB를 Claude Desktop과 연동하기 위한 상세한 설정 방법과 문제 해결 방법을 제공합니다.

## 🚨 중요 사항

MCP 서버를 Claude Desktop과 연동하려면 **반드시** 다음 순서를 따라야 합니다:

1. URL-DB 빌드 및 설치
2. Claude Desktop 설정 파일 수정
3. Claude Desktop 완전 재시작
4. 연결 확인

## 📦 1단계: URL-DB 빌드 및 설치

### 1.1 프로젝트 빌드

```bash
# 프로젝트 디렉토리로 이동
cd /Users/junwoobang/mcp/url-db

# 의존성 설치 및 빌드
make deps
make build

# 빌드 확인
ls -la ./bin/url-db
```

### 1.2 빌드된 바이너리 테스트

```bash
# MCP 모드로 실행 테스트
./bin/url-db -mcp-mode=stdio -h

# 버전 확인
./bin/url-db -version
```

### 1.3 실행 권한 확인

```bash
# 실행 권한 부여 (필요한 경우)
chmod +x ./bin/url-db
```

## 🔧 2단계: Claude Desktop 설정

### 2.1 설정 파일 위치

**macOS**: 
```bash
~/Library/Application Support/Claude/claude_desktop_config.json
```

### 2.2 설정 파일 백업

```bash
# 기존 설정 백업
cp "~/Library/Application Support/Claude/claude_desktop_config.json" \
   "~/Library/Application Support/Claude/claude_desktop_config.backup.json"
```

### 2.3 설정 파일 수정

설정 파일을 열어 다음과 같이 수정합니다:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "/Users/junwoobang/mcp/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/Users/junwoobang/mcp/url-db/url-db.sqlite"
      ],
      "env": {}
    }
  }
}
```

**⚠️ 주의사항:**
- `command` 경로는 **절대 경로**여야 합니다
- 홈 디렉토리 축약형(`~`)을 사용하지 마세요
- 경로에 공백이 있으면 큰따옴표로 감싸야 합니다

### 2.4 여러 MCP 서버가 있는 경우

기존 MCP 서버가 있다면 다음과 같이 추가합니다:

```json
{
  "mcpServers": {
    "existing-server": {
      // 기존 서버 설정
    },
    "url-db": {
      "command": "/Users/junwoobang/mcp/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-db-path=/Users/junwoobang/mcp/url-db/url-db.sqlite"
      ],
      "env": {}
    }
  }
}
```

## 🔄 3단계: Claude Desktop 재시작

### 3.1 완전한 재시작 방법

1. Claude Desktop 완전 종료:
   - 메뉴바의 Claude 아이콘 우클릭 → "Quit Claude"
   - 또는 `Cmd + Q`로 완전 종료

2. 프로세스 확인 (터미널):
   ```bash
   # Claude 프로세스가 완전히 종료되었는지 확인
   ps aux | grep -i claude
   
   # 남아있는 프로세스가 있다면 강제 종료
   pkill -f Claude
   ```

3. Claude Desktop 재실행

### 3.2 설정 적용 확인

Claude가 시작되면 개발자 도구를 열어 확인할 수 있습니다:
- `View` → `Developer Tools` → `Console`

## ✅ 4단계: 연결 확인

### 4.1 Claude에서 MCP 도구 확인

새 대화를 시작하고 다음을 입력해보세요:

```
url-db의 서버 정보를 확인해줘
```

또는

```
사용 가능한 MCP 도구들을 나열해줘
```

### 4.2 기본 기능 테스트

```
1. tech라는 도메인을 만들어줘
2. tech 도메인에 https://github.com URL을 추가해줘
3. tech 도메인의 URL 목록을 보여줘
```

## 🔍 문제 해결

### 문제 1: "MCP 서버에 연결할 수 없습니다"

**해결 방법:**

1. 바이너리 경로 확인:
   ```bash
   ls -la /Users/junwoobang/mcp/url-db/bin/url-db
   ```

2. 실행 권한 확인:
   ```bash
   chmod +x /Users/junwoobang/mcp/url-db/bin/url-db
   ```

3. 직접 실행 테스트:
   ```bash
   /Users/junwoobang/mcp/url-db/bin/url-db -mcp-mode=stdio
   ```

### 문제 2: "server not initialized" 오류

**해결 방법:**

MCP 프로토콜은 초기화 순서가 중요합니다. Claude Desktop은 자동으로 이를 처리하지만, 수동 테스트 시:

```bash
# 테스트용 초기화 시퀀스
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"clientInfo":{"name":"test","version":"1.0"}}}
{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | /Users/junwoobang/mcp/url-db/bin/url-db -mcp-mode=stdio
```

### 문제 3: 데이터베이스 접근 오류

**해결 방법:**

1. 데이터베이스 디렉토리 생성:
   ```bash
   mkdir -p /Users/junwoobang/mcp/url-db
   ```

2. 권한 설정:
   ```bash
   chmod 755 /Users/junwoobang/mcp/url-db
   ```

3. 데이터베이스 파일이 이미 있다면 권한 확인:
   ```bash
   chmod 644 /Users/junwoobang/mcp/url-db/url-db.sqlite
   ```

### 문제 4: JSON 파싱 오류

**일반적인 실수들:**

❌ 잘못된 예시들:
```json
{
  "mcpServers": {
    "url-db": {
      "command": "~/mcp/url-db/bin/url-db",  // ~ 사용 금지
      "args": ["-mcp-mode=stdio"],
      "env": {},  // 쉼표 오류
    }
  }
}
```

✅ 올바른 예시:
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/Users/junwoobang/mcp/url-db/bin/url-db",
      "args": ["-mcp-mode=stdio"],
      "env": {}
    }
  }
}
```

### 문제 5: Claude가 MCP 도구를 인식하지 못함

**해결 방법:**

1. Claude Desktop 완전 재시작 (위 3단계 참조)
2. 새 대화 시작 (기존 대화는 MCP 연결이 없을 수 있음)
3. 개발자 도구에서 오류 확인

## 📊 로그 확인

### Claude Desktop 로그 위치

```bash
# macOS
~/Library/Logs/Claude/
```

### URL-DB 로그 활성화

디버그 모드로 실행:
```json
{
  "mcpServers": {
    "url-db": {
      "command": "/Users/junwoobang/mcp/url-db/bin/url-db",
      "args": [
        "-mcp-mode=stdio",
        "-debug=true"
      ],
      "env": {
        "LOG_LEVEL": "debug"
      }
    }
  }
}
```

## 🎯 설정 검증 스크립트

다음 스크립트로 설정을 검증할 수 있습니다:

```bash
#!/bin/bash
# save as check-mcp-setup.sh

echo "🔍 URL-DB MCP 설정 검증 중..."

# 1. 바이너리 확인
if [ -f "/Users/junwoobang/mcp/url-db/bin/url-db" ]; then
    echo "✅ 바이너리 존재함"
else
    echo "❌ 바이너리를 찾을 수 없습니다"
    exit 1
fi

# 2. 실행 권한 확인
if [ -x "/Users/junwoobang/mcp/url-db/bin/url-db" ]; then
    echo "✅ 실행 권한 있음"
else
    echo "❌ 실행 권한 없음"
    exit 1
fi

# 3. MCP 모드 테스트
echo "🧪 MCP 모드 테스트 중..."
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | \
    /Users/junwoobang/mcp/url-db/bin/url-db -mcp-mode=stdio 2>/dev/null | \
    grep -q "result" && echo "✅ MCP 모드 정상 작동" || echo "❌ MCP 모드 오류"

# 4. Claude 설정 파일 확인
CONFIG_FILE="$HOME/Library/Application Support/Claude/claude_desktop_config.json"
if [ -f "$CONFIG_FILE" ]; then
    echo "✅ Claude 설정 파일 존재"
    if grep -q "url-db" "$CONFIG_FILE"; then
        echo "✅ url-db 설정 발견"
    else
        echo "⚠️  url-db 설정이 없습니다"
    fi
else
    echo "❌ Claude 설정 파일을 찾을 수 없습니다"
fi

echo "✨ 검증 완료!"
```

## 💡 추가 팁

1. **데이터베이스 백업**: 정기적으로 SQLite 파일을 백업하세요
   ```bash
   cp url-db.sqlite url-db.sqlite.backup
   ```

2. **다중 환경**: 개발/프로덕션 환경을 분리하려면 다른 데이터베이스 경로를 사용하세요

3. **성능**: 대량의 데이터가 있다면 인덱스를 확인하세요

## 🆘 추가 지원

문제가 지속되면 다음 정보와 함께 이슈를 생성해주세요:

1. `claude_desktop_config.json` 내용 (민감 정보 제외)
2. 오류 메시지 스크린샷
3. 다음 명령어 출력:
   ```bash
   /Users/junwoobang/mcp/url-db/bin/url-db -version
   ```

## 🎉 성공 확인

다음과 같은 응답을 받으면 성공입니다:

```
Claude: URL-DB MCP 서버에 성공적으로 연결되었습니다. 사용 가능한 도구는 다음과 같습니다:
- list_domains: 모든 도메인 조회
- create_domain: 새 도메인 생성
- list_nodes: URL 목록 조회
... (18개 도구)
```