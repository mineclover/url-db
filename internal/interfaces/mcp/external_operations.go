package mcp

import (
	"context"
	"fmt"
	"strconv"

	"url-db/internal/models"
)

// External Dependency Management Operations

// CreateSubscription creates a new subscription for node events
func (s *mcpService) CreateSubscription(ctx context.Context, req *MCPCreateSubscriptionRequest) (*models.NodeSubscription, error) {
	// Parse composite ID to get node ID
	_, _, nodeIDStr, err := s.converter.ParseCompositeID(req.CompositeID)
	if err != nil {
		return nil, fmt.Errorf("invalid composite ID: %w", err)
	}

	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	// Create subscription request
	subscriptionReq := &models.CreateNodeSubscriptionRequest{
		SubscriberService:  req.SubscriberService,
		SubscriberEndpoint: req.SubscriberEndpoint,
		EventTypes:         req.EventTypes,
		FilterConditions:   req.FilterConditions,
	}

	return s.subscriptionService.CreateSubscription(nodeID, subscriptionReq)
}

// ListSubscriptions lists subscriptions with optional service filter
func (s *mcpService) ListSubscriptions(ctx context.Context, serviceName string, page, size int) (*MCPSubscriptionListResponse, error) {
	var subscriptions []*models.NodeSubscription
	var total int
	var err error

	if serviceName != "" {
		subscriptions, err = s.subscriptionService.GetServiceSubscriptions(serviceName)
		total = len(subscriptions)

		// Apply pagination manually
		start := (page - 1) * size
		end := start + size
		if start >= total {
			subscriptions = []*models.NodeSubscription{}
		} else {
			if end > total {
				end = total
			}
			subscriptions = subscriptions[start:end]
		}
	} else {
		subscriptions, total, err = s.subscriptionService.GetAllSubscriptions(page, size)
	}

	if err != nil {
		return nil, err
	}

	return &MCPSubscriptionListResponse{
		Subscriptions: subscriptions,
		Total:         total,
		Page:          page,
		Size:          size,
	}, nil
}

// GetNodeSubscriptions gets all subscriptions for a specific node
func (s *mcpService) GetNodeSubscriptions(ctx context.Context, compositeID string) ([]*models.NodeSubscription, error) {
	// Parse composite ID to get node ID
	_, _, nodeIDStr, err := s.converter.ParseCompositeID(compositeID)
	if err != nil {
		return nil, fmt.Errorf("invalid composite ID: %w", err)
	}

	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	return s.subscriptionService.GetNodeSubscriptions(nodeID)
}

// DeleteSubscription deletes a subscription
func (s *mcpService) DeleteSubscription(ctx context.Context, subscriptionID int64) error {
	return s.subscriptionService.DeleteSubscription(subscriptionID)
}

// CreateDependency creates a new dependency relationship
func (s *mcpService) CreateDependency(ctx context.Context, req *MCPCreateDependencyRequest) (*models.NodeDependency, error) {
	// Parse dependent node composite ID
	_, _, dependentNodeIDStr, err := s.converter.ParseCompositeID(req.DependentNodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid dependent node composite ID: %w", err)
	}

	dependentNodeID, err := strconv.ParseInt(dependentNodeIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid dependent node ID: %w", err)
	}

	// Parse dependency node composite ID
	_, _, dependencyNodeIDStr, err := s.converter.ParseCompositeID(req.DependencyNodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid dependency node composite ID: %w", err)
	}

	dependencyNodeID, err := strconv.ParseInt(dependencyNodeIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid dependency node ID: %w", err)
	}

	// Create dependency request
	dependencyReq := &models.CreateNodeDependencyRequest{
		DependencyNodeID: dependencyNodeID,
		DependencyType:   req.DependencyType,
		CascadeDelete:    req.CascadeDelete,
		CascadeUpdate:    req.CascadeUpdate,
		Metadata:         req.Metadata,
	}

	return s.dependencyService.CreateDependency(dependentNodeID, dependencyReq)
}

// ListNodeDependencies gets all dependencies for a node
func (s *mcpService) ListNodeDependencies(ctx context.Context, compositeID string) ([]*models.NodeDependency, error) {
	// Parse composite ID to get node ID
	_, _, nodeIDStr, err := s.converter.ParseCompositeID(compositeID)
	if err != nil {
		return nil, fmt.Errorf("invalid composite ID: %w", err)
	}

	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	return s.dependencyService.GetNodeDependencies(nodeID)
}

// ListNodeDependents gets all nodes that depend on this node
func (s *mcpService) ListNodeDependents(ctx context.Context, compositeID string) ([]*models.NodeDependency, error) {
	// Parse composite ID to get node ID
	_, _, nodeIDStr, err := s.converter.ParseCompositeID(compositeID)
	if err != nil {
		return nil, fmt.Errorf("invalid composite ID: %w", err)
	}

	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	return s.dependencyService.GetNodeDependents(nodeID)
}

// DeleteDependency deletes a dependency relationship
func (s *mcpService) DeleteDependency(ctx context.Context, dependencyID int64) error {
	return s.dependencyService.DeleteDependency(dependencyID)
}

// GetNodeEvents gets event history for a node
func (s *mcpService) GetNodeEvents(ctx context.Context, compositeID string, limit int) ([]*models.NodeEvent, error) {
	// Parse composite ID to get node ID
	_, _, nodeIDStr, err := s.converter.ParseCompositeID(compositeID)
	if err != nil {
		return nil, fmt.Errorf("invalid composite ID: %w", err)
	}

	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	return s.eventService.GetNodeEvents(nodeID, limit)
}

// GetPendingEvents gets unprocessed events
func (s *mcpService) GetPendingEvents(ctx context.Context, limit int) ([]*models.NodeEvent, error) {
	return s.eventService.GetPendingEvents(limit)
}

// ProcessEvent marks an event as processed
func (s *mcpService) ProcessEvent(ctx context.Context, eventID int64) error {
	return s.eventService.ProcessEvent(eventID)
}

// GetEventStats gets system event statistics
func (s *mcpService) GetEventStats(ctx context.Context) (map[string]interface{}, error) {
	return s.eventService.GetEventStats()
}
