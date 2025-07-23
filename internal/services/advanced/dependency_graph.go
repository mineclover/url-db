package services

import (
	"fmt"
	"os"
	"time"
	"url-db/internal/models"
	"url-db/internal/repositories"
)

// DependencyGraphService provides graph operations for dependencies
type DependencyGraphService struct {
	dependencyRepo *repositories.DependencyRepository
	nodeRepo       repositories.NodeRepository
}

// NewDependencyGraphService creates a new dependency graph service
func NewDependencyGraphService(depRepo *repositories.DependencyRepository, nodeRepo repositories.NodeRepository) *DependencyGraphService {
	return &DependencyGraphService{
		dependencyRepo: depRepo,
		nodeRepo:       nodeRepo,
	}
}

// graphNode represents a node in the dependency graph for traversal
type graphNode struct {
	ID      int64
	Visited bool
	InStack bool
	LowLink int
	Index   int
	Parent  *graphNode
}

// DetectCycles uses Tarjan's algorithm to find all cycles in the dependency graph
func (s *DependencyGraphService) DetectCycles(domainID int64) ([]models.CircularDependency, error) {
	// Build adjacency list
	graph, err := s.buildGraph(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to build graph: %w", err)
	}

	// Initialize Tarjan's algorithm
	index := 0
	stack := make([]int64, 0)
	nodes := make(map[int64]*graphNode)
	cycles := make([]models.CircularDependency, 0)

	// Initialize all nodes
	for nodeID := range graph {
		nodes[nodeID] = &graphNode{
			ID:      nodeID,
			Index:   -1,
			LowLink: -1,
		}
	}

	// Run Tarjan's algorithm on all unvisited nodes
	for nodeID := range graph {
		if nodes[nodeID].Index == -1 {
			s.tarjanDFS(nodeID, graph, nodes, &stack, &index, &cycles)
		}
	}

	// Get node details for cycles
	for i := range cycles {
		details, err := s.getNodeDetails(cycles[i].Path)
		if err != nil {
			continue
		}
		cycles[i].NodeDetails = details
		cycles[i].Strength = s.calculateCycleStrength(cycles[i].Path, graph)
	}

	return cycles, nil
}

// tarjanDFS performs depth-first search for Tarjan's algorithm
func (s *DependencyGraphService) tarjanDFS(
	nodeID int64,
	graph map[int64][]dependencyEdge,
	nodes map[int64]*graphNode,
	stack *[]int64,
	index *int,
	cycles *[]models.CircularDependency,
) {
	node := nodes[nodeID]
	node.Index = *index
	node.LowLink = *index
	*index++
	*stack = append(*stack, nodeID)
	node.InStack = true

	// Visit all neighbors
	for _, edge := range graph[nodeID] {
		neighbor := nodes[edge.To]

		if neighbor.Index == -1 {
			// Neighbor not visited, recurse
			neighbor.Parent = node
			s.tarjanDFS(edge.To, graph, nodes, stack, index, cycles)

			// Update low link
			if neighbor.LowLink < node.LowLink {
				node.LowLink = neighbor.LowLink
			}
		} else if neighbor.InStack {
			// Neighbor is in stack, we found a cycle
			if neighbor.Index < node.LowLink {
				node.LowLink = neighbor.Index
			}

			// Extract cycle
			cycle := s.extractCycle(nodeID, edge.To, nodes)
			if len(cycle) > 0 {
				*cycles = append(*cycles, models.CircularDependency{
					Path: cycle,
				})
			}
		}
	}

	// Check if node is a root of SCC
	if node.LowLink == node.Index {
		// Pop from stack until we reach current node
		for {
			if len(*stack) == 0 {
				break
			}

			topID := (*stack)[len(*stack)-1]
			*stack = (*stack)[:len(*stack)-1]
			nodes[topID].InStack = false

			if topID == nodeID {
				break
			}
		}
	}
}

// extractCycle extracts a cycle path from the graph
func (s *DependencyGraphService) extractCycle(from, to int64, nodes map[int64]*graphNode) []int64 {
	path := []int64{from}
	current := nodes[from].Parent

	for current != nil && current.ID != to {
		path = append([]int64{current.ID}, path...)
		current = current.Parent

		// Prevent infinite loop
		if len(path) > 100 {
			break
		}
	}

	if current != nil && current.ID == to {
		path = append([]int64{to}, path...)
	}

	return path
}

// ValidateNewDependency checks if adding a dependency would create a cycle
func (s *DependencyGraphService) ValidateNewDependency(
	dependentID, dependencyID int64,
) (*models.DependencyValidationResult, error) {
	result := &models.DependencyValidationResult{
		IsValid:  true,
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
		Cycles:   make([]models.CircularDependency, 0),
	}

	// Check for self-dependency
	if dependentID == dependencyID {
		result.IsValid = false
		result.Errors = append(result.Errors, "Self-dependencies are not allowed")
		return result, nil
	}

	// Check if dependency already exists
	existing, err := s.dependencyRepo.GetDependency(dependentID, dependencyID)
	if err == nil && existing != nil {
		result.Warnings = append(result.Warnings, "Dependency already exists")
	}

	// Check if adding this would create a cycle
	wouldCreateCycle, cyclePath, err := s.wouldCreateCycle(dependentID, dependencyID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for cycles: %w", err)
	}

	if wouldCreateCycle {
		result.IsValid = false
		result.Errors = append(result.Errors, "Adding this dependency would create a circular dependency")

		details, _ := s.getNodeDetails(cyclePath)
		result.Cycles = append(result.Cycles, models.CircularDependency{
			Path:        cyclePath,
			NodeDetails: details,
		})
	}

	// Check domain rules
	rules, err := s.dependencyRepo.GetActiveRules(dependentID)
	if err == nil {
		for _, rule := range rules {
			if !s.validateRule(*rule, dependentID, dependencyID) {
				result.IsValid = false
				result.Errors = append(result.Errors, fmt.Sprintf("Violates rule: %s", rule.RuleName))
			}
		}
	}

	return result, nil
}

// wouldCreateCycle checks if adding a dependency would create a cycle
func (s *DependencyGraphService) wouldCreateCycle(dependentID, dependencyID int64) (bool, []int64, error) {
	// Use DFS to check if there's already a path from dependency to dependent
	visited := make(map[int64]bool)
	path := []int64{dependencyID}

	return s.dfsCheckPath(dependencyID, dependentID, visited, &path)
}

// dfsCheckPath performs depth-first search to find a path
func (s *DependencyGraphService) dfsCheckPath(
	current, target int64,
	visited map[int64]bool,
	path *[]int64,
) (bool, []int64, error) {
	if current == target {
		return true, *path, nil
	}

	if visited[current] {
		return false, nil, nil
	}

	visited[current] = true

	// Get all dependencies of current node
	deps, err := s.dependencyRepo.GetNodeDependencies(current)
	if err != nil {
		return false, nil, err
	}

	for _, dep := range deps {
		*path = append(*path, dep.DependencyNodeID)

		found, cyclePath, err := s.dfsCheckPath(dep.DependencyNodeID, target, visited, path)
		if err != nil {
			return false, nil, err
		}

		if found {
			return true, cyclePath, nil
		}

		// Backtrack
		*path = (*path)[:len(*path)-1]
	}

	return false, nil, nil
}

// GetDependencyGraph builds a complete dependency graph for a node
func (s *DependencyGraphService) GetDependencyGraph(
	nodeID int64,
	maxDepth int,
) (*models.DependencyGraph, error) {
	graph := &models.DependencyGraph{
		NodeID:       nodeID,
		Dependencies: make([]models.DependencyNode, 0),
		Dependents:   make([]models.DependencyNode, 0),
		Metadata:     make(map[string]interface{}),
	}

	// Check cache first
	cached, err := s.dependencyRepo.GetCachedGraph(nodeID)
	if err == nil && cached != nil && !s.isCacheExpired(cached) {
		// Parse and return cached graph
		return s.parseCachedGraph(cached)
	}

	// Build dependencies tree
	depVisited := make(map[int64]bool)
	dependencies, depthDeps, err := s.buildDependencyTree(nodeID, 0, maxDepth, depVisited)
	if err != nil {
		return nil, fmt.Errorf("failed to build dependency tree: %w", err)
	}
	graph.Dependencies = dependencies

	// Build dependents tree
	depVisited = make(map[int64]bool)
	dependents, depthDents, err := s.buildDependentTree(nodeID, 0, maxDepth, depVisited)
	if err != nil {
		return nil, fmt.Errorf("failed to build dependent tree: %w", err)
	}
	graph.Dependents = dependents

	// Calculate metadata
	graph.Depth = max(depthDeps, depthDents)
	graph.HasCircular = s.hasCircularDependency(nodeID)
	graph.Metadata["total_dependencies"] = len(dependencies)
	graph.Metadata["total_dependents"] = len(dependents)

	// Cache the result
	if err := s.cacheGraph(nodeID, graph); err != nil {
		// Log error but don't fail
		fmt.Fprintf(os.Stderr, "Failed to cache graph: %v\n", err)
	}

	return graph, nil
}

// buildDependencyTree recursively builds the dependency tree
func (s *DependencyGraphService) buildDependencyTree(
	nodeID int64,
	currentDepth, maxDepth int,
	visited map[int64]bool,
) ([]models.DependencyNode, int, error) {
	if currentDepth >= maxDepth || visited[nodeID] {
		return nil, currentDepth, nil
	}

	visited[nodeID] = true

	deps, err := s.dependencyRepo.GetNodeDependencies(nodeID)
	if err != nil {
		return nil, currentDepth, err
	}

	nodes := make([]models.DependencyNode, 0, len(deps))
	maxFoundDepth := currentDepth

	for _, dep := range deps {
		node, err := s.nodeRepo.GetByID(int(dep.DependencyNodeID))
		if err != nil {
			continue
		}

		depNode := models.DependencyNode{
			NodeID:         dep.DependencyNodeID,
			CompositeID:    s.buildCompositeID(node),
			Title:          node.Title,
			DependencyType: dep.DependencyType,
			Category:       s.getCategoryForType(dep.DependencyType),
			Strength:       s.getStrength(dep),
			Priority:       s.getPriority(dep),
			IsRequired:     s.getIsRequired(dep),
		}

		// Parse metadata
		if dep.Metadata != nil {
			metadata, _ := dep.ParseMetadata()
			depNode.Metadata = metadata
		}

		// Recursively get children
		children, childDepth, _ := s.buildDependencyTree(
			dep.DependencyNodeID,
			currentDepth+1,
			maxDepth,
			visited,
		)
		depNode.Children = children

		if childDepth > maxFoundDepth {
			maxFoundDepth = childDepth
		}

		nodes = append(nodes, depNode)
	}

	return nodes, maxFoundDepth, nil
}

// buildDependentTree recursively builds the dependent tree
func (s *DependencyGraphService) buildDependentTree(
	nodeID int64,
	currentDepth, maxDepth int,
	visited map[int64]bool,
) ([]models.DependencyNode, int, error) {
	if currentDepth >= maxDepth || visited[nodeID] {
		return nil, currentDepth, nil
	}

	visited[nodeID] = true

	deps, err := s.dependencyRepo.GetNodeDependents(nodeID)
	if err != nil {
		return nil, currentDepth, err
	}

	nodes := make([]models.DependencyNode, 0, len(deps))
	maxFoundDepth := currentDepth

	for _, dep := range deps {
		node, err := s.nodeRepo.GetByID(int(dep.DependentNodeID))
		if err != nil {
			continue
		}

		depNode := models.DependencyNode{
			NodeID:         dep.DependentNodeID,
			CompositeID:    s.buildCompositeID(node),
			Title:          node.Title,
			DependencyType: dep.DependencyType,
			Category:       s.getCategoryForType(dep.DependencyType),
			Strength:       s.getStrength(dep),
			Priority:       s.getPriority(dep),
			IsRequired:     s.getIsRequired(dep),
		}

		// Parse metadata
		if dep.Metadata != nil {
			metadata, _ := dep.ParseMetadata()
			depNode.Metadata = metadata
		}

		// Recursively get children
		children, childDepth, _ := s.buildDependentTree(
			dep.DependentNodeID,
			currentDepth+1,
			maxDepth,
			visited,
		)
		depNode.Children = children

		if childDepth > maxFoundDepth {
			maxFoundDepth = childDepth
		}

		nodes = append(nodes, depNode)
	}

	return nodes, maxFoundDepth, nil
}

// Helper types and methods

type dependencyEdge struct {
	To       int64
	Type     string
	Strength int
}

func (s *DependencyGraphService) buildGraph(domainID int64) (map[int64][]dependencyEdge, error) {
	// Implementation would fetch all dependencies for the domain
	// and build an adjacency list representation
	return nil, fmt.Errorf("not implemented")
}

func (s *DependencyGraphService) getNodeDetails(nodeIDs []int64) ([]string, error) {
	details := make([]string, 0, len(nodeIDs))
	for _, id := range nodeIDs {
		node, err := s.nodeRepo.GetByID(int(id))
		if err != nil {
			details = append(details, fmt.Sprintf("Node %d (unknown)", id))
			continue
		}
		details = append(details, fmt.Sprintf("%s (%s)", node.Title, s.buildCompositeID(node)))
	}
	return details, nil
}

func (s *DependencyGraphService) calculateCycleStrength(path []int64, graph map[int64][]dependencyEdge) int {
	minStrength := 100
	for i := 0; i < len(path)-1; i++ {
		from := path[i]
		to := path[i+1]

		for _, edge := range graph[from] {
			if edge.To == to && edge.Strength < minStrength {
				minStrength = edge.Strength
			}
		}
	}
	return minStrength
}

func (s *DependencyGraphService) validateRule(rule models.DependencyRule, dependentID, dependencyID int64) bool {
	// Implementation would validate based on rule type and config
	return true
}

func (s *DependencyGraphService) isCacheExpired(cache *models.DependencyGraphCache) bool {
	if cache.ExpiresAt == nil {
		return false
	}
	return cache.ExpiresAt.Before(time.Now())
}

func (s *DependencyGraphService) parseCachedGraph(cache *models.DependencyGraphCache) (*models.DependencyGraph, error) {
	// Implementation would parse JSON graph data
	return nil, fmt.Errorf("not implemented")
}

func (s *DependencyGraphService) cacheGraph(nodeID int64, graph *models.DependencyGraph) error {
	// Implementation would serialize and cache the graph
	return nil
}

func (s *DependencyGraphService) hasCircularDependency(nodeID int64) bool {
	// Quick check for circular dependencies involving this node
	return false
}

func (s *DependencyGraphService) buildCompositeID(node *models.Node) string {
	// Implementation would build composite ID
	return fmt.Sprintf("url-db:domain:%d", node.ID)
}

func (s *DependencyGraphService) getCategoryForType(depType string) string {
	switch depType {
	case models.DependencyTypeHard, models.DependencyTypeSoft, models.DependencyTypeReference:
		return models.CategoryStructural
	case models.DependencyTypeRuntime, models.DependencyTypeCompile, models.DependencyTypeOptional:
		return models.CategoryBehavioral
	case models.DependencyTypeSync, models.DependencyTypeAsync:
		return models.CategoryData
	default:
		return models.CategoryStructural
	}
}

func (s *DependencyGraphService) getStrength(dep *models.NodeDependency) int {
	// For V1, use default strengths based on type
	switch dep.DependencyType {
	case models.DependencyTypeHard:
		return 90
	case models.DependencyTypeSoft:
		return 50
	case models.DependencyTypeReference:
		return 30
	default:
		return 50
	}
}

func (s *DependencyGraphService) getPriority(dep *models.NodeDependency) int {
	// For V1, use default priorities
	return 50
}

func (s *DependencyGraphService) getIsRequired(dep *models.NodeDependency) bool {
	// For V1, check type
	return dep.DependencyType != models.DependencyTypeOptional &&
		dep.DependencyType != models.DependencyTypeReference
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
