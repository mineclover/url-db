package node

import (
	"context"
	"fmt"

	"url-db/internal/domain/attribute"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
)

// SetNodeAttributesUseCase handles setting attributes for a node with validation
type SetNodeAttributesUseCase struct {
	nodeRepo          repository.NodeRepository
	attributeRepo     repository.AttributeRepository
	nodeAttributeRepo repository.NodeAttributeRepository
	validatorRegistry *attribute.ValidatorRegistry
}

// NewSetNodeAttributesUseCase creates a new use case for setting node attributes
func NewSetNodeAttributesUseCase(
	nodeRepo repository.NodeRepository,
	attributeRepo repository.AttributeRepository,
	nodeAttributeRepo repository.NodeAttributeRepository,
) *SetNodeAttributesUseCase {
	return &SetNodeAttributesUseCase{
		nodeRepo:          nodeRepo,
		attributeRepo:     attributeRepo,
		nodeAttributeRepo: nodeAttributeRepo,
		validatorRegistry: attribute.NewValidatorRegistry(),
	}
}

// AttributeInput represents an attribute to be set
type AttributeInput struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	OrderIndex *int   `json:"order_index,omitempty"`
}

// Execute sets attributes for a node with validation
func (uc *SetNodeAttributesUseCase) Execute(ctx context.Context, nodeID int, attributes []AttributeInput) error {
	// Verify node exists
	node, err := uc.nodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil {
		return fmt.Errorf("node not found: %d", nodeID)
	}

	// Get domain to get domain-specific attributes
	domain, err := uc.nodeRepo.GetDomainByNodeID(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get domain for node: %w", err)
	}
	if domain == nil {
		return fmt.Errorf("domain not found for node: %d", nodeID)
	}

	// Process and validate each attribute
	var nodeAttributes []*entity.NodeAttribute
	for _, attrInput := range attributes {
		// Get attribute definition from domain
		attr, err := uc.attributeRepo.GetByName(ctx, domain.ID(), attrInput.Name)
		if err != nil {
			return fmt.Errorf("failed to get attribute '%s': %w", attrInput.Name, err)
		}
		if attr == nil {
			return fmt.Errorf("attribute '%s' not defined in domain '%s'", attrInput.Name, domain.Name())
		}

		// Create validated node attribute
		nodeAttr, err := entity.ValidatedNodeAttribute(
			nodeID,
			attr.ID(),
			attribute.AttributeType(attr.Type()),
			attrInput.Value,
			attrInput.OrderIndex,
			uc.validatorRegistry,
		)
		if err != nil {
			return fmt.Errorf("validation failed for attribute '%s': %w", attrInput.Name, err)
		}

		nodeAttributes = append(nodeAttributes, nodeAttr)
	}

	// Set all attributes (this will replace existing ones)
	err = uc.nodeAttributeRepo.SetNodeAttributes(ctx, nodeID, nodeAttributes)
	if err != nil {
		return fmt.Errorf("failed to set node attributes: %w", err)
	}

	return nil
}