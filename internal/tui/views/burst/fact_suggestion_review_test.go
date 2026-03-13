package burst_test

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	burstview "github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SuggestionReview (Fact)", func() {
	var (
		modal *burstview.SuggestionReview
		facts []display.Fact
		theme themes.Theme
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()

		f1 := fixtures.FactFromBurst("f1", "b1")
		f1.Text = "Led migration of monolith to microservices"
		f1.CompetencyCategories = []string{"technical_leadership"}
		f1.RoleFit = career.RoleFitStaff
		f1.StrengthSignal = "strong"

		f2 := fixtures.FactFromBurst("f2", "b2")
		f2.Text = "Mentored three junior engineers on Go best practices"
		f2.CompetencyCategories = []string{"mentoring", "technical_leadership"}
		f2.RoleFit = career.RoleFitSeniorIC
		f2.StrengthSignal = "moderate"

		f3 := fixtures.Fact("f3", "e1")
		f3.Text = "Designed event-driven architecture for payment processing"
		f3.CompetencyCategories = []string{"architecture"}
		f3.RoleFit = career.RoleFitPrincipal
		f3.StrengthSignal = "strong"

		facts = []display.Fact{
			display.FactFromDomain(f1),
			display.FactFromDomain(f2),
			display.FactFromDomain(f3),
		}
	})

	Describe("Creation", func() {
		It("should create modal with facts and theme", func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			Expect(modal).NotTo(BeNil())
			Expect(modal.HasSuggestions()).To(BeTrue())
			Expect(modal.GetSuggestionsCount()).To(Equal(3))
		})

		It("should create modal with nil theme using default", func() {
			modal = burstview.NewFactSuggestion(facts, nil)
			Expect(modal).NotTo(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should initialize accepted facts as empty", func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			accepted := modal.GetAcceptedFacts()
			Expect(accepted).To(BeEmpty())
		})

		It("should handle empty facts list", func() {
			modal = burstview.NewFactSuggestion([]display.Fact{}, theme)
			Expect(modal.HasSuggestions()).To(BeFalse())
			Expect(modal.GetSuggestionsCount()).To(Equal(0))
			Expect(modal.GetCurrentFact()).To(BeNil())
		})

		It("should handle single fact", func() {
			singleFact := []display.Fact{facts[0]}
			modal = burstview.NewFactSuggestion(singleFact, theme)
			Expect(modal.HasSuggestions()).To(BeTrue())
			Expect(modal.GetSuggestionsCount()).To(Equal(1))
		})

		It("should initialize at first fact", func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()

			current := modal.GetCurrentFact()
			Expect(current).NotTo(BeNil())
			Expect(current.ID).To(Equal("f1"))
		})
	})

	Describe("Content Rendering", func() {
		BeforeEach(func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()
		})

		It("should display fact text in the view", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("monolith"))
		})

		It("should display role fit in the view", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("staff"))
		})

		It("should display strength signal in the view", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("strong"))
		})

		It("should show title as Review Fact Suggestions", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Review Fact Suggestions"))
		})

		It("should show selected fact detail", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Selected:"))
		})

		It("should show footer badges", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Accept"))
			Expect(view).To(ContainSubstring("Reject"))
			Expect(view).To(ContainSubstring("Navigate"))
			Expect(view).To(ContainSubstring("Cancel"))
		})

		It("should display competency categories", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("technical_leadership"))
		})

		It("should handle empty facts with message", func() {
			modal = burstview.NewFactSuggestion([]display.Fact{}, theme)
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("No facts detected"))
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()
		})

		It("should navigate to next fact with 'j' key", func() {
			Expect(modal.GetCurrentFact().ID).To(Equal("f1"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentFact().ID).To(Equal("f2"))
		})

		It("should navigate to previous fact with 'k' key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentFact().ID).To(Equal("f2"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(modal.GetCurrentFact().ID).To(Equal("f1"))
		})

		It("should navigate with down arrow key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(modal.GetCurrentFact().ID).To(Equal("f2"))
		})

		It("should navigate with up arrow key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(modal.GetCurrentFact().ID).To(Equal("f1"))
		})

		It("should handle pgdown navigation", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyPgDown})
			Expect(modal.GetCurrentFact().ID).To(Equal("f3"))
		})

		It("should handle pgup navigation", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyPgDown})
			modal.Update(tea.KeyMsg{Type: tea.KeyPgUp})
			Expect(modal.GetCurrentFact().ID).To(Equal("f1"))
		})

		It("should handle 'n' key for page down", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.GetCurrentFact().ID).To(Equal("f3"))
		})

		It("should handle 'p' key for page up", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
			Expect(modal.GetCurrentFact().ID).To(Equal("f1"))
		})
	})

	Describe("Accept", func() {
		BeforeEach(func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()
		})

		It("should accept current fact with 'a' key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedFacts()
			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].ID).To(Equal("f1"))
		})

		It("should set action to accept", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionActionAccept))
		})

		It("should remove accepted fact from pending list", func() {
			Expect(modal.GetSuggestionsCount()).To(Equal(3))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetSuggestionsCount()).To(Equal(2))
		})

		It("should adjust index when accepting non-last fact", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			current := modal.GetCurrentFact()
			Expect(current).NotTo(BeNil())
			Expect(current.ID).To(Equal("f2"))
		})

		It("should close modal when last fact accepted", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should accept multiple facts in sequence", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedFacts()
			Expect(accepted).To(HaveLen(2))
			Expect(accepted[0].ID).To(Equal("f1"))
			Expect(accepted[1].ID).To(Equal("f2"))
		})

		It("should preserve all accepted facts when modal closes", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedFacts()
			Expect(accepted).To(HaveLen(3))
		})

		It("should adjust index when accepting last item in list", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentFact().ID).To(Equal("f3"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetCurrentFact().ID).To(Equal("f2"))
		})
	})

	Describe("Reject", func() {
		BeforeEach(func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()
		})

		It("should reject current fact with 'r' key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			accepted := modal.GetAcceptedFacts()
			Expect(accepted).To(BeEmpty())

			Expect(modal.GetCurrentFact().ID).To(Equal("f2"))
		})

		It("should set action to reject", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionActionReject))
		})

		It("should remove rejected fact from pending list", func() {
			Expect(modal.GetSuggestionsCount()).To(Equal(3))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.GetSuggestionsCount()).To(Equal(2))
		})

		It("should close modal when last fact rejected", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.HasSuggestions()).To(BeFalse())
		})

		It("should not add rejected to accepted list", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			accepted := modal.GetAcceptedFacts()
			Expect(accepted).To(BeEmpty())
		})
	})

	Describe("Mixed Accept/Reject", func() {
		BeforeEach(func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()
		})

		It("should handle accept, reject, accept sequence", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedFacts()
			Expect(accepted).To(HaveLen(2))
			Expect(accepted[0].ID).To(Equal("f1"))
			Expect(accepted[1].ID).To(Equal("f3"))
		})

		It("should handle reject, accept, reject sequence", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			accepted := modal.GetAcceptedFacts()
			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].ID).To(Equal("f2"))
		})

		It("should handle navigation between accept/reject", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			current := modal.GetCurrentFact()
			Expect(current.ID).To(Equal("f3"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			current = modal.GetCurrentFact()
			Expect(current.ID).To(Equal("f1"))

			accepted := modal.GetAcceptedFacts()
			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].ID).To(Equal("f2"))
		})
	})

	Describe("Cancel", func() {
		BeforeEach(func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()
		})

		It("should set action to cancel on Esc", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionActionCancel))
		})

		It("should hide modal on Esc", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should preserve accepted facts after cancel", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAcceptedFacts()).To(HaveLen(1))

			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(modal.GetAcceptedFacts()).To(HaveLen(1))
		})
	})

	Describe("View Events", func() {
		BeforeEach(func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()
		})

		It("should set view events action on Enter", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionActionViewEvents))
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Dimensions", func() {
		BeforeEach(func() {
			modal = burstview.NewFactSuggestion(facts, theme)
		})

		It("should accept dimension changes", func() {
			modal.SetDimensions(100, 50)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle small dimensions", func() {
			modal.SetDimensions(50, 10)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle window size message", func() {
			modal.Show()
			_, cmd := modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			_ = cmd
		})
	})

	Describe("Edge Cases", func() {
		It("should handle update when modal is not visible", func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Hide()

			_, cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).To(BeNil())
			Expect(modal.GetAcceptedFacts()).To(BeEmpty())
		})

		It("should not mutate original facts slice", func() {
			originalLen := len(facts)
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			Expect(facts).To(HaveLen(originalLen))
		})

		It("should report HasSuggestions correctly for facts", func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			Expect(modal.HasSuggestions()).To(BeTrue())
			Expect(modal.GetSuggestionsCount()).To(Equal(3))
		})

		It("should report HasSuggestions false when empty", func() {
			modal = burstview.NewFactSuggestion([]display.Fact{}, theme)
			Expect(modal.HasSuggestions()).To(BeFalse())
			Expect(modal.GetSuggestionsCount()).To(Equal(0))
		})

		It("should handle SetDimensions with fact type", func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.SetDimensions(100, 50)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle SetDimensions with very small width for facts", func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.SetDimensions(30, 50)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle SetDimensions with very large width for facts", func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.SetDimensions(200, 50)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle window size message with facts", func() {
			modal = burstview.NewFactSuggestion(facts, theme)
			modal.Show()
			_, cmd := modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(cmd).To(BeNil())
		})
	})
})
