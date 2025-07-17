# URL Database

## 개요
URL 기반 무제한 속성 태깅이 가능한 데이터베이스 시스템입니다. 
MCP (Model Context Protocol) 서버 지원으로 AI 모델과의 통합이 가능합니다.

## 주요 기능

### 🔗 URL 관리
- URL을 원본 형태 그대로 저장
- 도메인별로 URL 구조화 관리
- 중복 URL 방지 (도메인 내 UNIQUE 제약)
- POST 방식 URL 조회로 긴 URL 처리

### 🏷️ 속성 시스템
- 무제한 속성 정의 및 할당
- 6가지 속성 타입 지원: `tag`, `ordered_tag`, `number`, `string`, `markdown`, `image`
- 도메인별 속성 스키마 관리
- 속성 값 검증 및 타입 강제

### 🔑 합성키 시스템
- 외부 시스템과의 데이터 교환용 합성키 지원
- 형식: `tool_name:domain_name:id`
- 내부 ID 숨김으로 보안 강화
- MCP 클라이언트 친화적 식별자

### 🤖 MCP 서버 지원
- AI 모델과의 직접 통합 가능
- 표준 MCP 프로토콜 준수
- 배치 처리 및 메타데이터 API 제공

## 문서 구조

### API 문서 (`api/`)
- [01-domain-api.md](api/01-domain-api.md) - 도메인 관리 API
- [02-attribute-api.md](api/02-attribute-api.md) - 속성 관리 API
- [03-url-api.md](api/03-url-api.md) - 노드 관리 API (기존 API)
- [04-url-attribute-api.md](api/04-url-attribute-api.md) - 노드 속성 값 관리 API
- [05-url-attribute-validation-api.md](api/05-url-attribute-validation-api.md) - 노드 속성 확인 API
- [06-mcp-api.md](api/06-mcp-api.md) - **MCP 서버 API (새로운 기능)**

### 스펙 문서 (`spec/`)
- [error-codes.md](spec/error-codes.md) - 에러 코드 정의
- [composite-key-conventions.md](spec/composite-key-conventions.md) - **합성키 컨벤션 (새로운 기능)**
- [domain-errors.md](spec/domain-errors.md) - 도메인 관련 에러
- [attribute-errors.md](spec/attribute-errors.md) - 속성 관련 에러
- [node-errors.md](spec/node-errors.md) - 노드 관련 에러
- [node-attribute-errors.md](spec/node-attribute-errors.md) - 노드 속성 관련 에러

### 속성 타입 스펙 (`spec/attribute-types/`)
- [tag.md](spec/attribute-types/tag.md) - 일반 태그
- [ordered_tag.md](spec/attribute-types/ordered_tag.md) - 순서 태그
- [number.md](spec/attribute-types/number.md) - 숫자
- [string.md](spec/attribute-types/string.md) - 문자열
- [markdown.md](spec/attribute-types/markdown.md) - 마크다운
- [image.md](spec/attribute-types/image.md) - 이미지

## 데이터베이스 스키마
- [schema.sql](../schema.sql) - SQLite 데이터베이스 스키마

## 시작하기

### 1. 기본 사용법
```bash
# 도메인 생성
POST /api/domains
{
  "name": "tech-articles",
  "description": "기술 관련 아티클"
}

# URL 추가
POST /api/domains/1/urls
{
  "url": "https://example.com/article",
  "title": "Example Article"
}
```

### 2. MCP 서버 사용법
```bash
# MCP 방식으로 노드 생성
POST /api/mcp/nodes
{
  "domain_name": "tech-articles",
  "url": "https://example.com/article",
  "title": "Example Article"
}

# 합성키로 노드 조회
GET /api/mcp/nodes/url-db:tech-articles:123
```

### 3. 속성 관리
```bash
# 속성 정의
POST /api/domains/1/attributes
{
  "name": "category",
  "type": "tag",
  "description": "카테고리 태그"
}

# 속성 값 설정
POST /api/urls/1/attributes
{
  "attribute_id": 1,
  "value": "javascript"
}
```

## 주요 특징

### 🔒 데이터 무결성
- SQL 수준 UNIQUE 제약 조건
- 외래키 관계 및 CASCADE 삭제
- 속성 타입 강제 및 검증

### 🚀 성능 최적화
- 인덱스 기반 빠른 검색
- 배치 처리 지원
- 페이지네이션 내장

### 🔌 확장성
- 도메인별 독립적 관리
- 속성 시스템 유연성
- MCP 프로토콜 호환

### 🛡️ 보안
- 내부 ID 숨김 (합성키 사용)
- 도메인 격리
- 입력 검증 및 타입 체크

## 사용 사례

- **북마크 관리**: URL과 태그, 메모 관리
- **콘텐츠 큐레이션**: 아티클 수집 및 분류
- **연구 자료 관리**: 논문, 자료 링크 체계화
- **AI 모델 통합**: MCP를 통한 자동 콘텐츠 처리