package intents

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type BulkOperationsModel struct {
	data   *BulkOperationsContext
	result *IntentResult[*BulkOperationsResult]
}

func NewBulkOperationsIntent(data *BulkOperationsContext) *BulkOperationsModel {
	return &BulkOperationsModel{
		data:   data,
		result: nil,
	}
}

func (m *BulkOperationsModel) Init() tea.Cmd {
	// context already set in data
	m.data.CurrentState = BulkSelectOpState
	// Return a no-op command to satisfy the intent lifecycle
	return func() tea.Msg { return nil }
}

func (m *BulkOperationsModel) Update(msg tea.Msg) tea.Cmd {
	switch m.data.CurrentState {
	case BulkSelectOpState:
		return m.handleSelectOpState(msg)
	case BulkConfigureState:
		return m.handleConfigureState(msg)
	case BulkExecuteState:
		return m.handleExecuteState(msg)
	case BulkCompleteState:
		return tea.Quit
	}
	return nil
}

func (m *BulkOperationsModel) View() string {
	switch m.data.CurrentState {
	case BulkSelectOpState:
		return m.viewSelectOp()
	case BulkConfigureState:
		return m.viewConfigure()
	case BulkExecuteState:
		return m.viewExecute()
	case BulkCompleteState:
		return m.viewComplete()
	}
	return "Unknown state"
}

func (m *BulkOperationsModel) Result() *IntentResult[interface{}] {
	return &IntentResult[interface{}]{
		Status: m.result.Status,
		Error:  m.result.Error,
		Data:   m.result.Data,
	}
}

func (m *BulkOperationsModel) handleSelectOpState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.data.CurrentState = BulkCompleteState
			m.result = &IntentResult[*BulkOperationsResult]{
				Status: Cancelled,
				Data: &BulkOperationsResult{
					Operation: "cancelled",
				},
			}
			return tea.Quit

		case "enter":
			if m.data.SelectedOp != "" {
				m.data.CurrentState = BulkConfigureState
			}
		}
	}
	return nil
}

func (m *BulkOperationsModel) handleConfigureState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.data.StartExecution()
			m.data.CurrentState = BulkExecuteState

		case "esc":
			m.data.CurrentState = BulkSelectOpState
		}
	}
	return nil
}

func (m *BulkOperationsModel) handleExecuteState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			if m.data.IsExecuting && !m.data.IsPaused {
				m.data.PauseExecution()
			} else if m.data.IsPaused {
				m.data.ResumeExecution()
			}

		case "c":
			m.data.CompleteExecution()
			m.data.CurrentState = BulkCompleteState
			m.result = &IntentResult[*BulkOperationsResult]{
				Status: Completed,
				Data: &BulkOperationsResult{
					Operation:      m.data.SelectedOp,
					ProcessedCount: m.data.ProcessedCount,
					SuccessCount:   m.data.SuccessCount,
					FailureCount:   m.data.FailureCount,
					SkippedCount:   m.data.SkippedCount,
					Results:        m.data.Results,
					Errors:         m.data.Errors,
					Message:        "Operation completed",
				},
			}
			return tea.Quit
		}
	}
	return nil
}

func (m *BulkOperationsModel) viewSelectOp() string {
	output := "Select Operation\n"
	output += "=================================================================\n"
	output += "Available operations:\n"
	for i, op := range m.data.AvailableOps {
		prefix := "  "
		if m.data.SelectedOp == op {
			prefix = "> "
		}
		output += fmt.Sprintf("%s[%d] %s\n", prefix, i+1, op)
	}
	output += "\nOptions: enter (proceed), q (quit)\n"
	return output
}

func (m *BulkOperationsModel) viewConfigure() string {
	output := "Configure Operation\n"
	output += "=================================================================\n"
	output += fmt.Sprintf("Operation: %s\n", m.data.SelectedOp)
	output += fmt.Sprintf("Scope: %s\n", m.data.ScopeType)
	output += fmt.Sprintf("Affected items: %d\n", m.data.AffectedItemCount)
	output += "\nOptions: enter (execute), esc (back), q (quit)\n"
	return output
}

func (m *BulkOperationsModel) viewExecute() string {
	output := "Executing Operation\n"
	output += "=================================================================\n"
	progress := m.data.GetProgress()
	progressBar := fmt.Sprintf("[%-50s] %.0f%%", "=", progress*100)
	output += progressBar + "\n"
	output += fmt.Sprintf("Processed: %d/%d\n", m.data.ProcessedCount, m.data.AffectedItemCount)
	output += fmt.Sprintf("Success: %d | Failed: %d | Skipped: %d\n", m.data.SuccessCount, m.data.FailureCount, m.data.SkippedCount)

	if m.data.IsPaused {
		output += "\n[PAUSED] Options: p (resume), c (complete)\n"
	} else if m.data.IsExecuting {
		output += "\nOptions: p (pause), c (complete)\n"
	}
	return output
}

func (m *BulkOperationsModel) viewComplete() string {
	output := "Operation Complete\n"
	output += "=================================================================\n"
	output += fmt.Sprintf("Operation: %s\n", m.data.SelectedOp)
	output += fmt.Sprintf("Processed: %d items\n", m.data.ProcessedCount)
	output += fmt.Sprintf("Success: %d | Failed: %d | Skipped: %d\n", m.data.SuccessCount, m.data.FailureCount, m.data.SkippedCount)
	output += "\nPress any key to exit...\n"
	return output
}
