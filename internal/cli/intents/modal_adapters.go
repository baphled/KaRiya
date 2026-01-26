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

// NewErrorModalAdapter creates a new adapter for an error modal.
// If theme is nil, a default theme will be used for rendering.
func NewErrorModalAdapter(modal *feedback.Modal, width, height int, theme themes.Theme) *ErrorModalAdapter {
	return &ErrorModalAdapter{
		modal:  modal,
		width:  width,
		height: height,
		theme:  theme,
	}
}

// IsVisible returns true if the modal is not nil.
// Error modals are visible as long as they exist (nil check).
func (a *ErrorModalAdapter) IsVisible() bool {
	return a.modal != nil
}

// View returns the modal's rendered content using the stored dimensions.
func (a *ErrorModalAdapter) View() string {
	if a.modal == nil {
		return ""
	}
	return a.modal.Render(a.width, a.height)
}

// HandleUpdate processes keyboard input for the error modal.
// Error modals only respond to Esc to dismiss.
func (a *ErrorModalAdapter) HandleUpdate(msg tea.Msg) ModalUpdateResult {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
		return ModalUpdateResult{Closed: true}
	}
	return ModalUpdateResult{}
}

// RenderOverlay provides special rendering for error modals.
// Uses feedback.RenderOverlay for dimming and positioning behavior.
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

// NewFormModalAdapter creates a new adapter for a form-based modal.
//
// Parameters:
//   - isVisible: Function returning the modal's visibility state
//   - view: Function returning the modal's rendered content
//   - update: Function handling message updates with form data return
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

// IsVisible delegates to the isVisible function.
func (a *FormModalAdapter[T]) IsVisible() bool {
	return a.isVisible()
}

// View delegates to the view function.
func (a *FormModalAdapter[T]) View() string {
	return a.view()
}

// HandleUpdate calls the modal's Update and normalizes the result.
// The Closed state is detected by checking visibility after the update.
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

// NewConfirmModalAdapter creates a new adapter for a confirm modal.
func NewConfirmModalAdapter(modal *feedback.ConfirmModal) *ConfirmModalAdapter {
	return &ConfirmModalAdapter{modal: modal}
}

// IsVisible returns true if modal exists and is visible.
func (a *ConfirmModalAdapter) IsVisible() bool {
	return a.modal != nil && a.modal.IsVisible()
}

// View returns the modal's rendered content.
func (a *ConfirmModalAdapter) View() string {
	if a.modal == nil {
		return ""
	}
	return a.modal.View()
}

// HandleUpdate processes the confirm modal's Update and normalizes the result.
// The Data field contains the confirmation result (bool).
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

// NewViewModalAdapter creates a new adapter for a view-only modal.
//
// Parameters:
//   - isVisible: Function returning the modal's visibility state
//   - view: Function returning the modal's rendered content
//   - update: Function handling message updates (tea.Model pattern)
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

// IsVisible delegates to the isVisible function.
func (a *ViewModalAdapter) IsVisible() bool {
	return a.isVisible()
}

// View delegates to the view function.
func (a *ViewModalAdapter) View() string {
	return a.view()
}

// HandleUpdate calls the modal's Update and normalizes the result.
// View modals don't have Applied (no form submission).
func (a *ViewModalAdapter) HandleUpdate(msg tea.Msg) ModalUpdateResult {
	_, cmd := a.update(msg)
	closed := !a.isVisible()
	return ModalUpdateResult{
		Cmd:    cmd,
		Closed: closed,
	}
}
