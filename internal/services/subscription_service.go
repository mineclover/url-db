package services

import (
	"fmt"

	"url-db/internal/models"
	"url-db/internal/repositories"
)

// SubscriptionService handles business logic for subscriptions
type SubscriptionService struct {
	subscriptionRepo *repositories.SubscriptionRepository
	nodeRepo         repositories.NodeRepository
	eventRepo        *repositories.EventRepository
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(
	subscriptionRepo *repositories.SubscriptionRepository,
	nodeRepo repositories.NodeRepository,
	eventRepo *repositories.EventRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		nodeRepo:         nodeRepo,
		eventRepo:        eventRepo,
	}
}

// CreateSubscription creates a new subscription
func (s *SubscriptionService) CreateSubscription(nodeID int64, req *models.CreateNodeSubscriptionRequest) (*models.NodeSubscription, error) {
	// Verify node exists
	node, err := s.nodeRepo.GetByID(int(nodeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil {
		return nil, fmt.Errorf("node not found")
	}

	// Create subscription
	subscription := &models.NodeSubscription{
		SubscriberService:  req.SubscriberService,
		SubscriberEndpoint: req.SubscriberEndpoint,
		SubscribedNodeID:   nodeID,
		EventTypes:         models.EventTypeList(req.EventTypes),
		FilterConditions:   req.FilterConditions,
		IsActive:           true,
	}

	err = s.subscriptionRepo.Create(subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subscription, nil
}

// GetSubscription retrieves a subscription by ID
func (s *SubscriptionService) GetSubscription(id int64) (*models.NodeSubscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	if subscription == nil {
		return nil, fmt.Errorf("subscription not found")
	}

	return subscription, nil
}

// UpdateSubscription updates a subscription
func (s *SubscriptionService) UpdateSubscription(id int64, req *models.UpdateNodeSubscriptionRequest) (*models.NodeSubscription, error) {
	// Verify subscription exists
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	if subscription == nil {
		return nil, fmt.Errorf("subscription not found")
	}

	// Build updates
	updates := make(map[string]interface{})

	if req.SubscriberEndpoint != nil {
		updates["subscriber_endpoint"] = req.SubscriberEndpoint
	}

	if len(req.EventTypes) > 0 {
		updates["event_types"] = models.EventTypeList(req.EventTypes)
	}

	if req.FilterConditions != nil {
		updates["filter_conditions"] = req.FilterConditions
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Update subscription
	err = s.subscriptionRepo.Update(id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	// Get updated subscription
	return s.subscriptionRepo.GetByID(id)
}

// DeleteSubscription deletes a subscription
func (s *SubscriptionService) DeleteSubscription(id int64) error {
	err := s.subscriptionRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// GetNodeSubscriptions retrieves all subscriptions for a node
func (s *SubscriptionService) GetNodeSubscriptions(nodeID int64) ([]*models.NodeSubscription, error) {
	// Verify node exists
	node, err := s.nodeRepo.GetByID(int(nodeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil {
		return nil, fmt.Errorf("node not found")
	}

	return s.subscriptionRepo.GetByNode(nodeID)
}

// GetServiceSubscriptions retrieves all subscriptions for a service
func (s *SubscriptionService) GetServiceSubscriptions(service string) ([]*models.NodeSubscription, error) {
	return s.subscriptionRepo.GetByService(service)
}

// GetAllSubscriptions retrieves all subscriptions with pagination
func (s *SubscriptionService) GetAllSubscriptions(page, pageSize int) ([]*models.NodeSubscription, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.subscriptionRepo.GetAll(offset, pageSize)
}

// TriggerNodeEvent creates an event and notifies relevant subscribers
func (s *SubscriptionService) TriggerNodeEvent(nodeID int64, eventType string, eventData *models.EventData) error {
	// Create event
	event := &models.NodeEvent{
		NodeID:    nodeID,
		EventType: eventType,
		EventData: eventData,
	}

	err := s.eventRepo.Create(event)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	// Get active subscriptions for this node
	subscriptions, err := s.subscriptionRepo.GetByNode(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	// Filter subscriptions by event type
	for _, sub := range subscriptions {
		shouldNotify := false
		for _, et := range sub.EventTypes {
			if et == eventType {
				shouldNotify = true
				break
			}
		}

		if shouldNotify {
			// TODO: Implement webhook notification
			// This would typically be done asynchronously in a production system
			_ = sub
		}
	}

	return nil
}
