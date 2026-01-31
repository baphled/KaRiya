//nolint:errcheck // Test file - error handling for test setup is not relevant.
package timeline_test

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventListScreen", func() {
	var (
		screen *timeline.EventListScreen
		events []*career.Event
	)

	BeforeEach(func() {
		events = []*career.Event{
			{
				ID:      "event-1",
				Date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Text:    "Backend Developer at TechCorp - Built scalable APIs",
				Company: "TechCorp",
			},
			{
				ID:      "event-2",
				Date:    time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
				Text:    "DevOps Engineer at CloudInc - Managed Kubernetes clusters",
				Company: "CloudInc",
			},
			{
				ID:      "event-3",
				Date:    time.Date(2022, 3, 10, 0, 0, 0, 0, time.UTC),
				Text:    "Frontend Developer at WebSolutions - Created React applications",
				Company: "WebSolutions",
			},
		}
	})

	Describe("Construction", func() {
		It("should create an event list screen", func() {
			screen = timeline.NewTimelineEventListScreen(events)
			Expect(screen).NotTo(BeNil())
		})

		It("should store events", func() {
			screen = timeline.NewTimelineEventListScreen(events)
			Expect(screen.GetEvents()).To(Equal(events))
		})

		It("should handle empty event list", func() {
			screen = timeline.NewTimelineEventListScreen([]*career.Event{})
			Expect(screen).NotTo(BeNil())
			Expect(screen.GetEvents()).To(BeEmpty())
		})

		It("should start with first item selected", func() {
			screen = timeline.NewTimelineEventListScreen(events)
			Expect(screen.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventListScreen(events)
			screen.SetTerminalInfo(120, 40)
		})

		It("should move down with arrow key", func() {
			msg := tea.KeyMsg{Type: tea.KeyDown}
			screen.Update(msg)
			Expect(screen.GetSelectedIndex()).To(Equal(1))
		})

		It("should move up with arrow key", func() {
			// Move to second item first
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(screen.GetSelectedIndex()).To(Equal(1))

			// Move back up
			msg := tea.KeyMsg{Type: tea.KeyUp}
			screen.Update(msg)
			Expect(screen.GetSelectedIndex()).To(Equal(0))
		})

		It("should move down with j key (vim)", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			screen.Update(msg)
			Expect(screen.GetSelectedIndex()).To(Equal(1))
		})

		It("should move up with k key (vim)", func() {
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			screen.Update(msg)
			Expect(screen.GetSelectedIndex()).To(Equal(0))
		})

		It("should not go below first item", func() {
			msg := tea.KeyMsg{Type: tea.KeyUp}
			screen.Update(msg)
			Expect(screen.GetSelectedIndex()).To(Equal(0))
		})

		It("should not go beyond last item", func() {
			// Move to last item
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(screen.GetSelectedIndex()).To(Equal(2))

			// Try to move beyond
			msg := tea.KeyMsg{Type: tea.KeyDown}
			screen.Update(msg)
			Expect(screen.GetSelectedIndex()).To(Equal(2))
		})
	})

	Describe("Selection", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventListScreen(events)
			screen.SetTerminalInfo(120, 40)
		})

		It("should return NavigateResult with selected event on enter", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.ResultData).To(Equal(events[0]))
		})

		It("should return correct event after navigation", func() {
			// Navigate to second item
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Select it
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(result).NotTo(BeNil())
			navResult := result.(*screens.NavigateResult)
			Expect(navResult.ResultData).To(Equal(events[1]))
		})
	})

	Describe("Actions", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventListScreen(events)
			screen.SetTerminalInfo(120, 40)
		})

		It("should return add action on 'a' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("add"))
		})

		It("should return edit action on 'e' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("edit"))
			Expect(data["event"]).To(Equal(events[0]))
		})

		It("should return delete action on 'd' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("delete"))
			Expect(data["event"]).To(Equal(events[0]))
		})

		It("should return filter action on 'f' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("filter"))
		})

		Context("after navigating to different event", func() {
			BeforeEach(func() {
				// Navigate to second event
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(screen.GetSelectedIndex()).To(Equal(1))
			})

			It("should edit the navigated-to event", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				navResult := result.(*screens.NavigateResult)
				data := navResult.ResultData.(map[string]interface{})
				Expect(data["action"]).To(Equal("edit"))
				Expect(data["event"]).To(Equal(events[1])) // Second event
			})

			It("should delete the navigated-to event", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				navResult := result.(*screens.NavigateResult)
				data := navResult.ResultData.(map[string]interface{})
				Expect(data["action"]).To(Equal("delete"))
				Expect(data["event"]).To(Equal(events[1])) // Second event
			})

			It("should view details of the navigated-to event on Enter", func() {
				_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(result).NotTo(BeNil())
				navResult := result.(*screens.NavigateResult)
				Expect(navResult.ResultData).To(Equal(events[1])) // Second event
			})
		})

		Context("with empty event list", func() {
			BeforeEach(func() {
				screen = timeline.NewTimelineEventListScreen([]*career.Event{})
				screen.SetTerminalInfo(120, 40)
			})

			It("should still allow add action on empty list", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				navResult := result.(*screens.NavigateResult)
				data := navResult.ResultData.(map[string]interface{})
				Expect(data["action"]).To(Equal("add"))
			})

			It("should return nil for edit on empty list", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
				_, result := screen.Update(msg)

				Expect(result).To(BeNil())
			})

			It("should return nil for delete on empty list", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
				_, result := screen.Update(msg)

				Expect(result).To(BeNil())
			})

			It("should return nil for Enter on empty list", func() {
				_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(result).To(BeNil())
			})

			It("should still allow filter action on empty list", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
				_, result := screen.Update(msg)

				Expect(result).NotTo(BeNil())
				navResult := result.(*screens.NavigateResult)
				data := navResult.ResultData.(map[string]interface{})
				Expect(data["action"]).To(Equal("filter"))
			})
		})
	})

	Describe("Cancellation", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventListScreen(events)
		})

		It("should return CancelResult on escape", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should not handle 'q' key (handled by intent)", func() {
			// 'q' is a global key handled by the intent before delegation
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			_, result := screen.Update(msg)

			// Screen should not process 'q', it should return nil
			Expect(result).To(BeNil())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventListScreen(events)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render a view with StandardView structure", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include breadcrumbs", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Timeline"))
		})

		It("should show event text", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Backend Developer"))
			Expect(view).To(ContainSubstring("DevOps Engineer"))
		})

		It("should show event companies", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("CloudInc"))
		})

		It("should show event count", func() {
			view := screen.View()
			// New format: "Events: 3 | Page 1 of 1"
			Expect(view).To(ContainSubstring("Events: 3"))
		})

		It("should show help text in footer", func() {
			view := screen.View()
			// UIKit NavigateBadge uses "↑↓/jk" format
			Expect(view).To(ContainSubstring("↑↓/jk"))
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		Context("when list is empty", func() {
			BeforeEach(func() {
				screen = timeline.NewTimelineEventListScreen([]*career.Event{})
				screen.SetTerminalInfo(120, 40)
			})

			It("should show empty state message", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("No events"))
			})
		})
	})

	Describe("Terminal Handling", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventListScreen(events)
		})

		It("should update dimensions on SetTerminalInfo", func() {
			screen.SetTerminalInfo(100, 30)
			Expect(screen.Width()).To(Equal(100))
			Expect(screen.Height()).To(Equal(30))
		})

		It("should handle WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.Width()).To(Equal(80))
			Expect(screen.Height()).To(Equal(24))
		})
	})

	Describe("Screen Interface", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventListScreen(events)
		})

		It("should implement Screen interface", func() {
			var _ screens.Screen = screen
		})

		It("should support SetTheme", func() {
			theme := "test-theme"
			screen.SetTheme(theme)
			Expect(screen.Theme()).To(Equal(theme))
		})
	})

	Describe("SetContentHeight", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventListScreen(events)
			screen.SetTerminalInfo(140, 40)
		})

		It("should not stretch table to terminal width", func() {
			screen.SetContentHeight(20)
			content := screen.RenderContent()
			lines := strings.Split(content, "\n")

			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					Expect(len(line)).To(BeNumerically("<", 140),
						"table content should be narrower than terminal width for centering")
				}
			}
		})
	})

	Describe("UIKit Footer Rendering", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventListScreen(events)
			screen.SetTerminalInfo(120, 40)
		})

		It("should NOT use legacy colon-separated format in footer", func() {
			// Set a proper theme
			th := themes.NewDefaultTheme()
			screen.SetTheme(th)

			view := screen.View()

			// Legacy format uses "key: action" (e.g., "Enter: View", "a: Add")
			// UIKit primitives use styled "[key] action" format without colons
			// Verify the old format is NOT used
			Expect(view).NotTo(ContainSubstring("Enter: View"))
			Expect(view).NotTo(ContainSubstring("a: Add"))
			Expect(view).NotTo(ContainSubstring("e: Edit"))
			Expect(view).NotTo(ContainSubstring("d: Delete"))
			Expect(view).NotTo(ContainSubstring("q: Quit"))
			Expect(view).NotTo(ContainSubstring("?: Help"))
		})

		It("should render footer with UIKit badges containing action hints", func() {
			th := themes.NewDefaultTheme()
			screen.SetTheme(th)

			view := screen.View()

			// UIKit badges render action hints (the hint part of HelpKeyBadge)
			Expect(view).To(ContainSubstring("Navigate"))
			Expect(view).To(ContainSubstring("View"))
			Expect(view).To(ContainSubstring("Add"))
			Expect(view).To(ContainSubstring("Edit"))
			Expect(view).To(ContainSubstring("Delete"))
			Expect(view).To(ContainSubstring("Back"))
			Expect(view).To(ContainSubstring("Quit"))
			Expect(view).To(ContainSubstring("Help"))
		})

		It("should render footer with default theme when no theme is set", func() {
			// Don't set theme - should fall back to default and still use UIKit
			view := screen.View()

			// Should NOT have legacy colon format
			Expect(view).NotTo(ContainSubstring("Esc: Back"))

			// Should have action hints
			Expect(view).To(ContainSubstring("Navigate"))
			Expect(view).To(ContainSubstring("Back"))
		})

		It("should include page navigation badge", func() {
			th := themes.NewDefaultTheme()
			screen.SetTheme(th)

			view := screen.View()

			// Should have page navigation hint
			Expect(view).To(ContainSubstring("Page"))
		})
	})
})
