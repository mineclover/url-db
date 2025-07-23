package handler

import (
	"net/http"
	"strconv"

	"url-db/internal/application/dto/request"
	"url-db/internal/application/usecase/attribute"

	"github.com/gin-gonic/gin"
)

// AttributeHandler handles HTTP requests for attribute operations
type AttributeHandler struct {
	createUseCase *attribute.CreateAttributeUseCase
	listUseCase   *attribute.ListAttributesUseCase
}

// NewAttributeHandler creates a new attribute handler
func NewAttributeHandler(
	createUC *attribute.CreateAttributeUseCase,
	listUC *attribute.ListAttributesUseCase,
) *AttributeHandler {
	return &AttributeHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
	}
}

// CreateAttribute handles POST /api/domains/{domain_id}/attributes
func (h *AttributeHandler) CreateAttribute(c *gin.Context) {
	domainIDStr := c.Param("domain_id")

	domainID, err := strconv.Atoi(domainIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain ID"})
		return
	}

	var req request.CreateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set domain ID from URL parameter
	req.DomainID = domainID

	response, err := h.createUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		// In a real implementation, you'd have proper error handling/mapping
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ListAttributes handles GET /api/domains/{domain_id}/attributes
func (h *AttributeHandler) ListAttributes(c *gin.Context) {
	domainIDStr := c.Param("domain_id")

	domainID, err := strconv.Atoi(domainIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain ID"})
		return
	}

	response, err := h.listUseCase.Execute(c.Request.Context(), domainID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
