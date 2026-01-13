// Package skills provides screens for skill management workflows.
package skills

import (
	"fmt"
	"strconv"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
	table         table.Model
	listContainer *components.TableListContainer
	eventCounts   map[string]int // Skill ID -> event count
}

// NewSkillsListScreen creates a new skills list screen.
func NewSkillsListScreen(skills []*career.Skill) *SkillsListScreen {
	// Create table columns matching legacy format
	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
		{Title: "Years", Width: 8},
		{Title: "Events", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	// Apply default styles - theme will be applied via SetTheme
	t.SetStyles(table.DefaultStyles())

	screen := &SkillsListScreen{
		BaseScreen:    base.NewBaseScreen(),
		skills:        skills,
		selectedIndex: 0,
		table:         t,
		listContainer: components.NewTableListContainer(t, "Skills Management", 100),
		eventCounts:   make(map[string]int), // Will be populated later
	}

	// Update table rows with skills
	screen.updateTableRows()

	return screen
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

// updateTableRows updates the table rows based on skills (with pagination).
func (s *SkillsListScreen) updateTableRows() {
	pageSize := 15
	total := len(s.skills)

	// Determine which page current selection is on
	page := 0
	if pageSize > 0 && s.selectedIndex >= 0 {
		page = s.selectedIndex / pageSize
	}

	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	pageSkills := s.skills[start:end]

	rows := make([]table.Row, 0, len(pageSkills))
	for idx, skill := range pageSkills {
		realIdx := start + idx

		// Name with selection indicator
		name := skill.Name
		if realIdx == s.selectedIndex {
			name = "▶ " + name
		} else {
			name = "  " + name
		}

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
		if count, ok := s.eventCounts[skill.ID]; ok && count > 0 {
			eventCount = strconv.Itoa(count)
		}

		rows = append(rows, table.Row{name, category, level, years, eventCount})
	}

	s.table.SetRows(rows)

	// Calculate relative cursor position for this page
	relativeCursor := 0
	if s.selectedIndex >= start && s.selectedIndex < end {
		relativeCursor = s.selectedIndex - start
	}

	// Set table cursor to relative position
	s.table.SetCursor(relativeCursor)

	// Sync container's selected index
	s.listContainer.SetSelectedIdx(relativeCursor)

	// CRITICAL: Sync updated table back to container (fixes display bug)
	// The container stores a VALUE COPY of the table, so we must explicitly
	// update it after modifying rows/cursor, otherwise it renders stale data
	s.listContainer.SetTable(s.table)
}

// View renders the skills list screen.
func (s *SkillsListScreen) View() string {
	// Handle empty state
	if len(s.skills) == 0 {
		s.listContainer.SetEmptyStateMessage("No skills found. Press 'a' to add your first skill.")
		content := s.listContainer.Render()
		footer := "a: Add skill  Esc: Back"
		return s.CreateView([]string{"Main Menu", "Manage Skills"}, content, footer)
	}

	// Ensure table rows are synchronized
	s.updateTableRows()

	// Build pagination info matching legacy format
	pageSize := 15
	totalItems := len(s.skills)
	currentPage := (s.selectedIndex / pageSize) + 1
	totalPages := (totalItems + pageSize - 1) / pageSize
	paginationInfo := fmt.Sprintf("Skills: %d | Page %d of %d", totalItems, currentPage, totalPages)
	s.listContainer.SetPaginationInfo(paginationInfo)

	// Render table via container
	content := s.listContainer.Render()

	// Footer with actions (matching legacy)
	footer := "Enter: View  a: Add  e: Edit  d: Delete  ↑↓/jk: Navigate  g/G: Top/Bottom  Esc: Back"

	return s.CreateView([]string{"Main Menu", "Manage Skills"}, content, footer)
}

// SetTheme applies theme to the table (override BaseScreen).
func (s *SkillsListScreen) SetTheme(theme interface{}) {
	s.BaseScreen.SetTheme(theme)
	// Apply themed table styles if theme is available
	if t, ok := theme.(themes.Theme); ok && t != nil {
		s.table.SetStyles(themes.NewThemedTableStyles(t))
	}
}

// SetEventCounts sets the event counts for skills (used for displaying event count column).
func (s *SkillsListScreen) SetEventCounts(counts map[string]int) {
	s.eventCounts = counts
	// Refresh table rows with new counts
	s.updateTableRows()
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
