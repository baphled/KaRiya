//nolint:errcheck // Benchmark tests - error handling not relevant for performance measurement.
package intents

import (
	"context"
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
)

// CaptureEvent benchmarks have been moved to internal/cli/intents/captureevent/benchmarks_test.go

// BrowseTimeline benchmarks have been moved to internal/cli/intents/browsetimeline/benchmarks_test.go

// BenchmarkGenerateCVInit benchmarks GenerateCV intent initialization.
func BenchmarkGenerateCVInit(b *testing.B) {
	ctx := &GenerateCVContext{
		AvailableProfiles: make([]*CVProfile, 0),
		Events:            make([]*career.Event, 0),
		Facts:             make([]*career.Fact, 0),
		DefaultProfile:    nil,
	}

	b.ResetTimer()
	for range b.N {
		intent, _ := NewGenerateCVIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkGenerateCVView benchmarks GenerateCV intent view rendering.
func BenchmarkGenerateCVView(b *testing.B) {
	ctx := &GenerateCVContext{
		AvailableProfiles: make([]*CVProfile, 0),
		Events:            make([]*career.Event, 0),
		Facts:             make([]*career.Fact, 0),
		DefaultProfile:    nil,
	}
	intent, _ := NewGenerateCVIntent(ctx)

	b.ResetTimer()
	for range b.N {
		_ = intent.View()
	}
}

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
