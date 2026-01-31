package modals

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

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
