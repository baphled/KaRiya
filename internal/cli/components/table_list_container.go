package components

import (
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// TableListContainer provides a reusable container for displaying table-based lists
// with consistent header, footer, pagination, and help footer components.
// It handles the layout and rendering of table lists across the application.
// It also manages cursor/selection state and provides navigation methods.
type TableListContainer struct {
	table          table.Model
	footer         FooterModel
	paginationInfo string
	showPagination bool
	emptyMessage   string
	showEmptyState bool
	errorMessage   string
	showError      bool
	width          int
	height         int
	selectedIdx    int // Cursor position for selection
}

// NewTableListContainer creates a new TableListContainer with default styling
func NewTableListContainer(tableModel table.Model, headerTitle string, width int) *TableListContainer {
	return &TableListContainer{
		table:          tableModel,
		footer:         NewFooter(width),
		width:          width,
		height:         24,
		emptyMessage:   "No items to display",
		showEmptyState: false,
		errorMessage:   "",
		showError:      false,
		showPagination: false,
		selectedIdx:    0,
	}
}

// SetTable updates the table model being displayed and syncs cursor state
// Important: This pushes the container's selectedIdx to the table, not the other way around
func (tlc *TableListContainer) SetTable(tableModel table.Model) *TableListContainer {
	tlc.table = tableModel
	// Set the table's cursor to match the container's selected index
	// This ensures indicators in row text match the table's visual selection
	tlc.validateAndSyncIdx()
	return tlc
}

// SetDimensions sets the width and height of the container
func (tlc *TableListContainer) SetDimensions(width, height int) *TableListContainer {
	tlc.width = width
	tlc.height = height
	tlc.footer.SetWidth(width)
	tlc.table.SetWidth(width)
	tlc.table.SetHeight(height - 10)
	return tlc
}

// SetPaginationInfo sets the pagination information to display
func (tlc *TableListContainer) SetPaginationInfo(info string) *TableListContainer {
	tlc.paginationInfo = info
	tlc.showPagination = true
	return tlc
}

// HidePagination hides the pagination information
func (tlc *TableListContainer) HidePagination() *TableListContainer {
	tlc.showPagination = false
	return tlc
}

// SetEmptyStateMessage sets the message to display when the list is empty
func (tlc *TableListContainer) SetEmptyStateMessage(message string) *TableListContainer {
	tlc.emptyMessage = message
	tlc.showEmptyState = true
	return tlc
}

// HideEmptyState hides the empty state message
func (tlc *TableListContainer) HideEmptyState() *TableListContainer {
	tlc.showEmptyState = false
	return tlc
}

// SetErrorMessage sets an error message to display
func (tlc *TableListContainer) SetErrorMessage(message string) *TableListContainer {
	tlc.errorMessage = message
	tlc.showError = true
	return tlc
}

// ClearError clears the error message
func (tlc *TableListContainer) ClearError() *TableListContainer {
	tlc.errorMessage = ""
	tlc.showError = false
	return tlc
}

// GetWidth returns the current width
func (tlc *TableListContainer) GetWidth() int {
	return tlc.width
}

// GetHeight returns the current height
func (tlc *TableListContainer) GetHeight() int {
	return tlc.height
}

// validateAndSyncIdx ensures the selected index is valid and syncs to the table
// This is called after SetTable to push the container's index to the table
func (tlc *TableListContainer) validateAndSyncIdx() {
	rowCount := len(tlc.table.Rows())

	// If table is empty, set index to 0
	if rowCount == 0 {
		tlc.selectedIdx = 0
		tlc.table.SetCursor(0)
		return
	}

	// Ensure selectedIdx is within valid bounds
	if tlc.selectedIdx < 0 {
		tlc.selectedIdx = 0
	} else if tlc.selectedIdx >= rowCount {
		tlc.selectedIdx = rowCount - 1
	}

	// Sync the table cursor to match the container's index
	tlc.table.SetCursor(tlc.selectedIdx)
}

// syncIdx ensures the selected index is valid and within bounds
// Deprecated: Use validateAndSyncIdx instead which pushes to table
func (tlc *TableListContainer) syncIdx() {
	cursor := tlc.table.Cursor()
	// If cursor is negative (table empty or not set), default to 0
	if cursor < 0 {
		tlc.selectedIdx = 0
	} else if cursor >= len(tlc.table.Rows()) {
		// If cursor is out of bounds, set to last valid index
		tlc.selectedIdx = len(tlc.table.Rows()) - 1
		if tlc.selectedIdx < 0 {
			tlc.selectedIdx = 0
		}
	} else {
		tlc.selectedIdx = cursor
	}
}

// GetSelectedIdx returns the currently selected index
// Always returns a value >= 0, even if the table is empty
func (tlc *TableListContainer) GetSelectedIdx() int {
	// Ensure selectedIdx is always valid
	if tlc.selectedIdx < 0 {
		tlc.selectedIdx = 0
	}
	if len(tlc.table.Rows()) > 0 && tlc.selectedIdx >= len(tlc.table.Rows()) {
		tlc.selectedIdx = len(tlc.table.Rows()) - 1
	}
	return tlc.selectedIdx
}

// SetSelectedIdx sets the selected index directly
func (tlc *TableListContainer) SetSelectedIdx(idx int) *TableListContainer {
	if idx < 0 {
		idx = 0
	}
	if len(tlc.table.Rows()) > 0 && idx >= len(tlc.table.Rows()) {
		idx = len(tlc.table.Rows()) - 1
	}
	tlc.selectedIdx = idx
	tlc.table.SetCursor(idx)
	return tlc
}

// MoveUp moves the cursor up by one position
func (tlc *TableListContainer) MoveUp(count int) *TableListContainer {
	newIdx := tlc.selectedIdx - count
	if newIdx < 0 {
		newIdx = 0
	}
	tlc.selectedIdx = newIdx
	tlc.table.SetCursor(newIdx)
	return tlc
}

// MoveDown moves the cursor down by one position
func (tlc *TableListContainer) MoveDown(count int) *TableListContainer {
	newIdx := tlc.selectedIdx + count
	maxIdx := len(tlc.table.Rows()) - 1
	if maxIdx < 0 {
		maxIdx = 0
	}
	if newIdx > maxIdx {
		newIdx = maxIdx
	}
	tlc.selectedIdx = newIdx
	tlc.table.SetCursor(newIdx)
	return tlc
}

// MoveToFirst moves the cursor to the first item
func (tlc *TableListContainer) MoveToFirst() *TableListContainer {
	tlc.selectedIdx = 0
	tlc.table.SetCursor(0)
	return tlc
}

// MoveToLast moves the cursor to the last item
func (tlc *TableListContainer) MoveToLast() *TableListContainer {
	maxIdx := len(tlc.table.Rows()) - 1
	if maxIdx < 0 {
		maxIdx = 0
	}
	tlc.selectedIdx = maxIdx
	tlc.table.SetCursor(maxIdx)
	return tlc
}

// SyncCursorFromTable syncs the container's selectedIdx from the table's cursor
func (tlc *TableListContainer) SyncCursorFromTable() *TableListContainer {
	tlc.syncIdx()
	return tlc
}

// Render returns the complete table list view with all components
func (tlc *TableListContainer) Render() string {
	footerView := tlc.footer.View()

	// Handle error state
	if tlc.showError && tlc.errorMessage != "" {
		errorContent := styles.ErrorBox.Render(tlc.errorMessage)
		parts := []string{errorContent}
		if footerView != "" {
			parts = append(parts, "", footerView)
		}
		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	}

	// Handle empty state
	if tlc.showEmptyState && len(tlc.table.Rows()) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Italic(true)
		emptyContent := emptyStyle.Render(tlc.emptyMessage)

		parts := []string{emptyContent}

		// Add pagination even for empty state
		if tlc.showPagination && tlc.paginationInfo != "" {
			paginationStyle := lipgloss.NewStyle().
				Foreground(styles.ColorTextSecondary).
				MarginTop(1)
			paginationView := paginationStyle.Render(tlc.paginationInfo)
			parts = append(parts, "", paginationView)
		}

		if footerView != "" {
			parts = append(parts, "", footerView)
		}

		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	}

	// Render table
	tableView := tlc.table.View()

	parts := []string{tableView}

	// Add pagination if enabled
	if tlc.showPagination && tlc.paginationInfo != "" {
		paginationStyle := lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			MarginTop(1)
		paginationView := paginationStyle.Render(tlc.paginationInfo)
		parts = append(parts, "", paginationView)
	}

	if footerView != "" {
		parts = append(parts, "", footerView)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
