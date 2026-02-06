// Package skills provides screens for skill management workflows.
package skills

import (
	"strconv"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// SkillsListState represents the state constant for this screen.
const SkillsListState = "skills_list"

// skillRowFormatter formats a skill for table display.
func skillRowFormatter(skill *career.Skill, _ int, eventCounts map[string]int) []string {
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
// - internal/cli/screens/base/select_screen.go (SelectScreen pattern)
// - docs/workflows/MANAGE_SKILLS_WORKFLOW.md (Skills workflow guide).
type SkillsListScreen struct {
	*base.Screen

	skills        []*career.Skill
	tableBehavior *behaviors.TableBehavior[*career.Skill]
	eventCounts   map[string]int
}

// NewSkillsListScreen creates a new skills list screen.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A fully initialized SkillsListScreen ready for use.
//
// Side effects:
//   - None.
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
		Screen:        base.NewBaseScreen(),
		skills:        skills,
		tableBehavior: tableBehavior,
		eventCounts:   eventCounts,
	}

	return screen
}

// Update handles messages and returns result for actions.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - screens.ScreenResult: result indicating navigation or action.
//
// Side effects:
//   - May update table selection.
//   - May return CancelResult or NavigateResult.
func (s *SkillsListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		return s.handleKeyMsg(msg)
	}

	return nil, nil
}

// handleKeyMsg handles keyboard input and returns appropriate result.
func (s *SkillsListScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Cmd, screens.ScreenResult) {
	// Handle escape/back.
	if msg.Type == tea.KeyEsc {
		return nil, &screens.CancelResult{}
	}

	// Handle navigation keys - delegate to TableBehavior.
	if s.tableBehavior.HandleNavigation(msg.String()) {
		return nil, nil
	}

	// Handle action keys.
	return s.handleActionKey(msg)
}

// handleActionKey handles action-specific key presses.
func (s *SkillsListScreen) handleActionKey(msg tea.KeyMsg) (tea.Cmd, screens.ScreenResult) {
	// Handle enter key for view action.
	if msg.Type == tea.KeyEnter {
		return s.handleViewAction()
	}

	// Handle 's' key - view skill events.
	if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 && msg.Runes[0] == 's' {
		// Navigate to skill events screen for selected skill
		if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "view_events",
					"skill":  *selected,
				},
			}
		}
	}

	// Handle character keys for other actions.
	switch msg.String() {
	case "a":
		return s.handleAddAction()
	case "e":
		return s.handleEditAction()
	case "d":
		return s.handleDeleteAction()
	}
	return nil, nil
}

// handleViewAction handles viewing a selected skill.
func (s *SkillsListScreen) handleViewAction() (tea.Cmd, screens.ScreenResult) {
	if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
		return nil, &screens.NavigateResult{
			ResultData: map[string]interface{}{
				"action": "view",
				"skill":  *selected,
			},
		}
	}
	return nil, nil
}

// handleAddAction handles adding a new skill.
func (s *SkillsListScreen) handleAddAction() (tea.Cmd, screens.ScreenResult) {
	return nil, &screens.NavigateResult{
		ResultData: map[string]interface{}{
			"action": "add",
		},
	}
}

// handleEditAction handles editing a selected skill.
func (s *SkillsListScreen) handleEditAction() (tea.Cmd, screens.ScreenResult) {
	if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
		return nil, &screens.NavigateResult{
			ResultData: map[string]interface{}{
				"action": "edit",
				"skill":  *selected,
			},
		}
	}
	return nil, nil
}

// handleDeleteAction handles deleting a selected skill.
func (s *SkillsListScreen) handleDeleteAction() (tea.Cmd, screens.ScreenResult) {
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

// RenderContent returns just the content (table) without StandardView wrapper.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *SkillsListScreen) RenderContent() string {
	return s.tableBehavior.Render()
}

// View renders the skills list screen.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *SkillsListScreen) View() string {
	// Render table via behavior.
	content := s.RenderContent()

	// Get theme for UIKit primitives (fall back to default if not set).
	var th themes.Theme
	if screenTheme := s.Theme(); screenTheme != nil {
		if t, ok := screenTheme.(themes.Theme); ok {
			th = t
		}
	}
	if th == nil {
		th = themes.NewDefaultTheme()
	}

	// Build footer using UIKit primitives for consistent styling.
	var footer string
	if len(s.skills) == 0 {
		footer = primitives.RenderHelpFooter(th,
			primitives.AddBadge(th),
			primitives.BackBadge(th),
		)
	} else {
		footer = primitives.RenderHelpFooter(th,
			primitives.NavigateBadge(th),
			primitives.HelpKeyBadge("Enter", "View", th),
			primitives.AddBadge(th),
			primitives.EditBadge(th),
			primitives.DeleteBadge(th),
			primitives.HelpKeyBadge("f", "Filter", th),
			primitives.HelpKeyBadge("s", "Sort", th),
			primitives.HelpKeyBadge("/", "Search", th),
			primitives.BackBadge(th),
		)
	}

	return s.CreateView([]string{"Main Menu", "Manage Skills"}, content, footer)
}

// SetTheme applies theme to the table behavior.
//
// Expected:
//   - interface{} must be valid.
//
// Side effects:
//   - None.
func (s *SkillsListScreen) SetTheme(theme interface{}) {
	s.Screen.SetTheme(theme)
	// Apply theme to table behavior if available
	if t, ok := theme.(themes.Theme); ok && t != nil {
		s.tableBehavior.SetTheme(t)
	}
}

// SetEventCounts sets the event counts for skills (used for displaying event count column).
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
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
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (s *SkillsListScreen) GetSelectedIndex() int {
	return s.tableBehavior.GetSelectedIndex()
}

// SetSelectedIndex sets the currently selected index.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (s *SkillsListScreen) SetSelectedIndex(index int) {
	s.tableBehavior.SetSelectedIndex(index)
}
