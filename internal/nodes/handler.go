package nodes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
)

type NodeHandler struct {
	service NodeService
}

func NewNodeHandler(service NodeService) *NodeHandler {
	return &NodeHandler{service: service}
}

// CreateNode godoc
// @Summary      Create a new URL node
// @Description  Create a new URL node in a domain
// @Tags         nodes
// @Accept       json
// @Produce      json
// @Param        domain_id  path      int                        true   "Domain ID"
// @Param        node       body      models.CreateNodeRequest  true   "Node data"
// @Success      201        {object}  models.Node
// @Failure      400        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      409        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Router       /domains/{domain_id}/urls [post]
func (h *NodeHandler) CreateNode(c *gin.Context) {
	domainIDStr := c.Param("domain_id")
	domainID, err := strconv.Atoi(domainIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	var req models.CreateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}

	node, err := h.service.CreateNode(domainID, &req)
	if err != nil {
		switch err {
		case ErrNodeDomainNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Domain not found",
			})
		case ErrNodeAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error":   "conflict",
				"message": "Node already exists",
			})
		case ErrNodeURLInvalid:
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": "Invalid URL",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, node)
}

// GetNodesByDomain godoc
// @Summary      Get nodes by domain
// @Description  Get all nodes in a domain with pagination and optional search
// @Tags         nodes
// @Produce      json
// @Param        domain_id  path   int     true   "Domain ID"
// @Param        page       query  int     false  "Page number" default(1)
// @Param        size       query  int     false  "Page size" default(20)
// @Param        search     query  string  false  "Search term for URL content"
// @Success      200        {object}  models.NodeListResponse
// @Failure      400        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Router       /domains/{domain_id}/urls [get]
func (h *NodeHandler) GetNodesByDomain(c *gin.Context) {
	domainIDStr := c.Param("domain_id")
	domainID, err := strconv.Atoi(domainIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	search := c.Query("search")

	var response *models.NodeListResponse
	if search != "" {
		response, err = h.service.SearchNodes(domainID, search, page, size)
	} else {
		response, err = h.service.GetNodesByDomainID(domainID, page, size)
	}

	if err != nil {
		switch err {
		case ErrNodeDomainNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Domain not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetNode godoc
// @Summary      Get a node
// @Description  Get node by ID
// @Tags         nodes
// @Produce      json
// @Param        id  path      int  true  "Node ID"
// @Success      200 {object}  models.Node
// @Failure      400 {object}  map[string]interface{}
// @Failure      404 {object}  map[string]interface{}
// @Failure      500 {object}  map[string]interface{}
// @Router       /urls/{id} [get]
func (h *NodeHandler) GetNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid node ID",
		})
		return
	}

	node, err := h.service.GetNodeByID(id)
	if err != nil {
		switch err {
		case ErrNodeNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Node not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, node)
}

// UpdateNode godoc
// @Summary      Update a node
// @Description  Update node title and description by ID
// @Tags         nodes
// @Accept       json
// @Produce      json
// @Param        id    path      int                       true  "Node ID"
// @Param        node  body      models.UpdateNodeRequest  true  "Updated node data"
// @Success      200   {object}  models.Node
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /urls/{id} [put]
func (h *NodeHandler) UpdateNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid node ID",
		})
		return
	}

	var req models.UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}

	node, err := h.service.UpdateNode(id, &req)
	if err != nil {
		switch err {
		case ErrNodeNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Node not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, node)
}

// DeleteNode godoc
// @Summary      Delete a node
// @Description  Delete node by ID
// @Tags         nodes
// @Param        id  path  int  true  "Node ID"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /urls/{id} [delete]
func (h *NodeHandler) DeleteNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid node ID",
		})
		return
	}

	err = h.service.DeleteNode(id)
	if err != nil {
		switch err {
		case ErrNodeNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Node not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Internal server error",
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// FindNodeByURL godoc
// @Summary      Find node by URL
// @Description  Find a node by its URL within a domain
// @Tags         nodes
// @Accept       json
// @Produce      json
// @Param        domain_id  path      int                         true  "Domain ID"
// @Param        request    body      models.FindNodeByURLRequest true  "URL to search for"
// @Success      200        {object}  models.Node
// @Failure      400        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Router       /domains/{domain_id}/urls/find [post]
func (h *NodeHandler) FindNodeByURL(c *gin.Context) {
	domainIDStr := c.Param("domain_id")
	domainID, err := strconv.Atoi(domainIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid domain ID",
		})
		return
	}

	var req models.FindNodeByURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}

	node, err := h.service.FindNodeByURL(domainID, &req)
	if err != nil {
		switch err {
		case ErrNodeNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Node not found",
			})
		case ErrNodeDomainNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Domain not found",
			})
		case ErrNodeURLInvalid:
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": "Invalid URL",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Node creation by domain
		api.POST("/domains/:domain_id/urls", h.CreateNode)

		// Node list by domain
		api.GET("/domains/:domain_id/urls", h.GetNodesByDomain)

		// Find node by URL
		api.POST("/domains/:domain_id/urls/find", h.FindNodeByURL)

		// Individual node operations
		api.GET("/urls/:id", h.GetNode)
		api.PUT("/urls/:id", h.UpdateNode)
		api.DELETE("/urls/:id", h.DeleteNode)
	}
}
