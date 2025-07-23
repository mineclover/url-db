package setup

import (
	"database/sql"
	
	"github.com/jmoiron/sqlx"
	
	// Domain and services
	"url-db/internal/attributes"
	"url-db/internal/compositekey"
	"url-db/internal/domains"
	"url-db/internal/nodeattributes"
	"url-db/internal/nodes"
	"url-db/internal/repositories"
	"url-db/internal/services"
	
	// Interface adapters
	"url-db/internal/interface/adapter"
	"url-db/internal/interfaces/mcp"
	
	// Only imports needed for legacy dependencies
)

// Dependencies holds all application dependencies
// This is being gradually migrated to use Clean Architecture
type Dependencies struct {
	// Clean Architecture dependencies
	Clean *CleanDependencies
	
	// Legacy dependencies (to be removed gradually)
	Legacy *LegacyDependencies
}

// LegacyDependencies holds legacy dependencies during migration
type LegacyDependencies struct {
	DomainRepo           domains.DomainRepository
	NodeRepo             nodes.NodeRepository
	AttributeRepo        attributes.AttributeRepository
	NodeAttributeRepo    nodeattributes.Repository
	
	// Services
	DomainService        domains.DomainService
	NodeService          nodes.NodeService
	AttributeService     attributes.AttributeService
	NodeAttributeService nodeattributes.Service
	SubscriptionService  *services.SubscriptionService
	DependencyService    *services.DependencyService
	EventService         *services.EventService
	
	// MCP related
	MCPService           mcp.MCPService
	CompositeKeyService  *compositekey.Service
	MCPHandlerAdapter    *adapter.MCPHandlerServiceAdapter
}

// InitializeDependencies sets up all application dependencies
func InitializeDependencies(sqlDB *sql.DB, sqlxDB *sqlx.DB, toolName string) *Dependencies {
	// Create factories for Clean Architecture
	appFactory := NewApplicationFactory(sqlDB, sqlxDB, toolName)
	
	// Create Clean Architecture dependencies
	cleanDeps := appFactory.CreateCleanArchitectureDependencies()
	
	// Initialize legacy dependencies (for backward compatibility)
	legacyDeps := initializeLegacyDependencies(sqlDB, sqlxDB, toolName)
	
	return &Dependencies{
		Clean:  cleanDeps,
		Legacy: legacyDeps,
	}
}

// initializeLegacyDependencies initializes legacy dependencies
func initializeLegacyDependencies(sqlDB *sql.DB, sqlxDB *sqlx.DB, toolName string) *LegacyDependencies {
	legacy := &LegacyDependencies{}
	
	// Initialize legacy repositories
	legacy.DomainRepo = domains.NewDomainRepository(sqlDB)
	legacy.NodeRepo = nodes.NewSQLiteNodeRepository(sqlDB)
	legacy.AttributeRepo = attributes.NewSQLiteAttributeRepository(sqlDB)
	legacy.NodeAttributeRepo = nodeattributes.NewRepository(sqlxDB)
	
	// Use the repositories package for MCP service dependencies
	repos := repositories.NewRepositories(sqlDB)
	
	// Initialize external dependency repositories
	subscriptionRepo := repositories.NewSubscriptionRepository(sqlxDB)
	dependencyRepo := repositories.NewDependencyRepository(sqlxDB)
	eventRepo := repositories.NewEventRepository(sqlxDB)
	
	// Initialize external dependency services
	legacy.SubscriptionService = services.NewSubscriptionService(subscriptionRepo, repos.Node, eventRepo)
	legacy.DependencyService = services.NewDependencyService(dependencyRepo, repos.Node, eventRepo)
	legacy.EventService = services.NewEventService(eventRepo, repos.Node)
	
	// Initialize services
	legacy.DomainService = domains.NewDomainService(legacy.DomainRepo)
	basicNodeService := nodes.NewNodeService(legacy.NodeRepo)
	// Wrap node service with event tracking
	legacy.NodeService = nodes.NewNodeServiceWithEvents(basicNodeService, legacy.EventService)
	legacy.AttributeService = attributes.NewAttributeService(legacy.AttributeRepo, legacy.DomainService)
	
	// Initialize validators and managers for node attributes
	nodeAttributeValidator := nodeattributes.NewValidator()
	nodeAttributeOrderManager := nodeattributes.NewOrderManager(legacy.NodeAttributeRepo)
	legacy.NodeAttributeService = nodeattributes.NewService(legacy.NodeAttributeRepo, nodeAttributeValidator, nodeAttributeOrderManager)
	
	// Initialize composite key service
	legacy.CompositeKeyService = compositekey.NewService(toolName)
	compositeKeyAdapter := mcp.NewCompositeKeyAdapter(legacy.CompositeKeyService)
	
	// Initialize MCP converter and service
	mcpConverter := mcp.NewConverter(compositeKeyAdapter, toolName)
	legacy.MCPService = mcp.NewMCPService(
		adapter.NewMCPNodeServiceAdapter(legacy.NodeService), 
		adapter.NewMCPDomainServiceAdapter(legacy.DomainService, legacy.DomainRepo), 
		adapter.NewMCPAttributeServiceAdapter(legacy.NodeAttributeService, legacy.AttributeService), 
		adapter.NewNodeCountServiceAdapter(repos.Node), 
		legacy.SubscriptionService, 
		legacy.DependencyService, 
		legacy.EventService, 
		mcpConverter,
	)
	
	// Create MCP handler adapter
	legacy.MCPHandlerAdapter = adapter.NewMCPHandlerServiceAdapter(legacy.MCPService)
	
	return legacy
}