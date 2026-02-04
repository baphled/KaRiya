package forms

import (
	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/layout"
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
		if msg.String() == "esc" {
			fe.Cancelled = true
			return self, nil
		}

		if msg.String() == "q" || msg.String() == "ctrl+c" {
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
