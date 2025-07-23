package config

import (
	"os"
	"strconv"
	"strings"
	"url-db/internal/constants"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	ToolName             string
	AutoCreateAttributes bool
}

func Load() *Config {
	return &Config{
		Port:                 getEnv("PORT", strconv.Itoa(constants.DefaultPort)),
		DatabaseURL:          getEnv("DATABASE_URL", "file:./"+constants.DefaultDBPath),
		ToolName:             getEnv("TOOL_NAME", constants.DefaultServerName),
		AutoCreateAttributes: getBoolEnv("AUTO_CREATE_ATTRIBUTES", true),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		// 대소문자 구분 없이 true/false 파싱
		lowerValue := strings.ToLower(strings.TrimSpace(value))
		switch lowerValue {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		default:
			// 잘못된 값이면 기본값 반환
			return defaultValue
		}
	}
	return defaultValue
}
