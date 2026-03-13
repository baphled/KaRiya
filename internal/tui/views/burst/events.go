package burst

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/ui/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
)

// Events wraps feedback.DetailModal for viewing events in a burst.
type Events struct {
	modal   *feedback.DetailModal
	burstID string
	events  []display.Event
	theme   themes.Theme
}

// NewEvents creates a new burst events modal.
//
// Expected:
//   - burstid must be a valid string.
//   - burstname must be a valid string.
//   - events must be a valid slice of *career.Event.
//   - theme must be a valid Theme instance (can be nil).
//
// Returns:
//   - A fully initialized Events ready for use.
//
// Side effects:
//   - None.
func NewEvents(burstID string, burstName string, events []display.Event, theme themes.Theme) *Events {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	content := renderEventsContent(events, theme)
	title := fmt.Sprintf("Events in Burst: %s (%d)", burstName, len(events))

	modal := feedback.NewDetailModal(title, content)
	if theme != nil {
		modal = modal.WithTheme(theme)
	}

	return &Events{
		modal:   modal,
		burstID: burstID,
		events:  events,
		theme:   theme,
	}
}

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *Events) Init() tea.Cmd {
	return m.modal.Init()
}

// Update handles keyboard input and window sizing.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Model: the updated model.
//   - tea.Cmd: command to execute.
//
// Side effects:
//   - Delegates to underlying modal.
func (m *Events) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.modal.Update(msg)
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Events) View() string {
	return m.modal.View()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *Events) IsVisible() bool {
	return m.modal.IsVisible()
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *Events) Show() {
	m.modal.Show()
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *Events) Hide() {
	m.modal.Hide()
}

// SetDimensions sets the terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *Events) SetDimensions(width, height int) {
	m.modal.SetDimensions(width, height)
}

// SetEvents updates the events being displayed.
//
// Expected:
//   - event must be valid.
//
// Side effects:
//   - None.
func (m *Events) SetEvents(events []display.Event) {
	m.events = events
	content := renderEventsContent(events, m.theme)
	m.modal.SetContent(content)
	m.modal.SetTitle(fmt.Sprintf("Events (%d)", len(events)))
}

// GetBurstID returns the burst ID this modal is showing events for.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Events) GetBurstID() string {
	return m.burstID
}

// renderEventsContent renders the events as formatted text.
func renderEventsContent(events []display.Event, theme themes.Theme) string {
	if theme == nil {
		theme = themes2.Default()
	}

	if len(events) == 0 {
		return primitives.ErrorText("No events found for this burst.", theme).Render()
	}

	var b strings.Builder

	for idx := range events {
		event := events[idx]
		fmt.Fprintf(&b, "%d. %s\n", idx+1, event.Text)
		dateText := "   Date: " + event.Date.Format("2006-01-02")
		b.WriteString(primitives.NewText(dateText, theme).
			Foreground(theme.SecondaryColor()).Render())
		b.WriteString("\n")

		if len(event.Tags) > 0 {
			tagsText := "   Tags: " + strings.Join(event.Tags, ", ")
			b.WriteString(primitives.NewText(tagsText, theme).
				Foreground(theme.SecondaryColor()).Render())
			b.WriteString("\n")
		}

		if idx < len(events)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
