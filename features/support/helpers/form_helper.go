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
//
// Expected:
//   - ctx must contain a valid TestEnv set by the BDD environment.
//
// Returns:
//   - A FormHelper bound to the test environment, or error if env is nil.
//
// Side effects:
//   - None.
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
// Expected:
//   - fieldLabel must be a non-empty string matching a visible form label.
//
// Returns:
//   - error if field not found after max attempts.
//
// Side effects:
//   - Sends Tab key events to the test environment.
func (f *FormHelper) NavigateToField(fieldLabel string) error {
	const maxAttempts = 20

	for range maxAttempts {
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
//
// Expected:
//   - fieldLabel must match a visible form label.
//   - value must be a valid string to type into the field.
//
// Returns:
//   - error if the field cannot be found.
//
// Side effects:
//   - Sends Tab and text key events to the test environment.
func (f *FormHelper) SetFieldValue(fieldLabel, value string) error {
	if err := f.NavigateToField(fieldLabel); err != nil {
		return err
	}
	f.env.TypeText(value)
	return nil
}

// SelectFieldOption navigates to a select field and chooses an option.
// Centralizes select field interaction.
//
// Expected:
//   - fieldLabel must match a visible form label.
//   - option must match a visible option text in the select field.
//
// Returns:
//   - error if the field or option cannot be found.
//
// Side effects:
//   - Sends Tab, navigation, and confirm key events to the test environment.
func (f *FormHelper) SelectFieldOption(fieldLabel, option string) error {
	if err := f.NavigateToField(fieldLabel); err != nil {
		return err
	}

	// Navigate through options until we find the target
	const maxOptions = 10
	for range maxOptions {
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
//
// Expected:
//   - fieldLabel must match a visible boolean form field.
//
// Returns:
//   - error if the field cannot be found.
//
// Side effects:
//   - Sends Tab and confirm key events to the test environment.
func (f *FormHelper) ToggleBoolean(fieldLabel string) error {
	if err := f.NavigateToField(fieldLabel); err != nil {
		return err
	}
	f.env.Confirm() // Toggle with Enter/Space
	return nil
}

// SubmitForm submits the form using Ctrl+S.
// Single place to change submission method.
//
// Returns:
//   - nil always (error interface for consistency).
//
// Side effects:
//   - Sends Ctrl+S key event to the test environment.
func (f *FormHelper) SubmitForm() error {
	f.env.PressKey(tea.KeyCtrlS)
	return nil
}

// CancelForm cancels the form using Escape.
//
// Returns:
//   - nil always (error interface for consistency).
//
// Side effects:
//   - Sends Escape key event to the test environment.
func (f *FormHelper) CancelForm() error {
	f.env.PressKey(tea.KeyEscape)
	return nil
}

// GetFieldValue returns the current field value (visible in view).
//
// Expected:
//   - fieldLabel must match a visible form label.
//
// Returns:
//   - The current view content as a string, or error if field not found.
//
// Side effects:
//   - Sends Tab key events to navigate to the field.
func (f *FormHelper) GetFieldValue(fieldLabel string) (string, error) {
	if err := f.NavigateToField(fieldLabel); err != nil {
		return "", err
	}
	// Return the current view content - caller can parse as needed
	return f.env.GetView(), nil
}

// IsFieldVisible checks if a field label is currently visible.
//
// Expected:
//   - fieldLabel must be a non-empty string.
//
// Returns:
//   - true if the field label appears in the current view.
//
// Side effects:
//   - None.
func (f *FormHelper) IsFieldVisible(fieldLabel string) bool {
	view := f.env.GetView()
	return strings.Contains(view, fieldLabel)
}

// FieldExistsInForm checks if a field exists in the form by tabbing through all fields.
// This is less strict than IsFieldVisible - it will find fields even if off-screen.
//
// Expected:
//   - fieldLabel must be a non-empty string.
//   - maxFields must be a positive integer limiting the search depth.
//
// Returns:
//   - true if the field label is found within maxFields tab presses.
//
// Side effects:
//   - Sends Tab key events to the test environment.
func (f *FormHelper) FieldExistsInForm(fieldLabel string, maxFields int) bool {
	initialView := f.env.GetView()

	//nolint:intrange // Traditional loop required: i is used for cycle detection (i > 0 guard below)
	for i := 0; i < maxFields; i++ {
		view := f.env.GetView()
		if strings.Contains(view, fieldLabel) {
			return true
		}

		// If we've cycled back to the initial view, we've checked all fields
		if i > 0 && view == initialView {
			return false
		}

		f.env.Tab()
	}

	return false
}
