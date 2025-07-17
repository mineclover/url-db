package nodeattributes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"internal/models"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Node attribute routes
		api.POST("/urls/:url_id/attributes", h.CreateNodeAttribute)
		api.GET("/urls/:url_id/attributes", h.GetNodeAttributesByNodeID)
		api.GET("/url-attributes/:id", h.GetNodeAttributeByID)
		api.PUT("/url-attributes/:id", h.UpdateNodeAttribute)
		api.DELETE("/url-attributes/:id", h.DeleteNodeAttribute)
		api.DELETE("/urls/:url_id/attributes/:attribute_id", h.DeleteNodeAttributeByNodeIDAndAttributeID)
	}
}

func (h *Handler) CreateNodeAttribute(c *gin.Context) {
	urlID, err := strconv.Atoi(c.Param("url_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL ID",
			"code":  "validation_error",
		})
		return
	}

	var req models.CreateNodeAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "validation_error",
		})
		return
	}

	nodeAttribute, err := h.service.CreateNodeAttribute(urlID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		switch err {
		case ErrNodeAttributeDomainMismatch:
			statusCode = http.StatusBadRequest
			errorCode = "business_rule_violation"
		case ErrNodeAttributeExists:
			statusCode = http.StatusConflict
			errorCode = "conflict"
		case ErrNodeAttributeValueInvalid, ErrInvalidAttributeType, ErrOrderIndexRequired, ErrOrderIndexNotAllowed, ErrInvalidOrderIndex:
			statusCode = http.StatusBadRequest
			errorCode = "validation_error"
		case ErrDuplicateOrderIndex:
			statusCode = http.StatusConflict
			errorCode = "constraint_violation"
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"code":  errorCode,
		})
		return
	}

	c.JSON(http.StatusCreated, nodeAttribute)
}

func (h *Handler) GetNodeAttributesByNodeID(c *gin.Context) {
	urlID, err := strconv.Atoi(c.Param("url_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL ID",
			"code":  "validation_error",
		})
		return
	}

	response, err := h.service.GetNodeAttributesByNodeID(urlID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"code":  "internal_error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetNodeAttributeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
			"code":  "validation_error",
		})
		return
	}

	nodeAttribute, err := h.service.GetNodeAttributeByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		if err == ErrNodeAttributeNotFound {
			statusCode = http.StatusNotFound
			errorCode = "not_found"
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"code":  errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, nodeAttribute)
}

func (h *Handler) UpdateNodeAttribute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
			"code":  "validation_error",
		})
		return
	}

	var req models.UpdateNodeAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "validation_error",
		})
		return
	}

	nodeAttribute, err := h.service.UpdateNodeAttribute(id, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		switch err {
		case ErrNodeAttributeNotFound:
			statusCode = http.StatusNotFound
			errorCode = "not_found"
		case ErrNodeAttributeValueInvalid, ErrInvalidAttributeType, ErrOrderIndexRequired, ErrOrderIndexNotAllowed, ErrInvalidOrderIndex:
			statusCode = http.StatusBadRequest
			errorCode = "validation_error"
		case ErrDuplicateOrderIndex:
			statusCode = http.StatusConflict
			errorCode = "constraint_violation"
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"code":  errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, nodeAttribute)
}

func (h *Handler) DeleteNodeAttribute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
			"code":  "validation_error",
		})
		return
	}

	err = h.service.DeleteNodeAttribute(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		if err == ErrNodeAttributeNotFound {
			statusCode = http.StatusNotFound
			errorCode = "not_found"
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"code":  errorCode,
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteNodeAttributeByNodeIDAndAttributeID(c *gin.Context) {
	urlID, err := strconv.Atoi(c.Param("url_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL ID",
			"code":  "validation_error",
		})
		return
	}

	attributeID, err := strconv.Atoi(c.Param("attribute_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Attribute ID",
			"code":  "validation_error",
		})
		return
	}

	err = h.service.DeleteNodeAttributeByNodeIDAndAttributeID(urlID, attributeID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		if err == ErrNodeAttributeNotFound {
			statusCode = http.StatusNotFound
			errorCode = "not_found"
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"code":  errorCode,
		})
		return
	}

	c.Status(http.StatusNoContent)
}