package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"url-db/internal/models"
)

// SubscriptionRepository handles database operations for subscriptions
type SubscriptionRepository struct {
	db *sqlx.DB
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(subscription *models.NodeSubscription) error {
	query := `
		INSERT INTO node_subscriptions (
			subscriber_service, subscriber_endpoint, subscribed_node_id,
			event_types, filter_conditions, is_active
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.Exec(
		query,
		subscription.SubscriberService,
		subscription.SubscriberEndpoint,
		subscription.SubscribedNodeID,
		subscription.EventTypes,
		subscription.FilterConditions,
		subscription.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	subscription.ID = id
	return nil
}

// GetByID retrieves a subscription by ID
func (r *SubscriptionRepository) GetByID(id int64) (*models.NodeSubscription, error) {
	var subscription models.NodeSubscription
	query := `
		SELECT id, subscriber_service, subscriber_endpoint, subscribed_node_id,
			   event_types, filter_conditions, is_active, created_at, updated_at
		FROM node_subscriptions
		WHERE id = ?
	`
	
	err := r.db.Get(&subscription, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	
	return &subscription, nil
}

// GetByNode retrieves all subscriptions for a specific node
func (r *SubscriptionRepository) GetByNode(nodeID int64) ([]*models.NodeSubscription, error) {
	var subscriptions []*models.NodeSubscription
	query := `
		SELECT id, subscriber_service, subscriber_endpoint, subscribed_node_id,
			   event_types, filter_conditions, is_active, created_at, updated_at
		FROM node_subscriptions
		WHERE subscribed_node_id = ? AND is_active = true
	`
	
	err := r.db.Select(&subscriptions, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	
	return subscriptions, nil
}

// GetByService retrieves all subscriptions for a specific service
func (r *SubscriptionRepository) GetByService(service string) ([]*models.NodeSubscription, error) {
	var subscriptions []*models.NodeSubscription
	query := `
		SELECT id, subscriber_service, subscriber_endpoint, subscribed_node_id,
			   event_types, filter_conditions, is_active, created_at, updated_at
		FROM node_subscriptions
		WHERE subscriber_service = ?
		ORDER BY created_at DESC
	`
	
	err := r.db.Select(&subscriptions, query, service)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	
	return subscriptions, nil
}

// Update updates a subscription
func (r *SubscriptionRepository) Update(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	
	query := "UPDATE node_subscriptions SET "
	args := []interface{}{}
	
	for field, value := range updates {
		query += field + " = ?, "
		args = append(args, value)
	}
	
	query = query[:len(query)-2] + " WHERE id = ?"
	args = append(args, id)
	
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	
	return nil
}

// Delete deletes a subscription
func (r *SubscriptionRepository) Delete(id int64) error {
	query := "DELETE FROM node_subscriptions WHERE id = ?"
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

// GetAll retrieves all subscriptions with pagination
func (r *SubscriptionRepository) GetAll(offset, limit int) ([]*models.NodeSubscription, int, error) {
	var subscriptions []*models.NodeSubscription
	var total int
	
	countQuery := "SELECT COUNT(*) FROM node_subscriptions"
	err := r.db.Get(&total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count subscriptions: %w", err)
	}
	
	query := `
		SELECT id, subscriber_service, subscriber_endpoint, subscribed_node_id,
			   event_types, filter_conditions, is_active, created_at, updated_at
		FROM node_subscriptions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	err = r.db.Select(&subscriptions, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	
	return subscriptions, total, nil
}