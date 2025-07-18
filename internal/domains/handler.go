package domains

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
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

// CreateDomain godoc
// @Summary      Create a new domain
// @Description  Create a new domain with name and description
// @Tags         domains
// @Accept       json
// @Produce      json
// @Param        domain  body      models.CreateDomainRequest  true  "Domain data"
// @Success      201     {object}  models.Domain
// @Failure      400     {object}  map[string]interface{}
// @Failure      409     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /domains [post]
func (h *DomainHandler) CreateDomain(c *gin.Context) {
	var req models.CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}

	domain, err := h.service.CreateDomain(c.Request.Context(), &req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "validation_error:") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": strings.TrimPrefix(err.Error(), "validation_error: "),
			})
			return
		}
		if strings.HasPrefix(err.Error(), "conflict:") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "conflict",
				"message": strings.TrimPrefix(err.Error(), "conflict: "),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to create domain",
		})
		return
	}

	c.JSON(http.StatusCreated, domain)
}

// GetDomain godoc
// @Summary      Get a domain
// @Description  Get domain by ID
// @Tags         domains
// @Produce      json
// @Param        id   path      int  true  "Domain ID"
// @Success      200  {object}  models.Domain
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /domains/{id} [get]
func (h *DomainHandler) GetDomain(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	domain, err := h.service.GetDomain(c.Request.Context(), id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not_found:") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": strings.TrimPrefix(err.Error(), "not_found: "),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to get domain",
		})
		return
	}

	c.JSON(http.StatusOK, domain)
}

// ListDomains godoc
// @Summary      List domains
// @Description  Get all domains with pagination
// @Tags         domains
// @Produce      json
// @Param        page  query     int  false  "Page number"  default(1)
// @Param        size  query     int  false  "Page size"    default(20)
// @Success      200   {object}  models.DomainListResponse
// @Failure      500   {object}  map[string]interface{}
// @Router       /domains [get]
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
			"error":   "internal_error",
			"message": "Failed to list domains",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateDomain godoc
// @Summary      Update a domain
// @Description  Update domain description by ID
// @Tags         domains
// @Accept       json
// @Produce      json
// @Param        id      path      int                        true  "Domain ID"
// @Param        domain  body      models.UpdateDomainRequest true  "Updated domain data"
// @Success      200     {object}  models.Domain
// @Failure      400     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /domains/{id} [put]
func (h *DomainHandler) UpdateDomain(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	var req models.UpdateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}

	domain, err := h.service.UpdateDomain(c.Request.Context(), id, &req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "validation_error:") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": strings.TrimPrefix(err.Error(), "validation_error: "),
			})
			return
		}
		if strings.HasPrefix(err.Error(), "not_found:") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": strings.TrimPrefix(err.Error(), "not_found: "),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to update domain",
		})
		return
	}

	c.JSON(http.StatusOK, domain)
}

// DeleteDomain godoc
// @Summary      Delete a domain
// @Description  Delete domain by ID
// @Tags         domains
// @Param        id  path  int  true  "Domain ID"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /domains/{id} [delete]
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	err = h.service.DeleteDomain(c.Request.Context(), id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not_found:") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": strings.TrimPrefix(err.Error(), "not_found: "),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to delete domain",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
