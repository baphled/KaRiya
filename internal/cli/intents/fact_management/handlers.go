// Package fact_management implements the FactManagement intent for managing career facts.
package fact_management

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/intents"
	tea "github.com/charmbracelet/bubbletea"
)

// handleListState handles messages in the list state.
func (i *Intent) handleListState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first.
		switch intents.HandleGlobalKeys(msg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		case intents.KeyBack:
			// At root state, back means cancel and return to main menu.
			i.setCancelled()
			return nil
		}

		// Try TableBehavior navigation.
		if i.tableBehavior.HandleNavigation(msg.String()) {
			i.syncTableSelection()
			return nil
		}

		switch msg.String() {
		case "enter", " ":
			i.syncTableSelection()
			if i.context.SelectedFact != nil {
				i.state = StateView
			}
			return nil

		case "n":
			i.context.StartNewFact()
			i.editModal = intents.NewEditFactModal(i.context.EditingFact)
			i.state = StateEditor
			// Return form init command to properly initialize the huh form.
			return i.editModal.Init()

		case "e":
			i.syncTableSelection()
			if i.context.SelectedFact != nil {
				i.context.StartEditFact(i.context.SelectedFact)
				i.editModal = intents.NewEditFactModal(i.context.EditingFact)
				i.state = StateEditor
				// Return form init command to properly initialize the huh form.
				return i.editModal.Init()
			}
			return nil

		case "d":
			i.syncTableSelection()
			if i.context.SelectedFact != nil {
				i.context.FactToDelete = i.context.SelectedFact
				i.state = StateDeleteConfirm
			}
			return nil

		case "r":
			if err := i.context.LoadFacts(); err != nil {
				i.result = &intents.IntentResult[*Result]{
					Status: intents.Failed,
					Error: &intents.IntentError{
						Code:    "RELOAD_FAILED",
						Message: "Failed to reload facts",
						Cause:   err,
					},
				}
			} else {
				i.tableBehavior.SetItems(i.context.Facts)
				i.syncTableSelection()
			}
			return nil
		}
	}
	return nil
}

// handleViewState handles messages in the view state.
func (i *Intent) handleViewState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch intents.HandleGlobalKeys(msg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		case intents.KeyBack:
			// Go back to list state.
			i.state = StateList
			i.context.SelectedFact = nil
			return nil
		}

		switch msg.String() {
		case "e":
			if i.context.SelectedFact != nil {
				i.context.StartEditFact(i.context.SelectedFact)
				i.editModal = intents.NewEditFactModal(i.context.EditingFact)
				i.state = StateEditor
				// Return form init command to properly initialize the huh form.
				return i.editModal.Init()
			}

		case "d":
			if i.context.SelectedFact != nil {
				i.context.FactToDelete = i.context.SelectedFact
				i.state = StateDeleteConfirm
			}
		}
	}
	return nil
}

// handleEditorState handles messages in the editor state.
func (i *Intent) handleEditorState(msg tea.Msg) tea.Cmd {
	// Handle global keys FIRST (before delegating to modal)
	// This ensures esc, q, ? keys work even when modal has focus.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch intents.HandleGlobalKeys(msg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		case intents.KeyBack:
			// Close modal if active.
			if i.editModal != nil {
				i.editModal = nil
			}
			i.context.CancelEdit()
			if i.context.IsNewFact {
				i.state = StateList
			} else {
				i.state = StateView
			}
			return nil
		}
	}

	// If modal is not initialized, handle legacy behavior (fallback).
	if i.editModal == nil {
		return nil
	}

	// Delegate to the modal for form handling.
	cmd := i.editModal.Update(msg)

	// Check if modal completed (form submitted or cancelled).
	if result := i.editModal.Result(); result != nil {
		if result.Accepted {
			// Apply changes from the modal to the editing fact.
			i.context.EditingFact.Text = result.Modified.Text
			i.context.EditingFact.CompetencyCategories = result.Modified.CompetencyCategories
			i.context.EditingFact.RoleFit = result.Modified.RoleFit
			i.context.EditingFact.AudienceRelevance = result.Modified.AudienceRelevance
			i.context.EditingFact.StrengthSignal = result.Modified.StrengthSignal

			// Save the fact.
			if err := i.context.SaveEdit(); err != nil {
				i.context.SetFormError("general", fmt.Sprintf("Save failed: %v", err))
			} else {
				action := "updated"
				if i.context.IsNewFact {
					action = "created"
				}
				i.result = &intents.IntentResult[*Result]{
					Status: intents.Completed,
					Data: &Result{
						Action:  action,
						Fact:    i.context.EditingFact,
						Facts:   i.context.Facts,
						Message: fmt.Sprintf("Fact %s successfully", action),
					},
				}
				// Refresh table after save.
				i.tableBehavior.SetItems(i.context.Facts)
				i.syncTableSelection()
			}
		}

		// Clear modal and return to appropriate state.
		i.editModal = nil
		i.context.CancelEdit()
		if i.context.IsNewFact {
			i.state = StateList
		} else {
			i.state = StateView
		}
		return nil
	}

	return cmd
}

// handleDeleteConfirmState handles messages in the delete confirm state.
func (i *Intent) handleDeleteConfirmState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			if i.context.FactToDelete != nil {
				if err := i.context.DeleteFact(i.context.FactToDelete.ID); err != nil {
					i.result = &intents.IntentResult[*Result]{
						Status: intents.Failed,
						Error: &intents.IntentError{
							Code:    "DELETE_FAILED",
							Message: "Failed to delete fact",
							Cause:   err,
						},
					}
				} else {
					i.result = &intents.IntentResult[*Result]{
						Status: intents.Completed,
						Data: &Result{
							Action:  "deleted",
							Fact:    i.context.FactToDelete,
							Facts:   i.context.Facts,
							Message: "Fact deleted successfully",
						},
					}
					// Refresh table after delete.
					i.tableBehavior.SetItems(i.context.Facts)
					i.syncTableSelection()
				}
				i.context.FactToDelete = nil
				i.state = StateList
			}

		case "n", "esc":
			i.context.FactToDelete = nil
			i.state = StateView
		}
	}
	return nil
}

// handleResultsState handles messages in the results state.
func (i *Intent) handleResultsState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch intents.HandleGlobalKeys(msg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		case intents.KeyBack:
			i.state = StateList
			return nil
		}
	}
	return nil
}
