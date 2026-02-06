//nolint:errcheck // Benchmark tests - error handling not relevant for performance measurement.
package intents

import (
	"context"
	"testing"
)

// CaptureEvent benchmarks have been moved to internal/cli/intents/captureevent/benchmarks_test.go

// BrowseTimeline benchmarks have been moved to internal/cli/intents/browsetimeline/benchmarks_test.go

// GenerateCV benchmarks have been moved to internal/cli/intents/generatecv/benchmarks_test.go

// BenchmarkConfigureSystemInit benchmarks ConfigureSystem intent initialization.
func BenchmarkConfigureSystemInit(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		intent, _ := NewConfigureSystemIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkConfigureSystemView benchmarks ConfigureSystem intent view rendering.
func BenchmarkConfigureSystemView(b *testing.B) {
	ctx := context.Background()
	intent, _ := NewConfigureSystemIntent(ctx)

	b.ResetTimer()
	for range b.N {
		_ = intent.View()
	}
}

// BenchmarkIntentRouterActivation benchmarks IntentRouter intent activation.
func BenchmarkIntentRouterActivation(b *testing.B) {
	router := NewDefaultIntentRouter()
	//nolint:errcheck // Benchmark setup - error handling not relevant.
	router.RegisterIntent("test", func() Intent {
		intent, _ := NewConfigureSystemIntent(b.Context())
		return intent
	})

	b.ResetTimer()
	for range b.N {
		router.ActivateIntent("test", make(map[string]interface{}))
	}
}

// BenchmarkIntentResultCreation benchmarks IntentResult creation.
func BenchmarkIntentResultCreation(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_ = &IntentResult[interface{}]{
			Status: Completed,
			Data:   nil,
			Error:  nil,
		}
	}
}

// BenchmarkIntentResultMetadata benchmarks IntentResult metadata operations.
func BenchmarkIntentResultMetadata(b *testing.B) {
	result := &IntentResult[interface{}]{
		Status:   Completed,
		Data:     nil,
		Error:    nil,
		Metadata: make(map[string]interface{}),
	}

	b.ResetTimer()
	for range b.N {
		result.WithMetadata("key", "value")
		_, _ = result.GetMetadata("key")
	}
}
