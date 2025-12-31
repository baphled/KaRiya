package models

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmationDialog represents a confirmation dialog
type ConfirmationDialog struct {
	title       string
	message     string
	confirmText string
	cancelText  string
	focused     bool // true = confirm, false = cancel
	confirmed   bool
	cancelled   bool
	helpFooter  components.HelpFooterModel
}

// NewConfirmationDialog creates a new confirmation dialog
func NewConfirmationDialog(title, message string) *ConfirmationDialog {
	return &ConfirmationDialog{
		title:       title,
		message:     message,
		confirmText: "Yes, Delete",
		cancelText:  "Cancel",
		focused:     false, // Default to cancel for safety
		confirmed:   false,
		cancelled:   false,
		helpFooter:  components.NewHelpFooter("confirmation_dialog", 80),
	}
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
		case "esc", "q":
			d.cancelled = true
			return d, nil
		}
	}

	return d, nil
}

// View renders the confirmation dialog
func (d *ConfirmationDialog) View() string {
	// Title - use destructive style
	title := styles.ModalDestructiveTitle.Render(d.title)

	// Message - standard modal message style
	message := styles.ModalMessage.Render(d.message)

	// Buttons with consistent styling
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

	buttons := styles.ModalButtonContainer.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			cancelButton,
			"  ",
			confirmButton,
		),
	)

	// Instructions - standard modal instructions
	instructions := styles.ModalInstructions.Render(
		"Tab/←→: Switch | Enter: Confirm | Esc: Cancel",
	)

	// Combine all elements with consistent vertical spacing
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		message,
		buttons,
		instructions,
	)

	// Wrap in destructive modal box
	dialog := styles.ModalDestructive.
		Width(60).
		Render(content)

	d.helpFooter.SetWidth(80)
	return lipgloss.JoinVertical(lipgloss.Left, dialog, d.helpFooter.View())
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
