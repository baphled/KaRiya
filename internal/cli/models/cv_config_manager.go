package models

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
)

// CVConfigManagerModel manages the display and interaction with saved CV configurations.
// It follows the same pattern as other list-based views with table, pagination, and deletion handling.
type CVConfigManagerModel struct {
	*BaseStandardModel
	configManager cv.ConfigManager
	configs       []*career.CVConfig
	filtered      []*career.CVConfig
	table         table.Model
	listContainer *components.TableListContainer
	pagination    *PaginationHelper
	loading       bool
	header        components.HeaderModel
	helpFooter    components.HelpFooterModel
	selectedEvent *career.CareerEvent // Selected event for CV generation context
	width         int
	height        int
	breadcrumbs   []string
	sortBy        string
	sortOrder     string
	deletionState *ListDeletionState
}

// NewCVConfigManagerModel creates a new CV Config Manager model.
func NewCVConfigManagerModel(
	baseModel *BaseStandardModel,
	configManager cv.ConfigManager,
) *CVConfigManagerModel {
	// Create table with CV config columns
	// These are CONFIGURATION TEMPLATES - not generated CVs
	columns := []table.Column{
		{Title: "Template Name", Width: 30},
		{Title: "Target Role", Width: 20},
		{Title: "Audiences", Width: 25},
		{Title: "Created", Width: 12},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	// Apply styling to table
	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(styles.ColorAccentTeal).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(styles.ColorAccentTeal)
	s.Selected = s.Selected.
		Foreground(styles.ColorAccentTeal).
		Background(styles.ColorBackground).
		Bold(true)
	t.SetStyles(s)

	m := &CVConfigManagerModel{
		BaseStandardModel: baseModel,
		configManager:     configManager,
		configs:           make([]*career.CVConfig, 0),
		filtered:          make([]*career.CVConfig, 0),
		table:             t,
		listContainer:     components.NewTableListContainer(t, "CV Configuration Templates", 80),
		pagination:        NewPaginationHelper(10),
		loading:           true,
		header:            components.NewHeader("📋 CV Configuration Templates", 80),
		helpFooter:        components.NewHelpFooter("cv_config_manager", 80),
		selectedEvent:     nil,
		width:             80,
		height:            20,
		breadcrumbs:       []string{"Home", "CV Management", "Configuration Templates"},
		sortBy:            "updated",
		sortOrder:         "desc",
		deletionState:     NewListDeletionState(),
	}
	m.header.SetBreadcrumbs(m.breadcrumbs)
	m.listContainer.SetBreadcrumbs(m.breadcrumbs)
	return m
}

// NewCVConfigManagerModelWithEvent creates a new CV Config Manager model with a selected event.
// This displays CV CONFIGURATION TEMPLATES - not generated CVs.
func NewCVConfigManagerModelWithEvent(
	baseModel *BaseStandardModel,
	configManager cv.ConfigManager,
	event *career.CareerEvent,
) *CVConfigManagerModel {
	// Create table with CV config columns
	// These are CONFIGURATION TEMPLATES - not generated CVs
	columns := []table.Column{
		{Title: "Template Name", Width: 30},
		{Title: "Target Role", Width: 20},
		{Title: "Audiences", Width: 25},
		{Title: "Created", Width: 12},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	// Apply styling to table
	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(styles.ColorAccentTeal).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(styles.ColorAccentTeal)
	s.Selected = s.Selected.
		Foreground(styles.ColorAccentTeal).
		Background(styles.ColorBackground).
		Bold(true)
	t.SetStyles(s)

	m := &CVConfigManagerModel{
		BaseStandardModel: baseModel,
		configManager:     configManager,
		configs:           make([]*career.CVConfig, 0),
		filtered:          make([]*career.CVConfig, 0),
		table:             t,
		listContainer:     components.NewTableListContainer(t, "CV Configuration Templates", 80),
		pagination:        NewPaginationHelper(10),
		loading:           true,
		header:            components.NewHeader("📋 CV Configuration Templates", 80),
		helpFooter:        components.NewHelpFooter("cv_config_manager", 80),
		selectedEvent:     event,
		width:             80,
		height:            20,
		breadcrumbs:       []string{"Home", "Events", "CV Generation"},
		sortBy:            "updated",
		sortOrder:         "desc",
		deletionState:     NewListDeletionState(),
	}
	m.header.SetBreadcrumbs(m.breadcrumbs)
	m.listContainer.SetBreadcrumbs(m.breadcrumbs)
	return m
}

// Init initializes the model and loads configurations.
func (m *CVConfigManagerModel) Init() tea.Cmd {
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
	// Handle deletion confirmation if active
	if m.deletionState.IsConfirming() {
		cmd := m.deletionState.UpdateConfirmation(msg)

		if m.deletionState.IsConfirmed() {
			return m.performConfigDeletion()
		}

		if m.deletionState.IsCancelled() {
			m.deletionState.Clear()
			return m, nil
		}

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		m.listContainer.SetDimensions(msg.Width, msg.Height)
		m.updateTableRows()
		return m, nil

	case ConfigsLoadedMsg:
		m.loading = false
		m.configs = msg.configs
		m.applyFiltersAndSort()
		m.pagination.SetTotalCount(len(m.filtered))
		m.updateTableRows()
		return m, nil

	case ConfigLoadError:
		m.loading = false
		m.SetError(msg.err)
		m.listContainer.SetErrorMessage(fmt.Sprintf("Error loading configurations: %v", msg.err))
		return m, nil

	case tea.KeyMsg:
		// Allow quit/escape even while loading
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}

		// If still loading, don't process other keys
		if m.loading {
			return m, nil
		}

		// Handle deletion confirmation
		if m.deletionState.IsConfirming() {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.listContainer.GetSelectedIdx() > 0 {
				m.listContainer.MoveUp(1)
			}
			m.updateTableRows()
			return m, nil

		case "down", "j":
			if m.listContainer.GetSelectedIdx() < len(m.filtered)-1 {
				m.listContainer.MoveDown(1)
			}
			m.updateTableRows()
			return m, nil

		case "home", "g":
			m.listContainer.MoveToFirst()
			m.updateTableRows()
			return m, nil

		case "end", "G":
			m.listContainer.MoveToLast()
			m.updateTableRows()
			return m, nil

		case "enter":
			if len(m.filtered) > 0 {
				selectedIdx := m.listContainer.GetSelectedIdx()
				if selectedIdx < len(m.filtered) {
					selectedConfig := m.filtered[selectedIdx]
					return m, func() tea.Msg {
						return GenerateCVFromConfigMsg{Config: selectedConfig}
					}
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
			if len(m.filtered) > 0 {
				selectedIdx := m.listContainer.GetSelectedIdx()
				if selectedIdx < len(m.filtered) {
					selectedConfig := m.filtered[selectedIdx]
					return m, func() tea.Msg {
						return EditCVConfigMsg{config: selectedConfig}
					}
				}
			}
			return m, nil

		case "d":
			// Delete config (with confirmation)
			if len(m.filtered) > 0 {
				selectedIdx := m.listContainer.GetSelectedIdx()
				if selectedIdx < len(m.filtered) {
					selectedConfig := m.filtered[selectedIdx]
					m.deletionState.ShowConfirmation("Configuration", selectedConfig.Name)
					m.deletionState.DeletingItemID = selectedConfig.Name
					return m, nil
				}
			}
			return m, nil

		case "r":
			// Retry loading configs
			if m.GetLastError() != nil {
				m.ClearError()
				m.loading = true
				m.listContainer.ClearError()
				return m, m.loadConfigs()
			}
			return m, nil
		}
	}

	// Return model unchanged for unhandled messages
	return m, nil
}

// performConfigDeletion performs the actual deletion of a CV configuration
func (m *CVConfigManagerModel) performConfigDeletion() (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		err := m.configManager.DeleteConfig(context.Background(), m.deletionState.DeletingItemID)
		if err != nil {
			m.deletionState.SetErrorMsg(fmt.Sprintf("Failed to delete: %v", err))
			return ConfigDeletionError{err: err}
		}
		// Remove from filtered list
		for i, config := range m.filtered {
			if config.Name == m.deletionState.DeletingItemID {
				m.filtered = append(m.filtered[:i], m.filtered[i+1:]...)
				break
			}
		}
		// Remove from configs list
		for i, config := range m.configs {
			if config.Name == m.deletionState.DeletingItemID {
				m.configs = append(m.configs[:i], m.configs[i+1:]...)
				break
			}
		}
		m.pagination.SetTotalCount(len(m.filtered))
		m.updateTableRows()
		m.deletionState.SetSuccessMsg("Configuration deleted successfully")
		return ConfigDeletedMsg{}
	}
}

// applyFiltersAndSort applies filters and sorting to the configurations
func (m *CVConfigManagerModel) applyFiltersAndSort() {
	m.filtered = m.configs
	m.sortConfigs()
}

// sortConfigs sorts the filtered configurations
func (m *CVConfigManagerModel) sortConfigs() {
	// Sort by the specified field
	if m.sortBy == "updated" {
		if m.sortOrder == "asc" {
			// Sort ascending
			for i := 0; i < len(m.filtered)-1; i++ {
				for j := i + 1; j < len(m.filtered); j++ {
					if m.filtered[i].UpdatedAt.After(m.filtered[j].UpdatedAt) {
						m.filtered[i], m.filtered[j] = m.filtered[j], m.filtered[i]
					}
				}
			}
		} else {
			// Sort descending
			for i := 0; i < len(m.filtered)-1; i++ {
				for j := i + 1; j < len(m.filtered); j++ {
					if m.filtered[i].UpdatedAt.Before(m.filtered[j].UpdatedAt) {
						m.filtered[i], m.filtered[j] = m.filtered[j], m.filtered[i]
					}
				}
			}
		}
	}
}

// updateTableRows updates the table with rows from filtered configs
func (m *CVConfigManagerModel) updateTableRows() {
	var rows []table.Row
	pageConfigs := m.getPageConfigs()

	// Get the selected index within the current page
	selectedIdx := m.listContainer.GetSelectedIdx()

	// Ensure cursor is within valid bounds
	if len(pageConfigs) > 0 {
		if selectedIdx >= len(pageConfigs) {
			selectedIdx = len(pageConfigs) - 1
		}
		if selectedIdx < 0 {
			selectedIdx = 0
		}
	} else {
		selectedIdx = 0
	}

	for i, config := range pageConfigs {
		audiences := ""
		if len(config.TargetAudience) > 0 {
			audiences = config.TargetAudience[0]
			if len(config.TargetAudience) > 1 {
				audiences = fmt.Sprintf("%d audiences", len(config.TargetAudience))
			}
		}

		name := config.Name
		if len(name) > 28 {
			name = name[:25] + "..."
		}

		// Add focus indicator for the selected row
		if i == selectedIdx {
			name = "▶ " + name
		} else {
			name = "  " + name
		}

		row := table.Row{
			name,
			config.TargetRole,
			audiences,
			config.CreatedAt.Format("2006-01-02"),
		}
		rows = append(rows, row)
	}

	m.table.SetRows(rows)
	m.listContainer.SetSelectedIdx(selectedIdx)
	m.table.SetCursor(selectedIdx)
}

// getPageConfigs returns the configurations for the current page
func (m *CVConfigManagerModel) getPageConfigs() []*career.CVConfig {
	if len(m.filtered) == 0 {
		return []*career.CVConfig{}
	}

	startIdx := m.pagination.GetPageStartIndex()
	endIdx := m.pagination.GetPageEndIndex()

	if startIdx >= len(m.filtered) {
		startIdx = 0
	}
	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}

	if startIdx >= endIdx {
		return []*career.CVConfig{}
	}

	return m.filtered[startIdx:endIdx]
}

// View renders the CV Configuration Manager screen.
func (m *CVConfigManagerModel) View() string {
	// Always show something useful, even if loading
	if m.loading {
		headerView := m.header.View()
		loadingMsg := "Loading configuration templates...\n\nPress 'q' or 'esc' to cancel"
		return fmt.Sprintf("%s\n\n%s", headerView, loadingMsg)
	}

	// Handle error state
	var errorContent string
	if m.GetLastError() != nil {
		m.listContainer.SetErrorMessage(fmt.Sprintf("Error loading configurations: %v", m.GetLastError()))
		errorContent = styles.ErrorBox.Render(fmt.Sprintf("Error: %v. Press 'r' to retry", m.GetLastError()))
	}

	// Handle empty state
	if len(m.filtered) == 0 && m.GetLastError() == nil {
		m.listContainer.SetEmptyStateMessage("No CV configuration templates found.\n\nPress 'n' to create a new template")
	}

	// Update list container with current state
	m.listContainer.SetTable(m.table).
		SetDimensions(m.width, m.height).
		SetHelpFooterKey("cv_config_manager").
		SetBreadcrumbs(m.breadcrumbs)

	// Set pagination info
	startIdx := m.pagination.GetPageStartIndex()
	endIdx := m.pagination.GetPageEndIndex()
	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}
	paginationText := fmt.Sprintf("Showing %d-%d of %d templates", startIdx+1, endIdx, len(m.filtered))
	m.listContainer.SetPaginationInfo(paginationText)

	headerView := m.header.View()

	// Handle deletion confirmation
	if m.deletionState.IsConfirming() {
		contentView := m.listContainer.Render()
		confirmView := m.deletionState.ConfirmationDialog.View()
		return fmt.Sprintf("%s\n\n%s\n\n%s", headerView, contentView, confirmView)
	}

	if errorContent != "" {
		return fmt.Sprintf("%s\n\n%s", headerView, errorContent)
	}

	contentView := m.listContainer.Render()
	return fmt.Sprintf("%s\n\n%s", headerView, contentView)
}

// GetSelectedConfig returns the currently selected configuration.
func (m *CVConfigManagerModel) GetSelectedConfig() *career.CVConfig {
	if len(m.filtered) > 0 {
		selectedIdx := m.listContainer.GetSelectedIdx()
		if selectedIdx < len(m.filtered) {
			return m.filtered[selectedIdx]
		}
	}
	return nil
}

// RefreshConfigs reloads the configurations from the manager.
func (m *CVConfigManagerModel) RefreshConfigs() tea.Cmd {
	m.loading = true
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

// ConfigDeletedMsg is sent when a configuration has been deleted
type ConfigDeletedMsg struct{}

// ConfigDeletionError is sent when there's an error deleting a configuration
type ConfigDeletionError struct {
	err error
}

// NavigateToScreenMsg navigates to a specific screen by ID.
type NavigateToScreenMsg struct {
	screenID string
}

// ConfigLoadTimeoutMsg is sent when config loading times out
type ConfigLoadTimeoutMsg struct{}

