// Package forms provides common form utilities and configurations for the KaRiya TUI.
// It uses Charm's huh library for consistent, accessible form handling.
//
// # Architecture
//
// This package serves as the ONLY allowed import point for huh in the codebase.
// Screens, modals, and intents should import forms/ instead of huh directly.
//
// # Type Aliases
//
// The package provides type aliases for common huh types:
//   - forms.Form = *huh.Form
//   - forms.Group = *huh.Group
//   - forms.Field = huh.Field
//
// # Usage in Modals
//
//	modal := &MyModal{
//	    form: forms.NewFormWithDimensions(width, height, group),
//	}
//
//	// In Update:
//	updatedForm, cmd := forms.Update(m.form, msg)
//	m.form = updatedForm
//	if forms.IsCompleted(m.form) { ... }
package forms

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Form is the form type used throughout KaRiya.
// This is an alias for *huh.Form to avoid direct huh imports in screens/modals.
type Form = *huh.Form

// Group is a group of form fields.
// This is an alias for *huh.Group to avoid direct huh imports.
type Group = *huh.Group

// Field is a form field interface.
// This is an alias for huh.Field to avoid direct huh imports.
type Field = huh.Field

// Input is a text input field.
type Input = *huh.Input

// Text is a text area field.
type Text = *huh.Text

// Select is a single-select field.
type Select = *huh.Select[string]

// MultiSelect is a multi-select field.
type MultiSelect = *huh.MultiSelect[string]

// Confirm is a confirmation field.
type Confirm = *huh.Confirm

// NewGroup creates a new form group from fields.
//
// Expected:
//   - field must be valid.
//
// Returns:
//   - A Group value.
//
// Side effects:
//   - None.
func NewGroup(fields ...Field) Group {
	return huh.NewGroup(fields...)
}

// Update handles form message updates.
// Returns the updated form and any command to execute.
// Use this instead of calling form.Update directly.
//
// Expected:
//   - form must be a valid Form.
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - Form: the updated form.
//   - tea.Cmd: command to execute.
//
// Side effects:
//   - None.
func Update(form Form, msg tea.Msg) (Form, tea.Cmd) {
	model, cmd := form.Update(msg)
	f, ok := model.(*huh.Form)
	if !ok {
		return form, cmd
	}
	return f, cmd
}

// WithCustomSubmitKey configures a form to use a custom submit key binding.
//
// Expected:
//   - form must be a valid Form.
//   - keys must be valid key strings (e.g., "ctrl+s", "enter").
//
// Returns:
//   - A Form value.
//
// Side effects:
//   - None.
func WithCustomSubmitKey(form Form, keys ...string) Form {
	keymap := huh.NewDefaultKeyMap()
	keymap.Input.Submit = key.NewBinding(key.WithKeys(keys...))
	return form.WithKeyMap(keymap)
}

// Theme returns the Catppuccin theme configured for KaRiya forms.
//
// Returns:
//   - A fully initialized huh.Theme ready for use.
//
// Side effects:
//   - None.
func Theme() *huh.Theme {
	return huh.ThemeCatppuccin()
}

// newForm creates a new form with KaRiya's default theme and configuration.
func newForm(groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).
		WithTheme(Theme()).
		WithShowHelp(false).
		WithShowErrors(false)
}

// NewFormWithDimensions creates a new form with KaRiya's default theme and fixed dimensions.
//
// Expected:
//   - int must be valid.
//   - group must be valid.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewFormWithDimensions(width, height int, groups ...*huh.Group) *huh.Form {
	form := huh.NewForm(groups...).WithTheme(Theme())
	if height > 0 {
		form = form.WithHeight(height)
	}
	if width > 0 {
		form = form.WithWidth(width)
	}
	return form
}

// DefaultFormHeight calculates a reasonable form height based on terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func DefaultFormHeight(terminalHeight int) int {
	const overhead = 20
	const minHeight = 10

	height := terminalHeight - overhead
	if height < minHeight {
		height = minHeight
	}
	return height
}

// ModalFormHeight calculates the form height for forms displayed inside an
// overlay modal. The returned height is passed to huh's Group.WithHeight to
// enable viewport scrolling within the fields group.
//
// The calculation accounts for all chrome consumed by the overlay rendering
// pipeline: logo area (DefaultLogoHeight + 1 gap + 2 bottom margin = 12),
// modal border (2), modal padding (2), title with margin (2), and the
// footer separator line plus badge row (2). Total overhead = 20.
//
// Expected:
//   - terminalHeight must be a positive integer representing terminal rows.
//
// Returns:
//   - The maximum form height that fits within the overlay without truncation.
//
// Side effects:
//   - None.
func ModalFormHeight(terminalHeight int) int {
	const modalOverhead = 20
	const minHeight = 5

	height := terminalHeight - modalOverhead - HelpFooterHeight
	if height < minHeight {
		height = minHeight
	}
	return height
}

// ModalFormWidth calculates the usable form width inside an overlay modal.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func ModalFormWidth(modalWidth int) int {
	const chromeWidth = 6
	width := modalWidth - chromeWidth
	if width < 30 {
		width = 30
	}
	return width
}

// HelpFooterHeight is the space reserved for the help footer (blank line + badge row)
// rendered below the form inside modal views.
const HelpFooterHeight = 2

// ModalBoxChrome is the vertical space consumed by box border (2) + padding (4).
const ModalBoxChrome = 6

// newScrollableForm creates a scrollable form with a confirm field at the bottom.
//
// All fields (including the confirm) are placed in a single huh group so that
// the group's built-in viewport handles scrolling. This avoids the pagination
// behaviour of LayoutDefault (which shows one group per page) and the lack of
// scrolling in LayoutStack (which bypasses the viewport entirely).
//
// Expected:
//   - fields must contain at least one field.
//   - confirmValue must not be nil.
//   - width and height must be positive integers.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func newScrollableForm(fields []huh.Field, confirmValue *bool, width, height int) *huh.Form {
	confirmField := huh.NewConfirm().
		Key("submit").
		Title("Save Changes").
		Description("Submit the form to save your changes").
		Affirmative("Submit").
		Negative("Cancel").
		Value(confirmValue)

	fields = append(fields, confirmField)
	group := huh.NewGroup(fields...).WithHeight(height)

	form := huh.NewForm(group).
		WithTheme(Theme())

	if width > 0 {
		form = form.WithWidth(width)
	}

	return form
}

// IsCompleted checks if the form has been completed by the user.
//
// Expected:
//   - form must be valid.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsCompleted(form *huh.Form) bool {
	return form.State == huh.StateCompleted
}

// IsAborted checks if the form was cancelled/aborted by the user.
//
// Expected:
//   - form must be valid.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsAborted(form *huh.Form) bool {
	return form.State == huh.StateAborted
}

// IsTextInputFocused checks if the currently focused field is a text input
// (either Input or Text). This is used to determine whether vim-style j/k
// navigation should be disabled to allow typing those characters.
//
// Expected:
//   - form may be nil (returns false).
//
// Returns:
//   - true if the focused field is *huh.Input or *huh.Text.
//   - false otherwise (including nil form).
//
// Side effects:
//   - None.
func IsTextInputFocused(form Form) bool {
	if form == nil {
		return false
	}
	focused := form.GetFocusedField()
	switch focused.(type) {
	case *huh.Input, *huh.Text:
		return true
	}
	return false
}

// FieldConfig represents common field configuration options.
type FieldConfig struct {
	Key         string
	Title       string
	Description string
	Placeholder string
	CharLimit   int
	Required    bool
	Validate    func(string) error
}

// NewInput creates a pre-configured input field.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized huh.Input ready for use.
//
// Side effects:
//   - None.
func NewInput(config FieldConfig) *huh.Input {
	input := huh.NewInput().
		Key(config.Key).
		Title(config.Title).
		Prompt("> ")

	if config.Description != "" {
		input = input.Description(config.Description)
	}

	if config.Placeholder != "" {
		input = input.Placeholder(config.Placeholder)
	}

	if config.CharLimit > 0 {
		input = input.CharLimit(config.CharLimit)
	}

	if config.Validate != nil {
		input = input.Validate(config.Validate)
	}

	return input
}

// NewText creates a pre-configured text area field.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized huh.Text ready for use.
//
// Side effects:
//   - None.
func NewText(config FieldConfig) *huh.Text {
	text := huh.NewText().
		Key(config.Key).
		Title(config.Title)

	if config.Description != "" {
		text = text.Description(config.Description)
	}

	if config.Placeholder != "" {
		text = text.Placeholder(config.Placeholder)
	}

	if config.CharLimit > 0 {
		text = text.CharLimit(config.CharLimit)
	}

	if config.Validate != nil {
		text = text.Validate(config.Validate)
	}

	return text
}

// SelectOption represents a selectable option in a Select or MultiSelect field.
type SelectOption struct {
	Key   string
	Value string
}

// NewSelect creates a pre-configured select field.
//
// Expected:
//   - Must be a valid string.
//   - []selectoption must be valid.
//
// Returns:
//   - A fully initialized huh.Select[string] ready for use.
//
// Side effects:
//   - None.
func NewSelect(fieldKey, title, description string, options []SelectOption) *huh.Select[string] {
	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.Value, opt.Key)
	}

	sel := huh.NewSelect[string]().
		Key(fieldKey).
		Title(title).
		Options(huhOptions...)

	if description != "" {
		sel = sel.Description(description)
	}

	return sel
}

// NewMultiSelect creates a pre-configured multi-select field.
//
// Expected:
//   - Must be a valid string.
//   - []selectoption must be valid.
//   - int must be valid.
//
// Returns:
//   - A fully initialized huh.MultiSelect[string] ready for use.
//
// Side effects:
//   - None.
func NewMultiSelect(fieldKey, title, description string, options []SelectOption, limit int) *huh.MultiSelect[string] {
	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.Value, opt.Key)
	}

	multi := huh.NewMultiSelect[string]().
		Key(fieldKey).
		Title(title).
		Options(huhOptions...)

	if description != "" {
		multi = multi.Description(description)
	}

	if limit > 0 {
		multi = multi.Limit(limit)
	}

	return multi
}

// NewConfirm creates a pre-configured confirmation field.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized huh.Confirm ready for use.
//
// Side effects:
//   - None.
func NewConfirm(fieldKey, title, description, affirmative, negative string) *huh.Confirm {
	confirm := huh.NewConfirm().
		Key(fieldKey).
		Title(title)

	if description != "" {
		confirm = confirm.Description(description)
	}

	if affirmative != "" && negative != "" {
		confirm = confirm.Affirmative(affirmative).Negative(negative)
	}

	return confirm
}
