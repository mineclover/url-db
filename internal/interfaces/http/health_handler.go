package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthService interface {
	CheckDatabaseConnection() error
	GetSystemInfo() (*HealthInfo, error)
}

type HealthInfo struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Database  DatabaseHealth    `json:"database"`
	System    SystemHealth      `json:"system"`
	Checks    map[string]string `json:"checks"`
}

type DatabaseHealth struct {
	Status     string `json:"status"`
	Connection string `json:"connection"`
	Response   string `json:"response"`
}

type SystemHealth struct {
	Uptime    time.Duration `json:"uptime"`
	GoVersion string        `json:"go_version"`
	Platform  string        `json:"platform"`
}

type HealthHandler struct {
	*BaseHandler
	healthService HealthService
	startTime     time.Time
}

func NewHealthHandler(healthService HealthService) *HealthHandler {
	return &HealthHandler{
		BaseHandler:   NewBaseHandler(),
		healthService: healthService,
		startTime:     time.Now(),
	}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	checks := make(map[string]string)
	overallStatus := "healthy"

	// Check database connection
	dbStatus := "healthy"
	if err := h.healthService.CheckDatabaseConnection(); err != nil {
		dbStatus = "unhealthy"
		overallStatus = "unhealthy"
		checks["database"] = err.Error()
	} else {
		checks["database"] = "ok"
	}

	health := &HealthInfo{
		Status:    overallStatus,
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Database: DatabaseHealth{
			Status:     dbStatus,
			Connection: "sqlite",
			Response:   "ok",
		},
		System: SystemHealth{
			Uptime:    time.Since(h.startTime),
			GoVersion: "go1.21",
			Platform:  "linux/amd64",
		},
		Checks: checks,
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}

func (h *HealthHandler) GetLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now(),
	})
}

func (h *HealthHandler) GetReadiness(c *gin.Context) {
	// Check if the service is ready to serve traffic
	if err := h.healthService.CheckDatabaseConnection(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not ready",
			"reason":    "database connection failed",
			"error":     err.Error(),
			"timestamp": time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now(),
	})
}

func (h *HealthHandler) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    "1.0.0",
		"build_time": "2024-01-01T00:00:00Z",
		"commit":     "abc123",
		"go_version": "go1.21",
	})
}

func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	// Health check endpoints
	health := r.Group("/health")
	{
		health.GET("", h.GetHealth)
		health.GET("/live", h.GetLiveness)
		health.GET("/ready", h.GetReadiness)
	}

	// Version endpoint
	r.GET("/version", h.GetVersion)
}
