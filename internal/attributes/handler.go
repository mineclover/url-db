package attributes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
)

// AttributeHandler handles HTTP requests for attributes
type AttributeHandler struct {
	service AttributeService
}

// NewAttributeHandler creates a new attribute handler
func NewAttributeHandler(service AttributeService) *AttributeHandler {
	return &AttributeHandler{service: service}
}

// RegisterRoutes registers attribute routes
func (h *AttributeHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/domains/:domain_id/attributes", h.CreateAttribute)
	router.GET("/domains/:domain_id/attributes", h.ListAttributes)
	router.GET("/attributes/:id", h.GetAttribute)
	router.PUT("/attributes/:id", h.UpdateAttribute)
	router.DELETE("/attributes/:id", h.DeleteAttribute)
}

// CreateAttribute godoc
// @Summary      Create a new attribute
// @Description  Create a new attribute for a domain
// @Tags         attributes
// @Accept       json
// @Produce      json
// @Param        domain_id  path      int                           true  "Domain ID"
// @Param        attribute  body      models.CreateAttributeRequest true  "Attribute data"
// @Success      201        {object}  models.Attribute
// @Failure      400        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      409        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Router       /domains/{domain_id}/attributes [post]
func (h *AttributeHandler) CreateAttribute(c *gin.Context) {
	domainIDStr := c.Param("domain_id")
	domainID, err := strconv.Atoi(domainIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	var req models.CreateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}

	attribute, err := h.service.CreateAttribute(c.Request.Context(), domainID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, attribute)
}

// GetAttribute godoc
// @Summary      Get an attribute
// @Description  Get attribute by ID
// @Tags         attributes
// @Produce      json
// @Param        id  path      int  true  "Attribute ID"
// @Success      200 {object}  models.Attribute
// @Failure      400 {object}  map[string]interface{}
// @Failure      404 {object}  map[string]interface{}
// @Failure      500 {object}  map[string]interface{}
// @Router       /attributes/{id} [get]
func (h *AttributeHandler) GetAttribute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid attribute ID",
		})
		return
	}

	attribute, err := h.service.GetAttribute(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, attribute)
}

// ListAttributes godoc
// @Summary      List attributes
// @Description  Get all attributes for a domain
// @Tags         attributes
// @Produce      json
// @Param        domain_id  path      int  true  "Domain ID"
// @Success      200        {object}  models.AttributeListResponse
// @Failure      400        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Router       /domains/{domain_id}/attributes [get]
func (h *AttributeHandler) ListAttributes(c *gin.Context) {
	domainIDStr := c.Param("domain_id")
	domainID, err := strconv.Atoi(domainIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	response, err := h.service.ListAttributes(c.Request.Context(), domainID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateAttribute godoc
// @Summary      Update an attribute
// @Description  Update attribute description by ID
// @Tags         attributes
// @Accept       json
// @Produce      json
// @Param        id         path      int                           true  "Attribute ID"
// @Param        attribute  body      models.UpdateAttributeRequest true  "Updated attribute data"
// @Success      200        {object}  models.Attribute
// @Failure      400        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Router       /attributes/{id} [put]
func (h *AttributeHandler) UpdateAttribute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid attribute ID",
		})
		return
	}

	var req models.UpdateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}

	attribute, err := h.service.UpdateAttribute(c.Request.Context(), id, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, attribute)
}

// DeleteAttribute godoc
// @Summary      Delete an attribute
// @Description  Delete attribute by ID (only if no associated values exist)
// @Tags         attributes
// @Param        id  path  int  true  "Attribute ID"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /attributes/{id} [delete]
func (h *AttributeHandler) DeleteAttribute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid attribute ID",
		})
		return
	}

	err = h.service.DeleteAttribute(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleError handles service errors and converts them to appropriate HTTP responses
func (h *AttributeHandler) handleError(c *gin.Context, err error) {
	switch err {
	case ErrAttributeNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "Attribute not found",
		})
	case ErrAttributeAlreadyExists:
		c.JSON(http.StatusConflict, gin.H{
			"error":   "conflict",
			"message": "Attribute already exists",
		})
	case ErrAttributeTypeInvalid:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid attribute type",
		})
	case ErrAttributeHasValues:
		c.JSON(http.StatusConflict, gin.H{
			"error":   "conflict",
			"message": "Cannot delete attribute with existing values",
		})
	case ErrDomainNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "Domain not found",
		})
	case ErrAttributeNameRequired:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Attribute name is required",
		})
	case ErrAttributeNameTooLong:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Attribute name too long",
		})
	case ErrAttributeTypeRequired:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Attribute type is required",
		})
	case ErrDescriptionTooLong:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Description too long",
		})
	default:
		if strings.Contains(err.Error(), "validation error") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Internal server error",
			})
		}
	}
}
