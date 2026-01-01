package navigation

import (
	tea "github.com/charmbracelet/bubbletea"
)

// KeyHandlers provides factory methods for creating standard key handlers
// This is the central point of truth for key handling across the application
type KeyHandlers struct{}

// List returns a key handler for list views
func (kh *KeyHandlers) List() KeyHandler {
	return NewListKeyHandler()
}

// Form returns a key handler for form views
func (kh *KeyHandlers) Form() KeyHandler {
	return NewFormKeyHandler()
}

// Dialog returns a key handler for dialog/modal views
func (kh *KeyHandlers) Dialog() KeyHandler {
	return NewDialogKeyHandler()
}

// View returns a key handler for read-only view screens
func (kh *KeyHandlers) View() KeyHandler {
	return NewViewKeyHandler()
}

// Menu returns a key handler for menu selections
func (kh *KeyHandlers) Menu() KeyHandler {
	return NewMenuKeyHandler()
}

// NewKeyHandlers creates a new KeyHandlers instance
func NewKeyHandlers() *KeyHandlers {
	return &KeyHandlers{}
}

// GlobalKeyHandler handles keys that are relevant across all views
// These are typically quit, help, and home navigation
type GlobalKeyHandler struct{}

// NewGlobalKeyHandler creates a new global key handler
func NewGlobalKeyHandler() *GlobalKeyHandler {
	return &GlobalKeyHandler{}
}

// HandleKey processes global key presses
func (h *GlobalKeyHandler) HandleKey(keyMsg tea.KeyMsg) KeyAction {
	keyStr := keyMsg.String()

	switch keyStr {
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
	case "h":
		return KeyAction{
			NavigationKey: ptrNavigationKey(KeyHome),
			IsHandled:     true,
		}
	default:
		return KeyAction{IsHandled: false}
	}
}

// IsGlobalKey checks if a key is a global key that should be handled by the app
func IsGlobalKey(keyMsg tea.KeyMsg) bool {
	return NewGlobalKeyHandler().HandleKey(keyMsg).IsHandled
}

// IsNavigationKey checks if the key corresponds to a navigation constant
func IsNavigationKey(key string, navKey NavigationKey) bool {
	switch navKey {
	case KeyUp:
		return key == "up" || key == "k"
	case KeyDown:
		return key == "down" || key == "j"
	case KeyLeft:
		return key == "left" || key == "h"
	case KeyRight:
		return key == "right" || key == "l"
	case KeyBack:
		return key == "esc"
	case KeySelect:
		return key == "enter"
	case KeyToggle:
		return key == "space"
	case KeyFilter:
		return key == "f"
	case KeySort:
		return key == "s"
	case KeySearch:
		return key == "/"
	case KeyEdit:
		return key == "e"
	case KeyDelete:
		return key == "d"
	case KeyBulk:
		return key == "b"
	case KeyCapture:
		return key == "c"
	case KeyList:
		return key == "l"
	case KeyMetadata:
		return key == "m"
	case KeyHelp:
		return key == "?"
	case KeyHome:
		return key == "h"
	case KeyQuit:
		return key == "q" || key == "ctrl+c"
	case KeyPending:
		return key == "p"
	default:
		return false
	}
}

// GetNavigationKeyFromString maps a key string to its corresponding NavigationKey
// Returns nil if the key doesn't correspond to any NavigationKey
func GetNavigationKeyFromString(key string) *NavigationKey {
	for _, navKey := range AllNavigationKeys() {
		if IsNavigationKey(key, navKey) {
			return &navKey
		}
	}
	return nil
}

// MatchesNavigationKey checks if a key message matches a navigation key
func MatchesNavigationKey(keyMsg tea.KeyMsg, navKey NavigationKey) bool {
	return IsNavigationKey(keyMsg.String(), navKey)
}

