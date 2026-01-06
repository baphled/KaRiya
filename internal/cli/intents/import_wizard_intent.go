package intents

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ImportWizardModel struct {
	data   *ImportWizardContext
	result *IntentResult[*ImportWizardResult]
}

func NewImportWizardIntent(data *ImportWizardContext) *ImportWizardModel {
	return &ImportWizardModel{
		data:   data,
		result: nil,
	}
}

func (m *ImportWizardModel) Init() tea.Cmd {
	// context already set in data
	m.data.CurrentState = ImportFileSelectState
	// Return a no-op command to satisfy the intent lifecycle
	return func() tea.Msg { return nil }
}

func (m *ImportWizardModel) Update(msg tea.Msg) tea.Cmd {
	switch m.data.CurrentState {
	case ImportFileSelectState:
		return m.handleFileSelectState(msg)
	case ImportPreviewState:
		return m.handlePreviewState(msg)
	case ImportProgressState:
		return m.handleProgressState(msg)
	case ImportCompleteState:
		return tea.Quit
	}
	return nil
}

func (m *ImportWizardModel) View() string {
	switch m.data.CurrentState {
	case ImportFileSelectState:
		return m.viewFileSelect()
	case ImportPreviewState:
		return m.viewPreview()
	case ImportProgressState:
		return m.viewProgress()
	case ImportCompleteState:
		return m.viewComplete()
	}
	return "Unknown state"
}

func (m *ImportWizardModel) Result() *IntentResult[interface{}] {
	if m.result == nil {
		return nil
	}
	return &IntentResult[interface{}]{
		Status: m.result.Status,
		Error:  m.result.Error,
	}
}

func (m *ImportWizardModel) handleFileSelectState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.data.CurrentState = ImportCompleteState
			m.result = &IntentResult[*ImportWizardResult]{
				Status: Cancelled,
				Data: &ImportWizardResult{
					Action: "cancelled",
				},
			}
			return tea.Quit

		case "enter":
			if m.data.FilePath != "" {
				m.data.CurrentState = ImportPreviewState
			}
		}
	}
	return nil
}

func (m *ImportWizardModel) handlePreviewState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.data.CurrentState = ImportFileSelectState

		case "enter":
			m.data.StartImport()
			m.data.CurrentState = ImportProgressState
		}
	}
	return nil
}

func (m *ImportWizardModel) handleProgressState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			if m.data.IsImporting && !m.data.IsPaused {
				m.data.PauseImport()
			} else if m.data.IsPaused {
				m.data.ResumeImport()
			}

		case "c":
			m.data.CancelImport()
			m.data.CurrentState = ImportCompleteState
			m.result = &IntentResult[*ImportWizardResult]{
				Status: Cancelled,
				Data: &ImportWizardResult{
					Action:        "cancelled",
					ProcessedRows: m.data.ProcessedRows,
				},
			}
			return tea.Quit
		}
	}
	return nil
}

func (m *ImportWizardModel) viewFileSelect() string {
	output := "Select File for Import\n"
	output += "=================================================================\n"
	output += "Current file: "
	if m.data.FilePath != "" {
		output += m.data.FilePath
	} else {
		output += "(none selected)"
	}
	output += "\n\nOptions: enter (proceed), q (quit)\n"
	return output
}

func (m *ImportWizardModel) viewPreview() string {
	output := "Import Preview\n"
	output += "=================================================================\n"
	output += fmt.Sprintf("File: %s\n", m.data.FilePath)
	output += fmt.Sprintf("Size: %d bytes\n", m.data.FileSize)
	output += fmt.Sprintf("Total rows: %d\n", m.data.TotalRows)
	output += "\nOptions: enter (start import), esc (back), q (quit)\n"
	return output
}

func (m *ImportWizardModel) viewProgress() string {
	output := "Import Progress\n"
	output += "=================================================================\n"
	progress := m.data.GetProgress()
	progressBar := fmt.Sprintf("[%-50s] %.0f%%", "=", progress*100)
	output += progressBar + "\n"
	output += fmt.Sprintf("Processed: %d/%d\n", m.data.ProcessedRows, m.data.TotalRows)
	output += fmt.Sprintf("Successful: %d | Errors: %d\n", m.data.SuccessfulRows, m.data.ErrorRows)

	if m.data.IsPaused {
		output += "\n[PAUSED] Options: p (resume), c (cancel)\n"
	} else if m.data.IsImporting {
		output += "\nOptions: p (pause), c (cancel)\n"
	}
	return output
}

func (m *ImportWizardModel) viewComplete() string {
	output := "Import Complete\n"
	output += "=================================================================\n"
	output += fmt.Sprintf("Processed: %d rows\n", m.data.ProcessedRows)
	output += fmt.Sprintf("Successful: %d rows\n", m.data.SuccessfulRows)
	output += fmt.Sprintf("Errors: %d rows\n", m.data.ErrorRows)

	if len(m.data.Errors) > 0 {
		output += "\nFirst error: " + m.data.Errors[0] + "\n"
	}

	output += "\nPress any key to exit...\n"
	return output
}
