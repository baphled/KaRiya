package intents

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	tea "github.com/charmbracelet/bubbletea"
)

type ImportWizardModel struct {
	*BaseIntent
	data   *ImportWizardContext
	result *IntentResult[*ImportWizardResult]
	active bool
}

func NewImportWizardIntent(data *ImportWizardContext) *ImportWizardModel {
	return &ImportWizardModel{
		BaseIntent: NewBaseIntent(),
		data:       data,
		result:     nil,
		active:     false,
	}
}

func (m *ImportWizardModel) Init() tea.Cmd {
	// Mark intent as active
	m.active = true

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
	if !m.active {
		return "ImportWizard intent is not active"
	}

	// Create standard view with breadcrumbs
	view := m.CreateViewWithBreadcrumbs("Main Menu", "Import Data", m.getStateName())

	// Handle progress modal for import state
	if m.data.CurrentState == ImportProgressState && m.data.IsImporting {
		progress := m.data.GetProgress()
		message := fmt.Sprintf("Importing: %d/%d rows", m.data.ProcessedRows, m.data.TotalRows)
		if m.data.SuccessfulRows > 0 || m.data.ErrorRows > 0 {
			message += fmt.Sprintf("\nSuccessful: %d | Errors: %d", m.data.SuccessfulRows, m.data.ErrorRows)
		}
		if m.data.IsPaused {
			message += "\n\n⏸ PAUSED - Press 'p' to resume"
		}
		m.SetProgress("Importing CSV", message, progress)
	}

	// Handle errors (but not during complete state - we show them in content there)
	if len(m.data.Errors) > 0 && m.data.CurrentState != ImportCompleteState && m.data.CurrentState != ImportProgressState {
		// Show first error as modal
		m.SetError(fmt.Errorf("import error: %s", m.data.Errors[0]))
	}

	// Get content for current state
	content := m.getStateContent()
	view.WithContent(content)

	// Get context-aware help
	help := m.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
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

// getStateName returns the display name for the current state
func (m *ImportWizardModel) getStateName() string {
	switch m.data.CurrentState {
	case ImportFileSelectState:
		return "Select File"
	case ImportPreviewState:
		return "Preview"
	case ImportProgressState:
		return "Importing"
	case ImportCompleteState:
		return "Complete"
	default:
		return ""
	}
}

// getStateContent returns the content for the current state
func (m *ImportWizardModel) getStateContent() string {
	switch m.data.CurrentState {
	case ImportFileSelectState:
		return m.getFileSelectContent()
	case ImportPreviewState:
		return m.getPreviewContent()
	case ImportProgressState:
		// Progress shown in modal, content is minimal
		return m.getProgressContent()
	case ImportCompleteState:
		return m.getCompleteContent()
	default:
		return "Unknown state"
	}
}

// getContextHelp returns context-aware help text
func (m *ImportWizardModel) getContextHelp() string {
	theme := m.Theme()

	switch m.data.CurrentState {
	case ImportFileSelectState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Select file"),
				components.CancelBadge(),
			),
			ThemedCustomFooter(theme, components.QuitBadge()),
		)
	case ImportPreviewState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Start import"),
				components.BackBadge(),
			),
			ThemedCustomFooter(theme, components.QuitBadge()),
		)
	case ImportProgressState:
		if m.data.IsPaused {
			return CombineThemedFooters(
				ThemedCustomFooter(theme,
					components.NewKeyBadge("p", "Resume"),
					components.NewKeyBadge("c", "Cancel import"),
				),
				ThemedCustomFooter(theme, components.QuitBadge()),
			)
		}
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("p", "Pause"),
				components.NewKeyBadge("c", "Cancel import"),
			),
			ThemedCustomFooter(theme, components.QuitBadge()),
		)
	case ImportCompleteState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Done"),
			),
			ThemedCustomFooter(theme, components.QuitBadge()),
		)
	default:
		return ThemedCustomFooter(theme, components.QuitBadge())
	}
}

// Content methods - refactored from view methods

func (m *ImportWizardModel) getFileSelectContent() string {
	var content strings.Builder

	content.WriteString("📂 Select CSV File for Import\n\n")

	if m.data.FilePath != "" {
		content.WriteString(fmt.Sprintf("Selected file: %s\n", m.data.FilePath))
	} else {
		content.WriteString("❌ No file selected\n")
	}

	content.WriteString("\n")
	content.WriteString("Please ensure you have a CSV file ready with career event data.\n")
	content.WriteString("The file should contain columns for event details.\n")

	return content.String()
}

func (m *ImportWizardModel) getPreviewContent() string {
	var content strings.Builder

	content.WriteString("📋 Import Preview\n\n")
	content.WriteString(fmt.Sprintf("File: %s\n", m.data.FilePath))
	content.WriteString(fmt.Sprintf("Size: %s\n", humanizeBytes(m.data.FileSize)))
	content.WriteString(fmt.Sprintf("Total rows: %d\n", m.data.TotalRows))

	content.WriteString("\n")
	content.WriteString("Review the file details above.\n")
	content.WriteString("Press Enter to begin importing or Esc to go back.\n")

	return content.String()
}

func (m *ImportWizardModel) getProgressContent() string {
	// Most progress info is in the modal, just show a simple status
	var content strings.Builder

	content.WriteString("⏳ Import in Progress\n\n")

	if m.data.IsPaused {
		content.WriteString("Import is currently paused.\n")
		content.WriteString("Press 'p' to resume or 'c' to cancel.\n")
	} else {
		content.WriteString("Importing your CSV file...\n")
		content.WriteString("Progress details are shown above.\n")
	}

	return content.String()
}

func (m *ImportWizardModel) getCompleteContent() string {
	var content strings.Builder

	if m.data.ErrorRows > 0 {
		content.WriteString("⚠️  Import Complete with Errors\n\n")
	} else {
		content.WriteString("✅ Import Complete\n\n")
	}

	content.WriteString(fmt.Sprintf("Total processed: %d rows\n", m.data.ProcessedRows))
	content.WriteString(fmt.Sprintf("Successful: %d rows\n", m.data.SuccessfulRows))
	content.WriteString(fmt.Sprintf("Errors: %d rows\n", m.data.ErrorRows))

	if len(m.data.Errors) > 0 {
		content.WriteString("\n📝 Error Summary:\n")
		// Show first 5 errors
		maxErrors := 5
		if len(m.data.Errors) < maxErrors {
			maxErrors = len(m.data.Errors)
		}
		for i := 0; i < maxErrors; i++ {
			content.WriteString(fmt.Sprintf("  • %s\n", m.data.Errors[i]))
		}
		if len(m.data.Errors) > 5 {
			content.WriteString(fmt.Sprintf("  ... and %d more errors\n", len(m.data.Errors)-5))
		}
	}

	return content.String()
}

// Helper function to humanize byte sizes
func humanizeBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// State handlers - unchanged from original

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

		case "esc":
			// Cancel and return to menu
			m.result = &IntentResult[*ImportWizardResult]{
				Status: Cancelled,
				Data: &ImportWizardResult{
					Action: "cancelled",
				},
			}
			m.active = false
			return nil

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
		case "q", "ctrl+c":
			m.result = &IntentResult[*ImportWizardResult]{
				Status: Cancelled,
				Data: &ImportWizardResult{
					Action: "cancelled",
				},
			}
			return tea.Quit

		case "esc":
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
		case "q", "ctrl+c":
			// Cancel import and quit
			m.data.CancelImport()
			m.result = &IntentResult[*ImportWizardResult]{
				Status: Cancelled,
				Data: &ImportWizardResult{
					Action:         "cancelled",
					ProcessedRows:  m.data.ProcessedRows,
					SuccessfulRows: m.data.SuccessfulRows,
					ErrorRows:      m.data.ErrorRows,
				},
			}
			return tea.Quit

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
					Action:         "cancelled",
					ProcessedRows:  m.data.ProcessedRows,
					SuccessfulRows: m.data.SuccessfulRows,
					ErrorRows:      m.data.ErrorRows,
				},
			}
		}
	}
	return nil
}
