# URL-DB: AI 어시스턴트용 URL 관리 시스템

URL-DB는 Claude Desktop 등의 AI 어시스턴트가 URL을 효율적으로 저장하고 관리할 수 있게 해주는 MCP(Model Context Protocol) 서버입니다.

## 🎯 무엇을 할 수 있나요?

- **URL 저장 및 분류**: 웹사이트 주소를 도메인별로 체계적으로 관리
- **스마트 태깅**: URL에 태그, 카테고리, 메모 등 다양한 속성 추가
- **빠른 검색**: 저장된 URL을 키워드, 태그, 도메인으로 빠르게 찾기
- **AI 통합**: Claude Desktop에서 자연어로 URL 관리 가능
- **데이터 소유권**: 모든 데이터는 본인의 컴퓨터에 SQLite 파일로 저장

## 🚀 빠른 시작

### 1. Docker로 간단 설치 (권장)

```bash
# Docker 이미지 다운로드 및 실행
docker run -it --rm -v ~/url-db-data:/data asfdassdssa/url-db:latest
```

### 2. Claude Desktop 설정

Claude Desktop 설정 파일에 다음 내용을 추가하세요:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%AppData%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "~/url-db-data:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 3. Claude Desktop 재시작

설정을 저장한 후 Claude Desktop을 재시작하면 URL 관리 기능을 사용할 수 있습니다!

## 💡 사용 예시

Claude Desktop에서 다음과 같이 대화하면 됩니다:

```
👤 "이 URL을 개발 자료로 저장해줘: https://github.com/microsoft/vscode"

🤖 GitHub 개발 자료로 저장했습니다!
   - 도메인: github
   - 태그: development, editor, microsoft
   - 저장 위치: ~/url-db-data/url-db.sqlite

👤 "개발 관련 URL들 찾아줘"

🤖 개발 관련 URL 5개를 찾았습니다:
   1. https://github.com/microsoft/vscode (Visual Studio Code)
   2. https://nodejs.org (Node.js 공식사이트)
   ...
```

## 🗂️ 데이터베이스 위치 및 설정

### 기본 설정

위의 기본 설정을 사용하면 SQLite 데이터베이스가 다음 위치에 저장됩니다:

- **macOS/Linux**: `~/url-db-data/url-db.sqlite`
- **Windows**: `%UserProfile%\url-db-data\url-db.sqlite`

### 사용자 정의 위치

다른 위치에 저장하고 싶다면 설정에서 경로를 변경하세요:

```json
{
  "mcpServers": {
    "url-db": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "/your/custom/path:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 직접 데이터베이스 접근

SQLite 파일에 직접 접근할 수도 있습니다:

```bash
# 데이터베이스 내용 확인
sqlite3 ~/url-db-data/url-db.sqlite "SELECT * FROM domains;"

# 모든 URL 조회  
sqlite3 ~/url-db-data/url-db.sqlite "SELECT url, title FROM nodes LIMIT 10;"
```

## 🛠️ 고급 설정 옵션

### 1. 프로젝트별 데이터베이스

여러 프로젝트를 위해 별도의 데이터베이스를 사용할 수 있습니다:

```json
{
  "mcpServers": {
    "url-db-work": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-v", "~/work-urls:/data",
        "asfdassdssa/url-db:latest"
      ]
    },
    "url-db-personal": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm", 
        "-v", "~/personal-urls:/data",
        "asfdassdssa/url-db:latest"
      ]
    }
  }
}
```

### 2. 로컬 빌드 (Docker 없이)

Docker를 사용하지 않고 직접 빌드하려면:

```bash
# 소스코드 다운로드
git clone https://github.com/mineclover/url-db.git
cd url-db

# 빌드 및 실행
make build
./bin/url-db -mcp-mode=stdio
```

Claude Desktop 설정:
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

## 🔧 주요 기능

### 도메인 시스템
- URL을 도메인별로 자동 분류 (github.com, stackoverflow.com 등)
- 도메인별 커스텀 속성 정의 가능

### 속성 시스템
- **태그**: 키워드 기반 분류
- **카테고리**: 계층적 분류
- **평점**: 5점 척도 평가
- **메모**: 자유 텍스트 설명
- **날짜**: 생성/수정 시간 자동 기록

### 검색 기능
- 제목, URL, 설명으로 검색
- 태그별 필터링
- 도메인별 그룹화
- 날짜 범위 검색

## 🐳 Docker 배포 옵션

### 1. Docker Hub에서 바로 사용
```bash
docker run -it --rm -v url-db-data:/data asfdassdssa/url-db:latest
```

### 2. 여러 서비스 모드로 실행
```bash
# HTTP API 서버 (포트 8080)
docker run -d -p 8080:8080 -v url-db-data:/data asfdassdssa/url-db:latest -port=8080

# 모든 서비스 동시 실행
git clone https://github.com/mineclover/url-db.git
cd url-db
make docker-compose-up
```

### 3. 개발자용 빌드
```bash
git clone https://github.com/mineclover/url-db.git
cd url-db
make docker-build
make docker-run
```

## 🔍 트러블슈팅

### "command not found" 오류
- Docker가 설치되어 있는지 확인
- Docker 데몬이 실행 중인지 확인: `docker ps`

### 데이터가 저장되지 않음
- 볼륨 마운트 경로 확인
- 디렉토리 권한 확인: `ls -la ~/url-db-data/`

### Claude Desktop에서 인식 안됨
- 설정 파일 경로 확인
- JSON 문법 오류 확인 (따옴표, 쉼표 등)
- Claude Desktop 재시작 필요

### 자세한 로그 확인
```bash
# Docker 컨테이너 로그 확인
docker run -it --rm -v url-db-data:/data asfdassdssa/url-db:latest -mcp-mode=stdio
```

## 📚 추가 문서

- [Docker 배포 가이드](docker-deployment.md) - 상세한 Docker 설정 방법
- [Claude Desktop 설정 가이드](docker-hub-deploy.md) - 다양한 설정 예시
- [SQLite 호스트 저장 가이드](sqlite-host-storage-guide.md) - 데이터베이스 관리 방법
- [개발자 가이드](CLAUDE.md) - 코드 기여 및 개발 환경 설정

## 🤝 지원 및 기여

- **버그 리포트**: [GitHub Issues](https://github.com/mineclover/url-db/issues)
- **기능 요청**: [GitHub Discussions](https://github.com/mineclover/url-db/discussions)
- **기여하기**: [Contributing Guide](CONTRIBUTING.md)

## 📄 라이선스

Apache 2.0 License - 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

---

**💡 팁**: Claude Desktop에서 "URL 관리 도구가 있어?" 라고 물어보면 사용 가능한 모든 기능을 확인할 수 있습니다!