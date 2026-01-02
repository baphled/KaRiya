package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
)

// RoleSelectorModel handles role selection for CV configuration.
type RoleSelectorModel struct {
	*BaseStandardModel
	roles          []string
	selectedIdx    int
	width          int
	height         int
	header         components.HeaderModel
	helpFooter     components.HelpFooterModel
	onRoleSelected func(string)
	onBack         func()
}

// RoleDescription provides user-friendly descriptions for each role.
var RoleDescription = map[string]string{
	"principal": "Principal - Senior leadership and strategic influence",
	"staff":     "Staff - Broad technical impact and influence",
	"em":        "EM - Engineering Management and team leadership",
	"senior_ic": "Senior IC - Senior Individual Contributor",
}

// NewRoleSelectorModel creates a new role selector model.
func NewRoleSelectorModel(baseModel *BaseStandardModel) *RoleSelectorModel {
	return &RoleSelectorModel{
		BaseStandardModel: baseModel,
		roles:             []string{"principal", "staff", "em", "senior_ic"},
		selectedIdx:       0,
		width:             80,
		height:            20,
		header:            components.NewHeader("Select Target Role", 80),
		helpFooter:        components.NewHelpFooter("role_selector", 80),
		onRoleSelected:    nil,
		onBack:            nil,
	}
}

// Init initializes the role selector model.
func (m *RoleSelectorModel) Init() tea.Cmd {
	return m.BaseStandardModel.Init()
}

// Update handles messages and updates the model state.
func (m *RoleSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil

		case tea.KeyDown:
			if m.selectedIdx < len(m.roles)-1 {
				m.selectedIdx++
			}
			return m, nil

		case tea.KeyEnter:
			selectedRole := m.roles[m.selectedIdx]
			if m.onRoleSelected != nil {
				m.onRoleSelected(selectedRole)
			}
			return m, func() tea.Msg {
				return RoleSelectedMsg{Role: selectedRole}
			}

		case tea.KeyEscape:
			if m.onBack != nil {
				m.onBack()
			}
			return m, func() tea.Msg {
				return BackMsg{}
			}

		case tea.KeyRunes:
			// Handle letter keys for navigation or enter key
			if len(msg.Runes) > 0 {
				switch msg.Runes[0] {
				case '':
					// Handle enter key
					selectedRole := m.roles[m.selectedIdx]
					if m.onRoleSelected != nil {
						m.onRoleSelected(selectedRole)
					}
					return m, func() tea.Msg {
						return RoleSelectedMsg{Role: selectedRole}
					}
				case 'j':
					if m.selectedIdx < len(m.roles)-1 {
						m.selectedIdx++
					}
					return m, nil
				case 'k':
					if m.selectedIdx > 0 {
						m.selectedIdx--
					}
					return m, nil
				}
			}
		}
	}

	// Return model unchanged for unhandled messages
	return m, nil
}

// View renders the role selector screen.
func (m *RoleSelectorModel) View() string {
	headerContent := m.header.View()

	// Render role list
	var roleItems []string
	for i, role := range m.roles {
		indicator := "  "
		if i == m.selectedIdx {
			indicator = styles.ListItemSelected.Render("→")
		}

		desc := RoleDescription[role]
		if i == m.selectedIdx {
			roleItems = append(roleItems, fmt.Sprintf(
				"%s %s",
				indicator,
				styles.ListItemSelected.Render(desc),
			))
		} else {
			roleItems = append(roleItems, fmt.Sprintf("%s %s", indicator, desc))
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		roleItems...,
	)

	// Add padding and styling
	contentBox := lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)

	footerContent := m.helpFooter.View()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		headerContent,
		"",
		contentBox,
		"",
		footerContent,
	)
}

// GetSelectedRole returns the currently selected role.
func (m *RoleSelectorModel) GetSelectedRole() string {
	if m.selectedIdx < 0 || m.selectedIdx >= len(m.roles) {
		return ""
	}
	return m.roles[m.selectedIdx]
}

// GetWidth returns the model width.
func (m *RoleSelectorModel) GetWidth() int {
	return m.width
}

// GetHeight returns the model height.
func (m *RoleSelectorModel) GetHeight() int {
	return m.height
}

// SetOnRoleSelected sets the callback for role selection.
func (m *RoleSelectorModel) SetOnRoleSelected(fn func(string)) {
	m.onRoleSelected = fn
}

// SetOnBack sets the callback for back navigation.
func (m *RoleSelectorModel) SetOnBack(fn func()) {
	m.onBack = fn
}

// RoleSelectedMsg indicates a role has been selected.
type RoleSelectedMsg struct {
	Role string
}
