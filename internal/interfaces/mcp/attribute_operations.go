package mcp

import (
	"context"
	"fmt"

	"url-db/internal/models"
)

// Attribute operation methods for mcpService

func (s *mcpService) GetNodeAttributes(ctx context.Context, compositeID string) (*models.MCPNodeAttributeResponse, error) {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attributes, err := s.attributeService.GetNodeAttributes(ctx, nodeID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get node attributes: %v", err))
	}

	mcpAttributes := s.converter.NodeAttributesToMCPAttributes(attributes)

	return &models.MCPNodeAttributeResponse{
		CompositeID: compositeID,
		Attributes:  mcpAttributes,
	}, nil
}

func (s *mcpService) SetNodeAttributes(ctx context.Context, compositeID string, req *models.SetMCPNodeAttributesRequest) (*models.MCPNodeAttributeResponse, error) {
	if err := s.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := s.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domainName, err := s.converter.ExtractDomainNameFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	existingAttributes, err := s.attributeService.GetNodeAttributes(ctx, nodeID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get existing attributes: %v", err))
	}

	existingAttrMap := make(map[string]int)
	for _, attr := range existingAttributes {
		existingAttrMap[attr.Name] = attr.AttributeID
	}

	for _, attrReq := range req.Attributes {
		attribute, err := s.attributeService.GetAttributeByName(ctx, domain.ID, attrReq.Name)
		if err != nil {
			return nil, NewValidationError(fmt.Sprintf("attribute '%s' not found", attrReq.Name))
		}

		if existingAttrID, exists := existingAttrMap[attrReq.Name]; exists {
			if err := s.attributeService.DeleteNodeAttribute(ctx, nodeID, existingAttrID); err != nil {
				return nil, NewInternalServerError(fmt.Sprintf("failed to delete existing attribute: %v", err))
			}
		}

		createReq := &models.CreateNodeAttributeRequest{
			AttributeID: attribute.ID,
			Value:       attrReq.Value,
			OrderIndex:  attrReq.OrderIndex,
		}

		if _, err := s.attributeService.SetNodeAttribute(ctx, nodeID, createReq); err != nil {
			return nil, NewInternalServerError(fmt.Sprintf("failed to set attribute: %v", err))
		}
	}

	return s.GetNodeAttributes(ctx, compositeID)
}

// Domain attribute management methods
func (s *mcpService) ListDomainAttributes(ctx context.Context, domainName string) (*models.MCPDomainAttributeListResponse, error) {
	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	response, err := s.attributeService.ListAttributes(ctx, domain.ID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list attributes: %v", err))
	}

	mcpAttributes := make([]models.MCPDomainAttribute, len(response.Attributes))
	for i, attr := range response.Attributes {
		compositeID := s.converter.CreateAttributeCompositeID(domain.Name, attr.ID)
		mcpAttributes[i] = models.MCPDomainAttribute{
			CompositeID: compositeID,
			Name:        attr.Name,
			Type:        attr.Type,
			Description: attr.Description,
			CreatedAt:   attr.CreatedAt,
			UpdatedAt:   attr.CreatedAt, // Assuming no updated_at field in Attribute model
		}
	}

	return &models.MCPDomainAttributeListResponse{
		DomainName: domainName,
		Attributes: mcpAttributes,
		TotalCount: len(mcpAttributes),
	}, nil
}

func (s *mcpService) CreateDomainAttribute(ctx context.Context, domainName string, req *models.CreateAttributeRequest) (*models.MCPDomainAttribute, error) {
	domain, err := s.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	attribute, err := s.attributeService.CreateAttribute(ctx, domain.ID, req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to create attribute: %v", err))
	}

	compositeID := s.converter.CreateAttributeCompositeID(domain.Name, attribute.ID)
	return &models.MCPDomainAttribute{
		CompositeID: compositeID,
		Name:        attribute.Name,
		Type:        attribute.Type,
		Description: attribute.Description,
		CreatedAt:   attribute.CreatedAt,
		UpdatedAt:   attribute.CreatedAt,
	}, nil
}

func (s *mcpService) GetDomainAttribute(ctx context.Context, compositeID string) (*models.MCPDomainAttribute, error) {
	if err := s.converter.ValidateAttributeCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attributeID, err := s.converter.ExtractAttributeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attribute, err := s.attributeService.GetAttribute(ctx, attributeID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get attribute: %v", err))
	}

	return &models.MCPDomainAttribute{
		CompositeID: compositeID,
		Name:        attribute.Name,
		Type:        attribute.Type,
		Description: attribute.Description,
		CreatedAt:   attribute.CreatedAt,
		UpdatedAt:   attribute.CreatedAt,
	}, nil
}

func (s *mcpService) UpdateDomainAttribute(ctx context.Context, compositeID string, req *models.UpdateAttributeRequest) (*models.MCPDomainAttribute, error) {
	if err := s.converter.ValidateAttributeCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attributeID, err := s.converter.ExtractAttributeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attribute, err := s.attributeService.UpdateAttribute(ctx, attributeID, req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to update attribute: %v", err))
	}

	return &models.MCPDomainAttribute{
		CompositeID: compositeID,
		Name:        attribute.Name,
		Type:        attribute.Type,
		Description: attribute.Description,
		CreatedAt:   attribute.CreatedAt,
		UpdatedAt:   attribute.CreatedAt,
	}, nil
}

func (s *mcpService) DeleteDomainAttribute(ctx context.Context, compositeID string) error {
	if err := s.converter.ValidateAttributeCompositeID(compositeID); err != nil {
		return NewInvalidCompositeKeyError(compositeID)
	}

	attributeID, err := s.converter.ExtractAttributeIDFromCompositeID(compositeID)
	if err != nil {
		return NewInvalidCompositeKeyError(compositeID)
	}

	if err := s.attributeService.DeleteAttribute(ctx, attributeID); err != nil {
		return NewInternalServerError(fmt.Sprintf("failed to delete attribute: %v", err))
	}

	return nil
}
