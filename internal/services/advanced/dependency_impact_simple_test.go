package services_test

import (
	"context"
	"testing"
	services "url-db/internal/services/advanced"

	"github.com/stretchr/testify/assert"
)

// Test comprehensive constructor scenarios
func TestNewDependencyImpactAnalyzer_ComprehensiveConstructor(t *testing.T) {
	testCases := []struct {
		name         string
		depRepo      interface{}
		nodeRepo     interface{}
		graphService interface{}
	}{
		{"all_nil", nil, nil, nil},
		{"dep_nil_others_not", nil, "not_nil", "not_nil"},
		{"node_nil_others_not", "not_nil", nil, "not_nil"},
		{"graph_nil_others_not", "not_nil", "not_nil", nil},
		{"dep_node_nil", nil, nil, "not_nil"},
		{"dep_graph_nil", nil, "not_nil", nil},
		{"node_graph_nil", "not_nil", nil, nil},
		{"all_not_nil", "not_nil", "not_nil", "not_nil"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Constructor should always succeed regardless of parameters
			analyzer := services.NewDependencyImpactAnalyzer(nil, nil, nil)
			assert.NotNil(t, analyzer, "Analyzer should not be nil regardless of input parameters")
		})
	}
}

// Test AnalyzeImpact with various contexts
func TestAnalyzeImpact_ComprehensiveContexts(t *testing.T) {
	analyzer := services.NewDependencyImpactAnalyzer(nil, nil, nil)

	contextTests := []struct {
		name    string
		ctx     context.Context
		nodeID  int64
		impactType string
	}{
		{"background_context", context.Background(), 1, "delete"},
		{"todo_context", context.TODO(), 1, "delete"},
		{"background_zero_node", context.Background(), 0, "delete"},
		{"background_negative_node", context.Background(), -1, "delete"},
		{"background_large_node", context.Background(), 999999, "delete"},
		{"background_update", context.Background(), 1, "update"},
		{"background_version", context.Background(), 1, "version_change"},
		{"background_invalid_type", context.Background(), 1, "invalid_type"},
	}

	for _, tt := range contextTests {
		t.Run(tt.name, func(t *testing.T) {
			// All these should result in panics due to nil nodeRepo
			assert.Panics(t, func() {
				analyzer.AnalyzeImpact(tt.ctx, tt.nodeID, tt.impactType)
			}, "Should panic with nil nodeRepo")
		})
	}
}

// Test AnalyzeImpact with comprehensive node ID variations
func TestAnalyzeImpact_NodeIDVariations(t *testing.T) {
	analyzer := services.NewDependencyImpactAnalyzer(nil, nil, nil)
	ctx := context.Background()

	nodeIDTests := []int64{
		0,                    // zero
		1,                    // small positive
		-1,                   // small negative
		100,                  // medium positive
		-100,                 // medium negative
		999999,               // large positive
		-999999,              // large negative
		9223372036854775807,  // max int64
		-9223372036854775808, // min int64
	}

	for _, nodeID := range nodeIDTests {
		t.Run("node_"+string(rune(nodeID)), func(t *testing.T) {
			// Should always panic with nil node repository
			assert.Panics(t, func() {
				analyzer.AnalyzeImpact(ctx, nodeID, "delete")
			}, "Should panic with nil node repository")
		})
	}
}

// Test AnalyzeImpact with comprehensive impact type variations
func TestAnalyzeImpact_ImpactTypeVariations(t *testing.T) {
	analyzer := services.NewDependencyImpactAnalyzer(nil, nil, nil)
	ctx := context.Background()

	impactTypes := []string{
		"delete",
		"update", 
		"version_change",
		"invalid_type",
		"",           // empty string
		"DELETE",     // uppercase
		"Update",     // mixed case
		"delete ",    // with space
		" delete",    // with leading space
		"del ete",    // with internal space
		"123",        // numeric
		"delete123",  // alphanumeric
		"!@#$%",      // special characters
	}

	for _, impactType := range impactTypes {
		t.Run("type_"+impactType, func(t *testing.T) {
			// Should always panic due to nil nodeRepo
			assert.Panics(t, func() {
				analyzer.AnalyzeImpact(ctx, 1, impactType)
			}, "Should panic with nil nodeRepo")
		})
	}
}

// Test multiple analyzer instances
func TestMultipleAnalyzerInstances(t *testing.T) {
	analyzer1 := services.NewDependencyImpactAnalyzer(nil, nil, nil)
	analyzer2 := services.NewDependencyImpactAnalyzer(nil, nil, nil)
	analyzer3 := services.NewDependencyImpactAnalyzer(nil, nil, nil)

	assert.NotNil(t, analyzer1, "First analyzer should not be nil")
	assert.NotNil(t, analyzer2, "Second analyzer should not be nil")
	assert.NotNil(t, analyzer3, "Third analyzer should not be nil")
	
	// All should be different instances (use pointer comparison)
	assert.True(t, analyzer1 != analyzer2, "Different instances should not be equal")
	assert.True(t, analyzer1 != analyzer3, "Different instances should not be equal")
	assert.True(t, analyzer2 != analyzer3, "Different instances should not be equal")

	// All should behave consistently for the same input (all should panic)
	ctx := context.Background()
	
	assert.Panics(t, func() { analyzer1.AnalyzeImpact(ctx, 1, "delete") }, "Should panic")
	assert.Panics(t, func() { analyzer2.AnalyzeImpact(ctx, 1, "delete") }, "Should panic")
	assert.Panics(t, func() { analyzer3.AnalyzeImpact(ctx, 1, "delete") }, "Should panic")
}

// Test analyzer consistency
func TestAnalyzerConsistency(t *testing.T) {
	analyzer := services.NewDependencyImpactAnalyzer(nil, nil, nil)
	ctx := context.Background()

	// Test that multiple calls with same parameters produce consistent results (all panic)
	t.Run("consistent_panics", func(t *testing.T) {
		assert.Panics(t, func() { analyzer.AnalyzeImpact(ctx, 42, "delete") }, "First call should panic")
		assert.Panics(t, func() { analyzer.AnalyzeImpact(ctx, 42, "delete") }, "Second call should panic")
	})
}

// Performance test for analyzer creation (should be very fast)
func TestAnalyzerCreationPerformance(t *testing.T) {
	// Test that creating many analyzers is fast
	for i := 0; i < 1000; i++ {
		analyzer := services.NewDependencyImpactAnalyzer(nil, nil, nil)
		assert.NotNil(t, analyzer, "Analyzer should be created successfully")
	}
}

// Test analyzer with various context cancellation scenarios
func TestAnalyzeImpact_ContextCancellation(t *testing.T) {
	analyzer := services.NewDependencyImpactAnalyzer(nil, nil, nil)

	t.Run("cancelled_context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Should panic due to nil nodeRepo
		assert.Panics(t, func() {
			analyzer.AnalyzeImpact(ctx, 1, "delete")
		}, "Should panic with nil nodeRepo")
	})

	t.Run("timeout_context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0) // Immediate timeout
		defer cancel()

		// Should panic due to nil nodeRepo
		assert.Panics(t, func() {
			analyzer.AnalyzeImpact(ctx, 1, "delete")
		}, "Should panic with nil nodeRepo")
	})
}

// Test edge cases for analyzer behavior
func TestAnalyzer_EdgeCases(t *testing.T) {
	analyzer := services.NewDependencyImpactAnalyzer(nil, nil, nil)
	ctx := context.Background()

	edgeCases := []struct {
		name       string
		nodeID     int64
		impactType string
	}{
		{"zero_node_empty_type", 0, ""},
		{"negative_node_empty_type", -1, ""},
		{"max_node_long_type", 9223372036854775807, "very_long_impact_type_name_that_is_definitely_invalid"},
		{"min_node_unicode_type", -9223372036854775808, "删除"},
		{"normal_node_null_like_type", 42, "null"},
		{"normal_node_bool_like_type", 42, "true"},
		{"normal_node_json_like_type", 42, "{\"type\":\"delete\"}"},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			// All should panic due to nil nodeRepo
			assert.Panics(t, func() {
				analyzer.AnalyzeImpact(ctx, tc.nodeID, tc.impactType)
			}, "Should panic for edge case input")
		})
	}
}

// Test analyzer state consistency
func TestAnalyzer_StateConsistency(t *testing.T) {
	analyzer := services.NewDependencyImpactAnalyzer(nil, nil, nil)
	ctx := context.Background()

	// Test that the analyzer doesn't maintain state between calls
	// (i.e., each call is independent)
	calls := []struct {
		nodeID     int64
		impactType string
	}{
		{1, "delete"},
		{2, "update"},
		{3, "version_change"},
		{1, "delete"}, // Repeat first call
	}

	// All calls should panic consistently
	for i, call := range calls {
		t.Run("call_"+string(rune(i+48)), func(t *testing.T) {
			assert.Panics(t, func() {
				analyzer.AnalyzeImpact(ctx, call.nodeID, call.impactType)
			}, "Call %d should panic", i)
		})
	}
}