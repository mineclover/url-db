# URL-DB: AI 어시스턴트용 URL 관리 시스템

URL-DB는 HTTP 클라이언트와 AI 어시스턴트가 URL을 효율적으로 저장하고 관리할 수 있게 해주는 MCP(Model Context Protocol) 서버입니다.

## 🎯 무엇을 할 수 있나요?

- **URL 저장 및 분류**: 웹사이트 주소를 도메인별로 체계적으로 관리
- **스마트 태깅**: URL에 태그, 카테고리, 메모 등 다양한 속성 추가
- **빠른 검색**: 저장된 URL을 키워드, 태그, 도메인으로 빠르게 찾기
- **AI 통합**: HTTP 클라이언트와 SSE 연결로 다양한 환경에서 사용 가능
- **폴더별 데이터 관리**: Docker 볼륨으로 프로젝트별 데이터베이스 자동 분리

## 🚀 빠른 시작 (SSE 서버 모드)

### 1. SSE 서버 실행

```bash
# SSE 모드로 URL-DB 서버 시작 (포트 8080)
docker run -d -p 8080:8080 -v ~/url-db:/data --name url-db-server asfdassdssa/url-db:latest -mcp-mode=sse

# 서버 상태 확인
curl http://localhost:8080/health
```

### 2. HTTP 클라이언트 연결 설정

SSE 모드를 지원하는 HTTP 클라이언트에서 다음 엔드포인트로 연결:

```
서버 주소: http://localhost:8080/mcp
Health Check: http://localhost:8080/health  
데이터베이스: ~/url-db/url-db.sqlite (자동 생성)
```

### 3. 서버 상태 확인

```bash
# 서버가 정상 동작하는지 확인
curl http://localhost:8080/health
```

## 💡 사용 방법

SSE 서버가 실행되면 MCP 프로토콜을 지원하는 클라이언트에서 다음 엔드포인트로 연결하여 URL 관리 기능을 사용할 수 있습니다:

- **MCP 엔드포인트**: `http://localhost:8080/mcp`  
- **Health Check**: `http://localhost:8080/health`
- **사용 가능한 도구**: 18개의 MCP 도구 (아래 참조)

## 🗂️ 폴더별 데이터베이스 관리

### 자동 폴더 기반 분리

Docker 볼륨을 활용하여 프로젝트별로 데이터베이스를 자동 분리할 수 있습니다:

```bash
# 프로젝트별 서버 실행 - 각각 독립된 데이터베이스 생성
docker run -d -p 8080:8080 -v ~/work-project:/data --name work-db asfdassdssa/url-db:latest -mcp-mode=sse
docker run -d -p 8081:8081 -v ~/personal-urls:/data --name personal-db asfdassdssa/url-db:latest -mcp-mode=sse -port=8081
docker run -d -p 8082:8082 -v ~/research-links:/data --name research-db asfdassdssa/url-db:latest -mcp-mode=sse -port=8082
```

### 데이터베이스 위치 확인

각 폴더에 자동으로 `url-db.sqlite` 파일이 생성됩니다:

```bash
# 각 프로젝트의 데이터베이스 확인
ls -la ~/work-project/url-db.sqlite      # 작업 프로젝트 DB
ls -la ~/personal-urls/url-db.sqlite     # 개인 URL DB  
ls -la ~/research-links/url-db.sqlite    # 연구 자료 DB
```

### 직접 데이터베이스 접근

SQLite 파일에 직접 접근하여 데이터 확인:

```bash
# 작업 프로젝트 데이터베이스 조회
sqlite3 ~/work-project/url-db.sqlite "SELECT * FROM domains;"

# 개인 URL 데이터베이스 조회
sqlite3 ~/personal-urls/url-db.sqlite "SELECT url, title FROM nodes LIMIT 10;"
```

## 🛠️ 다중 서버 운영

### 포트별 서버 관리

여러 프로젝트를 동시에 운영할 수 있습니다:

```bash
# 서버 상태 확인
curl http://localhost:8080/health  # 작업 프로젝트
curl http://localhost:8081/health  # 개인 URL
curl http://localhost:8082/health  # 연구 자료

# 서버 중지
docker stop work-db personal-db research-db

# 서버 재시작
docker start work-db personal-db research-db
```

### 로컬 빌드 (개발자용)

Docker 없이 직접 빌드하여 사용:

```bash
# 소스코드 다운로드
git clone https://github.com/mineclover/url-db.git
cd url-db

# SSE 모드로 빌드 및 실행
make build
./bin/url-db -mcp-mode=sse -port=8080

# 서버 테스트
curl http://localhost:8080/health
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

## 🐳 실제 사용 시나리오

### 개발팀 협업 시나리오

```bash
# 팀별 URL 데이터베이스 서버 구축
docker run -d -p 8080:8080 -v ~/team-frontend:/data --name frontend-urls asfdassdssa/url-db:latest -mcp-mode=sse
docker run -d -p 8081:8081 -v ~/team-backend:/data --name backend-urls asfdassdssa/url-db:latest -mcp-mode=sse -port=8081
docker run -d -p 8082:8082 -v ~/team-design:/data --name design-urls asfdassdssa/url-db:latest -mcp-mode=sse -port=8082

# 각 팀에서 MCP 클라이언트로 접근
# Frontend팀: http://localhost:8080/mcp
# Backend팀: http://localhost:8081/mcp  
# Design팀: http://localhost:8082/mcp
```

### 연구자/학습자 시나리오

```bash
# 주제별 연구 자료 서버
docker run -d -p 8080:8080 -v ~/ai-research:/data --name ai-papers asfdassdssa/url-db:latest -mcp-mode=sse
docker run -d -p 8081:8081 -v ~/web-dev-learning:/data --name webdev-resources asfdassdssa/url-db:latest -mcp-mode=sse -port=8081

# MCP 클라이언트로 자료 정리
# AI 논문 및 자료: http://localhost:8080/mcp
# 웹 개발 학습 자료: http://localhost:8081/mcp
```

### Docker Compose로 한번에 관리

```bash
# 프로젝트 클론 후 전체 서비스 실행
git clone https://github.com/mineclover/url-db.git
cd url-db
make docker-compose-up

# 접근 포인트:
# http://localhost:8080 - HTTP API
# http://localhost:8081 - SSE MCP 
# http://localhost:8082 - HTTP MCP
```

## 🔍 트러블슈팅

### 서버 연결 문제
```bash
# 서버 상태 확인
docker ps | grep url-db
curl http://localhost:8080/health

# 로그 확인  
docker logs url-db-server
```

### 데이터가 저장되지 않음
```bash
# 볼륨 마운트 확인
docker inspect url-db-server | grep Mounts -A 10

# 데이터베이스 파일 확인
ls -la ~/url-db/url-db.sqlite
```

### 포트 충돌 문제
```bash
# 사용 중인 포트 확인
lsof -i :8080

# 다른 포트로 실행
docker run -d -p 8083:8083 -v ~/url-db:/data --name url-db-alt asfdassdssa/url-db:latest -mcp-mode=sse -port=8083
```

### 서버 재시작
```bash
# 컨테이너 재시작
docker restart url-db-server

# 완전 재생성
docker stop url-db-server
docker rm url-db-server
docker run -d -p 8080:8080 -v ~/url-db:/data --name url-db-server asfdassdssa/url-db:latest -mcp-mode=sse
```

## 📚 추가 문서

- [Docker 배포 가이드](docker-deployment.md) - 상세한 Docker 설정 방법
- [HTTP 클라이언트 연동 가이드](docker-hub-deploy.md) - 다양한 설정 예시
- [SQLite 호스트 저장 가이드](sqlite-host-storage-guide.md) - 데이터베이스 관리 방법
- [개발자 가이드](CLAUDE.md) - 코드 기여 및 개발 환경 설정

## 🤝 지원 및 기여

- **버그 리포트**: [GitHub Issues](https://github.com/mineclover/url-db/issues)
- **기능 요청**: [GitHub Discussions](https://github.com/mineclover/url-db/discussions)
- **기여하기**: [Contributing Guide](CONTRIBUTING.md)

## 📄 라이선스

Apache 2.0 License - 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 🛠️ MCP 도구 목록

URL-DB는 다음과 같은 MCP 도구들을 제공합니다:

### 도메인 관리
- **get_server_info**: Get server information
- **list_domains**: Get all domains
- **create_domain**: Create new domain for organizing URLs

### URL(노드) 관리
- **list_nodes**: List URLs in domain
- **create_node**: Add URL to domain
- **get_node**: Get URL details
- **update_node**: Update URL title or description
- **delete_node**: Remove URL
- **find_node_by_url**: Search by exact URL
- **scan_all_content**: Retrieve all URLs and their content from a domain using page-based navigation with token optimization for AI processing

### 속성 관리
- **get_node_attributes**: Get URL tags and attributes
- **set_node_attributes**: Add or update URL tags
- **list_domain_attributes**: Get available tag types for domain
- **create_domain_attribute**: Define new tag type for domain
- **get_domain_attribute**: Get details of a specific domain attribute
- **update_domain_attribute**: Update domain attribute description
- **delete_domain_attribute**: Remove domain attribute definition
- **filter_nodes_by_attributes**: Filter nodes by attribute values
- **get_node_with_attributes**: Get URL details with all attributes

### 의존성 관리
- **create_dependency**: Create dependency relationship between nodes
- **list_node_dependencies**: List what a node depends on
- **list_node_dependents**: List what depends on a node
- **delete_dependency**: Remove dependency relationship

### 템플릿 관리
- **list_templates**: List templates in domain
- **create_template**: Create new template in domain
- **get_template**: Get template details
- **update_template**: Update template
- **delete_template**: Delete template
- **clone_template**: Clone existing template
- **generate_template_scaffold**: Generate template scaffold for given type
- **validate_template**: Validate template data structure

---

**💡 팁**: MCP 클라이언트에서 18개의 도구를 통해 URL을 체계적으로 관리할 수 있습니다. 서버 상태는 health 엔드포인트로 확인하세요!