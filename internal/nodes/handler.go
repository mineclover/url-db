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