package main

import (
	"log"
	
	"url-db/internal/config"
	"url-db/internal/database"
	"url-db/internal/handlers"
	"url-db/internal/domains"
	"url-db/internal/nodes"
	"url-db/internal/attributes"
	"url-db/internal/nodeattributes"
	"url-db/internal/services"
	"url-db/internal/compositekey"
	"url-db/internal/mcp"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize repositories
	domainRepo := domains.NewRepository(db)
	nodeRepo := nodes.NewRepository(db)
	attributeRepo := attributes.NewRepository(db)
	nodeAttributeRepo := nodeattributes.NewRepository(db)

	// Initialize services
	domainService := domains.NewService(domainRepo)
	nodeService := nodes.NewService(nodeRepo, domainRepo)
	attributeService := attributes.NewService(attributeRepo, domainRepo)
	
	// Initialize validators and managers for node attributes
	nodeAttributeValidator := nodeattributes.NewValidator()
	nodeAttributeOrderManager := nodeattributes.NewOrderManager(nodeAttributeRepo)
	nodeAttributeService := nodeattributes.NewService(nodeAttributeRepo, nodeAttributeValidator, nodeAttributeOrderManager)
	
	// Initialize additional services
	nodeCountService := services.NewNodeCountService(nodeRepo)
	
	// Initialize composite key service
	compositeKeyService := compositekey.NewService(cfg.ToolName)
	compositeKeyAdapter := mcp.NewCompositeKeyAdapter(compositeKeyService)
	
	// Initialize MCP converter and service
	mcpConverter := mcp.NewConverter(compositeKeyAdapter)
	mcpService := mcp.NewMCPService(nodeService, domainService, attributeService, nodeCountService, mcpConverter)

	// Initialize handlers
	domainHandler := domains.NewHandler(domainService)
	nodeHandler := nodes.NewHandler(nodeService)
	attributeHandler := attributes.NewHandler(attributeService)
	nodeAttributeHandler := nodeattributes.NewHandler(nodeAttributeService)
	
	// Create MCP handler with service
	mcpHandler := handlers.NewMCPHandler(mcpService)

	// Setup routes
	r := gin.Default()
	handlers.SetupRoutes(r, domainHandler, nodeHandler, attributeHandler, nodeAttributeHandler, mcpHandler)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}