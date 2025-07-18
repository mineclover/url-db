package mcp

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
)

type MCPHandler struct {
	mcpService       MCPService
	batchProcessor   *BatchProcessor
	domainManager    *DomainManager
	attributeManager *AttributeManager
	metadataManager  *MetadataManager
}

func NewMCPHandler(
	mcpService MCPService,
	batchProcessor *BatchProcessor,
	domainManager *DomainManager,
	attributeManager *AttributeManager,
	metadataManager *MetadataManager,
) *MCPHandler {
	return &MCPHandler{
		mcpService:       mcpService,
		batchProcessor:   batchProcessor,
		domainManager:    domainManager,
		attributeManager: attributeManager,
		metadataManager:  metadataManager,
	}
}

func (h *MCPHandler) CreateNode(c *gin.Context) {
	var req models.CreateMCPNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	node, err := h.mcpService.CreateNode(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, node)
}

func (h *MCPHandler) GetNode(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.handleError(c, NewValidationError("composite_id is required"))
		return
	}

	node, err := h.mcpService.GetNode(c.Request.Context(), compositeID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *MCPHandler) UpdateNode(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.handleError(c, NewValidationError("composite_id is required"))
		return
	}

	var req models.UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	node, err := h.mcpService.UpdateNode(c.Request.Context(), compositeID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *MCPHandler) DeleteNode(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.handleError(c, NewValidationError("composite_id is required"))
		return
	}

	if err := h.mcpService.DeleteNode(c.Request.Context(), compositeID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *MCPHandler) ListNodes(c *gin.Context) {
	domainName := c.Query("domain_name")
	search := c.Query("search")

	page, err := h.parseIntQuery(c, "page", 1)
	if err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	size, err := h.parseIntQuery(c, "size", 20)
	if err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	if size > 100 {
		size = 100
	}

	response, err := h.mcpService.ListNodes(c.Request.Context(), domainName, page, size, search)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) FindNodeByURL(c *gin.Context) {
	var req models.FindMCPNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	node, err := h.mcpService.FindNodeByURL(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *MCPHandler) BatchGetNodes(c *gin.Context) {
	var req models.BatchMCPNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	if len(req.CompositeIDs) > 100 {
		h.handleError(c, NewValidationError("batch size cannot exceed 100"))
		return
	}

	response, err := h.batchProcessor.BatchGetNodes(c.Request.Context(), req.CompositeIDs)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) ListDomains(c *gin.Context) {
	response, err := h.domainManager.ListDomains(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) CreateDomain(c *gin.Context) {
	var req models.CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	domain, err := h.domainManager.CreateDomain(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, domain)
}

func (h *MCPHandler) GetDomain(c *gin.Context) {
	domainName := c.Param("domain_name")
	if domainName == "" {
		h.handleError(c, NewValidationError("domain_name is required"))
		return
	}

	domain, err := h.domainManager.GetDomain(c.Request.Context(), domainName)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain)
}

func (h *MCPHandler) UpdateDomain(c *gin.Context) {
	domainName := c.Param("domain_name")
	if domainName == "" {
		h.handleError(c, NewValidationError("domain_name is required"))
		return
	}

	var req models.UpdateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	domain, err := h.domainManager.UpdateDomain(c.Request.Context(), domainName, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain)
}

func (h *MCPHandler) DeleteDomain(c *gin.Context) {
	domainName := c.Param("domain_name")
	if domainName == "" {
		h.handleError(c, NewValidationError("domain_name is required"))
		return
	}

	if err := h.domainManager.DeleteDomain(c.Request.Context(), domainName); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *MCPHandler) GetNodeAttributes(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.handleError(c, NewValidationError("composite_id is required"))
		return
	}

	response, err := h.attributeManager.GetNodeAttributes(c.Request.Context(), compositeID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) SetNodeAttributes(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.handleError(c, NewValidationError("composite_id is required"))
		return
	}

	var req models.SetMCPNodeAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	response, err := h.attributeManager.SetNodeAttributes(c.Request.Context(), compositeID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) AddNodeAttribute(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.handleError(c, NewValidationError("composite_id is required"))
		return
	}

	var req AddAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	response, err := h.attributeManager.AddNodeAttribute(c.Request.Context(), compositeID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) RemoveNodeAttribute(c *gin.Context) {
	compositeID := c.Param("composite_id")
	if compositeID == "" {
		h.handleError(c, NewValidationError("composite_id is required"))
		return
	}

	attributeName := c.Param("attribute_name")
	if attributeName == "" {
		h.handleError(c, NewValidationError("attribute_name is required"))
		return
	}

	response, err := h.attributeManager.RemoveNodeAttribute(c.Request.Context(), compositeID, attributeName)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) GetServerInfo(c *gin.Context) {
	info, err := h.metadataManager.GetServerInfo(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, info)
}

func (h *MCPHandler) GetDetailedServerInfo(c *gin.Context) {
	info, err := h.metadataManager.GetDetailedServerInfo(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, info)
}

func (h *MCPHandler) GetHealthStatus(c *gin.Context) {
	status, err := h.metadataManager.GetHealthStatus(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *MCPHandler) GetStatistics(c *gin.Context) {
	stats, err := h.metadataManager.GetStatistics(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *MCPHandler) GetAPIDocumentation(c *gin.Context) {
	docs, err := h.metadataManager.GetAPIDocumentation(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

func (h *MCPHandler) BatchCreateNodes(c *gin.Context) {
	var requests []models.CreateMCPNodeRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	if len(requests) > 100 {
		h.handleError(c, NewValidationError("batch size cannot exceed 100"))
		return
	}

	response, err := h.batchProcessor.BatchCreateNodes(c.Request.Context(), requests)
	if err != nil {
		h.handleError(c, err)
		return
	}

	statusCode := http.StatusCreated
	if len(response.Failed) > 0 {
		statusCode = http.StatusMultiStatus
	}

	c.JSON(statusCode, response)
}

func (h *MCPHandler) BatchUpdateNodes(c *gin.Context) {
	var requests []BatchUpdateRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	if len(requests) > 100 {
		h.handleError(c, NewValidationError("batch size cannot exceed 100"))
		return
	}

	response, err := h.batchProcessor.BatchUpdateNodes(c.Request.Context(), requests)
	if err != nil {
		h.handleError(c, err)
		return
	}

	statusCode := http.StatusOK
	if len(response.Failed) > 0 {
		statusCode = http.StatusMultiStatus
	}

	c.JSON(statusCode, response)
}

func (h *MCPHandler) BatchDeleteNodes(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	if len(req.CompositeIDs) > 100 {
		h.handleError(c, NewValidationError("batch size cannot exceed 100"))
		return
	}

	response, err := h.batchProcessor.BatchDeleteNodes(c.Request.Context(), req.CompositeIDs)
	if err != nil {
		h.handleError(c, err)
		return
	}

	statusCode := http.StatusOK
	if len(response.Failed) > 0 {
		statusCode = http.StatusMultiStatus
	}

	c.JSON(statusCode, response)
}

func (h *MCPHandler) BatchSetNodeAttributes(c *gin.Context) {
	var requests []BatchAttributeRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	if len(requests) > 100 {
		h.handleError(c, NewValidationError("batch size cannot exceed 100"))
		return
	}

	response, err := h.attributeManager.BatchSetNodeAttributes(c.Request.Context(), requests)
	if err != nil {
		h.handleError(c, err)
		return
	}

	statusCode := http.StatusOK
	if len(response.Failed) > 0 {
		statusCode = http.StatusMultiStatus
	}

	c.JSON(statusCode, response)
}

func (h *MCPHandler) SearchDomains(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		h.handleError(c, NewValidationError("query parameter 'q' is required"))
		return
	}

	domains, err := h.domainManager.GetDomainByPartialName(c.Request.Context(), query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, MCPDomainListResponse{
		Domains: domains,
	})
}

func (h *MCPHandler) GetPopularDomains(c *gin.Context) {
	limit, err := h.parseIntQuery(c, "limit", 10)
	if err != nil {
		h.handleError(c, NewValidationError(err.Error()))
		return
	}

	domains, err := h.domainManager.GetPopularDomains(c.Request.Context(), limit)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, MCPDomainListResponse{
		Domains: domains,
	})
}

func (h *MCPHandler) GetDomainStats(c *gin.Context) {
	domainName := c.Param("domain_name")
	if domainName == "" {
		h.handleError(c, NewValidationError("domain_name is required"))
		return
	}

	stats, err := h.domainManager.GetDomainStats(c.Request.Context(), domainName)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *MCPHandler) parseIntQuery(c *gin.Context, key string, defaultValue int) (int, error) {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (h *MCPHandler) handleError(c *gin.Context, err error) {
	if mcpErr, ok := err.(*MCPError); ok {
		c.JSON(mcpErr.HTTPCode, mcpErr)
		return
	}

	internalErr := NewInternalServerError(err.Error())
	c.JSON(internalErr.HTTPCode, internalErr)
}

func (h *MCPHandler) RegisterRoutes(r *gin.RouterGroup) {
	api := r.Group("/mcp")
	{
		nodes := api.Group("/nodes")
		{
			nodes.GET("", h.ListNodes)
			nodes.POST("", h.CreateNode)
			nodes.POST("/find", h.FindNodeByURL)
			nodes.POST("/batch", h.BatchGetNodes)
			nodes.POST("/batch/create", h.BatchCreateNodes)
			nodes.PUT("/batch/update", h.BatchUpdateNodes)
			nodes.DELETE("/batch/delete", h.BatchDeleteNodes)
			nodes.GET("/:composite_id", h.GetNode)
			nodes.PUT("/:composite_id", h.UpdateNode)
			nodes.DELETE("/:composite_id", h.DeleteNode)

			nodes.GET("/:composite_id/attributes", h.GetNodeAttributes)
			nodes.PUT("/:composite_id/attributes", h.SetNodeAttributes)
			nodes.POST("/:composite_id/attributes", h.AddNodeAttribute)
			nodes.DELETE("/:composite_id/attributes/:attribute_name", h.RemoveNodeAttribute)
			nodes.PUT("/attributes/batch", h.BatchSetNodeAttributes)
		}

		domains := api.Group("/domains")
		{
			domains.GET("", h.ListDomains)
			domains.POST("", h.CreateDomain)
			domains.GET("/search", h.SearchDomains)
			domains.GET("/popular", h.GetPopularDomains)
			domains.GET("/:domain_name", h.GetDomain)
			domains.PUT("/:domain_name", h.UpdateDomain)
			domains.DELETE("/:domain_name", h.DeleteDomain)
			domains.GET("/:domain_name/stats", h.GetDomainStats)
		}

		server := api.Group("/server")
		{
			server.GET("/info", h.GetServerInfo)
			server.GET("/info/detailed", h.GetDetailedServerInfo)
			server.GET("/health", h.GetHealthStatus)
			server.GET("/stats", h.GetStatistics)
			server.GET("/docs", h.GetAPIDocumentation)
		}
	}
}

func (h *MCPHandler) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (h *MCPHandler) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[%s] %s %s %d %v",
			clientIP,
			method,
			path,
			statusCode,
			latency,
		)
	}
}
