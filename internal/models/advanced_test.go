package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"url-db/internal/models"
)

// Subscription model tests
func TestNodeSubscription(t *testing.T) {
	now := time.Now()
	endpoint := "https://webhook.example.com"
	
	subscription := models.NodeSubscription{
		ID:                 1,
		SubscriberService:  "test-service",
		SubscriberEndpoint: &endpoint,
		SubscribedNodeID:   100,
		EventTypes:         models.EventTypeList{"created", "updated"},
		IsActive:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	assert.Equal(t, int64(1), subscription.ID)
	assert.Equal(t, "test-service", subscription.SubscriberService)
	assert.NotNil(t, subscription.SubscriberEndpoint)
	assert.Equal(t, "https://webhook.example.com", *subscription.SubscriberEndpoint)
	assert.Equal(t, int64(100), subscription.SubscribedNodeID)
	assert.Len(t, subscription.EventTypes, 2)
	assert.True(t, subscription.IsActive)
}

func TestEventTypeList_ScanAndValue(t *testing.T) {
	// Test Value method
	eventTypes := models.EventTypeList{"created", "updated"}
	value, err := eventTypes.Value()
	assert.NoError(t, err)
	expected := `["created","updated"]`
	assert.Equal(t, expected, value)

	// Test empty list
	emptyTypes := models.EventTypeList{}
	value, err = emptyTypes.Value()
	assert.NoError(t, err)
	assert.Equal(t, "[]", value)

	// Test Scan method with string
	var scannedTypes models.EventTypeList
	err = scannedTypes.Scan(`["deleted","attribute_changed"]`)
	assert.NoError(t, err)
	assert.Len(t, scannedTypes, 2)
	assert.Contains(t, scannedTypes, "deleted")
	assert.Contains(t, scannedTypes, "attribute_changed")

	// Test Scan with []byte
	err = scannedTypes.Scan([]byte(`["created"]`))
	assert.NoError(t, err)
	assert.Len(t, scannedTypes, 1)
	assert.Equal(t, "created", scannedTypes[0])

	// Test Scan with nil
	err = scannedTypes.Scan(nil)
	assert.NoError(t, err)
	assert.Len(t, scannedTypes, 0)

	// Test Scan with other type
	err = scannedTypes.Scan(123)
	assert.NoError(t, err)
}

func TestFilterCondition_ScanAndValue(t *testing.T) {
	// Create filter condition
	filter := &models.FilterCondition{
		AttributeFilters: []models.AttributeFilter{
			{AttributeName: "category", Operator: "equals", Value: "tech"},
		},
		ChangeFilters: &models.ChangeFilter{
			Fields: []string{"title", "description"},
		},
	}

	// Test Value method
	value, err := filter.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	// Test nil Value
	var nilFilter *models.FilterCondition
	value, err = nilFilter.Value()
	assert.NoError(t, err)
	assert.Nil(t, value)

	// Test Scan method
	var scannedFilter models.FilterCondition
	jsonData := `{"attribute_filters":[{"attribute_name":"tag","operator":"contains","value":"test"}]}`
	err = scannedFilter.Scan(jsonData)
	assert.NoError(t, err)
	assert.Len(t, scannedFilter.AttributeFilters, 1)
	assert.Equal(t, "tag", scannedFilter.AttributeFilters[0].AttributeName)

	// Test Scan with []byte
	err = scannedFilter.Scan([]byte(jsonData))
	assert.NoError(t, err)

	// Test Scan with nil
	err = scannedFilter.Scan(nil)
	assert.NoError(t, err)

	// Test Scan with other type
	err = scannedFilter.Scan(123)
	assert.NoError(t, err)
}

func TestCreateNodeSubscriptionRequest(t *testing.T) {
	endpoint := "https://api.example.com/webhook"
	req := models.CreateNodeSubscriptionRequest{
		SubscriberService:  "analytics-service",
		SubscriberEndpoint: &endpoint,
		EventTypes:         []string{"created", "updated"},
		FilterConditions: &models.FilterCondition{
			AttributeFilters: []models.AttributeFilter{
				{AttributeName: "priority", Operator: "gte", Value: 5},
			},
		},
	}

	assert.Equal(t, "analytics-service", req.SubscriberService)
	assert.NotNil(t, req.SubscriberEndpoint)
	assert.Equal(t, "https://api.example.com/webhook", *req.SubscriberEndpoint)
	assert.Len(t, req.EventTypes, 2)
	assert.NotNil(t, req.FilterConditions)
	assert.Len(t, req.FilterConditions.AttributeFilters, 1)
}

func TestUpdateNodeSubscriptionRequest(t *testing.T) {
	newEndpoint := "https://new-webhook.example.com"
	isActive := false
	
	req := models.UpdateNodeSubscriptionRequest{
		SubscriberEndpoint: &newEndpoint,
		EventTypes:         []string{"deleted"},
		IsActive:           &isActive,
	}

	assert.NotNil(t, req.SubscriberEndpoint)
	assert.Equal(t, "https://new-webhook.example.com", *req.SubscriberEndpoint)
	assert.Len(t, req.EventTypes, 1)
	assert.Equal(t, "deleted", req.EventTypes[0])
	assert.NotNil(t, req.IsActive)
	assert.False(t, *req.IsActive)
}

// Dependency model tests
func TestNodeDependency(t *testing.T) {
	now := time.Now()
	metadata := &models.DependencyMetadata{
		Relationship: "parent-child",
		Description:  "Child node depends on parent",
	}
	
	dep := models.NodeDependency{
		ID:               1,
		DependentNodeID:  100,
		DependencyNodeID: 200,
		DependencyType:   models.DependencyTypeHard,
		CascadeDelete:    true,
		CascadeUpdate:    false,
		Metadata:         metadata,
		CreatedAt:        now,
	}

	assert.Equal(t, int64(1), dep.ID)
	assert.Equal(t, int64(100), dep.DependentNodeID)
	assert.Equal(t, int64(200), dep.DependencyNodeID)
	assert.Equal(t, models.DependencyTypeHard, dep.DependencyType)
	assert.True(t, dep.CascadeDelete)
	assert.False(t, dep.CascadeUpdate)
	assert.NotNil(t, dep.Metadata)
	assert.Equal(t, "parent-child", dep.Metadata.Relationship)
}

func TestNodeDependency_ParseMetadata(t *testing.T) {
	// Test with metadata
	dep := models.NodeDependency{
		Metadata: &models.DependencyMetadata{
			Relationship: "parent-child",
			Description:  "Test dependency",
		},
	}
	
	result, err := dep.ParseMetadata()
	assert.NoError(t, err)
	assert.Equal(t, "parent-child", result["relationship"])
	assert.Equal(t, "Test dependency", result["description"])

	// Test with nil metadata
	depNil := models.NodeDependency{Metadata: nil}
	result, err = depNil.ParseMetadata()
	assert.NoError(t, err)
	assert.Empty(t, result)

	// Test with empty metadata
	depEmpty := models.NodeDependency{
		Metadata: &models.DependencyMetadata{},
	}
	result, err = depEmpty.ParseMetadata()
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestDependencyMetadata_ScanAndValue(t *testing.T) {
	// Test Value method
	metadata := &models.DependencyMetadata{
		Relationship: "test-relationship",
		Description:  "test description",
	}
	
	value, err := metadata.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	// Test nil Value
	var nilMetadata *models.DependencyMetadata
	value, err = nilMetadata.Value()
	assert.NoError(t, err)
	assert.Nil(t, value)

	// Test Scan method
	var scannedMetadata models.DependencyMetadata
	jsonData := `{"relationship":"parent","description":"parent dependency"}`
	err = scannedMetadata.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, "parent", scannedMetadata.Relationship)
	assert.Equal(t, "parent dependency", scannedMetadata.Description)

	// Test Scan with []byte
	err = scannedMetadata.Scan([]byte(jsonData))
	assert.NoError(t, err)

	// Test Scan with nil
	err = scannedMetadata.Scan(nil)
	assert.NoError(t, err)

	// Test Scan with other type
	err = scannedMetadata.Scan(123)
	assert.NoError(t, err)
}

func TestDependencyConstants(t *testing.T) {
	// Test dependency type constants
	assert.Equal(t, "hard", models.DependencyTypeHard)
	assert.Equal(t, "soft", models.DependencyTypeSoft)
	assert.Equal(t, "reference", models.DependencyTypeReference)
	assert.Equal(t, "runtime", models.DependencyTypeRuntime)
	assert.Equal(t, "compile", models.DependencyTypeCompile)
	assert.Equal(t, "optional", models.DependencyTypeOptional)
	assert.Equal(t, "sync", models.DependencyTypeSync)
	assert.Equal(t, "async", models.DependencyTypeAsync)

	// Test category constants
	assert.Equal(t, "structural", models.CategoryStructural)
	assert.Equal(t, "behavioral", models.CategoryBehavioral)
	assert.Equal(t, "data", models.CategoryData)

	// Test impact level constants
	assert.Equal(t, "critical", models.ImpactLevelCritical)
	assert.Equal(t, "high", models.ImpactLevelHigh)
	assert.Equal(t, "medium", models.ImpactLevelMedium)
	assert.Equal(t, "low", models.ImpactLevelLow)

	// Test event type constants
	assert.Equal(t, "created", models.EventTypeCreated)
	assert.Equal(t, "updated", models.EventTypeUpdated)
	assert.Equal(t, "deleted", models.EventTypeDeleted)
	assert.Equal(t, "attribute_changed", models.EventTypeAttributeChanged)
	assert.Equal(t, "connection_changed", models.EventTypeConnectionChanged)
}

func TestCreateNodeDependencyRequest(t *testing.T) {
	metadata := &models.DependencyMetadata{
		Relationship: "child",
		Description:  "Child dependency",
	}
	
	req := models.CreateNodeDependencyRequest{
		DependencyNodeID: 300,
		DependencyType:   models.DependencyTypeSoft,
		CascadeDelete:    false,
		CascadeUpdate:    true,
		Metadata:         metadata,
	}

	assert.Equal(t, int64(300), req.DependencyNodeID)
	assert.Equal(t, models.DependencyTypeSoft, req.DependencyType)
	assert.False(t, req.CascadeDelete)
	assert.True(t, req.CascadeUpdate)
	assert.NotNil(t, req.Metadata)
	assert.Equal(t, "child", req.Metadata.Relationship)
}

func TestNodeEvent(t *testing.T) {
	now := time.Now()
	processedAt := now.Add(time.Minute)
	
	eventData := &models.EventData{
		EventID:   "evt-123",
		NodeID:    100,
		EventType: models.EventTypeCreated,
		Timestamp: now,
	}
	
	event := models.NodeEvent{
		ID:          1,
		NodeID:      100,
		EventType:   models.EventTypeCreated,
		EventData:   eventData,
		OccurredAt:  now,
		ProcessedAt: &processedAt,
	}

	assert.Equal(t, int64(1), event.ID)
	assert.Equal(t, int64(100), event.NodeID)
	assert.Equal(t, models.EventTypeCreated, event.EventType)
	assert.NotNil(t, event.EventData)
	assert.Equal(t, "evt-123", event.EventData.EventID)
	assert.NotNil(t, event.ProcessedAt)
}

func TestEventData_ScanAndValue(t *testing.T) {
	// Test Value method
	eventData := &models.EventData{
		EventID:   "test-event",
		NodeID:    100,
		EventType: "created",
		Timestamp: time.Now(),
	}
	
	value, err := eventData.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	// Test nil Value 
	var nilEventData *models.EventData
	value, err = nilEventData.Value()
	assert.NoError(t, err)
	assert.Nil(t, value)

	// Test Scan method
	var scannedEventData models.EventData
	jsonData := `{"event_id":"scan-test","node_id":200,"event_type":"updated"}`
	err = scannedEventData.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, "scan-test", scannedEventData.EventID)
	assert.Equal(t, int64(200), scannedEventData.NodeID)
	assert.Equal(t, "updated", scannedEventData.EventType)

	// Test Scan with []byte
	err = scannedEventData.Scan([]byte(jsonData))
	assert.NoError(t, err)

	// Test Scan with nil
	err = scannedEventData.Scan(nil)
	assert.NoError(t, err)

	// Test Scan with other type
	err = scannedEventData.Scan(123)
	assert.NoError(t, err)
}

func TestComplexDependencyStructs(t *testing.T) {
	// Test DependencyGraph
	graph := models.DependencyGraph{
		NodeID: 100,
		Dependencies: []models.DependencyNode{
			{NodeID: 200, DependencyType: "hard"},
		},
		Dependents: []models.DependencyNode{
			{NodeID: 300, DependencyType: "soft"},
		},
		Depth:       2,
		HasCircular: false,
	}
	
	assert.Equal(t, int64(100), graph.NodeID)
	assert.Len(t, graph.Dependencies, 1)
	assert.Len(t, graph.Dependents, 1)
	assert.Equal(t, 2, graph.Depth)
	assert.False(t, graph.HasCircular)

	// Test ImpactAnalysisResult
	impactResult := models.ImpactAnalysisResult{
		SourceNodeID:    100,
		SourceComposite: "url-db:test:100",
		ImpactType:      "deletion",
		AffectedNodes: []models.AffectedNode{
			{NodeID: 200, ImpactLevel: "high", Reason: "hard dependency"},
		},
		ImpactScore:     85,
		CascadeDepth:    3,
		EstimatedTime:   "5 minutes",
		Warnings:        []string{"Critical dependency detected"},
		Recommendations: []string{"Consider migration strategy"},
	}
	
	assert.Equal(t, int64(100), impactResult.SourceNodeID)
	assert.Equal(t, "url-db:test:100", impactResult.SourceComposite)
	assert.Equal(t, "deletion", impactResult.ImpactType)
	assert.Len(t, impactResult.AffectedNodes, 1)
	assert.Equal(t, 85, impactResult.ImpactScore)
	assert.Equal(t, 3, impactResult.CascadeDepth)
	assert.Len(t, impactResult.Warnings, 1)
	assert.Len(t, impactResult.Recommendations, 1)

	// Test CircularDependency
	circular := models.CircularDependency{
		Path:        []int64{100, 200, 300, 100},
		NodeDetails: []string{"Node A", "Node B", "Node C", "Node A"},
		Strength:    50,
	}
	
	assert.Len(t, circular.Path, 4)
	assert.Equal(t, int64(100), circular.Path[0])
	assert.Equal(t, int64(100), circular.Path[3])
	assert.Len(t, circular.NodeDetails, 4)
	assert.Equal(t, 50, circular.Strength)

	// Test DependencyValidationResult
	validationResult := models.DependencyValidationResult{
		IsValid:  false,
		Errors:   []string{"Circular dependency detected"},
		Warnings: []string{"Weak dependency strength"},
		Cycles:   []models.CircularDependency{circular},
	}
	
	assert.False(t, validationResult.IsValid)
	assert.Len(t, validationResult.Errors, 1)
	assert.Len(t, validationResult.Warnings, 1)
	assert.Len(t, validationResult.Cycles, 1)
}

func TestNodeDependencyV2(t *testing.T) {
	now := time.Now()
	validUntil := now.Add(24 * time.Hour)
	createdBy := "admin"
	versionConstraint := ">=1.0.0"
	
	metadata := &models.DependencyMetadataV2{
		Relationship:     "service-dependency",
		Description:      "Service depends on database",
		HealthCheckURL:   "http://db.example.com/health",
		SyncFrequency:    "5m",
		StartupOrder:     1,
		FallbackBehavior: "cache",
		CustomFields: map[string]interface{}{
			"timeout": 30,
			"retries": 3,
		},
	}
	
	depV2 := models.NodeDependencyV2{
		ID:                1,
		DependentNodeID:   100,
		DependencyNodeID:  200,
		DependencyType:    models.DependencyTypeRuntime,
		Category:          models.CategoryBehavioral,
		Strength:          85,
		Priority:          90,
		CascadeDelete:     false,
		CascadeUpdate:     true,
		Metadata:          metadata,
		VersionConstraint: &versionConstraint,
		IsRequired:        true,
		IsActive:          true,
		ValidFrom:         now,
		ValidUntil:        &validUntil,
		CreatedAt:         now,
		UpdatedAt:         now,
		CreatedBy:         &createdBy,
	}

	assert.Equal(t, int64(1), depV2.ID)
	assert.Equal(t, int64(100), depV2.DependentNodeID)
	assert.Equal(t, int64(200), depV2.DependencyNodeID)
	assert.Equal(t, models.DependencyTypeRuntime, depV2.DependencyType)
	assert.Equal(t, models.CategoryBehavioral, depV2.Category)
	assert.Equal(t, 85, depV2.Strength)
	assert.Equal(t, 90, depV2.Priority)
	assert.False(t, depV2.CascadeDelete)
	assert.True(t, depV2.CascadeUpdate)
	assert.NotNil(t, depV2.Metadata)
	assert.Equal(t, "service-dependency", depV2.Metadata.Relationship)
	assert.Equal(t, "http://db.example.com/health", depV2.Metadata.HealthCheckURL)
	assert.NotNil(t, depV2.VersionConstraint)
	assert.Equal(t, ">=1.0.0", *depV2.VersionConstraint)
	assert.True(t, depV2.IsRequired)
	assert.True(t, depV2.IsActive)
	assert.NotNil(t, depV2.ValidUntil)
	assert.NotNil(t, depV2.CreatedBy)
	assert.Equal(t, "admin", *depV2.CreatedBy)
}

func TestAdvancedDependencyStructs(t *testing.T) {
	// Test DependencyTypeConfig
	typeConfig := models.DependencyTypeConfig{
		TypeName:           "database",
		Category:           models.CategoryStructural,
		CascadeDelete:      true,
		CascadeUpdate:      false,
		ValidationRequired: true,
		DefaultStrength:    80,
		DefaultPriority:    90,
		MetadataSchema: map[string]interface{}{
			"connection_string": "string",
			"timeout":          "integer",
		},
		Description: "Database dependency configuration",
	}
	
	assert.Equal(t, "database", typeConfig.TypeName)
	assert.Equal(t, models.CategoryStructural, typeConfig.Category)
	assert.True(t, typeConfig.CascadeDelete)
	assert.False(t, typeConfig.CascadeUpdate)
	assert.True(t, typeConfig.ValidationRequired)
	assert.Equal(t, 80, typeConfig.DefaultStrength)
	assert.Equal(t, 90, typeConfig.DefaultPriority)
	assert.NotNil(t, typeConfig.MetadataSchema)
	assert.Equal(t, "string", typeConfig.MetadataSchema["connection_string"])

	// Test DependencyRule
	now := time.Now()
	domainID := int64(1)
	
	rule := models.DependencyRule{
		ID:       1,
		DomainID: &domainID,
		RuleName: "max-depth-rule",
		RuleType: "validation",
		RuleConfig: map[string]interface{}{
			"max_depth": 5,
			"allow_cycles": false,
		},
		IsActive:  true,
		CreatedAt: now,
	}
	
	assert.Equal(t, int64(1), rule.ID)
	assert.NotNil(t, rule.DomainID)
	assert.Equal(t, int64(1), *rule.DomainID)
	assert.Equal(t, "max-depth-rule", rule.RuleName)
	assert.Equal(t, "validation", rule.RuleType)
	assert.NotNil(t, rule.RuleConfig)
	assert.Equal(t, 5, rule.RuleConfig["max_depth"])
	assert.Equal(t, false, rule.RuleConfig["allow_cycles"])
	assert.True(t, rule.IsActive)

	// Test DependencyGraphCache
	expiresAt := now.Add(time.Hour)
	
	cache := models.DependencyGraphCache{
		ID:        1,
		NodeID:    100,
		GraphData: `{"dependencies":[],"dependents":[]}`,
		MaxDepth:  3,
		CreatedAt: now,
		ExpiresAt: &expiresAt,
	}
	
	assert.Equal(t, int64(1), cache.ID)
	assert.Equal(t, int64(100), cache.NodeID)
	assert.NotEmpty(t, cache.GraphData)
	assert.Equal(t, 3, cache.MaxDepth)
	assert.NotNil(t, cache.ExpiresAt)
}