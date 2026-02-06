package modals

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SkillAction represents the action to perform on skills.
type SkillAction int

const (
	// ActionNone represents no action.
	ActionNone SkillAction = iota
	// ActionAddExisting represents adding an existing skill.
	ActionAddExisting
	// ActionAddNew represents adding a new skill.
	ActionAddNew
	// ActionInfer represents inferring skills from event text.
	ActionInfer
	// ActionRemove represents removing a skill from the event.
	ActionRemove
)

// maxSkillNameLength defines the maximum length for skill names before truncation.
const maxSkillNameLength = 25

// SkillsDetailModal displays skills associated with an event in an interactive table.
// This modal uses TableBehavior for consistent table display and navigation.
//
// Features:
// - Data table display with Name, Category, and Level columns
// - Standard table navigation (j/k, up/down, page up/down)
// - Action keys: a (add existing), n (add new), i (infer), d (remove)
// - Close with Escape or 'q'
//
// Usage:
//
//	modal := NewSkillsDetailModal(eventID, skills, theme)
//	modal.SetDimensions(width, height)
//	modal.Show()
//
//	// In Update:
//	model, cmd := modal.Update(msg)
//
//	// Check for actions:
//	if modal.GetAction() == ActionAddExisting {
//	    // Open skill picker modal
//	    modal.ClearAction()
//	}
//
//	// In View:
//	if modal.IsVisible() {
//	    return renderModalOverlay(modal, background)
//	}
type SkillsDetailModal struct {
	eventID       string
	skills        []*career.Skill
	theme         themes.Theme
	visible       bool
	width         int
	height        int
	table         *behaviors.TableBehavior[*career.Skill]
	action        SkillAction
	selectedSkill *career.Skill
}

// skillRowFormatter formats a skill for table display.
func skillRowFormatter(skill *career.Skill, _ int) []string {
	if skill == nil {
		return []string{"-", "-", "-"}
	}

	name := skill.Name
	if name == "" {
		name = "(Unnamed)"
	}
	if len(name) > maxSkillNameLength {
		name = name[:maxSkillNameLength] + "..."
	}

	category := skill.Category
	if category == "" {
		category = "-"
	}

	level := skill.Level
	if level == "" {
		level = "-"
	}

	return []string{name, category, level}
}

// NewSkillsDetailModal creates a new skills detail modal.
//
// Expected:
//   - Must be a valid string.
//   - skill must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized SkillsDetailModal ready for use.
//
// Side effects:
//   - None.
func NewSkillsDetailModal(eventID string, skills []*career.Skill, theme themes.Theme) *SkillsDetailModal {
	filteredSkills := make([]*career.Skill, 0, len(skills))
	for _, s := range skills {
		if s != nil {
			filteredSkills = append(filteredSkills, s)
		}
	}

	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
	}

	tableBehavior := behaviors.NewTableBehavior(theme, columns, skillRowFormatter).
		PageSize(10).
		EmptyMessage("No skills associated with this event.").
		HidePagination()

	tableBehavior.SetItems(filteredSkills)

	m := &SkillsDetailModal{
		eventID: eventID,
		skills:  filteredSkills,
		theme:   theme,
		visible: false,
		width:   100,
		height:  24,
		table:   tableBehavior,
		action:  ActionNone,
	}
	return m
}

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) Init() tea.Cmd {
	return nil
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
//   - May update dimensions.
//   - May hide modal.
//   - May set action state.
func (m *SkillsDetailModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateTableDimensions()
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.Hide()
			return m, nil
		}

		switch msg.String() {
		case "q":
			m.Hide()
			return m, nil
		case "a":
			m.action = ActionAddExisting
			return m, nil
		case "n":
			m.action = ActionAddNew
			return m, nil
		case "i":
			m.action = ActionInfer
			return m, nil
		case "d":
			if selected := m.table.GetSelectedItem(); selected != nil {
				m.selectedSkill = *selected
				m.action = ActionRemove
			}
			return m, nil
		}

		if m.table.HandleNavigation(msg.String()) {
			return m, nil
		}
	}

	return m, nil
}

// updateTableDimensions updates the table dimensions based on modal size.
func (m *SkillsDetailModal) updateTableDimensions() {
	tableHeight := m.height - 14
	if tableHeight < 5 {
		tableHeight = 5
	}
	if tableHeight > 10 {
		tableHeight = 10
	}
	m.table.PageSize(tableHeight)
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) View() string {
	if !m.visible {
		return ""
	}

	theme := m.theme
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	maxModalHeight := 20
	terminalMaxHeight := int(float64(m.height) * 0.75)
	if terminalMaxHeight > maxModalHeight {
		maxModalHeight = terminalMaxHeight
	}
	if maxModalHeight < 16 {
		maxModalHeight = 16
	}

	modalWidth := m.width - 8
	if modalWidth > 70 {
		modalWidth = 70
	}
	if modalWidth < 55 {
		modalWidth = 55
	}

	title := primitives.Title(fmt.Sprintf("Skills (%d)", len(m.skills)), theme).Render()

	var content string
	if m.table.IsEmpty() {
		content = primitives.Muted("No skills associated with this event.", theme).
			Italic().
			MarginTop(2).
			MarginBottom(2).
			Render()
	} else {
		content = m.table.Render()
	}

	var footerText string
	if !m.table.IsEmpty() {
		footerText = fmt.Sprintf("a: Add Existing | n: New Skill | i: Infer | d: Remove | Esc: Close  [%d/%d]",
			m.table.GetSelectedIndex()+1, m.table.Count())
	} else {
		footerText = "a: Add Existing | n: New Skill | i: Infer | Esc: Close"
	}
	footer := primitives.Muted(footerText, theme).Render()

	modalContent := lipgloss.JoinVertical(lipgloss.Left, title, "", content, "", footer)

	return containers.NewBox(theme).
		Content(modalContent).
		Width(modalWidth).
		MaxHeight(maxModalHeight).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) Show() {
	m.visible = true
	m.action = ActionNone
	m.selectedSkill = nil
	m.table.SetSelectedIndex(0)
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) Hide() {
	m.visible = false
}

// SetDimensions sets the terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	m.updateTableDimensions()
}

// SetSkills updates the skills being displayed.
//
// Expected:
//   - skills must be a valid slice of career.Skill pointers.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) SetSkills(skills []*career.Skill) {
	filteredSkills := make([]*career.Skill, 0, len(skills))
	for _, s := range skills {
		if s != nil {
			filteredSkills = append(filteredSkills, s)
		}
	}

	m.skills = filteredSkills
	m.action = ActionNone
	m.selectedSkill = nil
	m.table.SetItems(filteredSkills)
}

// GetEventID returns the event ID this modal is showing skills for.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) GetEventID() string {
	return m.eventID
}

// GetAction returns the current action state.
//
// Returns:
//   - A SkillAction value.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) GetAction() SkillAction {
	return m.action
}

// ClearAction clears the action state.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) ClearAction() {
	m.action = ActionNone
	m.selectedSkill = nil
}

// GetSelectedSkill returns the currently selected skill in the table.
//
// Returns:
//   - A fully initialized career.Skill ready for use (nil if empty).
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) GetSelectedSkill() *career.Skill {
	if m.selectedSkill != nil {
		return m.selectedSkill
	}
	if selected := m.table.GetSelectedItem(); selected != nil {
		return *selected
	}
	return nil
}

// GetSelectedIndex returns the current selection index.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *SkillsDetailModal) GetSelectedIndex() int {
	return m.table.GetSelectedIndex()
}
