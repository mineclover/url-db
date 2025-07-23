# URL-DB 테스트 전략

## 개요

URL-DB 프로젝트는 Clean Architecture 기반으로 설계된 MCP(Management Control Plane) 서버로, 수만 URL 그래프를 실시간으로 전파·분석하는 핵심 컴포넌트입니다. 본 문서는 **6-단계 테스트 피라미드**를 기반으로 한 체계적인 테스트 전략을 제시합니다.

## 1. 테스트 피라미드

### 1.1 6-단계 피라미드 구조

```
    Performance & Load Tests (5%)
           ▲
    End-to-End Scenario Tests (10%)
           ▲
    Contract Tests (10%)
           ▲
    Integration Tests (20%)
           ▲
    Property-Based & Fuzz Tests (15%)
           ▲
    Unit Tests (40%)
```

### 1.2 각 단계별 목적과 도구

| 단계 | 비율 | 목적 | 도구 | 대상 |
|------|------|------|------|------|
| **Unit Tests** | 40% | 순수 함수, 값 객체, DTO | `testing`, `testify` | Domain, Application Layer |
| **Property-Based & Fuzz Tests** | 15% | 경계·인코딩 취약점 | `rapid`, `gopter`, Go 1.18+ fuzzing | Value Objects, Parsers |
| **Integration Tests** | 20% | gRPC bufconn, Testcontainers | `bufconn`, `testcontainers-go` | Infrastructure Layer |
| **Contract Tests** | 10% | Pact (Producer·Consumer) | `pact-go` | Interface Layer |
| **End-to-End Scenario Tests** | 10% | 실제 MCP workflow | `httptest`, `grpctest` | 전체 시스템 |
| **Performance & Load Tests** | 5% | benchmark + pprof | `testing.B`, `pprof` | 핫패스, 부하 |

## 2. 레이어별 테스트 전략

### 2.1 Domain Layer (40% - Unit Tests)

#### 엔티티 테스트
```go
// internal/domain/entity/domain_test.go
func TestNewDomain(t *testing.T) {
    tests := []struct {
        name        string
        domainName  string
        description string
        wantErr     bool
    }{
        {"valid domain", "test-domain", "Test description", false},
        {"empty name", "", "Test description", true},
        {"name too long", strings.Repeat("a", 101), "Test description", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            domain, err := entity.NewDomain(tt.domainName, tt.description)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, domain)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, domain)
                assert.Equal(t, tt.domainName, domain.Name())
            }
        })
    }
}
```

#### 값 객체 테스트
```go
// internal/domain/valueobject/composite_key_test.go
func TestCompositeKey_Normalization(t *testing.T) {
    key, err := NewCompositeKey("Test Tool", "Test Domain", "123")
    assert.NoError(t, err)
    
    // 자동 정규화 검증
    assert.Equal(t, "test-tool", key.ToolName())
    assert.Equal(t, "test-domain", key.DomainName())
}
```

#### 도메인 서비스 테스트
```go
// internal/domain/service/dependency_graph_test.go
func TestDependencyGraphService_DetectCycles(t *testing.T) {
    service := NewDependencyGraphService(mockRepo, nil)
    
    cycles, err := service.DetectCycles(context.Background(), dependencies)
    assert.NoError(t, err)
    assert.Len(t, cycles, expectedCycleCount)
}
```

### 2.2 Application Layer (40% - Unit Tests)

#### 유스케이스 테스트
```go
// internal/application/usecase/domain/create_test.go
func TestCreateDomainUseCase_Execute(t *testing.T) {
    mockRepo := &MockDomainRepository{}
    useCase := NewCreateDomainUseCase(mockRepo)
    
    req := &CreateDomainRequest{
        Name:        "test-domain",
        Description: "Test description",
    }
    
    // Mock 설정
    mockRepo.On("GetByName", mock.Anything, "test-domain").Return(nil, nil)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Domain")).Return(nil)
    
    result, err := useCase.Execute(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "test-domain", result.Name)
    mockRepo.AssertExpectations(t)
}
```

#### DTO 테스트
```go
// internal/application/dto/request/create_domain_test.go
func TestCreateDomainRequest_Validation(t *testing.T) {
    tests := []struct {
        name    string
        request CreateDomainRequest
        wantErr bool
    }{
        {"valid request", CreateDomainRequest{Name: "test", Description: "desc"}, false},
        {"empty name", CreateDomainRequest{Name: "", Description: "desc"}, true},
        {"name too long", CreateDomainRequest{Name: strings.Repeat("a", 101), Description: "desc"}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.request.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 2.3 Infrastructure Layer (20% - Integration Tests)

#### 리포지토리 테스트
```go
// internal/infrastructure/persistence/sqlite/repository/domain_test.go
func TestDomainRepository_Integration(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewDomainRepository(db)
    
    domain, err := entity.NewDomain("test-domain", "Test description")
    assert.NoError(t, err)
    
    // Create
    err = repo.Create(context.Background(), domain)
    assert.NoError(t, err)
    
    // GetByName
    retrieved, err := repo.GetByName(context.Background(), "test-domain")
    assert.NoError(t, err)
    assert.Equal(t, domain.Name(), retrieved.Name())
}
```

#### gRPC 서버 테스트 (bufconn)
```go
// internal/interface/mcp/grpc_integration_test.go
func TestMCPService_Integration(t *testing.T) {
    lis := bufconn.Listen(1 << 20)
    grpcServer := grpc.NewServer()
    mcp.RegisterMCPServiceServer(grpcServer, &mcp.Server{})
    go grpcServer.Serve(lis)
    
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet",
        grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
            return lis.Dial()
        }),
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    assert.NoError(t, err)
    defer conn.Close()
    
    client := mcp.NewMCPServiceClient(conn)
    resp, err := client.CreateDomain(ctx, &mcp.CreateDomainRequest{
        Name: "test-domain",
    })
    
    assert.NoError(t, err)
    assert.Equal(t, "test-domain", resp.Name)
}
```

### 2.4 Interface Layer (10% - Contract Tests)

#### HTTP 핸들러 테스트
```go
// internal/interface/http/handler/domain_test.go
func TestDomainHandler_CreateDomain(t *testing.T) {
    mockUseCase := &MockCreateDomainUseCase{}
    handler := NewDomainHandler(mockUseCase)
    
    reqBody := `{"name":"test-domain","description":"Test description"}`
    req := httptest.NewRequest("POST", "/domains", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("*CreateDomainRequest")).
        Return(&CreateDomainResponse{Name: "test-domain"}, nil)
    
    handler.CreateDomain(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    mockUseCase.AssertExpectations(t)
}
```

#### MCP 핸들러 테스트
```go
// internal/interface/mcp/handler_test.go
func TestMCPHandler_ListDomains(t *testing.T) {
    mockUseCase := &MockListDomainsUseCase{}
    handler := NewMCPHandler(mockUseCase)
    
    request := &mcp.ListDomainsRequest{}
    response, err := handler.ListDomains(context.Background(), request)
    
    assert.NoError(t, err)
    assert.NotNil(t, response)
    mockUseCase.AssertExpectations(t)
}
```

## 3. 고급 테스트 기법

### 3.1 속성 기반 테스트 (15% - Property-Based & Fuzz Tests)

#### Rapid를 사용한 속성 기반 테스트
```go
// internal/domain/valueobject/composite_key_property_test.go
func TestCompositeKey_PropertyBased(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        toolName := rapid.String().Draw(t, "toolName")
        domainName := rapid.String().Draw(t, "domainName")
        id := rapid.String().Draw(t, "id")
        
        key, err := NewCompositeKey(toolName, domainName, id)
        if err != nil {
            t.Skip()
        }
        
        // 속성: 정규화 후 파싱하면 동일한 결과
        normalized := key.String()
        parsed, err := ParseCompositeKey(normalized)
        
        if err != nil {
            t.Fatalf("Failed to parse normalized key: %v", err)
        }
        
        if !key.Equals(parsed) {
            t.Fatalf("Normalization property violated")
        }
    })
}
```

#### Go 1.18+ 내장 Fuzzing
```go
// internal/domain/valueobject/composite_key_fuzz_test.go
func FuzzCompositeKeyParse(f *testing.F) {
    f.Add("tool:domain:123")
    f.Add("test-tool:test-domain:456")
    f.Add("invalid-key")
    
    f.Fuzz(func(t *testing.T, input string) {
        key, err := ParseCompositeKey(input)
        
        if err == nil {
            normalized := key.String()
            parsed, err2 := ParseCompositeKey(normalized)
            
            if err2 != nil {
                t.Fatalf("Failed to parse normalized key: %v", err2)
            }
            
            if !key.Equals(parsed) {
                t.Fatalf("Parse-String roundtrip failed")
            }
        }
    })
}
```

### 3.2 동시성 테스트

#### Race Detector
```bash
# CI에서 레이스 감지
go test -race ./...
```

#### 고루틴 누수 방지
```go
// internal/application/usecase/domain/create_leak_test.go
func TestCreateDomainUseCase_NoGoroutineLeak(t *testing.T) {
    defer goleak.VerifyNone(t)
    
    mockRepo := &MockDomainRepository{}
    useCase := NewCreateDomainUseCase(mockRepo)
    
    req := &CreateDomainRequest{
        Name:        "test-domain",
        Description: "Test description",
    }
    
    _, err := useCase.Execute(context.Background(), req)
    assert.NoError(t, err)
}
```

### 3.3 계약 테스트 (10% - Contract Tests)

#### Pact를 사용한 계약 테스트
```go
// internal/interface/mcp/contract_test.go
func TestMCPService_Contract(t *testing.T) {
    pact := &dsl.Pact{
        Consumer: "url-db-client",
        Provider: "url-db-mcp-server",
    }
    defer pact.Teardown()
    
    pact.
        AddInteraction().
        Given("A domain exists").
        UponReceiving("A request to create a domain").
        WithRequest(dsl.Request{
            Method: "POST",
            Path:   "/mcp/domains",
            Headers: dsl.MapMatcher{
                "Content-Type": dsl.String("application/json"),
            },
            Body: map[string]interface{}{
                "name":        "test-domain",
                "description": "Test domain",
            },
        }).
        WillRespondWith(dsl.Response{
            Status: 201,
            Headers: dsl.MapMatcher{
                "Content-Type": dsl.String("application/json"),
            },
            Body: map[string]interface{}{
                "name":        "test-domain",
                "description": "Test domain",
                "created_at":  dsl.Like("2023-01-01T00:00:00Z"),
            },
        })
    
    err := pact.Verify(func() error {
        client := mcp.NewClient("http://localhost:8080")
        domain, err := client.CreateDomain("test-domain", "Test domain")
        if err != nil {
            return err
        }
        
        if domain.Name != "test-domain" {
            t.Errorf("Expected domain name 'test-domain', got '%s'", domain.Name)
        }
        
        return nil
    })
    
    assert.NoError(t, err)
}
```

### 3.4 E2E 시나리오 테스트 (10% - End-to-End Scenario Tests)

#### 도메인 관리 워크플로우
```go
// tests/e2e/domain_workflow_test.go
func TestDomainManagementWorkflow(t *testing.T) {
    // 1. 도메인 생성
    t.Run("Create Domain", func(t *testing.T) {
        server := httptest.NewServer(http.NewRouter())
        defer server.Close()
        
        client := mcp.NewClient(server.URL)
        domain, err := client.CreateDomain("test-domain", "Test description")
        
        assert.NoError(t, err)
        assert.Equal(t, "test-domain", domain.Name)
    })
    
    // 2. 도메인에 노드 추가
    t.Run("Add Node to Domain", func(t *testing.T) {
        // 노드 생성 로직...
    })
    
    // 3. 노드에 속성 할당
    t.Run("Assign Attributes to Node", func(t *testing.T) {
        // 속성 할당 로직...
    })
    
    // 4. 노드 간 의존성 생성
    t.Run("Create Dependencies Between Nodes", func(t *testing.T) {
        // 의존성 생성 로직...
    })
    
    // 5. 순환 의존성 감지
    t.Run("Detect Circular Dependencies", func(t *testing.T) {
        // 순환 의존성 감지 로직...
    })
    
    // 6. 영향도 분석
    t.Run("Analyze Impact", func(t *testing.T) {
        // 영향도 분석 로직...
    })
    
    // 7. 도메인 삭제 (의존성 정리)
    t.Run("Delete Domain with Cleanup", func(t *testing.T) {
        // 도메인 삭제 로직...
    })
}
```

### 3.5 성능 및 부하 테스트 (5% - Performance & Load Tests)

#### 벤치마크 테스트
```go
// internal/domain/service/benchmark_test.go
func BenchmarkDependencyGraphService_DetectCycles(b *testing.B) {
    service := NewDependencyGraphService(nil, nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := service.DetectCycles(context.Background(), nil)
        if err != nil {
            b.Fatalf("Benchmark failed: %v", err)
        }
    }
}

func BenchmarkCompositeKey_Normalization(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := NewCompositeKey("Test Tool", "Test Domain", "123")
        if err != nil {
            b.Fatalf("Benchmark failed: %v", err)
        }
    }
}
```

#### 부하 테스트
```go
// tests/performance/load_test.go
func TestLoadTest(t *testing.T) {
    server := httptest.NewServer(http.NewRouter())
    defer server.Close()
    
    scenarios := []struct {
        name     string
        requests int
        duration time.Duration
    }{
        {"Low Load", 100, 10 * time.Second},
        {"Medium Load", 1000, 30 * time.Second},
        {"High Load", 10000, 60 * time.Second},
    }
    
    for _, scenario := range scenarios {
        t.Run(scenario.name, func(t *testing.T) {
            start := time.Now()
            
            results := make(chan error, scenario.requests)
            for i := 0; i < scenario.requests; i++ {
                go func() {
                    // 실제 요청 로직...
                    results <- nil
                }()
            }
            
            successCount := 0
            for i := 0; i < scenario.requests; i++ {
                if <-results == nil {
                    successCount++
                }
            }
            
            duration := time.Since(start)
            tps := float64(successCount) / duration.Seconds()
            
            t.Logf("Success Rate: %d/%d (%.2f%%)", successCount, scenario.requests, float64(successCount)/float64(scenario.requests)*100)
            t.Logf("Throughput: %.2f requests/second", tps)
            
            // 성능 기준 검증
            if float64(successCount)/float64(scenario.requests) < 0.95 {
                t.Errorf("Success rate below 95%%: %.2f%%", float64(successCount)/float64(scenario.requests)*100)
            }
            
            if tps < 100 {
                t.Errorf("Throughput below 100 req/s: %.2f", tps)
            }
        })
    }
}
```

## 4. CI/CD 자동화

### 4.1 테스트 파이프라인

```yaml
# .github/workflows/test.yml
name: Test Pipeline

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Install dependencies
      run: go mod download
    
    - name: Lint
      run: |
        go install honnef.co/go/tools/cmd/staticcheck@latest
        staticcheck ./...
    
    - name: Build with race detection
      run: go build -race ./...
    
    - name: Unit tests
      run: go test -race -coverprofile=unit.out ./...
    
    - name: Property/Fuzz tests
      run: |
        go test -run Fuzz -fuzztime=10s ./...
        go test -run Property ./...
    
    - name: Integration tests
      run: go test -tags=integration -race ./...
    
    - name: Contract tests
      run: |
        go install github.com/pact-foundation/pact-go/v2/cmd/pact@latest
        pact-go verify
    
    - name: E2E tests
      run: go test ./tests/e2e/...
    
    - name: Performance tests
      run: |
        go test -bench=. -benchmem ./...
        go test -bench=. -cpuprofile=cpu.out -memprofile=mem.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: unit.out
        flags: unit
        name: codecov-umbrella
    
    - name: Upload performance profiles
      uses: actions/upload-artifact@v3
      with:
        name: performance-profiles
        path: |
          cpu.out
          mem.out
```

### 4.2 커버리지 거버넌스

```go
// internal/test/coverage/coverage.go
type CoverageThresholds struct {
    Overall         float64
    Domain          float64
    Application     float64
    Infrastructure  float64
    Interface       float64
}

func ValidateCoverage(profilePath string, thresholds CoverageThresholds) error {
    // 커버리지 파일 파싱
    data, err := os.ReadFile(profilePath)
    if err != nil {
        return fmt.Errorf("failed to read coverage file: %v", err)
    }
    
    // 패키지별 커버리지 계산
    packages := parseCoverageData(string(data))
    
    // 임계값 검증
    for pkg, coverage := range packages {
        var threshold float64
        
        switch {
        case strings.Contains(pkg, "domain"):
            threshold = thresholds.Domain
        case strings.Contains(pkg, "application"):
            threshold = thresholds.Application
        case strings.Contains(pkg, "infrastructure"):
            threshold = thresholds.Infrastructure
        case strings.Contains(pkg, "interface"):
            threshold = thresholds.Interface
        default:
            threshold = thresholds.Overall
        }
        
        if coverage < threshold {
            return fmt.Errorf("package %s coverage %.2f%% below threshold %.2f%%", pkg, coverage, threshold)
        }
    }
    
    return nil
}
```

## 5. 테스트 환경 설정

### 5.1 테스트 데이터베이스

```go
// internal/test/database/setup.go
func SetupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }
    
    // 스키마 생성
    schema, err := os.ReadFile("../../schema.sql")
    if err != nil {
        t.Fatalf("Failed to read schema file: %v", err)
    }
    
    _, err = db.Exec(string(schema))
    if err != nil {
        t.Fatalf("Failed to create schema: %v", err)
    }
    
    return db
}
```

### 5.2 테스트 컨테이너

```go
// internal/test/container/setup.go
func SetupTestContainers(t *testing.T) (string, func()) {
    ctx := context.Background()
    
    // Redis 컨테이너 시작
    req := testcontainers.ContainerRequest{
        Image:        "redis:7-alpine",
        ExposedPorts: []string{"6379/tcp"},
        WaitingFor:   wait.ForLog("Ready to accept connections"),
    }
    
    redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("Failed to start Redis container: %v", err)
    }
    
    host, err := redisContainer.Host(ctx)
    if err != nil {
        t.Fatalf("Failed to get Redis host: %v", err)
    }
    
    port, err := redisContainer.MappedPort(ctx, "6379")
    if err != nil {
        t.Fatalf("Failed to get Redis port: %v", err)
    }
    
    cleanup := func() {
        redisContainer.Terminate(ctx)
    }
    
    return fmt.Sprintf("%s:%s", host, port.Port()), cleanup
}
```

## 6. 테스트 모범 사례

### 6.1 AAA 패턴

```go
func TestCreateDomain_AAA_Pattern(t *testing.T) {
    // Arrange - 테스트 데이터 준비
    factory := factory.NewDomainFactory()
    mockRepo := &MockDomainRepository{}
    useCase := domain.NewCreateDomainUseCase(mockRepo)
    
    req := &domain.CreateDomainRequest{
        Name:        "test-domain",
        Description: "Test description",
    }
    
    // Act - 테스트 실행
    result, err := useCase.Execute(context.Background(), req)
    
    // Assert - 결과 검증
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "test-domain", result.Name)
    assert.Equal(t, "Test description", result.Description)
}
```

### 6.2 테스트 헬퍼 함수

```go
// internal/test/helper/domain_helper.go
func NewTestDomain(name string) *entity.Domain {
    domain, _ := entity.NewDomain(name, "Test description")
    return domain
}

func NewTestNode(url string) *entity.Node {
    node, _ := entity.NewNode(url, "Test title")
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
```

### 6.3 Example 테스트

```go
// internal/domain/valueobject/composite_key_example_test.go
func ExampleNewCompositeKey() {
    key, err := NewCompositeKey("test-tool", "test-domain", "123")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Tool: %s\n", key.ToolName())
    fmt.Printf("Domain: %s\n", key.DomainName())
    fmt.Printf("ID: %s\n", key.ID())
    fmt.Printf("String: %s\n", key.String())
    
    // Output:
    // Tool: test-tool
    // Domain: test-domain
    // ID: 123
    // String: test-tool:test-domain:123
}
```

## 7. 성능 고려사항

### 7.1 테스트 실행 속도

- **단위 테스트**: < 1초
- **통합 테스트**: < 10초
- **E2E 테스트**: < 60초

### 7.2 메모리 사용량

- **모킹 사용**: 실제 객체 대신 가벼운 모의 객체
- **테스트 격리**: 각 테스트 후 정리
- **리소스 관리**: 데이터베이스 연결 풀 관리

### 7.3 병렬 실행

- **독립적인 테스트**: 순서에 의존하지 않는 테스트
- **공유 상태 최소화**: 전역 변수 사용 금지
- **테스트 데이터 격리**: 각 테스트별 독립적인 데이터

## 결론

이 테스트 전략을 통해 URL-DB 프로젝트는:

1. **체계적인 테스트 피라미드**: 6단계 구조로 기능과 신뢰성을 동시에 담보
2. **고급 테스트 기법**: 속성 기반, 퍼즈, 계약 테스트로 예상치 못한 버그 탐지
3. **동시성 안전**: 레이스 감지, 고루틴 누수 방지로 안정성 확보
4. **성능 보장**: 벤치마크, 부하 테스트로 성능 회귀 방지
5. **자동화된 품질 관리**: CI/CD 파이프라인으로 지속적인 품질 보장

이 전략을 단계적으로 적용하면 **테스트 작성 ROI가 극대화**되고, 배포 사고·성능 회귀 위험을 실질적으로 감소시킬 수 있습니다. 