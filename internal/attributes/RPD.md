# 속성 관리 모듈 RPD

## 참조 문서
- [docs/api/02-attribute-api.md](../../docs/api/02-attribute-api.md)
- [docs/spec/attribute-errors.md](../../docs/spec/attribute-errors.md)
- [docs/spec/attribute-types/](../../docs/spec/attribute-types/) - 속성 타입별 스펙
- [schema.sql](../../schema.sql) - attributes 테이블

## 요구사항 분석

### 기능 요구사항
1. **속성 생성**: 도메인별 속성 정의 (이름 중복 방지)
2. **속성 목록 조회**: 도메인별 속성 목록
3. **속성 상세 조회**: ID로 단일 속성 조회
4. **속성 수정**: 설명 수정 (이름, 타입 수정 불가)
5. **속성 삭제**: CASCADE 삭제로 관련 속성 값 정리

### 지원하는 속성 타입
- `tag`: 일반 태그
- `ordered_tag`: 순서가 있는 태그
- `number`: 숫자 값
- `string`: 문자열
- `markdown`: 마크다운 텍스트
- `image`: 이미지 URL

### 비기능 요구사항
- 도메인 내 속성 이름 유일성 보장
- 속성 타입 검증
- 트랜잭션 안전성
- 에러 처리 및 로깅

## 데이터 모델

### Attribute 구조체
```go
type Attribute struct {
    ID          int           `json:"id" db:"id"`
    DomainID    int           `json:"domain_id" db:"domain_id"`
    Name        string        `json:"name" db:"name"`
    Type        AttributeType `json:"type" db:"type"`
    Description string        `json:"description" db:"description"`
    CreatedAt   time.Time     `json:"created_at" db:"created_at"`
}
```

### 속성 타입 정의
```go
type AttributeType string

const (
    AttributeTypeTag        AttributeType = "tag"
    AttributeTypeOrderedTag AttributeType = "ordered_tag"
    AttributeTypeNumber     AttributeType = "number"
    AttributeTypeString     AttributeType = "string"
    AttributeTypeMarkdown   AttributeType = "markdown"
    AttributeTypeImage      AttributeType = "image"
)
```

### 요청/응답 모델
- `CreateAttributeRequest`: name(필수), type(필수), description(선택)
- `UpdateAttributeRequest`: description만 수정 가능
- `AttributeListResponse`: 속성 목록

## 아키텍처 설계

### 계층 구조
```
Handler -> Service -> Repository -> Database
```

### 각 계층의 책임
1. **Repository**: 데이터 접근 로직
   - CRUD 연산
   - 도메인별 조회
   - 트랜잭션 관리

2. **Service**: 비즈니스 로직
   - 속성 타입 검증
   - 도메인 존재 확인
   - 중복 이름 검증

3. **Handler**: HTTP 요청 처리
   - 요청 파싱
   - 응답 생성
   - 상태 코드 관리

## 구현 계획

### Phase 1: Repository Layer
- [ ] AttributeRepository interface 정의
- [ ] SQLite 구현체 작성
- [ ] 단위 테스트 작성

### Phase 2: Service Layer
- [ ] AttributeService interface 정의
- [ ] 속성 타입 검증 로직
- [ ] 단위 테스트 작성

### Phase 3: Handler Layer
- [ ] HTTP 핸들러 구현
- [ ] 라우터 설정
- [ ] 통합 테스트 작성

### Phase 4: Type Validators
- [ ] 각 속성 타입별 검증기 구현
- [ ] 속성 값 검증 로직
- [ ] 타입별 테스트 작성

## 에러 처리

### 속성별 에러 코드
- `ATTRIBUTE_ALREADY_EXISTS`: 속성 이름 중복
- `ATTRIBUTE_NOT_FOUND`: 속성 존재하지 않음
- `ATTRIBUTE_TYPE_INVALID`: 지원하지 않는 속성 타입
- `ATTRIBUTE_HAS_VALUES`: 속성 값 존재로 삭제 불가
- `DOMAIN_NOT_FOUND`: 도메인 존재하지 않음

### 검증 규칙
- 속성 이름: 필수, 최대 255자, 도메인 내 고유
- 속성 타입: 필수, 지원하는 타입만 허용
- 설명: 선택, 최대 1000자

## 속성 타입별 검증

### tag (일반 태그)
- 값: 문자열, 최대 255자
- 특별한 검증 없음

### ordered_tag (순서 태그)
- 값: 문자열, 최대 255자
- order_index 필수

### number (숫자)
- 값: 숫자 형태 문자열
- 범위 검증 (선택적)

### string (문자열)
- 값: 문자열, 최대 2048자
- 특별한 검증 없음

### markdown (마크다운)
- 값: 마크다운 텍스트
- 기본 마크다운 문법 검증

### image (이미지)
- 값: 이미지 URL
- URL 형식 검증

## 테스트 전략

### 단위 테스트
- Repository 계층: CRUD 메서드 테스트
- Service 계층: 비즈니스 로직 테스트
- 속성 타입 검증기 테스트
- Handler 계층: HTTP 요청/응답 테스트

### 통합 테스트
- End-to-end API 테스트
- 속성 타입별 시나리오 테스트
- 에러 시나리오 테스트

## 파일 구조
```
internal/attributes/
├── RPD.md
├── repository.go       # AttributeRepository interface
├── repository_test.go  # Repository 테스트
├── service.go          # AttributeService interface
├── service_test.go     # Service 테스트
├── handler.go          # HTTP 핸들러
├── handler_test.go     # Handler 테스트
├── types.go            # 속성 타입 정의
├── validators.go       # 속성 타입별 검증기
├── validators_test.go  # 검증기 테스트
└── errors.go           # 속성별 에러 정의
```

## 의존성
- `internal/database`: 데이터베이스 연결
- `internal/models`: 공통 모델 정의
- `internal/domains`: 도메인 검증
- `github.com/gin-gonic/gin`: HTTP 라우터
- `github.com/mattn/go-sqlite3`: SQLite 드라이버
- `github.com/stretchr/testify`: 테스트 유틸리티