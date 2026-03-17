//nolint:errcheck // Test file - error handling for test setup is not relevant.
package event_test

import (
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/types"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StrategySelect", func() {
	var view *event.StrategySelect

	Describe("Construction", func() {
		It("should create a non-nil StrategySelect", func() {
			view = event.NewStrategySelect()
			Expect(view).NotTo(BeNil())
		})

		It("should have two items (Quick and Manual strategies)", func() {
			view = event.NewStrategySelect()
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Quick"))
			Expect(content).To(ContainSubstring("Manual"))
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			view = event.NewStrategySelect()
			cmd := view.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("WithInitialSelection", func() {
		It("should set selection to valid index", func() {
			view = event.NewStrategySelect()
			view.WithInitialSelection(1)
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(widgets.ResultNavigate))
			navResult, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			strategy, ok := navResult.ResultData.(types.CaptureStrategy)
			Expect(ok).To(BeTrue())
			Expect(strategy).To(Equal(types.StrategyManual))
		})

		It("should clamp negative index to 0", func() {
			view = event.NewStrategySelect()
			view.WithInitialSelection(-5)
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			navResult, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			strategy, ok := navResult.ResultData.(types.CaptureStrategy)
			Expect(ok).To(BeTrue())
			Expect(strategy).To(Equal(types.StrategyQuick))
		})

		It("should clamp too-large index to last item", func() {
			view = event.NewStrategySelect()
			view.WithInitialSelection(100)
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			navResult, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			strategy, ok := navResult.ResultData.(types.CaptureStrategy)
			Expect(ok).To(BeTrue())
			Expect(strategy).To(Equal(types.StrategyManual))
		})
	})

	Describe("Update with WindowSizeMsg", func() {
		BeforeEach(func() {
			view = event.NewStrategySelect()
		})

		It("should update terminal dimensions", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("Update with KeyMsg", func() {
		BeforeEach(func() {
			view = event.NewStrategySelect()
			view.SetTerminalInfo(120, 40)
		})

		Describe("Navigation keys", func() {
			It("should move down with down key", func() {
				msg := tea.KeyMsg{Type: tea.KeyDown}
				_, result := view.Update(msg)
				Expect(result).To(BeNil())
				msg = tea.KeyMsg{Type: tea.KeyEnter}
				_, result = view.Update(msg)
				Expect(result).NotTo(BeNil())
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				strategy, ok := navResult.ResultData.(types.CaptureStrategy)
				Expect(ok).To(BeTrue())
				Expect(strategy).To(Equal(types.StrategyManual))
			})

			It("should move up with up key", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyDown})
				msg := tea.KeyMsg{Type: tea.KeyUp}
				_, result := view.Update(msg)
				Expect(result).To(BeNil())
				msg = tea.KeyMsg{Type: tea.KeyEnter}
				_, result = view.Update(msg)
				Expect(result).NotTo(BeNil())
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				strategy, ok := navResult.ResultData.(types.CaptureStrategy)
				Expect(ok).To(BeTrue())
				Expect(strategy).To(Equal(types.StrategyQuick))
			})

			It("should move down with j key (vim)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
				_, result := view.Update(msg)
				Expect(result).To(BeNil())
				msg = tea.KeyMsg{Type: tea.KeyEnter}
				_, result = view.Update(msg)
				Expect(result).NotTo(BeNil())
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				strategy, ok := navResult.ResultData.(types.CaptureStrategy)
				Expect(ok).To(BeTrue())
				Expect(strategy).To(Equal(types.StrategyManual))
			})

			It("should move up with k key (vim)", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
				_, result := view.Update(msg)
				Expect(result).To(BeNil())
				msg = tea.KeyMsg{Type: tea.KeyEnter}
				_, result = view.Update(msg)
				Expect(result).NotTo(BeNil())
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				strategy, ok := navResult.ResultData.(types.CaptureStrategy)
				Expect(ok).To(BeTrue())
				Expect(strategy).To(Equal(types.StrategyQuick))
			})

			It("should jump to top with g key", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyDown})
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
				_, result := view.Update(msg)
				Expect(result).To(BeNil())
				msg = tea.KeyMsg{Type: tea.KeyEnter}
				_, result = view.Update(msg)
				Expect(result).NotTo(BeNil())
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				strategy, ok := navResult.ResultData.(types.CaptureStrategy)
				Expect(ok).To(BeTrue())
				Expect(strategy).To(Equal(types.StrategyQuick))
			})

			It("should jump to bottom with G key", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
				_, result := view.Update(msg)
				Expect(result).To(BeNil())
				msg = tea.KeyMsg{Type: tea.KeyEnter}
				_, result = view.Update(msg)
				Expect(result).NotTo(BeNil())
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				strategy, ok := navResult.ResultData.(types.CaptureStrategy)
				Expect(ok).To(BeTrue())
				Expect(strategy).To(Equal(types.StrategyManual))
			})
		})

		Describe("Enter key", func() {
			It("should return NavigateViewResult with selected strategy at index 0", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				strategy, ok := navResult.ResultData.(types.CaptureStrategy)
				Expect(ok).To(BeTrue())
				Expect(strategy).To(Equal(types.StrategyQuick))
			})

			It("should return NavigateViewResult with selected strategy at index 1", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyDown})
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				navResult, ok := result.(*widgets.NavigateViewResult)
				Expect(ok).To(BeTrue())
				strategy, ok := navResult.ResultData.(types.CaptureStrategy)
				Expect(ok).To(BeTrue())
				Expect(strategy).To(Equal(types.StrategyManual))
			})
		})

		Describe("Escape key", func() {
			It("should return CancelViewResult", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultCancel))
				_, ok := result.(*widgets.CancelViewResult)
				Expect(ok).To(BeTrue())
			})
		})

		Describe("Unknown keys", func() {
			It("should return nil for unknown key", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}}
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("RenderContent", func() {
		BeforeEach(func() {
			view = event.NewStrategySelect()
			view.SetTerminalInfo(120, 40)
		})

		It("should return non-empty string", func() {
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should contain Select Capture Strategy title", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Select Capture Strategy"))
		})

		It("should contain Quick text", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Quick"))
		})

		It("should contain Manual text", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Manual"))
		})

		It("should show selection indicator for selected item", func() {
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("▶"))
		})
	})

	Describe("HelpText", func() {
		BeforeEach(func() {
			view = event.NewStrategySelect()
		})

		It("should return non-empty string", func() {
			Expect(view.HelpText()).NotTo(BeEmpty())
		})

		It("should contain Navigate text", func() {
			Expect(view.HelpText()).To(ContainSubstring("Navigate"))
		})

		It("should contain Select text", func() {
			Expect(view.HelpText()).To(ContainSubstring("Select"))
		})

		It("should contain Back text", func() {
			Expect(view.HelpText()).To(ContainSubstring("Back"))
		})

		It("should contain Quit text", func() {
			Expect(view.HelpText()).To(ContainSubstring("Quit"))
		})

		It("should return minimal help text for a zero-value selector", func() {
			view = &event.StrategySelect{}

			Expect(view.HelpText()).To(Equal("Esc: Back  q: Quit"))
		})
	})

	Describe("Zero-value behavior", func() {
		BeforeEach(func() {
			view = &event.StrategySelect{}
		})

		It("should render an empty-state message", func() {
			Expect(view.RenderContent()).To(ContainSubstring("No items available"))
		})

		It("should not navigate when enter is pressed with no items", func() {
			cmd, result := view.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should allow reselection on an empty selector without panicking", func() {
			Expect(func() {
				view.WithInitialSelection(0)
				view.WithInitialSelection(-1)
			}).NotTo(Panic())
			Expect(view.HelpText()).To(Equal("Esc: Back  q: Quit"))
		})

		It("should ignore navigation keys when no items are available", func() {
			for _, msg := range []tea.KeyMsg{
				{Type: tea.KeyDown},
				{Type: tea.KeyUp},
				{Type: tea.KeyRunes, Runes: []rune{'g'}},
				{Type: tea.KeyRunes, Runes: []rune{'G'}},
			} {
				cmd, result := view.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			}

			Expect(view.RenderContent()).To(ContainSubstring("No items available"))
		})
	})

	Describe("Scroll behaviour with small terminal", func() {
		BeforeEach(func() {
			view = event.NewStrategySelect()
			view.SetTerminalInfo(80, 10)
			msg := tea.WindowSizeMsg{Width: 80, Height: 10}
			view.Update(msg)
		})

		It("should handle navigation within visible items", func() {
			view.Update(tea.KeyMsg{Type: tea.KeyDown})
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := view.Update(msg)
			Expect(result).NotTo(BeNil())
			navResult, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			strategy, ok := navResult.ResultData.(types.CaptureStrategy)
			Expect(ok).To(BeTrue())
			Expect(strategy).To(Equal(types.StrategyManual))
		})

		It("should clamp navigation at bottom", func() {
			view.Update(tea.KeyMsg{Type: tea.KeyDown})
			view.Update(tea.KeyMsg{Type: tea.KeyDown})
			view.Update(tea.KeyMsg{Type: tea.KeyDown})
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := view.Update(msg)
			navResult, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			strategy, ok := navResult.ResultData.(types.CaptureStrategy)
			Expect(ok).To(BeTrue())
			Expect(strategy).To(Equal(types.StrategyManual))
		})

		It("should clamp navigation at top", func() {
			view.Update(tea.KeyMsg{Type: tea.KeyUp})
			view.Update(tea.KeyMsg{Type: tea.KeyUp})
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := view.Update(msg)
			navResult, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			strategy, ok := navResult.ResultData.(types.CaptureStrategy)
			Expect(ok).To(BeTrue())
			Expect(strategy).To(Equal(types.StrategyQuick))
		})
	})

	Describe("Window resize updates visible items", func() {
		It("should update visible items on resize", func() {
			view = event.NewStrategySelect()
			msg := tea.WindowSizeMsg{Width: 120, Height: 50}
			view.Update(msg)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should handle very small terminal height", func() {
			view = event.NewStrategySelect()
			msg := tea.WindowSizeMsg{Width: 40, Height: 5}
			view.Update(msg)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})
})
