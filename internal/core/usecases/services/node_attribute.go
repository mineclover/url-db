package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"url-db/internal/models"
)

type NodeAttributeRepository interface {
	Create(ctx context.Context, nodeAttribute *models.NodeAttribute) error
	GetByID(ctx context.Context, id int) (*models.NodeAttribute, error)
	GetByNodeAndAttribute(ctx context.Context, nodeID, attributeID int) (*models.NodeAttribute, error)
	ListByNode(ctx context.Context, nodeID int) ([]*models.NodeAttribute, error)
	Update(ctx context.Context, nodeAttribute *models.NodeAttribute) error
	Delete(ctx context.Context, id int) error
	ExistsByNodeAndAttribute(ctx context.Context, nodeID, attributeID int) (bool, error)
}

type nodeAttributeService struct {
	nodeAttributeRepo NodeAttributeRepository
	nodeRepo          NodeRepository
	attributeRepo     AttributeRepository
	logger            *log.Logger
}

func NewNodeAttributeService(nodeAttributeRepo NodeAttributeRepository, nodeRepo NodeRepository, attributeRepo AttributeRepository, logger *log.Logger) NodeAttributeService {
	return &nodeAttributeService{
		nodeAttributeRepo: nodeAttributeRepo,
		nodeRepo:          nodeRepo,
		attributeRepo:     attributeRepo,
		logger:            logger,
	}
}

func (s *nodeAttributeService) CreateNodeAttribute(ctx context.Context, nodeID int, req *models.CreateNodeAttributeRequest) (*models.NodeAttribute, error) {
	if err := validatePositiveInteger(nodeID, "nodeID"); err != nil {
		return nil, err
	}

	if err := validatePositiveInteger(req.AttributeID, "attributeID"); err != nil {
		return nil, err
	}

	req.Value = normalizeString(req.Value)

	node, err := s.nodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNodeNotFoundError(nodeID)
		}
		s.logger.Printf("Failed to get node: %v", err)
		return nil, err
	}

	attribute, err := s.attributeRepo.GetByID(ctx, req.AttributeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAttributeNotFoundError(req.AttributeID)
		}
		s.logger.Printf("Failed to get attribute: %v", err)
		return nil, err
	}

	if node.DomainID != attribute.DomainID {
		return nil, NewBusinessLogicError("node and attribute must belong to the same domain")
	}

	if err := validateAttributeValue(string(attribute.Type), req.Value); err != nil {
		return nil, NewAttributeValueInvalidError(req.AttributeID, req.Value, err.Error())
	}

	exists, err := s.nodeAttributeRepo.ExistsByNodeAndAttribute(ctx, nodeID, req.AttributeID)
	if err != nil {
		s.logger.Printf("Failed to check node attribute existence: %v", err)
		return nil, err
	}
	if exists {
		return nil, NewBusinessLogicError("node attribute already exists")
	}

	nodeAttribute := &models.NodeAttribute{
		NodeID:      nodeID,
		AttributeID: req.AttributeID,
		Value:       req.Value,
		OrderIndex:  req.OrderIndex,
		CreatedAt:   time.Now(),
	}

	if err := s.nodeAttributeRepo.Create(ctx, nodeAttribute); err != nil {
		s.logger.Printf("Failed to create node attribute: %v", err)
		return nil, err
	}

	s.logger.Printf("Created node attribute: node %d, attribute %d (ID: %d)", nodeID, req.AttributeID, nodeAttribute.ID)
	return nodeAttribute, nil
}

func (s *nodeAttributeService) GetNodeAttribute(ctx context.Context, id int) (*models.NodeAttribute, error) {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return nil, err
	}

	nodeAttribute, err := s.nodeAttributeRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNodeAttributeNotFoundError(id)
		}
		s.logger.Printf("Failed to get node attribute: %v", err)
		return nil, err
	}

	return nodeAttribute, nil
}

func (s *nodeAttributeService) ListNodeAttributesByNode(ctx context.Context, nodeID int) (*models.NodeAttributeListResponse, error) {
	if err := validatePositiveInteger(nodeID, "nodeID"); err != nil {
		return nil, err
	}

	_, err := s.nodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNodeNotFoundError(nodeID)
		}
		s.logger.Printf("Failed to get node: %v", err)
		return nil, err
	}

	nodeAttributes, err := s.nodeAttributeRepo.ListByNode(ctx, nodeID)
	if err != nil {
		s.logger.Printf("Failed to list node attributes: %v", err)
		return nil, err
	}

	nodeAttributeList := make([]models.NodeAttributeWithInfo, len(nodeAttributes))
	for i, nodeAttribute := range nodeAttributes {
		// Convert NodeAttribute to NodeAttributeWithInfo
		// This is a placeholder - in reality you'd need to fetch attribute details
		nodeAttributeList[i] = models.NodeAttributeWithInfo{
			ID:          nodeAttribute.ID,
			NodeID:      nodeAttribute.NodeID,
			AttributeID: nodeAttribute.AttributeID,
			Name:        "unknown", // Would need to fetch from attribute table
			Type:        "string",  // Would need to fetch from attribute table
			Value:       nodeAttribute.Value,
			OrderIndex:  nodeAttribute.OrderIndex,
			CreatedAt:   nodeAttribute.CreatedAt,
		}
	}

	return &models.NodeAttributeListResponse{
		NodeAttributes: nodeAttributeList,
		Attributes:     nodeAttributeList, // For backward compatibility
	}, nil
}

func (s *nodeAttributeService) UpdateNodeAttribute(ctx context.Context, id int, req *models.UpdateNodeAttributeRequest) (*models.NodeAttribute, error) {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return nil, err
	}

	req.Value = normalizeString(req.Value)

	nodeAttribute, err := s.nodeAttributeRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNodeAttributeNotFoundError(id)
		}
		s.logger.Printf("Failed to get node attribute: %v", err)
		return nil, err
	}

	attribute, err := s.attributeRepo.GetByID(ctx, nodeAttribute.AttributeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAttributeNotFoundError(nodeAttribute.AttributeID)
		}
		s.logger.Printf("Failed to get attribute: %v", err)
		return nil, err
	}

	if err := validateAttributeValue(string(attribute.Type), req.Value); err != nil {
		return nil, NewAttributeValueInvalidError(nodeAttribute.AttributeID, req.Value, err.Error())
	}

	nodeAttribute.Value = req.Value
	nodeAttribute.OrderIndex = req.OrderIndex

	if err := s.nodeAttributeRepo.Update(ctx, nodeAttribute); err != nil {
		s.logger.Printf("Failed to update node attribute: %v", err)
		return nil, err
	}

	s.logger.Printf("Updated node attribute: node %d, attribute %d (ID: %d)", nodeAttribute.NodeID, nodeAttribute.AttributeID, nodeAttribute.ID)
	return nodeAttribute, nil
}

func (s *nodeAttributeService) DeleteNodeAttribute(ctx context.Context, id int) error {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return err
	}

	err := s.nodeAttributeRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewNodeAttributeNotFoundError(id)
		}
		s.logger.Printf("Failed to delete node attribute: %v", err)
		return err
	}

	s.logger.Printf("Deleted node attribute with ID: %d", id)
	return nil
}

func (s *nodeAttributeService) ValidateNodeAttributeValue(ctx context.Context, nodeID, attributeID int, value string) error {
	if err := validatePositiveInteger(nodeID, "nodeID"); err != nil {
		return err
	}

	if err := validatePositiveInteger(attributeID, "attributeID"); err != nil {
		return err
	}

	value = normalizeString(value)

	node, err := s.nodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewNodeNotFoundError(nodeID)
		}
		s.logger.Printf("Failed to get node: %v", err)
		return err
	}

	attribute, err := s.attributeRepo.GetByID(ctx, attributeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewAttributeNotFoundError(attributeID)
		}
		s.logger.Printf("Failed to get attribute: %v", err)
		return err
	}

	if node.DomainID != attribute.DomainID {
		return NewBusinessLogicError("node and attribute must belong to the same domain")
	}

	if err := validateAttributeValue(string(attribute.Type), value); err != nil {
		return NewAttributeValueInvalidError(attributeID, value, err.Error())
	}

	return nil
}
