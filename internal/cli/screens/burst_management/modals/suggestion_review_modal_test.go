package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SuggestionReviewModal", func() {
	var (
		modal       *modals.SuggestionReviewModal
		suggestions []burst_fact.BurstSuggestion
		theme       themes.Theme
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		suggestions = []burst_fact.BurstSuggestion{
			{
				Name:            "Backend Development",
				Description:     "API and infrastructure work",
				EventIDs:        []string{"e1", "e2", "e3"},
				ConfidenceScore: 0.85,
			},
			{
				Name:            "Frontend Updates",
				Description:     "UI improvements and React migration",
				EventIDs:        []string{"e4", "e5"},
				ConfidenceScore: 0.72,
			},
			{
				Name:            "DevOps Initiative",
				Description:     "CI/CD and infrastructure automation",
				EventIDs:        []string{"e6", "e7", "e8", "e9"},
				ConfidenceScore: 0.91,
			},
		}
	})

	Describe("Creation", func() {
		It("should create modal with suggestions and theme", func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			Expect(modal).NotTo(BeNil())
			Expect(modal.HasSuggestions()).To(BeTrue())
			Expect(modal.GetSuggestionsCount()).To(Equal(3))
		})

		It("should create modal with nil theme using default", func() {
			modal = modals.NewSuggestionReviewModal(suggestions, nil)
			Expect(modal).NotTo(BeNil())
			// Should not panic and render correctly.
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should initialize at first suggestion (index 0)", func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()

			current := modal.GetCurrentSuggestion()
			Expect(current).NotTo(BeNil())
			Expect(current.Name).To(Equal("Backend Development"))
		})

		It("should initialize accepted list as empty", func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(BeEmpty())
		})

		It("should handle empty suggestions list", func() {
			modal = modals.NewSuggestionReviewModal([]burst_fact.BurstSuggestion{}, theme)
			Expect(modal.HasSuggestions()).To(BeFalse())
			Expect(modal.GetSuggestionsCount()).To(Equal(0))
			Expect(modal.GetCurrentSuggestion()).To(BeNil())
		})

		It("should handle single suggestion", func() {
			singleSuggestion := []burst_fact.BurstSuggestion{suggestions[0]}
			modal = modals.NewSuggestionReviewModal(singleSuggestion, theme)
			Expect(modal.HasSuggestions()).To(BeTrue())
			Expect(modal.GetSuggestionsCount()).To(Equal(1))
		})
	})

	Describe("Visibility", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
		})

		It("should be visible by default after creation", func() {
			// DetailModal is visible by default.
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should remain visible after Show()", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should not be visible after Hide()", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()
		})

		It("should navigate to next suggestion with 'n' key", func() {
			// Start at first suggestion.
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))

			// Navigate to next.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))
		})

		It("should navigate to previous suggestion with 'p' key", func() {
			// Navigate to second first.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))

			// Navigate back.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))
		})

		It("should navigate with right arrow key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRight})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))
		})

		It("should navigate with left arrow key", func() {
			// Go to second first.
			modal.Update(tea.KeyMsg{Type: tea.KeyRight})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))

			// Go back.
			modal.Update(tea.KeyMsg{Type: tea.KeyLeft})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))
		})

		It("should not navigate before first suggestion (boundary)", func() {
			// Already at first, try to go back.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))

			// Try with left arrow.
			modal.Update(tea.KeyMsg{Type: tea.KeyLeft})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))
		})

		It("should not navigate after last suggestion (boundary)", func() {
			// Go to last suggestion.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))

			// Try to go further.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))

			// Try with right arrow.
			modal.Update(tea.KeyMsg{Type: tea.KeyRight})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))
		})

		It("should update content correctly when navigating", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
			Expect(view).To(ContainSubstring("1 of 3"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			view = modal.View()
			Expect(view).To(ContainSubstring("Frontend Updates"))
			Expect(view).To(ContainSubstring("2 of 3"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			view = modal.View()
			Expect(view).To(ContainSubstring("DevOps Initiative"))
			Expect(view).To(ContainSubstring("3 of 3"))
		})

		It("should show correct X of Y counter", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("1 of 3"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			view = modal.View()
			Expect(view).To(ContainSubstring("2 of 3"))
		})
	})

	Describe("Accept", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()
		})

		It("should accept current suggestion with 'a' key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].Name).To(Equal("Backend Development"))
		})

		It("should set action to accept", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAction()).To(Equal(modals.SuggestionActionAccept))
		})

		It("should remove accepted suggestion from pending list", func() {
			Expect(modal.GetSuggestionsCount()).To(Equal(3))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetSuggestionsCount()).To(Equal(2))
		})

		It("should adjust index when accepting non-last suggestion", func() {
			// Accept first, should show second (now first).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			current := modal.GetCurrentSuggestion()
			Expect(current).NotTo(BeNil())
			Expect(current.Name).To(Equal("Frontend Updates"))
		})

		It("should close modal when last suggestion accepted", func() {
			// Accept all three.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should accept multiple suggestions in sequence", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(2))
			Expect(accepted[0].Name).To(Equal("Backend Development"))
			Expect(accepted[1].Name).To(Equal("Frontend Updates"))
		})

		It("should preserve all accepted suggestions when modal closes", func() {
			// Accept all.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(3))
		})
	})

	Describe("Reject", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()
		})

		It("should reject current suggestion with 'r' key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			// Rejected should not be in accepted.
			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(BeEmpty())

			// Should move to next.
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))
		})

		It("should set action to reject", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.GetAction()).To(Equal(modals.SuggestionActionReject))
		})

		It("should remove rejected suggestion from pending list", func() {
			Expect(modal.GetSuggestionsCount()).To(Equal(3))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.GetSuggestionsCount()).To(Equal(2))
		})

		It("should adjust index when rejecting non-last suggestion", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			// Should now show second suggestion (originally at index 1).
			current := modal.GetCurrentSuggestion()
			Expect(current.Name).To(Equal("Frontend Updates"))
		})

		It("should close modal when last suggestion rejected", func() {
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

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(BeEmpty())
		})
	})

	Describe("Mixed Accept/Reject", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()
		})

		It("should handle accept, reject, accept sequence", func() {
			// Accept first (Backend Development).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			// Reject second (Frontend Updates).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			// Accept third (DevOps Initiative).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(2))
			Expect(accepted[0].Name).To(Equal("Backend Development"))
			Expect(accepted[1].Name).To(Equal("DevOps Initiative"))
		})

		It("should handle reject, accept, reject sequence", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].Name).To(Equal("Frontend Updates"))
		})

		It("should handle navigation between accept/reject", func() {
			// Navigate to second.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			// Accept second (Frontend Updates).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// Now showing DevOps (was third, now second).
			current := modal.GetCurrentSuggestion()
			Expect(current.Name).To(Equal("DevOps Initiative"))

			// Reject DevOps.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			// Now showing Backend (only one left).
			current = modal.GetCurrentSuggestion()
			Expect(current.Name).To(Equal("Backend Development"))

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].Name).To(Equal("Frontend Updates"))
		})
	})

	Describe("Cancel", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()
		})

		It("should set action to cancel on Esc", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.GetAction()).To(Equal(modals.SuggestionActionCancel))
		})

		It("should hide modal on Esc", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should preserve accepted suggestions after cancel", func() {
			// Accept one first.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAcceptedSuggestions()).To(HaveLen(1))

			// Cancel.
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Accepted should still be preserved.
			Expect(modal.GetAcceptedSuggestions()).To(HaveLen(1))
		})

		It("should clear action after Show() is called again", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.GetAction()).To(Equal(modals.SuggestionActionCancel))

			modal.Show()
			Expect(modal.GetAction()).To(Equal(modals.SuggestionAction("")))
		})
	})

	Describe("Content Rendering", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()
		})

		It("should display suggestion name", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
		})

		It("should display suggestion description", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("API and infrastructure work"))
		})

		It("should display events count", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Events:"))
			Expect(view).To(ContainSubstring("3"))
		})

		It("should display confidence percentage formatted", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("85.0%"))
		})

		It("should handle empty description", func() {
			emptyDescSuggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "No Description",
					Description:     "",
					EventIDs:        []string{"e1"},
					ConfidenceScore: 0.60,
				},
			}
			modal = modals.NewSuggestionReviewModal(emptyDescSuggestions, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("No Description"))
			// Should not crash or show weird formatting.
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle high confidence score", func() {
			highConfSuggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "High Confidence",
					Description:     "Very confident",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.999,
				},
			}
			modal = modals.NewSuggestionReviewModal(highConfSuggestions, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("99.9%"))
		})

		It("should handle low confidence score", func() {
			lowConfSuggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Low Confidence",
					Description:     "Less confident",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.50,
				},
			}
			modal = modals.NewSuggestionReviewModal(lowConfSuggestions, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("50.0%"))
		})

		It("should show footer badges", func() {
			view := modal.View()
			// Check for key hints in footer.
			Expect(view).To(ContainSubstring("Accept"))
			Expect(view).To(ContainSubstring("Reject"))
			Expect(view).To(ContainSubstring("Next"))
			Expect(view).To(ContainSubstring("Cancel"))
		})
	})

	Describe("Dimensions", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
		})

		It("should accept dimension changes", func() {
			modal.SetDimensions(100, 50)
			// Should not panic.
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Action Management", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()
		})

		It("should clear action with ClearAction()", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAction()).To(Equal(modals.SuggestionActionAccept))

			modal.ClearAction()
			Expect(modal.GetAction()).To(Equal(modals.SuggestionAction("")))
		})
	})

	Describe("Edge Cases", func() {
		It("should handle update when modal is not visible", func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			// Hide the modal first (it's visible by default).
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			// Should not panic and should not process.
			_, cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).To(BeNil())

			// Should not have accepted anything.
			Expect(modal.GetAcceptedSuggestions()).To(BeEmpty())
		})

		It("should handle init correctly", func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			cmd := modal.Init()
			// Should return nil or valid command.
			_ = cmd
		})

		It("should handle window size message", func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()

			_, cmd := modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			// Should not panic.
			_ = cmd
		})
	})
})
