package setup

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"url-db/internal/application/usecase/attribute"
	"url-db/internal/application/usecase/domain"
	"url-db/internal/application/usecase/node"
	domainAttribute "url-db/internal/domain/attribute"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
	"url-db/internal/domain/service"
	sqliteRepo "url-db/internal/infrastructure/persistence/sqlite/repository"
)

// RepositoryFactory creates repository instances
type RepositoryFactory interface {
	CreateDomainRepository() repository.DomainRepository
	CreateNodeRepository() repository.NodeRepository
	CreateAttributeRepository() repository.AttributeRepository
	CreateNodeAttributeRepository() repository.NodeAttributeRepository
	CreateTemplateRepository() repository.TemplateRepository
	CreateTemplateAttributeRepository() repository.TemplateAttributeRepository
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

func (f *ApplicationFactory) CreateTemplateRepository() repository.TemplateRepository {
	return sqliteRepo.NewTemplateRepository(f.db)
}

func (f *ApplicationFactory) CreateTemplateAttributeRepository() repository.TemplateAttributeRepository {
	return sqliteRepo.NewSQLiteTemplateAttributeRepository(f.db)
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
	templateRepo := f.CreateTemplateRepository()
	templateAttributeRepo := f.CreateTemplateAttributeRepository()

	// Create validation registry
	validatorRegistry := domainAttribute.NewValidatorRegistry()

	// Create services
	templateService, err := service.NewTemplateService(templateRepo, domainRepo, attributeRepo)
	if err != nil {
		panic("Failed to create template service: " + err.Error())
	}

	// Create use cases
	createDomainUC, listDomainsUC := f.CreateDomainUseCases(domainRepo)
	createNodeUC, listNodesUC := f.CreateNodeUseCases(nodeRepo, domainRepo)
	createAttributeUC, listAttributesUC := f.CreateAttributeUseCases(attributeRepo, domainRepo)
	setNodeAttributesUC := node.NewSetNodeAttributesUseCase(nodeRepo, attributeRepo, nodeAttributeRepo, templateService)
	filterNodesUC := node.NewFilterNodesByAttributesUseCase(nodeRepo)
	getNodeWithAttributesUC := node.NewGetNodeWithAttributesUseCase(nodeRepo, nodeAttributeRepo, attributeRepo)

	return &CleanDependencies{
		// Repositories
		DomainRepo:            domainRepo,
		NodeRepo:              nodeRepo,
		AttributeRepo:         attributeRepo,
		NodeAttributeRepo:     nodeAttributeRepo,
		TemplateRepo:          templateRepo,
		TemplateAttributeRepo: templateAttributeRepo,

		// Services
		TemplateService: templateService,

		// Validators
		ValidatorRegistry: validatorRegistry,

		// Use Cases
		CreateDomainUC:          createDomainUC,
		ListDomainsUC:           listDomainsUC,
		CreateNodeUC:            createNodeUC,
		ListNodesUC:             listNodesUC,
		CreateAttributeUC:       createAttributeUC,
		ListAttributesUC:        listAttributesUC,
		SetNodeAttributesUC:     setNodeAttributesUC,
		FilterNodesUC:           filterNodesUC,
		GetNodeWithAttributesUC: getNodeWithAttributesUC,
	}
}

// CleanDependencies holds Clean Architecture dependencies
type CleanDependencies struct {
	// Repositories
	DomainRepo            repository.DomainRepository
	NodeRepo              repository.NodeRepository
	AttributeRepo         repository.AttributeRepository
	NodeAttributeRepo     repository.NodeAttributeRepository
	TemplateRepo          repository.TemplateRepository
	TemplateAttributeRepo repository.TemplateAttributeRepository

	// Services
	TemplateService service.TemplateService

	// Validators
	ValidatorRegistry *domainAttribute.ValidatorRegistry

	// Use Cases
	CreateDomainUC          *domain.CreateDomainUseCase
	ListDomainsUC           *domain.ListDomainsUseCase
	CreateNodeUC            *node.CreateNodeUseCase
	ListNodesUC             *node.ListNodesUseCase
	CreateAttributeUC       *attribute.CreateAttributeUseCase
	ListAttributesUC        *attribute.ListAttributesUseCase
	SetNodeAttributesUC     *node.SetNodeAttributesUseCase
	FilterNodesUC           *node.FilterNodesByAttributesUseCase
	GetNodeWithAttributesUC *node.GetNodeWithAttributesUseCase
}

// Individual UseCase factory methods for MCP server
func (f *ApplicationFactory) CreateListDomainsUseCase() *domain.ListDomainsUseCase {
	domainRepo := f.CreateDomainRepository()
	_, listUseCase := f.CreateDomainUseCases(domainRepo)
	return listUseCase
}

func (f *ApplicationFactory) CreateCreateDomainUseCase() *domain.CreateDomainUseCase {
	domainRepo := f.CreateDomainRepository()
	createUseCase, _ := f.CreateDomainUseCases(domainRepo)
	return createUseCase
}

func (f *ApplicationFactory) CreateListNodesUseCase() *node.ListNodesUseCase {
	nodeRepo := f.CreateNodeRepository()
	domainRepo := f.CreateDomainRepository()
	_, listUseCase := f.CreateNodeUseCases(nodeRepo, domainRepo)
	return listUseCase
}

func (f *ApplicationFactory) CreateCreateNodeUseCase() *node.CreateNodeUseCase {
	nodeRepo := f.CreateNodeRepository()
	domainRepo := f.CreateDomainRepository()
	createUseCase, _ := f.CreateNodeUseCases(nodeRepo, domainRepo)
	return createUseCase
}

// Node attributes UseCase factory methods
func (f *ApplicationFactory) CreateSetNodeAttributesUseCase() *node.SetNodeAttributesUseCase {
	nodeRepo := f.CreateNodeRepository()
	attributeRepo := f.CreateAttributeRepository()
	nodeAttributeRepo := f.CreateNodeAttributeRepository()
	domainRepo := f.CreateDomainRepository()
	templateRepo := f.CreateTemplateRepository()
	templateService, err := service.NewTemplateService(templateRepo, domainRepo, attributeRepo)
	if err != nil {
		panic("Failed to create template service: " + err.Error())
	}
	return node.NewSetNodeAttributesUseCase(nodeRepo, attributeRepo, nodeAttributeRepo, templateService)
}

// Domain attributes (schema) UseCase factory methods
func (f *ApplicationFactory) CreateListAttributesUseCase() *attribute.ListAttributesUseCase {
	attributeRepo := f.CreateAttributeRepository()
	domainRepo := f.CreateDomainRepository()
	_, listUseCase := f.CreateAttributeUseCases(attributeRepo, domainRepo)
	return listUseCase
}

func (f *ApplicationFactory) CreateCreateAttributeUseCase() *attribute.CreateAttributeUseCase {
	attributeRepo := f.CreateAttributeRepository()
	domainRepo := f.CreateDomainRepository()
	createUseCase, _ := f.CreateAttributeUseCases(attributeRepo, domainRepo)
	return createUseCase
}

// Helper method to get domain by name
func (f *ApplicationFactory) GetDomainByName(ctx context.Context, domainName string) (*entity.Domain, error) {
	domainRepo := f.CreateDomainRepository()
	return domainRepo.GetByName(ctx, domainName)
}
