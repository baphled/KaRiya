package components

import (
	"strings"
	"testing"

	"github.com/baphled/kariya/internal/cli/terminal"
)

// Terminal size fixtures for testing
var (
	terminalSizeTiny      = &terminal.Info{Width: 80, Height: 24}  // Minimum viable
	terminalSizeCompact   = &terminal.Info{Width: 100, Height: 30} // Small laptop
	terminalSizeNormal    = &terminal.Info{Width: 120, Height: 40} // Standard
	terminalSizeLarge     = &terminal.Info{Width: 160, Height: 50} // Large desktop
	terminalSizeXLarge    = &terminal.Info{Width: 200, Height: 60} // Very large
	terminalSizeUltraWide = &terminal.Info{Width: 240, Height: 40} // Ultra-wide monitor
	terminalSizeTall      = &terminal.Info{Width: 120, Height: 80} // Tall terminal
)

// TestStandardView_AllTerminalSizes tests StandardView rendering at all common terminal sizes
func TestStandardView_AllTerminalSizes(t *testing.T) {
	sizes := map[string]*terminal.Info{
		"Tiny (80x24)":       terminalSizeTiny,
		"Compact (100x30)":   terminalSizeCompact,
		"Normal (120x40)":    terminalSizeNormal,
		"Large (160x50)":     terminalSizeLarge,
		"XLarge (200x60)":    terminalSizeXLarge,
		"UltraWide (240x40)": terminalSizeUltraWide,
		"Tall (120x80)":      terminalSizeTall,
	}

	for name, size := range sizes {
		t.Run(name, func(t *testing.T) {
			testStandardViewAtSize(t, size)
		})
	}
}

func testStandardViewAtSize(t *testing.T, size *terminal.Info) {
	logo := NewASCIILogo(false, size.Width)
	view := NewStandardView(size).
		WithLogo(logo, 2).
		WithBreadcrumbs("Home", "Settings", "Profile").
		WithContent("Test content\nLine 2\nLine 3").
		WithHelp("q Quit  h Help  m Menu").
		WithFooterSeparator(true)

	output := view.Render()

	// Basic assertions
	if output == "" {
		t.Errorf("Expected non-empty output for terminal %dx%d", size.Width, size.Height)
	}

	// Logo should be present
	if !strings.Contains(output, "KARIYA") && !strings.Contains(output, "██") {
		t.Error("Expected logo to be visible")
	}

	// Breadcrumbs should be present
	if !strings.Contains(output, "Home") {
		t.Error("Expected breadcrumbs to be visible")
	}

	// Content should be present
	if !strings.Contains(output, "Test content") {
		t.Error("Expected content to be visible")
	}

	// Help should be present
	if !strings.Contains(output, "q Quit") {
		t.Error("Expected help text to be visible")
	}

	// Footer separator should be present
	if !strings.Contains(output, "─") {
		t.Error("Expected footer separator to be visible")
	}
}

// TestStandardView_MinimumViableTerminal tests rendering at the minimum supported size
func TestStandardView_MinimumViableTerminal(t *testing.T) {
	size := terminalSizeTiny // 80x24
	logo := NewASCIILogo(false, size.Width)
	view := NewStandardView(size).
		WithLogo(logo, 2).
		WithBreadcrumbs("Home", "Settings").
		WithContent("Minimal content").
		WithHelp("q Quit").
		WithFooterSeparator(true)

	output := view.Render()

	// Should handle minimal size gracefully
	if output == "" {
		t.Error("Expected output for minimal terminal size")
	}

	// All essential elements should still be present
	if !strings.Contains(output, "Minimal content") {
		t.Error("Expected content visible at minimal size")
	}
}

// TestStandardView_ContentOverflow tests content that exceeds available space
func TestStandardView_ContentOverflow(t *testing.T) {
	size := terminalSizeTiny // 80x24 - small terminal
	logo := NewASCIILogo(false, size.Width)

	// Create content that exceeds terminal height
	longContent := strings.Repeat("Line of content\n", 50)

	view := NewStandardView(size).
		WithLogo(logo, 2).
		WithContent(longContent).
		WithHelp("q Quit")

	output := view.Render()

	// Should handle overflow gracefully
	if output == "" {
		t.Error("Expected output even with content overflow")
	}
}

// TestModal_AllTerminalSizes tests modal rendering at all common terminal sizes
func TestModal_AllTerminalSizes(t *testing.T) {
	sizes := map[string]*terminal.Info{
		"Tiny (80x24)":       terminalSizeTiny,
		"Compact (100x30)":   terminalSizeCompact,
		"Normal (120x40)":    terminalSizeNormal,
		"Large (160x50)":     terminalSizeLarge,
		"XLarge (200x60)":    terminalSizeXLarge,
		"UltraWide (240x40)": terminalSizeUltraWide,
		"Tall (120x80)":      terminalSizeTall,
	}

	modalTypes := map[string]*ModalContent{
		"Error":    NewErrorModal("Error", "An error occurred"),
		"Loading":  NewLoadingModal("Loading...", false),
		"Progress": NewProgressModal("Processing", "Processing items...", 0.5),
		"Success":  NewSuccessModal("Operation completed successfully"),
		"Warning":  NewWarningModal("Warning", "Please be careful"),
	}

	for sizeName, size := range sizes {
		for modalName, modal := range modalTypes {
			t.Run(sizeName+"_"+modalName, func(t *testing.T) {
				output := modal.Render(size.Width, size.Height)

				if output == "" {
					t.Errorf("Expected modal output for %s at size %s", modalName, sizeName)
				}
			})
		}
	}
}

// TestModal_ErrorModalSizing tests error modal sizing at different terminal sizes
func TestModal_ErrorModalSizing(t *testing.T) {
	modal := NewErrorModal("Error", "This is an error message that should be displayed properly")

	sizes := []*terminal.Info{
		terminalSizeTiny,
		terminalSizeNormal,
		terminalSizeXLarge,
	}

	for _, size := range sizes {
		output := modal.Render(size.Width, size.Height)

		if output == "" {
			t.Errorf("Expected error modal output for %dx%d", size.Width, size.Height)
		}

		// Should contain error indicator
		if !strings.Contains(output, "Error") && !strings.Contains(output, "✖") {
			t.Errorf("Expected error indicator at size %dx%d", size.Width, size.Height)
		}
	}
}

// TestModal_ProgressModalSizing tests progress modal sizing at different terminal sizes
func TestModal_ProgressModalSizing(t *testing.T) {
	modal := NewProgressModal("Processing", "Processing your request...", 0.5)

	sizes := []*terminal.Info{
		terminalSizeTiny,
		terminalSizeNormal,
		terminalSizeXLarge,
	}

	for _, size := range sizes {
		output := modal.Render(size.Width, size.Height)

		if output == "" {
			t.Errorf("Expected progress modal output for %dx%d", size.Width, size.Height)
		}

		// Should contain progress indicator
		if !strings.Contains(output, "50%") {
			t.Errorf("Expected progress percentage at size %dx%d", size.Width, size.Height)
		}
	}
}

// TestModal_SuccessModalSizing tests success modal sizing at different terminal sizes
func TestModal_SuccessModalSizing(t *testing.T) {
	modal := NewSuccessModal("Operation completed successfully")

	sizes := []*terminal.Info{
		terminalSizeTiny,
		terminalSizeNormal,
		terminalSizeXLarge,
	}

	for _, size := range sizes {
		output := modal.Render(size.Width, size.Height)

		if output == "" {
			t.Errorf("Expected success modal output for %dx%d", size.Width, size.Height)
		}

		// Should contain success indicator
		if !strings.Contains(output, "Success") && !strings.Contains(output, "✓") {
			t.Errorf("Expected success indicator at size %dx%d", size.Width, size.Height)
		}
	}
}

// TestStandardViewWithModal_AllSizes tests StandardView with modal overlay at all sizes
func TestStandardViewWithModal_AllSizes(t *testing.T) {
	sizes := []*terminal.Info{
		terminalSizeTiny,
		terminalSizeCompact,
		terminalSizeNormal,
		terminalSizeLarge,
		terminalSizeXLarge,
	}

	for _, size := range sizes {
		t.Run("ErrorModal", func(t *testing.T) {
			testStandardViewWithModalAtSize(t, size, NewErrorModal("Error", "Test error"))
		})

		t.Run("LoadingModal", func(t *testing.T) {
			testStandardViewWithModalAtSize(t, size, NewLoadingModal("Loading...", false))
		})

		t.Run("ProgressModal", func(t *testing.T) {
			testStandardViewWithModalAtSize(t, size, NewProgressModal("Processing", "Processing...", 0.75))
		})
	}
}

func testStandardViewWithModalAtSize(t *testing.T, size *terminal.Info, modal *ModalContent) {
	logo := NewASCIILogo(false, size.Width)
	view := NewStandardView(size).
		WithLogo(logo, 2).
		WithContent("Background content").
		WithHelp("q Quit").
		ShowModalOverlay(modal)

	output := view.Render()

	if output == "" {
		t.Errorf("Expected output for StandardView with modal at %dx%d", size.Width, size.Height)
	}

	// Modal should be visible
	// Background content should still be rendered (but may be dimmed)
}

// TestResponsiveLayout tests that layout responds to terminal size changes
func TestResponsiveLayout(t *testing.T) {
	logo := NewASCIILogo(false, 80)

	// Start with small terminal
	view := NewStandardView(terminalSizeTiny).
		WithLogo(logo, 2).
		WithContent("Content").
		WithHelp("Help")

	output1 := view.Render()

	// Resize to large terminal
	logo.SetWidth(160)
	view.TerminalInfo = terminalSizeLarge

	output2 := view.Render()

	// Both should produce valid output
	if output1 == "" || output2 == "" {
		t.Error("Expected valid output at both sizes")
	}

	// Outputs should differ due to different terminal sizes
	// (This is a basic check - actual content may vary)
}

// TestLayoutBreaking tests for potential layout breaking at edge sizes
func TestLayoutBreaking(t *testing.T) {
	// Test at various sizes to ensure no layout breaking
	testSizes := []*terminal.Info{
		{Width: 40, Height: 10},   // Very small
		{Width: 80, Height: 24},   // Minimum
		{Width: 100, Height: 30},  // Compact
		{Width: 120, Height: 40},  // Normal
		{Width: 200, Height: 60},  // Large
		{Width: 300, Height: 100}, // Very large
	}

	for _, size := range testSizes {
		t.Run("Size_"+string(rune(size.Width))+"x"+string(rune(size.Height)), func(t *testing.T) {
			logo := NewASCIILogo(false, size.Width)
			view := NewStandardView(size).
				WithLogo(logo, 2).
				WithBreadcrumbs("A", "B", "C", "D", "E").
				WithContent("Test content with multiple lines\nLine 2\nLine 3").
				WithHelp("Very long help text that might overflow").
				WithFooterSeparator(true)

			// Should not panic or crash
			output := view.Render()

			if output == "" {
				t.Errorf("Expected output at size %dx%d", size.Width, size.Height)
			}
		})
	}
}
