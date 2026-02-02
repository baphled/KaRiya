package navigation

// NavigationKey maps a user-facing keyboard shortcut label to the action it
// triggers throughout the TUI. Components display the string value in help
// bars and badges so the user knows which key to press.
//
//nolint:revive // "NavigationKey" name is intentional for clarity in external packages.
type NavigationKey string

const (
	// KeyBack dismisses the current screen and returns to its parent, such as
	// closing a modal or navigating from a detail view back to a list.
	KeyBack NavigationKey = "Esc"
	// KeyUp moves the cursor or highlight one row upward in a list, menu, or
	// form field group.
	KeyUp NavigationKey = "↑/k"
	// KeyDown moves the cursor or highlight one row downward in a list, menu,
	// or form field group.
	KeyDown NavigationKey = "↓/j"
	// KeyLeft moves focus one column or tab to the left in multi-column
	// layouts and form navigation.
	KeyLeft NavigationKey = "←/h"
	// KeyRight moves focus one column or tab to the right in multi-column
	// layouts and form navigation.
	KeyRight NavigationKey = "→/l"
	// KeySelect confirms the currently highlighted item or submits the active
	// form, advancing the intent to the next state.
	KeySelect NavigationKey = "Enter"
	// KeyToggle flips the checked state of a checkbox, tag, or multi-select
	// option without advancing the cursor.
	KeyToggle NavigationKey = "Space"

	// KeyAdd opens the creation form for a new item in the current context,
	// such as a new career event or skill entry.
	KeyAdd NavigationKey = "a"
	// KeyFilter opens or toggles the filter modal, allowing the user to
	// narrow the visible items by tag, date, or category.
	KeyFilter NavigationKey = "f"
	// KeySort opens or cycles the sort modal, allowing the user to reorder
	// items by date, name, or relevance.
	KeySort NavigationKey = "s"
	// KeySearch activates the search input overlay, enabling incremental
	// text matching against the current item list.
	KeySearch NavigationKey = "/"
	// KeyEdit opens the edit form for the currently selected item, loading
	// its existing values into the form fields.
	KeyEdit NavigationKey = "e"
	// KeyDelete initiates deletion of the currently selected item, typically
	// showing a confirmation modal before removing it.
	KeyDelete NavigationKey = "d"
	// KeyBulk enters bulk operations mode, enabling multi-select actions
	// such as batch tagging or batch deletion.
	KeyBulk NavigationKey = "b"
	// KeyCapture opens the event capture form for recording a new career
	// timeline entry with text, date, and metadata.
	KeyCapture NavigationKey = "c"
	// KeyList switches to the event list view, displaying all career
	// timeline entries in a scrollable table.
	KeyList NavigationKey = "l"
	// KeyMetadata opens the metadata review screen where the user can
	// inspect and edit tags, categories, and enrichment data.
	KeyMetadata NavigationKey = "m"
	// KeyHelp toggles the help overlay that displays all available keyboard
	// shortcuts for the current screen.
	KeyHelp NavigationKey = "?"
	// KeyQuit exits the application, returning control to the terminal.
	KeyQuit NavigationKey = "q"
	// KeyPending opens the pending review queue showing items that need
	// user confirmation, such as extracted facts awaiting approval.
	KeyPending NavigationKey = "p"
	// KeyFacts switches to the facts management view listing all curated
	// career facts extracted from timeline events.
	KeyFacts NavigationKey = "t"
	// KeyCV opens the CV configuration manager where the user can create,
	// edit, and select CV generation profiles.
	KeyCV NavigationKey = "v"
	// KeyGenerate triggers CV generation using the currently active
	// configuration profile and launches the output preview.
	KeyGenerate NavigationKey = "g"
)

// AllNavigationKeys returns every defined NavigationKey value. The returned
//
// Returns:
//   - A []NavigationKey value.
//
// Side effects:
//   - None.
func AllNavigationKeys() []NavigationKey {
	return []NavigationKey{
		KeyBack,
		KeyUp,
		KeyDown,
		KeyLeft,
		KeyRight,
		KeySelect,
		KeyToggle,
		KeyAdd,
		KeyFilter,
		KeySort,
		KeySearch,
		KeyEdit,
		KeyDelete,
		KeyHelp,
		KeyQuit,
		KeyBulk,
		KeyCapture,
		KeyList,
		KeyMetadata,
		KeyPending,
		KeyFacts,
		KeyGenerate,
		KeyCV,
	}
}

// KeyDescription maps each NavigationKey to a concise human-readable label
// displayed in help overlays and tooltip badges.
var KeyDescription = map[NavigationKey]string{
	KeyBack:     "Go back to previous screen",
	KeyUp:       "Navigate up (also j)",
	KeyDown:     "Navigate down (also k)",
	KeyLeft:     "Navigate left (also h)",
	KeyRight:    "Navigate right (also l)",
	KeySelect:   "Confirm selection or submit",
	KeyToggle:   "Toggle selection or expansion",
	KeyAdd:      "Add new item",
	KeyFilter:   "Show or toggle filters",
	KeySort:     "Show or toggle sort options",
	KeySearch:   "Show search functionality",
	KeyEdit:     "Edit selected item",
	KeyDelete:   "Delete selected item",
	KeyHelp:     "Show help information",
	KeyQuit:     "Quit application",
	KeyBulk:     "Enter bulk operations mode",
	KeyCapture:  "Capture new event",
	KeyList:     "List all events",
	KeyMetadata: "Open metadata review",
	KeyPending:  "Review pending items",
	KeyFacts:    "View all facts",
	KeyCV:       "Manage CV configurations",
	KeyGenerate: "Generate CV from configurations",
}
