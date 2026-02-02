package base

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// ItemRenderer is a function that converts an item of type T to a string for display.
//
// Example:
//
//	func renderProfile(p *CVProfile) string {
//	    return fmt.Sprintf("%s (%s)", p.Name, p.TargetRole)
//	}
type ItemRenderer[T any] func(T) string

// SelectScreen provides a reusable list selection screen for any item type.
//
// The generic type parameter T is the element type displayed in the list
// (e.g., string, *CVProfile, or any domain struct). T must satisfy the
// "any" constraint. An ItemRenderer[T] function converts each T value to a
// display string shown as one row in the list.
//
// The screen handles list navigation (arrow keys and j/k for single-step
// movement, g/G for jump-to-top and jump-to-bottom), item selection via
// Enter, cancellation via Esc, and automatic scrolling when the list is
// larger than the visible area.
//
// SelectScreen communicates results back to the parent intent through the
// ScreenResult interface returned from Update. On Enter it returns a
// NavigateResult whose Data field holds the selected T value and whose
// metadata map includes the "selected_index" key. On Esc it returns a
// CancelResult with "selected_index" metadata preserving the cursor
// position. The parent intent inspects these results in its Update loop to
// drive state transitions.
//
// Example usage:
//
//	screen := base.NewBaseSelectScreen[*CVProfile](
//	    profiles,
//	    renderer,
//	    []string{"Main Menu", "Generate CV", "Select Profile"},
//	    "Select CV Profile",
//	)
//
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    selectedProfile := result.Data().(*CVProfile)
//	}
type SelectScreen[T any] struct {
	*Screen

	items         []T
	renderer      ItemRenderer[T]
	breadcrumbs   []string
	title         string
	selectedIndex int
	scrollOffset  int
	visibleItems  int
}

// NewBaseSelectScreen creates a new SelectScreen with the given items.
//
// Parameters:
//   - items: The list of items to select from
//   - renderer: Function to convert items to display strings
//   - breadcrumbs: Breadcrumb trail for header (e.g., ["Main Menu", "Select Profile"])
//   - title: Title for the screen (e.g., "Select CV Profile")
//
// Expected:
//   - items must be a valid slice of T.
//   - renderer must be a valid ItemRenderer function.
//   - breadcrumbs must be a valid slice of strings.
//   - title must be a valid string.
//
// Returns:
//   - A fully initialized SelectScreen[T] ready for use with selection at index 0.
//
// Side effects:
//   - None.
func NewBaseSelectScreen[T any](
	items []T,
	renderer ItemRenderer[T],
	breadcrumbs []string,
	title string,
) *SelectScreen[T] {
	return &SelectScreen[T]{
		Screen:        NewBaseScreen(),
		items:         items,
		renderer:      renderer,
		breadcrumbs:   breadcrumbs,
		title:         title,
		selectedIndex: 0,
		scrollOffset:  0,
		visibleItems:  10,
	}
}

// WithInitialSelection sets the initial selection index.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized SelectScreen[T] ready for use.
//
// Side effects:
//   - None.
func (s *SelectScreen[T]) WithInitialSelection(index int) *SelectScreen[T] {
	if index < 0 {
		index = 0
	}
	if index >= len(s.items) {
		index = len(s.items) - 1
	}
	if index < 0 {
		index = 0
	}
	s.selectedIndex = index
	s.updateScrollOffset()
	return s
}

// Update implements the Screen interface.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - screens.ScreenResult: result indicating selection or cancellation.
//
// Side effects:
//   - May update selection index.
//   - May update scroll offset.
func (s *SelectScreen[T]) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	if cmd := s.HandleWindowSizeMsg(msg); cmd != nil {
		s.updateVisibleItems()
		return cmd, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "down", "j":
			s.navigateDown()
			return nil, nil

		case "up", "k":
			s.navigateUp()
			return nil, nil

		case "g":
			s.jumpToTop()
			return nil, nil

		case "G":
			s.jumpToBottom()
			return nil, nil

		case "enter":
			return nil, s.handleSelection()

		case "esc":
			return nil, s.handleCancellation()
		}
	}

	return nil, nil
}

// View implements the Screen interface.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *SelectScreen[T]) View() string {
	content := s.RenderContent()
	footer := s.RenderFooter()

	return s.CreateView(s.breadcrumbs, content, footer)
}

// navigateDown moves selection down by one item.
func (s *SelectScreen[T]) navigateDown() {
	if len(s.items) == 0 {
		return
	}

	if s.selectedIndex < len(s.items)-1 {
		s.selectedIndex++
		s.updateScrollOffset()
	}
}

// navigateUp moves selection up by one item.
func (s *SelectScreen[T]) navigateUp() {
	if len(s.items) == 0 {
		return
	}

	if s.selectedIndex > 0 {
		s.selectedIndex--
		s.updateScrollOffset()
	}
}

// jumpToTop moves selection to the first item.
func (s *SelectScreen[T]) jumpToTop() {
	if len(s.items) == 0 {
		return
	}

	s.selectedIndex = 0
	s.scrollOffset = 0
}

// jumpToBottom moves selection to the last item.
func (s *SelectScreen[T]) jumpToBottom() {
	if len(s.items) == 0 {
		return
	}

	s.selectedIndex = len(s.items) - 1
	s.updateScrollOffset()
}

// updateScrollOffset adjusts scroll position to keep selection visible.
func (s *SelectScreen[T]) updateScrollOffset() {
	if s.selectedIndex < s.scrollOffset {
		s.scrollOffset = s.selectedIndex
	}

	if s.selectedIndex >= s.scrollOffset+s.visibleItems {
		s.scrollOffset = s.selectedIndex - s.visibleItems + 1
	}
}

// updateVisibleItems recalculates how many items can be shown based on terminal height.
func (s *SelectScreen[T]) updateVisibleItems() {
	availableHeight := s.Height() - 15
	if availableHeight < 5 {
		availableHeight = 5
	}
	s.visibleItems = availableHeight
}

// handleSelection returns a NavigateResult with the selected item.
func (s *SelectScreen[T]) handleSelection() screens.ScreenResult {
	if len(s.items) == 0 {
		return nil
	}

	selectedItem := s.items[s.selectedIndex]

	result := &screens.NavigateResult{
		ResultData: selectedItem,
	}

	result.WithMetadata("selected_index", s.selectedIndex)

	return result
}

// handleCancellation returns a CancelResult with current state in metadata.
func (s *SelectScreen[T]) handleCancellation() screens.ScreenResult {
	result := &screens.CancelResult{}

	result.WithMetadata("selected_index", s.selectedIndex)

	return result
}

// RenderContent renders the list of items with selection indicator.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *SelectScreen[T]) RenderContent() string {
	if len(s.items) == 0 {
		return "\n  No items available\n"
	}

	var b strings.Builder

	if s.title != "" {
		b.WriteString("\n  ")
		b.WriteString(s.title)
		b.WriteString("\n\n")
	}

	start := s.scrollOffset
	end := s.scrollOffset + s.visibleItems
	if end > len(s.items) {
		end = len(s.items)
	}

	for i := start; i < end; i++ {
		prefix := "  "
		if i == s.selectedIndex {
			prefix = "▶ "
		}

		itemText := s.renderer(s.items[i])
		b.WriteString(prefix)
		b.WriteString(itemText)
		b.WriteString("\n")
	}

	if s.scrollOffset > 0 {
		b.WriteString("\n  ↑ More items above")
	}
	if end < len(s.items) {
		b.WriteString("\n  ↓ More items below")
	}

	return b.String()
}

// RenderFooter renders footer with navigation hints.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *SelectScreen[T]) RenderFooter() string {
	if len(s.items) == 0 {
		return "Esc: Back  q: Quit"
	}

	return "↑/↓/j/k: Navigate  g/G: Jump  Enter: Select  Esc: Back  q: Quit"
}
