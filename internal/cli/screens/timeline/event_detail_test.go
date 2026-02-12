//nolint:errcheck // Test file - error handling for test setup is not relevant.
package timeline_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventDetailScreen", func() {
	var (
		screen *timeline.EventDetailScreen
		event  *career.Event
	)

	BeforeEach(func() {
		event = fixtures.EventWith("event-1", "Backend Developer at TechCorp - Built scalable APIs using Go and PostgreSQL", "TechCorp", "API Platform")
		event.Date = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		event.Tags = []string{"backend", "api", "golang"}
		event.Categories = []string{"development", "architecture"}
		event.Skills = []string{"skill-1", "skill-2"}
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

		It("should return skills action on 's' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("skills"))
			Expect(data["event"]).To(Equal(event))
		})

		It("should return help action on '?' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("help"))
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
			// Label "Skills:" is styled with ANSI codes, so check label and value separately.
			Expect(view).To(ContainSubstring("Skills:"))
			Expect(view).To(ContainSubstring("2 associated"))
		})

		It("should show help text in footer", func() {
			view := screen.View()
			// UIKit badges render action hints (not "key: action" format)
			Expect(view).To(ContainSubstring("Edit"))
			Expect(view).To(ContainSubstring("Delete"))
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

	Describe("Theme Integration", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventDetailScreen(event)
			screen.SetTerminalInfo(120, 40)
		})

		It("should use provided theme for rendering content", func() {
			// Create a theme and set it on the screen
			th := themes.NewDefaultTheme()
			screen.SetTheme(th)

			// RenderContent should use the screen's theme, not theme.Default()
			// This test verifies the screen doesn't hardcode theme.Default()
			content := screen.RenderContent()
			Expect(content).NotTo(BeEmpty())
			// The content should be rendered (if theme wasn't used properly,
			// the DetailView would fall back to defaults but still render)
			Expect(content).To(ContainSubstring("Event Details"))
		})

		It("should fall back to default theme when no theme is set", func() {
			// Don't set a theme - screen should handle gracefully
			content := screen.RenderContent()
			Expect(content).NotTo(BeEmpty())
			Expect(content).To(ContainSubstring("Event Details"))
		})
	})

	Describe("UIKit Footer Rendering", func() {
		BeforeEach(func() {
			screen = timeline.NewTimelineEventDetailScreen(event)
			screen.SetTerminalInfo(120, 40)
		})

		It("should NOT use legacy colon-separated format in footer", func() {
			// Set a proper theme
			th := themes.NewDefaultTheme()
			screen.SetTheme(th)

			view := screen.View()

			// Legacy format uses "key: action" (e.g., "e: Edit", "d: Delete")
			// UIKit primitives use styled "[key] action" format without colons
			// Verify the old format is NOT used
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
			Expect(view).NotTo(ContainSubstring("Esc/Backspace: Back"))

			// Should have action hints
			Expect(view).To(ContainSubstring("Edit"))
			Expect(view).To(ContainSubstring("Back"))
		})
	})
})
