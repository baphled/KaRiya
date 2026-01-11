package models

import (
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// CaptureForm defines the interface for capture event forms.
// Both FormModel (legacy) and HuhFormModel (new) implement this interface.
type CaptureForm interface {
	// Init initializes the form.
	Init() tea.Cmd

	// Update handles messages and updates form state.
	Update(msg tea.Msg) (tea.Model, tea.Cmd)

	// View renders the form.
	View() string

	// SubmitForm triggers form submission.
	SubmitForm() tea.Cmd

	// SetStrategy sets the capture strategy ("quick" or "manual").
	SetStrategy(strategy string)

	// GetStrategy returns the current capture strategy.
	GetStrategy() string

	// Reset resets the form to its initial state.
	Reset()

	// LoadEventForEditing populates the form with an existing event.
	LoadEventForEditing(event *career.CareerEvent)

	// IsEditMode returns whether the form is in edit mode.
	IsEditMode() bool

	// GetEditEventID returns the ID of the event being edited.
	GetEditEventID() string

	// ShowOptionalFields returns whether optional fields are shown.
	ShowOptionalFields() bool

	// ToggleOptionalFields toggles optional field visibility.
	ToggleOptionalFields()
}

// Ensure both implementations satisfy the interface.
var _ CaptureForm = (*FormModel)(nil)
var _ CaptureForm = (*HuhFormModel)(nil)
