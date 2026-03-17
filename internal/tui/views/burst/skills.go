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

// Skills wraps feedback.DetailModal for viewing skills associated with a burst.
type Skills struct {
	modal   *feedback.DetailModal
	burstID string
	skills  []display.Skill
	theme   themes.Theme
}

// NewSkills creates a new burst skills modal.
//
// Expected:
//   - burstid must be a valid string.
//   - burstname must be a valid string.
//   - skills must be a valid slice of *career.Skill.
//   - theme must be a valid Theme instance (can be nil).
//
// Returns:
//   - A fully initialized Skills ready for use.
//
// Side effects:
//   - None.
func NewSkills(burstID string, burstName string, skills []display.Skill, theme themes.Theme) *Skills {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	content := renderSkillsContent(skills, theme)
	title := fmt.Sprintf("Skills in Burst: %s (%d)", burstName, len(skills))

	modal := feedback.NewDetailModal(title, content)
	if theme != nil {
		modal = modal.WithTheme(theme)
	}

	return &Skills{
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
func (m *Skills) Init() tea.Cmd {
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
func (m *Skills) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.modal.Update(msg)
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Skills) View() string {
	return m.modal.View()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *Skills) IsVisible() bool {
	return m.modal.IsVisible()
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *Skills) Show() {
	m.modal.Show()
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *Skills) Hide() {
	m.modal.Hide()
}

// SetDimensions sets the terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *Skills) SetDimensions(width, height int) {
	m.modal.SetDimensions(width, height)
}

// SetSkills updates the skills being displayed.
//
// Expected:
//   - skill must be valid.
//
// Side effects:
//   - None.
func (m *Skills) SetSkills(skills []display.Skill) {
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
func (m *Skills) GetBurstID() string {
	return m.burstID
}

// renderSkillsContent renders the skills as formatted text.
func renderSkillsContent(skills []display.Skill, theme themes.Theme) string {
	if theme == nil {
		theme = themes2.Default()
	}

	if len(skills) == 0 {
		return primitives.WarningText("No skills associated with this burst yet. Use 'i' to infer skills.", theme).Render()
	}

	var b strings.Builder

	for idx := range skills {
		skill := skills[idx]
		fmt.Fprintf(&b, "%d. %s\n", idx+1, skill.Name)

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
