package event

import (
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// Detail wraps feedback.DetailModal for viewing career event details.
// This uses the generic UIKit DetailModal with event-specific content rendering.
//
// Usage:
//
//	view := eventviews.NewDetail(event, theme)
//	modal.SetDimensions(width, height)
//	modal.Show()
//
//	// In Update:
//	model, cmd := modal.Update(msg)
//
//	// In View:
//	if modal.IsVisible() {
//	    return behaviors.RenderModalOverlay(modal, background)
//	}
type Detail struct {
	*widgets.ModalDelegate
	event            display.Event
	theme            themes.Theme
	showSkillsOption bool
}

// NewDetail creates a new event detail modal.
//
// Expected:
//   - event must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized Detail ready for use.
//
// Side effects:
//   - None.
func NewDetail(event display.Event, theme themes.Theme) *Detail {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	content := RenderEventDetailContent(event, theme)

	modal := feedback.NewDetailModal("Event Details", content)
	if theme != nil {
		modal = modal.WithTheme(theme)
	}

	m := &Detail{
		ModalDelegate:    widgets.NewModalDelegate(modal),
		event:            event,
		theme:            theme,
		showSkillsOption: true,
	}

	m.updateFooterBadges()

	return m
}

// WithShowSkillsOption sets whether to show the "s: Skills" option in the footer.
//
// Expected:
//   - bool must be valid.
//
// Returns:
//   - A fully initialized Detail ready for use.
//
// Side effects:
//   - None.
func (m *Detail) WithShowSkillsOption(show bool) *Detail {
	m.showSkillsOption = show
	m.updateFooterBadges()
	return m
}

func (m *Detail) updateFooterBadges() {
	theme := m.theme
	if theme == nil {
		theme = themes2.Default()
	}

	badges := []*primitives.Badge{
		primitives.HelpKeyBadge("e", "Edit", theme),
		primitives.HelpKeyBadge("d", "Delete", theme),
	}
	if m.showSkillsOption {
		badges = append(badges, primitives.HelpKeyBadge("s", "Skills", theme))
	}
	badges = append(badges, primitives.HelpKeyBadge("Enter/Esc", "Close", theme))

	m.SetModal(m.Modal().WithFooterBadges(badges...))
}

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *Detail) Init() tea.Cmd {
	m.updateFooterBadges()
	return m.Modal().Init()
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
func (m *Detail) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.ModalUpdate(msg)
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Detail) View() string {
	return m.ModalView()
}

// SetEvent updates the event being displayed.
//
// Expected:
//   - event must be valid.
//
// Side effects:
//   - None.
func (m *Detail) SetEvent(event display.Event) {
	m.event = event
	content := RenderEventDetailContent(event, m.theme)
	m.SetContent(content)
}

// GetEvent returns the event being displayed.
//
// Returns:
//   - A display.Event value.
//
// Side effects:
//   - None.
func (m *Detail) GetEvent() display.Event {
	return m.event
}
