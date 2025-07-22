# Enhanced Dependency System Design

## Overview

This document outlines an enhanced dependency tracking and relationship management system for URL-DB, providing advanced features for managing complex dependency graphs, circular dependency detection, impact analysis, and versioning.

## Current System Analysis

### Existing Schema
```sql
CREATE TABLE node_dependencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dependent_node_id INTEGER NOT NULL,
    dependency_node_id INTEGER NOT NULL,
    dependency_type TEXT NOT NULL,
    description TEXT,
    is_required BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Limitations
1. Limited dependency types (only text field)
2. No metadata support for complex relationships
3. No versioning or history tracking
4. No built-in circular dependency detection
5. No impact analysis capabilities
6. No dependency strength or priority

## Enhanced Dependency System

### 1. Dependency Type System

#### Core Dependency Types
```yaml
dependency_types:
  structural:
    hard:
      description: "Strong coupling, deletion cascades"
      cascade_delete: true
      cascade_update: true
      validation_required: true
      
    soft:
      description: "Loose coupling, no cascading"
      cascade_delete: false
      cascade_update: false
      validation_required: false
      
    reference:
      description: "Informational link only"
      cascade_delete: false
      cascade_update: false
      validation_required: false
      
  behavioral:
    runtime:
      description: "Required at runtime"
      health_check: true
      startup_order: true
      
    compile:
      description: "Required at build time"
      version_constraint: true
      
    optional:
      description: "Enhances functionality if present"
      fallback_behavior: true
      
  data:
    sync:
      description: "Data synchronization required"
      sync_frequency: "configurable"
      conflict_resolution: "configurable"
      
    async:
      description: "Eventually consistent data"
      queue_based: true
      retry_policy: "configurable"
```

### 2. Enhanced Database Schema

```sql
-- Dependency types registry
CREATE TABLE dependency_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type_name TEXT NOT NULL UNIQUE,
    category TEXT NOT NULL, -- 'structural', 'behavioral', 'data'
    cascade_delete BOOLEAN DEFAULT FALSE,
    cascade_update BOOLEAN DEFAULT FALSE,
    validation_required BOOLEAN DEFAULT TRUE,
    metadata_schema TEXT, -- JSON schema for type-specific metadata
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Enhanced node dependencies with metadata
CREATE TABLE node_dependencies_v2 (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dependent_node_id INTEGER NOT NULL,
    dependency_node_id INTEGER NOT NULL,
    dependency_type_id INTEGER NOT NULL,
    strength INTEGER DEFAULT 50, -- 0-100, dependency strength
    priority INTEGER DEFAULT 50, -- 0-100, resolution priority
    metadata TEXT, -- JSON: type-specific metadata
    version_constraint TEXT, -- Semantic versioning constraint
    is_required BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    valid_from DATETIME DEFAULT CURRENT_TIMESTAMP,
    valid_until DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    FOREIGN KEY (dependent_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (dependency_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (dependency_type_id) REFERENCES dependency_types(id),
    UNIQUE(dependent_node_id, dependency_node_id, dependency_type_id, valid_from)
);

-- Dependency history tracking
CREATE TABLE dependency_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dependency_id INTEGER NOT NULL,
    action TEXT NOT NULL, -- 'created', 'updated', 'deleted', 'activated', 'deactivated'
    previous_state TEXT, -- JSON: previous dependency state
    new_state TEXT, -- JSON: new dependency state
    change_reason TEXT,
    changed_by TEXT,
    changed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dependency_id) REFERENCES node_dependencies_v2(id)
);

-- Dependency graph cache for performance
CREATE TABLE dependency_graph_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER NOT NULL,
    graph_data TEXT NOT NULL, -- JSON: pre-computed dependency graph
    depth INTEGER DEFAULT 0, -- Max depth in dependency tree
    total_dependencies INTEGER DEFAULT 0,
    total_dependents INTEGER DEFAULT 0,
    has_circular BOOLEAN DEFAULT FALSE,
    computed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- Dependency validation rules
CREATE TABLE dependency_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id INTEGER,
    rule_name TEXT NOT NULL,
    rule_type TEXT NOT NULL, -- 'circular_prevention', 'max_depth', 'type_compatibility'
    rule_config TEXT NOT NULL, -- JSON: rule configuration
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

-- Impact analysis results
CREATE TABLE dependency_impact_analysis (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_node_id INTEGER NOT NULL,
    impact_type TEXT NOT NULL, -- 'delete', 'update', 'version_change'
    affected_nodes TEXT NOT NULL, -- JSON: array of affected node IDs with impact details
    impact_score INTEGER, -- 0-100, overall impact severity
    analysis_metadata TEXT, -- JSON: detailed analysis results
    analyzed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_deps_v2_dependent ON node_dependencies_v2(dependent_node_id);
CREATE INDEX idx_deps_v2_dependency ON node_dependencies_v2(dependency_node_id);
CREATE INDEX idx_deps_v2_type ON node_dependencies_v2(dependency_type_id);
CREATE INDEX idx_deps_v2_active ON node_dependencies_v2(is_active);
CREATE INDEX idx_deps_v2_valid_from ON node_dependencies_v2(valid_from);
CREATE INDEX idx_deps_v2_valid_until ON node_dependencies_v2(valid_until);
CREATE INDEX idx_deps_history_dep ON dependency_history(dependency_id);
CREATE INDEX idx_deps_history_action ON dependency_history(action);
CREATE INDEX idx_deps_cache_node ON dependency_graph_cache(node_id);
CREATE INDEX idx_deps_cache_expires ON dependency_graph_cache(expires_at);
CREATE INDEX idx_deps_impact_source ON dependency_impact_analysis(source_node_id);
CREATE INDEX idx_deps_impact_type ON dependency_impact_analysis(impact_type);
```

### 3. Advanced Features

#### A. Circular Dependency Detection
```go
type DependencyGraph struct {
    nodes map[int64]*Node
    edges map[int64][]int64
}

func (g *DependencyGraph) DetectCycles() ([][]int64, error) {
    // Tarjan's strongly connected components algorithm
    // Returns cycles if found
}

func (g *DependencyGraph) ValidateNewDependency(from, to int64) error {
    // Check if adding this dependency would create a cycle
    // Use DFS with path tracking
}
```

#### B. Impact Analysis
```go
type ImpactAnalysis struct {
    SourceNode      int64
    ImpactType      string
    AffectedNodes   []AffectedNode
    ImpactScore     int
    CascadeDepth    int
    EstimatedTime   time.Duration
}

type AffectedNode struct {
    NodeID       int64
    ImpactLevel  string // 'critical', 'high', 'medium', 'low'
    Reason       string
    ActionNeeded string
}

func AnalyzeImpact(nodeID int64, action string) (*ImpactAnalysis, error) {
    // Analyze cascading effects of an action
    // Consider dependency strength, type, and priority
}
```

#### C. Dependency Resolution
```go
type DependencyResolver struct {
    graph *DependencyGraph
    rules []DependencyRule
}

func (r *DependencyResolver) ResolveDependencyOrder(nodes []int64) ([]int64, error) {
    // Topological sort with priority consideration
    // Handle circular dependencies gracefully
}

func (r *DependencyResolver) GetMinimalDependencySet(nodeID int64) ([]int64, error) {
    // Return minimal set of dependencies needed
    // Consider optional vs required dependencies
}
```

#### D. Version Constraint Management
```go
type VersionConstraint struct {
    Constraint string // e.g., ">=1.2.0 <2.0.0"
    Strategy   string // 'strict', 'compatible', 'latest'
}

func (vc *VersionConstraint) Satisfies(version string) bool {
    // Check if version satisfies constraint
    // Support semantic versioning
}
```

### 4. API Enhancements

#### New MCP Tools
```yaml
tools:
  # Enhanced dependency creation
  create_dependency_v2:
    parameters:
      dependent_id: string
      dependency_id: string
      type: string
      strength: integer (0-100)
      priority: integer (0-100)
      metadata: object
      version_constraint: string
      valid_from: datetime
      valid_until: datetime
      
  # Dependency analysis
  analyze_dependency_impact:
    parameters:
      node_id: string
      action: string # 'delete', 'update', 'version_change'
      
  # Dependency graph operations
  get_dependency_graph:
    parameters:
      node_id: string
      depth: integer
      include_metadata: boolean
      
  # Circular dependency check
  validate_dependency:
    parameters:
      dependent_id: string
      dependency_id: string
      
  # Dependency resolution
  resolve_dependencies:
    parameters:
      node_ids: array[string]
      strategy: string # 'minimal', 'complete', 'optimal'
      
  # History and versioning
  get_dependency_history:
    parameters:
      dependency_id: string
      from_date: datetime
      to_date: datetime
```

### 5. Performance Optimizations

#### A. Graph Caching
- Pre-compute dependency graphs for frequently accessed nodes
- Invalidate cache on dependency changes
- Use materialized paths for fast ancestor/descendant queries

#### B. Batch Operations
- Bulk dependency creation/update
- Batch validation for multiple dependencies
- Optimized graph traversal algorithms

#### C. Query Optimization
- Indexed lookups for common queries
- Denormalized data for read-heavy operations
- Connection pooling for concurrent access

### 6. Monitoring and Metrics

```yaml
metrics:
  dependency_health:
    - total_dependencies
    - circular_dependencies_count
    - average_dependency_depth
    - orphaned_dependencies
    
  performance:
    - graph_computation_time
    - cache_hit_rate
    - validation_time
    - impact_analysis_duration
    
  usage:
    - most_depended_nodes
    - dependency_type_distribution
    - version_constraint_violations
```

## Implementation Roadmap

### Phase 1: Core Enhancement (Week 1-2)
- Implement enhanced schema
- Add dependency type registry
- Basic circular dependency detection

### Phase 2: Advanced Features (Week 3-4)
- Impact analysis
- Version constraint management
- Dependency history tracking

### Phase 3: Performance & Tools (Week 5-6)
- Graph caching system
- MCP tool implementation
- Monitoring and metrics

### Phase 4: Integration & Testing (Week 7-8)
- Integration with existing systems
- Comprehensive testing
- Documentation and examples

## Benefits

1. **Type Safety**: Strongly typed dependency relationships
2. **Flexibility**: Extensible metadata system
3. **Performance**: Optimized graph operations
4. **Reliability**: Built-in validation and error prevention
5. **Observability**: Comprehensive tracking and analysis
6. **Scalability**: Efficient handling of large dependency graphs