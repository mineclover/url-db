package setup

import (
	"database/sql"
	
	"github.com/jmoiron/sqlx"
	
	"url-db/internal/domain/repository"
	"url-db/internal/application/usecase/domain"
	"url-db/internal/application/usecase/node"
	"url-db/internal/application/usecase/attribute"
	sqliteRepo "url-db/internal/infrastructure/persistence/sqlite/repository"
	cleanHandler "url-db/internal/interface/http/handler"
	"url-db/internal/interface/adapter"
	"url-db/internal/compositekey"
	"url-db/internal/interfaces/mcp"
)

// RepositoryFactory creates repository instances
type RepositoryFactory interface {
	CreateDomainRepository() repository.DomainRepository
	CreateNodeRepository() repository.NodeRepository
	CreateAttributeRepository() repository.AttributeRepository
}

// UseCaseFactory creates use case instances
type UseCaseFactory interface {
	CreateDomainUseCases() (*domain.CreateDomainUseCase, *domain.ListDomainsUseCase)
	CreateNodeUseCases() (*node.CreateNodeUseCase, *node.ListNodesUseCase)
	CreateAttributeUseCases() (*attribute.CreateAttributeUseCase, *attribute.ListAttributesUseCase)
}

// HandlerFactory creates handler instances
type HandlerFactory interface {
	CreateDomainHandler() *cleanHandler.DomainHandler
	CreateNodeHandler() *cleanHandler.NodeHandler
	CreateAttributeHandler() *cleanHandler.AttributeHandler
}

// AdapterFactory creates adapter instances
type AdapterFactory interface {
	CreateMCPUseCaseAdapter() *adapter.MCPUseCaseAdapter
}

// sqliteRepositoryFactory implements RepositoryFactory for SQLite
type sqliteRepositoryFactory struct {
	db *sql.DB
}

// NewSQLiteRepositoryFactory creates a new SQLite repository factory
func NewSQLiteRepositoryFactory(db *sql.DB) RepositoryFactory {
	return &sqliteRepositoryFactory{db: db}
}

func (f *sqliteRepositoryFactory) CreateDomainRepository() repository.DomainRepository {
	return sqliteRepo.NewDomainRepository(f.db)
}

func (f *sqliteRepositoryFactory) CreateNodeRepository() repository.NodeRepository {
	return sqliteRepo.NewNodeRepository(f.db)
}

func (f *sqliteRepositoryFactory) CreateAttributeRepository() repository.AttributeRepository {
	return sqliteRepo.NewAttributeRepository(f.db)
}

// useCaseFactory implements UseCaseFactory
type useCaseFactory struct {
	repoFactory RepositoryFactory
}

// NewUseCaseFactory creates a new use case factory
func NewUseCaseFactory(repoFactory RepositoryFactory) UseCaseFactory {
	return &useCaseFactory{repoFactory: repoFactory}
}

func (f *useCaseFactory) CreateDomainUseCases() (*domain.CreateDomainUseCase, *domain.ListDomainsUseCase) {
	domainRepo := f.repoFactory.CreateDomainRepository()
	return domain.NewCreateDomainUseCase(domainRepo), 
		   domain.NewListDomainsUseCase(domainRepo)
}

func (f *useCaseFactory) CreateNodeUseCases() (*node.CreateNodeUseCase, *node.ListNodesUseCase) {
	nodeRepo := f.repoFactory.CreateNodeRepository()
	domainRepo := f.repoFactory.CreateDomainRepository()
	return node.NewCreateNodeUseCase(nodeRepo, domainRepo),
		   node.NewListNodesUseCase(nodeRepo)
}

func (f *useCaseFactory) CreateAttributeUseCases() (*attribute.CreateAttributeUseCase, *attribute.ListAttributesUseCase) {
	attributeRepo := f.repoFactory.CreateAttributeRepository()
	domainRepo := f.repoFactory.CreateDomainRepository()
	return attribute.NewCreateAttributeUseCase(attributeRepo, domainRepo),
		   attribute.NewListAttributesUseCase(attributeRepo, domainRepo)
}

// handlerFactory implements HandlerFactory
type handlerFactory struct {
	useCaseFactory UseCaseFactory
}

// NewHandlerFactory creates a new handler factory
func NewHandlerFactory(useCaseFactory UseCaseFactory) HandlerFactory {
	return &handlerFactory{useCaseFactory: useCaseFactory}
}

func (f *handlerFactory) CreateDomainHandler() *cleanHandler.DomainHandler {
	createUC, listUC := f.useCaseFactory.CreateDomainUseCases()
	return cleanHandler.NewDomainHandler(createUC, listUC)
}

func (f *handlerFactory) CreateNodeHandler() *cleanHandler.NodeHandler {
	createUC, listUC := f.useCaseFactory.CreateNodeUseCases()
	return cleanHandler.NewNodeHandler(createUC, listUC)
}

func (f *handlerFactory) CreateAttributeHandler() *cleanHandler.AttributeHandler {
	createUC, listUC := f.useCaseFactory.CreateAttributeUseCases()
	return cleanHandler.NewAttributeHandler(createUC, listUC)
}

// adapterFactory implements AdapterFactory
type adapterFactory struct {
	useCaseFactory UseCaseFactory
	toolName       string
}

// NewAdapterFactory creates a new adapter factory
func NewAdapterFactory(useCaseFactory UseCaseFactory, toolName string) AdapterFactory {
	return &adapterFactory{
		useCaseFactory: useCaseFactory,
		toolName:       toolName,
	}
}

func (f *adapterFactory) CreateMCPUseCaseAdapter() *adapter.MCPUseCaseAdapter {
	createDomainUC, listDomainsUC := f.useCaseFactory.CreateDomainUseCases()
	createNodeUC, listNodesUC := f.useCaseFactory.CreateNodeUseCases()
	createAttributeUC, listAttributesUC := f.useCaseFactory.CreateAttributeUseCases()
	
	// Create composite key service and converter
	compositeKeyService := compositekey.NewService(f.toolName)
	compositeKeyAdapter := mcp.NewCompositeKeyAdapter(compositeKeyService)
	mcpConverter := mcp.NewConverter(compositeKeyAdapter, f.toolName)
	
	return adapter.NewMCPUseCaseAdapter(
		createDomainUC,
		listDomainsUC,
		createNodeUC,
		listNodesUC,
		createAttributeUC,
		listAttributesUC,
		mcpConverter,
	)
}

// ApplicationFactory is the main factory that orchestrates all other factories
type ApplicationFactory struct {
	repoFactory    RepositoryFactory
	useCaseFactory UseCaseFactory
	handlerFactory HandlerFactory
	adapterFactory AdapterFactory
}

// NewApplicationFactory creates a new application factory
func NewApplicationFactory(sqlDB *sql.DB, sqlxDB *sqlx.DB, toolName string) *ApplicationFactory {
	repoFactory := NewSQLiteRepositoryFactory(sqlDB)
	useCaseFactory := NewUseCaseFactory(repoFactory)
	handlerFactory := NewHandlerFactory(useCaseFactory)
	adapterFactory := NewAdapterFactory(useCaseFactory, toolName)
	
	return &ApplicationFactory{
		repoFactory:    repoFactory,
		useCaseFactory: useCaseFactory,
		handlerFactory: handlerFactory,
		adapterFactory: adapterFactory,
	}
}

// CreateCleanArchitectureDependencies creates all Clean Architecture dependencies
func (f *ApplicationFactory) CreateCleanArchitectureDependencies() *CleanDependencies {
	return &CleanDependencies{
		DomainHandler:     f.handlerFactory.CreateDomainHandler(),
		NodeHandler:       f.handlerFactory.CreateNodeHandler(),
		AttributeHandler:  f.handlerFactory.CreateAttributeHandler(),
		MCPUseCaseAdapter: f.adapterFactory.CreateMCPUseCaseAdapter(),
	}
}

// CleanDependencies holds only the Clean Architecture dependencies
type CleanDependencies struct {
	DomainHandler     *cleanHandler.DomainHandler
	NodeHandler       *cleanHandler.NodeHandler
	AttributeHandler  *cleanHandler.AttributeHandler
	MCPUseCaseAdapter *adapter.MCPUseCaseAdapter
}