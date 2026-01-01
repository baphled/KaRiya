package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ListDeletionState manages the deletion confirmation and execution workflow
type ListDeletionState struct {
	ConfirmationDialog *ConfirmationDialog
	DeletingItemID    string
	SuccessMsg        string
	ErrorMsg          string
	ShowMessage       bool
}

// NewListDeletionState creates a new deletion state tracker
func NewListDeletionState() *ListDeletionState {
	return &ListDeletionState{
		ConfirmationDialog: nil,
		DeletingItemID:     "",
		SuccessMsg:         "",
		ErrorMsg:           "",
		ShowMessage:        false,
	}
}

// ShowConfirmation creates and shows a confirmation dialog for deletion
func (lds *ListDeletionState) ShowConfirmation(title, itemName string) {
	message := fmt.Sprintf("Delete %s: \"%s\"?\n\nThis action cannot be undone.", title, truncateText(itemName, 60))
	lds.ConfirmationDialog = NewConfirmationDialog("Delete "+title, message)
}

// IsConfirming returns true if a confirmation dialog is active
func (lds *ListDeletionState) IsConfirming() bool {
	return lds.ConfirmationDialog != nil
}

// IsConfirmed returns true if the user confirmed deletion
func (lds *ListDeletionState) IsConfirmed() bool {
	return lds.ConfirmationDialog != nil && lds.ConfirmationDialog.IsConfirmed()
}

// IsCancelled returns true if the user cancelled deletion
func (lds *ListDeletionState) IsCancelled() bool {
	return lds.ConfirmationDialog != nil && lds.ConfirmationDialog.IsCancelled()
}

// SetSuccessMsg sets the success message after deletion
func (lds *ListDeletionState) SetSuccessMsg(msg string) {
	lds.SuccessMsg = msg
	lds.ErrorMsg = ""
	lds.ShowMessage = true
}

// SetErrorMsg sets the error message if deletion failed
func (lds *ListDeletionState) SetErrorMsg(msg string) {
	lds.ErrorMsg = msg
	lds.SuccessMsg = ""
	lds.ShowMessage = true
}

// Clear resets the deletion state
func (lds *ListDeletionState) Clear() {
	lds.ConfirmationDialog = nil
	lds.DeletingItemID = ""
	lds.SuccessMsg = ""
	lds.ErrorMsg = ""
	lds.ShowMessage = false
}

// UpdateConfirmation handles messages for the confirmation dialog
func (lds *ListDeletionState) UpdateConfirmation(msg tea.Msg) tea.Cmd {
	if lds.ConfirmationDialog == nil {
		return nil
	}
	updatedDialog, cmd := lds.ConfirmationDialog.Update(msg)
	lds.ConfirmationDialog = updatedDialog
	return cmd
}

// ListNavigationKeyHandler provides common keyboard navigation for list models
// It handles up/down/pgup/pgdn/home/end keys and delegates item actions to callbacks
type ListNavigationKeyHandler struct {
	itemCallbacks ListItemCallbacks
}

// ListItemCallbacks provides callbacks for item-specific actions
type ListItemCallbacks interface {
	// OnView is called when user wants to view an item (enter/v)
	OnView(item interface{}) tea.Cmd

	// OnEdit is called when user wants to edit an item (e)
	OnEdit(item interface{}) tea.Cmd

	// OnDelete is called when user wants to delete an item (x/d)
	OnDelete(item interface{})

	// OnToggleSelection is called when user toggles item selection (space)
	OnToggleSelection(item interface{})

	// HasSelectedItem returns the currently selected item or nil
	HasSelectedItem() interface{}

	// MoveUp moves the selection up by count items
	MoveUp(count int)

	// MoveDown moves the selection down by count items
	MoveDown(count int)

	// MoveToFirst moves selection to first item
	MoveToFirst()

	// MoveToLast moves selection to last item
	MoveToLast()

	// UpdateDisplay refreshes the table display
	UpdateDisplay()

	// GetRowCount returns the number of visible rows
	GetRowCount() int

	// GetCurrentIndex returns the current selection index
	GetCurrentIndex() int

	// GetPageSize returns the page size for pagination
	GetPageSize() int
}

// NewListNavigationKeyHandler creates a new key handler
func NewListNavigationKeyHandler(callbacks ListItemCallbacks) *ListNavigationKeyHandler {
	return &ListNavigationKeyHandler{
		itemCallbacks: callbacks,
	}
}

// HandleNavigationKey processes navigation keys and returns a command if needed
func (lnkh *ListNavigationKeyHandler) HandleNavigationKey(keyStr string) tea.Cmd {
	switch keyStr {
	case "up", "k":
		lnkh.itemCallbacks.MoveUp(1)
		lnkh.itemCallbacks.UpdateDisplay()
	case "down", "j":
		lnkh.itemCallbacks.MoveDown(1)
		lnkh.itemCallbacks.UpdateDisplay()
	case "pgup", "ctrl+b":
		pageSize := lnkh.itemCallbacks.GetPageSize()
		if lnkh.itemCallbacks.GetCurrentIndex() >= pageSize {
			lnkh.itemCallbacks.MoveUp(pageSize)
		} else {
			lnkh.itemCallbacks.MoveToFirst()
		}
		lnkh.itemCallbacks.UpdateDisplay()
	case "pgdn", "ctrl+f":
		pageSize := lnkh.itemCallbacks.GetPageSize()
		if lnkh.itemCallbacks.GetCurrentIndex()+pageSize < lnkh.itemCallbacks.GetRowCount() {
			lnkh.itemCallbacks.MoveDown(pageSize)
		} else {
			lnkh.itemCallbacks.MoveToLast()
		}
		lnkh.itemCallbacks.UpdateDisplay()
	case "home", "g":
		lnkh.itemCallbacks.MoveToFirst()
		lnkh.itemCallbacks.UpdateDisplay()
	case "end", "G":
		lnkh.itemCallbacks.MoveToLast()
		lnkh.itemCallbacks.UpdateDisplay()
	case "enter", "v":
		item := lnkh.itemCallbacks.HasSelectedItem()
		if item != nil {
			return lnkh.itemCallbacks.OnView(item)
		}
	case " ", "space":
		item := lnkh.itemCallbacks.HasSelectedItem()
		if item != nil {
			lnkh.itemCallbacks.OnToggleSelection(item)
		}
	case "x", "d":
		item := lnkh.itemCallbacks.HasSelectedItem()
		if item != nil {
			lnkh.itemCallbacks.OnDelete(item)
		}
	case "e":
		item := lnkh.itemCallbacks.HasSelectedItem()
		if item != nil {
			return lnkh.itemCallbacks.OnEdit(item)
		}
	}
	return nil
}

// truncateText truncates text to maxLen with ellipsis if needed
func truncateText(text string, maxLen int) string {
	if len(text) > maxLen {
		return text[:maxLen-3] + "..."
	}
	return text
}


// truncateText truncates text to maxLen with ellipsis if needed
func truncateText(text string, maxLen int) string {
	if len(text) > maxLen {
		return text[:maxLen-3] + "..."
	}
	return text
}
