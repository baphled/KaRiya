package forms

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// WizardFormAdapter wraps a *huh.Form to provide step tracking and state
// queries without exposing huh types to consumers. It satisfies the
// behaviors.WizardForm interface via duck typing (no import needed).
//
// The adapter handles:
//   - Form initialization, update, and view delegation
//   - Completion and abort state checking via forms.IsCompleted/IsAborted
//   - Manual current step tracking (huh does not expose group index)
//   - Dimension updates for responsive resizing
type WizardFormAdapter struct {
	form        *huh.Form
	currentStep int
	totalSteps  int
}

// NewWizardFormAdapter creates an adapter wrapping the given huh form.
// totalSteps should match the number of groups in the form.
func NewWizardFormAdapter(form *huh.Form, totalSteps int) *WizardFormAdapter {
	return &WizardFormAdapter{
		form:       form,
		totalSteps: totalSteps,
	}
}

// Init initializes the underlying huh form.
func (a *WizardFormAdapter) Init() tea.Cmd {
	if a.form == nil {
		return nil
	}
	return a.form.Init()
}

// Update delegates the message to the huh form and returns the resulting command.
func (a *WizardFormAdapter) Update(msg tea.Msg) tea.Cmd {
	if a.form == nil {
		return nil
	}

	model, cmd := a.form.Update(msg)
	if f, ok := model.(*huh.Form); ok {
		a.form = f
	}

	return cmd
}

// View returns the huh form's rendered view.
func (a *WizardFormAdapter) View() string {
	if a.form == nil {
		return ""
	}
	return a.form.View()
}

// IsCompleted returns true when the huh form has been completed.
func (a *WizardFormAdapter) IsCompleted() bool {
	if a.form == nil {
		return false
	}
	return IsCompleted(a.form)
}

// IsAborted returns true when the huh form has been aborted.
func (a *WizardFormAdapter) IsAborted() bool {
	if a.form == nil {
		return false
	}
	return IsAborted(a.form)
}

// CurrentStep returns the manually-tracked current step index (0-based).
func (a *WizardFormAdapter) CurrentStep() int {
	return a.currentStep
}

// TotalSteps returns the total number of steps in the wizard.
func (a *WizardFormAdapter) TotalSteps() int {
	return a.totalSteps
}

// SetCurrentStep manually sets the current step index, clamped to valid range.
func (a *WizardFormAdapter) SetCurrentStep(step int) {
	if step < 0 {
		step = 0
	}
	maxStep := a.totalSteps - 1
	if maxStep < 0 {
		maxStep = 0
	}
	if step > maxStep {
		step = maxStep
	}
	a.currentStep = step
}

// SetDimensions updates the form's width and height.
func (a *WizardFormAdapter) SetDimensions(width, height int) {
	if a.form == nil {
		return
	}
	a.form = a.form.
		WithWidth(width).
		WithHeight(DefaultFormHeight(height))
}

// Form returns the underlying *huh.Form for direct access when needed.
func (a *WizardFormAdapter) Form() *huh.Form {
	return a.form
}

// SetForm replaces the underlying form and resets the step counter.
func (a *WizardFormAdapter) SetForm(form *huh.Form, totalSteps int) {
	a.form = form
	a.totalSteps = totalSteps
	a.currentStep = 0
}
