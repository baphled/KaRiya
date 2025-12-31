package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// ShortcutAction defines the action to be executed when a shortcut is triggered
type ShortcutAction func() tea.Cmd

// ShortcutHandler manages keyboard shortcuts for a model
type ShortcutHandler struct {
	shortcuts map[string]key.Binding
	actions   map[string]ShortcutAction
}

// NewShortcutHandler creates a new shortcut handler instance
func NewShortcutHandler() *ShortcutHandler {
	return &ShortcutHandler{
		shortcuts: make(map[string]key.Binding),
		actions:   make(map[string]ShortcutAction),
	}
}

// RegisterShortcut registers a keyboard shortcut with an associated action
func (sh *ShortcutHandler) RegisterShortcut(id string, binding key.Binding, action ShortcutAction) {
	sh.shortcuts[id] = binding
	sh.actions[id] = action
}

// UnregisterShortcut removes a shortcut from the handler
func (sh *ShortcutHandler) UnregisterShortcut(id string) {
	delete(sh.shortcuts, id)
	delete(sh.actions, id)
}

// GetShortcut retrieves a shortcut by ID
func (sh *ShortcutHandler) GetShortcut(id string) (key.Binding, bool) {
	shortcut, exists := sh.shortcuts[id]
	return shortcut, exists
}

// GetAllShortcuts returns all registered shortcuts
func (sh *ShortcutHandler) GetAllShortcuts() map[string]key.Binding {
	return sh.shortcuts
}

// HandleMsg processes a tea message and executes associated shortcut actions
func (sh *ShortcutHandler) HandleMsg(msg tea.Msg) (tea.Cmd, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil, false
	}

	keyStr := keyMsg.String()
	for id, binding := range sh.shortcuts {
		// Get the help key representation which includes the key information
		help := binding.Help()
		if keyStr == help.Key || keyStr == fmt.Sprintf("rune(%q)", help.Key) {
			if action, exists := sh.actions[id]; exists {
				return action(), true
			}
		}
	}

	return nil, false
}

// ClearShortcuts removes all registered shortcuts
func (sh *ShortcutHandler) ClearShortcuts() {
	sh.shortcuts = make(map[string]key.Binding)
	sh.actions = make(map[string]ShortcutAction)
}

// CommonShortcuts provides a standard set of shortcuts used across the application
type CommonShortcuts struct {
	Quit       key.Binding
	Back       key.Binding
	Help       key.Binding
	Enter      key.Binding
	Up         key.Binding
	Down       key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
}

// NewCommonShortcuts creates a standard set of keyboard shortcuts
func NewCommonShortcuts() CommonShortcuts {
	return CommonShortcuts{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Help: key.NewBinding(
			key.WithKeys("?", "h"),
			key.WithHelp("?", "help"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdn"),
			key.WithHelp("pgdn", "page down"),
		),
	}
}

