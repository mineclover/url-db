# 서비스 모듈 RPD

## 참조 문서
- [internal/domains/RPD.md](../domains/RPD.md) - 도메인 서비스
- [internal/nodes/RPD.md](../nodes/RPD.md) - 노드 서비스
- [internal/attributes/RPD.md](../attributes/RPD.md) - 속성 서비스
- [internal/nodeattributes/RPD.md](../nodeattributes/RPD.md) - 노드 속성 서비스
- [internal/compositekey/RPD.md](../compositekey/RPD.md) - 합성키 서비스
- [internal/mcp/RPD.md](../mcp/RPD.md) - MCP 서비스

## 요구사항 분석

### 기능 요구사항
1. **도메인 서비스**: 도메인 비즈니스 로직 및 검증
2. **노드 서비스**: 노드 비즈니스 로직 및 URL 처리
3. **속성 서비스**: 속성 비즈니스 로직 및 타입 검증
4. **노드 속성 서비스**: 노드 속성 비즈니스 로직 및 값 검증
5. **합성키 서비스**: 합성키 생성, 파싱, 검증
6. **MCP 서비스**: MCP 프로토콜 비즈니스 로직
7. **트랜잭션 조정**: 복합 비즈니스 로직 트랜잭션 관리

### 비기능 요구사항
- 인터페이스 기반 설계
- 의존성 주입 지원
- 에러 처리 및 로깅
- 비즈니스 규칙 검증
- 성능 최적화

## 아키텍처 설계

### 계층 구조
```
Handler -> Service -> Repository -> Database
```

### 서비스 인터페이스
1. **DomainService**: 도메인 비즈니스 로직
2. **NodeService**: 노드 비즈니스 로직
3. **AttributeService**: 속성 비즈니스 로직
4. **NodeAttributeService**: 노드 속성 비즈니스 로직
5. **CompositeKeyService**: 합성키 처리
6. **MCPService**: MCP 통합 서비스

## 구현 계획

### Phase 1: Core Services
- [ ] 각 서비스 인터페이스 정의
- [ ] 기본 비즈니스 로직 구현
- [ ] 검증 로직 구현
- [ ] 단위 테스트 작성

### Phase 2: Advanced Logic
- [ ] 복합 비즈니스 로직 구현
- [ ] 트랜잭션 조정 로직
- [ ] 에러 처리 로직
- [ ] 통합 테스트 작성

### Phase 3: Integration Services
- [ ] MCP 서비스 구현
- [ ] 서비스 간 통합
- [ ] 성능 최적화
- [ ] 종합 테스트

### Phase 4: Optimization
- [ ] 캐싱 구현
- [ ] 성능 모니터링
- [ ] 벤치마크 테스트

## 서비스 인터페이스

### DomainService
```go
type DomainService interface {
    CreateDomain(req *models.CreateDomainRequest) (*models.Domain, error)
    GetDomain(id int) (*models.Domain, error)
    GetDomainByName(name string) (*models.Domain, error)
    ListDomains(page, size int) (*models.DomainListResponse, error)
    UpdateDomain(id int, req *models.UpdateDomainRequest) (*models.Domain, error)
    DeleteDomain(id int) error
}
```

### NodeService
```go
type NodeService interface {
    CreateNode(domainID int, req *models.CreateNodeRequest) (*models.Node, error)
    GetNode(id int) (*models.Node, error)
    GetNodeByDomainAndURL(domainID int, url string) (*models.Node, error)
    ListNodesByDomain(domainID int, page, size int, search string) (*models.NodeListResponse, error)
    UpdateNode(id int, req *models.UpdateNodeRequest) (*models.Node, error)
    DeleteNode(id int) error
    FindNodeByURL(domainID int, req *models.FindNodeByURLRequest) (*models.Node, error)
}
```

### AttributeService
```go
type AttributeService interface {
    CreateAttribute(domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error)
    GetAttribute(id int) (*models.Attribute, error)
    ListAttributesByDomain(domainID int) (*models.AttributeListResponse, error)
    UpdateAttribute(id int, req *models.UpdateAttributeRequest) (*models.Attribute, error)
    DeleteAttribute(id int) error
    ValidateAttributeValue(attributeID int, value string) error
}
```

### NodeAttributeService
```go
type NodeAttributeService interface {
    CreateNodeAttribute(nodeID int, req *models.CreateNodeAttributeRequest) (*models.NodeAttribute, error)
    GetNodeAttribute(id int) (*models.NodeAttribute, error)
    ListNodeAttributesByNode(nodeID int) (*models.NodeAttributeListResponse, error)
    UpdateNodeAttribute(id int, req *models.UpdateNodeAttributeRequest) (*models.NodeAttribute, error)
    DeleteNodeAttribute(id int) error
    ValidateNodeAttributeValue(nodeID, attributeID int, value string) error
}
```

### CompositeKeyService
```go
type CompositeKeyService interface {
    Create(domainName string, id int) string
    Parse(compositeKey string) (*models.CompositeKey, error)
    Validate(compositeKey string) error
    GetToolName() string
}
```

### MCPService
```go
type MCPService interface {
    CreateNode(req *models.CreateMCPNodeRequest) (*models.MCPNode, error)
    GetNode(compositeID string) (*models.MCPNode, error)
    ListNodes(domainName string, page, size int, search string) (*models.MCPNodeListResponse, error)
    UpdateNode(compositeID string, req *models.UpdateMCPNodeRequest) (*models.MCPNode, error)
    DeleteNode(compositeID string) error
    FindNodeByURL(req *models.FindMCPNodeRequest) (*models.MCPNode, error)
    BatchGetNodes(req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error)
    
    // 도메인 관리
    ListDomains() (*models.MCPDomainListResponse, error)
    CreateDomain(req *models.CreateMCPDomainRequest) (*models.MCPDomain, error)
    
    // 속성 관리
    GetNodeAttributes(compositeID string) (*models.MCPNodeAttributeResponse, error)
    SetNodeAttributes(compositeID string, req *models.SetMCPNodeAttributesRequest) error
    
    // 서버 정보
    GetServerInfo() (*models.MCPServerInfo, error)
}
```

## 비즈니스 로직 구현

### 도메인 서비스 구현
```go
type domainService struct {
    domainRepo repositories.DomainRepository
    logger     *log.Logger
}

func (s *domainService) CreateDomain(req *models.CreateDomainRequest) (*models.Domain, error) {
    // 입력 검증
    if err := s.validateCreateDomainRequest(req); err != nil {
        return nil, err
    }
    
    // 중복 확인
    exists, err := s.domainRepo.ExistsByName(req.Name)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, NewDomainAlreadyExistsError(req.Name)
    }
    
    // 도메인 생성
    domain := &models.Domain{
        Name:        req.Name,
        Description: req.Description,
    }
    
    if err := s.domainRepo.Create(domain); err != nil {
        return nil, err
    }
    
    s.logger.Printf("Created domain: %s (ID: %d)", domain.Name, domain.ID)
    return domain, nil
}

func (s *domainService) validateCreateDomainRequest(req *models.CreateDomainRequest) error {
    if req.Name == "" {
        return NewValidationError("name", "domain name is required")
    }
    
    if len(req.Name) > 255 {
        return NewValidationError("name", "domain name too long")
    }
    
    if len(req.Description) > 1000 {
        return NewValidationError("description", "description too long")
    }
    
    return nil
}
```

### 노드 서비스 구현
```go
type nodeService struct {
    nodeRepo   repositories.NodeRepository
    domainRepo repositories.DomainRepository
    logger     *log.Logger
}

func (s *nodeService) CreateNode(domainID int, req *models.CreateNodeRequest) (*models.Node, error) {
    // 도메인 존재 확인
    domain, err := s.domainRepo.GetByID(domainID)
    if err != nil {
        return nil, err
    }
    
    // 입력 검증
    if err := s.validateCreateNodeRequest(req); err != nil {
        return nil, err
    }
    
    // 중복 확인
    exists, err := s.nodeRepo.ExistsByDomainAndContent(domainID, req.URL)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, NewNodeAlreadyExistsError(req.URL)
    }
    
    // 제목 자동 생성
    title := req.Title
    if title == "" {
        title = s.generateTitleFromURL(req.URL)
    }
    
    // 노드 생성
    node := &models.Node{
        Content:     req.URL,
        DomainID:    domainID,
        Title:       title,
        Description: req.Description,
    }
    
    if err := s.nodeRepo.Create(node); err != nil {
        return nil, err
    }
    
    s.logger.Printf("Created node: %s in domain %s (ID: %d)", 
        node.Content, domain.Name, node.ID)
    return node, nil
}

func (s *nodeService) generateTitleFromURL(url string) string {
    // URL에서 의미 있는 제목 생성
    // 예: https://example.com/article/123 -> "Article 123"
    // 간단한 구현 예시
    parts := strings.Split(url, "/")
    if len(parts) > 0 {
        return strings.Title(strings.ReplaceAll(parts[len(parts)-1], "-", " "))
    }
    return "Untitled"
}
```

### 합성키 서비스 구현
```go
type compositeKeyService struct {
    toolName string
}

func NewCompositeKeyService(toolName string) CompositeKeyService {
    return &compositeKeyService{
        toolName: toolName,
    }
}

func (s *compositeKeyService) Create(domainName string, id int) string {
    normalizedDomain := s.normalizeName(domainName)
    return fmt.Sprintf("%s:%s:%d", s.toolName, normalizedDomain, id)
}

func (s *compositeKeyService) Parse(compositeKey string) (*models.CompositeKey, error) {
    parts := strings.Split(compositeKey, ":")
    if len(parts) != 3 {
        return nil, NewInvalidCompositeKeyError(compositeKey, "invalid format")
    }
    
    toolName := parts[0]
    domainName := parts[1]
    idStr := parts[2]
    
    if toolName != s.toolName {
        return nil, NewInvalidCompositeKeyError(compositeKey, "invalid tool name")
    }
    
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return nil, NewInvalidCompositeKeyError(compositeKey, "invalid ID")
    }
    
    if id <= 0 {
        return nil, NewInvalidCompositeKeyError(compositeKey, "ID must be positive")
    }
    
    return &models.CompositeKey{
        ToolName:   toolName,
        DomainName: domainName,
        ID:         id,
    }, nil
}

func (s *compositeKeyService) normalizeName(name string) string {
    // 소문자 변환
    normalized := strings.ToLower(name)
    
    // 특수 문자를 하이픈으로 변환
    reg := regexp.MustCompile(`[^a-z0-9\-_]`)
    normalized = reg.ReplaceAllString(normalized, "-")
    
    // 연속된 하이픈 제거
    reg = regexp.MustCompile(`-+`)
    normalized = reg.ReplaceAllString(normalized, "-")
    
    // 앞뒤 하이픈 제거
    normalized = strings.Trim(normalized, "-")
    
    return normalized
}
```

### MCP 서비스 구현
```go
type mcpService struct {
    nodeService         NodeService
    domainService       DomainService
    attributeService    AttributeService
    compositeKeyService CompositeKeyService
    toolName            string
    version             string
}

func (s *mcpService) CreateNode(req *models.CreateMCPNodeRequest) (*models.MCPNode, error) {
    // 도메인 조회/생성
    domain, err := s.domainService.GetDomainByName(req.DomainName)
    if err != nil {
        // 도메인이 없으면 생성
        createDomainReq := &models.CreateDomainRequest{
            Name:        req.DomainName,
            Description: fmt.Sprintf("Auto-created domain for %s", req.DomainName),
        }
        domain, err = s.domainService.CreateDomain(createDomainReq)
        if err != nil {
            return nil, err
        }
    }
    
    // 노드 생성
    createNodeReq := &models.CreateNodeRequest{
        URL:         req.URL,
        Title:       req.Title,
        Description: req.Description,
    }
    
    node, err := s.nodeService.CreateNode(domain.ID, createNodeReq)
    if err != nil {
        return nil, err
    }
    
    // MCP 노드로 변환
    return s.convertToMCPNode(node, domain), nil
}

func (s *mcpService) convertToMCPNode(node *models.Node, domain *models.Domain) *models.MCPNode {
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

## 에러 처리

### 서비스 에러 정의
```go
type ServiceError struct {
    Code    string
    Message string
    Details interface{}
}

func (e *ServiceError) Error() string {
    return e.Message
}

func NewValidationError(field, message string) *ServiceError {
    return &ServiceError{
        Code:    "VALIDATION_ERROR",
        Message: fmt.Sprintf("Validation failed for field '%s': %s", field, message),
        Details: map[string]string{"field": field, "message": message},
    }
}

func NewDomainAlreadyExistsError(name string) *ServiceError {
    return &ServiceError{
        Code:    "DOMAIN_ALREADY_EXISTS",
        Message: fmt.Sprintf("Domain '%s' already exists", name),
        Details: map[string]string{"domain": name},
    }
}

func NewNodeAlreadyExistsError(url string) *ServiceError {
    return &ServiceError{
        Code:    "NODE_ALREADY_EXISTS",
        Message: fmt.Sprintf("Node with URL '%s' already exists", url),
        Details: map[string]string{"url": url},
    }
}
```

## 테스트 전략

### 단위 테스트
```go
func TestDomainService_CreateDomain(t *testing.T) {
    // 모킹 설정
    mockRepo := &MockDomainRepository{}
    service := NewDomainService(mockRepo, log.New(os.Stdout, "", 0))
    
    // 테스트 데이터
    req := &models.CreateDomainRequest{
        Name:        "test-domain",
        Description: "Test description",
    }
    
    // 모킹 동작 설정
    mockRepo.On("ExistsByName", req.Name).Return(false, nil)
    mockRepo.On("Create", mock.AnythingOfType("*models.Domain")).Return(nil)
    
    // 테스트 실행
    domain, err := service.CreateDomain(req)
    
    // 결과 검증
    assert.NoError(t, err)
    assert.Equal(t, req.Name, domain.Name)
    assert.Equal(t, req.Description, domain.Description)
    
    // 모킹 검증
    mockRepo.AssertExpectations(t)
}
```

### 통합 테스트
```go
func TestNodeService_Integration(t *testing.T) {
    // 실제 데이터베이스 설정
    db := setupTestDB(t)
    defer db.Close()
    
    // 리포지토리 생성
    domainRepo := repositories.NewSQLiteDomainRepository(db)
    nodeRepo := repositories.NewSQLiteNodeRepository(db)
    
    // 서비스 생성
    domainService := NewDomainService(domainRepo, log.New(os.Stdout, "", 0))
    nodeService := NewNodeService(nodeRepo, domainRepo, log.New(os.Stdout, "", 0))
    
    // 도메인 생성
    domain, err := domainService.CreateDomain(&models.CreateDomainRequest{
        Name:        "test-domain",
        Description: "Test domain",
    })
    require.NoError(t, err)
    
    // 노드 생성
    node, err := nodeService.CreateNode(domain.ID, &models.CreateNodeRequest{
        URL:         "https://example.com",
        Title:       "Test Node",
        Description: "Test description",
    })
    require.NoError(t, err)
    
    // 검증
    assert.Equal(t, "https://example.com", node.Content)
    assert.Equal(t, domain.ID, node.DomainID)
}
```

## 파일 구조
```
internal/services/
├── RPD.md
├── interfaces.go          # 서비스 인터페이스
├── domain.go              # 도메인 서비스 구현
├── domain_test.go         # 도메인 서비스 테스트
├── node.go                # 노드 서비스 구현
├── node_test.go           # 노드 서비스 테스트
├── attribute.go           # 속성 서비스 구현
├── attribute_test.go      # 속성 서비스 테스트
├── node_attribute.go      # 노드 속성 서비스 구현
├── node_attribute_test.go # 노드 속성 서비스 테스트
├── composite_key.go       # 합성키 서비스 구현
├── composite_key_test.go  # 합성키 서비스 테스트
├── mcp.go                 # MCP 서비스 구현
├── mcp_test.go            # MCP 서비스 테스트
├── errors.go              # 서비스 에러 정의
├── validators.go          # 공통 검증 로직
├── validators_test.go     # 검증 로직 테스트
└── testutils.go           # 테스트 유틸리티
```

## 성능 고려사항

### 캐싱 전략
- 도메인 정보 캐싱
- 자주 조회되는 노드 캐싱
- 속성 정보 캐싱

### 비동기 처리
- 배치 작업 비동기 처리
- 이벤트 기반 아키텍처

### 성능 모니터링
- 서비스 메서드 실행 시간 측정
- 메모리 사용량 모니터링
- 에러율 추적

## 의존성
- `internal/repositories`: 리포지토리 인터페이스
- `internal/models`: 데이터 모델
- `log`: 로깅
- `fmt`: 문자열 포맷팅
- `strings`: 문자열 처리
- `strconv`: 문자열 변환
- `regexp`: 정규 표현식
- `github.com/stretchr/testify`: 테스트 유틸리티

## 보안 고려사항

### 입력 검증
- 모든 사용자 입력 검증
- SQL 인젝션 방지
- XSS 방지

### 권한 관리
- 도메인별 접근 제어
- 리소스 소유권 확인
- 감사 로깅