// Package helpers provides Page Object Pattern abstractions for BDD step definitions.
// This layer ensures navigation changes require updates in ONE place only.
package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
)

// FormHelper provides form navigation and interaction abstraction.
// If form navigation changes (Tab → Arrow keys), update methods here only.
type FormHelper struct {
	env *e2e.TestEnv
	ctx context.Context
}

// NewFormHelper creates a form helper for the given context.
func NewFormHelper(ctx context.Context) (*FormHelper, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return nil, godog.ErrPending
	}
	return &FormHelper{env: env, ctx: ctx}, nil
}

// NavigateToField tabs through form fields until the target label is visible.
// This centralizes field navigation - if navigation mechanism changes, update here only.
//
// Parameters:
//   - fieldLabel: The visible label of the field to navigate to
//
// Returns:
//   - error if field not found after max attempts
func (f *FormHelper) NavigateToField(fieldLabel string) error {
	const maxAttempts = 20

	for i := 0; i < maxAttempts; i++ {
		view := f.env.GetView()
		if strings.Contains(view, fieldLabel) {
			return nil // Field found and focused
		}
		f.env.Tab() // Navigate to next field
	}

	return fmt.Errorf("field '%s' not found after %d attempts", fieldLabel, maxAttempts)
}

// SetFieldValue navigates to a field and enters the specified value.
// Single place to change how values are entered if UI changes.
func (f *FormHelper) SetFieldValue(fieldLabel, value string) error {
	if err := f.NavigateToField(fieldLabel); err != nil {
		return err
	}
	f.env.TypeText(value)
	return nil
}

// SelectFieldOption navigates to a select field and chooses an option.
// Centralizes select field interaction.
func (f *FormHelper) SelectFieldOption(fieldLabel, option string) error {
	if err := f.NavigateToField(fieldLabel); err != nil {
		return err
	}

	// Navigate through options until we find the target
	const maxOptions = 10
	for i := 0; i < maxOptions; i++ {
		view := f.env.GetView()
		if strings.Contains(view, option) {
			f.env.Confirm() // Select this option
			return nil
		}
		f.env.NavigateDown() // Move to next option
	}

	return fmt.Errorf("option '%s' not found in field '%s'", option, fieldLabel)
}

// ToggleBoolean toggles a boolean field (checkbox/switch).
func (f *FormHelper) ToggleBoolean(fieldLabel string) error {
	if err := f.NavigateToField(fieldLabel); err != nil {
		return err
	}
	f.env.Confirm() // Toggle with Enter/Space
	return nil
}

// SubmitForm submits the form using Ctrl+S.
// Single place to change submission method.
func (f *FormHelper) SubmitForm() error {
	f.env.PressKey(tea.KeyCtrlS)
	return nil
}

// CancelForm cancels the form using Escape.
func (f *FormHelper) CancelForm() error {
	f.env.PressKey(tea.KeyEscape)
	return nil
}

// GetFieldValue returns the current field value (visible in view).
func (f *FormHelper) GetFieldValue(fieldLabel string) (string, error) {
	if err := f.NavigateToField(fieldLabel); err != nil {
		return "", err
	}
	// Return the current view content - caller can parse as needed
	return f.env.GetView(), nil
}

// IsFieldVisible checks if a field label is currently visible.
func (f *FormHelper) IsFieldVisible(fieldLabel string) bool {
	view := f.env.GetView()
	return strings.Contains(view, fieldLabel)
}
