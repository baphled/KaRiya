package components

import (
	"strings"
	"testing"

	"github.com/baphled/kariya/internal/cli/terminal"
)

// Performance targets from Task 16 Phase 5:
// - StandardView render < 50ms
// - Modal render < 20ms
// - Full view render < 100ms

// BenchmarkStandardViewRender benchmarks StandardView rendering performance
func BenchmarkStandardViewRender(b *testing.B) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)

	view := NewStandardView(termInfo).
		WithLogo(logo, 2).
		WithBreadcrumbs("Main Menu", "Settings", "Profile").
		WithContent("Test content\nLine 2\nLine 3").
		WithHelp("q Quit  h Help  m Menu").
		WithFooterSeparator(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.Render()
	}
}

// BenchmarkStandardViewRenderWithLongContent benchmarks rendering with long content
func BenchmarkStandardViewRenderWithLongContent(b *testing.B) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)

	// Create long content
	longContent := strings.Repeat("Line of content\n", 100)

	view := NewStandardView(termInfo).
		WithLogo(logo, 2).
		WithBreadcrumbs("Main Menu", "Content").
		WithContent(longContent).
		WithHelp("q Quit").
		WithFooterSeparator(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.Render()
	}
}

// BenchmarkStandardViewRenderMinimal benchmarks minimal view rendering
func BenchmarkStandardViewRenderMinimal(b *testing.B) {
	termInfo := &terminal.Info{Width: 120, Height: 40}

	view := NewStandardView(termInfo).
		WithContent("Minimal content").
		WithHelp("q Quit")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.Render()
	}
}

// BenchmarkModalRenderError benchmarks error modal rendering
func BenchmarkModalRenderError(b *testing.B) {
	modal := NewErrorModal("Error", "An error occurred while processing your request")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modal.Render(120, 40)
	}
}

// BenchmarkModalRenderLoading benchmarks loading modal rendering
func BenchmarkModalRenderLoading(b *testing.B) {
	modal := NewLoadingModal("Processing your request...", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modal.Render(120, 40)
	}
}

// BenchmarkModalRenderProgress benchmarks progress modal rendering
func BenchmarkModalRenderProgress(b *testing.B) {
	modal := NewProgressModal("Processing", "Processing items...", 0.5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modal.Render(120, 40)
	}
}

// BenchmarkModalRenderSuccess benchmarks success modal rendering
func BenchmarkModalRenderSuccess(b *testing.B) {
	modal := NewSuccessModal("Operation completed successfully!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modal.Render(120, 40)
	}
}

// BenchmarkFullViewRender benchmarks full view with logo + content + modal
func BenchmarkFullViewRender(b *testing.B) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)
	modal := NewLoadingModal("Processing...", false)

	view := NewStandardView(termInfo).
		WithLogo(logo, 2).
		WithBreadcrumbs("Main Menu", "Operations", "Processing").
		WithContent("Background content while processing...").
		WithHelp("Esc Cancel  q Quit").
		WithFooterSeparator(true).
		ShowModalOverlay(modal)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.Render()
	}
}

// BenchmarkFullViewRenderWithProgress benchmarks full view with progress modal
func BenchmarkFullViewRenderWithProgress(b *testing.B) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)
	modal := NewProgressModal("Processing", "Analyzing data...", 0.75)

	view := NewStandardView(termInfo).
		WithLogo(logo, 2).
		WithBreadcrumbs("Main Menu", "Operations").
		WithContent("Processing your data...").
		WithHelp("q Quit").
		ShowModalOverlay(modal)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.Render()
	}
}

// BenchmarkLoadingMessageRotation benchmarks message rotation performance
func BenchmarkLoadingMessageRotation(b *testing.B) {
	messages := []string{
		"Processing request...",
		"Analyzing data...",
		"Generating results...",
		"Finalizing...",
	}
	rotator := NewLoadingMessageRotator(messages, 0) // No time-based rotation

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rotator.Rotate()
		_ = rotator.GetCurrent()
	}
}

// BenchmarkSpinnerAdvance benchmarks spinner advancement
func BenchmarkSpinnerAdvance(b *testing.B) {
	spinner := NewSimpleSpinner()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spinner.Advance()
		_ = spinner.GetFrame()
	}
}

// BenchmarkStandardViewRenderAtDifferentSizes benchmarks rendering at various terminal sizes
func BenchmarkStandardViewRenderAtDifferentSizes(b *testing.B) {
	sizes := []*terminal.Info{
		{Width: 80, Height: 24},  // Tiny
		{Width: 120, Height: 40}, // Normal
		{Width: 200, Height: 60}, // Large
	}

	for _, size := range sizes {
		b.Run("Size_"+string(rune(size.Width))+"x"+string(rune(size.Height)), func(b *testing.B) {
			logo := NewASCIILogo(false, size.Width)
			view := NewStandardView(size).
				WithLogo(logo, 2).
				WithBreadcrumbs("Main Menu", "Settings").
				WithContent("Test content").
				WithHelp("q Quit").
				WithFooterSeparator(true)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = view.Render()
			}
		})
	}
}

// BenchmarkModalRenderWithLongMessage benchmarks modal with long message
func BenchmarkModalRenderWithLongMessage(b *testing.B) {
	longMessage := strings.Repeat("Error detail. ", 50)
	modal := NewErrorModal("Error", longMessage)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modal.Render(120, 40)
	}
}

// BenchmarkASCIILogoRender benchmarks logo rendering
func BenchmarkASCIILogoRender(b *testing.B) {
	logo := NewASCIILogo(false, 120)
	logo.SetExternalCentering(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = logo.ViewStatic()
	}
}

// BenchmarkStandardViewWithAllFeatures benchmarks view with all features enabled
func BenchmarkStandardViewWithAllFeatures(b *testing.B) {
	termInfo := &terminal.Info{Width: 160, Height: 50}
	logo := NewASCIILogo(false, termInfo.Width)
	modal := NewProgressModal("Processing", "Complex operation in progress...", 0.65)

	view := NewStandardView(termInfo).
		WithLogo(logo, 2).
		WithBreadcrumbs("Main Menu", "Advanced", "Operations", "Complex Task").
		WithTitle("Complex Operation", "Performing advanced processing").
		WithContent(strings.Repeat("Data line\n", 20)).
		WithHelp("↑↓ Navigate  Esc Cancel  q Quit  h Help  m Menu").
		WithFooterSeparator(true).
		ShowModalOverlay(modal)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.Render()
	}
}

// BenchmarkConcurrentRendering benchmarks concurrent view rendering
func BenchmarkConcurrentRendering(b *testing.B) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)

	b.RunParallel(func(pb *testing.PB) {
		view := NewStandardView(termInfo).
			WithLogo(logo, 2).
			WithContent("Concurrent test").
			WithHelp("q Quit")

		for pb.Next() {
			_ = view.Render()
		}
	})
}

// Memory allocation benchmarks

// BenchmarkStandardViewRenderAllocs benchmarks memory allocations for StandardView
func BenchmarkStandardViewRenderAllocs(b *testing.B) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)

	view := NewStandardView(termInfo).
		WithLogo(logo, 2).
		WithBreadcrumbs("Main Menu", "Settings").
		WithContent("Test content").
		WithHelp("q Quit").
		WithFooterSeparator(true)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.Render()
	}
}

// BenchmarkModalRenderAllocs benchmarks memory allocations for modal
func BenchmarkModalRenderAllocs(b *testing.B) {
	modal := NewLoadingModal("Processing...", false)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modal.Render(120, 40)
	}
}

// BenchmarkFullViewRenderAllocs benchmarks memory allocations for full view
func BenchmarkFullViewRenderAllocs(b *testing.B) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)
	modal := NewLoadingModal("Processing...", false)

	view := NewStandardView(termInfo).
		WithLogo(logo, 2).
		WithContent("Content").
		WithHelp("q Quit").
		ShowModalOverlay(modal)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.Render()
	}
}
