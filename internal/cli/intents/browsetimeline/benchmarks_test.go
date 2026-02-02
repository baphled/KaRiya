package browsetimeline

import (
	"fmt"
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

// createBenchmarkEvents creates n test events for benchmarking.
func createBenchmarkEvents(n int) []*career.Event {
	events := make([]*career.Event, n)
	for i := range n {
		events[i] = fixtures.EventWith(fmt.Sprintf("event-%d", i), "Benchmark event", "BenchCorp", "")
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
	for range b.N {
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
	for range b.N {
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
	for range b.N {
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
	for range b.N {
		intent, _ := NewIntent(ctx)
		_ = intent.Init()
	}
}
