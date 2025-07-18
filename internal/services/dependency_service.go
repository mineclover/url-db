package services

import (
	"fmt"

	"url-db/internal/models"
	"url-db/internal/repositories"
)

// DependencyService handles business logic for dependencies
type DependencyService struct {
	dependencyRepo *repositories.DependencyRepository
	nodeRepo       *repositories.NodeRepository
	eventRepo      *repositories.EventRepository
}

// NewDependencyService creates a new dependency service
func NewDependencyService(
	dependencyRepo *repositories.DependencyRepository,
	nodeRepo *repositories.NodeRepository,
	eventRepo *repositories.EventRepository,
) *DependencyService {
	return &DependencyService{
		dependencyRepo: dependencyRepo,
		nodeRepo:       nodeRepo,
		eventRepo:      eventRepo,
	}
}

// CreateDependency creates a new dependency
func (s *DependencyService) CreateDependency(dependentNodeID int64, req *models.CreateNodeDependencyRequest) (*models.NodeDependency, error) {
	// Verify both nodes exist
	dependentNode, err := s.nodeRepo.GetByID(dependentNodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependent node: %w", err)
	}
	if dependentNode == nil {
		return nil, fmt.Errorf("dependent node not found")
	}
	
	dependencyNode, err := s.nodeRepo.GetByID(req.DependencyNodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependency node: %w", err)
	}
	if dependencyNode == nil {
		return nil, fmt.Errorf("dependency node not found")
	}
	
	// Check for circular dependency
	hasCircular, err := s.dependencyRepo.CheckCircularDependency(dependentNodeID, req.DependencyNodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check circular dependency: %w", err)
	}
	if hasCircular {
		return nil, fmt.Errorf("circular dependency detected")
	}
	
	// Create dependency
	dependency := &models.NodeDependency{
		DependentNodeID:  dependentNodeID,
		DependencyNodeID: req.DependencyNodeID,
		DependencyType:   req.DependencyType,
		CascadeDelete:    req.CascadeDelete,
		CascadeUpdate:    req.CascadeUpdate,
		Metadata:         req.Metadata,
	}
	
	err = s.dependencyRepo.Create(dependency)
	if err != nil {
		return nil, fmt.Errorf("failed to create dependency: %w", err)
	}
	
	return dependency, nil
}

// GetDependency retrieves a dependency by ID
func (s *DependencyService) GetDependency(id int64) (*models.NodeDependency, error) {
	dependency, err := s.dependencyRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependency: %w", err)
	}
	if dependency == nil {
		return nil, fmt.Errorf("dependency not found")
	}
	
	return dependency, nil
}

// DeleteDependency deletes a dependency
func (s *DependencyService) DeleteDependency(id int64) error {
	err := s.dependencyRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete dependency: %w", err)
	}
	
	return nil
}

// GetNodeDependencies retrieves all dependencies for a node (as dependent)
func (s *DependencyService) GetNodeDependencies(nodeID int64) ([]*models.NodeDependency, error) {
	// Verify node exists
	node, err := s.nodeRepo.GetByID(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil {
		return nil, fmt.Errorf("node not found")
	}
	
	return s.dependencyRepo.GetByDependentNode(nodeID)
}

// GetNodeDependents retrieves all nodes that depend on a given node
func (s *DependencyService) GetNodeDependents(nodeID int64) ([]*models.NodeDependency, error) {
	// Verify node exists
	node, err := s.nodeRepo.GetByID(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}
	if node == nil {
		return nil, fmt.Errorf("node not found")
	}
	
	return s.dependencyRepo.GetByDependencyNode(nodeID)
}

// HandleNodeDeletion handles cascading deletes when a node is deleted
func (s *DependencyService) HandleNodeDeletion(nodeID int64) error {
	// Get all dependents with cascade delete
	dependentIDs, err := s.dependencyRepo.GetDependentsWithCascadeDelete(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get cascade delete dependents: %w", err)
	}
	
	// Delete each dependent node (this will trigger their own cascades)
	for _, depID := range dependentIDs {
		err = s.nodeRepo.Delete(depID)
		if err != nil {
			return fmt.Errorf("failed to delete dependent node %d: %w", depID, err)
		}
	}
	
	return nil
}

// HandleNodeUpdate handles cascading updates when a node is updated
func (s *DependencyService) HandleNodeUpdate(nodeID int64, changes *models.EventChanges) error {
	// Get all dependents with cascade update
	dependentIDs, err := s.dependencyRepo.GetDependentsWithCascadeUpdate(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get cascade update dependents: %w", err)
	}
	
	// Create update event for each dependent
	for _, depID := range dependentIDs {
		eventData := &models.EventData{
			NodeID:    depID,
			EventType: "dependency_updated",
			Metadata: map[string]interface{}{
				"dependency_node_id": nodeID,
				"changes":           changes,
			},
		}
		
		event := &models.NodeEvent{
			NodeID:    depID,
			EventType: "dependency_updated",
			EventData: eventData,
		}
		
		err = s.eventRepo.Create(event)
		if err != nil {
			return fmt.Errorf("failed to create update event for dependent %d: %w", depID, err)
		}
	}
	
	return nil
}