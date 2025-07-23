# 고성능 Go MCP 서버 테스트 가이드

## 개요

URL-DB MCP 서버는 수만 URL 그래프를 실시간으로 전파·분석하는 핵심 컴포넌트입니다. 단순한 코드 커버리지 80%를 만족하는 수준으로는 장애·성능 사고를 예방하기 어렵습니다. 본 가이드는 기존 테스트 피라미드를 확장해 **경계·동시성·계약·성능·보안**까지 체계적으로 검증하는 실천 방법을 제시합니다.

## 1. 테스트 피라미드 재구성

최근 Go 생태계는 퍼즈·속성 기반 테스트와 컨테이너 통합 테스트 도구가 급속히 성숙했습니다. 이를 반영해 URL-DB MCP 서버는 다음과 같은 **6-단계 피라미드**를 권장합니다.

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

### 1.1 각 단계별 목적

| 단계 | 비율 | 목적 | 도구 |
|------|------|------|------|
| **Unit Tests** | 40% | 순수 함수, 값 객체, DTO | `testing`, `testify` |
| **Property-Based & Fuzz Tests** | 15% | 경계·인코딩 취약점 | `rapid`, `gopter`, Go 1.18+ fuzzing |
| **Integration Tests** | 20% | gRPC bufconn, Testcontainers | `bufconn`, `testcontainers-go` |
| **Contract Tests** | 10% | Pact (Producer·Consumer) | `pact-go` |
| **End-to-End Scenario Tests** | 10% | 실제 MCP workflow | `httptest`, `grpctest` |
| **Performance & Load Tests** | 5% | benchmark + pprof | `testing.B`, `pprof` |

위 계층은 **기능 완전성(1–4단계)과 시스템 신뢰성(5–6단계)**을 동시에 담보하며, 레이어별 실패 원인을 명확히 국소화합니다.

## 2. Unit Layer 강화 전략

### 2.1 속성 기반(Property-Based)·퍼즈 테스트

테이블 주도(Unit) 테스트만으로는 예상치 못한 입력으로 발생하는 패닉·무한루프를 잡기 어렵습니다.

#### Rapid를 사용한 속성 기반 테스트

```go
// internal/domain/valueobject/composite_key_property_test.go
package valueobject_test

import (
    "testing"
    "pgregory.net/rapid"
    "url-db/internal/domain/valueobject"
)

func TestCompositeKey_PropertyBased(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        // 속성: 정규화된 복합키는 파싱 후 다시 문자열화하면 동일해야 함
        toolName := rapid.String().Draw(t, "toolName")
        domainName := rapid.String().Draw(t, "domainName")
        id := rapid.String().Draw(t, "id")
        
        key, err := valueobject.NewCompositeKey(toolName, domainName, id)
        if err != nil {
            t.Skip() // 유효하지 않은 입력은 스킵
        }
        
        // 속성 검증: 정규화 후 파싱하면 동일한 결과
        normalized := key.String()
        parsed, err := valueobject.ParseCompositeKey(normalized)
        
        if err != nil {
            t.Fatalf("Failed to parse normalized key: %v", err)
        }
        
        if !key.Equals(parsed) {
            t.Fatalf("Normalization property violated: %s != %s", key.String(), parsed.String())
        }
    })
}

func TestURL_PropertyBased(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        // 속성: URL 정규화는 멱등성이어야 함
        rawURL := rapid.String().Draw(t, "rawURL")
        
        url1, err1 := valueobject.NewURL(rawURL)
        if err1 != nil {
            t.Skip()
        }
        
        url2, err2 := valueobject.NewURL(url1.String())
        if err2 != nil {
            t.Fatalf("Failed to create URL from normalized string: %v", err2)
        }
        
        if !url1.Equals(url2) {
            t.Fatalf("URL normalization is not idempotent: %s != %s", url1.String(), url2.String())
        }
    })
}
```

#### Go 1.18+ 내장 Fuzzing

```go
// internal/domain/valueobject/composite_key_fuzz_test.go
package valueobject_test

import (
    "testing"
    "url-db/internal/domain/valueobject"
)

func FuzzCompositeKeyParse(f *testing.F) {
    // 시드 데이터 추가
    f.Add("tool:domain:123")
    f.Add("test-tool:test-domain:456")
    f.Add("invalid-key")
    
    f.Fuzz(func(t *testing.T, input string) {
        key, err := valueobject.ParseCompositeKey(input)
        
        // 속성: 파싱 성공 시 문자열화하면 원본과 일치하거나 정규화됨
        if err == nil {
            normalized := key.String()
            parsed, err2 := valueobject.ParseCompositeKey(normalized)
            
            if err2 != nil {
                t.Fatalf("Failed to parse normalized key: %v", err2)
            }
            
            if !key.Equals(parsed) {
                t.Fatalf("Parse-String roundtrip failed: %s != %s", key.String(), parsed.String())
            }
        }
    })
}

func FuzzURLValidation(f *testing.F) {
    f.Add("https://example.com")
    f.Add("http://localhost:8080/path?query=value")
    f.Add("invalid-url")
    
    f.Fuzz(func(t *testing.T, input string) {
        url, err := valueobject.NewURL(input)
        
        if err == nil {
            // 유효한 URL은 정규화 후에도 유효해야 함
            normalized := url.String()
            url2, err2 := valueobject.NewURL(normalized)
            
            if err2 != nil {
                t.Fatalf("Normalized URL is invalid: %v", err2)
            }
            
            if !url.Equals(url2) {
                t.Fatalf("URL normalization failed: %s != %s", url.String(), url2.String())
            }
        }
    })
}
```

### 2.2 경쟁 조건 탐지

#### Race Detector 활용

```go
// internal/application/usecase/domain/create_race_test.go
package domain_test

import (
    "context"
    "sync"
    "testing"
    "url-db/internal/application/usecase/domain"
)

func TestCreateDomainUseCase_RaceCondition(t *testing.T) {
    mockRepo := &MockDomainRepository{}
    useCase := domain.NewCreateDomainUseCase(mockRepo)
    
    // 동시에 같은 도메인명으로 생성 시도
    const numGoroutines = 10
    var wg sync.WaitGroup
    results := make(chan error, numGoroutines)
    
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            req := &domain.CreateDomainRequest{
                Name:        "race-test-domain",
                Description: "Race condition test",
            }
            
            _, err := useCase.Execute(context.Background(), req)
            results <- err
        }()
    }
    
    wg.Wait()
    close(results)
    
    // 하나만 성공하고 나머지는 실패해야 함
    successCount := 0
    for err := range results {
        if err == nil {
            successCount++
        }
    }
    
    if successCount != 1 {
        t.Errorf("Expected exactly 1 success, got %d", successCount)
    }
}
```

#### 정적 레이스 스캐너 (Chronos)

```bash
# CI에서 정적 레이스 스캔
go install github.com/amit-davidson/Chronos/cmd/chronos@latest
chronos ./...
```

### 2.3 고루틴 누수 방지

```go
// internal/application/usecase/domain/create_leak_test.go
package domain_test

import (
    "testing"
    "go.uber.org/goleak"
    "url-db/internal/application/usecase/domain"
)

func TestCreateDomainUseCase_NoGoroutineLeak(t *testing.T) {
    defer goleak.VerifyNone(t)
    
    mockRepo := &MockDomainRepository{}
    useCase := domain.NewCreateDomainUseCase(mockRepo)
    
    req := &domain.CreateDomainRequest{
        Name:        "test-domain",
        Description: "Test description",
    }
    
    _, err := useCase.Execute(context.Background(), req)
    if err != nil {
        t.Fatalf("Failed to create domain: %v", err)
    }
}

func TestDependencyGraphService_NoGoroutineLeak(t *testing.T) {
    defer goleak.VerifyNone(t)
    
    mockRepo := &MockDependencyRepository{}
    service := domain.NewDependencyGraphService(mockRepo, nil)
    
    // 복잡한 그래프 연산 수행
    cycles, err := service.DetectCycles(context.Background(), nil)
    if err != nil {
        t.Fatalf("Failed to detect cycles: %v", err)
    }
    
    _ = cycles // 결과 사용
}
```

### 2.4 뮤테이션(돌연변이) 테스트

```bash
# 뮤테이션 테스트 실행
go install github.com/avito-tech/go-mutesting/...@latest
go-mutesting ./internal/domain/...
```

## 3. Integration Layer 고도화

### 3.1 gRPC 서버 – `bufconn` 인-메모리 테스트

```go
// internal/interface/mcp/grpc_integration_test.go
package mcp_test

import (
    "context"
    "net"
    "testing"
    "google.golang.org/grpc"
    "google.golang.org/grpc/test/bufconn"
    "google.golang.org/grpc/credentials/insecure"
    "url-db/internal/interface/mcp"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
    lis = bufconn.Listen(bufSize)
    s := grpc.NewServer()
    mcp.RegisterMCPServiceServer(s, &mcp.Server{})
    go func() {
        if err := s.Serve(lis); err != nil {
            panic(err)
        }
    }()
}

func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

func TestMCPService_Integration(t *testing.T) {
    ctx := context.Background()
    
    conn, err := grpc.DialContext(ctx, "bufnet",
        grpc.WithContextDialer(bufDialer),
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err)
    }
    defer conn.Close()
    
    client := mcp.NewMCPServiceClient(conn)
    
    // 실제 gRPC 호출 테스트
    resp, err := client.CreateDomain(ctx, &mcp.CreateDomainRequest{
        Name:        "test-domain",
        Description: "Integration test domain",
    })
    
    if err != nil {
        t.Fatalf("Failed to create domain: %v", err)
    }
    
    if resp.Name != "test-domain" {
        t.Errorf("Expected domain name 'test-domain', got '%s'", resp.Name)
    }
}
```

### 3.2 Testcontainers 도입

```go
// internal/infrastructure/persistence/sqlite/integration_test.go
package sqlite_test

import (
    "context"
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    "url-db/internal/infrastructure/persistence/sqlite"
)

func TestSQLiteRepository_WithTestcontainers(t *testing.T) {
    ctx := context.Background()
    
    // SQLite 컨테이너 시작
    req := testcontainers.ContainerRequest{
        Image:        "alpine:latest",
        ExposedPorts: []string{"8080/tcp"},
        WaitingFor:   wait.ForLog("Ready"),
        Cmd:          []string{"sh", "-c", "apk add --no-cache sqlite && sqlite3 /tmp/test.db 'CREATE TABLE domains (id INTEGER PRIMARY KEY, name TEXT);' && tail -f /dev/null"},
    }
    
    sqliteContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("Failed to start container: %v", err)
    }
    defer sqliteContainer.Terminate(ctx)
    
    // 컨테이너 내부에서 SQLite 파일 생성
    exec, _, err := sqliteContainer.Exec(ctx, []string{"sh", "-c", "sqlite3 /tmp/test.db 'CREATE TABLE domains (id INTEGER PRIMARY KEY, name TEXT);'"})
    if err != nil {
        t.Fatalf("Failed to create table: %v", err)
    }
    
    if exec != 0 {
        t.Fatalf("Table creation failed with exit code: %d", exec)
    }
    
    // 테스트 로직...
}
```

### 3.3 소비자 주도 계약(Contract) 테스트

```go
// internal/interface/mcp/contract_test.go
package mcp_test

import (
    "testing"
    "github.com/pact-foundation/pact-go/dsl"
    "url-db/internal/interface/mcp"
)

func TestMCPService_Contract(t *testing.T) {
    // Pact 설정
    pact := &dsl.Pact{
        Consumer: "url-db-client",
        Provider: "url-db-mcp-server",
    }
    defer pact.Teardown()
    
    // 소비자 테스트
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
    
    // 계약 검증
    err := pact.Verify(func() error {
        // 실제 서비스 호출
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
    
    if err != nil {
        t.Fatalf("Contract verification failed: %v", err)
    }
}
```

## 4. 전-시스템(E2E)·성능 검증

### 4.1 시나리오 드리븐 E2E

```go
// tests/e2e/domain_workflow_test.go
package e2e_test

import (
    "context"
    "testing"
    "net/http/httptest"
    "url-db/internal/interface/http"
    "url-db/internal/interface/mcp"
)

func TestDomainManagementWorkflow(t *testing.T) {
    // 1. 도메인 생성
    t.Run("Create Domain", func(t *testing.T) {
        server := httptest.NewServer(http.NewRouter())
        defer server.Close()
        
        client := mcp.NewClient(server.URL)
        domain, err := client.CreateDomain("test-domain", "Test description")
        
        if err != nil {
            t.Fatalf("Failed to create domain: %v", err)
        }
        
        if domain.Name != "test-domain" {
            t.Errorf("Expected domain name 'test-domain', got '%s'", domain.Name)
        }
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

### 4.2 벤치마크·프로파일

```go
// internal/domain/service/benchmark_test.go
package service_test

import (
    "context"
    "testing"
    "url-db/internal/domain/service"
)

func BenchmarkDependencyGraphService_DetectCycles(b *testing.B) {
    service := service.NewDependencyGraphService(nil, nil)
    
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
        _, err := valueobject.NewCompositeKey("Test Tool", "Test Domain", "123")
        if err != nil {
            b.Fatalf("Benchmark failed: %v", err)
        }
    }
}

// CPU 프로파일링
func BenchmarkWithProfiling(b *testing.B) {
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        // 벤치마크 로직...
    }
}
```

### 4.3 부하·스트레스 테스트

```go
// tests/performance/load_test.go
package performance_test

import (
    "context"
    "testing"
    "time"
    "url-db/internal/interface/http"
)

func TestLoadTest(t *testing.T) {
    server := httptest.NewServer(http.NewRouter())
    defer server.Close()
    
    // 부하 테스트 시나리오
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
            
            // 동시 요청 생성
            results := make(chan error, scenario.requests)
            for i := 0; i < scenario.requests; i++ {
                go func() {
                    // 실제 요청 로직...
                    results <- nil
                }()
            }
            
            // 결과 수집
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

## 5. CI/CD 자동화와 커버리지 거버넌스

### 5.1 CI/CD 파이프라인

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

### 5.2 커버리지 거버넌스

```go
// internal/test/coverage/coverage.go
package coverage

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type CoverageThresholds struct {
    Overall     float64
    Domain      float64
    Application float64
    Infrastructure float64
    Interface   float64
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

## 6. 동시성·컨텍스트·시간 의존성 테스트

### 6.1 컨텍스트 취소

```go
// internal/application/usecase/domain/context_test.go
package domain_test

import (
    "context"
    "testing"
    "time"
    "url-db/internal/application/usecase/domain"
)

func TestCreateDomainUseCase_ContextCancellation(t *testing.T) {
    mockRepo := &MockDomainRepository{}
    useCase := domain.NewCreateDomainUseCase(mockRepo)
    
    // 컨텍스트 취소 시뮬레이션
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // 즉시 취소
    cancel()
    
    req := &domain.CreateDomainRequest{
        Name:        "test-domain",
        Description: "Test description",
    }
    
    _, err := useCase.Execute(ctx, req)
    
    if err == nil {
        t.Error("Expected error due to context cancellation")
    }
}

func TestDependencyGraphService_ContextTimeout(t *testing.T) {
    mockRepo := &MockDependencyRepository{}
    service := domain.NewDependencyGraphService(mockRepo, nil)
    
    // 짧은 타임아웃 설정
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer cancel()
    
    // 시간이 오래 걸리는 연산
    _, err := service.DetectCycles(ctx, nil)
    
    if err == nil {
        t.Error("Expected timeout error")
    }
}
```

### 6.2 시계 주입

```go
// internal/domain/entity/clock_test.go
package entity_test

import (
    "testing"
    "time"
    "url-db/internal/domain/entity"
)

type MockClock struct {
    now time.Time
}

func (c *MockClock) Now() time.Time {
    return c.now
}

func TestDomain_WithMockClock(t *testing.T) {
    mockTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    clock := &MockClock{now: mockTime}
    
    domain, err := entity.NewDomainWithClock("test-domain", "Test description", clock)
    if err != nil {
        t.Fatalf("Failed to create domain: %v", err)
    }
    
    if !domain.CreatedAt().Equal(mockTime) {
        t.Errorf("Expected creation time %v, got %v", mockTime, domain.CreatedAt())
    }
}
```

### 6.3 패턴 기반 동시성 테스트

```go
// internal/domain/service/concurrency_test.go
package service_test

import (
    "context"
    "sync"
    "testing"
    "time"
    "url-db/internal/domain/service"
)

func TestDependencyGraphService_ConcurrentAccess(t *testing.T) {
    service := service.NewDependencyGraphService(nil, nil)
    
    const numGoroutines = 10
    var wg sync.WaitGroup
    results := make(chan error, numGoroutines)
    
    // 동시에 여러 그래프 연산 수행
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            
            _, err := service.DetectCycles(ctx, nil)
            results <- err
        }(i)
    }
    
    wg.Wait()
    close(results)
    
    // 모든 고루틴이 성공적으로 완료되었는지 확인
    for err := range results {
        if err != nil {
            t.Errorf("Concurrent operation failed: %v", err)
        }
    }
}

func TestWorkerPool_Pattern(t *testing.T) {
    const numWorkers = 5
    const numJobs = 100
    
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    
    // 워커 풀 시작
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for job := range jobs {
                // 작업 처리
                results <- job * 2
            }
        }(i)
    }
    
    // 작업 전송
    for i := 0; i < numJobs; i++ {
        jobs <- i
    }
    close(jobs)
    
    // 워커 완료 대기
    wg.Wait()
    close(results)
    
    // 결과 검증
    resultCount := 0
    for result := range results {
        if result%2 != 0 {
            t.Errorf("Expected even result, got %d", result)
        }
        resultCount++
    }
    
    if resultCount != numJobs {
        t.Errorf("Expected %d results, got %d", numJobs, resultCount)
    }
}
```

## 7. 유지보수 가능한 테스트 코드를 위한 패턴

### 7.1 AAA 패턴과 팩토리 헬퍼

```go
// internal/test/factory/domain_factory.go
package factory

import (
    "url-db/internal/domain/entity"
    "url-db/internal/domain/valueobject"
)

type DomainFactory struct{}

func NewDomainFactory() *DomainFactory {
    return &DomainFactory{}
}

func (f *DomainFactory) CreateValidDomain(name string) *entity.Domain {
    domain, _ := entity.NewDomain(name, "Test description")
    return domain
}

func (f *DomainFactory) CreateDomainWithAttributes(name string, attributes []*entity.Attribute) *entity.Domain {
    domain, _ := entity.NewDomain(name, "Test description")
    // 속성 추가 로직...
    return domain
}

// internal/application/usecase/domain/create_aaa_test.go
func TestCreateDomainUseCase_AAA_Pattern(t *testing.T) {
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

### 7.2 Mock Repository + GoMock

```go
// internal/test/mock/repository_mock.go
package mock

import (
    "context"
    "github.com/golang/mock/gomock"
    "url-db/internal/domain/entity"
)

//go:generate mockgen -destination=domain_repository_mock.go -package=mock url-db/internal/domain/repository DomainRepository

func NewMockDomainRepository(ctrl *gomock.Controller) *MockDomainRepository {
    return NewMockDomainRepository(ctrl)
}

// internal/application/usecase/domain/create_mock_test.go
func TestCreateDomainUseCase_WithGoMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := NewMockDomainRepository(ctrl)
    useCase := domain.NewCreateDomainUseCase(mockRepo)
    
    // Mock 설정
    mockRepo.EXPECT().
        GetByName(gomock.Any(), "test-domain").
        Return(nil, nil)
    
    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil)
    
    req := &domain.CreateDomainRequest{
        Name:        "test-domain",
        Description: "Test description",
    }
    
    result, err := useCase.Execute(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 7.3 Example 테스트

```go
// internal/domain/valueobject/composite_key_example_test.go
package valueobject_test

import (
    "fmt"
    "url-db/internal/domain/valueobject"
)

func ExampleNewCompositeKey() {
    key, err := valueobject.NewCompositeKey("test-tool", "test-domain", "123")
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

func ExampleParseCompositeKey() {
    key, err := valueobject.ParseCompositeKey("test-tool:test-domain:123")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Parsed tool: %s\n", key.ToolName())
    fmt.Printf("Parsed domain: %s\n", key.DomainName())
    fmt.Printf("Parsed ID: %s\n", key.ID())
    
    // Output:
    // Parsed tool: test-tool
    // Parsed domain: test-domain
    // Parsed ID: 123
}
```

## 결론

이 고급 테스트 가이드를 통해 URL-DB MCP 서버는:

1. **Layer 확장**: 속성·퍼즈·계약·성능 층을 피라미드에 통합함으로써 기능·안정성 양립 구조를 확보
2. **신규 도구 스택**: Rapid, bufconn, Pact, goleak, Testcontainers 등 **표준 라이브러리 친화 도구**를 활용해 네트워크·동시성 테스트를 경량화
3. **자동 게이트**: 커버리지·레이스·고루틴 누수·벤치마크 한계치를 CI에서 차단하여 **"깨끗한 메인 브랜치"**를 지속적으로 유지
4. **동시성 안전**: 컨텍스트·시계 주입·레이스 디텍터를 조합해 MCP 서버 특유의 고병렬 환경에서도 **데이터 정합성과 자원 누수를 예방**

위 가이드를 URL-DB 프로젝트에 단계적으로 적용하면 **테스트 작성 ROI가 극대화**되고, 배포 사고·성능 회귀 위험을 실질적으로 감소시킬 수 있습니다. 