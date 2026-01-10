package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/charmbracelet/lipgloss"
)

// ModalContainer renders a modal dialog with title, message, buttons, and instructions.
// It provides consistent modal layout with support for destructive actions and
// centered rendering.
// ModalContainer is a stateless rendering component.
type ModalContainer struct {
	title           string
	message         string
	buttons         []string
	instructions    string
	isDestructive   bool
	hasTitle        bool
	hasMessage      bool
	hasButtons      bool
	hasInstructions bool
	theme           themes.Theme
}

// NewModalContainer creates a new ModalContainer.
func NewModalContainer() *ModalContainer {
	return &ModalContainer{
		isDestructive:   false,
		hasTitle:        false,
		hasMessage:      false,
		hasButtons:      false,
		hasInstructions: false,
	}
}

// SetTitle sets the title for the modal.
// This method uses the builder pattern to allow method chaining.
func (mc *ModalContainer) SetTitle(title string) *ModalContainer {
	mc.title = title
	mc.hasTitle = true
	return mc
}

// SetMessage sets the message content for the modal.
// This method uses the builder pattern to allow method chaining.
func (mc *ModalContainer) SetMessage(message string) *ModalContainer {
	mc.message = message
	mc.hasMessage = true
	return mc
}

// SetButtons sets the button labels for the modal.
// This method uses the builder pattern to allow method chaining.
func (mc *ModalContainer) SetButtons(buttons []string) *ModalContainer {
	mc.buttons = buttons
	mc.hasButtons = true
	return mc
}

// SetInstructions sets the instruction text for the modal.
// This method uses the builder pattern to allow method chaining.
func (mc *ModalContainer) SetInstructions(instructions string) *ModalContainer {
	mc.instructions = instructions
	mc.hasInstructions = true
	return mc
}

// WithDestructiveStyle marks the modal as destructive (e.g., for delete confirmations).
// This changes the styling to use error colors.
// This method uses the builder pattern to allow method chaining.
func (mc *ModalContainer) WithDestructiveStyle() *ModalContainer {
	mc.isDestructive = true
	return mc
}

// WithTheme sets the theme for the modal.
// This method uses the builder pattern to allow method chaining.
func (mc *ModalContainer) WithTheme(theme themes.Theme) *ModalContainer {
	mc.theme = theme
	return mc
}

// Theme helper methods for consistent themed styling.

// getPrimaryColor returns the primary text color from theme or fallback.
func (mc *ModalContainer) getPrimaryColor() lipgloss.Color {
	if mc.theme != nil {
		return mc.theme.ForegroundColor()
	}
	return styles.ColorTextPrimary
}

// getAccentColor returns the accent color from theme or fallback.
func (mc *ModalContainer) getAccentColor() lipgloss.Color {
	if mc.theme != nil {
		return mc.theme.PrimaryColor()
	}
	return styles.ColorAccentTeal
}

// getMutedColor returns the muted text color from theme or fallback.
func (mc *ModalContainer) getMutedColor() lipgloss.Color {
	if mc.theme != nil {
		return mc.theme.MutedColor()
	}
	return styles.ColorTextMuted
}

// getErrorColor returns the error color from theme or fallback.
func (mc *ModalContainer) getErrorColor() lipgloss.Color {
	if mc.theme != nil {
		return mc.theme.ErrorColor()
	}
	return styles.ColorError
}

// Render returns the styled modal container as a string.
// It combines title, message, buttons, and instructions with appropriate styling.
func (mc *ModalContainer) Render() string {
	var parts []string

	// Determine styles based on destructive flag
	var titleStyle, messageStyle, buttonStyle, instructionStyle lipgloss.Style

	if mc.isDestructive {
		titleStyle = styles.ModalDestructiveTitle.
			Foreground(mc.getErrorColor())
		messageStyle = lipgloss.NewStyle().
			Foreground(mc.getPrimaryColor()).
			MarginBottom(2)
		buttonStyle = lipgloss.NewStyle().
			Foreground(mc.getErrorColor()).
			Bold(true)
	} else {
		titleStyle = styles.ModalTitle.
			Foreground(mc.getPrimaryColor())
		messageStyle = styles.ModalMessage.
			Foreground(mc.getPrimaryColor())
		buttonStyle = lipgloss.NewStyle().
			Foreground(mc.getAccentColor()).
			Bold(true)
	}

	instructionStyle = styles.ModalInstructions.
		Foreground(mc.getMutedColor())

	// Render title if present
	if mc.hasTitle {
		parts = append(parts, titleStyle.Render(mc.title))
	}

	// Render message if present
	if mc.hasMessage {
		parts = append(parts, messageStyle.Render(mc.message))
	}

	// Render buttons if present
	if mc.hasButtons && len(mc.buttons) > 0 {
		buttonTexts := []string{}
		for _, btn := range mc.buttons {
			buttonTexts = append(buttonTexts, buttonStyle.Render(btn))
		}
		buttonContainer := styles.ModalButtonContainer.
			Render(strings.Join(buttonTexts, " "))
		parts = append(parts, buttonContainer)
	}

	// Render instructions if present
	if mc.hasInstructions {
		parts = append(parts, instructionStyle.Render(mc.instructions))
	}

	// Combine all parts
	content := strings.Join(parts, "\n")

	// Apply modal styling with border and background
	var modalStyle lipgloss.Style
	if mc.isDestructive {
		modalStyle = styles.ModalDestructive.
			Foreground(mc.getPrimaryColor())
	} else {
		modalStyle = styles.ModalBase.
			Foreground(mc.getPrimaryColor())
	}

	return modalStyle.Render(content)
}
