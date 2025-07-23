package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine,
	domainHandler *DomainHandler,
	nodeHandler *NodeHandler,
	attributeHandler *AttributeHandler,
	nodeAttributeHandler *NodeAttributeHandler,
	mcpHandler *MCPHandler,
	healthHandler *HealthHandler,
	subscriptionHandler *SubscriptionHandler,
	dependencyHandler *DependencyHandler,
	eventHandler *EventHandler) {

	// Setup middleware
	SetupMiddleware(r)

	// Health check routes (no auth required)
	healthHandler.RegisterRoutes(r)

	// API routes
	api := r.Group("/api")
	{
		// Domain routes
		domains := api.Group("/domains")
		{
			domains.POST("", domainHandler.CreateDomain)
			domains.GET("", domainHandler.GetDomains)
			domains.GET("/:id", domainHandler.GetDomain)
			domains.PUT("/:id", domainHandler.UpdateDomain)
			domains.DELETE("/:id", domainHandler.DeleteDomain)

			// Domain-specific node routes
			domainNodes := domains.Group("/:domain_id/urls")
			{
				domainNodes.POST("", nodeHandler.CreateNode)
				domainNodes.GET("", nodeHandler.GetNodesByDomain)
				domainNodes.POST("/find", nodeHandler.FindNodeByURL)
			}

			// Domain-specific attribute routes
			domainAttributes := domains.Group("/:domain_id/attributes")
			{
				domainAttributes.POST("", attributeHandler.CreateAttribute)
				domainAttributes.GET("", attributeHandler.GetAttributesByDomain)
			}
		}

		// Node routes
		urls := api.Group("/urls")
		{
			urls.GET("/:id", nodeHandler.GetNode)
			urls.PUT("/:id", nodeHandler.UpdateNode)
			urls.DELETE("/:id", nodeHandler.DeleteNode)

			// Node-specific attribute routes
			urlAttributes := urls.Group("/:url_id/attributes")
			{
				urlAttributes.POST("", nodeAttributeHandler.CreateNodeAttribute)
				urlAttributes.GET("", nodeAttributeHandler.GetNodeAttributesByNode)
				urlAttributes.DELETE("/:attribute_id", nodeAttributeHandler.DeleteNodeAttributeByNodeAndAttribute)
			}
		}

		// Attribute routes
		attributes := api.Group("/attributes")
		{
			attributes.GET("/:id", attributeHandler.GetAttribute)
			attributes.PUT("/:id", attributeHandler.UpdateAttribute)
			attributes.DELETE("/:id", attributeHandler.DeleteAttribute)
		}

		// Node attribute routes
		urlAttributesGlobal := api.Group("/url-attributes")
		{
			urlAttributesGlobal.GET("/:id", nodeAttributeHandler.GetNodeAttribute)
			urlAttributesGlobal.PUT("/:id", nodeAttributeHandler.UpdateNodeAttribute)
			urlAttributesGlobal.DELETE("/:id", nodeAttributeHandler.DeleteNodeAttribute)
		}

		// External dependency management routes
		// Subscription routes
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.GET("/:id", subscriptionHandler.GetSubscription)
			subscriptions.PUT("/:id", subscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:id", subscriptionHandler.DeleteSubscription)
			subscriptions.GET("", subscriptionHandler.GetServiceSubscriptions) // with ?service= query param
		}

		// Node subscription routes
		nodes := api.Group("/nodes")
		{
			nodes.POST("/:nodeId/subscriptions", subscriptionHandler.CreateSubscription)
			nodes.GET("/:nodeId/subscriptions", subscriptionHandler.GetNodeSubscriptions)

			// Node dependency routes
			nodes.POST("/:nodeId/dependencies", dependencyHandler.CreateDependency)
			nodes.GET("/:nodeId/dependencies", dependencyHandler.GetNodeDependencies)
			nodes.GET("/:nodeId/dependents", dependencyHandler.GetNodeDependents)

			// Node event routes
			nodes.GET("/:nodeId/events", eventHandler.GetNodeEvents)
		}

		// Dependency routes
		dependencies := api.Group("/dependencies")
		{
			dependencies.GET("/:id", dependencyHandler.GetDependency)
			dependencies.DELETE("/:id", dependencyHandler.DeleteDependency)
		}

		// Event routes
		events := api.Group("/events")
		{
			events.GET("/pending", eventHandler.GetPendingEvents)
			events.POST("/:eventId/process", eventHandler.ProcessEvent)
			events.GET("", eventHandler.GetEventsByType) // with query params
			events.GET("/stats", eventHandler.GetEventStats)
			events.POST("/cleanup", eventHandler.CleanupEvents)
		}

		// MCP routes
		mcp := api.Group("/mcp")
		{
			// Node operations
			nodes := mcp.Group("/nodes")
			{
				nodes.POST("", mcpHandler.CreateMCPNode)
				nodes.GET("", mcpHandler.GetMCPNodes)
				nodes.GET("/:composite_id", mcpHandler.GetMCPNode)
				nodes.PUT("/:composite_id", mcpHandler.UpdateMCPNode)
				nodes.DELETE("/:composite_id", mcpHandler.DeleteMCPNode)
				nodes.POST("/find", mcpHandler.FindMCPNodeByURL)
				nodes.POST("/batch", mcpHandler.BatchGetMCPNodes)

				// Node attributes
				nodes.GET("/:composite_id/attributes", mcpHandler.GetMCPNodeAttributes)
				nodes.PUT("/:composite_id/attributes", mcpHandler.SetMCPNodeAttributes)
			}

			// Domain operations
			domains := mcp.Group("/domains")
			{
				domains.GET("", mcpHandler.GetMCPDomains)
				domains.POST("", mcpHandler.CreateMCPDomain)
			}

			// Server info
			server := mcp.Group("/server")
			{
				server.GET("/info", mcpHandler.GetMCPServerInfo)
			}
		}
	}
}

func NewRouter() *gin.Engine {
	// Set gin mode based on environment
	gin.SetMode(gin.ReleaseMode) // Can be configured via environment variable

	r := gin.New()

	// Add basic middleware that should always be present
	r.Use(gin.Recovery())

	return r
}

type RouterConfig struct {
	DomainHandler        *DomainHandler
	NodeHandler          *NodeHandler
	AttributeHandler     *AttributeHandler
	NodeAttributeHandler *NodeAttributeHandler
	MCPHandler           *MCPHandler
	HealthHandler        *HealthHandler
	SubscriptionHandler  *SubscriptionHandler
	DependencyHandler    *DependencyHandler
	EventHandler         *EventHandler
}

func SetupRouter(config *RouterConfig) *gin.Engine {
	r := NewRouter()

	SetupRoutes(r,
		config.DomainHandler,
		config.NodeHandler,
		config.AttributeHandler,
		config.NodeAttributeHandler,
		config.MCPHandler,
		config.HealthHandler,
		config.SubscriptionHandler,
		config.DependencyHandler,
		config.EventHandler,
	)

	return r
}
