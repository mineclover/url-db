package service

import (
	"context"
	"fmt"

	"url-db/internal/application/dto/response"
	"url-db/internal/constants"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
)

// ContentScanner provides token-based pagination for scanning domain content
type ContentScanner struct {
	nodeRepo      repository.NodeRepository
	attributeRepo repository.NodeAttributeRepository
	domainRepo    repository.DomainRepository
}

// NewContentScanner creates a new ContentScanner instance
func NewContentScanner(
	nodeRepo repository.NodeRepository,
	attributeRepo repository.NodeAttributeRepository,
	domainRepo repository.DomainRepository,
) *ContentScanner {
	return &ContentScanner{
		nodeRepo:      nodeRepo,
		attributeRepo: attributeRepo,
		domainRepo:    domainRepo,
	}
}

// ScanRequest represents the parameters for content scanning
type ScanRequest struct {
	DomainName         string `json:"domain_name"`
	MaxTokensPerPage   int    `json:"max_tokens_per_page"`
	Page               int    `json:"page"`               // Page number (1-based)
	IncludeAttributes  bool   `json:"include_attributes"`
	CompressAttributes bool   `json:"compress_attributes"` // Remove duplicate attribute values
}

// ScanResponse represents the response from content scanning
type ScanResponse struct {
	Items      []response.NodeWithAttributes `json:"items"`
	Pagination PaginationInfo                `json:"pagination"`
	Metadata   ScanMetadata                  `json:"metadata"`
}

// PaginationInfo contains pagination details
type PaginationInfo struct {
	CurrentPage   int  `json:"current_page"`
	TotalPages    int  `json:"total_pages"`
	CurrentTokens int  `json:"current_tokens"`
	HasMore       bool `json:"has_more"`
	HasPrevious   bool `json:"has_previous"`
}

// ScanMetadata contains scanning metadata
type ScanMetadata struct {
	TotalNodes         int                    `json:"total_nodes"`
	ProcessedNodes     int                    `json:"processed_nodes"`
	EstimatedTokens    int                    `json:"estimated_tokens"`
	EstimatedPages     int                    `json:"estimated_pages"`
	AttributeSummary   *AttributeSummary      `json:"attribute_summary,omitempty"`
	CompressedOutput   bool                   `json:"compressed_output"`
}

// AttributeSummary contains compressed attribute information
type AttributeSummary struct {
	UniqueValues       map[string][]string `json:"unique_values"`        // attribute_name -> unique values
	ValueCounts        map[string]int      `json:"value_counts"`         // "attr_name:value" -> count
	MostCommonValues   map[string]string   `json:"most_common_values"`   // attribute_name -> most common value
	TotalDuplicatesRemoved int             `json:"total_duplicates_removed"`
}

// PageInfo represents page calculation information
type PageInfo struct {
	CurrentPage      int `json:"current_page"`
	NodesPerPage     int `json:"nodes_per_page"`
	TotalNodes       int `json:"total_nodes"`
	TotalPages       int `json:"total_pages"`
	StartIndex       int `json:"start_index"`
	EndIndex         int `json:"end_index"`
}

// ScanAllContent performs page-based scanning of domain content with token optimization
func (cs *ContentScanner) ScanAllContent(ctx context.Context, req ScanRequest) (*ScanResponse, error) {
	// Validate domain exists
	domain, err := cs.domainRepo.GetByName(ctx, req.DomainName)
	if err != nil {
		return nil, fmt.Errorf("domain not found: %w", err)
	}

	// Set default values
	if req.MaxTokensPerPage <= 0 {
		req.MaxTokensPerPage = constants.DefaultMaxTokensPerPage
	}
	if req.MaxTokensPerPage > constants.MaxTokensPerPage {
		req.MaxTokensPerPage = constants.MaxTokensPerPage
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	// Get total node count
	totalNodes, err := cs.nodeRepo.CountByDomain(ctx, domain.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to count nodes: %w", err)
	}

	if totalNodes == 0 {
		return &ScanResponse{
			Items: []response.NodeWithAttributes{},
			Pagination: PaginationInfo{
				CurrentPage:   1,
				TotalPages:    1,
				CurrentTokens: 0,
				HasMore:       false,
				HasPrevious:   false,
			},
			Metadata: ScanMetadata{
				TotalNodes:       0,
				ProcessedNodes:   0,
				EstimatedTokens:  0,
				EstimatedPages:   1,
				CompressedOutput: req.CompressAttributes,
			},
		}, nil
	}

	// Calculate page information based on estimated nodes per page
	avgTokensPerNode := constants.AvgTokensPerNode
	if req.IncludeAttributes {
		avgTokensPerNode = int(float64(avgTokensPerNode) * 1.5)
	}
	estimatedNodesPerPage := req.MaxTokensPerPage / avgTokensPerNode
	if estimatedNodesPerPage < 1 {
		estimatedNodesPerPage = 1
	}

	pageInfo := cs.calculatePageInfo(req.Page, estimatedNodesPerPage, totalNodes)

	// Fetch nodes for the current page
	nodes, err := cs.fetchNodesForPage(ctx, domain.ID(), pageInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nodes: %w", err)
	}

	// Build response with token optimization
	result, actualTokens, attributesSummary, err := cs.buildOptimizedResponse(ctx, nodes, req)
	if err != nil {
		return nil, fmt.Errorf("failed to build response: %w", err)
	}

	// Calculate total estimated tokens
	estimatedTotalTokens := totalNodes * avgTokensPerNode
	estimatedPages := (estimatedTotalTokens / req.MaxTokensPerPage) + 1

	response := &ScanResponse{
		Items: result,
		Pagination: PaginationInfo{
			CurrentPage:   req.Page,
			TotalPages:    pageInfo.TotalPages,
			CurrentTokens: actualTokens,
			HasMore:       req.Page < pageInfo.TotalPages,
			HasPrevious:   req.Page > 1,
		},
		Metadata: ScanMetadata{
			TotalNodes:         totalNodes,
			ProcessedNodes:     len(result),
			EstimatedTokens:    estimatedTotalTokens,
			EstimatedPages:     estimatedPages,
			AttributeSummary:   attributesSummary,
			CompressedOutput:   req.CompressAttributes,
		},
	}

	return response, nil
}

// calculatePageInfo calculates page boundaries and metadata
func (cs *ContentScanner) calculatePageInfo(currentPage, nodesPerPage, totalNodes int) PageInfo {
	totalPages := (totalNodes + nodesPerPage - 1) / nodesPerPage // Ceiling division
	if totalPages < 1 {
		totalPages = 1
	}

	startIndex := (currentPage - 1) * nodesPerPage
	endIndex := startIndex + nodesPerPage
	if endIndex > totalNodes {
		endIndex = totalNodes
	}

	return PageInfo{
		CurrentPage:  currentPage,
		NodesPerPage: nodesPerPage,
		TotalNodes:   totalNodes,
		TotalPages:   totalPages,
		StartIndex:   startIndex,
		EndIndex:     endIndex,
	}
}

// fetchNodesForPage fetches nodes for a specific page
func (cs *ContentScanner) fetchNodesForPage(ctx context.Context, domainID int, pageInfo PageInfo) ([]*entity.Node, error) {
	// For now, we'll use the existing cursor-based method and slice the results
	// This is not optimal for large datasets but works for the current implementation
	allNodes, err := cs.nodeRepo.GetByDomainFromCursor(ctx, domainID, 0, pageInfo.TotalNodes)
	if err != nil {
		return nil, err
	}

	// Slice the nodes for the current page
	if pageInfo.StartIndex >= len(allNodes) {
		return []*entity.Node{}, nil
	}

	endIdx := pageInfo.EndIndex
	if endIdx > len(allNodes) {
		endIdx = len(allNodes)
	}

	return allNodes[pageInfo.StartIndex:endIdx], nil
}

// buildOptimizedResponse builds the response with token optimization and attribute compression
func (cs *ContentScanner) buildOptimizedResponse(ctx context.Context, nodes []*entity.Node, req ScanRequest) ([]response.NodeWithAttributes, int, *AttributeSummary, error) {
	result := make([]response.NodeWithAttributes, 0, len(nodes))
	totalTokens := 0
	var attributeSummary *AttributeSummary

	if req.CompressAttributes && req.IncludeAttributes {
		attributeSummary = &AttributeSummary{
			UniqueValues:       make(map[string][]string),
			ValueCounts:        make(map[string]int),
			MostCommonValues:   make(map[string]string),
			TotalDuplicatesRemoved: 0,
		}
	}

	// First pass: collect all attributes for compression analysis
	allAttributes := make(map[int][]*entity.NodeAttribute)
	if req.IncludeAttributes {
		for _, node := range nodes {
			attributes, err := cs.attributeRepo.GetByNodeID(ctx, node.ID())
			if err != nil {
				return nil, 0, nil, fmt.Errorf("failed to get attributes for node %d: %w", node.ID(), err)
			}
			allAttributes[node.ID()] = attributes
		}

		// Analyze attributes for compression if requested
		if req.CompressAttributes {
			cs.analyzeAttributesForCompression(allAttributes, attributeSummary)
		}
	}

	// Second pass: build response with optimized attributes
	for _, node := range nodes {
		nodeResp := response.NodeWithAttributes{
			ID:        node.ID(),
			Content:   node.Content(),
			CreatedAt: node.CreatedAt(),
			UpdatedAt: node.UpdatedAt(),
		}

		// Handle optional title and description
		if title := node.Title(); title != "" {
			nodeResp.Title = &title
		}
		if desc := node.Description(); desc != "" {
			nodeResp.Description = &desc
		}

		// Add attributes with compression
		if req.IncludeAttributes {
			attributes := allAttributes[node.ID()]
			if req.CompressAttributes {
				nodeResp.Attributes = cs.compressAttributes(attributes, attributeSummary)
			} else {
				nodeResp.Attributes = make([]response.AttributeValue, len(attributes))
				for j, attr := range attributes {
					nodeResp.Attributes[j] = response.AttributeValue{
						Name:          attr.Name(),
						Value:         attr.Value(),
						AttributeType: attr.AttributeType(),
						OrderIndex:    attr.OrderIndex(),
					}
				}
			}
		}

		// Estimate tokens for this node
		nodeTokens := cs.estimateNodeTokens(nodeResp, req.IncludeAttributes)
		totalTokens += nodeTokens

		result = append(result, nodeResp)
	}

	return result, totalTokens, attributeSummary, nil
}

// analyzeAttributesForCompression analyzes all attributes to build compression metadata
func (cs *ContentScanner) analyzeAttributesForCompression(allAttributes map[int][]*entity.NodeAttribute, summary *AttributeSummary) {
	attributeValueCounts := make(map[string]map[string]int) // attr_name -> value -> count

	// Count all attribute values
	for _, attributes := range allAttributes {
		for _, attr := range attributes {
			attrName := attr.Name()
			attrValue := attr.Value()

			if attributeValueCounts[attrName] == nil {
				attributeValueCounts[attrName] = make(map[string]int)
			}
			attributeValueCounts[attrName][attrValue]++
		}
	}

	// Build summary
	for attrName, valueCounts := range attributeValueCounts {
		// Get unique values
		uniqueValues := make([]string, 0, len(valueCounts))
		maxCount := 0
		mostCommonValue := ""

		for value, count := range valueCounts {
			uniqueValues = append(uniqueValues, value)
			summary.ValueCounts[attrName+":"+value] = count

			if count > maxCount {
				maxCount = count
				mostCommonValue = value
			}

			// Count duplicates (anything beyond the first occurrence)
			if count > 1 {
				summary.TotalDuplicatesRemoved += count - 1
			}
		}

		summary.UniqueValues[attrName] = uniqueValues
		summary.MostCommonValues[attrName] = mostCommonValue
	}
}

// compressAttributes applies compression to attributes by removing duplicates
func (cs *ContentScanner) compressAttributes(attributes []*entity.NodeAttribute, summary *AttributeSummary) []response.AttributeValue {
	if len(attributes) == 0 {
		return nil
	}

	// Group by attribute name and only show unique values
	seen := make(map[string]bool)
	compressed := make([]response.AttributeValue, 0)

	for _, attr := range attributes {
		key := attr.Name() + ":" + attr.Value()
		if seen[key] {
			continue // Skip duplicate
		}
		seen[key] = true

		compressed = append(compressed, response.AttributeValue{
			Name:          attr.Name(),
			Value:         attr.Value(),
			AttributeType: attr.AttributeType(),
			OrderIndex:    attr.OrderIndex(),
		})
	}

	return compressed
}

// estimateNodeTokens estimates tokens for a node with attributes
func (cs *ContentScanner) estimateNodeTokens(node response.NodeWithAttributes, includeAttributes bool) int {
	tokens := 0
	
	// Base content tokens (URL)
	tokens += len(node.Content) / 4
	
	// Title tokens
	if node.Title != nil {
		tokens += len(*node.Title) / 4
	}
	
	// Description tokens
	if node.Description != nil {
		tokens += len(*node.Description) / 4
	}
	
	// Attribute tokens (if included)
	if includeAttributes && node.Attributes != nil {
		for _, attr := range node.Attributes {
			tokens += len(attr.Name) / 4
			tokens += len(attr.Value) / 4
			if attr.AttributeType != nil {
				tokens += len(*attr.AttributeType) / 4
			}
		}
	}
	
	// JSON structure overhead (~20% additional tokens)
	tokens = int(float64(tokens) * 1.2)
	
	// Minimum tokens per node
	if tokens < constants.MinTokensPerNode {
		tokens = constants.MinTokensPerNode
	}
	
	return tokens
}

// SmartChunker handles token-based chunking of content
type SmartChunker struct {
	TargetTokens      int                             `json:"target_tokens"`
	BufferTokens      int                             `json:"buffer_tokens"`
	CurrentChunk      []response.NodeWithAttributes   `json:"current_chunk"`
	CurrentTokens     int                             `json:"current_tokens"`
	IncludeAttributes bool                            `json:"include_attributes"`
}

// NewSmartChunker creates a new SmartChunker instance
func NewSmartChunker(targetTokens int, includeAttributes bool) *SmartChunker {
	bufferTokens := targetTokens / 8 // 12.5% buffer
	if bufferTokens < 500 {
		bufferTokens = 500
	}

	return &SmartChunker{
		TargetTokens:      targetTokens - bufferTokens,
		BufferTokens:      bufferTokens,
		CurrentChunk:      make([]response.NodeWithAttributes, 0),
		CurrentTokens:     0,
		IncludeAttributes: includeAttributes,
	}
}

// CanAddNode checks if a node can be added without exceeding token limit
func (sc *SmartChunker) CanAddNode(node response.NodeWithAttributes) bool {
	nodeTokens := sc.EstimateNodeTokens(node)
	return sc.CurrentTokens+nodeTokens <= sc.TargetTokens
}

// AddNode adds a node to the current chunk
func (sc *SmartChunker) AddNode(node response.NodeWithAttributes) {
	nodeTokens := sc.EstimateNodeTokens(node)
	sc.CurrentChunk = append(sc.CurrentChunk, node)
	sc.CurrentTokens += nodeTokens
}

// EstimateNodeTokens estimates the token count for a node
func (sc *SmartChunker) EstimateNodeTokens(node response.NodeWithAttributes) int {
	// Base estimation: ~4 characters per token (conservative for multilingual)
	tokens := 0
	
	// URL/Content tokens
	tokens += len(node.Content) / 4
	
	// Title tokens
	if node.Title != nil {
		tokens += len(*node.Title) / 4
	}
	
	// Description tokens
	if node.Description != nil {
		tokens += len(*node.Description) / 4
	}
	
	// Attribute tokens (if included)
	if sc.IncludeAttributes && node.Attributes != nil {
		for _, attr := range node.Attributes {
			tokens += len(attr.Name) / 4
			tokens += len(attr.Value) / 4
			if attr.AttributeType != nil {
				tokens += len(*attr.AttributeType) / 4
			}
		}
	}
	
	// JSON structure overhead (~20% additional tokens)
	tokens = int(float64(tokens) * 1.2)
	
	// Minimum tokens per node
	if tokens < constants.MinTokensPerNode {
		tokens = constants.MinTokensPerNode
	}
	
	return tokens
}

// Legacy SmartChunker methods kept for backwards compatibility but not used in page-based scanning