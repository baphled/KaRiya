// Package widgets provides higher-level composite components that combine
// primitives and behaviors into reusable UI patterns.
package widgets

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// FooterRenderer is an interface for components that can render a footer.
// This allows TableList to work with any footer implementation without
// creating import cycles.
type FooterRenderer interface {
	View() string
}

// TableList provides a reusable container for displaying table-based lists
// with consistent header, footer, pagination, and help footer components.
// It handles the layout and rendering of table lists across the application.
// It also manages cursor/selection state and provides navigation methods.
//
// This is the UIKit version that uses theme-based styling exclusively.
type TableList struct {
	table          table.Model
	footer         FooterRenderer
	paginationInfo string
	showPagination bool
	emptyMessage   string
	showEmptyState bool
	errorMessage   string
	showError      bool
	width          int
	height         int
	selectedIdx    int // Cursor position for selection
	theme          themes.Theme
}

// NewTableList creates a new TableList with default styling
// Note: Footer must be set separately using SetFooter() to avoid import cycles.
func NewTableList(tableModel table.Model, width int) *TableList {
	return &TableList{
		table:          tableModel,
		footer:         nil, // Must be set via SetFooter()
		width:          width,
		height:         24,
		emptyMessage:   "No items to display",
		showEmptyState: false,
		errorMessage:   "",
		showError:      false,
		showPagination: false,
		selectedIdx:    0,
		theme:          themes.NewDefaultTheme(),
	}
}

// getTheme returns the theme or default if nil
func (tl *TableList) getTheme() themes.Theme {
	if tl.theme != nil {
		return tl.theme
	}
	return themes.NewDefaultTheme()
}

// SetTable updates the table model being displayed and syncs cursor state
// Important: This pushes the container's selectedIdx to the table, not the other way around
func (tl *TableList) SetTable(tableModel table.Model) *TableList {
	tl.table = tableModel
	// Set the table's cursor to match the container's selected index
	// This ensures indicators in row text match the table's visual selection
	tl.validateAndSyncIdx()
	return tl
}

// SetDimensions sets the width and height of the container
// Note: If you have a footer, update its width separately.
func (tl *TableList) SetDimensions(width, height int) *TableList {
	tl.width = width
	tl.height = height
	tl.table.SetWidth(width)
	tl.table.SetHeight(height - 10)
	return tl
}

// SetPaginationInfo sets the pagination information to display
func (tl *TableList) SetPaginationInfo(info string) *TableList {
	tl.paginationInfo = info
	tl.showPagination = true
	return tl
}

// HidePagination hides the pagination information
func (tl *TableList) HidePagination() *TableList {
	tl.showPagination = false
	return tl
}

// SetEmptyStateMessage sets the message to display when the list is empty
func (tl *TableList) SetEmptyStateMessage(message string) *TableList {
	tl.emptyMessage = message
	tl.showEmptyState = true
	return tl
}

// HideEmptyState hides the empty state message
func (tl *TableList) HideEmptyState() *TableList {
	tl.showEmptyState = false
	return tl
}

// SetErrorMessage sets an error message to display
func (tl *TableList) SetErrorMessage(message string) *TableList {
	tl.errorMessage = message
	tl.showError = true
	return tl
}

// ClearError clears the error message
func (tl *TableList) ClearError() *TableList {
	tl.errorMessage = ""
	tl.showError = false
	return tl
}

// GetWidth returns the current width
func (tl *TableList) GetWidth() int {
	return tl.width
}

// GetHeight returns the current height
func (tl *TableList) GetHeight() int {
	return tl.height
}

// validateAndSyncIdx ensures the selected index is valid and syncs to the table
// This is called after SetTable to push the container's index to the table
func (tl *TableList) validateAndSyncIdx() {
	rowCount := len(tl.table.Rows())

	// If table is empty, set index to 0
	if rowCount == 0 {
		tl.selectedIdx = 0
		tl.table.SetCursor(0)
		return
	}

	// Ensure selectedIdx is within valid bounds
	if tl.selectedIdx < 0 {
		tl.selectedIdx = 0
	} else if tl.selectedIdx >= rowCount {
		tl.selectedIdx = rowCount - 1
	}

	// Sync the table cursor to match the container's index
	tl.table.SetCursor(tl.selectedIdx)
}

// syncIdx ensures the selected index is valid and within bounds
// Deprecated: Use validateAndSyncIdx instead which pushes to table
func (tl *TableList) syncIdx() {
	cursor := tl.table.Cursor()
	// If cursor is negative (table empty or not set), default to 0
	if cursor < 0 {
		tl.selectedIdx = 0
	} else if cursor >= len(tl.table.Rows()) {
		// If cursor is out of bounds, set to last valid index
		tl.selectedIdx = len(tl.table.Rows()) - 1
		if tl.selectedIdx < 0 {
			tl.selectedIdx = 0
		}
	} else {
		tl.selectedIdx = cursor
	}
}

// GetSelectedIdx returns the currently selected index
// Always returns a value >= 0, even if the table is empty
func (tl *TableList) GetSelectedIdx() int {
	// Ensure selectedIdx is always valid
	if tl.selectedIdx < 0 {
		tl.selectedIdx = 0
	}
	if len(tl.table.Rows()) > 0 && tl.selectedIdx >= len(tl.table.Rows()) {
		tl.selectedIdx = len(tl.table.Rows()) - 1
	}
	return tl.selectedIdx
}

// SetSelectedIdx sets the selected index directly
func (tl *TableList) SetSelectedIdx(idx int) *TableList {
	if idx < 0 {
		idx = 0
	}
	if len(tl.table.Rows()) > 0 && idx >= len(tl.table.Rows()) {
		idx = len(tl.table.Rows()) - 1
	}
	tl.selectedIdx = idx
	tl.table.SetCursor(idx)
	return tl
}

// MoveUp moves the cursor up by one position
func (tl *TableList) MoveUp(count int) *TableList {
	newIdx := tl.selectedIdx - count
	if newIdx < 0 {
		newIdx = 0
	}
	tl.selectedIdx = newIdx
	tl.table.SetCursor(newIdx)
	return tl
}

// MoveDown moves the cursor down by one position
func (tl *TableList) MoveDown(count int) *TableList {
	newIdx := tl.selectedIdx + count
	maxIdx := len(tl.table.Rows()) - 1
	if maxIdx < 0 {
		maxIdx = 0
	}
	if newIdx > maxIdx {
		newIdx = maxIdx
	}
	tl.selectedIdx = newIdx
	tl.table.SetCursor(newIdx)
	return tl
}

// MoveToFirst moves the cursor to the first item
func (tl *TableList) MoveToFirst() *TableList {
	tl.selectedIdx = 0
	tl.table.SetCursor(0)
	return tl
}

// MoveToLast moves the cursor to the last item
func (tl *TableList) MoveToLast() *TableList {
	maxIdx := len(tl.table.Rows()) - 1
	if maxIdx < 0 {
		maxIdx = 0
	}
	tl.selectedIdx = maxIdx
	tl.table.SetCursor(maxIdx)
	return tl
}

// SyncCursorFromTable syncs the container's selectedIdx from the table's cursor
func (tl *TableList) SyncCursorFromTable() *TableList {
	tl.syncIdx()
	return tl
}

// WithTheme sets the theme for the container
// Note: If you have a footer, set its theme separately before calling SetFooter().
func (tl *TableList) WithTheme(theme themes.Theme) *TableList {
	tl.theme = theme
	return tl
}

// Footer returns the footer for external configuration
func (tl *TableList) Footer() FooterRenderer {
	return tl.footer
}

// SetFooter sets a custom footer
func (tl *TableList) SetFooter(footer FooterRenderer) *TableList {
	tl.footer = footer
	return tl
}

// Table returns the underlying table model for external configuration
func (tl *TableList) Table() table.Model {
	return tl.table
}

// Render returns the complete table list view with all components
func (tl *TableList) Render() string {
	theme := tl.getTheme()

	footerView := ""
	if tl.footer != nil {
		footerView = tl.footer.View()
	}

	// Handle error state with theme-based error box
	if tl.showError && tl.errorMessage != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(theme.ErrorColor()).
			Background(lipgloss.Color("#2D1F1F")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.ErrorColor()).
			Padding(0, 1)
		errorContent := errorStyle.Render(tl.errorMessage)

		parts := []string{errorContent}
		if footerView != "" {
			parts = append(parts, "", footerView)
		}
		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	}

	// Handle empty state
	if tl.showEmptyState && len(tl.table.Rows()) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(theme.MutedColor()).
			Italic(true)
		emptyContent := emptyStyle.Render(tl.emptyMessage)

		parts := []string{emptyContent}

		// Add pagination even for empty state
		if tl.showPagination && tl.paginationInfo != "" {
			paginationStyle := lipgloss.NewStyle().
				Foreground(theme.MutedColor()).
				MarginTop(1)
			paginationView := paginationStyle.Render(tl.paginationInfo)
			parts = append(parts, "", paginationView)
		}

		if footerView != "" {
			parts = append(parts, "", footerView)
		}

		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	}

	// Render table
	tableView := tl.table.View()

	parts := []string{tableView}

	// Add pagination if enabled
	if tl.showPagination && tl.paginationInfo != "" {
		paginationStyle := lipgloss.NewStyle().
			Foreground(theme.MutedColor()).
			MarginTop(1)
		paginationView := paginationStyle.Render(tl.paginationInfo)
		parts = append(parts, "", paginationView)
	}

	if footerView != "" {
		parts = append(parts, "", footerView)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
