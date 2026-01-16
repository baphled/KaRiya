package behaviors

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/lipgloss"
)

// SortOption[T] represents a sort option with a label, comparator, and reverse flag.
type SortOption[T any] struct {
	Label      string
	Comparator SortComparator[T]
	Reverse    bool
}

// SortMenuBehavior[T] provides a menu UI for sorting table data.
// It manages navigation, selection state, and applies sort comparators to the table.
//
// Usage:
//
//	sortMenu := behaviors.NewSortMenuBehavior(theme, table, "Sort Items").
//	    AddOption("Name (A-Z)", nameAscComparator, false).
//	    AddOption("Name (Z-A)", nameAscComparator, true).
//	    AddOption("Age (Desc)", ageComparator, true).
//	    OnApply(func() { intent.returnToList() })
//
//	// In Update:
//	if sortMenu.IsActive() && sortMenu.HandleKey(key) {
//	    return nil
//	}
//
//	// In View:
//	if sortMenu.IsActive() {
//	    return sortMenu.Render()
//	}
type SortMenuBehavior[T any] struct {
	theme.Aware

	// References
	table *TableBehavior[T]

	// Configuration
	title   string
	options []SortOption[T]

	// State
	isActive      bool
	focusedIndex  int // Current focus position
	selectedIndex int // Which option is currently applied (-1 if none)

	// Callback
	onApply func()
}

// NewSortMenuBehavior creates a new sort menu behavior.
// Panics if table is nil.
func NewSortMenuBehavior[T any](themeObj theme.Theme, table *TableBehavior[T], title string) *SortMenuBehavior[T] {
	if table == nil {
		panic("table cannot be nil")
	}

	sortMenu := &SortMenuBehavior[T]{
		table:         table,
		title:         title,
		options:       []SortOption[T]{},
		isActive:      false,
		focusedIndex:  0,
		selectedIndex: -1, // No option selected initially
	}

	sortMenu.SetTheme(themeObj)
	return sortMenu
}

// AddOption adds a sort option to the menu.
func (s *SortMenuBehavior[T]) AddOption(label string, comparator SortComparator[T], reverse bool) *SortMenuBehavior[T] {
	s.options = append(s.options, SortOption[T]{
		Label:      label,
		Comparator: comparator,
		Reverse:    reverse,
	})
	return s
}

// OnApply sets the callback to invoke when a sort is applied.
func (s *SortMenuBehavior[T]) OnApply(callback func()) *SortMenuBehavior[T] {
	s.onApply = callback
	return s
}

// Show activates the sort menu.
func (s *SortMenuBehavior[T]) Show() {
	s.isActive = true
}

// Hide deactivates the sort menu.
func (s *SortMenuBehavior[T]) Hide() {
	s.isActive = false
}

// IsActive returns true if the sort menu is currently active.
func (s *SortMenuBehavior[T]) IsActive() bool {
	return s.isActive
}

// HandleKey processes key presses for navigation and selection.
// Returns true if the key was handled.
func (s *SortMenuBehavior[T]) HandleKey(key string) bool {
	if !s.isActive {
		return false
	}

	switch key {
	case "down", "j":
		s.moveFocusDown()
		return true
	case "up", "k":
		s.moveFocusUp()
		return true
	case "enter":
		s.applySort()
		return true
	case "esc":
		s.Hide()
		return true
	default:
		return false
	}
}

// moveFocusDown increments focus index with wrapping.
func (s *SortMenuBehavior[T]) moveFocusDown() {
	if len(s.options) == 0 {
		return
	}

	s.focusedIndex++
	if s.focusedIndex >= len(s.options) {
		s.focusedIndex = 0
	}
}

// moveFocusUp decrements focus index with wrapping.
func (s *SortMenuBehavior[T]) moveFocusUp() {
	if len(s.options) == 0 {
		return
	}

	s.focusedIndex--
	if s.focusedIndex < 0 {
		s.focusedIndex = len(s.options) - 1
	}
}

// applySort applies the focused option's sort to the table.
func (s *SortMenuBehavior[T]) applySort() {
	if s.focusedIndex < 0 || s.focusedIndex >= len(s.options) {
		return
	}

	option := s.options[s.focusedIndex]

	// Apply sort to table
	s.table.SetSort(option.Comparator, option.Reverse)

	// Update selected index
	s.selectedIndex = s.focusedIndex

	// Invoke callback
	if s.onApply != nil {
		s.onApply()
	}

	// Hide menu
	s.Hide()
}

// Render returns the themed menu view.
func (s *SortMenuBehavior[T]) Render() string {
	if !s.isActive {
		return ""
	}

	// Use Catppuccin Macchiato colors for consistent theming
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CBA6F7")). // Catppuccin Mauve
		Bold(true).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A6E3A1")). // Catppuccin Green
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CDD6F4")) // Catppuccin Text

	var lines []string
	lines = append(lines, titleStyle.Render(s.title))
	lines = append(lines, "")

	for idx, option := range s.options {
		marker := "  "
		if idx == s.focusedIndex {
			marker = "▶ "
		}

		text := option.Label

		// Show checkmark if this option is currently selected
		if idx == s.selectedIndex {
			text += " ✓"
		}

		style := normalStyle
		if idx == s.focusedIndex {
			style = selectedStyle
		}

		lines = append(lines, style.Render(marker+text))
	}

	return strings.Join(lines, "\n")
}
