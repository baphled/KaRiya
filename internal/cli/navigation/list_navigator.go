package navigation

// ListNavigator provides minimal callbacks for list navigation.
// Implementers are responsible for managing their own data and table updates.
type ListNavigator interface {
	// GetTotalItems returns the total number of items in the list
	GetTotalItems() int

	// GetSelectedIndex returns the current selection index (0-based, absolute)
	GetSelectedIndex() int

	// SetSelectedIndex sets the selection index and updates display.
	// Implementers should sync table cursor and call updateTableRows() here.
	SetSelectedIndex(idx int)

	// GetPageSize returns items per page (for page up/down navigation)
	GetPageSize() int
}

// ListNavigationHandler handles keyboard navigation for lists.
// It is the single source of truth for navigation operations.
type ListNavigationHandler struct {
	navigator ListNavigator
}

// NewListNavigationHandler creates a new list navigation handler.
//
// Expected:
//   - listnavigator must be valid.
//
// Returns:
//   - A fully initialized ListNavigationHandler ready for use.
//
// Side effects:
//   - None.
func NewListNavigationHandler(navigator ListNavigator) *ListNavigationHandler {
	return &ListNavigationHandler{
		navigator: navigator,
	}
}

// HandleKey processes navigation keys and updates the selection accordingly.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (h *ListNavigationHandler) HandleKey(keyStr string) bool {
	totalItems := h.navigator.GetTotalItems()
	if totalItems == 0 {
		return false
	}

	currentIndex := h.navigator.GetSelectedIndex()
	pageSize := h.navigator.GetPageSize()

	switch keyStr {
	case "up", "k":
		// Move up by 1
		newIndex := currentIndex - 1
		if newIndex < 0 {
			newIndex = 0
		}
		h.navigator.SetSelectedIndex(newIndex)
		return true

	case "down", "j":
		// Move down by 1
		newIndex := currentIndex + 1
		if newIndex >= totalItems {
			newIndex = totalItems - 1
		}
		h.navigator.SetSelectedIndex(newIndex)
		return true

	case "pgup", "ctrl+u":
		// Page up
		newIndex := currentIndex - pageSize
		if newIndex < 0 {
			newIndex = 0
		}
		h.navigator.SetSelectedIndex(newIndex)
		return true

	case "pgdn", "ctrl+d":
		// Page down
		newIndex := currentIndex + pageSize
		if newIndex >= totalItems {
			newIndex = totalItems - 1
		}
		h.navigator.SetSelectedIndex(newIndex)
		return true

	case "home", "g":
		// Go to first item
		h.navigator.SetSelectedIndex(0)
		return true

	case "end", "G":
		// Go to last item
		h.navigator.SetSelectedIndex(totalItems - 1)
		return true

	default:
		return false
	}
}

// FormatRowText returns text with a selection indicator if this is the selected row.
//
// Expected:
//   - int must be valid.
//   - Must be a valid string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (h *ListNavigationHandler) FormatRowText(rowIndex int, text string) string {
	selectedIndex := h.navigator.GetSelectedIndex()
	if rowIndex == selectedIndex {
		return "▶ " + text
	}
	return "  " + text
}
