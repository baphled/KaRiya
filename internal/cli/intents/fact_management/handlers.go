// Package fact_management implements the FactManagement intent for managing career facts.
package fact_management

import (
	"github.com/baphled/kariya/internal/cli/intents"
	factmodals "github.com/baphled/kariya/internal/cli/screens/fact/modals"
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
				i.openDetailModal(i.context.SelectedFact)
			}
			return nil

		case "n":
			i.context.StartNewFact()
			i.openEditModal(i.context.EditingFact, true)
			return nil

		case "e":
			i.syncTableSelection()
			if i.context.SelectedFact != nil {
				i.context.StartEditFact(i.context.SelectedFact)
				i.openEditModal(i.context.EditingFact, false)
			}
			return nil

		case "d":
			i.syncTableSelection()
			if i.context.SelectedFact != nil {
				i.openDeleteConfirm(i.context.SelectedFact)
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
			i.state = StateList
			i.context.SelectedFact = nil
			return nil
		}

		switch msg.String() {
		case "e":
			if i.context.SelectedFact != nil {
				i.context.StartEditFact(i.context.SelectedFact)
				i.openEditModal(i.context.EditingFact, false)
			}
			return nil

		case "d":
			if i.context.SelectedFact != nil {
				i.openDeleteConfirm(i.context.SelectedFact)
			}
			return nil
		}
	}
	return nil
}

// handleEditorState handles messages in the editor state.
func (i *Intent) handleEditorState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch intents.HandleGlobalKeys(msg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		case intents.KeyBack:
			i.editModal = nil
			i.context.CancelEdit()
			i.state = StateList
			return nil
		}
	}
	return nil
}

// handleDeleteConfirmState handles messages in the delete confirm state.
func (i *Intent) handleDeleteConfirmState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			return i.performDelete()
		case "n", "esc":
			i.context.FactToDelete = nil
			i.state = StateList
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

// handleModalUpdates handles updates for all modals.
func (i *Intent) handleModalUpdates(msg tea.Msg) tea.Cmd {
	// Error modal has highest priority.
	if i.errorModal != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
			i.errorModal = nil
			return noopCmd
		}
		return noopCmd
	}

	// Delete confirmation modal.
	if i.deleteModal != nil && i.deleteModal.IsVisible() {
		cmd, confirmed := i.deleteModal.Update(msg)
		if !i.deleteModal.IsVisible() {
			if confirmed && i.context.FactToDelete != nil {
				deleteCmd := i.performDelete()
				i.deleteModal = nil
				return deleteCmd
			}
			i.context.FactToDelete = nil
			i.deleteModal = nil
			return noopCmd
		}
		return cmd
	}

	// Edit modal.
	if i.editModal != nil && i.editModal.IsVisible() {
		cmd, result := i.editModal.Update(msg)
		if !i.editModal.IsVisible() {
			if result != nil && result.Accepted {
				i.applyEditResult(result)
			}
			i.editModal = nil
			i.context.CancelEdit()
			i.state = StateList
			return noopCmd
		}
		return cmd
	}

	// Detail modal.
	if i.detailModal != nil && i.detailModal.IsVisible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "e":
				if i.context.SelectedFact != nil {
					i.detailModal.Hide()
					i.detailModal = nil
					i.context.StartEditFact(i.context.SelectedFact)
					i.openEditModal(i.context.EditingFact, false)
					return nil
				}
			case "d":
				if i.context.SelectedFact != nil {
					i.detailModal.Hide()
					i.detailModal = nil
					i.openDeleteConfirm(i.context.SelectedFact)
					return nil
				}
			}
		}
		_, cmd := i.detailModal.Update(msg)
		if !i.detailModal.IsVisible() {
			i.detailModal = nil
		}
		return cmd
	}

	return nil
}

// noopCmd is a sentinel command to indicate a message was consumed.
func noopCmd() tea.Msg { return nil }

// applyEditResult applies the result of an edit operation.
func (i *Intent) applyEditResult(result *factmodals.EditResult) {
	if result == nil || !result.Accepted {
		return
	}

	// Apply changes from the modal to the editing fact.
	i.context.EditingFact.Text = result.Fact.Text
	i.context.EditingFact.CompetencyCategories = result.Fact.CompetencyCategories
	i.context.EditingFact.RoleFit = result.Fact.RoleFit
	i.context.EditingFact.AudienceRelevance = result.Fact.AudienceRelevance
	i.context.EditingFact.StrengthSignal = result.Fact.StrengthSignal

	// Save the fact.
	if err := i.context.SaveEdit(); err != nil {
		i.context.SetFormError("general", "Save failed: "+err.Error())
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
				Message: "Fact " + action + " successfully",
			},
		}
		// Refresh table after save.
		i.tableBehavior.SetItems(i.context.Facts)
		i.syncTableSelection()
	}
}
