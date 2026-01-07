package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ListDeletionState manages the deletion confirmation and execution workflow
type ListDeletionState struct {
	ConfirmationDialog *ConfirmationDialog
	DeletingItemID     string
	SuccessMsg         string
	ErrorMsg           string
	ShowMessage        bool
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

// PaginationHelper manages pagination state and calculations for list models
// It provides methods for calculating pages, navigating between pages,
// and getting page bounds for data slicing.
type PaginationHelper struct {
	currentPage int // Current page number (1-based)
	pageSize    int // Number of items per page
	totalCount  int // Total number of items
}

// NewPaginationHelper creates a new pagination helper
func NewPaginationHelper(pageSize int) *PaginationHelper {
	return &PaginationHelper{
		currentPage: 1,
		pageSize:    pageSize,
		totalCount:  0,
	}
}

// SetTotalCount sets the total number of items
func (ph *PaginationHelper) SetTotalCount(count int) *PaginationHelper {
	ph.totalCount = count
	// Ensure current page is still valid
	if ph.currentPage > ph.GetTotalPages() && ph.totalCount > 0 {
		ph.currentPage = ph.GetTotalPages()
	}
	if ph.totalCount == 0 {
		ph.currentPage = 1
	}
	return ph
}

// GetCurrentPage returns the current page number
func (ph *PaginationHelper) GetCurrentPage() int {
	return ph.currentPage
}

// GetPageSize returns the page size
func (ph *PaginationHelper) GetPageSize() int {
	return ph.pageSize
}

// GetTotalCount returns the total number of items
func (ph *PaginationHelper) GetTotalCount() int {
	return ph.totalCount
}

// GetTotalPages calculates the total number of pages
func (ph *PaginationHelper) GetTotalPages() int {
	if ph.totalCount == 0 {
		return 1
	}
	return (ph.totalCount + ph.pageSize - 1) / ph.pageSize
}

// GetPageStartIndex returns the starting index for the current page (0-based)
func (ph *PaginationHelper) GetPageStartIndex() int {
	return (ph.currentPage - 1) * ph.pageSize
}

// GetPageEndIndex returns the ending index for the current page (exclusive, 0-based)
func (ph *PaginationHelper) GetPageEndIndex() int {
	return ph.GetPageStartIndex() + ph.pageSize
}

// GetPaginationInfo returns formatted pagination info for display
// Example: "Showing 1-10 of 25 items"
func (ph *PaginationHelper) GetPaginationInfo(itemName string) string {
	if ph.totalCount == 0 {
		return fmt.Sprintf("Showing 0 of 0 %s", itemName)
	}
	startIdx := ph.GetPageStartIndex() + 1
	endIdx := ph.GetPageEndIndex()
	if endIdx > ph.totalCount {
		endIdx = ph.totalCount
	}
	return fmt.Sprintf("Showing %d-%d of %d %s", startIdx, endIdx, ph.totalCount, itemName)
}

// NextPage moves to the next page if available
func (ph *PaginationHelper) NextPage() *PaginationHelper {
	if ph.currentPage < ph.GetTotalPages() {
		ph.currentPage++
	}
	return ph
}

// PrevPage moves to the previous page if available
func (ph *PaginationHelper) PrevPage() *PaginationHelper {
	if ph.currentPage > 1 {
		ph.currentPage--
	}
	return ph
}

// GoToFirstPage moves to the first page
func (ph *PaginationHelper) GoToFirstPage() *PaginationHelper {
	ph.currentPage = 1
	return ph
}

// GoToLastPage moves to the last page
func (ph *PaginationHelper) GoToLastPage() *PaginationHelper {
	ph.currentPage = ph.GetTotalPages()
	return ph
}

// IsFirstPage returns true if on the first page
func (ph *PaginationHelper) IsFirstPage() bool {
	return ph.currentPage == 1
}

// IsLastPage returns true if on the last page
func (ph *PaginationHelper) IsLastPage() bool {
	return ph.currentPage >= ph.GetTotalPages()
}

// CanGoNext returns true if there's a next page
func (ph *PaginationHelper) CanGoNext() bool {
	return ph.currentPage < ph.GetTotalPages()
}

// CanGoPrev returns true if there's a previous page
func (ph *PaginationHelper) CanGoPrev() bool {
	return ph.currentPage > 1
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

	// NextPage moves to the next page
	NextPage()

	// PrevPage moves to the previous page
	PrevPage()

	// GoToFirstPage moves to the first page
	GoToFirstPage()

	// GoToLastPage moves to the last page
	GoToLastPage()

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
// Returns a command for all handled keys to prevent the table from processing them
func (lnkh *ListNavigationKeyHandler) HandleNavigationKey(keyStr string) tea.Cmd {
	switch keyStr {
	case "up", "k":
		lnkh.itemCallbacks.MoveUp(1)
		lnkh.itemCallbacks.UpdateDisplay()
		// Return a command to prevent table from processing this key
		return func() tea.Msg { return nil }
	case "down", "j":
		lnkh.itemCallbacks.MoveDown(1)
		lnkh.itemCallbacks.UpdateDisplay()
		// Return a command to prevent table from processing this key
		return func() tea.Msg { return nil }
	case "pgup", "ctrl+u":
		lnkh.itemCallbacks.PrevPage()
		lnkh.itemCallbacks.UpdateDisplay()
		// Return a command to prevent table from processing this key
		return func() tea.Msg { return nil }
	case "pgdn", "ctrl+d":
		lnkh.itemCallbacks.NextPage()
		lnkh.itemCallbacks.UpdateDisplay()
		// Return a command to prevent table from processing this key
		return func() tea.Msg { return nil }
	case "home", "g":
		lnkh.itemCallbacks.GoToFirstPage()
		lnkh.itemCallbacks.UpdateDisplay()
		// Return a command to prevent table from processing this key
		return func() tea.Msg { return nil }
	case "end", "G":
		lnkh.itemCallbacks.GoToLastPage()
		lnkh.itemCallbacks.UpdateDisplay()
		// Return a command to prevent table from processing this key
		return func() tea.Msg { return nil }
	case "enter", "v":
		item := lnkh.itemCallbacks.HasSelectedItem()
		if item != nil {
			return lnkh.itemCallbacks.OnView(item)
		}
	case " ", "space":
		item := lnkh.itemCallbacks.HasSelectedItem()
		if item != nil {
			lnkh.itemCallbacks.OnToggleSelection(item)
			// Return a command to prevent table from processing this key
			return func() tea.Msg { return nil }
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
