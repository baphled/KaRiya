package behaviors

import (
	"github.com/baphled/kariya/internal/cli/uikit/theme"
)

// SortOption[T] represents a sort option with a label, comparator, and reverse flag.
type SortOption[T any] struct {
	Label      string
	Comparator SortComparator[T]
	Reverse    bool
}

// SortMenuBehavior[T] provides a menu UI for sorting table data.
// It manages navigation, selection state, and applies sort comparators to the table.
type SortMenuBehavior[T any] struct {
	theme.Aware
}

// NewSortMenuBehavior creates a new sort menu behavior.
// Panics if table is nil.
func NewSortMenuBehavior[T any](themeObj theme.Theme, table *TableBehavior[T], title string) *SortMenuBehavior[T] {
	if table == nil {
		panic("table cannot be nil")
	}

	return &SortMenuBehavior[T]{}
}

// AddOption adds a sort option to the menu.
func (s *SortMenuBehavior[T]) AddOption(label string, comparator SortComparator[T], reverse bool) *SortMenuBehavior[T] {
	return s
}

// OnApply sets the callback to invoke when a sort is applied.
func (s *SortMenuBehavior[T]) OnApply(callback func()) *SortMenuBehavior[T] {
	return s
}

// Show activates the sort menu.
func (s *SortMenuBehavior[T]) Show() {
}

// Hide deactivates the sort menu.
func (s *SortMenuBehavior[T]) Hide() {
}

// IsActive returns true if the sort menu is currently active.
func (s *SortMenuBehavior[T]) IsActive() bool {
	return false
}

// HandleKey processes key presses for navigation and selection.
// Returns true if the key was handled.
func (s *SortMenuBehavior[T]) HandleKey(key string) bool {
	return false
}

// Render returns the themed menu view.
func (s *SortMenuBehavior[T]) Render() string {
	return ""
}
