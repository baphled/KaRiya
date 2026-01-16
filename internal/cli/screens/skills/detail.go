package skills

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SkillDetailState represents the state constant for this screen
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
// - internal/cli/screens/base/detail_screen.go (BaseDetailScreen pattern)
type SkillDetailScreen struct {
	*base.BaseScreen

	skill *career.Skill
}

// NewSkillDetailScreen creates a new skill detail screen.
func NewSkillDetailScreen(skill *career.Skill) *SkillDetailScreen {
	return &SkillDetailScreen{
		BaseScreen: base.NewBaseScreen(),
		skill:      skill,
	}
}

// Update handles messages and returns result for actions.
func (s *SkillDetailScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, &screens.CancelResult{}

		case "enter":
			return nil, &screens.NavigateResult{
				ResultData: "back",
			}

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
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(titleStyle.Render(s.skill.Name))
	b.WriteString("\n\n")

	// Category badge
	categoryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(s.getCategoryColor(s.skill.Category)).
		Padding(0, 1)
	b.WriteString(categoryStyle.Render(s.skill.Category))
	b.WriteString("\n\n")

	// Details section
	s.renderField(&b, "Level", s.skill.Level)

	if s.skill.YearsUsed != nil {
		s.renderField(&b, "Years of Experience", fmt.Sprintf("%d years", *s.skill.YearsUsed))
	}

	if s.skill.LastUsed != nil {
		s.renderField(&b, "Last Used", s.skill.LastUsed.Format("2006-01-02"))
	}

	b.WriteString("\n")
	s.renderField(&b, "Created", s.skill.CreatedAt.Format("2006-01-02 15:04:05"))
	s.renderField(&b, "Updated", s.skill.UpdatedAt.Format("2006-01-02 15:04:05"))

	content := b.String()
	footer := "Enter/Esc: Back  e: Edit  d: Delete"

	return s.CreateView([]string{"Main Menu", "Manage Skills", "Skill Details"}, content, footer)
}

// renderField renders a labeled field.
func (s *SkillDetailScreen) renderField(b *strings.Builder, label, value string) {
	if value == "" {
		return
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Bold(true)

	b.WriteString(labelStyle.Render(label + ":"))
	b.WriteString(" ")
	b.WriteString(value)
	b.WriteString("\n")
}

// getCategoryColor returns a color for the skill category.
func (s *SkillDetailScreen) getCategoryColor(category string) lipgloss.Color {
	colors := map[string]lipgloss.Color{
		"backend":  lipgloss.Color("10"), // Green
		"frontend": lipgloss.Color("12"), // Blue
		"devops":   lipgloss.Color("11"), // Yellow
		"database": lipgloss.Color("13"), // Magenta
		"cloud":    lipgloss.Color("14"), // Cyan
		"mobile":   lipgloss.Color("9"),  // Red
		"tooling":  lipgloss.Color("8"),  // Gray
	}

	if color, ok := colors[category]; ok {
		return color
	}
	return lipgloss.Color("7") // Default white
}

// GetSkill returns the skill being displayed.
func (s *SkillDetailScreen) GetSkill() *career.Skill {
	return s.skill
}
