package models

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
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
	headerModel   components.HeaderModel
	footer        string
	selectedEvent *career.CareerEvent // Selected event for CV generation context
}

// NewCVConfigManagerModel creates a new CV Config Manager model.
func NewCVConfigManagerModel(
	baseModel *BaseStandardModel,
	configManager cv.ConfigManager,
) *CVConfigManagerModel {
	m := &CVConfigManagerModel{
		BaseStandardModel: baseModel,
		configManager:     configManager,
		configs:           make([]*career.CVConfig, 0),
		selectedIdx:       0,
		loading:           true,
		header:            "CV Configuration Manager",
		headerModel:       components.NewHeader("CV Configuration Manager", 80),
		footer:            "",
		selectedEvent:     nil,
	}
	m.headerModel.SetBreadcrumbs([]string{"Home", "CV Management", "Configurations"})
	return m
}

// NewCVConfigManagerModelWithEvent creates a new CV Config Manager model with a selected event.
func NewCVConfigManagerModelWithEvent(
	baseModel *BaseStandardModel,
	configManager cv.ConfigManager,
	event *career.CareerEvent,
) *CVConfigManagerModel {
	m := &CVConfigManagerModel{
		BaseStandardModel: baseModel,
		configManager:     configManager,
		configs:           make([]*career.CVConfig, 0),
		selectedIdx:       0,
		loading:           true,
		header:            "CV Configuration Manager",
		headerModel:       components.NewHeader("CV Configuration Manager", 80),
		footer:            "",
		selectedEvent:     event,
	}
	m.headerModel.SetBreadcrumbs([]string{"Home", "Events", "CV Generation"})
	return m
}

// Init initializes the model and loads configurations.
func (m *CVConfigManagerModel) Init() tea.Cmd {
	// Don't batch with nil - just return the load command
	return m.loadConfigs()
}

// loadConfigs loads all CV configurations from the manager.
func (m *CVConfigManagerModel) loadConfigs() tea.Cmd {
	return func() tea.Msg {
		// Create a channel to handle timeout
		type result struct {
			configs []*career.CVConfig
			err     error
		}
		
		resultChan := make(chan result, 1)
		
		// Load configs in a goroutine
		go func() {
			configs, err := m.configManager.ListConfigs(context.Background())
			resultChan <- result{configs, err}
		}()
		
		// Wait for result with 3-second timeout
		select {
		case res := <-resultChan:
			if res.err != nil {
				return ConfigLoadError{err: res.err}
			}
			return ConfigsLoadedMsg{configs: res.configs}
		case <-time.After(3 * time.Second):
			// Timeout - return empty list instead of error to allow UI to continue
			return ConfigsLoadedMsg{configs: []*career.CVConfig{}}
		}
	}
}

// Update handles messages and updates the model state.
func (m *CVConfigManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.headerModel.SetWidth(msg.Width)
		return m, nil

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
		// Allow quit/escape even while loading
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg {
				return BackToMainMenuMsg{}
			}
		}

		// If still loading, don't process other keys
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

		case "r":
			// Retry loading configs
			if m.GetLastError() != nil {
				m.ClearError()
				m.loading = true
				return m, m.loadConfigs()
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
	headerContent := m.headerModel.View()
	
	// Always show something useful, even if loading
	if m.loading {
		return fmt.Sprintf("%s\n\nLoading configurations...\n\nPress 'q' or 'esc' to cancel", headerContent)
	}

	if m.GetLastError() != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)
		return fmt.Sprintf("%s\n\n%s\n\nError: %v\n\nPress 'r' to retry or 'q' to go back",
			headerContent,
			m.renderConfigList(),
			errorStyle.Render(m.GetLastError().Error()))
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s",
		headerContent,
		m.renderConfigList(),
		m.renderFooter())
}

// renderConfigList renders the configuration list as a table.
func (m *CVConfigManagerModel) renderConfigList() string {
	if len(m.configs) == 0 {
		msg := "No CV configurations found.\n\nYou can:"
		if m.GetLastError() != nil {
			msg += "\n  • Press 'r' to retry loading"
		}
		msg += "\n  • Press 'n' to create a new configuration"
		msg += "\n  • Press 'q' to go back"
		return msg
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

// ConfigLoadTimeoutMsg is sent when config loading times out
type ConfigLoadTimeoutMsg struct{}
