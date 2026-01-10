package intents

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/navigation"
	tea "github.com/charmbracelet/bubbletea"
)

type BulkOperationsModel struct {
	*BaseIntent
	data       *BulkOperationsContext
	result     *IntentResult[*BulkOperationsResult]
	active     bool
	navHandler *navigation.ListNavigationHandler
}

func NewBulkOperationsIntent(data *BulkOperationsContext) *BulkOperationsModel {
	model := &BulkOperationsModel{
		BaseIntent: NewBaseIntent(),
		data:       data,
		result:     nil,
		active:     false,
	}
	// Initialize navigation handler
	model.navHandler = navigation.NewListNavigationHandler(model)
	return model
}

func (m *BulkOperationsModel) Init() tea.Cmd {
	// Mark intent as active
	m.active = true

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
	if !m.active {
		return "BulkOperations intent is not active"
	}

	// Create standard view with breadcrumbs
	breadcrumbs := []string{"Main Menu", "Bulk Operations"}
	if m.data.SelectedOp != "" {
		breadcrumbs = append(breadcrumbs, m.getOperationDisplayName())
	}
	view := CreateStandardViewWithBreadcrumbs(m.BaseIntent, breadcrumbs...)

	// Handle progress modal for execution state
	if m.data.CurrentState == BulkExecuteState && m.data.IsExecuting {
		progress := m.data.GetProgress()
		message := fmt.Sprintf("Processing: %d/%d items", m.data.ProcessedCount, m.data.AffectedItemCount)
		message += fmt.Sprintf("\nSuccess: %d | Failed: %d | Skipped: %d",
			m.data.SuccessCount, m.data.FailureCount, m.data.SkippedCount)

		if m.data.IsPaused {
			message += "\n\n⏸ PAUSED - Press 'p' to resume or 'c' to complete"
		}

		m.SetProgress(fmt.Sprintf("Bulk %s", m.getOperationDisplayName()), message, progress)
	}

	// Handle errors
	if len(m.data.Errors) > 0 && m.data.CurrentState != BulkCompleteState {
		m.SetError(fmt.Errorf("operation error: %s", m.data.Errors[0]))
	}

	// Get content for current state
	content := m.getStateContent()
	view.WithContent(content)

	// Get context-aware help
	help := m.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

func (m *BulkOperationsModel) Result() *IntentResult[interface{}] {
	if m.result == nil {
		return nil
	}
	return &IntentResult[interface{}]{
		Status: m.result.Status,
		Error:  m.result.Error,
		Data:   m.result.Data,
	}
}

// Helper methods

func (m *BulkOperationsModel) getOperationDisplayName() string {
	switch m.data.SelectedOp {
	case "delete":
		return "Delete"
	case "tag":
		return "Tag"
	case "archive":
		return "Archive"
	case "export":
		return "Export"
	default:
		// Capitalize first letter only
		if len(m.data.SelectedOp) == 0 {
			return ""
		}
		return strings.ToUpper(m.data.SelectedOp[:1]) + m.data.SelectedOp[1:]
	}
}

func (m *BulkOperationsModel) getStateContent() string {
	switch m.data.CurrentState {
	case BulkSelectOpState:
		return m.getSelectOpContent()
	case BulkConfigureState:
		return m.getConfigureContent()
	case BulkExecuteState:
		return m.getExecuteContent()
	case BulkCompleteState:
		return m.getCompleteContent()
	default:
		return "Unknown state"
	}
}

func (m *BulkOperationsModel) getContextHelp() string {
	theme := m.Theme()

	switch m.data.CurrentState {
	case BulkSelectOpState:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Select operation"),
				components.BackBadge(),
			),
			ThemedCustomFooter(theme, components.QuitBadge()),
		)
	case BulkConfigureState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Execute"),
				components.BackBadge(),
			),
			ThemedCustomFooter(theme, components.QuitBadge()),
		)
	case BulkExecuteState:
		if m.data.IsPaused {
			return CombineThemedFooters(
				ThemedCustomFooter(theme,
					components.NewKeyBadge("p", "Resume"),
					components.NewKeyBadge("c", "Complete now"),
				),
				ThemedCustomFooter(theme, components.QuitBadge()),
			)
		}
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("p", "Pause"),
				components.NewKeyBadge("c", "Complete now"),
			),
			ThemedCustomFooter(theme, components.QuitBadge()),
		)
	case BulkCompleteState:
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

// Content methods

func (m *BulkOperationsModel) getSelectOpContent() string {
	var content strings.Builder

	content.WriteString("🔧 Select Bulk Operation\n\n")

	if len(m.data.AvailableOps) > 0 {
		content.WriteString("Available operations:\n\n")

		for i, op := range m.data.AvailableOps {
			selected := ""
			if m.data.SelectedOp == op {
				selected = " ▶ "
			} else {
				selected = "   "
			}

			// Add operation descriptions
			desc := m.getOperationDescription(op)
			// Capitalize first letter
			opTitle := op
			if len(op) > 0 {
				opTitle = strings.ToUpper(op[:1]) + op[1:]
			}
			content.WriteString(fmt.Sprintf("%s%d. %s - %s\n", selected, i+1, opTitle, desc))
		}
	} else {
		content.WriteString("No operations available.\n")
	}

	content.WriteString("\nUse arrow keys to navigate and Enter to select.\n")

	return content.String()
}

func (m *BulkOperationsModel) getOperationDescription(op string) string {
	switch op {
	case "delete":
		return "Permanently remove selected items"
	case "tag":
		return "Add or modify tags on selected items"
	case "archive":
		return "Move items to archive"
	case "export":
		return "Export selected items to file"
	default:
		return "Perform operation on selected items"
	}
}

func (m *BulkOperationsModel) getConfigureContent() string {
	var content strings.Builder

	content.WriteString("⚙️  Configure Operation\n\n")

	content.WriteString(fmt.Sprintf("Operation: %s\n", m.getOperationDisplayName()))
	content.WriteString(fmt.Sprintf("Scope: %s\n", m.data.ScopeType))
	content.WriteString(fmt.Sprintf("Items to process: %d\n\n", m.data.AffectedItemCount))

	// Add operation-specific configuration info
	switch m.data.SelectedOp {
	case "delete":
		content.WriteString("⚠️  Warning: This operation cannot be undone!\n")
		content.WriteString("All selected items will be permanently removed.\n")
	case "tag":
		content.WriteString("Tags will be added to all selected items.\n")
		content.WriteString("Existing tags will be preserved.\n")
	case "archive":
		content.WriteString("Items will be moved to the archive.\n")
		content.WriteString("You can restore them later if needed.\n")
	case "export":
		content.WriteString("Items will be exported in CSV format.\n")
		content.WriteString("You'll be prompted for the output location.\n")
	}

	content.WriteString("\nPress Enter to execute or Esc to go back.\n")

	return content.String()
}

func (m *BulkOperationsModel) getExecuteContent() string {
	// Most execution info is in the progress modal
	var content strings.Builder

	content.WriteString("⚙️  Bulk Operation in Progress\n\n")

	if m.data.IsPaused {
		content.WriteString("Operation is currently paused.\n")
		content.WriteString("Press 'p' to resume or 'c' to complete with current progress.\n")
	} else {
		content.WriteString("Processing items...\n")
		content.WriteString("Progress details are shown above.\n")
	}

	return content.String()
}

func (m *BulkOperationsModel) getCompleteContent() string {
	var content strings.Builder

	// Determine icon based on results
	if m.data.FailureCount > 0 {
		content.WriteString("⚠️  Operation Complete with Errors\n\n")
	} else {
		content.WriteString("✅ Operation Complete\n\n")
	}

	content.WriteString(fmt.Sprintf("Operation: %s\n", m.getOperationDisplayName()))
	content.WriteString(fmt.Sprintf("Total processed: %d items\n\n", m.data.ProcessedCount))

	// Results breakdown
	content.WriteString("Results:\n")
	content.WriteString(fmt.Sprintf("  ✓ Successful: %d\n", m.data.SuccessCount))
	content.WriteString(fmt.Sprintf("  ✗ Failed: %d\n", m.data.FailureCount))
	content.WriteString(fmt.Sprintf("  ⊘ Skipped: %d\n", m.data.SkippedCount))

	// Show errors if any
	if len(m.data.Errors) > 0 {
		content.WriteString("\n📝 Error Summary:\n")
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

// State handlers

func (m *BulkOperationsModel) handleSelectOpState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			m.ToggleHelp()
			return nil
		case KeyBack:
			// At root state, back means cancel and return to main menu
			m.result = &IntentResult[*BulkOperationsResult]{
				Status: Cancelled,
				Data: &BulkOperationsResult{
					Operation: "cancelled",
				},
			}
			m.active = false
			return nil
		}

		// Try list navigation handler first (handles up/down/j/k/pgup/pgdn/home/end/g/G)
		if m.navHandler.HandleKey(msg.String()) {
			return nil
		}

		switch msg.String() {
		case "enter":
			if m.data.SelectedOp != "" {
				m.data.CurrentState = BulkConfigureState
			}

		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Allow number selection
			idx := int(msg.String()[0] - '1')
			if idx >= 0 && idx < len(m.data.AvailableOps) {
				m.data.SelectOperation(m.data.AvailableOps[idx])
			}
		}
	}
	return nil
}

func (m *BulkOperationsModel) handleConfigureState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			m.ToggleHelp()
			return nil
		case KeyBack:
			// Go back to select operation state
			m.data.CurrentState = BulkSelectOpState
			return nil
		}

		switch msg.String() {
		case "enter":
			m.data.StartExecution()
			m.data.CurrentState = BulkExecuteState
		}
	}
	return nil
}

func (m *BulkOperationsModel) handleExecuteState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help)
		// Note: esc doesn't go back during execution - use 'c' to cancel
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			m.data.CompleteExecution()
			m.result = &IntentResult[*BulkOperationsResult]{
				Status: Cancelled,
				Data: &BulkOperationsResult{
					Operation:      m.data.SelectedOp,
					ProcessedCount: m.data.ProcessedCount,
					SuccessCount:   m.data.SuccessCount,
					FailureCount:   m.data.FailureCount,
					SkippedCount:   m.data.SkippedCount,
					Message:        "Operation cancelled",
				},
			}
			return tea.Quit
		case KeyHelp:
			m.ToggleHelp()
			return nil
		}

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

// ListNavigator interface implementation

// GetTotalItems returns the total number of available operations.
func (m *BulkOperationsModel) GetTotalItems() int {
	return len(m.data.AvailableOps)
}

// GetSelectedIndex returns the current selection index.
func (m *BulkOperationsModel) GetSelectedIndex() int {
	if m.data.SelectedOp == "" {
		return 0
	}
	for i, op := range m.data.AvailableOps {
		if op == m.data.SelectedOp {
			return i
		}
	}
	return 0
}

// SetSelectedIndex sets the selection index and updates the selected operation.
func (m *BulkOperationsModel) SetSelectedIndex(idx int) {
	if len(m.data.AvailableOps) == 0 {
		m.data.SelectedOp = ""
		return
	}

	// Clamp index to valid range
	if idx < 0 {
		idx = 0
	}
	if idx >= len(m.data.AvailableOps) {
		idx = len(m.data.AvailableOps) - 1
	}

	m.data.SelectedOp = m.data.AvailableOps[idx]
}

// GetPageSize returns the page size for pagination.
func (m *BulkOperationsModel) GetPageSize() int {
	return 10
}
