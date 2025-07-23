package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"url-db/internal/models"
)

// EventRepository handles database operations for events
type EventRepository struct {
	db *sqlx.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create creates a new event
func (r *EventRepository) Create(event *models.NodeEvent) error {
	query := `
		INSERT INTO node_events (
			node_id, event_type, event_data
		) VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		event.NodeID,
		event.EventType,
		event.EventData,
	)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	event.ID = id
	return nil
}

// GetByID retrieves an event by ID
func (r *EventRepository) GetByID(id int64) (*models.NodeEvent, error) {
	var event models.NodeEvent
	query := `
		SELECT id, node_id, event_type, event_data, occurred_at, processed_at
		FROM node_events
		WHERE id = ?
	`

	err := r.db.Get(&event, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

// GetByNode retrieves all events for a specific node
func (r *EventRepository) GetByNode(nodeID int64, limit int) ([]*models.NodeEvent, error) {
	var events []*models.NodeEvent
	query := `
		SELECT id, node_id, event_type, event_data, occurred_at, processed_at
		FROM node_events
		WHERE node_id = ?
		ORDER BY occurred_at DESC
		LIMIT ?
	`

	err := r.db.Select(&events, query, nodeID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}

// GetPendingEvents retrieves unprocessed events
func (r *EventRepository) GetPendingEvents(limit int) ([]*models.NodeEvent, error) {
	var events []*models.NodeEvent
	query := `
		SELECT id, node_id, event_type, event_data, occurred_at, processed_at
		FROM node_events
		WHERE processed_at IS NULL
		ORDER BY occurred_at ASC
		LIMIT ?
	`

	err := r.db.Select(&events, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending events: %w", err)
	}

	return events, nil
}

// MarkAsProcessed marks an event as processed
func (r *EventRepository) MarkAsProcessed(id int64) error {
	query := "UPDATE node_events SET processed_at = ? WHERE id = ?"

	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark event as processed: %w", err)
	}

	return nil
}

// GetByTypeAndDateRange retrieves events by type within a date range
func (r *EventRepository) GetByTypeAndDateRange(eventType string, start, end time.Time) ([]*models.NodeEvent, error) {
	var events []*models.NodeEvent
	query := `
		SELECT id, node_id, event_type, event_data, occurred_at, processed_at
		FROM node_events
		WHERE event_type = ? AND occurred_at BETWEEN ? AND ?
		ORDER BY occurred_at DESC
	`

	err := r.db.Select(&events, query, eventType, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by type and date: %w", err)
	}

	return events, nil
}

// DeleteOldEvents deletes events older than the specified duration
func (r *EventRepository) DeleteOldEvents(olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)
	query := "DELETE FROM node_events WHERE occurred_at < ? AND processed_at IS NOT NULL"

	result, err := r.db.Exec(query, cutoffTime)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old events: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetEventStats retrieves statistics about events
func (r *EventRepository) GetEventStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total events
	var totalEvents int
	err := r.db.Get(&totalEvents, "SELECT COUNT(*) FROM node_events")
	if err != nil {
		return nil, fmt.Errorf("failed to get total events: %w", err)
	}
	stats["total_events"] = totalEvents

	// Pending events
	var pendingEvents int
	err = r.db.Get(&pendingEvents, "SELECT COUNT(*) FROM node_events WHERE processed_at IS NULL")
	if err != nil {
		return nil, fmt.Errorf("failed to get pending events: %w", err)
	}
	stats["pending_events"] = pendingEvents

	// Events by type
	type eventTypeCount struct {
		EventType string `db:"event_type"`
		Count     int    `db:"count"`
	}
	var eventTypeCounts []eventTypeCount

	query := `
		SELECT event_type, COUNT(*) as count
		FROM node_events
		GROUP BY event_type
		ORDER BY count DESC
	`
	err = r.db.Select(&eventTypeCounts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by type: %w", err)
	}

	eventsByType := make(map[string]int)
	for _, etc := range eventTypeCounts {
		eventsByType[etc.EventType] = etc.Count
	}
	stats["events_by_type"] = eventsByType

	return stats, nil
}
