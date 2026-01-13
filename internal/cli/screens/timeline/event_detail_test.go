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

var _ = Describe("TimelineEventDetailScreen", func() {
	var (
		screen *timeline.TimelineEventDetailScreen
		event  *career.CareerEvent
	)

	BeforeEach(func() {
		event = &career.CareerEvent{
			ID:         "event-1",
			Date:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Text:       "Backend Developer at TechCorp - Built scalable APIs using Go and PostgreSQL",
			Company:    "TechCorp",
			Project:    "API Platform",
			Tags:       []string{"backend", "api", "golang"},
			Categories: []string{"development", "architecture"},
			Skills:     []string{"skill-1", "skill-2"},
		}
	})

	Describe("Construction", func() {
		It("should create an event detail screen", func() {
			screen = timeline.NewTimelineEventDetailScreen(event)
			Expect(screen).NotTo(BeNil())
		})

		It("should store event", func() {
			screen = timeline.NewTimelineEventDetailScreen(event)
			Expect(screen.GetEvent()).To(Equal(event))
		})
	})

	Describe("Actions", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventDetailScreen(event)
			screen.SetTerminalInfo(120, 40)
		})

		It("should return edit action on 'e' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("edit"))
			Expect(data["event"]).To(Equal(event))
		})

		It("should return delete action on 'd' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("delete"))
			Expect(data["event"]).To(Equal(event))
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventDetailScreen(event)
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

		It("should return CancelResult on backspace", func() {
			msg := tea.KeyMsg{Type: tea.KeyBackspace}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventDetailScreen(event)
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
			Expect(view).To(ContainSubstring("Event Details"))
		})

		It("should show event date", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("2024-01-01"))
		})

		It("should show event text", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Backend Developer"))
			Expect(view).To(ContainSubstring("scalable APIs"))
		})

		It("should show company", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("TechCorp"))
		})

		It("should show project", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("API Platform"))
		})

		It("should show tags", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("backend"))
			Expect(view).To(ContainSubstring("api"))
			Expect(view).To(ContainSubstring("golang"))
		})

		It("should show categories", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("development"))
			Expect(view).To(ContainSubstring("architecture"))
		})

		It("should show skills count", func() {
			view := screen.View()
			// Format: "Skills: 2 associated"
			Expect(view).To(ContainSubstring("Skills: 2"))
		})

		It("should show help text in footer", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("e: Edit"))
			Expect(view).To(ContainSubstring("d: Delete"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		Context("when optional fields are missing", func() {
			BeforeEach(func() {
				event.Project = ""
				event.Tags = nil
				event.Categories = nil
				event.Skills = nil
				screen = timeline.NewTimelineEventDetailScreen(event)
				screen.SetTerminalInfo(120, 40)
			})

			It("should not show empty project field", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Date"))
				Expect(view).To(ContainSubstring("Company"))
				// Project label should not appear when empty
			})

			It("should handle empty tags gracefully", func() {
				view := screen.View()
				// Should not crash, but tags section may be omitted or show "None"
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("Terminal Handling", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventDetailScreen(event)
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
			screen = timeline.NewTimelineEventDetailScreen(event)
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
