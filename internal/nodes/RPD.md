# 노드 관리 모듈 RPD

## 참조 문서
- [docs/api/03-url-api.md](../../docs/api/03-url-api.md)
- [docs/spec/node-errors.md](../../docs/spec/node-errors.md)
- [schema.sql](../../schema.sql) - nodes 테이블

## 요구사항 분석

### 기능 요구사항
1. **노드 생성**: 도메인별 URL 저장 (중복 방지)
2. **노드 목록 조회**: 도메인별, 페이지네이션, 검색 지원
3. **노드 상세 조회**: ID로 단일 노드 조회
4. **노드 수정**: 제목, 설명 수정 (URL, 도메인 수정 불가)
5. **노드 삭제**: CASCADE 삭제로 관련 속성 값 정리
6. **URL로 노드 찾기**: POST 방식으로 URL 전체 전송하여 조회

### 비기능 요구사항
- 도메인 내 URL 유일성 보장 (UNIQUE 제약)
- URL 원본 형태 보존
- 긴 URL 처리 (최대 2048자)
- 페이지네이션 성능 최적화
- 검색 기능 (제목, 컨텐츠)

## 데이터 모델

### Node 구조체
```go
type Node struct {
    ID          int       `json:"id" db:"id"`
    Content     string    `json:"content" db:"content"`
    DomainID    int       `json:"domain_id" db:"domain_id"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

### 요청/응답 모델
- `CreateNodeRequest`: url(필수), title(선택), description(선택)
- `UpdateNodeRequest`: title, description 수정 가능
- `FindNodeByURLRequest`: url(필수)
- `NodeListResponse`: 페이지네이션 정보 포함

## 아키텍처 설계

### 계층 구조
```
Handler -> Service -> Repository -> Database
```

### 각 계층의 책임
1. **Repository**: 데이터 접근 로직
   - CRUD 연산
   - 도메인별 조회
   - URL 검색
   - 페이지네이션

2. **Service**: 비즈니스 로직
   - URL 검증
   - 도메인 존재 확인
   - 중복 URL 검증
   - 제목 자동 생성

3. **Handler**: HTTP 요청 처리
   - 요청 파싱
   - 응답 생성
   - 상태 코드 관리

## 구현 계획

### Phase 1: Repository Layer
- [ ] NodeRepository interface 정의
- [ ] SQLite 구현체 작성
- [ ] 단위 테스트 작성

### Phase 2: Service Layer
- [ ] NodeService interface 정의
- [ ] URL 검증 로직
- [ ] 제목 자동 생성 로직
- [ ] 단위 테스트 작성

### Phase 3: Handler Layer
- [ ] HTTP 핸들러 구현
- [ ] 라우터 설정
- [ ] 통합 테스트 작성

### Phase 4: Search & Pagination
- [ ] 검색 기능 구현
- [ ] 페이지네이션 최적화
- [ ] 성능 테스트

## 에러 처리

### 노드별 에러 코드
- `NODE_ALREADY_EXISTS`: 도메인 내 URL 중복
- `NODE_NOT_FOUND`: 노드 존재하지 않음
- `NODE_URL_INVALID`: URL 형식 오류
- `NODE_DOMAIN_NOT_FOUND`: 도메인 존재하지 않음
- `NODE_HAS_ATTRIBUTES`: 속성 존재로 삭제 불가

### 검증 규칙
- URL: 필수, 최대 2048자 (형식 검증 없음)
- 제목: 선택, 최대 255자 (비어있으면 컨텐츠에서 자동 생성)
- 설명: 선택, 최대 1000자
- 컨텐츠와 domain_id는 수정 불가

## 특별 기능

### 제목 자동 생성
- 제목이 비어있을 경우 URL에서 자동 생성
- 도메인명과 경로를 조합하여 의미 있는 제목 생성

### URL 검색 최적화
- POST 방식으로 전체 URL 전송
- GET 방식의 URL 길이 제한 회피
- 인덱스 활용한 빠른 검색

### 페이지네이션
- 커서 기반 페이지네이션 지원
- 대용량 데이터 처리 최적화
- 검색과 페이지네이션 조합

## 테스트 전략

### 단위 테스트
- Repository 계층: CRUD 메서드 테스트
- Service 계층: 비즈니스 로직 테스트
- URL 검증 로직 테스트
- Handler 계층: HTTP 요청/응답 테스트

### 통합 테스트
- End-to-end API 테스트
- 긴 URL 처리 테스트
- 페이지네이션 성능 테스트
- 검색 기능 테스트

### 성능 테스트
- 대용량 URL 데이터 처리
- 검색 성능 측정
- 페이지네이션 성능 측정

## 파일 구조
```
internal/nodes/
├── RPD.md
├── repository.go       # NodeRepository interface
├── repository_test.go  # Repository 테스트
├── service.go          # NodeService interface
├── service_test.go     # Service 테스트
├── handler.go          # HTTP 핸들러
├── handler_test.go     # Handler 테스트
├── url_utils.go        # URL 처리 유틸리티
├── url_utils_test.go   # URL 유틸리티 테스트
└── errors.go           # 노드별 에러 정의
```

## 의존성
- `internal/database`: 데이터베이스 연결
- `internal/models`: 공통 모델 정의
- `internal/domains`: 도메인 검증
- `github.com/gin-gonic/gin`: HTTP 라우터
- `github.com/mattn/go-sqlite3`: SQLite 드라이버
- `github.com/stretchr/testify`: 테스트 유틸리티

## 성능 고려사항

### 인덱스 활용
- `idx_nodes_domain`: 도메인별 조회 최적화
- `idx_nodes_content`: URL 검색 최적화

### 메모리 최적화
- 스트리밍 기반 대용량 데이터 처리
- 적절한 페이지 크기 설정

### 캐시 전략
- 자주 조회되는 노드 캐싱
- 검색 결과 캐싱 (선택적)