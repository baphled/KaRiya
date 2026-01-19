package intents

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

type MetadataEditorModel struct {
	*BaseIntent
	data   *MetadataEditorContext
	result *IntentResult[*MetadataEditorResult]
	active bool
}

func NewMetadataEditorIntent(data *MetadataEditorContext) *MetadataEditorModel {
	return &MetadataEditorModel{
		BaseIntent: NewBaseIntent(),
		data:       data,
		result:     nil,
		active:     false,
	}
}

func (m *MetadataEditorModel) Init() tea.Cmd {
	// Mark intent as active
	m.active = true

	// context already set in data
	m.data.CurrentState = MetadataReviewState

	// Return a no-op command to satisfy the intent lifecycle
	return func() tea.Msg { return nil }
}

func (m *MetadataEditorModel) Update(msg tea.Msg) tea.Cmd {
	switch m.data.CurrentState {
	case MetadataReviewState:
		return m.handleReviewState(msg)
	case MetadataEditState:
		return m.handleEditState(msg)
	case MetadataConfirmState:
		return m.handleConfirmState(msg)
	}
	return nil
}

func (m *MetadataEditorModel) View() string {
	if !m.active {
		return "MetadataEditor intent is not active"
	}

	// Create standard view with breadcrumbs
	view := m.CreateViewWithBreadcrumbs("Main Menu", "Edit Metadata", m.getStateName())

	// Handle validation errors
	if m.data.HasFormErrors() {
		var errorMessages []string
		for field, err := range m.data.FormErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", field, err))
		}
		m.SetError(fmt.Errorf("validation errors:\n%s", strings.Join(errorMessages, "\n")))
	}

	// Get content for current state
	content := m.getStateContent()
	view.WithContent(content)

	// Get context-aware help
	help := m.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

func (m *MetadataEditorModel) Result() *IntentResult[interface{}] {
	// Return nil when intent hasn't completed yet
	// Only return non-nil result when the intent has explicitly completed
	if m.result == nil {
		return nil
	}
	return &IntentResult[interface{}]{
		Status: m.result.Status,
		Error:  m.result.Error,
	}
}

// getStateName returns the display name for the current state
func (m *MetadataEditorModel) getStateName() string {
	switch m.data.CurrentState {
	case MetadataReviewState:
		return "Review"
	case MetadataEditState:
		return "Edit"
	case MetadataConfirmState:
		return "Confirm"
	default:
		return ""
	}
}

// getStateContent returns the content for the current state
func (m *MetadataEditorModel) getStateContent() string {
	switch m.data.CurrentState {
	case MetadataReviewState:
		return m.getReviewContent()
	case MetadataEditState:
		return m.getEditContent()
	case MetadataConfirmState:
		return m.getConfirmContent()
	default:
		return "Unknown state"
	}
}

// getContextHelp returns context-aware help text
func (m *MetadataEditorModel) getContextHelp() string {
	theme := m.Theme()

	switch m.data.CurrentState {
	case MetadataReviewState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("e", "Edit metadata", theme),
				primitives.BackBadge(theme),
			),
			ThemedCustomFooter(theme, primitives.QuitBadge(theme)),
		)
	case MetadataEditState:
		if m.data.HasChanges() {
			return CombineThemedFooters(
				ThemedFormFooter(theme),
				ThemedCustomFooter(theme,
					primitives.HelpKeyBadge("Ctrl+S", "Save changes", theme),
					primitives.CancelBadge(theme),
				),
				ThemedCustomFooter(theme, primitives.QuitBadge(theme)),
			)
		}
		return CombineThemedFooters(
			ThemedFormFooter(theme),
			ThemedCustomFooter(theme, primitives.BackBadge(theme)),
			ThemedCustomFooter(theme, primitives.QuitBadge(theme)),
		)
	case MetadataConfirmState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("y/Enter", "Confirm", theme),
				primitives.HelpKeyBadge("n/Esc", "Back", theme),
			),
			ThemedCustomFooter(theme, primitives.QuitBadge(theme)),
		)
	default:
		return ThemedCustomFooter(theme, primitives.QuitBadge(theme))
	}
}

// Content methods

func (m *MetadataEditorModel) getReviewContent() string {
	var content strings.Builder

	content.WriteString("📋 Metadata Review\n\n")

	if m.data.EntityType != "" && m.data.EntityID != "" {
		content.WriteString(fmt.Sprintf("Entity: %s\n", m.data.EntityType))
		content.WriteString(fmt.Sprintf("ID: %s\n\n", m.data.EntityID))
	}

	if len(m.data.OriginalMetadata) > 0 {
		content.WriteString(fmt.Sprintf("📝 Metadata Fields: %d\n\n", len(m.data.OriginalMetadata)))

		// Show current metadata
		for field, value := range m.data.OriginalMetadata {
			content.WriteString(fmt.Sprintf("  • %s: %v\n", field, value))
		}
	} else {
		content.WriteString("No metadata available.\n")
	}

	content.WriteString("\nPress 'e' to edit metadata or Esc to go back.\n")

	return content.String()
}

func (m *MetadataEditorModel) getEditContent() string {
	var content strings.Builder

	content.WriteString("✏️  Edit Metadata\n\n")

	if m.data.EntityType != "" && m.data.EntityID != "" {
		content.WriteString(fmt.Sprintf("Entity: %s (%s)\n\n", m.data.EntityType, m.data.EntityID))
	}

	// Show changed fields
	if m.data.HasChanges() {
		content.WriteString(fmt.Sprintf("📝 Changed fields: %d\n\n", len(m.data.ChangedFields)))

		changes := m.data.GetChanges()
		for field, newValue := range changes {
			originalValue := m.data.OriginalMetadata[field]
			content.WriteString(fmt.Sprintf("  • %s:\n", field))
			content.WriteString(fmt.Sprintf("    Before: %v\n", originalValue))
			content.WriteString(fmt.Sprintf("    After:  %v\n", newValue))
		}

		content.WriteString("\nPress Ctrl+S to save or Esc to cancel changes.\n")
	} else {
		content.WriteString("No changes made yet.\n\n")

		// Show all editable fields
		if len(m.data.EditedMetadata) > 0 {
			content.WriteString("Available fields:\n")
			for field, value := range m.data.EditedMetadata {
				content.WriteString(fmt.Sprintf("  • %s: %v\n", field, value))
			}
		}

		content.WriteString("\nMake changes and press Ctrl+S to save, or Esc to cancel.\n")
	}

	return content.String()
}

func (m *MetadataEditorModel) getConfirmContent() string {
	var content strings.Builder

	content.WriteString("💾 Confirm Changes\n\n")

	if len(m.data.ChangedFields) > 0 {
		content.WriteString(fmt.Sprintf("You are about to save %d changed field(s):\n\n", len(m.data.ChangedFields)))

		changes := m.data.GetChanges()
		for field, newValue := range changes {
			originalValue := m.data.OriginalMetadata[field]
			content.WriteString(fmt.Sprintf("  • %s: %v → %v\n", field, originalValue, newValue))
		}

		content.WriteString("\nDo you want to save these changes?\n")
	} else {
		content.WriteString("No changes to save.\n")
	}

	return content.String()
}

// State handlers

func (m *MetadataEditorModel) handleReviewState(msg tea.Msg) tea.Cmd {
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
			m.result = &IntentResult[*MetadataEditorResult]{
				Status: Cancelled,
				Data: &MetadataEditorResult{
					Action: "cancelled",
				},
			}
			m.active = false
			return nil
		}

		switch msg.String() {
		case "e":
			m.data.CurrentState = MetadataEditState
		}
	}
	return nil
}

func (m *MetadataEditorModel) handleEditState(msg tea.Msg) tea.Cmd {
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
			// Go back to review state
			m.data.ResetChanges()
			m.data.CurrentState = MetadataReviewState
			return nil
		}

		switch msg.String() {
		case "ctrl+s":
			if m.data.HasChanges() {
				m.data.CurrentState = MetadataConfirmState
			}
		}
	}
	return nil
}

func (m *MetadataEditorModel) handleConfirmState(msg tea.Msg) tea.Cmd {
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
			// Go back to edit state
			m.data.CurrentState = MetadataEditState
			return nil
		}

		switch msg.String() {
		case "y", "enter":
			m.result = &IntentResult[*MetadataEditorResult]{
				Status: Completed,
				Data: &MetadataEditorResult{
					Action:         "saved",
					EntityType:     m.data.EntityType,
					EntityID:       m.data.EntityID,
					Changes:        m.data.GetChanges(),
					OriginalValues: m.data.OriginalMetadata,
					Message:        "Metadata saved successfully",
				},
			}
			return tea.Quit

		case "n":
			m.data.CurrentState = MetadataEditState
		}
	}
	return nil
}
