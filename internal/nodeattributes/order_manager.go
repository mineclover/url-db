package nodeattributes

import (
	"fmt"
	"internal/models"
)

type OrderManager interface {
	AssignOrderIndex(nodeID, attributeID int, orderIndex *int) (*int, error)
	ValidateOrderIndex(nodeID, attributeID int, orderIndex *int, excludeID *int) error
	ReorderAfterDeletion(nodeID, attributeID int, deletedOrderIndex int) error
}

type orderManager struct {
	repo Repository
}

func NewOrderManager(repo Repository) OrderManager {
	return &orderManager{repo: repo}
}

func (om *orderManager) AssignOrderIndex(nodeID, attributeID int, orderIndex *int) (*int, error) {
	if orderIndex == nil {
		// Auto-assign next available order index
		maxIndex, err := om.repo.GetMaxOrderIndex(nodeID, attributeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get max order index: %w", err)
		}
		nextIndex := maxIndex + 1
		return &nextIndex, nil
	}
	
	// Validate provided order index
	if *orderIndex < 1 {
		return nil, ErrInvalidOrderIndex
	}
	
	// Check if order index already exists
	existing, err := om.getNodeAttributeByOrderIndex(nodeID, attributeID, *orderIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing order index: %w", err)
	}
	
	if existing != nil {
		// Shift existing items up
		if err := om.repo.ReorderAfterIndex(nodeID, attributeID, *orderIndex-1); err != nil {
			return nil, fmt.Errorf("failed to reorder existing items: %w", err)
		}
	}
	
	return orderIndex, nil
}

func (om *orderManager) ValidateOrderIndex(nodeID, attributeID int, orderIndex *int, excludeID *int) error {
	if orderIndex == nil {
		return nil
	}
	
	if *orderIndex < 1 {
		return ErrInvalidOrderIndex
	}
	
	// Check if order index already exists (excluding current item if updating)
	existing, err := om.getNodeAttributeByOrderIndex(nodeID, attributeID, *orderIndex)
	if err != nil {
		return fmt.Errorf("failed to check existing order index: %w", err)
	}
	
	if existing != nil {
		// If we're updating and this is the same item, it's okay
		if excludeID != nil && existing.ID == *excludeID {
			return nil
		}
		return ErrDuplicateOrderIndex
	}
	
	return nil
}

func (om *orderManager) ReorderAfterDeletion(nodeID, attributeID int, deletedOrderIndex int) error {
	// Get all items with order index greater than deleted item
	nodeAttrs, err := om.repo.GetByNodeID(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get node attributes: %w", err)
	}
	
	// Update order indexes for items after the deleted one
	for _, attr := range nodeAttrs {
		if attr.AttributeID == attributeID && attr.OrderIndex != nil && *attr.OrderIndex > deletedOrderIndex {
			newIndex := *attr.OrderIndex - 1
			updateReq := &models.UpdateNodeAttributeRequest{
				Value:      attr.Value,
				OrderIndex: &newIndex,
			}
			
			if _, err := om.repo.Update(attr.ID, updateReq); err != nil {
				return fmt.Errorf("failed to update order index after deletion: %w", err)
			}
		}
	}
	
	return nil
}

func (om *orderManager) getNodeAttributeByOrderIndex(nodeID, attributeID int, orderIndex int) (*models.NodeAttributeWithInfo, error) {
	nodeAttrs, err := om.repo.GetByNodeID(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node attributes: %w", err)
	}
	
	for _, attr := range nodeAttrs {
		if attr.AttributeID == attributeID && attr.OrderIndex != nil && *attr.OrderIndex == orderIndex {
			return &attr, nil
		}
	}
	
	return nil, nil
}