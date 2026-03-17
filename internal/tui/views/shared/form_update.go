//nolint:revive,nolintlint // "shared" is intentional — cross-domain reusable views
package shared

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// HandleFormUpdate processes standard form messages and returns updated form state.
// It handles WindowSizeMsg (calls onResize and returns form.Init()), KeyEsc (returns
// CancelViewResult), form completion (calls onSubmit and returns SubmitViewResult),
// and form abortion (returns CancelViewResult).
//
// Expected:
//   - form must be a valid forms.Form pointer.
//   - msg must be a valid tea.Msg.
//   - onResize is called on WindowSizeMsg with width and height.
//   - onSubmit is called when form completes; must return the form data payload.
//
// Returns:
//   - Updated forms.Form, tea.Cmd, and ViewResult (nil while in progress).
//
// Side effects:
//   - Calls onResize and onSubmit callbacks when appropriate.
func HandleFormUpdate(
	form forms.Form,
	msg tea.Msg,
	onResize func(width, height int),
	onSubmit func() interface{},
) (forms.Form, tea.Cmd, widgets.ViewResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if onResize != nil {
			onResize(msg.Width, msg.Height)
		}
		cmd := form.Init()
		return form, cmd, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return form, nil, &widgets.CancelViewResult{}
		}
	}

	var cmd tea.Cmd
	form, cmd = forms.Update(form, msg)

	if forms.IsCompleted(form) {
		data := onSubmit()
		return form, cmd, &widgets.SubmitViewResult{FormData: data}
	}

	if forms.IsAborted(form) {
		return form, cmd, &widgets.CancelViewResult{}
	}

	return form, cmd, nil
}
