package setup

import (
	"url-db/internal/attributes"
	"url-db/internal/constants"
	"url-db/internal/domains"
	handlers "url-db/internal/interfaces/http"
	"url-db/internal/nodeattributes"
	"url-db/internal/nodes"
	
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all routes for the application
func SetupRouter(deps *Dependencies) *gin.Engine {
	r := gin.Default()
	
	// Setup CORS middleware
	r.Use(corsMiddleware())
	
	// Health check endpoint
	r.GET("/health", healthCheck)
	
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Setup API routes
	api := r.Group("/api")
	
	// Setup Clean Architecture routes
	setupCleanRoutes(api, deps)
	
	// Setup legacy routes (for backward compatibility)
	setupLegacyRoutes(api, deps)
	
	// Setup MCP routes if needed
	mcpHandler := handlers.NewMCPHandler(deps.Legacy.MCPHandlerAdapter)
	mcpHandler.RegisterRoutes(r)
	
	return r
}

// corsMiddleware returns a CORS middleware handler
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// healthCheck handles the health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(constants.StatusOK, gin.H{
		"status":  "healthy",
		"service": constants.DefaultServerName,
	})
}

// setupLegacyRoutes sets up all legacy routes for backward compatibility
func setupLegacyRoutes(api *gin.RouterGroup, deps *Dependencies) {
	// Initialize handlers
	domainHandler := domains.NewDomainHandler(deps.Legacy.DomainService)
	nodeHandler := nodes.NewNodeHandler(deps.Legacy.NodeService)
	attributeHandler := attributes.NewAttributeHandler(deps.Legacy.AttributeService)
	nodeAttributeHandler := nodeattributes.NewHandler(deps.Legacy.NodeAttributeService)
	
	// Initialize external dependency handlers
	subscriptionHandler := handlers.NewSubscriptionHandler(deps.Legacy.SubscriptionService)
	dependencyHandler := handlers.NewDependencyHandler(deps.Legacy.DependencyService)
	eventHandler := handlers.NewEventHandler(deps.Legacy.EventService)
	
	// Domain routes
	api.POST("/domains", domainHandler.CreateDomain)
	api.GET("/domains", domainHandler.ListDomains)
	api.GET("/domains/:domain_id", func(c *gin.Context) {
		// Convert :domain_id to :id for the domain handler
		domainID := c.Param("domain_id")
		c.Params = []gin.Param{{Key: "id", Value: domainID}}
		domainHandler.GetDomain(c)
	})
	api.PUT("/domains/:domain_id", func(c *gin.Context) {
		// Convert :domain_id to :id for the domain handler
		domainID := c.Param("domain_id")
		c.Params = []gin.Param{{Key: "id", Value: domainID}}
		domainHandler.UpdateDomain(c)
	})
	api.DELETE("/domains/:domain_id", func(c *gin.Context) {
		// Convert :domain_id to :id for the domain handler
		domainID := c.Param("domain_id")
		c.Params = []gin.Param{{Key: "id", Value: domainID}}
		domainHandler.DeleteDomain(c)
	})
	
	// Domain attribute routes
	api.POST("/domains/:domain_id/attributes", attributeHandler.CreateAttribute)
	api.GET("/domains/:domain_id/attributes", attributeHandler.ListAttributes)
	
	// Attribute detail routes
	api.GET("/attributes/:id", attributeHandler.GetAttribute)
	api.PUT("/attributes/:id", attributeHandler.UpdateAttribute)
	api.DELETE("/attributes/:id", attributeHandler.DeleteAttribute)
	
	// Node routes
	api.POST("/domains/:domain_id/urls", nodeHandler.CreateNode)
	api.GET("/domains/:domain_id/urls", nodeHandler.GetNodesByDomain)
	api.POST("/domains/:domain_id/urls/find", nodeHandler.FindNodeByURL)
	api.GET("/urls/:id", nodeHandler.GetNode)
	api.PUT("/urls/:id", nodeHandler.UpdateNode)
	api.DELETE("/urls/:id", nodeHandler.DeleteNode)
	
	// Node attribute routes
	api.POST("/urls/:id/attributes", func(c *gin.Context) {
		// Convert :id to :url_id for the node attribute handler
		urlID := c.Param("id")
		c.Params = []gin.Param{{Key: "url_id", Value: urlID}}
		nodeAttributeHandler.CreateNodeAttribute(c)
	})
	api.GET("/urls/:id/attributes", func(c *gin.Context) {
		// Convert :id to :url_id for the node attribute handler
		urlID := c.Param("id")
		c.Params = []gin.Param{{Key: "url_id", Value: urlID}}
		nodeAttributeHandler.GetNodeAttributesByNodeID(c)
	})
	api.GET("/url-attributes/:id", nodeAttributeHandler.GetNodeAttributeByID)
	api.PUT("/url-attributes/:id", nodeAttributeHandler.UpdateNodeAttribute)
	api.DELETE("/url-attributes/:id", nodeAttributeHandler.DeleteNodeAttribute)
	api.DELETE("/urls/:id/attributes/:attribute_id", func(c *gin.Context) {
		// Convert :id to :url_id for the node attribute handler
		urlID := c.Param("id")
		attributeID := c.Param("attribute_id")
		c.Params = []gin.Param{
			{Key: "url_id", Value: urlID},
			{Key: "attribute_id", Value: attributeID},
		}
		nodeAttributeHandler.DeleteNodeAttributeByNodeIDAndAttributeID(c)
	})
	
	// External dependency routes
	// Subscription routes
	api.POST("/nodes/:nodeId/subscriptions", subscriptionHandler.CreateSubscription)
	api.GET("/nodes/:nodeId/subscriptions", subscriptionHandler.GetNodeSubscriptions)
	api.GET("/subscriptions", subscriptionHandler.GetServiceSubscriptions)
	api.GET("/subscriptions/:id", subscriptionHandler.GetSubscription)
	api.PUT("/subscriptions/:id", subscriptionHandler.UpdateSubscription)
	api.DELETE("/subscriptions/:id", subscriptionHandler.DeleteSubscription)
	
	// Dependency routes
	api.POST("/nodes/:nodeId/dependencies", dependencyHandler.CreateDependency)
	api.GET("/nodes/:nodeId/dependencies", dependencyHandler.GetNodeDependencies)
	api.GET("/nodes/:nodeId/dependents", dependencyHandler.GetNodeDependents)
	api.GET("/dependencies/:id", dependencyHandler.GetDependency)
	api.DELETE("/dependencies/:id", dependencyHandler.DeleteDependency)
	
	// Event routes
	api.GET("/nodes/:nodeId/events", eventHandler.GetNodeEvents)
	api.GET("/events/pending", eventHandler.GetPendingEvents)
	api.GET("/events", eventHandler.GetEventsByType)
	api.GET("/events/stats", eventHandler.GetEventStats)
	api.POST("/events/:eventId/process", eventHandler.ProcessEvent)
	api.POST("/events/cleanup", eventHandler.CleanupEvents)
}

// setupCleanRoutes sets up Clean Architecture routes
func setupCleanRoutes(api *gin.RouterGroup, deps *Dependencies) {
	// Clean Architecture routes - using proper HTTP handlers
	v2 := api.Group("/v2")
	
	// Domain routes
	v2.POST("/domains", func(c *gin.Context) {
		deps.Clean.DomainHandler.CreateDomain(c.Writer, c.Request)
	})
	v2.GET("/domains", func(c *gin.Context) {
		deps.Clean.DomainHandler.ListDomains(c.Writer, c.Request)
	})
	
	// Node routes
	v2.POST("/nodes", func(c *gin.Context) {
		deps.Clean.NodeHandler.CreateNode(c.Writer, c.Request)
	})
	v2.GET("/nodes", func(c *gin.Context) {
		deps.Clean.NodeHandler.ListNodes(c.Writer, c.Request)
	})
	
	// Attribute routes
	v2.POST("/domains/:domain_id/attributes", deps.Clean.AttributeHandler.CreateAttribute)
	v2.GET("/domains/:domain_id/attributes", deps.Clean.AttributeHandler.ListAttributes)
}