package nodeattributes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"url-db/internal/models"
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

// CreateNodeAttribute godoc
// @Summary      Create a node attribute
// @Description  Create a new attribute value for a URL node
// @Tags         node-attributes
// @Accept       json
// @Produce      json
// @Param        url_id     path      int                               true  "URL Node ID"
// @Param        attribute  body      models.CreateNodeAttributeRequest true  "Node attribute data"
// @Success      201        {object}  models.NodeAttribute
// @Failure      400        {object}  map[string]interface{}
// @Failure      409        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Router       /urls/{url_id}/attributes [post]
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

// GetNodeAttributesByNodeID godoc
// @Summary      Get node attributes
// @Description  Get all attribute values for a URL node
// @Tags         node-attributes
// @Produce      json
// @Param        url_id  path      int  true  "URL Node ID"
// @Success      200     {object}  models.NodeAttributeListResponse
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /urls/{url_id}/attributes [get]
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

// GetNodeAttributeByID godoc
// @Summary      Get a node attribute
// @Description  Get a specific node attribute by ID
// @Tags         node-attributes
// @Produce      json
// @Param        id  path      int  true  "Node Attribute ID"
// @Success      200 {object}  models.NodeAttribute
// @Failure      400 {object}  map[string]interface{}
// @Failure      404 {object}  map[string]interface{}
// @Failure      500 {object}  map[string]interface{}
// @Router       /url-attributes/{id} [get]
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

// UpdateNodeAttribute godoc
// @Summary      Update a node attribute
// @Description  Update a node attribute value and order index by ID
// @Tags         node-attributes
// @Accept       json
// @Produce      json
// @Param        id         path      int                               true  "Node Attribute ID"
// @Param        attribute  body      models.UpdateNodeAttributeRequest true  "Updated node attribute data"
// @Success      200        {object}  models.NodeAttribute
// @Failure      400        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      409        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Router       /url-attributes/{id} [put]
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

// DeleteNodeAttribute godoc
// @Summary      Delete a node attribute
// @Description  Delete a node attribute by ID
// @Tags         node-attributes
// @Param        id  path  int  true  "Node Attribute ID"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /url-attributes/{id} [delete]
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

// DeleteNodeAttributeByNodeIDAndAttributeID godoc
// @Summary      Delete a node attribute by node and attribute ID
// @Description  Delete a specific node attribute by URL node ID and attribute ID
// @Tags         node-attributes
// @Param        url_id       path  int  true  "URL Node ID"
// @Param        attribute_id path  int  true  "Attribute ID"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /urls/{url_id}/attributes/{attribute_id} [delete]
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
