package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SuggestionReviewModal", func() {
	var (
		modal       *modals.SuggestionReviewModal
		suggestions []burstfact.BurstSuggestion
		theme       themes.Theme
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		suggestions = []burstfact.BurstSuggestion{
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

	Describe("Display All Suggestions", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()
		})

		It("should show all suggestion names in the view", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
			Expect(view).To(ContainSubstring("Frontend Updates"))
			Expect(view).To(ContainSubstring("DevOps Initiative"))
		})

		It("should show all confidence scores in the view", func() {
			view := modal.View()
			// All confidence scores should be visible as progress bars.
			Expect(view).To(ContainSubstring("85%"))
			Expect(view).To(ContainSubstring("72%"))
			Expect(view).To(ContainSubstring("91%"))
		})
	})

	Describe("Confidence Sorting", func() {
		BeforeEach(func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()
		})

		It("should sort suggestions by confidence (highest first)", func() {
			// GetAllSuggestions should return sorted by confidence descending.
			sorted := modal.GetAllSuggestions()
			Expect(sorted).To(HaveLen(3))
			Expect(sorted[0].Name).To(Equal("DevOps Initiative"))   // 0.91
			Expect(sorted[1].Name).To(Equal("Backend Development")) // 0.85
			Expect(sorted[2].Name).To(Equal("Frontend Updates"))    // 0.72
		})

		It("should select highest confidence suggestion first", func() {
			current := modal.GetCurrentSuggestion()
			Expect(current).NotTo(BeNil())
			// First selected should be highest confidence.
			Expect(current.Name).To(Equal("DevOps Initiative"))
			Expect(current.ConfidenceScore).To(Equal(0.91))
		})

		It("should navigate to next suggestion in confidence order", func() {
			// First is DevOps Initiative (highest).
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))

			// Navigate to next (second highest).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))

			// Navigate to next (third highest).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))
		})
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

		It("should initialize at first suggestion (highest confidence)", func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			modal.Show()

			current := modal.GetCurrentSuggestion()
			Expect(current).NotTo(BeNil())
			// First is highest confidence after sorting.
			Expect(current.Name).To(Equal("DevOps Initiative"))
		})

		It("should initialize accepted list as empty", func() {
			modal = modals.NewSuggestionReviewModal(suggestions, theme)
			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(BeEmpty())
		})

		It("should handle empty suggestions list", func() {
			modal = modals.NewSuggestionReviewModal([]burstfact.BurstSuggestion{}, theme)
			Expect(modal.HasSuggestions()).To(BeFalse())
			Expect(modal.GetSuggestionsCount()).To(Equal(0))
			Expect(modal.GetCurrentSuggestion()).To(BeNil())
		})

		It("should handle single suggestion", func() {
			singleSuggestion := []burstfact.BurstSuggestion{suggestions[0]}
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

		It("should navigate to next suggestion with 'j' key", func() {
			// Start at first suggestion (highest confidence after sorting).
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))

			// Navigate to next.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))
		})

		It("should navigate to previous suggestion with 'k' key", func() {
			// Navigate to second first.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))

			// Navigate back.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))
		})

		It("should navigate with down arrow key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))
		})

		It("should navigate with up arrow key", func() {
			// Go to second first.
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))

			// Go back.
			modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))
		})

		It("should not navigate before first suggestion (boundary)", func() {
			// Already at first, try to go back.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))

			// Try with up arrow.
			modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))
		})

		It("should not navigate after last suggestion (boundary)", func() {
			// Go to last suggestion.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))

			// Try to go further.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))

			// Try with down arrow.
			modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))
		})

		It("should show selected suggestion details when navigating", func() {
			// All suggestions visible in table, selected details shown below.
			view := modal.View()
			Expect(view).To(ContainSubstring("DevOps Initiative"))
			Expect(view).To(ContainSubstring("Selected:"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			view = modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
		})

		It("should handle pgdown navigation", func() {
			// With only 3 items and page size 10, pgdown goes to last item.
			modal.Update(tea.KeyMsg{Type: tea.KeyPgDown})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))
		})

		It("should handle pgup navigation", func() {
			// Navigate to last first.
			modal.Update(tea.KeyMsg{Type: tea.KeyPgDown})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))

			// Page up should go back to first.
			modal.Update(tea.KeyMsg{Type: tea.KeyPgUp})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))
		})

		It("should handle 'n' key for page down", func() {
			// With only 3 items and page size 10, n goes to last item.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))
		})

		It("should handle 'p' key for page up", func() {
			// Navigate to last first.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))

			// 'p' should go back to first.
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("DevOps Initiative"))
		})

		It("should show pagination info", func() {
			view := modal.View()
			// Should show "Suggestions: X | Page Y of Z".
			Expect(view).To(ContainSubstring("Suggestions:"))
			Expect(view).To(ContainSubstring("Page"))
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
			// First suggestion is highest confidence (DevOps Initiative).
			Expect(accepted[0].Name).To(Equal("DevOps Initiative"))
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
			// Accept first (DevOps), should show second (Backend, now first).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			current := modal.GetCurrentSuggestion()
			Expect(current).NotTo(BeNil())
			Expect(current.Name).To(Equal("Backend Development"))
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
			// Sorted order: DevOps (0.91), Backend (0.85), Frontend (0.72).
			Expect(accepted[0].Name).To(Equal("DevOps Initiative"))
			Expect(accepted[1].Name).To(Equal("Backend Development"))
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

			// Should move to next (Backend Development, now first after DevOps removed).
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))
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

			// Should now show Backend Development (was second, now first).
			current := modal.GetCurrentSuggestion()
			Expect(current.Name).To(Equal("Backend Development"))
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
			// Sorted order: DevOps (0.91), Backend (0.85), Frontend (0.72).
			// Accept first (DevOps Initiative).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			// Reject second (Backend Development).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			// Accept third (Frontend Updates).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(2))
			Expect(accepted[0].Name).To(Equal("DevOps Initiative"))
			Expect(accepted[1].Name).To(Equal("Frontend Updates"))
		})

		It("should handle reject, accept, reject sequence", func() {
			// Sorted order: DevOps (0.91), Backend (0.85), Frontend (0.72).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}) // Reject DevOps
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}) // Accept Backend
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}) // Reject Frontend

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].Name).To(Equal("Backend Development"))
		})

		It("should handle navigation between accept/reject", func() {
			// Sorted order: DevOps (0.91), Backend (0.85), Frontend (0.72).
			// Navigate to second (Backend).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			// Accept second (Backend Development).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// After removing Backend, index stays at 1. Remaining: [DevOps, Frontend].
			// Index 1 now points to Frontend (was third, now second).
			current := modal.GetCurrentSuggestion()
			Expect(current.Name).To(Equal("Frontend Updates"))

			// Reject Frontend (currently selected).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			// Now showing DevOps (only one left).
			current = modal.GetCurrentSuggestion()
			Expect(current.Name).To(Equal("DevOps Initiative"))

			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(1))
			Expect(accepted[0].Name).To(Equal("Backend Development"))
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

		It("should display selected suggestion description", func() {
			view := modal.View()
			// First selected is DevOps Initiative.
			Expect(view).To(ContainSubstring("CI/CD and infrastructure automation"))
		})

		It("should display events count in table", func() {
			view := modal.View()
			// Events column shows counts.
			Expect(view).To(ContainSubstring("4")) // DevOps has 4
			Expect(view).To(ContainSubstring("3")) // Backend has 3
			Expect(view).To(ContainSubstring("2")) // Frontend has 2
		})

		It("should display confidence percentage formatted", func() {
			view := modal.View()
			// All confidences visible in table as progress bars.
			Expect(view).To(ContainSubstring("91%"))
			Expect(view).To(ContainSubstring("85%"))
			Expect(view).To(ContainSubstring("72%"))
		})

		It("should handle empty description", func() {
			emptyDescSuggestions := []burstfact.BurstSuggestion{
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
			highConfSuggestions := []burstfact.BurstSuggestion{
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
			Expect(view).To(ContainSubstring("99%"))
		})

		It("should handle low confidence score", func() {
			lowConfSuggestions := []burstfact.BurstSuggestion{
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
			Expect(view).To(ContainSubstring("50%"))
		})

		It("should show footer badges with correct keys", func() {
			view := modal.View()
			// Check for key hints in footer.
			Expect(view).To(ContainSubstring("Accept"))
			Expect(view).To(ContainSubstring("Reject"))
			Expect(view).To(ContainSubstring("Navigate"))
			Expect(view).To(ContainSubstring("Page"))
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
