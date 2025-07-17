package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
)

type NodeAttributeService interface {
	CreateNodeAttribute(nodeID int, req *models.CreateNodeAttributeRequest) (*models.NodeAttribute, error)
	GetNodeAttributesByNodeID(nodeID int) ([]models.NodeAttributeWithInfo, error)
	GetNodeAttributeByID(id int) (*models.NodeAttributeWithInfo, error)
	UpdateNodeAttribute(id int, req *models.UpdateNodeAttributeRequest) (*models.NodeAttribute, error)
	DeleteNodeAttribute(id int) error
	DeleteNodeAttributeByNodeAndAttribute(nodeID, attributeID int) error
}

type NodeAttributeHandler struct {
	*BaseHandler
	nodeAttributeService NodeAttributeService
}

func NewNodeAttributeHandler(nodeAttributeService NodeAttributeService) *NodeAttributeHandler {
	return &NodeAttributeHandler{
		BaseHandler:          NewBaseHandler(),
		nodeAttributeService: nodeAttributeService,
	}
}

func (h *NodeAttributeHandler) CreateNodeAttribute(c *gin.Context) {
	nodeID, err := h.ParseIntParam(c, "url_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	var req models.CreateNodeAttributeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	nodeAttribute, err := h.nodeAttributeService.CreateNodeAttribute(nodeID, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, nodeAttribute)
}

func (h *NodeAttributeHandler) GetNodeAttributesByNode(c *gin.Context) {
	nodeID, err := h.ParseIntParam(c, "url_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	attributes, err := h.nodeAttributeService.GetNodeAttributesByNodeID(nodeID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attributes": attributes,
	})
}

func (h *NodeAttributeHandler) GetNodeAttribute(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	attribute, err := h.nodeAttributeService.GetNodeAttributeByID(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, attribute)
}

func (h *NodeAttributeHandler) UpdateNodeAttribute(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	var req models.UpdateNodeAttributeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	attribute, err := h.nodeAttributeService.UpdateNodeAttribute(id, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, attribute)
}

func (h *NodeAttributeHandler) DeleteNodeAttribute(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	err = h.nodeAttributeService.DeleteNodeAttribute(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NodeAttributeHandler) DeleteNodeAttributeByNodeAndAttribute(c *gin.Context) {
	nodeID, err := h.ParseIntParam(c, "url_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	attributeID, err := h.ParseIntParam(c, "attribute_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	err = h.nodeAttributeService.DeleteNodeAttributeByNodeAndAttribute(nodeID, attributeID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NodeAttributeHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Individual node attribute operations
		urlAttributes := api.Group("/url-attributes")
		{
			urlAttributes.GET("/:id", h.GetNodeAttribute)
			urlAttributes.PUT("/:id", h.UpdateNodeAttribute)
			urlAttributes.DELETE("/:id", h.DeleteNodeAttribute)
		}

		// Node-specific attribute operations
		urls := api.Group("/urls")
		{
			urlNodeAttributes := urls.Group("/:url_id/attributes")
			{
				urlNodeAttributes.POST("", h.CreateNodeAttribute)
				urlNodeAttributes.GET("", h.GetNodeAttributesByNode)
				urlNodeAttributes.DELETE("/:attribute_id", h.DeleteNodeAttributeByNodeAndAttribute)
			}
		}
	}
}