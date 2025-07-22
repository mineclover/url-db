# MCP LLM Judge Scenarios - External Dependency API Testing

> **Test Status**: 🚧 READY FOR TESTING (2025-07-22)  
> **API Coverage**: External Dependency API (Subscriptions, Dependencies, Events)  
> **Test Suite**: Comprehensive validation scenarios for new external dependency functionality

## Overview

This document provides comprehensive testing scenarios for evaluating the newly implemented external dependency API through LLM-based testing. These scenarios validate the subscription management, dependency relationships, and event tracking systems integrated into the URL-DB MCP server.

## New API Components

The external dependency API adds three main components to URL-DB:

### 1. **Subscription Management**
- Event-based notifications when nodes change
- Webhook endpoints for external service integration
- Filtered event subscriptions by event type and criteria

### 2. **Dependency Management** 
- Directed relationships between nodes (A depends on B)
- Cascade delete and update behaviors
- Support for different dependency types (hard, soft, reference)

### 3. **Event Tracking**
- Automatic event generation for all node CRUD operations
- Event processing and status tracking
- Historical audit trail with before/after state capture

## Test Environment Setup

```bash
# Start MCP server with external dependency support
./bin/url-db -mcp-mode=stdio -db-path=test_external_deps.db

# Or HTTP mode for REST API verification
./bin/url-db -mcp-mode=sse -port=8080 -db-path=test_external_deps.db
```

## New MCP Tools Available

**Subscription Management**:
- `create_subscription` - Subscribe to node events
- `list_subscriptions` - List subscriptions by service
- `get_node_subscriptions` - Get subscriptions for specific node
- `delete_subscription` - Cancel subscription

**Dependency Management**:
- `create_dependency` - Create dependency relationship
- `list_node_dependencies` - Get node's dependencies
- `list_node_dependents` - Get nodes that depend on this node
- `delete_dependency` - Remove dependency relationship

**Event Management**:
- `get_node_events` - Get event history for node
- `get_pending_events` - Get unprocessed events
- `process_event` - Mark event as processed
- `get_event_stats` - Get system event statistics

## Scenario Categories

### 1. Basic External Dependency Functionality

#### Scenario 1.1: Subscription Lifecycle Management
**Objective**: Test complete subscription creation, management, and deletion workflow
**LLM Instructions**: Create subscriptions, verify filtering, and test cleanup

```
Test Steps:
1. Create domain "microservices" and add nodes: api-gateway, user-service, payment-service
2. Create subscription to monitor "api-gateway" for all event types
   - Service: "monitoring-system"
   - Events: ["created", "updated", "deleted"]
3. Create filtered subscription for "user-service" monitoring only "updated" events
   - Service: "alert-system" 
   - Events: ["updated"]
   - Include webhook endpoint: "https://alerts.example.com/webhook"
4. List all subscriptions and verify both were created
5. Get node-specific subscriptions for "api-gateway"
6. Delete the first subscription and verify removal
```

**Expected MCP Tools Used**:
- `create_domain`, `create_node`
- `create_subscription`
- `list_subscriptions`
- `get_node_subscriptions`
- `delete_subscription`

**Success Criteria**:
- Subscriptions created with correct event type filters
- Service-based subscription listing works
- Node-specific subscription queries return accurate results
- Subscription deletion removes records without affecting others

#### Scenario 1.2: Dependency Relationship Creation and Querying
**Objective**: Test creation and management of various dependency types
**LLM Instructions**: Create complex dependency relationships and verify bidirectional queries

```
Test Steps:
1. Use nodes from previous scenario: api-gateway, user-service, payment-service
2. Create hard dependency: api-gateway depends on user-service
   - Type: "hard"
   - Cascade delete: true
   - Cascade update: true
3. Create soft dependency: user-service depends on payment-service
   - Type: "soft"
   - Cascade delete: false
   - Cascade update: true
4. Create reference dependency: api-gateway references payment-service
   - Type: "reference"
   - Include metadata: {"relationship": "external API", "description": "Payment processing endpoint"}
5. List dependencies for api-gateway (should show 2)
6. List dependents for user-service (should show api-gateway)
7. List dependents for payment-service (should show user-service and api-gateway)
```

**Expected MCP Tools Used**:
- `create_dependency`
- `list_node_dependencies`
- `list_node_dependents`

**Success Criteria**:
- All dependency types created successfully
- Cascade options stored correctly
- Metadata preserved in reference dependency
- Bidirectional dependency queries work correctly

#### Scenario 1.3: Automatic Event Generation and Tracking
**Objective**: Verify automatic event creation during node operations
**LLM Instructions**: Perform node operations and verify events are automatically captured

```
Test Steps:
1. Create new node "database-service" in "microservices" domain
2. Immediately check for "node.created" event using get_node_events
3. Update the node title to "Primary Database Service"
4. Check for "node.updated" event with before/after data
5. Update description to "Main PostgreSQL database instance"
6. Verify second "node.updated" event was created
7. Get all pending events to see unprocessed events
8. Process one event and verify it's marked as processed
9. Delete the node and check for "node.deleted" event
```

**Expected MCP Tools Used**:
- `create_node`
- `update_node`
- `delete_node`
- `get_node_events`
- `get_pending_events`
- `process_event`

**Success Criteria**:
- Events automatically generated for create, update, delete operations
- Event data contains proper before/after state information
- Events appear in both node-specific and pending event lists
- Event processing correctly marks events as processed
- Deleted node events still accessible in history

### 2. Advanced Integration Scenarios

#### Scenario 2.1: Complex Dependency Chains with Cascade Behavior
**Objective**: Test multi-level dependencies and cascade operations
**LLM Instructions**: Create dependency chains and test cascade behavior

```
Test Steps:
1. Create domain "infrastructure" with nodes: load-balancer, web-server, app-server, database
2. Create dependency chain with cascade delete enabled:
   - load-balancer depends on web-server (cascade delete: true)
   - web-server depends on app-server (cascade delete: true)  
   - app-server depends on database (cascade delete: false)
3. Create additional cross-dependencies:
   - load-balancer depends on database (type: "reference", cascade delete: false)
4. Verify dependency structure:
   - List dependencies for each node
   - List dependents for each node
   - Verify cascade settings are correct
5. Test dependency integrity:
   - Attempt to delete database (should succeed as no cascade delete dependencies)
   - Verify dependent services remain intact
   - Check events generated during dependency operations
```

**Expected MCP Tools Used**:
- `create_dependency` (multiple calls)
- `list_node_dependencies`
- `list_node_dependents`
- `delete_node`
- `get_node_events`

**Success Criteria**:
- Complex dependency chain created successfully
- Cascade settings properly configured and enforced
- Cross-dependencies handled correctly
- Dependency integrity maintained during operations
- Appropriate events generated for all dependency changes

#### Scenario 2.2: Event Processing and System Statistics
**Objective**: Test event processing workflows and system monitoring
**LLM Instructions**: Generate events, process them systematically, and analyze statistics

```
Test Steps:
1. Perform multiple node operations to generate events:
   - Create 5 new nodes in different domains
   - Update each node 2 times (title and description changes)
   - Delete 2 nodes
2. Get comprehensive system event statistics
3. Analyze statistics for:
   - Total event counts by type
   - Processing status distribution
   - Event generation rate
4. Get list of all pending events (should be substantial)
5. Process events in batches:
   - Process first 5 events individually
   - Verify each is marked as processed
6. Get updated statistics and verify processed counts increased
7. Get remaining pending events and verify count decreased
```

**Expected MCP Tools Used**:
- `create_node`, `update_node`, `delete_node` (multiple calls)
- `get_event_stats`
- `get_pending_events`
- `process_event` (multiple calls)

**Success Criteria**:
- Event statistics provide accurate counts and breakdowns
- Pending event lists reflect current processing status
- Event processing updates statistics in real-time
- Processing timestamps recorded correctly
- System handles batch processing efficiently

### 3. Integration and Workflow Tests

#### Scenario 3.1: Complete Project Management Workflow
**Objective**: Simulate realistic project management using all external dependency features
**LLM Instructions**: Execute end-to-end project lifecycle with monitoring and dependencies

```
Workflow Steps:
1. **Project Setup**:
   - Create domain "project-alpha" 
   - Add nodes: frontend, backend, database, docs, config
   
2. **Dependency Architecture**:
   - frontend depends on backend (hard, cascade update: true)
   - backend depends on database (hard, cascade delete: false)
   - frontend references docs (reference type)
   - All services reference config (reference type)
   
3. **Monitoring Setup**:
   - Create subscription for "devops-team" monitoring all components
   - Create subscription for "qa-team" monitoring only update events
   - Create subscription for "docs-team" monitoring docs changes only
   
4. **Development Simulation**:
   - Update backend API (should trigger cascade update event to frontend)
   - Update database schema (should generate dependency-related events)
   - Update documentation (should notify docs-team subscription)
   - Add new configuration (should affect all services via reference dependencies)
   
5. **Audit and Reporting**:
   - Generate comprehensive event history for entire project
   - Show subscription activity and notifications
   - Document dependency relationships and their cascade behaviors
   - Provide system statistics for the workflow
```

**Integration Verification**:
- All three external dependency systems work together seamlessly
- Events capture complete audit trail of project changes
- Subscriptions provide appropriate notifications for different teams
- Dependencies maintain referential integrity throughout workflow
- Cascade behaviors work correctly across dependency types

#### Scenario 3.2: Cross-Domain Resource Management
**Objective**: Test external dependencies across multiple organizational domains
**LLM Instructions**: Manage resources and dependencies spanning different domains

```
Workflow Steps:
1. **Multi-Domain Setup**:
   - Create domain "applications" with nodes: web-app, mobile-app
   - Create domain "infrastructure" with nodes: auth-service, cdn, monitoring
   - Create domain "data" with nodes: user-db, analytics-db, cache
   
2. **Cross-Domain Dependencies**:
   - web-app depends on auth-service (cross-domain hard dependency)
   - mobile-app depends on auth-service (cross-domain hard dependency)
   - Both apps depend on cdn (cross-domain soft dependency)
   - auth-service depends on user-db (cross-domain hard dependency)
   - monitoring references all apps and infrastructure (cross-domain references)
   
3. **Distributed Monitoring**:
   - Create subscriptions for "app-team" monitoring applications domain
   - Create subscriptions for "infra-team" monitoring infrastructure domain
   - Create subscriptions for "data-team" monitoring data domain
   - Create cross-domain subscription for "architect-team" monitoring all domains
   
4. **Cross-Domain Operations**:
   - Update auth-service (should affect multiple dependent apps)
   - Scale cdn configuration (soft dependency, should notify apps)
   - Database maintenance on user-db (hard dependency, should cascade to auth-service)
   - Monitor event propagation across domain boundaries
   
5. **System Analysis**:
   - Verify cross-domain dependency resolution
   - Check event tracking across all domains
   - Validate subscription notifications for relevant teams
   - Analyze cascade behavior across domain boundaries
```

**Cross-Domain Verification**:
- Dependencies work correctly across different domains
- Event tracking maintains domain context while enabling cross-domain analysis
- Subscription filtering respects domain boundaries when configured
- Cascade operations work across domains with proper safeguards

### 4. Error Handling and Edge Cases

#### Scenario 4.1: Invalid Input Handling
**Objective**: Test robust error handling for malformed requests
**LLM Instructions**: Attempt various invalid operations and verify error responses

```
Error Test Cases:
1. **Invalid Composite IDs**:
   - Create subscription with malformed composite ID "invalid-format"
   - Create dependency with non-existent node "url-db:fake:999"
   - Get events for empty composite ID ""
   
2. **Invalid Parameters**:
   - Create subscription with empty event types array
   - Create dependency with invalid dependency type "invalid-type"
   - Process non-existent event ID 99999
   
3. **Circular Dependencies**:
   - Create dependency A → B
   - Create dependency B → C
   - Attempt to create dependency C → A (should fail)
   - Attempt direct self-dependency A → A (should fail)
   
4. **Resource Conflicts**:
   - Delete node that has active subscriptions
   - Delete node that has dependencies with cascade delete enabled
   - Create duplicate subscriptions for same service and node
```

**Error Validation**:
- All invalid requests return appropriate error codes
- Error messages are informative and actionable
- System maintains stability under error conditions
- No data corruption occurs from failed operations

#### Scenario 4.2: Concurrency and Race Conditions
**Objective**: Test system behavior under concurrent operations
**LLM Instructions**: Simulate concurrent access and verify data consistency

```
Concurrency Test Cases:
1. **Simultaneous Node Operations**:
   - Create multiple subscriptions for same node simultaneously
   - Update same node from multiple operations concurrently
   - Process same events from multiple clients
   
2. **Dependency Race Conditions**:
   - Create competing dependencies simultaneously
   - Delete node while creating dependencies for it
   - Update node while processing its events
   
3. **Event System Stress**:
   - Generate high volume of events rapidly
   - Process events while new ones are being created
   - Verify event ordering and consistency
```

**Concurrency Validation**:
- Data consistency maintained under concurrent access
- Event ordering preserved correctly
- No race conditions cause data corruption
- System performance remains stable under load

### 5. Performance and Scalability

#### Scenario 5.1: High-Volume Event Processing
**Objective**: Test system performance with large numbers of events
**LLM Instructions**: Generate substantial event volumes and measure performance

```
Performance Test Steps:
1. **Event Generation Load Test**:
   - Create 100 nodes rapidly across multiple domains
   - Perform bulk updates (200+ operations)
   - Create 50+ dependency relationships
   - Create 20+ subscriptions for various services
   
2. **Event Processing Performance**:
   - Measure event generation latency per operation
   - Process events in large batches (50+ events)
   - Monitor memory and CPU usage during processing
   - Test event querying performance with large datasets
   
3. **Subscription System Scale**:
   - Create subscriptions for many services (100+ subscriptions)
   - Test subscription listing performance
   - Verify event filtering performance with complex criteria
   
4. **Dependency Network Performance**:
   - Create complex dependency networks (500+ relationships)
   - Test dependency traversal performance
   - Measure response times for dependency queries
   - Test cascade operation performance
```

**Performance Metrics**:
- Event generation adds <50ms latency to node operations
- Event processing throughput >100 events/second
- Dependency queries complete within 100ms for networks <1000 nodes
- Memory usage remains stable under sustained load
- Database performance scales appropriately with data volume

### 6. Integration Verification

#### Scenario 6.1: MCP Tool Consistency and Completeness
**Objective**: Verify all external dependency MCP tools work correctly
**LLM Instructions**: Test every new MCP tool systematically

```
Tool Verification Steps:
1. **Subscription Tools**:
   - create_subscription: Test all parameter combinations
   - list_subscriptions: Verify pagination and filtering
   - get_node_subscriptions: Test with various node types
   - delete_subscription: Verify cleanup and error handling
   
2. **Dependency Tools**:
   - create_dependency: Test all dependency types and cascade options
   - list_node_dependencies: Verify complete dependency resolution
   - list_node_dependents: Test bidirectional relationship queries
   - delete_dependency: Verify safe dependency removal
   
3. **Event Tools**:
   - get_node_events: Test event filtering and pagination
   - get_pending_events: Verify processing status accuracy
   - process_event: Test batch and individual processing
   - get_event_stats: Verify statistical accuracy and completeness
```

**Tool Verification Criteria**:
- All tools respond with correct data structures
- Parameter validation works correctly
- Error responses follow MCP standards
- Tool documentation matches actual behavior
- Composite ID handling consistent across all tools

## Success Criteria and Evaluation Framework

### Functional Requirements ✅
- **Event Integration**: Events automatically generated for all node CRUD operations
- **Subscription Management**: Complete lifecycle management with filtering support
- **Dependency Relationships**: All dependency types work with proper cascade behavior
- **Cross-Domain Support**: External dependencies work across different domains
- **Data Consistency**: Referential integrity maintained throughout all operations

### Performance Requirements ✅
- **Low Latency**: Event generation adds minimal overhead to node operations
- **High Throughput**: Event processing handles substantial event volumes
- **Scalable Queries**: Dependency and event queries perform well at scale
- **Memory Efficiency**: System remains stable under sustained load

### Integration Requirements ✅
- **MCP Compliance**: All new tools follow MCP protocol standards
- **REST API Parity**: External dependency features available via both MCP and REST
- **Composite ID Consistency**: Uniform composite ID handling across all tools
- **Event System Integration**: Seamless integration with existing node operations

### Error Handling Requirements ✅
- **Graceful Degradation**: System handles errors without data corruption
- **Informative Messages**: Error responses provide actionable information
- **Input Validation**: Robust validation prevents invalid state creation
- **Consistency**: Error behavior consistent across all components

## LLM Judge Configuration

### Evaluation Criteria
```yaml
evaluation_criteria:
  functional_correctness: 30%    # Do the tools work as specified?
  integration_quality: 25%      # How well do components work together?
  error_handling: 20%           # How robust is error handling?
  performance: 15%              # Does the system perform adequately?
  consistency: 10%              # Is behavior consistent across tools?
```

### Scoring Rubric
- **10/10**: Exceptional - Exceeds all requirements
- **8-9/10**: Excellent - Meets all requirements with minor issues
- **6-7/10**: Good - Meets most requirements with some gaps
- **4-5/10**: Fair - Basic functionality works but has significant issues
- **1-3/10**: Poor - Major functionality broken or missing
- **0/10**: Failing - Does not work at all

### Expected Output Format
```json
{
  "test_suite": "external_dependency_api",
  "execution_date": "2025-07-22",
  "scenarios_executed": 18,
  "scenarios_passed": 17,
  "scenarios_failed": 1,
  "overall_score": 94.4,
  "detailed_results": [
    {
      "scenario_id": "1.1",
      "name": "Subscription Lifecycle Management",
      "status": "PASSED",
      "score": 9.2,
      "execution_time": 3.4,
      "tools_tested": ["create_subscription", "list_subscriptions", "delete_subscription"],
      "issues": [],
      "performance_metrics": {
        "avg_response_time": 45,
        "max_response_time": 120
      }
    }
  ],
  "recommendations": [
    "Optimize event query performance for large datasets",
    "Add batch processing capabilities for dependency creation",
    "Improve error message specificity for validation failures"
  ]
}
```

## Conclusion

This comprehensive testing framework ensures the external dependency API meets production standards and integrates seamlessly with the existing URL-DB system. The scenarios cover:

- **Complete Feature Coverage**: All subscription, dependency, and event functionality
- **Realistic Workflows**: End-to-end scenarios that mirror actual usage patterns  
- **Robust Error Handling**: Edge cases and error conditions thoroughly tested
- **Performance Validation**: System performance under various load conditions
- **Integration Quality**: Seamless integration with existing MCP and REST APIs

The LLM judge approach provides thorough evaluation while simulating real-world usage patterns, ensuring the external dependency API is ready for production deployment with confidence in its reliability, performance, and usability.