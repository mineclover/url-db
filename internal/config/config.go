package config

import (
	"os"
	"strconv"
	"url-db/internal/constants"
)

type Config struct {
	Port        string
	DatabaseURL string
	ToolName    string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", strconv.Itoa(constants.DefaultPort)),
		DatabaseURL: getEnv("DATABASE_URL", "file:./"+constants.DefaultDBPath),
		ToolName:    getEnv("TOOL_NAME", constants.DefaultServerName),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
