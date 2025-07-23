# 테스트 친화적 설계 실제 적용 예시

## 개요

이 문서는 URL-DB 프로젝트에서 테스트 친화적 설계 원칙을 실제로 적용하는 구체적인 예시를 보여줍니다.

## 1. 도메인 엔티티 설계

### Before: 테스트하기 어려운 설계

```go
// internal/models/domain.go
type Domain struct {
    ID          int       `db:"id"`
    Name        string    `db:"name"`
    Description string    `db:"description"`
    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
}

func (d *Domain) UpdateDescription(desc string) {
    d.Description = desc
    d.UpdatedAt = time.Now() // 부수 효과
}
```

### After: 테스트하기 쉬운 설계

```go
// internal/domain/entity/domain.go
type Domain struct {
    id          int
    name        string
    description string
    createdAt   time.Time
}

// 생성자 - 비즈니스 규칙 검증
func NewDomain(name, description string) (*Domain, error) {
    if name == "" {
        return nil, errors.New("domain name cannot be empty")
    }
    if len(name) > 100 {
        return nil, errors.New("domain name too long")
    }
    
    return &Domain{
        name:        name,
        description: description,
        createdAt:   time.Now(),
    }, nil
}

// 불변 메서드들
func (d *Domain) ID() int           { return d.id }
func (d *Domain) Name() string      { return d.name }
func (d *Domain) Description() string { return d.description }
func (d *Domain) CreatedAt() time.Time { return d.createdAt }

// 새로운 인스턴스 반환 (불변성 유지)
func (d *Domain) WithDescription(description string) *Domain {
    return &Domain{
        id:          d.id,
        name:        d.name,
        description: description,
        createdAt:   d.createdAt,
    }
}

// 비즈니스 로직
func (d *Domain) IsActive() bool {
    return d.name != ""
}

func (d *Domain) CanBeDeleted() bool {
    return d.IsActive() && d.createdAt.Before(time.Now().Add(-24*time.Hour))
}
```

### 테스트 예시

```go
// internal/domain/entity/domain_test.go
func TestNewDomain(t *testing.T) {
    tests := []struct {
        name        string
        domainName  string
        description string
        wantErr     bool
        errMessage  string
    }{
        {
            name:        "valid domain",
            domainName:  "test-domain",
            description: "Test description",
            wantErr:     false,
        },
        {
            name:        "empty name",
            domainName:  "",
            description: "Test description",
            wantErr:     true,
            errMessage:  "domain name cannot be empty",
        },
        {
            name:        "name too long",
            domainName:  strings.Repeat("a", 101),
            description: "Test description",
            wantErr:     true,
            errMessage:  "domain name too long",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            domain, err := NewDomain(tt.domainName, tt.description)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMessage)
                assert.Nil(t, domain)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, domain)
                assert.Equal(t, tt.domainName, domain.Name())
                assert.Equal(t, tt.description, domain.Description())
                assert.True(t, domain.IsActive())
            }
        })
    }
}

func TestDomain_WithDescription(t *testing.T) {
    original, _ := NewDomain("test", "original")
    updated := original.WithDescription("updated")
    
    // 원본은 변경되지 않음
    assert.Equal(t, "original", original.Description())
    // 새로운 인스턴스가 생성됨
    assert.Equal(t, "updated", updated.Description())
    assert.NotEqual(t, original, updated)
}
```

## 2. 값 객체 설계

### Before: 단순한 문자열

```go
// internal/models/composite_key.go
type CompositeKey string

func (ck CompositeKey) String() string {
    return string(ck)
}
```

### After: 강력한 값 객체

```go
// internal/domain/valueobject/composite_key.go
type CompositeKey struct {
    toolName   string
    domainName string
    id         string
}

func NewCompositeKey(toolName, domainName, id string) (*CompositeKey, error) {
    if toolName == "" {
        return nil, errors.New("tool name cannot be empty")
    }
    if domainName == "" {
        return nil, errors.New("domain name cannot be empty")
    }
    if id == "" {
        return nil, errors.New("id cannot be empty")
    }
    
    return &CompositeKey{
        toolName:   normalizeToolName(toolName),
        domainName: normalizeDomainName(domainName),
        id:         id,
    }, nil
}

func ParseCompositeKey(key string) (*CompositeKey, error) {
    parts := strings.Split(key, ":")
    if len(parts) != 3 {
        return nil, errors.New("invalid composite key format")
    }
    
    return NewCompositeKey(parts[0], parts[1], parts[2])
}

func (ck *CompositeKey) ToolName() string   { return ck.toolName }
func (ck *CompositeKey) DomainName() string { return ck.domainName }
func (ck *CompositeKey) ID() string         { return ck.id }

func (ck *CompositeKey) String() string {
    return fmt.Sprintf("%s:%s:%s", ck.toolName, ck.domainName, ck.id)
}

func (ck *CompositeKey) Equals(other *CompositeKey) bool {
    return ck.toolName == other.toolName &&
           ck.domainName == other.domainName &&
           ck.id == other.id
}

// 순수 함수들
func normalizeToolName(name string) string {
    return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

func normalizeDomainName(name string) string {
    return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
```

### 테스트 예시

```go
// internal/domain/valueobject/composite_key_test.go
func TestNewCompositeKey(t *testing.T) {
    tests := []struct {
        name       string
        toolName   string
        domainName string
        id         string
        wantErr    bool
    }{
        {
            name:       "valid composite key",
            toolName:   "test-tool",
            domainName: "test-domain",
            id:         "123",
            wantErr:    false,
        },
        {
            name:       "empty tool name",
            toolName:   "",
            domainName: "test-domain",
            id:         "123",
            wantErr:    true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            key, err := NewCompositeKey(tt.toolName, tt.domainName, tt.id)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, key)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, key)
                assert.Equal(t, tt.toolName, key.ToolName())
                assert.Equal(t, tt.domainName, key.DomainName())
                assert.Equal(t, tt.id, key.ID())
            }
        })
    }
}

func TestCompositeKey_Normalization(t *testing.T) {
    key, err := NewCompositeKey("Test Tool", "Test Domain", "123")
    assert.NoError(t, err)
    
    // 자동 정규화
    assert.Equal(t, "test-tool", key.ToolName())
    assert.Equal(t, "test-domain", key.DomainName())
}
```

## 3. 서비스 레이어 설계

### Before: 강한 결합

```go
// internal/services/domain.go
type DomainService struct {
    db     *sql.DB
    logger *log.Logger
    cache  *redis.Client
}

func (s *DomainService) CreateDomain(name, description string) error {
    // 직접 데이터베이스 접근
    _, err := s.db.Exec("INSERT INTO domains (name, description) VALUES (?, ?)", name, description)
    if err != nil {
        s.logger.Printf("Failed to create domain: %v", err)
        return err
    }
    
    // 직접 캐시 무효화
    s.cache.Del("domains")
    
    return nil
}
```

### After: 의존성 주입과 인터페이스 분리

```go
// internal/domain/repository/domain.go
type DomainRepository interface {
    Create(ctx context.Context, domain *entity.Domain) error
    GetByName(ctx context.Context, name string) (*entity.Domain, error)
    List(ctx context.Context, limit, offset int) ([]*entity.Domain, error)
    Update(ctx context.Context, domain *entity.Domain) error
    Delete(ctx context.Context, name string) error
}

// internal/application/usecase/domain/create.go
type CreateDomainUseCase struct {
    repo DomainRepository
}

func NewCreateDomainUseCase(repo DomainRepository) *CreateDomainUseCase {
    return &CreateDomainUseCase{repo: repo}
}

func (uc *CreateDomainUseCase) Execute(ctx context.Context, req *CreateDomainRequest) (*CreateDomainResponse, error) {
    // 비즈니스 로직 검증
    domain, err := entity.NewDomain(req.Name, req.Description)
    if err != nil {
        return nil, err
    }
    
    // 중복 검사
    existing, _ := uc.repo.GetByName(ctx, req.Name)
    if existing != nil {
        return nil, errors.New("domain already exists")
    }
    
    // 저장
    if err := uc.repo.Create(ctx, domain); err != nil {
        return nil, err
    }
    
    return &CreateDomainResponse{
        Name:        domain.Name(),
        Description: domain.Description(),
        CreatedAt:   domain.CreatedAt(),
    }, nil
}

// internal/application/dto/request/create_domain.go
type CreateDomainRequest struct {
    Name        string `json:"name" validate:"required,max=100"`
    Description string `json:"description" validate:"max=500"`
}

// internal/application/dto/response/domain.go
type CreateDomainResponse struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### 테스트 예시

```go
// internal/application/usecase/domain/create_test.go
func TestCreateDomainUseCase_Execute(t *testing.T) {
    tests := []struct {
        name        string
        request     *CreateDomainRequest
        setupMock   func(*MockDomainRepository)
        wantErr     bool
        errMessage  string
    }{
        {
            name: "successful creation",
            request: &CreateDomainRequest{
                Name:        "test-domain",
                Description: "Test description",
            },
            setupMock: func(mock *MockDomainRepository) {
                mock.On("GetByName", mock.Anything, "test-domain").Return(nil, nil)
                mock.On("Create", mock.Anything, mock.AnythingOfType("*entity.Domain")).Return(nil)
            },
            wantErr: false,
        },
        {
            name: "domain already exists",
            request: &CreateDomainRequest{
                Name:        "existing-domain",
                Description: "Test description",
            },
            setupMock: func(mock *MockDomainRepository) {
                existingDomain, _ := entity.NewDomain("existing-domain", "Existing")
                mock.On("GetByName", mock.Anything, "existing-domain").Return(existingDomain, nil)
            },
            wantErr:    true,
            errMessage: "domain already exists",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &MockDomainRepository{}
            tt.setupMock(mockRepo)
            
            useCase := NewCreateDomainUseCase(mockRepo)
            
            result, err := useCase.Execute(context.Background(), tt.request)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMessage)
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                assert.Equal(t, tt.request.Name, result.Name)
                assert.Equal(t, tt.request.Description, result.Description)
            }
            
            mockRepo.AssertExpectations(t)
        })
    }
}
```

## 4. 팩토리 패턴 적용

### 테스트 데이터 팩토리

```go
// internal/test/factory/domain_factory.go
type DomainFactory struct{}

func NewDomainFactory() *DomainFactory {
    return &DomainFactory{}
}

func (f *DomainFactory) CreateValidDomain(name string) *entity.Domain {
    domain, _ := entity.NewDomain(name, "Test description")
    return domain
}

func (f *DomainFactory) CreateDomainWithDescription(name, description string) *entity.Domain {
    domain, _ := entity.NewDomain(name, description)
    return domain
}

func (f *DomainFactory) CreateInvalidDomain() *entity.Domain {
    // 의도적으로 잘못된 도메인 생성
    return &entity.Domain{}
}

// internal/test/factory/node_factory.go
type NodeFactory struct{}

func NewNodeFactory() *NodeFactory {
    return &NodeFactory{}
}

func (f *NodeFactory) CreateValidNode(url string) *entity.Node {
    node, _ := entity.NewNode(url, "Test title")
    return node
}

func (f *NodeFactory) CreateNodeWithAttributes(url string, attributes []*entity.NodeAttribute) *entity.Node {
    node, _ := entity.NewNode(url, "Test title")
    // 속성 추가 로직
    return node
}

// internal/test/factory/attribute_factory.go
type AttributeFactory struct{}

func NewAttributeFactory() *AttributeFactory {
    return &AttributeFactory{}
}

func (f *AttributeFactory) CreateTagAttribute(name, value string) *entity.Attribute {
    attr, _ := entity.NewAttribute(name, valueobject.TagType, "Test attribute")
    return attr
}

func (f *AttributeFactory) CreateNumberAttribute(name string, value float64) *entity.Attribute {
    attr, _ := entity.NewAttribute(name, valueobject.NumberType, fmt.Sprintf("%f", value))
    return attr
}
```

### 테스트에서 사용

```go
// internal/application/usecase/domain/create_integration_test.go
func TestCreateDomainUseCase_Integration(t *testing.T) {
    // 팩토리 사용
    domainFactory := NewDomainFactory()
    nodeFactory := NewNodeFactory()
    attrFactory := NewAttributeFactory()
    
    // 테스트 데이터 생성
    domain := domainFactory.CreateValidDomain("test-domain")
    node := nodeFactory.CreateValidNode("https://example.com")
    attribute := attrFactory.CreateTagAttribute("category", "tech")
    
    // 테스트 로직...
}
```

## 5. 이벤트 기반 설계

### 도메인 이벤트 정의

```go
// internal/domain/event/events.go
type DomainCreatedEvent struct {
    DomainID   int       `json:"domain_id"`
    DomainName string    `json:"domain_name"`
    CreatedAt  time.Time `json:"created_at"`
}

type DomainDeletedEvent struct {
    DomainID   int       `json:"domain_id"`
    DomainName string    `json:"domain_name"`
    DeletedAt  time.Time `json:"deleted_at"`
}

type NodeCreatedEvent struct {
    NodeID     int       `json:"node_id"`
    DomainID   int       `json:"domain_id"`
    URL        string    `json:"url"`
    CreatedAt  time.Time `json:"created_at"`
}

// internal/domain/event/publisher.go
type EventPublisher interface {
    Publish(eventType string, event interface{}) error
}

// internal/application/usecase/domain/create_with_events.go
type CreateDomainUseCaseWithEvents struct {
    repo   DomainRepository
    events EventPublisher
}

func NewCreateDomainUseCaseWithEvents(repo DomainRepository, events EventPublisher) *CreateDomainUseCaseWithEvents {
    return &CreateDomainUseCaseWithEvents{
        repo:   repo,
        events: events,
    }
}

func (uc *CreateDomainUseCaseWithEvents) Execute(ctx context.Context, req *CreateDomainRequest) (*CreateDomainResponse, error) {
    domain, err := entity.NewDomain(req.Name, req.Description)
    if err != nil {
        return nil, err
    }
    
    if err := uc.repo.Create(ctx, domain); err != nil {
        return nil, err
    }
    
    // 이벤트 발행
    event := DomainCreatedEvent{
        DomainID:   domain.ID(),
        DomainName: domain.Name(),
        CreatedAt:  domain.CreatedAt(),
    }
    
    if err := uc.events.Publish("domain.created", event); err != nil {
        // 로그만 남기고 계속 진행
        log.Printf("Failed to publish domain.created event: %v", err)
    }
    
    return &CreateDomainResponse{
        Name:        domain.Name(),
        Description: domain.Description(),
        CreatedAt:   domain.CreatedAt(),
    }, nil
}
```

### 테스트 예시

```go
// internal/application/usecase/domain/create_with_events_test.go
func TestCreateDomainUseCaseWithEvents_Execute(t *testing.T) {
    mockRepo := &MockDomainRepository{}
    mockEvents := &MockEventPublisher{}
    
    useCase := NewCreateDomainUseCaseWithEvents(mockRepo, mockEvents)
    
    req := &CreateDomainRequest{
        Name:        "test-domain",
        Description: "Test description",
    }
    
    // Mock 설정
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Domain")).Return(nil)
    mockEvents.On("Publish", "domain.created", mock.AnythingOfType("DomainCreatedEvent")).Return(nil)
    
    result, err := useCase.Execute(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "test-domain", result.Name)
    
    mockRepo.AssertExpectations(t)
    mockEvents.AssertExpectations(t)
}
```

## 6. 설정 주입 패턴

### 설정 구조체

```go
// internal/config/config.go
type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Cache    CacheConfig    `yaml:"cache"`
}

type DatabaseConfig struct {
    URL      string        `yaml:"url"`
    MaxConns int           `yaml:"max_conns"`
    Timeout  time.Duration `yaml:"timeout"`
}

type ServerConfig struct {
    Port    int           `yaml:"port"`
    Timeout time.Duration `yaml:"timeout"`
}

type CacheConfig struct {
    URL      string        `yaml:"url"`
    Timeout  time.Duration `yaml:"timeout"`
}

// internal/application/service/domain_service.go
type DomainService struct {
    repo   DomainRepository
    config Config
}

func NewDomainService(repo DomainRepository, config Config) *DomainService {
    return &DomainService{
        repo:   repo,
        config: config,
    }
}

func (s *DomainService) CreateDomain(ctx context.Context, name string) error {
    // 설정 기반 로직
    if len(name) > s.config.Server.MaxNameLength {
        return errors.New("name too long")
    }
    
    domain, err := entity.NewDomain(name, "")
    if err != nil {
        return err
    }
    
    return s.repo.Create(ctx, domain)
}
```

### 테스트에서 사용

```go
// internal/application/service/domain_service_test.go
func TestDomainService_CreateDomain(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        domainName string
        wantErr bool
    }{
        {
            name: "valid domain",
            config: Config{
                Server: ServerConfig{MaxNameLength: 100},
            },
            domainName: "test-domain",
            wantErr:    false,
        },
        {
            name: "name too long",
            config: Config{
                Server: ServerConfig{MaxNameLength: 10},
            },
            domainName: "very-long-domain-name",
            wantErr:    true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &MockDomainRepository{}
            service := NewDomainService(mockRepo, tt.config)
            
            err := service.CreateDomain(context.Background(), tt.domainName)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 결론

이러한 설계 원칙들을 적용하면:

1. **테스트 작성이 쉬워집니다** - 의존성 주입과 인터페이스 분리
2. **비즈니스 로직이 명확해집니다** - 값 객체와 순수 함수
3. **코드 재사용성이 높아집니다** - 팩토리 패턴과 헬퍼 함수
4. **유지보수가 쉬워집니다** - 이벤트 기반 설계와 설정 주입

이러한 원칙들은 Clean Architecture와 완벽하게 맞으며, 테스트 가능한 코드를 자연스럽게 만들어줍니다. 