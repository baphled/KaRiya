package behaviors

import (
	"github.com/baphled/kariya/internal/tui/navigation"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// BadgeFunc is a function that creates a Badge given a theme.
// Domain views pass these to ListBehavior so badge rendering stays generic.
type BadgeFunc func(theme.Theme) *primitives.Badge

// ListBehavior is a generic, domain-agnostic table and help-footer widget.
// It wraps TableBehavior[T] and handles its own navigation key messages.
// Domain views compose it and supply column definitions, a row formatter,
// and badge functions — keeping all domain knowledge in the view layer.
type ListBehavior[T any] struct {
	table      *TableBehavior[T]
	keys       navigation.ListKeyMap
	badgeFuncs []BadgeFunc
	th         theme.Theme
}

// NewListBehavior constructs a ListBehavior with the given column definitions,
// row formatter, and badge functions.
//
// Returns:
//   - A pointer to the initialized ListBehavior[T].
//
// Expected:
//   - columns must be a non-empty slice of ColumnDef.
//   - formatter must be a valid function that converts items to string slices.
//   - badgeFuncs may be empty or contain badge factory functions.
//
// Side effects:
//   - Allocates a new TableBehavior and initializes internal state.
func NewListBehavior[T any](
	columns []ColumnDef,
	formatter func(T, int) []string,
	badgeFuncs []BadgeFunc,
) *ListBehavior[T] {
	tb := NewTableBehavior[T](nil, columns, formatter).
		PageSize(15).
		PaginationPrefix("Items").
		EmptyMessage("No items found.")
	return &ListBehavior[T]{
		table:      tb,
		keys:       navigation.DefaultListKeyMap(),
		badgeFuncs: badgeFuncs,
	}
}

// PaginationPrefix sets the label shown before item counts in pagination.
//
// Returns:
//   - A pointer to the ListBehavior[T] for method chaining.
//
// Expected:
//   - prefix should be a non-empty string describing the item type.
//
// Side effects:
//   - Updates the internal TableBehavior's pagination prefix.
func (w *ListBehavior[T]) PaginationPrefix(prefix string) *ListBehavior[T] {
	w.table.PaginationPrefix(prefix)
	return w
}

// EmptyMessage sets the placeholder text displayed when the list is empty.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized ListBehavior[T] ready for use.
//
// Side effects:
//   - None.
func (w *ListBehavior[T]) EmptyMessage(msg string) *ListBehavior[T] {
	w.table.EmptyMessage(msg)
	return w
}

// SetItems replaces the widget's item list.
//
// Expected:
//   - []t must be valid.
//
// Side effects:
//   - None.
func (w *ListBehavior[T]) SetItems(items []T) {
	w.table.SetItems(items)
}

// SetTheme updates the theme used for help-text rendering.
//
// Expected:
//   - th must be a valid theme instance (can be nil).
//
// Side effects:
//   - None.
func (w *ListBehavior[T]) SetTheme(th theme.Theme) {
	w.th = th
}

// GetSelectedIndex returns the current selection index.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (w *ListBehavior[T]) GetSelectedIndex() int {
	return w.table.GetSelectedIndex()
}

// GetSelectedItem returns a pointer to the currently selected item, or nil.
//
// Returns:
//   - A fully initialized T ready for use.
//
// Side effects:
//   - None.
func (w *ListBehavior[T]) GetSelectedItem() *T {
	return w.table.GetSelectedItem()
}

// HandleKey processes a key message and returns true if the key was consumed.
//
// Expected:
//   - keymsg must be valid.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (w *ListBehavior[T]) HandleKey(msg tea.KeyMsg) bool {
	switch {
	case key.Matches(msg, w.keys.Up),
		key.Matches(msg, w.keys.Down),
		key.Matches(msg, w.keys.PageUp),
		key.Matches(msg, w.keys.PageDown),
		key.Matches(msg, w.keys.GoToStart),
		key.Matches(msg, w.keys.GoToEnd):
		w.table.HandleNavigation(msg.String())
		return true
	case key.Matches(msg, w.keys.Select):
		return true
	}
	return false
}

// RenderContent renders the table as a string — no chrome.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (w *ListBehavior[T]) RenderContent() string {
	return w.table.Render()
}

// HelpText renders the key-binding footer using the stored badge functions.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (w *ListBehavior[T]) HelpText() string {
	th := w.th
	if th == nil {
		th = themes.NewDefaultTheme()
	}
	badges := make([]*primitives.Badge, len(w.badgeFuncs))
	for i, fn := range w.badgeFuncs {
		badges[i] = fn(th)
	}
	return primitives.RenderHelpFooter(th, badges...)
}
