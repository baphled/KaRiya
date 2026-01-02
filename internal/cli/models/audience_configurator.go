package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
)

// AudienceConfiguratorModel handles audience selection for CV configuration.
type AudienceConfiguratorModel struct {
	*BaseStandardModel
	audiences           []string
	selected            map[string]bool
	focusIdx            int
	width               int
	height              int
	header              components.HeaderModel
	helpFooter          components.HelpFooterModel
	onAudiencesSelected func([]string)
	onBack              func()
}

// AudienceDescription provides user-friendly descriptions for each audience.
var AudienceDescription = map[string]string{
	"hiring_manager": "Hiring Manager - Hiring managers evaluating for the role",
	"recruiter":      "Recruiter - Recruiters evaluating general fit and skills",
	"peer":           "Peer - Technical peers reviewing technical depth",
}

// NewAudienceConfiguratorModel creates a new audience configurator model.
func NewAudienceConfiguratorModel(
	baseModel *BaseStandardModel,
	selectedAudiences []string,
) *AudienceConfiguratorModel {
	selected := make(map[string]bool)
	for _, a := range selectedAudiences {
		selected[a] = true
	}

	return &AudienceConfiguratorModel{
		BaseStandardModel:   baseModel,
		audiences:           []string{"hiring_manager", "recruiter", "peer"},
		selected:            selected,
		focusIdx:            0,
		width:               80,
		height:              20,
		header:              components.NewHeader("Select Target Audiences", 80),
		helpFooter:          components.NewHelpFooter("audience_configurator", 80),
		onAudiencesSelected: nil,
		onBack:              nil,
	}
}

// Init initializes the audience configurator model.
func (m *AudienceConfiguratorModel) Init() tea.Cmd {
	return m.BaseStandardModel.Init()
}

// Update handles messages and updates the model state.
func (m *AudienceConfiguratorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.focusIdx > 0 {
				m.focusIdx--
			}
			return m, nil

		case tea.KeyDown:
			if m.focusIdx < len(m.audiences)-1 {
				m.focusIdx++
			}
			return m, nil

		case tea.KeySpace:
			// Toggle selection of current item
			audience := m.audiences[m.focusIdx]
			m.selected[audience] = !m.selected[audience]
			m.ClearError()
			return m, nil

		case tea.KeyEnter:
			// Confirm selection - collect selected audiences
			var selectedAudiences []string
			for _, a := range m.audiences {
				if m.selected[a] {
					selectedAudiences = append(selectedAudiences, a)
				}
			}

			// Must select at least one audience
			if len(selectedAudiences) == 0 {
				m.SetError(fmt.Errorf("must select at least one audience"))
				return m, nil
			}

			if m.onAudiencesSelected != nil {
				m.onAudiencesSelected(selectedAudiences)
			}

			return m, func() tea.Msg {
				return AudiencesSelectedMsg{Audiences: selectedAudiences}
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
				case '\r':
					// Handle enter key
					var selectedAudiences []string
					for _, a := range m.audiences {
						if m.selected[a] {
							selectedAudiences = append(selectedAudiences, a)
						}
					}
					if len(selectedAudiences) == 0 {
						m.SetError(fmt.Errorf("must select at least one audience"))
						return m, nil
					}
					if m.onAudiencesSelected != nil {
						m.onAudiencesSelected(selectedAudiences)
					}
					return m, func() tea.Msg {
						return AudiencesSelectedMsg{Audiences: selectedAudiences}
					}
				case 'j':
					if m.focusIdx < len(m.audiences)-1 {
						m.focusIdx++
					}
					return m, nil
				case 'k':
					if m.focusIdx > 0 {
						m.focusIdx--
					}
					return m, nil
				}
			}
		}
	}

	// Return model unchanged for unhandled messages
	return m, nil
}

// View renders the audience configurator screen.
func (m *AudienceConfiguratorModel) View() string {
	headerContent := m.header.View()

	// Render audience list
	var audienceItems []string
	for i, audience := range m.audiences {
		checkbox := "[ ]"
		if m.selected[audience] {
			checkbox = styles.ListItemSelected.Render("[✓]")
		}

		indicator := "  "
		if i == m.focusIdx {
			indicator = styles.ListItemSelected.Render("→")
		}

		desc := AudienceDescription[audience]
		if i == m.focusIdx {
			audienceItems = append(audienceItems, fmt.Sprintf(
				"%s %s %s",
				indicator,
				checkbox,
				styles.ListItemSelected.Render(desc),
			))
		} else {
			audienceItems = append(audienceItems, fmt.Sprintf("%s %s %s", indicator, checkbox, desc))
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		audienceItems...,
	)

	// Add error message if present
	var fullContent string
	if m.GetLastError() != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)
		errorMsg := errorStyle.Render(fmt.Sprintf("Error: %v", m.GetLastError()))
		fullContent = lipgloss.JoinVertical(
			lipgloss.Top,
			content,
			"",
			errorMsg,
		)
	} else {
		fullContent = content
	}

	// Add padding and styling
	contentBox := lipgloss.NewStyle().
		Padding(1, 2).
		Render(fullContent)

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

// GetSelectedAudiences returns the currently selected audiences.
func (m *AudienceConfiguratorModel) GetSelectedAudiences() []string {
	var result []string
	for _, a := range m.audiences {
		if m.selected[a] {
			result = append(result, a)
		}
	}
	return result
}

// GetFocusedAudience returns the currently focused audience.
func (m *AudienceConfiguratorModel) GetFocusedAudience() string {
	if m.focusIdx < 0 || m.focusIdx >= len(m.audiences) {
		return ""
	}
	return m.audiences[m.focusIdx]
}

// IsSelected returns whether an audience is selected.
func (m *AudienceConfiguratorModel) IsSelected(audience string) bool {
	return m.selected[audience]
}

// GetWidth returns the model width.
func (m *AudienceConfiguratorModel) GetWidth() int {
	return m.width
}

// GetHeight returns the model height.
func (m *AudienceConfiguratorModel) GetHeight() int {
	return m.height
}

// SetOnAudiencesSelected sets the callback for audiences selection.
func (m *AudienceConfiguratorModel) SetOnAudiencesSelected(fn func([]string)) {
	m.onAudiencesSelected = fn
}

// SetOnBack sets the callback for back navigation.
func (m *AudienceConfiguratorModel) SetOnBack(fn func()) {
	m.onBack = fn
}

// AudiencesSelectedMsg indicates audiences have been selected.
type AudiencesSelectedMsg struct {
	Audiences []string
}
