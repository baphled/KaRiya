// Package factmanagement implements the FactManagement intent for managing career facts.
package factmanagement

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/intents"
	factmodals "github.com/baphled/kariya/internal/cli/screens/facts/modals"
	tea "github.com/charmbracelet/bubbletea"
)

// handleListState handles messages in the list state.
func (i *Intent) handleListState(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	// Handle global keys first.
	if cmd := i.handleListGlobalKeys(keyMsg); cmd != nil {
		return cmd
	}

	// Try TableBehavior navigation.
	if i.tableBehavior.HandleNavigation(keyMsg.String()) {
		i.syncTableSelection()
		return nil
	}

	return i.handleListKeyActions(keyMsg)
}

// handleListGlobalKeys handles global key bindings in list state.
func (i *Intent) handleListGlobalKeys(msg tea.KeyMsg) tea.Cmd {
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
	return nil
}

// handleListKeyActions handles specific key actions in list state.
func (i *Intent) handleListKeyActions(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter", " ":
		i.syncTableSelection()
		if i.context.SelectedFact != nil {
			i.state = StateView
		}
		return nil

	case "n":
		i.context.StartNewFact()
		i.editModal = factmodals.NewEditFactModal(i.context.EditingFact)
		i.state = StateEditor
		return i.editModal.Init()

	case "e":
		return i.handleEditKeyInList()

	case "d":
		i.syncTableSelection()
		if i.context.SelectedFact != nil {
			i.context.FactToDelete = i.context.SelectedFact
			i.state = StateDeleteConfirm
		}
		return nil

	case "r":
		return i.handleRefreshKey()
	}
	return nil
}

// handleEditKeyInList handles the edit key press in list state.
func (i *Intent) handleEditKeyInList() tea.Cmd {
	i.syncTableSelection()
	if i.context.SelectedFact != nil {
		i.context.StartEditFact(i.context.SelectedFact)
		i.editModal = factmodals.NewEditFactModal(i.context.EditingFact)
		i.state = StateEditor
		return i.editModal.Init()
	}
	return nil
}

// handleRefreshKey handles the refresh key press.
func (i *Intent) handleRefreshKey() tea.Cmd {
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

// handleViewState handles messages in the view state.
func (i *Intent) handleViewState(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch intents.HandleGlobalKeys(keyMsg) {
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

	switch keyMsg.String() {
	case "e":
		if i.context.SelectedFact != nil {
			i.context.StartEditFact(i.context.SelectedFact)
			i.editModal = factmodals.NewEditFactModal(i.context.EditingFact)
			i.state = StateEditor
			return i.editModal.Init()
		}

	case "d":
		if i.context.SelectedFact != nil {
			i.context.FactToDelete = i.context.SelectedFact
			i.state = StateDeleteConfirm
		}
	}
	return nil
}

// handleEditorState handles messages in the editor state.
func (i *Intent) handleEditorState(msg tea.Msg) tea.Cmd {
	// Handle global keys FIRST (before delegating to modal).
	// This ensures esc, q, ? keys work even when modal has focus.
	if cmd, handled := i.handleEditorGlobalKeys(msg); handled {
		return cmd
	}

	// If modal is not initialized, handle legacy behavior (fallback).
	if i.editModal == nil {
		return nil
	}

	// Delegate to the modal for form handling.
	cmd := i.editModal.Update(msg)

	// Check if modal completed (form submitted or cancelled).
	if result := i.editModal.Result(); result != nil {
		i.handleEditorModalResult(result)
		return nil
	}

	return cmd
}

// handleEditorGlobalKeys handles global key bindings in editor state.
// Returns (cmd, handled) where handled indicates if the key was processed.
func (i *Intent) handleEditorGlobalKeys(msg tea.Msg) (tea.Cmd, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil, false
	}

	switch intents.HandleGlobalKeys(keyMsg) {
	case intents.KeyQuit:
		return tea.Quit, true
	case intents.KeyHelp:
		i.ToggleHelp()
		return nil, true
	case intents.KeyBack:
		i.cancelEditorAndReturn()
		return nil, true
	}
	return nil, false
}

// cancelEditorAndReturn closes the editor and returns to the appropriate state.
func (i *Intent) cancelEditorAndReturn() {
	// Save IsNewFact before CancelEdit clears it.
	wasNewFact := i.context.IsNewFact
	// Close modal if active.
	if i.editModal != nil {
		i.editModal = nil
	}
	i.context.CancelEdit()
	if wasNewFact {
		i.state = StateList
	} else {
		i.state = StateView
	}
}

// handleEditorModalResult handles the result from the editor modal.
func (i *Intent) handleEditorModalResult(result *factmodals.EditResult) {
	if result.Accepted {
		i.applyEditorChanges(result)
	}

	// Clear modal and return to appropriate state.
	// Save IsNewFact before CancelEdit clears it.
	wasNewFact := i.context.IsNewFact
	i.editModal = nil
	i.context.CancelEdit()
	if wasNewFact {
		i.state = StateList
	} else {
		i.state = StateView
	}
}

// applyEditorChanges applies changes from the modal and saves the fact.
func (i *Intent) applyEditorChanges(result *factmodals.EditResult) {
	// Apply changes from the modal to the editing fact.
	i.context.EditingFact.Text = result.Modified.Text
	i.context.EditingFact.CompetencyCategories = result.Modified.CompetencyCategories
	i.context.EditingFact.RoleFit = result.Modified.RoleFit
	i.context.EditingFact.AudienceRelevance = result.Modified.AudienceRelevance
	i.context.EditingFact.StrengthSignal = result.Modified.StrengthSignal

	// Save the fact.
	if err := i.context.SaveEdit(); err != nil {
		i.context.SetFormError("general", fmt.Sprintf("Save failed: %v", err))
		return
	}

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

// handleDeleteConfirmState handles messages in the delete confirm state.
func (i *Intent) handleDeleteConfirmState(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch keyMsg.String() {
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
	return nil
}

// handleResultsState handles messages in the results state.
func (i *Intent) handleResultsState(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch intents.HandleGlobalKeys(keyMsg) {
	case intents.KeyQuit:
		return tea.Quit
	case intents.KeyHelp:
		i.ToggleHelp()
		return nil
	case intents.KeyBack:
		i.state = StateList
		return nil
	}
	return nil
}
