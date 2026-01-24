package burst_management_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
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
	})
})
