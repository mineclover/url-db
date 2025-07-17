package repositories

import (
	"context"
	"database/sql"
	"url-db/internal/models"
)

// DomainRepository 는 도메인 데이터 접근을 위한 인터페이스입니다.
type DomainRepository interface {
	Create(domain *models.Domain) error
	GetByID(id int) (*models.Domain, error)
	GetByName(name string) (*models.Domain, error)
	List(offset, limit int) ([]models.Domain, int, error)
	Update(domain *models.Domain) error
	Delete(id int) error
	ExistsByName(name string) (bool, error)
	
	// 트랜잭션 지원 메서드
	CreateTx(tx *sql.Tx, domain *models.Domain) error
	UpdateTx(tx *sql.Tx, domain *models.Domain) error
	DeleteTx(tx *sql.Tx, id int) error
}

// NodeRepository 는 노드 데이터 접근을 위한 인터페이스입니다.
type NodeRepository interface {
	Create(node *models.Node) error
	GetByID(id int) (*models.Node, error)
	GetByDomainAndContent(domainID int, content string) (*models.Node, error)
	ListByDomain(domainID int, offset, limit int) ([]models.Node, int, error)
	Search(domainID int, query string, offset, limit int) ([]models.Node, int, error)
	Update(node *models.Node) error
	Delete(id int) error
	ExistsByDomainAndContent(domainID int, content string) (bool, error)
	CountNodesByDomain(ctx context.Context, domainID int) (int, error)
	
	// 배치 처리 메서드
	BatchCreate(nodes []models.Node) error
	BatchUpdate(nodes []models.Node) error
	BatchDelete(ids []int) error
	
	// 트랜잭션 지원 메서드
	CreateTx(tx *sql.Tx, node *models.Node) error
	UpdateTx(tx *sql.Tx, node *models.Node) error
	DeleteTx(tx *sql.Tx, id int) error
}

// AttributeRepository 는 속성 데이터 접근을 위한 인터페이스입니다.
type AttributeRepository interface {
	Create(attribute *models.Attribute) error
	GetByID(id int) (*models.Attribute, error)
	GetByDomainAndName(domainID int, name string) (*models.Attribute, error)
	ListByDomain(domainID int) ([]models.Attribute, error)
	Update(attribute *models.Attribute) error
	Delete(id int) error
	ExistsByDomainAndName(domainID int, name string) (bool, error)
	
	// 트랜잭션 지원 메서드
	CreateTx(tx *sql.Tx, attribute *models.Attribute) error
	UpdateTx(tx *sql.Tx, attribute *models.Attribute) error
	DeleteTx(tx *sql.Tx, id int) error
}

// NodeAttributeRepository 는 노드 속성 데이터 접근을 위한 인터페이스입니다.
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
	
	// 배치 처리 메서드
	BatchCreate(nodeAttributes []models.NodeAttribute) error
	BatchUpdate(nodeAttributes []models.NodeAttribute) error
	BatchDeleteByNode(nodeID int) error
	BatchDeleteByAttribute(attributeID int) error
	
	// 트랜잭션 지원 메서드
	CreateTx(tx *sql.Tx, nodeAttribute *models.NodeAttribute) error
	UpdateTx(tx *sql.Tx, nodeAttribute *models.NodeAttribute) error
	DeleteTx(tx *sql.Tx, id int) error
}

// NodeConnectionRepository 는 노드 연결 데이터 접근을 위한 인터페이스입니다.
type NodeConnectionRepository interface {
	Create(ctx context.Context, connection *models.NodeConnection) error
	GetByID(ctx context.Context, id int) (*models.NodeConnection, error)
	GetBySourceAndTarget(ctx context.Context, sourceNodeID, targetNodeID int, relationshipType string) (*models.NodeConnection, error)
	ListBySourceNode(ctx context.Context, sourceNodeID int, offset, limit int) ([]models.NodeConnectionWithInfo, int, error)
	ListByTargetNode(ctx context.Context, targetNodeID int, offset, limit int) ([]models.NodeConnectionWithInfo, int, error)
	ListByRelationshipType(ctx context.Context, relationshipType string, offset, limit int) ([]models.NodeConnectionWithInfo, int, error)
	Update(ctx context.Context, connection *models.NodeConnection) error
	Delete(ctx context.Context, id int) error
	DeleteBySourceNode(ctx context.Context, sourceNodeID int) error
	DeleteByTargetNode(ctx context.Context, targetNodeID int) error
	ExistsBySourceAndTarget(ctx context.Context, sourceNodeID, targetNodeID int, relationshipType string) (bool, error)
	
	// 배치 처리 메서드
	BatchCreate(ctx context.Context, connections []models.NodeConnection) error
	BatchDelete(ctx context.Context, ids []int) error
	
	// 트랜잭션 지원 메서드
	CreateTx(tx *sql.Tx, connection *models.NodeConnection) error
	UpdateTx(tx *sql.Tx, connection *models.NodeConnection) error
	DeleteTx(tx *sql.Tx, id int) error
}

// Transactional 은 트랜잭션을 지원하는 인터페이스입니다.
type Transactional interface {
	WithTransaction(fn func(tx *sql.Tx) error) error
}

// Repositories 는 모든 리포지토리를 포함하는 구조체입니다.
type Repositories struct {
	Domain        DomainRepository
	Node          NodeRepository
	Attribute     AttributeRepository
	NodeAttribute NodeAttributeRepository
}

// NewRepositories 는 새로운 리포지토리 컬렉션을 생성합니다.
func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Domain:        NewSQLiteDomainRepository(db),
		Node:          NewSQLiteNodeRepository(db),
		Attribute:     NewSQLiteAttributeRepository(db),
		NodeAttribute: NewSQLiteNodeAttributeRepository(db),
	}
}