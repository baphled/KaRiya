// Package skills provides screens for skill management workflows.
package skills

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SkillsListState represents the state constant for this screen
const SkillsListState = "skills_list"

// SkillsListScreen displays a list of skills with actions.
//
// This screen provides:
// - List of skills with name, category, and level
// - Navigation (↑↓/jk/g/G)
// - Actions: View (enter), Add (a), Edit (e), Delete (d)
// - Empty state handling
//
// Example usage:
//
//	skills := []*career.Skill{...}
//	screen := skills.NewSkillsListScreen(skills)
//
//	// In intent Update:
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    data := result.Data().(map[string]interface{})
//	    action := data["action"].(string)
//	    switch action {
//	    case "view":
//	        skill := data["skill"].(*career.Skill)
//	        // Navigate to detail screen
//	    case "add":
//	        // Navigate to add form
//	    case "edit":
//	        skill := data["skill"].(*career.Skill)
//	        // Navigate to edit form
//	    case "delete":
//	        skill := data["skill"].(*career.Skill)
//	        // Navigate to delete confirmation
//	    }
//	}
//
// Related:
// - internal/cli/screens/base/select_screen.go (BaseSelectScreen pattern)
// - docs/workflows/MANAGE_SKILLS_WORKFLOW.md (Skills workflow guide)
type SkillsListScreen struct {
	*base.BaseScreen

	skills        []*career.Skill
	selectedIndex int
}

// NewSkillsListScreen creates a new skills list screen.
func NewSkillsListScreen(skills []*career.Skill) *SkillsListScreen {
	return &SkillsListScreen{
		BaseScreen:    base.NewBaseScreen(),
		skills:        skills,
		selectedIndex: 0,
	}
}

// Update handles messages and returns result for actions.
func (s *SkillsListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, &screens.CancelResult{}

		case "down", "j":
			if len(s.skills) > 0 && s.selectedIndex < len(s.skills)-1 {
				s.selectedIndex++
			}
			return nil, nil

		case "up", "k":
			if s.selectedIndex > 0 {
				s.selectedIndex--
			}
			return nil, nil

		case "g":
			s.selectedIndex = 0
			return nil, nil

		case "G":
			if len(s.skills) > 0 {
				s.selectedIndex = len(s.skills) - 1
			}
			return nil, nil

		case "enter":
			if len(s.skills) > 0 {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "view",
						"skill":  s.skills[s.selectedIndex],
					},
				}
			}
			return nil, nil

		case "a":
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "add",
				},
			}

		case "e":
			if len(s.skills) > 0 {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "edit",
						"skill":  s.skills[s.selectedIndex],
					},
				}
			}
			return nil, nil

		case "d":
			if len(s.skills) > 0 {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "delete",
						"skill":  s.skills[s.selectedIndex],
					},
				}
			}
			return nil, nil
		}
	}

	return nil, nil
}

// View renders the skills list screen.
func (s *SkillsListScreen) View() string {
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(titleStyle.Render("Skills Management"))
	b.WriteString("\n\n")

	// Empty state
	if len(s.skills) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		b.WriteString(emptyStyle.Render("No skills found. Press 'a' to add your first skill."))
		b.WriteString("\n")
	} else {
		// Render skill list
		for i, skill := range s.skills {
			s.renderSkillItem(&b, i, skill)
		}
	}

	content := b.String()
	footer := "Enter: View  a: Add  e: Edit  d: Delete  ↑↓/jk: Navigate  g/G: Top/Bottom  Esc: Back"

	return s.CreateView([]string{"Main Menu", "Manage Skills"}, content, footer)
}

// renderSkillItem renders a single skill list item.
func (s *SkillsListScreen) renderSkillItem(b *strings.Builder, index int, skill *career.Skill) {
	// Selection indicator
	if index == s.selectedIndex {
		b.WriteString("▶ ")
	} else {
		b.WriteString("  ")
	}

	// Skill name (bold if selected)
	nameStyle := lipgloss.NewStyle()
	if index == s.selectedIndex {
		nameStyle = nameStyle.Bold(true).Foreground(lipgloss.Color("12"))
	}
	b.WriteString(nameStyle.Render(skill.Name))

	// Category badge
	categoryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(s.getCategoryColor(skill.Category)).
		Padding(0, 1)
	b.WriteString("  ")
	b.WriteString(categoryStyle.Render(skill.Category))

	// Level badge (if set)
	if skill.Level != "" {
		levelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true)
		b.WriteString("  ")
		b.WriteString(levelStyle.Render(skill.Level))
	}

	b.WriteString("\n")
}

// getCategoryColor returns a color for the skill category.
func (s *SkillsListScreen) getCategoryColor(category string) lipgloss.Color {
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

// GetSelectedIndex returns the currently selected index.
func (s *SkillsListScreen) GetSelectedIndex() int {
	return s.selectedIndex
}

// SetSelectedIndex sets the currently selected index.
func (s *SkillsListScreen) SetSelectedIndex(index int) {
	if index >= 0 && index < len(s.skills) {
		s.selectedIndex = index
	}
}
