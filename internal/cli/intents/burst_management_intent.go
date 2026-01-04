package intents

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BurstManagementModel implements the Intent interface for burst management
type BurstManagementModel struct {
	data          *BurstManagementContext
	table         *table.Model
	listContainer *components.TableListContainer
	result        *IntentResult[*BurstManagementResult]
}

// NewBurstManagementIntent creates a new BurstManagement intent
func NewBurstManagementIntent(data *BurstManagementContext) *BurstManagementModel {
	// Create table model for bursts
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Competency", Width: 25},
		{Title: "Events", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

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

	return &BurstManagementModel{
		data:          data,
		table:         &t,
		listContainer: components.NewTableListContainer(t, "Manage Bursts", 100),
		result: nil,
	}
}

// Init initializes the intent and loads bursts
func (m *BurstManagementModel) Init() tea.Cmd {
	// context already set in data
	m.data.CurrentState = BurstListState

	// Load bursts from repository
	if err := m.data.LoadBursts(); err != nil {
		m.result = &IntentResult[*BurstManagementResult]{
			Status: Failed,
			Error: &IntentError{
				Code:    "LOAD_BURSTS_FAILED",
				Message: "Failed to load bursts",
				Cause:   err,
			},
		}
		return tea.Quit
	}

	m.updateTableRows()
	return nil
}

// updateTableRows updates the table rows based on bursts
func (m *BurstManagementModel) updateTableRows() {
	rows := make([]table.Row, 0, len(m.data.Bursts))
	for _, burst := range m.data.Bursts {
		competency := burst.CompetencyFocus
		if competency == "" {
			competency = "-"
		}
		eventCount := fmt.Sprintf("%d", len(burst.EventIDs))
		rows = append(rows, table.Row{burst.Name, competency, eventCount})
	}
	m.table.SetRows(rows)
	m.listContainer.SetTable(*m.table)
}

// Update processes messages and updates the intent state
func (m *BurstManagementModel) Update(msg tea.Msg) tea.Cmd {
	switch m.data.CurrentState {
	case BurstListState:
		return m.handleListState(msg)
	case BurstViewState:
		return m.handleViewState(msg)
	case BurstEditorState:
		return m.handleEditorState(msg)
	case BurstDeleteConfirmState:
		return m.handleDeleteConfirmState(msg)
	case BurstSuggestState:
		return m.handleSuggestState(msg)
	case BurstCompletedState:
		return tea.Quit
	}
	return nil
}

// View renders the current state
func (m *BurstManagementModel) View() string {
	switch m.data.CurrentState {
	case BurstListState:
		return m.viewList()
	case BurstViewState:
		return m.viewBurst()
	case BurstEditorState:
		return m.viewEditor()
	case BurstDeleteConfirmState:
		return m.viewDeleteConfirm()
	case BurstSuggestState:
		return m.viewSuggest()
	case BurstCompletedState:
		return "Burst management completed"
	}
	return "Unknown state"
}

// Result returns the intent result
func (m *BurstManagementModel) Result() *IntentResult[interface{}] {
	return &IntentResult[interface{}]{
		Status: m.result.Status,
		Error:  m.result.Error,
	}
}

// handleListState handles messages in list state
func (m *BurstManagementModel) handleListState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.data.CurrentState = BurstCompletedState
			m.result = &IntentResult[*BurstManagementResult]{
				Status: Cancelled,
				Data: &BurstManagementResult{
					Action: "none",
					Bursts: m.data.Bursts,
				},
			}
			return tea.Quit

		case "j", "down":
			cursor := m.table.Cursor()
			if cursor < len(m.data.Bursts)-1 {
				m.table.SetCursor(cursor + 1)
				m.data.SelectedBurstIndex = cursor + 1
				if m.data.SelectedBurstIndex < len(m.data.Bursts) {
					m.data.SelectedBurst = m.data.Bursts[m.data.SelectedBurstIndex]
				}
			}

		case "k", "up":
			cursor := m.table.Cursor()
			if cursor > 0 {
				m.table.SetCursor(cursor - 1)
				m.data.SelectedBurstIndex = cursor - 1
				if m.data.SelectedBurstIndex >= 0 && m.data.SelectedBurstIndex < len(m.data.Bursts) {
					m.data.SelectedBurst = m.data.Bursts[m.data.SelectedBurstIndex]
				}
			}

		case "enter", " ":
			if m.data.SelectedBurst != nil {
				m.data.CurrentState = BurstViewState
			}

		case "n":
			m.data.StartNewBurst()
			m.data.CurrentState = BurstEditorState

		case "r":
			if err := m.data.LoadBursts(); err != nil {
				m.result = &IntentResult[*BurstManagementResult]{
					Status: Failed,
					Error: &IntentError{
						Code:    "RELOAD_FAILED",
						Message: "Failed to reload bursts",
						Cause:   err,
					},
				}
			} else {
				m.updateTableRows()
			}

		case "tab", "right":
			if m.data.CurrentPage < (m.data.TotalBursts / m.data.PageSize) {
				m.data.CurrentPage++
				m.updateTableRows()
			}

		case "shift+tab", "left":
			if m.data.CurrentPage > 0 {
				m.data.CurrentPage--
				m.updateTableRows()
			}
		}
	}
	return nil
}

// handleViewState handles messages in view state
func (m *BurstManagementModel) handleViewState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.data.CurrentState = BurstListState
			m.data.SelectedBurst = nil

		case "e":
			if m.data.SelectedBurst != nil {
				m.data.StartEditBurst(m.data.SelectedBurst)
				m.data.CurrentState = BurstEditorState
			}

		case "d":
			if m.data.SelectedBurst != nil {
				m.data.BurstToDelete = m.data.SelectedBurst
				m.data.CurrentState = BurstDeleteConfirmState
			}
		}
	}
	return nil
}

// handleEditorState handles messages in editor state
func (m *BurstManagementModel) handleEditorState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			// Save the burst
			if err := m.data.SaveEdit(); err != nil {
				m.data.SetFormError("general", fmt.Sprintf("Save failed: %v", err))
			} else {
				action := "updated"
				if m.data.IsNewBurst {
					action = "created"
				}
				m.result = &IntentResult[*BurstManagementResult]{
					Status: Completed,
					Data: &BurstManagementResult{
						Action:  action,
						Burst:   m.data.EditingBurst,
						Bursts:  m.data.Bursts,
						Message: fmt.Sprintf("Burst %s successfully", action),
					},
				}
				m.data.CurrentState = BurstListState
				m.data.CancelEdit()
				m.updateTableRows()
			}

		case "esc":
			// Cancel editing
			m.data.CancelEdit()
			if m.data.IsNewBurst {
				m.data.CurrentState = BurstListState
			} else {
				m.data.CurrentState = BurstViewState
			}
		}
	}
	return nil
}

// handleDeleteConfirmState handles messages in delete confirm state
func (m *BurstManagementModel) handleDeleteConfirmState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			// Confirm deletion
			if m.data.BurstToDelete != nil {
				if err := m.data.DeleteBurst(m.data.BurstToDelete.ID); err != nil {
					m.result = &IntentResult[*BurstManagementResult]{
						Status: Failed,
						Error: &IntentError{
							Code:    "DELETE_FAILED",
							Message: "Failed to delete burst",
							Cause:   err,
						},
					}
				} else {
					m.result = &IntentResult[*BurstManagementResult]{
						Status: Completed,
						Data: &BurstManagementResult{
							Action:  "deleted",
							Burst:   m.data.BurstToDelete,
							Bursts:  m.data.Bursts,
							Message: "Burst deleted successfully",
						},
					}
					m.updateTableRows()
				}
				m.data.BurstToDelete = nil
				m.data.CurrentState = BurstListState
			}

		case "n", "esc":
			// Cancel deletion
			m.data.BurstToDelete = nil
			m.data.CurrentState = BurstViewState
		}
	}
	return nil
}

// handleSuggestState handles messages in suggest state
func (m *BurstManagementModel) handleSuggestState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.data.SelectedSuggestionIndex < len(m.data.Suggestions)-1 {
				m.data.SelectedSuggestionIndex++
			}

		case "k", "up":
			if m.data.SelectedSuggestionIndex > 0 {
				m.data.SelectedSuggestionIndex--
			}

		case "enter", " ":
			// Accept suggestion
			if m.data.SelectedSuggestionIndex >= 0 && m.data.SelectedSuggestionIndex < len(m.data.Suggestions) {
				suggestion := m.data.Suggestions[m.data.SelectedSuggestionIndex]
				burst := &domain.Burst{
					Name:            suggestion.Title,
					Description:     suggestion.Description,
					CompetencyFocus: suggestion.CompetencyFocus,
					CreatedAt:       suggestion.RecommendedStartDate,
					UpdatedAt:       suggestion.RecommendedEndDate,
				}
				if err := m.data.CreateBurst(burst); err != nil {
					m.result = &IntentResult[*BurstManagementResult]{
						Status: Failed,
						Error: &IntentError{
							Code:    "CREATE_BURST_FAILED",
							Message: "Failed to create burst from suggestion",
							Cause:   err,
						},
					}
				} else {
					m.result = &IntentResult[*BurstManagementResult]{
						Status: Completed,
						Data: &BurstManagementResult{
							Action:  "created",
							Burst:   burst,
							Bursts:  m.data.Bursts,
							Message: "Burst created from suggestion",
						},
					}
					m.updateTableRows()
				}
				m.data.CurrentState = BurstListState
			}

		case "esc", "q":
			m.data.CurrentState = BurstListState
		}
	}
	return nil
}

// View rendering methods

func (m *BurstManagementModel) viewList() string {
	if len(m.data.Bursts) == 0 {
		m.listContainer.SetEmptyStateMessage("No bursts found. Press 'n' to create a new burst, 'r' to refresh, or 'q' to quit.")
		return m.listContainer.Render()
	}

	// Build pagination info
	paginationInfo := fmt.Sprintf("Total: %d bursts", m.data.TotalBursts)
	m.listContainer.SetPaginationInfo(paginationInfo)

	// Set breadcrumbs if needed
	m.listContainer.SetBreadcrumbs([]string{"Home", "Bursts"})

	// Set help footer
	m.listContainer.SetHelpFooterKey("burst_management")

	return m.listContainer.Render()
}

func (m *BurstManagementModel) viewBurst() string {
	if m.data.SelectedBurst == nil {
		return "No burst selected"
	}

	burst := m.data.SelectedBurst
	output := fmt.Sprintf("Burst: %s\n", burst.Name)
	output += "=================================================================\n"
	output += fmt.Sprintf("Description: %s\n", burst.Description)
	output += fmt.Sprintf("Competency Focus: %s\n", burst.CompetencyFocus)
	output += fmt.Sprintf("Events: %d\n", len(burst.EventIDs))
	output += fmt.Sprintf("Created: %s\n", burst.CreatedAt.Format("2006-01-02"))
	output += fmt.Sprintf("Updated: %s\n", burst.UpdatedAt.Format("2006-01-02"))
	output += "\nOptions: e (edit), d (delete), esc (back), q (quit)\n"

	return output
}

func (m *BurstManagementModel) viewEditor() string {
	output := "Edit Burst\n"
	output += "=================================================================\n"

	if m.data.EditingBurst != nil {
		output += fmt.Sprintf("Name: %s\n", m.data.EditingBurst.Name)
		output += fmt.Sprintf("Description: %s\n", m.data.EditingBurst.Description)
		output += fmt.Sprintf("Competency Focus: %s\n", m.data.EditingBurst.CompetencyFocus)

		if m.data.HasFormErrors() {
			output += "\nErrors:\n"
			for field, err := range m.data.FormErrors {
				output += fmt.Sprintf("  %s: %s\n", field, err)
			}
		}
	}

	output += "\nOptions: Ctrl+S (save), Esc (cancel)\n"

	return output
}

func (m *BurstManagementModel) viewDeleteConfirm() string {
	if m.data.BurstToDelete == nil {
		return "No burst to delete"
	}

	output := fmt.Sprintf("Delete Burst: %s?\n", m.data.BurstToDelete.Name)
	output += "=================================================================\n"
	output += "This action cannot be undone.\n"
	output += "\nOptions: y (confirm), n (cancel), esc (back)\n"

	return output
}

func (m *BurstManagementModel) viewSuggest() string {
	output := "Burst Suggestions\n"
	output += "=================================================================\n"

	if len(m.data.Suggestions) == 0 {
		output += "No suggestions available\n"
	} else {
		for i, suggestion := range m.data.Suggestions {
			prefix := "  "
			if i == m.data.SelectedSuggestionIndex {
				prefix = "> "
			}
			output += fmt.Sprintf("%s[%d] %s (Confidence: %.1f%%)\n", prefix, i+1, suggestion.Title, suggestion.ConfidenceScore*100)
			output += fmt.Sprintf("    Events: %d, Skills: %v\n", len(suggestion.Events), suggestion.Skills)
		}
	}

	output += "\nOptions: j/k (navigate), enter (accept), esc (cancel), q (quit)\n"

	return output
}
