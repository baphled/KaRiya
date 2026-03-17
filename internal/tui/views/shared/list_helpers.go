// Package shared provides shared helper functions for view implementations.
package shared

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// HandleListViewUpdate processes common list view messages (WindowSizeMsg, KeyMsg).
// It delegates terminal resizing and key handling to the caller's functions.
//
// Expected:
//   - baseView must be a valid BaseView pointer.
//   - msg must be a valid tea.Msg type.
//   - handleKey must be a function that processes KeyMsg and returns command and result.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - widgets.ViewResult: result indicating user action, or nil for internal state changes.
//
// Side effects:
//   - May update terminal dimensions via baseView.SetTerminalInfo.
func HandleListViewUpdate(
	baseView *widgets.BaseView,
	msg tea.Msg,
	handleKey func(tea.KeyMsg) (tea.Cmd, widgets.ViewResult),
) (tea.Cmd, widgets.ViewResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		baseView.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil
	case tea.KeyMsg:
		return handleKey(msg)
	}
	return nil, nil
}
