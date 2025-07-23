package handlers

import (
	"net/http"

	"url-db/internal/models"

	"github.com/gin-gonic/gin"
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
	// Domain attribute management methods
	ListDomainAttributes(domainName string) (*models.MCPDomainAttributeListResponse, error)
	CreateDomainAttribute(domainName string, req *models.CreateAttributeRequest) (*models.MCPDomainAttribute, error)
	GetDomainAttribute(compositeID string) (*models.MCPDomainAttribute, error)
	UpdateDomainAttribute(compositeID string, req *models.UpdateAttributeRequest) (*models.MCPDomainAttribute, error)
	DeleteDomainAttribute(compositeID string) error
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

// CreateMCPNode godoc
// @Summary      Create a new MCP node
// @Description  Create a new URL node using MCP composite key format
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        node  body      models.CreateMCPNodeRequest  true  "MCP node data"
// @Success      201   {object}  models.MCPNode
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /mcp/nodes [post]
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

// GetMCPNode godoc
// @Summary      Get an MCP node
// @Description  Get MCP node by composite ID
// @Tags         mcp
// @Produce      json
// @Param        composite_id  path      string  true  "Composite ID (domain::url)"
// @Success      200           {object}  models.MCPNode
// @Failure      400           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /mcp/nodes/{composite_id} [get]
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

// GetMCPNodes godoc
// @Summary      List MCP nodes
// @Description  Get all MCP nodes with pagination, domain filtering, and search
// @Tags         mcp
// @Produce      json
// @Param        domain_name  query  string  false  "Filter by domain name"
// @Param        page         query  int     false  "Page number" default(1)
// @Param        size         query  int     false  "Page size (max 100)" default(20)
// @Param        search       query  string  false  "Search term for URL content"
// @Success      200          {object}  models.MCPNodeListResponse
// @Failure      500          {object}  map[string]interface{}
// @Router       /mcp/nodes [get]
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

// UpdateMCPNode godoc
// @Summary      Update an MCP node
// @Description  Update MCP node title and description by composite ID
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        composite_id  path      string                       true  "Composite ID (domain::url)"
// @Param        node          body      models.UpdateMCPNodeRequest  true  "Updated MCP node data"
// @Success      200           {object}  models.MCPNode
// @Failure      400           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /mcp/nodes/{composite_id} [put]
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

// DeleteMCPNode godoc
// @Summary      Delete an MCP node
// @Description  Delete MCP node by composite ID
// @Tags         mcp
// @Param        composite_id  path  string  true  "Composite ID (domain::url)"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /mcp/nodes/{composite_id} [delete]
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

// FindMCPNodeByURL godoc
// @Summary      Find MCP node by URL
// @Description  Find an MCP node by its URL within a domain
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        request  body      models.FindMCPNodeRequest  true  "Domain and URL to search for"
// @Success      200      {object}  models.MCPNode
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /mcp/nodes/find [post]
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

// BatchGetMCPNodes godoc
// @Summary      Batch get MCP nodes
// @Description  Get multiple MCP nodes by their composite IDs in a single request
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        request  body      models.BatchMCPNodeRequest  true  "List of composite IDs"
// @Success      200      {object}  models.BatchMCPNodeResponse
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /mcp/nodes/batch [post]
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

// GetMCPDomains godoc
// @Summary      List MCP domains
// @Description  Get all domains with their node counts in MCP format
// @Tags         mcp
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /mcp/domains [get]
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

// CreateMCPDomain godoc
// @Summary      Create a new MCP domain
// @Description  Create a new domain using MCP format
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        domain  body      models.CreateMCPDomainRequest  true  "MCP domain data"
// @Success      201     {object}  models.MCPDomain
// @Failure      400     {object}  map[string]interface{}
// @Failure      409     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /mcp/domains [post]
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

// GetMCPNodeAttributes godoc
// @Summary      Get MCP node attributes
// @Description  Get all attribute values for an MCP node by composite ID
// @Tags         mcp
// @Produce      json
// @Param        composite_id  path      string  true  "Composite ID (domain::url)"
// @Success      200           {object}  models.MCPNodeAttributesResponse
// @Failure      400           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /mcp/nodes/{composite_id}/attributes [get]
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

// SetMCPNodeAttributes godoc
// @Summary      Set MCP node attributes
// @Description  Set or update attribute values for an MCP node by composite ID
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        composite_id  path      string                              true  "Composite ID (domain::url)"
// @Param        attributes    body      models.SetMCPNodeAttributesRequest  true  "Attribute data to set"
// @Success      200           {object}  models.MCPNodeAttributesResponse
// @Failure      400           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /mcp/nodes/{composite_id}/attributes [put]
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

// GetMCPServerInfo godoc
// @Summary      Get MCP server information
// @Description  Get server capabilities and configuration information
// @Tags         mcp
// @Produce      json
// @Success      200  {object}  models.MCPServerInfo
// @Failure      500  {object}  map[string]interface{}
// @Router       /mcp/server/info [get]
func (h *MCPHandler) GetMCPServerInfo(c *gin.Context) {
	info, err := h.mcpService.GetMCPServerInfo()
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, info)
}

// ListDomainAttributes godoc
// @Summary      List domain attributes
// @Description  Get all attribute definitions for a domain
// @Tags         mcp
// @Produce      json
// @Param        domain_name  path      string  true  "Domain name"
// @Success      200          {object}  models.MCPDomainAttributeListResponse
// @Failure      400          {object}  map[string]interface{}
// @Failure      404          {object}  map[string]interface{}
// @Failure      500          {object}  map[string]interface{}
// @Router       /mcp/domains/{domain_name}/attributes [get]
func (h *MCPHandler) ListDomainAttributes(c *gin.Context) {
	domainName := c.Param("domain_name")
	if domainName == "" {
		h.HandleError(c, NewValidationError("Missing domain_name parameter", nil))
		return
	}

	response, err := h.mcpService.ListDomainAttributes(domainName)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateDomainAttribute godoc
// @Summary      Create domain attribute
// @Description  Create a new attribute definition for a domain
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        domain_name  path      string                        true  "Domain name"
// @Param        attribute    body      models.CreateAttributeRequest  true  "Attribute data"
// @Success      201          {object}  models.MCPDomainAttribute
// @Failure      400          {object}  map[string]interface{}
// @Failure      404          {object}  map[string]interface{}
// @Failure      409          {object}  map[string]interface{}
// @Failure      500          {object}  map[string]interface{}
// @Router       /mcp/domains/{domain_name}/attributes [post]
func (h *MCPHandler) CreateDomainAttribute(c *gin.Context) {
	domainName := c.Param("domain_name")
	if domainName == "" {
		h.HandleError(c, NewValidationError("Missing domain_name parameter", nil))
		return
	}

	var req models.CreateAttributeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	attribute, err := h.mcpService.CreateDomainAttribute(domainName, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, attribute)
}

// GetDomainAttribute godoc
// @Summary      Get domain attribute
// @Description  Get a specific attribute definition by composite ID
// @Tags         mcp
// @Produce      json
// @Param        composite_id  path      string  true  "Composite ID (domain:attribute_id)"
// @Success      200           {object}  models.MCPDomainAttribute
// @Failure      400           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /mcp/attributes/{composite_id} [get]
func (h *MCPHandler) GetDomainAttribute(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.HandleError(c, NewValidationError("Missing composite_id parameter", nil))
		return
	}

	attribute, err := h.mcpService.GetDomainAttribute(compositeID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, attribute)
}

// UpdateDomainAttribute godoc
// @Summary      Update domain attribute
// @Description  Update an attribute definition by composite ID
// @Tags         mcp
// @Accept       json
// @Produce      json
// @Param        composite_id  path      string                        true  "Composite ID (domain:attribute_id)"
// @Param        attribute     body      models.UpdateAttributeRequest  true  "Attribute update data"
// @Success      200           {object}  models.MCPDomainAttribute
// @Failure      400           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /mcp/attributes/{composite_id} [put]
func (h *MCPHandler) UpdateDomainAttribute(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.HandleError(c, NewValidationError("Missing composite_id parameter", nil))
		return
	}

	var req models.UpdateAttributeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	attribute, err := h.mcpService.UpdateDomainAttribute(compositeID, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, attribute)
}

// DeleteDomainAttribute godoc
// @Summary      Delete domain attribute
// @Description  Delete an attribute definition by composite ID
// @Tags         mcp
// @Produce      json
// @Param        composite_id  path      string  true  "Composite ID (domain:attribute_id)"
// @Success      200           {object}  map[string]interface{}
// @Failure      400           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /mcp/attributes/{composite_id} [delete]
func (h *MCPHandler) DeleteDomainAttribute(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.HandleError(c, NewValidationError("Missing composite_id parameter", nil))
		return
	}

	err := h.mcpService.DeleteDomainAttribute(compositeID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attribute deleted successfully"})
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
				domains.GET("/:domain_name/attributes", h.ListDomainAttributes)
				domains.POST("/:domain_name/attributes", h.CreateDomainAttribute)
			}

			// Domain attributes (by composite ID)
			attributes := mcp.Group("/attributes")
			{
				attributes.GET("/:composite_id", h.GetDomainAttribute)
				attributes.PUT("/:composite_id", h.UpdateDomainAttribute)
				attributes.DELETE("/:composite_id", h.DeleteDomainAttribute)
			}

			// Server info
			server := mcp.Group("/server")
			{
				server.GET("/info", h.GetMCPServerInfo)
			}
		}
	}
}
