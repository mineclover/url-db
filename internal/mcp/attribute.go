package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/url-db/internal/models"
)

type AttributeManager struct {
	attributeService AttributeService
	domainService    DomainService
	converter        *Converter
}

func NewAttributeManager(attributeService AttributeService, domainService DomainService, converter *Converter) *AttributeManager {
	return &AttributeManager{
		attributeService: attributeService,
		domainService:    domainService,
		converter:        converter,
	}
}

func (am *AttributeManager) GetNodeAttributes(ctx context.Context, compositeID string) (*models.MCPNodeAttributeResponse, error) {
	if err := am.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := am.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attributes, err := am.attributeService.GetNodeAttributes(ctx, nodeID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get node attributes: %v", err))
	}

	mcpAttributes := am.converter.NodeAttributesToMCPAttributes(attributes)

	return &models.MCPNodeAttributeResponse{
		CompositeID: compositeID,
		Attributes:  mcpAttributes,
	}, nil
}

func (am *AttributeManager) SetNodeAttributes(ctx context.Context, compositeID string, req *models.SetMCPNodeAttributesRequest) (*models.MCPNodeAttributeResponse, error) {
	if err := am.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := am.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domainName, err := am.converter.ExtractDomainNameFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domain, err := am.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	if err := am.validateAttributeRequests(req.Attributes); err != nil {
		return nil, NewValidationError(err.Error())
	}

	existingAttributes, err := am.attributeService.GetNodeAttributes(ctx, nodeID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get existing attributes: %v", err))
	}

	existingAttrMap := make(map[string]models.NodeAttributeWithInfo)
	for _, attr := range existingAttributes {
		existingAttrMap[attr.Name] = attr
	}

	requestedAttrMap := make(map[string]string)
	for _, attrReq := range req.Attributes {
		requestedAttrMap[attrReq.Name] = attrReq.Value
	}

	for _, existingAttr := range existingAttributes {
		if _, exists := requestedAttrMap[existingAttr.Name]; !exists {
			if err := am.attributeService.DeleteNodeAttribute(ctx, nodeID, existingAttr.AttributeID); err != nil {
				return nil, NewInternalServerError(fmt.Sprintf("failed to delete attribute: %v", err))
			}
		}
	}

	for _, attrReq := range req.Attributes {
		attribute, err := am.attributeService.GetAttributeByName(ctx, domain.ID, attrReq.Name)
		if err != nil {
			return nil, NewValidationError(fmt.Sprintf("attribute '%s' not found in domain '%s'", attrReq.Name, domainName))
		}

		if err := am.validateAttributeValue(attribute.Type, attrReq.Value); err != nil {
			return nil, NewValidationError(fmt.Sprintf("invalid value for attribute '%s': %v", attrReq.Name, err))
		}

		if existingAttr, exists := existingAttrMap[attrReq.Name]; exists {
			if err := am.attributeService.DeleteNodeAttribute(ctx, nodeID, existingAttr.AttributeID); err != nil {
				return nil, NewInternalServerError(fmt.Sprintf("failed to delete existing attribute: %v", err))
			}
		}

		createReq := &models.CreateNodeAttributeRequest{
			AttributeID: attribute.ID,
			Value:       attrReq.Value,
		}

		if _, err := am.attributeService.SetNodeAttribute(ctx, nodeID, createReq); err != nil {
			return nil, NewInternalServerError(fmt.Sprintf("failed to set attribute: %v", err))
		}
	}

	return am.GetNodeAttributes(ctx, compositeID)
}

func (am *AttributeManager) AddNodeAttribute(ctx context.Context, compositeID string, req *AddAttributeRequest) (*models.MCPNodeAttributeResponse, error) {
	if err := am.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := am.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domainName, err := am.converter.ExtractDomainNameFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domain, err := am.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	attribute, err := am.attributeService.GetAttributeByName(ctx, domain.ID, req.Name)
	if err != nil {
		return nil, NewValidationError(fmt.Sprintf("attribute '%s' not found in domain '%s'", req.Name, domainName))
	}

	if err := am.validateAttributeValue(attribute.Type, req.Value); err != nil {
		return nil, NewValidationError(fmt.Sprintf("invalid value for attribute '%s': %v", req.Name, err))
	}

	createReq := &models.CreateNodeAttributeRequest{
		AttributeID: attribute.ID,
		Value:       req.Value,
		OrderIndex:  req.OrderIndex,
	}

	if _, err := am.attributeService.SetNodeAttribute(ctx, nodeID, createReq); err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to add attribute: %v", err))
	}

	return am.GetNodeAttributes(ctx, compositeID)
}

func (am *AttributeManager) RemoveNodeAttribute(ctx context.Context, compositeID string, attributeName string) (*models.MCPNodeAttributeResponse, error) {
	if err := am.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := am.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domainName, err := am.converter.ExtractDomainNameFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	domain, err := am.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	attribute, err := am.attributeService.GetAttributeByName(ctx, domain.ID, attributeName)
	if err != nil {
		return nil, NewValidationError(fmt.Sprintf("attribute '%s' not found in domain '%s'", attributeName, domainName))
	}

	if err := am.attributeService.DeleteNodeAttribute(ctx, nodeID, attribute.ID); err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to remove attribute: %v", err))
	}

	return am.GetNodeAttributes(ctx, compositeID)
}

func (am *AttributeManager) GetNodeAttributesByType(ctx context.Context, compositeID string, attributeType models.AttributeType) ([]models.MCPAttribute, error) {
	if err := am.converter.ValidateCompositeID(compositeID); err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	nodeID, err := am.converter.ExtractNodeIDFromCompositeID(compositeID)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeID)
	}

	attributes, err := am.attributeService.GetNodeAttributes(ctx, nodeID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get node attributes: %v", err))
	}

	var filteredAttributes []models.MCPAttribute
	for _, attr := range attributes {
		if attr.Type == attributeType {
			mcpAttr := am.converter.NodeAttributeToMCPAttribute(&attr)
			if mcpAttr != nil {
				filteredAttributes = append(filteredAttributes, *mcpAttr)
			}
		}
	}

	return filteredAttributes, nil
}

func (am *AttributeManager) validateAttributeRequests(attributes []struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}) error {
	if len(attributes) == 0 {
		return fmt.Errorf("attributes cannot be empty")
	}

	nameSet := make(map[string]bool)
	for _, attr := range attributes {
		if attr.Name == "" {
			return fmt.Errorf("attribute name cannot be empty")
		}

		if attr.Value == "" {
			return fmt.Errorf("attribute value cannot be empty")
		}

		if nameSet[attr.Name] {
			return fmt.Errorf("duplicate attribute name: %s", attr.Name)
		}
		nameSet[attr.Name] = true
	}

	return nil
}

func (am *AttributeManager) validateAttributeValue(attributeType models.AttributeType, value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	switch attributeType {
	case models.AttributeTypeTag:
		return am.validateTagValue(value)
	case models.AttributeTypeOrderedTag:
		return am.validateOrderedTagValue(value)
	case models.AttributeTypeNumber:
		return am.validateNumberValue(value)
	case models.AttributeTypeString:
		return am.validateStringValue(value)
	case models.AttributeTypeMarkdown:
		return am.validateMarkdownValue(value)
	case models.AttributeTypeImage:
		return am.validateImageValue(value)
	default:
		return fmt.Errorf("unknown attribute type: %s", attributeType)
	}
}

func (am *AttributeManager) validateTagValue(value string) error {
	if len(value) > 50 {
		return fmt.Errorf("tag value cannot exceed 50 characters")
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("tag value cannot be empty after trimming")
	}

	return nil
}

func (am *AttributeManager) validateOrderedTagValue(value string) error {
	if len(value) > 50 {
		return fmt.Errorf("ordered tag value cannot exceed 50 characters")
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("ordered tag value cannot be empty after trimming")
	}

	return nil
}

func (am *AttributeManager) validateNumberValue(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("number value cannot be empty")
	}

	if _, err := fmt.Sscanf(value, "%f", new(float64)); err != nil {
		return fmt.Errorf("invalid number format: %s", value)
	}

	return nil
}

func (am *AttributeManager) validateStringValue(value string) error {
	if len(value) > 2048 {
		return fmt.Errorf("string value cannot exceed 2048 characters")
	}

	return nil
}

func (am *AttributeManager) validateMarkdownValue(value string) error {
	if len(value) > 10000 {
		return fmt.Errorf("markdown value cannot exceed 10000 characters")
	}

	return nil
}

func (am *AttributeManager) validateImageValue(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("image value cannot be empty")
	}

	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") && !strings.HasPrefix(value, "data:image/") {
		return fmt.Errorf("image value must be a valid URL or data URI")
	}

	if len(value) > 2048 {
		return fmt.Errorf("image value cannot exceed 2048 characters")
	}

	return nil
}

func (am *AttributeManager) BatchSetNodeAttributes(ctx context.Context, requests []BatchAttributeRequest) (*BatchAttributeResponse, error) {
	if len(requests) == 0 {
		return &BatchAttributeResponse{
			Success: []BatchAttributeSuccess{},
			Failed:  []BatchAttributeFailure{},
		}, nil
	}

	var successes []BatchAttributeSuccess
	var failures []BatchAttributeFailure

	for _, req := range requests {
		response, err := am.SetNodeAttributes(ctx, req.CompositeID, &req.AttributesRequest)
		if err != nil {
			failures = append(failures, BatchAttributeFailure{
				CompositeID: req.CompositeID,
				Error:       err.Error(),
			})
		} else {
			successes = append(successes, BatchAttributeSuccess{
				CompositeID: req.CompositeID,
				Attributes:  response.Attributes,
			})
		}
	}

	return &BatchAttributeResponse{
		Success: successes,
		Failed:  failures,
	}, nil
}

type AddAttributeRequest struct {
	Name       string `json:"name" binding:"required"`
	Value      string `json:"value" binding:"required"`
	OrderIndex *int   `json:"order_index"`
}

type BatchAttributeRequest struct {
	CompositeID       string                                `json:"composite_id"`
	AttributesRequest models.SetMCPNodeAttributesRequest `json:"attributes_request"`
}

type BatchAttributeResponse struct {
	Success []BatchAttributeSuccess `json:"success"`
	Failed  []BatchAttributeFailure `json:"failed"`
}

type BatchAttributeSuccess struct {
	CompositeID string                `json:"composite_id"`
	Attributes  []models.MCPAttribute `json:"attributes"`
}

type BatchAttributeFailure struct {
	CompositeID string `json:"composite_id"`
	Error       string `json:"error"`
}

type AttributeFilter struct {
	Name  string                `json:"name,omitempty"`
	Type  models.AttributeType `json:"type,omitempty"`
	Value string                `json:"value,omitempty"`
}

type AttributeStats struct {
	AttributeName string `json:"attribute_name"`
	TotalCount    int    `json:"total_count"`
	UniqueValues  int    `json:"unique_values"`
	MostCommon    string `json:"most_common"`
}