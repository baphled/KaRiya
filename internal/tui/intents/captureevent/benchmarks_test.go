package captureevent

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// BenchmarkCaptureEventInit benchmarks CaptureEvent intent initialization.
func BenchmarkCaptureEventInit(b *testing.B) {
	ctx := &IntentValidator{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}

	b.ResetTimer()
	for range b.N {
		intent, err := NewIntent(ctx)
		if err != nil {
			b.Fatal(err)
		}
		_ = intent.Init()
	}
}

// BenchmarkCaptureEventView benchmarks CaptureEvent intent view rendering.
func BenchmarkCaptureEventView(b *testing.B) {
	ctx := &IntentValidator{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewIntent(ctx)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for range b.N {
		_ = intent.View()
	}
}

// BenchmarkCaptureEventUpdate benchmarks CaptureEvent intent message handling.
func BenchmarkCaptureEventUpdate(b *testing.B) {
	ctx := &IntentValidator{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewIntent(ctx)
	if err != nil {
		b.Fatal(err)
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	b.ResetTimer()
	for range b.N {
		_ = intent.Update(msg)
	}
}
