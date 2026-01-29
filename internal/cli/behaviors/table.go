package behaviors

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// TableBehavior provides data binding, pagination, navigation, filtering,
// and sorting for table-based list views.
//
// The generic type parameter T is the domain item displayed in each row
// (e.g., *career.Event or *career.Skill). T must satisfy the "any"
// constraint; pointer types are typical.
//
// TableBehavior implements the ListNavigator interface so that screens can
// delegate all keyboard-driven list movement (up, down, j, k, page-up,
// page-down, home, end, g, G) without custom logic. Pagination is
// calculated automatically from the configured PageSize. Sorting and
// filtering are applied lazily via SetSort and SetFilter; the display list
// is recalculated on the next read when items are marked dirty.
//
// Screens embed a *TableBehavior[T], call SetItems to bind domain data,
// HandleNavigation in their Update method to process key events, and Render
// in their View method to obtain the final table string including a
// pagination footer line.
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

	columns      []ColumnDef
	rowFormatter RowFormatter[T]
	pageSize     int

	allItems      []T
	displayItems  []T
	selectedIndex int

	filterPredicate FilterPredicate[T]
	sortComparator  SortComparator[T]
	sortReverse     bool

	emptyMessage     string
	paginationPrefix string
	showPagination   bool

	table         table.Model
	navHandler    *navigation.ListNavigationHandler
	width, height int
	needsRefresh  bool
}

// NewTableBehavior creates a new table behavior with the given configuration.
func NewTableBehavior[T any](themeObj themes.Theme, columns []ColumnDef, formatter RowFormatter[T]) *TableBehavior[T] {
	bubbleColumns := make([]table.Column, len(columns))
	for i, col := range columns {
		bubbleColumns[i] = table.Column{
			Title: col.Title,
			Width: col.Width,
		}
	}

	bubblesTable := table.New(
		table.WithColumns(bubbleColumns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
	)

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

	behavior.SetTheme(themeObj)
	behavior.navHandler = navigation.NewListNavigationHandler(behavior)

	return behavior
}

// PageSize sets items per page (default: 15).
func (tb *TableBehavior[T]) PageSize(size int) *TableBehavior[T] {
	tb.pageSize = size
	return tb
}

// EmptyMessage sets the message shown when no items exist.
func (tb *TableBehavior[T]) EmptyMessage(msg string) *TableBehavior[T] {
	tb.emptyMessage = msg
	return tb
}

// PaginationPrefix sets the prefix for pagination info ("Events", "Skills").
func (tb *TableBehavior[T]) PaginationPrefix(prefix string) *TableBehavior[T] {
	tb.paginationPrefix = prefix
	return tb
}

// Dimensions sets the container width and height.
func (tb *TableBehavior[T]) Dimensions(width, height int) *TableBehavior[T] {
	tb.width = width
	tb.height = height
	tb.table.SetWidth(width)
	tb.table.SetHeight(height - 10)
	return tb
}

// HidePagination hides pagination info.
func (tb *TableBehavior[T]) HidePagination() *TableBehavior[T] {
	tb.showPagination = false
	return tb
}

// ShowPagination shows pagination info (default).
func (tb *TableBehavior[T]) ShowPagination() *TableBehavior[T] {
	tb.showPagination = true
	return tb
}

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

// IsEmpty returns true if there are no items to display.
func (tb *TableBehavior[T]) IsEmpty() bool {
	tb.refreshDisplayItems()
	return len(tb.displayItems) == 0
}

// Count returns the number of displayed items.
func (tb *TableBehavior[T]) Count() int {
	tb.refreshDisplayItems()
	return len(tb.displayItems)
}

// TotalCount returns the total number of items (before filtering).
func (tb *TableBehavior[T]) TotalCount() int {
	return len(tb.allItems)
}

// GetTotalItems implements ListNavigator.
func (tb *TableBehavior[T]) GetTotalItems() int {
	return tb.Count()
}

// SetSelectedIndex implements ListNavigator.
func (tb *TableBehavior[T]) SetSelectedIndex(idx int) {
	tb.refreshDisplayItems()

	if len(tb.displayItems) == 0 {
		tb.selectedIndex = 0
		return
	}

	if idx < 0 {
		idx = 0
	}
	if idx >= len(tb.displayItems) {
		idx = len(tb.displayItems) - 1
	}

	tb.selectedIndex = idx
	tb.updateTableRows()
}

// GetPageSize implements ListNavigator.
func (tb *TableBehavior[T]) GetPageSize() int {
	return tb.pageSize
}

// HandleNavigation processes navigation keys (up, down, j, k, pgup, pgdn, home, end, g, G).
// Returns true if the key was handled.
func (tb *TableBehavior[T]) HandleNavigation(keyStr string) bool {
	if len(tb.displayItems) == 0 {
		return false
	}
	return tb.navHandler.HandleKey(keyStr)
}

// SetFilter applies a filter predicate.
// Pass nil to clear the filter.
// Selection is preserved if the previously selected item passes the filter.
func (tb *TableBehavior[T]) SetFilter(pred FilterPredicate[T]) *TableBehavior[T] {
	var selectedItem *T
	if tb.selectedIndex >= 0 && tb.selectedIndex < len(tb.displayItems) {
		selectedItem = &tb.displayItems[tb.selectedIndex]
	}

	tb.filterPredicate = pred
	tb.needsRefresh = true
	tb.refreshDisplayItems()

	if selectedItem != nil {
		for i, item := range tb.displayItems {
			if compareItems(*selectedItem, item) {
				tb.selectedIndex = i
				tb.updateTableRows()
				return tb
			}
		}
	}

	tb.selectedIndex = 0
	tb.updateTableRows()
	return tb
}

// ClearFilter removes the active filter.
func (tb *TableBehavior[T]) ClearFilter() *TableBehavior[T] {
	tb.filterPredicate = nil
	tb.needsRefresh = true
	tb.refreshDisplayItems()
	tb.updateTableRows()
	return tb
}

// HasFilter returns true if a filter is active.
func (tb *TableBehavior[T]) HasFilter() bool {
	return tb.filterPredicate != nil
}

// SetSort applies a sort comparator.
// Pass nil to use the original order.
func (tb *TableBehavior[T]) SetSort(cmp SortComparator[T], reverse bool) *TableBehavior[T] {
	var selectedItem *T
	if tb.selectedIndex >= 0 && tb.selectedIndex < len(tb.displayItems) {
		selectedItem = &tb.displayItems[tb.selectedIndex]
	}

	tb.sortComparator = cmp
	tb.sortReverse = reverse
	tb.needsRefresh = true
	tb.refreshDisplayItems()

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

// ClearSort removes the active sort.
func (tb *TableBehavior[T]) ClearSort() *TableBehavior[T] {
	tb.sortComparator = nil
	tb.sortReverse = false
	tb.needsRefresh = true
	tb.refreshDisplayItems()
	tb.updateTableRows()
	return tb
}

// HasSort returns true if a sort is active.
func (tb *TableBehavior[T]) HasSort() bool {
	return tb.sortComparator != nil
}

// Render returns the complete table view with pagination.
func (tb *TableBehavior[T]) Render() string {
	tb.refreshDisplayItems()

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

	tb.updateTableRows()

	parts := []string{tb.table.View()}

	if tb.showPagination {
		parts = append(parts, "", tb.RenderPaginationInfo())
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// RenderPaginationInfo returns the pagination string.
// Example: "Events: 42 | Page 2 of 5".
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

// refreshDisplayItems recalculates displayItems from allItems applying filter and sort.
func (tb *TableBehavior[T]) refreshDisplayItems() {
	if !tb.needsRefresh {
		return
	}

	result := make([]T, len(tb.allItems))
	copy(result, tb.allItems)

	if tb.filterPredicate != nil {
		filtered := make([]T, 0, len(result))
		for _, item := range result {
			if tb.filterPredicate(item) {
				filtered = append(filtered, item)
			}
		}
		result = filtered
	}

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

	if tb.selectedIndex >= len(tb.displayItems) {
		if len(tb.displayItems) > 0 {
			tb.selectedIndex = len(tb.displayItems) - 1
		} else {
			tb.selectedIndex = 0
		}
	}
}

// updateTableRows syncs the bubbles/table model with current page.
func (tb *TableBehavior[T]) updateTableRows() {
	if len(tb.displayItems) == 0 {
		tb.table.SetRows([]table.Row{})
		return
	}

	page := tb.selectedIndex / tb.pageSize
	start := page * tb.pageSize
	end := start + tb.pageSize
	if end > len(tb.displayItems) {
		end = len(tb.displayItems)
	}

	pageItems := tb.displayItems[start:end]

	rows := make([]table.Row, len(pageItems))
	for i, item := range pageItems {
		realIdx := start + i
		cells := tb.rowFormatter(item, realIdx)

		if realIdx == tb.selectedIndex {
			cells[0] = tb.navHandler.FormatRowText(realIdx, cells[0])
		} else {
			cells[0] = "  " + cells[0]
		}

		rows[i] = table.Row(cells)
	}

	tb.table.SetRows(rows)

	relativeCursor := tb.selectedIndex - start
	if relativeCursor < 0 {
		relativeCursor = 0
	}
	if relativeCursor >= len(rows) {
		relativeCursor = len(rows) - 1
	}
	tb.table.SetCursor(relativeCursor)
}

// compareItems checks whether two items of type T are deeply equal.
func compareItems[T any](a, b T) bool {
	return reflect.DeepEqual(a, b)
}
