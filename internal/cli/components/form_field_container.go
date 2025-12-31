package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/baphled/kariya/internal/cli/styles"
)

// FormFieldContainer renders a form field with label, input, error, and hint sections.
// It provides consistent form field layout across the application with support for
// validation states and helpful hints.
// FormFieldContainer is a stateless rendering component.
type FormFieldContainer struct {
	label    string
	input    string
	error    string
	hint     string
	isFocused bool
	hasLabel bool
	hasInput bool
	hasError bool
	hasHint  bool
}

// NewFormFieldContainer creates a new FormFieldContainer.
func NewFormFieldContainer() *FormFieldContainer {
	return &FormFieldContainer{
		isFocused: false,
		hasLabel:  false,
		hasInput:  false,
		hasError:  false,
		hasHint:   false,
	}
}

// SetLabel sets the label for the form field.
// This method uses the builder pattern to allow method chaining.
func (ffc *FormFieldContainer) SetLabel(label string) *FormFieldContainer {
	ffc.label = label
	ffc.hasLabel = true
	return ffc
}

// SetInput sets the input field content.
// This method uses the builder pattern to allow method chaining.
func (ffc *FormFieldContainer) SetInput(input string) *FormFieldContainer {
	ffc.input = input
	ffc.hasInput = true
	return ffc
}

// SetError sets the error message for the field.
// This method uses the builder pattern to allow method chaining.
func (ffc *FormFieldContainer) SetError(error string) *FormFieldContainer {
	ffc.error = error
	ffc.hasError = true
	return ffc
}

// SetHint sets the hint/help text for the field.
// This method uses the builder pattern to allow method chaining.
func (ffc *FormFieldContainer) SetHint(hint string) *FormFieldContainer {
	ffc.hint = hint
	ffc.hasHint = true
	return ffc
}

// SetFocused sets whether the field is currently focused.
// This method uses the builder pattern to allow method chaining.
func (ffc *FormFieldContainer) SetFocused(focused bool) *FormFieldContainer {
	ffc.isFocused = focused
	return ffc
}

// Render returns the styled form field container as a string.
// It combines label, input, error, and hint sections with appropriate styling.
func (ffc *FormFieldContainer) Render() string {
	var parts []string

	// Render label if present
	if ffc.hasLabel {
		labelStyle := styles.InputLabel.Copy().
			Foreground(styles.ColorTextSecondary)
		parts = append(parts, labelStyle.Render(ffc.label))
	}

	// Render input if present
	if ffc.hasInput {
		var inputStyle lipgloss.Style
		if ffc.hasError {
			inputStyle = styles.InputError.Copy().
				Foreground(styles.ColorTextPrimary)
		} else if ffc.isFocused {
			inputStyle = styles.InputFocused.Copy().
				Foreground(styles.ColorTextPrimary)
		} else {
			inputStyle = styles.InputBase.Copy().
				Foreground(styles.ColorTextPrimary)
		}
		parts = append(parts, inputStyle.Render(ffc.input))
	}

	// Render error if present
	if ffc.hasError {
		errorStyle := styles.ErrorText.Copy().
			Foreground(styles.ColorError).
			MarginTop(1)
		parts = append(parts, errorStyle.Render(ffc.error))
	}

	// Render hint if present
	if ffc.hasHint {
		hintStyle := styles.InputHint.Copy().
			Foreground(styles.ColorTextMuted)
		parts = append(parts, hintStyle.Render(ffc.hint))
	}

	return strings.Join(parts, "\n")
}

