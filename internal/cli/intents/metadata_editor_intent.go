package intents

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MetadataEditorModel struct {
	data   *MetadataEditorContext
	result *IntentResult[*MetadataEditorResult]
}

func NewMetadataEditorIntent(data *MetadataEditorContext) *MetadataEditorModel {
	return &MetadataEditorModel{
		data: data,
		result: nil,
	}
}

func (m *MetadataEditorModel) Init() tea.Cmd {
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
	switch m.data.CurrentState {
	case MetadataReviewState:
		return m.viewReview()
	case MetadataEditState:
		return m.viewEdit()
	case MetadataConfirmState:
		return m.viewConfirm()
	}
	return "Unknown state"
}

func (m *MetadataEditorModel) Result() *IntentResult[interface{}] {
	return &IntentResult[interface{}]{
		Status: m.result.Status,
		Error:  m.result.Error,
	}
}

func (m *MetadataEditorModel) handleReviewState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.result = &IntentResult[*MetadataEditorResult]{
				Status: Cancelled,
				Data: &MetadataEditorResult{
					Action: "cancelled",
				},
			}
			return tea.Quit

		case "e":
			m.data.CurrentState = MetadataEditState

		case "esc":
			return tea.Quit
		}
	}
	return nil
}

func (m *MetadataEditorModel) handleEditState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			if m.data.HasChanges() {
				m.data.CurrentState = MetadataConfirmState
			}

		case "esc":
			m.data.ResetChanges()
			m.data.CurrentState = MetadataReviewState
		}
	}
	return nil
}

func (m *MetadataEditorModel) handleConfirmState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
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

		case "n", "esc":
			m.data.CurrentState = MetadataEditState
		}
	}
	return nil
}

func (m *MetadataEditorModel) viewReview() string {
	output := "Metadata Review\n"
	output += "=================================================================\n"
	output += fmt.Sprintf("Entity: %s (%s)\n", m.data.EntityType, m.data.EntityID)
	output += fmt.Sprintf("Fields: %d\n", len(m.data.OriginalMetadata))
	output += "\nOptions: e (edit), esc (back), q (quit)\n"
	return output
}

func (m *MetadataEditorModel) viewEdit() string {
	output := "Edit Metadata\n"
	output += "=================================================================\n"
	output += fmt.Sprintf("Entity: %s (%s)\n", m.data.EntityType, m.data.EntityID)
	output += fmt.Sprintf("Changed fields: %d\n", len(m.data.ChangedFields))

	if m.data.HasFormErrors() {
		output += "\nErrors:\n"
		for field, err := range m.data.FormErrors {
			output += fmt.Sprintf("  %s: %s\n", field, err)
		}
	}

	output += "\nOptions: Ctrl+S (save), Esc (cancel)\n"
	return output
}

func (m *MetadataEditorModel) viewConfirm() string {
	output := "Confirm Changes\n"
	output += "=================================================================\n"
	output += fmt.Sprintf("Save %d changed field(s)?\n", len(m.data.ChangedFields))
	output += "\nOptions: y (confirm), n (cancel), esc (back)\n"
	return output
}
