package modals

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// BurstSkillsModal wraps feedback.DetailModal for viewing skills associated with a burst.
type BurstSkillsModal struct {
	modal   *feedback.DetailModal
	burstID string
	skills  []*career.Skill
	theme   themes.Theme
}

// NewBurstSkillsModal creates a new burst skills modal.
//
// Expected:
//   - burstid must be a valid string.
//   - burstname must be a valid string.
//   - skills must be a valid slice of *career.Skill.
//   - theme must be a valid Theme instance (can be nil).
//
// Returns:
//   - A fully initialized BurstSkillsModal ready for use.
//
// Side effects:
//   - None.
func NewBurstSkillsModal(burstID string, burstName string, skills []*career.Skill, theme themes.Theme) *BurstSkillsModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	content := renderSkillsContent(skills, theme)
	title := fmt.Sprintf("Skills in Burst: %s (%d)", burstName, len(skills))

	modal := feedback.NewDetailModal(title, content)
	if theme != nil {
		modal = modal.WithTheme(theme)
	}

	return &BurstSkillsModal{
		modal:   modal,
		burstID: burstID,
		skills:  skills,
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
func (m *BurstSkillsModal) Init() tea.Cmd {
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
func (m *BurstSkillsModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.modal.Update(msg)
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstSkillsModal) View() string {
	return m.modal.View()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *BurstSkillsModal) IsVisible() bool {
	return m.modal.IsVisible()
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *BurstSkillsModal) Show() {
	m.modal.Show()
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *BurstSkillsModal) Hide() {
	m.modal.Hide()
}

// SetDimensions sets the terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *BurstSkillsModal) SetDimensions(width, height int) {
	m.modal.SetDimensions(width, height)
}

// SetSkills updates the skills being displayed.
//
// Expected:
//   - skill must be valid.
//
// Side effects:
//   - None.
func (m *BurstSkillsModal) SetSkills(skills []*career.Skill) {
	m.skills = skills
	content := renderSkillsContent(skills, m.theme)
	m.modal.SetContent(content)
	m.modal.SetTitle(fmt.Sprintf("Skills (%d)", len(skills)))
}

// GetBurstID returns the burst ID this modal is showing skills for.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstSkillsModal) GetBurstID() string {
	return m.burstID
}

// renderSkillsContent renders the skills as formatted text.
func renderSkillsContent(skills []*career.Skill, theme themes.Theme) string {
	if theme == nil {
		theme = themes2.Default()
	}

	if len(skills) == 0 {
		return primitives.WarningText("No skills associated with this burst yet. Use 'i' to infer skills.", theme).Render()
	}

	var b strings.Builder

	for idx, skill := range skills {
		b.WriteString(fmt.Sprintf("%d. %s\n", idx+1, skill.Name))

		if skill.Category != "" {
			catText := "   Category: " + skill.Category
			b.WriteString(primitives.NewText(catText, theme).
				Foreground(theme.SecondaryColor()).Render())
			b.WriteString("\n")
		}

		if skill.Level != "" {
			levelText := "   Level: " + skill.Level
			b.WriteString(primitives.NewText(levelText, theme).
				Foreground(theme.SecondaryColor()).Render())
			b.WriteString("\n")
		}

		if idx < len(skills)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
