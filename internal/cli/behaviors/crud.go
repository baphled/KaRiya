package behaviors

import (
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
)

// CRUDBehavior[T] provides create/edit/delete operations for a table.
// It handles key bindings (n/e/d), delete confirmation modals, and mode management.
type CRUDBehavior[T any] struct {
	theme.Aware
}

// NewCRUDBehavior creates a new CRUD behavior for the given table.
// Panics if table is nil.
func NewCRUDBehavior[T any](themeObj theme.Theme, table *TableBehavior[T]) *CRUDBehavior[T] {
	if table == nil {
		panic("table cannot be nil")
	}

	return &CRUDBehavior[T]{
		Aware: theme.Aware{},
	}
}

// Mode returns the current CRUD mode.
func (c *CRUDBehavior[T]) Mode() CRUDMode {
	return ModeList
}

// IsInListMode returns true if in list mode.
func (c *CRUDBehavior[T]) IsInListMode() bool {
	return false
}

// ReturnToList returns to list mode.
func (c *CRUDBehavior[T]) ReturnToList() {
}

// ItemNamer sets the function to get an item's name for display.
func (c *CRUDBehavior[T]) ItemNamer(namer func(T) string) *CRUDBehavior[T] {
	return c
}

// OnCreate sets the create callback and enables create operations.
func (c *CRUDBehavior[T]) OnCreate(callback func() tea.Cmd) *CRUDBehavior[T] {
	return c
}

// OnEdit sets the edit callback and enables edit operations.
func (c *CRUDBehavior[T]) OnEdit(callback func(T) tea.Cmd) *CRUDBehavior[T] {
	return c
}

// OnDelete sets the delete callback and enables delete operations.
func (c *CRUDBehavior[T]) OnDelete(callback func(T) tea.Cmd) *CRUDBehavior[T] {
	return c
}

// HandleKey processes CRUD operation keys (n/e/d).
// Returns (cmd, true) if key was handled, (nil, false) otherwise.
func (c *CRUDBehavior[T]) HandleKey(key string) (tea.Cmd, bool) {
	return nil, false
}

// Update processes messages when in delete confirmation mode.
// Returns (cmd, true) if message was handled, (nil, false) otherwise.
func (c *CRUDBehavior[T]) Update(msg tea.Msg) (tea.Cmd, bool) {
	return nil, false
}

// RenderDeleteConfirm renders the delete confirmation modal.
func (c *CRUDBehavior[T]) RenderDeleteConfirm() string {
	return ""
}

// RenderHelpKeys renders help text for enabled CRUD operations.
func (c *CRUDBehavior[T]) RenderHelpKeys() string {
	return ""
}
