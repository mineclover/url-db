package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
)

type MCPService interface {
	CreateMCPNode(req *models.CreateMCPNodeRequest) (*models.MCPNode, error)
	GetMCPNodeByCompositeID(compositeID string) (*models.MCPNode, error)
	GetMCPNodes(domainName string, page, size int, search string) (*models.MCPNodeListResponse, error)
	UpdateMCPNode(compositeID string, req *models.UpdateMCPNodeRequest) (*models.MCPNode, error)
	DeleteMCPNode(compositeID string) error
	FindMCPNodeByURL(req *models.FindMCPNodeRequest) (*models.MCPNode, error)
	BatchGetMCPNodes(req *models.BatchMCPNodeRequest) (*models.BatchMCPNodeResponse, error)
	GetMCPDomains() ([]models.MCPDomain, error)
	CreateMCPDomain(req *models.CreateMCPDomainRequest) (*models.MCPDomain, error)
	GetMCPNodeAttributes(compositeID string) (*models.MCPNodeAttributesResponse, error)
	SetMCPNodeAttributes(compositeID string, req *models.SetMCPNodeAttributesRequest) (*models.MCPNodeAttributesResponse, error)
	GetMCPServerInfo() (*models.MCPServerInfo, error)
}

type MCPHandler struct {
	*BaseHandler
	mcpService MCPService
}

func NewMCPHandler(mcpService MCPService) *MCPHandler {
	return &MCPHandler{
		BaseHandler: NewBaseHandler(),
		mcpService:  mcpService,
	}
}

func (h *MCPHandler) CreateMCPNode(c *gin.Context) {
	var req models.CreateMCPNodeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	node, err := h.mcpService.CreateMCPNode(&req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, node)
}

func (h *MCPHandler) GetMCPNode(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.HandleError(c, NewValidationError("Missing composite_id parameter", nil))
		return
	}

	node, err := h.mcpService.GetMCPNodeByCompositeID(compositeID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *MCPHandler) GetMCPNodes(c *gin.Context) {
	domainName := h.GetStringQuery(c, "domain_name")
	page := h.ParseIntQuery(c, "page", 1)
	size := h.ParseIntQuery(c, "size", 20)
	search := h.GetStringQuery(c, "search")

	if size > 100 {
		size = 100
	}

	response, err := h.mcpService.GetMCPNodes(domainName, page, size, search)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) UpdateMCPNode(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.HandleError(c, NewValidationError("Missing composite_id parameter", nil))
		return
	}

	var req models.UpdateMCPNodeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	node, err := h.mcpService.UpdateMCPNode(compositeID, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *MCPHandler) DeleteMCPNode(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.HandleError(c, NewValidationError("Missing composite_id parameter", nil))
		return
	}

	err := h.mcpService.DeleteMCPNode(compositeID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *MCPHandler) FindMCPNodeByURL(c *gin.Context) {
	var req models.FindMCPNodeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	node, err := h.mcpService.FindMCPNodeByURL(&req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *MCPHandler) BatchGetMCPNodes(c *gin.Context) {
	var req models.BatchMCPNodeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	response, err := h.mcpService.BatchGetMCPNodes(&req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) GetMCPDomains(c *gin.Context) {
	domains, err := h.mcpService.GetMCPDomains()
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domains": domains,
	})
}

func (h *MCPHandler) CreateMCPDomain(c *gin.Context) {
	var req models.CreateMCPDomainRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	domain, err := h.mcpService.CreateMCPDomain(&req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, domain)
}

func (h *MCPHandler) GetMCPNodeAttributes(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.HandleError(c, NewValidationError("Missing composite_id parameter", nil))
		return
	}

	response, err := h.mcpService.GetMCPNodeAttributes(compositeID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) SetMCPNodeAttributes(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.HandleError(c, NewValidationError("Missing composite_id parameter", nil))
		return
	}

	var req models.SetMCPNodeAttributesRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	response, err := h.mcpService.SetMCPNodeAttributes(compositeID, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) GetMCPServerInfo(c *gin.Context) {
	info, err := h.mcpService.GetMCPServerInfo()
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, info)
}

func (h *MCPHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		mcp := api.Group("/mcp")
		{
			// Node operations
			nodes := mcp.Group("/nodes")
			{
				nodes.POST("", h.CreateMCPNode)
				nodes.GET("", h.GetMCPNodes)
				nodes.GET("/:composite_id", h.GetMCPNode)
				nodes.PUT("/:composite_id", h.UpdateMCPNode)
				nodes.DELETE("/:composite_id", h.DeleteMCPNode)
				nodes.POST("/find", h.FindMCPNodeByURL)
				nodes.POST("/batch", h.BatchGetMCPNodes)
				
				// Node attributes
				nodes.GET("/:composite_id/attributes", h.GetMCPNodeAttributes)
				nodes.PUT("/:composite_id/attributes", h.SetMCPNodeAttributes)
			}
			
			// Domain operations
			domains := mcp.Group("/domains")
			{
				domains.GET("", h.GetMCPDomains)
				domains.POST("", h.CreateMCPDomain)
			}
			
			// Server info
			server := mcp.Group("/server")
			{
				server.GET("/info", h.GetMCPServerInfo)
			}
		}
	}
}