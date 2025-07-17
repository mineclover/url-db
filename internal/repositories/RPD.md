# 리포지토리 모듈 RPD

## 참조 문서
- [schema.sql](../../schema.sql) - 데이터베이스 스키마
- [internal/database/RPD.md](../database/RPD.md) - 데이터베이스 모듈

## 요구사항 분석

### 기능 요구사항
1. **도메인 리포지토리**: 도메인 CRUD 및 검색
2. **노드 리포지토리**: 노드 CRUD, 검색, 페이지네이션
3. **속성 리포지토리**: 속성 CRUD 및 도메인별 조회
4. **노드 속성 리포지토리**: 노드 속성 CRUD, 조인 쿼리
5. **트랜잭션 지원**: 모든 리포지토리에서 트랜잭션 처리
6. **성능 최적화**: 인덱스 활용, 쿼리 최적화

### 비기능 요구사항
- 인터페이스 기반 설계
- 테스트 가능한 구조
- 트랜잭션 안전성
- 에러 처리 및 로깅
- 성능 모니터링

## 아키텍처 설계

### 계층 구조
```
Service -> Repository -> Database -> SQLite
```

### 리포지토리 인터페이스
1. **DomainRepository**: 도메인 데이터 접근
2. **NodeRepository**: 노드 데이터 접근
3. **AttributeRepository**: 속성 데이터 접근
4. **NodeAttributeRepository**: 노드 속성 데이터 접근

## 구현 계획

### Phase 1: Interface Definition
- [ ] 각 리포지토리 인터페이스 정의
- [ ] 공통 베이스 인터페이스 정의
- [ ] 에러 타입 정의

### Phase 2: SQLite Implementation
- [ ] SQLite 기반 구현체 작성
- [ ] 쿼리 최적화
- [ ] 트랜잭션 처리
- [ ] 단위 테스트 작성

### Phase 3: Advanced Features
- [ ] 페이지네이션 구현
- [ ] 검색 기능 구현
- [ ] 배치 처리 구현
- [ ] 성능 테스트

### Phase 4: Integration
- [ ] 서비스 계층 통합
- [ ] 통합 테스트
- [ ] 성능 벤치마크

## 리포지토리 인터페이스

### DomainRepository
```go
type DomainRepository interface {
    Create(domain *models.Domain) error
    GetByID(id int) (*models.Domain, error)
    GetByName(name string) (*models.Domain, error)
    List(offset, limit int) ([]models.Domain, int, error)
    Update(domain *models.Domain) error
    Delete(id int) error
    ExistsByName(name string) (bool, error)
}
```

### NodeRepository
```go
type NodeRepository interface {
    Create(node *models.Node) error
    GetByID(id int) (*models.Node, error)
    GetByDomainAndContent(domainID int, content string) (*models.Node, error)
    ListByDomain(domainID int, offset, limit int) ([]models.Node, int, error)
    Search(domainID int, query string, offset, limit int) ([]models.Node, int, error)
    Update(node *models.Node) error
    Delete(id int) error
    ExistsByDomainAndContent(domainID int, content string) (bool, error)
}
```

### AttributeRepository
```go
type AttributeRepository interface {
    Create(attribute *models.Attribute) error
    GetByID(id int) (*models.Attribute, error)
    GetByDomainAndName(domainID int, name string) (*models.Attribute, error)
    ListByDomain(domainID int) ([]models.Attribute, error)
    Update(attribute *models.Attribute) error
    Delete(id int) error
    ExistsByDomainAndName(domainID int, name string) (bool, error)
}
```

### NodeAttributeRepository
```go
type NodeAttributeRepository interface {
    Create(nodeAttribute *models.NodeAttribute) error
    GetByID(id int) (*models.NodeAttribute, error)
    GetByNodeAndAttribute(nodeID, attributeID int) (*models.NodeAttribute, error)
    ListByNode(nodeID int) ([]models.NodeAttributeWithInfo, error)
    ListByAttribute(attributeID int) ([]models.NodeAttribute, error)
    Update(nodeAttribute *models.NodeAttribute) error
    Delete(id int) error
    DeleteByNode(nodeID int) error
    DeleteByAttribute(attributeID int) error
    ExistsByNodeAndAttribute(nodeID, attributeID int) (bool, error)
}
```

## SQLite 구현

### 도메인 리포지토리 구현
```go
type sqliteDomainRepository struct {
    db *sql.DB
}

func (r *sqliteDomainRepository) Create(domain *models.Domain) error {
    query := `
        INSERT INTO domains (name, description, created_at, updated_at)
        VALUES (?, ?, datetime('now'), datetime('now'))
        RETURNING id, created_at, updated_at
    `
    
    return r.db.QueryRow(query, domain.Name, domain.Description).Scan(
        &domain.ID, &domain.CreatedAt, &domain.UpdatedAt,
    )
}

func (r *sqliteDomainRepository) GetByID(id int) (*models.Domain, error) {
    query := `
        SELECT id, name, description, created_at, updated_at
        FROM domains
        WHERE id = ?
    `
    
    domain := &models.Domain{}
    err := r.db.QueryRow(query, id).Scan(
        &domain.ID, &domain.Name, &domain.Description,
        &domain.CreatedAt, &domain.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, ErrDomainNotFound
    }
    
    return domain, err
}
```

### 노드 리포지토리 구현
```go
type sqliteNodeRepository struct {
    db *sql.DB
}

func (r *sqliteNodeRepository) Create(node *models.Node) error {
    query := `
        INSERT INTO nodes (content, domain_id, title, description, created_at, updated_at)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
        RETURNING id, created_at, updated_at
    `
    
    return r.db.QueryRow(query, node.Content, node.DomainID, 
        node.Title, node.Description).Scan(
        &node.ID, &node.CreatedAt, &node.UpdatedAt,
    )
}

func (r *sqliteNodeRepository) Search(domainID int, query string, offset, limit int) ([]models.Node, int, error) {
    searchQuery := `
        SELECT id, content, domain_id, title, description, created_at, updated_at
        FROM nodes
        WHERE domain_id = ? AND (title LIKE ? OR content LIKE ?)
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `
    
    searchPattern := "%" + query + "%"
    rows, err := r.db.Query(searchQuery, domainID, searchPattern, searchPattern, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    var nodes []models.Node
    for rows.Next() {
        var node models.Node
        err := rows.Scan(&node.ID, &node.Content, &node.DomainID,
            &node.Title, &node.Description, &node.CreatedAt, &node.UpdatedAt)
        if err != nil {
            return nil, 0, err
        }
        nodes = append(nodes, node)
    }
    
    // 총 개수 조회
    countQuery := `
        SELECT COUNT(*)
        FROM nodes
        WHERE domain_id = ? AND (title LIKE ? OR content LIKE ?)
    `
    
    var total int
    err = r.db.QueryRow(countQuery, domainID, searchPattern, searchPattern).Scan(&total)
    if err != nil {
        return nil, 0, err
    }
    
    return nodes, total, nil
}
```

### 조인 쿼리 구현
```go
func (r *sqliteNodeAttributeRepository) ListByNode(nodeID int) ([]models.NodeAttributeWithInfo, error) {
    query := `
        SELECT na.id, na.node_id, na.attribute_id, na.value, na.order_index, na.created_at,
               a.name, a.type
        FROM node_attributes na
        JOIN attributes a ON na.attribute_id = a.id
        WHERE na.node_id = ?
        ORDER BY a.name, na.order_index
    `
    
    rows, err := r.db.Query(query, nodeID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var attributes []models.NodeAttributeWithInfo
    for rows.Next() {
        var attr models.NodeAttributeWithInfo
        err := rows.Scan(&attr.ID, &attr.NodeID, &attr.AttributeID,
            &attr.Value, &attr.OrderIndex, &attr.CreatedAt,
            &attr.Name, &attr.Type)
        if err != nil {
            return nil, err
        }
        attributes = append(attributes, attr)
    }
    
    return attributes, nil
}
```

## 트랜잭션 처리

### 트랜잭션 인터페이스
```go
type Transactional interface {
    WithTransaction(fn func(tx *sql.Tx) error) error
}

type transactionalRepository struct {
    db *sql.DB
}

func (r *transactionalRepository) WithTransaction(fn func(tx *sql.Tx) error) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()
    
    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit()
}
```

### 트랜잭션 사용 예제
```go
func (s *domainService) CreateDomainWithAttributes(domain *models.Domain, attributes []models.Attribute) error {
    return s.domainRepo.WithTransaction(func(tx *sql.Tx) error {
        // 도메인 생성
        if err := s.domainRepo.CreateTx(tx, domain); err != nil {
            return err
        }
        
        // 속성들 생성
        for _, attr := range attributes {
            attr.DomainID = domain.ID
            if err := s.attributeRepo.CreateTx(tx, &attr); err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

## 에러 처리

### 리포지토리 에러 정의
```go
var (
    ErrDomainNotFound      = errors.New("domain not found")
    ErrNodeNotFound        = errors.New("node not found")
    ErrAttributeNotFound   = errors.New("attribute not found")
    ErrNodeAttributeNotFound = errors.New("node attribute not found")
    ErrDuplicateEntry      = errors.New("duplicate entry")
    ErrForeignKeyConstraint = errors.New("foreign key constraint")
)
```

### SQLite 에러 매핑
```go
func mapSQLiteError(err error) error {
    if err == nil {
        return nil
    }
    
    if err == sql.ErrNoRows {
        return ErrNotFound
    }
    
    if sqliteErr, ok := err.(sqlite3.Error); ok {
        switch sqliteErr.Code {
        case sqlite3.ErrConstraintUnique:
            return ErrDuplicateEntry
        case sqlite3.ErrConstraintForeignKey:
            return ErrForeignKeyConstraint
        }
    }
    
    return err
}
```

## 테스트 전략

### 단위 테스트
```go
func TestDomainRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := NewSQLiteDomainRepository(db)
    
    domain := &models.Domain{
        Name:        "test-domain",
        Description: "Test domain",
    }
    
    err := repo.Create(domain)
    assert.NoError(t, err)
    assert.NotZero(t, domain.ID)
    assert.NotZero(t, domain.CreatedAt)
    assert.NotZero(t, domain.UpdatedAt)
}

func TestDomainRepository_GetByID(t *testing.T) {
    db := setupTestDB(t)
    repo := NewSQLiteDomainRepository(db)
    
    // 테스트 데이터 생성
    domain := createTestDomain(t, repo)
    
    // 조회 테스트
    found, err := repo.GetByID(domain.ID)
    assert.NoError(t, err)
    assert.Equal(t, domain.Name, found.Name)
    assert.Equal(t, domain.Description, found.Description)
}
```

### 통합 테스트
```go
func TestNodeRepository_Integration(t *testing.T) {
    db := setupTestDB(t)
    domainRepo := NewSQLiteDomainRepository(db)
    nodeRepo := NewSQLiteNodeRepository(db)
    
    // 도메인 생성
    domain := createTestDomain(t, domainRepo)
    
    // 노드 생성
    node := &models.Node{
        Content:     "https://example.com",
        DomainID:    domain.ID,
        Title:       "Test Node",
        Description: "Test description",
    }
    
    err := nodeRepo.Create(node)
    assert.NoError(t, err)
    
    // 노드 조회
    found, err := nodeRepo.GetByID(node.ID)
    assert.NoError(t, err)
    assert.Equal(t, node.Content, found.Content)
}
```

## 파일 구조
```
internal/repositories/
├── RPD.md
├── interfaces.go          # 리포지토리 인터페이스
├── domain.go              # 도메인 리포지토리 구현
├── domain_test.go         # 도메인 리포지토리 테스트
├── node.go                # 노드 리포지토리 구현
├── node_test.go           # 노드 리포지토리 테스트
├── attribute.go           # 속성 리포지토리 구현
├── attribute_test.go      # 속성 리포지토리 테스트
├── node_attribute.go      # 노드 속성 리포지토리 구현
├── node_attribute_test.go # 노드 속성 리포지토리 테스트
├── base.go                # 공통 베이스 리포지토리
├── transaction.go         # 트랜잭션 헬퍼
├── transaction_test.go    # 트랜잭션 테스트
├── errors.go              # 리포지토리 에러 정의
├── testutils.go           # 테스트 유틸리티
└── benchmark_test.go      # 성능 벤치마크
```

## 성능 최적화

### 인덱스 활용
- 적절한 WHERE 절 사용
- 복합 인덱스 활용
- 인덱스 힌트 제공

### 쿼리 최적화
- 불필요한 조인 제거
- 적절한 LIMIT 사용
- 서브쿼리 최적화

### 배치 처리
```go
func (r *sqliteNodeRepository) BatchCreate(nodes []models.Node) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.Prepare(`
        INSERT INTO nodes (content, domain_id, title, description, created_at, updated_at)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, node := range nodes {
        _, err := stmt.Exec(node.Content, node.DomainID, node.Title, node.Description)
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}
```

## 의존성
- `database/sql`: 표준 데이터베이스 인터페이스
- `github.com/mattn/go-sqlite3`: SQLite 드라이버
- `internal/models`: 데이터 모델
- `internal/database`: 데이터베이스 연결
- `github.com/stretchr/testify`: 테스트 유틸리티

## 보안 고려사항

### SQL 인젝션 방지
- 프리페어드 스테이트먼트 사용
- 사용자 입력 검증
- 쿼리 파라미터 바인딩

### 접근 제어
- 리포지토리 레벨 권한 검증
- 도메인별 접근 제한
- 감사 로깅