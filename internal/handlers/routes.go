package handlers

import (
	"github.com/gin-gonic/gin"
	"internal/domains"
	"internal/nodes"
	"internal/attributes"
	"internal/nodeattributes"
)

type DomainHandler interface {
	RegisterRoutes(r *gin.Engine)
}

type NodeHandler interface {
	RegisterRoutes(r *gin.Engine)
}

type AttributeHandler interface {
	RegisterRoutes(r *gin.Engine)
}

type NodeAttributeHandler interface {
	RegisterRoutes(r *gin.Engine)
}

type MCPHandlerInterface interface {
	RegisterRoutes(r *gin.Engine)
}

func SetupRoutes(r *gin.Engine, 
	domainHandler DomainHandler,
	nodeHandler NodeHandler,
	attributeHandler AttributeHandler,
	nodeAttributeHandler NodeAttributeHandler,
	mcpHandler MCPHandlerInterface) {

	// Setup CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "url-db",
		})
	})

	// Register all handler routes
	domainHandler.RegisterRoutes(r)
	nodeHandler.RegisterRoutes(r)
	attributeHandler.RegisterRoutes(r)
	nodeAttributeHandler.RegisterRoutes(r)
	mcpHandler.RegisterRoutes(r)
}