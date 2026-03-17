package burst_test

import (
	"github.com/baphled/kariya/internal/service/career/burstfact"
	burstview "github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SuggestionReview", func() {
	var (
		modal       *burstview.SuggestionReview
		suggestions []display.BurstSuggestion
		theme       themes.Theme
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		suggestions = display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{
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
		})
	})

	Describe("Display All Suggestions", func() {
		BeforeEach(func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			modal = burstview.NewSuggestionReview(suggestions, theme)
			Expect(modal).NotTo(BeNil())
			Expect(modal.HasSuggestions()).To(BeTrue())
			Expect(modal.GetSuggestionsCount()).To(Equal(3))
		})

		It("should create modal with nil theme using default", func() {
			modal = burstview.NewSuggestionReview(suggestions, nil)
			Expect(modal).NotTo(BeNil())
			// Should not panic and render correctly.
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should initialize at first suggestion (highest confidence)", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			current := modal.GetCurrentSuggestion()
			Expect(current).NotTo(BeNil())
			// First is highest confidence after sorting.
			Expect(current.Name).To(Equal("DevOps Initiative"))
		})

		It("should initialize accepted list as empty", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(BeEmpty())
		})

		It("should handle empty suggestions list", func() {
			modal = burstview.NewSuggestionReview([]display.BurstSuggestion{}, theme)
			Expect(modal.HasSuggestions()).To(BeFalse())
			Expect(modal.GetSuggestionsCount()).To(Equal(0))
			Expect(modal.GetCurrentSuggestion()).To(BeNil())
		})

		It("should handle single suggestion", func() {
			singleSuggestion := []display.BurstSuggestion{suggestions[0]}
			modal = burstview.NewSuggestionReview(singleSuggestion, theme)
			Expect(modal.HasSuggestions()).To(BeTrue())
			Expect(modal.GetSuggestionsCount()).To(Equal(1))
		})
	})

	Describe("Visibility", func() {
		BeforeEach(func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionActionAccept))
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

		It("should adjust index when accepting last item in list", func() {
			// Navigate to last (Frontend Updates, sorted third).
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))
		})
	})

	Describe("Reject", func() {
		BeforeEach(func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionActionReject))
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
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionActionCancel))

			modal.Show()
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionAction("")))
		})
	})

	Describe("Content Rendering", func() {
		BeforeEach(func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			emptyDescSuggestions := display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{
				{
					Name:            "No Description",
					Description:     "",
					EventIDs:        []string{"e1"},
					ConfidenceScore: 0.60,
				},
			})
			modal = burstview.NewSuggestionReview(emptyDescSuggestions, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("No Description"))
			// Should not crash or show weird formatting.
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle high confidence score", func() {
			highConfSuggestions := display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{
				{
					Name:            "High Confidence",
					Description:     "Very confident",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.999,
				},
			})
			modal = burstview.NewSuggestionReview(highConfSuggestions, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("99%"))
		})

		It("should handle low confidence score", func() {
			lowConfSuggestions := display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{
				{
					Name:            "Low Confidence",
					Description:     "Less confident",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.50,
				},
			})
			modal = burstview.NewSuggestionReview(lowConfSuggestions, theme)
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
			modal = burstview.NewSuggestionReview(suggestions, theme)
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
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()
		})

		It("should clear action with ClearAction()", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionActionAccept))

			modal.ClearAction()
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionAction("")))
		})
	})

	Describe("Edge Cases", func() {
		It("should handle update when modal is not visible", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			_, cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).To(BeNil())
			Expect(modal.GetAcceptedSuggestions()).To(BeEmpty())
		})

		It("should handle init correctly", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			cmd := modal.Init()
			_ = cmd
		})

		It("should render empty suggestion list view", func() {
			modal = burstview.NewSuggestionReview([]display.BurstSuggestion{}, theme)
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("No suggestions available"))
		})

		It("should ignore unknown rune keys", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			_, cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(cmd).To(BeNil())
			Expect(modal.GetSuggestionsCount()).To(Equal(3))
		})

		It("should handle window size message", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			_, cmd := modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			// Should not panic.
			_ = cmd
		})

		It("should handle very small window dimensions in View()", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.SetDimensions(40, 10)
			modal.Show()

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle window resize with small height", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			_, cmd := modal.Update(tea.WindowSizeMsg{Width: 80, Height: 8})
			Expect(cmd).To(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle Enter key to view events", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.GetAction()).To(Equal(burstview.SuggestionActionViewEvents))
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should handle SetDimensions with very small width", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.SetDimensions(30, 50)
			modal.Show()

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle SetDimensions with very large width", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.SetDimensions(200, 50)
			modal.Show()

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle removing middle suggestion", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.GetSuggestionsCount()).To(Equal(2))
		})

		It("should handle removing last suggestion and adjusting index", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))
		})

		It("should handle accepting and rejecting same suggestion", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAcceptedSuggestions()).To(HaveLen(1))

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.GetSuggestionsCount()).To(Equal(1))
		})

		It("should handle View with nil theme", func() {
			modal = burstview.NewSuggestionReview(suggestions, nil)
			modal.Show()

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Review Burst Suggestions"))
		})

		It("should handle very small height in View", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.SetDimensions(100, 5)
			modal.Show()

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle Update with unknown message type", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			type UnknownMsg struct{}
			_, cmd := modal.Update(UnknownMsg{})
			Expect(cmd).To(BeNil())
		})

		It("should handle accepting first suggestion and then navigating", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Frontend Updates"))
		})

		It("should handle rejecting first suggestion and then navigating", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			Expect(modal.GetCurrentSuggestion().Name).To(Equal("Backend Development"))
		})

		It("should handle window resize with small width", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			_, cmd := modal.Update(tea.WindowSizeMsg{Width: 50, Height: 40})
			Expect(cmd).To(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle window resize with large width", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			_, cmd := modal.Update(tea.WindowSizeMsg{Width: 200, Height: 40})
			Expect(cmd).To(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle multiple window resize messages", func() {
			modal = burstview.NewSuggestionReview(suggestions, theme)
			modal.Show()

			modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
			modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			modal.Update(tea.WindowSizeMsg{Width: 80, Height: 25})

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
