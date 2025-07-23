package services

import (
	"context"
	"fmt"
	"time"
	"url-db/internal/models"
	"url-db/internal/repositories"
)

// DependencyImpactAnalyzer provides impact analysis for dependency changes
type DependencyImpactAnalyzer struct {
	dependencyRepo *repositories.DependencyRepository
	nodeRepo       repositories.NodeRepository
	graphService   *DependencyGraphService
}

// NewDependencyImpactAnalyzer creates a new impact analyzer
func NewDependencyImpactAnalyzer(
	depRepo *repositories.DependencyRepository,
	nodeRepo repositories.NodeRepository,
	graphService *DependencyGraphService,
) *DependencyImpactAnalyzer {
	return &DependencyImpactAnalyzer{
		dependencyRepo: depRepo,
		nodeRepo:       nodeRepo,
		graphService:   graphService,
	}
}

// AnalyzeImpact performs comprehensive impact analysis for a dependency change
func (a *DependencyImpactAnalyzer) AnalyzeImpact(
	ctx context.Context,
	sourceNodeID int64,
	impactType string,
) (*models.ImpactAnalysisResult, error) {
	// Get source node
	sourceNode, err := a.nodeRepo.GetByID(int(sourceNodeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get source node: %w", err)
	}

	result := &models.ImpactAnalysisResult{
		SourceNodeID:    sourceNodeID,
		SourceComposite: a.buildCompositeID(sourceNode),
		ImpactType:      impactType,
		AffectedNodes:   make([]models.AffectedNode, 0),
		Warnings:        make([]string, 0),
		Recommendations: make([]string, 0),
	}

	// Analyze based on impact type
	switch impactType {
	case "delete":
		return a.analyzeDeleteImpact(ctx, sourceNodeID, result)
	case "update":
		return a.analyzeUpdateImpact(ctx, sourceNodeID, result)
	case "version_change":
		return a.analyzeVersionChangeImpact(ctx, sourceNodeID, result)
	default:
		return nil, fmt.Errorf("unsupported impact type: %s", impactType)
	}
}

// analyzeDeleteImpact analyzes the impact of deleting a node
func (a *DependencyImpactAnalyzer) analyzeDeleteImpact(
	ctx context.Context,
	nodeID int64,
	result *models.ImpactAnalysisResult,
) (*models.ImpactAnalysisResult, error) {
	// Find all nodes that depend on this node
	dependents, err := a.dependencyRepo.GetNodeDependents(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependents: %w", err)
	}

	visited := make(map[int64]bool)
	maxDepth := 0
	totalImpactScore := 0

	for _, dep := range dependents {
		impact, depth := a.analyzeNodeDeleteImpact(ctx, dep, []int64{nodeID}, visited)
		if impact != nil {
			result.AffectedNodes = append(result.AffectedNodes, *impact)
			totalImpactScore += a.getImpactScore(impact.ImpactLevel)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}

	// Calculate overall impact score and metadata
	result.ImpactScore = a.normalizeImpactScore(totalImpactScore, len(result.AffectedNodes))
	result.CascadeDepth = maxDepth
	result.EstimatedTime = a.estimateDeleteTime(len(result.AffectedNodes), maxDepth)

	// Add recommendations
	a.addDeleteRecommendations(result)

	return result, nil
}

// analyzeUpdateImpact analyzes the impact of updating a node
func (a *DependencyImpactAnalyzer) analyzeUpdateImpact(
	ctx context.Context,
	nodeID int64,
	result *models.ImpactAnalysisResult,
) (*models.ImpactAnalysisResult, error) {
	// Find nodes that depend on this node with cascade_update=true
	dependents, err := a.dependencyRepo.GetNodeDependents(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependents: %w", err)
	}

	visited := make(map[int64]bool)
	maxDepth := 0
	totalImpactScore := 0

	for _, dep := range dependents {
		if dep.CascadeUpdate {
			impact, depth := a.analyzeNodeUpdateImpact(ctx, dep, []int64{nodeID}, visited)
			if impact != nil {
				result.AffectedNodes = append(result.AffectedNodes, *impact)
				totalImpactScore += a.getImpactScore(impact.ImpactLevel)
				if depth > maxDepth {
					maxDepth = depth
				}
			}
		}
	}

	result.ImpactScore = a.normalizeImpactScore(totalImpactScore, len(result.AffectedNodes))
	result.CascadeDepth = maxDepth
	result.EstimatedTime = a.estimateUpdateTime(len(result.AffectedNodes), maxDepth)

	a.addUpdateRecommendations(result)

	return result, nil
}

// analyzeVersionChangeImpact analyzes the impact of changing a node's version
func (a *DependencyImpactAnalyzer) analyzeVersionChangeImpact(
	ctx context.Context,
	nodeID int64,
	result *models.ImpactAnalysisResult,
) (*models.ImpactAnalysisResult, error) {
	// Find nodes with version constraints on this node
	dependents, err := a.dependencyRepo.GetNodeDependents(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependents: %w", err)
	}

	visited := make(map[int64]bool)
	maxDepth := 0
	totalImpactScore := 0

	for _, dep := range dependents {
		impact, depth := a.analyzeNodeVersionImpact(ctx, dep, []int64{nodeID}, visited)
		if impact != nil {
			result.AffectedNodes = append(result.AffectedNodes, *impact)
			totalImpactScore += a.getImpactScore(impact.ImpactLevel)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}

	result.ImpactScore = a.normalizeImpactScore(totalImpactScore, len(result.AffectedNodes))
	result.CascadeDepth = maxDepth
	result.EstimatedTime = a.estimateVersionChangeTime(len(result.AffectedNodes))

	a.addVersionChangeRecommendations(result)

	return result, nil
}

// analyzeNodeDeleteImpact analyzes impact on a specific node from deletion
func (a *DependencyImpactAnalyzer) analyzeNodeDeleteImpact(
	ctx context.Context,
	dependency *models.NodeDependency,
	path []int64,
	visited map[int64]bool,
) (*models.AffectedNode, int) {
	nodeID := dependency.DependentNodeID

	if visited[nodeID] {
		return nil, 0
	}
	visited[nodeID] = true

	node, err := a.nodeRepo.GetByID(int(nodeID))
	if err != nil {
		return nil, 0
	}

	impactLevel := a.calculateDeleteImpactLevel(dependency)
	actionNeeded := a.determineDeleteAction(dependency, impactLevel)

	affected := &models.AffectedNode{
		NodeID:       nodeID,
		CompositeID:  a.buildCompositeID(node),
		Title:        node.Title,
		ImpactLevel:  impactLevel,
		Reason:       fmt.Sprintf("Depends on deleted node with %s dependency", dependency.DependencyType),
		ActionNeeded: actionNeeded,
		Path:         append([]int64(nil), path...),
	}

	// Recursively analyze cascade effects if cascade_delete is enabled
	depth := 1
	if dependency.CascadeDelete {
		childDeps, _ := a.dependencyRepo.GetNodeDependents(nodeID)
		for _, childDep := range childDeps {
			childPath := append(path, nodeID)
			childAffected, childDepth := a.analyzeNodeDeleteImpact(ctx, childDep, childPath, visited)
			if childAffected != nil {
				depth = max(depth, childDepth+1)
			}
		}
	}

	return affected, depth
}

// analyzeNodeUpdateImpact analyzes impact on a specific node from update
func (a *DependencyImpactAnalyzer) analyzeNodeUpdateImpact(
	ctx context.Context,
	dependency *models.NodeDependency,
	path []int64,
	visited map[int64]bool,
) (*models.AffectedNode, int) {
	nodeID := dependency.DependentNodeID

	if visited[nodeID] {
		return nil, 0
	}
	visited[nodeID] = true

	node, err := a.nodeRepo.GetByID(int(nodeID))
	if err != nil {
		return nil, 0
	}

	impactLevel := a.calculateUpdateImpactLevel(dependency)
	actionNeeded := a.determineUpdateAction(dependency, impactLevel)

	affected := &models.AffectedNode{
		NodeID:       nodeID,
		CompositeID:  a.buildCompositeID(node),
		Title:        node.Title,
		ImpactLevel:  impactLevel,
		Reason:       fmt.Sprintf("May need update due to %s dependency change", dependency.DependencyType),
		ActionNeeded: actionNeeded,
		Path:         append([]int64(nil), path...),
	}

	// Recursively analyze cascade effects
	depth := 1
	if dependency.CascadeUpdate {
		childDeps, _ := a.dependencyRepo.GetNodeDependents(nodeID)
		for _, childDep := range childDeps {
			if childDep.CascadeUpdate {
				childPath := append(path, nodeID)
				childAffected, childDepth := a.analyzeNodeUpdateImpact(ctx, childDep, childPath, visited)
				if childAffected != nil {
					depth = max(depth, childDepth+1)
				}
			}
		}
	}

	return affected, depth
}

// analyzeNodeVersionImpact analyzes version constraint impact
func (a *DependencyImpactAnalyzer) analyzeNodeVersionImpact(
	ctx context.Context,
	dependency *models.NodeDependency,
	path []int64,
	visited map[int64]bool,
) (*models.AffectedNode, int) {
	nodeID := dependency.DependentNodeID

	if visited[nodeID] {
		return nil, 0
	}
	visited[nodeID] = true

	node, err := a.nodeRepo.GetByID(int(nodeID))
	if err != nil {
		return nil, 0
	}

	// Check if this dependency has version constraints
	hasVersionConstraint := false
	if metadata, err := dependency.ParseMetadata(); err == nil && metadata != nil {
		if _, exists := metadata["version_constraint"]; exists {
			hasVersionConstraint = true
		}
	}

	if !hasVersionConstraint {
		return nil, 0
	}

	impactLevel := a.calculateVersionImpactLevel(dependency)
	actionNeeded := a.determineVersionAction(dependency, impactLevel)

	affected := &models.AffectedNode{
		NodeID:       nodeID,
		CompositeID:  a.buildCompositeID(node),
		Title:        node.Title,
		ImpactLevel:  impactLevel,
		Reason:       "Has version constraint that may be affected",
		ActionNeeded: actionNeeded,
		Path:         append([]int64(nil), path...),
	}

	return affected, 1
}

// Helper methods for impact calculation

func (a *DependencyImpactAnalyzer) calculateDeleteImpactLevel(dep *models.NodeDependency) string {
	if dep.CascadeDelete {
		return models.ImpactLevelCritical
	}

	switch dep.DependencyType {
	case models.DependencyTypeHard:
		return models.ImpactLevelHigh
	case models.DependencyTypeSoft:
		return models.ImpactLevelMedium
	case models.DependencyTypeReference:
		return models.ImpactLevelLow
	default:
		return models.ImpactLevelMedium
	}
}

func (a *DependencyImpactAnalyzer) calculateUpdateImpactLevel(dep *models.NodeDependency) string {
	if dep.CascadeUpdate {
		return models.ImpactLevelHigh
	}

	switch dep.DependencyType {
	case models.DependencyTypeHard:
		return models.ImpactLevelMedium
	case models.DependencyTypeSoft:
		return models.ImpactLevelLow
	case models.DependencyTypeReference:
		return models.ImpactLevelLow
	default:
		return models.ImpactLevelLow
	}
}

func (a *DependencyImpactAnalyzer) calculateVersionImpactLevel(dep *models.NodeDependency) string {
	switch dep.DependencyType {
	case models.DependencyTypeHard, models.DependencyTypeRuntime:
		return models.ImpactLevelHigh
	case models.DependencyTypeCompile:
		return models.ImpactLevelMedium
	default:
		return models.ImpactLevelLow
	}
}

func (a *DependencyImpactAnalyzer) determineDeleteAction(dep *models.NodeDependency, impactLevel string) string {
	if dep.CascadeDelete {
		return "Will be automatically deleted"
	}

	switch impactLevel {
	case models.ImpactLevelCritical, models.ImpactLevelHigh:
		return "Review and update or remove dependency"
	case models.ImpactLevelMedium:
		return "Consider updating configuration"
	case models.ImpactLevelLow:
		return "Update documentation if needed"
	default:
		return "Monitor for issues"
	}
}

func (a *DependencyImpactAnalyzer) determineUpdateAction(dep *models.NodeDependency, impactLevel string) string {
	if dep.CascadeUpdate {
		return "Will be automatically updated"
	}

	switch impactLevel {
	case models.ImpactLevelHigh:
		return "Review and test changes"
	case models.ImpactLevelMedium:
		return "Monitor for compatibility issues"
	case models.ImpactLevelLow:
		return "Update documentation if needed"
	default:
		return "No action required"
	}
}

func (a *DependencyImpactAnalyzer) determineVersionAction(dep *models.NodeDependency, impactLevel string) string {
	switch impactLevel {
	case models.ImpactLevelHigh:
		return "Verify version compatibility and update constraints"
	case models.ImpactLevelMedium:
		return "Test with new version"
	case models.ImpactLevelLow:
		return "Update version constraint if needed"
	default:
		return "Monitor for issues"
	}
}

func (a *DependencyImpactAnalyzer) getImpactScore(level string) int {
	switch level {
	case models.ImpactLevelCritical:
		return 100
	case models.ImpactLevelHigh:
		return 75
	case models.ImpactLevelMedium:
		return 50
	case models.ImpactLevelLow:
		return 25
	default:
		return 25
	}
}

func (a *DependencyImpactAnalyzer) normalizeImpactScore(total, count int) int {
	if count == 0 {
		return 0
	}

	avg := total / count

	// Apply curve based on count of affected nodes
	if count > 10 {
		avg = int(float64(avg) * 1.2) // Amplify for many affected nodes
	} else if count > 5 {
		avg = int(float64(avg) * 1.1)
	}

	if avg > 100 {
		return 100
	}
	return avg
}

func (a *DependencyImpactAnalyzer) estimateDeleteTime(affectedCount, depth int) string {
	base := time.Duration(affectedCount) * 5 * time.Minute // 5 minutes per affected node
	cascade := time.Duration(depth) * 10 * time.Minute     // 10 minutes per depth level
	total := base + cascade

	if total < time.Hour {
		return fmt.Sprintf("%d minutes", int(total.Minutes()))
	}
	return fmt.Sprintf("%.1f hours", total.Hours())
}

func (a *DependencyImpactAnalyzer) estimateUpdateTime(affectedCount, depth int) string {
	base := time.Duration(affectedCount) * 2 * time.Minute // 2 minutes per affected node
	cascade := time.Duration(depth) * 5 * time.Minute      // 5 minutes per depth level
	total := base + cascade

	if total < time.Hour {
		return fmt.Sprintf("%d minutes", int(total.Minutes()))
	}
	return fmt.Sprintf("%.1f hours", total.Hours())
}

func (a *DependencyImpactAnalyzer) estimateVersionChangeTime(affectedCount int) string {
	total := time.Duration(affectedCount) * 3 * time.Minute // 3 minutes per affected node

	if total < time.Hour {
		return fmt.Sprintf("%d minutes", int(total.Minutes()))
	}
	return fmt.Sprintf("%.1f hours", total.Hours())
}

func (a *DependencyImpactAnalyzer) addDeleteRecommendations(result *models.ImpactAnalysisResult) {
	criticalCount := 0
	highCount := 0

	for _, node := range result.AffectedNodes {
		switch node.ImpactLevel {
		case models.ImpactLevelCritical:
			criticalCount++
		case models.ImpactLevelHigh:
			highCount++
		}
	}

	if criticalCount > 0 {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("⚠️  %d nodes will be automatically deleted due to cascade delete", criticalCount))
		result.Recommendations = append(result.Recommendations,
			"Consider backing up data before proceeding")
	}

	if highCount > 0 {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("🔍 Review %d nodes with high impact dependencies", highCount))
	}

	if len(result.AffectedNodes) > 5 {
		result.Recommendations = append(result.Recommendations,
			"Consider phased deletion approach to minimize disruption")
	}
}

func (a *DependencyImpactAnalyzer) addUpdateRecommendations(result *models.ImpactAnalysisResult) {
	if len(result.AffectedNodes) > 0 {
		result.Recommendations = append(result.Recommendations,
			"Test in staging environment before production deployment")
	}

	if result.CascadeDepth > 3 {
		result.Recommendations = append(result.Recommendations,
			"Deep cascade detected - consider gradual rollout")
	}
}

func (a *DependencyImpactAnalyzer) addVersionChangeRecommendations(result *models.ImpactAnalysisResult) {
	if len(result.AffectedNodes) > 0 {
		result.Recommendations = append(result.Recommendations,
			"Verify version constraints are compatible")
		result.Recommendations = append(result.Recommendations,
			"Run compatibility tests before deployment")
	}
}

func (a *DependencyImpactAnalyzer) buildCompositeID(node *models.Node) string {
	return fmt.Sprintf("url-db:%d", node.ID)
}
