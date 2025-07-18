package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"url-db/internal/services"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
	eventService *services.EventService
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

// GetNodeEvents retrieves events for a node
// @Summary Get node events
// @Description Get events for a specific node
// @Tags events
// @Produce json
// @Param nodeId path int true "Node ID"
// @Param limit query int false "Maximum number of events to return" default(50)
// @Success 200 {array} models.NodeEvent
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/nodes/{nodeId}/events [get]
func (h *EventHandler) GetNodeEvents(c *gin.Context) {
	nodeID, err := strconv.ParseInt(c.Param("nodeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid node ID",
		})
		return
	}
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	
	events, err := h.eventService.GetNodeEvents(nodeID, limit)
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
			Message: "Failed to get events",
		})
		return
	}
	
	c.JSON(http.StatusOK, events)
}

// GetPendingEvents retrieves unprocessed events
// @Summary Get pending events
// @Description Get unprocessed events for processing
// @Tags events
// @Produce json
// @Param limit query int false "Maximum number of events to return" default(100)
// @Success 200 {array} models.NodeEvent
// @Failure 500 {object} ErrorResponse
// @Router /api/events/pending [get]
func (h *EventHandler) GetPendingEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	
	events, err := h.eventService.GetPendingEvents(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get pending events",
		})
		return
	}
	
	c.JSON(http.StatusOK, events)
}

// ProcessEvent marks an event as processed
// @Summary Process event
// @Description Mark an event as processed
// @Tags events
// @Param eventId path int true "Event ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/events/{eventId}/process [post]
func (h *EventHandler) ProcessEvent(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid event ID",
		})
		return
	}
	
	err = h.eventService.ProcessEvent(eventID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Event not found",
			})
			return
		}
		if err.Error() == "event already processed" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Event already processed",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to process event",
		})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// GetEventsByType retrieves events by type and date range
// @Summary Get events by type
// @Description Get events by type within a date range
// @Tags events
// @Produce json
// @Param type query string true "Event type"
// @Param start query string true "Start date (RFC3339)"
// @Param end query string true "End date (RFC3339)"
// @Success 200 {array} models.NodeEvent
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/events [get]
func (h *EventHandler) GetEventsByType(c *gin.Context) {
	eventType := c.Query("type")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Event type is required",
		})
		return
	}
	
	startStr := c.Query("start")
	endStr := c.Query("end")
	
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid start date format",
		})
		return
	}
	
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid end date format",
		})
		return
	}
	
	events, err := h.eventService.GetEventsByTypeAndDateRange(eventType, start, end)
	if err != nil {
		if err.Error() == "end date must be after start date" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get events",
		})
		return
	}
	
	c.JSON(http.StatusOK, events)
}

// GetEventStats retrieves event statistics
// @Summary Get event statistics
// @Description Get statistics about events in the system
// @Tags events
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/events/stats [get]
func (h *EventHandler) GetEventStats(c *gin.Context) {
	stats, err := h.eventService.GetEventStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get event statistics",
		})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// CleanupEvents deletes old processed events
// @Summary Cleanup old events
// @Description Delete processed events older than specified days
// @Tags events
// @Param days query int true "Number of days to retain events" minimum(1)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/events/cleanup [post]
func (h *EventHandler) CleanupEvents(c *gin.Context) {
	days, err := strconv.Atoi(c.Query("days"))
	if err != nil || days < 1 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid days parameter",
		})
		return
	}
	
	retention := time.Duration(days) * 24 * time.Hour
	deleted, err := h.eventService.CleanupOldEvents(retention)
	if err != nil {
		if err.Error() == "minimum retention period is 24 hours" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to cleanup events",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"deleted_events": deleted,
		"message":       "Events cleaned up successfully",
	})
}