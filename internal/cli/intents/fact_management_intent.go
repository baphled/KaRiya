package intents

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type FactManagementModel struct {
	data   *FactManagementContext
	result *IntentResult[*FactManagementResult]
}

func NewFactManagementIntent(data *FactManagementContext) *FactManagementModel {
	return &FactManagementModel{
		data: data,
		result: &IntentResult[*FactManagementResult]{
			Status: Cancelled,
		},
	}
}

func (m *FactManagementModel) Init() tea.Cmd {
	// context already set in data
	m.data.CurrentState = FactListState

	if err := m.data.LoadFacts(); err != nil {
		m.result = &IntentResult[*FactManagementResult]{
			Status: Failed,
			Error: &IntentError{
				Code:    "LOAD_FACTS_FAILED",
				Message: "Failed to load facts",
				Cause:   err,
			},
		}
		return tea.Quit
	}

	return nil
}

func (m *FactManagementModel) Update(msg tea.Msg) tea.Cmd {
	switch m.data.CurrentState {
	case FactListState:
		return m.handleListState(msg)
	case FactViewState:
		return m.handleViewState(msg)
	case FactEditorState:
		return m.handleEditorState(msg)
	case FactDeleteConfirmState:
		return m.handleDeleteConfirmState(msg)
	case FactResultsState:
		return m.handleResultsState(msg)
	case FactCompletedState:
		return tea.Quit
	}
	return nil
}

func (m *FactManagementModel) View() string {
	switch m.data.CurrentState {
	case FactListState:
		return m.viewList()
	case FactViewState:
		return m.viewFact()
	case FactEditorState:
		return m.viewEditor()
	case FactDeleteConfirmState:
		return m.viewDeleteConfirm()
	case FactResultsState:
		return m.viewResults()
	case FactCompletedState:
		return "Fact management completed"
	}
	return "Unknown state"
}

func (m *FactManagementModel) Result() *IntentResult[interface{}] {
	return &IntentResult[interface{}]{
		Status: m.result.Status,
		Error:  m.result.Error,
	}
}

func (m *FactManagementModel) handleListState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.data.CurrentState = FactCompletedState
			m.result = &IntentResult[*FactManagementResult]{
				Status: Cancelled,
				Data: &FactManagementResult{
					Action: "none",
					Facts:  m.data.Facts,
				},
			}
			return tea.Quit

		case "j", "down":
			if m.data.SelectedFactIndex < len(m.data.GetPageFacts())-1 {
				m.data.SelectFact(m.data.SelectedFactIndex + 1)
			}

		case "k", "up":
			if m.data.SelectedFactIndex > 0 {
				m.data.SelectFact(m.data.SelectedFactIndex - 1)
			}

		case "enter", " ":
			if m.data.SelectedFact != nil {
				m.data.CurrentState = FactViewState
			}

		case "n":
			m.data.StartNewFact()
			m.data.CurrentState = FactEditorState

		case "r":
			if err := m.data.LoadFacts(); err != nil {
				m.result = &IntentResult[*FactManagementResult]{
					Status: Failed,
					Error: &IntentError{
						Code:    "RELOAD_FAILED",
						Message: "Failed to reload facts",
						Cause:   err,
					},
				}
			}

		case "tab", "right":
			if m.data.CurrentPage < (m.data.TotalFacts / m.data.PageSize) {
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

func (m *FactManagementModel) handleViewState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.data.CurrentState = FactListState
			m.data.SelectedFact = nil

		case "e":
			if m.data.SelectedFact != nil {
				m.data.StartEditFact(m.data.SelectedFact)
				m.data.CurrentState = FactEditorState
			}

		case "d":
			if m.data.SelectedFact != nil {
				m.data.FactToDelete = m.data.SelectedFact
				m.data.CurrentState = FactDeleteConfirmState
			}
		}
	}
	return nil
}

func (m *FactManagementModel) handleEditorState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			if err := m.data.SaveEdit(); err != nil {
				m.data.SetFormError("general", fmt.Sprintf("Save failed: %v", err))
			} else {
				action := "updated"
				if m.data.IsNewFact {
					action = "created"
				}
				m.result = &IntentResult[*FactManagementResult]{
					Status: Completed,
					Data: &FactManagementResult{
						Action:  action,
						Fact:    m.data.EditingFact,
						Facts:   m.data.Facts,
						Message: fmt.Sprintf("Fact %s successfully", action),
					},
				}
				m.data.CurrentState = FactListState
				m.data.CancelEdit()
			}

		case "esc":
			m.data.CancelEdit()
			if m.data.IsNewFact {
				m.data.CurrentState = FactListState
			} else {
				m.data.CurrentState = FactViewState
			}
		}
	}
	return nil
}

func (m *FactManagementModel) handleDeleteConfirmState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			if m.data.FactToDelete != nil {
				if err := m.data.DeleteFact(m.data.FactToDelete.ID); err != nil {
					m.result = &IntentResult[*FactManagementResult]{
						Status: Failed,
						Error: &IntentError{
							Code:    "DELETE_FAILED",
							Message: "Failed to delete fact",
							Cause:   err,
						},
					}
				} else {
					m.result = &IntentResult[*FactManagementResult]{
						Status: Completed,
						Data: &FactManagementResult{
							Action:  "deleted",
							Fact:    m.data.FactToDelete,
							Facts:   m.data.Facts,
							Message: "Fact deleted successfully",
						},
					}
				}
				m.data.FactToDelete = nil
				m.data.CurrentState = FactListState
			}

		case "n", "esc":
			m.data.FactToDelete = nil
			m.data.CurrentState = FactViewState
		}
	}
	return nil
}

func (m *FactManagementModel) handleResultsState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.data.CurrentState = FactListState
		}
	}
	return nil
}

// View rendering methods

func (m *FactManagementModel) viewList() string {
	if len(m.data.Facts) == 0 {
		return "No facts found. Press 'n' to create a new fact, 'r' to refresh, or 'q' to quit."
	}

	output := "Facts (j/k: navigate, enter: view, n: new, r: refresh, q: quit)\n"
	output += "=================================================================\n"

	pageFacts := m.data.GetPageFacts()
	for i, fact := range pageFacts {
		prefix := "  "
		if i == m.data.SelectedFactIndex {
			prefix = "> "
		}
		output += fmt.Sprintf("%s[%d] %s\n", prefix, i+1, truncate(fact.Text, 50))
		if m.data.IsRowExpanded(i) {
			output += fmt.Sprintf("    Categories: %v\n", fact.CompetencyCategories)
			output += fmt.Sprintf("    Strength: %s\n", fact.StrengthSignal)
		}
	}

	output += fmt.Sprintf("\nPage %d of %d | Total: %d facts\n", m.data.CurrentPage+1, (m.data.TotalFacts/m.data.PageSize)+1, m.data.TotalFacts)

	return output
}

func (m *FactManagementModel) viewFact() string {
	if m.data.SelectedFact == nil {
		return "No fact selected"
	}

	fact := m.data.SelectedFact
	output := fmt.Sprintf("Fact: %s\n", truncate(fact.Text, 70))
	output += "=================================================================\n"
	output += fmt.Sprintf("Categories: %v\n", fact.CompetencyCategories)
	output += fmt.Sprintf("Strength Signal: %s\n", fact.StrengthSignal)
	output += fmt.Sprintf("Role Fit: %v\n", fact.RoleFit)
	output += fmt.Sprintf("Created: %s\n", fact.CreatedAt.Format("2006-01-02"))
	output += "\nOptions: e (edit), d (delete), esc (back), q (quit)\n"

	return output
}

func (m *FactManagementModel) viewEditor() string {
	output := "Edit Fact\n"
	output += "=================================================================\n"

	if m.data.EditingFact != nil {
		output += fmt.Sprintf("Text: %s\n", truncate(m.data.EditingFact.Text, 70))
		output += fmt.Sprintf("Categories: %v\n", m.data.EditingFact.CompetencyCategories)
		output += fmt.Sprintf("Strength: %s\n", m.data.EditingFact.StrengthSignal)

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

func (m *FactManagementModel) viewDeleteConfirm() string {
	if m.data.FactToDelete == nil {
		return "No fact to delete"
	}

	output := fmt.Sprintf("Delete Fact: %s?\n", truncate(m.data.FactToDelete.Text, 50))
	output += "=================================================================\n"
	output += "This action cannot be undone.\n"
	output += "\nOptions: y (confirm), n (cancel), esc (back)\n"

	return output
}

func (m *FactManagementModel) viewResults() string {
	output := "Fact Results\n"
	output += "=================================================================\n"
	output += fmt.Sprintf("Total facts: %d\n", m.data.TotalFacts)
	output += "\nOptions: esc (back), q (quit)\n"

	return output
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
