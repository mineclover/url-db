package mcp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"url-db/internal/models"
)

type Converter struct {
	compositeKeyService CompositeKeyService
	toolName            string
}

type CompositeKeyService interface {
	Create(domainName string, nodeID int) string
	Parse(compositeID string) (*CompositeKey, error)
	Validate(compositeID string) error
}

type CompositeKey struct {
	ToolName   string
	DomainName string
	ID         int
}

func NewConverter(compositeKeyService CompositeKeyService, toolName string) *Converter {
	return &Converter{
		compositeKeyService: compositeKeyService,
		toolName:            toolName,
	}
}

// ParseCompositeID parses a composite ID and returns tool name, domain name, and node ID
func (c *Converter) ParseCompositeID(compositeID string) (string, string, string, error) {
	compositeKey, err := c.compositeKeyService.Parse(compositeID)
	if err != nil {
		return "", "", "", err
	}
	
	return compositeKey.ToolName, compositeKey.DomainName, strconv.Itoa(compositeKey.ID), nil
}

func (c *Converter) NodeToMCPNode(node *models.Node, domain *models.Domain) (*models.MCPNode, error) {
	if node == nil || domain == nil {
		return nil, fmt.Errorf("node and domain cannot be nil")
	}

	compositeID := c.compositeKeyService.Create(domain.Name, node.ID)

	return &models.MCPNode{
		CompositeID: compositeID,
		URL:         node.Content,
		DomainName:  domain.Name,
		Title:       node.Title,
		Description: node.Description,
		CreatedAt:   node.CreatedAt,
		UpdatedAt:   node.UpdatedAt,
	}, nil
}

func (c *Converter) MCPNodeToNode(mcpNode *models.MCPNode) (*models.Node, int, error) {
	if mcpNode == nil {
		return nil, 0, fmt.Errorf("mcpNode cannot be nil")
	}

	compositeKey, err := c.compositeKeyService.Parse(mcpNode.CompositeID)
	if err != nil {
		return nil, 0, err
	}

	node := &models.Node{
		ID:          compositeKey.ID,
		Content:     mcpNode.URL,
		Title:       mcpNode.Title,
		Description: mcpNode.Description,
		CreatedAt:   mcpNode.CreatedAt,
		UpdatedAt:   mcpNode.UpdatedAt,
	}

	return node, compositeKey.ID, nil
}

func (c *Converter) CreateMCPNodeRequestToCreateNodeRequest(req *models.CreateMCPNodeRequest) *models.CreateNodeRequest {
	if req == nil {
		return nil
	}

	return &models.CreateNodeRequest{
		URL:         req.URL,
		Title:       req.Title,
		Description: req.Description,
	}
}

func (c *Converter) UpdateMCPNodeRequestToUpdateNodeRequest(req *models.UpdateNodeRequest) *models.UpdateNodeRequest {
	if req == nil {
		return nil
	}

	return &models.UpdateNodeRequest{
		Title:       req.Title,
		Description: req.Description,
	}
}

func (c *Converter) DomainToMCPDomain(domain *models.Domain, nodeCount int) *MCPDomain {
	if domain == nil {
		return nil
	}

	return &MCPDomain{
		Name:        domain.Name,
		Description: domain.Description,
		NodeCount:   nodeCount,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}

func (c *Converter) NodeAttributeToMCPAttribute(attr *models.NodeAttributeWithInfo) *models.MCPAttribute {
	if attr == nil {
		return nil
	}

	return &models.MCPAttribute{
		Name:  attr.Name,
		Type:  string(attr.Type),
		Value: attr.Value,
	}
}

func (c *Converter) NodeAttributesToMCPAttributes(attrs []models.NodeAttributeWithInfo) []models.MCPAttribute {
	if attrs == nil {
		return nil
	}

	mcpAttrs := make([]models.MCPAttribute, 0, len(attrs))
	for _, attr := range attrs {
		if mcpAttr := c.NodeAttributeToMCPAttribute(&attr); mcpAttr != nil {
			mcpAttrs = append(mcpAttrs, *mcpAttr)
		}
	}

	return mcpAttrs
}

func (c *Converter) ExtractDomainNameFromCompositeID(compositeID string) (string, error) {
	compositeKey, err := c.compositeKeyService.Parse(compositeID)
	if err != nil {
		return "", err
	}
	return compositeKey.DomainName, nil
}

func (c *Converter) ExtractNodeIDFromCompositeID(compositeID string) (int, error) {
	compositeKey, err := c.compositeKeyService.Parse(compositeID)
	if err != nil {
		return 0, err
	}
	return compositeKey.ID, nil
}

func (c *Converter) ValidateCompositeID(compositeID string) error {
	return c.compositeKeyService.Validate(compositeID)
}

// Attribute composite key methods
func (c *Converter) CreateAttributeCompositeID(domainName string, attributeID int) string {
	// Use a different format for attributes: tool-name:domain:attr-{id}
	return fmt.Sprintf("%s:%s:attr-%d", c.toolName, domainName, attributeID)
}

func (c *Converter) ExtractAttributeIDFromCompositeID(compositeID string) (int, error) {
	// Don't use the standard validation because it expects numeric IDs
	parts := strings.Split(compositeID, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid composite ID format")
	}

	// Validate tool name
	if parts[0] != c.toolName {
		return 0, fmt.Errorf("invalid tool name in composite ID")
	}

	// Check if this is an attribute ID (has "attr-" prefix)
	if !strings.HasPrefix(parts[2], "attr-") {
		return 0, fmt.Errorf("not an attribute composite ID")
	}

	// Extract the numeric ID after "attr-"
	idStr := strings.TrimPrefix(parts[2], "attr-")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid attribute ID: %v", err)
	}

	if id <= 0 {
		return 0, fmt.Errorf("attribute ID must be positive")
	}

	return id, nil
}

func (c *Converter) ValidateAttributeCompositeID(compositeID string) error {
	parts := strings.Split(compositeID, ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid composite ID format")
	}

	if parts[0] != c.toolName {
		return fmt.Errorf("invalid tool name in composite ID")
	}

	if parts[1] == "" {
		return fmt.Errorf("domain name is empty")
	}

	if !strings.HasPrefix(parts[2], "attr-") {
		return fmt.Errorf("not an attribute composite ID")
	}

	return nil
}

type MCPDomain struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	NodeCount   int       `json:"node_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MCPDomainListResponse struct {
	Domains []MCPDomain `json:"domains"`
}

type MCPServerInfo struct {
	Name               string   `json:"name"`
	Version            string   `json:"version"`
	Description        string   `json:"description"`
	Capabilities       []string `json:"capabilities"`
	CompositeKeyFormat string   `json:"composite_key_format"`
}
