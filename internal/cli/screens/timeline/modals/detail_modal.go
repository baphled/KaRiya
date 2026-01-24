package modals

import (
	"fmt"

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
	event            *career.CareerEvent
	theme            themes.Theme
	showSkillsOption bool
}

// NewEventDetailModal creates a new event detail modal.
// By default, shows the "s: Skills" option in the footer.
func NewEventDetailModal(event *career.CareerEvent, theme themes.Theme) *EventDetailModal {
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

	// Initialize footer badges immediately.
	m.updateFooterBadges()

	return m
}

// WithShowSkillsOption sets whether to show the "s: Skills" option in the footer.
func (m *EventDetailModal) WithShowSkillsOption(show bool) *EventDetailModal {
	m.showSkillsOption = show
	m.updateFooterBadges()
	return m
}

// updateFooterBadges updates the footer badges based on current options.
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
func (m *EventDetailModal) SetEvent(event *career.CareerEvent) {
	m.event = event
	content := RenderEventDetailContent(event, m.theme)
	m.modal.SetContent(content)
}

// GetEvent returns the event being displayed.
func (m *EventDetailModal) GetEvent() *career.CareerEvent {
	return m.event
}

// SkillsDetailModal wraps feedback.DetailModal for viewing skills associated with an event.
type SkillsDetailModal struct {
	modal   *feedback.DetailModal
	eventID string
	skills  []*career.Skill
	theme   themes.Theme
}

// NewSkillsDetailModal creates a new skills detail modal.
func NewSkillsDetailModal(eventID string, skills []*career.Skill, theme themes.Theme) *SkillsDetailModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	content := RenderSkillsContent(skills, theme)
	title := fmt.Sprintf("Skills (%d)", len(skills))

	modal := feedback.NewDetailModal(title, content)
	if theme != nil {
		modal = modal.WithTheme(theme)
	}

	return &SkillsDetailModal{
		modal:   modal,
		eventID: eventID,
		skills:  skills,
		theme:   theme,
	}
}

// Init initializes the modal.
func (m *SkillsDetailModal) Init() tea.Cmd {
	return m.modal.Init()
}

// Update handles keyboard input and window sizing.
func (m *SkillsDetailModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.modal.Update(msg)
}

// View renders the modal content.
func (m *SkillsDetailModal) View() string {
	return m.modal.View()
}

// IsVisible returns whether the modal is currently visible.
func (m *SkillsDetailModal) IsVisible() bool {
	return m.modal.IsVisible()
}

// Show makes the modal visible.
func (m *SkillsDetailModal) Show() {
	m.modal.Show()
}

// Hide hides the modal.
func (m *SkillsDetailModal) Hide() {
	m.modal.Hide()
}

// SetDimensions sets the terminal dimensions.
func (m *SkillsDetailModal) SetDimensions(width, height int) {
	m.modal.SetDimensions(width, height)
}

// SetSkills updates the skills being displayed.
func (m *SkillsDetailModal) SetSkills(skills []*career.Skill) {
	m.skills = skills
	content := RenderSkillsContent(skills, m.theme)
	m.modal.SetContent(content)
	m.modal.SetTitle(fmt.Sprintf("Skills (%d)", len(skills)))
}

// GetEventID returns the event ID this modal is showing skills for.
func (m *SkillsDetailModal) GetEventID() string {
	return m.eventID
}
