package base

import (
	"github.com/baphled/kariya/internal/cli/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// ContentRenderer is a function type that renders content for given dimensions.
//
// Parameters:
//   - data: The data structure to render
//   - width: Terminal width for the content
//   - height: Terminal height for the content
//
// The renderer should generate a string representation of the data
// appropriate for the given dimensions.
type ContentRenderer[T any] func(data T, width, height int) string

// BaseDetailScreen provides a reusable screen for displaying detail views.
//
// This screen handles:
// - Content rendering with proper dimensions
// - Scrolling support (up/down, vim keys j/k, g/G for top/bottom)
// - Navigation (escape to go back, enter to confirm)
// - Custom action keys (e.g., 'e' for edit, 'd' for delete)
// - Scroll position preservation in metadata
// - StandardView integration
//
// Type parameter T should be a pointer to your data structure.
//
// Example usage:
//
//	type EventDetail struct {
//	    Title       string
//	    Description string
//	    Date        time.Time
//	}
//
//	renderer := func(data *EventDetail, width, height int) string {
//	    var b strings.Builder
//	    b.WriteString(fmt.Sprintf("Title: %s\n", data.Title))
//	    b.WriteString(fmt.Sprintf("Description: %s\n", data.Description))
//	    b.WriteString(fmt.Sprintf("Date: %s\n", data.Date.Format("2006-01-02")))
//	    return b.String()
//	}
//
//	detail := &EventDetail{...}
//	screen := base.NewBaseDetailScreen(
//	    []string{"Main Menu", "View Event"},
//	    renderer,
//	    detail,
//	)
//	screen.AddAction("e", "edit")  // Optional custom action
//
// Related:
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
// - internal/cli/screens/cv/preview.go (Example detail screen)
type BaseDetailScreen[T any] struct {
	*BaseScreen

	// breadcrumbs for navigation context
	breadcrumbs []string

	// contentRenderer renders the data as a string
	contentRenderer ContentRenderer[T]

	// data is the data structure to display
	data T

	// scrollOffset tracks vertical scroll position
	scrollOffset int

	// footer is the help text shown at the bottom
	footer string

	// actions maps key strings to action names (for custom actions)
	// e.g., "e" -> "edit", "d" -> "delete"
	actions map[string]string
}

// NewBaseDetailScreen creates a new detail screen.
//
// Parameters:
//   - breadcrumbs: Navigation breadcrumb trail (e.g., []string{"Main Menu", "View Item"})
//   - renderer: Function that renders the data as a string for given dimensions
//   - data: The data structure to display
//
// Default behavior:
//   - Escape key goes back (CancelResult)
//   - Enter key confirms (NavigateResult with "confirm")
//   - Up/Down and j/k scroll the content
//   - g/G jump to top/bottom
func NewBaseDetailScreen[T any](
	breadcrumbs []string,
	renderer ContentRenderer[T],
	data T,
) *BaseDetailScreen[T] {
	return &BaseDetailScreen[T]{
		BaseScreen:      NewBaseScreen(),
		breadcrumbs:     breadcrumbs,
		contentRenderer: renderer,
		data:            data,
		scrollOffset:    0,
		footer:          "Enter: Confirm  Esc: Back  ↑↓/jk: Scroll  g/G: Top/Bottom",
		actions:         make(map[string]string),
	}
}

// Update handles messages and returns result when user takes action.
func (s *BaseDetailScreen[T]) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back - preserve scroll position
			result := &screens.CancelResult{}
			result.WithMetadata("scroll_offset", s.scrollOffset)
			return nil, result

		case "enter":
			// Confirm action
			return nil, &screens.NavigateResult{
				ResultData: "confirm",
			}

		case "down", "j":
			// Scroll down
			s.scrollOffset++
			return nil, nil

		case "up", "k":
			// Scroll up (but not below 0)
			if s.scrollOffset > 0 {
				s.scrollOffset--
			}
			return nil, nil

		case "g":
			// Jump to top
			s.scrollOffset = 0
			return nil, nil

		case "G":
			// Jump to bottom (large number, will be capped by view)
			s.scrollOffset = 9999
			return nil, nil

		default:
			// Check for custom actions
			if action, ok := s.actions[msg.String()]; ok {
				return nil, &screens.NavigateResult{
					ResultData: action,
				}
			}
		}
	}

	return nil, nil
}

// View renders the detail screen using StandardView.
func (s *BaseDetailScreen[T]) View() string {
	// Render content with current dimensions
	content := ""
	if s.contentRenderer != nil {
		content = s.contentRenderer(s.data, s.Width(), s.Height())
	}

	// TODO: Apply scroll offset to content
	// For now, just render the full content
	// In future, we can slice by lines based on scrollOffset

	// Use BaseScreen's CreateView helper for StandardView integration
	return s.CreateView(s.breadcrumbs, content, s.footer)
}

// SetFooter updates the footer help text.
func (s *BaseDetailScreen[T]) SetFooter(footer string) {
	s.footer = footer
}

// GetData returns the data structure being displayed.
//
// This allows the intent to access the data after the screen completes.
func (s *BaseDetailScreen[T]) GetData() T {
	return s.data
}

// GetScrollOffset returns the current scroll offset.
func (s *BaseDetailScreen[T]) GetScrollOffset() int {
	return s.scrollOffset
}

// SetScrollOffset sets the scroll offset.
//
// This is useful for restoring scroll position when navigating back.
func (s *BaseDetailScreen[T]) SetScrollOffset(offset int) {
	if offset < 0 {
		offset = 0
	}
	s.scrollOffset = offset
}

// AddAction registers a custom action key.
//
// When the user presses the specified key, a NavigateResult is returned
// with the action name as data.
//
// Example:
//
//	screen.AddAction("e", "edit")   // Pressing 'e' returns NavigateResult{Data: "edit"}
//	screen.AddAction("d", "delete") // Pressing 'd' returns NavigateResult{Data: "delete"}
//
// The intent can then handle these actions:
//
//	if result.Type() == screens.ResultNavigate {
//	    action := result.Data().(string)
//	    switch action {
//	    case "edit":
//	        // Transition to edit screen
//	    case "delete":
//	        // Show delete confirmation
//	    }
//	}
func (s *BaseDetailScreen[T]) AddAction(key, action string) {
	s.actions[key] = action
}

// RestoreFromMetadata restores screen state from metadata.
//
// This is useful when navigating back to this screen from another screen.
// The metadata should contain the scroll_offset key.
func (s *BaseDetailScreen[T]) RestoreFromMetadata(metadata map[string]interface{}) {
	if offset, ok := metadata["scroll_offset"]; ok {
		if offsetInt, ok := offset.(int); ok {
			s.SetScrollOffset(offsetInt)
		}
	}
}
