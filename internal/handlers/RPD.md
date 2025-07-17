# HTTP 핸들러 모듈 RPD

## 참조 문서
- [docs/api/01-domain-api.md](../../docs/api/01-domain-api.md) - 도메인 API
- [docs/api/02-attribute-api.md](../../docs/api/02-attribute-api.md) - 속성 API
- [docs/api/03-url-api.md](../../docs/api/03-url-api.md) - 노드 API
- [docs/api/04-url-attribute-api.md](../../docs/api/04-url-attribute-api.md) - 노드 속성 API
- [docs/api/05-url-attribute-validation-api.md](../../docs/api/05-url-attribute-validation-api.md) - 검증 API
- [docs/api/06-mcp-api.md](../../docs/api/06-mcp-api.md) - MCP API
- [docs/spec/error-codes.md](../../docs/spec/error-codes.md) - 에러 코드

## 요구사항 분석

### 기능 요구사항
1. **HTTP 요청 처리**: 각 API 엔드포인트별 HTTP 요청 처리
2. **요청 검증**: 입력 데이터 검증 및 바인딩
3. **응답 생성**: JSON 응답 생성 및 상태 코드 설정
4. **에러 처리**: 에러 응답 표준화
5. **미들웨어**: 로깅, 인증, CORS 등
6. **라우팅**: RESTful API 라우팅 설정

### 비기능 요구사항
- 일관된 API 응답 형식
- 적절한 HTTP 상태 코드
- 에러 메시지 표준화
- 성능 모니터링
- 보안 헤더 설정

## 아키텍처 설계

### 계층 구조
```
HTTP Request -> Middleware -> Handler -> Service -> Repository -> Database
```

### 핸들러 구조
1. **DomainHandler**: 도메인 관리 API
2. **AttributeHandler**: 속성 관리 API
3. **NodeHandler**: 노드 관리 API
4. **NodeAttributeHandler**: 노드 속성 관리 API
5. **MCPHandler**: MCP 서버 API
6. **HealthHandler**: 헬스 체크 API

## 구현 계획

### Phase 1: Core Handlers
- [ ] 각 핸들러 구조체 정의
- [ ] 기본 CRUD 핸들러 구현
- [ ] 요청/응답 모델 바인딩
- [ ] 단위 테스트 작성

### Phase 2: Error Handling
- [ ] 에러 응답 표준화
- [ ] 상태 코드 매핑
- [ ] 에러 미들웨어 구현
- [ ] 에러 처리 테스트

### Phase 3: Middleware
- [ ] 로깅 미들웨어
- [ ] CORS 미들웨어
- [ ] 인증 미들웨어 (선택)
- [ ] 레이트 리미팅 (선택)

### Phase 4: Routing & Integration
- [ ] 라우터 설정
- [ ] 미들웨어 적용
- [ ] 통합 테스트
- [ ] API 문서화

## 핸들러 상세 설계

### DomainHandler
```go
type DomainHandler struct {
    domainService services.DomainService
}

func (h *DomainHandler) CreateDomain(c *gin.Context) {
    var req models.CreateDomainRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    domain, err := h.domainService.CreateDomain(&req)
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, domain)
}
```

### NodeHandler
```go
type NodeHandler struct {
    nodeService services.NodeService
}

func (h *NodeHandler) CreateNode(c *gin.Context) {
    domainID, err := strconv.Atoi(c.Param("domain_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain ID"})
        return
    }
    
    var req models.CreateNodeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    node, err := h.nodeService.CreateNode(domainID, &req)
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, node)
}
```

### MCPHandler
```go
type MCPHandler struct {
    mcpService services.MCPService
}

func (h *MCPHandler) CreateMCPNode(c *gin.Context) {
    var req models.CreateMCPNodeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    node, err := h.mcpService.CreateNode(&req)
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, node)
}
```

## 에러 처리

### 에러 응답 표준화
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

func (h *BaseHandler) handleError(c *gin.Context, err error) {
    switch e := err.(type) {
    case *ValidationError:
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error:   "VALIDATION_ERROR",
            Message: e.Message,
            Details: e.Details,
        })
    case *NotFoundError:
        c.JSON(http.StatusNotFound, ErrorResponse{
            Error:   "NOT_FOUND",
            Message: e.Message,
        })
    case *ConflictError:
        c.JSON(http.StatusConflict, ErrorResponse{
            Error:   "CONFLICT",
            Message: e.Message,
        })
    default:
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error:   "INTERNAL_SERVER_ERROR",
            Message: "An unexpected error occurred",
        })
    }
}
```

### 상태 코드 매핑
- 200: 성공
- 201: 생성 성공
- 204: 삭제 성공
- 400: 잘못된 요청
- 404: 리소스 없음
- 409: 충돌 (중복 등)
- 500: 서버 오류

## 미들웨어

### 로깅 미들웨어
```go
func LoggingMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC3339),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
    })
}
```

### CORS 미들웨어
```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

### 에러 복구 미들웨어
```go
func RecoveryMiddleware() gin.HandlerFunc {
    return gin.RecoveryWithWriter(os.Stdout, func(c *gin.Context, recovered interface{}) {
        if err, ok := recovered.(error); ok {
            log.Printf("Panic recovered: %s", err.Error())
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "INTERNAL_SERVER_ERROR",
            "message": "An unexpected error occurred",
        })
    })
}
```

## 라우팅

### 라우터 설정
```go
func SetupRoutes(r *gin.Engine, domainHandler *DomainHandler, nodeHandler *NodeHandler, 
    attributeHandler *AttributeHandler, mcpHandler *MCPHandler) {
    
    // API 그룹
    api := r.Group("/api")
    
    // 도메인 라우트
    domains := api.Group("/domains")
    {
        domains.POST("", domainHandler.CreateDomain)
        domains.GET("", domainHandler.GetDomains)
        domains.GET("/:id", domainHandler.GetDomain)
        domains.PUT("/:id", domainHandler.UpdateDomain)
        domains.DELETE("/:id", domainHandler.DeleteDomain)
    }
    
    // 노드 라우트
    nodes := api.Group("/urls")
    {
        nodes.GET("/:id", nodeHandler.GetNode)
        nodes.PUT("/:id", nodeHandler.UpdateNode)
        nodes.DELETE("/:id", nodeHandler.DeleteNode)
    }
    
    // 도메인별 노드 라우트
    domainNodes := domains.Group("/:domain_id/urls")
    {
        domainNodes.POST("", nodeHandler.CreateNode)
        domainNodes.GET("", nodeHandler.GetNodesByDomain)
        domainNodes.POST("/find", nodeHandler.FindNodeByURL)
    }
    
    // MCP 라우트
    mcp := api.Group("/mcp")
    {
        mcp.POST("/nodes", mcpHandler.CreateMCPNode)
        mcp.GET("/nodes", mcpHandler.GetMCPNodes)
        mcp.GET("/nodes/:composite_id", mcpHandler.GetMCPNode)
        mcp.PUT("/nodes/:composite_id", mcpHandler.UpdateMCPNode)
        mcp.DELETE("/nodes/:composite_id", mcpHandler.DeleteMCPNode)
        mcp.POST("/nodes/find", mcpHandler.FindMCPNodeByURL)
        mcp.POST("/nodes/batch", mcpHandler.BatchGetMCPNodes)
    }
}
```

## 테스트 전략

### 단위 테스트
- 각 핸들러 메서드 테스트
- 요청 바인딩 테스트
- 응답 생성 테스트
- 에러 처리 테스트

### 통합 테스트
- End-to-end API 테스트
- 미들웨어 통합 테스트
- 라우팅 테스트

### 테스트 헬퍼
```go
func TestHandler(t *testing.T) {
    // 테스트용 서비스 모킹
    mockService := &MockDomainService{}
    handler := NewDomainHandler(mockService)
    
    // 테스트용 라우터 설정
    r := gin.New()
    r.POST("/domains", handler.CreateDomain)
    
    // 테스트 요청 생성
    req := httptest.NewRequest("POST", "/domains", strings.NewReader(`{"name":"test"}`))
    req.Header.Set("Content-Type", "application/json")
    
    // 테스트 실행
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    
    // 결과 검증
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

## 파일 구조
```
internal/handlers/
├── RPD.md
├── domain_handler.go      # 도메인 핸들러
├── domain_handler_test.go # 도메인 핸들러 테스트
├── node_handler.go        # 노드 핸들러
├── node_handler_test.go   # 노드 핸들러 테스트
├── attribute_handler.go   # 속성 핸들러
├── attribute_handler_test.go # 속성 핸들러 테스트
├── node_attribute_handler.go # 노드 속성 핸들러
├── node_attribute_handler_test.go # 노드 속성 핸들러 테스트
├── mcp_handler.go         # MCP 핸들러
├── mcp_handler_test.go    # MCP 핸들러 테스트
├── health_handler.go      # 헬스 체크 핸들러
├── health_handler_test.go # 헬스 체크 핸들러 테스트
├── router.go              # 라우터 설정
├── router_test.go         # 라우터 테스트
├── middleware.go          # 미들웨어 모음
├── middleware_test.go     # 미들웨어 테스트
├── errors.go              # 에러 처리 헬퍼
├── errors_test.go         # 에러 처리 테스트
└── base_handler.go        # 공통 핸들러 기능
```

## 의존성
- `github.com/gin-gonic/gin`: HTTP 프레임워크
- `internal/services`: 서비스 계층
- `internal/models`: 데이터 모델
- `net/http`: HTTP 상태 코드
- `github.com/stretchr/testify`: 테스트 유틸리티

## 성능 고려사항

### 응답 시간 최적화
- 적절한 타임아웃 설정
- 스트리밍 응답 (대용량 데이터)
- 압축 미들웨어

### 메모리 관리
- 요청 크기 제한
- 적절한 버퍼 크기
- 메모리 누수 방지

### 보안
- 입력 검증 강화
- SQL 인젝션 방지
- XSS 방지
- CSRF 보호 (필요시)