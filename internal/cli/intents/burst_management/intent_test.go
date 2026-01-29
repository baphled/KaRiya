package burst_management_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/mocks"
	tea "github.com/charmbracelet/bubbletea"
)

// executeBatchCmd executes a tea.Cmd and returns all non-tick messages.
// This helper handles both single commands and batch commands.
// For batch commands, it waits for all goroutines to complete.
func executeBatchCmd(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}

	msg := cmd()
	if batchMsg, ok := msg.(tea.BatchMsg); ok {
		return executeBatchCommands(batchMsg)
	}
	return filterTickMessage(msg)
}

// executeBatchCommands executes all commands in a batch and collects results.
func executeBatchCommands(batchMsg tea.BatchMsg) []tea.Msg {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var messages []tea.Msg

	for _, batchCmd := range batchMsg {
		if batchCmd == nil {
			continue
		}
		wg.Add(1)
		go func(c tea.Cmd) {
			defer wg.Done()
			if result := c(); result != nil {
				if filtered := filterTickMessage(result); len(filtered) > 0 {
					mu.Lock()
					messages = append(messages, filtered...)
					mu.Unlock()
				}
			}
		}(batchCmd)
	}

	// Wait for all goroutines to complete with timeout.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All goroutines completed, messages slice is now populated.
	case <-time.After(5 * time.Second):
		// Timeout - return what we have.
	}

	return messages
}

// filterTickMessage returns the message if it's not a spinner tick.
func filterTickMessage(msg tea.Msg) []tea.Msg {
	if msg == nil {
		return nil
	}
	if _, isTick := msg.(feedback.ModalSpinnerTickMsg); isTick {
		return nil
	}
	return []tea.Msg{msg}
}

// executeAsyncCmd is a convenience wrapper that returns the first non-tick message.
func executeAsyncCmd(cmd tea.Cmd) tea.Msg {
	messages := executeBatchCmd(cmd)
	if len(messages) > 0 {
		return messages[0]
	}
	return nil
}

var _ = Describe("Intent Methods", func() {
	Describe("Init", func() {
		It("should initialize the intent successfully", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{
					{ID: "burst-1", Name: "Burst 1"},
				},
			}
			ctx.Validate()

			intent, err := burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			cmd := intent.Init()
			Expect(cmd).To(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
		})

		It("should set selected burst if bursts available", func() {
			bursts := []*career.Burst{
				{ID: "burst-1", Name: "First Burst"},
				{ID: "burst-2", Name: "Second Burst"},
			}
			ctx := &burst_management.IntentContext{
				Bursts: bursts,
			}
			ctx.Validate()

			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()

			Expect(intent.GetSelectedBurst()).NotTo(BeNil())
			Expect(intent.GetSelectedBurst().ID).To(Equal("burst-1"))
		})

		It("should handle empty burst list", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()

			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()

			Expect(intent.GetFilteredBursts()).To(BeEmpty())
			Expect(intent.GetSelectedBurst()).To(BeNil())
		})
	})

	Describe("Update", func() {
		var intent *burst_management.Intent

		BeforeEach(func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{
					{ID: "burst-1", Name: "Burst 1"},
				},
			}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
		})

		It("should return nil when intent is not active", func() {
			intent.Deactivate()
			cmd := intent.Update(tea.KeyMsg{})
			Expect(cmd).To(BeNil())
		})

		It("should handle quit key in list view", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			cmd := intent.Update(msg)
			// Quit is handled by the app, not the intent, so cmd should be nil
			// The screen doesn't handle 'q' - that's a global shortcut
			Expect(cmd).To(BeNil())
		})

		It("should handle escape key to cancel in list view", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			cmd := intent.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(intent.IsActive()).To(BeFalse())
			Expect(intent.Result()).NotTo(BeNil())
		})

		It("should delegate to appropriate state handler", func() {
			intent.SetState(burst_management.StateDetail)
			cmd := intent.Update(tea.KeyMsg{})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("View", func() {
		var intent *burst_management.Intent

		BeforeEach(func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{
					{ID: "burst-1", Name: "Burst 1"},
				},
			}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
		})

		It("should render when active", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show inactive message when not active", func() {
			intent.Deactivate()
			view := intent.View()
			Expect(view).To(ContainSubstring("not active"))
		})

		It("should render list view in StateList", func() {
			intent.SetState(burst_management.StateList)
			view := intent.View()
			// Now renders actual table with burst data, not placeholder text
			Expect(view).To(ContainSubstring("Burst 1"))
			Expect(view).To(ContainSubstring("Name"))
		})

		It("should render detail view in StateDetail", func() {
			intent.SetState(burst_management.StateDetail)
			intent.SetSelectedBurst(&career.Burst{ID: "burst-1", Name: "Test"})
			view := intent.View()
			// StateDetail doesn't have a screen yet, so it falls back to state-based rendering
			// which shows the fallback placeholder
			Expect(view).To(ContainSubstring("Burst"))
		})

		It("should show no burst selected when nil in detail", func() {
			intent.SetState(burst_management.StateDetail)
			intent.SetSelectedBurst(nil)
			view := intent.View()
			// StateDetail doesn't have a screen yet, so it falls back to state-based rendering
			Expect(view).To(ContainSubstring("Burst"))
		})

		It("should show no bursts found when empty in list", func() {
			emptyCtx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			emptyCtx.Validate()
			emptyIntent, _ := burst_management.NewIntent(emptyCtx)
			emptyIntent.Init()

			view := emptyIntent.View()
			Expect(view).To(ContainSubstring("No bursts found"))
		})
	})

	Describe("Modal Rendering with Registry", func() {
		var intent *burst_management.Intent

		BeforeEach(func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{
					{ID: "burst-1", Name: "Test Burst", Description: "Test Description"},
				},
			}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
		})

		Describe("Error Modal", func() {
			It("should have no error modal visible initially", func() {
				Expect(intent.HasVisibleErrorModal()).To(BeFalse())
			})

			It("should show error modal when ShowErrorModal is called", func() {
				intent.ShowErrorModal("Operation Failed", "Something went wrong")

				Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			})

			It("should render error modal content in view", func() {
				intent.ShowErrorModal("Delete Failed", "Database connection error")

				view := intent.View()

				Expect(view).To(ContainSubstring("Delete Failed"))
				Expect(view).To(ContainSubstring("Database connection error"))
			})

			It("should dismiss error modal when user presses Esc", func() {
				intent.ShowErrorModal("Test Error", "Test message")
				Expect(intent.HasVisibleErrorModal()).To(BeTrue())

				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.HasVisibleErrorModal()).To(BeFalse())
			})

			It("should prioritize error modal over delete modal in view", func() {
				// Create both error and delete modals.
				intent.ShowErrorModal("Error Modal", "This is an error")
				// Note: delete modal is created internally during delete flow
				// For now, just verify error modal is visible.

				view := intent.View()
				Expect(view).To(ContainSubstring("Error Modal"))
			})
		})

		Describe("Delete Modal", func() {
			It("should have no delete modal visible initially", func() {
				Expect(intent.HasVisibleDeleteModal()).To(BeFalse())
			})

			// Delete modal is created internally during delete flow.
			// We'll test this through the delete action in handlers.
		})

		Describe("Modal Registry Integration", func() {
			It("should have modal registry initialized", func() {
				registry := intent.GetModalRegistry()
				Expect(registry).NotTo(BeNil())
			})

			It("should rebuild registry when modals change", func() {
				// Initially no modals.
				Expect(intent.HasActiveModal()).To(BeFalse())

				// Show error modal.
				intent.ShowErrorModal("Test", "Message")

				// Rebuild registry.
				intent.RebuildModalRegistry()

				// Should have active modal.
				Expect(intent.HasActiveModal()).To(BeTrue())
			})

			It("should render base view without modals when none active", func() {
				view := intent.View()

				// Should show list view.
				Expect(view).To(ContainSubstring("Test Burst"))
				// Should not have modal content.
				Expect(view).NotTo(ContainSubstring("Error"))
			})

			It("should render modal overlay when modal is active", func() {
				intent.ShowErrorModal("Test Error", "Error message")

				view := intent.View()

				// Should have both base view and modal.
				Expect(view).To(ContainSubstring("Test Error"))
				Expect(view).To(ContainSubstring("Error message"))
			})
		})

		Describe("HasActiveModal using Registry", func() {
			It("should return false when no modals are visible", func() {
				// No modals created.
				Expect(intent.HasActiveModal()).To(BeFalse())
			})

			It("should return true when error modal is visible", func() {
				intent.ShowErrorModal("Error", "Message")

				// HasActiveModal should detect the error modal.
				Expect(intent.HasActiveModal()).To(BeTrue())
			})

			It("should return false after error modal is dismissed", func() {
				intent.ShowErrorModal("Error", "Message")
				Expect(intent.HasActiveModal()).To(BeTrue())

				// Dismiss with Esc.
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Should no longer have active modal.
				Expect(intent.HasActiveModal()).To(BeFalse())
			})

			It("should use registry to check for visible modals", func() {
				// Create error modal.
				intent.ShowErrorModal("Test", "Message")

				// HasActiveModal should return true (rebuilds registry internally).
				Expect(intent.HasActiveModal()).To(BeTrue())

				// Registry should also have visible modal.
				registry := intent.GetModalRegistry()
				Expect(registry.HasVisibleModal()).To(BeTrue())

				// HasActiveModal should match registry state.
				Expect(intent.HasActiveModal()).To(Equal(registry.HasVisibleModal()))
			})

			It("should rebuild registry automatically when checking for active modals", func() {
				// Start with no modals.
				Expect(intent.HasActiveModal()).To(BeFalse())

				// Create error modal directly (simulating internal state change).
				intent.ShowErrorModal("Test", "Message")

				// HasActiveModal should detect the modal without manual rebuild.
				Expect(intent.HasActiveModal()).To(BeTrue())
			})
		})
	})

	Describe("Detail Screen Navigation", func() {
		var intent *burst_management.Intent
		var burst *career.Burst

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test Description",
				EventIDs:    []string{"event-1", "event-2"},
			}
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
		})

		It("should transition to detail screen when burst is selected", func() {
			// Select burst via NavigateResult (simulating screen's Enter handling).
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)

			// State should stay as List (modals don't change state).
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			// Detail modal should be visible.
			Expect(intent.GetDetailModal()).NotTo(BeNil())
			Expect(intent.GetSelectedBurst()).To(Equal(burst))
		})

		It("should return to list when escape pressed in detail view", func() {
			// Navigate to detail modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			// Press escape to close detail modal.
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Detail modal should be closed.
			Expect(intent.GetDetailModal()).To(BeNil())
			// Should still be in list state.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should handle view_events action from detail screen", func() {
			// Navigate to detail modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			// Press 'v' for view events (triggers async load).
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			// State stays as List (modals don't change state).
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			// Detail modal still exists but is hidden (waiting for events to load).
			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})

		It("should handle view_facts action from detail screen", func() {
			// Navigate to detail modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			// Press 'f' for view facts (triggers async load).
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// State stays as List (modals don't change state).
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			// Detail modal still exists but is hidden (waiting for facts to load).
			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})

		It("should handle edit action from detail screen", func() {
			// Navigate to detail modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			// Press 'e' for edit.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// State stays as List.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			// Edit modal should be visible.
			Expect(intent.HasVisibleEditModal()).To(BeTrue())
		})

		It("should handle delete action from detail screen", func() {
			// Navigate to detail modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			// Press 'd' for delete.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			// State stays as List.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			// Delete confirmation modal should be visible.
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())
		})
	})

	Describe("Delete Burst Flow", func() {
		var intent *burst_management.Intent
		var burst *career.Burst

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test Description",
			}
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
		})

		It("should delete burst when 'y' is pressed in delete confirm state", func() {
			// Navigate to detail modal and then delete.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())

			// Press 'y' to confirm deletion.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Burst should be deleted - filtered bursts should be empty.
			Expect(intent.GetFilteredBursts()).To(HaveLen(0))
			// Should return to list view with no modals.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleDeleteModal()).To(BeFalse())
		})

		It("should cancel deletion when 'n' is pressed", func() {
			// Navigate to detail modal and then delete.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())

			// Press 'n' to cancel.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			// Burst should still exist.
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
			// Should return to list view with no delete modal.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleDeleteModal()).To(BeFalse())
		})

		It("should return to list view on escape from delete confirm", func() {
			// Navigate to detail modal and then delete.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())

			// Press escape.
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should close modal and return to list view.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleDeleteModal()).To(BeFalse())
		})

		It("should handle confirm action from detail screen", func() {
			// Navigate to detail modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			// Press 'c' for confirm.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			// State stays as List.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			// Confirm modal should be visible.
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())
		})
	})

	Describe("Confirm Burst Flow", func() {
		var intent *burst_management.Intent
		var burst *career.Burst

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test Description",
				Confirmed:   false, // Not confirmed yet
			}
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
		})

		It("should mark burst as confirmed when 'y' is pressed in confirm state", func() {
			// Navigate to detail modal and then confirm.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())

			// Press 'y' to confirm.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Burst should now be confirmed.
			Expect(intent.GetSelectedBurst().Confirmed).To(BeTrue())
			Expect(cmd).NotTo(BeNil(), "confirmBurst should return extraction command")
		})

		It("should cancel confirmation when 'n' is pressed", func() {
			// Navigate to confirm modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())

			// Press 'n' to cancel.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			// Burst should remain unconfirmed.
			Expect(intent.GetSelectedBurst().Confirmed).To(BeFalse())
			// Should show detail modal again.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})

		It("should return to detail view on escape", func() {
			// Navigate to confirm modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())

			// Press escape.
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should close confirm modal and show detail modal again.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleConfirmModal()).To(BeFalse())
			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})
	})

	Describe("Edit Burst Flow", func() {
		var intent *burst_management.Intent
		var burst *career.Burst

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Original Name",
				Description: "Original Description",
			}
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
		})

		It("should transition to edit state when 'e' is pressed on detail screen", func() {
			// Navigate to detail modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			// Press 'e' to edit.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// State stays as List with edit modal visible.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleEditModal()).To(BeTrue())
		})

		It("should update burst name and description on form submit", func() {
			// Navigate to detail modal and edit.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleEditModal()).To(BeTrue())

			// Simulate form completion with new values.
			// The edit modal will send an EditBurstMsg with updated values.
			editMsg := burst_management.EditBurstMsg{
				BurstID:     "burst-1",
				Name:        "Updated Name",
				Description: "Updated Description",
			}
			intent.Update(editMsg)

			// Burst should be updated.
			Expect(intent.GetSelectedBurst().Name).To(Equal("Updated Name"))
			Expect(intent.GetSelectedBurst().Description).To(Equal("Updated Description"))

			// Should show detail modal again with updated burst.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleEditModal()).To(BeFalse())
			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})

		It("should cancel edit and return to detail on escape", func() {
			originalName := burst.Name
			originalDesc := burst.Description

			// Navigate to edit modal.
			result := &screens.NavigateResult{
				ResultData: burst,
			}
			intent.HandleNavigate(result)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleEditModal()).To(BeTrue())

			// Press escape to cancel.
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Burst should remain unchanged.
			Expect(intent.GetSelectedBurst().Name).To(Equal(originalName))
			Expect(intent.GetSelectedBurst().Description).To(Equal(originalDesc))

			// Should close edit modal and return to list view.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleEditModal()).To(BeFalse())
		})

		It("should handle validation errors in edit form", func() {
			// Navigate to edit state.
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Attempt to submit with empty name.
			editMsg := burst_management.EditBurstMsg{
				BurstID:     "burst-1",
				Name:        "", // Invalid - empty name
				Description: "Some description",
			}
			intent.Update(editMsg)

			// Should remain in edit state with error visible.
			Expect(intent.GetState()).To(Equal(burst_management.StateEdit))
			// Modal should still be visible (form has validation error).
			Expect(intent.HasVisibleEditModal()).To(BeTrue())
		})
	})

	Describe("Suggestion Modal Message Flow", func() {
		It("should create modal when suggestions are loaded", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()

			// Load suggestions.
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{
						Name:            "Suggestion 1",
						Description:     "First suggestion",
						EventIDs:        []string{"e1"},
						ConfidenceScore: 0.9,
					},
				},
				Error: nil,
			}

			intent.Update(msg)

			// Verify modal is created and visible.
			Expect(intent.GetSuggestionModal()).NotTo(BeNil())
			Expect(intent.GetSuggestionModal().IsVisible()).To(BeTrue())
		})

		It("should forward key to modal when modal is visible", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()

			// Load suggestions.
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test", EventIDs: []string{"e1"}},
				},
				Error: nil,
			}
			intent.Update(msg)

			// Send key to intent.
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			intent.Update(keyMsg)

			// Modal should have processed the key and closed.
			Expect(intent.GetSuggestionModal()).To(BeNil())
		})

		It("should send SuggestionReviewCompleteMsg when modal closes with accepted suggestions", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()

			initialCount := len(intent.GetFilteredBursts())

			// Load suggestions.
			loadMsg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{
						Name:        "Accepted Burst",
						Description: "Test burst",
						EventIDs:    []string{"e1"},
					},
				},
				Error: nil,
			}
			intent.Update(loadMsg)

			// Accept suggestion (modal closes and sends message).
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// Process the completion message (in tests, we need to send it manually).
			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{loadMsg.Suggestions[0]},
				Cancelled:           false,
			}
			intent.Update(completeMsg)

			// Verify burst was created.
			Expect(intent.GetFilteredBursts()).To(HaveLen(initialCount + 1))
			Expect(intent.GetFilteredBursts()[0].Name).To(Equal("Accepted Burst"))
		})

		It("should clear modal when cancelled", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()

			// Load suggestions.
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test", EventIDs: []string{"e1"}},
				},
				Error: nil,
			}
			intent.Update(msg)

			// Cancel modal.
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Modal should be cleared.
			Expect(intent.GetSuggestionModal()).To(BeNil())
		})

		It("should not create bursts when review is cancelled", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()

			initialCount := len(intent.GetFilteredBursts())

			// Load suggestions.
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test", EventIDs: []string{"e1"}},
				},
				Error: nil,
			}
			intent.Update(msg)

			// Cancel and process completion.
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			completeMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: nil,
				Cancelled:           true,
			}
			intent.Update(completeMsg)

			// No bursts should be created.
			Expect(intent.GetFilteredBursts()).To(HaveLen(initialCount))
		})
	})

	Describe("View Rendering", func() {
		var (
			intent *burst_management.Intent
			ctx    *burst_management.IntentContext
		)

		BeforeEach(func() {
			testBurst := &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test description",
				EventIDs:    []string{"e1", "e2"},
			}

			ctx = &burst_management.IntentContext{
				Bursts: []*career.Burst{testBurst},
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should render inactive message when not active", func() {
			intent.Deactivate()
			view := intent.View()
			Expect(view).To(ContainSubstring("not active"))
		})

		It("should render view in StateList", func() {
			intent.SetState(burst_management.StateList)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateDetail with selected burst", func() {
			intent.SetState(burst_management.StateDetail)
			intent.SetSelectedBurst(ctx.Bursts[0])
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateDetailEvents", func() {
			intent.SetState(burst_management.StateDetailEvents)
			intent.SetSelectedBurst(ctx.Bursts[0])
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateDetailFacts", func() {
			intent.SetState(burst_management.StateDetailFacts)
			intent.SetSelectedBurst(ctx.Bursts[0])
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateEdit", func() {
			intent.SetState(burst_management.StateEdit)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateDeleteConfirm", func() {
			intent.SetState(burst_management.StateDeleteConfirm)
			intent.SetSelectedBurst(ctx.Bursts[0])
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateConfirm", func() {
			intent.SetState(burst_management.StateConfirm)
			intent.SetSelectedBurst(ctx.Bursts[0])
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateExtractingFacts", func() {
			intent.SetState(burst_management.StateExtractingFacts)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateSuggesting", func() {
			intent.SetState(burst_management.StateSuggesting)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render view in StateSuggestionReview", func() {
			intent.SetState(burst_management.StateSuggestionReview)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		Context("state-specific content fallback", func() {
			// These tests force the getStateContent path by having no activeScreen.

			It("should show 'No bursts found' when list is empty", func() {
				emptyCtx := &burst_management.IntentContext{
					Bursts: []*career.Burst{},
				}
				emptyCtx.Validate()
				emptyIntent, err := burst_management.NewIntent(emptyCtx)
				Expect(err).NotTo(HaveOccurred())
				emptyIntent.Init()

				// Clear active screen to force fallback.
				emptyIntent.SetState(burst_management.StateList)
				view := emptyIntent.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle StateDetail with no selected burst", func() {
				intent.SetState(burst_management.StateDetail)
				intent.SetSelectedBurst(nil)
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("Handler Functions", func() {
		var (
			intent *burst_management.Intent
			ctx    *burst_management.IntentContext
		)

		BeforeEach(func() {
			testBurst := &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test description",
				EventIDs:    []string{"e1", "e2"},
			}

			ctx = &burst_management.IntentContext{
				Bursts: []*career.Burst{testBurst},
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should handle HandleSubmit gracefully", func() {
			result := &screens.SubmitResult{
				FormData: map[string]interface{}{"test": "data"},
			}
			cmd := intent.HandleSubmit(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle HandleError and show error modal", func() {
			result := &screens.ErrorResult{
				Err: errors.New("test error"),
			}
			cmd := intent.HandleError(result)
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should return extracted facts count", func() {
			// The count should be 0 initially.
			count := intent.GetExtractedFactsCount()
			Expect(count).To(Equal(0))
		})
	})

	Describe("Action Data Handling via HandleNavigate", func() {
		var (
			intent *burst_management.Intent
			ctx    *burst_management.IntentContext
			burst  *career.Burst
		)

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test description",
				EventIDs:    []string{"e1", "e2"},
			}

			ctx = &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		// Note: "add" action test removed - bursts are created via AI suggestions only.

		It("should handle 'edit' action with burst and open edit modal", func() {
			actionData := map[string]interface{}{
				"action": "edit",
				"burst":  burst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateEdit))
			Expect(intent.GetSelectedBurst()).To(Equal(burst))
			Expect(intent.HasVisibleEditModal()).To(BeTrue())
		})

		It("should handle 'edit' action without burst gracefully", func() {
			actionData := map[string]interface{}{
				"action": "edit",
				// No burst provided.
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle 'delete' action with burst and open delete modal", func() {
			actionData := map[string]interface{}{
				"action": "delete",
				"burst":  burst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateDeleteConfirm))
			Expect(intent.GetSelectedBurst()).To(Equal(burst))
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())
		})

		It("should handle 'delete' action with long burst name (truncated)", func() {
			longNameBurst := &career.Burst{
				ID:          "burst-long",
				Name:        "This is a very long burst name that exceeds fifty characters and should be truncated",
				Description: "Test",
				EventIDs:    []string{"e1", "e2"},
			}
			actionData := map[string]interface{}{
				"action": "delete",
				"burst":  longNameBurst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateDeleteConfirm))
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())
		})

		It("should handle 'delete' action without burst gracefully", func() {
			actionData := map[string]interface{}{
				"action": "delete",
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle 'suggest' action and attempt burst detection", func() {
			actionData := map[string]interface{}{
				"action": "suggest",
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)
			// Without a service, startBurstDetection shows error and returns to list.
			// The state is set to StateSuggesting first, then reset to StateList on error.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			// Error modal should be visible since no service is available.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should handle 'view_events' action with burst", func() {
			actionData := map[string]interface{}{
				"action": "view_events",
				"burst":  burst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateDetailEvents))
			Expect(intent.GetSelectedBurst()).To(Equal(burst))
		})

		It("should handle 'view_events' action without burst gracefully", func() {
			actionData := map[string]interface{}{
				"action": "view_events",
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle 'view_facts' action with burst", func() {
			actionData := map[string]interface{}{
				"action": "view_facts",
				"burst":  burst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateDetailFacts))
			Expect(intent.GetSelectedBurst()).To(Equal(burst))
		})

		It("should handle 'view_facts' action without burst gracefully", func() {
			actionData := map[string]interface{}{
				"action": "view_facts",
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle 'confirm' action with burst and open confirm modal", func() {
			actionData := map[string]interface{}{
				"action": "confirm",
				"burst":  burst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateConfirm))
			Expect(intent.GetSelectedBurst()).To(Equal(burst))
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())
		})

		It("should handle 'confirm' action with long burst name (truncated)", func() {
			longNameBurst := &career.Burst{
				ID:          "burst-long",
				Name:        "This is a very long burst name that exceeds fifty characters and should be truncated",
				Description: "Test",
				EventIDs:    []string{"e1", "e2"},
			}
			actionData := map[string]interface{}{
				"action": "confirm",
				"burst":  longNameBurst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateConfirm))
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())
		})

		It("should handle 'confirm' action without burst gracefully", func() {
			actionData := map[string]interface{}{
				"action": "confirm",
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle unknown action gracefully", func() {
			actionData := map[string]interface{}{
				"action": "unknown_action",
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle action data without action key gracefully", func() {
			actionData := map[string]interface{}{
				"burst": burst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("View Fallback Rendering", func() {
		It("should render fallback view when activeScreen is nil", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()

			intent, err := burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			// Activate intent without calling Init() to have no activeScreen.
			intent.SetState(burst_management.StateList)

			// View should still render (fallback path).
			view := intent.View()
			Expect(view).To(ContainSubstring("Loading..."))
		})
	})

	Describe("HandleCancel State Transitions", func() {
		var (
			intent *burst_management.Intent
			ctx    *burst_management.IntentContext
			burst  *career.Burst
		)

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test description",
				EventIDs:    []string{"e1", "e2"},
			}

			ctx = &burst_management.IntentContext{
				Bursts: []*career.Burst{burst},
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should cancel intent when in StateList with no active modals", func() {
			intent.SetState(burst_management.StateList)
			result := &screens.CancelResult{}
			intent.HandleCancel(result)
			// Intent should be cancelled.
			intentResult := intent.Result()
			Expect(intentResult).NotTo(BeNil())
			Expect(intentResult.IsCancelled()).To(BeTrue())
		})

		It("should ignore cancel when in StateList with detail modal active", func() {
			intent.SetState(burst_management.StateList)
			// Navigate to detail to activate detail modal.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			// Try to cancel - should be ignored.
			cancelResult := &screens.CancelResult{}
			intent.HandleCancel(cancelResult)

			// Intent should NOT be cancelled (modal is active).
			// Result is nil when intent hasn't been completed/cancelled.
			intentResult := intent.Result()
			Expect(intentResult).To(BeNil())
		})

		It("should return to list from StateDetail", func() {
			intent.SetState(burst_management.StateDetail)
			result := &screens.CancelResult{}
			intent.HandleCancel(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should return to list from StateDetailEvents", func() {
			intent.SetState(burst_management.StateDetailEvents)
			result := &screens.CancelResult{}
			intent.HandleCancel(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should return to list from StateDetailFacts", func() {
			intent.SetState(burst_management.StateDetailFacts)
			result := &screens.CancelResult{}
			intent.HandleCancel(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should return to list and clear delete modal from StateDeleteConfirm", func() {
			// Set up delete state with modal.
			intent.SetState(burst_management.StateDeleteConfirm)
			actionData := map[string]interface{}{
				"action": "delete",
				"burst":  burst,
			}
			navResult := &screens.NavigateResult{ResultData: actionData}
			intent.HandleNavigate(navResult)
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())

			// Now cancel.
			cancelResult := &screens.CancelResult{}
			intent.HandleCancel(cancelResult)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleDeleteModal()).To(BeFalse())
		})

		It("should return to list from default/unknown state", func() {
			// Use StateEdit as an example of a state that falls into default.
			intent.SetState(burst_management.StateEdit)
			result := &screens.CancelResult{}
			intent.HandleCancel(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should ignore cancel when loadingEvents is true", func() {
			intent.SetState(burst_management.StateList)
			// Directly set loading state for testing.
			// This simulates the state when events are being loaded asynchronously.
			intent.SetLoadingEventsForTesting(true)

			// Try to cancel.
			cancelResult := &screens.CancelResult{}
			intent.HandleCancel(cancelResult)

			// Intent should NOT be cancelled (loading in progress).
			// Result is nil when intent hasn't been completed/cancelled.
			intentResult := intent.Result()
			Expect(intentResult).To(BeNil())
		})
	})

	Describe("Context CRUD Operations", func() {
		var ctx *burst_management.IntentContext

		BeforeEach(func() {
			ctx = &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()
		})

		It("should handle CreateBurst without repository", func() {
			burst := &career.Burst{
				ID:          "new-burst",
				Name:        "New Burst",
				Description: "Test",
				EventIDs:    []string{"e1", "e2"},
			}
			err := ctx.CreateBurst(burst)
			// Without repository, should return nil (no-op).
			Expect(err).To(BeNil())
		})

		It("should handle UpdateBurst without repository", func() {
			burst := &career.Burst{
				ID:          "burst-1",
				Name:        "Updated Burst",
				Description: "Test",
				EventIDs:    []string{"e1", "e2"},
			}
			err := ctx.UpdateBurst(burst)
			// Without repository, should return nil (no-op).
			Expect(err).To(BeNil())
		})

		It("should handle DeleteBurst without repository", func() {
			err := ctx.DeleteBurst("burst-1")
			// Without repository, should return nil (no-op).
			Expect(err).To(BeNil())
		})
	})

	Describe("Service Integration Tests", func() {
		var (
			intent      *burst_management.Intent
			ctx         *burst_management.IntentContext
			mockService *mocks.BurstServiceMock
			burst       *career.Burst
			events      []*career.Event
		)

		BeforeEach(func() {
			events = []*career.Event{
				{ID: "e1", Text: "First event text"},
				{ID: "e2", Text: "Second event text"},
			}

			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test description",
				EventIDs:    []string{"e1", "e2"},
			}

			mockService = mocks.NewBurstServiceMock().SetEvents(events)

			ctx = &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should load burst events via view_events action with service", func() {
			actionData := map[string]interface{}{
				"action": "view_events",
				"burst":  burst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)

			Expect(intent.GetState()).To(Equal(burst_management.StateDetailEvents))
			Expect(intent.GetSelectedBurst()).To(Equal(burst))
		})

		It("should load burst facts via view_facts action with service", func() {
			// Set up facts for the burst.
			testFacts := []*career.Fact{
				{ID: "f1", Text: "Fact 1"},
				{ID: "f2", Text: "Fact 2"},
			}
			mockService.SetFactsForBurst(burst.ID, testFacts)

			actionData := map[string]interface{}{
				"action": "view_facts",
				"burst":  burst,
			}
			result := &screens.NavigateResult{
				ResultData: actionData,
			}
			intent.HandleNavigate(result)

			Expect(intent.GetState()).To(Equal(burst_management.StateDetailFacts))
			Expect(intent.GetSelectedBurst()).To(Equal(burst))
		})

		It("should start burst detection with service and return suggestions", func() {
			suggestions := []burstfact.BurstSuggestion{
				{Name: "Suggested Burst", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
			}
			mockService.SetSuggestions(suggestions)

			// Trigger via 's' key shortcut.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			Expect(intent.GetState()).To(Equal(burst_management.StateSuggesting))
		})

		It("should handle BurstEventsLoadedMsg with events", func() {
			// First navigate to detail to set selectedBurst.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			// Press 'v' to trigger events loading.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			// Simulate async message returned.
			msg := burst_management.BurstEventsLoadedMsg{
				Events: events,
			}
			intent.Update(msg)

			// Events modal should be visible.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should handle BurstFactsLoadedMsg with facts", func() {
			testFacts := []*career.Fact{
				{ID: "f1", Text: "Fact 1"},
			}

			// First navigate to detail to set selectedBurst.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			// Press 'f' to trigger facts loading.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Simulate async message returned.
			msg := burst_management.BurstFactsLoadedMsg{
				Facts: testFacts,
			}
			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should handle fact extraction complete with facts", func() {
			extractedFacts := []*career.Fact{
				{ID: "f1", Text: "Extracted fact 1"},
				{ID: "f2", Text: "Extracted fact 2"},
			}

			msg := burst_management.FactExtractionCompleteMsg{
				Facts: extractedFacts,
				Error: nil,
			}
			intent.Update(msg)

			Expect(intent.GetExtractedFactsCount()).To(Equal(2))
		})

		It("should handle fact extraction error", func() {
			msg := burst_management.FactExtractionCompleteMsg{
				Facts: nil,
				Error: errors.New("extraction failed"),
			}
			intent.Update(msg)

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should handle EditBurstMsg with valid update", func() {
			mockRepo := careermemory.NewBurstRepository()
			// Add the burst to the repo first.
			_ = mockRepo.Create(ctx.Context, burst)

			ctx.BurstRepository = mockRepo

			// Create new intent with repo.
			newIntent, _ := burst_management.NewIntent(ctx)
			newIntent.Init()

			// Navigate to detail and then edit.
			navResult := &screens.NavigateResult{ResultData: burst}
			newIntent.HandleNavigate(navResult)
			newIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Send edit message.
			editMsg := burst_management.EditBurstMsg{
				BurstID:     burst.ID,
				Name:        "Updated Name",
				Description: "Updated Description",
			}
			newIntent.Update(editMsg)

			// Verify the burst was updated.
			Expect(newIntent.GetSelectedBurst().Name).To(Equal("Updated Name"))
		})
	})

	Describe("Context CRUD with Repository", func() {
		var (
			ctx  *burst_management.IntentContext
			repo *careermemory.BurstRepository
		)

		BeforeEach(func() {
			repo = careermemory.NewBurstRepository()
			ctx = &burst_management.IntentContext{
				Bursts:          []*career.Burst{},
				BurstRepository: repo,
			}
			ctx.Validate()
		})

		It("should create burst with repository", func() {
			burst := &career.Burst{
				ID:          "new-burst",
				Name:        "New Burst",
				Description: "Test burst",
				EventIDs:    []string{"e1", "e2"},
			}
			err := ctx.CreateBurst(burst)
			Expect(err).To(BeNil())

			// Verify it was created.
			bursts, _ := repo.List(ctx.Context, careerrepo.BurstListFilters{})
			Expect(bursts).To(HaveLen(1))
			Expect(bursts[0].Name).To(Equal("New Burst"))
		})

		It("should fail CreateBurst with invalid burst", func() {
			burst := &career.Burst{
				ID:          "",
				Name:        "", // Invalid - name required.
				Description: "",
				EventIDs:    []string{},
			}
			err := ctx.CreateBurst(burst)
			Expect(err).NotTo(BeNil())
		})

		It("should update burst with repository", func() {
			// Create first.
			burst := &career.Burst{
				ID:          "burst-1",
				Name:        "Original Name",
				Description: "Original",
				EventIDs:    []string{"e1", "e2"},
			}
			_ = repo.Create(ctx.Context, burst)

			// Update.
			burst.Name = "Updated Name"
			err := ctx.UpdateBurst(burst)
			Expect(err).To(BeNil())

			// Verify update.
			updated, _ := repo.GetByID(ctx.Context, "burst-1")
			Expect(updated.Name).To(Equal("Updated Name"))
		})

		It("should fail UpdateBurst with invalid burst", func() {
			burst := &career.Burst{
				ID:          "burst-1",
				Name:        "", // Invalid - name required.
				Description: "",
				EventIDs:    []string{},
			}
			err := ctx.UpdateBurst(burst)
			Expect(err).NotTo(BeNil())
		})

		It("should delete burst with repository", func() {
			// Create first.
			burst := &career.Burst{
				ID:          "burst-1",
				Name:        "To Delete",
				Description: "Test",
				EventIDs:    []string{"e1", "e2"},
			}
			_ = repo.Create(ctx.Context, burst)

			// Delete.
			err := ctx.DeleteBurst("burst-1")
			Expect(err).To(BeNil())

			// Verify deletion.
			bursts, _ := repo.List(ctx.Context, careerrepo.BurstListFilters{})
			Expect(bursts).To(HaveLen(0))
		})

		It("should load bursts from repository", func() {
			// Create some bursts.
			_ = repo.Create(ctx.Context, &career.Burst{ID: "b1", Name: "Burst 1", EventIDs: []string{"e1", "e2"}})
			_ = repo.Create(ctx.Context, &career.Burst{ID: "b2", Name: "Burst 2", EventIDs: []string{"e3", "e4"}})

			err := ctx.LoadBursts()
			Expect(err).To(BeNil())
			Expect(ctx.Bursts).To(HaveLen(2))
		})
	})

	Describe("Delete and Confirm Flows with Repository", func() {
		var (
			intent      *burst_management.Intent
			ctx         *burst_management.IntentContext
			mockService *mocks.BurstServiceMock
			repo        *careermemory.BurstRepository
			burst       *career.Burst
		)

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test description",
				EventIDs:    []string{"e1", "e2"},
			}

			repo = careermemory.NewBurstRepository()
			_ = repo.Create(context.Background(), burst)

			mockService = mocks.NewBurstServiceMock()

			ctx = &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				Service:         mockService,
				BurstRepository: repo,
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should delete burst through confirm modal", func() {
			// Navigate to detail and delete.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			// Confirm deletion with 'y'.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Burst should be deleted from repository.
			bursts, _ := repo.List(ctx.Context, careerrepo.BurstListFilters{})
			Expect(bursts).To(HaveLen(0))
		})

		It("should handle delete error gracefully", func() {
			// Remove burst from repo to simulate delete error.
			_ = repo.Delete(ctx.Context, burst.ID)

			// Navigate to detail and delete.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			// Try to confirm deletion - should fail.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Should show error modal.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should confirm burst and trigger fact extraction", func() {
			extractedFacts := []career.Fact{
				{ID: "f1", Text: "Extracted fact", SourceBurstID: burst.ID,
					CompetencyCategories: []string{"technical"},
					RoleFit:              "senior_ic", AudienceRelevance: []string{"peer"},
					StrengthSignal: "technical"},
			}
			mockService.SetExtractedFacts(extractedFacts)

			// Navigate to detail and confirm.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			_ = cmd

			// Confirm with 'y' - this should trigger both confirmation AND fact extraction.
			cmd = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Burst should be confirmed.
			confirmed, _ := repo.GetByID(ctx.Context, burst.ID)
			Expect(confirmed.Confirmed).To(BeTrue())

			// Fact extraction should have been triggered (cmd returned).
			Expect(cmd).NotTo(BeNil(), "confirmBurst should return a command for fact extraction")

			// Execute the returned command to trigger fact extraction.
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					// Process the extraction result message.
					intent.Update(msg)
				}
			}

			// Verify extraction was called.
			Expect(mockService.GetExtractCallCount()).To(BeNumerically(">=", 1),
				"fact extraction should be triggered when confirming a burst")
		})

		It("should handle confirm error gracefully", func() {
			// Configure service to fail on confirm.
			mockService.SetConfirmError(fmt.Errorf("database unavailable"))

			// Navigate to detail and confirm.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			// Try to confirm - should fail.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Should show error modal.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("Async Command Execution", func() {
		var (
			intent      *burst_management.Intent
			ctx         *burst_management.IntentContext
			mockService *mocks.BurstServiceMock
			burst       *career.Burst
			events      []*career.Event
		)

		BeforeEach(func() {
			events = []*career.Event{
				{ID: "e1", Text: "Event 1"},
				{ID: "e2", Text: "Event 2"},
			}

			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test",
				EventIDs:    []string{"e1", "e2"},
			}

			mockService = mocks.NewBurstServiceMock().SetEvents(events)

			ctx = &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should execute showBurstEventsModal command", func() {
			// Navigate to detail.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			// Press 'v' to show events modal.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command to get the message.
			msg := cmd()
			Expect(msg).NotTo(BeNil())

			// Should be a BurstEventsLoadedMsg.
			eventsMsg, ok := msg.(burst_management.BurstEventsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(eventsMsg.Events).To(HaveLen(2))
		})

		It("should execute showBurstFactsModal command", func() {
			testFacts := []*career.Fact{
				{ID: "f1", Text: "Fact 1"},
			}
			mockService.SetFactsForBurst(burst.ID, testFacts)

			// Navigate to detail.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			// Press 'f' to show facts modal.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command to get the message.
			msg := cmd()
			Expect(msg).NotTo(BeNil())

			// Should be a BurstFactsLoadedMsg.
			factsMsg, ok := msg.(burst_management.BurstFactsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(factsMsg.Facts).To(HaveLen(1))
		})

		It("should execute startBurstDetection command", func() {
			// Add extra events that are NOT in any existing burst.
			extraEvents := []*career.Event{
				{ID: "e1", Text: "Event 1"},
				{ID: "e2", Text: "Event 2"},
				{ID: "e3", Text: "Event 3"},
				{ID: "e4", Text: "Event 4"},
			}
			mockService.SetEvents(extraEvents)

			suggestions := []burstfact.BurstSuggestion{
				{Name: "Suggestion", EventIDs: []string{"e3", "e4"}, ConfidenceScore: 0.8},
			}
			mockService.SetSuggestions(suggestions)

			// Press 's' to start detection.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil(), "Command should not be nil - 's' key should trigger startBurstDetection")

			// Execute the command to get the message (handles batch commands).
			msg := executeAsyncCmd(cmd)
			Expect(msg).NotTo(BeNil())

			// Should be a BurstSuggestionsLoadedMsg with suggestions from unassigned events.
			suggestionsMsg, ok := msg.(burst_management.BurstSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(suggestionsMsg.Suggestions).To(HaveLen(1))
			Expect(suggestionsMsg.Error).To(BeNil())
		})

		It("should handle startBurstDetection with list events error", func() {
			mockService.SetListEventsError(errors.New("list events failed"))

			// Press 's' to start detection.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command (handles batch commands).
			msg := executeAsyncCmd(cmd)

			// Should be error message.
			suggestionsMsg, ok := msg.(burst_management.BurstSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(suggestionsMsg.Error).NotTo(BeNil())
		})

		It("should handle startBurstDetection with suggest error", func() {
			mockService.SetSuggestError(errors.New("suggest failed"))

			// Press 's' to start detection.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command (handles batch commands).
			msg := executeAsyncCmd(cmd)

			// Should be error message.
			suggestionsMsg, ok := msg.(burst_management.BurstSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(suggestionsMsg.Error).NotTo(BeNil())
		})

		It("should handle startBurstDetection with no events", func() {
			mockService.SetEvents([]*career.Event{})

			// Press 's' to start detection.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command (handles batch commands).
			msg := executeAsyncCmd(cmd)

			// Should be error message (no unassigned events available).
			suggestionsMsg, ok := msg.(burst_management.BurstSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(suggestionsMsg.Error).NotTo(BeNil())
			Expect(suggestionsMsg.Error.Error()).To(ContainSubstring("no unassigned events"))
		})

		It("should filter out events that are already in existing bursts", func() {
			// Create events where e1 and e2 are already in an existing burst.
			events := []*career.Event{
				{ID: "e1", Text: "Event 1"},
				{ID: "e2", Text: "Event 2"},
				{ID: "e3", Text: "Event 3"},
				{ID: "e4", Text: "Event 4"},
			}
			mockService.SetEvents(events)

			// Create an existing burst that uses e1 and e2.
			existingBurst := &career.Burst{
				ID:       "existing-burst",
				Name:     "Existing Burst",
				EventIDs: []string{"e1", "e2"},
			}

			// Update context with existing burst.
			ctx.Bursts = []*career.Burst{existingBurst}

			suggestions := []burstfact.BurstSuggestion{
				{Name: "New Suggestion", EventIDs: []string{"e3", "e4"}, ConfidenceScore: 0.8},
			}
			mockService.SetSuggestions(suggestions)

			// Press 's' to start detection.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command (handles batch commands).
			msg := executeAsyncCmd(cmd)

			// Should succeed with suggestions from unassigned events only.
			suggestionsMsg, ok := msg.(burst_management.BurstSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(suggestionsMsg.Error).To(BeNil())
			Expect(suggestionsMsg.Suggestions).To(HaveLen(1))
		})

		It("should show error when all events are already in bursts", func() {
			// All events are already in an existing burst.
			events := []*career.Event{
				{ID: "e1", Text: "Event 1"},
				{ID: "e2", Text: "Event 2"},
			}
			mockService.SetEvents(events)

			// Create an existing burst that uses all events.
			existingBurst := &career.Burst{
				ID:       "existing-burst",
				Name:     "Existing Burst",
				EventIDs: []string{"e1", "e2"},
			}

			// Update context with existing burst.
			ctx.Bursts = []*career.Burst{existingBurst}

			// Press 's' to start detection.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command (handles batch commands).
			msg := executeAsyncCmd(cmd)

			// Should be error message (no unassigned events).
			suggestionsMsg, ok := msg.(burst_management.BurstSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(suggestionsMsg.Error).NotTo(BeNil())
			Expect(suggestionsMsg.Error.Error()).To(ContainSubstring("no unassigned events"))
		})
	})

	Describe("Modal Update Edge Cases", func() {
		var (
			intent      *burst_management.Intent
			ctx         *burst_management.IntentContext
			mockService *mocks.BurstServiceMock
			burst       *career.Burst
		)

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test",
				EventIDs:    []string{"e1", "e2"},
			}

			mockService = mocks.NewBurstServiceMock()

			ctx = &burst_management.IntentContext{
				Bursts:  []*career.Burst{burst},
				Service: mockService,
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should consume messages when error modal is visible", func() {
			// Show error modal.
			intent.ShowErrorModal("Test", "Error")
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())

			// Any other key should be consumed by error modal.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// Error modal should still be visible.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should cancel loading modal with Esc", func() {
			// Navigate to detail.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			// Press 's' to start detection - this shows loading modal.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			// Press Esc to cancel loading.
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			// Should return to list state.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should dismiss error modal with Esc", func() {
			intent.ShowErrorModal("Test Error", "Error message")
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())

			// Press Esc to dismiss.
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			Expect(intent.HasVisibleErrorModal()).To(BeFalse())
		})

		It("should return to detail after cancelling events modal", func() {
			// Navigate to detail.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			// Press 'v' for events.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			// Simulate events loaded.
			eventsMsg := burst_management.BurstEventsLoadedMsg{
				Events: []*career.Event{},
			}
			intent.Update(eventsMsg)

			// Now close events modal with Esc or Enter.
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			// Should return to detail modal.
			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})

		It("should return to detail after cancelling facts modal", func() {
			// Navigate to detail.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)

			// Press 'f' for facts.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Simulate facts loaded.
			factsMsg := burst_management.BurstFactsLoadedMsg{
				Facts: []*career.Fact{},
			}
			intent.Update(factsMsg)

			// Now close facts modal with Enter.
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Should return to detail modal.
			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})

		It("should close detail modal with Enter", func() {
			// Navigate to detail.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			Expect(intent.GetDetailModal()).NotTo(BeNil())

			// Press Enter to close.
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.GetDetailModal()).To(BeNil())
		})

		It("should handle edit burst msg with burst ID not found", func() {
			repo := careermemory.NewBurstRepository()
			_ = repo.Create(ctx.Context, burst)
			ctx.BurstRepository = repo

			newIntent, _ := burst_management.NewIntent(ctx)
			newIntent.Init()

			// Navigate to detail and edit.
			navResult := &screens.NavigateResult{ResultData: burst}
			newIntent.HandleNavigate(navResult)
			newIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Send edit with wrong burst ID.
			editMsg := burst_management.EditBurstMsg{
				BurstID:     "non-existent-id",
				Name:        "Updated",
				Description: "Test",
			}
			newIntent.Update(editMsg)

			// Should show error since burst not found.
			Expect(newIntent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should cancel edit modal with Esc", func() {
			// Navigate to detail and edit.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.HasVisibleEditModal()).To(BeTrue())

			// Press Esc to cancel.
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			Expect(intent.HasVisibleEditModal()).To(BeFalse())
		})

		It("should cancel confirm modal and return to detail", func() {
			// Navigate to detail and confirm.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())

			// Press 'n' to cancel.
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			// Should return to detail modal.
			Expect(intent.HasVisibleConfirmModal()).To(BeFalse())
			Expect(intent.GetDetailModal()).NotTo(BeNil())
		})
	})

	Describe("Fact Extraction", func() {
		var (
			intent      *burst_management.Intent
			ctx         *burst_management.IntentContext
			mockService *mocks.BurstServiceMock
			repo        *careermemory.BurstRepository
			burst       *career.Burst
		)

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test",
				EventIDs:    []string{"e1", "e2"},
			}

			repo = careermemory.NewBurstRepository()
			_ = repo.Create(context.Background(), burst)

			mockService = mocks.NewBurstServiceMock()

			ctx = &burst_management.IntentContext{
				Bursts:          []*career.Burst{burst},
				Service:         mockService,
				BurstRepository: repo,
			}
			ctx.Validate()

			var err error
			intent, err = burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should execute extractFactsForBurst and save facts via suggestion acceptance", func() {
			// Set up extracted facts.
			extractedFacts := []career.Fact{
				{ID: "f1", Text: "Fact 1"},
				{ID: "f2", Text: "Fact 2"},
			}
			mockService.SetExtractedFacts(extractedFacts)

			// Create a suggestion and accept it - this triggers fact extraction.
			suggestion := burstfact.BurstSuggestion{
				Name:     "Test Suggestion",
				EventIDs: []string{"e1", "e2"},
			}
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{suggestion},
			})

			// Press 'a' to accept - this saves burst and triggers fact extraction.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command.
			msg := cmd()
			Expect(msg).NotTo(BeNil())

			// Should be FactExtractionCompleteMsg.
			extractionMsg, ok := msg.(burst_management.FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(extractionMsg.Facts).To(HaveLen(2))
			Expect(extractionMsg.Error).To(BeNil())
		})

		It("should handle extractFactsForBurst with extraction error via suggestion acceptance", func() {
			mockService.SetExtractError(errors.New("extraction failed"))

			// Create a suggestion and accept it.
			suggestion := burstfact.BurstSuggestion{
				Name:     "Test Suggestion",
				EventIDs: []string{"e1", "e2"},
			}
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{suggestion},
			})

			// Press 'a' to accept.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command.
			msg := cmd()

			// Should be error.
			extractionMsg, ok := msg.(burst_management.FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(extractionMsg.Error).NotTo(BeNil())
		})

		It("should handle extractFactsForBurst with save error via suggestion acceptance (skips failed saves)", func() {
			extractedFacts := []career.Fact{
				{ID: "f1", Text: "Fact 1"},
			}
			mockService.SetExtractedFacts(extractedFacts)
			mockService.SetSaveFactError(errors.New("save failed"))

			// Create a suggestion and accept it.
			suggestion := burstfact.BurstSuggestion{
				Name:     "Test Suggestion",
				EventIDs: []string{"e1", "e2"},
			}
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{suggestion},
			})

			// Press 'a' to accept.
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the command.
			msg := cmd()

			// Save errors are skipped, so Facts will be empty but no error.
			extractionMsg, ok := msg.(burst_management.FactExtractionCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(extractionMsg.Error).To(BeNil())
			Expect(extractionMsg.Facts).To(HaveLen(0)) // No facts saved due to error.
		})

		It("should show re-extract confirmation when facts already exist", func() {
			// Set up existing facts for the burst.
			existingFacts := []*career.Fact{
				{ID: "f1", Text: "Existing fact"},
			}
			mockService.SetFactsForBurst(burst.ID, existingFacts)

			// Navigate to detail and confirm.
			navResult := &screens.NavigateResult{ResultData: burst}
			intent.HandleNavigate(navResult)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			// Should show confirm modal (with re-extract message).
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())
		})
	})
})
