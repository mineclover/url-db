package handlers

import (
	"github.com/gin-gonic/gin"
)

// This file provides compatibility interfaces for legacy routing setup
// The actual implementations are in the individual handler files

func SetupLegacyRoutes(r *gin.Engine,
	domainHandler interface{ RegisterRoutes(r *gin.Engine) },
	nodeHandler interface{ RegisterRoutes(r *gin.Engine) },
	attributeHandler interface{ RegisterRoutes(r *gin.Engine) },
	nodeAttributeHandler interface{ RegisterRoutes(r *gin.Engine) },
	mcpHandler interface{ RegisterRoutes(r *gin.Engine) }) {

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
			"status":  "healthy",
			"service": "url-db",
		})
	})

	// Register all handler routes if they implement RegisterRoutes
	if domainHandler != nil {
		domainHandler.RegisterRoutes(r)
	}
	if nodeHandler != nil {
		nodeHandler.RegisterRoutes(r)
	}
	if attributeHandler != nil {
		attributeHandler.RegisterRoutes(r)
	}
	if nodeAttributeHandler != nil {
		nodeAttributeHandler.RegisterRoutes(r)
	}
	if mcpHandler != nil {
		mcpHandler.RegisterRoutes(r)
	}
}
