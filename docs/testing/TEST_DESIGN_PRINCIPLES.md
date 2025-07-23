# 테스트 친화적 설계 원칙 (Test-Friendly Design Principles)

## 개요

테스트 코드 작성의 어려움을 해결하기 위한 설계 원칙들을 정의합니다. 이 원칙들은 코드를 테스트하기 쉽게 만들면서도 비즈니스 가치를 높이는 것을 목표로 합니다.

## 핵심 원칙

### 1. 의존성 주입 (Dependency Injection)

**문제**: 하드코딩된 의존성으로 인한 테스트 어려움
**해결**: 인터페이스를 통한 의존성 주입

```go
// ❌ 나쁜 예 - 하드코딩된 의존성
type DomainService struct {
    db *sql.DB // 직접 의존
}

func (s *DomainService) CreateDomain(name string) error {
    _, err := s.db.Exec("INSERT INTO domains (name) VALUES (?)", name)
    return err
}

// ✅ 좋은 예 - 의존성 주입
type DomainRepository interface {
    Create(ctx context.Context, domain *entity.Domain) error
}

type DomainService struct {
    repo DomainRepository // 인터페이스 의존
}

func (s *DomainService) CreateDomain(ctx context.Context, name string) error {
    domain, err := entity.NewDomain(name)
    if err != nil {
        return err
    }
    return s.repo.Create(ctx, domain)
}
```

### 2. 순수 함수 (Pure Functions)

**문제**: 부수 효과로 인한 테스트 복잡성
**해결**: 순수 함수로 비즈니스 로직 분리

```go
// ❌ 나쁜 예 - 부수 효과가 있는 함수
func (s *Service) ProcessNode(node *Node) error {
    // 데이터베이스 접근
    // 로그 기록
    // 외부 API 호출
    // 파일 시스템 접근
    return nil
}

// ✅ 좋은 예 - 순수 함수 분리
func ValidateNode(node *Node) error {
    // 순수한 검증 로직만
    if node.URL == "" {
        return errors.New("URL is required")
    }
    return nil
}

func (s *Service) ProcessNode(node *Node) error {
    // 검증
    if err := ValidateNode(node); err != nil {
        return err
    }
    
    // 저장
    return s.repo.Save(node)
}
```

### 3. 값 객체 (Value Objects)

**문제**: 복잡한 상태 관리로 인한 테스트 어려움
**해결**: 불변 값 객체 사용

```go
// ❌ 나쁜 예 - 가변 상태
type Domain struct {
    Name        string
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (d *Domain) UpdateDescription(desc string) {
    d.Description = desc
    d.UpdatedAt = time.Now() // 부수 효과
}

// ✅ 좋은 예 - 불변 값 객체
type Domain struct {
    name        string
    description string
    createdAt   time.Time
}

func NewDomain(name, description string) (*Domain, error) {
    if name == "" {
        return nil, errors.New("name cannot be empty")
    }
    return &Domain{
        name:        name,
        description: description,
        createdAt:   time.Now(),
    }, nil
}

func (d *Domain) WithDescription(description string) *Domain {
    return &Domain{
        name:        d.name,
        description: description,
        createdAt:   d.createdAt,
    }
}
```

### 4. 팩토리 패턴 (Factory Pattern)

**문제**: 복잡한 객체 생성으로 인한 테스트 어려움
**해결**: 팩토리 함수 사용

```go
// ❌ 나쁜 예 - 복잡한 생성자
func NewComplexNode(url, title string, attributes map[string]interface{}, metadata map[string]string) *Node {
    // 복잡한 초기화 로직
    return &Node{...}
}

// ✅ 좋은 예 - 팩토리 패턴
type NodeBuilder struct {
    url        string
    title      string
    attributes []Attribute
    metadata   map[string]string
}

func NewNodeBuilder() *NodeBuilder {
    return &NodeBuilder{
        attributes: make([]Attribute, 0),
        metadata:   make(map[string]string),
    }
}

func (b *NodeBuilder) WithURL(url string) *NodeBuilder {
    b.url = url
    return b
}

func (b *NodeBuilder) WithTitle(title string) *NodeBuilder {
    b.title = title
    return b
}

func (b *NodeBuilder) Build() (*Node, error) {
    return NewNode(b.url, b.title, b.attributes, b.metadata)
}

// 테스트에서 사용
func TestNodeCreation(t *testing.T) {
    node, err := NewNodeBuilder().
        WithURL("https://example.com").
        WithTitle("Example").
        Build()
    
    assert.NoError(t, err)
    assert.Equal(t, "https://example.com", node.URL())
}
```

### 5. 인터페이스 분리 (Interface Segregation)

**문제**: 거대한 인터페이스로 인한 모킹 어려움
**해결**: 작은 인터페이스로 분리

```go
// ❌ 나쁜 예 - 거대한 인터페이스
type Repository interface {
    Create(ctx context.Context, entity interface{}) error
    Get(ctx context.Context, id int) (interface{}, error)
    Update(ctx context.Context, entity interface{}) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, limit, offset int) ([]interface{}, error)
    Count(ctx context.Context) (int, error)
    // ... 수십 개의 메서드
}

// ✅ 좋은 예 - 작은 인터페이스
type Reader interface {
    Get(ctx context.Context, id int) (interface{}, error)
    List(ctx context.Context, limit, offset int) ([]interface{}, error)
}

type Writer interface {
    Create(ctx context.Context, entity interface{}) error
    Update(ctx context.Context, entity interface{}) error
    Delete(ctx context.Context, id int) error
}

type Repository interface {
    Reader
    Writer
}
```

### 6. 이벤트 기반 설계 (Event-Driven Design)

**문제**: 강한 결합으로 인한 테스트 어려움
**해결**: 이벤트를 통한 느슨한 결합

```go
// ❌ 나쁜 예 - 강한 결합
type DomainService struct {
    repo     DomainRepository
    logger   Logger
    notifier Notifier
    cache    Cache
}

func (s *DomainService) CreateDomain(name string) error {
    domain := &Domain{Name: name}
    
    // 직접 호출 - 테스트 어려움
    if err := s.repo.Create(domain); err != nil {
        s.logger.Error("Failed to create domain", err)
        return err
    }
    
    s.notifier.Notify("domain.created", domain)
    s.cache.Invalidate("domains")
    
    return nil
}

// ✅ 좋은 예 - 이벤트 기반
type DomainCreatedEvent struct {
    DomainID   int
    DomainName string
    CreatedAt  time.Time
}

type DomainService struct {
    repo   DomainRepository
    events EventPublisher
}

func (s *DomainService) CreateDomain(name string) error {
    domain := &Domain{Name: name}
    
    if err := s.repo.Create(domain); err != nil {
        return err
    }
    
    // 이벤트 발행 - 느슨한 결합
    s.events.Publish("domain.created", DomainCreatedEvent{
        DomainID:   domain.ID,
        DomainName: domain.Name,
        CreatedAt:  domain.CreatedAt,
    })
    
    return nil
}
```

### 7. 설정 주입 (Configuration Injection)

**문제**: 하드코딩된 설정으로 인한 테스트 어려움
**해결**: 설정을 의존성으로 주입

```go
// ❌ 나쁜 예 - 하드코딩된 설정
func NewService() *Service {
    return &Service{
        dbURL: "sqlite:///prod.db",
        cache: redis.NewClient(&redis.Options{
            Addr: "localhost:6379",
        }),
        timeout: 30 * time.Second,
    }
}

// ✅ 좋은 예 - 설정 주입
type Config struct {
    DatabaseURL string
    CacheAddr   string
    Timeout     time.Duration
}

func NewService(config Config) *Service {
    return &Service{
        dbURL:   config.DatabaseURL,
        cache:   redis.NewClient(&redis.Options{Addr: config.CacheAddr}),
        timeout: config.Timeout,
    }
}

// 테스트에서 사용
func TestService(t *testing.T) {
    config := Config{
        DatabaseURL: "sqlite:///test.db",
        CacheAddr:   "localhost:6379",
        Timeout:     5 * time.Second,
    }
    
    service := NewService(config)
    // 테스트 로직...
}
```

### 8. 시간 의존성 분리 (Time Dependency Separation)

**문제**: 시간에 의존적인 로직으로 인한 테스트 어려움
**해결**: 시간을 의존성으로 주입

```go
// ❌ 나쁜 예 - 직접 시간 사용
func (s *Service) CreateDomain(name string) *Domain {
    return &Domain{
        Name:      name,
        CreatedAt: time.Now(), // 테스트 어려움
    }
}

// ✅ 좋은 예 - 시간 주입
type Clock interface {
    Now() time.Time
}

type RealClock struct{}

func (c *RealClock) Now() time.Time {
    return time.Now()
}

type Service struct {
    clock Clock
}

func (s *Service) CreateDomain(name string) *Domain {
    return &Domain{
        Name:      name,
        CreatedAt: s.clock.Now(),
    }
}

// 테스트에서 사용
type MockClock struct {
    now time.Time
}

func (c *MockClock) Now() time.Time {
    return c.now
}

func TestCreateDomain(t *testing.T) {
    mockClock := &MockClock{now: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}
    service := &Service{clock: mockClock}
    
    domain := service.CreateDomain("test")
    assert.Equal(t, mockClock.now, domain.CreatedAt)
}
```

### 9. 에러 처리 표준화 (Error Handling Standardization)

**문제**: 다양한 에러 타입으로 인한 테스트 복잡성
**해결**: 표준화된 에러 타입 사용

```go
// ❌ 나쁜 예 - 다양한 에러 타입
func (s *Service) CreateDomain(name string) error {
    if name == "" {
        return errors.New("name is empty")
    }
    if len(name) > 100 {
        return fmt.Errorf("name too long: %d", len(name))
    }
    if s.repo.Exists(name) {
        return errors.New("domain already exists")
    }
    return s.repo.Create(&Domain{Name: name})
}

// ✅ 좋은 예 - 표준화된 에러
type DomainError struct {
    Code    string
    Message string
    Cause   error
}

func (e *DomainError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
    ErrDomainNameEmpty     = &DomainError{Code: "DOMAIN_NAME_EMPTY", Message: "Domain name cannot be empty"}
    ErrDomainNameTooLong   = &DomainError{Code: "DOMAIN_NAME_TOO_LONG", Message: "Domain name too long"}
    ErrDomainAlreadyExists = &DomainError{Code: "DOMAIN_ALREADY_EXISTS", Message: "Domain already exists"}
)

func (s *Service) CreateDomain(name string) error {
    if name == "" {
        return ErrDomainNameEmpty
    }
    if len(name) > 100 {
        return ErrDomainNameTooLong
    }
    if s.repo.Exists(name) {
        return ErrDomainAlreadyExists
    }
    return s.repo.Create(&Domain{Name: name})
}

// 테스트에서 사용
func TestCreateDomain_WithEmptyName_ShouldReturnError(t *testing.T) {
    service := NewService(mockRepo)
    
    err := service.CreateDomain("")
    
    assert.Error(t, err)
    assert.IsType(t, &DomainError{}, err)
    domainErr := err.(*DomainError)
    assert.Equal(t, "DOMAIN_NAME_EMPTY", domainErr.Code)
}
```

### 10. 테스트 헬퍼 함수 (Test Helper Functions)

**문제**: 반복적인 테스트 설정 코드
**해결**: 재사용 가능한 헬퍼 함수

```go
// 테스트 헬퍼 함수들
func NewTestDomain(name string) *Domain {
    domain, _ := NewDomain(name, "Test description")
    return domain
}

func NewTestNode(url string) *Node {
    node, _ := NewNode(url, "Test title")
    return node
}

func NewMockRepository() *MockDomainRepository {
    return &MockDomainRepository{}
}

func NewTestService() *DomainService {
    return &DomainService{
        repo: NewMockRepository(),
    }
}

// 테스트에서 사용
func TestCreateDomain(t *testing.T) {
    service := NewTestService()
    domain := NewTestDomain("test-domain")
    
    err := service.Create(domain)
    
    assert.NoError(t, err)
}
```

## 테스트 작성 가이드라인

### 1. 테스트 구조화

```go
func Test[FunctionName]_[Scenario]_[ExpectedResult](t *testing.T) {
    // Arrange - 테스트 데이터 준비
    service := NewTestService()
    input := NewTestInput()
    
    // Act - 테스트 실행
    result, err := service.Function(input)
    
    // Assert - 결과 검증
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### 2. 테스트 데이터 관리

```go
// 테스트 데이터 팩토리
type TestDataFactory struct{}

func (f *TestDataFactory) CreateDomain(name string) *Domain {
    return NewTestDomain(name)
}

func (f *TestDataFactory) CreateNode(url string) *Node {
    return NewTestNode(url)
}

func (f *TestDataFactory) CreateAttribute(name, value string) *Attribute {
    return NewTestAttribute(name, value)
}

// 테스트에서 사용
func TestComplexScenario(t *testing.T) {
    factory := &TestDataFactory{}
    
    domain := factory.CreateDomain("test-domain")
    node := factory.CreateNode("https://example.com")
    attribute := factory.CreateAttribute("category", "tech")
    
    // 테스트 로직...
}
```

### 3. 모킹 전략

```go
// 인터페이스 기반 모킹
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, domain *Domain) error {
    args := m.Called(ctx, domain)
    return args.Error(0)
}

func (m *MockRepository) Get(ctx context.Context, id int) (*Domain, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*Domain), args.Error(1)
}

// 테스트에서 사용
func TestCreateDomain(t *testing.T) {
    mockRepo := &MockRepository{}
    service := &DomainService{repo: mockRepo}
    
    domain := NewTestDomain("test")
    mockRepo.On("Create", mock.Anything, domain).Return(nil)
    
    err := service.Create(domain)
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

## 결론

이러한 설계 원칙들을 적용하면:

1. **테스트 작성이 쉬워집니다** - 의존성 주입과 인터페이스 분리
2. **테스트 실행이 빨라집니다** - 순수 함수와 모킹
3. **테스트 유지보수가 쉬워집니다** - 표준화된 에러 처리와 헬퍼 함수
4. **비즈니스 로직이 명확해집니다** - 값 객체와 팩토리 패턴

이 원칙들은 Clean Architecture와 잘 맞으며, 테스트 가능한 코드를 자연스럽게 만들어줍니다. 