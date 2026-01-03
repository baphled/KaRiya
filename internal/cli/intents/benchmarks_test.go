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

// BenchmarkBrowseTimelineInit benchmarks BrowseTimeline intent initialization
func BenchmarkBrowseTimelineInit(b *testing.B) {
	ctx := &BrowseTimelineContext{
		Events:          make([]*career.CareerEvent, 0),
		InitialFilters:  nil,
		SelectedEventID: "",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewBrowseTimelineIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkBrowseTimelineView benchmarks BrowseTimeline intent view rendering
func BenchmarkBrowseTimelineView(b *testing.B) {
	ctx := &BrowseTimelineContext{
		Events:          make([]*career.CareerEvent, 0),
		InitialFilters:  nil,
		SelectedEventID: "",
	}
	intent, _ := NewBrowseTimelineIntent(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.View()
	}
}

// BenchmarkGenerateCVInit benchmarks GenerateCV intent initialization
func BenchmarkGenerateCVInit(b *testing.B) {
	ctx := &GenerateCVContext{
		AvailableProfiles: make([]*CVProfile, 0),
		Events:            make([]*career.CareerEvent, 0),
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
		Events:            make([]*career.CareerEvent, 0),
		Facts:             make([]*career.Fact, 0),
		DefaultProfile:    nil,
	}
	intent, _ := NewGenerateCVIntent(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.View()
	}
}

// BenchmarkExportArtifactInit benchmarks ExportArtifact intent initialization
func BenchmarkExportArtifactInit(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewExportArtifactIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkExportArtifactView benchmarks ExportArtifact intent view rendering
func BenchmarkExportArtifactView(b *testing.B) {
	ctx := context.Background()
	intent, _ := NewExportArtifactIntent(ctx)

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

// BenchmarkIntentRouterActivation benchmarks IntentRouter intent activation
func BenchmarkIntentRouterActivation(b *testing.B) {
	router := NewDefaultIntentRouter()
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
		_, _ = router.ActivateIntent("test", make(map[string]interface{}))
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
