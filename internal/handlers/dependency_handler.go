package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
	"url-db/internal/services"
)

// DependencyHandler handles HTTP requests for dependencies
type DependencyHandler struct {
	dependencyService *services.DependencyService
}

// NewDependencyHandler creates a new dependency handler
func NewDependencyHandler(dependencyService *services.DependencyService) *DependencyHandler {
	return &DependencyHandler{
		dependencyService: dependencyService,
	}
}

// CreateDependency creates a new dependency
// @Summary Create dependency
// @Description Create a new dependency relationship between nodes
// @Tags dependencies
// @Accept json
// @Produce json
// @Param nodeId path int true "Dependent Node ID"
// @Param request body models.CreateNodeDependencyRequest true "Dependency details"
// @Success 201 {object} models.NodeDependency
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Circular dependency detected"
// @Failure 500 {object} ErrorResponse
// @Router /api/nodes/{nodeId}/dependencies [post]
func (h *DependencyHandler) CreateDependency(c *gin.Context) {
	nodeID, err := strconv.ParseInt(c.Param("nodeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid node ID",
		})
		return
	}
	
	var req models.CreateNodeDependencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body",
		})
		return
	}
	
	dependency, err := h.dependencyService.CreateDependency(nodeID, &req)
	if err != nil {
		if err.Error() == "dependent node not found" || err.Error() == "dependency node not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "circular dependency detected" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "conflict",
				Message: "Circular dependency detected",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create dependency",
		})
		return
	}
	
	c.JSON(http.StatusCreated, dependency)
}

// GetDependency retrieves a dependency
// @Summary Get dependency
// @Description Get a dependency by ID
// @Tags dependencies
// @Produce json
// @Param id path int true "Dependency ID"
// @Success 200 {object} models.NodeDependency
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/dependencies/{id} [get]
func (h *DependencyHandler) GetDependency(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid dependency ID",
		})
		return
	}
	
	dependency, err := h.dependencyService.GetDependency(id)
	if err != nil {
		if err.Error() == "dependency not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Dependency not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get dependency",
		})
		return
	}
	
	c.JSON(http.StatusOK, dependency)
}

// DeleteDependency deletes a dependency
// @Summary Delete dependency
// @Description Delete a dependency relationship
// @Tags dependencies
// @Param id path int true "Dependency ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/dependencies/{id} [delete]
func (h *DependencyHandler) DeleteDependency(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid dependency ID",
		})
		return
	}
	
	err = h.dependencyService.DeleteDependency(id)
	if err != nil {
		if err.Error() == "dependency not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Dependency not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete dependency",
		})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// GetNodeDependencies retrieves all dependencies for a node
// @Summary Get node dependencies
// @Description Get all dependencies where the node is dependent
// @Tags dependencies
// @Produce json
// @Param nodeId path int true "Node ID"
// @Success 200 {array} models.NodeDependency
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/nodes/{nodeId}/dependencies [get]
func (h *DependencyHandler) GetNodeDependencies(c *gin.Context) {
	nodeID, err := strconv.ParseInt(c.Param("nodeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid node ID",
		})
		return
	}
	
	dependencies, err := h.dependencyService.GetNodeDependencies(nodeID)
	if err != nil {
		if err.Error() == "node not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Node not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get dependencies",
		})
		return
	}
	
	c.JSON(http.StatusOK, dependencies)
}

// GetNodeDependents retrieves all nodes that depend on this node
// @Summary Get node dependents
// @Description Get all nodes that depend on this node
// @Tags dependencies
// @Produce json
// @Param nodeId path int true "Node ID"
// @Success 200 {array} models.NodeDependency
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/nodes/{nodeId}/dependents [get]
func (h *DependencyHandler) GetNodeDependents(c *gin.Context) {
	nodeID, err := strconv.ParseInt(c.Param("nodeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid node ID",
		})
		return
	}
	
	dependents, err := h.dependencyService.GetNodeDependents(nodeID)
	if err != nil {
		if err.Error() == "node not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Node not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get dependents",
		})
		return
	}
	
	c.JSON(http.StatusOK, dependents)
}