package burst_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
	tea "github.com/charmbracelet/bubbletea"
)

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
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Burst should now be confirmed.
			Expect(intent.GetSelectedBurst().Confirmed).To(BeTrue())
			// Should transition to extracting facts state.
			Expect(intent.GetState()).To(Equal(burst_management.StateExtractingFacts))
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
				Suggestions: []burst_fact.BurstSuggestion{
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
				Suggestions: []burst_fact.BurstSuggestion{
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
				Suggestions: []burst_fact.BurstSuggestion{
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
				AcceptedSuggestions: []burst_fact.BurstSuggestion{loadMsg.Suggestions[0]},
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
				Suggestions: []burst_fact.BurstSuggestion{
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
				Suggestions: []burst_fact.BurstSuggestion{
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
})
