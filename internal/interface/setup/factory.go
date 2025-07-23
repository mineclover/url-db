package setup

import (
	"database/sql"
	
	"url-db/internal/domain/repository"
	"url-db/internal/application/usecase/domain"
	"url-db/internal/application/usecase/node"
	"url-db/internal/application/usecase/attribute"
	sqliteRepo "url-db/internal/infrastructure/persistence/sqlite/repository"
)

// RepositoryFactory creates repository instances
type RepositoryFactory interface {
	CreateDomainRepository() repository.DomainRepository
	CreateNodeRepository() repository.NodeRepository
	CreateAttributeRepository() repository.AttributeRepository
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
	toolName string
}

// NewApplicationFactory creates a new application factory
func NewApplicationFactory(db *sql.DB, sqlxDB interface{}, toolName string) *ApplicationFactory {
	return &ApplicationFactory{
		db:       db,
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

	// Create use cases
	createDomainUC, listDomainsUC := f.CreateDomainUseCases(domainRepo)
	createNodeUC, listNodesUC := f.CreateNodeUseCases(nodeRepo, domainRepo)
	createAttributeUC, listAttributesUC := f.CreateAttributeUseCases(attributeRepo, domainRepo)

	return &CleanDependencies{
		// Repositories
		DomainRepo:    domainRepo,
		NodeRepo:      nodeRepo,
		AttributeRepo: attributeRepo,

		// Use Cases
		CreateDomainUC:    createDomainUC,
		ListDomainsUC:     listDomainsUC,
		CreateNodeUC:      createNodeUC,
		ListNodesUC:       listNodesUC,
		CreateAttributeUC: createAttributeUC,
		ListAttributesUC:  listAttributesUC,
	}
}

// CleanDependencies holds Clean Architecture dependencies
type CleanDependencies struct {
	// Repositories
	DomainRepo    repository.DomainRepository
	NodeRepo      repository.NodeRepository
	AttributeRepo repository.AttributeRepository

	// Use Cases
	CreateDomainUC    *domain.CreateDomainUseCase
	ListDomainsUC     *domain.ListDomainsUseCase
	CreateNodeUC      *node.CreateNodeUseCase
	ListNodesUC       *node.ListNodesUseCase
	CreateAttributeUC *attribute.CreateAttributeUseCase
	ListAttributesUC  *attribute.ListAttributesUseCase
}