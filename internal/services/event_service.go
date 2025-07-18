package services

import (
	"fmt"
	"time"

	"url-db/internal/models"
	"url-db/internal/repositories"
)

// EventService handles business logic for events
type EventService struct {
	eventRepo *repositories.EventRepository
	nodeRepo  *repositories.NodeRepository
}

// NewEventService creates a new event service
func NewEventService(
	eventRepo *repositories.EventRepository,
	nodeRepo *repositories.NodeRepository,
) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		nodeRepo:  nodeRepo,
	}
}

// GetNodeEvents retrieves events for a specific node
func (s *EventService) GetNodeEvents(nodeID int64, limit int) ([]*models.NodeEvent, error) {
	// Verify node exists
	node, err := s.nodeRepo.GetByID(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil {
		return nil, fmt.Errorf("node not found")
	}
	
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	
	return s.eventRepo.GetByNode(nodeID, limit)
}

// GetPendingEvents retrieves unprocessed events
func (s *EventService) GetPendingEvents(limit int) ([]*models.NodeEvent, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	
	return s.eventRepo.GetPendingEvents(limit)
}

// ProcessEvent marks an event as processed
func (s *EventService) ProcessEvent(eventID int64) error {
	// Verify event exists
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}
	if event == nil {
		return fmt.Errorf("event not found")
	}
	
	if event.ProcessedAt != nil {
		return fmt.Errorf("event already processed")
	}
	
	return s.eventRepo.MarkAsProcessed(eventID)
}

// GetEventsByTypeAndDateRange retrieves events by type within a date range
func (s *EventService) GetEventsByTypeAndDateRange(eventType string, start, end time.Time) ([]*models.NodeEvent, error) {
	// Validate date range
	if end.Before(start) {
		return nil, fmt.Errorf("end date must be after start date")
	}
	
	return s.eventRepo.GetByTypeAndDateRange(eventType, start, end)
}

// CleanupOldEvents deletes processed events older than the specified duration
func (s *EventService) CleanupOldEvents(olderThan time.Duration) (int64, error) {
	if olderThan < 24*time.Hour {
		return 0, fmt.Errorf("minimum retention period is 24 hours")
	}
	
	return s.eventRepo.DeleteOldEvents(olderThan)
}

// GetEventStats retrieves statistics about events
func (s *EventService) GetEventStats() (map[string]interface{}, error) {
	return s.eventRepo.GetEventStats()
}

// CreateNodeEvent creates a new event for a node
func (s *EventService) CreateNodeEvent(nodeID int64, eventType string, changes *models.EventChanges) error {
	// Verify node exists
	node, err := s.nodeRepo.GetByID(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil {
		return fmt.Errorf("node not found")
	}
	
	// Create event data
	eventData := &models.EventData{
		NodeID:    nodeID,
		EventType: eventType,
		Timestamp: time.Now(),
		Changes:   changes,
		Metadata: map[string]interface{}{
			"domain_id": node.DomainID,
		},
	}
	
	// Create event
	event := &models.NodeEvent{
		NodeID:    nodeID,
		EventType: eventType,
		EventData: eventData,
	}
	
	return s.eventRepo.Create(event)
}