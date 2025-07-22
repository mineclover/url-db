package nodes

import (
	"fmt"

	"url-db/internal/models"
	"url-db/internal/services"
)

// NodeServiceWithEvents wraps the basic node service and adds event tracking
type NodeServiceWithEvents struct {
	nodeService  NodeService
	eventService *services.EventService
}

// NewNodeServiceWithEvents creates a new node service with event tracking
func NewNodeServiceWithEvents(nodeService NodeService, eventService *services.EventService) NodeService {
	return &NodeServiceWithEvents{
		nodeService:  nodeService,
		eventService: eventService,
	}
}

func (s *NodeServiceWithEvents) CreateNode(domainID int, req *models.CreateNodeRequest) (*models.Node, error) {
	// Create the node using the wrapped service
	node, err := s.nodeService.CreateNode(domainID, req)
	if err != nil {
		return nil, err
	}

	// Create event for node creation
	changes := &models.EventChanges{
		After: map[string]interface{}{
			"url":         node.Content,
			"title":       node.Title,
			"description": node.Description,
			"domain_id":   node.DomainID,
		},
	}

	if eventErr := s.eventService.CreateNodeEvent(int64(node.ID), "node.created", changes); eventErr != nil {
		// Log the error but don't fail the operation
		fmt.Printf("Warning: Failed to create event for node creation: %v\n", eventErr)
	}

	return node, nil
}

func (s *NodeServiceWithEvents) GetNodeByID(id int) (*models.Node, error) {
	return s.nodeService.GetNodeByID(id)
}

func (s *NodeServiceWithEvents) GetNodesByDomainID(domainID, page, size int) (*models.NodeListResponse, error) {
	return s.nodeService.GetNodesByDomainID(domainID, page, size)
}

func (s *NodeServiceWithEvents) FindNodeByURL(domainID int, req *models.FindNodeByURLRequest) (*models.Node, error) {
	return s.nodeService.FindNodeByURL(domainID, req)
}

func (s *NodeServiceWithEvents) UpdateNode(id int, req *models.UpdateNodeRequest) (*models.Node, error) {
	// Get the existing node to capture before state
	existingNode, err := s.nodeService.GetNodeByID(id)
	if err != nil {
		return nil, err
	}

	// Update the node using the wrapped service
	updatedNode, err := s.nodeService.UpdateNode(id, req)
	if err != nil {
		return nil, err
	}

	// Create event for node update
	changes := &models.EventChanges{
		Before: map[string]interface{}{
			"title":       existingNode.Title,
			"description": existingNode.Description,
		},
		After: map[string]interface{}{
			"title":       updatedNode.Title,
			"description": updatedNode.Description,
		},
	}

	if eventErr := s.eventService.CreateNodeEvent(int64(updatedNode.ID), "node.updated", changes); eventErr != nil {
		// Log the error but don't fail the operation
		fmt.Printf("Warning: Failed to create event for node update: %v\n", eventErr)
	}

	return updatedNode, nil
}

func (s *NodeServiceWithEvents) DeleteNode(id int) error {
	// Get the existing node to capture before state
	existingNode, err := s.nodeService.GetNodeByID(id)
	if err != nil {
		return err
	}

	// Delete the node using the wrapped service
	if err := s.nodeService.DeleteNode(id); err != nil {
		return err
	}

	// Create event for node deletion
	changes := &models.EventChanges{
		Before: map[string]interface{}{
			"url":         existingNode.Content,
			"title":       existingNode.Title,
			"description": existingNode.Description,
			"domain_id":   existingNode.DomainID,
		},
	}

	if eventErr := s.eventService.CreateNodeEvent(int64(existingNode.ID), "node.deleted", changes); eventErr != nil {
		// Log the error but don't fail the operation
		fmt.Printf("Warning: Failed to create event for node deletion: %v\n", eventErr)
	}

	return nil
}

func (s *NodeServiceWithEvents) SearchNodes(domainID int, query string, page, size int) (*models.NodeListResponse, error) {
	return s.nodeService.SearchNodes(domainID, query, page, size)
}