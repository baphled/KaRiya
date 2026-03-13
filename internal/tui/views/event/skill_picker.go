package event

import (
	"fmt"

	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/containers"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxSkillPickerNameLength = 30

// SkillPicker allows selecting an existing skill to link to an event.
// Follows the EventsModal pattern for consistency.
//
// Usage:
//
//	modal := NewSkillPicker(skills, theme)
//	modal.SetDimensions(width, height)
//	modal.Show()
//
//	// In Update:
//	model, cmd := modal.Update(msg)
//
//	// Check for selection:
//	if modal.HasSelection() {
//	    selectedSkill := modal.GetSelectedSkill()
//	    // Link skill to event
//	    modal.ClearSelection()
//	}
type SkillPicker struct {
	skills        []display.Skill
	theme         themes.Theme
	visible       bool
	width         int
	height        int
	table         *behaviors.TableBehavior[display.Skill]
	selectedSkill display.Skill
}

// skillPickerRowFormatter formats a skill for table display.
func skillPickerRowFormatter(skill display.Skill, _ int) []string {
	if isEmptyDisplaySkill(skill) {
		return []string{"-", "-", "-"}
	}

	name := skill.Name
	if name == "" {
		name = "(Unnamed)"
	}
	if len(name) > maxSkillPickerNameLength {
		name = name[:maxSkillPickerNameLength] + "..."
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

// NewSkillPicker creates a new skill picker modal.
//
// Expected:
//   - skills must be a valid slice of career.Skill pointers.
//   - theme must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized SkillPicker ready for use.
//
// Side effects:
//   - None.
func NewSkillPicker(skills []display.Skill, theme themes.Theme) *SkillPicker {
	filteredSkills := make([]display.Skill, 0, len(skills))
	for i := range skills {
		skill := skills[i]
		if isEmptyDisplaySkill(skill) {
			continue
		}
		filteredSkills = append(filteredSkills, skill)
	}

	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 30},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
	}

	tableBehavior := behaviors.NewTableBehavior(theme, columns, skillPickerRowFormatter).
		PageSize(12).
		EmptyMessage("No skills available.").
		HidePagination()

	tableBehavior.SetItems(filteredSkills)

	m := &SkillPicker{
		skills:  filteredSkills,
		theme:   theme,
		visible: false,
		width:   100,
		height:  24,
		table:   tableBehavior,
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
func (m *SkillPicker) Init() tea.Cmd {
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
//   - May set selected skill.
func (m *SkillPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case tea.KeyEsc, tea.KeyBackspace:
			m.Hide()
			return m, nil

		case tea.KeyEnter:
			if selected := m.table.GetSelectedItem(); selected != nil {
				m.selectedSkill = *selected
				m.Hide()
			}
			return m, nil
		}

		if msg.String() == "q" {
			m.Hide()
			return m, nil
		}

		if m.table.HandleNavigation(msg.String()) {
			return m, nil
		}
	}

	return m, nil
}

// updateTableDimensions updates the table dimensions based on modal size.
func (m *SkillPicker) updateTableDimensions() {
	tableHeight := m.height - 14
	if tableHeight < 5 {
		tableHeight = 5
	}
	if tableHeight > 12 {
		tableHeight = 12
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
func (m *SkillPicker) View() string {
	if !m.visible {
		return ""
	}

	theme := m.theme
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	maxModalHeight := 24
	terminalMaxHeight := int(float64(m.height) * 0.85)
	if terminalMaxHeight > maxModalHeight {
		maxModalHeight = terminalMaxHeight
	}
	if maxModalHeight < 16 {
		maxModalHeight = 16
	}

	modalWidth := m.width - 8
	if modalWidth > 75 {
		modalWidth = 75
	}
	if modalWidth < 60 {
		modalWidth = 60
	}

	title := primitives.Title(fmt.Sprintf("Select Skill (%d available)", len(m.skills)), theme).Render()

	var content string
	if m.table.IsEmpty() {
		content = primitives.Muted("No skills available.", theme).
			Italic().
			MarginTop(2).
			MarginBottom(2).
			Render()
	} else {
		content = m.table.Render()
	}

	var footerText string
	if !m.table.IsEmpty() {
		footerText = fmt.Sprintf("Enter: Select | ↑↓/j/k: Navigate | Esc: Cancel  [%d/%d]",
			m.table.GetSelectedIndex()+1, m.table.Count())
	} else {
		footerText = "Esc: Cancel"
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

// SetDimensions updates the modal's available dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *SkillPicker) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	m.updateTableDimensions()
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *SkillPicker) Show() {
	m.visible = true
	m.selectedSkill = display.Skill{}
	m.table.SetSelectedIndex(0)
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *SkillPicker) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *SkillPicker) IsVisible() bool {
	return m.visible
}

// SetSkills updates the skills being displayed.
//
// Expected:
//   - skill must be valid.
//
// Side effects:
//   - None.
func (m *SkillPicker) SetSkills(skills []display.Skill) {
	filteredSkills := make([]display.Skill, 0, len(skills))
	for i := range skills {
		skill := skills[i]
		if isEmptyDisplaySkill(skill) {
			continue
		}
		filteredSkills = append(filteredSkills, skill)
	}

	m.skills = filteredSkills
	m.selectedSkill = display.Skill{}
	m.table.SetItems(filteredSkills)
}

// HasSelection returns true if the user selected a skill.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *SkillPicker) HasSelection() bool {
	return !isEmptyDisplaySkill(m.selectedSkill)
}

// GetSelectedSkill returns the selected skill (nil if none selected).
//
// Returns:
//   - A display.Skill value.
//
// Side effects:
//   - None.
func (m *SkillPicker) GetSelectedSkill() display.Skill {
	return m.selectedSkill
}

// ClearSelection clears any previous selection.
//
// Side effects:
//   - None.
func (m *SkillPicker) ClearSelection() {
	m.selectedSkill = display.Skill{}
}
