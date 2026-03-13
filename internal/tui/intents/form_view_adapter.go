package intents

import (
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// FormViewAdapter adapts a widgets.View to the ManagedModal interface.
// It bridges the View content layer with the modal registry, allowing
// form-content Views to be displayed as managed modals.
type FormViewAdapter struct {
	view    widgets.View
	visible bool
	width   int
	height  int
}

// NewFormViewAdapter creates a FormViewAdapter wrapping the given view.
// The adapter starts hidden and returns empty results when view is nil.
//
// Expected: View may be nil, returning "" from View() and empty ModalUpdateResult from HandleUpdate.
//
// Returns: A configured FormViewAdapter wrapping the given view.
//
// Side effects: None.
func NewFormViewAdapter(view widgets.View) *FormViewAdapter {
	return &FormViewAdapter{view: view}
}

// IsVisible implements ManagedModal returning true when the adapter is shown.
//
// Returns: A bool indicating visibility state.
//
// Side effects: None.
func (a *FormViewAdapter) IsVisible() bool { return a.visible }

// Show makes the adapter visible in the registry overlay.
//
// Side effects: Sets the visible flag to true.
func (a *FormViewAdapter) Show() { a.visible = true }

// Hide makes the adapter invisible without affecting underlying view state.
//
// Side effects: Sets the visible flag to false.
func (a *FormViewAdapter) Hide() { a.visible = false }

// SetDimensions sets width and height hints for content rendering.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (a *FormViewAdapter) SetDimensions(w, h int) {
	a.width = w
	a.height = h
}

// View implements ManagedModal delegating to the wrapped view's RenderContent.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (a *FormViewAdapter) View() string {
	if a.view == nil {
		return ""
	}
	return a.view.RenderContent()
}

// HandleUpdate implements ManagedModal translating ViewResult to ModalUpdateResult.
//
// Expected: Msg is the Bubble Tea message to forward to the wrapped view.
//
// Returns: A ModalUpdateResult with fields populated based on the ViewResult type.
//
// Side effects: May hide the adapter on Submit or Cancel results.
func (a *FormViewAdapter) HandleUpdate(msg tea.Msg) ModalUpdateResult {
	if a.view == nil {
		return ModalUpdateResult{}
	}
	cmd, result := a.view.Update(msg)
	if result == nil {
		return ModalUpdateResult{Cmd: cmd}
	}
	switch result.Type() {
	case widgets.ResultSubmit:
		a.visible = false
		return ModalUpdateResult{Cmd: cmd, Applied: true, Closed: true, Data: result.Data()}
	case widgets.ResultCancel:
		a.visible = false
		return ModalUpdateResult{Cmd: cmd, Closed: true, Applied: false}
	case widgets.ResultNavigate:
		return ModalUpdateResult{Cmd: cmd, Data: result.Data()}
	case widgets.ResultError:
		return ModalUpdateResult{Cmd: cmd, Data: result.Data()}
	default:
		return ModalUpdateResult{Cmd: cmd}
	}
}
