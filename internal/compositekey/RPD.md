# 합성키 시스템 모듈 RPD

## 참조 문서
- [docs/spec/composite-key-conventions.md](../../docs/spec/composite-key-conventions.md)
- [docs/api/06-mcp-api.md](../../docs/api/06-mcp-api.md)

## 요구사항 분석

### 기능 요구사항
1. **합성키 생성**: `tool_name:domain_name:id` 형식으로 생성
2. **합성키 파싱**: 합성키를 구성 요소로 분해
3. **합성키 검증**: 형식 및 구성 요소 검증
4. **문자열 정규화**: 도메인명과 도구명 정규화
5. **에러 처리**: 잘못된 합성키 형식 처리

### 비기능 요구사항
- 스레드 안전성
- 고성능 파싱
- 메모리 효율성
- 확장 가능성

## 데이터 모델

### CompositeKey 구조체
```go
type CompositeKey struct {
    ToolName   string `json:"tool_name"`
    DomainName string `json:"domain_name"`
    ID         int    `json:"id"`
}
```

### 합성키 형식
```
{tool_name}:{domain_name}:{id}
```

## 아키텍처 설계

### 계층 구조
```
Service -> Utilities
```

### 주요 컴포넌트
1. **CompositeKeyService**: 합성키 생성/파싱 서비스
2. **Normalizer**: 문자열 정규화 유틸리티
3. **Validator**: 합성키 검증 로직

## 구현 계획

### Phase 1: Core Functions
- [ ] CompositeKey 구조체 정의
- [ ] 합성키 생성 함수
- [ ] 합성키 파싱 함수
- [ ] 단위 테스트 작성

### Phase 2: Validation & Normalization
- [ ] 문자열 정규화 함수
- [ ] 합성키 검증 함수
- [ ] 에러 처리 로직
- [ ] 단위 테스트 작성

### Phase 3: Service Layer
- [ ] CompositeKeyService 구현
- [ ] 의존성 주입 지원
- [ ] 통합 테스트 작성

### Phase 4: Performance Optimization
- [ ] 성능 최적화
- [ ] 메모리 풀링
- [ ] 벤치마크 테스트

## 형식 규칙

### 구분자
- 콜론(`:`)을 사용하여 구성 요소 구분

### 문자 제한
- 영문자, 숫자, 하이픈(`-`), 언더스코어(`_`)만 허용
- 공백, 특수문자는 하이픈으로 변환

### 대소문자
- 모든 구성 요소는 소문자로 변환

### 길이 제한
- tool_name: 최대 50자
- domain_name: 최대 50자
- id: 최대 20자 (정수)

## 정규화 규칙

### 문자 변환
```go
// 예시
"My Domain Name" -> "my-domain-name"
"TechArticles" -> "tech-articles"
"special@chars!" -> "special-chars"
```

### 연속 구분자 처리
```go
// 예시
"my--domain" -> "my-domain"
"tech___articles" -> "tech-articles"
```

### 앞뒤 공백 제거
```go
// 예시
"  my domain  " -> "my-domain"
```

## 에러 처리

### 합성키 에러 코드
- `COMPOSITE_KEY_INVALID_FORMAT`: 잘못된 형식
- `COMPOSITE_KEY_INVALID_TOOL_NAME`: 도구명 오류
- `COMPOSITE_KEY_INVALID_DOMAIN_NAME`: 도메인명 오류
- `COMPOSITE_KEY_INVALID_ID`: ID 오류

### 검증 규칙
- 정확히 3개의 구성 요소
- 각 구성 요소 길이 제한
- ID는 양의 정수
- 허용된 문자만 사용

## 테스트 전략

### 단위 테스트
- 합성키 생성 테스트
- 합성키 파싱 테스트
- 문자열 정규화 테스트
- 검증 로직 테스트
- 에러 처리 테스트

### 성능 테스트
- 파싱 성능 벤치마크
- 생성 성능 벤치마크
- 메모리 사용량 테스트

### 통합 테스트
- 다른 모듈과의 통합 테스트
- MCP 서비스와의 통합 테스트

## 파일 구조
```
internal/compositekey/
├── RPD.md
├── composite_key.go      # CompositeKey 구조체
├── service.go            # CompositeKeyService
├── service_test.go       # Service 테스트
├── normalizer.go         # 문자열 정규화
├── normalizer_test.go    # 정규화 테스트
├── validator.go          # 검증 로직
├── validator_test.go     # 검증 테스트
├── errors.go             # 에러 정의
├── benchmark_test.go     # 성능 벤치마크
└── examples_test.go      # 사용 예제
```

## 사용 예제

### 합성키 생성
```go
service := compositekey.NewService("url-db")
key, err := service.Create("Tech Articles", 123)
// 결과: "url-db:tech-articles:123"
```

### 합성키 파싱
```go
compositeKey, err := service.Parse("url-db:tech-articles:123")
// 결과: CompositeKey{ToolName: "url-db", DomainName: "tech-articles", ID: 123}
```

### 검증
```go
isValid := service.Validate("url-db:tech-articles:123")
// 결과: true
```

## 성능 고려사항

### 메모리 최적화
- 문자열 풀링
- 재사용 가능한 버퍼
- 가비지 컬렉션 최소화

### 처리 속도
- 정규 표현식 미사용
- 인라인 함수 활용
- 조기 검증

### 확장성
- 플러그인 가능한 정규화 규칙
- 다양한 도구명 지원
- 설정 가능한 제한값

## 의존성
- `strings`: 문자열 처리
- `strconv`: 숫자 변환
- `fmt`: 문자열 포맷팅
- `github.com/stretchr/testify`: 테스트 유틸리티

## 보안 고려사항

### 입력 검증
- 악의적인 입력 필터링
- 길이 제한 강제
- 특수 문자 처리

### 정보 노출 방지
- 내부 ID 숨김
- 도메인 정보 보호
- 에러 메시지 최소화