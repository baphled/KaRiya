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

// BaseSelectScreen[T] provides a reusable list selection screen.
//
// This generic screen handles:
// - List navigation (↑/↓/j/k/g/G)
// - Item selection (Enter)
// - Cancellation (Esc)
// - Large list handling (scrolling)
// - Empty list handling
//
// Type parameter T can be any type (string, int, struct, pointer to struct, etc.)
//
// Example usage:
//
//	type CVProfile struct {
//	    Name string
//	    Role string
//	}
//
//	profiles := []*CVProfile{...}
//	renderer := func(p *CVProfile) string {
//	    return fmt.Sprintf("%s (%s)", p.Name, p.Role)
//	}
//
//	screen := base.NewBaseSelectScreen[*CVProfile](
//	    profiles,
//	    renderer,
//	    []string{"Main Menu", "Generate CV", "Select Profile"},
//	    "Select CV Profile",
//	)
//
//	// In intent's Update:
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    selectedProfile := result.Data().(*CVProfile)
//	    // ... use selected profile
//	}
//
// Related:
// - internal/cli/screens/contract.go (Screen interface, ScreenResult types)
// - internal/cli/screens/base/base_screen.go (BaseScreen)
type BaseSelectScreen[T any] struct {
	*BaseScreen

	// items is the list of items to select from
	items []T

	// renderer converts items to display strings
	renderer ItemRenderer[T]

	// breadcrumbs for the view header
	breadcrumbs []string

	// title for the view (e.g., "Select CV Profile")
	title string

	// selectedIndex is the currently selected item index
	selectedIndex int

	// scrollOffset for large lists (top visible item index)
	scrollOffset int

	// visibleItems is how many items can be shown at once
	visibleItems int
}

// NewBaseSelectScreen creates a new BaseSelectScreen with the given items.
//
// Parameters:
//   - items: The list of items to select from
//   - renderer: Function to convert items to display strings
//   - breadcrumbs: Breadcrumb trail for header (e.g., ["Main Menu", "Select Profile"])
//   - title: Title for the screen (e.g., "Select CV Profile")
//
// Returns a BaseSelectScreen with selection at index 0 (if items exist).
func NewBaseSelectScreen[T any](
	items []T,
	renderer ItemRenderer[T],
	breadcrumbs []string,
	title string,
) *BaseSelectScreen[T] {
	return &BaseSelectScreen[T]{
		BaseScreen:    NewBaseScreen(),
		items:         items,
		renderer:      renderer,
		breadcrumbs:   breadcrumbs,
		title:         title,
		selectedIndex: 0,
		scrollOffset:  0,
		visibleItems:  10, // Default to showing 10 items
	}
}

// WithInitialSelection sets the initial selection index.
//
// This is useful for restoring state when navigating back.
// If the index is out of bounds, it will be clamped to valid range.
func (s *BaseSelectScreen[T]) WithInitialSelection(index int) *BaseSelectScreen[T] {
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
// Handles:
// - Navigation keys (↑/↓/j/k/g/G)
// - Selection (Enter) → returns NavigateResult with selected item
// - Cancellation (Esc) → returns CancelResult
// - Window resize → updates dimensions
func (s *BaseSelectScreen[T]) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	// Handle window size via BaseScreen
	if cmd := s.BaseScreen.HandleWindowSizeMsg(msg); cmd != nil {
		s.updateVisibleItems()
		return cmd, nil
	}

	// Handle key messages
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		// Navigation - Down
		case "down", "j":
			s.navigateDown()
			return nil, nil

		// Navigation - Up
		case "up", "k":
			s.navigateUp()
			return nil, nil

		// Navigation - Jump to top
		case "g":
			s.jumpToTop()
			return nil, nil

		// Navigation - Jump to bottom
		case "G":
			s.jumpToBottom()
			return nil, nil

		// Selection
		case "enter":
			return nil, s.handleSelection()

		// Cancellation
		case "esc":
			return nil, s.handleCancellation()
		}
	}

	return nil, nil
}

// View implements the Screen interface.
//
// Renders the list with:
// - StandardView layout (logo, breadcrumbs, footer)
// - Selection indicator (▶) on current item
// - Scroll indicator if list is larger than visible area
func (s *BaseSelectScreen[T]) View() string {
	content := s.RenderContent()
	footer := s.RenderFooter()

	return s.CreateView(s.breadcrumbs, content, footer)
}

// navigateDown moves selection down by one item.
func (s *BaseSelectScreen[T]) navigateDown() {
	if len(s.items) == 0 {
		return
	}

	if s.selectedIndex < len(s.items)-1 {
		s.selectedIndex++
		s.updateScrollOffset()
	}
}

// navigateUp moves selection up by one item.
func (s *BaseSelectScreen[T]) navigateUp() {
	if len(s.items) == 0 {
		return
	}

	if s.selectedIndex > 0 {
		s.selectedIndex--
		s.updateScrollOffset()
	}
}

// jumpToTop moves selection to the first item.
func (s *BaseSelectScreen[T]) jumpToTop() {
	if len(s.items) == 0 {
		return
	}

	s.selectedIndex = 0
	s.scrollOffset = 0
}

// jumpToBottom moves selection to the last item.
func (s *BaseSelectScreen[T]) jumpToBottom() {
	if len(s.items) == 0 {
		return
	}

	s.selectedIndex = len(s.items) - 1
	s.updateScrollOffset()
}

// updateScrollOffset adjusts scroll position to keep selection visible.
func (s *BaseSelectScreen[T]) updateScrollOffset() {
	// If selected item is above visible window, scroll up
	if s.selectedIndex < s.scrollOffset {
		s.scrollOffset = s.selectedIndex
	}

	// If selected item is below visible window, scroll down
	if s.selectedIndex >= s.scrollOffset+s.visibleItems {
		s.scrollOffset = s.selectedIndex - s.visibleItems + 1
	}
}

// updateVisibleItems recalculates how many items can be shown based on terminal height.
func (s *BaseSelectScreen[T]) updateVisibleItems() {
	// Reserve space for logo, header, footer, spacing
	// Rough estimate: 10 lines for header/footer, rest for content
	availableHeight := s.Height() - 15
	if availableHeight < 5 {
		availableHeight = 5
	}
	s.visibleItems = availableHeight
}

// handleSelection returns a NavigateResult with the selected item.
func (s *BaseSelectScreen[T]) handleSelection() screens.ScreenResult {
	// Handle empty list
	if len(s.items) == 0 {
		return nil
	}

	selectedItem := s.items[s.selectedIndex]

	result := &screens.NavigateResult{
		ResultData: selectedItem,
	}

	// Store selection index in metadata for state preservation
	result.WithMetadata("selected_index", s.selectedIndex)

	return result
}

// handleCancellation returns a CancelResult with current state in metadata.
func (s *BaseSelectScreen[T]) handleCancellation() screens.ScreenResult {
	result := &screens.CancelResult{}

	// Store current index so it can be restored if user comes back
	result.WithMetadata("selected_index", s.selectedIndex)

	return result
}

// RenderContent renders the list of items with selection indicator.
// This is public so intents can get raw content for custom layouts.
func (s *BaseSelectScreen[T]) RenderContent() string {
	// Handle empty list
	if len(s.items) == 0 {
		return "\n  No items available\n"
	}

	var b strings.Builder

	// Add title if provided
	if s.title != "" {
		b.WriteString("\n  ")
		b.WriteString(s.title)
		b.WriteString("\n\n")
	}

	// Calculate visible range
	start := s.scrollOffset
	end := s.scrollOffset + s.visibleItems
	if end > len(s.items) {
		end = len(s.items)
	}

	// Render visible items
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

	// Add scroll indicators if needed
	if s.scrollOffset > 0 {
		b.WriteString("\n  ↑ More items above")
	}
	if end < len(s.items) {
		b.WriteString("\n  ↓ More items below")
	}

	return b.String()
}

// RenderFooter renders footer with navigation hints.
// This is public so intents can get raw footer for custom layouts.
func (s *BaseSelectScreen[T]) RenderFooter() string {
	if len(s.items) == 0 {
		return "Esc: Back  q: Quit"
	}

	return "↑/↓/j/k: Navigate  g/G: Jump  Enter: Select  Esc: Back  q: Quit"
}
