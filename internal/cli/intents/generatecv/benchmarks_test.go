//nolint:errcheck // Benchmark tests - error handling not relevant for performance measurement.
package generatecv

import (
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
)

// BenchmarkGenerateCVInit benchmarks GenerateCV intent initialization.
func BenchmarkGenerateCVInit(b *testing.B) {
	ctx := &IntentContext{
		AvailableProfiles: make([]*CVProfile, 0),
		Events:            make([]*career.Event, 0),
		Facts:             make([]*career.Fact, 0),
		DefaultProfile:    nil,
	}

	b.ResetTimer()
	for range b.N {
		intent, _ := NewIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkGenerateCVView benchmarks GenerateCV intent view rendering.
func BenchmarkGenerateCVView(b *testing.B) {
	ctx := &IntentContext{
		AvailableProfiles: make([]*CVProfile, 0),
		Events:            make([]*career.Event, 0),
		Facts:             make([]*career.Fact, 0),
		DefaultProfile:    nil,
	}
	intent, _ := NewIntent(ctx)

	b.ResetTimer()
	for range b.N {
		_ = intent.View()
	}
}
