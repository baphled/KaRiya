// Package behaviors provides reusable UI behavior components for the KaRiya TUI.

package behaviors

import (
	tea "github.com/charmbracelet/bubbletea"
)

// WizardForm abstracts multi-step form operations without exposing huh.
// Implementations wrap *huh.Form (via forms package) to provide step tracking,
// state queries, and dimension management.
//
// The interface enables:
//   - Testing with mock implementations
//   - Decoupling behaviors/ from huh library
//   - Consistent wizard form handling across wizard modals
type WizardForm interface {
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	View() string
	IsCompleted() bool
	IsAborted() bool
	CurrentStep() int
	TotalSteps() int
	SetDimensions(width, height int)
}

// WizardBehavior[T] provides state management, visibility control, and form
// delegation for multi-step wizard modals. It handles the common wizard lifecycle
// (visible -> completed/cancelled) while delegating form-specific logic to a
// WizardForm implementation.
//
// Usage:
//
//	form := forms.NewWizardFormAdapter(huhForm, 3)
//	data := &MyWizardData{}
//	wizard := behaviors.NewWizardBehavior(form, data)
//
//	// In Update:
//	cmd := wizard.Update(msg)
//	if wizard.IsCompleted() { config := wizard.Data() }
//
//	// In View:
//	return wizard.View()
type WizardBehavior[T any] struct {
	form      WizardForm
	data      *T
	visible   bool
	completed bool
	cancelled bool
	skipped   bool
}

// NewWizardBehavior creates a new wizard behavior wrapping the given form and data.
// The wizard starts visible and in an incomplete state.
func NewWizardBehavior[T any](form WizardForm, data *T) *WizardBehavior[T] {
	return &WizardBehavior[T]{
		form:    form,
		data:    data,
		visible: true,
	}
}

// Init initializes the wizard's form. Returns nil if the wizard is not visible.
func (w *WizardBehavior[T]) Init() tea.Cmd {
	if !w.visible {
		return nil
	}
	return w.form.Init()
}

// Update delegates the message to the form and checks for completion/abort.
// Returns nil without processing if the wizard is not visible.
func (w *WizardBehavior[T]) Update(msg tea.Msg) tea.Cmd {
	if !w.visible {
		return nil
	}

	cmd := w.form.Update(msg)

	if w.form.IsCompleted() {
		w.completed = true
		w.visible = false
	}

	if w.form.IsAborted() {
		w.cancelled = true
		w.visible = false
	}

	return cmd
}

// View returns the form's view, or empty string if the wizard is not visible.
func (w *WizardBehavior[T]) View() string {
	if !w.visible {
		return ""
	}
	return w.form.View()
}

// Data returns the wizard's data pointer.
func (w *WizardBehavior[T]) Data() *T {
	return w.data
}

// Form returns the current WizardForm.
func (w *WizardBehavior[T]) Form() WizardForm {
	return w.form
}

// SetForm replaces the current form with a new one.
func (w *WizardBehavior[T]) SetForm(form WizardForm) {
	w.form = form
}

// IsVisible returns whether the wizard is currently visible.
func (w *WizardBehavior[T]) IsVisible() bool {
	return w.visible
}

// IsCompleted returns whether the wizard finished successfully.
func (w *WizardBehavior[T]) IsCompleted() bool {
	return w.completed
}

// IsCancelled returns whether the wizard was cancelled.
func (w *WizardBehavior[T]) IsCancelled() bool {
	return w.cancelled
}

// IsSkipped returns whether the wizard was skipped.
func (w *WizardBehavior[T]) IsSkipped() bool {
	return w.skipped
}

// CurrentStep returns the form's current step index (0-based).
func (w *WizardBehavior[T]) CurrentStep() int {
	return w.form.CurrentStep()
}

// TotalSteps returns the form's total number of steps.
func (w *WizardBehavior[T]) TotalSteps() int {
	return w.form.TotalSteps()
}

// Show makes the wizard visible.
func (w *WizardBehavior[T]) Show() {
	w.visible = true
}

// Hide makes the wizard invisible.
func (w *WizardBehavior[T]) Hide() {
	w.visible = false
}

// Complete marks the wizard as completed and hides it.
func (w *WizardBehavior[T]) Complete() {
	w.completed = true
	w.visible = false
}

// Cancel marks the wizard as cancelled and hides it.
func (w *WizardBehavior[T]) Cancel() {
	w.cancelled = true
	w.visible = false
}

// Skip marks the wizard as skipped and completed, then hides it.
func (w *WizardBehavior[T]) Skip() {
	w.skipped = true
	w.completed = true
	w.visible = false
}

// Reset clears completion, cancellation, and skip flags, and shows the wizard.
func (w *WizardBehavior[T]) Reset() {
	w.completed = false
	w.cancelled = false
	w.skipped = false
	w.visible = true
}
