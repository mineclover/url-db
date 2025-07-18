package nodeattributes

import (
	"fmt"
	"time"

	"url-db/internal/models"
)

type Service interface {
	CreateNodeAttribute(nodeID int, req *models.CreateNodeAttributeRequest) (*models.NodeAttribute, error)
	GetNodeAttributeByID(id int) (*models.NodeAttributeWithInfo, error)
	GetNodeAttributesByNodeID(nodeID int) (*models.NodeAttributeListResponse, error)
	UpdateNodeAttribute(id int, req *models.UpdateNodeAttributeRequest) (*models.NodeAttribute, error)
	DeleteNodeAttribute(id int) error
	DeleteNodeAttributeByNodeIDAndAttributeID(nodeID, attributeID int) error
	ValidateNodeAttribute(nodeID int, req *models.CreateNodeAttributeRequest) error
}

type service struct {
	repo         Repository
	validator    Validator
	orderManager OrderManager
}

func NewService(repo Repository, validator Validator, orderManager OrderManager) Service {
	return &service{
		repo:         repo,
		validator:    validator,
		orderManager: orderManager,
	}
}

func (s *service) CreateNodeAttribute(nodeID int, req *models.CreateNodeAttributeRequest) (*models.NodeAttribute, error) {
	// Validate request
	if err := s.ValidateNodeAttribute(nodeID, req); err != nil {
		return nil, err
	}

	// Check if node and attribute belong to the same domain
	if err := s.repo.ValidateNodeAndAttributeDomain(nodeID, req.AttributeID); err != nil {
		return nil, ErrNodeAttributeDomainMismatch
	}

	// Get attribute type for validation
	attrInfo, err := s.repo.GetByNodeIDAndAttributeID(nodeID, req.AttributeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute info: %w", err)
	}

	// For non-ordered attributes, check if attribute already exists for this node
	if attrInfo != nil {
		// If it's not an ordered tag, we don't allow duplicates
		if attrInfo.Type != models.AttributeTypeOrderedTag {
			return nil, ErrNodeAttributeExists
		}
	}

	// Get attribute type from database (we need to query attributes table)
	attributeType, err := s.getAttributeType(req.AttributeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute type: %w", err)
	}

	// Validate attribute value and order index
	if err := s.validator.Validate(attributeType, req.Value, req.OrderIndex); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Handle order index for ordered tags
	var finalOrderIndex *int
	if attributeType == models.AttributeTypeOrderedTag {
		finalOrderIndex, err = s.orderManager.AssignOrderIndex(nodeID, req.AttributeID, req.OrderIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to assign order index: %w", err)
		}
	}

	// Create node attribute
	nodeAttribute := &models.NodeAttribute{
		NodeID:      nodeID,
		AttributeID: req.AttributeID,
		Value:       req.Value,
		OrderIndex:  finalOrderIndex,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(nodeAttribute); err != nil {
		return nil, fmt.Errorf("failed to create node attribute: %w", err)
	}

	return nodeAttribute, nil
}

func (s *service) GetNodeAttributeByID(id int) (*models.NodeAttributeWithInfo, error) {
	nodeAttr, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get node attribute: %w", err)
	}

	if nodeAttr == nil {
		return nil, ErrNodeAttributeNotFound
	}

	return nodeAttr, nil
}

func (s *service) GetNodeAttributesByNodeID(nodeID int) (*models.NodeAttributeListResponse, error) {
	nodeAttrs, err := s.repo.GetByNodeID(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node attributes: %w", err)
	}

	return &models.NodeAttributeListResponse{
		Attributes: nodeAttrs,
	}, nil
}

func (s *service) UpdateNodeAttribute(id int, req *models.UpdateNodeAttributeRequest) (*models.NodeAttribute, error) {
	// Get existing node attribute
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing node attribute: %w", err)
	}

	if existing == nil {
		return nil, ErrNodeAttributeNotFound
	}

	// Validate attribute value and order index
	if err := s.validator.Validate(existing.Type, req.Value, req.OrderIndex); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Handle order index for ordered tags
	var finalOrderIndex *int = req.OrderIndex
	if existing.Type == models.AttributeTypeOrderedTag {
		if req.OrderIndex != nil {
			// Validate order index (excluding current item)
			if err := s.orderManager.ValidateOrderIndex(existing.NodeID, existing.AttributeID, req.OrderIndex, &id); err != nil {
				return nil, fmt.Errorf("order index validation failed: %w", err)
			}

			// If order index changed, need to reorder
			if existing.OrderIndex == nil || *existing.OrderIndex != *req.OrderIndex {
				// Move other items to make space
				if err := s.repo.ReorderAfterIndex(existing.NodeID, existing.AttributeID, *req.OrderIndex-1); err != nil {
					return nil, fmt.Errorf("failed to reorder items: %w", err)
				}
			}
		}
	}

	// Update node attribute
	updateReq := &models.UpdateNodeAttributeRequest{
		Value:      req.Value,
		OrderIndex: finalOrderIndex,
	}

	updatedAttr, err := s.repo.Update(id, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update node attribute: %w", err)
	}

	return updatedAttr, nil
}

func (s *service) DeleteNodeAttribute(id int) error {
	// Get existing node attribute for order management
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get existing node attribute: %w", err)
	}

	if existing == nil {
		return ErrNodeAttributeNotFound
	}

	// Delete the node attribute
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete node attribute: %w", err)
	}

	// Handle order reordering for ordered tags
	if existing.Type == models.AttributeTypeOrderedTag && existing.OrderIndex != nil {
		if err := s.orderManager.ReorderAfterDeletion(existing.NodeID, existing.AttributeID, *existing.OrderIndex); err != nil {
			return fmt.Errorf("failed to reorder after deletion: %w", err)
		}
	}

	return nil
}

func (s *service) DeleteNodeAttributeByNodeIDAndAttributeID(nodeID, attributeID int) error {
	// Get existing node attribute for order management
	existing, err := s.repo.GetByNodeIDAndAttributeID(nodeID, attributeID)
	if err != nil {
		return fmt.Errorf("failed to get existing node attribute: %w", err)
	}

	if existing == nil {
		return ErrNodeAttributeNotFound
	}

	// Delete the node attribute
	if err := s.repo.DeleteByNodeIDAndAttributeID(nodeID, attributeID); err != nil {
		return fmt.Errorf("failed to delete node attribute: %w", err)
	}

	// Handle order reordering for ordered tags
	if existing.Type == models.AttributeTypeOrderedTag && existing.OrderIndex != nil {
		if err := s.orderManager.ReorderAfterDeletion(existing.NodeID, existing.AttributeID, *existing.OrderIndex); err != nil {
			return fmt.Errorf("failed to reorder after deletion: %w", err)
		}
	}

	return nil
}

func (s *service) ValidateNodeAttribute(nodeID int, req *models.CreateNodeAttributeRequest) error {
	if req.AttributeID <= 0 {
		return fmt.Errorf("attribute_id must be a positive integer")
	}

	if req.Value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	if len(req.Value) > 2048 {
		return fmt.Errorf("value must be 2048 characters or less")
	}

	return nil
}

func (s *service) getAttributeType(attributeID int) (models.AttributeType, error) {
	return s.repo.GetAttributeType(attributeID)
}
