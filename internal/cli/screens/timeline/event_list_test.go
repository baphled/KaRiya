package timeline_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimelineEventListScreen", func() {
	var (
		screen *timeline.TimelineEventListScreen
		events []*career.CareerEvent
	)

	BeforeEach(func() {
		events = []*career.CareerEvent{
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
			screen = timeline.NewTimelineEventListScreen([]*career.CareerEvent{})
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

		It("should return CancelResult on 'q' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
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
			Expect(view).To(ContainSubstring("3 events"))
		})

		It("should show help text in footer", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("↑/↓"))
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		Context("when list is empty", func() {
			BeforeEach(func() {
				screen = timeline.NewTimelineEventListScreen([]*career.CareerEvent{})
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
})
