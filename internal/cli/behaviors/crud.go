package behaviors

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
)

// CRUDBehavior[T] provides create/edit/delete operations for a table.
// It handles key bindings (n/e/d), delete confirmation modals, and mode management.
//
// Usage:
//
//	crud := behaviors.NewCRUDBehavior(theme, table).
//	    ItemNamer(func(item Skill) string { return item.Name }).
//	    OnCreate(func() tea.Cmd { return showCreateForm }).
//	    OnEdit(func(item Skill) tea.Cmd { return showEditForm(item) }).
//	    OnDelete(func(item Skill) tea.Cmd { return deleteSkill(item) })
//
//	// In Update:
//	if cmd, handled := crud.HandleKey(key); handled {
//	    return cmd
//	}
//
//	// In View:
//	if crud.Mode() == ModeDelete {
//	    return crud.RenderDeleteConfirm()
//	}
type CRUDBehavior[T any] struct {
	theme.Aware

	// References
	table *TableBehavior[T]

	// Configuration
	itemNamer func(T) string

	// Callbacks
	onCreate func() tea.Cmd
	onEdit   func(T) tea.Cmd
	onDelete func(T) tea.Cmd

	// State
	mode         CRUDMode
	confirmModal *components.ModalContent
	itemToDelete *T

	// Feature flags
	enableCreate bool
	enableEdit   bool
	enableDelete bool

	// Dimensions for modal rendering
	width  int
	height int
}

// NewCRUDBehavior creates a new CRUD behavior for the given table.
// Panics if table is nil.
func NewCRUDBehavior[T any](themeObj theme.Theme, table *TableBehavior[T]) *CRUDBehavior[T] {
	if table == nil {
		panic("table cannot be nil")
	}

	crud := &CRUDBehavior[T]{
		table:  table,
		mode:   ModeList,
		width:  80,
		height: 24,
	}

	crud.SetTheme(themeObj)
	return crud
}

// Mode returns the current CRUD mode.
func (c *CRUDBehavior[T]) Mode() CRUDMode {
	return c.mode
}

// IsInListMode returns true if in list mode.
func (c *CRUDBehavior[T]) IsInListMode() bool {
	return c.mode == ModeList
}

// ReturnToList returns to list mode.
func (c *CRUDBehavior[T]) ReturnToList() {
	c.mode = ModeList
	c.confirmModal = nil
	c.itemToDelete = nil
}

// ItemNamer sets the function to get an item's name for display.
// Used in delete confirmation messages.
func (c *CRUDBehavior[T]) ItemNamer(namer func(T) string) *CRUDBehavior[T] {
	c.itemNamer = namer
	return c
}

// OnCreate sets the create callback and enables create operations.
func (c *CRUDBehavior[T]) OnCreate(callback func() tea.Cmd) *CRUDBehavior[T] {
	c.onCreate = callback
	c.enableCreate = callback != nil
	return c
}

// OnEdit sets the edit callback and enables edit operations.
func (c *CRUDBehavior[T]) OnEdit(callback func(T) tea.Cmd) *CRUDBehavior[T] {
	c.onEdit = callback
	c.enableEdit = callback != nil
	return c
}

// OnDelete sets the delete callback and enables delete operations.
func (c *CRUDBehavior[T]) OnDelete(callback func(T) tea.Cmd) *CRUDBehavior[T] {
	c.onDelete = callback
	c.enableDelete = callback != nil
	return c
}

// SetDimensions sets the width and height for modal rendering.
func (c *CRUDBehavior[T]) SetDimensions(width, height int) {
	c.width = width
	c.height = height
}

// HandleKey processes CRUD operation keys (n/e/d).
// Returns (cmd, true) if key was handled, (nil, false) otherwise.
func (c *CRUDBehavior[T]) HandleKey(key string) (tea.Cmd, bool) {
	switch key {
	case "n":
		return c.handleCreate()
	case "e":
		return c.handleEdit()
	case "d":
		return c.handleDelete()
	default:
		return nil, false
	}
}

// handleCreate processes the create key.
func (c *CRUDBehavior[T]) handleCreate() (tea.Cmd, bool) {
	if !c.enableCreate {
		return nil, false
	}

	c.mode = ModeCreate
	cmd := c.onCreate()
	return cmd, true
}

// handleEdit processes the edit key.
func (c *CRUDBehavior[T]) handleEdit() (tea.Cmd, bool) {
	if !c.enableEdit {
		return nil, false
	}

	// Check if an item is selected
	if c.table.IsEmpty() {
		return nil, false
	}

	item := c.table.GetSelectedItem()
	if item == nil {
		return nil, false
	}

	c.mode = ModeEdit
	cmd := c.onEdit(*item)
	return cmd, true
}

// handleDelete processes the delete key.
func (c *CRUDBehavior[T]) handleDelete() (tea.Cmd, bool) {
	if !c.enableDelete {
		return nil, false
	}

	// Check if an item is selected
	if c.table.IsEmpty() {
		return nil, false
	}

	item := c.table.GetSelectedItem()
	if item == nil {
		return nil, false
	}

	// Store the item to delete
	c.itemToDelete = item

	// Create delete confirmation modal
	itemName := c.getItemName(*item)
	message := fmt.Sprintf("Are you sure you want to delete %s?\n\nThis action cannot be undone.", itemName)

	c.confirmModal = components.NewWarningModal("Delete Confirmation", message)
	c.mode = ModeDelete

	return nil, true
}

// getItemName returns the item name using ItemNamer if set, otherwise "this item".
func (c *CRUDBehavior[T]) getItemName(item T) string {
	if c.itemNamer != nil {
		return fmt.Sprintf("'%s'", c.itemNamer(item))
	}
	return "this item"
}

// Update processes messages when in delete confirmation mode.
// Returns (cmd, true) if message was handled, (nil, false) otherwise.
func (c *CRUDBehavior[T]) Update(msg tea.Msg) (tea.Cmd, bool) {
	if c.mode != ModeDelete {
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			// User confirmed delete
			if c.itemToDelete != nil && c.onDelete != nil {
				cmd := c.onDelete(*c.itemToDelete)
				c.ReturnToList()
				return cmd, true
			}
			c.ReturnToList()
			return nil, true

		case "n", "N", "esc":
			// User canceled delete
			c.ReturnToList()
			return nil, true
		}
	}

	return nil, false
}

// RenderDeleteConfirm renders the delete confirmation modal.
func (c *CRUDBehavior[T]) RenderDeleteConfirm() string {
	if c.confirmModal == nil {
		return ""
	}

	return c.confirmModal.Render(c.width, c.height)
}

// RenderHelpKeys renders help text for enabled CRUD operations.
// Returns a string like "n New • e Edit • d Delete" for enabled operations.
func (c *CRUDBehavior[T]) RenderHelpKeys() string {
	var keys []string

	if c.enableCreate {
		keys = append(keys, "n New")
	}
	if c.enableEdit {
		keys = append(keys, "e Edit")
	}
	if c.enableDelete {
		keys = append(keys, "d Delete")
	}

	if len(keys) == 0 {
		return ""
	}

	result := ""
	for i, key := range keys {
		if i > 0 {
			result += " • "
		}
		result += key
	}

	return result
}
