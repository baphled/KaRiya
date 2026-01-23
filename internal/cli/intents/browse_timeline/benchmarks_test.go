package browse_timeline

import (
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
)

// BenchmarkBrowseTimelineInit benchmarks BrowseTimeline intent initialization
func BenchmarkBrowseTimelineInit(b *testing.B) {
	ctx := &IntentContext{
		Events:          make([]*career.CareerEvent, 0),
		InitialFilters:  nil,
		SelectedEventID: "",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkBrowseTimelineView benchmarks BrowseTimeline intent view rendering
func BenchmarkBrowseTimelineView(b *testing.B) {
	ctx := &IntentContext{
		Events:          make([]*career.CareerEvent, 0),
		InitialFilters:  nil,
		SelectedEventID: "",
	}
	intent, _ := NewIntent(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.View()
	}
}
