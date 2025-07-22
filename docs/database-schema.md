# URL Database Schema Documentation

## Overview

This document provides comprehensive documentation for the URL Database SQLite schema. The schema is designed for immutable, versioned data structures that support URL management, attribute tagging, and external dependency tracking.

**Schema Version**: 1.0  
**Database Engine**: SQLite 3  
**Character Set**: UTF-8  

## Core Design Principles

1. **Immutable Schema**: This schema represents the production-ready, stable design
2. **Referential Integrity**: All foreign key relationships with proper CASCADE behaviors
3. **Performance Optimized**: Strategic indexing for common query patterns
4. **Extensible**: Support for future features through JSON metadata fields
5. **Audit Trail**: Automatic timestamp tracking with triggers

## Table Structure

### Core Tables

#### 1. `domains` - Domain Management
```sql
CREATE TABLE domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Purpose**: Top-level organization for URL collections  
**Constraints**:
- `name`: Unique across all domains, used in composite keys
- Primary key auto-increment for internal references

**Business Rules**:
- Domain names must be unique
- Deletion cascades to all child nodes and attributes
- Used in MCP composite key format: `url-db:{domain_name}:{node_id}`

#### 2. `nodes` - URL/Content Storage
```sql
CREATE TABLE nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    domain_id INTEGER NOT NULL,
    title TEXT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
    UNIQUE(content, domain_id)
);
```

**Purpose**: Store URLs and associated metadata  
**Constraints**:
- `content`: The actual URL or content string
- `domain_id`: Must reference valid domain
- Unique constraint prevents duplicate URLs within same domain

**Business Rules**:
- URLs must be unique within each domain
- Support for POST-method URL storage (long URLs)
- Title and description are optional metadata

#### 3. `attributes` - Attribute Schema Definition
```sql
CREATE TABLE attributes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('tag', 'ordered_tag', 'number', 'string', 'markdown', 'image')),
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
    UNIQUE(domain_id, name)
);
```

**Purpose**: Define attribute schemas per domain  
**Supported Types**:
- `tag`: Simple string tags
- `ordered_tag`: Tags with ordering significance
- `number`: Numeric values
- `string`: Text values
- `markdown`: Markdown-formatted text
- `image`: Image URLs or references

**Business Rules**:
- Attribute names unique within each domain
- Type enforcement through CHECK constraint
- Domain-specific attribute schemas

#### 4. `node_attributes` - Attribute Value Assignment
```sql
CREATE TABLE node_attributes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER NOT NULL,
    attribute_id INTEGER NOT NULL,
    value TEXT NOT NULL,
    order_index INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE
);
```

**Purpose**: Store actual attribute values for nodes  
**Features**:
- `order_index`: For ordered_tag types to maintain sequence
- Value stored as TEXT with type validation at application layer
- Many-to-many relationship between nodes and attributes

#### 5. `node_connections` - Inter-Node Relationships
```sql
CREATE TABLE node_connections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_node_id INTEGER NOT NULL,
    target_node_id INTEGER NOT NULL,
    relationship_type TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (target_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    UNIQUE(source_node_id, target_node_id, relationship_type)
);
```

**Purpose**: Model relationships between nodes  
**Relationship Types**:
- `parent` / `child`: Hierarchical relationships
- `related`: General associations
- `next` / `prev`: Sequential relationships
- Custom types as needed

### External Dependency Management

#### 6. `node_subscriptions` - Event Subscription System
```sql
CREATE TABLE node_subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subscriber_service TEXT NOT NULL,
    subscriber_endpoint TEXT,
    subscribed_node_id INTEGER NOT NULL,
    event_types TEXT NOT NULL,
    filter_conditions TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (subscribed_node_id) REFERENCES nodes(id) ON DELETE CASCADE
);
```

**Purpose**: Track external service subscriptions to node events  
**Features**:
- `event_types`: JSON array of event types to monitor
- `filter_conditions`: JSON object for subscription filtering
- `subscriber_endpoint`: Callback URL for notifications

#### 7. `node_dependencies` - Dependency Tracking
```sql
CREATE TABLE node_dependencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dependent_node_id INTEGER NOT NULL,
    dependency_node_id INTEGER NOT NULL,
    dependency_type TEXT NOT NULL,
    cascade_delete BOOLEAN DEFAULT FALSE,
    cascade_update BOOLEAN DEFAULT FALSE,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dependent_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (dependency_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    UNIQUE(dependent_node_id, dependency_node_id)
);
```

**Purpose**: Model dependencies between nodes  
**Dependency Types**:
- `hard`: Strong dependency, affects lifecycle
- `soft`: Loose dependency, informational
- `reference`: Reference-only relationship

#### 8. `node_events` - Event Audit Log
```sql
CREATE TABLE node_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER NOT NULL,
    event_type TEXT NOT NULL,
    event_data TEXT,
    occurred_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);
```

**Purpose**: Audit trail for all node lifecycle events  
**Event Types**:
- `created`: Node creation
- `updated`: Node modification
- `deleted`: Node deletion
- `attribute_changed`: Attribute modifications

## Indexes and Performance

### Primary Indexes
```sql
CREATE INDEX idx_nodes_domain ON nodes(domain_id);
CREATE INDEX idx_nodes_content ON nodes(content);
CREATE INDEX idx_attributes_domain ON attributes(domain_id);
CREATE INDEX idx_node_attributes_node ON node_attributes(node_id);
CREATE INDEX idx_node_attributes_attribute ON node_attributes(attribute_id);
CREATE INDEX idx_node_connections_source ON node_connections(source_node_id);
CREATE INDEX idx_node_connections_target ON node_connections(target_node_id);
```

### External Dependency Indexes
```sql
CREATE INDEX idx_subscriptions_service ON node_subscriptions(subscriber_service);
CREATE INDEX idx_subscriptions_node ON node_subscriptions(subscribed_node_id);
CREATE INDEX idx_subscriptions_active ON node_subscriptions(is_active);
CREATE INDEX idx_dependencies_dependent ON node_dependencies(dependent_node_id);
CREATE INDEX idx_dependencies_dependency ON node_dependencies(dependency_node_id);
CREATE INDEX idx_events_node ON node_events(node_id);
CREATE INDEX idx_events_type ON node_events(event_type);
CREATE INDEX idx_events_occurred ON node_events(occurred_at);
CREATE INDEX idx_events_unprocessed ON node_events(processed_at) WHERE processed_at IS NULL;
```

**Performance Characteristics**:
- Fast domain-based queries through `idx_nodes_domain`
- Efficient URL lookups via `idx_nodes_content`
- Quick attribute value retrieval with composite indexes
- Optimized event processing with filtered index on unprocessed events

## Triggers and Automation

### Automatic Timestamp Updates
```sql
CREATE TRIGGER domains_updated_at 
    AFTER UPDATE ON domains 
    FOR EACH ROW 
    BEGIN 
        UPDATE domains SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER nodes_updated_at 
    AFTER UPDATE ON nodes 
    FOR EACH ROW 
    BEGIN 
        UPDATE nodes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER node_subscriptions_updated_at 
    AFTER UPDATE ON node_subscriptions 
    FOR EACH ROW 
    BEGIN 
        UPDATE node_subscriptions SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;
```

**Purpose**: Maintain accurate modification timestamps without application logic

## Data Integrity Rules

### Foreign Key Constraints
- All child tables use `ON DELETE CASCADE` for proper cleanup
- Referential integrity enforced at database level
- No orphaned records possible

### Check Constraints
- Attribute types limited to supported values
- Prevents invalid type assignments

### Unique Constraints
- Domain names globally unique
- URLs unique within domains
- Attribute names unique within domains
- Node relationships prevent duplicates

## MCP Integration Patterns

### Composite Key Format
- Pattern: `url-db:{domain_name}:{node_id}`
- External systems use composite keys for all node references
- Internal IDs hidden from external APIs

### Resource URI Patterns
- `mcp://server/info` - Server metadata
- `mcp://domains/{domain_name}` - Domain information
- `mcp://domains/{domain_name}/nodes` - Domain node listings
- `mcp://nodes/{composite_id}` - Individual node resources

## Migration and Versioning

### Schema Version Control
- Schema changes require version increment
- Migration scripts for schema updates
- Backward compatibility considerations

### Data Migration
- Export/import procedures for data preservation
- Schema validation before migrations
- Rollback procedures for failed migrations

## Security Considerations

### Data Protection
- No sensitive data stored in plaintext
- Foreign key constraints prevent data corruption
- Audit trail for all modifications

### Access Control
- Application-level access control implementation
- Database-level constraints for data integrity
- No direct database access for external systems

## Backup and Recovery

### Backup Strategy
- Regular SQLite database file backups
- Point-in-time recovery capabilities
- Schema and data backup separation

### Recovery Procedures
- Database file restoration
- Transaction log replay
- Data consistency verification

## Performance Tuning

### Query Optimization
- Strategic indexing for common patterns
- Compound indexes for multi-column queries
- Partial indexes for filtered queries

### Maintenance Tasks
- Regular VACUUM operations
- Index statistics updates
- Query plan analysis

## Future Extensions

### Planned Features
- Full-text search capabilities
- JSON attribute value validation
- Advanced relationship types
- Temporal data tracking

### Extension Points
- Custom attribute types
- Additional relationship types
- Extended metadata fields
- External system integrations

---

**Important**: This schema is considered immutable for production use. Any modifications require careful planning, testing, and migration procedures to ensure data integrity and system stability.