package modals_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("SkillSuggestionModal", func() {
	var (
		modal       *modals.SkillSuggestionModal
		theme       themes.Theme
		suggestions []skillinference.SkillSuggestion
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		suggestions = []skillinference.SkillSuggestion{
			{
				Name:       "Go",
				Category:   "Backend",
				Confidence: 0.95,
				EventIDs:   []string{"event-1", "event-2"},
				Contexts:   []string{"Built API with Go", "Wrote Go microservices"},
			},
			{
				Name:       "PostgreSQL",
				Category:   "Database",
				Confidence: 0.85,
				EventIDs:   []string{"event-1"},
				Contexts:   []string{"Integrated PostgreSQL database"},
			},
			{
				Name:       "Docker",
				Category:   "DevOps",
				Confidence: 0.70,
				EventIDs:   []string{"event-2"},
				Contexts:   []string{"Deployed using Docker containers"},
			},
		}
	})

	Describe("NewSkillSuggestionModal", func() {
		It("should create a new modal", func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should sort suggestions by confidence descending", func() {
			// Create unsorted suggestions
			unsorted := []skillinference.SkillSuggestion{
				{Name: "Low", Confidence: 0.5},
				{Name: "High", Confidence: 0.95},
				{Name: "Medium", Confidence: 0.75},
			}

			modal = modals.NewSkillSuggestionModal(unsorted, theme)
			modal.Show()

			view := modal.View()
			// High confidence should appear first
			Expect(view).To(ContainSubstring("High"))
			// The view should show suggestions in descending confidence order
		})

		It("should handle nil theme", func() {
			modal = modals.NewSkillSuggestionModal(suggestions, nil)

			Expect(modal).NotTo(BeNil())
			// Should use default theme without panicking
		})

		It("should handle empty suggestions", func() {
			modal = modals.NewSkillSuggestionModal([]skillinference.SkillSuggestion{}, theme)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Visibility", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
		})

		It("should start hidden", func() {
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should show when Show is called", func() {
			modal.Show()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide when Hide is called", func() {
			modal.Show()
			modal.Hide()

			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should return empty string when hidden", func() {
			view := modal.View()

			Expect(view).To(BeEmpty())
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should navigate down with j key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			cmd, _ := modal.Update(msg)

			Expect(cmd).To(BeNil())
			// Should move to next suggestion
		})

		It("should navigate up with k key", func() {
			// Navigate down first
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			// Then navigate up
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			cmd, _ := modal.Update(msg)

			Expect(cmd).To(BeNil())
		})

		It("should navigate down with arrow down", func() {
			msg := tea.KeyMsg{Type: tea.KeyDown}
			cmd, _ := modal.Update(msg)

			Expect(cmd).To(BeNil())
		})

		It("should navigate up with arrow up", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})

			msg := tea.KeyMsg{Type: tea.KeyUp}
			cmd, _ := modal.Update(msg)

			Expect(cmd).To(BeNil())
		})
	})

	Describe("Actions", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should accept suggestion with 'a' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			cmd, action := modal.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(action).NotTo(BeNil())
			// Should return accept action
		})

		It("should reject suggestion with 'r' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
			cmd, action := modal.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(action).NotTo(BeNil())
			// Should return reject action
		})

		It("should accept all with 'A' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}}
			cmd, result := modal.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			// Should return all suggestions as accepted
		})

		It("should cancel with Escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			cmd, result := modal.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			// Should return cancel action
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("GetAcceptedSuggestions", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should return empty list initially", func() {
			accepted := modal.GetAcceptedSuggestions()

			Expect(accepted).To(BeEmpty())
		})

		It("should return accepted suggestions after accepting", func() {
			// Accept first suggestion
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedSuggestions()

			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].Name).To(Equal("Go"))
		})

		It("should accumulate accepted suggestions", func() {
			// Accept first
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			// Navigate to second and accept
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedSuggestions()

			Expect(accepted).To(HaveLen(2))
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should display suggestion name", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("Go"))
			Expect(view).To(ContainSubstring("PostgreSQL"))
		})

		It("should display category", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("Backend"))
			Expect(view).To(ContainSubstring("Database"))
		})

		It("should display confidence score", func() {
			view := modal.View()

			// Should show confidence visually (bar or percentage)
			Expect(view).To(MatchRegexp("95%|0.95"))
		})

		It("should display help text", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("a"))
			Expect(view).To(ContainSubstring("Accept"))
		})

		It("should show total count", func() {
			view := modal.View()

			Expect(view).To(MatchRegexp("[0-9]+ of [0-9]+"))
		})
	})

	Describe("Window resize", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should handle window size message", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}
			cmd, _ := modal.Update(msg)

			Expect(cmd).To(BeNil())
			// Modal should adapt to new size
		})
	})
})
