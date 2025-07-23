package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"url-db/internal/models"
)

// DependencyRepository handles database operations for dependencies
type DependencyRepository struct {
	db *sqlx.DB
}

// NewDependencyRepository creates a new dependency repository
func NewDependencyRepository(db *sqlx.DB) *DependencyRepository {
	return &DependencyRepository{db: db}
}

// Create creates a new dependency
func (r *DependencyRepository) Create(dependency *models.NodeDependency) error {
	query := `
		INSERT INTO node_dependencies (
			dependent_node_id, dependency_node_id, dependency_type,
			cascade_delete, cascade_update, metadata
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		dependency.DependentNodeID,
		dependency.DependencyNodeID,
		dependency.DependencyType,
		dependency.CascadeDelete,
		dependency.CascadeUpdate,
		dependency.Metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to create dependency: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	dependency.ID = id
	return nil
}

// GetByID retrieves a dependency by ID
func (r *DependencyRepository) GetByID(id int64) (*models.NodeDependency, error) {
	var dependency models.NodeDependency
	query := `
		SELECT id, dependent_node_id, dependency_node_id, dependency_type,
			   cascade_delete, cascade_update, metadata, created_at
		FROM node_dependencies
		WHERE id = ?
	`

	err := r.db.Get(&dependency, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get dependency: %w", err)
	}

	return &dependency, nil
}

// GetByDependentNode retrieves all dependencies for a dependent node
func (r *DependencyRepository) GetByDependentNode(nodeID int64) ([]*models.NodeDependency, error) {
	var dependencies []*models.NodeDependency
	query := `
		SELECT id, dependent_node_id, dependency_node_id, dependency_type,
			   cascade_delete, cascade_update, metadata, created_at
		FROM node_dependencies
		WHERE dependent_node_id = ?
		ORDER BY created_at DESC
	`

	err := r.db.Select(&dependencies, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}

	return dependencies, nil
}

// GetByDependencyNode retrieves all dependencies where a node is the dependency
func (r *DependencyRepository) GetByDependencyNode(nodeID int64) ([]*models.NodeDependency, error) {
	var dependencies []*models.NodeDependency
	query := `
		SELECT id, dependent_node_id, dependency_node_id, dependency_type,
			   cascade_delete, cascade_update, metadata, created_at
		FROM node_dependencies
		WHERE dependency_node_id = ?
		ORDER BY created_at DESC
	`

	err := r.db.Select(&dependencies, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}

	return dependencies, nil
}

// CheckCircularDependency checks if creating a dependency would create a circular reference
func (r *DependencyRepository) CheckCircularDependency(dependentID, dependencyID int64) (bool, error) {
	// Recursive CTE to check for circular dependencies
	query := `
		WITH RECURSIVE dependency_chain AS (
			SELECT dependent_node_id, dependency_node_id
			FROM node_dependencies
			WHERE dependent_node_id = ?
			
			UNION
			
			SELECT nd.dependent_node_id, nd.dependency_node_id
			FROM node_dependencies nd
			INNER JOIN dependency_chain dc ON nd.dependent_node_id = dc.dependency_node_id
		)
		SELECT COUNT(*) FROM dependency_chain WHERE dependency_node_id = ?
	`

	var count int
	err := r.db.Get(&count, query, dependencyID, dependentID)
	if err != nil {
		return false, fmt.Errorf("failed to check circular dependency: %w", err)
	}

	return count > 0, nil
}

// Delete deletes a dependency
func (r *DependencyRepository) Delete(id int64) error {
	query := "DELETE FROM node_dependencies WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete dependency: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetDependentsWithCascadeDelete retrieves all dependent nodes that should be deleted
func (r *DependencyRepository) GetDependentsWithCascadeDelete(nodeID int64) ([]int64, error) {
	var dependentIDs []int64
	query := `
		SELECT dependent_node_id
		FROM node_dependencies
		WHERE dependency_node_id = ? AND cascade_delete = true
	`

	err := r.db.Select(&dependentIDs, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cascade delete dependents: %w", err)
	}

	return dependentIDs, nil
}

// GetDependentsWithCascadeUpdate retrieves all dependent nodes that should be notified
func (r *DependencyRepository) GetDependentsWithCascadeUpdate(nodeID int64) ([]int64, error) {
	var dependentIDs []int64
	query := `
		SELECT dependent_node_id
		FROM node_dependencies
		WHERE dependency_node_id = ? AND cascade_update = true
	`

	err := r.db.Select(&dependentIDs, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cascade update dependents: %w", err)
	}

	return dependentIDs, nil
}

// GetDependency retrieves a specific dependency by dependent and dependency node IDs
func (r *DependencyRepository) GetDependency(dependentID, dependencyID int64) (*models.NodeDependency, error) {
	var dependency models.NodeDependency
	query := `
		SELECT id, dependent_node_id, dependency_node_id, dependency_type,
			   cascade_delete, cascade_update, metadata, created_at
		FROM node_dependencies
		WHERE dependent_node_id = ? AND dependency_node_id = ?
	`

	err := r.db.Get(&dependency, query, dependentID, dependencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get dependency: %w", err)
	}

	return &dependency, nil
}

// GetNodeDependencies retrieves all dependencies for a node
func (r *DependencyRepository) GetNodeDependencies(nodeID int64) ([]*models.NodeDependency, error) {
	return r.GetByDependentNode(nodeID)
}

// GetNodeDependents retrieves all nodes that depend on the given node
func (r *DependencyRepository) GetNodeDependents(nodeID int64) ([]*models.NodeDependency, error) {
	return r.GetByDependencyNode(nodeID)
}

// GetActiveRules retrieves active dependency rules for a node
func (r *DependencyRepository) GetActiveRules(nodeID int64) ([]*models.DependencyRule, error) {
	var rules []*models.DependencyRule
	// For now, return empty slice as dependency rules are not fully implemented
	return rules, nil
}

// GetCachedGraph retrieves cached dependency graph
func (r *DependencyRepository) GetCachedGraph(nodeID int64) (*models.DependencyGraphCache, error) {
	var cache models.DependencyGraphCache
	query := `
		SELECT id, node_id, graph_data, max_depth, created_at, expires_at
		FROM dependency_graph_cache
		WHERE node_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.Get(&cache, query, nodeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cached graph: %w", err)
	}

	return &cache, nil
}
