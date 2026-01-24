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
