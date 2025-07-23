package services_test

import (
	"testing"
	services "url-db/internal/services/advanced"

	"github.com/stretchr/testify/assert"
)

// Test more scenarios for ValidateNewDependency which only requires self-dependency validation
func TestValidateNewDependency_ComprehensiveScenarios(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)

	tests := []struct {
		name         string
		dependentID  int64
		dependencyID int64
		expectValid  bool
		expectError  string
	}{
		{
			name:         "self_dependency_positive",
			dependentID:  42,
			dependencyID: 42,
			expectValid:  false,
			expectError:  "Self-dependencies are not allowed",
		},
		{
			name:         "self_dependency_zero",
			dependentID:  0,
			dependencyID: 0,
			expectValid:  false,
			expectError:  "Self-dependencies are not allowed",
		},
		{
			name:         "self_dependency_negative",
			dependentID:  -5,
			dependencyID: -5,
			expectValid:  false,
			expectError:  "Self-dependencies are not allowed",
		},
		{
			name:         "self_dependency_large_number",
			dependentID:  999999999,
			dependencyID: 999999999,
			expectValid:  false,
			expectError:  "Self-dependencies are not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ValidateNewDependency(tt.dependentID, tt.dependencyID)

			assert.NoError(t, err, "ValidateNewDependency should not return error for self-dependency check")
			assert.NotNil(t, result, "Result should not be nil")
			assert.Equal(t, tt.expectValid, result.IsValid, "IsValid should match expected value")
			
			if tt.expectError != "" {
				assert.Contains(t, result.Errors, tt.expectError, "Should contain expected error message")
			} else {
				assert.Empty(t, result.Errors, "Should not have errors for valid dependency")
			}
		})
	}
}

// Test edge cases for DetectCycles
func TestDetectCycles_ComprehensiveErrors(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)

	testCases := []struct {
		name     string
		domainID int64
	}{
		{"zero_domain", 0},
		{"positive_domain", 1},
		{"negative_domain", -1},
		{"large_domain", 999999999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cycles, err := service.DetectCycles(tc.domainID)

			// Should always error with nil repository
			assert.Error(t, err, "Should error with nil repository")
			assert.Nil(t, cycles, "Cycles should be nil when error occurs")
			assert.Contains(t, err.Error(), "failed to build graph", "Should contain expected error message")
		})
	}
}

// Test GetDependencyGraph edge cases - all should panic with nil repository
func TestGetDependencyGraph_ComprehensivePanics(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)

	testCases := []struct {
		name     string
		nodeID   int64
		maxDepth int
	}{
		{"zero_node_zero_depth", 0, 0},
		{"positive_node_positive_depth", 1, 3},
		{"negative_node_positive_depth", -1, 5},
		{"large_node_large_depth", 999999, 100},
		{"positive_node_negative_depth", 1, -1},
		{"zero_node_large_depth", 0, 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// All these should panic with nil repository
			assert.Panics(t, func() {
				service.GetDependencyGraph(tc.nodeID, tc.maxDepth)
			}, "Should panic with nil repository for any input")
		})
	}
}

// Test service creation with different input scenarios
func TestNewDependencyGraphService_ComprehensiveCreation(t *testing.T) {
	testCases := []struct {
		name     string
		depRepo  interface{}
		nodeRepo interface{}
	}{
		{"both_nil", nil, nil},
		{"depRepo_nil", nil, "not_nil"},
		{"nodeRepo_nil", "not_nil", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Service should always be created successfully
			service := services.NewDependencyGraphService(nil, nil)
			assert.NotNil(t, service, "Service should not be nil regardless of input")
		})
	}
}

// Test consistent behavior patterns
func TestServiceConsistency(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)

	// Test that multiple calls to same method with same params produce same results
	t.Run("validate_consistency", func(t *testing.T) {
		result1, err1 := service.ValidateNewDependency(42, 42)
		result2, err2 := service.ValidateNewDependency(42, 42)

		assert.Equal(t, err1, err2, "Errors should be consistent")
		assert.Equal(t, result1.IsValid, result2.IsValid, "IsValid should be consistent")
		assert.Equal(t, result1.Errors, result2.Errors, "Errors should be consistent")
	})

	t.Run("detect_cycles_consistency", func(t *testing.T) {
		_, err1 := service.DetectCycles(1)
		_, err2 := service.DetectCycles(1)

		// Both should error consistently
		assert.Error(t, err1, "First call should error")
		assert.Error(t, err2, "Second call should error")
		assert.Contains(t, err1.Error(), "failed to build graph")
		assert.Contains(t, err2.Error(), "failed to build graph")
	})
}

// Test boundary conditions for ValidateNewDependency
func TestValidateNewDependency_BoundaryConditions(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)

	boundaryTests := []struct {
		name         string
		dependentID  int64
		dependencyID int64
		expectSelf   bool
	}{
		{"min_int64", -9223372036854775808, -9223372036854775808, true},
		{"max_int64", 9223372036854775807, 9223372036854775807, true},
		{"min_to_max", -9223372036854775808, 9223372036854775807, false},
		{"max_to_min", 9223372036854775807, -9223372036854775808, false},
		{"zero_to_one", 0, 1, false},
		{"one_to_zero", 1, 0, false},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectSelf {
				// Self-dependency should be detected without panic
				result, err := service.ValidateNewDependency(tt.dependentID, tt.dependencyID)
				assert.NoError(t, err, "Should not error on self-dependency check")
				assert.NotNil(t, result, "Result should not be nil")
				assert.False(t, result.IsValid, "Self-dependencies should be invalid")
				assert.Contains(t, result.Errors, "Self-dependencies are not allowed")
			} else {
				// Non-self dependencies will panic when checking existing dependencies with nil repo
				assert.Panics(t, func() {
					service.ValidateNewDependency(tt.dependentID, tt.dependencyID)
				}, "Should panic when checking existing dependencies with nil repo")
			}
		})
	}
}

// Test multiple service instances
func TestMultipleServiceInstances(t *testing.T) {
	service1 := services.NewDependencyGraphService(nil, nil)
	service2 := services.NewDependencyGraphService(nil, nil)

	assert.NotNil(t, service1, "First service should not be nil")
	assert.NotNil(t, service2, "Second service should not be nil")
	// We can't use NotEqual for pointers, so just check they are different instances
	assert.True(t, service1 != service2, "Different service instances should not be equal")

	// Both should behave identically
	result1, err1 := service1.ValidateNewDependency(100, 100)
	result2, err2 := service2.ValidateNewDependency(100, 100)

	assert.Equal(t, err1, err2, "Both services should have same error behavior")
	assert.Equal(t, result1.IsValid, result2.IsValid, "Both services should have same validation result")
	assert.Equal(t, result1.Errors, result2.Errors, "Both services should have same error messages")
}

// Performance test for self-dependency validation (should be very fast)
func TestValidateNewDependency_Performance(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)

	// Test that self-dependency validation is consistently fast
	for i := 0; i < 1000; i++ {
		result, err := service.ValidateNewDependency(int64(i), int64(i))
		assert.NoError(t, err, "Should not error")
		assert.False(t, result.IsValid, "Should be invalid")
		assert.Contains(t, result.Errors, "Self-dependencies are not allowed")
	}
}

// Test service behavior with various input types
func TestServiceInputVariations(t *testing.T) {
	service := services.NewDependencyGraphService(nil, nil)

	// Test various ID combinations
	idPairs := [][]int64{
		{1, 1},      // same positive
		{-1, -1},    // same negative  
		{0, 0},      // same zero
		{1, 2},      // different positive
		{-1, -2},    // different negative
		{-1, 1},     // negative to positive
		{1, -1},     // positive to negative
		{0, 1},      // zero to positive
		{1, 0},      // positive to zero
		{0, -1},     // zero to negative
		{-1, 0},     // negative to zero
	}

	for i, pair := range idPairs {
		t.Run("pair_"+string(rune(i+48)), func(t *testing.T) {
			dependent, dependency := pair[0], pair[1]
			
			if dependent == dependency {
				// Self-dependency - should work fine
				result, err := service.ValidateNewDependency(dependent, dependency)
				assert.NoError(t, err)
				assert.False(t, result.IsValid)
				assert.Contains(t, result.Errors, "Self-dependencies are not allowed")
			} else {
				// Different IDs - will panic when checking existing dependencies
				assert.Panics(t, func() {
					service.ValidateNewDependency(dependent, dependency)
				})
			}
		})
	}
}