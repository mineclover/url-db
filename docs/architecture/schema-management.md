# Schema Management - Single Source of Truth

## Overview

The database schema is now managed through a **single source of truth** approach to eliminate inconsistencies between schema definitions and ensure easier maintenance.

## Architecture

### Schema File
- **Location**: `/schema.sql`
- **Purpose**: Single authoritative source for all database structures
- **Content**: Complete schema including tables, indexes, triggers, and initial data

### Database Initialization
- **File**: `/internal/database/database.go`
- **Method**: `createSchema()` loads schema from external `schema.sql` file
- **Fallback**: Maintains inline schema for backward compatibility if file loading fails

### Key Features

1. **Automatic Detection**: Project root is detected by searching for `go.mod` file
2. **Fallback Strategy**: If `schema.sql` cannot be loaded, falls back to inline schema
3. **Error Handling**: Clear error messages for schema loading issues
4. **Compatibility**: Maintains existing API and behavior

## Schema Contents

The unified schema includes:

### Core Tables
- `domains` - Domain/folder organization
- `nodes` - URL storage with domain association
- `attributes` - Attribute type definitions per domain
- `node_attributes` - Attribute values for nodes

### Advanced Dependency System
- `dependency_types` - Type registry (8 built-in types)
- `node_dependencies_v2` - Enhanced dependency management
- `dependency_history` - Change tracking
- `dependency_graph_cache` - Performance optimization
- `dependency_rules` - Validation rules
- `dependency_impact_analysis` - Impact analysis results

### Event System
- `node_events` - Event logging
- `node_subscriptions` - External service subscriptions

### Performance Optimization
- 25+ indexes for optimal query performance
- Composite indexes for complex queries
- Conditional indexes for filtered operations

## Benefits

1. **Consistency**: Single schema definition eliminates version drift
2. **Maintainability**: Schema changes made in one place
3. **Reliability**: Fallback ensures system continues working
4. **Performance**: Optimized indexes and constraints
5. **Extensibility**: Easy to add new tables and features

## Usage

The schema is automatically loaded when the database is initialized:

```go
db, err := database.New(config)
// Schema is loaded from schema.sql automatically
```

## Advanced Features

### Dependency Types (8 built-in)
- **Structural**: hard, soft, reference
- **Behavioral**: runtime, compile, optional  
- **Data**: sync, async

### Performance Features
- Graph caching for dependency traversal
- Impact analysis with scoring
- History tracking for auditing
- Rule-based validation

## Migration Path

The system maintains backward compatibility while providing a path forward for enhanced dependency management features that are implemented in the advanced services layer.