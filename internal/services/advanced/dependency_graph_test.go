package services_test

import (
	"testing"
	services "url-db/internal/services/advanced"

	"github.com/stretchr/testify/assert"
)

func TestNewDependencyGraphService(t *testing.T) {
	tests := []struct {
		name         string
		expectNotNil bool
	}{
		{
			name:         "with_nil_repositories",
			expectNotNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewDependencyGraphService(nil, nil)
			if tt.expectNotNil {
				assert.NotNil(t, service)
			} else {
				assert.Nil(t, service)
			}
		})
	}
}

func TestValidateNewDependency_SelfDependency(t *testing.T) {
	tests := []struct {
		name         string
		dependentID  int64
		dependencyID int64
		expectValid  bool
		expectError  string
	}{
		{
			name:         "self_dependency_same_positive_id",
			dependentID:  1,
			dependencyID: 1,
			expectValid:  false,
			expectError:  "Self-dependencies are not allowed",
		},
		{
			name:         "self_dependency_zero_id",
			dependentID:  0,
			dependencyID: 0,
			expectValid:  false,
			expectError:  "Self-dependencies are not allowed",
		},
		{
			name:         "self_dependency_negative_id",
			dependentID:  -1,
			dependencyID: -1,
			expectValid:  false,
			expectError:  "Self-dependencies are not allowed",
		},
		{
			name:         "self_dependency_large_id",
			dependentID:  999999,
			dependencyID: 999999,
			expectValid:  false,
			expectError:  "Self-dependencies are not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test self-dependency validation which doesn't require database
			service := services.NewDependencyGraphService(nil, nil)

			result, err := service.ValidateNewDependency(tt.dependentID, tt.dependencyID)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectValid, result.IsValid)
			
			if tt.expectError != "" {
				assert.Contains(t, result.Errors, tt.expectError)
			} else {
				assert.Empty(t, result.Errors)
			}
		})
	}
}

func TestDetectCycles_NilRepository(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)

	cycles, err := service.DetectCycles(1)

	// Should error gracefully when repository is nil or buildGraph fails
	assert.Error(t, err)
	assert.Nil(t, cycles)
	assert.Contains(t, err.Error(), "failed to build graph")
}

func TestGetDependencyGraph_NilRepository(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)

	// The current implementation panics with nil repository
	// This test documents the current behavior
	assert.Panics(t, func() {
		service.GetDependencyGraph(1, 3)
	})
}

// Test basic service structure and behavior
func TestServiceBasicBehavior(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)
	
	// Test service creation
	assert.NotNil(t, service)
	
	// Test self-dependency validation (doesn't require database)
	result, err := service.ValidateNewDependency(42, 42)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "Self-dependencies are not allowed")
	
	// Test detect cycles with nil repo (should error gracefully)
	cycles, err := service.DetectCycles(1)
	assert.Error(t, err)
	assert.Nil(t, cycles)
	assert.Contains(t, err.Error(), "failed to build graph")
	
	// Test GetDependencyGraph with nil repo (currently panics)
	assert.Panics(t, func() {
		service.GetDependencyGraph(1, 3)
	})
}

// Test edge cases for ValidateNewDependency with nil repository
func TestValidateNewDependency_EdgeCases(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)
	
	tests := []struct {
		name         string
		dependentID  int64
		dependencyID int64
	}{
		{
			name:         "zero_to_positive",
			dependentID:  0,
			dependencyID: 1,
		},
		{
			name:         "negative_to_positive",
			dependentID:  -1,
			dependencyID: 1,
		},
		{
			name:         "large_ids",
			dependentID:  999999,
			dependencyID: 888888,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// With nil repository, ValidateNewDependency panics when checking existing dependencies
			assert.Panics(t, func() {
				service.ValidateNewDependency(tt.dependentID, tt.dependencyID)
			})
		})
	}
}

// Test DetectCycles edge cases
func TestDetectCycles_EdgeCases(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)
	
	testCases := []int64{0, 1, -1, 999999}
	
	for _, domainID := range testCases {
		t.Run("domain_"+string(rune(domainID+48)), func(t *testing.T) {
			cycles, err := service.DetectCycles(domainID)
			
			// Should always error with nil repository
			assert.Error(t, err)
			assert.Nil(t, cycles)
			assert.Contains(t, err.Error(), "failed to build graph")
		})
	}
}

// Test GetDependencyGraph edge cases
func TestGetDependencyGraph_EdgeCases(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)
	
	tests := []struct {
		name     string
		nodeID   int64
		maxDepth int
	}{
		{"zero_node_zero_depth", 0, 0},
		{"positive_node_zero_depth", 1, 0},
		{"negative_node_positive_depth", -1, 3},
		{"large_node_large_depth", 999999, 100},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Current implementation panics with nil repository
			assert.Panics(t, func() {
				service.GetDependencyGraph(tt.nodeID, tt.maxDepth)
			})
		})
	}
}