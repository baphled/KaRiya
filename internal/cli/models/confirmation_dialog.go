package models

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmationDialog represents a confirmation dialog
type ConfirmationDialog struct {
	*BaseStandardModel
	title       string
	message     string
	confirmText string
	cancelText  string
	focused     bool // true = confirm, false = cancel
	confirmed   bool
	cancelled   bool
	helpFooter  components.HelpFooterModel
	modal       *components.ModalContainer
}

// NewConfirmationDialog creates a new confirmation dialog
func NewConfirmationDialog(title, message string) *ConfirmationDialog {
	d := &ConfirmationDialog{
		BaseStandardModel: NewBaseStandardModel(),
		title:             title,
		message:           message,
		confirmText:       "Yes, Delete",
		cancelText:        "Cancel",
		focused:           false, // Default to cancel for safety
		confirmed:         false,
		cancelled:         false,
		helpFooter:        components.NewHelpFooter("confirmation_dialog", 80),
	}
	d.modal = components.NewModalContainer()
	return d
}

// Init initializes the confirmation dialog
func (d *ConfirmationDialog) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (d *ConfirmationDialog) Update(msg tea.Msg) (*ConfirmationDialog, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "shift+tab":
			d.focused = false
		case "right", "l", "tab":
			d.focused = true
		case "enter":
			if d.focused {
				d.confirmed = true
			} else {
				d.cancelled = true
			}
			return d, nil
		case "esc", "q", "ctrl+c":
			d.cancelled = true
			return d, nil
		}
	}

	return d, nil
}

// View renders the confirmation dialog using ModalContainer
func (d *ConfirmationDialog) View() string {
	// Prepare buttons with focus styling
	buttons := d.renderButtons()

	// Configure modal container with destructive style
	d.modal.
		SetTitle(d.title).
		SetMessage(d.message).
		SetButtons(buttons).
		SetInstructions("Tab/←→: Switch | Enter: Confirm | Esc: Cancel").
		WithDestructiveStyle()

	// Render modal container
	dialog := d.modal.Render()

	// Add help footer
	d.helpFooter.SetWidth(80)
	return lipgloss.JoinVertical(lipgloss.Left, dialog, d.helpFooter.View())
}

// renderButtons returns the button labels with focus styling applied
func (d *ConfirmationDialog) renderButtons() []string {
	cancelStyle := styles.ButtonSecondary
	confirmStyle := styles.ButtonSecondary
	if d.focused {
		confirmStyle = styles.ButtonFocused
	} else {
		cancelStyle = styles.ButtonFocused
	}

	cancelButton := cancelStyle.Render(d.cancelText)
	confirmButton := confirmStyle.
		Foreground(styles.ColorError).
		Render(d.confirmText)

	return []string{cancelButton, confirmButton}
}

// IsConfirmed returns true if the user confirmed
func (d *ConfirmationDialog) IsConfirmed() bool {
	return d.confirmed
}

// IsCancelled returns true if the user cancelled
func (d *ConfirmationDialog) IsCancelled() bool {
	return d.cancelled
}

// Reset resets the dialog state
func (d *ConfirmationDialog) Reset() {
	d.focused = false
	d.confirmed = false
	d.cancelled = false
}
