package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
)

type AttributeService interface {
	CreateAttribute(domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error)
	GetAttributesByDomainID(domainID int) ([]models.Attribute, error)
	GetAttributeByID(id int) (*models.Attribute, error)
	UpdateAttribute(id int, req *models.UpdateAttributeRequest) (*models.Attribute, error)
	DeleteAttribute(id int) error
}

type AttributeHandler struct {
	*BaseHandler
	attributeService AttributeService
}

func NewAttributeHandler(attributeService AttributeService) *AttributeHandler {
	return &AttributeHandler{
		BaseHandler:      NewBaseHandler(),
		attributeService: attributeService,
	}
}

func (h *AttributeHandler) CreateAttribute(c *gin.Context) {
	domainID, err := h.ParseIntParam(c, "domain_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	var req models.CreateAttributeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	attribute, err := h.attributeService.CreateAttribute(domainID, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, attribute)
}

func (h *AttributeHandler) GetAttributesByDomain(c *gin.Context) {
	domainID, err := h.ParseIntParam(c, "domain_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	attributes, err := h.attributeService.GetAttributesByDomainID(domainID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attributes": attributes,
	})
}

func (h *AttributeHandler) GetAttribute(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	attribute, err := h.attributeService.GetAttributeByID(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, attribute)
}

func (h *AttributeHandler) UpdateAttribute(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	var req models.UpdateAttributeRequest
	if err := h.BindJSON(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	attribute, err := h.attributeService.UpdateAttribute(id, &req)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, attribute)
}

func (h *AttributeHandler) DeleteAttribute(c *gin.Context) {
	id, err := h.ParseIntParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	err = h.attributeService.DeleteAttribute(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AttributeHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Individual attribute operations
		attributes := api.Group("/attributes")
		{
			attributes.GET("/:id", h.GetAttribute)
			attributes.PUT("/:id", h.UpdateAttribute)
			attributes.DELETE("/:id", h.DeleteAttribute)
		}

		// Domain-specific attribute operations
		domains := api.Group("/domains")
		{
			domainAttributes := domains.Group("/:domain_id/attributes")
			{
				domainAttributes.POST("", h.CreateAttribute)
				domainAttributes.GET("", h.GetAttributesByDomain)
			}
		}
	}
}
