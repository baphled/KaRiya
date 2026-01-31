package modals

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// EventDetailModal wraps feedback.DetailModal for viewing career event details.
// This uses the generic UIKit DetailModal with event-specific content rendering.
//
// Usage:
//
//	modal := modals.NewEventDetailModal(event, theme)
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
type EventDetailModal struct {
	modal            *feedback.DetailModal
	event            *career.Event
	theme            themes.Theme
	showSkillsOption bool
}

// NewEventDetailModal creates a new event detail modal.
// By default, shows the "s: Skills" option in the footer.
func NewEventDetailModal(event *career.Event, theme themes.Theme) *EventDetailModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	content := RenderEventDetailContent(event, theme)

	modal := feedback.NewDetailModal("Event Details", content)
	if theme != nil {
		modal = modal.WithTheme(theme)
	}

	m := &EventDetailModal{
		modal:            modal,
		event:            event,
		theme:            theme,
		showSkillsOption: true,
	}

	m.updateFooterBadges()

	return m
}

// WithShowSkillsOption sets whether to show the "s: Skills" option in the footer.
func (m *EventDetailModal) WithShowSkillsOption(show bool) *EventDetailModal {
	m.showSkillsOption = show
	m.updateFooterBadges()
	return m
}

func (m *EventDetailModal) updateFooterBadges() {
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

	m.modal = m.modal.WithFooterBadges(badges...)
}

// Init initializes the modal.
func (m *EventDetailModal) Init() tea.Cmd {
	m.updateFooterBadges()
	return m.modal.Init()
}

// Update handles keyboard input and window sizing.
func (m *EventDetailModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.modal.Update(msg)
}

// View renders the modal content.
func (m *EventDetailModal) View() string {
	return m.modal.View()
}

// IsVisible returns whether the modal is currently visible.
func (m *EventDetailModal) IsVisible() bool {
	return m.modal.IsVisible()
}

// Show makes the modal visible.
func (m *EventDetailModal) Show() {
	m.modal.Show()
}

// Hide hides the modal.
func (m *EventDetailModal) Hide() {
	m.modal.Hide()
}

// SetDimensions sets the terminal dimensions.
func (m *EventDetailModal) SetDimensions(width, height int) {
	m.modal.SetDimensions(width, height)
}

// SetEvent updates the event being displayed.
func (m *EventDetailModal) SetEvent(event *career.Event) {
	m.event = event
	content := RenderEventDetailContent(event, m.theme)
	m.modal.SetContent(content)
}

// GetEvent returns the event being displayed.
func (m *EventDetailModal) GetEvent() *career.Event {
	return m.event
}
