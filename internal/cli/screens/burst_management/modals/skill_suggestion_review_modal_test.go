package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SuggestionReviewModal (Skill)", func() {
	var (
		modal       *modals.SuggestionReviewModal
		suggestions []skillinference.SkillSuggestion
		theme       themes.Theme
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		suggestions = []skillinference.SkillSuggestion{
			{
				Name:       "Go",
				Category:   "backend",
				Confidence: 0.92,
				EventIDs:   []string{"e1", "e2", "e3"},
				Contexts:   []string{"Built API in Go", "Goroutine-based worker", "Go testing patterns"},
			},
			{
				Name:       "PostgreSQL",
				Category:   "database",
				Confidence: 0.78,
				EventIDs:   []string{"e4", "e5"},
				Contexts:   []string{"Query optimisation"},
			},
			{
				Name:       "Docker",
				Category:   "devops",
				Confidence: 0.85,
				EventIDs:   []string{"e6", "e7", "e8"},
				Contexts:   []string{"Container builds", "Multi-stage Dockerfile", "Docker Compose setup", "CI integration"},
			},
		}
	})

	Describe("Creation", func() {
		It("should create modal with skill suggestions", func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			Expect(modal).NotTo(BeNil())
			Expect(modal.HasSuggestions()).To(BeTrue())
			Expect(modal.GetSuggestionsCount()).To(Equal(3))
		})

		It("should create modal with nil theme using default", func() {
			modal = modals.NewSkillSuggestionModal(suggestions, nil)
			Expect(modal).NotTo(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should sort by confidence descending", func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
			all := modal.GetAllSkills()
			Expect(all[0].Name).To(Equal("Go"))
			Expect(all[1].Name).To(Equal("Docker"))
			Expect(all[2].Name).To(Equal("PostgreSQL"))
		})

		It("should initialize accepted skills as empty", func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			Expect(modal.GetAcceptedSkills()).To(BeEmpty())
		})

		It("should handle empty suggestions", func() {
			modal = modals.NewSkillSuggestionModal([]skillinference.SkillSuggestion{}, theme)
			Expect(modal.HasSuggestions()).To(BeFalse())
			Expect(modal.GetSuggestionsCount()).To(Equal(0))
			Expect(modal.GetCurrentSkill()).To(BeNil())
		})
	})

	Describe("Content Rendering", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should show Review Skill Suggestions title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Review Skill Suggestions"))
		})

		It("should display skill names", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Go"))
			Expect(view).To(ContainSubstring("Docker"))
			Expect(view).To(ContainSubstring("PostgreSQL"))
		})

		It("should display confidence scores", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("92%"))
			Expect(view).To(ContainSubstring("85%"))
			Expect(view).To(ContainSubstring("78%"))
		})

		It("should show selected skill detail with category badge", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Selected:"))
			Expect(view).To(ContainSubstring("backend"))
		})

		It("should show usage contexts for selected skill", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Usage Contexts:"))
			Expect(view).To(ContainSubstring("Built API in Go"))
		})

		It("should show overflow indicator for 4+ contexts", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			view := modal.View()
			Expect(view).To(ContainSubstring("and 1 more"))
		})

		It("should handle empty skills with message", func() {
			modal = modals.NewSkillSuggestionModal([]skillinference.SkillSuggestion{}, theme)
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("No skills detected"))
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should navigate with j/k keys", func() {
			Expect(modal.GetCurrentSkill().Name).To(Equal("Go"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSkill().Name).To(Equal("Docker"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(modal.GetCurrentSkill().Name).To(Equal("Go"))
		})

		It("should navigate with arrow keys", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(modal.GetCurrentSkill().Name).To(Equal("Docker"))

			modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(modal.GetCurrentSkill().Name).To(Equal("Go"))
		})
	})

	Describe("Accept", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should accept current skill with 'a' key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedSkills()
			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].Name).To(Equal("Go"))
		})

		It("should remove accepted skill from pending list", func() {
			Expect(modal.GetSuggestionsCount()).To(Equal(3))
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetSuggestionsCount()).To(Equal(2))
		})

		It("should close modal when last skill accepted", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Reject", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should reject current skill with 'r' key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			Expect(modal.GetAcceptedSkills()).To(BeEmpty())
			Expect(modal.GetCurrentSkill().Name).To(Equal("Docker"))
		})

		It("should close modal when last skill rejected", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.HasSuggestions()).To(BeFalse())
		})
	})

	Describe("Mixed Accept/Reject", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should handle accept, navigate, reject sequence", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			Expect(modal.GetAcceptedSkills()).To(HaveLen(1))
			Expect(modal.GetSuggestionsCount()).To(Equal(1))
		})

		It("should adjust index when accepting last item in list", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSkill().Name).To(Equal("PostgreSQL"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetCurrentSkill().Name).To(Equal("Docker"))
		})

		It("should adjust index when rejecting last item in list", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.GetCurrentSkill().Name).To(Equal("Docker"))
		})
	})

	Describe("Skill with no contexts", func() {
		It("should render without contexts section", func() {
			noCtxSuggestions := []skillinference.SkillSuggestion{
				{
					Name:       "Rust",
					Category:   "systems",
					Confidence: 0.70,
					EventIDs:   []string{"e1"},
					Contexts:   []string{},
				},
			}
			modal = modals.NewSkillSuggestionModal(noCtxSuggestions, theme)
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("Rust"))
			Expect(view).NotTo(ContainSubstring("Usage Contexts:"))
		})
	})

	Describe("View Events", func() {
		BeforeEach(func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
		})

		It("should set view events action on Enter", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.GetAction()).To(Equal(modals.SuggestionActionViewEvents))
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Dimensions", func() {
		It("should handle SetDimensions", func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.SetDimensions(100, 50)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle small dimensions", func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.SetDimensions(50, 10)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle window size message", func() {
			modal = modals.NewSkillSuggestionModal(suggestions, theme)
			modal.Show()
			_, cmd := modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			_ = cmd
		})
	})
})
