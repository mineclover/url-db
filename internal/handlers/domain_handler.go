package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
)

type DomainService interface {
	CreateDomain(req *models.CreateDomainRequest) (*models.Domain, error)
	GetDomains(page, size int) (*models.DomainListResponse, error)
	GetDomainByID(id int) (*models.Domain, error)
	UpdateDomain(id int, req *models.UpdateDomainRequest) (*models.Domain, error)
	DeleteDomain(id int) error
}

type DomainHandler struct {
	*BaseHandler
	domainService DomainService
}

func NewDomainHandler(domainService DomainService) *DomainHandler {
	return &DomainHandler{
		BaseHandler:   NewBaseHandler(),
		domainService: domainService,
	}
}

func (h *DomainHandler) CreateDomain(c *gin.Context) {
	var req models.CreateDomainRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	domain, err := h.domainService.CreateDomain(&req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, domain)
}

func (h *DomainHandler) GetDomains(c *gin.Context) {
	page := h.ParseIntQuery(c, "page", 1)
	size := h.ParseIntQuery(c, "size", 20)

	if size > 100 {
		size = 100
	}

	response, err := h.domainService.GetDomains(page, size)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DomainHandler) GetDomain(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	domain, err := h.domainService.GetDomainByID(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain)
}

func (h *DomainHandler) UpdateDomain(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	var req models.UpdateDomainRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	domain, err := h.domainService.UpdateDomain(id, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain)
}

func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	err = h.domainService.DeleteDomain(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DomainHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		domains := api.Group("/domains")
		{
			domains.POST("", h.CreateDomain)
			domains.GET("", h.GetDomains)
			domains.GET("/:id", h.GetDomain)
			domains.PUT("/:id", h.UpdateDomain)
			domains.DELETE("/:id", h.DeleteDomain)
		}
	}
}
