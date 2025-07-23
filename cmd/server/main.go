package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// Clean Architecture imports
	"url-db/internal/config"
	"url-db/internal/constants"
	"url-db/internal/database"
	"url-db/internal/interface/setup"
)

// @title           URL Database API
// @version         1.0
// @description     A URL management system with Clean Architecture and MCP integration.

func main() {
	// Parse command line flags
	var (
		dbPath   = flag.String("db-path", "", "Path to the database file")
		toolName = flag.String("tool-name", constants.DefaultServerName, "Tool name for composite keys")
		port     = flag.String("port", "8080", "Port for HTTP server")
		showHelp = flag.Bool("help", false, "Show help message")
		version  = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *showHelp {
		fmt.Println("URL Database Server - Clean Architecture")
		fmt.Println("Usage:")
		fmt.Println("  url-db [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -db-path string    Path to the database file")
		fmt.Println("  -tool-name string  Tool name for composite keys")
		fmt.Println("  -port string       Port for HTTP server (default: 8080)")
		fmt.Println("  -help             Show help message")
		fmt.Println("  -version          Show version information")
		os.Exit(0)
	}

	if *version {
		fmt.Println("URL Database Server v" + constants.DefaultServerVersion)
		fmt.Println("Clean Architecture Implementation")
		os.Exit(0)
	}

	// Load configuration
	cfg := config.Load()

	// Override with command-line flags
	if *dbPath != "" {
		cfg.DatabaseURL = "file:" + *dbPath
	}
	if *toolName != "" {
		cfg.ToolName = *toolName
	}
	if *port != "" {
		cfg.Port = *port
	}

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Clean Architecture factory
	factory := setup.NewApplicationFactory(db.DB(), nil, cfg.ToolName)
	
	// Create router (this would be implemented in setup package)
	router := setup.SetupCleanRouter(factory)
	
	// Start server
	log.Printf("Starting Clean Architecture server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}