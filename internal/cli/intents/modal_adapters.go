// Package intents provides intent implementations for the KaRiya TUI.
package intents

import (
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	themes "github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
)

// ============================================================================
// ErrorModalAdapter - Wraps feedback.Modal to implement ManagedModal
// ============================================================================

// ErrorModalAdapter wraps feedback.Modal to implement ManagedModal.
// This normalizes the Render(width, height) API to the standard View() pattern.
//
// The error modal uses a different API than other modals:
//   - feedback.Modal.Render(width, height) string
//
// This adapter stores dimensions and provides a View() method.
type ErrorModalAdapter struct {
	modal  *feedback.Modal
	width  int
	height int
	theme  themes.Theme
}

// NewErrorModalAdapter wraps an error modal with the ManagedModal interface for use with the modal registry.
// If theme is nil, a default theme will be used for rendering.
//
// Expected:
//   - modal may be nil, in which case IsVisible returns false.
//   - width and height are the terminal dimensions for rendering.
//   - theme may be nil; a default theme will be used if so.
//
// Returns:
//   - A configured ErrorModalAdapter wrapping the given modal.
//
// Side effects:
//   - None.
func NewErrorModalAdapter(modal *feedback.Modal, width, height int, theme themes.Theme) *ErrorModalAdapter {
	return &ErrorModalAdapter{
		modal:  modal,
		width:  width,
		height: height,
		theme:  theme,
	}
}

// IsVisible checks whether this modal should be rendered as an overlay in the current view cycle.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (a *ErrorModalAdapter) IsVisible() bool {
	return a.modal != nil
}

// View renders the modal content for overlay composition using the stored dimensions.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (a *ErrorModalAdapter) View() string {
	if a.modal == nil {
		return ""
	}
	return a.modal.Render(a.width, a.height)
}

// HandleUpdate processes keyboard input for the error modal.
// Error modals only respond to Esc to dismiss.
//
// Expected:
//   - msg is forwarded to the wrapped modal for processing.
//
// Returns:
//   - A ModalUpdateResult with Closed set to true if Esc was pressed.
//
// Side effects:
//   - None.
func (a *ErrorModalAdapter) HandleUpdate(msg tea.Msg) ModalUpdateResult {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
		return ModalUpdateResult{Closed: true}
	}
	return ModalUpdateResult{}
}

// RenderOverlay provides special rendering for error modals.
// Uses feedback.RenderOverlay for dimming and positioning behavior.
//
// Expected:
//   - baseView is the background content to overlay the modal onto.
//
// Returns:
//   - The composited view with the modal overlaid, or baseView if the modal is nil.
//
// Side effects:
//   - None.
func (a *ErrorModalAdapter) RenderOverlay(baseView string) string {
	if a.modal == nil {
		return baseView
	}
	modalContent := a.modal.Render(a.width, a.height)
	return feedback.RenderOverlay(baseView, modalContent, a.width, a.height, a.theme)
}

// ============================================================================
// FormModalAdapter - Generic adapter for form-based modals
// ============================================================================

// FormModalAdapter wraps form-based modals with signature:
//
//	Update(msg) -> (tea.Cmd, bool, *T)
//
// This is used for search, filter, sort, quickAdd, and edit modals.
// The generic type T represents the form's data structure.
type FormModalAdapter[T any] struct {
	isVisible func() bool
	view      func() string
	update    func(tea.Msg) (tea.Cmd, bool, T)
}

// NewFormModalAdapter wraps a form-based modal with the ManagedModal interface for use with the modal registry.
//
// Expected:
//   - isVisible must be non-nil; returns the modal's visibility state.
//   - view must be non-nil; returns the modal's rendered content.
//   - update must be non-nil; handles message updates and returns form data.
//
// Returns:
//   - A configured FormModalAdapter wrapping the provided functions.
//
// Side effects:
//   - None.
func NewFormModalAdapter[T any](
	isVisible func() bool,
	view func() string,
	update func(tea.Msg) (tea.Cmd, bool, T),
) *FormModalAdapter[T] {
	return &FormModalAdapter[T]{
		isVisible: isVisible,
		view:      view,
		update:    update,
	}
}

// IsVisible checks whether this modal should be rendered as an overlay in the current view cycle.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (a *FormModalAdapter[T]) IsVisible() bool {
	return a.isVisible()
}

// View renders the modal content for overlay composition.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (a *FormModalAdapter[T]) View() string {
	return a.view()
}

// HandleUpdate delegates the message to the wrapped modal and normalizes the result.
// The Closed state is detected by checking visibility after the update.
//
// Expected:
//   - msg is forwarded to the wrapped modal for processing.
//
// Returns:
//   - A ModalUpdateResult with Cmd, Closed, Applied, and Data fields populated.
//
// Side effects:
//   - May mutate the underlying modal's state via the update function.
func (a *FormModalAdapter[T]) HandleUpdate(msg tea.Msg) ModalUpdateResult {
	cmd, applied, data := a.update(msg)
	closed := !a.isVisible()
	return ModalUpdateResult{
		Cmd:     cmd,
		Closed:  closed,
		Applied: applied,
		Data:    data,
	}
}

// ============================================================================
// ConfirmModalAdapter - Wraps feedback.ConfirmModal
// ============================================================================

// ConfirmModalAdapter wraps feedback.ConfirmModal to implement ManagedModal.
// Confirm modals have signature: Update(msg) -> (tea.Cmd, bool)
type ConfirmModalAdapter struct {
	modal *feedback.ConfirmModal
}

// NewConfirmModalAdapter wraps a confirm modal with the ManagedModal interface for use with the modal registry.
//
// Expected:
//   - modal must be non-nil.
//
// Returns:
//   - A configured ConfirmModalAdapter wrapping the given modal.
//
// Side effects:
//   - None.
func NewConfirmModalAdapter(modal *feedback.ConfirmModal) *ConfirmModalAdapter {
	return &ConfirmModalAdapter{modal: modal}
}

// IsVisible checks whether this modal should be rendered as an overlay in the current view cycle.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (a *ConfirmModalAdapter) IsVisible() bool {
	return a.modal != nil && a.modal.IsVisible()
}

// View renders the modal content for overlay composition.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (a *ConfirmModalAdapter) View() string {
	if a.modal == nil {
		return ""
	}
	return a.modal.View()
}

// HandleUpdate delegates the message to the wrapped confirm modal and normalizes the result.
// The Data field contains the confirmation result (bool).
//
// Expected:
//   - msg is forwarded to the wrapped modal for processing.
//
// Returns:
//   - A ModalUpdateResult with Cmd, Closed, Applied, and Data fields populated.
//
// Side effects:
//   - May mutate the underlying confirm modal's state.
func (a *ConfirmModalAdapter) HandleUpdate(msg tea.Msg) ModalUpdateResult {
	cmd, confirmed := a.modal.Update(msg)
	closed := !a.modal.IsVisible()
	return ModalUpdateResult{
		Cmd:     cmd,
		Closed:  closed,
		Applied: confirmed,
		Data:    confirmed,
	}
}

// ============================================================================
// ViewModalAdapter - Wraps view-only modals (tea.Model pattern)
// ============================================================================

// ViewModalAdapter wraps view-only modals with signature:
//
//	Update(msg) -> (tea.Model, tea.Cmd)
//
// This is used for ViewEventDetailModal and ViewEventSkillsModal.
type ViewModalAdapter struct {
	isVisible func() bool
	view      func() string
	update    func(tea.Msg) (tea.Model, tea.Cmd)
}

// NewViewModalAdapter wraps a view-only modal with the ManagedModal interface for use with the modal registry.
//
// Expected:
//   - isVisible must be non-nil; returns the modal's visibility state.
//   - view must be non-nil; returns the modal's rendered content.
//   - update must be non-nil; handles message updates using the tea.Model pattern.
//
// Returns:
//   - A configured ViewModalAdapter wrapping the provided functions.
//
// Side effects:
//   - None.
func NewViewModalAdapter(
	isVisible func() bool,
	view func() string,
	update func(tea.Msg) (tea.Model, tea.Cmd),
) *ViewModalAdapter {
	return &ViewModalAdapter{
		isVisible: isVisible,
		view:      view,
		update:    update,
	}
}

// IsVisible checks whether this modal should be rendered as an overlay in the current view cycle.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (a *ViewModalAdapter) IsVisible() bool {
	return a.isVisible()
}

// View renders the modal content for overlay composition.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (a *ViewModalAdapter) View() string {
	return a.view()
}

// HandleUpdate delegates the message to the wrapped view modal and normalizes the result.
// View modals don't have Applied (no form submission).
//
// Expected:
//   - msg is forwarded to the wrapped modal for processing.
//
// Returns:
//   - A ModalUpdateResult with Cmd and Closed fields populated. Applied is always false.
//
// Side effects:
//   - May mutate the underlying modal's state via the update function.
func (a *ViewModalAdapter) HandleUpdate(msg tea.Msg) ModalUpdateResult {
	_, cmd := a.update(msg)
	closed := !a.isVisible()
	return ModalUpdateResult{
		Cmd:    cmd,
		Closed: closed,
	}
}
