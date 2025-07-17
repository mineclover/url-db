package handlers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.RecoveryWithWriter(os.Stdout, func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(error); ok {
			log.Printf("Panic recovered: %s", err.Error())
		} else {
			log.Printf("Panic recovered: %v", recovered)
		}
		
		c.JSON(500, gin.H{
			"error":   "internal_error",
			"message": "An unexpected error occurred",
		})
	})
}

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		c.Next()
	}
}

func RequestSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(413, gin.H{
				"error":   "request_too_large",
				"message": "Request body too large",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Request-Timeout", timeout.String())
		c.Next()
	}
}

func JSONOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
				c.JSON(415, gin.H{
					"error":   "unsupported_media_type",
					"message": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		
		c.Next()
	}
}

func APIVersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-API-Version", "v1")
		c.Next()
	}
}

func SetupMiddleware(r *gin.Engine) {
	// Recovery middleware should be first
	r.Use(RecoveryMiddleware())
	
	// Request ID for tracing
	r.Use(RequestIDMiddleware())
	
	// Logging middleware
	r.Use(LoggingMiddleware())
	
	// CORS middleware
	r.Use(CORSMiddleware())
	
	// Security headers
	r.Use(SecurityHeadersMiddleware())
	
	// Request size limit (10MB)
	r.Use(RequestSizeMiddleware(10 * 1024 * 1024))
	
	// JSON content type enforcement for API routes
	apiGroup := r.Group("/api")
	apiGroup.Use(JSONOnlyMiddleware())
	
	// API version header
	r.Use(APIVersionMiddleware())
	
	// Timeout middleware (30 seconds)
	r.Use(TimeoutMiddleware(30 * time.Second))
}