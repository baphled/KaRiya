package models

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
)

// CVConfigManagerModel manages the display and interaction with saved CV configurations.
type CVConfigManagerModel struct {
	*BaseStandardModel
	configManager cv.ConfigManager
	configs       []*career.CVConfig
	selectedIdx   int
	loading       bool
	header        string
	footer        string
	selectedEvent *career.CareerEvent // Selected event for CV generation context
}

// NewCVConfigManagerModel creates a new CV Config Manager model.
func NewCVConfigManagerModel(
	baseModel *BaseStandardModel,
	configManager cv.ConfigManager,
) *CVConfigManagerModel {
	return &CVConfigManagerModel{
		BaseStandardModel: baseModel,
		configManager:     configManager,
		configs:           make([]*career.CVConfig, 0),
		selectedIdx:       0,
		loading:           true,
		header:            "CV Configuration Manager",
		footer:            "",
		selectedEvent:     nil,
	}
}

// NewCVConfigManagerModelWithEvent creates a new CV Config Manager model with a selected event.
func NewCVConfigManagerModelWithEvent(
	baseModel *BaseStandardModel,
	configManager cv.ConfigManager,
	event *career.CareerEvent,
) *CVConfigManagerModel {
	return &CVConfigManagerModel{
		BaseStandardModel: baseModel,
		configManager:     configManager,
		configs:           make([]*career.CVConfig, 0),
		selectedIdx:       0,
		loading:           true,
		header:            "CV Configuration Manager",
		footer:            "",
		selectedEvent:     event,
	}
}

// Init initializes the model and loads configurations.
func (m *CVConfigManagerModel) Init() tea.Cmd {
	return tea.Batch(
		m.BaseStandardModel.Init(),
		m.loadConfigs(),
	)
}

// loadConfigs loads all CV configurations from the manager.
func (m *CVConfigManagerModel) loadConfigs() tea.Cmd {
	return func() tea.Msg {
		configs, err := m.configManager.ListConfigs(context.Background())
		if err != nil {
			return ConfigLoadError{err: err}
		}
		return ConfigsLoadedMsg{configs: configs}
	}
}

// Update handles messages and updates the model state.
func (m *CVConfigManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConfigsLoadedMsg:
		m.loading = false
		m.configs = msg.configs
		if len(m.configs) > 0 {
			m.selectedIdx = 0
		}
		return m, nil

	case ConfigLoadError:
		m.loading = false
		m.SetError(msg.err)
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil

		case "down", "j":
			if m.selectedIdx < len(m.configs)-1 {
				m.selectedIdx++
			}
			return m, nil

		case "enter":
			if len(m.configs) > 0 {
				selectedConfig := m.configs[m.selectedIdx]
				return m, func() tea.Msg {
					return GenerateCVFromConfigMsg{Config: selectedConfig}
				}
			}
			return m, nil

		case "n":
			// New config - navigate to editor
			return m, func() tea.Msg {
				return NavigateToScreenMsg{screenID: "cv_config_editor"}
			}

		case "e":
			// Edit config
			if len(m.configs) > 0 {
				selectedConfig := m.configs[m.selectedIdx]
				return m, func() tea.Msg {
					return EditCVConfigMsg{config: selectedConfig}
				}
			}
			return m, nil

		case "d":
			// Delete config (with confirmation)
			if len(m.configs) > 0 {
				selectedConfig := m.configs[m.selectedIdx]
				return m, func() tea.Msg {
					return ConfirmDeleteCVConfigMsg{config: selectedConfig}
				}
			}
			return m, nil

		case "esc", "q":
			return m, func() tea.Msg {
				return BackToMainMenuMsg{}
			}
		}
	}

	// Return model unchanged for unhandled messages
	return m, nil
}

// View renders the CV Configuration Manager screen.
func (m *CVConfigManagerModel) View() string {
	if m.loading {
		return fmt.Sprintf("%s\n\nLoading configurations...", m.header)
	}

	if m.GetLastError() != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)
		return fmt.Sprintf("%s\n\n%s\n\nError: %v",
			m.header,
			m.renderConfigList(),
			errorStyle.Render(m.GetLastError().Error()))
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s",
		m.header,
		m.renderConfigList(),
		m.renderFooter())
}

// renderConfigList renders the configuration list as a table.
func (m *CVConfigManagerModel) renderConfigList() string {
	if len(m.configs) == 0 {
		return "No CV configurations found.\nPress 'n' to create a new configuration."
	}

	var output string
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6"))

	output += headerStyle.Render("Name") + " | " +
		headerStyle.Render("Role") + " | " +
		headerStyle.Render("Audiences") + " | " +
		headerStyle.Render("Updated") + "\n"
	output += lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render("─────────────────────────────────────────────────────────────────\n")

	for i, config := range m.configs {
		selectedStyle := lipgloss.NewStyle()
		if i == m.selectedIdx {
			selectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("4")).
				Foreground(lipgloss.Color("15"))
		}

		audiences := ""
		if len(config.TargetAudience) > 0 {
			audiences = config.TargetAudience[0]
			if len(config.TargetAudience) > 1 {
				audiences = fmt.Sprintf("%d audiences", len(config.TargetAudience))
			}
		}

		line := fmt.Sprintf("  %-25s | %-15s | %-20s | %s\n",
			config.Name,
			config.TargetRole,
			audiences,
			config.UpdatedAt.Format("2006-01-02"))

		output += selectedStyle.Render(line)
	}

	return output
}

// renderFooter renders the footer with keyboard shortcuts.
func (m *CVConfigManagerModel) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	shortcuts := []string{
		"↑/k: Up",
		"↓/j: Down",
		"Enter: Generate",
		"n: New",
		"e: Edit",
		"d: Delete",
		"q: Back",
	}

	footer := ""
	for _, shortcut := range shortcuts {
		footer += footerStyle.Render(shortcut) + "  "
	}

	return footer
}

// GetSelectedConfig returns the currently selected configuration.
func (m *CVConfigManagerModel) GetSelectedConfig() *career.CVConfig {
	if len(m.configs) > 0 && m.selectedIdx >= 0 && m.selectedIdx < len(m.configs) {
		return m.configs[m.selectedIdx]
	}
	return nil
}

// RefreshConfigs reloads the configurations from the manager.
func (m *CVConfigManagerModel) RefreshConfigs() tea.Cmd {
	m.loading = true
	m.selectedIdx = 0
	return m.loadConfigs()
}

// Messages for CV Config Manager

// ConfigsLoadedMsg is sent when configurations have been successfully loaded.
type ConfigsLoadedMsg struct {
	configs []*career.CVConfig
}

// ConfigLoadError is sent when there's an error loading configurations.
type ConfigLoadError struct {
	err error
}

// GenerateCVFromConfigMsg triggers CV generation from a selected configuration.
type GenerateCVFromConfigMsg struct {
	Config *career.CVConfig
}

// EditCVConfigMsg triggers editing of a CV configuration.
type EditCVConfigMsg struct {
	config *career.CVConfig
}

// ConfirmDeleteCVConfigMsg triggers confirmation dialog for deleting a configuration.
type ConfirmDeleteCVConfigMsg struct {
	config *career.CVConfig
}

// BackToMainMenuMsg navigates back to the main menu.
type BackToMainMenuMsg struct{}

// NavigateToScreenMsg navigates to a specific screen by ID.
type NavigateToScreenMsg struct {
	screenID string
}

