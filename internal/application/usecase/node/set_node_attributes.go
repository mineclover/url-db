package node

import (
	"context"
	"fmt"

	"url-db/internal/domain/attribute"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
	"url-db/internal/domain/service"
)

// SetNodeAttributesUseCase handles setting attributes for a node with validation
type SetNodeAttributesUseCase struct {
	nodeRepo          repository.NodeRepository
	attributeRepo     repository.AttributeRepository
	nodeAttributeRepo repository.NodeAttributeRepository
	templateService   service.TemplateService
	validatorRegistry *attribute.ValidatorRegistry
}

// NewSetNodeAttributesUseCase creates a new use case for setting node attributes
func NewSetNodeAttributesUseCase(
	nodeRepo repository.NodeRepository,
	attributeRepo repository.AttributeRepository,
	nodeAttributeRepo repository.NodeAttributeRepository,
	templateService service.TemplateService,
) *SetNodeAttributesUseCase {
	return &SetNodeAttributesUseCase{
		nodeRepo:          nodeRepo,
		attributeRepo:     attributeRepo,
		nodeAttributeRepo: nodeAttributeRepo,
		templateService:   templateService,
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

		// Validate attribute value against templates (진입점 제약)
		templateValidation, err := uc.templateService.ValidateAttributeValue(ctx, domain.Name(), attrInput.Name, attrInput.Value)
		if err != nil {
			return fmt.Errorf("template validation error for attribute '%s': %w", attrInput.Name, err)
		}

		// Reject if template validation fails
		if !templateValidation.IsValid {
			return &TemplateValidationError{
				AttributeName: attrInput.Name,
				Value:         attrInput.Value,
				ErrorCode:     templateValidation.ErrorCode,
				ErrorMessage:  templateValidation.ErrorMessage,
				AllowedValues: templateValidation.AllowedValues,
				TemplateUsed:  templateValidation.TemplateUsed,
			}
		}

		// Create validated node attribute (기존 검증 유지)
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

// TemplateValidationError represents a template-based validation error
type TemplateValidationError struct {
	AttributeName string   `json:"attribute_name"`
	Value         string   `json:"value"`
	ErrorCode     string   `json:"error_code"`
	ErrorMessage  string   `json:"error_message"`
	AllowedValues []string `json:"allowed_values,omitempty"`
	TemplateUsed  string   `json:"template_used,omitempty"`
}

func (e *TemplateValidationError) Error() string {
	return fmt.Sprintf("Template validation failed for attribute '%s' with value '%s': %s", 
		e.AttributeName, e.Value, e.ErrorMessage)
}
