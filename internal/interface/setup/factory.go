package setup

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"url-db/internal/application/usecase/attribute"
	"url-db/internal/application/usecase/domain"
	"url-db/internal/application/usecase/node"
	"url-db/internal/domain/repository"
	domainAttribute "url-db/internal/domain/attribute"
	sqliteRepo "url-db/internal/infrastructure/persistence/sqlite/repository"
)

// RepositoryFactory creates repository instances
type RepositoryFactory interface {
	CreateDomainRepository() repository.DomainRepository
	CreateNodeRepository() repository.NodeRepository
	CreateAttributeRepository() repository.AttributeRepository
	CreateNodeAttributeRepository() repository.NodeAttributeRepository
}

// UseCaseFactory creates use case instances
type UseCaseFactory interface {
	CreateDomainUseCases(domainRepo repository.DomainRepository) (*domain.CreateDomainUseCase, *domain.ListDomainsUseCase)
	CreateNodeUseCases(nodeRepo repository.NodeRepository, domainRepo repository.DomainRepository) (*node.CreateNodeUseCase, *node.ListNodesUseCase)
	CreateAttributeUseCases(attributeRepo repository.AttributeRepository, domainRepo repository.DomainRepository) (*attribute.CreateAttributeUseCase, *attribute.ListAttributesUseCase)
}

// ApplicationFactory coordinates all factories
type ApplicationFactory struct {
	db       *sql.DB
	sqlxDB   *sqlx.DB
	toolName string
}

// NewApplicationFactory creates a new application factory
func NewApplicationFactory(db *sql.DB, sqlxDB *sqlx.DB, toolName string) *ApplicationFactory {
	return &ApplicationFactory{
		db:       db,
		sqlxDB:   sqlxDB,
		toolName: toolName,
	}
}

// Repository Factory Implementation
func (f *ApplicationFactory) CreateDomainRepository() repository.DomainRepository {
	return sqliteRepo.NewDomainRepository(f.db)
}

func (f *ApplicationFactory) CreateNodeRepository() repository.NodeRepository {
	return sqliteRepo.NewNodeRepository(f.db)
}

func (f *ApplicationFactory) CreateAttributeRepository() repository.AttributeRepository {
	return sqliteRepo.NewAttributeRepository(f.db)
}

func (f *ApplicationFactory) CreateNodeAttributeRepository() repository.NodeAttributeRepository {
	return sqliteRepo.NewSQLiteNodeAttributeRepository(f.sqlxDB)
}

// Use Case Factory Implementation
func (f *ApplicationFactory) CreateDomainUseCases(domainRepo repository.DomainRepository) (*domain.CreateDomainUseCase, *domain.ListDomainsUseCase) {
	createUC := domain.NewCreateDomainUseCase(domainRepo)
	listUC := domain.NewListDomainsUseCase(domainRepo)
	return createUC, listUC
}

func (f *ApplicationFactory) CreateNodeUseCases(nodeRepo repository.NodeRepository, domainRepo repository.DomainRepository) (*node.CreateNodeUseCase, *node.ListNodesUseCase) {
	createUC := node.NewCreateNodeUseCase(nodeRepo, domainRepo)
	listUC := node.NewListNodesUseCase(nodeRepo)
	return createUC, listUC
}

func (f *ApplicationFactory) CreateAttributeUseCases(attributeRepo repository.AttributeRepository, domainRepo repository.DomainRepository) (*attribute.CreateAttributeUseCase, *attribute.ListAttributesUseCase) {
	createUC := attribute.NewCreateAttributeUseCase(attributeRepo, domainRepo)
	listUC := attribute.NewListAttributesUseCase(attributeRepo, domainRepo)
	return createUC, listUC
}

// CreateCleanArchitectureDependencies creates all dependencies for Clean Architecture
func (f *ApplicationFactory) CreateCleanArchitectureDependencies() *CleanDependencies {
	// Create repositories
	domainRepo := f.CreateDomainRepository()
	nodeRepo := f.CreateNodeRepository()
	attributeRepo := f.CreateAttributeRepository()
	nodeAttributeRepo := f.CreateNodeAttributeRepository()

	// Create validation registry
	validatorRegistry := domainAttribute.NewValidatorRegistry()

	// Create use cases
	createDomainUC, listDomainsUC := f.CreateDomainUseCases(domainRepo)
	createNodeUC, listNodesUC := f.CreateNodeUseCases(nodeRepo, domainRepo)
	createAttributeUC, listAttributesUC := f.CreateAttributeUseCases(attributeRepo, domainRepo)
	setNodeAttributesUC := node.NewSetNodeAttributesUseCase(nodeRepo, attributeRepo, nodeAttributeRepo)

	return &CleanDependencies{
		// Repositories
		DomainRepo:        domainRepo,
		NodeRepo:          nodeRepo,
		AttributeRepo:     attributeRepo,
		NodeAttributeRepo: nodeAttributeRepo,

		// Validators
		ValidatorRegistry: validatorRegistry,

		// Use Cases
		CreateDomainUC:      createDomainUC,
		ListDomainsUC:       listDomainsUC,
		CreateNodeUC:        createNodeUC,
		ListNodesUC:         listNodesUC,
		CreateAttributeUC:   createAttributeUC,
		ListAttributesUC:    listAttributesUC,
		SetNodeAttributesUC: setNodeAttributesUC,
	}
}

// CleanDependencies holds Clean Architecture dependencies
type CleanDependencies struct {
	// Repositories
	DomainRepo        repository.DomainRepository
	NodeRepo          repository.NodeRepository
	AttributeRepo     repository.AttributeRepository
	NodeAttributeRepo repository.NodeAttributeRepository

	// Validators
	ValidatorRegistry *domainAttribute.ValidatorRegistry

	// Use Cases
	CreateDomainUC      *domain.CreateDomainUseCase
	ListDomainsUC       *domain.ListDomainsUseCase
	CreateNodeUC        *node.CreateNodeUseCase
	ListNodesUC         *node.ListNodesUseCase
	CreateAttributeUC   *attribute.CreateAttributeUseCase
	ListAttributesUC    *attribute.ListAttributesUseCase
	SetNodeAttributesUC *node.SetNodeAttributesUseCase
}
