package behaviors

import (
	"github.com/baphled/kariya/internal/cli/uikit/theme"
)

// FilterMenuBehavior[T] provides a sectioned menu UI for filtering table data.
// It manages navigation, selection state, and applies filter predicates to the table.
type FilterMenuBehavior[T any] struct {
	theme.Aware
}

// NewFilterMenuBehavior creates a new filter menu behavior.
// Panics if table is nil.
func NewFilterMenuBehavior[T any](themeObj theme.Theme, table *TableBehavior[T], title string) *FilterMenuBehavior[T] {
	if table == nil {
		panic("table cannot be nil")
	}

	return &FilterMenuBehavior[T]{}
}

// AddSection adds a section with options to the menu.
func (f *FilterMenuBehavior[T]) AddSection(section MenuSection) *FilterMenuBehavior[T] {
	return f
}

// OnApply sets the callback to invoke when a filter is applied.
func (f *FilterMenuBehavior[T]) OnApply(callback func()) *FilterMenuBehavior[T] {
	return f
}

// Show activates the filter menu.
func (f *FilterMenuBehavior[T]) Show() {
}

// Hide deactivates the filter menu.
func (f *FilterMenuBehavior[T]) Hide() {
}

// IsActive returns true if the filter menu is currently active.
func (f *FilterMenuBehavior[T]) IsActive() bool {
	return false
}

// HandleKey processes key presses for navigation and selection.
// Returns true if the key was handled.
func (f *FilterMenuBehavior[T]) HandleKey(key string) bool {
	return false
}

// Render returns the themed menu view.
func (f *FilterMenuBehavior[T]) Render() string {
	return ""
}
