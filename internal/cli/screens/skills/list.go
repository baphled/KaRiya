// Package skills provides screens for skill management workflows.
package skills

import (
	"strconv"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// SkillsListState represents the state constant for this screen
const SkillsListState = "skills_list"

// skillRowFormatter formats a skill for table display.
func skillRowFormatter(skill *career.Skill, index int, eventCounts map[string]int) []string {
	// Name with selection indicator (handled by TableBehavior)
	name := skill.Name

	// Category
	category := skill.Category
	if category == "" {
		category = "-"
	}

	// Level
	level := skill.Level
	if level == "" {
		level = "-"
	}

	// Years of experience
	years := "-"
	if skill.YearsUsed != nil && *skill.YearsUsed > 0 {
		years = strconv.Itoa(*skill.YearsUsed)
	}

	// Event count (from eventCounts map)
	eventCount := "-"
	if count, ok := eventCounts[skill.ID]; ok && count > 0 {
		eventCount = strconv.Itoa(count)
	}

	return []string{name, category, level, years, eventCount}
}

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
	tableBehavior *behaviors.TableBehavior[*career.Skill]
	eventCounts   map[string]int // Skill ID -> event count
}

// NewSkillsListScreen creates a new skills list screen.
func NewSkillsListScreen(skills []*career.Skill) *SkillsListScreen {
	// Create table columns matching legacy format
	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
		{Title: "Years", Width: 8},
		{Title: "Events", Width: 8},
	}

	// Initialize event counts (empty, will be populated via SetEventCounts)
	eventCounts := make(map[string]int)

	// Create row formatter that uses the eventCounts
	rowFormatter := func(skill *career.Skill, index int) []string {
		return skillRowFormatter(skill, index, eventCounts)
	}

	// Create table behavior (nil theme, will be set via SetTheme)
	tableBehavior := behaviors.NewTableBehavior[*career.Skill](nil, columns, rowFormatter).
		PageSize(15).
		PaginationPrefix("Skills").
		EmptyMessage("No skills found. Press 'a' to add your first skill.")

	// Set initial items
	tableBehavior.SetItems(skills)

	screen := &SkillsListScreen{
		BaseScreen:    base.NewBaseScreen(),
		skills:        skills,
		tableBehavior: tableBehavior,
		eventCounts:   eventCounts,
	}

	return screen
}

// Update handles messages and returns result for actions.
func (s *SkillsListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		key := msg.String()

		// Handle escape
		if key == "esc" {
			return nil, &screens.CancelResult{}
		}

		// Handle navigation keys - delegate to TableBehavior
		if key == "down" || key == "j" || key == "up" || key == "k" ||
			key == "g" || key == "G" || key == "pgup" || key == "pgdown" ||
			key == "home" || key == "end" {
			s.tableBehavior.HandleNavigation(key)
			return nil, nil
		}

		// Handle action keys
		switch key {
		case "enter":
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "view",
						"skill":  *selected,
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
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "edit",
						"skill":  *selected,
					},
				}
			}
			return nil, nil

		case "d":
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "delete",
						"skill":  *selected,
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
	// Render table via behavior
	content := s.tableBehavior.Render()

	// Footer with actions (matching legacy)
	footer := "Enter: View  a: Add  e: Edit  d: Delete  ↑↓/jk: Navigate  g/G: Top/Bottom  Esc: Back"

	// Handle empty state footer
	if len(s.skills) == 0 {
		footer = "a: Add skill  Esc: Back"
	}

	return s.CreateView([]string{"Main Menu", "Manage Skills"}, content, footer)
}

// SetTheme applies theme to the table behavior.
func (s *SkillsListScreen) SetTheme(theme interface{}) {
	s.BaseScreen.SetTheme(theme)
	// Apply theme to table behavior if available
	if t, ok := theme.(themes.Theme); ok && t != nil {
		s.tableBehavior.SetTheme(t)
	}
}

// SetEventCounts sets the event counts for skills (used for displaying event count column).
func (s *SkillsListScreen) SetEventCounts(counts map[string]int) {
	// Update the shared eventCounts map
	for k, v := range counts {
		s.eventCounts[k] = v
	}

	// Create new row formatter with updated counts
	rowFormatter := func(skill *career.Skill, index int) []string {
		return skillRowFormatter(skill, index, s.eventCounts)
	}

	// Rebuild table behavior with new formatter
	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
		{Title: "Years", Width: 8},
		{Title: "Events", Width: 8},
	}

	s.tableBehavior = behaviors.NewTableBehavior[*career.Skill](nil, columns, rowFormatter).
		PageSize(15).
		PaginationPrefix("Skills").
		EmptyMessage("No skills found. Press 'a' to add your first skill.")

	s.tableBehavior.SetItems(s.skills)
}

// GetSelectedIndex returns the currently selected index.
func (s *SkillsListScreen) GetSelectedIndex() int {
	return s.tableBehavior.GetSelectedIndex()
}

// SetSelectedIndex sets the currently selected index.
func (s *SkillsListScreen) SetSelectedIndex(index int) {
	s.tableBehavior.SetSelectedIndex(index)
}
