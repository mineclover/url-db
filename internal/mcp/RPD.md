# MCP 서버 모듈 RPD

## 참조 문서
- [docs/api/06-mcp-api.md](../../docs/api/06-mcp-api.md)
- [docs/spec/composite-key-conventions.md](../../docs/spec/composite-key-conventions.md)

## 요구사항 분석

### 기능 요구사항
1. **MCP 노드 생성**: 합성키 기반 노드 생성
2. **MCP 노드 조회**: 합성키로 노드 조회
3. **MCP 노드 목록**: 도메인별 필터링, 페이지네이션
4. **MCP 노드 수정**: 합성키 기반 노드 수정
5. **MCP 노드 삭제**: 합성키 기반 노드 삭제
6. **URL 검색**: 도메인명과 URL로 노드 찾기
7. **배치 조회**: 여러 합성키 동시 조회
8. **도메인 관리**: MCP 방식 도메인 CRUD
9. **속성 관리**: 합성키 기반 속성 관리
10. **서버 메타데이터**: MCP 서버 정보 제공

### 비기능 요구사항
- 합성키 기반 식별자 사용
- 기존 서비스와의 통합
- 에러 처리 및 로깅
- 성능 최적화

## 데이터 모델

### MCPNode 구조체
```go
type MCPNode struct {
    CompositeID string    `json:"composite_id"`
    URL         string    `json:"url"`
    DomainName  string    `json:"domain_name"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### MCP 요청/응답 모델
- `CreateMCPNodeRequest`: domain_name, url, title, description
- `FindMCPNodeRequest`: domain_name, url
- `BatchMCPNodeRequest`: composite_ids 배열
- `BatchMCPNodeResponse`: nodes, not_found 배열

## 아키텍처 설계

### 계층 구조
```
Handler -> MCPService -> BaseServices -> Repository -> Database
```

### 서비스 통합
1. **MCPService**: MCP 전용 비즈니스 로직
2. **NodeService**: 기존 노드 서비스 활용
3. **DomainService**: 기존 도메인 서비스 활용
4. **AttributeService**: 기존 속성 서비스 활용
5. **CompositeKeyService**: 합성키 처리

## 구현 계획

### Phase 1: Core MCP Service
- [ ] MCPService interface 정의
- [ ] 기존 서비스 통합 로직
- [ ] 합성키 변환 로직
- [ ] 단위 테스트 작성

### Phase 2: Handler Layer
- [ ] MCP HTTP 핸들러 구현
- [ ] 라우터 설정
- [ ] 에러 처리
- [ ] 통합 테스트 작성

### Phase 3: Advanced Features
- [ ] 배치 처리 구현
- [ ] 도메인 관리 구현
- [ ] 속성 관리 구현
- [ ] 서버 메타데이터 구현

### Phase 4: Performance & Optimization
- [ ] 배치 처리 최적화
- [ ] 캐싱 구현
- [ ] 성능 테스트

## 데이터 변환

### 내부 모델 → MCP 모델
```go
func (s *MCPService) convertToMCPNode(node *models.Node, domain *models.Domain) *models.MCPNode {
    compositeID := s.compositeKeyService.Create(domain.Name, node.ID)
    return &models.MCPNode{
        CompositeID: compositeID,
        URL:         node.Content,
        DomainName:  domain.Name,
        Title:       node.Title,
        Description: node.Description,
        CreatedAt:   node.CreatedAt,
        UpdatedAt:   node.UpdatedAt,
    }
}
```

### MCP 모델 → 내부 모델
```go
func (s *MCPService) convertFromMCPNode(mcpNode *models.MCPNode) (*models.Node, error) {
    compositeKey, err := s.compositeKeyService.Parse(mcpNode.CompositeID)
    if err != nil {
        return nil, err
    }
    
    return &models.Node{
        ID:          compositeKey.ID,
        Content:     mcpNode.URL,
        Title:       mcpNode.Title,
        Description: mcpNode.Description,
    }, nil
}
```

## 에러 처리

### MCP 특화 에러 코드
- `INVALID_COMPOSITE_KEY`: 잘못된 합성키 형식
- `DOMAIN_NOT_FOUND`: 도메인 존재하지 않음
- `RESOURCE_NOT_FOUND`: 합성키로 리소스 찾을 수 없음
- `ACCESS_DENIED`: 도메인 접근 권한 없음
- `BATCH_PARTIAL_FAILURE`: 배치 처리 일부 실패

### 에러 응답 형식
```json
{
  "error": "INVALID_COMPOSITE_KEY",
  "message": "합성키 형식이 올바르지 않습니다",
  "expected_format": "tool_name:domain_name:id",
  "provided": "invalid-key"
}
```

## 배치 처리

### 배치 조회 최적화
- 도메인별 그룹화
- 병렬 처리
- 부분 실패 허용

### 배치 응답 처리
- 성공한 항목과 실패한 항목 분리
- 에러 정보 포함
- 성능 통계 제공

## 도메인 관리

### MCP 도메인 API
- 도메인 목록 조회 (노드 수 포함)
- 도메인 생성 (이름 정규화)
- 도메인 수정
- 도메인 삭제

### 도메인 정보 확장
```go
type MCPDomain struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    NodeCount   int       `json:"node_count"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## 속성 관리

### MCP 속성 API
- 노드 속성 조회 (합성키 기반)
- 노드 속성 설정 (배치 처리)
- 속성 값 검증

### 속성 정보 변환
```go
type MCPAttribute struct {
    Name  string `json:"name"`
    Type  string `json:"type"`
    Value string `json:"value"`
}
```

## 서버 메타데이터

### 서버 정보 제공
```go
type MCPServerInfo struct {
    Name               string   `json:"name"`
    Version            string   `json:"version"`
    Description        string   `json:"description"`
    Capabilities       []string `json:"capabilities"`
    CompositeKeyFormat string   `json:"composite_key_format"`
}
```

## 테스트 전략

### 단위 테스트
- MCPService 메서드 테스트
- 데이터 변환 로직 테스트
- 에러 처리 테스트
- 합성키 처리 테스트

### 통합 테스트
- End-to-end MCP API 테스트
- 배치 처리 테스트
- 도메인 관리 테스트
- 속성 관리 테스트

### 성능 테스트
- 배치 처리 성능 측정
- 대용량 데이터 처리 테스트
- 동시성 테스트

## 파일 구조
```
internal/mcp/
├── RPD.md
├── service.go          # MCPService 구현
├── service_test.go     # Service 테스트
├── handler.go          # MCP HTTP 핸들러
├── handler_test.go     # Handler 테스트
├── converter.go        # 데이터 변환 로직
├── converter_test.go   # 변환 테스트
├── batch.go            # 배치 처리 로직
├── batch_test.go       # 배치 테스트
├── domain.go           # 도메인 관리 로직
├── domain_test.go      # 도메인 테스트
├── attribute.go        # 속성 관리 로직
├── attribute_test.go   # 속성 테스트
├── metadata.go         # 서버 메타데이터
├── metadata_test.go    # 메타데이터 테스트
└── errors.go           # MCP 에러 정의
```

## 의존성
- `internal/nodes`: 노드 서비스
- `internal/domains`: 도메인 서비스
- `internal/attributes`: 속성 서비스
- `internal/compositekey`: 합성키 서비스
- `internal/models`: 공통 모델
- `github.com/gin-gonic/gin`: HTTP 라우터
- `github.com/stretchr/testify`: 테스트 유틸리티

## 성능 고려사항

### 배치 처리 최적화
- 병렬 처리 활용
- 데이터베이스 연결 풀링
- 메모리 사용량 최적화

### 캐싱 전략
- 도메인 정보 캐싱
- 합성키 매핑 캐싱
- 속성 정보 캐싱

### 확장성
- 수평 확장 지원
- 로드 밸런싱 고려
- 상태 없는 서비스 설계