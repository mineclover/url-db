package domains

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/url-db/internal/models"
)

type DomainHandler struct {
	service DomainService
}

func NewDomainHandler(service DomainService) *DomainHandler {
	return &DomainHandler{service: service}
}

func (h *DomainHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/domains", h.CreateDomain)
	router.GET("/domains", h.ListDomains)
	router.GET("/domains/:id", h.GetDomain)
	router.PUT("/domains/:id", h.UpdateDomain)
	router.DELETE("/domains/:id", h.DeleteDomain)
}

func (h *DomainHandler) CreateDomain(c *gin.Context) {
	var req models.CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"message": err.Error(),
		})
		return
	}

	domain, err := h.service.CreateDomain(c.Request.Context(), &req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "validation_error:") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "validation_error",
				"message": strings.TrimPrefix(err.Error(), "validation_error: "),
			})
			return
		}
		if strings.HasPrefix(err.Error(), "conflict:") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "conflict",
				"message": strings.TrimPrefix(err.Error(), "conflict: "),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_error",
			"message": "Failed to create domain",
		})
		return
	}

	c.JSON(http.StatusCreated, domain)
}

func (h *DomainHandler) GetDomain(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	domain, err := h.service.GetDomain(c.Request.Context(), id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not_found:") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "not_found",
				"message": strings.TrimPrefix(err.Error(), "not_found: "),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_error",
			"message": "Failed to get domain",
		})
		return
	}

	c.JSON(http.StatusOK, domain)
}

func (h *DomainHandler) ListDomains(c *gin.Context) {
	page := 1
	size := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := c.Query("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			size = s
		}
	}

	response, err := h.service.ListDomains(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_error",
			"message": "Failed to list domains",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DomainHandler) UpdateDomain(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	var req models.UpdateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"message": err.Error(),
		})
		return
	}

	domain, err := h.service.UpdateDomain(c.Request.Context(), id, &req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "validation_error:") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "validation_error",
				"message": strings.TrimPrefix(err.Error(), "validation_error: "),
			})
			return
		}
		if strings.HasPrefix(err.Error(), "not_found:") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "not_found",
				"message": strings.TrimPrefix(err.Error(), "not_found: "),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_error",
			"message": "Failed to update domain",
		})
		return
	}

	c.JSON(http.StatusOK, domain)
}

func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	err = h.service.DeleteDomain(c.Request.Context(), id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not_found:") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "not_found",
				"message": strings.TrimPrefix(err.Error(), "not_found: "),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_error",
			"message": "Failed to delete domain",
		})
		return
	}

	c.Status(http.StatusNoContent)
}