//nolint:errcheck // Benchmark tests - error handling not relevant for performance measurement.
package intents_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/configure"
	"github.com/baphled/kariya/internal/config"
)

// CaptureEvent benchmarks have been moved to internal/cli/intents/captureevent/benchmarks_test.go

// BrowseTimeline benchmarks have been moved to internal/cli/intents/browsetimeline/benchmarks_test.go

// GenerateCV benchmarks have been moved to internal/cli/intents/generatecv/benchmarks_test.go

// BenchmarkConfigureSystemInit benchmarks ConfigureSystem intent initialization.
func BenchmarkConfigureSystemInit(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		cfg := config.DefaultConfig()
		intentCtx := &configure.IntentContext{
			Cfg:      cfg,
			Settings: configure.SettingsFromConfig(cfg),
		}
		intent, _ := configure.NewIntent(intentCtx)
		_ = intent.Init()
	}
}

// BenchmarkConfigureSystemView benchmarks ConfigureSystem intent view rendering.
func BenchmarkConfigureSystemView(b *testing.B) {
	cfg := config.DefaultConfig()
	intentCtx := &configure.IntentContext{
		Cfg:      cfg,
		Settings: configure.SettingsFromConfig(cfg),
	}
	intent, _ := configure.NewIntent(intentCtx)

	b.ResetTimer()
	for range b.N {
		_ = intent.View()
	}
}

// BenchmarkIntentRouterActivation benchmarks IntentRouter intent activation.
func BenchmarkIntentRouterActivation(b *testing.B) {
	router := intents.NewDefaultIntentRouter()
	//nolint:errcheck // Benchmark setup - error handling not relevant.
	router.RegisterIntent("test", func() intents.Intent {
		cfg := config.DefaultConfig()
		intentCtx := &configure.IntentContext{
			Cfg:      cfg,
			Settings: configure.SettingsFromConfig(cfg),
		}
		intent, _ := configure.NewIntent(intentCtx)
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
		_ = &intents.IntentResult[interface{}]{
			Status: intents.Completed,
			Data:   nil,
			Error:  nil,
		}
	}
}

// BenchmarkIntentResultMetadata benchmarks IntentResult metadata operations.
func BenchmarkIntentResultMetadata(b *testing.B) {
	result := &intents.IntentResult[interface{}]{
		Status:   intents.Completed,
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
