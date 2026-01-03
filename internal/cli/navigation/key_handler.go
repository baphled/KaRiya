package navigation

import (
	tea "github.com/charmbracelet/bubbletea"
)

// KeyAction represents the standardized action triggered by a key press
type KeyAction struct {
	// The navigation key that was triggered (e.g., KeyUp, KeyDown, KeySelect)
	// If nil, the key was not recognized
	NavigationKey *NavigationKey
	// Custom action type for non-standard keys
	ActionType string
	// Whether this key was handled by the handler
	IsHandled bool
	// Optional metadata about the action
	Metadata map[string]interface{}
}

// KeyHandler provides methods to handle key messages in a standardized way
// Different handler types exist for different view contexts
type KeyHandler interface {
	// HandleKey processes a key message and returns a standardized action
	HandleKey(keyMsg tea.KeyMsg) KeyAction
}

// ListKeyHandler handles key presses for list-like views
// Lists support: navigation (up/down), pagination, first/last, selection
type ListKeyHandler struct {
	allowVimKeys bool
}

// NewListKeyHandler creates a new list key handler
func NewListKeyHandler() *ListKeyHandler {
	return &ListKeyHandler{
		allowVimKeys: true, // Enable vim-like keys (j/k, h/l) by default
	}
}

// HandleKey processes key messages for list views
// Recognizes standard list navigation patterns
func (h *ListKeyHandler) HandleKey(keyMsg tea.KeyMsg) KeyAction {
	keyStr := keyMsg.String()

	switch keyStr {
	// Up navigation
	case "up", "k", "shift+tab":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyUp),
			IsHandled:     true,
		}
	// Down navigation
	case "down", "j", "tab":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyDown),
			IsHandled:     true,
		}
	// Left navigation (for horizontal lists)
	case "left", "h":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyLeft),
			IsHandled:     true,
		}
	// Right navigation (for horizontal lists)
	case "right", "l":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyRight),
			IsHandled:     true,
		}
	// First item
	case "home", "g":
		return KeyAction{
			ActionType: "navigate:first",
			IsHandled:  true,
		}
	// Last item
	case "end", "G", "shift+g":
		return KeyAction{
			ActionType: "navigate:last",
			IsHandled:  true,
		}
	// Page up
	case "pgup", "ctrl+b", "ctrl+u":
		return KeyAction{
			ActionType: "navigate:page_up",
			IsHandled:  true,
		}
	// Page down
	case "pgdn", "ctrl+f", "ctrl+d":
		return KeyAction{
			ActionType: "navigate:page_down",
			IsHandled:  true,
		}
	// Select/Enter
	case "enter":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeySelect),
			IsHandled:     true,
		}
	// Toggle selection (for multi-select lists)
	case "space":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyToggle),
			IsHandled:     true,
		}
	// Standard action keys
	case "/":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeySearch),
			IsHandled:     true,
		}
	case "f":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyFilter),
			IsHandled:     true,
		}
	case "s":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeySort),
			IsHandled:     true,
		}
	case "e":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyEdit),
			IsHandled:     true,
		}
	case "d":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyDelete),
			IsHandled:     true,
		}
	// Global actions (handled but app may process)
	case "q", "ctrl+c":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyQuit),
			IsHandled:     true,
		}
	case "esc":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyBack),
			IsHandled:     true,
		}
	case "?":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyHelp),
			IsHandled:     true,
		}
	default:
		return KeyAction{IsHandled: false}
	}
}

// FormKeyHandler handles key presses for form-like views
// Forms support: field navigation, submission, cancellation
type FormKeyHandler struct {
	allowVimKeys bool
}

// NewFormKeyHandler creates a new form key handler
func NewFormKeyHandler() *FormKeyHandler {
	return &FormKeyHandler{
		allowVimKeys: true, // Enable vim-like keys for form navigation
	}
}

// HandleKey processes key messages for form views
// Recognizes form-specific navigation patterns
func (h *FormKeyHandler) HandleKey(keyMsg tea.KeyMsg) KeyAction {
	keyStr := keyMsg.String()

	switch keyStr {
	// Field navigation - Tab/Shift+Tab
	case "tab":
		return KeyAction{
			ActionType: "field:next",
			IsHandled:  true,
		}
	case "shift+tab":
		return KeyAction{
			ActionType: "field:previous",
			IsHandled:  true,
		}
	// Vim-style field navigation (down/up move between fields)
	case "down", "j":
		return KeyAction{
			ActionType: "field:next",
			IsHandled:  true,
		}
	case "up", "k":
		return KeyAction{
			ActionType: "field:previous",
			IsHandled:  true,
		}
	// Submit form
	case "ctrl+s":
		return KeyAction{
			ActionType: "form:submit",
			IsHandled:  true,
		}
	// Cancel form
	case "esc":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyBack),
			IsHandled:     true,
		}
	// Global actions
	case "q", "ctrl+c":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyQuit),
			IsHandled:     true,
		}
	case "?":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyHelp),
			IsHandled:     true,
		}
	default:
		return KeyAction{IsHandled: false}
	}
}

// DialogKeyHandler handles key presses for dialog/modal views
// Dialogs support: yes/no confirmation, cancellation
type DialogKeyHandler struct{}

// NewDialogKeyHandler creates a new dialog key handler
func NewDialogKeyHandler() *DialogKeyHandler {
	return &DialogKeyHandler{}
}

// HandleKey processes key messages for dialog views
// Recognizes dialog-specific confirmation patterns
func (h *DialogKeyHandler) HandleKey(keyMsg tea.KeyMsg) KeyAction {
	keyStr := keyMsg.String()

	switch keyStr {
	// Button switching (left/right for confirmation dialogs)
	case "left", "h", "shift+tab":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyLeft),
			IsHandled:     true,
		}
	case "right", "l", "tab":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyRight),
			IsHandled:     true,
		}
	// Confirmation
	case "y":
		return KeyAction{
			ActionType: "dialog:yes",
			IsHandled:  true,
		}
	case "n":
		return KeyAction{
			ActionType: "dialog:no",
			IsHandled:  true,
		}
	// Alternative confirmation
	case "enter":
		return KeyAction{
			ActionType: "dialog:ok",
			IsHandled:  true,
		}
	case "esc":
		return KeyAction{
			ActionType: "dialog:cancel",
			IsHandled:  true,
		}
	// Global actions
	case "q", "ctrl+c":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyQuit),
			IsHandled:     true,
		}
	default:
		return KeyAction{IsHandled: false}
	}
}

// ViewKeyHandler handles key presses for read-only view screens
// Views support: scrolling, navigation, back
type ViewKeyHandler struct{}

// NewViewKeyHandler creates a new view key handler
func NewViewKeyHandler() *ViewKeyHandler {
	return &ViewKeyHandler{}
}

// HandleKey processes key messages for view screens
// Recognizes view-specific scrolling and navigation
func (h *ViewKeyHandler) HandleKey(keyMsg tea.KeyMsg) KeyAction {
	keyStr := keyMsg.String()

	switch keyStr {
	// Scrolling/Navigation - Up
	case "up", "k":
		return KeyAction{
			ActionType: "view:scroll_up",
			IsHandled:  true,
		}
	// Scrolling/Navigation - Down
	case "down", "j":
		return KeyAction{
			ActionType: "view:scroll_down",
			IsHandled:  true,
		}
	// Page navigation
	case "pgup", "ctrl+b", "ctrl+u":
		return KeyAction{
			ActionType: "view:page_up",
			IsHandled:  true,
		}
	case "pgdn", "ctrl+f", "ctrl+d":
		return KeyAction{
			ActionType: "view:page_down",
			IsHandled:  true,
		}
	// First/Last
	case "home", "g":
		return KeyAction{
			ActionType: "view:first",
			IsHandled:  true,
		}
	case "end", "G", "shift+g":
		return KeyAction{
			ActionType: "view:last",
			IsHandled:  true,
		}
	// Navigation
	case "enter":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeySelect),
			IsHandled:     true,
		}
	case "esc":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyBack),
			IsHandled:     true,
		}
	// Global actions
	case "q", "ctrl+c":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyQuit),
			IsHandled:     true,
		}
	case "?":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyHelp),
			IsHandled:     true,
		}
	default:
		return KeyAction{IsHandled: false}
	}
}

// MenuKeyHandler handles key presses for menu selections
// Menus support: navigation, selection
type MenuKeyHandler struct{}

// NewMenuKeyHandler creates a new menu key handler
func NewMenuKeyHandler() *MenuKeyHandler {
	return &MenuKeyHandler{}
}

// HandleKey processes key messages for menu views
// Recognizes menu-specific navigation and selection
func (h *MenuKeyHandler) HandleKey(keyMsg tea.KeyMsg) KeyAction {
	keyStr := keyMsg.String()

	switch keyStr {
	// Navigation
	case "up", "k":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyUp),
			IsHandled:     true,
		}
	case "down", "j":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyDown),
			IsHandled:     true,
		}
	// Selection
	case "enter":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeySelect),
			IsHandled:     true,
		}
	case "esc":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyBack),
			IsHandled:     true,
		}
	// Global actions
	case "q", "ctrl+c":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyQuit),
			IsHandled:     true,
		}
	case "?":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyHelp),
			IsHandled:     true,
		}
	default:
		return KeyAction{IsHandled: false}
	}
}

// ptrNavigationKey is a helper to create a pointer to a NavigationKey
func ptrNavigationKey(k NavigationKey) *NavigationKey {
	return &k
}
