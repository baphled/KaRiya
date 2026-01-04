package intents

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FactManagementModel struct {
	data          *FactManagementContext
	table         *table.Model
	listContainer *components.TableListContainer
	result        *IntentResult[*FactManagementResult]
}

func NewFactManagementIntent(data *FactManagementContext) *FactManagementModel {
	// Create table model for facts
	columns := []table.Column{
		{Title: "Fact", Width: 50},
		{Title: "Strength", Width: 15},
		{Title: "Categories", Width: 30},
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

	return &FactManagementModel{
		data:          data,
		table:         &t,
		listContainer: components.NewTableListContainer(t, "Manage Facts", 100),
		result: nil,
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

	m.updateTableRows()
	return nil
}

// updateTableRows updates the table rows based on facts
func (m *FactManagementModel) updateTableRows() {
	pageSize := 15
	total := len(m.data.Facts)

	// Determine which page current selection is on
	page := 0
	if pageSize > 0 && m.data.SelectedFactIndex >= 0 {
		page = m.data.SelectedFactIndex / pageSize
	}

	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	pageFacts := m.data.Facts[start:end]

	rows := make([]table.Row, 0, len(pageFacts))
	for idx, fact := range pageFacts {
		realIdx := start + idx
		text := truncate(fact.Text, 50)

		// Add visual indicator for selected row
		if realIdx == m.data.SelectedFactIndex {
			text = "▶ " + text
		} else {
			text = "  " + text
		}

		strength := fact.StrengthSignal
		if strength == "" {
			strength = "-"
		}
		categories := fmt.Sprintf("%v", fact.CompetencyCategories)
		if len(categories) > 30 {
			categories = categories[:27] + "..."
		}
		rows = append(rows, table.Row{text, strength, categories})
	}

	m.table.SetRows(rows)

	// Set table cursor relative to page
	if m.data.SelectedFactIndex >= start && m.data.SelectedFactIndex < end {
		m.table.SetCursor(m.data.SelectedFactIndex - start)
	} else if len(rows) > 0 {
		m.table.SetCursor(0)
	}

	m.listContainer.SetTable(*m.table)
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
		// Do not send tea.Quit, just return nil and let router manage exit.
		return nil
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
	if m.result == nil {
		return nil
	}
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
			return nil

		case "j", "down":
			// Move down by 1 within the current list bounds
			if m.data.SelectedFactIndex < len(m.data.Facts)-1 {
				m.data.SelectedFactIndex++
				m.data.SelectedFact = m.data.Facts[m.data.SelectedFactIndex]
				m.updateTableRows()
			}
			return nil

		case "k", "up":
			// Move up by 1 within the current list bounds
			if m.data.SelectedFactIndex > 0 {
				m.data.SelectedFactIndex--
				m.data.SelectedFact = m.data.Facts[m.data.SelectedFactIndex]
				m.updateTableRows()
			}
			return nil

		case "pgup", "b":
			// Page up (move up by page height).
			pageSize := 15
			newIndex := m.data.SelectedFactIndex - pageSize
			if newIndex < 0 {
				newIndex = 0
			}
			m.data.SelectedFactIndex = newIndex
			m.table.SetCursor(newIndex)
			if m.data.SelectedFactIndex >= 0 && m.data.SelectedFactIndex < len(m.data.Facts) {
				m.data.SelectedFact = m.data.Facts[m.data.SelectedFactIndex]
			}
			m.updateTableRows()
			return nil

		case "pgdn", "f":
			// Page down (move down by page height).
			pageSize := 15
			newIndex := m.data.SelectedFactIndex + pageSize
			if newIndex >= len(m.data.Facts) {
				newIndex = len(m.data.Facts) - 1
			}
			if newIndex < 0 {
				newIndex = 0
			}
			m.data.SelectedFactIndex = newIndex
			m.table.SetCursor(newIndex)
			if m.data.SelectedFactIndex >= 0 && m.data.SelectedFactIndex < len(m.data.Facts) {
				m.data.SelectedFact = m.data.Facts[m.data.SelectedFactIndex]
			}
			m.updateTableRows()
			return nil

		case "home", "g":
			// Go to first fact.
			m.table.SetCursor(0)
			m.data.SelectedFactIndex = 0
			if len(m.data.Facts) > 0 {
				m.data.SelectedFact = m.data.Facts[0]
			}
			m.updateTableRows()
			return nil

		case "end", "G":
			// Go to last fact.
			if len(m.data.Facts) > 0 {
				lastIndex := len(m.data.Facts) - 1
				m.table.SetCursor(lastIndex)
				m.data.SelectedFactIndex = lastIndex
				m.data.SelectedFact = m.data.Facts[lastIndex]
			}
			m.updateTableRows()
			return nil

		case "enter", " ":
			if m.data.SelectedFact != nil {
				m.data.CurrentState = FactViewState
			}
			return nil

		case "n":
			m.data.StartNewFact()
			m.data.CurrentState = FactEditorState

			return nil

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
			} else {
				m.updateTableRows()
			}

			return nil

		case "tab", "right":
			if m.data.CurrentPage < (m.data.TotalFacts / m.data.PageSize) {
				m.data.CurrentPage++
				m.updateTableRows()
			}

			return nil

		case "shift+tab", "left":
			if m.data.CurrentPage > 0 {
				m.data.CurrentPage--
				m.updateTableRows()
			}

			return nil
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
				m.updateTableRows()
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
					m.updateTableRows()
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
		m.listContainer.SetEmptyStateMessage("No facts found. Press 'n' to create a new fact, 'r' to refresh, or 'q' to quit.")
		return m.listContainer.Render()
	}

	// Ensure table rows are synchronized with current state
	m.updateTableRows()

	// Build pagination info with page number indicator
	pageSize := 15
	totalItems := len(m.data.Facts)
	currentPage := (m.data.SelectedFactIndex / pageSize) + 1
	totalPages := (totalItems + pageSize - 1) / pageSize
	paginationInfo := fmt.Sprintf("Total: %d facts | Page %d of %d", totalItems, currentPage, totalPages)
	m.listContainer.SetPaginationInfo(paginationInfo)

	// Set breadcrumbs if needed
	m.listContainer.SetBreadcrumbs([]string{"Home", "Facts"})

	// Set help footer
	m.listContainer.SetHelpFooterKey("fact_management")

	return m.listContainer.Render()
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
