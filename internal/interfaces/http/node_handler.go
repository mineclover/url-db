package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
)

type NodeService interface {
	CreateNode(domainID int, req *models.CreateNodeRequest) (*models.Node, error)
	GetNodeByID(id int) (*models.Node, error)
	GetNodesByDomainID(domainID, page, size int) (*models.NodeListResponse, error)
	FindNodeByURL(domainID int, req *models.FindNodeByURLRequest) (*models.Node, error)
	UpdateNode(id int, req *models.UpdateNodeRequest) (*models.Node, error)
	DeleteNode(id int) error
	SearchNodes(domainID int, query string, page, size int) (*models.NodeListResponse, error)
}

type NodeHandler struct {
	*BaseHandler
	nodeService NodeService
}

func NewNodeHandler(nodeService NodeService) *NodeHandler {
	return &NodeHandler{
		BaseHandler: NewBaseHandler(),
		nodeService: nodeService,
	}
}

func (h *NodeHandler) CreateNode(c *gin.Context) {
	domainID, err := h.ParseIntParam(c, "domain_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	var req models.CreateNodeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	node, err := h.nodeService.CreateNode(domainID, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, node)
}

func (h *NodeHandler) GetNodesByDomain(c *gin.Context) {
	domainID, err := h.ParseIntParam(c, "domain_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	page := h.ParseIntQuery(c, "page", 1)
	size := h.ParseIntQuery(c, "size", 20)
	search := h.GetStringQuery(c, "search")

	if size > 100 {
		size = 100
	}

	var response *models.NodeListResponse
	if search != "" {
		response, err = h.nodeService.SearchNodes(domainID, search, page, size)
	} else {
		response, err = h.nodeService.GetNodesByDomainID(domainID, page, size)
	}

	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *NodeHandler) GetNode(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	node, err := h.nodeService.GetNodeByID(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) UpdateNode(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	var req models.UpdateNodeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	node, err := h.nodeService.UpdateNode(id, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) DeleteNode(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	err = h.nodeService.DeleteNode(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NodeHandler) FindNodeByURL(c *gin.Context) {
	domainID, err := h.ParseIntParam(c, "domain_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	var req models.FindNodeByURLRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	node, err := h.nodeService.FindNodeByURL(domainID, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Individual node operations
		urls := api.Group("/urls")
		{
			urls.GET("/:id", h.GetNode)
			urls.PUT("/:id", h.UpdateNode)
			urls.DELETE("/:id", h.DeleteNode)
		}

		// Domain-specific node operations
		domains := api.Group("/domains")
		{
			domainUrls := domains.Group("/:domain_id/urls")
			{
				domainUrls.POST("", h.CreateNode)
				domainUrls.GET("", h.GetNodesByDomain)
				domainUrls.POST("/find", h.FindNodeByURL)
			}
		}
	}
}
