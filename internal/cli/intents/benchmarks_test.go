//nolint:errcheck // Benchmark tests - error handling not relevant for performance measurement.
package intents

import (
	"context"
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// BenchmarkCaptureEventInit benchmarks CaptureEvent intent initialization
func BenchmarkCaptureEventInit(b *testing.B) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewCaptureEventIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkCaptureEventView benchmarks CaptureEvent intent view rendering
func BenchmarkCaptureEventView(b *testing.B) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, _ := NewCaptureEventIntent(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.View()
	}
}

// BenchmarkCaptureEventUpdate benchmarks CaptureEvent intent message handling
func BenchmarkCaptureEventUpdate(b *testing.B) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, _ := NewCaptureEventIntent(ctx)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.Update(msg)
	}
}

// BrowseTimeline benchmarks have been moved to internal/cli/intents/browsetimeline/benchmarks_test.go

// BenchmarkGenerateCVInit benchmarks GenerateCV intent initialization
func BenchmarkGenerateCVInit(b *testing.B) {
	ctx := &GenerateCVContext{
		AvailableProfiles: make([]*CVProfile, 0),
		Events:            make([]*career.Event, 0),
		Facts:             make([]*career.Fact, 0),
		DefaultProfile:    nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewGenerateCVIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkGenerateCVView benchmarks GenerateCV intent view rendering
func BenchmarkGenerateCVView(b *testing.B) {
	ctx := &GenerateCVContext{
		AvailableProfiles: make([]*CVProfile, 0),
		Events:            make([]*career.Event, 0),
		Facts:             make([]*career.Fact, 0),
		DefaultProfile:    nil,
	}
	intent, _ := NewGenerateCVIntent(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.View()
	}
}

// BenchmarkConfigureSystemInit benchmarks ConfigureSystem intent initialization
func BenchmarkConfigureSystemInit(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewConfigureSystemIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkConfigureSystemView benchmarks ConfigureSystem intent view rendering
func BenchmarkConfigureSystemView(b *testing.B) {
	ctx := context.Background()
	intent, _ := NewConfigureSystemIntent(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.View()
	}
}

// BenchmarkIntentRouterActivation benchmarks IntentRouter intent activation.
func BenchmarkIntentRouterActivation(b *testing.B) {
	router := NewDefaultIntentRouter()
	//nolint:errcheck // Benchmark setup - error handling not relevant.
	router.RegisterIntent("test", func() Intent {
		ctx := &CaptureEventContext{
			CaptureStrategy: "manual",
			PreviousEvent:   nil,
			Metadata:        make(map[string]string),
		}
		intent, _ := NewCaptureEventIntent(ctx)
		return intent
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//nolint:errcheck // Benchmark loop - error handling not relevant.
		router.ActivateIntent("test", make(map[string]interface{}))
	}
}

// BenchmarkIntentResultCreation benchmarks IntentResult creation
func BenchmarkIntentResultCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &IntentResult[interface{}]{
			Status: Completed,
			Data:   nil,
			Error:  nil,
		}
	}
}

// BenchmarkIntentResultMetadata benchmarks IntentResult metadata operations
func BenchmarkIntentResultMetadata(b *testing.B) {
	result := &IntentResult[interface{}]{
		Status:   Completed,
		Data:     nil,
		Error:    nil,
		Metadata: make(map[string]interface{}),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result.WithMetadata("key", "value")
		_, _ = result.GetMetadata("key")
	}
}
