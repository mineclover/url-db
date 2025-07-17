# 노드 속성 관리 모듈 RPD

## 참조 문서
- [docs/api/04-url-attribute-api.md](../../docs/api/04-url-attribute-api.md)
- [docs/api/05-url-attribute-validation-api.md](../../docs/api/05-url-attribute-validation-api.md)
- [docs/spec/node-attribute-errors.md](../../docs/spec/node-attribute-errors.md)
- [docs/spec/attribute-types/](../../docs/spec/attribute-types/) - 속성 타입별 검증
- [schema.sql](../../schema.sql) - node_attributes 테이블

## 요구사항 분석

### 기능 요구사항
1. **노드 속성 생성**: 노드에 속성 값 할당
2. **노드 속성 목록 조회**: 노드별 속성 값 목록
3. **노드 속성 상세 조회**: ID로 단일 속성 값 조회
4. **노드 속성 수정**: 속성 값 및 순서 수정
5. **노드 속성 삭제**: 특정 속성 값 삭제
6. **속성 값 검증**: 속성 타입별 값 검증
7. **순서 관리**: ordered_tag 타입의 순서 관리

### 비기능 요구사항
- 속성 타입별 값 검증
- 순서 태그 순서 관리
- 트랜잭션 안전성
- 성능 최적화 (인덱스 활용)

## 데이터 모델

### NodeAttribute 구조체
```go
type NodeAttribute struct {
    ID          int       `json:"id" db:"id"`
    NodeID      int       `json:"node_id" db:"node_id"`
    AttributeID int       `json:"attribute_id" db:"attribute_id"`
    Value       string    `json:"value" db:"value"`
    OrderIndex  *int      `json:"order_index" db:"order_index"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
```

### 조인된 정보 포함 모델
```go
type NodeAttributeWithInfo struct {
    ID          int           `json:"id"`
    NodeID      int           `json:"node_id"`
    AttributeID int           `json:"attribute_id"`
    Name        string        `json:"name"`
    Type        AttributeType `json:"type"`
    Value       string        `json:"value"`
    OrderIndex  *int          `json:"order_index"`
    CreatedAt   time.Time     `json:"created_at"`
}
```

### 요청/응답 모델
- `CreateNodeAttributeRequest`: attribute_id, value, order_index(선택)
- `UpdateNodeAttributeRequest`: value, order_index(선택)
- `NodeAttributeListResponse`: 속성 정보 포함 목록

## 아키텍처 설계

### 계층 구조
```
Handler -> Service -> Repository -> Database
```

### 각 계층의 책임
1. **Repository**: 데이터 접근 로직
   - CRUD 연산
   - 노드별 속성 조회
   - 순서 관리
   - 조인 쿼리

2. **Service**: 비즈니스 로직
   - 속성 값 검증
   - 순서 관리 로직
   - 속성 타입별 처리
   - 노드/속성 존재 확인

3. **Handler**: HTTP 요청 처리
   - 요청 파싱
   - 응답 생성
   - 상태 코드 관리

## 구현 계획

### Phase 1: Repository Layer
- [ ] NodeAttributeRepository interface 정의
- [ ] SQLite 구현체 작성
- [ ] 조인 쿼리 구현
- [ ] 단위 테스트 작성

### Phase 2: Service Layer
- [ ] NodeAttributeService interface 정의
- [ ] 속성 타입별 검증 로직
- [ ] 순서 관리 로직
- [ ] 단위 테스트 작성

### Phase 3: Handler Layer
- [ ] HTTP 핸들러 구현
- [ ] 라우터 설정
- [ ] 통합 테스트 작성

### Phase 4: Validation System
- [ ] 속성 타입별 검증기 구현
- [ ] 검증 결과 응답
- [ ] 검증 테스트

## 에러 처리

### 노드 속성별 에러 코드
- `NODE_ATTRIBUTE_NOT_FOUND`: 노드 속성 존재하지 않음
- `NODE_NOT_FOUND`: 노드 존재하지 않음
- `ATTRIBUTE_NOT_FOUND`: 속성 존재하지 않음
- `NODE_ATTRIBUTE_VALUE_INVALID`: 속성 값 형식 오류
- `NODE_ATTRIBUTE_ORDER_INVALID`: 순서 인덱스 오류
- `NODE_ATTRIBUTE_DOMAIN_MISMATCH`: 노드와 속성의 도메인 불일치

### 검증 규칙
- 속성 값: 필수, 최대 2048자, 타입별 형식 검증
- 순서 인덱스: 선택, ordered_tag 타입에서만 사용
- 노드와 속성의 도메인 일치 확인

## 속성 타입별 검증

### tag (일반 태그)
- 값: 문자열, 최대 255자
- 순서 인덱스: 무시

### ordered_tag (순서 태그)
- 값: 문자열, 최대 255자
- 순서 인덱스: 필수, 양의 정수

### number (숫자)
- 값: 숫자 형태 문자열
- 정수/실수 검증
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
- 이미지 파일 확장자 검증

## 순서 관리

### 순서 자동 관리
- 순서 인덱스 미제공 시 자동 할당
- 기존 순서 재정렬
- 중복 순서 방지

### 순서 변경
- 순서 인덱스 수정 시 자동 재정렬
- 트랜잭션 내에서 순서 보장

## 검증 API

### 속성 값 검증
- 실시간 검증 API 제공
- 검증 결과 상세 정보 반환
- 타입별 검증 규칙 적용

### 배치 검증
- 여러 속성 값 동시 검증
- 검증 결과 요약 제공

## 테스트 전략

### 단위 테스트
- Repository 계층: CRUD 메서드 테스트
- Service 계층: 비즈니스 로직 테스트
- 속성 타입별 검증 테스트
- 순서 관리 로직 테스트
- Handler 계층: HTTP 요청/응답 테스트

### 통합 테스트
- End-to-end API 테스트
- 속성 타입별 시나리오 테스트
- 순서 관리 시나리오 테스트
- 검증 API 테스트

### 성능 테스트
- 대용량 속성 데이터 처리
- 조인 쿼리 성능 측정
- 검증 성능 측정

## 파일 구조
```
internal/nodeattributes/
├── RPD.md
├── repository.go       # NodeAttributeRepository interface
├── repository_test.go  # Repository 테스트
├── service.go          # NodeAttributeService interface
├── service_test.go     # Service 테스트
├── handler.go          # HTTP 핸들러
├── handler_test.go     # Handler 테스트
├── validators.go       # 속성 타입별 검증기
├── validators_test.go  # 검증기 테스트
├── order_manager.go    # 순서 관리 로직
├── order_manager_test.go # 순서 관리 테스트
└── errors.go           # 노드 속성별 에러 정의
```

## 의존성
- `internal/database`: 데이터베이스 연결
- `internal/models`: 공통 모델 정의
- `internal/nodes`: 노드 검증
- `internal/attributes`: 속성 검증
- `github.com/gin-gonic/gin`: HTTP 라우터
- `github.com/mattn/go-sqlite3`: SQLite 드라이버
- `github.com/stretchr/testify`: 테스트 유틸리티

## 성능 고려사항

### 인덱스 활용
- `idx_node_attributes_node`: 노드별 속성 조회 최적화
- `idx_node_attributes_attribute`: 속성별 값 조회 최적화

### 쿼리 최적화
- 조인 쿼리 최적화
- 속성 정보 포함 조회 최적화

### 메모리 관리
- 대용량 속성 데이터 스트리밍
- 적절한 페이지 크기 설정