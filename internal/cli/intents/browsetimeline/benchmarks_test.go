package browsetimeline

import (
	"fmt"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// createBenchmarkEvents creates n test events for benchmarking.
func createBenchmarkEvents(n int) []*career.CareerEvent {
	events := make([]*career.CareerEvent, n)
	now := time.Now()
	for i := 0; i < n; i++ {
		events[i] = &career.CareerEvent{
			ID:        fmt.Sprintf("event-%d", i),
			Text:      "Benchmark event",
			Company:   "BenchCorp",
			Date:      now,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return events
}

// BenchmarkBrowseTimelineInit benchmarks BrowseTimeline intent initialization.
func BenchmarkBrowseTimelineInit(b *testing.B) {
	events := createBenchmarkEvents(10)
	ctx := &IntentContext{
		Events:          events,
		InitialFilters:  nil,
		SelectedEventID: "",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkBrowseTimelineView benchmarks BrowseTimeline intent view rendering.
func BenchmarkBrowseTimelineView(b *testing.B) {
	events := createBenchmarkEvents(10)
	ctx := &IntentContext{
		Events:          events,
		InitialFilters:  nil,
		SelectedEventID: "",
	}
	intent, _ := NewIntent(ctx)
	intent.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.View()
	}
}

// BenchmarkBrowseTimelineInitEmpty benchmarks intent initialization with empty events.
func BenchmarkBrowseTimelineInitEmpty(b *testing.B) {
	ctx := &IntentContext{
		Events:          nil,
		InitialFilters:  nil,
		SelectedEventID: "",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkBrowseTimelineInitLarge benchmarks intent initialization with many events.
func BenchmarkBrowseTimelineInitLarge(b *testing.B) {
	events := createBenchmarkEvents(100)
	ctx := &IntentContext{
		Events:          events,
		InitialFilters:  nil,
		SelectedEventID: "",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewIntent(ctx)
		_ = intent.Init()
	}
}
