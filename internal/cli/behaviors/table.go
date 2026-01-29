package behaviors

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// TableBehavior[T] provides data binding, pagination, navigation, filtering, and sorting
// for table-based list views. It implements the ListNavigator interface and can be
// embedded in intents to eliminate boilerplate table management code.
//
// Usage:
//
//	table := behaviors.NewTableBehavior(theme, columns, formatter).
//	    PageSize(20).
//	    EmptyMessage("No items found")
//	table.SetItems(myItems)
//
//	// In Update:
//	if table.HandleNavigation(keyStr) {
//	    return nil
//	}
//
//	// In View:
//	return table.Render()
type TableBehavior[T any] struct {
	theme.Aware

	// Configuration (immutable after creation)
	columns      []ColumnDef
	rowFormatter RowFormatter[T]
	pageSize     int

	// Data (mutable)
	allItems      []T // Original unfiltered/unsorted items
	displayItems  []T // Items after filter/sort applied
	selectedIndex int // Index in displayItems (0-based)

	// Filter/Sort state
	filterPredicate FilterPredicate[T]
	sortComparator  SortComparator[T]
	sortReverse     bool

	// Display options
	emptyMessage     string
	paginationPrefix string
	showPagination   bool

	// Internal
	table         table.Model
	viewport      viewport.Model
	useViewport   bool
	navHandler    *navigation.ListNavigationHandler
	width, height int
	needsRefresh  bool
}

// NewTableBehavior creates a new table behavior with the given configuration.
func NewTableBehavior[T any](themeObj themes.Theme, columns []ColumnDef, formatter RowFormatter[T]) *TableBehavior[T] {
	// Convert columns to bubbles table columns
	bubbleColumns := make([]table.Column, len(columns))
	for i, col := range columns {
		bubbleColumns[i] = table.Column{
			Title: col.Title,
			Width: col.Width,
		}
	}

	// Create bubbles table
	bubblesTable := table.New(
		table.WithColumns(bubbleColumns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	// Create behavior
	behavior := &TableBehavior[T]{
		columns:          columns,
		rowFormatter:     formatter,
		pageSize:         15,
		emptyMessage:     "No items to display",
		paginationPrefix: "Items",
		showPagination:   true,
		table:            bubblesTable,
		width:            100,
		height:           24,
		allItems:         []T{},
		displayItems:     []T{},
		selectedIndex:    0,
	}

	// Set theme
	behavior.SetTheme(themeObj)

	// Create navigation handler
	behavior.navHandler = navigation.NewListNavigationHandler(behavior)

	return behavior
}

// ============= Configuration (fluent, chainable) =============

// PageSize sets items per page (default: 15)
func (tb *TableBehavior[T]) PageSize(size int) *TableBehavior[T] {
	tb.pageSize = size
	return tb
}

// EmptyMessage sets the message shown when no items exist
func (tb *TableBehavior[T]) EmptyMessage(msg string) *TableBehavior[T] {
	tb.emptyMessage = msg
	return tb
}

// PaginationPrefix sets the prefix for pagination info ("Events", "Skills")
func (tb *TableBehavior[T]) PaginationPrefix(prefix string) *TableBehavior[T] {
	tb.paginationPrefix = prefix
	return tb
}

// Dimensions sets the container width and height
func (tb *TableBehavior[T]) Dimensions(width, height int) *TableBehavior[T] {
	tb.width = width
	tb.height = height
	tb.table.SetWidth(width)
	tb.table.SetHeight(height - 10) // Reserve space for pagination/footer
	return tb
}

// SetHeight configures the table to use viewport with the specified height.
// This enables scrolling when content exceeds the available height.
// The height should be the available content height from ScreenLayout.GetAvailableContentHeight().
func (tb *TableBehavior[T]) SetHeight(height int) *TableBehavior[T] {
	tb.height = height
	tb.useViewport = true

	// Initialize or update viewport
	if tb.viewport.Width == 0 {
		// First time - create viewport
		tb.viewport = viewport.New(tb.width, height)
	} else {
		// Update existing viewport
		tb.viewport.Width = tb.width
		tb.viewport.Height = height
	}

	// Update table to not limit its own height - viewport will handle scrolling
	tb.table.SetHeight(1000) // Large height so table doesn't truncate

	return tb
}

// HidePagination hides pagination info
func (tb *TableBehavior[T]) HidePagination() *TableBehavior[T] {
	tb.showPagination = false
	return tb
}

// ShowPagination shows pagination info (default)
func (tb *TableBehavior[T]) ShowPagination() *TableBehavior[T] {
	tb.showPagination = true
	return tb
}

// ============= Data Management =============

// SetItems replaces all items and resets to first item.
func (tb *TableBehavior[T]) SetItems(items []T) *TableBehavior[T] {
	if items == nil {
		items = []T{}
	}
	tb.allItems = items
	tb.selectedIndex = 0
	tb.needsRefresh = true
	tb.refreshDisplayItems()
	return tb
}

// GetItems returns the current items after filter/sort.
func (tb *TableBehavior[T]) GetItems() []T {
	tb.refreshDisplayItems()
	return tb.displayItems
}

// GetSelectedItem returns the currently selected item, or nil if empty.
func (tb *TableBehavior[T]) GetSelectedItem() *T {
	tb.refreshDisplayItems()
	if len(tb.displayItems) == 0 {
		return nil
	}
	if tb.selectedIndex < 0 || tb.selectedIndex >= len(tb.displayItems) {
		return nil
	}
	return &tb.displayItems[tb.selectedIndex]
}

// GetSelectedIndex returns the current selection index (0-based in displayItems).
func (tb *TableBehavior[T]) GetSelectedIndex() int {
	return tb.selectedIndex
}

// IsEmpty returns true if there are no items to display
func (tb *TableBehavior[T]) IsEmpty() bool {
	tb.refreshDisplayItems()
	return len(tb.displayItems) == 0
}

// Count returns the number of displayed items
func (tb *TableBehavior[T]) Count() int {
	tb.refreshDisplayItems()
	return len(tb.displayItems)
}

// TotalCount returns the total number of items (before filtering)
func (tb *TableBehavior[T]) TotalCount() int {
	return len(tb.allItems)
}

// ============= ListNavigator Interface =============

// GetTotalItems implements ListNavigator
func (tb *TableBehavior[T]) GetTotalItems() int {
	return tb.Count()
}

// SetSelectedIndex implements ListNavigator
func (tb *TableBehavior[T]) SetSelectedIndex(idx int) {
	tb.refreshDisplayItems()

	if len(tb.displayItems) == 0 {
		tb.selectedIndex = 0
		return
	}

	// Clamp to valid range
	if idx < 0 {
		idx = 0
	}
	if idx >= len(tb.displayItems) {
		idx = len(tb.displayItems) - 1
	}

	tb.selectedIndex = idx
	tb.updateTableRows()

	// Sync viewport scrolling if viewport is enabled
	if tb.useViewport && tb.viewport.Height > 0 {
		// Each row is approximately 1 line (table header + rows)
		// Scroll viewport to keep selected row visible
		rowHeight := 1
		selectedLinePosition := idx * rowHeight

		// If selected line is below viewport bottom, scroll down
		if selectedLinePosition >= tb.viewport.YOffset+tb.viewport.Height {
			tb.viewport.SetYOffset(selectedLinePosition - tb.viewport.Height + 1)
		}

		// If selected line is above viewport top, scroll up
		if selectedLinePosition < tb.viewport.YOffset {
			tb.viewport.SetYOffset(selectedLinePosition)
		}
	}
}

// GetPageSize implements ListNavigator
func (tb *TableBehavior[T]) GetPageSize() int {
	return tb.pageSize
}

// ============= Navigation =============

// HandleNavigation processes navigation keys (up, down, j, k, pgup, pgdn, home, end, g, G).
// Returns true if the key was handled.
func (tb *TableBehavior[T]) HandleNavigation(keyStr string) bool {
	if len(tb.displayItems) == 0 {
		return false
	}
	return tb.navHandler.HandleKey(keyStr)
}

// ============= Filtering =============

// SetFilter applies a filter predicate.
// Pass nil to clear the filter.
// Selection is preserved if the previously selected item passes the filter.
func (tb *TableBehavior[T]) SetFilter(pred FilterPredicate[T]) *TableBehavior[T] {
	// Capture currently selected item if any
	var selectedItem *T
	if tb.selectedIndex >= 0 && tb.selectedIndex < len(tb.displayItems) {
		selectedItem = &tb.displayItems[tb.selectedIndex]
	}

	tb.filterPredicate = pred
	tb.needsRefresh = true
	tb.refreshDisplayItems()

	// Try to preserve selection
	if selectedItem != nil {
		for i, item := range tb.displayItems {
			// Check if this is the same item (by comparing all fields)
			if compareItems(*selectedItem, item) {
				tb.selectedIndex = i
				tb.updateTableRows()
				return tb
			}
		}
	}

	// Selection couldn't be preserved, reset to first
	tb.selectedIndex = 0
	tb.updateTableRows()
	return tb
}

// ClearFilter removes the active filter
func (tb *TableBehavior[T]) ClearFilter() *TableBehavior[T] {
	tb.filterPredicate = nil
	tb.needsRefresh = true
	tb.refreshDisplayItems()
	tb.updateTableRows()
	return tb
}

// HasFilter returns true if a filter is active
func (tb *TableBehavior[T]) HasFilter() bool {
	return tb.filterPredicate != nil
}

// ============= Sorting =============

// SetSort applies a sort comparator.
// Pass nil to use the original order.
func (tb *TableBehavior[T]) SetSort(cmp SortComparator[T], reverse bool) *TableBehavior[T] {
	// Capture currently selected item if any
	var selectedItem *T
	if tb.selectedIndex >= 0 && tb.selectedIndex < len(tb.displayItems) {
		selectedItem = &tb.displayItems[tb.selectedIndex]
	}

	tb.sortComparator = cmp
	tb.sortReverse = reverse
	tb.needsRefresh = true
	tb.refreshDisplayItems()

	// Try to preserve selection
	if selectedItem != nil {
		for i, item := range tb.displayItems {
			if compareItems(*selectedItem, item) {
				tb.selectedIndex = i
				tb.updateTableRows()
				return tb
			}
		}
	}

	tb.updateTableRows()
	return tb
}

// ClearSort removes the active sort
func (tb *TableBehavior[T]) ClearSort() *TableBehavior[T] {
	tb.sortComparator = nil
	tb.sortReverse = false
	tb.needsRefresh = true
	tb.refreshDisplayItems()
	tb.updateTableRows()
	return tb
}

// HasSort returns true if a sort is active
func (tb *TableBehavior[T]) HasSort() bool {
	return tb.sortComparator != nil
}

// ============= Rendering =============

// Render returns the complete table view with pagination.
func (tb *TableBehavior[T]) Render() string {
	tb.refreshDisplayItems()

	// Handle empty state
	if len(tb.displayItems) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(tb.MutedColor()).
			Italic(true)
		emptyContent := emptyStyle.Render(tb.emptyMessage)

		parts := []string{emptyContent}

		if tb.showPagination {
			parts = append(parts, "", tb.RenderPaginationInfo())
		}

		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	}

	// Update table rows for current page
	tb.updateTableRows()

	// Render table
	tableView := tb.table.View()
	parts := []string{tableView}

	// Add pagination if enabled
	if tb.showPagination {
		parts = append(parts, "", tb.RenderPaginationInfo())
	}

	combined := lipgloss.JoinVertical(lipgloss.Left, parts...)

	// Use viewport if enabled
	if tb.useViewport {
		tb.viewport.SetContent(combined)
		return tb.viewport.View()
	}

	return combined
}

// RenderPaginationInfo returns the pagination string.
// Example: "Events: 42 | Page 2 of 5"
func (tb *TableBehavior[T]) RenderPaginationInfo() string {
	tb.refreshDisplayItems()

	totalItems := len(tb.displayItems)
	if totalItems == 0 {
		return fmt.Sprintf("%s: 0", tb.paginationPrefix)
	}

	currentPage := (tb.selectedIndex / tb.pageSize) + 1
	totalPages := (totalItems + tb.pageSize - 1) / tb.pageSize

	return fmt.Sprintf("%s: %d | Page %d of %d", tb.paginationPrefix, totalItems, currentPage, totalPages)
}

// ============= Internal =============

// refreshDisplayItems recalculates displayItems from allItems applying filter and sort
func (tb *TableBehavior[T]) refreshDisplayItems() {
	if !tb.needsRefresh {
		return
	}

	// Start with all items
	result := make([]T, len(tb.allItems))
	copy(result, tb.allItems)

	// Apply filter if active
	if tb.filterPredicate != nil {
		filtered := make([]T, 0, len(result))
		for _, item := range result {
			if tb.filterPredicate(item) {
				filtered = append(filtered, item)
			}
		}
		result = filtered
	}

	// Apply sort if active
	if tb.sortComparator != nil {
		sort.Slice(result, func(i, j int) bool {
			cmp := tb.sortComparator(result[i], result[j])
			if tb.sortReverse {
				return cmp > 0
			}
			return cmp < 0
		})
	}

	tb.displayItems = result
	tb.needsRefresh = false

	// Ensure selection is valid
	if tb.selectedIndex >= len(tb.displayItems) {
		if len(tb.displayItems) > 0 {
			tb.selectedIndex = len(tb.displayItems) - 1
		} else {
			tb.selectedIndex = 0
		}
	}
}

// updateTableRows syncs the bubbles/table model with current page
func (tb *TableBehavior[T]) updateTableRows() {
	if len(tb.displayItems) == 0 {
		tb.table.SetRows([]table.Row{})
		return
	}

	// Calculate current page
	page := tb.selectedIndex / tb.pageSize
	start := page * tb.pageSize
	end := start + tb.pageSize
	if end > len(tb.displayItems) {
		end = len(tb.displayItems)
	}

	// Get items for this page
	pageItems := tb.displayItems[start:end]

	// Generate rows with selection indicator
	rows := make([]table.Row, len(pageItems))
	for i, item := range pageItems {
		realIdx := start + i
		cells := tb.rowFormatter(item, realIdx)

		// Add selection indicator to first column
		if realIdx == tb.selectedIndex {
			cells[0] = tb.navHandler.FormatRowText(realIdx, cells[0])
		} else {
			cells[0] = "  " + cells[0]
		}

		rows[i] = table.Row(cells)
	}

	tb.table.SetRows(rows)

	// Set cursor to relative position within page
	relativeCursor := tb.selectedIndex - start
	if relativeCursor < 0 {
		relativeCursor = 0
	}
	if relativeCursor >= len(rows) {
		relativeCursor = len(rows) - 1
	}
	tb.table.SetCursor(relativeCursor)
}

// compareItems is a helper to check if two items are the same.
// Uses reflection for deep equality checking to handle all types.
func compareItems[T any](a, b T) bool {
	return reflect.DeepEqual(a, b)
}
