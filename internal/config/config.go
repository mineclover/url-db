package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	ToolName    string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "file:./url-db.sqlite"),
		ToolName:    getEnv("TOOL_NAME", "url-db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}