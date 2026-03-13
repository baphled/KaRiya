package forms

import (
	"github.com/baphled/kariya/internal/tui/navigation"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/layout"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// EditorFields holds the common mutable state shared by huh-based form editor modals.
type EditorFields struct {
	Form      *huh.Form
	Err       error
	Cancelled bool
	Width     int
	Height    int
	Theme     themes.Theme
}

// EditorUpdate handles common Update logic for huh form editor modals.
// The onComplete callback is invoked when the form completes successfully.
// The quitMsg factory produces the message to send when the user presses q/ctrl+c.
//
// Expected:
//   - fe must be a valid EditorFields pointer.
//   - self must be a valid tea.Model.
//   - onComplete must be a valid callback function.
//   - quitMsg must be a valid message factory function.
//
// Returns:
//   - tea.Model: the updated model.
//   - tea.Cmd: command to execute.
//
// Side effects:
//   - May update fe.Width, fe.Height, fe.Form, or fe.Cancelled fields.
func EditorUpdate(
	fe *EditorFields, self tea.Model, msg tea.Msg,
	onComplete func() (tea.Model, tea.Cmd),
	quitMsg func() tea.Msg,
) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		fe.Width = msg.Width
		fe.Height = msg.Height
		return self, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyEscape {
			fe.Cancelled = true
			return self, nil
		}

		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			return self, quitMsg
		}
	}

	form, cmd := fe.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		fe.Form = f
	}

	if IsCompleted(fe.Form) {
		return onComplete()
	}

	if IsAborted(fe.Form) {
		fe.Cancelled = true
		return self, nil
	}

	return self, cmd
}

// RenderEditorView renders the common form editor layout with header, form content, and footer.
//
// Expected:
//   - fe must be a valid EditorFields pointer.
//   - title must be a valid string.
//   - helpContext must be a valid string.
//
// Returns:
//   - A rendered string containing the form editor view.
//
// Side effects:
//   - None.
func RenderEditorView(fe *EditorFields, title string, helpContext string) string {
	formView := fe.Form.View()

	if fe.Err != nil {
		errorColor := fe.Theme.ErrorColor()
		errorStyle := lipgloss.NewStyle().
			Foreground(errorColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			MarginTop(1)

		formView += "\n\n" + errorStyle.Render(fe.Err.Error())
	}

	headerView := layout.NewHeader(title, fe.Width).
		WithTheme(fe.Theme).
		View()
	footerView := layout.NewFooter(fe.Width).
		WithTheme(fe.Theme).
		WithHelp(navigation.GetContextualHelp(helpContext)).
		View()

	contentStyle := lipgloss.NewStyle().
		Width(fe.Width).
		Padding(1, 2)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		contentStyle.Render(formView),
		"",
		footerView,
	)
}
