package mcp

import (
	"url-db/internal/models"
)

// External Dependency Management MCP Models

// MCPCreateSubscriptionRequest represents a request to create a subscription via MCP
type MCPCreateSubscriptionRequest struct {
	CompositeID         string                    `json:"composite_id"`
	SubscriberService   string                    `json:"subscriber_service"`
	SubscriberEndpoint  *string                   `json:"subscriber_endpoint,omitempty"`
	EventTypes          []string                  `json:"event_types"`
	FilterConditions    *models.FilterCondition   `json:"filter_conditions,omitempty"`
}

// MCPSubscriptionListResponse represents a paginated list of subscriptions
type MCPSubscriptionListResponse struct {
	Subscriptions []*models.NodeSubscription `json:"subscriptions"`
	Total         int                        `json:"total"`
	Page          int                        `json:"page"`
	Size          int                        `json:"size"`
}

// MCPCreateDependencyRequest represents a request to create a dependency via MCP
type MCPCreateDependencyRequest struct {
	DependentNodeID  string                      `json:"dependent_node_id"`
	DependencyNodeID string                      `json:"dependency_node_id"`
	DependencyType   string                      `json:"dependency_type"`
	CascadeDelete    bool                        `json:"cascade_delete"`
	CascadeUpdate    bool                        `json:"cascade_update"`
	Metadata         *models.DependencyMetadata  `json:"metadata,omitempty"`
}