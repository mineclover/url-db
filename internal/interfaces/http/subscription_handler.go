package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"url-db/internal/models"
	"url-db/internal/services"
)

// SubscriptionHandler handles HTTP requests for subscriptions
type SubscriptionHandler struct {
	subscriptionService *services.SubscriptionService
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(subscriptionService *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// CreateSubscription creates a new subscription
// @Summary Create subscription
// @Description Create a new subscription for node events
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param nodeId path int true "Node ID"
// @Param request body models.CreateNodeSubscriptionRequest true "Subscription details"
// @Success 201 {object} models.NodeSubscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/nodes/{nodeId}/subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	nodeID, err := strconv.ParseInt(c.Param("nodeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid node ID",
		})
		return
	}

	var req models.CreateNodeSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body",
		})
		return
	}

	subscription, err := h.subscriptionService.CreateSubscription(nodeID, &req)
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
			Message: "Failed to create subscription",
		})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetSubscription retrieves a subscription
// @Summary Get subscription
// @Description Get a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.NodeSubscription
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid subscription ID",
		})
		return
	}

	subscription, err := h.subscriptionService.GetSubscription(id)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Subscription not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get subscription",
		})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// UpdateSubscription updates a subscription
// @Summary Update subscription
// @Description Update a subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param request body models.UpdateNodeSubscriptionRequest true "Update details"
// @Success 200 {object} models.NodeSubscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid subscription ID",
		})
		return
	}

	var req models.UpdateNodeSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body",
		})
		return
	}

	subscription, err := h.subscriptionService.UpdateSubscription(id, &req)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Subscription not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update subscription",
		})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// DeleteSubscription deletes a subscription
// @Summary Delete subscription
// @Description Delete a subscription
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid subscription ID",
		})
		return
	}

	err = h.subscriptionService.DeleteSubscription(id)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Subscription not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete subscription",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetNodeSubscriptions retrieves all subscriptions for a node
// @Summary Get node subscriptions
// @Description Get all subscriptions for a specific node
// @Tags subscriptions
// @Produce json
// @Param nodeId path int true "Node ID"
// @Success 200 {array} models.NodeSubscription
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/nodes/{nodeId}/subscriptions [get]
func (h *SubscriptionHandler) GetNodeSubscriptions(c *gin.Context) {
	nodeID, err := strconv.ParseInt(c.Param("nodeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid node ID",
		})
		return
	}

	subscriptions, err := h.subscriptionService.GetNodeSubscriptions(nodeID)
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
			Message: "Failed to get subscriptions",
		})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// GetServiceSubscriptions retrieves all subscriptions for a service
// @Summary Get service subscriptions
// @Description Get all subscriptions for a specific service
// @Tags subscriptions
// @Produce json
// @Param service query string true "Service name"
// @Success 200 {array} models.NodeSubscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions [get]
func (h *SubscriptionHandler) GetServiceSubscriptions(c *gin.Context) {
	service := c.Query("service")
	if service == "" {
		// Get all subscriptions with pagination
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

		subscriptions, total, err := h.subscriptionService.GetAllSubscriptions(page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to get subscriptions",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"subscriptions": subscriptions,
			"total":         total,
			"page":          page,
			"page_size":     pageSize,
		})
		return
	}

	subscriptions, err := h.subscriptionService.GetServiceSubscriptions(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get subscriptions",
		})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}
