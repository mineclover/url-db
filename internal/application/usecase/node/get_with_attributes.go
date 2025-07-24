package node

import (
	"context"
	"url-db/internal/application/dto/response"
	"url-db/internal/domain/repository"
)

// GetNodeWithAttributesUseCase handles getting a node with its attributes
type GetNodeWithAttributesUseCase struct {
	nodeRepo          repository.NodeRepository
	nodeAttributeRepo repository.NodeAttributeRepository
	attributeRepo     repository.AttributeRepository
}

// NewGetNodeWithAttributesUseCase creates a new instance of GetNodeWithAttributesUseCase
func NewGetNodeWithAttributesUseCase(nodeRepo repository.NodeRepository, nodeAttributeRepo repository.NodeAttributeRepository, attributeRepo repository.AttributeRepository) *GetNodeWithAttributesUseCase {
	return &GetNodeWithAttributesUseCase{
		nodeRepo:          nodeRepo,
		nodeAttributeRepo: nodeAttributeRepo,
		attributeRepo:     attributeRepo,
	}
}

// NodeWithAttributesResponse represents a node with its attributes
type NodeWithAttributesResponse struct {
	Node       response.NodeResponse                 `json:"node"`
	Attributes []response.NodeAttributeResponse      `json:"attributes"`
}

// Execute performs the get node with attributes use case
func (uc *GetNodeWithAttributesUseCase) Execute(ctx context.Context, nodeID int) (*NodeWithAttributesResponse, error) {
	// Get the node
	node, err := uc.nodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	// Get domain information
	domain, err := uc.nodeRepo.GetDomainByNodeID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	// Get node attributes
	nodeAttributes, err := uc.nodeAttributeRepo.GetByNodeID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	// Convert node to response
	nodeResponse := response.NodeResponse{
		ID:          node.ID(),
		URL:         node.URL(),
		DomainName:  domain.Name(),
		Title:       node.Title(),
		Description: node.Description(),
		CreatedAt:   node.CreatedAt(),
		UpdatedAt:   node.UpdatedAt(),
	}

	// Convert attributes to response
	var attributeResponses []response.NodeAttributeResponse
	for _, nodeAttr := range nodeAttributes {
		// Get attribute definition to show name and type
		attr, err := uc.attributeRepo.GetByID(ctx, nodeAttr.AttributeID())
		if err != nil {
			continue // Skip if attribute definition not found
		}

		attrResponse := response.NodeAttributeResponse{
			AttributeName: attr.Name(),
			AttributeType: attr.Type(),
			Value:         nodeAttr.Value(),
		}

		if nodeAttr.OrderIndex() != nil {
			attrResponse.OrderIndex = nodeAttr.OrderIndex()
		}

		attributeResponses = append(attributeResponses, attrResponse)
	}

	return &NodeWithAttributesResponse{
		Node:       nodeResponse,
		Attributes: attributeResponses,
	}, nil
}