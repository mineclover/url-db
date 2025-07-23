package repositories

import (
	"database/sql"
	"testing"
	"time"
	"url-db/internal/models"
	"url-db/internal/shared/testdb"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// TestDB 는 테스트용 데이터베이스 설정을 위한 구조체입니다.
type TestDB struct {
	DB *sql.DB
}

// SetupTestDB 는 테스트용 인메모리 SQLite 데이터베이스를 설정합니다.
func SetupTestDB(t *testing.T) *TestDB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// 테스트용 스키마 생성
	testdb.LoadSchema(t, db)

	return &TestDB{DB: db}
}

// Close 는 테스트 데이터베이스를 닫습니다.
func (tdb *TestDB) Close() {
	tdb.DB.Close()
}


// TestDomainBuilder 는 테스트용 도메인 빌더입니다.
type TestDomainBuilder struct {
	domain *models.Domain
}

// NewTestDomainBuilder 는 새로운 테스트 도메인 빌더를 생성합니다.
func NewTestDomainBuilder() *TestDomainBuilder {
	return &TestDomainBuilder{
		domain: &models.Domain{
			Name:        "test-domain",
			Description: "Test domain",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// WithName 은 도메인 이름을 설정합니다.
func (b *TestDomainBuilder) WithName(name string) *TestDomainBuilder {
	b.domain.Name = name
	return b
}

// WithDescription 은 도메인 설명을 설정합니다.
func (b *TestDomainBuilder) WithDescription(description string) *TestDomainBuilder {
	b.domain.Description = description
	return b
}

// WithID 는 도메인 ID를 설정합니다.
func (b *TestDomainBuilder) WithID(id int) *TestDomainBuilder {
	b.domain.ID = id
	return b
}

// Build 는 도메인을 빌드합니다.
func (b *TestDomainBuilder) Build() *models.Domain {
	return b.domain
}

// TestNodeBuilder 는 테스트용 노드 빌더입니다.
type TestNodeBuilder struct {
	node *models.Node
}

// NewTestNodeBuilder 는 새로운 테스트 노드 빌더를 생성합니다.
func NewTestNodeBuilder() *TestNodeBuilder {
	return &TestNodeBuilder{
		node: &models.Node{
			Content:     "https://example.com",
			DomainID:    1,
			Title:       "Test Node",
			Description: "Test node description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// WithContent 는 노드 콘텐츠를 설정합니다.
func (b *TestNodeBuilder) WithContent(content string) *TestNodeBuilder {
	b.node.Content = content
	return b
}

// WithDomainID 는 도메인 ID를 설정합니다.
func (b *TestNodeBuilder) WithDomainID(domainID int) *TestNodeBuilder {
	b.node.DomainID = domainID
	return b
}

// WithTitle 은 노드 제목을 설정합니다.
func (b *TestNodeBuilder) WithTitle(title string) *TestNodeBuilder {
	b.node.Title = title
	return b
}

// WithDescription 은 노드 설명을 설정합니다.
func (b *TestNodeBuilder) WithDescription(description string) *TestNodeBuilder {
	b.node.Description = description
	return b
}

// WithID 는 노드 ID를 설정합니다.
func (b *TestNodeBuilder) WithID(id int) *TestNodeBuilder {
	b.node.ID = id
	return b
}

// Build 는 노드를 빌드합니다.
func (b *TestNodeBuilder) Build() *models.Node {
	return b.node
}

// TestAttributeBuilder 는 테스트용 속성 빌더입니다.
type TestAttributeBuilder struct {
	attribute *models.Attribute
}

// NewTestAttributeBuilder 는 새로운 테스트 속성 빌더를 생성합니다.
func NewTestAttributeBuilder() *TestAttributeBuilder {
	return &TestAttributeBuilder{
		attribute: &models.Attribute{
			DomainID:    1,
			Name:        "test-attribute",
			Type:        "tag",
			Description: "Test attribute",
			CreatedAt:   time.Now(),
		},
	}
}

// WithDomainID 는 도메인 ID를 설정합니다.
func (b *TestAttributeBuilder) WithDomainID(domainID int) *TestAttributeBuilder {
	b.attribute.DomainID = domainID
	return b
}

// WithName 은 속성 이름을 설정합니다.
func (b *TestAttributeBuilder) WithName(name string) *TestAttributeBuilder {
	b.attribute.Name = name
	return b
}

// WithType 은 속성 타입을 설정합니다.
func (b *TestAttributeBuilder) WithType(attrType string) *TestAttributeBuilder {
	b.attribute.Type = models.AttributeType(attrType)
	return b
}

// WithDescription 은 속성 설명을 설정합니다.
func (b *TestAttributeBuilder) WithDescription(description string) *TestAttributeBuilder {
	b.attribute.Description = description
	return b
}

// WithID 는 속성 ID를 설정합니다.
func (b *TestAttributeBuilder) WithID(id int) *TestAttributeBuilder {
	b.attribute.ID = id
	return b
}

// Build 는 속성을 빌드합니다.
func (b *TestAttributeBuilder) Build() *models.Attribute {
	return b.attribute
}

// TestNodeAttributeBuilder 는 테스트용 노드 속성 빌더입니다.
type TestNodeAttributeBuilder struct {
	nodeAttribute *models.NodeAttribute
}

// NewTestNodeAttributeBuilder 는 새로운 테스트 노드 속성 빌더를 생성합니다.
func NewTestNodeAttributeBuilder() *TestNodeAttributeBuilder {
	return &TestNodeAttributeBuilder{
		nodeAttribute: &models.NodeAttribute{
			NodeID:      1,
			AttributeID: 1,
			Value:       "test-value",
			OrderIndex:  nil,
			CreatedAt:   time.Now(),
		},
	}
}

// WithNodeID 는 노드 ID를 설정합니다.
func (b *TestNodeAttributeBuilder) WithNodeID(nodeID int) *TestNodeAttributeBuilder {
	b.nodeAttribute.NodeID = nodeID
	return b
}

// WithAttributeID 는 속성 ID를 설정합니다.
func (b *TestNodeAttributeBuilder) WithAttributeID(attributeID int) *TestNodeAttributeBuilder {
	b.nodeAttribute.AttributeID = attributeID
	return b
}

// WithValue 는 속성 값을 설정합니다.
func (b *TestNodeAttributeBuilder) WithValue(value string) *TestNodeAttributeBuilder {
	b.nodeAttribute.Value = value
	return b
}

// WithOrderIndex 는 순서 인덱스를 설정합니다.
func (b *TestNodeAttributeBuilder) WithOrderIndex(orderIndex int) *TestNodeAttributeBuilder {
	b.nodeAttribute.OrderIndex = &orderIndex
	return b
}

// WithID 는 노드 속성 ID를 설정합니다.
func (b *TestNodeAttributeBuilder) WithID(id int) *TestNodeAttributeBuilder {
	b.nodeAttribute.ID = id
	return b
}

// Build 는 노드 속성을 빌드합니다.
func (b *TestNodeAttributeBuilder) Build() *models.NodeAttribute {
	return b.nodeAttribute
}

// CreateTestDomain 은 테스트용 도메인을 생성합니다.
func CreateTestDomain(t *testing.T, repo DomainRepository) *models.Domain {
	domain := NewTestDomainBuilder().Build()
	err := repo.Create(domain)
	require.NoError(t, err)
	return domain
}

// CreateTestNode 는 테스트용 노드를 생성합니다.
func CreateTestNode(t *testing.T, repo NodeRepository, domainID int) *models.Node {
	node := NewTestNodeBuilder().WithDomainID(domainID).Build()
	err := repo.Create(node)
	require.NoError(t, err)
	return node
}

// CreateTestAttribute 는 테스트용 속성을 생성합니다.
func CreateTestAttribute(t *testing.T, repo AttributeRepository, domainID int) *models.Attribute {
	attribute := NewTestAttributeBuilder().WithDomainID(domainID).Build()
	err := repo.Create(attribute)
	require.NoError(t, err)
	return attribute
}

// CreateTestNodeAttribute 는 테스트용 노드 속성을 생성합니다.
func CreateTestNodeAttribute(t *testing.T, repo NodeAttributeRepository, nodeID, attributeID int) *models.NodeAttribute {
	nodeAttribute := NewTestNodeAttributeBuilder().WithNodeID(nodeID).WithAttributeID(attributeID).Build()
	err := repo.Create(nodeAttribute)
	require.NoError(t, err)
	return nodeAttribute
}

// AssertDomainEqual 은 두 도메인이 같은지 확인합니다.
func AssertDomainEqual(t *testing.T, expected, actual *models.Domain) {
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Description, actual.Description)
	if expected.ID != 0 {
		require.Equal(t, expected.ID, actual.ID)
	}
}

// AssertNodeEqual 은 두 노드가 같은지 확인합니다.
func AssertNodeEqual(t *testing.T, expected, actual *models.Node) {
	require.Equal(t, expected.Content, actual.Content)
	require.Equal(t, expected.DomainID, actual.DomainID)
	require.Equal(t, expected.Title, actual.Title)
	require.Equal(t, expected.Description, actual.Description)
	if expected.ID != 0 {
		require.Equal(t, expected.ID, actual.ID)
	}
}

// AssertAttributeEqual 은 두 속성이 같은지 확인합니다.
func AssertAttributeEqual(t *testing.T, expected, actual *models.Attribute) {
	require.Equal(t, expected.DomainID, actual.DomainID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Type, actual.Type)
	require.Equal(t, expected.Description, actual.Description)
	if expected.ID != 0 {
		require.Equal(t, expected.ID, actual.ID)
	}
}

// AssertNodeAttributeEqual 은 두 노드 속성이 같은지 확인합니다.
func AssertNodeAttributeEqual(t *testing.T, expected, actual *models.NodeAttribute) {
	require.Equal(t, expected.NodeID, actual.NodeID)
	require.Equal(t, expected.AttributeID, actual.AttributeID)
	require.Equal(t, expected.Value, actual.Value)
	require.Equal(t, expected.OrderIndex, actual.OrderIndex)
	if expected.ID != 0 {
		require.Equal(t, expected.ID, actual.ID)
	}
}
