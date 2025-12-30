package navigation

// NavigationKey represents standardized keyboard shortcuts used throughout the TUI
type NavigationKey string

const (
	// Primary navigation keys
	KeyBack   NavigationKey = "Esc"    // Go back to previous screen
	KeyUp     NavigationKey = "↑/k"    // Navigate up in lists, menus, or forms
	KeyDown   NavigationKey = "↓/j"    // Navigate down in lists, menus, or forms
	KeyLeft   NavigationKey = "←/h"    // Navigate left in forms or menus
	KeyRight  NavigationKey = "→/l"    // Navigate right in forms or menus
	KeySelect NavigationKey = "Enter"  // Confirm selection or submit form
	KeyToggle NavigationKey = "Space"  // Toggle checkbox, tag, or item selection

	// Action keys
	KeyFilter   NavigationKey = "f"  // Show/toggle filters
	KeySort     NavigationKey = "s"  // Show/toggle sort options
	KeySearch   NavigationKey = "/" // Show/toggle search
	KeyEdit     NavigationKey = "e"  // Edit selected item
	KeyDelete   NavigationKey = "d"  // Delete selected item
	KeyBulk     NavigationKey = "b"  // Enter bulk operations mode
	KeyCapture  NavigationKey = "c"  // Capture new event
	KeyList     NavigationKey = "l"  // List events
	KeyMetadata NavigationKey = "m"  // Open metadata review
	KeyHelp     NavigationKey = "?"  // Show help
	KeyHome     NavigationKey = "h"  // Go to home screen
	KeyQuit     NavigationKey = "q"  // Quit application
)

// AllNavigationKeys returns a slice of all defined navigation keys
func AllNavigationKeys() []NavigationKey {
	return []NavigationKey{
		KeyBack,
		KeyUp,
		KeyDown,
		KeyLeft,
		KeyRight,
		KeySelect,
		KeyToggle,
		KeyFilter,
		KeySort,
		KeySearch,
		KeyEdit,
		KeyDelete,
		KeyHelp,
		KeyHome,
		KeyQuit,
		KeyBulk,
		KeyCapture,
		KeyList,
		KeyMetadata,
	}
}

// KeyDescription provides a human-readable description of each navigation key
var KeyDescription = map[NavigationKey]string{
	KeyBack:     "Go back to previous screen",
	KeyUp:       "Navigate up (also j)",
	KeyDown:     "Navigate down (also k)",
	KeyLeft:     "Navigate left (also h)",
	KeyRight:    "Navigate right (also l)",
	KeySelect:   "Confirm selection or submit",
	KeyToggle:   "Toggle selection or expansion",
	KeyFilter:   "Show or toggle filters",
	KeySort:     "Show or toggle sort options",
	KeySearch:   "Show search functionality",
	KeyEdit:     "Edit selected item",
	KeyDelete:   "Delete selected item",
	KeyHelp:     "Show help information",
	KeyHome:     "Go to home screen",
	KeyQuit:     "Quit application",
	KeyBulk:     "Enter bulk operations mode",
	KeyCapture:  "Capture new event",
	KeyList:     "List all events",
	KeyMetadata: "Open metadata review",
}

