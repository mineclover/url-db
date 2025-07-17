# 도메인 관리 모듈 RPD

## 참조 문서
- [docs/api/01-domain-api.md](../../docs/api/01-domain-api.md)
- [docs/spec/domain-errors.md](../../docs/spec/domain-errors.md)
- [schema.sql](../../schema.sql) - domains 테이블

## 요구사항 분석

### 기능 요구사항
1. **도메인 생성**: 새로운 도메인 생성 (이름 중복 방지)
2. **도메인 목록 조회**: 페이지네이션 지원
3. **도메인 상세 조회**: ID로 단일 도메인 조회
4. **도메인 수정**: 설명 수정 (이름 수정 불가)
5. **도메인 삭제**: CASCADE 삭제로 관련 데이터 정리

### 비기능 요구사항
- 도메인 이름 유일성 보장
- 트랜잭션 안전성
- 에러 처리 및 로깅
- 성능: 인덱스 활용

## 데이터 모델

### Domain 구조체
```go
type Domain struct {
    ID          int       `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

### 요청/응답 모델
- `CreateDomainRequest`: name(필수), description(선택)
- `UpdateDomainRequest`: description만 수정 가능
- `DomainListResponse`: 페이지네이션 정보 포함

## 아키텍처 설계

### 계층 구조
```
Handler -> Service -> Repository -> Database
```

### 각 계층의 책임
1. **Repository**: 데이터 접근 로직
   - CRUD 연산
   - 트랜잭션 관리
   - SQL 쿼리 실행

2. **Service**: 비즈니스 로직
   - 도메인 검증
   - 에러 처리
   - 트랜잭션 조정

3. **Handler**: HTTP 요청 처리
   - 요청 파싱
   - 응답 생성
   - 상태 코드 관리

## 구현 계획

### Phase 1: Repository Layer
- [ ] DomainRepository interface 정의
- [ ] SQLite 구현체 작성
- [ ] 단위 테스트 작성

### Phase 2: Service Layer
- [ ] DomainService interface 정의
- [ ] 비즈니스 로직 구현
- [ ] 단위 테스트 작성

### Phase 3: Handler Layer
- [ ] HTTP 핸들러 구현
- [ ] 라우터 설정
- [ ] 통합 테스트 작성

## 에러 처리

### 도메인별 에러 코드
- `DOMAIN_ALREADY_EXISTS`: 도메인 이름 중복
- `DOMAIN_NOT_FOUND`: 도메인 존재하지 않음
- `DOMAIN_NAME_INVALID`: 도메인 이름 형식 오류
- `DOMAIN_HAS_DEPENDENCIES`: 종속성 존재로 삭제 불가

### 검증 규칙
- 도메인 이름: 필수, 최대 255자, 고유
- 설명: 선택, 최대 1000자

## 테스트 전략

### 단위 테스트
- Repository 계층: 각 CRUD 메서드 테스트
- Service 계층: 비즈니스 로직 테스트
- Handler 계층: HTTP 요청/응답 테스트

### 통합 테스트
- End-to-end API 테스트
- 데이터베이스 트랜잭션 테스트
- 에러 시나리오 테스트

## 파일 구조
```
internal/domains/
├── RPD.md
├── repository.go       # DomainRepository interface
├── repository_test.go  # Repository 테스트
├── service.go          # DomainService interface
├── service_test.go     # Service 테스트
├── handler.go          # HTTP 핸들러
├── handler_test.go     # Handler 테스트
└── errors.go           # 도메인별 에러 정의
```

## 의존성
- `internal/database`: 데이터베이스 연결
- `internal/models`: 공통 모델 정의
- `github.com/gin-gonic/gin`: HTTP 라우터
- `github.com/mattn/go-sqlite3`: SQLite 드라이버
- `github.com/stretchr/testify`: 테스트 유틸리티