package capture_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/capture"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCaptureScreens(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Capture Screens Suite")
}

// StrategySelectScreen Tests
//
// StrategySelectScreen allows users to choose between Quick and Manual capture strategies.
//
// Related:
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration)
// - internal/cli/intents/capture_event.go (CaptureStrategy constants)
// - internal/cli/screens/base/select_screen.go (BaseSelectScreen[T])

var _ = Describe("StrategySelectScreen", func() {
	var (
		screen *capture.StrategySelectScreen
	)

	BeforeEach(func() {
		// Create screen with standard breadcrumbs
		breadcrumbs := []string{"Main Menu", "Capture Event"}
		screen = capture.NewStrategySelectScreen(breadcrumbs)
	})

	Describe("Creation", func() {
		It("should create with default selection at index 0 (Quick)", func() {
			Expect(screen).NotTo(BeNil())
			// View should show Quick strategy selected
			view := screen.View()
			Expect(view).To(ContainSubstring("▶"))
			Expect(view).To(ContainSubstring("Quick"))
		})

		It("should have exactly 2 strategies (Quick and Manual)", func() {
			// Navigate down once, should be at Manual (index 1)
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			view := screen.View()
			Expect(view).To(ContainSubstring("Manual"))

			// Navigate down again, should still be at Manual (boundary check)
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			view = screen.View()
			Expect(view).To(ContainSubstring("Manual"))
		})

		It("should display strategy descriptions", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Capture with minimal fields"))
			Expect(view).To(ContainSubstring("Full form with optional fields"))
		})

		It("should initialize with breadcrumbs", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Capture Event"))
		})
	})

	Describe("Navigation - Arrow Keys", func() {
		It("should navigate down with Down arrow key", func() {
			// Start at index 0 (Quick)
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			// Should now be at index 1 (Manual)
			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Manual`))
		})

		It("should navigate up with Up arrow key", func() {
			// Navigate to Manual first
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Navigate back up
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			// Should be back at Quick
			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Quick`))
		})

		It("should not go above first item (boundary check)", func() {
			// Try to navigate up from index 0
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			// Should still be at Quick
			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Quick`))
		})

		It("should not go below last item (boundary check)", func() {
			// Navigate to Manual (index 1)
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Try to navigate down from last item
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			// Should still be at Manual
			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Manual`))
		})
	})

	Describe("Navigation - Vim Keys", func() {
		It("should navigate down with 'j' key", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Manual`))
		})

		It("should navigate up with 'k' key", func() {
			// Navigate to Manual first
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Navigate back up with k
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Quick`))
		})
	})

	Describe("Selection - Enter Key", func() {
		It("should return NavigateResult when Enter is pressed on Quick", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
		})

		It("should include Quick strategy in result data", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())

			navResult, ok := result.(*screens.NavigateResult)
			Expect(ok).To(BeTrue())
			Expect(navResult.Data()).To(Equal(intents.StrategyQuick))
		})

		It("should include Manual strategy in result data when selected", func() {
			// Navigate to Manual
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Select Manual
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())

			navResult, ok := result.(*screens.NavigateResult)
			Expect(ok).To(BeTrue())
			Expect(navResult.Data()).To(Equal(intents.StrategyManual))
		})

		It("should include selected index in metadata", func() {
			// Navigate to Manual (index 1)
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Select Manual
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).NotTo(BeNil())

			navResult, ok := result.(*screens.NavigateResult)
			Expect(ok).To(BeTrue())

			metadata := navResult.Metadata()
			index, exists := metadata["selected_index"]
			Expect(exists).To(BeTrue())
			Expect(index).To(Equal(1))
		})

		It("should not return result for navigation keys", func() {
			// Up/Down/j/k should return nil result
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(result).To(BeNil())

			_, result = screen.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(result).To(BeNil())

			_, result = screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(result).To(BeNil())

			_, result = screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(result).To(BeNil())
		})
	})

	Describe("Cancellation - Escape Key", func() {
		It("should return CancelResult when Escape is pressed", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should include current index in cancel result metadata", func() {
			// Navigate to Manual
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Press Escape
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())

			cancelResult, ok := result.(*screens.CancelResult)
			Expect(ok).To(BeTrue())

			metadata := cancelResult.Metadata()
			index, exists := metadata["selected_index"]
			Expect(exists).To(BeTrue())
			Expect(index).To(Equal(1))
		})
	})

	Describe("View Rendering", func() {
		It("should render title", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Select Capture Strategy"))
		})

		It("should highlight selected strategy with ▶", func() {
			// Quick should be selected by default
			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Quick`))
		})

		It("should show footer with navigation hints", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Navigate"))
			Expect(view).To(ContainSubstring("Select"))
			Expect(view).To(ContainSubstring("Back"))
		})

		It("should update highlight when selection changes", func() {
			// Navigate to Manual
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Manual`))
		})
	})

	Describe("Window Resize", func() {
		It("should handle WindowSizeMsg", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			// View should still render correctly
			view := screen.View()
			Expect(view).To(ContainSubstring("Quick"))
		})

		It("should maintain selection after resize", func() {
			// Navigate to Manual
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Resize
			screen.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

			// Selection should be preserved
			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Manual`))
		})
	})

	Describe("State Preservation", func() {
		It("should restore selection with WithInitialSelection", func() {
			// Create screen with initial selection at Manual (index 1)
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewStrategySelectScreen(breadcrumbs).WithInitialSelection(1)

			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Manual`))
		})

		It("should clamp initial selection to valid range", func() {
			// Try to set selection out of bounds
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewStrategySelectScreen(breadcrumbs).WithInitialSelection(99)

			// Should be clamped to last item (Manual, index 1)
			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Manual`))
		})

		It("should handle negative initial selection", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewStrategySelectScreen(breadcrumbs).WithInitialSelection(-5)

			// Should be clamped to first item (Quick, index 0)
			view := screen.View()
			Expect(view).To(MatchRegexp(`▶.*Quick`))
		})
	})
})
