package skills

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SkillDetailState represents the state constant for this screen.
const SkillDetailState = "skill_detail"

// SkillDetailScreen displays detailed information about a skill.
//
// This screen provides:
// - Full skill information display
// - Actions: Edit (e), Delete (d), Back (enter/esc)
// - Formatted display with labels and values
//
// Example usage:
//
//	skill := &career.Skill{...}
//	screen := skills.NewSkillDetailScreen(skill)
//
//	// In intent Update:
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    action := result.Data().(string)
//	    switch action {
//	    case "edit":
//	        // Navigate to edit form
//	    case "delete":
//	        // Navigate to delete confirmation
//	    case "back":
//	        // Go back to list
//	    }
//	}
//
// Related:
// - internal/cli/screens/base/detail_screen.go (DetailScreen pattern).
type SkillDetailScreen struct {
	*base.Screen

	skill *career.Skill
}

// NewSkillDetailScreen creates a new skill detail screen.
func NewSkillDetailScreen(skill *career.Skill) *SkillDetailScreen {
	return &SkillDetailScreen{
		Screen: base.NewBaseScreen(),
		skill:  skill,
	}
}

// Update handles messages and returns result for actions.
func (s *SkillDetailScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		// Handle special keys.
		switch msg.Type {
		case tea.KeyEsc:
			return nil, &screens.CancelResult{}

		case tea.KeyEnter:
			return nil, &screens.NavigateResult{
				ResultData: "back",
			}
		}

		// Handle character keys.
		switch msg.String() {
		case "e":
			return nil, &screens.NavigateResult{
				ResultData: "edit",
			}

		case "d":
			return nil, &screens.NavigateResult{
				ResultData: "delete",
			}
		}
	}

	return nil, nil
}

// View renders the skill detail screen.
func (s *SkillDetailScreen) View() string {
	// Get theme for styling (fall back to default if not set).
	var th themes.Theme
	if screenTheme := s.Theme(); screenTheme != nil {
		if t, ok := screenTheme.(themes.Theme); ok {
			th = t
		}
	}
	if th == nil {
		th = themes.NewDefaultTheme()
	}

	var b strings.Builder

	// Title using theme colors.
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(th.PrimaryColor())
	b.WriteString(titleStyle.Render(s.skill.Name))
	b.WriteString("\n\n")

	// Category badge using theme.
	categoryStyle := lipgloss.NewStyle().
		Foreground(th.BackgroundColor()).
		Background(s.getCategoryColor(s.skill.Category)).
		Padding(0, 1)
	b.WriteString(categoryStyle.Render(s.skill.Category))
	b.WriteString("\n\n")

	// Details section.
	s.renderField(&b, "Level", s.skill.Level, th)

	if s.skill.YearsUsed != nil {
		s.renderField(&b, "Years of Experience", fmt.Sprintf("%d years", *s.skill.YearsUsed), th)
	}

	if s.skill.LastUsed != nil {
		s.renderField(&b, "Last Used", s.skill.LastUsed.Format("2006-01-02"), th)
	}

	b.WriteString("\n")
	s.renderField(&b, "Created", s.skill.CreatedAt.Format("2006-01-02 15:04:05"), th)
	s.renderField(&b, "Updated", s.skill.UpdatedAt.Format("2006-01-02 15:04:05"), th)

	content := b.String()

	// Build footer using UIKit primitives for consistent styling.
	footer := primitives.RenderHelpFooter(th,
		primitives.BackBadge(th),
		primitives.EditBadge(th),
		primitives.DeleteBadge(th),
		primitives.HelpKeyBadge("Ctrl+E", "Events", th),
	)

	return s.CreateView([]string{"Main Menu", "Manage Skills", "Skill Details"}, content, footer)
}

// renderField renders a labeled field using theme colors.
func (s *SkillDetailScreen) renderField(b *strings.Builder, label, value string, th themes.Theme) {
	if value == "" {
		return
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(th.MutedColor()).
		Bold(true)

	b.WriteString(labelStyle.Render(label + ":"))
	b.WriteString(" ")
	b.WriteString(value)
	b.WriteString("\n")
}

// getCategoryColor returns a color for the skill category.
func (s *SkillDetailScreen) getCategoryColor(category string) lipgloss.Color {
	colors := map[string]lipgloss.Color{
		"backend":  lipgloss.Color("10"),
		"frontend": lipgloss.Color("12"),
		"devops":   lipgloss.Color("11"),
		"database": lipgloss.Color("13"),
		"cloud":    lipgloss.Color("14"),
		"mobile":   lipgloss.Color("9"),
		"tooling":  lipgloss.Color("8"),
	}

	if color, ok := colors[category]; ok {
		return color
	}
	return lipgloss.Color("7")
}

// GetSkill returns the skill being displayed.
func (s *SkillDetailScreen) GetSkill() *career.Skill {
	return s.skill
}
