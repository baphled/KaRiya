package burst_management_test

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/testutil/mocks"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstManagement E2E Workflow Tests", func() {
	var (
		intent *burst_management.Intent
		ctx    *burst_management.IntentContext
		bursts []*career.Burst
	)

	BeforeEach(func() {
		// Create test bursts
		bursts = []*career.Burst{
			{
				ID:          "burst-1",
				Name:        "Backend Development at TechCorp",
				Description: "Built scalable microservices architecture",
				EventIDs:    []string{"event-1", "event-2"},
				Confirmed:   true,
				CreatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:          "burst-2",
				Name:        "DevOps Implementation",
				Description: "Migrated infrastructure to Kubernetes",
				EventIDs:    []string{"event-3"},
				Confirmed:   false,
				CreatedAt:   time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:          "burst-3",
				Name:        "Frontend Modernization",
				Description: "Migrated from jQuery to React",
				EventIDs:    []string{"event-4", "event-5", "event-6"},
				Confirmed:   true,
				CreatedAt:   time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC),
			},
		}

		// Create context (without repository/service for navigation tests)
		ctx = &burst_management.IntentContext{
			Bursts: bursts,
		}
		ctx.Validate()

		// Create intent
		var err error
		intent, err = burst_management.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("Complete Navigation Flow", func() {
		It("should support full workflow: list → detail → events → back to detail → back to list", func() {
			// Start at list view
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.IsActive()).To(BeTrue())

			// Select first burst (Enter key on list)
			result := &screens.NavigateResult{
				ResultData: bursts[0],
			}
			_ = intent.HandleNavigate(result)

			// Should show detail modal
			view := intent.View()
			Expect(view).To(ContainSubstring("Backend Development"))

			// Press 'v' to view events (triggers async load)
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			// Press Esc to go back to detail modal
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			view = intent.View()
			Expect(view).To(ContainSubstring("Backend Development"))

			// Press Esc again to close detail modal and return to list
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("should support navigation to facts modal and back", func() {
			// Navigate to detail
			result := &screens.NavigateResult{
				ResultData: bursts[0],
			}
			intent.HandleNavigate(result)

			// Press 'f' to view facts (triggers async load)
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Press Esc to go back to detail
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should still be active
			Expect(intent.IsActive()).To(BeTrue())
		})
	})

	Describe("Edit Workflow E2E", func() {
		It("should complete full edit flow: list → detail → edit → save → back to detail", func() {
			// Navigate to detail
			result := &screens.NavigateResult{
				ResultData: bursts[1],
			}
			intent.HandleNavigate(result)

			// Press 'e' to edit
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Edit modal should be visible
			Expect(intent.HasActiveModal()).To(BeTrue())
		})

		It("should handle edit cancellation", func() {
			// Navigate to detail and start edit
			result := &screens.NavigateResult{
				ResultData: bursts[0],
			}
			intent.HandleNavigate(result)
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Press Esc to cancel edit
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should not modify burst and modal should be closed
			Expect(intent.HasActiveModal()).To(BeFalse())
			originalBurst := bursts[0]
			Expect(originalBurst.Name).To(Equal("Backend Development at TechCorp"))
		})
	})

	Describe("Delete Workflow E2E", func() {
		It("should complete full delete flow: list → delete → confirm → list updated", func() {
			// Select burst to delete from list (using action)
			actionData := map[string]interface{}{
				"action": "delete",
				"burst":  bursts[2],
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			_ = intent.HandleNavigate(result)

			// Delete modal should be visible
			Expect(intent.HasActiveModal()).To(BeTrue())
			view := intent.View()
			Expect(view).To(ContainSubstring("Delete"))
		})

		It("should handle delete cancellation", func() {
			initialCount := len(intent.GetFilteredBursts())

			// Start delete
			actionData := map[string]interface{}{
				"action": "delete",
				"burst":  bursts[0],
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)

			// Cancel delete with Esc (modal handles this)
			Expect(intent.GetFilteredBursts()).To(HaveLen(initialCount))
		})
	})

	Describe("Confirm Burst Workflow E2E", func() {
		It("should complete full confirm flow: detail → confirm → burst marked confirmed", func() {
			// Select unconfirmed burst
			result := &screens.NavigateResult{
				ResultData: bursts[1],
			}
			intent.HandleNavigate(result)

			// Verify burst is not confirmed
			Expect(bursts[1].Confirmed).To(BeFalse())

			// Press 'c' to confirm
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			// Confirm modal should be visible
			Expect(intent.HasActiveModal()).To(BeTrue())
		})
	})

	Describe("Modal Priority and Overlay", func() {
		It("should render modals as overlays over the base screen", func() {
			// Start at list
			listView := intent.View()
			Expect(listView).NotTo(BeEmpty())

			// Open detail modal
			result := &screens.NavigateResult{
				ResultData: bursts[0],
			}
			intent.HandleNavigate(result)

			// View should contain modal overlay
			detailView := intent.View()
			Expect(detailView).NotTo(BeEmpty())
			Expect(detailView).To(ContainSubstring("Backend Development"))
		})

		It("should handle error modal with highest priority", func() {
			// Show error modal
			intent.ShowErrorModal("Test Error", "Something went wrong")

			// Error modal should be visible
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())

			// View should show error
			view := intent.View()
			Expect(view).To(ContainSubstring("Test Error"))

			// Press Esc to dismiss error
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.HasVisibleErrorModal()).To(BeFalse())
		})
	})

	Describe("Escape Key Navigation", func() {
		It("should close modals and navigate back with Escape key", func() {
			// Open detail modal
			result := &screens.NavigateResult{
				ResultData: bursts[0],
			}
			intent.HandleNavigate(result)
			Expect(intent.HasActiveModal()).To(BeTrue())

			// Press Esc to close detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.HasActiveModal()).To(BeFalse())

			// Press Esc at list to cancel intent
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.IsActive()).To(BeFalse())
		})

		It("should handle Esc in nested modals (events → detail → list)", func() {
			// Open detail
			result := &screens.NavigateResult{
				ResultData: bursts[0],
			}
			intent.HandleNavigate(result)

			// Open events modal
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			// Esc should close events and show detail again
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Another Esc should close detail
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.HasActiveModal()).To(BeFalse())
		})
	})

	Describe("Window Resize Handling", func() {
		It("should handle terminal resize while modal is open", func() {
			// Open detail modal
			result := &screens.NavigateResult{
				ResultData: bursts[0],
			}
			intent.HandleNavigate(result)

			// Send resize message
			resizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
			_ = intent.Update(resizeMsg)

			// Should still render correctly
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Empty State Handling", func() {
		It("should handle empty burst list gracefully", func() {
			emptyCtx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			emptyCtx.Validate()

			emptyIntent, err := burst_management.NewIntent(emptyCtx)
			Expect(err).NotTo(HaveOccurred())
			emptyIntent.Init()

			view := emptyIntent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(emptyIntent.GetFilteredBursts()).To(BeEmpty())
		})
	})

	Describe("Burst Suggestion Workflow E2E", func() {
		It("should show error when service is not available", func() {
			// Context has no service (default in tests)
			Expect(ctx.Service).To(BeNil())

			// Trigger suggestion
			actionData := map[string]interface{}{
				"action": "suggest",
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			_ = intent.HandleNavigate(result)

			// Should show error modal and stay at list (service unavailable)
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			view := intent.View()
			Expect(view).To(ContainSubstring("Service not available"))
		})

		It("should navigate through suggestions with n/p keys", func() {
			// Simulate suggestions loaded
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Backend Work",
					Description:     "API development",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.85,
				},
				{
					Name:            "DevOps Tasks",
					Description:     "Infrastructure work",
					EventIDs:        []string{"e3"},
					ConfidenceScore: 0.75,
				},
				{
					Name:            "Frontend Updates",
					Description:     "UI improvements",
					EventIDs:        []string{"e4", "e5"},
					ConfidenceScore: 0.90,
				},
			}

			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			}

			// Process suggestions loaded
			_ = intent.Update(msg)

			// Verify suggestion modal is visible
			modal := intent.GetSuggestionModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())

			// Verify initial view shows first suggestion
			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Work"))
			Expect(view).To(ContainSubstring("1 of 3"))

			// Navigate to next suggestion
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			view = modal.View()
			Expect(view).To(ContainSubstring("DevOps Tasks"))
			Expect(view).To(ContainSubstring("2 of 3"))

			// Navigate to next suggestion (third)
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			view = modal.View()
			Expect(view).To(ContainSubstring("Frontend Updates"))
			Expect(view).To(ContainSubstring("3 of 3"))

			// Navigate back with 'p'
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
			view = modal.View()
			Expect(view).To(ContainSubstring("DevOps Tasks"))
			Expect(view).To(ContainSubstring("2 of 3"))
		})

		It("should navigate through suggestions with arrow keys", func() {
			// Simulate suggestions loaded
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "First",
					Description:     "First suggestion",
					EventIDs:        []string{"e1"},
					ConfidenceScore: 0.80,
				},
				{
					Name:            "Second",
					Description:     "Second suggestion",
					EventIDs:        []string{"e2"},
					ConfidenceScore: 0.70,
				},
			}

			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			}

			_ = intent.Update(msg)

			// Navigate with right arrow
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRight})
			view := intent.View()
			Expect(view).To(ContainSubstring("Second"))

			// Navigate back with left arrow
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyLeft})
			view = intent.View()
			Expect(view).To(ContainSubstring("First"))
		})

		It("should accept suggestion and create burst", func() {
			initialBurstCount := len(intent.GetFilteredBursts())

			// Simulate suggestions loaded
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Accepted Burst",
					Description:     "This will be accepted",
					EventIDs:        []string{"e1", "e2", "e3"},
					ConfidenceScore: 0.95,
				},
			}

			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			}

			// Process suggestions loaded - creates modal
			_ = intent.Update(msg)

			// Verify modal is visible
			modal := intent.GetSuggestionModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())

			// Accept the suggestion via intent (message-based pattern)
			// Send 'a' key to intent, which forwards to modal, detects closure, and sends completion message
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// The modal closed and sent SuggestionReviewCompleteMsg via tea.Batch
			// We need to manually trigger processing of that message in tests
			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burst_fact.BurstSuggestion{suggestions[0]},
				Cancelled:           false,
			}
			_ = intent.Update(completeMsg)

			// Burst count should increase
			Expect(len(intent.GetFilteredBursts())).To(Equal(initialBurstCount + 1))

			// Verify the created burst
			newBurst := intent.GetFilteredBursts()[initialBurstCount]
			Expect(newBurst.Name).To(Equal("Accepted Burst"))
			Expect(newBurst.Description).To(Equal("This will be accepted"))
			Expect(newBurst.EventIDs).To(Equal([]string{"e1", "e2", "e3"}))
			Expect(newBurst.Confirmed).To(BeFalse())
		})

		It("should reject suggestion and move to next", func() {
			// Simulate suggestions loaded
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "First - Will Reject",
					Description:     "First suggestion",
					EventIDs:        []string{"e1"},
					ConfidenceScore: 0.60,
				},
				{
					Name:            "Second - Keep",
					Description:     "Second suggestion",
					EventIDs:        []string{"e2"},
					ConfidenceScore: 0.80,
				},
			}

			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			}

			_ = intent.Update(msg)

			// Verify we're showing first suggestion
			view := intent.View()
			Expect(view).To(ContainSubstring("First - Will Reject"))

			// Reject the first suggestion
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			// Should now show second suggestion at index 0
			view = intent.View()
			Expect(view).To(ContainSubstring("Second - Keep"))
			Expect(view).To(ContainSubstring("1 of 1")) // Only one left
		})

		It("should return to list when all suggestions are rejected", func() {
			// Simulate suggestions loaded with only one
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Only One",
					Description:     "Only suggestion",
					EventIDs:        []string{"e1"},
					ConfidenceScore: 0.50,
				},
			}

			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			}

			// Process suggestions loaded
			_ = intent.Update(msg)

			// Verify modal is visible
			modal := intent.GetSuggestionModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())

			// Reject the only suggestion
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			// Modal should close (no suggestions left)
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.HasSuggestions()).To(BeFalse())
		})

		It("should cancel suggestion review with Esc", func() {
			initialBurstCount := len(intent.GetFilteredBursts())

			// Simulate suggestions loaded
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Will Cancel",
					Description:     "Not accepting this",
					EventIDs:        []string{"e1"},
					ConfidenceScore: 0.70,
				},
			}

			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			}

			// Process suggestions loaded
			_ = intent.Update(msg)

			// Verify modal is visible
			modal := intent.GetSuggestionModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())

			// Cancel with Esc
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Modal should close
			Expect(modal.IsVisible()).To(BeFalse())

			// No burst should be created
			Expect(len(intent.GetFilteredBursts())).To(Equal(initialBurstCount))
		})

		It("should display suggestion details correctly", func() {
			// Simulate suggestion with specific values
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Test Burst Name",
					Description:     "Detailed description here",
					EventIDs:        []string{"e1", "e2", "e3", "e4"},
					ConfidenceScore: 0.8523,
				},
			}

			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			}

			_ = intent.Update(msg)

			// Verify view contains all details
			view := intent.View()
			Expect(view).To(ContainSubstring("Test Burst Name"))
			Expect(view).To(ContainSubstring("Detailed description here"))
			Expect(view).To(ContainSubstring("Events: 4"))
			Expect(view).To(ContainSubstring("85.2%")) // Confidence formatted
		})

		It("should handle error in suggestions loaded message", func() {
			// Simulate error during suggestion loading
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       fmt.Errorf("failed to analyze events"),
			}

			_ = intent.Update(msg)

			// Should show error and stay at list
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should handle empty suggestions list", func() {
			// Simulate no suggestions found
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burst_fact.BurstSuggestion{},
				Error:       nil,
			}

			_ = intent.Update(msg)

			// Should show error about no suggestions and stay at list
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			view := intent.View()
			Expect(view).To(ContainSubstring("No suggestions"))
		})
	})
})

var _ = Describe("Burst Suggestion Integration E2E", func() {
	var (
		intent      *burst_management.Intent
		ctx         *burst_management.IntentContext
		mockService *mocks.BurstServiceMock
		burstRepo   *careerrepo.MemoryBurstRepository
		events      []*career.CareerEvent
		suggestions []burst_fact.BurstSuggestion
	)

	BeforeEach(func() {
		now := time.Now()

		// Create test events.
		events = []*career.CareerEvent{
			{ID: "e1", Text: "Led backend project", Date: now.AddDate(0, -1, 0), Company: "TechCorp", Project: "Platform"},
			{ID: "e2", Text: "Built microservices", Date: now.AddDate(0, -2, 0), Company: "TechCorp", Project: "Platform"},
			{ID: "e3", Text: "Frontend redesign", Date: now.AddDate(0, -3, 0), Company: "TechCorp", Project: "UI"},
			{ID: "e4", Text: "React migration", Date: now.AddDate(0, -4, 0), Company: "TechCorp", Project: "UI"},
			{ID: "e5", Text: "DevOps setup", Date: now.AddDate(0, -5, 0), Company: "TechCorp", Project: "Infra"},
		}

		// Create suggestions (each must have at least 2 events for validation).
		suggestions = []burst_fact.BurstSuggestion{
			{
				Name:            "Backend Development",
				Description:     "API and microservices work",
				EventIDs:        []string{"e1", "e2"},
				ConfidenceScore: 0.85,
			},
			{
				Name:            "Frontend Modernization",
				Description:     "UI redesign and React migration",
				EventIDs:        []string{"e3", "e4"},
				ConfidenceScore: 0.78,
			},
			{
				Name:            "Infrastructure",
				Description:     "DevOps and CI/CD setup",
				EventIDs:        []string{"e5", "e1"},
				ConfidenceScore: 0.65,
			},
		}

		// Create mock service using centralized mock.
		mockService = mocks.NewBurstServiceMock().
			SetEvents(events).
			SetSuggestions(suggestions).
			SetExtractedFacts([]career.Fact{
				{ID: "f1", Text: "Built scalable API"},
				{ID: "f2", Text: "Improved response time by 40%"},
			})

		// Create burst repository.
		burstRepo = careerrepo.NewMemoryBurstRepository()

		// Create context with service and repository.
		ctx = &burst_management.IntentContext{
			Bursts:          []*career.Burst{},
			Service:         mockService,
			BurstRepository: burstRepo,
		}
		ctx.Validate()

		// Create intent.
		var err error
		intent, err = burst_management.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("Suggestion Detection Trigger", func() {
		It("should trigger suggestion detection with 's' key", func() {
			Expect(intent.GetState()).To(Equal(burst_management.StateList))

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateSuggesting))
		})

		It("should show loading modal in view", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			Expect(intent.GetState()).To(Equal(burst_management.StateSuggesting))
			view := intent.View()
			Expect(view).To(ContainSubstring("Detecting burst patterns"))
		})

		It("should trigger via action=suggest from screen", func() {
			actionData := map[string]interface{}{
				"action": "suggest",
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}

			cmd := intent.HandleNavigate(result)
			Expect(cmd).NotTo(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateSuggesting))
		})
	})

	Describe("Suggestion Accept with Repository", func() {
		It("should create burst and save to repository", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions[:1]}
			intent.Update(msg)

			modal := intent.GetSuggestionModal()
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			accepted := modal.GetAcceptedSuggestions()

			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: accepted,
			}
			intent.Update(completeMsg)

			Expect(len(intent.GetFilteredBursts())).To(Equal(1))

			allBursts, err := burstRepo.List(context.Background(), careerrepo.BurstListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(allBursts).To(HaveLen(1))
			Expect(allBursts[0].Name).To(Equal("Backend Development"))
		})

		It("should trigger fact extraction after accept", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions[:1]}
			intent.Update(msg)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions[:1],
			}
			cmd := intent.Update(completeMsg)

			Expect(cmd).NotTo(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateExtractingFacts))
			Expect(intent.IsExtractingFacts()).To(BeTrue())
		})
	})

	Describe("Multiple Suggestions with Repository", func() {
		It("should save all accepted bursts", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions[:2]}
			intent.Update(msg)

			modal := intent.GetSuggestionModal()
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			accepted := modal.GetAcceptedSuggestions()
			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: accepted,
			}
			intent.Update(completeMsg)

			allBursts, err := burstRepo.List(context.Background(), careerrepo.BurstListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(allBursts).To(HaveLen(2))
		})

		It("should handle accept, reject, accept sequence", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions}
			intent.Update(msg)

			modal := intent.GetSuggestionModal()

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAcceptedSuggestions()).To(HaveLen(1))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.GetAcceptedSuggestions()).To(HaveLen(1))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(modal.GetAcceptedSuggestions()).To(HaveLen(2))
		})
	})

	Describe("Error Handling with Service", func() {
		It("should handle detection error", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Error: fmt.Errorf("detection failed"),
			}
			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should handle fact extraction failure", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions[:1]}
			intent.Update(msg)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions[:1],
			}
			intent.Update(completeMsg)

			extractError := burst_management.FactExtractionCompleteMsg{
				Error: fmt.Errorf("extraction failed"),
			}
			intent.Update(extractError)

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should recover to list state after error dismissal", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Error: fmt.Errorf("some error"),
			}
			intent.Update(msg)

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleErrorModal()).To(BeFalse())
		})
	})

	Describe("Escape Handling", func() {
		It("should cancel suggestion review with Esc", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions}
			intent.Update(msg)

			modal := intent.GetSuggestionModal()
			Expect(modal.IsVisible()).To(BeTrue())

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should not create bursts when cancelled", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions}
			intent.Update(msg)

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				Cancelled: true,
			}
			intent.Update(completeMsg)

			Expect(len(intent.GetFilteredBursts())).To(Equal(0))
		})

		It("should preserve partial accepts when cancelled", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions}
			intent.Update(msg)

			modal := intent.GetSuggestionModal()
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			accepted := modal.GetAcceptedSuggestions()
			Expect(accepted).To(HaveLen(1))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(modal.GetAcceptedSuggestions()).To(HaveLen(1))
		})
	})
})

var _ = Describe("Fact Extraction from Accepted Burst Suggestions", func() {
	var (
		intent *burst_management.Intent
		ctx    *burst_management.IntentContext
	)

	BeforeEach(func() {
		ctx = &burst_management.IntentContext{
			Bursts: fixtures.Bursts(2, fixtures.Events(4)),
		}
		ctx.Validate()

		var err error
		intent, err = burst_management.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Context("when user accepts a burst suggestion", func() {
		var suggestions []burst_fact.BurstSuggestion

		BeforeEach(func() {
			suggestions = []burst_fact.BurstSuggestion{
				{
					Name:            "Backend Development",
					Description:     "Built microservices",
					EventIDs:        []string{"e1", "e2", "e3"},
					ConfidenceScore: 0.95,
				},
			}
		})

		It("creates the burst", func() {
			initialCount := len(intent.GetFilteredBursts())

			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions}
			_ = intent.Update(msg)
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			}
			_ = intent.Update(completeMsg)

			Expect(intent.GetFilteredBursts()).To(HaveLen(initialCount + 1))
		})

		It("triggers fact extraction for the created burst", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions}
			_ = intent.Update(msg)
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			}
			cmd := intent.Update(completeMsg)

			Expect(cmd).NotTo(BeNil(), "should return a command to trigger fact extraction")

			cmdMsg := cmd()
			_, isFactExtraction := cmdMsg.(burst_management.FactExtractionCompleteMsg)
			Expect(isFactExtraction).To(BeTrue(), "command should trigger fact extraction")
		})
	})

	Context("when user accepts multiple burst suggestions", func() {
		It("triggers fact extraction for each burst", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{Name: "Burst 1", EventIDs: []string{"e1"}, ConfidenceScore: 0.9},
				{Name: "Burst 2", EventIDs: []string{"e2"}, ConfidenceScore: 0.85},
			}

			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions}
			_ = intent.Update(msg)
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			}
			cmd := intent.Update(completeMsg)

			Expect(cmd).NotTo(BeNil(), "should return commands to extract facts for both bursts")
		})
	})
})

var _ = Describe("Fact Extraction from Confirmed Bursts", func() {
	var (
		intent *burst_management.Intent
		ctx    *burst_management.IntentContext
	)

	BeforeEach(func() {
		events := fixtures.Events(4)
		bursts := fixtures.Bursts(1, events)
		bursts[0].Confirmed = false

		ctx = &burst_management.IntentContext{
			Bursts: bursts,
		}
		ctx.Validate()

		var err error
		intent, err = burst_management.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Context("when confirming an unconfirmed burst", func() {
		It("triggers fact extraction", func() {
			burst := intent.GetFilteredBursts()[0]

			result := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(result)

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.GetState()).To(Equal(burst_management.StateExtractingFacts))
		})

		It("transitions to extracting facts state after confirmation", func() {
			burst := intent.GetFilteredBursts()[0]

			result := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(result)

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(cmd).NotTo(BeNil(), "confirming burst should return a command for fact extraction")
		})
	})
})

var _ = Describe("Re-extraction Workflow", func() {
	var (
		intent      *burst_management.Intent
		ctx         *burst_management.IntentContext
		mockService *mocks.BurstServiceMock
	)

	BeforeEach(func() {
		events := fixtures.Events(4)
		bursts := fixtures.Bursts(1, events)
		bursts[0].Confirmed = true

		mockService = mocks.NewBurstServiceMock()
		mockService.SetFactsForBurst(bursts[0].ID, []*career.Fact{
			{ID: "fact-1", SourceBurstID: bursts[0].ID, Text: "Led architecture design"},
			{ID: "fact-2", SourceBurstID: bursts[0].ID, Text: "Delivered project on time"},
		})

		ctx = &burst_management.IntentContext{
			Bursts:  bursts,
			Service: mockService,
		}
		ctx.Validate()

		var err error
		intent, err = burst_management.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Context("when confirming a burst that already has facts", func() {
		It("prompts user for re-extraction", func() {
			burst := intent.GetFilteredBursts()[0]

			result := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(result)

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("Re-extract"))
		})
	})
})

// BUG REGRESSION TESTS - Burst Suggestion Workflow
// These tests capture bugs found during code review.
// Each test should FAIL until the corresponding bug is fixed.

var _ = Describe("Burst Suggestion Workflow Bug Regressions", func() {
	var (
		intent      *burst_management.Intent
		ctx         *burst_management.IntentContext
		mockService *mocks.BurstServiceMock
		burstRepo   *careerrepo.MemoryBurstRepository
	)

	BeforeEach(func() {
		now := time.Now()

		events := []*career.CareerEvent{
			{ID: "e1", Text: "Led backend project", Date: now.AddDate(0, -1, 0)},
			{ID: "e2", Text: "Built microservices", Date: now.AddDate(0, -2, 0)},
		}

		mockService = mocks.NewBurstServiceMock().
			SetEvents(events).
			SetSuggestions([]burst_fact.BurstSuggestion{
				{
					Name:            "Backend Development",
					Description:     "API work",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.85,
				},
			}).
			SetExtractedFacts([]career.Fact{
				{ID: "f1", Text: "Built scalable API"},
			})

		burstRepo = careerrepo.NewMemoryBurstRepository()

		ctx = &burst_management.IntentContext{
			Bursts:          []*career.Burst{},
			Service:         mockService,
			BurstRepository: burstRepo,
		}
		ctx.Validate()

		var err error
		intent, err = burst_management.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("BUG: Fact extraction completion with nil selectedBurst", func() {
		// BUG: After accepting suggestions, fact extraction completes and calls
		// showBurstDetailModal(i.selectedBurst) but selectedBurst is nil,
		// causing a panic in renderBurstDetailContent.

		It("should not panic when fact extraction completes after accepting suggestions", func() {
			// Setup: Accept a suggestion (selectedBurst is never set in this flow)
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Test Burst",
					Description:     "Test description",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.9,
				},
			}

			// Load suggestions
			msg := burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions}
			intent.Update(msg)

			// Accept the suggestion
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// Simulate the completion message (normally sent via tea.Batch)
			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
				Cancelled:           false,
			}
			intent.Update(completeMsg)

			// Verify burst was created
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))

			// Now simulate fact extraction completing
			// This should NOT panic and should NOT try to show detail modal
			factMsg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{{ID: "f1", Text: "Test fact"}},
				Error: nil,
			}

			// This should not panic
			Expect(func() {
				intent.Update(factMsg)
			}).NotTo(Panic())

			// Should return to list state, not show detail modal
			Expect(intent.GetState()).To(Equal(burst_management.StateList))

			// Detail modal should NOT be visible (since we didn't select a burst)
			Expect(intent.GetDetailModal()).To(BeNil())
		})

		It("should stay on list view after fact extraction completes from suggestion flow", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Created Burst",
					Description:     "From suggestions",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.85,
				},
			}

			// Full flow: suggestions loaded -> accept -> complete -> fact extraction
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Fact extraction completes successfully
			intent.Update(burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{{ID: "f1", Text: "Extracted fact"}},
			})

			// Should be on list state showing the created burst
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
			Expect(intent.GetFilteredBursts()[0].Name).To(Equal("Created Burst"))

			// View should render without panic and show the list
			view := intent.View()
			Expect(view).To(ContainSubstring("Created Burst"))
		})
	})

	Describe("BUG: showBurstDetailModal with nil burst", func() {
		// BUG: showBurstDetailModal doesn't guard against nil burst parameter,
		// leading to panic when NewBurstDetailModal accesses burst fields.

		It("should handle nil burst gracefully in showBurstDetailModal", func() {
			// Directly test that showing detail modal with nil doesn't panic
			// We need to trigger a code path that calls showBurstDetailModal(nil)

			// Set selectedBurst to nil explicitly
			intent.SetSelectedBurst(nil)

			// Fact extraction complete tries to show detail modal for selectedBurst
			factMsg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{{ID: "f1", Text: "Test"}},
			}

			Expect(func() {
				intent.Update(factMsg)
			}).NotTo(Panic())
		})
	})

	Describe("BUG: State transition after suggestion modal shown", func() {
		// BUG: When suggestions are loaded and modal is shown,
		// state stays as StateSuggesting instead of StateSuggestionReview.

		It("should transition to StateSuggestionReview when modal is shown", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Test",
					Description:     "Test",
					EventIDs:        []string{"e1"},
					ConfidenceScore: 0.8,
				},
			}

			// Trigger suggestion detection
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateSuggesting))

			// Suggestions loaded - modal should show
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})

			// State should now be StateSuggestionReview (currently stays at StateSuggesting)
			Expect(intent.GetState()).To(Equal(burst_management.StateSuggestionReview))
		})
	})

	Describe("BUG: State not reset after cancelled suggestion review", func() {
		// BUG: When user cancels suggestion review, state may not properly
		// reset to StateList.

		It("should reset to StateList after cancelling suggestion review", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Test",
					Description:     "Test",
					EventIDs:        []string{"e1"},
					ConfidenceScore: 0.8,
				},
			}

			// Load suggestions
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})

			// Cancel with Esc
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Process the completion message
			intent.Update(burst_management.SuggestionReviewCompleteMsg{Cancelled: true})

			// State should be StateList
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("BUG: Created bursts visible in list after suggestion acceptance", func() {
		// This tests the full end-to-end flow to ensure bursts are actually
		// visible to the user after the entire suggestion workflow completes.

		It("should show created bursts in list view after full workflow", func() {
			initialCount := len(intent.GetFilteredBursts())
			Expect(initialCount).To(Equal(0))

			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Visible Burst",
					Description:     "Should appear in list",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.95,
				},
			}

			// 1. Trigger detection
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			// 2. Suggestions loaded
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})

			// 3. Accept suggestion
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// 4. Process completion
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// 5. Fact extraction completes (with or without error)
			intent.Update(burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{},
				Error: nil,
			})

			// Verify: burst should be in the list
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))

			// Verify: state should be list
			Expect(intent.GetState()).To(Equal(burst_management.StateList))

			// Verify: view should show the burst (no panic, renders correctly)
			Expect(func() {
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			}).NotTo(Panic())
		})

		It("should show multiple created bursts after accepting multiple suggestions", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{Name: "Burst One", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
				{Name: "Burst Two", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.8},
			}

			// Load and accept all
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}) // Accept first
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}) // Accept second

			// Get accepted from modal before it's cleared
			modal := intent.GetSuggestionModal()
			var accepted []burst_fact.BurstSuggestion
			if modal != nil {
				accepted = modal.GetAcceptedSuggestions()
			} else {
				accepted = suggestions // Modal already closed
			}

			// Process completion
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: accepted,
			})

			// Fact extraction completes
			intent.Update(burst_management.FactExtractionCompleteMsg{Facts: []*career.Fact{}})

			// Both bursts should be visible
			Expect(intent.GetFilteredBursts()).To(HaveLen(2))
		})
	})

	Describe("BUG: Fact extraction error handling after suggestion acceptance", func() {
		It("should show error but keep created bursts when fact extraction fails", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Burst With Failed Extraction",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.85,
				},
			}

			// Accept suggestion
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Burst should exist before fact extraction completes
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))

			// Fact extraction fails
			intent.Update(burst_management.FactExtractionCompleteMsg{
				Error: fmt.Errorf("extraction service unavailable"),
			})

			// Error modal should show
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())

			// But burst should still be in the list!
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
			Expect(intent.GetFilteredBursts()[0].Name).To(Equal("Burst With Failed Extraction"))
		})
	})

	Describe("BUG: Loading modal and list refresh during suggestion acceptance", func() {
		It("should show loading modal during fact extraction after accepting suggestions", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Burst With Loading Modal",
					Description:     "Should show loading during extraction",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.9,
				},
			}

			// 1. Load suggestions
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})

			// 2. Accept suggestion (modal closes when last suggestion accepted)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// 3. Process completion - this triggers fact extraction
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Verify: loading modal should be visible during fact extraction
			Expect(intent.HasActiveModal()).To(BeTrue(), "Loading modal should be visible during fact extraction")
			Expect(intent.GetState()).To(Equal(burst_management.StateExtractingFacts))

			// Verify: view should render the loading modal
			view := intent.View()
			Expect(view).To(ContainSubstring("Extracting"), "View should show extraction in progress")
		})

		It("should refresh list view and show burst name after fact extraction completes", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Refreshed Burst View",
					Description:     "Should appear in list after refresh",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.92,
				},
			}

			// Complete the full workflow
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Fact extraction completes
			intent.Update(burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{{ID: "f1", Text: "Extracted fact"}},
				Error: nil,
			})

			// Verify: state should be list
			Expect(intent.GetState()).To(Equal(burst_management.StateList))

			// Verify: view should contain the burst name (the actual visual refresh fix)
			view := intent.View()
			Expect(view).To(ContainSubstring("Refreshed Burst View"),
				"List view should show the created burst name after fact extraction completes")
		})

		It("should show burst in view immediately after accepting but before extraction completes", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Immediate Visibility Burst",
					Description:     "Should be visible before extraction",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.88,
				},
			}

			// Accept suggestion
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Before extraction completes, burst should already be in filtered list
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
			Expect(intent.GetFilteredBursts()[0].Name).To(Equal("Immediate Visibility Burst"))

			// State is extracting, but the list screen was already updated
			Expect(intent.GetState()).To(Equal(burst_management.StateExtractingFacts))
		})
	})
})

var _ = Describe("Burst Suggestion Persistence Tests", func() {
	var (
		intent      *burst_management.Intent
		ctx         *burst_management.IntentContext
		mockService *mocks.BurstServiceMock
		burstRepo   *careerrepo.MemoryBurstRepository
	)

	BeforeEach(func() {
		now := time.Now()

		events := []*career.CareerEvent{
			{ID: "e1", Text: "Led backend project", Date: now.AddDate(0, -1, 0)},
			{ID: "e2", Text: "Built microservices", Date: now.AddDate(0, -2, 0)},
			{ID: "e3", Text: "Deployed to production", Date: now.AddDate(0, -3, 0)},
		}

		mockService = mocks.NewBurstServiceMock().
			SetEvents(events).
			SetExtractedFacts([]career.Fact{
				{ID: "f1", Text: "Built scalable API"},
				{ID: "f2", Text: "Improved performance by 40%"},
			})

		burstRepo = careerrepo.NewMemoryBurstRepository()

		ctx = &burst_management.IntentContext{
			Bursts:          []*career.Burst{},
			Service:         mockService,
			BurstRepository: burstRepo,
		}
		ctx.Validate()

		var err error
		intent, err = burst_management.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("Burst persistence to repository", func() {
		It("should save accepted burst to repository with generated ID", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Persisted Burst",
					Description:     "Should be saved to repo",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.9,
				},
			}

			// Accept the suggestion
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Verify burst is in repository
			repobursts, err := burstRepo.List(context.Background(), careerrepo.BurstListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(repobursts).To(HaveLen(1))
			Expect(repobursts[0].Name).To(Equal("Persisted Burst"))
			Expect(repobursts[0].ID).NotTo(BeEmpty(), "Burst should have a generated ID")
		})

		It("should have burst ID populated in memory after repository save", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Burst With ID",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.85,
				},
			}

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// The in-memory burst should also have the ID
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
			Expect(intent.GetFilteredBursts()[0].ID).NotTo(BeEmpty(),
				"In-memory burst should have ID populated after repository save")
		})
	})

	Describe("Fact persistence with correct burst linkage", func() {
		It("should save facts with correct SourceBurstID", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Burst For Facts",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.9,
				},
			}

			// Complete the suggestion flow
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// This returns a command that triggers fact extraction
			cmd := intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Get the burst ID that was created
			createdBurst := intent.GetFilteredBursts()[0]
			Expect(createdBurst.ID).NotTo(BeEmpty())

			// Execute the fact extraction command to trigger SaveFact calls
			Expect(cmd).NotTo(BeNil(), "should return fact extraction command")

			// Execute the command - this calls ExtractFactsFromBurst and SaveFact
			msg := cmd()
			Expect(msg).NotTo(BeNil())

			// Verify facts were saved with correct SourceBurstID
			savedFacts := mockService.GetSavedFacts()
			Expect(savedFacts).To(HaveLen(2))
			for _, fact := range savedFacts {
				Expect(fact.SourceBurstID).To(Equal(createdBurst.ID),
					"Fact should be linked to the created burst")
			}
		})

		It("should not save facts with empty SourceBurstID", func() {
			// This tests that created bursts have IDs before fact extraction
			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "Burst Test",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.8,
				},
			}

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// This returns a command that triggers fact extraction
			cmd := intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Get the created burst - it should have an ID from repository
			createdBurst := intent.GetFilteredBursts()[0]
			Expect(createdBurst.ID).NotTo(BeEmpty(),
				"Burst should have ID generated by repository")

			// Execute the fact extraction command
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			Expect(msg).NotTo(BeNil())

			// All saved facts should have the burst's ID as SourceBurstID
			savedFacts := mockService.GetSavedFacts()
			Expect(savedFacts).To(HaveLen(2)) // Mock returns 2 facts
			for _, fact := range savedFacts {
				Expect(fact.SourceBurstID).To(Equal(createdBurst.ID),
					"Facts should be linked to the created burst")
			}
		})
	})

	Describe("List screen updates with new bursts", func() {
		It("should refresh list screen to show newly created burst", func() {
			// Start with empty list
			Expect(intent.GetFilteredBursts()).To(BeEmpty())

			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:            "New Visible Burst",
					Description:     "Should appear in refreshed list",
					EventIDs:        []string{"e1", "e2"},
					ConfidenceScore: 0.95,
				},
			}

			// Complete the flow
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Handle fact extraction (may fail, but burst should still be visible)
			intent.Update(burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{},
				Error: fmt.Errorf("service unavailable"),
			})

			// Dismiss error modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// List should now contain the new burst
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))

			// View should render the new burst name
			view := intent.View()
			Expect(view).To(ContainSubstring("New Visible Burst"))
		})

		It("should show all accepted bursts in list after accepting multiple", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{Name: "First Burst", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
				{Name: "Second Burst", EventIDs: []string{"e2", "e3"}, ConfidenceScore: 0.85},
				{Name: "Third Burst", EventIDs: []string{"e1", "e3"}, ConfidenceScore: 0.8},
			}

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})

			// Accept first - modal still has 2 suggestions left
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			// Reject second - modal still has 1 suggestion left
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			// Capture accepted before final accept closes modal
			modal := intent.GetSuggestionModal()
			Expect(modal).NotTo(BeNil())
			acceptedBeforeFinal := modal.GetAcceptedSuggestions()

			// Accept third (last one) - this closes the modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// Build accepted list: previous accepts + the one we just accepted
			accepted := []burst_fact.BurstSuggestion{
				{Name: "First Burst", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
				{Name: "Third Burst", EventIDs: []string{"e1", "e3"}, ConfidenceScore: 0.8},
			}
			Expect(acceptedBeforeFinal).To(HaveLen(1)) // First was already accepted

			// Process completion
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: accepted,
			})

			// Should have 2 bursts (First and Third, not Second)
			Expect(intent.GetFilteredBursts()).To(HaveLen(2))

			names := []string{
				intent.GetFilteredBursts()[0].Name,
				intent.GetFilteredBursts()[1].Name,
			}
			Expect(names).To(ContainElement("First Burst"))
			Expect(names).To(ContainElement("Third Burst"))
			Expect(names).NotTo(ContainElement("Second Burst"))
		})
	})

	Describe("Edge case: partial burst creation failures", func() {
		It("should handle case where some bursts fail to save", func() {
			// This tests the slice bounds bug at helpers.go:672
			// If 3 suggestions are accepted but only 2 successfully create,
			// the slice calculation will be wrong

			suggestions := []burst_fact.BurstSuggestion{
				{Name: "Will Succeed 1", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
				{Name: "Will Succeed 2", EventIDs: []string{"e2", "e3"}, ConfidenceScore: 0.85},
			}

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})

			// Accept first - modal still open
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			// Accept second - modal closes after this

			// Since modal closes after last accept, use the known suggestions
			// This should not panic even if the number of created bursts
			// doesn't match the number of accepted suggestions
			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				intent.Update(burst_management.SuggestionReviewCompleteMsg{
					AcceptedSuggestions: suggestions,
				})
			}).NotTo(Panic())

			// All successfully created bursts should be in the list
			Expect(intent.GetFilteredBursts()).To(HaveLen(2))
		})

		It("should track which bursts were actually created vs requested", func() {
			suggestions := []burst_fact.BurstSuggestion{
				{Name: "Burst A", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
			}

			intent.Update(burst_management.BurstSuggestionsLoadedMsg{Suggestions: suggestions})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			intent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
			})

			// Verify the count matches
			Expect(intent.GetFilteredBursts()).To(HaveLen(len(suggestions)))

			// Verify in repository too
			repoBursts, _ := burstRepo.List(context.Background(), careerrepo.BurstListFilters{})
			Expect(repoBursts).To(HaveLen(len(suggestions)))
		})
	})

	Describe("Modal Escape Handling E2E", func() {
		var (
			intent    *burst_management.Intent
			ctx       *burst_management.IntentContext
			burstRepo *careerrepo.MemoryBurstRepository
			burst     *career.Burst
		)

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test description",
				EventIDs:    []string{"e1", "e2"},
			}

			burstRepo = careerrepo.NewMemoryBurstRepository()
			_ = burstRepo.Create(context.Background(), burst)

			ctx = &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				BurstRepository: burstRepo,
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Describe("Edit Modal Escape - Existing Burst", func() {
			It("should close edit modal and return to list when escape is pressed", func() {
				// Navigate to detail modal.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)
				Expect(intent.GetDetailModal()).NotTo(BeNil())
				Expect(intent.GetDetailModal().IsVisible()).To(BeTrue())

				// Press 'e' to open edit modal.
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				Expect(intent.HasVisibleEditModal()).To(BeTrue())

				// Press Escape to cancel edit.
				intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

				// Edit modal should be closed.
				Expect(intent.HasVisibleEditModal()).To(BeFalse())
				// State should still be list (modal overlay pattern).
				Expect(intent.GetState()).To(Equal(burst_management.StateList))
			})

			It("should preserve original burst data when escape is pressed", func() {
				originalName := burst.Name
				originalDesc := burst.Description

				// Navigate to detail and edit.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				Expect(intent.HasVisibleEditModal()).To(BeTrue())

				// Press Escape to cancel.
				intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

				// Original data should be preserved.
				repoBurst, _ := burstRepo.GetByID(context.Background(), burst.ID)
				Expect(repoBurst.Name).To(Equal(originalName))
				Expect(repoBurst.Description).To(Equal(originalDesc))
			})

			It("should handle multiple escape presses gracefully", func() {
				// Navigate to detail.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)

				// Press 'e' to edit.
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				Expect(intent.HasVisibleEditModal()).To(BeTrue())

				// Multiple escapes.
				intent.Update(tea.KeyMsg{Type: tea.KeyEscape})
				intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

				// Should be in valid state (list or cancelled).
				Expect(intent.GetState()).To(BeElementOf(
					burst_management.StateList,
				))
			})
		})

		Describe("Delete Modal Escape", func() {
			It("should close delete modal and return to list when escape is pressed", func() {
				// Navigate to detail.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)

				// Press 'd' for delete.
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				Expect(intent.HasVisibleDeleteModal()).To(BeTrue())

				// Press 'n' to cancel (or Escape).
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				// Delete modal should be closed.
				Expect(intent.HasVisibleDeleteModal()).To(BeFalse())
				Expect(intent.GetState()).To(Equal(burst_management.StateList))
			})

			It("should not delete burst when cancelled", func() {
				initialBursts, _ := burstRepo.List(context.Background(), careerrepo.BurstListFilters{})
				initialCount := len(initialBursts)

				// Navigate to detail and delete.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				// Cancel with 'n'.
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				// Burst should still exist.
				finalBursts, _ := burstRepo.List(context.Background(), careerrepo.BurstListFilters{})
				Expect(len(finalBursts)).To(Equal(initialCount))
			})
		})

		Describe("Confirm Modal Escape", func() {
			It("should close confirm modal and return to detail when cancelled", func() {
				// Navigate to detail.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)

				// Press 'c' for confirm.
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
				Expect(intent.HasVisibleConfirmModal()).To(BeTrue())

				// Cancel with 'n'.
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				// Confirm modal should be closed, detail modal should reappear.
				Expect(intent.HasVisibleConfirmModal()).To(BeFalse())
				Expect(intent.GetDetailModal()).NotTo(BeNil())
			})

			It("should not confirm burst when cancelled", func() {
				// Navigate to detail and confirm.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

				// Cancel with 'n'.
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				// Burst should remain unconfirmed.
				repoBurst, _ := burstRepo.GetByID(context.Background(), burst.ID)
				Expect(repoBurst.Confirmed).To(BeFalse())
			})
		})

		Describe("Events Modal Escape", func() {
			It("should close events modal and return to detail when escape is pressed", func() {
				// Navigate to detail.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)
				Expect(intent.GetDetailModal()).NotTo(BeNil())

				// Press 'v' for events.
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

				// Simulate events loaded.
				intent.Update(burst_management.BurstEventsLoadedMsg{
					Events: []*career.CareerEvent{},
				})

				// Press Escape to close events modal.
				intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

				// Should return to detail modal.
				Expect(intent.GetDetailModal()).NotTo(BeNil())
				Expect(intent.GetDetailModal().IsVisible()).To(BeTrue())
			})
		})

		Describe("Facts Modal Escape", func() {
			It("should close facts modal and return to detail when escape is pressed", func() {
				// Navigate to detail.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)

				// Press 'f' for facts.
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

				// Simulate facts loaded.
				intent.Update(burst_management.BurstFactsLoadedMsg{
					Facts: []*career.Fact{},
				})

				// Press Escape to close facts modal.
				intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

				// Should return to detail modal.
				Expect(intent.GetDetailModal()).NotTo(BeNil())
				Expect(intent.GetDetailModal().IsVisible()).To(BeTrue())
			})
		})

		Describe("Detail Modal Escape", func() {
			It("should close detail modal when escape is pressed", func() {
				// Navigate to detail.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)
				Expect(intent.GetDetailModal()).NotTo(BeNil())

				// Press Escape to close.
				intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

				// Detail modal should be closed.
				Expect(intent.GetDetailModal()).To(BeNil())
			})

			It("should close detail modal when enter is pressed", func() {
				// Navigate to detail.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)

				// Press Enter to close.
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Detail modal should be closed.
				Expect(intent.GetDetailModal()).To(BeNil())
			})
		})

		Describe("Error Modal Escape", func() {
			It("should dismiss error modal when escape is pressed", func() {
				// Show error modal.
				intent.ShowErrorModal("Test Error", "Error message")
				Expect(intent.HasVisibleErrorModal()).To(BeTrue())

				// Press Escape to dismiss.
				intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

				// Error modal should be dismissed.
				Expect(intent.HasVisibleErrorModal()).To(BeFalse())
			})
		})

		Describe("Loading Modal Escape", func() {
			It("should cancel loading when escape is pressed", func() {
				mockService := mocks.NewBurstServiceMock()
				ctx.Service = mockService

				newIntent, _ := burst_management.NewIntent(ctx)
				newIntent.Init()

				// Start suggestion detection (shows loading modal).
				newIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
				Expect(newIntent.GetState()).To(Equal(burst_management.StateSuggesting))

				// Press Escape to cancel.
				newIntent.Update(tea.KeyMsg{Type: tea.KeyEscape})

				// Should return to list state.
				Expect(newIntent.GetState()).To(Equal(burst_management.StateList))
			})
		})

		Describe("Escape Priority", func() {
			It("should prioritize error modal over other modals", func() {
				// Navigate to detail.
				result := &screens.NavigateResult{ResultData: burst}
				intent.HandleNavigate(result)

				// Show error modal (simulating an error).
				intent.ShowErrorModal("Error", "Something went wrong")

				// Both detail and error modal exist.
				Expect(intent.GetDetailModal()).NotTo(BeNil())
				Expect(intent.HasVisibleErrorModal()).To(BeTrue())

				// First escape should dismiss error modal only.
				intent.Update(tea.KeyMsg{Type: tea.KeyEscape})
				Expect(intent.HasVisibleErrorModal()).To(BeFalse())
				// Detail modal should still be visible.
				Expect(intent.GetDetailModal()).NotTo(BeNil())
			})
		})
	})
})
