package intents

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	domain "github.com/baphled/kariya/internal/domain/career"
)

// BurstManagementModel implements the Intent interface for burst management
type BurstManagementModel struct {
	data *BurstManagementContext
	result *IntentResult[*BurstManagementResult]
}

// NewBurstManagementIntent creates a new BurstManagement intent
func NewBurstManagementIntent(data *BurstManagementContext) *BurstManagementModel {
	return &BurstManagementModel{
		data: data,
		result: &IntentResult[*BurstManagementResult]{
			Status: Cancelled,
		},
	}
}

// Init initializes the intent and loads bursts
func (m *BurstManagementModel) Init(ctx context.Context) tea.Cmd {
	m.data.Context = ctx
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

	return nil
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
			if m.data.SelectedBurstIndex < len(m.data.GetPageBursts())-1 {
				m.data.SelectBurst(m.data.SelectedBurstIndex + 1)
			}

		case "k", "up":
			if m.data.SelectedBurstIndex > 0 {
				m.data.SelectBurst(m.data.SelectedBurstIndex - 1)
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
			}

		case "tab", "right":
			if m.data.CurrentPage < (m.data.TotalBursts / m.data.PageSize) {
				m.data.CurrentPage++
			}

		case "shift+tab", "left":
			if m.data.CurrentPage > 0 {
				m.data.CurrentPage--
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
					Name:               suggestion.Title,
					Description:        suggestion.Description,
					CompetencyFocus:    suggestion.CompetencyFocus,
					CreatedAt:          suggestion.RecommendedStartDate,
					UpdatedAt:          suggestion.RecommendedEndDate,
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
		return "No bursts found. Press 'n' to create a new burst, 'r' to refresh, or 'q' to quit."
	}

	output := "Bursts (j/k: navigate, enter: view, n: new, r: refresh, q: quit)\n"
	output += "=================================================================\n"

	pageBursts := m.data.GetPageBursts()
	for i, burst := range pageBursts {
		prefix := "  "
		if i == m.data.SelectedBurstIndex {
			prefix = "> "
		}
		output += fmt.Sprintf("%s[%d] %s (Events: %d)\n", prefix, i+1, burst.Name, len(burst.EventIDs))
		if m.data.IsRowExpanded(i) {
			output += fmt.Sprintf("    Description: %s\n", burst.Description)
			output += fmt.Sprintf("    Competency: %s\n", burst.CompetencyFocus)
		}
	}

	output += fmt.Sprintf("\nPage %d of %d | Total: %d bursts\n", m.data.CurrentPage+1, (m.data.TotalBursts/m.data.PageSize)+1, m.data.TotalBursts)

	return output
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

