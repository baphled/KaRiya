//nolint:errcheck // Benchmark tests - error handling not relevant for performance measurement.
package capture_event

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// BenchmarkCaptureEventInit benchmarks CaptureEvent intent initialization.
func BenchmarkCaptureEventInit(b *testing.B) {
	ctx := &IntentContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intent, _ := NewIntent(ctx)
		_ = intent.Init()
	}
}

// BenchmarkCaptureEventView benchmarks CaptureEvent intent view rendering.
func BenchmarkCaptureEventView(b *testing.B) {
	ctx := &IntentContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, _ := NewIntent(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.View()
	}
}

// BenchmarkCaptureEventUpdate benchmarks CaptureEvent intent message handling.
func BenchmarkCaptureEventUpdate(b *testing.B) {
	ctx := &IntentContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, _ := NewIntent(ctx)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intent.Update(msg)
	}
}
